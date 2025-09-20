package common

import (
	"fmt"
	"sync"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
	"github.com/spf13/cobra"
)

// PluginRegistry manages plugin registration and lifecycle
type PluginRegistry struct {
	plugins      map[string]interfaces.CommandPlugin
	validator    *PluginValidator
	hotReload    *HotReloadManager
	deps         *Dependencies
	mutex        sync.RWMutex
}

// NewPluginRegistry creates a new plugin registry
func NewPluginRegistry(deps *Dependencies) *PluginRegistry {
	registry := &PluginRegistry{
		plugins: make(map[string]interfaces.CommandPlugin),
		deps:    deps,
	}

	// Initialize validator
	registry.validator = NewPluginValidator(registry, deps)

	// Initialize hot-reload manager
	if hotReload, err := NewHotReloadManager(registry, deps); err != nil {
		if deps.Logger != nil {
			deps.Logger.Warn("Failed to initialize hot reload manager: %v", err)
		}
	} else {
		registry.hotReload = hotReload
	}

	return registry
}

// Register registers a plugin
func (r *PluginRegistry) Register(plugin interfaces.CommandPlugin) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	name := plugin.Name()
	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}

	r.plugins[name] = plugin
	return nil
}

// GetPlugin retrieves a plugin by name
func (r *PluginRegistry) GetPlugin(name string) interfaces.CommandPlugin {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return nil
	}

	return plugin
}

// GetPluginWithError retrieves a plugin by name with error handling
func (r *PluginRegistry) GetPluginWithError(name string) (interfaces.CommandPlugin, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return plugin, nil
}

// ListPlugins returns all registered plugins
func (r *PluginRegistry) ListPlugins() []interfaces.PluginInfo {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var plugins []interfaces.PluginInfo
	for _, plugin := range r.plugins {
		info := interfaces.PluginInfo{
			Metadata: interfaces.PluginMetadata{
				Name:        plugin.Name(),
				Description: plugin.Description(),
				Version:     plugin.Version(),
			},
			Status: interfaces.PluginStatusActive,
		}
		plugins = append(plugins, info)
	}

	return plugins
}

// InitializeAll initializes all registered plugins
func (r *PluginRegistry) InitializeAll() error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for name, plugin := range r.plugins {
		if err := plugin.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize plugin %s: %w", name, err)
		}
	}

	return nil
}

// CreateCommands creates cobra commands for all registered plugins
func (r *PluginRegistry) CreateCommands() ([]*cobra.Command, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var commands []*cobra.Command
	for _, plugin := range r.plugins {
		cmd := plugin.Create(r.deps)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	}

	return commands, nil
}

// Cleanup cleans up all plugins
func (r *PluginRegistry) Cleanup() error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Cleanup hot-reload manager first
	if r.hotReload != nil {
		if err := r.hotReload.Cleanup(); err != nil {
			if r.deps.Logger != nil {
				r.deps.Logger.Error("Failed to cleanup hot reload manager: %v", err)
			}
		}
	}

	// Cleanup all plugins
	for name, plugin := range r.plugins {
		if err := plugin.Cleanup(); err != nil {
			return fmt.Errorf("failed to cleanup plugin %s: %w", name, err)
		}
	}

	return nil
}

// GetValidator returns the plugin validator instance
func (r *PluginRegistry) GetValidator() *PluginValidator {
	return r.validator
}

// GetHotReloadManager returns the hot reload manager instance
func (r *PluginRegistry) GetHotReloadManager() *HotReloadManager {
	return r.hotReload
}

// EnableHotReload enables the hot reload system
func (r *PluginRegistry) EnableHotReload() error {
	if r.hotReload == nil {
		return fmt.Errorf("hot reload manager not initialized")
	}
	return r.hotReload.Enable()
}

// DisableHotReload disables the hot reload system
func (r *PluginRegistry) DisableHotReload() error {
	if r.hotReload == nil {
		return fmt.Errorf("hot reload manager not initialized")
	}
	return r.hotReload.Disable()
}

// IsHotReloadEnabled returns whether hot reload is enabled
func (r *PluginRegistry) IsHotReloadEnabled() bool {
	if r.hotReload == nil {
		return false
	}
	return r.hotReload.IsEnabled()
}