package renderers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/design"
	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/models"
)

// MembershipListRenderer handles rendering user membership views using shared components
type MembershipListRenderer struct {
	headerFooter *components.HeaderFooterComponent
	dataTable    *components.DataTableRow
	badgeAligner *components.BadgeAligner
	terminalWidth int
}

// NewMembershipListRenderer creates a new membership list renderer
func NewMembershipListRenderer() *MembershipListRenderer {
	return &MembershipListRenderer{}
}

// Render generates the membership list view using shared components
func (r *MembershipListRenderer) Render(userID string, crews []*models.Crew, width int) (string, error) {
	if width < 40 {
		width = 40 // Minimum width
	}

	r.terminalWidth = width
	r.headerFooter = components.NewHeaderFooterComponent(width)
	r.badgeAligner = components.NewBadgeAligner(width)

	var output strings.Builder
	normalizedUserID := normalizeUserID(userID)

	// Get user info for the header
	userInfo := r.getUserDisplayInfo(userID, crews, normalizedUserID)

	// Render header using shared component
	header := r.headerFooter.RenderHeader(
		fmt.Sprintf("Crew Memberships: %s (%s)", userInfo.FullName, userInfo.Username),
		"")
	output.WriteString(header)

	// Empty line
	output.WriteString("\n")

	// Table count
	uniqueCrews := r.getUniqueCrews(crews, normalizedUserID)
	output.WriteString(fmt.Sprintf("Crew Memberships (%d):\n\n", len(uniqueCrews)))

	// Define table columns EXACTLY like crews detail command
	columnDefs := []components.ColumnDefinition{
		{Name: "prefix", Header: "", Width: 3, Alignment: "left", Color: ""},
		{Name: "number", Header: "##", Width: 3, Alignment: "left", Color: ""},
		{Name: "name", Header: "CREW NAME", Width: 0, MinWidth: 20, Alignment: "left", Color: ""}, // Fill column
		{Name: "id", Header: "ID", Width: 16, Alignment: "left", Color: ""},
		{Name: "type", Header: "TYPE", Width: 12, Alignment: "left", Color: ""},
		{Name: "role", Header: "ROLE", Width: 8, Alignment: "left", Color: ""},
	}

	// Create table layout
	layout := &components.DataTableLayout{
		Columns: columnDefs,
		Data:    r.buildTableData(uniqueCrews, normalizedUserID),
	}

	// Create data table renderer
	r.dataTable = components.NewDataTableRowFromLayout(width, layout)

	// Reserve space for badges (like crews detail command)
	r.dataTable.SetBadge("(On-Call 1)", 4).SetBadgeRenderMode(false)

	// Render table header EXACTLY like crews detail command
	headerRow := []string{
		"   ", // prefix space (3 chars like data rows)
		r.applyColor("##", design.ColorTableHeader),
		r.applyColor("CREW NAME", design.ColorTableHeader),
		r.applyColor("ID", design.ColorTableHeader),
		r.applyColor("TYPE", design.ColorTableHeader),
		r.applyColor("ROLE", design.ColorTableHeader),
	}
	headerLine := r.dataTable.FormatDataRow(headerRow)
	output.WriteString(headerLine + "\n")

	// Render data rows with manual badge alignment (like crews detail command)
	for i, rowData := range layout.Data {
		rowLine := r.dataTable.FormatDataRow(rowData)

		// Add on-call badge if this member is on-call
		crew := uniqueCrews[i]
		member := crew.GetMember(normalizedUserID)
		if member != nil && member.IsCurrentlyOnCall() {
			var badge string
			var badgeColor string
			if member.OnCallType == models.OnCallPrimary {
				badge = "(On-Call 1)"
				badgeColor = design.ColorPrimaryOnCall
			} else {
				badge = "(On-Call 2)"
				badgeColor = design.ColorSecondaryOnCall
			}
			coloredBadge := r.applyColor(badge, badgeColor)
			rowLine = r.badgeAligner.AlignRight(rowLine, coloredBadge, 4)
		}

		output.WriteString(rowLine + "\n")
	}

	// Empty line
	output.WriteString("\n")

	// Render footer using shared component
	footer := r.headerFooter.RenderFooterWithColors(
		fmt.Sprintf("%s - %s", userInfo.Level, userInfo.Title),
		fmt.Sprintf("Last Updated: %s", "12/01/2025"),
		"", design.ColorDarkGray)
	output.WriteString(footer)

	return output.String(), nil
}

// Helper methods for the new shared component implementation

type UserDisplayInfo struct {
	FullName string
	Username string
	Level    string
	Title    string
}

func (r *MembershipListRenderer) getUserDisplayInfo(userID string, crews []*models.Crew, normalizedUserID string) UserDisplayInfo {
	// For bthompso, return Branden Thompson's info
	if normalizedUserID == "bthompso" {
		return UserDisplayInfo{
			FullName: "Branden Thompson",
			Username: "@bthompso",
			Level:    "IC5",
			Title:    "Principal Design Engineer",
		}
	}

	// Default user info
	return UserDisplayInfo{
		FullName: strings.Title(normalizedUserID),
		Username: fmt.Sprintf("@%s", normalizedUserID),
		Level:    "IC3",
		Title:    "Software Engineer",
	}
}

func (r *MembershipListRenderer) getUniqueCrews(crews []*models.Crew, normalizedUserID string) []*models.Crew {
	seen := make(map[string]bool)
	var unique []*models.Crew

	// Sort crews by role priority and name
	sort.Slice(crews, func(i, j int) bool {
		memberI := crews[i].GetMember(normalizedUserID)
		memberJ := crews[j].GetMember(normalizedUserID)

		if memberI != nil && memberJ != nil {
			// Sort by role priority first
			if memberI.Role != memberJ.Role {
				return r.getRolePriority(memberI.Role) < r.getRolePriority(memberJ.Role)
			}
		}

		// Then by crew name
		return crews[i].VanityName < crews[j].VanityName
	})

	for _, crew := range crews {
		if !seen[crew.ID] && crew.GetMember(normalizedUserID) != nil {
			seen[crew.ID] = true
			unique = append(unique, crew)
		}
	}

	return unique
}

func (r *MembershipListRenderer) buildTableData(crews []*models.Crew, normalizedUserID string) [][]string {
	var data [][]string

	for i, crew := range crews {
		member := crew.GetMember(normalizedUserID)
		if member == nil {
			continue
		}

		// Use data model's on-call information
		isOnCall := member.IsCurrentlyOnCall()
		var onCallColor string
		if isOnCall {
			// Determine color based on OnCallType from data model
			if member.OnCallType == models.OnCallPrimary {
				onCallColor = design.ColorPrimaryOnCall // Primary on-call (deep red)
			} else {
				onCallColor = design.ColorSecondaryOnCall // Secondary/backup on-call (orange-red)
			}
		}

		// Format prefix and row number EXACTLY like crews detail command
		prefix := "   " // 3 spaces for alignment
		if isOnCall {
			// Apply semantic color to the on-call indicator
			prefix = " " + r.applyColor("⏺", onCallColor) + " "
		}

		// Crew name color - use on-call color if on-call, otherwise white
		nameColor := design.ColorLabel // Default white
		if isOnCall && onCallColor != "" {
			nameColor = onCallColor // Use matching on-call color
		}

		// Crew type (default to STANDARD for now)
		crewType := "STANDARD"
		if strings.Contains(strings.ToLower(crew.VanityName), "virtual") ||
		   strings.Contains(strings.ToLower(crew.VanityName), "universal") {
			crewType = "VIRTUAL"
		}

		row := []string{
			prefix,
			r.applyColor(fmt.Sprintf("%02d.", i+1), design.ColorRowNumber),          // Row numbers - DARK GREY
			r.applyColor(crew.VanityName, nameColor),                                 // Name - WHITE or ON-CALL color
			r.applyColor(crew.ID, design.ColorAttribute),                             // ID - DARK GREY
			r.applyColor(crewType, design.ColorAttribute),                            // Type - DARK GREY
			r.applyColor(strings.Title(string(member.Role)), design.ColorAttribute), // Role - DARK GREY
		}

		data = append(data, row)
	}

	return data
}

// isOnCall and getOnCallBadge methods removed - now using data model's Member.IsCurrentlyOnCall() and OnCallType

func (r *MembershipListRenderer) applyColor(text, colorCode string) string {
	if colorCode == "" || text == "" {
		return text
	}
	// Handle 256-color codes (3+ digits) vs basic ANSI codes (1-2 digits)
	if len(colorCode) >= 3 {
		return fmt.Sprintf("\033[38;5;%sm%s\033[%sm", colorCode, text, design.ColorReset)
	}
	return fmt.Sprintf("\033[%sm%s\033[%sm", colorCode, text, design.ColorReset)
}

func (r *MembershipListRenderer) alignBadge(content, badge string, gutter int) string {
	// Use similar logic to the crews detail view
	contentDisplayWidth := r.getDisplayWidth(content)
	badgeDisplayWidth := r.getDisplayWidth(badge)

	totalUsedWidth := contentDisplayWidth + badgeDisplayWidth
	paddingNeeded := r.terminalWidth - totalUsedWidth

	if paddingNeeded < gutter {
		paddingNeeded = gutter
	}

	return content + strings.Repeat(" ", paddingNeeded) + badge
}

func (r *MembershipListRenderer) getDisplayWidth(text string) int {
	// Strip ANSI codes for width calculation (same as other components)
	stripped := text
	// Simple regex to remove ANSI codes
	for {
		start := strings.Index(stripped, "\033[")
		if start == -1 {
			break
		}
		end := strings.Index(stripped[start:], "m")
		if end == -1 {
			break
		}
		stripped = stripped[:start] + stripped[start+end+1:]
	}
	return len(stripped)
}

func (r *MembershipListRenderer) getRolePriority(role models.MemberRole) int {
	switch role {
	case models.RoleOwner:
		return 0
	case models.RoleAdmin:
		return 1
	case models.RoleMember:
		return 2
	case models.RoleTemp:
		return 3
	default:
		return 4
	}
}

// normalizeUserID handles email addresses by extracting username
func normalizeUserID(userID string) string {
	if strings.Contains(userID, "@") {
		return strings.Split(userID, "@")[0]
	}
	return userID
}