package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	textTemplate "text/template"
	"time"
)

// PluginGenerator provides plugin scaffolding and generation capabilities
type PluginGenerator struct {
	deps     *Dependencies
	registry *PluginRegistry
}

// PluginTemplate contains the configuration for generating a new plugin
type PluginTemplate struct {
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Version         string            `json:"version"`
	Author          string            `json:"author"`
	PackageName     string            `json:"package_name"`
	CommandName     string            `json:"command_name"`
	RequiredServices []string         `json:"required_services"`
	OptionalServices []string         `json:"optional_services"`
	Advanced        bool              `json:"advanced"`
	Flags           []PluginFlag      `json:"flags"`
	Examples        []string          `json:"examples"`
	CustomFields    map[string]string `json:"custom_fields"`
}

// PluginFlag represents a command-line flag for the plugin
type PluginFlag struct {
	Name        string `json:"name"`
	Type        string `json:"type"`        // "string", "bool", "int", "int64"
	Default     string `json:"default"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// GenerationResult contains the results of plugin generation
type GenerationResult struct {
	PluginPath   string   `json:"plugin_path"`
	FilesCreated []string `json:"files_created"`
	Success      bool     `json:"success"`
	Message      string   `json:"message"`
}

// NewPluginGenerator creates a new plugin generator
func NewPluginGenerator(deps *Dependencies, registry *PluginRegistry) *PluginGenerator {
	return &PluginGenerator{
		deps:     deps,
		registry: registry,
	}
}

// GeneratePlugin creates a new plugin from a template
func (pg *PluginGenerator) GeneratePlugin(template PluginTemplate) (*GenerationResult, error) {
	result := &GenerationResult{
		FilesCreated: make([]string, 0),
		Success:      false,
	}

	// Validate template
	if err := pg.validateTemplate(template); err != nil {
		result.Message = fmt.Sprintf("Template validation failed: %v", err)
		return result, err
	}

	// Create plugin directory
	pluginDir := filepath.Join("plugins", template.PackageName)
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		result.Message = fmt.Sprintf("Failed to create plugin directory: %v", err)
		return result, err
	}

	result.PluginPath = pluginDir

	// Generate plugin.go file
	pluginFile := filepath.Join(pluginDir, "plugin.go")
	if err := pg.generatePluginFile(pluginFile, template); err != nil {
		result.Message = fmt.Sprintf("Failed to generate plugin file: %v", err)
		return result, err
	}
	result.FilesCreated = append(result.FilesCreated, pluginFile)

	// Generate test file if advanced
	if template.Advanced {
		testFile := filepath.Join(pluginDir, "plugin_test.go")
		if err := pg.generateTestFile(testFile, template); err != nil {
			result.Message = fmt.Sprintf("Failed to generate test file: %v", err)
			return result, err
		}
		result.FilesCreated = append(result.FilesCreated, testFile)
	}

	// Generate README.md
	readmeFile := filepath.Join(pluginDir, "README.md")
	if err := pg.generateReadmeFile(readmeFile, template); err != nil {
		result.Message = fmt.Sprintf("Failed to generate README file: %v", err)
		return result, err
	}
	result.FilesCreated = append(result.FilesCreated, readmeFile)

	// Generate example configuration if advanced
	if template.Advanced {
		configFile := filepath.Join(pluginDir, "config.example.yaml")
		if err := pg.generateConfigFile(configFile, template); err != nil {
			result.Message = fmt.Sprintf("Failed to generate config file: %v", err)
			return result, err
		}
		result.FilesCreated = append(result.FilesCreated, configFile)
	}

	result.Success = true
	result.Message = fmt.Sprintf("Plugin '%s' generated successfully", template.Name)

	return result, nil
}

// validateTemplate validates the plugin template
func (pg *PluginGenerator) validateTemplate(template PluginTemplate) error {
	if template.Name == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}

	if template.PackageName == "" {
		return fmt.Errorf("package name cannot be empty")
	}

	if template.CommandName == "" {
		return fmt.Errorf("command name cannot be empty")
	}

	if template.Description == "" {
		return fmt.Errorf("description cannot be empty")
	}

	if template.Version == "" {
		template.Version = "1.0.0"
	}

	// Check if plugin already exists
	existingPlugin := pg.registry.GetPlugin(template.CommandName)
	if existingPlugin != nil {
		return fmt.Errorf("plugin with command name '%s' already exists", template.CommandName)
	}

	// Validate flag types
	for _, flag := range template.Flags {
		if !pg.isValidFlagType(flag.Type) {
			return fmt.Errorf("invalid flag type '%s' for flag '%s'", flag.Type, flag.Name)
		}
	}

	return nil
}

// isValidFlagType checks if a flag type is valid
func (pg *PluginGenerator) isValidFlagType(flagType string) bool {
	validTypes := []string{"string", "bool", "int", "int64", "float64"}
	for _, validType := range validTypes {
		if flagType == validType {
			return true
		}
	}
	return false
}

// generatePluginFile generates the main plugin.go file
func (pg *PluginGenerator) generatePluginFile(filename string, pluginTemplate PluginTemplate) error {
	templateContent := `package {{.PackageName}}

import (
	"fmt"
	{{if .Flags}}"os"{{end}}

	"github.com/bthompso/engx-ergonomics-poc/pkg/common"
	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
	"github.com/spf13/cobra"
)

// Plugin implements the CommandPlugin interface for the {{.CommandName}} command
type Plugin struct {
	deps *common.Dependencies
}

// NewPlugin creates a new {{.Name}} plugin
func NewPlugin(deps *common.Dependencies) interfaces.CommandPlugin {
	return &Plugin{deps: deps}
}

// Name returns the plugin name
func (p *Plugin) Name() string {
	return "{{.CommandName}}"
}

// Description returns the plugin description
func (p *Plugin) Description() string {
	return "{{.Description}}"
}

// Version returns the plugin version
func (p *Plugin) Version() string {
	return "{{.Version}}"
}

// Create returns the cobra command for this plugin
func (p *Plugin) Create(deps interface{}) *cobra.Command {
	// Cast deps back to our Dependencies type
	var dependencies *common.Dependencies
	if d, ok := deps.(*common.Dependencies); ok {
		dependencies = d
		p.deps = d
	} else {
		// Fallback to our stored deps
		dependencies = p.deps
	}

	{{range .Flags}}var {{.Name}} {{.Type}}
	{{end}}

	cmd := &cobra.Command{
		Use:   "{{.CommandName}}{{if .Flags}} [OPTIONS]{{end}}",
		Short: "{{.Description}}",
		Long: ` + "`" + `{{.Description}}

{{if .Examples}}Examples:
{{range .Examples}}  {{.}}
{{end}}{{end}}` + "`" + `,
		RunE: func(cmd *cobra.Command, args []string) error {
			return p.executeCommand(cmd, args, dependencies{{range .Flags}}, {{.Name}}{{end}})
		},
	}

	{{range .Flags}}// Add {{.Name}} flag
	{{if eq .Type "bool"}}cmd.Flags().BoolVar(&{{.Name}}, "{{.Name}}", {{.Default}}, "{{.Description}}")
	{{else if eq .Type "string"}}cmd.Flags().StringVar(&{{.Name}}, "{{.Name}}", "{{.Default}}", "{{.Description}}")
	{{else if eq .Type "int"}}cmd.Flags().IntVar(&{{.Name}}, "{{.Name}}", {{.Default}}, "{{.Description}}")
	{{else if eq .Type "int64"}}cmd.Flags().Int64Var(&{{.Name}}, "{{.Name}}", {{.Default}}, "{{.Description}}")
	{{else if eq .Type "float64"}}cmd.Flags().Float64Var(&{{.Name}}, "{{.Name}}", {{.Default}}, "{{.Description}}")
	{{end}}{{if .Required}}cmd.MarkFlagRequired("{{.Name}}")
	{{end}}
	{{end}}

	return cmd
}

// Initialize initializes the plugin
func (p *Plugin) Initialize() error {
	if p.deps == nil {
		return fmt.Errorf("dependencies not provided")
	}
	return nil
}

// Cleanup performs any necessary cleanup
func (p *Plugin) Cleanup() error {
	return nil
}

// RequiredServices returns required service names
func (p *Plugin) RequiredServices() []string {
	return []string{{"{"}}{{range $i, $service := .RequiredServices}}{{if $i}}, {{end}}"{{$service}}"{{end}}{{"}"}}
}

// OptionalServices returns optional service names
func (p *Plugin) OptionalServices() []string {
	return []string{{"{"}}{{range $i, $service := .OptionalServices}}{{if $i}}, {{end}}"{{$service}}"{{end}}{{"}"}}
}

// executeCommand implements the actual command logic
func (p *Plugin) executeCommand(cmd *cobra.Command, args []string, deps *common.Dependencies{{range .Flags}}, {{.Name}} {{.Type}}{{end}}) error {
	// TODO: Implement your command logic here

	fmt.Printf("🔧 Executing {{.Name}} command...\n")
	{{range .Flags}}fmt.Printf("   {{.Name}}: %v\n", {{.Name}})
	{{end}}

	// Example: Use dependencies
	if deps.Logger != nil {
		deps.Logger.Info("{{.Name}} command executed successfully")
	}

	// TODO: Add your implementation here
	fmt.Println("✅ {{.Name}} command completed successfully!")

	return nil
}

{{if .Advanced}}
// Advanced plugin functionality (implements AdvancedCommandPlugin interface)

// GetMetadata returns detailed plugin metadata
func (p *Plugin) GetMetadata() interfaces.PluginMetadata {
	return interfaces.PluginMetadata{
		Name:        p.Name(),
		Description: p.Description(),
		Version:     p.Version(),
		Author:      "{{.Author}}",
		License:     "MIT",
		Tags:        []string{"{{.PackageName}}", "command"},
		Dependencies: map[string]string{
			{{range .RequiredServices}}"{{.}}": "1.0.0",
			{{end}}
		},
	}
}

// Validate validates plugin configuration and dependencies
func (p *Plugin) Validate() error {
	// TODO: Add custom validation logic
	return nil
}

// HealthCheck performs health check of plugin components
func (p *Plugin) HealthCheck() error {
	// TODO: Add health check logic
	if p.deps == nil {
		return fmt.Errorf("dependencies not initialized")
	}
	return nil
}

// LoadConfig loads plugin-specific configuration
func (p *Plugin) LoadConfig(configData map[string]interface{}) error {
	// TODO: Implement configuration loading
	return nil
}

// GetConfigSchema returns JSON schema for configuration
func (p *Plugin) GetConfigSchema() map[string]interface{} {
	// TODO: Return configuration schema
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"enabled": map[string]interface{}{
				"type": "boolean",
				"default": true,
			},
		},
	}
}

// OnBeforeExecute is called before command execution
func (p *Plugin) OnBeforeExecute(cmd *cobra.Command, args []string) error {
	if p.deps.Logger != nil {
		p.deps.Logger.Debug("{{.Name}} plugin: before execute")
	}
	return nil
}

// OnAfterExecute is called after command execution
func (p *Plugin) OnAfterExecute(cmd *cobra.Command, args []string) error {
	if p.deps.Logger != nil {
		p.deps.Logger.Debug("{{.Name}} plugin: after execute")
	}
	return nil
}

// OnError is called when command encounters error
func (p *Plugin) OnError(cmd *cobra.Command, err error) error {
	if p.deps.Logger != nil {
		p.deps.Logger.Error("{{.Name}} plugin error: %v", err)
	}
	return nil
}
{{end}}
`

	tmpl, err := textTemplate.New("plugin").Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse plugin template: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create plugin file: %w", err)
	}
	defer file.Close()

	return tmpl.Execute(file, pluginTemplate)
}

// generateTestFile generates a test file for the plugin
func (pg *PluginGenerator) generateTestFile(filename string, pluginTemplate PluginTemplate) error {
	testTemplateContent := `package {{.PackageName}}

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
	if plugin.Name() != "{{.CommandName}}" {
		t.Errorf("Expected name '{{.CommandName}}', got '%s'", plugin.Name())
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
	expectedRequired := []string{{"{"}}{{range $i, $service := .RequiredServices}}{{if $i}}, {{end}}"{{$service}}"{{end}}{{"}"}}

	if len(requiredServices) != len(expectedRequired) {
		t.Errorf("Expected %d required services, got %d", len(expectedRequired), len(requiredServices))
	}

	optionalServices := plugin.OptionalServices()
	expectedOptional := []string{{"{"}}{{range $i, $service := .OptionalServices}}{{if $i}}, {{end}}"{{$service}}"{{end}}{{"}"}}

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

{{if .Advanced}}
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
{{end}}
`

	tmpl, err := textTemplate.New("test").Parse(testTemplateContent)
	if err != nil {
		return fmt.Errorf("failed to parse test template: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	defer file.Close()

	return tmpl.Execute(file, pluginTemplate)
}

// generateReadmeFile generates a README.md file for the plugin
func (pg *PluginGenerator) generateReadmeFile(filename string, pluginTemplate PluginTemplate) error {
	readmeTemplateContent := `# {{.Name}} Plugin

{{.Description}}

## Installation

This plugin is part of the ENGX CLI tool plugin system.

## Usage

` + "```bash" + `
engx {{.CommandName}} [OPTIONS]
` + "```" + `

{{if .Flags}}### Options

{{range .Flags}}- ` + "`--{{.Name}}`" + ` ({{.Type}}): {{.Description}}{{if .Default}} (default: {{.Default}}){{end}}{{if .Required}} **[Required]**{{end}}
{{end}}
{{end}}

{{if .Examples}}### Examples

{{range .Examples}}` + "```bash" + `
{{.}}
` + "```" + `

{{end}}{{end}}

## Dependencies

{{if .RequiredServices}}### Required Services

{{range .RequiredServices}}- {{.}}
{{end}}
{{end}}

{{if .OptionalServices}}### Optional Services

{{range .OptionalServices}}- {{.}}
{{end}}
{{end}}

## Development

### Building

` + "```bash" + `
go build -o dist/engx ./cmd/engx
` + "```" + `

### Testing

` + "```bash" + `
go test ./plugins/{{.PackageName}}/...
` + "```" + `

### Validation

` + "```bash" + `
./dist/engx validate plugin {{.CommandName}}
` + "```" + `

## Plugin Information

- **Version**: {{.Version}}
- **Author**: {{.Author}}
- **Package**: {{.PackageName}}
- **Command**: {{.CommandName}}

## License

MIT License - see LICENSE file for details.

---

Generated by ENGX Plugin Generator on {{.CustomFields.GeneratedDate}}
`

	data := pluginTemplate
	data.CustomFields = map[string]string{
		"GeneratedDate": time.Now().Format("2006-01-02 15:04:05"),
	}

	tmpl, err := textTemplate.New("readme").Parse(readmeTemplateContent)
	if err != nil {
		return fmt.Errorf("failed to parse readme template: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create readme file: %w", err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateConfigFile generates an example configuration file
func (pg *PluginGenerator) generateConfigFile(filename string, pluginTemplate PluginTemplate) error {
	configTemplateContent := `# {{.Name}} Plugin Configuration Example
# Copy this file to your project and customize as needed

# Plugin settings
{{.PackageName}}:
  enabled: true
  version: "{{.Version}}"

  # Custom configuration options
  settings:
    # Add your plugin-specific settings here
    debug: false
    timeout: 30s

  {{if .Flags}}# Command line flag defaults
  defaults:
    {{range .Flags}}{{.Name}}: {{.Default}}  # {{.Description}}
    {{end}}
  {{end}}

# Dependencies configuration
dependencies:
  {{range .RequiredServices}}{{.}}: {}
  {{end}}
  {{range .OptionalServices}}{{.}}: {}
  {{end}}
`

	tmpl, err := textTemplate.New("config").Parse(configTemplateContent)
	if err != nil {
		return fmt.Errorf("failed to parse config template: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	return tmpl.Execute(file, pluginTemplate)
}

// GetAvailableTemplates returns predefined plugin templates
func (pg *PluginGenerator) GetAvailableTemplates() []PluginTemplate {
	return []PluginTemplate{
		{
			Name:        "Basic Command",
			Description: "A simple command plugin with basic functionality",
			Version:     "1.0.0",
			Advanced:    false,
			RequiredServices: []string{"logger"},
			OptionalServices: []string{"config"},
			Examples: []string{
				"engx mycommand",
				"engx mycommand --help",
			},
		},
		{
			Name:        "Advanced Command",
			Description: "An advanced command plugin with full functionality",
			Version:     "1.0.0",
			Advanced:    true,
			RequiredServices: []string{"config", "logger"},
			OptionalServices: []string{"filesystem", "tui"},
			Flags: []PluginFlag{
				{
					Name:        "output",
					Type:        "string",
					Default:     "text",
					Description: "Output format (text, json, yaml)",
					Required:    false,
				},
				{
					Name:        "verbose",
					Type:        "bool",
					Default:     "false",
					Description: "Enable verbose output",
					Required:    false,
				},
			},
			Examples: []string{
				"engx mycommand --output json",
				"engx mycommand --verbose",
			},
		},
		{
			Name:        "TUI Command",
			Description: "A command plugin with Terminal User Interface",
			Version:     "1.0.0",
			Advanced:    true,
			RequiredServices: []string{"config", "logger", "tui"},
			OptionalServices: []string{"filesystem"},
			Flags: []PluginFlag{
				{
					Name:        "interactive",
					Type:        "bool",
					Default:     "true",
					Description: "Run in interactive mode",
					Required:    false,
				},
			},
			Examples: []string{
				"engx mycommand",
				"engx mycommand --interactive=false",
			},
		},
	}
}

// ListPluginTemplates returns a formatted list of available templates
func (pg *PluginGenerator) ListPluginTemplates() string {
	templates := pg.GetAvailableTemplates()
	var result strings.Builder

	result.WriteString("📋 Available Plugin Templates:\n\n")

	for i, template := range templates {
		result.WriteString(fmt.Sprintf("%d. %s\n", i+1, template.Name))
		result.WriteString(fmt.Sprintf("   Description: %s\n", template.Description))
		result.WriteString(fmt.Sprintf("   Advanced: %t\n", template.Advanced))
		result.WriteString(fmt.Sprintf("   Required Services: %v\n", template.RequiredServices))
		result.WriteString(fmt.Sprintf("   Optional Services: %v\n", template.OptionalServices))
		if len(template.Flags) > 0 {
			result.WriteString(fmt.Sprintf("   Flags: %d\n", len(template.Flags)))
		}
		result.WriteString("\n")
	}

	return result.String()
}