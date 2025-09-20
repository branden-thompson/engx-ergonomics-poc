package common

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
	"github.com/fsnotify/fsnotify"
)

// HotReloadManager handles runtime plugin reloading and management
type HotReloadManager struct {
	registry        *PluginRegistry
	deps            *Dependencies
	watcher         *fsnotify.Watcher
	pluginPaths     map[string]string // plugin name -> file path
	pluginStates    map[string]PluginState
	enabled         bool
	watchPaths      []string
	mutex           sync.RWMutex
	reloadCallbacks []ReloadCallback
}

// PluginState represents the runtime state of a plugin
type PluginState struct {
	Name        string                     `json:"name"`
	Enabled     bool                       `json:"enabled"`
	Loaded      bool                       `json:"loaded"`
	LastReload  time.Time                  `json:"last_reload"`
	ReloadCount int                        `json:"reload_count"`
	Instance    interfaces.CommandPlugin   `json:"-"`
	Error       string                     `json:"error,omitempty"`
	FilePath    string                     `json:"file_path"`
}

// ReloadCallback is called when a plugin is reloaded
type ReloadCallback func(pluginName string, state PluginState, event ReloadEvent) error

// ReloadEvent represents the type of reload event
type ReloadEvent int

const (
	ReloadEventLoaded ReloadEvent = iota
	ReloadEventUnloaded
	ReloadEventReloaded
	ReloadEventError
	ReloadEventEnabled
	ReloadEventDisabled
)

// String returns a string representation of the reload event
func (re ReloadEvent) String() string {
	switch re {
	case ReloadEventLoaded:
		return "loaded"
	case ReloadEventUnloaded:
		return "unloaded"
	case ReloadEventReloaded:
		return "reloaded"
	case ReloadEventError:
		return "error"
	case ReloadEventEnabled:
		return "enabled"
	case ReloadEventDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}

// NewHotReloadManager creates a new hot reload manager
func NewHotReloadManager(registry *PluginRegistry, deps *Dependencies) (*HotReloadManager, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	manager := &HotReloadManager{
		registry:        registry,
		deps:            deps,
		watcher:         watcher,
		pluginPaths:     make(map[string]string),
		pluginStates:    make(map[string]PluginState),
		enabled:         false,
		watchPaths:      []string{"./plugins", "./internal/plugins"},
		reloadCallbacks: make([]ReloadCallback, 0),
	}

	// Initialize plugin states for existing plugins
	manager.initializePluginStates()

	return manager, nil
}

// Enable starts the hot reload system
func (hrm *HotReloadManager) Enable() error {
	hrm.mutex.Lock()
	defer hrm.mutex.Unlock()

	if hrm.enabled {
		return nil // Already enabled
	}

	// Add watch paths
	for _, path := range hrm.watchPaths {
		if err := hrm.addWatchPath(path); err != nil {
			if hrm.deps.Logger != nil {
				hrm.deps.Logger.Warn("Failed to watch path %s: %v", path, err)
			}
		}
	}

	hrm.enabled = true

	// Start watching in a goroutine
	go hrm.watchForChanges()

	if hrm.deps.Logger != nil {
		hrm.deps.Logger.Info("Hot reload system enabled")
	}

	return nil
}

// Disable stops the hot reload system
func (hrm *HotReloadManager) Disable() error {
	hrm.mutex.Lock()
	defer hrm.mutex.Unlock()

	if !hrm.enabled {
		return nil // Already disabled
	}

	hrm.enabled = false

	if err := hrm.watcher.Close(); err != nil {
		return fmt.Errorf("failed to close file watcher: %w", err)
	}

	if hrm.deps.Logger != nil {
		hrm.deps.Logger.Info("Hot reload system disabled")
	}

	return nil
}

// IsEnabled returns whether hot reload is enabled
func (hrm *HotReloadManager) IsEnabled() bool {
	hrm.mutex.RLock()
	defer hrm.mutex.RUnlock()
	return hrm.enabled
}

// AddWatchPath adds a new path to watch for plugin changes
func (hrm *HotReloadManager) AddWatchPath(path string) error {
	hrm.mutex.Lock()
	defer hrm.mutex.Unlock()

	hrm.watchPaths = append(hrm.watchPaths, path)

	if hrm.enabled {
		return hrm.addWatchPath(path)
	}

	return nil
}

// addWatchPath internal method to add watch path (must be called with lock)
func (hrm *HotReloadManager) addWatchPath(path string) error {
	// Check if path exists
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// Create directory if it doesn't exist
			if err := os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("failed to create watch directory %s: %w", path, err)
			}
		} else {
			return fmt.Errorf("failed to stat watch path %s: %w", path, err)
		}
	}

	if err := hrm.watcher.Add(path); err != nil {
		return fmt.Errorf("failed to add watch path %s: %w", path, err)
	}

	// Watch subdirectories recursively
	return filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && walkPath != path {
			if err := hrm.watcher.Add(walkPath); err != nil {
				if hrm.deps.Logger != nil {
					hrm.deps.Logger.Warn("Failed to watch subdirectory %s: %v", walkPath, err)
				}
			}
		}

		return nil
	})
}

// EnablePlugin enables a plugin at runtime
func (hrm *HotReloadManager) EnablePlugin(pluginName string) error {
	hrm.mutex.Lock()
	defer hrm.mutex.Unlock()

	state, exists := hrm.pluginStates[pluginName]
	if !exists {
		// Try refreshing states first in case they weren't initialized
		hrm.mutex.Unlock()
		hrm.RefreshPluginStates()
		hrm.mutex.Lock()
		state, exists = hrm.pluginStates[pluginName]
		if !exists {
			return fmt.Errorf("plugin %s not found", pluginName)
		}
	}

	if state.Enabled {
		return nil // Already enabled
	}

	state.Enabled = true
	state.LastReload = time.Now()
	hrm.pluginStates[pluginName] = state

	if hrm.deps.Logger != nil {
		hrm.deps.Logger.Info("Plugin %s enabled", pluginName)
	}

	// Call callbacks
	for _, callback := range hrm.reloadCallbacks {
		if err := callback(pluginName, state, ReloadEventEnabled); err != nil {
			if hrm.deps.Logger != nil {
				hrm.deps.Logger.Error("Reload callback error for %s: %v", pluginName, err)
			}
		}
	}

	return nil
}

// DisablePlugin disables a plugin at runtime
func (hrm *HotReloadManager) DisablePlugin(pluginName string) error {
	hrm.mutex.Lock()
	defer hrm.mutex.Unlock()

	state, exists := hrm.pluginStates[pluginName]
	if !exists {
		// Try refreshing states first in case they weren't initialized
		hrm.mutex.Unlock()
		hrm.RefreshPluginStates()
		hrm.mutex.Lock()
		state, exists = hrm.pluginStates[pluginName]
		if !exists {
			return fmt.Errorf("plugin %s not found", pluginName)
		}
	}

	if !state.Enabled {
		return nil // Already disabled
	}

	state.Enabled = false
	state.LastReload = time.Now()
	hrm.pluginStates[pluginName] = state

	if hrm.deps.Logger != nil {
		hrm.deps.Logger.Info("Plugin %s disabled", pluginName)
	}

	// Call callbacks
	for _, callback := range hrm.reloadCallbacks {
		if err := callback(pluginName, state, ReloadEventDisabled); err != nil {
			if hrm.deps.Logger != nil {
				hrm.deps.Logger.Error("Reload callback error for %s: %v", pluginName, err)
			}
		}
	}

	return nil
}

// GetPluginState returns the current state of a plugin
func (hrm *HotReloadManager) GetPluginState(pluginName string) (PluginState, bool) {
	hrm.mutex.RLock()
	defer hrm.mutex.RUnlock()

	state, exists := hrm.pluginStates[pluginName]
	return state, exists
}

// ListPluginStates returns all plugin states
func (hrm *HotReloadManager) ListPluginStates() map[string]PluginState {
	hrm.mutex.RLock()
	defer hrm.mutex.RUnlock()

	states := make(map[string]PluginState)
	for name, state := range hrm.pluginStates {
		states[name] = state
	}

	return states
}

// AddReloadCallback adds a callback to be called on reload events
func (hrm *HotReloadManager) AddReloadCallback(callback ReloadCallback) {
	hrm.mutex.Lock()
	defer hrm.mutex.Unlock()

	hrm.reloadCallbacks = append(hrm.reloadCallbacks, callback)
}

// initializePluginStates sets up initial plugin states
func (hrm *HotReloadManager) initializePluginStates() {
	plugins := hrm.registry.ListPlugins()
	for _, pluginInfo := range plugins {
		state := PluginState{
			Name:        pluginInfo.Metadata.Name,
			Enabled:     true,
			Loaded:      true,
			LastReload:  time.Now(),
			ReloadCount: 0,
			Instance:    hrm.registry.GetPlugin(pluginInfo.Metadata.Name),
		}
		hrm.pluginStates[pluginInfo.Metadata.Name] = state
	}
}

// RefreshPluginStates refreshes plugin states from the registry
func (hrm *HotReloadManager) RefreshPluginStates() {
	hrm.mutex.Lock()
	defer hrm.mutex.Unlock()

	plugins := hrm.registry.ListPlugins()
	for _, pluginInfo := range plugins {
		// Check if plugin state already exists
		if _, exists := hrm.pluginStates[pluginInfo.Metadata.Name]; !exists {
			// Add new plugin state
			state := PluginState{
				Name:        pluginInfo.Metadata.Name,
				Enabled:     true,
				Loaded:      true,
				LastReload:  time.Now(),
				ReloadCount: 0,
				Instance:    hrm.registry.GetPlugin(pluginInfo.Metadata.Name),
			}
			hrm.pluginStates[pluginInfo.Metadata.Name] = state
		}
	}
}

// watchForChanges monitors file system changes
func (hrm *HotReloadManager) watchForChanges() {
	debounceMap := make(map[string]time.Time)
	const debounceInterval = 500 * time.Millisecond

	for {
		select {
		case event, ok := <-hrm.watcher.Events:
			if !ok {
				return
			}

			if hrm.shouldProcessEvent(event) {
				// Debounce rapid file changes
				now := time.Now()
				if lastTime, exists := debounceMap[event.Name]; exists {
					if now.Sub(lastTime) < debounceInterval {
						continue
					}
				}
				debounceMap[event.Name] = now

				hrm.handleFileChange(event)
			}

		case err, ok := <-hrm.watcher.Errors:
			if !ok {
				return
			}

			if hrm.deps.Logger != nil {
				hrm.deps.Logger.Error("File watcher error: %v", err)
			}
		}
	}
}

// shouldProcessEvent determines if a file system event should trigger a reload
func (hrm *HotReloadManager) shouldProcessEvent(event fsnotify.Event) bool {
	// Only process .go files in plugin directories
	if filepath.Ext(event.Name) != ".go" {
		return false
	}

	// Skip test files
	if filepath.Base(event.Name) == "plugin_test.go" {
		return false
	}

	// Only process write and create events
	return event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create
}

// handleFileChange processes a file change event
func (hrm *HotReloadManager) handleFileChange(event fsnotify.Event) {
	if hrm.deps.Logger != nil {
		hrm.deps.Logger.Debug("Plugin file changed: %s", event.Name)
	}

	// Extract plugin name from file path
	pluginName := hrm.extractPluginName(event.Name)
	if pluginName == "" {
		return
	}

	// Note: In a real implementation, this would involve recompiling the plugin
	// For this simulation, we'll just update the state and call callbacks
	hrm.simulatePluginReload(pluginName, event.Name)
}

// extractPluginName extracts plugin name from file path
func (hrm *HotReloadManager) extractPluginName(filePath string) string {
	// Extract from path like ./plugins/create/plugin.go -> "create"
	dir := filepath.Dir(filePath)
	return filepath.Base(dir)
}

// simulatePluginReload simulates reloading a plugin (for this POC)
func (hrm *HotReloadManager) simulatePluginReload(pluginName string, filePath string) {
	hrm.mutex.Lock()
	defer hrm.mutex.Unlock()

	state, exists := hrm.pluginStates[pluginName]
	if !exists {
		// New plugin detected
		state = PluginState{
			Name:    pluginName,
			Enabled: true,
			Loaded:  false,
		}
	}

	state.LastReload = time.Now()
	state.ReloadCount++
	state.FilePath = filePath

	// In a real implementation, this would involve:
	// 1. Unloading the old plugin
	// 2. Recompiling the plugin
	// 3. Loading the new plugin
	// For this simulation, we'll just update the state

	if hrm.deps.Logger != nil {
		hrm.deps.Logger.Info("Simulated reload of plugin %s (reload #%d)", pluginName, state.ReloadCount)
	}

	hrm.pluginStates[pluginName] = state

	// Call callbacks
	event := ReloadEventReloaded
	if !exists {
		event = ReloadEventLoaded
	}

	for _, callback := range hrm.reloadCallbacks {
		if err := callback(pluginName, state, event); err != nil {
			if hrm.deps.Logger != nil {
				hrm.deps.Logger.Error("Reload callback error for %s: %v", pluginName, err)
			}
		}
	}
}

// Cleanup properly shuts down the hot reload manager
func (hrm *HotReloadManager) Cleanup() error {
	return hrm.Disable()
}

// GetStats returns hot reload statistics
func (hrm *HotReloadManager) GetStats() map[string]interface{} {
	hrm.mutex.RLock()
	defer hrm.mutex.RUnlock()

	totalReloads := 0
	enabledCount := 0
	loadedCount := 0

	for _, state := range hrm.pluginStates {
		totalReloads += state.ReloadCount
		if state.Enabled {
			enabledCount++
		}
		if state.Loaded {
			loadedCount++
		}
	}

	return map[string]interface{}{
		"enabled":        hrm.enabled,
		"total_plugins":  len(hrm.pluginStates),
		"enabled_plugins": enabledCount,
		"loaded_plugins":  loadedCount,
		"total_reloads":   totalReloads,
		"watch_paths":     hrm.watchPaths,
		"callback_count":  len(hrm.reloadCallbacks),
	}
}