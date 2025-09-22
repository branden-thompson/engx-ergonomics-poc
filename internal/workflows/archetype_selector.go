package workflows

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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

	// Try interactive mode first, fallback to simple mode if needed
	selected, err := as.interactiveSelection(orderedArchetypes, formatter, columnWidths)
	if err != nil {
		// Fallback to simple selection mode
		return as.fallbackSelection(orderedArchetypes, formatter, columnWidths)
	}

	return selected, nil
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

// interactiveSelection implements the full interactive archetype selection with arrow keys
func (as *ArchetypeSelectionStage) interactiveSelection(archetypes []*ArchetypeDefinition, formatter *components.TableFormatter, columnWidths []int) (*ArchetypeDefinition, error) {
	// Temporarily disable interactive mode to ensure clean table rendering
	// TODO: Re-implement interactive features once table formatting is stable
	return nil, fmt.Errorf("interactive mode disabled, using fallback")
}

// handleNumberInput processes debounced number input
func (as *ArchetypeSelectionStage) handleNumberInput(currentIndex *int, numberBuffer *string, lastInputTime *time.Time, maxIndex int) {
	time.Sleep(200 * time.Millisecond) // Debounce delay

	// Check if input has stopped
	if time.Since(*lastInputTime) >= 200*time.Millisecond && *numberBuffer != "" {
		if num, err := strconv.Atoi(*numberBuffer); err == nil && num >= 1 && num <= maxIndex {
			*currentIndex = num - 1
			*numberBuffer = ""
			as.printPrompt(*currentIndex+1, maxIndex, "")
		} else {
			// Invalid number, clear buffer
			*numberBuffer = ""
			as.printPrompt(*currentIndex+1, maxIndex, "")
		}
	}
}

// navigateUp handles upward navigation with wraparound
func (as *ArchetypeSelectionStage) navigateUp(currentIndex, maxIndex int) int {
	if currentIndex <= 0 {
		return maxIndex - 1 // Wrap to last item
	}
	return currentIndex - 1
}

// navigateDown handles downward navigation with wraparound
func (as *ArchetypeSelectionStage) navigateDown(currentIndex, maxIndex int) int {
	if currentIndex >= maxIndex-1 {
		return 0 // Wrap to first item
	}
	return currentIndex + 1
}

// updateDisplay refreshes the table and prompt
func (as *ArchetypeSelectionStage) updateDisplay(archetypes []*ArchetypeDefinition, formatter *components.TableFormatter, columnWidths []int, currentIndex int, numberBuffer string) {
	// Move cursor up to table start and re-render
	fmt.Printf("\033[%dA", len(archetypes)+3) // Move up past table + separator + prompt
	as.renderTableWithHighlight(archetypes, formatter, columnWidths, currentIndex, false)
	as.renderPrompt(currentIndex+1, len(archetypes), numberBuffer)
}

// renderTableWithHighlight renders the table with current selection highlighted
func (as *ArchetypeSelectionStage) renderTableWithHighlight(archetypes []*ArchetypeDefinition, formatter *components.TableFormatter, columnWidths []int, highlightIndex int, isSelected bool) {
	for i, archetype := range archetypes {
		var selectionCol string
		var nameColor string
		var nameStyle string

		if isSelected && i == highlightIndex {
			// Final selection: green with checkmark
			selectionCol = " ✓"
			nameColor = "92" // Bright green
			nameStyle = ""
		} else if i == highlightIndex {
			// Current highlight: blue italic with arrow
			selectionCol = " >"
			nameColor = "94" // Bright blue
			nameStyle = "\033[3m" // Italic
		} else {
			// Default: normal white
			selectionCol = "  "
			nameColor = "97" // White
			nameStyle = ""
		}

		number := fmt.Sprintf("%02d.", i+1)
		framework := as.getFrameworkDisplay(archetype)
		language := as.getLanguageDisplay(archetype)
		appType := fmt.Sprintf("'%s'", archetype.ID)

		// Apply styling to the name column
		styledName := fmt.Sprintf("%s\033[%sm%s\033[0m", nameStyle, nameColor, archetype.Name)

		values := []string{selectionCol, number, styledName, framework, language, appType}
		colors := []string{"92", "90", "", "90", "90", "94"} // Empty color for name since we handle it manually

		rowText := formatter.FormatRowWithWidths(values, colors, columnWidths)
		fmt.Printf("\033[K%s\n", rowText) // Clear line and print row
	}

	// Print separator
	width := getTerminalWidthForSelector()
	separator := strings.Repeat("-", width)
	fmt.Printf("\033[K\033[90m%s\033[0m\n", separator)
}

// renderPrompt displays the selection prompt
func (as *ArchetypeSelectionStage) renderPrompt(currentSelection, total int, numberBuffer string) {
	if numberBuffer != "" {
		fmt.Printf("\r\033[K? Select an archetype (1-%d): %s", total, numberBuffer)
	} else {
		fmt.Printf("\r\033[K? Select an archetype (1-%d): %d", total, currentSelection)
	}
}

// fallbackSelection implements the original simple selection mode
func (as *ArchetypeSelectionStage) fallbackSelection(archetypes []*ArchetypeDefinition, formatter *components.TableFormatter, columnWidths []int) (*ArchetypeDefinition, error) {
	// Print each archetype as a numbered row using calculated widths (original implementation)
	for i, archetype := range archetypes {
		selectionCol := "  " // Empty for simple mode
		number := fmt.Sprintf("%02d.", i+1)

		// Map archetype to display values
		framework := as.getFrameworkDisplay(archetype)
		language := as.getLanguageDisplay(archetype)
		appType := fmt.Sprintf("'%s'", archetype.ID)

		// Format the row using calculated widths (same pattern as templates view)
		values := []string{selectionCol, number, archetype.Name, framework, language, appType}
		colors := []string{"90", "90", "97", "90", "90", "94"}
		rowText := formatter.FormatRowWithWidths(values, colors, columnWidths)

		fmt.Println(rowText)
	}

	// Width-aware separator
	width := getTerminalWidthForSelector()
	separator := strings.Repeat("-", width)
	fmt.Printf("\n\033[90m%s\033[0m\n\n", separator)

	// Simple selection loop (original implementation)
	for {
		fmt.Print("? Select an archetype (1-" + strconv.Itoa(len(archetypes)) + "): ")

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
		if err != nil || selectedIndex < 1 || selectedIndex > len(archetypes) {
			fmt.Printf("\033[91mInvalid selection. Please enter a number between 1 and %d.\033[0m\n", len(archetypes))
			continue
		}

		// Get the selected archetype
		selected := archetypes[selectedIndex-1]

		// Display selection confirmation with new styling
		fmt.Printf("\n\033[92m[✓] %s\033[0m\n", selected.Name)

		return selected, nil
	}
}

// renderInteractiveTable renders the table with highlighting for the current selection
func (as *ArchetypeSelectionStage) renderInteractiveTable(archetypes []*ArchetypeDefinition, formatter *components.TableFormatter, columnWidths []int, highlightIndex int) {
	// Clear screen completely and start fresh
	fmt.Printf("\033[2J\033[H")

	// Re-render the wizard header
	header := components.NewHeader("ENGX NEW APPLICATION WIZARD")
	fmt.Print(header.Render())

	// Re-render the intro text
	fmt.Printf("\033[90mLet's get your application configured and setup for development, this is the \033[93mGUIDED MODE\033[90m that is the default with \033[38;5;198mengx\033[0m \033[38;5;208mcreate\033[0m\033[90m. You can bypass this by specifying an \033[38;5;48m--app-type\033[0m\033[90m directly in the create command:\033[0m\n\n")

	// Smart command formatting with syntax highlighting
	cmdFormatter := components.NewCommandFormatter()
	fmt.Printf("%s\n\n", cmdFormatter.FormatCommandInBackticks("engx create <AppName> --app-type <app-type>"))

	// Re-render the selection header
	selectionHeader := components.NewHeader("SELECT YOUR APPLICATION TYPE")
	fmt.Print(selectionHeader.Render())
	fmt.Printf("%s\n", formatter.FormatHeaderWithWidths(columnWidths))

	// Render each row with highlighting
	for i, archetype := range archetypes {
		var selectionCol string
		var nameColor string
		var nameStyle string

		if i == highlightIndex {
			// Current highlight: blue italic with arrow
			selectionCol = " >"
			nameColor = "94" // Bright blue
			nameStyle = "\033[3m" // Italic
		} else {
			// Default: normal white
			selectionCol = "  "
			nameColor = "97" // White
			nameStyle = ""
		}

		number := fmt.Sprintf("%02d.", i+1)
		framework := as.getFrameworkDisplay(archetype)
		language := as.getLanguageDisplay(archetype)
		appType := fmt.Sprintf("'%s'", archetype.ID)

		// Apply styling to the name column
		styledName := fmt.Sprintf("%s\033[%sm%s\033[0m", nameStyle, nameColor, archetype.Name)

		values := []string{selectionCol, number, styledName, framework, language, appType}
		colors := []string{"90", "90", "", "90", "90", "94"} // Empty color for name since we handle it manually

		rowText := formatter.FormatRowWithWidths(values, colors, columnWidths)
		fmt.Println(rowText)
	}

	// Width-aware separator
	width := getTerminalWidthForSelector()
	separator := strings.Repeat("-", width)
	fmt.Printf("\n\033[90m%s\033[0m\n", separator)

	// Initial prompt
	fmt.Printf("\n ? Select an archetype (1-%d): %d", len(archetypes), highlightIndex+1)
}

// updateInteractiveDisplay updates the highlighting and prompt by clearing and re-rendering
func (as *ArchetypeSelectionStage) updateInteractiveDisplay(archetypes []*ArchetypeDefinition, formatter *components.TableFormatter, columnWidths []int, currentIndex int, numberBuffer string) {
	// Simple approach: clear screen and re-render everything
	as.renderInteractiveTable(archetypes, formatter, columnWidths, currentIndex)

	// Update the prompt with current state
	if numberBuffer != "" {
		fmt.Printf("\r ? Select an archetype (1-%d): %s", len(archetypes), numberBuffer)
	}
}

// updatePromptOnly updates only the prompt line for number input
func (as *ArchetypeSelectionStage) updatePromptOnly(currentSelection, total int, numberBuffer string) {
	if numberBuffer != "" {
		fmt.Printf("\r\033[K ? Select an archetype (1-%d): %s", total, numberBuffer)
	} else {
		fmt.Printf("\r\033[K ? Select an archetype (1-%d): %d", total, currentSelection)
	}
}

// handleNumberInputInteractive processes debounced number input for interactive mode
func (as *ArchetypeSelectionStage) handleNumberInputInteractive(currentIndex *int, numberBuffer *string, lastInputTime *time.Time, maxIndex int, archetypes []*ArchetypeDefinition, formatter *components.TableFormatter, columnWidths []int) {
	time.Sleep(200 * time.Millisecond) // Debounce delay

	// Check if input has stopped
	if time.Since(*lastInputTime) >= 200*time.Millisecond && *numberBuffer != "" {
		if num, err := strconv.Atoi(*numberBuffer); err == nil && num >= 1 && num <= maxIndex {
			*currentIndex = num - 1
			*numberBuffer = ""
			as.updateInteractiveDisplay(archetypes, formatter, columnWidths, *currentIndex, "")
		} else {
			// Invalid number, clear buffer
			*numberBuffer = ""
			as.updatePromptOnly(*currentIndex+1, maxIndex, "")
		}
	}
}

// printStaticTable renders the table once without dynamic updates (for fallback mode)
func (as *ArchetypeSelectionStage) printStaticTable(archetypes []*ArchetypeDefinition, formatter *components.TableFormatter, columnWidths []int) {
	for i, archetype := range archetypes {
		selectionCol := "  " // Empty for static display
		number := fmt.Sprintf("%02d.", i+1)
		framework := as.getFrameworkDisplay(archetype)
		language := as.getLanguageDisplay(archetype)
		appType := fmt.Sprintf("'%s'", archetype.ID)

		values := []string{selectionCol, number, archetype.Name, framework, language, appType}
		colors := []string{"90", "90", "97", "90", "90", "94"}
		rowText := formatter.FormatRowWithWidths(values, colors, columnWidths)

		fmt.Println(rowText)
	}

	// Width-aware separator
	width := getTerminalWidthForSelector()
	separator := strings.Repeat("-", width)
	fmt.Printf("\n\033[90m%s\033[0m\n\n", separator)
}

// printPrompt displays and updates only the prompt line (for fallback mode)
func (as *ArchetypeSelectionStage) printPrompt(currentSelection, total int, numberBuffer string) {
	if numberBuffer != "" {
		fmt.Printf("\r\033[K ? Select an archetype (1-%d): %s", total, numberBuffer)
	} else {
		fmt.Printf("\r\033[K ? Select an archetype (1-%d): %d", total, currentSelection)
	}
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