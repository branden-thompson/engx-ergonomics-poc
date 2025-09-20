package styles

import (
	"fmt"
	"strings"
)

// TableColumn defines a column configuration for terminal tables
type TableColumn struct {
	Header    string
	Width     int
	Color     string // ANSI color code (e.g., "90", "97", "95")
	Alignment string // "left", "right", "center"
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
		// Apply color and padding
		headerText := tf.formatCell(col.Header, col.Width, col.Alignment, col.Color)
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
		cellText := tf.formatCell(value, col.Width, col.Alignment, color)
		parts = append(parts, cellText)

		// Add spacing between columns (except last)
		if i < len(tf.columns)-1 {
			parts = append(parts, " ")
		}
	}

	return strings.Join(parts, "")
}

// formatCell formats a single cell with color, padding, and alignment
func (tf *TableFormatter) formatCell(text string, width int, alignment string, color string) string {
	// Truncate if too long
	if len(text) > width {
		text = text[:width-3] + "..."
	}

	// Apply padding based on alignment
	var paddedText string
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

	// Apply ANSI color
	if color != "" {
		return fmt.Sprintf("\033[%sm%s\033[0m", color, paddedText)
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