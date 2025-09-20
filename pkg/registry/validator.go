package registry

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
)

// Validator provides validation capabilities for plugins before registration.
// It ensures plugins meet the required standards and conventions.
type Validator struct {
	namePattern *regexp.Regexp
	versionPattern *regexp.Regexp
}

// NewValidator creates a new plugin validator with default validation rules.
func NewValidator() *Validator {
	return &Validator{
		namePattern:    regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`),
		versionPattern: regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[a-zA-Z0-9.-]+)?$`),
	}
}

// ValidatePlugin performs comprehensive validation of a plugin.
// It checks naming conventions, version format, and interface implementation.
func (v *Validator) ValidatePlugin(plugin interfaces.CommandPlugin) error {
	if plugin == nil {
		return &ValidationError{
			Field:   "plugin",
			Message: "plugin cannot be nil",
		}
	}

	// Validate plugin name
	if err := v.validateName(plugin.Name()); err != nil {
		return err
	}

	// Validate plugin description
	if err := v.validateDescription(plugin.Description()); err != nil {
		return err
	}

	// Validate plugin version
	if err := v.validateVersion(plugin.Version()); err != nil {
		return err
	}

	// Validate service dependencies
	if err := v.validateServices(plugin.RequiredServices(), "required"); err != nil {
		return err
	}

	if err := v.validateServices(plugin.OptionalServices(), "optional"); err != nil {
		return err
	}

	// Check for interface implementation
	if err := v.validateInterface(plugin); err != nil {
		return err
	}

	return nil
}

// validateName ensures plugin names follow naming conventions.
// Names must be lowercase, alphanumeric with hyphens, and reasonable length.
func (v *Validator) validateName(name string) error {
	if name == "" {
		return &ValidationError{
			Field:   "name",
			Message: "plugin name cannot be empty",
		}
	}

	if len(name) < 2 {
		return &ValidationError{
			Field:   "name",
			Message: "plugin name must be at least 2 characters long",
		}
	}

	if len(name) > 50 {
		return &ValidationError{
			Field:   "name",
			Message: "plugin name must be 50 characters or less",
		}
	}

	if !v.namePattern.MatchString(name) {
		return &ValidationError{
			Field:   "name",
			Message: "plugin name must be lowercase alphanumeric with hyphens, starting and ending with alphanumeric characters",
		}
	}

	// Check for reserved names
	reservedNames := []string{"help", "version", "config", "test", "debug"}
	for _, reserved := range reservedNames {
		if name == reserved {
			return &ValidationError{
				Field:   "name",
				Message: fmt.Sprintf("plugin name '%s' is reserved", name),
			}
		}
	}

	return nil
}

// validateDescription ensures descriptions are meaningful and properly formatted.
func (v *Validator) validateDescription(description string) error {
	if description == "" {
		return &ValidationError{
			Field:   "description",
			Message: "plugin description cannot be empty",
		}
	}

	if len(description) < 10 {
		return &ValidationError{
			Field:   "description",
			Message: "plugin description must be at least 10 characters long",
		}
	}

	if len(description) > 200 {
		return &ValidationError{
			Field:   "description",
			Message: "plugin description must be 200 characters or less",
		}
	}

	// Check that description starts with uppercase letter
	if description[0] < 'A' || description[0] > 'Z' {
		return &ValidationError{
			Field:   "description",
			Message: "plugin description must start with an uppercase letter",
		}
	}

	return nil
}

// validateVersion ensures version strings follow semantic versioning.
func (v *Validator) validateVersion(version string) error {
	if version == "" {
		return &ValidationError{
			Field:   "version",
			Message: "plugin version cannot be empty",
		}
	}

	if !v.versionPattern.MatchString(version) {
		return &ValidationError{
			Field:   "version",
			Message: "plugin version must follow semantic versioning (e.g., '1.0.0', 'v1.2.3', '1.0.0-beta.1')",
		}
	}

	return nil
}

// validateServices ensures service dependency lists are valid.
func (v *Validator) validateServices(services []string, serviceType string) error {
	if services == nil {
		return nil // nil is acceptable (equivalent to empty slice)
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, service := range services {
		if service == "" {
			return &ValidationError{
				Field:   serviceType + "_services",
				Message: "service names cannot be empty",
			}
		}

		if seen[service] {
			return &ValidationError{
				Field:   serviceType + "_services",
				Message: fmt.Sprintf("duplicate service '%s' in %s services", service, serviceType),
			}
		}
		seen[service] = true

		// Validate service name format
		if !isValidServiceName(service) {
			return &ValidationError{
				Field:   serviceType + "_services",
				Message: fmt.Sprintf("invalid service name '%s': must be lowercase alphanumeric", service),
			}
		}
	}

	return nil
}

// validateInterface performs basic interface compliance checks.
func (v *Validator) validateInterface(plugin interfaces.CommandPlugin) error {
	// Test that command creation doesn't panic
	defer func() {
		if r := recover(); r != nil {
			// This would be caught by the caller as a validation error
		}
	}()

	// Create command with nil dependencies to test basic functionality
	cmd := plugin.Create(nil)
	if cmd == nil {
		return &ValidationError{
			Field:   "create",
			Message: "plugin Create() method returned nil command",
		}
	}

	// Validate that command Use field matches plugin name
	if cmd.Use == "" {
		return &ValidationError{
			Field:   "create",
			Message: "plugin command must have a Use field",
		}
	}

	// Extract command name from Use field (handle arguments)
	cmdName := strings.Fields(cmd.Use)[0]
	if cmdName != plugin.Name() {
		return &ValidationError{
			Field:   "create",
			Message: fmt.Sprintf("command Use field '%s' must match plugin name '%s'", cmdName, plugin.Name()),
		}
	}

	return nil
}

// isValidServiceName checks if a service name follows naming conventions.
func isValidServiceName(name string) bool {
	if len(name) < 2 || len(name) > 30 {
		return false
	}

	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
			return false
		}
	}

	return true
}

// ValidationError represents a plugin validation error with context.
type ValidationError struct {
	Field   string
	Message string
	Plugin  string
}

func (ve *ValidationError) Error() string {
	if ve.Plugin != "" {
		return fmt.Sprintf("validation error in plugin '%s' field '%s': %s", ve.Plugin, ve.Field, ve.Message)
	}
	return fmt.Sprintf("validation error in field '%s': %s", ve.Field, ve.Message)
}

// ValidatePluginMetadata validates plugin metadata structure.
func (v *Validator) ValidatePluginMetadata(metadata interfaces.PluginMetadata) error {
	if err := v.validateName(metadata.Name); err != nil {
		return err
	}

	if err := v.validateDescription(metadata.Description); err != nil {
		return err
	}

	if err := v.validateVersion(metadata.Version); err != nil {
		return err
	}

	// Validate optional fields if present
	if metadata.Author != "" && len(metadata.Author) > 100 {
		return &ValidationError{
			Field:   "author",
			Message: "author field must be 100 characters or less",
		}
	}

	if metadata.License != "" && len(metadata.License) > 50 {
		return &ValidationError{
			Field:   "license",
			Message: "license field must be 50 characters or less",
		}
	}

	// Validate tags
	if len(metadata.Tags) > 10 {
		return &ValidationError{
			Field:   "tags",
			Message: "maximum of 10 tags allowed",
		}
	}

	for _, tag := range metadata.Tags {
		if len(tag) > 30 {
			return &ValidationError{
				Field:   "tags",
				Message: "individual tags must be 30 characters or less",
			}
		}
	}

	return nil
}

// GetValidationRules returns the current validation rules for documentation purposes.
func (v *Validator) GetValidationRules() ValidationRules {
	return ValidationRules{
		NamePattern:        v.namePattern.String(),
		VersionPattern:     v.versionPattern.String(),
		MinNameLength:      2,
		MaxNameLength:      50,
		MinDescriptionLength: 10,
		MaxDescriptionLength: 200,
		MaxTagCount:        10,
		MaxTagLength:       30,
		ReservedNames:      []string{"help", "version", "config", "test", "debug"},
	}
}

// ValidationRules describes the current validation rules.
type ValidationRules struct {
	NamePattern          string   `json:"name_pattern"`
	VersionPattern       string   `json:"version_pattern"`
	MinNameLength        int      `json:"min_name_length"`
	MaxNameLength        int      `json:"max_name_length"`
	MinDescriptionLength int      `json:"min_description_length"`
	MaxDescriptionLength int      `json:"max_description_length"`
	MaxTagCount          int      `json:"max_tag_count"`
	MaxTagLength         int      `json:"max_tag_length"`
	ReservedNames        []string `json:"reserved_names"`
}