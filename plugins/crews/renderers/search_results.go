package renderers

import (
	"fmt"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/models"
)

// SearchResultsRenderer handles rendering search results
type SearchResultsRenderer struct {
	colorScheme ColorScheme
}

// NewSearchResultsRenderer creates a new search results renderer
func NewSearchResultsRenderer() *SearchResultsRenderer {
	return &SearchResultsRenderer{
		colorScheme: DefaultColorScheme(),
	}
}

// Render generates the search results view
func (r *SearchResultsRenderer) Render(query string, crews []*models.Crew, width int) (string, error) {
	if width < 40 {
		width = 40 // Minimum width
	}

	var output strings.Builder

	// Header
	r.renderHeader(&output, query, len(crews), width)

	if len(crews) == 0 {
		r.renderNoResults(&output, query, width)
		r.renderFooter(&output, width)
		return output.String(), nil
	}

	// Results
	for i, crew := range crews {
		r.renderCrewResult(&output, crew, i+1, width)
	}

	// Footer with usage hints
	r.renderUsageHints(&output, width)
	r.renderFooter(&output, width)

	return output.String(), nil
}

func (r *SearchResultsRenderer) renderHeader(output *strings.Builder, query string, resultCount int, width int) {
	header := fmt.Sprintf(" SEARCH RESULTS: %s (%d found) ", query, resultCount)

	// Calculate padding
	totalPadding := width - len(header) - 4
	if totalPadding < 0 {
		totalPadding = 0
	}
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	output.WriteString(fmt.Sprintf("┌─%s%s%s%s─┐\n",
		strings.Repeat("─", leftPadding),
		r.colorScheme.HeaderColor + header + r.colorScheme.ResetColor,
		strings.Repeat("─", rightPadding),
		""))
}

func (r *SearchResultsRenderer) renderNoResults(output *strings.Builder, query string, width int) {
	noResultsMsg := fmt.Sprintf("No crews found matching: %s", query)
	output.WriteString(fmt.Sprintf("│ %s%s │\n",
		noResultsMsg,
		strings.Repeat(" ", width-len(noResultsMsg)-4)))

	output.WriteString(fmt.Sprintf("│%s│\n", strings.Repeat(" ", width-2)))

	// Suggestions
	suggestions := []string{
		"Try different search terms:",
		"• Crew ID format: CREW-1234",
		"• User email: user@company.com",
		"• Asset URN: asset://service/name",
		"• Keywords from crew names or descriptions",
	}

	for _, suggestion := range suggestions {
		output.WriteString(fmt.Sprintf("│ %s%s │\n",
			suggestion,
			strings.Repeat(" ", width-len(suggestion)-4)))
	}
}

func (r *SearchResultsRenderer) renderCrewResult(output *strings.Builder, crew *models.Crew, index int, width int) {
	// Crew header line
	crewLine := fmt.Sprintf("%d. %s%s%s (%s)",
		index,
		r.colorScheme.CrewIDColor,
		crew.VanityName,
		r.colorScheme.ResetColor,
		crew.ID)

	output.WriteString(fmt.Sprintf("│ %s%s │\n",
		crewLine,
		strings.Repeat(" ", width-len(crewLine)+len(r.colorScheme.CrewIDColor)+len(r.colorScheme.ResetColor)-4)))

	// Description (truncated if too long)
	if crew.Description != "" {
		descLines := r.wrapText(crew.Description, width-6) // Leave space for indentation
		// Show only first line for search results
		if len(descLines) > 0 {
			descLine := "   " + descLines[0]
			if len(descLines) > 1 {
				descLine += "..."
			}

			output.WriteString(fmt.Sprintf("│ %s%s │\n",
				descLine,
				strings.Repeat(" ", width-len(descLine)-4)))
		}
	}

	// Quick info line
	memberCount := len(crew.GetActiveMembers())
	assetCount := len(crew.OwnedAssets)
	onCallMembers := len(crew.GetOnCallMembers())

	infoLine := fmt.Sprintf("   Members: %d | Assets: %d", memberCount, assetCount)
	if onCallMembers > 0 {
		infoLine += fmt.Sprintf(" | On-call: %d", onCallMembers)
	}

	output.WriteString(fmt.Sprintf("│ %s%s%s%s │\n",
		r.colorScheme.TextColor,
		infoLine,
		r.colorScheme.ResetColor,
		strings.Repeat(" ", width-len(infoLine)+len(r.colorScheme.TextColor)+len(r.colorScheme.ResetColor)-4)))

	// Spacer line between results (except for last result)
	output.WriteString(fmt.Sprintf("│%s│\n", strings.Repeat(" ", width-2)))
}

func (r *SearchResultsRenderer) renderUsageHints(output *strings.Builder, width int) {
	// Usage hints section
	output.WriteString(fmt.Sprintf("├─ USAGE HINTS %s─┤\n",
		strings.Repeat("─", width-15)))

	hints := []string{
		"Use crew ID for details: engx crews CREW-1234",
		"Check user memberships: engx crews user@company.com",
		"Find asset owner: engx crews asset://service/name",
	}

	for _, hint := range hints {
		// Color the command parts
		coloredHint := strings.ReplaceAll(hint, "engx crews",
			fmt.Sprintf("%sengx crews%s", r.colorScheme.AccentColor, r.colorScheme.ResetColor))

		output.WriteString(fmt.Sprintf("│ %s%s │\n",
			coloredHint,
			strings.Repeat(" ", width-len(hint)-4))) // Use original hint length for spacing
	}
}

func (r *SearchResultsRenderer) renderFooter(output *strings.Builder, width int) {
	output.WriteString(fmt.Sprintf("└%s┘\n", strings.Repeat("─", width-2)))
}

// Helper method for text wrapping
func (r *SearchResultsRenderer) wrapText(text string, maxWidth int) []string {
	if len(text) <= maxWidth {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	var currentLine string

	for _, word := range words {
		if len(currentLine)+len(word)+1 <= maxWidth {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}