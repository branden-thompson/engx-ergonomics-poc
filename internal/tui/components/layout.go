package components

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/design"
	"golang.org/x/term"
)

// getTerminalWidthLayout returns the current terminal width for layout components
func getTerminalWidthLayout() int {
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

// Component-based header architecture constants
// Pattern: [Lead][Left-Label][Flex-Separator][Right-Label][Tail]
const (
	MinLeftDashes  = 4 // Lead component: ----
	MinRightDashes = 4 // Tail component: ----
)

// Header creates a formatted section header with dividers
type Header struct {
	Title        string
	RightText    string
	Subtitle     string
	Width        int
	AutoWidth    bool // If true, use terminal width
	ShowLines    bool
}

// NewHeader creates a new header component with terminal width awareness
func NewHeader(title string) *Header {
	return &Header{
		Title:     title,
		Width:     getTerminalWidthLayout(),
		AutoWidth: true,
		ShowLines: true,
	}
}

// WithRightText adds right-aligned text to the header
func (h *Header) WithRightText(rightText string) *Header {
	h.RightText = rightText
	return h
}

// WithSubtitle adds a subtitle to the header
func (h *Header) WithSubtitle(subtitle string) *Header {
	h.Subtitle = subtitle
	return h
}

// WithWidth sets a fixed header width (disables auto-width)
func (h *Header) WithWidth(width int) *Header {
	h.Width = width
	h.AutoWidth = false
	return h
}

// WithoutLines disables separator lines
func (h *Header) WithoutLines() *Header {
	h.ShowLines = false
	return h
}

// Render outputs the formatted header with proper dash spacing rules
func (h *Header) Render() string {
	var result strings.Builder

	if h.ShowLines {
		// Get current width (auto-detect if enabled)
		width := h.Width
		if h.AutoWidth {
			width = getTerminalWidthLayout()
		}

		headerLine := h.formatHeaderLine(width)
		result.WriteString(headerLine)
		result.WriteString("\n")
	} else {
		result.WriteString(fmt.Sprintf("\033[%sm%s\033[0m\n", design.ColorBrightWhite, h.Title))
	}

	if h.Subtitle != "" {
		result.WriteString("\n")
		result.WriteString(fmt.Sprintf("\033[%sm%s\033[0m\n", design.ColorDarkGray, h.Subtitle))
	}

	result.WriteString("\n")
	return result.String()
}

// formatHeaderLine implements the new component-based architecture:
// Lead (4 dashes) + Left-Label ( TEXT ) + Flex-separator (dashes) + Right-Label ( TEXT ) + Tail (4 dashes)
// Lead (4 dashes) + Left-Label ( TEXT ) + Flex-separator (dashes to fill) + Tail (4 dashes)
func (h *Header) formatHeaderLine(width int) string {
	// Component lengths
	leadLength := MinLeftDashes        // ----
	tailLength := MinRightDashes       // ----
	leftLabelLength := 0
	rightLabelLength := 0

	// Calculate label lengths (including spaces)
	if h.Title != "" {
		leftLabelLength = 1 + len(h.Title) + 1 // " TEXT "
	}
	if h.RightText != "" {
		rightLabelLength = 1 + len(h.RightText) + 1 // " TEXT "
	}

	// Calculate minimum required width
	minRequired := leadLength + leftLabelLength + rightLabelLength + tailLength
	// Add safety margin to prevent edge case wrapping
	safetyMargin := 15
	minRequired += safetyMargin

	// If terminal is too narrow, use simplified format
	if width < minRequired {
		if h.RightText != "" {
			return fmt.Sprintf("\033[%sm%s\033[0m \033[%sm%s\033[0m",
				design.ColorBrightWhite, h.Title,
				design.ColorBrightWhite, h.RightText)
		} else {
			return fmt.Sprintf("\033[%sm%s\033[0m",
				design.ColorBrightWhite, h.Title)
		}
	}

	// Calculate flex-separator length to fill remaining space
	flexSeparatorLength := width - leadLength - leftLabelLength - rightLabelLength - tailLength

	// Build components
	lead := strings.Repeat(design.BorderHorizontal, leadLength)
	tail := strings.Repeat(design.BorderHorizontal, tailLength)
	flexSeparator := strings.Repeat(design.BorderHorizontal, flexSeparatorLength)

	var result strings.Builder

	// Lead dashes
	result.WriteString(fmt.Sprintf("\033[%sm%s\033[0m", design.ColorDarkGray, lead))

	// Left label
	if h.Title != "" {
		result.WriteString(fmt.Sprintf(" \033[%sm%s\033[0m ", design.ColorBrightWhite, h.Title))
	}

	// Flex separator
	result.WriteString(fmt.Sprintf("\033[%sm%s\033[0m", design.ColorDarkGray, flexSeparator))

	// Right label
	if h.RightText != "" {
		result.WriteString(fmt.Sprintf(" \033[%sm%s\033[0m ", design.ColorBrightWhite, h.RightText))
	}

	// Tail dashes
	result.WriteString(fmt.Sprintf("\033[%sm%s\033[0m", design.ColorDarkGray, tail))

	return result.String()
}

// Separator creates a full-width separator line
type Separator struct {
	Length    int
	AutoWidth bool // If true, use terminal width
	Char      string
	Color     string
}

// NewSeparator creates a new terminal width aware separator
func NewSeparator() *Separator {
	return &Separator{
		Length:    getTerminalWidthLayout(),
		AutoWidth: true,
		Char:      design.BorderHorizontal,
		Color:     design.ColorDarkGray,
	}
}

// WithLength sets separator length
func (s *Separator) WithLength(length int) *Separator {
	s.Length = length
	return s
}

// WithColor sets separator color
func (s *Separator) WithColor(color string) *Separator {
	s.Color = color
	return s
}

// Render outputs the separator
func (s *Separator) Render() string {
	length := s.Length
	if s.AutoWidth {
		length = getTerminalWidthLayout()
	}
	line := strings.Repeat(s.Char, length)
	return fmt.Sprintf("\033[%sm%s\033[0m\n", s.Color, line)
}

// Panel creates a content panel with optional borders
type Panel struct {
	Content    string
	Title      string
	Width      int
	ShowBorder bool
	Padding    int
}

// NewPanel creates a new panel
func NewPanel(content string) *Panel {
	return &Panel{
		Content:    content,
		Width:      design.WidthWide,
		ShowBorder: false,
		Padding:    design.SpaceXS,
	}
}

// WithTitle adds a title to the panel
func (p *Panel) WithTitle(title string) *Panel {
	p.Title = title
	return p
}

// WithBorder enables border display
func (p *Panel) WithBorder() *Panel {
	p.ShowBorder = true
	return p
}

// Render outputs the formatted panel
func (p *Panel) Render() string {
	var result strings.Builder

	if p.Title != "" {
		result.WriteString(fmt.Sprintf("\033[%sm%s\033[0m\n", design.ColorBrightWhite, p.Title))
		result.WriteString("\n")
	}

	// Add padding to content
	lines := strings.Split(p.Content, "\n")
	padding := strings.Repeat(" ", p.Padding)

	for _, line := range lines {
		result.WriteString(padding + line + "\n")
	}

	return result.String()
}

// List creates a formatted list component
type List struct {
	Items      []string
	Numbered   bool
	BulletChar string
	IndentSize int
}

// NewList creates a new list
func NewList(items []string) *List {
	return &List{
		Items:      items,
		Numbered:   false,
		BulletChar: "•",
		IndentSize: design.SpaceSM,
	}
}

// AsNumbered makes the list numbered
func (l *List) AsNumbered() *List {
	l.Numbered = true
	return l
}

// WithBullet sets a custom bullet character
func (l *List) WithBullet(bullet string) *List {
	l.BulletChar = bullet
	l.Numbered = false
	return l
}

// Render outputs the formatted list
func (l *List) Render() string {
	var result strings.Builder
	indent := strings.Repeat(" ", l.IndentSize)

	for i, item := range l.Items {
		var prefix string
		if l.Numbered {
			prefix = fmt.Sprintf("%d.", i+1)
		} else {
			prefix = l.BulletChar
		}

		result.WriteString(fmt.Sprintf("%s%s %s\n", indent, prefix, item))
	}

	return result.String()
}