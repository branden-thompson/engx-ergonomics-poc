package crews

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/data"
	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/subcommands"
)

// CrewsCommand represents the crews command
type CrewsCommand struct {
	dataStore    *data.SimulationDataStore
	resolver     *ParameterResolver
}

// NewCrewsCommand creates a new crews command
func NewCrewsCommand() *cobra.Command {
	crewsCmd := &CrewsCommand{
		dataStore: data.NewSimulationDataStore(),
		resolver:  NewParameterResolver(),
	}

	cmd := &cobra.Command{
		Use:   "crews [CrewID|UserID|AssetID|AssetURN]",
		Short: "Manage and query crew information, membership, and ownership",
		Long: `Manage and query crew information, membership, and ownership.

The crews command provides comprehensive crew management functionality:

• Smart parameter detection automatically routes to appropriate views
• Crew management operations for administrators
• Membership tracking and role management
• Asset ownership queries and management
• On-call rotation information

Smart Parameter Examples:
  engx crews CREW-1234              # Show crew details
  engx crews bthompso@company.com   # Show user's crew memberships
  engx crews AC123456               # Show asset ownership details
  engx crews asset://web-app/dashboard # Show asset's owning crew

Subcommands:
  create       Create a new crew with guided workflow
  details      Show detailed crew information
  membership   Manage user crew memberships
  assets       List assets owned by a crew
  oncall       Manage on-call rotations
  transfer     Transfer crew ownership
  manage       Edit crew details (admin only)`,
		RunE: crewsCmd.Execute,
		Example: `  # Smart parameter detection
  engx crews CREW-1234
  engx crews sarah@company.com
  engx crews AC123456
  engx crews asset://service/api-gateway

  # Explicit subcommands
  engx crews create
  engx crews details CREW-1234
  engx crews membership bthompso
  engx crews assets CREW-1234
  engx crews oncall CREW-1234`,
	}

	// Add subcommands
	cmd.AddCommand(subcommands.NewCreateCommand(crewsCmd.dataStore))
	cmd.AddCommand(subcommands.NewDetailsCommand(crewsCmd.dataStore))
	cmd.AddCommand(subcommands.NewMembershipCommand(crewsCmd.dataStore))
	cmd.AddCommand(subcommands.NewAssetsCommand(crewsCmd.dataStore))
	cmd.AddCommand(subcommands.NewOnCallCommand(crewsCmd.dataStore))
	cmd.AddCommand(subcommands.NewTransferCommand(crewsCmd.dataStore))
	cmd.AddCommand(subcommands.NewManageCommand(crewsCmd.dataStore))
	cmd.AddCommand(subcommands.NewOwnerCommand(crewsCmd.dataStore))
	cmd.AddCommand(subcommands.SmartSearchCommand(crewsCmd.dataStore))

	return cmd
}

// Execute handles the main crews command with smart parameter detection
func (c *CrewsCommand) Execute(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	// Use smart parameter resolution for single argument
	if len(args) == 1 {
		param := args[0]
		_, handlerType, err := c.resolver.Resolve(param)
		if err != nil {
			return fmt.Errorf("parameter resolution failed: %w", err)
		}

		// Execute the resolved handler with fallback logic
		switch handlerType {
		case "details":
			err := c.tryCrewDetails(param)
			if err != nil {
				return c.fallbackToSearch(param, "crew")
			}
			return nil
		case "membership":
			err := c.tryUserMembership(param)
			if err != nil {
				return c.fallbackToSearch(param, "user")
			}
			return nil
		case "owner":
			err := c.tryAssetOwner(param)
			if err != nil {
				return c.fallbackToSearch(param, "asset")
			}
			return nil
		case "search":
			return subcommands.HandleSearch(c.dataStore, param)
		default:
			return fmt.Errorf("unknown handler type: %s", handlerType)
		}
	}

	// For multiple arguments, show help
	return cmd.Help()
}

// tryCrewDetails attempts to show crew details, returns error if crew not found
func (c *CrewsCommand) tryCrewDetails(param string) error {
	// Check if crew exists
	_, err := c.dataStore.GetCrew(param)
	if err != nil {
		return err // Crew not found
	}

	// Crew exists, show details
	return subcommands.HandleDetails(c.dataStore, param)
}

// tryUserMembership attempts to show user memberships, returns error if user not found
func (c *CrewsCommand) tryUserMembership(param string) error {
	// Normalize and check if user has any crew memberships
	normalizedUserID := c.normalizeUserID(param)
	crews, err := c.dataStore.GetCrewsByUser(normalizedUserID)
	if err != nil {
		return err
	}

	if len(crews) == 0 {
		return fmt.Errorf("user %s not found", param)
	}

	// User found, show memberships
	return subcommands.HandleMembership(c.dataStore, param)
}

// tryAssetOwner attempts to show asset ownership, returns error if asset not found
func (c *CrewsCommand) tryAssetOwner(param string) error {
	// Try to resolve the asset parameter
	_, err := c.dataStore.ResolveCatalogAssetParameter(param)
	if err != nil {
		return err // Asset not found
	}

	// Asset exists, show ownership
	return subcommands.HandleOwner(c.dataStore, param)
}

// fallbackToSearch falls back to SmartSearch with a helpful message
func (c *CrewsCommand) fallbackToSearch(param, expectedType string) error {
	fmt.Printf("No exact %s match found for: %s\n", expectedType, param)
	fmt.Printf("Let me search for similar crews, assets, and users...\n\n")

	return subcommands.HandleSearch(c.dataStore, param)
}

// normalizeUserID normalizes user identifiers (same logic as in membership.go)
func (c *CrewsCommand) normalizeUserID(userID string) string {
	// Handle email addresses by extracting username
	if strings.Contains(userID, "@") {
		return strings.Split(userID, "@")[0]
	}
	return userID
}

// ParameterResolver handles smart parameter detection and routing
type ParameterResolver struct {
	patterns map[string]*regexp.Regexp
}

// NewParameterResolver creates a new parameter resolver
func NewParameterResolver() *ParameterResolver {
	return &ParameterResolver{
		patterns: map[string]*regexp.Regexp{
			"crewID":     regexp.MustCompile(`^CREW-\d{4}$`),
			"email":      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
			"assetURN":   regexp.MustCompile(`^(asset://|urn:asset:)`),
			"catalogID":  regexp.MustCompile(`^AC[0-9A-F]{6}$`), // AC + 6 hex characters
			"ldapUser":   regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._-]*[a-zA-Z0-9]?$`), // Must be <= 10 chars, mostly alphabetical
		},
	}
}

// Resolve determines the appropriate handler for a given parameter using 5-step logic
func (r *ParameterResolver) Resolve(param string) (interface{}, string, error) {
	// Step 1: Check exact crew ID match
	if r.patterns["crewID"].MatchString(param) {
		return nil, "details", nil
	}

	// Step 2: Check asset patterns first (more specific than user patterns)
	if r.patterns["assetURN"].MatchString(param) {
		return nil, "owner", nil
	}
	if r.patterns["catalogID"].MatchString(param) {
		return nil, "owner", nil
	}

	// Step 3: Check user patterns (email or LDAP)
	if r.patterns["email"].MatchString(param) {
		return nil, "membership", nil
	}
	if r.patterns["ldapUser"].MatchString(param) && len(param) >= 5 && len(param) <= 10 {
		return nil, "membership", nil
	}

	// Step 4: If nothing matched exactly, try smart search
	return nil, "search", nil
}

// GetTerminalWidth returns the current terminal width
func GetTerminalWidth() int {
	// Try to get terminal width from environment or tty
	if width := os.Getenv("COLUMNS"); width != "" {
		if w := parseWidth(width); w > 0 {
			return w
		}
	}

	// Fallback to reasonable default
	return 80
}

func parseWidth(width string) int {
	var w int
	if _, err := fmt.Sscanf(width, "%d", &w); err == nil && w > 0 {
		return w
	}
	return 0
}