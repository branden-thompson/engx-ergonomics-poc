package renderers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/design"
	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/models"
)

// AssetOwnershipRenderer handles rendering comprehensive asset ownership views
type AssetOwnershipRenderer struct {
	headerFooter *components.HeaderFooterComponent
	dataStore    DataStoreInterface // Interface for data access
	terminalWidth int
}

// DataStoreInterface provides access to crew and catalog data
type DataStoreInterface interface {
	GetCrew(id string) (*models.Crew, error)
	GetCatalogAsset(id string) (*models.CatalogAsset, error)
}

// NewAssetOwnershipRenderer creates a new asset ownership renderer
func NewAssetOwnershipRenderer(dataStore DataStoreInterface) *AssetOwnershipRenderer {
	return &AssetOwnershipRenderer{
		dataStore: dataStore,
	}
}

// Render generates the comprehensive asset ownership view
func (r *AssetOwnershipRenderer) Render(asset *models.CatalogAsset, width int) (string, error) {
	if width < 40 {
		width = 40 // Minimum width
	}

	r.terminalWidth = width
	r.headerFooter = components.NewHeaderFooterComponent(width)

	var output strings.Builder

	// Get owning crew information
	owningCrew, err := r.dataStore.GetCrew(asset.OwningCrewID)
	if err != nil {
		return "", fmt.Errorf("failed to get owning crew: %w", err)
	}

	// Render main header
	header := r.headerFooter.RenderHeader(
		asset.VanityName,
		string(asset.AssetType))
	output.WriteString(header)
	output.WriteString("\n")

	// Asset Owner section
	ownerCrewIDColor := r.getCrewTypeColor(owningCrew.Type)
	ownerIDDisplay := r.applyColor("(", design.ColorAttribute) + r.applyColor(asset.OwningCrewID, ownerCrewIDColor) + r.applyColor(")", design.ColorAttribute)
	output.WriteString(fmt.Sprintf("Asset Owner     : %s%50s%s\n",
		r.applyColor(owningCrew.VanityName, design.ColorLabel),
		"",
		ownerIDDisplay))
	output.WriteString("\n")

	// Current On-Call section
	onCallMembers := r.getOnCallMembers(owningCrew)
	output.WriteString("Current On-Call:\n\n")

	if len(onCallMembers) > 0 {
		err := r.renderOnCallTable(&output, onCallMembers, width)
		if err != nil {
			return "", fmt.Errorf("failed to render on-call table: %w", err)
		}
	} else {
		output.WriteString("   No current on-call assignments\n")
	}
	output.WriteString("\n")

	// Crews with Access section
	accessCrews, err := r.getCrewsWithAccess(asset)
	if err != nil {
		return "", fmt.Errorf("failed to get crews with access: %w", err)
	}

	output.WriteString(fmt.Sprintf("Crews with Access (%02d):\n", len(accessCrews)))
	output.WriteString("\n")

	err = r.renderAccessTable(&output, accessCrews, width)
	if err != nil {
		return "", fmt.Errorf("failed to render access table: %w", err)
	}
	output.WriteString("\n")

	// Asset Dependencies section
	output.WriteString(fmt.Sprintf("Asset Dependencies & Owning Crews (%02d):\n", len(asset.Dependencies)))
	output.WriteString("\n")

	err = r.renderDependenciesTable(&output, asset.Dependencies, width)
	if err != nil {
		return "", fmt.Errorf("failed to render dependencies table: %w", err)
	}

	output.WriteString("\n")

	// Render footer with asset ID and last updated
	assetIDText := r.applyColor("ASSET ID: ", design.ColorDarkGray) + asset.ID
	footer := r.headerFooter.RenderFooterWithColors(
		assetIDText,
		fmt.Sprintf("Last Updated: %s", asset.LastModified.Format("01/02/2006")),
		"", design.ColorDarkGray)
	output.WriteString(footer)

	return output.String(), nil
}

// renderOnCallTable renders the Current On-Call section
func (r *AssetOwnershipRenderer) renderOnCallTable(output *strings.Builder, members []OnCallInfo, width int) error {
	// Define columns for on-call table
	columnDefs := []components.ColumnDefinition{
		{Name: "prefix", Header: "", Width: 3, Alignment: "left", Color: ""},
		{Name: "number", Header: "##", Width: 3, Alignment: "left", Color: ""},
		{Name: "name", Header: "FULL NAME", Width: 0, MinWidth: 25, Alignment: "left", Color: ""},
		{Name: "ldap", Header: "LDAP", Width: 16, Alignment: "left", Color: ""},
		{Name: "level", Header: "LEVEL", Width: 8, Alignment: "left", Color: ""},
		{Name: "priority", Header: "ROTATION", Width: 12, Alignment: "left", Color: ""},
	}

	// Build table data
	var data [][]string
	for i, member := range members {
		// Get on-call color based on priority
		onCallColor := design.ColorSecondaryOnCall // Default
		if member.Priority == models.OnCallPrimary {
			onCallColor = design.ColorPrimaryOnCall
		}

		prefix := " " + r.applyColor("⏺", onCallColor) + " "
		rotation := "Secondary"
		if member.Priority == models.OnCallPrimary {
			rotation = "Primary"
		}

		row := []string{
			prefix,
			r.applyColor(fmt.Sprintf("%02d.", i+1), design.ColorRowNumber),
			r.applyColor(member.FullName, onCallColor), // Match on-call color
			r.applyColor("@"+member.UserID, design.ColorAttribute),
			r.applyColor(member.Level, design.ColorAttribute),
			r.applyColor(rotation, design.ColorAttribute),
		}
		data = append(data, row)
	}

	// Create and render table
	layout := &components.DataTableLayout{
		Columns: columnDefs,
		Data:    data,
	}

	dataTable := components.NewDataTableRowFromLayout(width, layout)

	// Render header
	headerRow := []string{
		"   ",
		r.applyColor("##", design.ColorTableHeader),
		r.applyColor("FULL NAME", design.ColorTableHeader),
		r.applyColor("LDAP", design.ColorTableHeader),
		r.applyColor("LEVEL", design.ColorTableHeader),
		r.applyColor("ROTATION", design.ColorTableHeader),
	}
	headerLine := dataTable.FormatDataRow(headerRow)
	output.WriteString(headerLine + "\n")

	// Render data rows
	for _, rowData := range data {
		rowLine := dataTable.FormatDataRow(rowData)
		output.WriteString(rowLine + "\n")
	}

	return nil
}

// renderAccessTable renders the Crews with Access section
func (r *AssetOwnershipRenderer) renderAccessTable(output *strings.Builder, accessInfo []AccessInfo, width int) error {
	// Define columns for access table
	columnDefs := []components.ColumnDefinition{
		{Name: "prefix", Header: "", Width: 3, Alignment: "left", Color: ""},
		{Name: "number", Header: "##", Width: 3, Alignment: "left", Color: ""},
		{Name: "name", Header: "CREW NAME", Width: 0, MinWidth: 30, Alignment: "left", Color: ""},
		{Name: "id", Header: "ID", Width: 12, Alignment: "left", Color: ""}, // Increased from 8 to 12 for proper spacing
		{Name: "access", Header: "ACCESS-LEVEL", Width: 15, Alignment: "left", Color: ""},
		{Name: "oncall", Header: "ON-CALL", Width: 12, Alignment: "left", Color: ""},
	}

	// Build table data
	var data [][]string
	for i, access := range accessInfo {
		onCallUser := access.OnCallUser
		if onCallUser == "" {
			onCallUser = "-"
		} else {
			onCallUser = "@" + onCallUser
		}

		row := []string{
			"   ", // Empty prefix for alignment (3 spaces)
			r.applyColor(fmt.Sprintf("%02d.", i+1), design.ColorRowNumber),
			r.applyColor(access.CrewName, design.ColorLabel),
			r.applyColor(access.CrewID, design.ColorAttribute), // Keep crew IDs grey in access table
			r.applyColor(string(access.AccessLevel), design.ColorAttribute),
			r.applyColor(onCallUser, design.ColorAttribute),
		}
		data = append(data, row)
	}

	// Create and render table
	layout := &components.DataTableLayout{
		Columns: columnDefs,
		Data:    data,
	}

	dataTable := components.NewDataTableRowFromLayout(width, layout)

	// Render header
	headerRow := []string{
		"   ",
		r.applyColor("##", design.ColorTableHeader),
		r.applyColor("CREW NAME", design.ColorTableHeader),
		r.applyColor("ID", design.ColorTableHeader),
		r.applyColor("ACCESS-LEVEL", design.ColorTableHeader),
		r.applyColor("ON-CALL", design.ColorTableHeader),
	}
	headerLine := dataTable.FormatDataRow(headerRow)
	output.WriteString(headerLine + "\n")

	// Render data rows
	for _, rowData := range data {
		rowLine := dataTable.FormatDataRow(rowData)
		output.WriteString(rowLine + "\n")
	}

	return nil
}

// renderDependenciesTable renders the Asset Dependencies section
func (r *AssetOwnershipRenderer) renderDependenciesTable(output *strings.Builder, dependencies []models.AssetDependency, width int) error {
	// Define columns for dependencies table
	columnDefs := []components.ColumnDefinition{
		{Name: "prefix", Header: "", Width: 3, Alignment: "left", Color: ""},
		{Name: "number", Header: "##", Width: 3, Alignment: "left", Color: ""},
		{Name: "name", Header: "DEPENDENCY NAME", Width: 0, MinWidth: 25, Alignment: "left", Color: ""},
		{Name: "version", Header: "VERSION", Width: 10, Alignment: "left", Color: ""},
		{Name: "id", Header: "ID", Width: 10, Alignment: "left", Color: ""},
		{Name: "owner", Header: "OWNER", Width: 12, Alignment: "left", Color: ""},
	}

	// Build table data
	var data [][]string
	for i, dep := range dependencies {
		healthColor := dep.Health.GetColor()
		prefix := " " + r.applyColor("⏺", healthColor) + " "

		row := []string{
			prefix,
			r.applyColor(fmt.Sprintf("%02d.", i+1), design.ColorRowNumber),
			r.applyColor(dep.Name, design.ColorLabel),
			r.applyColor(dep.Version, design.ColorAttribute),
			r.applyColor(dep.DependencyID, design.ColorAttribute),
			r.applyColor(dep.OwningCrewID, design.ColorAttribute),
		}
		data = append(data, row)
	}

	// Create and render table
	layout := &components.DataTableLayout{
		Columns: columnDefs,
		Data:    data,
	}

	dataTable := components.NewDataTableRowFromLayout(width, layout)

	// Render header
	headerRow := []string{
		"   ",
		r.applyColor("##", design.ColorTableHeader),
		r.applyColor("DEPENDENCY NAME", design.ColorTableHeader),
		r.applyColor("VERSION", design.ColorTableHeader),
		r.applyColor("ID", design.ColorTableHeader),
		r.applyColor("OWNER", design.ColorTableHeader),
	}
	headerLine := dataTable.FormatDataRow(headerRow)
	output.WriteString(headerLine + "\n")

	// Render data rows
	for _, rowData := range data {
		rowLine := dataTable.FormatDataRow(rowData)
		output.WriteString(rowLine + "\n")
	}

	return nil
}

// Helper types and methods

type OnCallInfo struct {
	UserID   string
	FullName string
	Level    string
	Priority models.OnCallType
}

type AccessInfo struct {
	CrewID      string
	CrewName    string
	AccessLevel models.AssetAccessLevel
	OnCallUser  string // Current on-call user for this crew
}

// getOnCallMembers extracts on-call information from a crew using rotation system
func (r *AssetOwnershipRenderer) getOnCallMembers(crew *models.Crew) []OnCallInfo {
	var onCallMembers []OnCallInfo
	currentRotations := crew.GetCurrentOnCallRotations()

	for _, rotation := range currentRotations {
		// Find the member details for this rotation
		for _, member := range crew.Members {
			if member.UserID == rotation.OnCallMember {
				priority := models.OnCallSecondary // Default
				if rotation.Name == "Primary" {
					priority = models.OnCallPrimary
				}

				onCallMembers = append(onCallMembers, OnCallInfo{
					UserID:   member.UserID,
					FullName: member.FullName,
					Level:    member.Level,
					Priority: priority,
				})
				break
			}
		}
	}

	// Sort by priority (Primary first)
	sort.Slice(onCallMembers, func(i, j int) bool {
		if onCallMembers[i].Priority != onCallMembers[j].Priority {
			return onCallMembers[i].Priority == models.OnCallPrimary
		}
		return onCallMembers[i].FullName < onCallMembers[j].FullName
	})

	return onCallMembers
}

// getCrewsWithAccess builds access information including on-call status
func (r *AssetOwnershipRenderer) getCrewsWithAccess(asset *models.CatalogAsset) ([]AccessInfo, error) {
	var accessInfo []AccessInfo

	for _, grant := range asset.AccessGrants {
		if grant.Status != models.GrantActive {
			continue
		}

		crew, err := r.dataStore.GetCrew(grant.CrewID)
		if err != nil {
			// Skip crews we can't find (might be placeholders)
			continue
		}

		// Find current on-call user for this crew using rotation system
		onCallUser := ""
		currentRotations := crew.GetCurrentOnCallRotations()

		// Look for Primary first
		for _, rotation := range currentRotations {
			if rotation.Name == "Primary" {
				onCallUser = rotation.OnCallMember
				break
			}
		}

		// If no primary, look for Secondary
		if onCallUser == "" {
			for _, rotation := range currentRotations {
				if rotation.Name == "Secondary" {
					onCallUser = rotation.OnCallMember
					break
				}
			}
		}

		// If still no one, take any available rotation
		if onCallUser == "" && len(currentRotations) > 0 {
			onCallUser = currentRotations[0].OnCallMember
		}

		accessInfo = append(accessInfo, AccessInfo{
			CrewID:      grant.CrewID,
			CrewName:    crew.VanityName,
			AccessLevel: grant.AccessLevel,
			OnCallUser:  onCallUser,
		})
	}

	// Sort by access level priority (Owner first, then by name)
	sort.Slice(accessInfo, func(i, j int) bool {
		if accessInfo[i].AccessLevel != accessInfo[j].AccessLevel {
			return r.getAccessPriority(accessInfo[i].AccessLevel) < r.getAccessPriority(accessInfo[j].AccessLevel)
		}
		return accessInfo[i].CrewName < accessInfo[j].CrewName
	})

	return accessInfo, nil
}

// getAccessPriority returns sort priority for access levels
func (r *AssetOwnershipRenderer) getAccessPriority(level models.AssetAccessLevel) int {
	switch level {
	case models.AccessLevelOwner:
		return 0
	case models.AccessLevelPublish:
		return 1
	case models.AccessLevelDeploy:
		return 2
	case models.AccessLevelContrib:
		return 3
	case models.AccessLevelReadOnly:
		return 4
	default:
		return 5
	}
}

// getCrewTypeColor returns the appropriate color code for crew IDs based on crew type
func (r *AssetOwnershipRenderer) getCrewTypeColor(crewType models.CrewType) string {
	switch crewType {
	case models.CrewTypeStandard:
		return design.ColorCrewStandard
	case models.CrewTypeVirtual:
		return design.ColorCrewVirtual
	default:
		return design.ColorAttribute // fallback to default
	}
}

// applyColor applies ANSI color codes with 256-color support
func (r *AssetOwnershipRenderer) applyColor(text, colorCode string) string {
	if colorCode == "" || text == "" {
		return text
	}
	// Handle 256-color codes (3+ digits) vs basic ANSI codes (1-2 digits)
	if len(colorCode) >= 3 {
		return fmt.Sprintf("\033[38;5;%sm%s\033[%sm", colorCode, text, design.ColorReset)
	}
	return fmt.Sprintf("\033[%sm%s\033[%sm", colorCode, text, design.ColorReset)
}