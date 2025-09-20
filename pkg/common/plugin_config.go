package common

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// PluginConfigManager manages plugin-specific configurations
type PluginConfigManager struct {
	configPath   string
	configs      map[string]*PluginConfig
	schemas      map[string]*PluginConfigSchema
	deps         *Dependencies
	mutex        sync.RWMutex
	watchers     map[string]func(*PluginConfig) error
	autoSave     bool
	lastModified map[string]time.Time
}

// PluginConfig represents configuration for a single plugin
type PluginConfig struct {
	PluginName   string                 `yaml:"plugin_name" json:"plugin_name"`
	Version      string                 `yaml:"version" json:"version"`
	Enabled      bool                   `yaml:"enabled" json:"enabled"`
	Settings     map[string]interface{} `yaml:"settings" json:"settings"`
	Defaults     map[string]interface{} `yaml:"defaults,omitempty" json:"defaults,omitempty"`
	Environment  map[string]interface{} `yaml:"environment,omitempty" json:"environment,omitempty"`
	Flags        map[string]interface{} `yaml:"flags,omitempty" json:"flags,omitempty"`
	Metadata     PluginConfigMetadata   `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	LastUpdated  time.Time              `yaml:"last_updated" json:"last_updated"`
	UpdatedBy    string                 `yaml:"updated_by,omitempty" json:"updated_by,omitempty"`
}

// PluginConfigMetadata contains metadata about the plugin configuration
type PluginConfigMetadata struct {
	Description   string            `yaml:"description,omitempty" json:"description,omitempty"`
	Documentation string            `yaml:"documentation,omitempty" json:"documentation,omitempty"`
	Tags          []string          `yaml:"tags,omitempty" json:"tags,omitempty"`
	Categories    []string          `yaml:"categories,omitempty" json:"categories,omitempty"`
	Dependencies  []string          `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
	Conflicts     []string          `yaml:"conflicts,omitempty" json:"conflicts,omitempty"`
	MinVersion    string            `yaml:"min_version,omitempty" json:"min_version,omitempty"`
	MaxVersion    string            `yaml:"max_version,omitempty" json:"max_version,omitempty"`
	OS            []string          `yaml:"os,omitempty" json:"os,omitempty"`
	Arch          []string          `yaml:"arch,omitempty" json:"arch,omitempty"`
	Custom        map[string]string `yaml:"custom,omitempty" json:"custom,omitempty"`
}

// PluginConfigSchema defines the structure and validation rules for plugin configuration
type PluginConfigSchema struct {
	PluginName  string                            `yaml:"plugin_name" json:"plugin_name"`
	Version     string                            `yaml:"version" json:"version"`
	Description string                            `yaml:"description,omitempty" json:"description,omitempty"`
	Fields      map[string]*PluginConfigField     `yaml:"fields" json:"fields"`
	Groups      map[string]*PluginConfigGroup     `yaml:"groups,omitempty" json:"groups,omitempty"`
	Presets     map[string]*PluginConfigPreset    `yaml:"presets,omitempty" json:"presets,omitempty"`
	Validation  *PluginConfigValidation           `yaml:"validation,omitempty" json:"validation,omitempty"`
	UI          *PluginConfigUI                   `yaml:"ui,omitempty" json:"ui,omitempty"`
	Examples    map[string]map[string]interface{} `yaml:"examples,omitempty" json:"examples,omitempty"`
}

// PluginConfigField defines a single configuration field
type PluginConfigField struct {
	Name         string                 `yaml:"name" json:"name"`
	Type         string                 `yaml:"type" json:"type"` // string, int, bool, float, array, object, duration, url, file, enum
	Description  string                 `yaml:"description,omitempty" json:"description,omitempty"`
	Default      interface{}            `yaml:"default,omitempty" json:"default,omitempty"`
	Required     bool                   `yaml:"required,omitempty" json:"required,omitempty"`
	Sensitive    bool                   `yaml:"sensitive,omitempty" json:"sensitive,omitempty"`
	Validation   *FieldValidation       `yaml:"validation,omitempty" json:"validation,omitempty"`
	UI           *FieldUI               `yaml:"ui,omitempty" json:"ui,omitempty"`
	Group        string                 `yaml:"group,omitempty" json:"group,omitempty"`
	Dependencies []string               `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
	Conflicts    []string               `yaml:"conflicts,omitempty" json:"conflicts,omitempty"`
	Examples     []interface{}          `yaml:"examples,omitempty" json:"examples,omitempty"`
	Deprecated   bool                   `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
	Since        string                 `yaml:"since,omitempty" json:"since,omitempty"`
	Tags         []string               `yaml:"tags,omitempty" json:"tags,omitempty"`
}

// FieldValidation defines validation rules for a configuration field
type FieldValidation struct {
	Min        interface{} `yaml:"min,omitempty" json:"min,omitempty"`
	Max        interface{} `yaml:"max,omitempty" json:"max,omitempty"`
	MinLength  int         `yaml:"min_length,omitempty" json:"min_length,omitempty"`
	MaxLength  int         `yaml:"max_length,omitempty" json:"max_length,omitempty"`
	Pattern    string      `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	Enum       []string    `yaml:"enum,omitempty" json:"enum,omitempty"`
	Format     string      `yaml:"format,omitempty" json:"format,omitempty"` // email, url, uuid, date, time, etc.
	Custom     string      `yaml:"custom,omitempty" json:"custom,omitempty"` // custom validation function name
}

// FieldUI defines UI-specific properties for configuration fields
type FieldUI struct {
	Label       string            `yaml:"label,omitempty" json:"label,omitempty"`
	Placeholder string            `yaml:"placeholder,omitempty" json:"placeholder,omitempty"`
	Help        string            `yaml:"help,omitempty" json:"help,omitempty"`
	Widget      string            `yaml:"widget,omitempty" json:"widget,omitempty"` // input, textarea, select, checkbox, etc.
	Order       int               `yaml:"order,omitempty" json:"order,omitempty"`
	Hidden      bool              `yaml:"hidden,omitempty" json:"hidden,omitempty"`
	Readonly    bool              `yaml:"readonly,omitempty" json:"readonly,omitempty"`
	Advanced    bool              `yaml:"advanced,omitempty" json:"advanced,omitempty"`
	Columns     int               `yaml:"columns,omitempty" json:"columns,omitempty"`
	Attrs       map[string]string `yaml:"attrs,omitempty" json:"attrs,omitempty"`
}

// PluginConfigGroup defines a logical grouping of configuration fields
type PluginConfigGroup struct {
	Name        string   `yaml:"name" json:"name"`
	Label       string   `yaml:"label,omitempty" json:"label,omitempty"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Order       int      `yaml:"order,omitempty" json:"order,omitempty"`
	Collapsible bool     `yaml:"collapsible,omitempty" json:"collapsible,omitempty"`
	Collapsed   bool     `yaml:"collapsed,omitempty" json:"collapsed,omitempty"`
	Advanced    bool     `yaml:"advanced,omitempty" json:"advanced,omitempty"`
	Fields      []string `yaml:"fields,omitempty" json:"fields,omitempty"`
}

// PluginConfigPreset defines a preset configuration for quick setup
type PluginConfigPreset struct {
	Name        string                 `yaml:"name" json:"name"`
	Label       string                 `yaml:"label,omitempty" json:"label,omitempty"`
	Description string                 `yaml:"description,omitempty" json:"description,omitempty"`
	Settings    map[string]interface{} `yaml:"settings" json:"settings"`
	Tags        []string               `yaml:"tags,omitempty" json:"tags,omitempty"`
}

// PluginConfigValidation defines global validation rules
type PluginConfigValidation struct {
	RequiredGroups []string `yaml:"required_groups,omitempty" json:"required_groups,omitempty"`
	MutuallyExclusive [][]string `yaml:"mutually_exclusive,omitempty" json:"mutually_exclusive,omitempty"`
	ConditionalRequired map[string][]string `yaml:"conditional_required,omitempty" json:"conditional_required,omitempty"`
	Custom []string `yaml:"custom,omitempty" json:"custom,omitempty"`
}

// PluginConfigUI defines global UI configuration
type PluginConfigUI struct {
	Theme       string            `yaml:"theme,omitempty" json:"theme,omitempty"`
	Layout      string            `yaml:"layout,omitempty" json:"layout,omitempty"` // tabs, accordion, single-page
	Sections    []string          `yaml:"sections,omitempty" json:"sections,omitempty"`
	ShowAdvanced bool             `yaml:"show_advanced,omitempty" json:"show_advanced,omitempty"`
	Wizard      bool              `yaml:"wizard,omitempty" json:"wizard,omitempty"`
	Attrs       map[string]string `yaml:"attrs,omitempty" json:"attrs,omitempty"`
}

// PluginConfigChangeEvent represents a configuration change event
type PluginConfigChangeEvent struct {
	PluginName  string                 `json:"plugin_name"`
	Action      string                 `json:"action"` // create, update, delete, reload
	OldConfig   *PluginConfig          `json:"old_config,omitempty"`
	NewConfig   *PluginConfig          `json:"new_config,omitempty"`
	Changes     map[string]interface{} `json:"changes,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	User        string                 `json:"user,omitempty"`
	Source      string                 `json:"source,omitempty"` // ui, cli, api, file
}

// NewPluginConfigManager creates a new plugin configuration manager
func NewPluginConfigManager(deps *Dependencies) *PluginConfigManager {
	configDir := filepath.Join(os.Getenv("HOME"), ".engx", "plugins")
	if workingDir, err := os.Getwd(); err == nil {
		if _, err := os.Stat(filepath.Join(workingDir, ".engx")); err == nil {
			configDir = filepath.Join(workingDir, ".engx", "plugins")
		}
	}

	return &PluginConfigManager{
		configPath:   configDir,
		configs:      make(map[string]*PluginConfig),
		schemas:      make(map[string]*PluginConfigSchema),
		deps:         deps,
		watchers:     make(map[string]func(*PluginConfig) error),
		autoSave:     true,
		lastModified: make(map[string]time.Time),
	}
}

// LoadAll loads all plugin configurations from the configuration directory
func (pcm *PluginConfigManager) LoadAll() error {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	if err := os.MkdirAll(pcm.configPath, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.WalkDir(pcm.configPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || (!strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml")) {
			return nil
		}

		// Skip schema files
		if strings.HasSuffix(path, ".schema.yaml") || strings.HasSuffix(path, ".schema.yml") {
			return pcm.loadSchema(path)
		}

		return pcm.loadConfigFile(path)
	})
}

// LoadConfig loads configuration for a specific plugin
func (pcm *PluginConfigManager) LoadConfig(pluginName string) (*PluginConfig, error) {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	if config, exists := pcm.configs[pluginName]; exists {
		return config, nil
	}

	configFile := filepath.Join(pcm.configPath, pluginName+".yaml")
	if err := pcm.loadConfigFile(configFile); err != nil {
		if os.IsNotExist(err) {
			// Create default config
			config := pcm.createDefaultConfig(pluginName)
			pcm.configs[pluginName] = config
			return config, nil
		}
		return nil, err
	}

	return pcm.configs[pluginName], nil
}

// SaveConfig saves configuration for a specific plugin
func (pcm *PluginConfigManager) SaveConfig(pluginName string, config *PluginConfig) error {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	config.LastUpdated = time.Now()
	pcm.configs[pluginName] = config

	if pcm.autoSave {
		return pcm.saveConfigFile(pluginName, config)
	}

	return nil
}

// GetConfig returns configuration for a specific plugin
func (pcm *PluginConfigManager) GetConfig(pluginName string) (*PluginConfig, error) {
	pcm.mutex.RLock()
	defer pcm.mutex.RUnlock()

	if config, exists := pcm.configs[pluginName]; exists {
		return pcm.copyConfig(config), nil
	}

	return nil, fmt.Errorf("configuration for plugin %s not found", pluginName)
}

// SetConfigValue sets a specific configuration value for a plugin
func (pcm *PluginConfigManager) SetConfigValue(pluginName, key string, value interface{}) error {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	config, exists := pcm.configs[pluginName]
	if !exists {
		config = pcm.createDefaultConfig(pluginName)
		pcm.configs[pluginName] = config
	}

	if config.Settings == nil {
		config.Settings = make(map[string]interface{})
	}

	oldValue := config.Settings[key]
	config.Settings[key] = value
	config.LastUpdated = time.Now()

	// Validate the new value
	if err := pcm.validateConfigValue(pluginName, key, value); err != nil {
		config.Settings[key] = oldValue // Rollback
		return fmt.Errorf("validation failed for %s.%s: %w", pluginName, key, err)
	}

	// Trigger watchers
	if watcher, exists := pcm.watchers[pluginName]; exists {
		if err := watcher(config); err != nil {
			if pcm.deps.Logger != nil {
				pcm.deps.Logger.Warn("Config watcher error for %s: %v", pluginName, err)
			}
		}
	}

	if pcm.autoSave {
		return pcm.saveConfigFile(pluginName, config)
	}

	return nil
}

// GetConfigValue returns a specific configuration value for a plugin
func (pcm *PluginConfigManager) GetConfigValue(pluginName, key string) (interface{}, error) {
	config, err := pcm.GetConfig(pluginName)
	if err != nil {
		return nil, err
	}

	if value, exists := config.Settings[key]; exists {
		return value, nil
	}

	// Check defaults
	if value, exists := config.Defaults[key]; exists {
		return value, nil
	}

	// Check schema defaults
	if schema, exists := pcm.schemas[pluginName]; exists {
		if field, exists := schema.Fields[key]; exists && field.Default != nil {
			return field.Default, nil
		}
	}

	return nil, fmt.Errorf("configuration key %s not found for plugin %s", key, pluginName)
}

// ListConfigs returns all loaded plugin configurations
func (pcm *PluginConfigManager) ListConfigs() map[string]*PluginConfig {
	pcm.mutex.RLock()
	defer pcm.mutex.RUnlock()

	result := make(map[string]*PluginConfig)
	for name, config := range pcm.configs {
		result[name] = pcm.copyConfig(config)
	}

	return result
}

// LoadSchema loads configuration schema for a plugin
func (pcm *PluginConfigManager) LoadSchema(pluginName string) (*PluginConfigSchema, error) {
	pcm.mutex.RLock()
	defer pcm.mutex.RUnlock()

	if schema, exists := pcm.schemas[pluginName]; exists {
		return schema, nil
	}

	return nil, fmt.Errorf("schema for plugin %s not found", pluginName)
}

// RegisterSchema registers a configuration schema for a plugin
func (pcm *PluginConfigManager) RegisterSchema(schema *PluginConfigSchema) error {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	pcm.schemas[schema.PluginName] = schema

	if pcm.deps.Logger != nil {
		pcm.deps.Logger.Debug("Registered config schema for plugin: %s", schema.PluginName)
	}

	return nil
}

// ValidateConfig validates a plugin configuration against its schema
func (pcm *PluginConfigManager) ValidateConfig(pluginName string, config *PluginConfig) error {
	schema, exists := pcm.schemas[pluginName]
	if !exists {
		// No schema available, basic validation only
		return pcm.validateBasicConfig(config)
	}

	return pcm.validateAgainstSchema(config, schema)
}

// AddWatcher adds a configuration change watcher for a plugin
func (pcm *PluginConfigManager) AddWatcher(pluginName string, watcher func(*PluginConfig) error) {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	pcm.watchers[pluginName] = watcher
}

// RemoveWatcher removes a configuration change watcher for a plugin
func (pcm *PluginConfigManager) RemoveWatcher(pluginName string) {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	delete(pcm.watchers, pluginName)
}

// SetAutoSave enables or disables automatic saving of configuration changes
func (pcm *PluginConfigManager) SetAutoSave(enabled bool) {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	pcm.autoSave = enabled
}

// SaveAll saves all modified configurations to disk
func (pcm *PluginConfigManager) SaveAll() error {
	pcm.mutex.RLock()
	configs := make(map[string]*PluginConfig)
	for name, config := range pcm.configs {
		configs[name] = config
	}
	pcm.mutex.RUnlock()

	for name, config := range configs {
		if err := pcm.saveConfigFile(name, config); err != nil {
			return fmt.Errorf("failed to save config for %s: %w", name, err)
		}
	}

	return nil
}

// ResetConfig resets a plugin configuration to its defaults
func (pcm *PluginConfigManager) ResetConfig(pluginName string) error {
	pcm.mutex.Lock()
	defer pcm.mutex.Unlock()

	config := pcm.createDefaultConfig(pluginName)
	pcm.configs[pluginName] = config

	if pcm.autoSave {
		return pcm.saveConfigFile(pluginName, config)
	}

	return nil
}

// ExportConfig exports a plugin configuration to a file
func (pcm *PluginConfigManager) ExportConfig(pluginName, format, outputPath string) error {
	config, err := pcm.GetConfig(pluginName)
	if err != nil {
		return err
	}

	var data []byte
	switch format {
	case "yaml", "yml":
		data, err = yaml.Marshal(config)
	case "json":
		data, err = json.MarshalIndent(config, "", "  ")
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(outputPath, data, 0644)
}

// ImportConfig imports a plugin configuration from a file
func (pcm *PluginConfigManager) ImportConfig(pluginName, inputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config PluginConfig

	// Try YAML first, then JSON
	if err := yaml.Unmarshal(data, &config); err != nil {
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse config file (tried YAML and JSON): %w", err)
		}
	}

	config.PluginName = pluginName
	return pcm.SaveConfig(pluginName, &config)
}

// GetConfigPath returns the configuration file path for a plugin
func (pcm *PluginConfigManager) GetConfigPath(pluginName string) string {
	return filepath.Join(pcm.configPath, pluginName+".yaml")
}

// GetStats returns statistics about loaded configurations
func (pcm *PluginConfigManager) GetStats() map[string]interface{} {
	pcm.mutex.RLock()
	defer pcm.mutex.RUnlock()

	totalConfigs := len(pcm.configs)
	enabledConfigs := 0
	totalSettings := 0
	schemasLoaded := len(pcm.schemas)

	for _, config := range pcm.configs {
		if config.Enabled {
			enabledConfigs++
		}
		totalSettings += len(config.Settings)
	}

	return map[string]interface{}{
		"total_configs":    totalConfigs,
		"enabled_configs":  enabledConfigs,
		"disabled_configs": totalConfigs - enabledConfigs,
		"total_settings":   totalSettings,
		"schemas_loaded":   schemasLoaded,
		"auto_save":        pcm.autoSave,
		"config_path":      pcm.configPath,
		"watchers":         len(pcm.watchers),
	}
}

// Private helper methods

func (pcm *PluginConfigManager) loadConfigFile(configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	var config PluginConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file %s: %w", configFile, err)
	}

	if config.PluginName == "" {
		// Derive plugin name from filename
		filename := filepath.Base(configFile)
		config.PluginName = filename[:len(filename)-len(filepath.Ext(filename))]
	}

	pcm.configs[config.PluginName] = &config
	pcm.lastModified[config.PluginName] = time.Now()

	return nil
}

func (pcm *PluginConfigManager) loadSchema(schemaFile string) error {
	data, err := os.ReadFile(schemaFile)
	if err != nil {
		return err
	}

	var schema PluginConfigSchema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		return fmt.Errorf("failed to parse schema file %s: %w", schemaFile, err)
	}

	pcm.schemas[schema.PluginName] = &schema
	return nil
}

func (pcm *PluginConfigManager) saveConfigFile(pluginName string, config *PluginConfig) error {
	configFile := filepath.Join(pcm.configPath, pluginName+".yaml")

	if err := os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	pcm.lastModified[pluginName] = time.Now()
	return nil
}

func (pcm *PluginConfigManager) createDefaultConfig(pluginName string) *PluginConfig {
	config := &PluginConfig{
		PluginName:  pluginName,
		Version:     "1.0.0",
		Enabled:     true,
		Settings:    make(map[string]interface{}),
		Defaults:    make(map[string]interface{}),
		Environment: make(map[string]interface{}),
		Flags:       make(map[string]interface{}),
		LastUpdated: time.Now(),
		UpdatedBy:   "system",
	}

	// Apply schema defaults if available
	if schema, exists := pcm.schemas[pluginName]; exists {
		for fieldName, field := range schema.Fields {
			if field.Default != nil {
				config.Defaults[fieldName] = field.Default
			}
		}
	}

	return config
}

func (pcm *PluginConfigManager) copyConfig(config *PluginConfig) *PluginConfig {
	copied := *config

	// Deep copy maps
	if config.Settings != nil {
		copied.Settings = make(map[string]interface{})
		for k, v := range config.Settings {
			copied.Settings[k] = v
		}
	}

	if config.Defaults != nil {
		copied.Defaults = make(map[string]interface{})
		for k, v := range config.Defaults {
			copied.Defaults[k] = v
		}
	}

	if config.Environment != nil {
		copied.Environment = make(map[string]interface{})
		for k, v := range config.Environment {
			copied.Environment[k] = v
		}
	}

	if config.Flags != nil {
		copied.Flags = make(map[string]interface{})
		for k, v := range config.Flags {
			copied.Flags[k] = v
		}
	}

	return &copied
}

func (pcm *PluginConfigManager) validateConfigValue(pluginName, key string, value interface{}) error {
	schema, exists := pcm.schemas[pluginName]
	if !exists {
		return nil // No schema, no validation
	}

	field, exists := schema.Fields[key]
	if !exists {
		return fmt.Errorf("unknown configuration field: %s", key)
	}

	return pcm.validateFieldValue(field, value)
}

func (pcm *PluginConfigManager) validateFieldValue(field *PluginConfigField, value interface{}) error {
	// Type validation
	if err := pcm.validateFieldType(field, value); err != nil {
		return err
	}

	// Custom validation rules
	if field.Validation != nil {
		return pcm.validateFieldRules(field, value)
	}

	return nil
}

func (pcm *PluginConfigManager) validateFieldType(field *PluginConfigField, value interface{}) error {
	valueType := reflect.TypeOf(value)

	switch field.Type {
	case "string":
		if valueType.Kind() != reflect.String {
			return fmt.Errorf("expected string, got %T", value)
		}
	case "int":
		if valueType.Kind() != reflect.Int && valueType.Kind() != reflect.Int64 && valueType.Kind() != reflect.Float64 {
			return fmt.Errorf("expected integer, got %T", value)
		}
	case "bool":
		if valueType.Kind() != reflect.Bool {
			return fmt.Errorf("expected boolean, got %T", value)
		}
	case "float":
		if valueType.Kind() != reflect.Float64 && valueType.Kind() != reflect.Int && valueType.Kind() != reflect.Int64 {
			return fmt.Errorf("expected float, got %T", value)
		}
	case "array":
		if valueType.Kind() != reflect.Slice && valueType.Kind() != reflect.Array {
			return fmt.Errorf("expected array, got %T", value)
		}
	case "object":
		if valueType.Kind() != reflect.Map {
			return fmt.Errorf("expected object, got %T", value)
		}
	}

	return nil
}

func (pcm *PluginConfigManager) validateFieldRules(field *PluginConfigField, value interface{}) error {
	validation := field.Validation

	// Enum validation
	if len(validation.Enum) > 0 {
		strValue := fmt.Sprintf("%v", value)
		for _, enum := range validation.Enum {
			if enum == strValue {
				return nil
			}
		}
		return fmt.Errorf("value must be one of: %v", validation.Enum)
	}

	// String length validation
	if field.Type == "string" {
		strValue := value.(string)
		if validation.MinLength > 0 && len(strValue) < validation.MinLength {
			return fmt.Errorf("string length must be at least %d", validation.MinLength)
		}
		if validation.MaxLength > 0 && len(strValue) > validation.MaxLength {
			return fmt.Errorf("string length must be at most %d", validation.MaxLength)
		}
	}

	return nil
}

func (pcm *PluginConfigManager) validateBasicConfig(config *PluginConfig) error {
	if config.PluginName == "" {
		return fmt.Errorf("plugin name is required")
	}

	if config.Version == "" {
		return fmt.Errorf("plugin version is required")
	}

	return nil
}

func (pcm *PluginConfigManager) validateAgainstSchema(config *PluginConfig, schema *PluginConfigSchema) error {
	// Basic validation first
	if err := pcm.validateBasicConfig(config); err != nil {
		return err
	}

	// Validate required fields
	for fieldName, field := range schema.Fields {
		if field.Required {
			if _, exists := config.Settings[fieldName]; !exists {
				return fmt.Errorf("required field %s is missing", fieldName)
			}
		}
	}

	// Validate each setting against schema
	for key, value := range config.Settings {
		if err := pcm.validateConfigValue(config.PluginName, key, value); err != nil {
			return fmt.Errorf("validation failed for field %s: %w", key, err)
		}
	}

	return nil
}