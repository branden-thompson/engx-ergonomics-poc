package workflows

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
	"golang.org/x/term"
)

// ArchetypeRegistry interface to avoid circular imports
type ArchetypeRegistry interface {
	GetAvailable() []*ArchetypeDefinition
}

// ArchetypeSelectionStage implements archetype selection workflow stage
type ArchetypeSelectionStage struct {
	registry ArchetypeRegistry
	reader   *bufio.Reader
}

// NewArchetypeSelectionStage creates a new archetype selection stage
func NewArchetypeSelectionStage(registry ArchetypeRegistry) *ArchetypeSelectionStage {
	return &ArchetypeSelectionStage{
		registry: registry,
		reader:   bufio.NewReader(os.Stdin),
	}
}

// GetName returns the stage name
func (as *ArchetypeSelectionStage) GetName() string {
	return "Archetype Selection"
}

// GetDescription returns the stage description
func (as *ArchetypeSelectionStage) GetDescription() string {
	return "Interactive selection of application archetype"
}

// CanSkip determines if this stage should be skipped
func (as *ArchetypeSelectionStage) CanSkip(ctx *WorkflowContext) bool {
	// Skip if archetype already selected (direct mode or previous selection)
	return ctx.SelectedArchetype != nil
}

// Execute runs the archetype selection process
func (as *ArchetypeSelectionStage) Execute(ctx *WorkflowContext) (*StageResult, error) {
	// Get available archetypes from registry
	archetypes := as.registry.GetAvailable()
	if len(archetypes) == 0 {
		return nil, fmt.Errorf("no archetypes available in registry")
	}

	// Display selection UI
	selected, err := as.showSelectionUI(archetypes)
	if err != nil {
		return nil, fmt.Errorf("archetype selection failed: %w", err)
	}

	// Return result with selected archetype
	return &StageResult{
		StageType: StageTypeArchetypeSelection,
		Data: map[string]interface{}{
			"selectedArchetype":   selected,
			"availableArchetypes": archetypes,
		},
	}, nil
}

// showSelectionUI displays the interactive archetype selection interface
func (as *ArchetypeSelectionStage) showSelectionUI(archetypes []*ArchetypeDefinition) (*ArchetypeDefinition, error) {
	// Reorder archetypes to match templates view: prod-web, dev-web, hackday, engx-cmd, cli, service
	orderedArchetypes := as.reorderArchetypes(archetypes)

	// Display styled header using the proper component (width-aware)
	header := components.NewHeader("SELECT YOUR APPLICATION TYPE")
	fmt.Print(header.Render())

	// Create table formatter with archetype-specific columns including selection column
	formatter := as.createArchetypeSelectionTable()

	// Prepare data for width calculation (same as templates view)
	data := make([][]string, len(orderedArchetypes))
	for i, archetype := range orderedArchetypes {
		selectionCol := "  " // 2-char selection column (will be used for > and ✓)
		number := fmt.Sprintf("%02d.", i+1)
		framework := as.getFrameworkDisplay(archetype)
		language := as.getLanguageDisplay(archetype)
		appType := fmt.Sprintf("'%s'", archetype.ID)

		data[i] = []string{
			selectionCol,
			number,
			archetype.Name,
			framework,
			language,
			appType,
		}
	}

	// Calculate flexible widths without badge space reservation
	badgeWidth := 0
	columnWidths := formatter.CalculateFlexibleWidths(data, badgeWidth)

	// Print header with calculated widths - add selection column header
	fmt.Printf("%s\n", formatter.FormatHeaderWithWidths(columnWidths))

	// Print each archetype as a numbered row using calculated widths
	for i, archetype := range orderedArchetypes {
		selectionCol := "  " // Empty for now, will be used for interactive selection
		number := fmt.Sprintf("%02d.", i+1)

		// Map archetype to display values
		framework := as.getFrameworkDisplay(archetype)
		language := as.getLanguageDisplay(archetype)
		appType := fmt.Sprintf("'%s'", archetype.ID)

		// Format the row using calculated widths (same pattern as templates view)
		values := []string{selectionCol, number, archetype.Name, framework, language, appType}
		colors := []string{"90", "90", "97", "90", "90", "94"} // Added color for selection column
		rowText := formatter.FormatRowWithWidths(values, colors, columnWidths)

		// Removed default indicator to prevent line wrapping
		fmt.Println(rowText)
	}

	// Width-aware separator
	width := getTerminalWidthForSelector()
	separator := strings.Repeat("-", width)
	fmt.Printf("\n\033[90m%s\033[0m\n\n", separator)

	// Interactive selection loop
	for {
		fmt.Print("? Select an archetype (1-" + strconv.Itoa(len(orderedArchetypes)) + "): ")

		input, err := as.reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read selection: %w", err)
		}

		choice := strings.TrimSpace(input)
		if choice == "" {
			continue
		}

		// Parse the selection
		selectedIndex, err := strconv.Atoi(choice)
		if err != nil || selectedIndex < 1 || selectedIndex > len(orderedArchetypes) {
			fmt.Printf("\033[91mInvalid selection. Please enter a number between 1 and %d.\033[0m\n", len(orderedArchetypes))
			continue
		}

		// Get the selected archetype
		selected := orderedArchetypes[selectedIndex-1]

		// Display selection confirmation with new styling
		fmt.Printf("\n\033[92m[✓] %s\033[0m\n", selected.Name)

		return selected, nil
	}
}

// reorderArchetypes reorders archetypes to match the templates view
// Order: prod-web, dev-web, hackday, engx-cmd, cli, service (agent not in registry yet)
func (as *ArchetypeSelectionStage) reorderArchetypes(archetypes []*ArchetypeDefinition) []*ArchetypeDefinition {
	// Define the desired order matching templates view
	desiredOrder := []string{"prod-web", "dev-web", "hackday", "engx-cmd", "cli", "service"}

	// Create a map for quick lookup
	archetypeMap := make(map[string]*ArchetypeDefinition)
	for _, archetype := range archetypes {
		archetypeMap[archetype.ID] = archetype
	}

	// Build ordered slice
	ordered := make([]*ArchetypeDefinition, 0, len(archetypes))

	// Add archetypes in desired order
	for _, id := range desiredOrder {
		if archetype, exists := archetypeMap[id]; exists {
			ordered = append(ordered, archetype)
		}
	}

	// Add any remaining archetypes that weren't in the desired order
	for _, archetype := range archetypes {
		found := false
		for _, id := range desiredOrder {
			if archetype.ID == id {
				found = true
				break
			}
		}
		if !found {
			ordered = append(ordered, archetype)
		}
	}

	return ordered
}

// createArchetypeSelectionTable creates a custom table formatter for archetype selection with selection column
func (as *ArchetypeSelectionStage) createArchetypeSelectionTable() *components.TableFormatter {
	columns := []components.TableColumn{
		{Header: "  ", Width: 2, MinWidth: 2, Color: "90", Alignment: "left"},                              // Selection column (  > or  ✓)
		{Header: "##", Width: 3, MinWidth: 3, Color: "135", Alignment: "left"},                            // Number column
		{Header: "NAME", Width: 18, MinWidth: 12, Flexible: true, Color: "135", Alignment: "left"},        // Name column (flexible)
		{Header: "FRAMEWORK", Width: 16, MinWidth: 10, Color: "135", Alignment: "left"},                   // Framework column
		{Header: "LANGUAGE", Width: 12, MinWidth: 8, Color: "135", Alignment: "left"},                     // Language column
		{Header: "--app-type", Width: 12, MinWidth: 10, Color: "135", Alignment: "left"},                  // App type column
	}
	return components.NewTableFormatter(columns)
}

// getFrameworkDisplay returns the framework display string for an archetype
func (as *ArchetypeSelectionStage) getFrameworkDisplay(archetype *ArchetypeDefinition) string {
	switch archetype.Category {
	case CategoryWebApplication:
		return "React Router 7"
	case CategoryCLITool, CategoryPlugin:
		return "Bubble"
	case CategoryBackendService:
		return "Bubble"
	case CategoryPrototype:
		return "React Router 7"
	default:
		return "Variable"
	}
}

// getLanguageDisplay returns the language display string for an archetype
func (as *ArchetypeSelectionStage) getLanguageDisplay(archetype *ArchetypeDefinition) string {
	switch archetype.Category {
	case CategoryWebApplication, CategoryPrototype:
		return "TypeScript"
	case CategoryCLITool, CategoryPlugin, CategoryBackendService:
		return "GOLANG"
	default:
		return "Variable"
	}
}

// PromptForAppName prompts the user for an application name
func (as *ArchetypeSelectionStage) PromptForAppName() (string, error) {
	fmt.Print(" ?  What do you want to call your App: ")

	input, err := as.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read app name: %w", err)
	}

	appName := strings.TrimSpace(input)
	if appName == "" {
		return "", fmt.Errorf("application name cannot be empty")
	}

	// Basic validation
	if strings.ContainsAny(appName, " \t\n\r/\\:*?\"<>|") {
		return "", fmt.Errorf("application name contains invalid characters")
	}

	return appName, nil
}

// getTerminalWidthForSelector gets the current terminal width for responsive separator
func getTerminalWidthForSelector() int {
	// Try to get terminal width from stdout
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && width > 0 {
		return width
	}
	// Fallback to a reasonable default
	return 80
}