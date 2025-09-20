package registry

import (
	"fmt"
	"sort"
	"sync"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common"
	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
	"github.com/spf13/cobra"
)

// Manager handles the lifecycle and registration of command plugins.
// It provides thread-safe operations for registering, discovering,
// and managing plugins throughout the application lifecycle.
type Manager struct {
	plugins    map[string]interfaces.CommandPlugin
	deps       *common.Dependencies
	registered []string
	validator  *Validator
	logger     interfaces.Logger
	mutex      sync.RWMutex
}

// NewManager creates a new plugin registry manager with the provided dependencies.
func NewManager(deps *common.Dependencies) *Manager {
	return &Manager{
		plugins:   make(map[string]interfaces.CommandPlugin),
		deps:      deps,
		validator: NewValidator(),
		logger:    deps.Logger.WithComponent("registry"),
		mutex:     sync.RWMutex{},
	}
}

// Register adds a plugin to the registry after validation and initialization.
// This method is thread-safe and handles the complete plugin lifecycle.
func (m *Manager) Register(plugin interfaces.CommandPlugin) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Validate plugin before registration
	if err := m.validator.ValidatePlugin(plugin); err != nil {
		return fmt.Errorf("plugin validation failed for %s: %w", plugin.Name(), err)
	}

	// Check for duplicate registration
	if _, exists := m.plugins[plugin.Name()]; exists {
		return fmt.Errorf("plugin %s is already registered", plugin.Name())
	}

	// Validate dependencies are available
	if err := m.validateDependencies(plugin); err != nil {
		return fmt.Errorf("dependency validation failed for %s: %w", plugin.Name(), err)
	}

	// Initialize the plugin
	m.logger.Info("Initializing plugin: %s", plugin.Name())
	if err := plugin.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %w", plugin.Name(), err)
	}

	// Register the plugin
	m.plugins[plugin.Name()] = plugin
	m.registered = append(m.registered, plugin.Name())

	m.logger.Info("Successfully registered plugin: %s (version %s)", plugin.Name(), plugin.Version())
	return nil
}

// SafeRegister provides error recovery during plugin registration.
// If a plugin fails to register, it logs the error but doesn't crash the system.
func (m *Manager) SafeRegister(plugin interfaces.CommandPlugin) error {
	defer func() {
		if r := recover(); r != nil {
			m.logger.Error("Plugin %s panicked during registration: %v", plugin.Name(), r)
		}
	}()

	return m.Register(plugin)
}

// Unregister removes a plugin from the registry and cleans up its resources.
func (m *Manager) Unregister(pluginName string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	plugin, exists := m.plugins[pluginName]
	if !exists {
		return fmt.Errorf("plugin %s is not registered", pluginName)
	}

	// Cleanup plugin resources
	if err := plugin.Cleanup(); err != nil {
		m.logger.Warn("Error during cleanup of plugin %s: %v", pluginName, err)
	}

	// Remove from registry
	delete(m.plugins, pluginName)

	// Remove from registered list
	for i, name := range m.registered {
		if name == pluginName {
			m.registered = append(m.registered[:i], m.registered[i+1:]...)
			break
		}
	}

	m.logger.Info("Unregistered plugin: %s", pluginName)
	return nil
}

// GetCommands returns all registered plugin commands as cobra commands.
// Commands are returned in registration order for consistent behavior.
func (m *Manager) GetCommands() []*cobra.Command {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var commands []*cobra.Command
	for _, pluginName := range m.registered {
		if plugin, exists := m.plugins[pluginName]; exists {
			cmd := plugin.Create(m.deps)
			if cmd != nil {
				commands = append(commands, cmd)
			}
		}
	}

	return commands
}

// GetCommandsSafely returns commands with error recovery.
// If a plugin fails to create a command, it's logged but doesn't affect other plugins.
func (m *Manager) GetCommandsSafely() []*cobra.Command {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var commands []*cobra.Command
	for _, pluginName := range m.registered {
		func() {
			defer func() {
				if r := recover(); r != nil {
					m.logger.Error("Plugin %s panicked during command creation: %v", pluginName, r)
				}
			}()

			if plugin, exists := m.plugins[pluginName]; exists {
				cmd := plugin.Create(m.deps)
				if cmd != nil {
					commands = append(commands, cmd)
				}
			}
		}()
	}

	return commands
}

// GetPlugin returns a specific plugin by name.
func (m *Manager) GetPlugin(name string) (interfaces.CommandPlugin, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	plugin, exists := m.plugins[name]
	return plugin, exists
}

// ListPlugins returns information about all registered plugins.
func (m *Manager) ListPlugins() []interfaces.PluginInfo {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var plugins []interfaces.PluginInfo
	for _, name := range m.registered {
		if plugin, exists := m.plugins[name]; exists {
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
	}

	return plugins
}

// GetRegisteredNames returns the names of all registered plugins in registration order.
func (m *Manager) GetRegisteredNames() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy to prevent external modification
	names := make([]string, len(m.registered))
	copy(names, m.registered)
	return names
}

// Count returns the number of registered plugins.
func (m *Manager) Count() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.plugins)
}

// Cleanup performs cleanup of all registered plugins.
// This should be called during application shutdown.
func (m *Manager) Cleanup() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var errors []error

	// Cleanup plugins in reverse registration order
	for i := len(m.registered) - 1; i >= 0; i-- {
		pluginName := m.registered[i]
		if plugin, exists := m.plugins[pluginName]; exists {
			if err := plugin.Cleanup(); err != nil {
				errors = append(errors, fmt.Errorf("cleanup failed for plugin %s: %w", pluginName, err))
			}
		}
	}

	// Clear all registrations
	m.plugins = make(map[string]interfaces.CommandPlugin)
	m.registered = nil

	if len(errors) > 0 {
		return &MultiError{Errors: errors}
	}

	return nil
}

// validateDependencies checks if all required dependencies for a plugin are available.
func (m *Manager) validateDependencies(plugin interfaces.CommandPlugin) error {
	required := plugin.RequiredServices()
	for _, serviceName := range required {
		if !m.deps.HasService(serviceName) {
			return fmt.Errorf("required service %s is not available", serviceName)
		}
	}

	return nil
}

// GetStats returns statistics about the plugin registry.
func (m *Manager) GetStats() RegistryStats {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := RegistryStats{
		TotalPlugins:      len(m.plugins),
		RegisteredPlugins: len(m.registered),
		PluginNames:       make([]string, len(m.registered)),
	}

	copy(stats.PluginNames, m.registered)
	sort.Strings(stats.PluginNames)

	return stats
}

// RegistryStats provides statistics about the plugin registry state.
type RegistryStats struct {
	TotalPlugins      int      `json:"total_plugins"`
	RegisteredPlugins int      `json:"registered_plugins"`
	PluginNames       []string `json:"plugin_names"`
}

// MultiError represents multiple errors that occurred during operation.
type MultiError struct {
	Errors []error
}

func (me *MultiError) Error() string {
	if len(me.Errors) == 0 {
		return "no errors"
	}
	if len(me.Errors) == 1 {
		return me.Errors[0].Error()
	}
	return fmt.Sprintf("multiple errors occurred: %s (and %d more)",
		me.Errors[0].Error(), len(me.Errors)-1)
}