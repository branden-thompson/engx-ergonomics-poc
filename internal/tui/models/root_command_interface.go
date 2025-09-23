package models

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/term"

	"github.com/bthompso/engx-ergonomics-poc/internal/config"
	"github.com/bthompso/engx-ergonomics-poc/pkg/common"
)

// ANSI colors for styling (matching create command)
const (
	colorReset         = "\033[0m"
	colorWhite         = "\033[97m"  // Bright white
	colorBrightBlue    = "\033[94m"  // Bright blue for highlighted items
	colorItalic        = "\033[3m"   // Italic style
	colorDarkGray      = "\033[90m"  // Dark gray for numbers/separators
	colorLightPurple   = "\033[38;5;135m" // Light purple for headers
	colorBrightMagenta = "\033[95m"  // Bright magenta for backticks
	colorHotPink       = "\033[38;5;198m" // Hot pink for "engx"
	colorOrange        = "\033[38;5;208m" // Orange for command names
	colorGreen         = "\033[92m"  // Green for success/selected
	colorGray          = "\033[37m"  // Gray for descriptions
)

// Column widths (exact specifications from requirements)
const (
	iconColWidth    = 3  // '   ', ' > ', ' ✓ '
	numberColWidth  = 4  // '01. '
	commandColWidth = 20 // '`engx longcommand`'
	crewColWidth    = 13 // 'CREW-1234    '
	// useToColWidth is flexible (fit to space)
)

// InteractiveCommandSelector represents the root command interactive interface
type InteractiveCommandSelector struct {
	roadmapConfig    *config.RoadmapConfig
	commandDiscovery *common.CommandDiscovery
	commands         []common.CommandEntry
	categories       map[string][]common.CommandEntry

	// UI state
	highlightedIndex int
	selectedIndex    int // -1 if none selected yet
	showHelp         bool
	showDevTools     bool
	terminalWidth    int
}

// NewInteractiveCommandSelector creates a new interactive command selector
func NewInteractiveCommandSelector(roadmapConfig *config.RoadmapConfig, commandDiscovery *common.CommandDiscovery, showDevTools bool) (*InteractiveCommandSelector, error) {
	// Discover commands
	commands, err := commandDiscovery.DiscoverCommands()
	if err != nil {
		return nil, fmt.Errorf("failed to discover commands: %w", err)
	}

	categories := commandDiscovery.GetCommandCategories(commands)

	selector := &InteractiveCommandSelector{
		roadmapConfig:    roadmapConfig,
		commandDiscovery: commandDiscovery,
		commands:         commands,
		categories:       categories,
		highlightedIndex: 0,
		selectedIndex:    -1, // No selection yet
		showHelp:         false,
		showDevTools:     showDevTools,
		terminalWidth:    80, // Default
	}

	// Get terminal width
	selector.terminalWidth = getTerminalWidth()

	// Find first selectable command for initial highlight
	selector.findFirstSelectableCommand()

	return selector, nil
}

// getTerminalWidth gets the current terminal width (same logic as layout components)
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

// findFirstSelectableCommand sets highlight to first available command
func (s *InteractiveCommandSelector) findFirstSelectableCommand() {
	for i, cmd := range s.commands {
		if cmd.Selectable {
			s.highlightedIndex = i
			return
		}
	}
	// If no selectable commands, stay at 0
	s.highlightedIndex = 0
}

// calculateUseToWidth calculates the flexible USE TO column width
func (s *InteractiveCommandSelector) calculateUseToWidth() int {
	usedWidth := iconColWidth + numberColWidth + commandColWidth + crewColWidth + 8 // padding
	remaining := s.terminalWidth - usedWidth
	if remaining < 20 {
		return 20 // minimum width
	}
	return remaining
}

// Run starts the command selection process
func (s *InteractiveCommandSelector) Run() (*common.CommandEntry, error) {
	// Render the interface
	s.renderInterface()

	// Simple prompt at the end
	fmt.Print("? Where would you like to start? (1-n): ")

	// Read simple input
	var choice int
	_, err := fmt.Scanf("%d", &choice)
	if err != nil {
		return nil, fmt.Errorf("invalid input")
	}

	// Find command by number
	for _, cmd := range s.commands {
		if cmd.Number == choice {
			if cmd.Selectable {
				return &cmd, nil
			} else {
				return nil, fmt.Errorf("command %d is not yet available", choice)
			}
		}
	}

	return nil, fmt.Errorf("command %d not found", choice)
}

// renderInterface displays the command selection interface
func (s *InteractiveCommandSelector) renderInterface() {
	// Render animated header
	s.renderAnimatedHeader()

	// Render commands by category
	for _, category := range s.roadmapConfig.CommandCategories {
		categoryName := s.roadmapConfig.GetFormattedCategoryName(category.Name)
		commands := s.categories[categoryName]

		// Skip developer tools category unless showDevTools is true
		if categoryName == "EngX CLI Developer Tools" && !s.showDevTools {
			continue
		}

		if len(commands) > 0 {
			// Format category header with blue @username highlighting
			formattedCategoryName := s.formatCategoryHeader(categoryName)
			categoryHeader := "\n" + formattedCategoryName + ":" + colorReset + "\n"
			fmt.Print(categoryHeader)
			s.renderTableHeader()

			for _, cmd := range commands {
				s.renderCommand(cmd)
			}
		}
	}

	// Render footer separator
	s.renderFooter()
}

// renderTableHeader renders the table header row with dynamic width
func (s *InteractiveCommandSelector) renderTableHeader() {
	// Calculate dynamic column widths
	numberColWidth := 7      // "   01. "
	commandColWidth := 22    // "`engx templates` " (with padding)
	crewColWidth := 13       // "CREW-1234    "

	// USE TO column fills remaining space with 4-character gutter before CREW ID
	useToColWidth := s.terminalWidth - numberColWidth - commandColWidth - crewColWidth - 4 // 4-char gutter
	if useToColWidth < 20 {
		useToColWidth = 20 // minimum width
	}

	// Build header with proper spacing and 4-character gutter
	headerRow := "\n" + colorLightPurple +
		"   ##  " +
		"COMMAND" + strings.Repeat(" ", commandColWidth-7) + // 7 = len("COMMAND")
		"USE TO" + strings.Repeat(" ", useToColWidth-6) +    // 6 = len("USE TO")
		strings.Repeat(" ", 4) +                             // 4-character gutter
		"OWNER CREW ID" +
		colorReset + "\n"
	fmt.Print(headerRow)
}

// renderCommand renders a single command row in the template format
func (s *InteractiveCommandSelector) renderCommand(cmd common.CommandEntry) {
	// Format number with leading zero and proper spacing: "   01. "
	number := fmt.Sprintf("   %02d. ", cmd.Number)

	// Format command with syntax coloring for ALL commands
	coloredCommandText := colorBrightMagenta + "`" +
		colorHotPink + "engx" + " " +
		colorOrange + cmd.Command +
		colorBrightMagenta + "`" + colorReset

	// Format description
	description := cmd.UseCase

	// Calculate dynamic USE TO column width first
	numberColWidth := 7      // "   01. "
	commandColWidth := 22    // "`engx templates` " (with padding)
	crewColWidth := 13       // "CREW-1234    "

	// USE TO column fills remaining space with 4-character gutter before CREW ID
	useToColWidth := s.terminalWidth - numberColWidth - commandColWidth - crewColWidth - 4 // 4-char gutter
	if useToColWidth < 20 {
		useToColWidth = 20 // minimum width
	}

	// Wrap description text for the USE TO column
	wrappedDesc := s.wrapText(description, useToColWidth)

	// Get colors for text - consistent across ALL rows regardless of status
	numberColor := colorDarkGray      // Grey numbers for all rows
	descColor := colorDarkGray        // Dark grey descriptions for all rows (more subdued than headers)
	crewColor := colorBrightBlue      // Light blue CREW IDs for all rows

	// Use the column widths calculated above

	// Calculate padding for command column
	plainCommandText := "`engx " + cmd.Command + "`"
	commandPadding := commandColWidth - len(plainCommandText)
	if commandPadding < 0 {
		commandPadding = 0
	}

	// Print first line
	firstLine := ""
	if len(wrappedDesc) > 0 {
		firstLine = wrappedDesc[0]
	}

	// Calculate padding for USE TO column to right-align OWNER CREW ID
	descPadding := useToColWidth - len(firstLine)
	if descPadding < 0 {
		descPadding = 0
	}

	// Add explicit 4-character gutter between USE TO and CREW ID
	gutter := strings.Repeat(" ", 4)

	// Build the row with dynamic spacing
	row := numberColor + number + colorReset +
		coloredCommandText + strings.Repeat(" ", commandPadding) +
		descColor + firstLine + colorReset + strings.Repeat(" ", descPadding) +
		gutter +
		crewColor + cmd.OwnerCrew + colorReset + "\n"
	fmt.Print(row)

	// Print additional lines for wrapped description with proper indentation
	if len(wrappedDesc) > 1 {
		// Indent to align with description column
		indentSize := numberColWidth + commandColWidth
		indent := strings.Repeat(" ", indentSize)

		for _, line := range wrappedDesc[1:] {
			// Calculate padding for this line to maintain column alignment
			linePadding := useToColWidth - len(line)
			if linePadding < 0 {
				linePadding = 0
			}

			wrappedRow := indent +
				descColor + line + strings.Repeat(" ", linePadding) + colorReset + "\n"
			fmt.Print(wrappedRow)
		}

		// Add blank line after multi-line descriptions for readability
		fmt.Print("\n")
	}
}

// wrapText wraps text to fit within specified width
func (s *InteractiveCommandSelector) wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		// Check if adding this word would exceed width
		if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > width {
			// Start new line
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}

		// Add word to current line
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}

	// Add final line
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// formatCategoryHeader formats category headers with blue @username highlighting
func (s *InteractiveCommandSelector) formatCategoryHeader(categoryName string) string {
	// Check if this category contains @username
	if strings.Contains(categoryName, "@") {
		// Split on @ to separate the text and username
		parts := strings.Split(categoryName, "@")
		if len(parts) == 2 {
			// Extract username (remove closing parenthesis if present)
			username := parts[1]
			if strings.HasSuffix(username, ")") {
				username = strings.TrimSuffix(username, ")")
				// Format: "Specific to you (" + colorGreen + "@username" + colorWhite + ")"
				return colorWhite + parts[0] + colorGreen + "@" + username + colorWhite + ")"
			} else {
				// Format: "Specific to you (" + colorGreen + "@username" + colorWhite
				return colorWhite + parts[0] + colorGreen + "@" + username + colorWhite
			}
		}
	}

	// Check if this is the main EngX CLI header
	if strings.Contains(categoryName, "EngX CLI") {
		// Create animated multi-colored EngX CLI
		return s.formatAnimatedEngXCLI(categoryName)
	}

	// Default: just white
	return colorWhite + categoryName
}

// renderAnimatedHeader renders the main header with multi-colored EngX CLI
func (s *InteractiveCommandSelector) renderAnimatedHeader() {
	// Create smooth gradient interpolation for "EngX CLI"
	animatedEngX := s.createGradientText("EngX CLI")

	// Build the header line manually to include animated text
	leftPadding := colorDarkGray + "----" + colorReset + " "
	rightText := " " + colorWhite + "OCT 2025 Release" + colorReset + " " + colorDarkGray + "----" + colorReset
	version := colorWhite + " v0.8.0 " + colorReset

	// Calculate padding
	totalWidth := s.terminalWidth
	usedWidth := 4 + 1 + len("EngX CLI") + len(" v0.8.0 ") + 1 + len("OCT 2025 Release") + 1 + 4 // rough calculation
	middlePadding := totalWidth - usedWidth
	if middlePadding < 0 {
		middlePadding = 0
	}

	headerLine := leftPadding + animatedEngX + version + strings.Repeat(colorDarkGray + "-" + colorReset, middlePadding) + rightText + "\n"
	fmt.Print(headerLine)
}

// createGradientText creates a smooth color gradient across the given text
func (s *InteractiveCommandSelector) createGradientText(text string) string {
	// Define the three key colors for interpolation
	// #DD51D6 (221, 81, 214) -> #378FE9 (55, 143, 233) -> #7CE3B3 (124, 227, 179)

	textRunes := []rune(text)
	textLength := len(textRunes)
	if textLength == 0 {
		return ""
	}

	result := ""
	for i, char := range textRunes {
		// Calculate position as a float from 0.0 to 1.0
		position := float64(i) / float64(textLength-1)

		var r, g, b int

		if position <= 0.5 {
			// Interpolate between color1 (#DD51D6) and color2 (#378FE9)
			t := position * 2.0 // Scale to 0.0-1.0 for first half
			r = int(221.0*(1.0-t) + 55.0*t)
			g = int(81.0*(1.0-t) + 143.0*t)
			b = int(214.0*(1.0-t) + 233.0*t)
		} else {
			// Interpolate between color2 (#378FE9) and color3 (#7CE3B3)
			t := (position - 0.5) * 2.0 // Scale to 0.0-1.0 for second half
			r = int(55.0*(1.0-t) + 124.0*t)
			g = int(143.0*(1.0-t) + 227.0*t)
			b = int(233.0*(1.0-t) + 179.0*t)
		}

		// Create ANSI escape sequence for this color
		colorCode := fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
		result += colorCode + string(char)
	}

	result += colorReset
	return result
}

// formatAnimatedEngXCLI creates a multi-colored animated EngX CLI text
func (s *InteractiveCommandSelector) formatAnimatedEngXCLI(categoryName string) string {
	if strings.Contains(categoryName, "EngX CLI Developer Tools") {
		// Create smooth gradient for "EngX CLI" in developer tools
		gradientEngX := s.createGradientText("EngX CLI")
		return gradientEngX + colorWhite + " Developer Tools"
	}

	// Default white for other EngX CLI references
	return colorWhite + categoryName
}

// renderFooter renders just the separator line
func (s *InteractiveCommandSelector) renderFooter() {
	// Add empty line before footer separator for readability
	fmt.Print("\n")

	// Simple separator line using our ANSI colors directly
	separatorLine := strings.Repeat("-", s.terminalWidth)
	fmt.Print(colorDarkGray + separatorLine + colorReset + "\n")
}

// handleInput processes user input and returns the selected command
func (s *InteractiveCommandSelector) handleInput() (*common.CommandEntry, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		// Read single character
		char, _, err := reader.ReadRune()
		if err != nil {
			return nil, fmt.Errorf("error reading input: %w", err)
		}

		switch char {
		case 'q', '\x03': // q or Ctrl+C
			fmt.Println("\nGoodbye!")
			return nil, fmt.Errorf("user quit")

		case 'j', '\x1b': // j or down arrow (simplified)
			s.moveNext()
			s.updateDisplay()

		case 'k': // k or up arrow
			s.movePrevious()
			s.updateDisplay()

		case '\r', '\n': // Enter
			if s.highlightedIndex >= 0 && s.highlightedIndex < len(s.commands) {
				cmd := s.commands[s.highlightedIndex]
				if cmd.Selectable {
					s.selectedIndex = s.highlightedIndex
					s.updateDisplay() // Show selection
					return &cmd, nil
				} else {
					fmt.Printf("\nCommand '%s' is not yet available.\n", cmd.Command)
				}
			}

		case 'h', '?':
			s.showHelp = !s.showHelp
			s.updateDisplay()

		default:
			// Handle number selection
			if char >= '0' && char <= '9' {
				num := int(char - '0')
				if num == 0 {
					num = 10 // 0 maps to command 10
				}

				// Find command by number
				for i, cmd := range s.commands {
					if cmd.Number == num {
						if cmd.Selectable {
							s.selectedIndex = i
							s.highlightedIndex = i
							s.updateDisplay() // Show selection
							return &cmd, nil
						} else {
							fmt.Printf("\nCommand %d (%s) is not yet available.\n", num, cmd.Command)
						}
						break
					}
				}
			}
		}
	}
}

// moveNext moves highlight to next selectable command
func (s *InteractiveCommandSelector) moveNext() {
	for i := s.highlightedIndex + 1; i < len(s.commands); i++ {
		if s.commands[i].Selectable {
			s.highlightedIndex = i
			return
		}
	}

	// Wrap to beginning
	for i := 0; i < s.highlightedIndex; i++ {
		if s.commands[i].Selectable {
			s.highlightedIndex = i
			return
		}
	}
}

// movePrevious moves highlight to previous selectable command
func (s *InteractiveCommandSelector) movePrevious() {
	for i := s.highlightedIndex - 1; i >= 0; i-- {
		if s.commands[i].Selectable {
			s.highlightedIndex = i
			return
		}
	}

	// Wrap to end
	for i := len(s.commands) - 1; i > s.highlightedIndex; i-- {
		if s.commands[i].Selectable {
			s.highlightedIndex = i
			return
		}
	}
}

// updateDisplay refreshes the display
func (s *InteractiveCommandSelector) updateDisplay() {
	s.renderInterface()
}

// ExecuteCommand executes the selected command
func (s *InteractiveCommandSelector) ExecuteCommand(cmd common.CommandEntry) error {
	fmt.Printf("\n%sExecuting: %sengx %s%s\n", colorGreen, colorWhite, cmd.Command, colorReset)

	// Execute the command
	execCmd := exec.Command("engx", s.commandDiscovery.FormatCommandForExecution(cmd)...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	return execCmd.Run()
}

// Legacy compatibility - keep old function name for root_router.go
func NewInteractiveCommandModel(roadmapConfig *config.RoadmapConfig, commandDiscovery *common.CommandDiscovery) (*InteractiveCommandSelector, error) {
	return NewInteractiveCommandSelector(roadmapConfig, commandDiscovery, false)
}