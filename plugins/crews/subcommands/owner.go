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
		Use:   "owner <AssetIdentifier>",
		Short: "Show comprehensive asset ownership and access information",
		Long: `Show comprehensive asset ownership information including:

• Asset information and metadata
• Owning crew details and current on-call members
• All crews with access and their permission levels
• Asset dependencies and their health status
• Complete ownership chain and contact information

Supports multiple asset identifier formats:
• Catalog ID: AC123456
• Asset URN: asset://web-app/dashboard
• Vanity Name: "EngX Web Application" or partial matches like "engx"

This command provides complete asset ownership visibility for
contact, access management, and dependency tracking.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return HandleOwner(dataStore, args[0])
		},
		Example: `  engx crews owner AC123456
  engx crews owner asset://web-app/dashboard
  engx crews owner "EngX Web Application"
  engx crews owner engx`,
	}

	return cmd
}

// HandleOwner processes the owner command with multiple input format support
func HandleOwner(dataStore *data.SimulationDataStore, assetParam string) error {
	// Resolve asset using multiple input formats
	asset, err := dataStore.ResolveCatalogAssetParameter(assetParam)
	if err != nil {
		return fmt.Errorf("failed to find asset '%s': %w", assetParam, err)
	}

	// Create renderer and display comprehensive ownership view
	renderer := renderers.NewAssetOwnershipRenderer(dataStore)
	output, err := renderer.Render(asset, getTerminalWidth())
	if err != nil {
		return fmt.Errorf("failed to render asset ownership information: %w", err)
	}

	fmt.Print(output)
	return nil
}