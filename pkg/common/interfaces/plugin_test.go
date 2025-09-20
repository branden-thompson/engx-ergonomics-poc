package interfaces

import (
	"testing"

	"github.com/spf13/cobra"
)

// MockCommandPlugin provides a mock implementation of CommandPlugin for testing
type MockCommandPlugin struct {
	NameValue            string
	DescriptionValue     string
	VersionValue         string
	RequiredServicesList []string
	OptionalServicesList []string
	InitializeError      error
	CleanupError         error
	CommandCreated       *cobra.Command
}

func (m *MockCommandPlugin) Name() string {
	return m.NameValue
}

func (m *MockCommandPlugin) Description() string {
	return m.DescriptionValue
}

func (m *MockCommandPlugin) Version() string {
	return m.VersionValue
}

func (m *MockCommandPlugin) Create(deps interface{}) *cobra.Command {
	if m.CommandCreated == nil {
		m.CommandCreated = &cobra.Command{
			Use:   m.NameValue,
			Short: m.DescriptionValue,
		}
	}
	return m.CommandCreated
}

func (m *MockCommandPlugin) Initialize() error {
	return m.InitializeError
}

func (m *MockCommandPlugin) Cleanup() error {
	return m.CleanupError
}

func (m *MockCommandPlugin) RequiredServices() []string {
	return m.RequiredServicesList
}

func (m *MockCommandPlugin) OptionalServices() []string {
	return m.OptionalServicesList
}

// TestCommandPluginInterface verifies that the CommandPlugin interface
// can be implemented correctly
func TestCommandPluginInterface(t *testing.T) {
	plugin := &MockCommandPlugin{
		NameValue:            "test",
		DescriptionValue:     "Test plugin",
		VersionValue:         "1.0.0",
		RequiredServicesList: []string{"config", "logger"},
		OptionalServicesList: []string{"tui"},
	}

	// Test basic interface methods
	if plugin.Name() != "test" {
		t.Errorf("Expected name 'test', got '%s'", plugin.Name())
	}

	if plugin.Description() != "Test plugin" {
		t.Errorf("Expected description 'Test plugin', got '%s'", plugin.Description())
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", plugin.Version())
	}

	// Test service dependencies
	required := plugin.RequiredServices()
	if len(required) != 2 || required[0] != "config" || required[1] != "logger" {
		t.Errorf("Expected required services [config, logger], got %v", required)
	}

	optional := plugin.OptionalServices()
	if len(optional) != 1 || optional[0] != "tui" {
		t.Errorf("Expected optional services [tui], got %v", optional)
	}

	// Test lifecycle methods
	if err := plugin.Initialize(); err != nil {
		t.Errorf("Expected no error from Initialize(), got %v", err)
	}

	if err := plugin.Cleanup(); err != nil {
		t.Errorf("Expected no error from Cleanup(), got %v", err)
	}

	// Test command creation
	cmd := plugin.Create(nil) // Dependencies can be nil for this test
	if cmd == nil {
		t.Error("Expected command to be created, got nil")
	}

	if cmd.Use != "test" {
		t.Errorf("Expected command Use to be 'test', got '%s'", cmd.Use)
	}

	if cmd.Short != "Test plugin" {
		t.Errorf("Expected command Short to be 'Test plugin', got '%s'", cmd.Short)
	}
}

// TestPluginStatus tests the PluginStatus enum and string conversion
func TestPluginStatus(t *testing.T) {
	tests := []struct {
		status   PluginStatus
		expected string
	}{
		{PluginStatusUnknown, "unknown"},
		{PluginStatusRegistered, "registered"},
		{PluginStatusInitialized, "initialized"},
		{PluginStatusActive, "active"},
		{PluginStatusError, "error"},
		{PluginStatusDisabled, "disabled"},
	}

	for _, test := range tests {
		if test.status.String() != test.expected {
			t.Errorf("Expected status %d to be '%s', got '%s'",
				test.status, test.expected, test.status.String())
		}
	}
}

// TestVerbosityLevel tests the VerbosityLevel enum and string conversion
func TestVerbosityLevel(t *testing.T) {
	tests := []struct {
		level    VerbosityLevel
		expected string
	}{
		{VerbosityQuiet, "quiet"},
		{VerbosityConcise, "concise"},
		{VerbosityNormal, "normal"},
		{VerbosityVerbose, "verbose"},
		{VerbosityDebug, "debug"},
	}

	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Errorf("Expected verbosity %d to be '%s', got '%s'",
				test.level, test.expected, test.level.String())
		}
	}
}

// TestLogLevel tests the LogLevel enum and string conversion
func TestLogLevel(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LogLevelDebug, "debug"},
		{LogLevelInfo, "info"},
		{LogLevelWarn, "warn"},
		{LogLevelError, "error"},
	}

	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Errorf("Expected log level %d to be '%s', got '%s'",
				test.level, test.expected, test.level.String())
		}
	}
}

// TestAggressivenessLevel tests the AggressivenessLevel enum and string conversion
func TestAggressivenessLevel(t *testing.T) {
	tests := []struct {
		level    AggressivenessLevel
		expected string
	}{
		{AggressivenessOff, "off"},
		{AggressivenessDefault, "default"},
		{AggressivenessScout, "scout"},
		{AggressivenessAggressive, "aggressive"},
		{AggressivenessInvasive, "invasive"},
		{AggressivenessApocalyptic, "apocalyptic"},
	}

	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Errorf("Expected aggressiveness %d to be '%s', got '%s'",
				test.level, test.expected, test.level.String())
		}
	}
}

// TestPluginInfo tests the PluginInfo structure
func TestPluginInfo(t *testing.T) {
	metadata := PluginMetadata{
		Name:        "test-plugin",
		Description: "Test plugin for unit testing",
		Version:     "1.0.0",
		Author:      "Test Author",
		License:     "MIT",
		Tags:        []string{"test", "example"},
	}

	info := PluginInfo{
		Metadata: metadata,
		Status:   PluginStatusActive,
		Error:    "",
	}

	if info.Metadata.Name != "test-plugin" {
		t.Errorf("Expected plugin name 'test-plugin', got '%s'", info.Metadata.Name)
	}

	if info.Status != PluginStatusActive {
		t.Errorf("Expected status %d, got %d", PluginStatusActive, info.Status)
	}

	if info.Error != "" {
		t.Errorf("Expected empty error, got '%s'", info.Error)
	}
}