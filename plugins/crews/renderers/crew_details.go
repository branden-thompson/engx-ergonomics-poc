package renderers

import (
	"fmt"
	"os"
	"os/user"
	"sort"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/design"
	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/models"
)

// CrewDetailsRenderer handles rendering crew detail views using modular components
type CrewDetailsRenderer struct {
	colorScheme       ColorScheme
	headerFooter      *components.HeaderFooterComponent
	dataTableRow      *components.DataTableRow
	badgeAligner      *components.BadgeAligner
	width             int
}

// NewCrewDetailsRenderer creates a new crew details renderer
func NewCrewDetailsRenderer() *CrewDetailsRenderer {
	return &CrewDetailsRenderer{
		colorScheme: DefaultColorScheme(),
	}
}

// Render generates the crew details view with the exact template format
func (r *CrewDetailsRenderer) Render(crew *models.Crew, assets []*models.AssetOwnership, width int) (string, error) {
	if width < 80 {
		width = 80 // Minimum width for proper formatting
	}

	r.width = width
	r.headerFooter = components.NewHeaderFooterComponent(width)
	r.badgeAligner = components.NewBadgeAligner(width)

	var output strings.Builder

	// Header section using modular component
	r.renderHeader(&output, crew)

	// Members section using modular table component
	r.renderMembers(&output, crew)

	// Assets and description section
	r.renderAssetsAndDescription(&output, crew, assets)

	// Footer section using modular component
	r.renderFooter(&output, crew)

	return output.String(), nil
}

// createMembersDataTableLayout creates the table layout for crew members
func (r *CrewDetailsRenderer) createMembersDataTableLayout(width int, members []models.Member) *components.DataTableLayout {
	layout := &components.DataTableLayout{
		Columns: []components.ColumnDefinition{
			{Name: "prefix", Header: "", Width: 3, Alignment: "left", Color: ""},
			{Name: "number", Header: "##", Width: 3, Alignment: "left", Color: ""},
			{Name: "fullName", Header: "FULL NAME", Width: 0, MinWidth: 20, Alignment: "left", Color: ""}, // Fill column
			{Name: "ldap", Header: "LDAP", Width: 16, Alignment: "left", Color: ""},
			{Name: "level", Header: "LEVEL", Width: 8, Alignment: "left", Color: ""},
			{Name: "role", Header: "ROLE", Width: 8, Alignment: "left", Color: ""},
		},
		Data: [][]string{},
	}

	// Convert members to data rows with proper colors
	for i, member := range members {
		prefix := "   " // 3 spaces for alignment
		var onCallColor string

		if member.IsOnCall {
			// Assign specific colors based on member
			if member.UserID == "hbacot" { // Hunter Bacot
				onCallColor = design.ColorPrimaryOnCall // #D62 (deep red)
			} else if member.UserID == "osurtiz" { // Olivier Surtiz
				onCallColor = design.ColorSecondaryOnCall // #DA2 (orange-red)
			} else {
				onCallColor = design.ColorPrimaryOnCall // Default for other on-call members
			}

			// Apply semantic color to the on-call indicator
			prefix = " " + r.applyColor("⏺", onCallColor) + " "
		}

		// Determine full name color - use on-call color if on-call, otherwise white
		nameColor := design.ColorLabel // Default white
		if member.IsOnCall && onCallColor != "" {
			nameColor = onCallColor // Use matching on-call color
		}

		row := []string{
			prefix,
			r.applyColor(fmt.Sprintf("%02d.", i+1), design.ColorRowNumber),          // Row numbers - DARK GREY
			r.applyColor(member.FullName, nameColor),                                 // Full name - WHITE or ON-CALL color
			r.applyColor(fmt.Sprintf("@%s", member.UserID), design.ColorAttribute),  // LDAP - DARK GREY
			r.applyColor(member.Level, design.ColorAttribute),                        // Level - DARK GREY
			r.applyColor(strings.Title(string(member.Role)), design.ColorAttribute), // Role - DARK GREY
		}
		layout.Data = append(layout.Data, row)
	}

	return layout
}

// renderHeader creates the header line with crew name and ID
func (r *CrewDetailsRenderer) renderHeader(output *strings.Builder, crew *models.Crew) {
	// Apply crew type color to the crew ID while keeping "CREW ID:" white
	crewTypeColor := r.getCrewTypeColor("STANDARD CREW") // For now, assume all crews are standard
	rightLabel := fmt.Sprintf("CREW ID: %s", r.applyColor(crew.ID, crewTypeColor))

	header := r.headerFooter.RenderHeader(crew.VanityName, rightLabel)
	output.WriteString(header)
	output.WriteString("\n")
}

// renderMembers creates the members table using the new DataTableLayout system
func (r *CrewDetailsRenderer) renderMembers(output *strings.Builder, crew *models.Crew) {
	members := crew.GetActiveMembers()
	output.WriteString(fmt.Sprintf("Crew Members (%d):\n\n", len(members)))

	// Sort members by role priority (owner, admin, member, auto, temp)
	sort.Slice(members, func(i, j int) bool {
		return r.getRolePriority(members[i].Role) < r.getRolePriority(members[j].Role)
	})

	// Create table layout with sorted members
	layout := r.createMembersDataTableLayout(r.width, members)

	// Create data table renderer from layout
	dataTableRow := components.NewDataTableRowFromLayout(r.width, layout)

	// Reserve space for badges by setting badge space on the renderer (use longer badge for space calculation)
	dataTableRow.SetBadge("(On-Call 1)", 4).SetBadgeRenderMode(false)

	// Create header row using the same structure as data rows with colors
	headerRowData := []string{
		"   ", // prefix space (3 chars like data rows)
		r.applyColor("##", design.ColorTableHeader),      // number column - LIGHT PURPLE
		r.applyColor("FULL NAME", design.ColorTableHeader), // full name column - LIGHT PURPLE
		r.applyColor("LDAP", design.ColorTableHeader),       // ldap column - LIGHT PURPLE
		r.applyColor("LEVEL", design.ColorTableHeader),      // level column - LIGHT PURPLE
		r.applyColor("ROLE", design.ColorTableHeader),       // role column - LIGHT PURPLE
	}

	// Render header using the exact same logic as data rows
	headerLine := dataTableRow.FormatDataRow(headerRowData)
	output.WriteString(headerLine + "\n")

	// Render each data row
	for i, rowData := range layout.Data {
		rowLine := dataTableRow.FormatDataRow(rowData)

		// Add badge only for on-call members using BadgeAligner with matching color
		if strings.Contains(rowData[0], "⏺") { // Check if row has on-call indicator
			// Get the corresponding member to determine badge color and priority
			member := members[i]
			var badgeColor string
			var badgeText string
			if member.UserID == "hbacot" { // Hunter Bacot
				badgeColor = design.ColorPrimaryOnCall // #D62 (deep red)
				badgeText = "(On-Call 1)" // Primary
			} else if member.UserID == "osurtiz" { // Olivier Surtiz
				badgeColor = design.ColorSecondaryOnCall // #DA2 (orange-red)
				badgeText = "(On-Call 2)" // Secondary
			} else {
				badgeColor = design.ColorPrimaryOnCall // Default for other on-call members
				badgeText = "(On-Call 1)" // Default to primary
			}

			// Apply matching color to badge as the icon
			coloredBadge := r.applyColor(badgeText, badgeColor)
			rowLine = r.badgeAligner.AlignRight(rowLine, coloredBadge, 4)
		}

		output.WriteString(rowLine + "\n")
	}

	output.WriteString("\n")
}

// memberToDataMap converts a member to data map for the modular data table
func (r *CrewDetailsRenderer) memberToDataMap(member *models.Member, index int) map[string]string {
	return map[string]string{
		"number":   fmt.Sprintf("%02d.", index),
		"fullName": member.FullName,
		"ldap":     fmt.Sprintf("@%s", member.UserID),
		"level":    member.Level,
		"role":     strings.Title(string(member.Role)),
	}
}

// renderAssetsAndDescription creates the assets and description section
func (r *CrewDetailsRenderer) renderAssetsAndDescription(output *strings.Builder, crew *models.Crew, assets []*models.AssetOwnership) {
	output.WriteString(fmt.Sprintf("Assets Owned (%d)\n\n", len(assets)))

	// Get crew owner info
	owner := crew.GetOwner()
	managerName := "Not assigned"
	if owner != nil {
		managerName = fmt.Sprintf("%s (@%s)", owner.FullName, owner.UserID)
	}

	// Calculate available width for description (left column)
	managerText := fmt.Sprintf("Crew Manager of Record:\n%s", managerName)
	managerWidth := r.getMaxLineLength(managerText)

	// Reserve space for manager column and some padding
	descWidth := r.width - managerWidth - 8 // 8 chars for comfortable spacing
	if descWidth < 30 {
		descWidth = 30 // Minimum width for readable text
	}

	// Wrap description text to fit in left column
	descLines := r.wrapText(crew.Description, descWidth)

	// Print headers with proper spacing
	headerSpacing := r.width - 12 - 23 // 12=len("Description:"), 23=len("Crew Manager of Record:")
	if headerSpacing < 4 {
		headerSpacing = 4
	}
	output.WriteString(fmt.Sprintf("Description:%s%s\n",
		strings.Repeat(" ", headerSpacing),
		"Crew Manager of Record:"))

	// Print content lines
	maxLines := len(descLines)
	if maxLines < 1 {
		maxLines = 1 // At least one line for manager
	}

	for i := 0; i < maxLines; i++ {
		descLine := ""
		managerLine := ""

		if i < len(descLines) {
			// Apply description color (dark grey) to description text
			descLine = r.applyColor(descLines[i], design.ColorDescription)
		}
		if i == 0 { // Manager name only on first content line
			// Apply username-aware coloring to manager name
			if owner != nil {
				// Color the name part as label (white) and username part with appropriate color
				nameColor := design.ColorLabel
				usernameColor := r.getUsernameColor(owner.UserID)
				managerLine = fmt.Sprintf("%s (%s)",
					r.applyColor(owner.FullName, nameColor),
					r.applyColor(fmt.Sprintf("@%s", owner.UserID), usernameColor))
			} else {
				// "Not assigned" case
				managerLine = r.applyColor(managerName, design.ColorLabel)
			}
		}

		// Right-align manager text to terminal edge
		// Note: Use uncolored text length for width calculations
		plainDescLine := ""
		plainManagerLine := ""
		if i < len(descLines) {
			plainDescLine = descLines[i]
		}
		if i == 0 {
			plainManagerLine = managerName
		}

		contentWidth := len(plainDescLine) + len(plainManagerLine)
		padding := r.width - contentWidth
		if padding < 4 {
			padding = 4 // minimum spacing
		}

		output.WriteString(fmt.Sprintf("%s%s%s\n", descLine, strings.Repeat(" ", padding), managerLine))
	}

	output.WriteString("\n")
}

// renderFooter creates the footer line with crew type and last updated
func (r *CrewDetailsRenderer) renderFooter(output *strings.Builder, crew *models.Crew) {
	// Apply crew type color to the crew type text
	crewType := "STANDARD CREW" // For now, assume all crews are standard
	crewTypeColor := r.getCrewTypeColor(crewType)
	coloredCrewType := r.applyColor(crewType, crewTypeColor)

	footer := r.headerFooter.RenderFooterWithColors(coloredCrewType, fmt.Sprintf("Last Updated: %s", crew.LastModified.Format("01/02/2006")), "", design.ColorDarkGray)
	output.WriteString(footer)
}

// Helper methods

func (r *CrewDetailsRenderer) getRolePriority(role models.MemberRole) int {
	switch role {
	case models.RoleOwner:
		return 0
	case models.RoleAdmin:
		return 1
	case models.RoleMember:
		return 2
	case models.RoleAuto:
		return 3
	case models.RoleTemp:
		return 4
	default:
		return 5
	}
}

// wrapText wraps text to the specified width, breaking on word boundaries
func (r *CrewDetailsRenderer) wrapText(text string, width int) []string {
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
		// If adding this word would exceed the width, start a new line
		if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > width {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}

		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}

	// Add the last line if it has content
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// getMaxLineLength returns the maximum length of any line in a multi-line string
func (r *CrewDetailsRenderer) getMaxLineLength(text string) int {
	lines := strings.Split(text, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}
	return maxLen
}

// applyColor wraps text with ANSI color codes (handles both basic and 256-color codes)
func (r *CrewDetailsRenderer) applyColor(text, colorCode string) string {
	if colorCode == "" || text == "" {
		return text
	}

	// Check if it's a 256-color code (3-digit numbers like "196", "202")
	if len(colorCode) == 3 && colorCode[0] >= '1' && colorCode[0] <= '9' {
		// Use 256-color format for extended colors
		return fmt.Sprintf("\033[38;5;%sm%s\033[%sm", colorCode, text, design.ColorReset)
	}

	// Use basic ANSI format for standard colors (90, 91, 92, etc.)
	return fmt.Sprintf("\033[%sm%s\033[%sm", colorCode, text, design.ColorReset)
}

// getCurrentUsername attempts to get the current user's username
func (r *CrewDetailsRenderer) getCurrentUsername() string {
	// Try environment variable first (common in CI/containers)
	if envUser := os.Getenv("USER"); envUser != "" {
		return envUser
	}

	// Try os/user package
	if currentUser, err := user.Current(); err == nil && currentUser.Username != "" {
		return currentUser.Username
	}

	// Fallback to empty string (will default to "other user" color)
	return ""
}

// getUsernameColor returns the appropriate color for a username
func (r *CrewDetailsRenderer) getUsernameColor(username string) string {
	currentUser := r.getCurrentUsername()
	if currentUser != "" && username == currentUser {
		return design.ColorCurrentUser // Bright Green
	}
	return design.ColorOtherUser // Light Blue/Cyan
}

// getRoleColor returns the appropriate color for a member role
func (r *CrewDetailsRenderer) getRoleColor(role models.MemberRole) string {
	switch role {
	case models.RoleOwner:
		return design.ColorRoleOwner // Bright Magenta (#A4C)
	case models.RoleAdmin:
		return design.ColorRoleAdmin // Bright Yellow (#FA0)
	case models.RoleMember:
		return design.ColorRoleMember // Bright Green (#3AD)
	case models.RoleAuto:
		return design.ColorRoleAuto // Bright Blue (#69A)
	case models.RoleTemp:
		return design.ColorRoleTemp // Bright Red (#C66)
	default:
		return design.ColorAttribute // Default to dark grey
	}
}

// getLevelColor returns the appropriate color for a member level
func (r *CrewDetailsRenderer) getLevelColor(level string) string {
	switch level {
	// IC (Individual Contributor) levels
	case "IC1":
		return design.ColorIC1 // Light grey (#999)
	case "IC2":
		return design.ColorIC2 // Light blue (#9CF)
	case "IC3":
		return design.ColorIC3 // Bright blue (#59D)
	case "IC4":
		return design.ColorIC4 // Blue (#38E)
	case "IC5":
		return design.ColorIC5 // Cyan (#3BE)
	case "IC6":
		return design.ColorIC6 // Bright cyan (#0BF)
	case "IC7":
		return design.ColorIC7 // Cyan-green (#5DD)
	case "IC8":
		return design.ColorIC8 // Green-cyan (#2DC)
	case "IC9":
		return design.ColorIC9 // Bright cyan (#0FD)
	// MR (Management) levels
	case "MR2":
		return design.ColorMR2 // Brown/tan (#8B7)
	case "MR3":
		return design.ColorMR3 // Olive (#6B3)
	case "MR4":
		return design.ColorMR4 // Dark yellow (#5C1)
	case "MR5":
		return design.ColorMR5 // Green (#3D3)
	case "MR6":
		return design.ColorMR6 // Teal (#3D8)
	case "MR7":
		return design.ColorMR7 // Dark cyan (#1E8)
	case "MR8":
		return design.ColorMR8 // Cyan (#0F8)
	case "MR9":
		return design.ColorMR9 // Bright cyan (#0FD)
	default:
		return design.ColorAttribute // Default to dark grey for unknown levels
	}
}

// getCrewTypeColor returns the appropriate color for a crew type
func (r *CrewDetailsRenderer) getCrewTypeColor(crewType string) string {
	switch strings.ToUpper(crewType) {
	case "STANDARD CREW", "STANDARD":
		return design.ColorStandardCrew // Blue (#08F)
	case "VIRTUAL CREW", "VIRTUAL":
		return design.ColorVirtualCrew // Light purple (#D0F)
	default:
		return design.ColorStandardCrew // Default to standard crew color
	}
}


