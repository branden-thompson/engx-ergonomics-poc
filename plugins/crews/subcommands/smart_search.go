package subcommands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/data"
	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/renderers"
	"github.com/spf13/cobra"
)

// SmartSearchCommand handles interactive smart search with selection
func SmartSearchCommand(dataStore *data.SimulationDataStore) *cobra.Command {
	return &cobra.Command{
		Use:     "search [searchTerm]",
		Short:   "Smart search across crews, assets, and users",
		Long:    "Performs intelligent search across crews, assets, and users with interactive selection",
		Example: "engx crews search web\nengx crews search branden\nengx crews search shadcn",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			searchTerm := strings.Join(args, " ")
			return runSmartSearch(dataStore, searchTerm, cmd)
		},
	}
}

// InteractiveSmartSearch handles the case where no exact match is found
func InteractiveSmartSearch(dataStore *data.SimulationDataStore, searchTerm string, cmd *cobra.Command) error {
	return runSmartSearch(dataStore, searchTerm, cmd)
}

// runSmartSearch executes the smart search and interactive selection
func runSmartSearch(dataStore *data.SimulationDataStore, searchTerm string, cmd *cobra.Command) error {
	// Measure search execution time
	start := time.Now()

	renderer := renderers.NewSmartSearchRenderer(dataStore)

	// Perform search
	results, err := renderer.Search(searchTerm, "")
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	// Calculate execution time
	executionTime := time.Since(start)
	results.ExecutionTime = fmt.Sprintf("%.5fs", executionTime.Seconds())

	// Get terminal width
	width := getTerminalWidth()

	// Render search results
	output, err := renderer.Render(results, width)
	if err != nil {
		return fmt.Errorf("failed to render search results: %w", err)
	}

	fmt.Print(output)

	totalResults := len(results.CrewMatches) + len(results.AssetMatches) + len(results.UserMatches)

	// If no results, just return
	if totalResults == 0 {
		fmt.Println("\n No matching crews, assets, or users found.")
		return nil
	}

	// Interactive selection loop
	return handleInteractiveSelection(dataStore, renderer, results, cmd)
}

// handleInteractiveSelection manages the interactive selection process
func handleInteractiveSelection(dataStore *data.SimulationDataStore, renderer *renderers.SmartSearchRenderer, results *renderers.SmartSearchResults, cmd *cobra.Command) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		inputBytes, _, err := reader.ReadLine()
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		inputStr := strings.TrimSpace(string(inputBytes))
		if inputStr == "" {
			continue
		}

		// Check if input is a number (selection)
		if number, err := strconv.Atoi(inputStr); err == nil {
			result := renderer.GetResultByNumber(results, number)
			if result == nil {
				fmt.Printf("Invalid selection. Please choose 1-%d or enter a new search term: ",
					len(results.CrewMatches)+len(results.AssetMatches)+len(results.UserMatches))
				continue
			}

			// Route to appropriate view based on result type
			return routeToView(dataStore, result, cmd)
		}

		// Otherwise, treat as new search term
		fmt.Printf("\nSearching for: %s\n\n", inputStr)
		return runSmartSearch(dataStore, inputStr, cmd)
	}
}

// routeToView routes the selected result to the appropriate detailed view
func routeToView(dataStore *data.SimulationDataStore, result *renderers.SearchResult, cmd *cobra.Command) error {
	switch result.Type {
	case "crew":
		// Route to crew view
		fmt.Printf("\n--- Routing to crew view for: %s ---\n\n", result.ID)
		return HandleDetails(dataStore, result.ID)

	case "asset":
		// Route to asset view
		fmt.Printf("\n--- Routing to asset view for: %s ---\n\n", result.ID)
		return HandleOwner(dataStore, result.ID)

	case "user":
		// Route to user view
		fmt.Printf("\n--- Routing to user view for: %s ---\n\n", result.ID)
		return HandleMembership(dataStore, result.ID)

	default:
		return fmt.Errorf("unknown result type: %s", result.Type)
	}
}

