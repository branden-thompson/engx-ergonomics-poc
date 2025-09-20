package interfaces

import (
	"github.com/spf13/cobra"
)

// CommandPlugin defines the interface that all command plugins must implement.
// This interface provides a standardized way to create, register, and manage commands
// within the engx CLI tool ecosystem.
type CommandPlugin interface {
	// Core plugin information
	Name() string        // Unique identifier for the plugin (e.g., "create", "update")
	Description() string // Human-readable description of what the command does
	Version() string     // Plugin version for compatibility checking

	// Command creation and lifecycle
	Create(deps interface{}) *cobra.Command // Creates the cobra command with dependencies injected
	Initialize() error                        // Initialize plugin resources (called once during registration)
	Cleanup() error                          // Cleanup plugin resources (called during shutdown)

	// Dependency management
	RequiredServices() []string // List of required service names (e.g., "config", "tui", "chaos")
	OptionalServices() []string // List of optional service names that enhance functionality
}

// PluginMetadata provides additional information about a plugin
// for discovery, validation, and management purposes.
type PluginMetadata struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Version      string            `json:"version"`
	Author       string            `json:"author,omitempty"`
	Homepage     string            `json:"homepage,omitempty"`
	License      string            `json:"license,omitempty"`
	Tags         []string          `json:"tags,omitempty"`
	Dependencies map[string]string `json:"dependencies,omitempty"` // service -> min version
}

// PluginStatus represents the current state of a plugin
type PluginStatus int

const (
	PluginStatusUnknown PluginStatus = iota
	PluginStatusRegistered
	PluginStatusInitialized
	PluginStatusActive
	PluginStatusError
	PluginStatusDisabled
)

func (ps PluginStatus) String() string {
	switch ps {
	case PluginStatusRegistered:
		return "registered"
	case PluginStatusInitialized:
		return "initialized"
	case PluginStatusActive:
		return "active"
	case PluginStatusError:
		return "error"
	case PluginStatusDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}

// PluginInfo combines plugin metadata with runtime status information
type PluginInfo struct {
	Metadata PluginMetadata `json:"metadata"`
	Status   PluginStatus   `json:"status"`
	Error    string         `json:"error,omitempty"`
}

// AdvancedCommandPlugin extends the basic CommandPlugin interface
// with additional capabilities for more sophisticated plugins.
type AdvancedCommandPlugin interface {
	CommandPlugin

	// Metadata and validation
	GetMetadata() PluginMetadata // Returns detailed plugin metadata
	Validate() error            // Validates plugin configuration and dependencies
	HealthCheck() error         // Performs health check of plugin components

	// Configuration management
	LoadConfig(configData map[string]interface{}) error // Load plugin-specific configuration
	GetConfigSchema() map[string]interface{}            // Returns JSON schema for configuration

	// Event handling
	OnBeforeExecute(cmd *cobra.Command, args []string) error // Called before command execution
	OnAfterExecute(cmd *cobra.Command, args []string) error  // Called after command execution
	OnError(cmd *cobra.Command, err error) error            // Called when command encounters error
}