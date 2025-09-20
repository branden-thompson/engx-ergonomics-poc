package example_advanced

import (
	"testing"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common"
	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
)

func TestPlugin_Basic(t *testing.T) {
	// Create test dependencies
	deps := common.NewDependencies()

	// Create plugin
	plugin := NewPlugin(deps)

	// Test basic interface compliance
	if plugin.Name() != "example-advanced" {
		t.Errorf("Expected name 'example-advanced', got '%s'", plugin.Name())
	}

	if plugin.Description() == "" {
		t.Error("Description should not be empty")
	}

	if plugin.Version() == "" {
		t.Error("Version should not be empty")
	}

	// Test initialization
	if err := plugin.Initialize(); err != nil {
		t.Errorf("Initialize failed: %v", err)
	}

	// Test cleanup
	if err := plugin.Cleanup(); err != nil {
		t.Errorf("Cleanup failed: %v", err)
	}
}

func TestPlugin_Dependencies(t *testing.T) {
	deps := common.NewDependencies()
	plugin := NewPlugin(deps)

	requiredServices := plugin.RequiredServices()
	expectedRequired := []string{"config", "logger"}

	if len(requiredServices) != len(expectedRequired) {
		t.Errorf("Expected %d required services, got %d", len(expectedRequired), len(requiredServices))
	}

	optionalServices := plugin.OptionalServices()
	expectedOptional := []string{"filesystem", "tui"}

	if len(optionalServices) != len(expectedOptional) {
		t.Errorf("Expected %d optional services, got %d", len(expectedOptional), len(optionalServices))
	}
}

func TestPlugin_Command(t *testing.T) {
	deps := common.NewDependencies()
	plugin := NewPlugin(deps)

	cmd := plugin.Create(deps)
	if cmd == nil {
		t.Error("Create should return a valid command")
	}

	if cmd.Use == "" {
		t.Error("Command Use should not be empty")
	}

	if cmd.Short == "" {
		t.Error("Command Short should not be empty")
	}
}


func TestPlugin_AdvancedInterface(t *testing.T) {
	deps := common.NewDependencies()
	plugin := NewPlugin(deps)

	// Test if plugin implements AdvancedCommandPlugin
	advancedPlugin, ok := plugin.(interfaces.AdvancedCommandPlugin)
	if !ok {
		t.Error("Plugin should implement AdvancedCommandPlugin interface")
		return
	}

	// Test metadata
	metadata := advancedPlugin.GetMetadata()
	if metadata.Name == "" {
		t.Error("Metadata name should not be empty")
	}

	// Test validation
	if err := advancedPlugin.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}

	// Test health check
	if err := advancedPlugin.HealthCheck(); err != nil {
		t.Errorf("Health check failed: %v", err)
	}

	// Test config schema
	schema := advancedPlugin.GetConfigSchema()
	if schema == nil {
		t.Error("Config schema should not be nil")
	}
}

