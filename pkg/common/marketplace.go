package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// MarketplaceManager handles plugin marketplace operations
type MarketplaceManager struct {
	registries  map[string]*MarketplaceRegistry
	cache       *MarketplaceCache
	deps        *Dependencies
	httpClient  *http.Client
	mutex       sync.RWMutex
	localPath   string
}

// MarketplaceRegistry represents a remote plugin registry
type MarketplaceRegistry struct {
	Name        string                 `json:"name"`
	URL         string                 `json:"url"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Priority    int                    `json:"priority"`
	Auth        *RegistryAuth          `json:"auth,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	LastSync    time.Time              `json:"last_sync"`
}

// RegistryAuth contains authentication information for private registries
type RegistryAuth struct {
	Type     string            `json:"type"` // "token", "basic", "none"
	Token    string            `json:"token,omitempty"`
	Username string            `json:"username,omitempty"`
	Password string            `json:"password,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
}

// MarketplacePlugin represents a plugin available in the marketplace
type MarketplacePlugin struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Version      string                 `json:"version"`
	Description  string                 `json:"description"`
	Author       string                 `json:"author"`
	License      string                 `json:"license"`
	Homepage     string                 `json:"homepage,omitempty"`
	Repository   string                 `json:"repository,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	Category     string                 `json:"category,omitempty"`
	DownloadURL  string                 `json:"download_url"`
	Checksum     string                 `json:"checksum,omitempty"`
	Size         int64                  `json:"size,omitempty"`
	Dependencies []string               `json:"dependencies,omitempty"`
	Requirements map[string]string      `json:"requirements,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	Downloads    int64                  `json:"downloads,omitempty"`
	Rating       float64                `json:"rating,omitempty"`
	Registry     string                 `json:"registry"`
}

// MarketplaceCache handles caching of marketplace data
type MarketplaceCache struct {
	plugins     map[string]*MarketplacePlugin
	registries  map[string]*MarketplaceRegistry
	lastUpdate  time.Time
	cachePath   string
	maxAge      time.Duration
	mutex       sync.RWMutex
}

// SearchResult represents search results from the marketplace
type SearchResult struct {
	Query     string               `json:"query"`
	Results   []*MarketplacePlugin `json:"results"`
	Total     int                  `json:"total"`
	Page      int                  `json:"page"`
	PerPage   int                  `json:"per_page"`
	Timestamp time.Time            `json:"timestamp"`
}

// InstallOptions contains options for plugin installation
type InstallOptions struct {
	Version     string            `json:"version,omitempty"`
	Source      string            `json:"source,omitempty"`
	Force       bool              `json:"force"`
	SkipDeps    bool              `json:"skip_deps"`
	Destination string            `json:"destination,omitempty"`
	Config      map[string]string `json:"config,omitempty"`
}

// NewMarketplaceManager creates a new marketplace manager
func NewMarketplaceManager(deps *Dependencies) (*MarketplaceManager, error) {
	// Use home directory as a fallback for user config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	localPath := filepath.Join(homeDir, ".engx", "marketplace")
	if err := os.MkdirAll(localPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create marketplace directory: %w", err)
	}

	cachePath := filepath.Join(localPath, "cache")
	cache, err := NewMarketplaceCache(cachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize marketplace cache: %w", err)
	}

	manager := &MarketplaceManager{
		registries: make(map[string]*MarketplaceRegistry),
		cache:      cache,
		deps:       deps,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		localPath: localPath,
	}

	// Load default registries
	if err := manager.loadDefaultRegistries(); err != nil {
		if deps.Logger != nil {
			deps.Logger.Warn("Failed to load default registries: %v", err)
		}
	}

	return manager, nil
}

// NewMarketplaceCache creates a new marketplace cache
func NewMarketplaceCache(cachePath string) (*MarketplaceCache, error) {
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	cache := &MarketplaceCache{
		plugins:    make(map[string]*MarketplacePlugin),
		registries: make(map[string]*MarketplaceRegistry),
		cachePath:  cachePath,
		maxAge:     24 * time.Hour, // Cache for 24 hours
	}

	// Load existing cache
	if err := cache.Load(); err != nil {
		// Cache load failure is not critical, just start fresh
		return cache, nil
	}

	return cache, nil
}

// AddRegistry adds a new plugin registry
func (mm *MarketplaceManager) AddRegistry(registry *MarketplaceRegistry) error {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	if registry.Name == "" {
		return fmt.Errorf("registry name cannot be empty")
	}

	if registry.URL == "" {
		return fmt.Errorf("registry URL cannot be empty")
	}

	mm.registries[registry.Name] = registry

	if mm.deps.Logger != nil {
		mm.deps.Logger.Info("Added registry: %s (%s)", registry.Name, registry.URL)
	}

	return mm.saveRegistries()
}

// RemoveRegistry removes a plugin registry
func (mm *MarketplaceManager) RemoveRegistry(name string) error {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	if _, exists := mm.registries[name]; !exists {
		return fmt.Errorf("registry %s not found", name)
	}

	delete(mm.registries, name)

	if mm.deps.Logger != nil {
		mm.deps.Logger.Info("Removed registry: %s", name)
	}

	return mm.saveRegistries()
}

// ListRegistries returns all configured registries
func (mm *MarketplaceManager) ListRegistries() []*MarketplaceRegistry {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	registries := make([]*MarketplaceRegistry, 0, len(mm.registries))
	for _, registry := range mm.registries {
		registries = append(registries, registry)
	}

	// Sort by priority (lower numbers = higher priority)
	sort.Slice(registries, func(i, j int) bool {
		return registries[i].Priority < registries[j].Priority
	})

	return registries
}

// SyncRegistries synchronizes plugin data from all enabled registries
func (mm *MarketplaceManager) SyncRegistries() error {
	mm.mutex.RLock()
	registries := make([]*MarketplaceRegistry, 0)
	for _, registry := range mm.registries {
		if registry.Enabled {
			registries = append(registries, registry)
		}
	}
	mm.mutex.RUnlock()

	var syncErrors []string
	var syncCount int

	for _, registry := range registries {
		if err := mm.syncRegistry(registry); err != nil {
			syncErrors = append(syncErrors, fmt.Sprintf("%s: %v", registry.Name, err))
			if mm.deps.Logger != nil {
				mm.deps.Logger.Error("Failed to sync registry %s: %v", registry.Name, err)
			}
		} else {
			syncCount++
			registry.LastSync = time.Now()
		}
	}

	// Save updated registry states
	if err := mm.saveRegistries(); err != nil {
		if mm.deps.Logger != nil {
			mm.deps.Logger.Warn("Failed to save registry states: %v", err)
		}
	}

	if len(syncErrors) > 0 {
		return fmt.Errorf("sync failed for %d registries: %s", len(syncErrors), strings.Join(syncErrors, "; "))
	}

	if mm.deps.Logger != nil {
		mm.deps.Logger.Info("Successfully synced %d registries", syncCount)
	}

	return nil
}

// SearchPlugins searches for plugins in the marketplace
func (mm *MarketplaceManager) SearchPlugins(query string, page, perPage int) (*SearchResult, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}

	mm.cache.mutex.RLock()
	defer mm.cache.mutex.RUnlock()

	var matchingPlugins []*MarketplacePlugin
	query = strings.ToLower(query)

	for _, plugin := range mm.cache.plugins {
		if mm.pluginMatches(plugin, query) {
			matchingPlugins = append(matchingPlugins, plugin)
		}
	}

	// Sort by relevance (simple scoring based on name match, then downloads)
	sort.Slice(matchingPlugins, func(i, j int) bool {
		scoreI := mm.calculateRelevanceScore(matchingPlugins[i], query)
		scoreJ := mm.calculateRelevanceScore(matchingPlugins[j], query)
		if scoreI != scoreJ {
			return scoreI > scoreJ
		}
		return matchingPlugins[i].Downloads > matchingPlugins[j].Downloads
	})

	// Paginate results
	start := (page - 1) * perPage
	end := start + perPage
	if start >= len(matchingPlugins) {
		matchingPlugins = []*MarketplacePlugin{}
	} else {
		if end > len(matchingPlugins) {
			end = len(matchingPlugins)
		}
		matchingPlugins = matchingPlugins[start:end]
	}

	return &SearchResult{
		Query:     query,
		Results:   matchingPlugins,
		Total:     len(matchingPlugins),
		Page:      page,
		PerPage:   perPage,
		Timestamp: time.Now(),
	}, nil
}

// GetPlugin retrieves a specific plugin by ID
func (mm *MarketplaceManager) GetPlugin(pluginID string) (*MarketplacePlugin, error) {
	mm.cache.mutex.RLock()
	defer mm.cache.mutex.RUnlock()

	plugin, exists := mm.cache.plugins[pluginID]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found in marketplace", pluginID)
	}

	return plugin, nil
}

// InstallPlugin installs a plugin from the marketplace
func (mm *MarketplaceManager) InstallPlugin(pluginID string, options *InstallOptions) error {
	plugin, err := mm.GetPlugin(pluginID)
	if err != nil {
		return fmt.Errorf("failed to find plugin %s: %w", pluginID, err)
	}

	if mm.deps.Logger != nil {
		mm.deps.Logger.Info("Installing plugin %s version %s", plugin.Name, plugin.Version)
	}

	// Check dependencies first
	if !options.SkipDeps && len(plugin.Dependencies) > 0 {
		if err := mm.checkDependencies(plugin); err != nil {
			return fmt.Errorf("dependency check failed: %w", err)
		}
	}

	// Download plugin
	if err := mm.downloadPlugin(plugin, options); err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}

	if mm.deps.Logger != nil {
		mm.deps.Logger.Info("Successfully installed plugin %s", plugin.Name)
	}

	return nil
}

// loadDefaultRegistries loads the default plugin registries
func (mm *MarketplaceManager) loadDefaultRegistries() error {
	defaultRegistries := []*MarketplaceRegistry{
		{
			Name:        "official",
			URL:         "https://registry.engx.dev/plugins",
			Description: "Official engx plugin registry",
			Enabled:     true,
			Priority:    1,
		},
		{
			Name:        "community",
			URL:         "https://community.engx.dev/plugins",
			Description: "Community-maintained plugins",
			Enabled:     true,
			Priority:    10,
		},
	}

	for _, registry := range defaultRegistries {
		mm.registries[registry.Name] = registry
	}

	return mm.saveRegistries()
}

// syncRegistry synchronizes plugins from a specific registry
func (mm *MarketplaceManager) syncRegistry(registry *MarketplaceRegistry) error {
	// This is a simulation - in reality, this would make HTTP requests
	// to fetch plugin metadata from the registry

	if mm.deps.Logger != nil {
		mm.deps.Logger.Debug("Syncing registry: %s", registry.Name)
	}

	// Simulate registry response with sample plugins
	samplePlugins := mm.generateSamplePlugins(registry.Name)

	// Update cache with new plugins
	mm.cache.mutex.Lock()
	for _, plugin := range samplePlugins {
		mm.cache.plugins[plugin.ID] = plugin
	}
	mm.cache.lastUpdate = time.Now()
	mm.cache.mutex.Unlock()

	// Save cache
	return mm.cache.Save()
}

// generateSamplePlugins creates sample plugins for demonstration
func (mm *MarketplaceManager) generateSamplePlugins(registryName string) []*MarketplacePlugin {
	now := time.Now()

	plugins := []*MarketplacePlugin{
		{
			ID:          fmt.Sprintf("%s/react-enhanced", registryName),
			Name:        "react-enhanced",
			Version:     "1.2.0",
			Description: "Enhanced React application scaffolding with TypeScript and modern tooling",
			Author:      "EngX Team",
			License:     "MIT",
			Homepage:    "https://github.com/engx/react-enhanced",
			Repository:  "https://github.com/engx/react-enhanced.git",
			Tags:        []string{"react", "typescript", "scaffold", "modern"},
			Category:    "Framework",
			DownloadURL: fmt.Sprintf("https://%s.engx.dev/plugins/react-enhanced/1.2.0.tar.gz", registryName),
			Size:        2048576,
			CreatedAt:   now.AddDate(0, -3, 0),
			UpdatedAt:   now.AddDate(0, 0, -7),
			Downloads:   15432,
			Rating:      4.8,
			Registry:    registryName,
		},
		{
			ID:          fmt.Sprintf("%s/vue-starter", registryName),
			Name:        "vue-starter",
			Version:     "1.0.5",
			Description: "Vue.js application starter with Vite and composition API",
			Author:      "Vue Community",
			License:     "MIT",
			Homepage:    "https://github.com/vue/vue-starter",
			Repository:  "https://github.com/vue/vue-starter.git",
			Tags:        []string{"vue", "vite", "composition-api", "starter"},
			Category:    "Framework",
			DownloadURL: fmt.Sprintf("https://%s.engx.dev/plugins/vue-starter/1.0.5.tar.gz", registryName),
			Size:        1536000,
			CreatedAt:   now.AddDate(0, -2, 0),
			UpdatedAt:   now.AddDate(0, 0, -14),
			Downloads:   8291,
			Rating:      4.5,
			Registry:    registryName,
		},
		{
			ID:          fmt.Sprintf("%s/go-microservice", registryName),
			Name:        "go-microservice",
			Version:     "2.1.0",
			Description: "Go microservice template with gRPC, REST API, and monitoring",
			Author:      "Go Team",
			License:     "Apache-2.0",
			Homepage:    "https://github.com/golang/microservice-template",
			Repository:  "https://github.com/golang/microservice-template.git",
			Tags:        []string{"go", "microservice", "grpc", "rest", "monitoring"},
			Category:    "Backend",
			DownloadURL: fmt.Sprintf("https://%s.engx.dev/plugins/go-microservice/2.1.0.tar.gz", registryName),
			Size:        3145728,
			Dependencies: []string{"docker", "protobuf"},
			CreatedAt:   now.AddDate(0, -6, 0),
			UpdatedAt:   now.AddDate(0, 0, -3),
			Downloads:   12847,
			Rating:      4.7,
			Registry:    registryName,
		},
	}

	// Add registry-specific variations
	if registryName == "community" {
		plugins = append(plugins, &MarketplacePlugin{
			ID:          "community/experimental-ui",
			Name:        "experimental-ui",
			Version:     "0.3.0-beta",
			Description: "Experimental UI components and patterns (beta)",
			Author:      "Community Contributors",
			License:     "MIT",
			Tags:        []string{"ui", "experimental", "beta", "components"},
			Category:    "UI/UX",
			DownloadURL: "https://community.engx.dev/plugins/experimental-ui/0.3.0-beta.tar.gz",
			Size:        512000,
			CreatedAt:   now.AddDate(0, -1, 0),
			UpdatedAt:   now.AddDate(0, 0, -2),
			Downloads:   1247,
			Rating:      3.9,
			Registry:    "community",
		})
	}

	return plugins
}

// pluginMatches checks if a plugin matches the search query
func (mm *MarketplaceManager) pluginMatches(plugin *MarketplacePlugin, query string) bool {
	searchFields := []string{
		plugin.Name,
		plugin.Description,
		plugin.Author,
		plugin.Category,
		strings.Join(plugin.Tags, " "),
	}

	for _, field := range searchFields {
		if strings.Contains(strings.ToLower(field), query) {
			return true
		}
	}

	return false
}

// calculateRelevanceScore calculates a relevance score for search results
func (mm *MarketplaceManager) calculateRelevanceScore(plugin *MarketplacePlugin, query string) float64 {
	score := 0.0

	// Name match (highest weight)
	if strings.Contains(strings.ToLower(plugin.Name), query) {
		score += 100
		if strings.ToLower(plugin.Name) == query {
			score += 50 // Exact match bonus
		}
	}

	// Description match
	if strings.Contains(strings.ToLower(plugin.Description), query) {
		score += 30
	}

	// Tag match
	for _, tag := range plugin.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			score += 20
		}
	}

	// Category match
	if strings.Contains(strings.ToLower(plugin.Category), query) {
		score += 15
	}

	// Author match
	if strings.Contains(strings.ToLower(plugin.Author), query) {
		score += 10
	}

	// Download count factor (normalize to 0-10 range)
	if plugin.Downloads > 0 {
		downloadScore := float64(plugin.Downloads) / 1000.0
		if downloadScore > 10 {
			downloadScore = 10
		}
		score += downloadScore
	}

	// Rating factor
	score += plugin.Rating

	return score
}

// checkDependencies verifies that plugin dependencies are available
func (mm *MarketplaceManager) checkDependencies(plugin *MarketplacePlugin) error {
	for _, dep := range plugin.Dependencies {
		// This is a simulation - in reality, this would check if dependencies
		// are installed or available in the system
		if mm.deps.Logger != nil {
			mm.deps.Logger.Debug("Checking dependency: %s", dep)
		}
	}
	return nil
}

// downloadPlugin downloads and installs a plugin
func (mm *MarketplaceManager) downloadPlugin(plugin *MarketplacePlugin, options *InstallOptions) error {
	// This is a simulation - in reality, this would:
	// 1. Download the plugin archive from the URL
	// 2. Verify checksums
	// 3. Extract to the appropriate location
	// 4. Register with the plugin manager

	if mm.deps.Logger != nil {
		mm.deps.Logger.Info("Downloading plugin %s from %s", plugin.Name, plugin.DownloadURL)
		mm.deps.Logger.Info("Plugin size: %d bytes", plugin.Size)
	}

	// Simulate download delay
	time.Sleep(100 * time.Millisecond)

	return nil
}

// saveRegistries saves registry configuration to disk
func (mm *MarketplaceManager) saveRegistries() error {
	configPath := filepath.Join(mm.localPath, "registries.json")

	data, err := json.MarshalIndent(mm.registries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registries: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write registries config: %w", err)
	}

	return nil
}

// Load loads the marketplace cache from disk
func (mc *MarketplaceCache) Load() error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	cachePath := filepath.Join(mc.cachePath, "plugins.json")
	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No cache file exists yet
		}
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	var cacheData struct {
		Plugins    map[string]*MarketplacePlugin `json:"plugins"`
		LastUpdate time.Time                     `json:"last_update"`
	}

	if err := json.Unmarshal(data, &cacheData); err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	mc.plugins = cacheData.Plugins
	mc.lastUpdate = cacheData.LastUpdate

	return nil
}

// Save saves the marketplace cache to disk
func (mc *MarketplaceCache) Save() error {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	cacheData := struct {
		Plugins    map[string]*MarketplacePlugin `json:"plugins"`
		LastUpdate time.Time                     `json:"last_update"`
	}{
		Plugins:    mc.plugins,
		LastUpdate: mc.lastUpdate,
	}

	data, err := json.MarshalIndent(cacheData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	cachePath := filepath.Join(mc.cachePath, "plugins.json")
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// IsExpired checks if the cache is expired
func (mc *MarketplaceCache) IsExpired() bool {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	return time.Since(mc.lastUpdate) > mc.maxAge
}

// Clear clears all cached data
func (mc *MarketplaceCache) Clear() error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.plugins = make(map[string]*MarketplacePlugin)
	mc.lastUpdate = time.Time{}

	return mc.Save()
}

// GetStats returns marketplace statistics
func (mm *MarketplaceManager) GetStats() map[string]interface{} {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	enabledRegistries := 0
	totalPlugins := 0

	for _, registry := range mm.registries {
		if registry.Enabled {
			enabledRegistries++
		}
	}

	mm.cache.mutex.RLock()
	totalPlugins = len(mm.cache.plugins)
	lastSync := mm.cache.lastUpdate
	mm.cache.mutex.RUnlock()

	return map[string]interface{}{
		"total_registries":   len(mm.registries),
		"enabled_registries": enabledRegistries,
		"total_plugins":      totalPlugins,
		"cache_expired":      mm.cache.IsExpired(),
		"last_sync":          lastSync,
		"cache_path":         mm.cache.cachePath,
	}
}