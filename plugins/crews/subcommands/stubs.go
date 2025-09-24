package subcommands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/data"
)

// NewCreateCommand creates the create subcommand
func NewCreateCommand(dataStore *data.SimulationDataStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new crew with guided workflow",
		Long: `Create a new crew with guided workflow.

This command will walk you through creating a new crew including:
• Basic crew information (name, description)
• Initial membership and roles
• On-call rotation setup (optional)
• Asset ownership assignment (optional)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return HandleCreate(dataStore)
		},
		Example: `  engx crews create`,
	}

	return cmd
}

// NewAssetsCommand creates the assets subcommand
func NewAssetsCommand(dataStore *data.SimulationDataStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assets <CrewID>",
		Short: "List assets owned by a crew",
		Long: `List all assets owned by a crew including:
• Asset names and types
• Asset URNs and metadata
• Asset status and health
• Delegated access permissions`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return HandleAssets(dataStore, args[0])
		},
		Example: `  engx crews assets CREW-1234`,
	}

	return cmd
}

// NewOnCallCommand creates the oncall subcommand
func NewOnCallCommand(dataStore *data.SimulationDataStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oncall <CrewID>",
		Short: "Show on-call information for a crew",
		Long: `Show on-call rotation and current status for a crew including:
• Current on-call members
• Upcoming rotation schedule
• Escalation path and backup contacts
• On-call rotation type and frequency`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return HandleOnCall(dataStore, args[0])
		},
		Example: `  engx crews oncall CREW-1234`,
	}

	return cmd
}

// NewTransferCommand creates the transfer subcommand
func NewTransferCommand(dataStore *data.SimulationDataStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer <CrewID> --to <UserID>",
		Short: "Transfer crew ownership to another user",
		Long: `Transfer ownership of a crew to another user.

This operation requires current owner privileges and will:
• Transfer full ownership rights
• Update all associated permissions
• Notify relevant stakeholders
• Maintain audit trail of the transfer`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			toUser, _ := cmd.Flags().GetString("to")
			if toUser == "" {
				return fmt.Errorf("--to flag is required")
			}
			return HandleTransfer(dataStore, args[0], toUser)
		},
		Example: `  engx crews transfer CREW-1234 --to sarah@company.com`,
	}

	cmd.Flags().String("to", "", "New owner (UserID or email)")
	cmd.MarkFlagRequired("to")

	return cmd
}

// NewManageCommand creates the manage subcommand
func NewManageCommand(dataStore *data.SimulationDataStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manage <CrewID>",
		Short: "Edit crew details (admin only)",
		Long: `Edit crew details including:
• Crew name and description
• Member roles and permissions
• On-call rotation configuration
• Asset assignments

This command requires admin or owner privileges.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return HandleManage(dataStore, args[0])
		},
		Example: `  engx crews manage CREW-1234`,
	}

	return cmd
}

// Placeholder handlers - these would be fully implemented in a real system

func HandleCreate(dataStore *data.SimulationDataStore) error {
	fmt.Printf("🚧 Crew Creation Workflow\n\n")
	fmt.Printf("This feature will provide a guided workflow for creating new crews,\n")
	fmt.Printf("similar to the 'engx create' experience.\n\n")
	fmt.Printf("The workflow would include:\n")
	fmt.Printf("• Basic crew information collection\n")
	fmt.Printf("• Initial membership setup\n")
	fmt.Printf("• On-call rotation configuration\n")
	fmt.Printf("• Asset ownership assignment\n\n")
	fmt.Printf("📋 Coming Soon - Interactive crew creation\n")
	return nil
}

func HandleAssets(dataStore *data.SimulationDataStore, crewID string) error {
	// Get crew to verify it exists
	crew, err := dataStore.GetCrew(crewID)
	if err != nil {
		return fmt.Errorf("failed to get crew: %w", err)
	}

	// Get assets
	assets, err := dataStore.GetAssetsByCrewID(crewID)
	if err != nil {
		return fmt.Errorf("failed to get assets: %w", err)
	}

	fmt.Printf("📦 Assets owned by %s (%s)\n\n", crew.VanityName, crew.ID)

	if len(assets) == 0 {
		fmt.Printf("No assets currently owned by this crew.\n")
		return nil
	}

	for _, asset := range assets {
		fmt.Printf("• %s (%s)\n", asset.AssetName, asset.AssetType)
		fmt.Printf("  URN: %s\n", asset.AssetURN)
		fmt.Printf("  Status: %s\n", asset.Status)
		fmt.Printf("\n")
	}

	return nil
}

func HandleOnCall(dataStore *data.SimulationDataStore, crewID string) error {
	crew, err := dataStore.GetCrew(crewID)
	if err != nil {
		return fmt.Errorf("failed to get crew: %w", err)
	}

	fmt.Printf("🚨 On-Call Status for %s (%s)\n\n", crew.VanityName, crew.ID)

	if !crew.OnCallSchedule.Enabled {
		fmt.Printf("On-call rotation is not enabled for this crew.\n")
		return nil
	}

	fmt.Printf("Current On-Call: ")
	if len(crew.OnCallSchedule.CurrentOnCall) == 0 {
		fmt.Printf("None\n")
	} else {
		for i, userID := range crew.OnCallSchedule.CurrentOnCall {
			if i > 0 {
				fmt.Printf(", ")
			}
			// Find member details
			for _, member := range crew.Members {
				if member.UserID == userID {
					fmt.Printf("%s (%s)", member.FullName, userID)
					break
				}
			}
		}
		fmt.Printf("\n")
	}

	fmt.Printf("Rotation Type: %s\n", crew.OnCallSchedule.RotationType)

	if len(crew.OnCallSchedule.EscalationPath) > 0 {
		fmt.Printf("Escalation Path: ")
		for i, userID := range crew.OnCallSchedule.EscalationPath {
			if i > 0 {
				fmt.Printf(" → ")
			}
			fmt.Printf("%s", userID)
		}
		fmt.Printf("\n")
	}

	return nil
}

func HandleTransfer(dataStore *data.SimulationDataStore, crewID, toUser string) error {
	crew, err := dataStore.GetCrew(crewID)
	if err != nil {
		return fmt.Errorf("failed to get crew: %w", err)
	}

	fmt.Printf("🔄 Crew Ownership Transfer\n\n")
	fmt.Printf("Crew: %s (%s)\n", crew.VanityName, crew.ID)
	fmt.Printf("Current Owner: %s\n", crew.GetOwner().FullName)
	fmt.Printf("New Owner: %s\n", toUser)
	fmt.Printf("\nThis operation would:\n")
	fmt.Printf("• Transfer full ownership rights\n")
	fmt.Printf("• Update all associated permissions\n")
	fmt.Printf("• Send notifications to stakeholders\n")
	fmt.Printf("• Create audit trail entry\n\n")
	fmt.Printf("✅ Transfer request would be processed (simulation mode)\n")
	return nil
}

func HandleManage(dataStore *data.SimulationDataStore, crewID string) error {
	crew, err := dataStore.GetCrew(crewID)
	if err != nil {
		return fmt.Errorf("failed to get crew: %w", err)
	}

	fmt.Printf("⚙️  Crew Management for %s (%s)\n\n", crew.VanityName, crew.ID)
	fmt.Printf("This feature would provide an interactive interface for:\n\n")
	fmt.Printf("• Editing crew name and description\n")
	fmt.Printf("• Managing member roles and permissions\n")
	fmt.Printf("• Configuring on-call rotations\n")
	fmt.Printf("• Assigning asset ownership\n")
	fmt.Printf("• Setting up access delegations\n\n")
	fmt.Printf("📋 Coming Soon - Interactive crew management\n")
	return nil
}