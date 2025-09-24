package components

import (
	"regexp"
	"strings"
)

// BadgeAligner provides consistent right-alignment of badges/labels to terminal edge
type BadgeAligner struct {
	width int
}

// NewBadgeAligner creates a new badge aligner for the specified terminal width
func NewBadgeAligner(width int) *BadgeAligner {
	if width < 20 {
		width = 20 // Minimum width
	}
	return &BadgeAligner{width: width}
}

// AlignRight takes a content string and a badge string, and returns the content
// with the badge right-aligned to the terminal edge with proper spacing
func (b *BadgeAligner) AlignRight(content, badge string, minGutter int) string {
	if badge == "" {
		return content
	}

	if minGutter < 1 {
		minGutter = 1 // Minimum 1 space gutter
	}

	// Calculate display width without ANSI codes
	contentDisplayWidth := b.getDisplayWidth(content)
	badgeDisplayWidth := b.getDisplayWidth(badge) // Badge may contain ANSI codes

	// Calculate padding needed: total_width - content_width - badge_width
	paddingNeeded := b.width - contentDisplayWidth - badgeDisplayWidth

	// Ensure minimum gutter
	if paddingNeeded < minGutter {
		paddingNeeded = minGutter
	}

	return content + strings.Repeat(" ", paddingNeeded) + badge
}

// getDisplayWidth returns the display width of a string, stripping ANSI color codes
// Uses the same proven approach as the header/footer component
func (b *BadgeAligner) getDisplayWidth(text string) int {
	// Comprehensive regex to match all ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	stripped := ansiRegex.ReplaceAllString(text, "")
	return len(stripped)
}

// GetWidth returns the configured terminal width
func (b *BadgeAligner) GetWidth() int {
	return b.width
}