package subcommands

import (
	"fmt"

	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/data"
	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/renderers"
)

// HandleSearch processes search queries for ambiguous parameters
func HandleSearch(dataStore *data.SimulationDataStore, query string) error {
	// Search crews by name, description, or ID
	crews, err := dataStore.SearchCrews(query)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	// If no results, show helpful message
	if len(crews) == 0 {
		fmt.Printf("🔍 No crews found matching: %s\n\n", query)
		fmt.Printf("Try:\n")
		fmt.Printf("• Crew ID format: CREW-1234\n")
		fmt.Printf("• User email: user@company.com\n")
		fmt.Printf("• Asset URN: asset://service/name\n")
		fmt.Printf("• Search terms: team name or description keywords\n\n")

		// Show available crews as suggestions
		allCrews, _ := dataStore.GetAllCrews()
		if len(allCrews) > 0 {
			fmt.Printf("Available crews:\n")
			for _, crew := range allCrews {
				fmt.Printf("• %s - %s\n", crew.ID, crew.VanityName)
			}
		}
		return nil
	}

	// Create renderer and display results
	renderer := renderers.NewSearchResultsRenderer()
	output, err := renderer.Render(query, crews, getTerminalWidth())
	if err != nil {
		return fmt.Errorf("failed to render search results: %w", err)
	}

	fmt.Print(output)
	return nil
}