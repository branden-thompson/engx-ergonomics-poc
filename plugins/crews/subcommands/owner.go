package subcommands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/data"
	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/renderers"
)

// NewOwnerCommand creates the owner subcommand
func NewOwnerCommand(dataStore *data.SimulationDataStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner <AssetURN>",
		Short: "Show the owning crew for an asset",
		Long: `Show the crew that owns a specific asset including:

• Asset information and metadata
• Owning crew details and contact information
• Current on-call members for the owning crew
• Delegated access and permissions
• Asset status and management information

This command helps identify who is responsible for a particular
asset and how to contact them for issues or requests.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return HandleOwner(dataStore, args[0])
		},
		Example: `  engx crews owner asset://web-app/dashboard
  engx crews owner asset://service/api-gateway
  engx crews owner asset://database/user-store`,
	}

	return cmd
}

// HandleOwner processes the owner command
func HandleOwner(dataStore *data.SimulationDataStore, assetURN string) error {
	// Get asset information
	asset, err := dataStore.GetAsset(assetURN)
	if err != nil {
		return fmt.Errorf("failed to get asset information: %w", err)
	}

	// Get owning crew
	crew, err := dataStore.GetCrewByAsset(assetURN)
	if err != nil {
		return fmt.Errorf("failed to get owning crew: %w", err)
	}

	// Create renderer and display
	renderer := renderers.NewAssetOwnerRenderer()
	output, err := renderer.Render(asset, crew, getTerminalWidth())
	if err != nil {
		return fmt.Errorf("failed to render asset owner information: %w", err)
	}

	fmt.Print(output)
	return nil
}