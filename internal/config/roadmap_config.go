package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// CommandStatus represents the availability status of a command
type CommandStatus int

const (
	StatusAvailable CommandStatus = iota
	StatusComingSoon
	StatusExperimental
	StatusDeprecated
)

func (cs CommandStatus) String() string {
	switch cs {
	case StatusAvailable:
		return "available"
	case StatusComingSoon:
		return "coming_soon"
	case StatusExperimental:
		return "experimental"
	case StatusDeprecated:
		return "deprecated"
	default:
		return "unknown"
	}
}

// ParseCommandStatus converts string to CommandStatus
func ParseCommandStatus(status string) CommandStatus {
	switch strings.ToLower(status) {
	case "available":
		return StatusAvailable
	case "coming_soon":
		return StatusComingSoon
	case "experimental":
		return StatusExperimental
	case "deprecated":
		return StatusDeprecated
	default:
		return StatusComingSoon
	}
}

// RoadmapMetadata contains version and update information
type RoadmapMetadata struct {
	Version    string `yaml:"version"`
	Updated    string `yaml:"updated"`
	CLIVersion string `yaml:"cli_version"`
}

// UserContext contains user-specific configuration
type UserContext struct {
	DefaultCrew   string `yaml:"default_crew"`
	LDAPSource    string `yaml:"ldap_source"`
	LDAPOverride  string `yaml:"ldap_override"`
	ResolvedLDAP  string `yaml:"-"` // Computed at runtime
}

// RoadmapCommand represents a single command entry
type RoadmapCommand struct {
	Cmd         string `yaml:"cmd"`
	Description string `yaml:"description"`
	OwnerCrew   string `yaml:"owner_crew"`
	Status      string `yaml:"status"`
	// Computed fields
	Available    bool          `yaml:"-"`
	StatusEnum   CommandStatus `yaml:"-"`
	Number       int           `yaml:"-"`
	CategoryName string        `yaml:"-"`
}

// CommandCategory represents a group of related commands
type CommandCategory struct {
	Name     string           `yaml:"name"`
	Commands []RoadmapCommand `yaml:"commands"`
}

// StatusIndicator defines how status is displayed
type StatusIndicator struct {
	Symbol string `yaml:"symbol"`
	Style  string `yaml:"style"`
}

// InterfaceConfig defines UI layout parameters
type InterfaceConfig struct {
	HeaderText           string `yaml:"header_text"`
	MinimumWidth         int    `yaml:"minimum_width"`
	DefaultWidth         int    `yaml:"default_width"`
	NumberColumnWidth    int    `yaml:"number_column_width"`
	CommandColumnWidth   int    `yaml:"command_column_width"`
	CrewColumnWidth      int    `yaml:"crew_column_width"`
	ExpandableColumn     string `yaml:"expandable_column"`
}

// RoadmapConfig represents the complete roadmap configuration
type RoadmapConfig struct {
	Metadata         RoadmapMetadata            `yaml:"metadata"`
	UserContext      UserContext                `yaml:"user_context"`
	CommandCategories []CommandCategory         `yaml:"command_categories"`
	ExcludedCommands []string                   `yaml:"excluded_commands"`
	Interface        InterfaceConfig            `yaml:"interface"`
	StatusIndicators map[string]StatusIndicator `yaml:"status_indicators"`
}

// LoadRoadmapConfig loads the roadmap configuration from file
func LoadRoadmapConfig() (*RoadmapConfig, error) {
	configPath := filepath.Join(".engx", "roadmap.yaml")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("roadmap configuration not found at %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read roadmap config: %w", err)
	}

	var config RoadmapConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse roadmap config: %w", err)
	}

	// Post-process configuration
	if err := config.PostProcess(); err != nil {
		return nil, fmt.Errorf("failed to process roadmap config: %w", err)
	}

	return &config, nil
}

// PostProcess performs post-loading processing and validation
func (rc *RoadmapConfig) PostProcess() error {
	// Resolve LDAP username
	rc.UserContext.ResolvedLDAP = rc.ResolveLDAPUsername()

	// Process commands and assign numbers
	commandNumber := 1
	for categoryIdx := range rc.CommandCategories {
		category := &rc.CommandCategories[categoryIdx]

		for cmdIdx := range category.Commands {
			cmd := &category.Commands[cmdIdx]

			// Parse status
			cmd.StatusEnum = ParseCommandStatus(cmd.Status)
			cmd.Available = cmd.StatusEnum == StatusAvailable

			// Assign sequential number
			cmd.Number = commandNumber
			commandNumber++

			// Set category name for reference
			cmd.CategoryName = category.Name
		}
	}

	// Set defaults for interface config if not specified
	if rc.Interface.MinimumWidth == 0 {
		rc.Interface.MinimumWidth = 40
	}
	if rc.Interface.DefaultWidth == 0 {
		rc.Interface.DefaultWidth = 80
	}
	if rc.Interface.NumberColumnWidth == 0 {
		rc.Interface.NumberColumnWidth = 4
	}
	if rc.Interface.CommandColumnWidth == 0 {
		rc.Interface.CommandColumnWidth = 20
	}
	if rc.Interface.CrewColumnWidth == 0 {
		rc.Interface.CrewColumnWidth = 12
	}
	if rc.Interface.ExpandableColumn == "" {
		rc.Interface.ExpandableColumn = "description"
	}

	return nil
}

// ResolveLDAPUsername determines the LDAP username based on configuration
func (rc *RoadmapConfig) ResolveLDAPUsername() string {
	switch rc.UserContext.LDAPSource {
	case "config":
		if rc.UserContext.LDAPOverride != "" {
			return rc.UserContext.LDAPOverride
		}
		return "unknown-user"
	case "system":
		fallthrough
	default:
		// Get from environment
		if user := os.Getenv("USER"); user != "" {
			return user
		}
		if user := os.Getenv("USERNAME"); user != "" {
			return user
		}
		return "unknown-user"
	}
}

// GetAllCommands returns all commands across categories
func (rc *RoadmapConfig) GetAllCommands() []RoadmapCommand {
	var allCommands []RoadmapCommand

	for _, category := range rc.CommandCategories {
		allCommands = append(allCommands, category.Commands...)
	}

	return allCommands
}

// GetAvailableCommands returns only available commands
func (rc *RoadmapConfig) GetAvailableCommands() []RoadmapCommand {
	var availableCommands []RoadmapCommand

	for _, cmd := range rc.GetAllCommands() {
		if cmd.Available {
			availableCommands = append(availableCommands, cmd)
		}
	}

	return availableCommands
}

// GetCommandByNumber finds a command by its assigned number
func (rc *RoadmapConfig) GetCommandByNumber(number int) *RoadmapCommand {
	for _, cmd := range rc.GetAllCommands() {
		if cmd.Number == number {
			return &cmd
		}
	}
	return nil
}

// GetCommandByName finds a command by its name
func (rc *RoadmapConfig) GetCommandByName(name string) *RoadmapCommand {
	for _, cmd := range rc.GetAllCommands() {
		if cmd.Cmd == name {
			return &cmd
		}
	}
	return nil
}

// IsCommandExcluded checks if a command should be excluded from the interface
func (rc *RoadmapConfig) IsCommandExcluded(commandName string) bool {
	for _, excluded := range rc.ExcludedCommands {
		if excluded == commandName {
			return true
		}
	}
	return false
}

// GetHeaderText returns the formatted header text with version substitution
func (rc *RoadmapConfig) GetHeaderText() string {
	header := rc.Interface.HeaderText

	// Replace version placeholder if present
	if strings.Contains(header, "##.##.####") && rc.Metadata.CLIVersion != "" {
		header = strings.ReplaceAll(header, "##.##.####", rc.Metadata.CLIVersion)
	}

	// Replace LDAP username placeholder
	if strings.Contains(header, "@LDAP-username") {
		header = strings.ReplaceAll(header, "@LDAP-username", "@"+rc.UserContext.ResolvedLDAP)
	}

	return header
}

// GetStatusIndicator returns the display configuration for a status
func (rc *RoadmapConfig) GetStatusIndicator(status CommandStatus) StatusIndicator {
	statusStr := status.String()
	if indicator, exists := rc.StatusIndicators[statusStr]; exists {
		return indicator
	}

	// Default indicators if not configured
	switch status {
	case StatusAvailable:
		return StatusIndicator{Symbol: "✓", Style: "success"}
	case StatusComingSoon:
		return StatusIndicator{Symbol: "⏳", Style: "muted"}
	case StatusExperimental:
		return StatusIndicator{Symbol: "🧪", Style: "warning"}
	case StatusDeprecated:
		return StatusIndicator{Symbol: "⚠️", Style: "error"}
	default:
		return StatusIndicator{Symbol: "?", Style: "muted"}
	}
}

// Validate performs comprehensive validation of the configuration
func (rc *RoadmapConfig) Validate() error {
	// Check required fields
	if rc.Metadata.Version == "" {
		return fmt.Errorf("metadata.version is required")
	}

	if len(rc.CommandCategories) == 0 {
		return fmt.Errorf("at least one command category is required")
	}

	// Validate each category
	for i, category := range rc.CommandCategories {
		if category.Name == "" {
			return fmt.Errorf("category %d: name is required", i)
		}

		if len(category.Commands) == 0 {
			return fmt.Errorf("category '%s': at least one command is required", category.Name)
		}

		// Validate commands within category
		for j, cmd := range category.Commands {
			if cmd.Cmd == "" {
				return fmt.Errorf("category '%s', command %d: cmd is required", category.Name, j)
			}

			if cmd.Description == "" {
				return fmt.Errorf("category '%s', command '%s': description is required", category.Name, cmd.Cmd)
			}

			if cmd.OwnerCrew == "" {
				return fmt.Errorf("category '%s', command '%s': owner_crew is required", category.Name, cmd.Cmd)
			}
		}
	}

	// Validate interface configuration
	if rc.Interface.MinimumWidth < 30 {
		return fmt.Errorf("interface.minimum_width must be at least 30")
	}

	if rc.Interface.DefaultWidth < rc.Interface.MinimumWidth {
		return fmt.Errorf("interface.default_width must be >= minimum_width")
	}

	return nil
}

// GetFormattedCategoryName returns category name with LDAP substitution
func (rc *RoadmapConfig) GetFormattedCategoryName(categoryName string) string {
	formatted := categoryName
	if strings.Contains(formatted, "@LDAP-username") {
		formatted = strings.ReplaceAll(formatted, "@LDAP-username", "@"+rc.UserContext.ResolvedLDAP)
	}
	return formatted
}