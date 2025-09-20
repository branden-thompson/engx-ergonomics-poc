package registry

import (
	"testing"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common"
	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
	"github.com/spf13/cobra"
)

// MockPlugin provides a test implementation of CommandPlugin
type MockPlugin struct {
	name             string
	description      string
	version          string
	requiredServices []string
	optionalServices []string
	initError        error
	cleanupError     error
	command          *cobra.Command
}

func (m *MockPlugin) Name() string        { return m.name }
func (m *MockPlugin) Description() string { return m.description }
func (m *MockPlugin) Version() string     { return m.version }

func (m *MockPlugin) Create(deps interface{}) *cobra.Command {
	if m.command == nil {
		m.command = &cobra.Command{
			Use:   m.name,
			Short: m.description,
		}
	}
	return m.command
}

func (m *MockPlugin) Initialize() error        { return m.initError }
func (m *MockPlugin) Cleanup() error           { return m.cleanupError }
func (m *MockPlugin) RequiredServices() []string { return m.requiredServices }
func (m *MockPlugin) OptionalServices() []string { return m.optionalServices }

// CreateMockDependencies creates mock dependencies for testing
func CreateMockDependencies() *common.Dependencies {
	return &common.Dependencies{
		Logger: &MockLogger{},
	}
}

type MockLogger struct{}

func (ml *MockLogger) Debug(msg string, args ...interface{}) {}
func (ml *MockLogger) Info(msg string, args ...interface{})  {}
func (ml *MockLogger) Warn(msg string, args ...interface{})  {}
func (ml *MockLogger) Error(msg string, args ...interface{}) {}
func (ml *MockLogger) WithContext(ctx map[string]interface{}) interfaces.Logger { return ml }
func (ml *MockLogger) WithPlugin(pluginName string) interfaces.Logger           { return ml }
func (ml *MockLogger) WithComponent(componentName string) interfaces.Logger     { return ml }
func (ml *MockLogger) SetLevel(level interfaces.LogLevel)                       {}
func (ml *MockLogger) GetLevel() interfaces.LogLevel { return interfaces.LogLevelInfo }

func TestNewManager(t *testing.T) {
	deps := CreateMockDependencies()
	manager := NewManager(deps)

	if manager == nil {
		t.Fatal("NewManager returned nil")
	}

	if manager.deps != deps {
		t.Error("Manager dependencies not set correctly")
	}

	if manager.Count() != 0 {
		t.Errorf("Expected 0 plugins initially, got %d", manager.Count())
	}
}

func TestRegisterPlugin(t *testing.T) {
	deps := CreateMockDependencies()
	manager := NewManager(deps)

	plugin := &MockPlugin{
		name:        "test-plugin",
		description: "Test plugin for unit testing",
		version:     "1.0.0",
	}

	err := manager.Register(plugin)
	if err != nil {
		t.Fatalf("Failed to register plugin: %v", err)
	}

	if manager.Count() != 1 {
		t.Errorf("Expected 1 plugin after registration, got %d", manager.Count())
	}

	// Test getting the plugin
	retrievedPlugin, exists := manager.GetPlugin("test-plugin")
	if !exists {
		t.Error("Plugin not found after registration")
	}

	if retrievedPlugin != plugin {
		t.Error("Retrieved plugin is not the same as registered plugin")
	}
}

func TestRegisterDuplicatePlugin(t *testing.T) {
	deps := CreateMockDependencies()
	manager := NewManager(deps)

	plugin1 := &MockPlugin{
		name:        "duplicate",
		description: "First plugin with duplicate name",
		version:     "1.0.0",
	}

	plugin2 := &MockPlugin{
		name:        "duplicate",
		description: "Second plugin with duplicate name",
		version:     "2.0.0",
	}

	// Register first plugin
	err := manager.Register(plugin1)
	if err != nil {
		t.Fatalf("Failed to register first plugin: %v", err)
	}

	// Try to register second plugin with same name
	err = manager.Register(plugin2)
	if err == nil {
		t.Error("Expected error when registering duplicate plugin, but got none")
	}

	if manager.Count() != 1 {
		t.Errorf("Expected 1 plugin after duplicate registration attempt, got %d", manager.Count())
	}
}

func TestRegisterPluginWithInitError(t *testing.T) {
	deps := CreateMockDependencies()
	manager := NewManager(deps)

	plugin := &MockPlugin{
		name:        "failing-plugin",
		description: "Plugin that fails initialization",
		version:     "1.0.0",
		initError:   &ValidationError{Field: "test", Message: "initialization failed"},
	}

	err := manager.Register(plugin)
	if err == nil {
		t.Error("Expected error when registering plugin with init error, but got none")
	}

	if manager.Count() != 0 {
		t.Errorf("Expected 0 plugins after failed registration, got %d", manager.Count())
	}
}

func TestUnregisterPlugin(t *testing.T) {
	deps := CreateMockDependencies()
	manager := NewManager(deps)

	plugin := &MockPlugin{
		name:        "removable-plugin",
		description: "Plugin to be removed",
		version:     "1.0.0",
	}

	// Register plugin
	err := manager.Register(plugin)
	if err != nil {
		t.Fatalf("Failed to register plugin: %v", err)
	}

	// Unregister plugin
	err = manager.Unregister("removable-plugin")
	if err != nil {
		t.Fatalf("Failed to unregister plugin: %v", err)
	}

	if manager.Count() != 0 {
		t.Errorf("Expected 0 plugins after unregistration, got %d", manager.Count())
	}

	// Try to get unregistered plugin
	_, exists := manager.GetPlugin("removable-plugin")
	if exists {
		t.Error("Plugin still exists after unregistration")
	}
}

func TestUnregisterNonexistentPlugin(t *testing.T) {
	deps := CreateMockDependencies()
	manager := NewManager(deps)

	err := manager.Unregister("nonexistent-plugin")
	if err == nil {
		t.Error("Expected error when unregistering nonexistent plugin, but got none")
	}
}

func TestGetCommands(t *testing.T) {
	deps := CreateMockDependencies()
	manager := NewManager(deps)

	plugin1 := &MockPlugin{
		name:        "command1",
		description: "First command plugin",
		version:     "1.0.0",
	}

	plugin2 := &MockPlugin{
		name:        "command2",
		description: "Second command plugin",
		version:     "1.0.0",
	}

	// Register plugins
	manager.Register(plugin1)
	manager.Register(plugin2)

	commands := manager.GetCommands()
	if len(commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(commands))
	}

	// Check that commands are created correctly
	for i, cmd := range commands {
		if cmd == nil {
			t.Errorf("Command %d is nil", i)
		}
	}
}

func TestListPlugins(t *testing.T) {
	deps := CreateMockDependencies()
	manager := NewManager(deps)

	plugin := &MockPlugin{
		name:        "listed-plugin",
		description: "Plugin to be listed",
		version:     "1.2.3",
	}

	manager.Register(plugin)

	plugins := manager.ListPlugins()
	if len(plugins) != 1 {
		t.Errorf("Expected 1 plugin in list, got %d", len(plugins))
	}

	info := plugins[0]
	if info.Metadata.Name != "listed-plugin" {
		t.Errorf("Expected plugin name 'listed-plugin', got '%s'", info.Metadata.Name)
	}

	if info.Status != interfaces.PluginStatusActive {
		t.Errorf("Expected plugin status %d, got %d", interfaces.PluginStatusActive, info.Status)
	}
}

func TestGetRegisteredNames(t *testing.T) {
	deps := CreateMockDependencies()
	manager := NewManager(deps)

	plugin1 := &MockPlugin{name: "alpha", description: "Alpha plugin", version: "1.0.0"}
	plugin2 := &MockPlugin{name: "beta", description: "Beta plugin", version: "1.0.0"}

	manager.Register(plugin1)
	manager.Register(plugin2)

	names := manager.GetRegisteredNames()
	if len(names) != 2 {
		t.Errorf("Expected 2 registered names, got %d", len(names))
	}

	// Check that order is preserved (registration order)
	if names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("Expected names [alpha, beta], got %v", names)
	}

	// Modify returned slice should not affect internal state
	names[0] = "modified"
	internalNames := manager.GetRegisteredNames()
	if internalNames[0] != "alpha" {
		t.Error("Internal names were modified by external slice modification")
	}
}

func TestGetStats(t *testing.T) {
	deps := CreateMockDependencies()
	manager := NewManager(deps)

	plugin1 := &MockPlugin{name: "stats1", description: "Stats plugin 1", version: "1.0.0"}
	plugin2 := &MockPlugin{name: "stats2", description: "Stats plugin 2", version: "1.0.0"}

	manager.Register(plugin1)
	manager.Register(plugin2)

	stats := manager.GetStats()
	if stats.TotalPlugins != 2 {
		t.Errorf("Expected 2 total plugins in stats, got %d", stats.TotalPlugins)
	}

	if stats.RegisteredPlugins != 2 {
		t.Errorf("Expected 2 registered plugins in stats, got %d", stats.RegisteredPlugins)
	}

	if len(stats.PluginNames) != 2 {
		t.Errorf("Expected 2 plugin names in stats, got %d", len(stats.PluginNames))
	}
}

func TestCleanup(t *testing.T) {
	deps := CreateMockDependencies()
	manager := NewManager(deps)

	plugin1 := &MockPlugin{name: "cleanup1", description: "Cleanup plugin 1", version: "1.0.0"}
	plugin2 := &MockPlugin{name: "cleanup2", description: "Cleanup plugin 2", version: "1.0.0"}

	manager.Register(plugin1)
	manager.Register(plugin2)

	err := manager.Cleanup()
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	if manager.Count() != 0 {
		t.Errorf("Expected 0 plugins after cleanup, got %d", manager.Count())
	}

	names := manager.GetRegisteredNames()
	if len(names) != 0 {
		t.Errorf("Expected 0 registered names after cleanup, got %d", len(names))
	}
}

func TestSafeRegister(t *testing.T) {
	deps := CreateMockDependencies()
	manager := NewManager(deps)

	// Create a plugin that will cause validation error
	invalidPlugin := &MockPlugin{
		name:        "", // Invalid empty name
		description: "Invalid plugin",
		version:     "1.0.0",
	}

	// SafeRegister should handle the error gracefully
	err := manager.SafeRegister(invalidPlugin)
	if err == nil {
		t.Error("Expected error from SafeRegister with invalid plugin, but got none")
	}

	// Manager should still be functional
	validPlugin := &MockPlugin{
		name:        "valid-plugin",
		description: "Valid plugin after error",
		version:     "1.0.0",
	}

	err = manager.SafeRegister(validPlugin)
	if err != nil {
		t.Fatalf("Failed to register valid plugin after error: %v", err)
	}

	if manager.Count() != 1 {
		t.Errorf("Expected 1 plugin after safe register operations, got %d", manager.Count())
	}
}