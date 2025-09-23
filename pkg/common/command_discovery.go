package common

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/bthompso/engx-ergonomics-poc/internal/config"
)

// CommandEntry represents a command for the interactive interface
type CommandEntry struct {
	Number      int                    `json:"number"`
	Command     string                 `json:"command"`
	UseCase     string                 `json:"use_case"`
	OwnerCrew   string                 `json:"owner_crew"`
	Status      config.CommandStatus   `json:"status"`
	Category    string                 `json:"category"`
	Available   bool                   `json:"available"`
	Selectable  bool                   `json:"selectable"`
}

// CommandDiscovery handles dynamic command discovery and filtering
type CommandDiscovery struct {
	roadmapConfig   *config.RoadmapConfig
	cobraCommands   []*cobra.Command
	pluginRegistry  *PluginRegistry
}

// NewCommandDiscovery creates a new command discovery instance
func NewCommandDiscovery(roadmapConfig *config.RoadmapConfig, rootCmd *cobra.Command, pluginRegistry *PluginRegistry) *CommandDiscovery {
	// Extract subcommands from cobra root command
	var cobraCommands []*cobra.Command
	if rootCmd != nil {
		cobraCommands = rootCmd.Commands()
	}

	return &CommandDiscovery{
		roadmapConfig:  roadmapConfig,
		cobraCommands:  cobraCommands,
		pluginRegistry: pluginRegistry,
	}
}

// DiscoverCommands performs complete command discovery and returns organized command entries
func (cd *CommandDiscovery) DiscoverCommands() ([]CommandEntry, error) {
	if cd.roadmapConfig == nil {
		return nil, fmt.Errorf("roadmap configuration is required")
	}

	var allEntries []CommandEntry

	// Process roadmap commands
	roadmapEntries := cd.processRoadmapCommands()
	allEntries = append(allEntries, roadmapEntries...)

	// Validate and enrich with cobra command data
	cd.enrichWithCobraData(allEntries)

	// Filter out excluded commands
	filteredEntries := cd.filterExcludedCommands(allEntries)

	return filteredEntries, nil
}

// processRoadmapCommands converts roadmap configuration to command entries
func (cd *CommandDiscovery) processRoadmapCommands() []CommandEntry {
	var entries []CommandEntry

	for _, category := range cd.roadmapConfig.CommandCategories {
		categoryName := cd.roadmapConfig.GetFormattedCategoryName(category.Name)

		for _, cmd := range category.Commands {
			entry := CommandEntry{
				Number:     cmd.Number,
				Command:    cmd.Cmd,
				UseCase:    cmd.Description,
				OwnerCrew:  cmd.OwnerCrew,
				Status:     cmd.StatusEnum,
				Category:   categoryName,
				Available:  cmd.Available,
				Selectable: cmd.Available, // Only available commands are selectable
			}

			entries = append(entries, entry)
		}
	}

	return entries
}

// enrichWithCobraData validates commands against cobra registry and marks availability
func (cd *CommandDiscovery) enrichWithCobraData(entries []CommandEntry) {
	cobraCommandMap := cd.buildCobraCommandMap()

	for i := range entries {
		entry := &entries[i]

		// Check if command exists in cobra registry
		if _, exists := cobraCommandMap[entry.Command]; exists {
			// Command exists and is available
			if entry.Status == config.StatusComingSoon {
				// Override roadmap status if command actually exists
				entry.Status = config.StatusAvailable
				entry.Available = true
				entry.Selectable = true
			}
		} else {
			// Command doesn't exist in cobra - mark as coming soon if available
			if entry.Status == config.StatusAvailable {
				entry.Status = config.StatusComingSoon
				entry.Available = false
				entry.Selectable = false
			}
		}
	}
}

// buildCobraCommandMap creates a map of command names to cobra commands
func (cd *CommandDiscovery) buildCobraCommandMap() map[string]*cobra.Command {
	commandMap := make(map[string]*cobra.Command)

	for _, cmd := range cd.cobraCommands {
		commandMap[cmd.Name()] = cmd

		// Also add aliases
		for _, alias := range cmd.Aliases {
			commandMap[alias] = cmd
		}
	}

	return commandMap
}

// filterExcludedCommands removes commands that should not appear in the interface
func (cd *CommandDiscovery) filterExcludedCommands(entries []CommandEntry) []CommandEntry {
	var filtered []CommandEntry

	for _, entry := range entries {
		if !cd.roadmapConfig.IsCommandExcluded(entry.Command) {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

// GetCommandByNumber finds a command entry by its number
func (cd *CommandDiscovery) GetCommandByNumber(entries []CommandEntry, number int) *CommandEntry {
	for _, entry := range entries {
		if entry.Number == number {
			return &entry
		}
	}
	return nil
}

// GetSelectableCommands returns only commands that can be selected
func (cd *CommandDiscovery) GetSelectableCommands(entries []CommandEntry) []CommandEntry {
	var selectable []CommandEntry

	for _, entry := range entries {
		if entry.Selectable {
			selectable = append(selectable, entry)
		}
	}

	return selectable
}

// ValidateCommand checks if a command can be executed
func (cd *CommandDiscovery) ValidateCommand(commandName string) (bool, error) {
	// Check if command exists in cobra registry
	cobraCommandMap := cd.buildCobraCommandMap()
	if _, exists := cobraCommandMap[commandName]; exists {
		return true, nil
	}

	// Check if command is excluded
	if cd.roadmapConfig.IsCommandExcluded(commandName) {
		return false, fmt.Errorf("command '%s' is excluded from the interface", commandName)
	}

	// Check roadmap status
	if roadmapCmd := cd.roadmapConfig.GetCommandByName(commandName); roadmapCmd != nil {
		if !roadmapCmd.Available {
			return false, fmt.Errorf("command '%s' is not yet available (status: %s)", commandName, roadmapCmd.Status)
		}
	}

	return false, fmt.Errorf("command '%s' not found", commandName)
}

// GetCommandCategories organizes commands by their categories for display
func (cd *CommandDiscovery) GetCommandCategories(entries []CommandEntry) map[string][]CommandEntry {
	categories := make(map[string][]CommandEntry)

	for _, entry := range entries {
		categories[entry.Category] = append(categories[entry.Category], entry)
	}

	return categories
}

// FormatCommandForExecution prepares a command string for execution
func (cd *CommandDiscovery) FormatCommandForExecution(entry CommandEntry) []string {
	// Handle multi-word commands (e.g., "my apps" -> ["my", "apps"])
	parts := strings.Fields(entry.Command)
	if len(parts) == 0 {
		return []string{entry.Command}
	}

	return parts
}

// GetAvailableCommandCount returns the number of available commands
func (cd *CommandDiscovery) GetAvailableCommandCount(entries []CommandEntry) int {
	count := 0
	for _, entry := range entries {
		if entry.Available {
			count++
		}
	}
	return count
}

// GetStatusSummary provides a summary of command statuses
func (cd *CommandDiscovery) GetStatusSummary(entries []CommandEntry) map[config.CommandStatus]int {
	summary := make(map[config.CommandStatus]int)

	for _, entry := range entries {
		summary[entry.Status]++
	}

	return summary
}

// SearchCommands performs search filtering on commands
func (cd *CommandDiscovery) SearchCommands(entries []CommandEntry, query string) []CommandEntry {
	if query == "" {
		return entries
	}

	query = strings.ToLower(query)
	var results []CommandEntry

	for _, entry := range entries {
		// Search in command name
		if strings.Contains(strings.ToLower(entry.Command), query) {
			results = append(results, entry)
			continue
		}

		// Search in description
		if strings.Contains(strings.ToLower(entry.UseCase), query) {
			results = append(results, entry)
			continue
		}

		// Search in category
		if strings.Contains(strings.ToLower(entry.Category), query) {
			results = append(results, entry)
			continue
		}
	}

	return results
}

// GetMaxCommandNumber returns the highest command number for validation
func (cd *CommandDiscovery) GetMaxCommandNumber(entries []CommandEntry) int {
	maxNumber := 0
	for _, entry := range entries {
		if entry.Number > maxNumber {
			maxNumber = entry.Number
		}
	}
	return maxNumber
}

// Debug information helpers

// GetDebugInfo returns debug information about the discovery process
func (cd *CommandDiscovery) GetDebugInfo() map[string]interface{} {
	info := make(map[string]interface{})

	if cd.roadmapConfig != nil {
		info["roadmap_version"] = cd.roadmapConfig.Metadata.Version
		info["roadmap_categories"] = len(cd.roadmapConfig.CommandCategories)
		info["excluded_commands"] = len(cd.roadmapConfig.ExcludedCommands)
		info["ldap_user"] = cd.roadmapConfig.UserContext.ResolvedLDAP
	}

	info["cobra_commands"] = len(cd.cobraCommands)

	if cd.pluginRegistry != nil {
		pluginInfo := cd.pluginRegistry.ListPlugins()
		info["plugin_count"] = len(pluginInfo)
	}

	return info
}

// GetCobraCommandNames returns a list of all cobra command names for debugging
func (cd *CommandDiscovery) GetCobraCommandNames() []string {
	var names []string
	for _, cmd := range cd.cobraCommands {
		names = append(names, cmd.Name())
	}
	return names
}