package subcommands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/data"
	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/renderers"
)

// NewDetailsCommand creates the details subcommand
func NewDetailsCommand(dataStore *data.SimulationDataStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "details <CrewID>",
		Short: "Show detailed information about a crew",
		Long: `Show detailed information about a crew including:

• Basic crew information and description
• Member list with roles and on-call status
• Owned assets and delegated access
• On-call rotation schedule
• Access grants and permissions

The details view provides comprehensive information about a crew's
structure, responsibilities, and current operational status.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return HandleDetails(dataStore, args[0])
		},
		Example: `  engx crews details CREW-1234
  engx crews details CREW-2345`,
	}

	return cmd
}

// HandleDetails processes the details command
func HandleDetails(dataStore *data.SimulationDataStore, crewID string) error {
	// Get crew information
	crew, err := dataStore.GetCrew(crewID)
	if err != nil {
		return fmt.Errorf("failed to get crew details: %w", err)
	}

	// Get assets owned by this crew
	assets, err := dataStore.GetAssetsByCrewID(crewID)
	if err != nil {
		return fmt.Errorf("failed to get crew assets: %w", err)
	}

	// Create renderer and display
	renderer := renderers.NewCrewDetailsRenderer()
	output, err := renderer.Render(crew, assets, getTerminalWidth())
	if err != nil {
		return fmt.Errorf("failed to render crew details: %w", err)
	}

	fmt.Print(output)
	return nil
}

// getTerminalWidth returns the current terminal width using the same logic as the modular table component
func getTerminalWidth() int {
	// Try to get terminal width from stdout
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && width > 0 {
		return width
	}

	// Try COLUMNS environment variable
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width > 0 {
			return width
		}
	}

	// Fallback to reasonable default
	return 80
}