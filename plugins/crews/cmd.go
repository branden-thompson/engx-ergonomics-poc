package crews

import (
	"fmt"
	"os"
	"regexp"

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
		Use:   "crews [CrewID|UserID|AssetURN]",
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

	return cmd
}

// Execute handles the main crews command with smart parameter detection
func (c *CrewsCommand) Execute(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	// Use smart parameter resolution for single argument
	if len(args) == 1 {
		_, handlerType, err := c.resolver.Resolve(args[0])
		if err != nil {
			return fmt.Errorf("parameter resolution failed: %w", err)
		}

		// Execute the resolved handler
		switch handlerType {
		case "details":
			return subcommands.HandleDetails(c.dataStore, args[0])
		case "membership":
			return subcommands.HandleMembership(c.dataStore, args[0])
		case "owner":
			return subcommands.HandleOwner(c.dataStore, args[0])
		case "search":
			return subcommands.HandleSearch(c.dataStore, args[0])
		default:
			return fmt.Errorf("unknown handler type: %s", handlerType)
		}
	}

	// For multiple arguments, show help
	return cmd.Help()
}

// ParameterResolver handles smart parameter detection and routing
type ParameterResolver struct {
	patterns map[string]*regexp.Regexp
}

// NewParameterResolver creates a new parameter resolver
func NewParameterResolver() *ParameterResolver {
	return &ParameterResolver{
		patterns: map[string]*regexp.Regexp{
			"crewID":    regexp.MustCompile(`^CREW-\d{4}$`),
			"email":     regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
			"assetURN":  regexp.MustCompile(`^(asset://|urn:asset:)`),
			"ldapUser":  regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._-]*$`),
		},
	}
}

// Resolve determines the appropriate handler for a given parameter
func (r *ParameterResolver) Resolve(param string) (interface{}, string, error) {
	switch {
	case r.patterns["crewID"].MatchString(param):
		return nil, "details", nil
	case r.patterns["email"].MatchString(param):
		return nil, "membership", nil
	case r.patterns["assetURN"].MatchString(param):
		return nil, "owner", nil
	case r.patterns["ldapUser"].MatchString(param):
		return nil, "membership", nil
	default:
		return nil, "search", nil
	}
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