package renderers

import (
	"fmt"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/models"
)

// AssetOwnerRenderer handles rendering asset ownership information
type AssetOwnerRenderer struct {
	colorScheme ColorScheme
}

// NewAssetOwnerRenderer creates a new asset owner renderer
func NewAssetOwnerRenderer() *AssetOwnerRenderer {
	return &AssetOwnerRenderer{
		colorScheme: DefaultColorScheme(),
	}
}

// Render generates the asset ownership view
func (r *AssetOwnerRenderer) Render(asset *models.AssetOwnership, crew *models.Crew, width int) (string, error) {
	if width < 40 {
		width = 40 // Minimum width
	}

	var output strings.Builder

	// Header with asset information
	r.renderAssetHeader(&output, asset, width)

	// Asset details
	r.renderAssetDetails(&output, asset, width)

	// Owning crew information
	r.renderOwnerCrew(&output, crew, width)

	// Contact information (on-call members)
	r.renderContactInfo(&output, crew, width)

	// Footer
	r.renderFooter(&output, width)

	return output.String(), nil
}

func (r *AssetOwnerRenderer) renderAssetHeader(output *strings.Builder, asset *models.AssetOwnership, width int) {
	header := fmt.Sprintf(" %s ASSET OWNERSHIP ", r.colorScheme.AssetsIcon)

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

func (r *AssetOwnerRenderer) renderAssetDetails(output *strings.Builder, asset *models.AssetOwnership, width int) {
	// Asset name
	assetNameLine := fmt.Sprintf("Asset: %s%s%s",
		r.colorScheme.AccentColor,
		asset.AssetName,
		r.colorScheme.ResetColor)

	output.WriteString(fmt.Sprintf("│ %s%s │\n",
		assetNameLine,
		strings.Repeat(" ", width-len(assetNameLine)+len(r.colorScheme.AccentColor)+len(r.colorScheme.ResetColor)-4)))

	// Asset URN
	urnLine := fmt.Sprintf("URN: %s", asset.AssetURN)
	if len(urnLine) > width-4 {
		// Wrap long URNs
		urnLines := r.wrapText(urnLine, width-4)
		for _, line := range urnLines {
			output.WriteString(fmt.Sprintf("│ %s%s │\n",
				line,
				strings.Repeat(" ", width-len(line)-4)))
		}
	} else {
		output.WriteString(fmt.Sprintf("│ %s%s │\n",
			urnLine,
			strings.Repeat(" ", width-len(urnLine)-4)))
	}

	// Asset type and status
	metaLine := fmt.Sprintf("Type: %s | Status: %s | Modified: %s",
		asset.AssetType,
		asset.Status,
		asset.LastModified.Format("2006-01-02"))

	output.WriteString(fmt.Sprintf("│ %s%s │\n",
		metaLine,
		strings.Repeat(" ", width-len(metaLine)-4)))
}

func (r *AssetOwnerRenderer) renderOwnerCrew(output *strings.Builder, crew *models.Crew, width int) {
	// Ownership section header
	output.WriteString(fmt.Sprintf("├─%s OWNED BY %s─┤\n",
		r.colorScheme.MembershipIcon,
		strings.Repeat("─", width-12)))

	// Crew identification
	crewLine := fmt.Sprintf("Crew: %s%s%s (%s)",
		r.colorScheme.CrewIDColor,
		crew.VanityName,
		r.colorScheme.ResetColor,
		crew.ID)

	output.WriteString(fmt.Sprintf("│ %s%s │\n",
		crewLine,
		strings.Repeat(" ", width-len(crewLine)+len(r.colorScheme.CrewIDColor)+len(r.colorScheme.ResetColor)-4)))

	// Crew description (truncated)
	if crew.Description != "" {
		descLines := r.wrapText(crew.Description, width-4)
		// Show only first line if multiple lines
		if len(descLines) > 0 {
			descLine := descLines[0]
			if len(descLines) > 1 {
				descLine += "..."
			}
			output.WriteString(fmt.Sprintf("│ %s%s │\n",
				descLine,
				strings.Repeat(" ", width-len(descLine)-4)))
		}
	}

	// Owner and admin information
	owner := crew.GetOwner()
	if owner != nil {
		ownerLine := fmt.Sprintf("Owner: %s%s%s (%s)",
			r.colorScheme.OwnerColor,
			owner.FullName,
			r.colorScheme.ResetColor,
			owner.Email)

		output.WriteString(fmt.Sprintf("│ %s%s │\n",
			ownerLine,
			strings.Repeat(" ", width-len(ownerLine)+len(r.colorScheme.OwnerColor)+len(r.colorScheme.ResetColor)-4)))
	}

	// Show member count
	memberCount := len(crew.GetActiveMembers())
	memberLine := fmt.Sprintf("Members: %d active", memberCount)
	output.WriteString(fmt.Sprintf("│ %s%s │\n",
		memberLine,
		strings.Repeat(" ", width-len(memberLine)-4)))
}

func (r *AssetOwnerRenderer) renderContactInfo(output *strings.Builder, crew *models.Crew, width int) {
	// Contact section header
	output.WriteString(fmt.Sprintf("├─%s CONTACT INFORMATION %s─┤\n",
		r.colorScheme.OnCallIcon,
		strings.Repeat("─", width-22)))

	// On-call information
	if crew.OnCallSchedule.Enabled {
		if len(crew.OnCallSchedule.CurrentOnCall) > 0 {
			onCallLine := fmt.Sprintf("Current On-Call: %s%s%s",
				r.colorScheme.OnCallActiveColor,
				strings.Join(crew.OnCallSchedule.CurrentOnCall, ", "),
				r.colorScheme.ResetColor)

			output.WriteString(fmt.Sprintf("│ %s%s │\n",
				onCallLine,
				strings.Repeat(" ", width-len(onCallLine)+len(r.colorScheme.OnCallActiveColor)+len(r.colorScheme.ResetColor)-4)))
		} else {
			noOnCallLine := "Current On-Call: None"
			output.WriteString(fmt.Sprintf("│ %s%s │\n",
				noOnCallLine,
				strings.Repeat(" ", width-len(noOnCallLine)-4)))
		}

		// Escalation path
		if len(crew.OnCallSchedule.EscalationPath) > 0 {
			escalationLine := fmt.Sprintf("Escalation: %s",
				strings.Join(crew.OnCallSchedule.EscalationPath, " → "))

			if len(escalationLine) > width-4 {
				escalationLine = escalationLine[:width-7] + "..."
			}

			output.WriteString(fmt.Sprintf("│ %s%s │\n",
				escalationLine,
				strings.Repeat(" ", width-len(escalationLine)-4)))
		}
	} else {
		noOnCallLine := "On-call rotation: Disabled"
		output.WriteString(fmt.Sprintf("│ %s%s │\n",
			noOnCallLine,
			strings.Repeat(" ", width-len(noOnCallLine)-4)))

		// Show owner as primary contact
		owner := crew.GetOwner()
		if owner != nil {
			contactLine := fmt.Sprintf("Primary Contact: %s", owner.Email)
			output.WriteString(fmt.Sprintf("│ %s%s │\n",
				contactLine,
				strings.Repeat(" ", width-len(contactLine)-4)))
		}
	}
}

func (r *AssetOwnerRenderer) renderFooter(output *strings.Builder, width int) {
	output.WriteString(fmt.Sprintf("└%s┘\n", strings.Repeat("─", width-2)))
}

// Helper method for text wrapping
func (r *AssetOwnerRenderer) wrapText(text string, maxWidth int) []string {
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