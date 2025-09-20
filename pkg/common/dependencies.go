package common

import (
	"fmt"
	"log"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
)

// Dependencies represents the collection of shared services available to plugins.
// This struct implements the dependency injection pattern, providing plugins
// with access to common functionality while maintaining loose coupling.
type Dependencies struct {
	Config     interfaces.ConfigManager
	TUI        interfaces.TUIRegistry
	Chaos      interfaces.ChaosInjector
	AAR        interfaces.AARGenerator
	Logger     interfaces.Logger
	Filesystem interfaces.FilesystemManager
}

// NewDependencies creates a new Dependencies instance with all services initialized.
// This function serves as the main dependency injection factory for the plugin system.
func NewDependencies() *Dependencies {
	deps := &Dependencies{
		Logger: NewDefaultLogger(),
	}

	// Initialize services in dependency order
	deps.Config = NewConfigManager(deps.Logger)
	deps.Filesystem = NewFilesystemManager(deps.Logger)
	deps.TUI = NewTUIRegistry(deps.Logger)
	deps.Chaos = NewChaosInjector(deps.Config, deps.Logger)
	deps.AAR = NewAARGenerator(deps.Config, deps.Logger)

	return deps
}

// Validate ensures all required dependencies are properly initialized
func (d *Dependencies) Validate() error {
	if d.Config == nil {
		return ErrMissingDependency{Service: "Config"}
	}
	if d.TUI == nil {
		return ErrMissingDependency{Service: "TUI"}
	}
	if d.Chaos == nil {
		return ErrMissingDependency{Service: "Chaos"}
	}
	if d.AAR == nil {
		return ErrMissingDependency{Service: "AAR"}
	}
	if d.Logger == nil {
		return ErrMissingDependency{Service: "Logger"}
	}
	if d.Filesystem == nil {
		return ErrMissingDependency{Service: "Filesystem"}
	}

	return nil
}

// GetService returns a service by name for dynamic dependency resolution
func (d *Dependencies) GetService(name string) interface{} {
	switch name {
	case "config":
		return d.Config
	case "tui":
		return d.TUI
	case "chaos":
		return d.Chaos
	case "aar":
		return d.AAR
	case "logger":
		return d.Logger
	case "filesystem":
		return d.Filesystem
	default:
		return nil
	}
}

// HasService checks if a service is available
func (d *Dependencies) HasService(name string) bool {
	return d.GetService(name) != nil
}

// ListServices returns a list of all available service names
func (d *Dependencies) ListServices() []string {
	return []string{"config", "tui", "chaos", "aar", "logger", "filesystem"}
}

// Cleanup performs cleanup of all managed services
func (d *Dependencies) Cleanup() error {
	var errors []error

	// Cleanup services in reverse dependency order
	if cleaner, ok := d.AAR.(interface{ Cleanup() error }); ok {
		if err := cleaner.Cleanup(); err != nil {
			errors = append(errors, err)
		}
	}

	if cleaner, ok := d.Chaos.(interface{ Cleanup() error }); ok {
		if err := cleaner.Cleanup(); err != nil {
			errors = append(errors, err)
		}
	}

	if cleaner, ok := d.TUI.(interface{ Cleanup() error }); ok {
		if err := cleaner.Cleanup(); err != nil {
			errors = append(errors, err)
		}
	}

	if cleaner, ok := d.Filesystem.(interface{ Cleanup() error }); ok {
		if err := cleaner.Cleanup(); err != nil {
			errors = append(errors, err)
		}
	}

	if cleaner, ok := d.Config.(interface{ Cleanup() error }); ok {
		if err := cleaner.Cleanup(); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return &MultiError{Errors: errors}
	}

	return nil
}

// ErrMissingDependency indicates a required dependency is missing
type ErrMissingDependency struct {
	Service string
}

func (e ErrMissingDependency) Error() string {
	return "missing required dependency: " + e.Service
}

// MultiError represents multiple errors that occurred during operation
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
	return "multiple errors occurred: " + me.Errors[0].Error() + " (and " +
		   fmt.Sprintf("%d", len(me.Errors)-1) + " more)"
}

// Factory functions for creating service implementations
// These will be implemented as we migrate existing services

func NewDefaultLogger() interfaces.Logger {
	return &DefaultLogger{logger: log.Default()}
}


// DefaultLogger provides a basic logger implementation
type DefaultLogger struct {
	logger *log.Logger
}

func (dl *DefaultLogger) Debug(msg string, args ...interface{}) {
	dl.logger.Printf("[DEBUG] "+msg, args...)
}

func (dl *DefaultLogger) Info(msg string, args ...interface{}) {
	dl.logger.Printf("[INFO] "+msg, args...)
}

func (dl *DefaultLogger) Warn(msg string, args ...interface{}) {
	dl.logger.Printf("[WARN] "+msg, args...)
}

func (dl *DefaultLogger) Error(msg string, args ...interface{}) {
	dl.logger.Printf("[ERROR] "+msg, args...)
}

func (dl *DefaultLogger) WithContext(ctx map[string]interface{}) interfaces.Logger {
	return dl // Simple implementation - could be enhanced
}

func (dl *DefaultLogger) WithPlugin(pluginName string) interfaces.Logger {
	return dl // Simple implementation - could be enhanced
}

func (dl *DefaultLogger) WithComponent(componentName string) interfaces.Logger {
	return dl // Simple implementation - could be enhanced
}

func (dl *DefaultLogger) SetLevel(level interfaces.LogLevel) {
	// Simple implementation - could be enhanced to actually filter logs
}

func (dl *DefaultLogger) GetLevel() interfaces.LogLevel {
	return interfaces.LogLevelInfo // Default level
}