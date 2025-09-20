package registry

import (
	"fmt"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
)

// DiscoveryResult represents the result of plugin discovery.
type DiscoveryResult struct {
	Plugins []interfaces.CommandPlugin
	Errors  []error
}

// PluginDiscoverer defines the interface for plugin discovery mechanisms.
type PluginDiscoverer interface {
	DiscoverPlugins() DiscoveryResult
	GetName() string
}

// BuiltinDiscoverer handles discovery of built-in plugins.
// This is the primary discovery mechanism for compile-time plugins.
type BuiltinDiscoverer struct {
	plugins []interfaces.CommandPlugin
}

// NewBuiltinDiscoverer creates a new built-in plugin discoverer.
func NewBuiltinDiscoverer() *BuiltinDiscoverer {
	return &BuiltinDiscoverer{
		plugins: make([]interfaces.CommandPlugin, 0),
	}
}

// AddPlugin adds a plugin to the built-in discovery list.
// This is used during application initialization to register compile-time plugins.
func (bd *BuiltinDiscoverer) AddPlugin(plugin interfaces.CommandPlugin) {
	bd.plugins = append(bd.plugins, plugin)
}

// DiscoverPlugins returns all built-in plugins.
func (bd *BuiltinDiscoverer) DiscoverPlugins() DiscoveryResult {
	return DiscoveryResult{
		Plugins: bd.plugins,
		Errors:  nil,
	}
}

// GetName returns the discoverer name.
func (bd *BuiltinDiscoverer) GetName() string {
	return "builtin"
}

// CompositeDiscoverer combines multiple discovery mechanisms.
type CompositeDiscoverer struct {
	discoverers []PluginDiscoverer
}

// NewCompositeDiscoverer creates a new composite discoverer.
func NewCompositeDiscoverer(discoverers ...PluginDiscoverer) *CompositeDiscoverer {
	return &CompositeDiscoverer{
		discoverers: discoverers,
	}
}

// AddDiscoverer adds a discoverer to the composite.
func (cd *CompositeDiscoverer) AddDiscoverer(discoverer PluginDiscoverer) {
	cd.discoverers = append(cd.discoverers, discoverer)
}

// DiscoverPlugins runs all discoverers and combines results.
func (cd *CompositeDiscoverer) DiscoverPlugins() DiscoveryResult {
	var allPlugins []interfaces.CommandPlugin
	var allErrors []error

	for _, discoverer := range cd.discoverers {
		result := discoverer.DiscoverPlugins()
		allPlugins = append(allPlugins, result.Plugins...)
		allErrors = append(allErrors, result.Errors...)
	}

	return DiscoveryResult{
		Plugins: allPlugins,
		Errors:  allErrors,
	}
}

// GetName returns the discoverer name.
func (cd *CompositeDiscoverer) GetName() string {
	return "composite"
}

// DiscoverBuiltinPlugins provides a convenient function to discover
// all built-in plugins for the engx application.
func DiscoverBuiltinPlugins() []interfaces.CommandPlugin {
	discoverer := NewBuiltinDiscoverer()

	// TODO: This will be populated as we convert existing commands to plugins
	// discoverer.AddPlugin(&create.Plugin{})
	// discoverer.AddPlugin(&test.Plugin{})

	result := discoverer.DiscoverPlugins()
	return result.Plugins
}

// AutoDiscoverPlugins provides automatic plugin discovery using all available mechanisms.
func AutoDiscoverPlugins() DiscoveryResult {
	builtinDiscoverer := NewBuiltinDiscoverer()

	// TODO: Add converted plugins here
	// builtinDiscoverer.AddPlugin(&create.Plugin{})
	// builtinDiscoverer.AddPlugin(&test.Plugin{})

	composite := NewCompositeDiscoverer(builtinDiscoverer)

	// TODO: Add other discovery mechanisms in the future
	// - File-based discovery
	// - Environment-based discovery
	// - Network-based discovery (for advanced scenarios)

	return composite.DiscoverPlugins()
}

// DiscoveryStats provides statistics about plugin discovery.
type DiscoveryStats struct {
	TotalDiscovered int                    `json:"total_discovered"`
	SuccessfulPlugins int                  `json:"successful_plugins"`
	FailedPlugins   int                    `json:"failed_plugins"`
	DiscovererStats map[string]interface{} `json:"discoverer_stats"`
}

// GetDiscoveryStats returns statistics about the last discovery operation.
func GetDiscoveryStats(result DiscoveryResult) DiscoveryStats {
	return DiscoveryStats{
		TotalDiscovered:   len(result.Plugins) + len(result.Errors),
		SuccessfulPlugins: len(result.Plugins),
		FailedPlugins:     len(result.Errors),
		DiscovererStats:   make(map[string]interface{}),
	}
}

// ValidateDiscoveredPlugins validates all discovered plugins and returns
// a filtered list of valid plugins along with validation errors.
func ValidateDiscoveredPlugins(plugins []interfaces.CommandPlugin) ([]interfaces.CommandPlugin, []error) {
	validator := NewValidator()
	var validPlugins []interfaces.CommandPlugin
	var errors []error

	for _, plugin := range plugins {
		if err := validator.ValidatePlugin(plugin); err != nil {
			errors = append(errors, fmt.Errorf("plugin %s validation failed: %w", plugin.Name(), err))
		} else {
			validPlugins = append(validPlugins, plugin)
		}
	}

	return validPlugins, errors
}

// DiscoveryConfiguration holds configuration for plugin discovery.
type DiscoveryConfiguration struct {
	EnableBuiltinDiscovery bool     `json:"enable_builtin_discovery"`
	EnableFileDiscovery    bool     `json:"enable_file_discovery"`
	SearchPaths           []string `json:"search_paths"`
	PluginPrefix          string   `json:"plugin_prefix"`
	ValidationEnabled     bool     `json:"validation_enabled"`
}

// DefaultDiscoveryConfiguration returns default discovery configuration.
func DefaultDiscoveryConfiguration() DiscoveryConfiguration {
	return DiscoveryConfiguration{
		EnableBuiltinDiscovery: true,
		EnableFileDiscovery:    false, // Disabled by default for security
		SearchPaths:           []string{},
		PluginPrefix:          "engx-",
		ValidationEnabled:     true,
	}
}

// Future expansion capabilities:

// FileDiscoverer would handle discovery of plugins from filesystem
// type FileDiscoverer struct {
//     searchPaths []string
//     prefix      string
// }

// NetworkDiscoverer would handle discovery of plugins from network sources
// type NetworkDiscoverer struct {
//     registryURL string
//     credentials string
// }

// EnvironmentDiscoverer would handle discovery based on environment variables
// type EnvironmentDiscoverer struct {
//     envPrefix string
// }