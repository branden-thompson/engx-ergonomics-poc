package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Loader handles configuration loading with inheritance
type Loader struct {
	globalPath  string
	projectPath string
}

// NewLoader creates a new configuration loader
func NewLoader() *Loader {
	homeDir, _ := os.UserHomeDir()
	globalPath := filepath.Join(homeDir, ".dpx-web", "config.yaml")
	projectPath := ".dpx-web/config.yaml"

	return &Loader{
		globalPath:  globalPath,
		projectPath: projectPath,
	}
}

// Load loads configuration with inheritance priority:
// 1. Built-in defaults (lowest priority)
// 2. Global config (~/.dpx-web/config.yaml)
// 3. Project config (.dpx-web/config.yaml)
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

// GetConfigPaths returns the paths that will be checked for configuration
func (l *Loader) GetConfigPaths() []string {
	return []string{
		l.globalPath,
		l.projectPath,
	}
}

// EnsureGlobalConfigDir creates the global config directory if it doesn't exist
func (l *Loader) EnsureGlobalConfigDir() error {
	dir := filepath.Dir(l.globalPath)
	return os.MkdirAll(dir, 0755)
}

// SaveGlobalConfig saves a configuration to the global config file
func (l *Loader) SaveGlobalConfig(config *Config) error {
	if err := l.EnsureGlobalConfigDir(); err != nil {
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

// SaveProjectConfig saves a configuration to the project config file
func (l *Loader) SaveProjectConfig(config *Config) error {
	dir := filepath.Dir(l.projectPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create project config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(l.projectPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write project config: %w", err)
	}

	return nil
}