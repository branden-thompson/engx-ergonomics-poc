package components

import (
	"regexp"
	"strings"
)

// ColumnDefinition defines how a column should be rendered in a data table row
type ColumnDefinition struct {
	Name              string // Column name/identifier
	Header            string // Column header text
	ShowName          bool   // Whether to show the column name
	Width             int    // Fixed width (0 = fill column)
	MinWidth          int    // Minimum width for fill columns
	MaxWidth          int    // Maximum width for fill columns (0 = no limit)
	SupportsWrap      bool   // Whether content can wrap within the column width
	Alignment         string // "left", "right", "center"
	Color             string // ANSI color code
	Fill              bool   // Whether this column should fill remaining space (alternative to Width=0)
	Truncatable       bool   // Whether this column can be truncated in narrow widths
	TruncatedMinWidth int    // Minimum width this column can be truncated to
	TruncationTail    string // Truncation indicator (default: "...")
}

// DataTableLayout defines the complete table structure with data
type DataTableLayout struct {
	Columns     []ColumnDefinition
	Data        [][]string // Each row is a slice of strings matching column order
	GutterWidth int        // Space between columns (default: 1, recommended: 4)
}

// DataTableRow handles rendering of table rows with automatic width calculation
type DataTableRow struct {
	terminalWidth int
	columns       []ColumnDefinition
	gutterWidth   int    // Space between columns (configurable)
	prefix        string // Icon/indicator prefix (e.g., " ⏺ ")
	badge         string // Right-aligned badge (e.g., "(On-Call)")
	badgeGutter   int    // Minimum spaces before badge
	renderBadge   bool   // Whether to actually render the badge (default true)
}

// NewDataTableRow creates a new data table row renderer
func NewDataTableRow(terminalWidth int, columns []ColumnDefinition) *DataTableRow {
	if terminalWidth < 20 {
		terminalWidth = 20
	}
	return &DataTableRow{
		terminalWidth: terminalWidth,
		columns:       columns,
		gutterWidth:   1, // Default 1-space gutter (backward compatible)
		badgeGutter:   4, // Default 4-space gutter
		renderBadge:   true, // Default to rendering badges
	}
}

// NewDataTableRowFromLayout creates a new data table row renderer from a layout
func NewDataTableRowFromLayout(terminalWidth int, layout *DataTableLayout) *DataTableRow {
	if terminalWidth < 20 {
		terminalWidth = 20
	}

	// Use layout's gutter width if specified, otherwise default to 1
	gutterWidth := 1
	if layout.GutterWidth > 0 {
		gutterWidth = layout.GutterWidth
	}

	return &DataTableRow{
		terminalWidth: terminalWidth,
		columns:       layout.Columns,
		gutterWidth:   gutterWidth,
		badgeGutter:   4, // Default 4-space gutter
		renderBadge:   true, // Default to rendering badges
	}
}

// SetPrefix sets the prefix (icon/indicator) for the row
func (d *DataTableRow) SetPrefix(prefix string) *DataTableRow {
	d.prefix = prefix
	return d
}

// SetBadge sets the right-aligned badge and optional custom gutter
func (d *DataTableRow) SetBadge(badge string, gutter int) *DataTableRow {
	d.badge = badge
	if gutter > 0 {
		d.badgeGutter = gutter
	}
	return d
}

// SetBadgeRenderMode sets whether badges should be rendered or just reserved for space
func (d *DataTableRow) SetBadgeRenderMode(render bool) *DataTableRow {
	d.renderBadge = render
	return d
}

// RenderHeader renders the table header row
func (d *DataTableRow) RenderHeader() string {
	widths := d.calculateColumnWidths(nil) // No data needed for header

	var headerParts []string
	for i, col := range d.columns {
		width := widths[i]
		header := col.Header

		// Apply alignment and padding FIRST (without colors)
		formatted := d.formatCell(header, width, col.Alignment, true)

		// Apply color AFTER formatting to avoid width calculation issues
		if col.Color != "" {
			formatted = "\033[" + col.Color + "m" + formatted + "\033[0m" // Full ANSI escape sequence
		}

		headerParts = append(headerParts, formatted)
	}

	// Join with configurable gutter width following pattern: [prefix][number][name][gutter][attr1][gutter]...
	var headerLine strings.Builder
	headerLine.WriteString(d.prefix)
	gutter := strings.Repeat(" ", d.gutterWidth)

	for i, part := range headerParts {
		headerLine.WriteString(part)

		// Add gutter after column 2 (index 2) and onwards, but not after the last column
		if i >= 2 && i < len(headerParts)-1 {
			headerLine.WriteString(gutter)
		}
	}

	return headerLine.String()
}

// RenderRow renders a data row using the column definitions (legacy method for backward compatibility)
func (d *DataTableRow) RenderRow(data map[string]string) string {
	// Convert map data to slice format for consistency
	var rowData []string
	for _, col := range d.columns {
		value := data[col.Name] // Use Name instead of Key
		if value == "" && col.Name == "Key" { // Fallback for old Key field
			// Try to find by header name for backward compatibility
			for k, v := range data {
				if k == col.Header {
					value = v
					break
				}
			}
		}
		rowData = append(rowData, value)
	}

	return d.FormatDataRow(rowData)
}

// FormatDataRow renders a data row from a layout's data slice
func (d *DataTableRow) FormatDataRow(rowData []string) string {
	if len(rowData) != len(d.columns) {
		// Handle mismatched data gracefully
		return "Error: Data length does not match column count"
	}

	widths := d.calculateColumnWidthsFromData(rowData)

	var rowParts []string
	for i, col := range d.columns {
		width := widths[i]
		value := ""
		if i < len(rowData) {
			value = rowData[i]
		}

		// Handle truncation if content is too wide AND column is truncatable
		if d.getDisplayWidth(value) > width && col.Truncatable {
			// Use configurable truncation tail from column definition
			truncationTail := col.TruncationTail
			if truncationTail == "" {
				truncationTail = "..." // Default
			}
			value = d.truncateText(value, width, truncationTail)
		}

		// Apply alignment and padding FIRST (without colors)
		formatted := d.formatCell(value, width, col.Alignment, false)

		// Apply color AFTER formatting to avoid width calculation issues
		if col.Color != "" {
			formatted = "\033[" + col.Color + "m" + formatted + "\033[0m"
		}

		rowParts = append(rowParts, formatted)
	}

	// Join with configurable gutter width following pattern: [prefix][number][name][gutter][attr1][gutter]...
	var rowLine strings.Builder
	gutter := strings.Repeat(" ", d.gutterWidth)

	for i, part := range rowParts {
		rowLine.WriteString(part)

		// Add gutter after column 2 (index 2) and onwards, but not after the last column
		if i >= 2 && i < len(rowParts)-1 {
			rowLine.WriteString(gutter)
		}
	}

	finalRowLine := rowLine.String()

	// Add badge if specified AND rendering is enabled
	if d.badge != "" && d.renderBadge {
		finalRowLine = d.alignBadge(finalRowLine, d.badge)
	}

	return finalRowLine
}

// calculateColumnWidths determines the width of each column based on terminal width
func (d *DataTableRow) calculateColumnWidths(data map[string]string) []int {
	var widths []int
	var fillColumns []int
	usedWidth := 0

	// Calculate prefix width (without ANSI codes)
	prefixWidth := d.getDisplayWidth(d.prefix)
	usedWidth += prefixWidth

	// Calculate gutter spaces between columns using configurable gutter width
	// Following your design: no gutter at beginning, between ICON/##, or end
	// So for N columns, we have (N-3) gutters assuming first 2 are prefix columns
	gutterCount := len(d.columns) - 3
	if gutterCount < 0 {
		gutterCount = 0
	}
	gutterSpace := d.gutterWidth * gutterCount
	usedWidth += gutterSpace

	// Calculate badge space if present
	badgeSpace := 0
	if d.badge != "" {
		badgeSpace = len(d.badge) + d.badgeGutter
		usedWidth += badgeSpace
	}

	// First pass: calculate fixed column widths and identify fill columns
	for i, col := range d.columns {
		if col.Fill || col.Width == 0 {
			// Fill column - will be calculated later
			widths = append(widths, 0)
			fillColumns = append(fillColumns, i)
		} else {
			// Fixed width column
			widths = append(widths, col.Width)
			usedWidth += col.Width
		}
	}

	// Second pass: distribute remaining width among fill columns
	if len(fillColumns) > 0 {
		remainingWidth := d.terminalWidth - usedWidth
		if remainingWidth > 0 {
			fillWidthPerColumn := remainingWidth / len(fillColumns)

			for _, colIndex := range fillColumns {
				col := d.columns[colIndex]
				width := fillWidthPerColumn

				// Apply minimum width constraint
				if col.MinWidth > 0 && width < col.MinWidth {
					width = col.MinWidth
				}

				// Apply maximum width constraint
				if col.MaxWidth > 0 && width > col.MaxWidth {
					width = col.MaxWidth
				}

				widths[colIndex] = width
			}
		} else {
			// Not enough space - use minimum widths
			for _, colIndex := range fillColumns {
				col := d.columns[colIndex]
				if col.MinWidth > 0 {
					widths[colIndex] = col.MinWidth
				} else {
					widths[colIndex] = 10 // Default minimum
				}
			}
		}
	}

	return widths
}

// calculateColumnWidthsFromData determines column widths from data slice instead of map
func (d *DataTableRow) calculateColumnWidthsFromData(rowData []string) []int {
	var widths []int
	var fillColumns []int
	usedWidth := 0

	// Calculate gutter spaces between columns using configurable gutter width
	// Following your design: no gutter at beginning, between ICON/##, or end
	// So for N columns, we have (N-3) gutters assuming first 2 are prefix columns
	gutterCount := len(d.columns) - 3
	if gutterCount < 0 {
		gutterCount = 0
	}
	gutterSpace := d.gutterWidth * gutterCount
	usedWidth += gutterSpace

	// Calculate badge space if present
	if d.badge != "" {
		badgeSpace := len(d.badge) + d.badgeGutter
		usedWidth += badgeSpace
	}

	// First pass: calculate fixed column widths and identify fill columns
	for i, col := range d.columns {
		if col.Fill || col.Width == 0 {
			// Fill column - will be calculated later
			widths = append(widths, 0)
			fillColumns = append(fillColumns, i)
		} else {
			// Fixed width column
			widths = append(widths, col.Width)
			usedWidth += col.Width
		}
	}

	// Second pass: distribute remaining width among fill columns
	if len(fillColumns) > 0 {
		remainingWidth := d.terminalWidth - usedWidth
		if remainingWidth > 0 {
			fillWidthPerColumn := remainingWidth / len(fillColumns)

			for _, colIndex := range fillColumns {
				col := d.columns[colIndex]
				width := fillWidthPerColumn

				// Apply minimum width constraint
				if col.MinWidth > 0 && width < col.MinWidth {
					width = col.MinWidth
				}

				// Apply maximum width constraint
				if col.MaxWidth > 0 && width > col.MaxWidth {
					width = col.MaxWidth
				}

				widths[colIndex] = width
			}
		} else {
			// Not enough space - use minimum widths
			for _, colIndex := range fillColumns {
				col := d.columns[colIndex]
				if col.MinWidth > 0 {
					widths[colIndex] = col.MinWidth
				} else {
					widths[colIndex] = 10 // Default minimum
				}
			}
		}
	}

	// Third pass: Apply intelligent truncation if table is still too wide
	d.applyIntelligentTruncation(widths, rowData)

	return widths
}

// applyIntelligentTruncation implements the intelligent truncation algorithm
func (d *DataTableRow) applyIntelligentTruncation(widths []int, rowData []string) {
	// Calculate total used width including gutters
	totalUsedWidth := 0
	for _, width := range widths {
		totalUsedWidth += width
	}

	// Add gutter space calculation
	gutterCount := len(d.columns) - 3
	if gutterCount < 0 {
		gutterCount = 0
	}
	totalUsedWidth += d.gutterWidth * gutterCount

	// Add badge space if present
	if d.badge != "" {
		totalUsedWidth += len(d.badge) + d.badgeGutter
	}

	// Check if we need truncation
	if totalUsedWidth <= d.terminalWidth {
		return // No truncation needed
	}

	// Find truncatable columns
	var truncatableColumns []int
	for i, col := range d.columns {
		if col.Truncatable {
			truncatableColumns = append(truncatableColumns, i)
		}
	}

	if len(truncatableColumns) == 0 {
		return // No columns can be truncated
	}

	// Calculate excess width that needs to be trimmed
	excessWidth := totalUsedWidth - d.terminalWidth

	// Distribute truncation across truncatable columns
	for excessWidth > 0 && len(truncatableColumns) > 0 {
		truncationPerColumn := maxInt(1, excessWidth/len(truncatableColumns))
		var remainingTruncatable []int

		for _, colIndex := range truncatableColumns {
			col := d.columns[colIndex]
			currentWidth := widths[colIndex]

			// Determine minimum width for this column
			minWidth := col.TruncatedMinWidth
			if minWidth == 0 {
				minWidth = 5 // Default minimum to fit "..." + 2 chars
			}

			// Calculate how much we can truncate this column
			canTruncate := currentWidth - minWidth
			if canTruncate <= 0 {
				continue // This column is already at minimum
			}

			// Truncate this column
			toTruncate := minInt(truncationPerColumn, canTruncate)
			widths[colIndex] -= toTruncate
			excessWidth -= toTruncate

			// Keep this column in the list if it can be truncated further
			if widths[colIndex] > minWidth {
				remainingTruncatable = append(remainingTruncatable, colIndex)
			}
		}

		truncatableColumns = remainingTruncatable
		if len(remainingTruncatable) == 0 {
			break // No more columns can be truncated
		}
	}
}

// minInt returns the smaller of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// maxInt returns the larger of two integers
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// truncateText truncates text to fit within the specified width with configurable tail
func (d *DataTableRow) truncateText(text string, width int, truncationTail string) string {
	if truncationTail == "" {
		truncationTail = "..." // Default truncation tail
	}

	tailWidth := len(truncationTail)
	if width <= tailWidth {
		return strings.Repeat(".", width)
	}
	if d.getDisplayWidth(text) <= width {
		return text
	}

	// Truncate and add configurable tail
	for i := len(text) - 1; i > 0; i-- {
		truncated := text[:i]
		if d.getDisplayWidth(truncated)+tailWidth <= width {
			return truncated + truncationTail
		}
	}
	return strings.Repeat(".", width)
}

// formatCell formats a cell value with proper alignment and padding
func (d *DataTableRow) formatCell(value string, width int, alignment string, isHeader bool) string {
	displayWidth := d.getDisplayWidth(value)

	if displayWidth >= width {
		// Value is too long - truncate for now (could be enhanced with ellipsis)
		return value // TODO: Implement smart truncation
	}

	padding := width - displayWidth

	switch alignment {
	case "right":
		return strings.Repeat(" ", padding) + value
	case "center":
		leftPad := padding / 2
		rightPad := padding - leftPad
		return strings.Repeat(" ", leftPad) + value + strings.Repeat(" ", rightPad)
	default: // "left" or unspecified
		return value + strings.Repeat(" ", padding)
	}
}

// alignBadge aligns the badge to the right edge of the terminal using header/footer logic
func (d *DataTableRow) alignBadge(content, badge string) string {
	// Use the exact same calculation as HeaderFooterComponent
	contentDisplayWidth := d.getDisplayWidth(content)
	badgeDisplayWidth := len(badge)

	// Calculate total used width: content + badge (no lead/tail like header/footer)
	totalUsedWidth := contentDisplayWidth + badgeDisplayWidth

	// Calculate padding needed to fill terminal width
	paddingNeeded := d.terminalWidth - totalUsedWidth

	// Ensure minimum gutter
	if paddingNeeded < d.badgeGutter {
		paddingNeeded = d.badgeGutter
	}

	return content + strings.Repeat(" ", paddingNeeded) + badge
}

// getDisplayWidth returns the display width of a string, stripping ANSI color codes
func (d *DataTableRow) getDisplayWidth(text string) int {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	stripped := ansiRegex.ReplaceAllString(text, "")
	return len(stripped)
}

// CalculateColumnWidthsFromData exposes the column width calculation method
func (d *DataTableRow) CalculateColumnWidthsFromData(rowData []string) []int {
	return d.calculateColumnWidthsFromData(rowData)
}

// FormatCell exposes the cell formatting method
func (d *DataTableRow) FormatCell(value string, width int, alignment string, isHeader bool) string {
	return d.formatCell(value, width, alignment, isHeader)
}

// NewEnhancedTableLayout creates a new table layout with configurable gutters
func NewEnhancedTableLayout(gutterWidth int, columns []ColumnDefinition) *DataTableLayout {
	return &DataTableLayout{
		Columns:     columns,
		GutterWidth: gutterWidth,
	}
}