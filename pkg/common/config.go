package common

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
	"gopkg.in/yaml.v3"
)

// ConfigManager implements the interfaces.ConfigManager interface
type ConfigManager struct {
	config     *Config
	loader     *Loader
	logger     interfaces.Logger
	globalPath string
}

// Config represents the complete configuration structure
type Config struct {
	Project      *ProjectConfig              `yaml:"project,omitempty"`
	Defaults     *DefaultsConfig             `yaml:"defaults,omitempty"`
	Environments map[string]*EnvConfig       `yaml:"environments,omitempty"`
	Commands     map[string]*CmdConfig       `yaml:"custom_commands,omitempty"`
}

// ProjectConfig contains project-specific settings
type ProjectConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description,omitempty"`
	Repository  string `yaml:"repository,omitempty"`
}

// DefaultsConfig contains default behavior settings
type DefaultsConfig struct {
	Verbosity        string        `yaml:"verbosity"`         // normal, verbose, quiet
	DeploymentTarget string        `yaml:"deployment_target"` // development, staging, production
	Timeout          time.Duration `yaml:"timeout"`
	Theme            string        `yaml:"theme"`             // auto, dark, light
	Template         string        `yaml:"template,omitempty"` // typescript, javascript, minimal
}

// EnvConfig contains environment-specific settings
type EnvConfig struct {
	APIEndpoint string            `yaml:"api_endpoint"`
	Variables   map[string]string `yaml:"variables,omitempty"`
	Timeout     time.Duration     `yaml:"timeout,omitempty"`
}

// CmdConfig contains custom command configurations
type CmdConfig struct {
	Flags []string          `yaml:"flags,omitempty"`
	Env   map[string]string `yaml:"env,omitempty"`
}

// Loader handles configuration loading with inheritance
type Loader struct {
	globalPath  string
	projectPath string
}

// NewConfigManager creates a new ConfigManager instance
func NewConfigManager(logger interfaces.Logger) interfaces.ConfigManager {
	homeDir, _ := os.UserHomeDir()
	globalPath := filepath.Join(homeDir, ".engx", "config.yaml")
	projectPath := ".engx/config.yaml"

	loader := &Loader{
		globalPath:  globalPath,
		projectPath: projectPath,
	}

	return &ConfigManager{
		loader:     loader,
		logger:     logger.WithComponent("config"),
		globalPath: globalPath,
	}
}

// LoadConfig loads configuration from the specified path or default locations
func (cm *ConfigManager) LoadConfig(path string) (*interfaces.Config, error) {
	var config *Config
	var err error

	if path != "" {
		config, err = cm.loader.LoadWithCustomPath(path)
	} else {
		config, err = cm.loader.Load()
	}

	if err != nil {
		cm.logger.Error("Failed to load configuration: %v", err)
		return nil, err
	}

	cm.config = config
	cm.logger.Debug("Configuration loaded successfully")

	return cm.convertToInterface(config), nil
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() *interfaces.Config {
	if cm.config == nil {
		// Return a default config if none is loaded
		return &interfaces.Config{
			Project:      &interfaces.ProjectConfig{},
			Defaults:     &interfaces.DefaultsConfig{},
			Environments: make(map[string]*interfaces.EnvConfig),
			Commands:     make(map[string]*interfaces.CmdConfig),
		}
	}
	return cm.convertToInterface(cm.config)
}

// ReloadConfig reloads configuration from default locations
func (cm *ConfigManager) ReloadConfig() error {
	config, err := cm.loader.Load()
	if err != nil {
		return err
	}
	cm.config = config
	return nil
}

// ValidateConfig validates the current configuration
func (cm *ConfigManager) ValidateConfig() error {
	if cm.config == nil {
		return fmt.Errorf("no configuration loaded")
	}
	return nil
}

// GetConfigSchema returns the configuration schema
func (cm *ConfigManager) GetConfigSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"project": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name":        map[string]interface{}{"type": "string"},
					"version":     map[string]interface{}{"type": "string"},
					"description": map[string]interface{}{"type": "string"},
				},
			},
		},
	}
}

// GetPluginConfig returns configuration for a specific plugin
func (cm *ConfigManager) GetPluginConfig(pluginName string) (map[string]interface{}, error) {
	// For now, return empty config
	return make(map[string]interface{}), nil
}

// SetPluginConfig sets configuration for a specific plugin
func (cm *ConfigManager) SetPluginConfig(pluginName string, config map[string]interface{}) error {
	// For now, just log the operation
	cm.logger.Debug("Setting plugin config for %s: %v", pluginName, config)
	return nil
}

// GetGlobalFlags returns global flags
func (cm *ConfigManager) GetGlobalFlags() *interfaces.GlobalFlags {
	return &interfaces.GlobalFlags{}
}

// GetVerbosity returns the current verbosity level
func (cm *ConfigManager) GetVerbosity() interfaces.VerbosityLevel {
	if cm.config == nil || cm.config.Defaults == nil {
		return interfaces.VerbosityNormal
	}

	switch cm.config.Defaults.Verbosity {
	case "quiet":
		return interfaces.VerbosityQuiet
	case "concise":
		return interfaces.VerbosityConcise
	case "verbose":
		return interfaces.VerbosityVerbose
	case "debug":
		return interfaces.VerbosityDebug
	default:
		return interfaces.VerbosityNormal
	}
}

// SetVerbosity sets the verbosity level
func (cm *ConfigManager) SetVerbosity(level interfaces.VerbosityLevel) error {
	if cm.config == nil {
		return fmt.Errorf("configuration not loaded")
	}

	if cm.config.Defaults == nil {
		cm.config.Defaults = &DefaultsConfig{}
	}

	switch level {
	case interfaces.VerbosityQuiet:
		cm.config.Defaults.Verbosity = "quiet"
	case interfaces.VerbosityConcise:
		cm.config.Defaults.Verbosity = "concise"
	case interfaces.VerbosityNormal:
		cm.config.Defaults.Verbosity = "normal"
	case interfaces.VerbosityVerbose:
		cm.config.Defaults.Verbosity = "verbose"
	case interfaces.VerbosityDebug:
		cm.config.Defaults.Verbosity = "debug"
	default:
		return fmt.Errorf("invalid verbosity level")
	}

	cm.logger.Debug("Verbosity level set to: %s", cm.config.Defaults.Verbosity)
	return nil
}

// GetGlobalConfigPath returns the path to the global configuration file
func (cm *ConfigManager) GetGlobalConfigPath() string {
	return cm.globalPath
}

// SaveConfig saves the current configuration to the specified path
func (cm *ConfigManager) SaveConfig(path string) error {
	if cm.config == nil {
		return fmt.Errorf("no configuration to save")
	}

	var err error
	if path == cm.globalPath {
		err = cm.loader.SaveGlobalConfig(cm.config)
	} else {
		err = cm.saveConfigToPath(path)
	}

	if err != nil {
		cm.logger.Error("Failed to save configuration to %s: %v", path, err)
		return err
	}

	cm.logger.Info("Configuration saved to: %s", path)
	return nil
}

// GetProjectName returns the project name from configuration
func (cm *ConfigManager) GetProjectName() string {
	if cm.config != nil && cm.config.Project != nil {
		return cm.config.Project.Name
	}
	return ""
}

// GetDeploymentTarget returns the current deployment target
func (cm *ConfigManager) GetDeploymentTarget() string {
	if cm.config != nil && cm.config.Defaults != nil {
		return cm.config.Defaults.DeploymentTarget
	}
	return "production"
}

// GetTimeout returns the configured timeout duration
func (cm *ConfigManager) GetTimeout() time.Duration {
	if cm.config != nil && cm.config.Defaults != nil && cm.config.Defaults.Timeout > 0 {
		return cm.config.Defaults.Timeout
	}
	return 5 * time.Minute
}

// GetTheme returns the configured theme
func (cm *ConfigManager) GetTheme() string {
	if cm.config != nil && cm.config.Defaults != nil {
		return cm.config.Defaults.Theme
	}
	return "auto"
}

// GetTemplate returns the configured template
func (cm *ConfigManager) GetTemplate() string {
	if cm.config != nil && cm.config.Defaults != nil {
		return cm.config.Defaults.Template
	}
	return "typescript"
}

// GetEnvironment returns configuration for a specific environment
func (cm *ConfigManager) GetEnvironment(name string) map[string]interface{} {
	if cm.config == nil || cm.config.Environments == nil {
		return nil
	}

	envConfig := cm.config.Environments[name]
	if envConfig == nil {
		return nil
	}

	result := make(map[string]interface{})
	result["api_endpoint"] = envConfig.APIEndpoint
	result["timeout"] = envConfig.Timeout
	if envConfig.Variables != nil {
		result["variables"] = envConfig.Variables
	}

	return result
}

// convertToInterface converts internal Config to interface Config
func (cm *ConfigManager) convertToInterface(config *Config) *interfaces.Config {
	result := &interfaces.Config{
		Project:      &interfaces.ProjectConfig{},
		Defaults:     &interfaces.DefaultsConfig{},
		Environments: make(map[string]*interfaces.EnvConfig),
		Commands:     make(map[string]*interfaces.CmdConfig),
	}

	// Convert project config
	if config.Project != nil {
		result.Project.Name = config.Project.Name
		result.Project.Version = config.Project.Version
		result.Project.Description = config.Project.Description
	}

	// Convert defaults config
	if config.Defaults != nil {
		// Convert string verbosity to enum
		switch config.Defaults.Verbosity {
		case "quiet":
			result.Defaults.Verbosity = interfaces.VerbosityQuiet
		case "concise":
			result.Defaults.Verbosity = interfaces.VerbosityConcise
		case "verbose":
			result.Defaults.Verbosity = interfaces.VerbosityVerbose
		case "debug":
			result.Defaults.Verbosity = interfaces.VerbosityDebug
		default:
			result.Defaults.Verbosity = interfaces.VerbosityNormal
		}
	}

	// Convert environments
	if config.Environments != nil {
		for name, env := range config.Environments {
			result.Environments[name] = &interfaces.EnvConfig{
				Variables: env.Variables,
			}
		}
	}

	// Convert commands
	if config.Commands != nil {
		for name := range config.Commands {
			result.Commands[name] = &interfaces.CmdConfig{}
		}
	}

	return result
}

// saveConfigToPath saves configuration to a specific path
func (cm *ConfigManager) saveConfigToPath(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := yaml.Marshal(cm.config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Load loads configuration with inheritance priority:
// 1. Built-in defaults (lowest priority)
// 2. Global config (~/.engx/config.yaml)
// 3. Project config (.engx/config.yaml)
// 4. Command-line flags (highest priority - handled by caller)
func (l *Loader) Load() (*Config, error) {
	// Start with defaults
	config := NewDefaultConfig()

	// Load and merge global config
	if globalConfig, err := l.loadConfigFile(l.globalPath); err == nil {
		config.Merge(globalConfig)
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load global config: %w", err)
	}

	// Load and merge project config
	if projectConfig, err := l.loadConfigFile(l.projectPath); err == nil {
		config.Merge(projectConfig)
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load project config: %w", err)
	}

	return config, nil
}

// LoadWithCustomPath loads configuration from a custom path
func (l *Loader) LoadWithCustomPath(customPath string) (*Config, error) {
	config := NewDefaultConfig()

	// Load global config first (if exists)
	if globalConfig, err := l.loadConfigFile(l.globalPath); err == nil {
		config.Merge(globalConfig)
	}

	// Load custom config
	if customConfig, err := l.loadConfigFile(customPath); err != nil {
		return nil, fmt.Errorf("failed to load custom config from %s: %w", customPath, err)
	} else {
		config.Merge(customConfig)
	}

	return config, nil
}

// SaveGlobalConfig saves a configuration to the global config file
func (l *Loader) SaveGlobalConfig(config *Config) error {
	dir := filepath.Dir(l.globalPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(l.globalPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write global config: %w", err)
	}

	return nil
}

// loadConfigFile loads a single YAML configuration file
func (l *Loader) loadConfigFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config at %s: %w", path, err)
	}

	return &config, nil
}

// NewDefaultConfig returns a configuration with sensible defaults
func NewDefaultConfig() *Config {
	return &Config{
		Defaults: &DefaultsConfig{
			Verbosity:        "normal",
			DeploymentTarget: "production",
			Timeout:          5 * time.Minute,
			Theme:            "auto",
			Template:         "typescript",
		},
		Environments: map[string]*EnvConfig{
			"development": {
				APIEndpoint: "http://localhost:3000",
				Timeout:     30 * time.Second,
			},
			"staging": {
				APIEndpoint: "https://staging.company.com",
				Timeout:     2 * time.Minute,
			},
			"production": {
				APIEndpoint: "https://api.company.com",
				Timeout:     5 * time.Minute,
			},
		},
		Commands: make(map[string]*CmdConfig),
	}
}

// Merge merges another config into this one, with the other config taking precedence
func (c *Config) Merge(other *Config) {
	if other == nil {
		return
	}

	// Merge project config
	if other.Project != nil {
		c.Project = other.Project
	}

	// Merge defaults
	if other.Defaults != nil {
		if c.Defaults == nil {
			c.Defaults = &DefaultsConfig{}
		}
		c.mergeDefaults(other.Defaults)
	}

	// Merge environments
	if other.Environments != nil {
		if c.Environments == nil {
			c.Environments = make(map[string]*EnvConfig)
		}
		for name, env := range other.Environments {
			c.Environments[name] = env
		}
	}

	// Merge commands
	if other.Commands != nil {
		if c.Commands == nil {
			c.Commands = make(map[string]*CmdConfig)
		}
		for name, cmd := range other.Commands {
			c.Commands[name] = cmd
		}
	}
}

func (c *Config) mergeDefaults(other *DefaultsConfig) {
	if other.Verbosity != "" {
		c.Defaults.Verbosity = other.Verbosity
	}
	if other.DeploymentTarget != "" {
		c.Defaults.DeploymentTarget = other.DeploymentTarget
	}
	if other.Timeout != 0 {
		c.Defaults.Timeout = other.Timeout
	}
	if other.Theme != "" {
		c.Defaults.Theme = other.Theme
	}
	if other.Template != "" {
		c.Defaults.Template = other.Template
	}
}