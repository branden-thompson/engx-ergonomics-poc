package subcommands

import (
	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/data"
	"github.com/spf13/cobra"
)

// HandleSearch processes search queries for ambiguous parameters
func HandleSearch(dataStore *data.SimulationDataStore, query string) error {
	// Create a mock cobra command for compatibility with the routing system
	mockCmd := &cobra.Command{}

	// Use the interactive smart search for non-exact matches
	return InteractiveSmartSearch(dataStore, query, mockCmd)
}