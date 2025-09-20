package common

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
)

// SimpleConfigUI provides basic configuration management functionality
type SimpleConfigUI struct {
	configManager *PluginConfigManager
	logger        interfaces.Logger
}

// NewSimpleConfigUI creates a new simple configuration UI
func NewSimpleConfigUI(configManager *PluginConfigManager, deps *Dependencies) *SimpleConfigUI {
	return &SimpleConfigUI{
		configManager: configManager,
		logger:        deps.Logger.WithComponent("config-ui"),
	}
}

// ShowPluginList displays the plugin configuration list
func (cui *SimpleConfigUI) ShowPluginList() error {
	configs := cui.configManager.ListConfigs()

	fmt.Println("📦 Plugin Configurations")
	fmt.Println(strings.Repeat("=", 80))

	if len(configs) == 0 {
		fmt.Println("No plugin configurations found.")
		fmt.Println("")
		fmt.Println("Available commands:")
		fmt.Println("  engx dev config create <plugin> - Create new plugin configuration")
		fmt.Println("  engx dev config wizard <plugin> - Run configuration wizard")
		return nil
	}

	// Table header
	fmt.Printf("%-20s %-10s %-10s %-20s %-15s\n", "PLUGIN", "VERSION", "ENABLED", "LAST UPDATED", "SETTINGS")
	fmt.Println(strings.Repeat("-", 80))

	// Sort plugins by name
	var pluginNames []string
	for name := range configs {
		pluginNames = append(pluginNames, name)
	}
	sort.Strings(pluginNames)

	// Display plugins
	for _, name := range pluginNames {
		config := configs[name]
		enabledStr := "❌"
		if config.Enabled {
			enabledStr = "✅"
		}

		lastUpdated := config.LastUpdated.Format("2006-01-02 15:04")
		settingsCount := len(config.Settings)

		fmt.Printf("%-20s %-10s %-10s %-20s %-15s\n",
			name, config.Version, enabledStr, lastUpdated,
			fmt.Sprintf("%d settings", settingsCount))
	}

	fmt.Println("")
	fmt.Println("💡 Use 'engx dev config edit <plugin>' to configure a plugin")
	fmt.Println("💡 Use 'engx dev config wizard <plugin>' for guided setup")

	return nil
}

// ShowStats displays configuration statistics
func (cui *SimpleConfigUI) ShowStats() error {
	stats := cui.configManager.GetStats()

	fmt.Println("📊 Plugin Configuration Statistics")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("")

	fmt.Printf("Total Configurations: %v\n", stats["total_configs"])
	fmt.Printf("Enabled Plugins: %v\n", stats["enabled_configs"])
	fmt.Printf("Disabled Plugins: %v\n", stats["disabled_configs"])
	fmt.Printf("Total Settings: %v\n", stats["total_settings"])
	fmt.Printf("Schemas Loaded: %v\n", stats["schemas_loaded"])
	fmt.Printf("Auto Save: %v\n", stats["auto_save"])
	fmt.Printf("Config Path: %v\n", stats["config_path"])
	fmt.Printf("Active Watchers: %v\n", stats["watchers"])

	return nil
}

// CreatePluginConfig creates a new plugin configuration with basic setup
func (cui *SimpleConfigUI) CreatePluginConfig(pluginName string) error {
	fmt.Printf("Creating configuration for plugin: %s\n", pluginName)

	// Check if config already exists
	if _, err := cui.configManager.GetConfig(pluginName); err == nil {
		fmt.Printf("⚠️  Configuration for %s already exists. Use 'edit' to modify.\n", pluginName)
		return nil
	}

	// Create basic configuration
	config := &PluginConfig{
		PluginName:  pluginName,
		Version:     "1.0.0",
		Enabled:     true,
		Settings:    make(map[string]interface{}),
		Defaults:    make(map[string]interface{}),
		Environment: make(map[string]interface{}),
		Flags:       make(map[string]interface{}),
		Metadata: PluginConfigMetadata{
			Description: fmt.Sprintf("Configuration for %s plugin", pluginName),
		},
	}

	if err := cui.configManager.SaveConfig(pluginName, config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✅ Configuration created successfully for %s!\n", pluginName)
	fmt.Printf("📄 Config file: %s\n", cui.configManager.GetConfigPath(pluginName))
	fmt.Printf("💡 Use 'engx dev config edit %s' to add settings\n", pluginName)

	return nil
}

// ShowPluginConfig displays configuration for a specific plugin
func (cui *SimpleConfigUI) ShowPluginConfig(pluginName string) error {
	config, err := cui.configManager.GetConfig(pluginName)
	if err != nil {
		return fmt.Errorf("failed to get config for %s: %w", pluginName, err)
	}

	fmt.Printf("📦 Plugin Configuration: %s\n", pluginName)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Version: %s\n", config.Version)
	fmt.Printf("Enabled: %v\n", config.Enabled)
	fmt.Printf("Last Updated: %s\n", config.LastUpdated.Format("2006-01-02 15:04:05"))
	if config.UpdatedBy != "" {
		fmt.Printf("Updated By: %s\n", config.UpdatedBy)
	}

	if config.Metadata.Description != "" {
		fmt.Printf("Description: %s\n", config.Metadata.Description)
	}

	fmt.Println("")

	// Show settings
	if len(config.Settings) > 0 {
		fmt.Println("⚙️  Settings:")
		fmt.Println(strings.Repeat("-", 40))
		for key, value := range config.Settings {
			fmt.Printf("  %s: %v\n", key, value)
		}
		fmt.Println("")
	} else {
		fmt.Println("⚙️  No settings configured")
		fmt.Println("")
	}

	// Show defaults if any
	if len(config.Defaults) > 0 {
		fmt.Println("🔧 Defaults:")
		fmt.Println(strings.Repeat("-", 40))
		for key, value := range config.Defaults {
			fmt.Printf("  %s: %v\n", key, value)
		}
		fmt.Println("")
	}

	fmt.Printf("💡 Use 'engx dev config set %s <key> <value>' to add settings\n", pluginName)
	fmt.Printf("💡 Use 'engx dev config get %s <key>' to view specific settings\n", pluginName)

	return nil
}

// InteractiveConfigWizard provides a simplified wizard experience
func (cui *SimpleConfigUI) InteractiveConfigWizard(pluginName string) error {
	fmt.Printf("🧙 Configuration Wizard for %s\n", pluginName)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("")

	// Check if config exists
	config, err := cui.configManager.GetConfig(pluginName)
	if err != nil {
		// Create new config
		config = &PluginConfig{
			PluginName:  pluginName,
			Version:     "1.0.0",
			Enabled:     true,
			Settings:    make(map[string]interface{}),
			Defaults:    make(map[string]interface{}),
			Environment: make(map[string]interface{}),
			Flags:       make(map[string]interface{}),
			Metadata: PluginConfigMetadata{
				Description: fmt.Sprintf("Configuration for %s plugin", pluginName),
			},
		}
	}

	fmt.Printf("Current configuration for %s:\n", pluginName)
	fmt.Printf("  Enabled: %v\n", config.Enabled)
	fmt.Printf("  Settings: %d configured\n", len(config.Settings))
	fmt.Println("")

	fmt.Println("The configuration wizard will help you set up common settings.")
	fmt.Println("For advanced configuration, use the 'edit' command or modify the config file directly.")
	fmt.Println("")

	// Add a few common settings examples
	commonSettings := map[string]interface{}{
		"debug":   false,
		"timeout": 30,
		"output":  "text",
	}

	fmt.Println("Adding common default settings:")
	for key, value := range commonSettings {
		if _, exists := config.Settings[key]; !exists {
			config.Settings[key] = value
			fmt.Printf("  ✅ %s: %v\n", key, value)
		} else {
			fmt.Printf("  ⚠️  %s: already configured\n", key)
		}
	}

	// Save configuration
	if err := cui.configManager.SaveConfig(pluginName, config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("")
	fmt.Printf("✅ Configuration wizard completed for %s!\n", pluginName)
	fmt.Printf("📄 Config file: %s\n", cui.configManager.GetConfigPath(pluginName))
	fmt.Printf("💡 Use 'engx dev config get %s <key>' to view settings\n", pluginName)
	fmt.Printf("💡 Use 'engx dev config set %s <key> <value>' to modify settings\n", pluginName)

	return nil
}