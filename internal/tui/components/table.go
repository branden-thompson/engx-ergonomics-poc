package components

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/design"
	"golang.org/x/term"
)

// TableColumn defines a column configuration for terminal tables
type TableColumn struct {
	Header       string
	Width        int
	MinWidth     int    // Minimum width for this column
	MaxWidth     int    // Maximum width for this column (0 = no limit)
	Flexible     bool   // If true, this column can expand to fill available space
	Color        string // ANSI color code (e.g., "90", "97", "95")
	Alignment    string // "left", "right", "center"
	NoTruncate   bool   // If true, don't truncate text that exceeds width
}

// TableFormatter handles consistent table formatting with exact spacing
type TableFormatter struct {
	columns []TableColumn
}

// NewTableFormatter creates a new table formatter with column definitions
func NewTableFormatter(columns []TableColumn) *TableFormatter {
	return &TableFormatter{columns: columns}
}

// FormatHeader formats the table header row with proper spacing
func (tf *TableFormatter) FormatHeader() string {
	var parts []string

	for i, col := range tf.columns {
		// Apply color and padding (headers should not be truncated)
		headerText := tf.formatCell(col.Header, col.Width, col.Alignment, col.Color, false)
		parts = append(parts, headerText)

		// Add spacing between columns (except last)
		if i < len(tf.columns)-1 {
			parts = append(parts, " ")
		}
	}

	return strings.Join(parts, "")
}

// FormatRow formats a data row with proper spacing and colors
func (tf *TableFormatter) FormatRow(values []string, colors []string) string {
	var parts []string

	for i, col := range tf.columns {
		value := ""
		if i < len(values) {
			value = values[i]
		}

		color := col.Color
		if i < len(colors) && colors[i] != "" {
			color = colors[i]
		}

		// Apply color and padding
		cellText := tf.formatCell(value, col.Width, col.Alignment, color, col.NoTruncate)
		parts = append(parts, cellText)

		// Add spacing between columns (except last)
		if i < len(tf.columns)-1 {
			parts = append(parts, " ")
		}
	}

	return strings.Join(parts, "")
}

// FormatRowWithWidths formats a data row using custom column widths
func (tf *TableFormatter) FormatRowWithWidths(values []string, colors []string, widths []int) string {
	var parts []string

	for i, col := range tf.columns {
		value := ""
		if i < len(values) {
			value = values[i]
		}

		color := col.Color
		if i < len(colors) && colors[i] != "" {
			color = colors[i]
		}

		// Use dynamic width if provided, otherwise fall back to column width
		width := col.Width
		if i < len(widths) {
			width = widths[i]
		}

		// Apply color and padding
		cellText := tf.formatCell(value, width, col.Alignment, color, col.NoTruncate)
		parts = append(parts, cellText)

		// Add spacing between columns (except last)
		if i < len(tf.columns)-1 {
			parts = append(parts, " ")
		}
	}

	return strings.Join(parts, "")
}

// FormatHeaderWithWidths formats the table header with custom column widths
func (tf *TableFormatter) FormatHeaderWithWidths(widths []int) string {
	var parts []string

	for i, col := range tf.columns {
		// Use dynamic width if provided, otherwise fall back to column width
		width := col.Width
		if i < len(widths) {
			width = widths[i]
		}

		// Apply color and padding (headers should not be truncated)
		headerText := tf.formatCell(col.Header, width, col.Alignment, col.Color, false)
		parts = append(parts, headerText)

		// Add spacing between columns (except last)
		if i < len(tf.columns)-1 {
			parts = append(parts, " ")
		}
	}

	return strings.Join(parts, "")
}

// formatCell formats a single cell with color, padding, and alignment
func (tf *TableFormatter) formatCell(text string, width int, alignment string, color string, noTruncate bool) string {
	// Truncate if too long (unless noTruncate is true)
	if !noTruncate && len(text) > width {
		text = text[:width-3] + "..."
	}

	// For truncated cells, apply normal padding
	// For non-truncated cells, use minimum width but allow overflow
	var paddedText string
	if noTruncate {
		// Don't pad to exact width, just apply color
		paddedText = text
	} else {
		// Apply padding based on alignment
		switch alignment {
		case "right":
			paddedText = fmt.Sprintf("%*s", width, text)
		case "center":
			padding := width - len(text)
			leftPad := padding / 2
			rightPad := padding - leftPad
			paddedText = strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
		default: // "left"
			paddedText = fmt.Sprintf("%-*s", width, text)
		}
	}

	// Apply ANSI color
	if color != "" {
		// Handle 256-color codes (3-digit numbers) vs standard codes
		if len(color) >= 3 {
			return fmt.Sprintf("\033[38;5;%sm%s\033[0m", color, paddedText)
		} else {
			return fmt.Sprintf("\033[%sm%s\033[0m", color, paddedText)
		}
	}
	return paddedText
}

// CalculateOptimalWidths calculates column widths based on content
func CalculateOptimalWidths(headers []string, rows [][]string, minWidths []int) []int {
	widths := make([]int, len(headers))

	// Start with header lengths
	for i, header := range headers {
		widths[i] = len(header)
	}

	// Check all data rows
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Apply minimum widths
	for i, minWidth := range minWidths {
		if i < len(widths) && widths[i] < minWidth {
			widths[i] = minWidth
		}
	}

	return widths
}

// getTerminalWidth returns the current terminal width, with fallback to default
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

// CalculateFlexibleWidths calculates column widths with flexible columns that expand to fill available space
func (tf *TableFormatter) CalculateFlexibleWidths(data [][]string, reservedSpace int) []int {
	terminalWidth := getTerminalWidth()

	// Calculate space available for content (leave some margin)
	availableWidth := terminalWidth - reservedSpace - 2 // 2 for margins

	// Calculate minimum required space for fixed columns
	fixedWidth := 0
	flexibleColumns := 0
	spacingWidth := len(tf.columns) - 1 // Space between columns

	for _, col := range tf.columns {
		if col.Flexible {
			flexibleColumns++
			if col.MinWidth > 0 {
				fixedWidth += col.MinWidth
			}
		} else {
			fixedWidth += col.Width
		}
	}

	// Calculate remaining space for flexible columns
	remainingWidth := availableWidth - fixedWidth - spacingWidth
	if remainingWidth < 0 {
		remainingWidth = 0
	}

	// Calculate actual widths
	widths := make([]int, len(tf.columns))
	flexibleWidthPerColumn := remainingWidth / max(flexibleColumns, 1)

	for i, col := range tf.columns {
		if col.Flexible {
			calculatedWidth := col.MinWidth + flexibleWidthPerColumn
			if col.MaxWidth > 0 && calculatedWidth > col.MaxWidth {
				calculatedWidth = col.MaxWidth
			}
			widths[i] = calculatedWidth
		} else {
			widths[i] = col.Width
		}
	}

	return widths
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Preset table configurations for common use cases

// NewArchetypeTable creates a table formatter for application archetype listings
func NewArchetypeTable() *TableFormatter {
	columns := []TableColumn{
		{Header: "##", Width: design.ColXS, Color: design.ColorLightPurple, Alignment: "left"},
		{Header: "NAME", Width: design.ColLG, Color: design.ColorLightPurple, Alignment: "left"},
		{Header: "FRAMEWORK", Width: design.ColLG, Color: design.ColorLightPurple, Alignment: "left"},
		{Header: "LANGUAGE", Width: design.ColMD, Color: design.ColorLightPurple, Alignment: "left"},
		{Header: "--app-type", Width: design.ColMD, Color: design.ColorLightPurple, Alignment: "left"},
	}
	return NewTableFormatter(columns)
}

// NewFlexibleArchetypeTable creates a terminal width aware table formatter for application archetype listings
func NewFlexibleArchetypeTable() *TableFormatter {
	columns := []TableColumn{
		{Header: "##", Width: 3, MinWidth: 3, Color: design.ColorLightPurple, Alignment: "left"},
		{Header: "NAME", Width: 18, MinWidth: 12, Flexible: true, Color: design.ColorLightPurple, Alignment: "left"},
		{Header: "FRAMEWORK", Width: 16, MinWidth: 10, Color: design.ColorLightPurple, Alignment: "left"},
		{Header: "LANGUAGE", Width: 12, MinWidth: 8, Color: design.ColorLightPurple, Alignment: "left"},
		{Header: "--app-type", Width: 12, MinWidth: 10, Color: design.ColorLightPurple, Alignment: "left"},
	}
	return NewTableFormatter(columns)
}

// NewCommandTable creates a table formatter for command listings
func NewCommandTable() *TableFormatter {
	columns := []TableColumn{
		{Header: "COMMAND", Width: 70, Color: design.ColorBrightMagenta, Alignment: "left", NoTruncate: true}, // Don't truncate commands
		{Header: "DESCRIPTION", Width: design.WidthNarrow, Color: design.ColorDarkGray, Alignment: "left"},
	}
	return NewTableFormatter(columns)
}

// NewStatusTable creates a table formatter for status displays
func NewStatusTable() *TableFormatter {
	columns := []TableColumn{
		{Header: "ITEM", Width: design.ColLG, Color: design.ColorBrightWhite, Alignment: "left"},
		{Header: "STATUS", Width: design.ColMD, Color: design.ColorDarkGray, Alignment: "left"},
		{Header: "DETAILS", Width: design.ColXXL, Color: design.ColorDarkGray, Alignment: "left"},
	}
	return NewTableFormatter(columns)
}