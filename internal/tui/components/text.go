package components

import (
	"fmt"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/design"
)

// Text provides styled text formatting utilities
type Text struct {
	content string
	style   design.TextStyle
}

// NewText creates a new text component
func NewText(content string) *Text {
	return &Text{
		content: content,
		style:   design.StyleBody,
	}
}

// Style applies a text style
func (t *Text) Style(style design.TextStyle) *Text {
	t.style = style
	return t
}

// AsTitle applies title styling
func (t *Text) AsTitle() *Text {
	t.style = design.StyleTitle
	return t
}

// AsHeader applies header styling
func (t *Text) AsHeader() *Text {
	t.style = design.StyleHeader
	return t
}

// AsSubheader applies subheader styling
func (t *Text) AsSubheader() *Text {
	t.style = design.StyleSubheader
	return t
}

// AsMuted applies muted styling
func (t *Text) AsMuted() *Text {
	t.style = design.StyleMuted
	return t
}

// AsCommand applies command styling
func (t *Text) AsCommand() *Text {
	t.style = design.StyleCommand
	return t
}

// AsSuccess applies success styling
func (t *Text) AsSuccess() *Text {
	t.style = design.StyleSuccess
	return t
}

// AsWarning applies warning styling
func (t *Text) AsWarning() *Text {
	t.style = design.StyleWarning
	return t
}

// AsError applies error styling
func (t *Text) AsError() *Text {
	t.style = design.StyleError
	return t
}

// Render outputs the styled text
func (t *Text) Render() string {
	var codes []string

	// Add color
	if t.style.Color != "" {
		codes = append(codes, t.style.Color)
	}

	// Add formatting
	if t.style.Bold {
		codes = append(codes, "1")
	}
	if t.style.Underline {
		codes = append(codes, "4")
	}

	if len(codes) > 0 {
		return fmt.Sprintf("\033[%sm%s\033[0m", strings.Join(codes, ";"), t.content)
	}

	return t.content
}

// String returns the rendered text (satisfies Stringer interface)
func (t *Text) String() string {
	return t.Render()
}

// Convenience functions for quick text styling

// Title creates title-styled text
func Title(text string) string {
	return NewText(text).AsTitle().Render()
}

// HeaderText creates header-styled text
func HeaderText(text string) string {
	return NewText(text).AsHeader().Render()
}

// Subheader creates subheader-styled text
func Subheader(text string) string {
	return NewText(text).AsSubheader().Render()
}

// Muted creates muted-styled text
func Muted(text string) string {
	return NewText(text).AsMuted().Render()
}

// Command creates command-styled text
func Command(text string) string {
	return NewText(text).AsCommand().Render()
}

// Success creates success-styled text
func Success(text string) string {
	return NewText(text).AsSuccess().Render()
}

// Warning creates warning-styled text
func Warning(text string) string {
	return NewText(text).AsWarning().Render()
}

// Error creates error-styled text
func Error(text string) string {
	return NewText(text).AsError().Render()
}

// Badge creates a styled badge/indicator
type Badge struct {
	text  string
	style design.TextStyle
}

// NewBadge creates a new badge
func NewBadge(text string) *Badge {
	return &Badge{
		text:  text,
		style: design.StyleMuted,
	}
}

// AsSuccess makes the badge green
func (b *Badge) AsSuccess() *Badge {
	b.style = design.StyleSuccess
	return b
}

// AsWarning makes the badge yellow
func (b *Badge) AsWarning() *Badge {
	b.style = design.StyleWarning
	return b
}

// AsError makes the badge red
func (b *Badge) AsError() *Badge {
	b.style = design.StyleError
	return b
}

// AsCommand makes the badge magenta
func (b *Badge) AsCommand() *Badge {
	b.style = design.StyleCommand
	return b
}

// Render outputs the styled badge
func (b *Badge) Render() string {
	return fmt.Sprintf("\033[%sm%s\033[0m", b.style.Color, b.text)
}

// Highlight creates emphasized text with background color (simulated)
func Highlight(text string) string {
	// Use bright background effect with contrasting text
	return fmt.Sprintf("\033[%s;7m %s \033[0m", design.ColorBrightWhite, text)
}

// Code creates code-formatted text (monospace simulation)
func Code(text string) string {
	return fmt.Sprintf("`%s`", text)
}

// Link creates link-styled text
func Link(text string) string {
	return fmt.Sprintf("\033[%s;4m%s\033[0m", design.ColorBrightBlue, text)
}