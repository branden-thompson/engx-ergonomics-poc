package subcommands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/data"
	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/renderers"
)

// NewMembershipCommand creates the membership subcommand
func NewMembershipCommand(dataStore *data.SimulationDataStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "membership <UserID|Email>",
		Short: "Show crew memberships for a user",
		Long: `Show all crew memberships for a user including:

• List of crews the user belongs to
• Role in each crew (owner, admin, member, temp)
• Current on-call status for each crew
• Join date and membership duration
• Access level and permissions

This command accepts either LDAP username or email address
and displays a comprehensive view of the user's crew affiliations.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return HandleMembership(dataStore, args[0])
		},
		Example: `  engx crews membership bthompso
  engx crews membership sarah@company.com
  engx crews membership john.frontend`,
	}

	// Add subcommand for changing membership
	changeCmd := &cobra.Command{
		Use:   "change <UserID|Email> --to <role>",
		Short: "Change a user's role in crews",
		Long: `Change a user's role in one or more crews.

Available roles:
• admin   - Management privileges
• member  - Standard access
• temp    - Temporary access
• removed - Revoke access

This command requires appropriate permissions to modify crew membership.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			toRole, _ := cmd.Flags().GetString("to")
			if toRole == "" {
				return fmt.Errorf("--to flag is required")
			}
			return HandleMembershipChange(dataStore, args[0], toRole)
		},
		Example: `  engx crews membership change john@company.com --to admin
  engx crews membership change temp.user --to removed`,
	}

	changeCmd.Flags().String("to", "", "New role: admin, member, temp, or removed")
	changeCmd.MarkFlagRequired("to")

	cmd.AddCommand(changeCmd)
	return cmd
}

// HandleMembership processes the membership command
func HandleMembership(dataStore *data.SimulationDataStore, userID string) error {
	// Normalize user ID (handle email addresses)
	normalizedUserID := normalizeUserID(userID)

	// Get crews for this user
	crews, err := dataStore.GetCrewsByUser(normalizedUserID)
	if err != nil {
		return fmt.Errorf("failed to get user crews: %w", err)
	}

	if len(crews) == 0 {
		fmt.Printf("No crew memberships found for user: %s\n", userID)
		return nil
	}

	// Create renderer and display
	renderer := renderers.NewMembershipListRenderer()
	output, err := renderer.Render(userID, crews, getTerminalWidth())
	if err != nil {
		return fmt.Errorf("failed to render membership list: %w", err)
	}

	fmt.Print(output)
	return nil
}

// HandleMembershipChange processes membership role changes
func HandleMembershipChange(dataStore *data.SimulationDataStore, userID, newRole string) error {
	// Normalize user ID
	normalizedUserID := normalizeUserID(userID)

	// Get crews for this user
	crews, err := dataStore.GetCrewsByUser(normalizedUserID)
	if err != nil {
		return fmt.Errorf("failed to get user crews: %w", err)
	}

	if len(crews) == 0 {
		return fmt.Errorf("user %s is not a member of any crews", userID)
	}

	// Validate new role
	validRoles := []string{"admin", "member", "temp", "removed"}
	isValid := false
	for _, role := range validRoles {
		if role == newRole {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("invalid role '%s'. Valid roles: %s", newRole, strings.Join(validRoles, ", "))
	}

	// In a real implementation, this would update the database
	// For simulation, we'll just show what would happen
	fmt.Printf("🔄 Membership Change Request\n")
	fmt.Printf("User: %s\n", userID)
	fmt.Printf("New Role: %s\n", newRole)
	fmt.Printf("Affected Crews: %d\n\n", len(crews))

	for _, crew := range crews {
		member := crew.GetMember(normalizedUserID)
		if member != nil {
			fmt.Printf("• %s (%s): %s → %s\n",
				crew.VanityName, crew.ID, member.Role, newRole)
		}
	}

	fmt.Printf("\n✅ Membership changes would be applied (simulation mode)\n")
	return nil
}

// normalizeUserID handles email addresses by extracting username
func normalizeUserID(userID string) string {
	if strings.Contains(userID, "@") {
		return strings.Split(userID, "@")[0]
	}
	return userID
}