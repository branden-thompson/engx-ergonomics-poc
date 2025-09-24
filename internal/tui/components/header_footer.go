package components

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/design"
)

// HeaderFooterComponent provides consistent terminal-width-aware headers, footers, and separators
// with built-in default colors for consistent styling across all plugins
type HeaderFooterComponent struct {
	width int
}

// NewHeaderFooterComponent creates a new header/footer component
func NewHeaderFooterComponent(width int) *HeaderFooterComponent {
	if width < 40 {
		width = 40 // Minimum width
	}
	return &HeaderFooterComponent{width: width}
}

// RenderHeader creates a header line with optional left and right labels
// Format: ----<left-label>---{...}---<right-label>----
// Uses default colors: grey dashes, white labels
func (h *HeaderFooterComponent) RenderHeader(leftLabel, rightLabel string) string {
	return h.renderLine(leftLabel, rightLabel, "", "")
}

// RenderHeaderWithColors creates a header line with custom label colors
// Pass empty string for default white color
func (h *HeaderFooterComponent) RenderHeaderWithColors(leftLabel, rightLabel, leftColor, rightColor string) string {
	return h.renderLine(leftLabel, rightLabel, leftColor, rightColor)
}

// RenderFooter creates a footer line with optional left and right labels
// Format: ----<left-label>---{...}---<right-label>----
// Uses default colors: grey dashes, white labels
func (h *HeaderFooterComponent) RenderFooter(leftLabel, rightLabel string) string {
	return h.renderLine(leftLabel, rightLabel, "", "")
}

// RenderFooterWithColors creates a footer line with custom label colors
// Pass empty string for default white color
func (h *HeaderFooterComponent) RenderFooterWithColors(leftLabel, rightLabel, leftColor, rightColor string) string {
	return h.renderLine(leftLabel, rightLabel, leftColor, rightColor)
}

// RenderSeparator creates a plain separator line
// Format: ----{...}----
// Uses default grey color for dashes
func (h *HeaderFooterComponent) RenderSeparator() string {
	return h.renderLine("", "", "", "")
}

// renderLine is the core implementation for all line types with color support
func (h *HeaderFooterComponent) renderLine(leftLabel, rightLabel, leftColor, rightColor string) string {
	// Apply default colors
	dashColor := design.ColorDarkGray
	if leftColor == "" {
		leftColor = design.ColorBrightWhite // Default white for left label
	}
	if rightColor == "" {
		rightColor = design.ColorBrightWhite // Default white for right label
	}

	// Build colored components
	lead := h.colorize("----", dashColor)
	tail := h.colorize("----", dashColor)

	// Prepare labels with spaces and colors if they exist
	leftLabelWithSpaces := ""
	if leftLabel != "" {
		coloredLabel := h.colorize(leftLabel, leftColor)
		leftLabelWithSpaces = fmt.Sprintf(" %s ", coloredLabel)
	}

	rightLabelWithSpaces := ""
	if rightLabel != "" {
		coloredLabel := h.colorize(rightLabel, rightColor)
		rightLabelWithSpaces = fmt.Sprintf(" %s ", coloredLabel)
	}

	// Calculate remaining width for fill dashes (using display width for calculation to handle ANSI codes)
	usedWidth := 4 + h.getDisplayWidth(leftLabel) + h.getDisplayWidth(rightLabel) + 4 // lead + labels + tail (without color codes)
	if leftLabel != "" {
		usedWidth += 2 // spaces around left label
	}
	if rightLabel != "" {
		usedWidth += 2 // spaces around right label
	}

	remainingWidth := h.width - usedWidth
	if remainingWidth < 4 {
		remainingWidth = 4 // Minimum fill
	}

	// Create colored fill dashes
	fillDashes := h.colorize(strings.Repeat("-", remainingWidth), dashColor)

	// Build the complete line
	return fmt.Sprintf("%s%s%s%s%s\n",
		lead,
		leftLabelWithSpaces,
		fillDashes,
		rightLabelWithSpaces,
		tail)
}

// colorize applies ANSI color codes to text
func (h *HeaderFooterComponent) colorize(text, color string) string {
	if color == "" || text == "" {
		return text
	}
	return fmt.Sprintf("\033[%sm%s\033[%sm", color, text, design.ColorReset)
}

// GetWidth returns the configured width
func (h *HeaderFooterComponent) GetWidth() int {
	return h.width
}

// getDisplayWidth returns the display width of a string, stripping ANSI color codes
func (h *HeaderFooterComponent) getDisplayWidth(text string) int {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	stripped := ansiRegex.ReplaceAllString(text, "")
	return len(stripped)
}