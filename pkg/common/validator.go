package common

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
)

// PluginValidator provides comprehensive plugin validation capabilities
type PluginValidator struct {
	registry          *PluginRegistry
	deps              *Dependencies
	knownServices     map[string]ServiceInfo
	validationRules   []ValidationRule
	mutex             sync.RWMutex
}

// ServiceInfo contains information about available services
type ServiceInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Available   bool   `json:"available"`
	Description string `json:"description"`
}

// ValidationRule defines a custom validation rule for plugins
type ValidationRule struct {
	Name        string
	Description string
	Validator   func(plugin interfaces.CommandPlugin, metadata *interfaces.PluginMetadata) error
}

// ValidationResult contains the results of plugin validation
type ValidationResult struct {
	PluginName   string           `json:"plugin_name"`
	IsValid      bool             `json:"is_valid"`
	Errors       []string         `json:"errors,omitempty"`
	Warnings     []string         `json:"warnings,omitempty"`
	Dependencies DependencyStatus `json:"dependencies"`
	Metadata     ValidationInfo   `json:"metadata"`
}

// DependencyStatus tracks the status of plugin dependencies
type DependencyStatus struct {
	Required         []ServiceDependency `json:"required"`
	Optional         []ServiceDependency `json:"optional"`
	Missing          []string            `json:"missing,omitempty"`
	Satisfied        bool                `json:"satisfied"`
	SatisfiedCount   int                 `json:"satisfied_count"`
	RequiredCount    int                 `json:"required_count"`
}

// ServiceDependency represents a dependency on a service
type ServiceDependency struct {
	Name        string `json:"name"`
	Required    bool   `json:"required"`
	Available   bool   `json:"available"`
	Version     string `json:"version,omitempty"`
	MinVersion  string `json:"min_version,omitempty"`
	Compatible  bool   `json:"compatible"`
	Description string `json:"description,omitempty"`
}

// ValidationInfo contains metadata validation results
type ValidationInfo struct {
	NameValid        bool     `json:"name_valid"`
	VersionValid     bool     `json:"version_valid"`
	DescriptionValid bool     `json:"description_valid"`
	Issues           []string `json:"issues,omitempty"`
}

// NewPluginValidator creates a new plugin validator instance
func NewPluginValidator(registry *PluginRegistry, deps *Dependencies) *PluginValidator {
	validator := &PluginValidator{
		registry:      registry,
		deps:          deps,
		knownServices: make(map[string]ServiceInfo),
		mutex:         sync.RWMutex{},
	}

	// Initialize known services
	validator.initializeKnownServices()

	// Setup default validation rules
	validator.setupDefaultValidationRules()

	return validator
}

// initializeKnownServices populates the known services map
func (pv *PluginValidator) initializeKnownServices() {
	services := []ServiceInfo{
		{"config", "1.0.0", true, "Configuration management service"},
		{"filesystem", "1.0.0", true, "Filesystem operations service"},
		{"logger", "1.0.0", true, "Logging service"},
		{"aar", "1.0.0", true, "After Action Report service"},
		{"tui", "1.0.0", true, "Terminal User Interface service"},
		{"chaos", "1.0.0", true, "Chaos engineering service"},
		{"simulation", "1.0.0", true, "Progress simulation service"},
		{"prompts", "1.0.0", true, "Interactive prompts service"},
	}

	for _, service := range services {
		pv.knownServices[service.Name] = service
	}
}

// setupDefaultValidationRules sets up the default validation rules
func (pv *PluginValidator) setupDefaultValidationRules() {
	pv.validationRules = []ValidationRule{
		{
			Name:        "name_format",
			Description: "Plugin name must be non-empty and follow naming conventions",
			Validator:   pv.validateNameFormat,
		},
		{
			Name:        "version_format",
			Description: "Plugin version must follow semantic versioning",
			Validator:   pv.validateVersionFormat,
		},
		{
			Name:        "description_length",
			Description: "Plugin description must be meaningful and appropriately sized",
			Validator:   pv.validateDescriptionLength,
		},
		{
			Name:        "interface_compliance",
			Description: "Plugin must implement required interface methods",
			Validator:   pv.validateInterfaceCompliance,
		},
		{
			Name:        "dependency_availability",
			Description: "Required dependencies must be available",
			Validator:   pv.validateDependencyAvailability,
		},
	}
}

// ValidatePlugin performs comprehensive validation of a plugin
func (pv *PluginValidator) ValidatePlugin(plugin interfaces.CommandPlugin) (*ValidationResult, error) {
	if plugin == nil {
		return nil, fmt.Errorf("plugin cannot be nil")
	}

	result := &ValidationResult{
		PluginName: plugin.Name(),
		IsValid:    true,
		Errors:     make([]string, 0),
		Warnings:   make([]string, 0),
	}

	// Get plugin metadata if available
	var metadata *interfaces.PluginMetadata
	if advancedPlugin, ok := plugin.(interfaces.AdvancedCommandPlugin); ok {
		pluginMetadata := advancedPlugin.GetMetadata()
		metadata = &pluginMetadata
	} else {
		// Create basic metadata from plugin interface
		metadata = &interfaces.PluginMetadata{
			Name:        plugin.Name(),
			Description: plugin.Description(),
			Version:     plugin.Version(),
		}
	}

	// Validate metadata
	result.Metadata = pv.validateMetadata(metadata)

	// Validate dependencies
	result.Dependencies = pv.validateDependencies(plugin)

	// Run custom validation rules
	for _, rule := range pv.validationRules {
		if err := rule.Validator(plugin, metadata); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", rule.Name, err))
			result.IsValid = false
		}
	}

	// Check if plugin implements AdvancedCommandPlugin and run additional validation
	if advancedPlugin, ok := plugin.(interfaces.AdvancedCommandPlugin); ok {
		if err := advancedPlugin.Validate(); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("plugin_validation: %v", err))
			result.IsValid = false
		}

		if err := advancedPlugin.HealthCheck(); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("health_check: %v", err))
		}
	}

	// Final validation check
	if !result.Dependencies.Satisfied {
		result.IsValid = false
		result.Errors = append(result.Errors, "required dependencies not satisfied")
	}

	return result, nil
}

// validateMetadata validates plugin metadata
func (pv *PluginValidator) validateMetadata(metadata *interfaces.PluginMetadata) ValidationInfo {
	info := ValidationInfo{
		NameValid:        true,
		VersionValid:     true,
		DescriptionValid: true,
		Issues:          make([]string, 0),
	}

	// Validate name
	if metadata.Name == "" {
		info.NameValid = false
		info.Issues = append(info.Issues, "name cannot be empty")
	} else if !pv.isValidPluginName(metadata.Name) {
		info.NameValid = false
		info.Issues = append(info.Issues, "name contains invalid characters")
	}

	// Validate version
	if metadata.Version == "" {
		info.VersionValid = false
		info.Issues = append(info.Issues, "version cannot be empty")
	} else if !pv.isValidVersion(metadata.Version) {
		info.VersionValid = false
		info.Issues = append(info.Issues, "version format is invalid")
	}

	// Validate description
	if metadata.Description == "" {
		info.DescriptionValid = false
		info.Issues = append(info.Issues, "description cannot be empty")
	} else if len(metadata.Description) < 10 {
		info.DescriptionValid = false
		info.Issues = append(info.Issues, "description too short (minimum 10 characters)")
	} else if len(metadata.Description) > 200 {
		info.DescriptionValid = false
		info.Issues = append(info.Issues, "description too long (maximum 200 characters)")
	}

	return info
}

// validateDependencies validates plugin dependencies
func (pv *PluginValidator) validateDependencies(plugin interfaces.CommandPlugin) DependencyStatus {
	status := DependencyStatus{
		Required:      make([]ServiceDependency, 0),
		Optional:      make([]ServiceDependency, 0),
		Missing:       make([]string, 0),
		Satisfied:     true,
	}

	// Check required services
	requiredServices := plugin.RequiredServices()
	status.RequiredCount = len(requiredServices)

	for _, serviceName := range requiredServices {
		dependency := pv.checkServiceDependency(serviceName, true)
		status.Required = append(status.Required, dependency)

		if !dependency.Available {
			status.Missing = append(status.Missing, serviceName)
			status.Satisfied = false
		} else {
			status.SatisfiedCount++
		}
	}

	// Check optional services
	optionalServices := plugin.OptionalServices()
	for _, serviceName := range optionalServices {
		dependency := pv.checkServiceDependency(serviceName, false)
		status.Optional = append(status.Optional, dependency)
	}

	return status
}

// checkServiceDependency checks if a service dependency is satisfied
func (pv *PluginValidator) checkServiceDependency(serviceName string, required bool) ServiceDependency {
	pv.mutex.RLock()
	service, exists := pv.knownServices[serviceName]
	pv.mutex.RUnlock()

	return ServiceDependency{
		Name:        serviceName,
		Required:    required,
		Available:   exists && service.Available,
		Version:     service.Version,
		Compatible:  exists && service.Available, // For now, assume compatibility
		Description: service.Description,
	}
}

// Validation rule implementations

func (pv *PluginValidator) validateNameFormat(plugin interfaces.CommandPlugin, metadata *interfaces.PluginMetadata) error {
	name := plugin.Name()
	if name == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}
	if !pv.isValidPluginName(name) {
		return fmt.Errorf("plugin name '%s' contains invalid characters", name)
	}
	return nil
}

func (pv *PluginValidator) validateVersionFormat(plugin interfaces.CommandPlugin, metadata *interfaces.PluginMetadata) error {
	version := plugin.Version()
	if version == "" {
		return fmt.Errorf("plugin version cannot be empty")
	}
	if !pv.isValidVersion(version) {
		return fmt.Errorf("plugin version '%s' is not valid semantic version", version)
	}
	return nil
}

func (pv *PluginValidator) validateDescriptionLength(plugin interfaces.CommandPlugin, metadata *interfaces.PluginMetadata) error {
	description := plugin.Description()
	if description == "" {
		return fmt.Errorf("plugin description cannot be empty")
	}
	if len(description) < 10 {
		return fmt.Errorf("plugin description too short (minimum 10 characters)")
	}
	if len(description) > 200 {
		return fmt.Errorf("plugin description too long (maximum 200 characters)")
	}
	return nil
}

func (pv *PluginValidator) validateInterfaceCompliance(plugin interfaces.CommandPlugin, metadata *interfaces.PluginMetadata) error {
	// Check if all required methods are properly implemented
	if plugin.Name() == "" {
		return fmt.Errorf("Name() method returns empty string")
	}
	if plugin.Description() == "" {
		return fmt.Errorf("Description() method returns empty string")
	}
	if plugin.Version() == "" {
		return fmt.Errorf("Version() method returns empty string")
	}
	if plugin.Create(pv.deps) == nil {
		return fmt.Errorf("Create() method returns nil command")
	}
	return nil
}

func (pv *PluginValidator) validateDependencyAvailability(plugin interfaces.CommandPlugin, metadata *interfaces.PluginMetadata) error {
	requiredServices := plugin.RequiredServices()
	var missingServices []string

	for _, serviceName := range requiredServices {
		pv.mutex.RLock()
		service, exists := pv.knownServices[serviceName]
		pv.mutex.RUnlock()

		if !exists || !service.Available {
			missingServices = append(missingServices, serviceName)
		}
	}

	if len(missingServices) > 0 {
		return fmt.Errorf("missing required services: %s", strings.Join(missingServices, ", "))
	}

	return nil
}

// Helper methods

func (pv *PluginValidator) isValidPluginName(name string) bool {
	// Plugin names should be lowercase, alphanumeric, with hyphens allowed
	match, _ := regexp.MatchString(`^[a-z][a-z0-9-]*[a-z0-9]$|^[a-z]$`, name)
	return match
}

func (pv *PluginValidator) isValidVersion(version string) bool {
	// Basic semantic versioning pattern (major.minor.patch)
	match, _ := regexp.MatchString(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`, version)
	return match
}

// AddValidationRule adds a custom validation rule
func (pv *PluginValidator) AddValidationRule(rule ValidationRule) {
	pv.mutex.Lock()
	defer pv.mutex.Unlock()
	pv.validationRules = append(pv.validationRules, rule)
}

// RemoveValidationRule removes a validation rule by name
func (pv *PluginValidator) RemoveValidationRule(name string) {
	pv.mutex.Lock()
	defer pv.mutex.Unlock()

	for i, rule := range pv.validationRules {
		if rule.Name == name {
			pv.validationRules = append(pv.validationRules[:i], pv.validationRules[i+1:]...)
			break
		}
	}
}

// ValidateAllPlugins validates all registered plugins
func (pv *PluginValidator) ValidateAllPlugins() (map[string]*ValidationResult, error) {
	results := make(map[string]*ValidationResult)

	plugins := pv.registry.ListPlugins()
	for _, pluginInfo := range plugins {
		// Get the actual plugin instance
		plugin := pv.registry.GetPlugin(pluginInfo.Metadata.Name)
		if plugin == nil {
			results[pluginInfo.Metadata.Name] = &ValidationResult{
				PluginName: pluginInfo.Metadata.Name,
				IsValid:    false,
				Errors:     []string{"plugin instance not found"},
			}
			continue
		}

		result, err := pv.ValidatePlugin(plugin)
		if err != nil {
			results[pluginInfo.Metadata.Name] = &ValidationResult{
				PluginName: pluginInfo.Metadata.Name,
				IsValid:    false,
				Errors:     []string{fmt.Sprintf("validation error: %v", err)},
			}
			continue
		}

		results[pluginInfo.Metadata.Name] = result
	}

	return results, nil
}

// GetServiceInfo returns information about a service
func (pv *PluginValidator) GetServiceInfo(serviceName string) (ServiceInfo, bool) {
	pv.mutex.RLock()
	defer pv.mutex.RUnlock()

	service, exists := pv.knownServices[serviceName]
	return service, exists
}

// RegisterService registers a new service for dependency checking
func (pv *PluginValidator) RegisterService(service ServiceInfo) {
	pv.mutex.Lock()
	defer pv.mutex.Unlock()

	pv.knownServices[service.Name] = service
}

// GetKnownServices returns all known services
func (pv *PluginValidator) GetKnownServices() map[string]ServiceInfo {
	pv.mutex.RLock()
	defer pv.mutex.RUnlock()

	services := make(map[string]ServiceInfo)
	for name, service := range pv.knownServices {
		services[name] = service
	}

	return services
}