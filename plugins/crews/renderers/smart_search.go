package renderers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/design"
	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/data"
)

// SmartSearchRenderer handles rendering the interstitial search results view
type SmartSearchRenderer struct {
	headerFooter *components.HeaderFooterComponent
	dataStore    DataStoreInterface // Interface for data access
	terminalWidth int
}

// SearchResult represents a unified search result that can be crew, asset, or user
type SearchResult struct {
	Type        string // "crew", "asset", "user"
	ID          string
	DisplayName string
	SubInfo1    string // For crew: ID, for asset: VERSION, for user: LDAP
	SubInfo2    string // For crew: MEMBERS, for asset: ID, for user: LEVEL
	SubInfo3    string // For crew: ASSETS, for asset: OWNER, for user: MEMBERSHIPS
	OriginalRef interface{} // Reference to original object
}

// SmartSearchResults holds categorized search results
type SmartSearchResults struct {
	CrewMatches  []SearchResult
	AssetMatches []SearchResult
	UserMatches  []SearchResult
	SearchTerm   string
	ExecutionTime string
}

// NewSmartSearchRenderer creates a new smart search renderer
func NewSmartSearchRenderer(dataStore DataStoreInterface) *SmartSearchRenderer {
	return &SmartSearchRenderer{
		dataStore: dataStore,
	}
}

// Search performs smart search across crews, assets, and users
func (r *SmartSearchRenderer) Search(searchTerm string, executionTime string) (*SmartSearchResults, error) {
	results := &SmartSearchResults{
		SearchTerm:    searchTerm,
		ExecutionTime: executionTime,
	}

	// Search crews
	crews, err := r.searchCrews(searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search crews: %w", err)
	}
	results.CrewMatches = crews

	// Search assets
	assets, err := r.searchAssets(searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search assets: %w", err)
	}
	results.AssetMatches = assets

	// Search users
	users, err := r.searchUsers(searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	results.UserMatches = users

	return results, nil
}

// Render generates the smart search results view
func (r *SmartSearchRenderer) Render(results *SmartSearchResults, width int) (string, error) {
	if width < 40 {
		width = 40 // Minimum width
	}

	r.terminalWidth = width
	r.headerFooter = components.NewHeaderFooterComponent(width)

	var output strings.Builder

	totalResults := len(results.CrewMatches) + len(results.AssetMatches) + len(results.UserMatches)

	// Render main header
	header := r.headerFooter.RenderHeader(
		"EngX CLI Crews SmartSearch",
		fmt.Sprintf("%05d Total", totalResults))
	output.WriteString(header)
	output.WriteString("\n")

	itemNumber := 1

	// Render crew matches
	if len(results.CrewMatches) > 0 {
		output.WriteString(fmt.Sprintf("Potential Crew Matches (%02d):\n\n", len(results.CrewMatches)))

		err := r.renderCrewMatches(&output, results.CrewMatches, &itemNumber, width)
		if err != nil {
			return "", fmt.Errorf("failed to render crew matches: %w", err)
		}
		output.WriteString("\n")
	}

	// Render asset matches
	if len(results.AssetMatches) > 0 {
		output.WriteString(fmt.Sprintf("Potential Asset Matches (%02d):\n\n", len(results.AssetMatches)))

		err := r.renderAssetMatches(&output, results.AssetMatches, &itemNumber, width)
		if err != nil {
			return "", fmt.Errorf("failed to render asset matches: %w", err)
		}
		output.WriteString("\n")
	}

	// Render user matches
	if len(results.UserMatches) > 0 {
		output.WriteString(fmt.Sprintf("Potential User Matches (%02d):\n\n", len(results.UserMatches)))

		err := r.renderUserMatches(&output, results.UserMatches, &itemNumber, width)
		if err != nil {
			return "", fmt.Errorf("failed to render user matches: %w", err)
		}
		output.WriteString("\n")
	}

	// Render footer with search execution time
	footer := r.headerFooter.RenderFooterWithColors(
		fmt.Sprintf("Search Executed in: %s", results.ExecutionTime),
		"",
		"", design.ColorDarkGray)
	output.WriteString(footer)
	output.WriteString("\n")

	// Render selection prompt
	output.WriteString(r.applyColor(" ! If none of these are correct, you try a different query below.\n", design.ColorAttribute))
	rangeText := ""
	if totalResults > 0 {
		rangeText = fmt.Sprintf(" (1-%d) or ", totalResults)
	}
	output.WriteString(r.applyColor(fmt.Sprintf(" ? Which one would you like details for%s<SearchString>: ", rangeText), design.ColorAttribute))

	return output.String(), nil
}

// renderCrewMatches renders the crew matches table
func (r *SmartSearchRenderer) renderCrewMatches(output *strings.Builder, matches []SearchResult, itemNumber *int, width int) error {
	// Define columns for crew matches table using enhanced gutter system
	layout := components.NewEnhancedTableLayout(4, []components.ColumnDefinition{
		{Name: "prefix", Header: "", Width: 3, Alignment: "left", Color: ""},
		{Name: "number", Header: "##", Width: 3, Alignment: "left", Color: ""},
		{Name: "name", Header: "FULL CREW NAME", Fill: true, MinWidth: 20, Alignment: "left", Color: "", Truncatable: true, TruncatedMinWidth: 12, TruncationTail: "..."},
		{Name: "id", Header: "ID", Width: 4, Alignment: "left", Color: ""},        // Just ID content width
		{Name: "members", Header: "MEMBERS", Width: 7, Alignment: "right", Color: ""}, // Header width
		{Name: "assets", Header: "ASSETS", Width: 6, Alignment: "right", Color: ""},   // Header width
	})

	// Build table data
	var data [][]string
	for _, match := range matches {
		row := []string{
			"",
			r.applyColor(fmt.Sprintf("%02d.", *itemNumber), design.ColorRowNumber),
			r.applyColor(match.DisplayName, design.ColorLabel),
			r.applyColor(match.SubInfo1, design.ColorAttribute),
			r.applyColor(match.SubInfo2, design.ColorAttribute),
			r.applyColor(match.SubInfo3, design.ColorAttribute),
		}
		data = append(data, row)
		*itemNumber++
	}

	// Add data to layout
	layout.Data = data

	dataTable := components.NewDataTableRowFromLayout(width, layout)

	// Render header
	headerRow := []string{
		"",
		r.applyColor("##", design.ColorTableHeader),
		r.applyColor("FULL CREW NAME", design.ColorTableHeader),
		r.applyColor("ID", design.ColorTableHeader),
		r.applyColor("MEMBERS", design.ColorTableHeader),
		r.applyColor("ASSETS", design.ColorTableHeader),
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

// renderAssetMatches renders the asset matches table
func (r *SmartSearchRenderer) renderAssetMatches(output *strings.Builder, matches []SearchResult, itemNumber *int, width int) error {
	// Define columns for asset matches table using enhanced gutter system
	layout := components.NewEnhancedTableLayout(4, []components.ColumnDefinition{
		{Name: "prefix", Header: "", Width: 3, Alignment: "left", Color: ""},
		{Name: "number", Header: "##", Width: 3, Alignment: "left", Color: ""},
		{Name: "name", Header: "ASSET NAME", Fill: true, MinWidth: 20, Alignment: "left", Color: "", Truncatable: true, TruncatedMinWidth: 12, TruncationTail: "..."},
		{Name: "version", Header: "VERSION", Width: 7, Alignment: "left", Color: ""}, // Header width
		{Name: "id", Header: "ID", Width: 8, Alignment: "left", Color: ""},           // AC123456 width
		{Name: "owner", Header: "OWNER", Width: 9, Alignment: "left", Color: ""},     // CREW-1234 width
	})

	// Build table data
	var data [][]string
	for _, match := range matches {
		row := []string{
			"",
			r.applyColor(fmt.Sprintf("%02d.", *itemNumber), design.ColorRowNumber),
			r.applyColor(match.DisplayName, design.ColorLabel),
			r.applyColor(match.SubInfo1, design.ColorAttribute),
			r.applyColor(match.SubInfo2, design.ColorAttribute),
			r.applyColor(match.SubInfo3, design.ColorAttribute),
		}
		data = append(data, row)
		*itemNumber++
	}

	// Add data to layout
	layout.Data = data

	dataTable := components.NewDataTableRowFromLayout(width, layout)

	// Render header
	headerRow := []string{
		"",
		r.applyColor("##", design.ColorTableHeader),
		r.applyColor("ASSET NAME", design.ColorTableHeader),
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

// renderUserMatches renders the user matches table
func (r *SmartSearchRenderer) renderUserMatches(output *strings.Builder, matches []SearchResult, itemNumber *int, width int) error {
	// Define columns for user matches table using enhanced gutter system
	layout := components.NewEnhancedTableLayout(4, []components.ColumnDefinition{
		{Name: "prefix", Header: "", Width: 3, Alignment: "left", Color: ""},
		{Name: "number", Header: "##", Width: 3, Alignment: "left", Color: ""},
		{Name: "name", Header: "FULL NAME", Fill: true, MinWidth: 20, Alignment: "left", Color: "", Truncatable: true, TruncatedMinWidth: 12, TruncationTail: "..."},
		{Name: "ldap", Header: "LDAP", Width: 14, Alignment: "left", Color: ""},     // @webinfra.lead width
		{Name: "level", Header: "LEVEL", Width: 9, Alignment: "left", Color: ""},    // Principal width
		{Name: "memberships", Header: "CREWS", Width: 5, Alignment: "right", Color: ""}, // Header width
	})

	// Build table data
	var data [][]string
	for _, match := range matches {
		row := []string{
			"",
			r.applyColor(fmt.Sprintf("%02d.", *itemNumber), design.ColorRowNumber),
			r.applyColor(match.DisplayName, design.ColorLabel),
			r.applyColor(match.SubInfo1, design.ColorAttribute),
			r.applyColor(match.SubInfo2, design.ColorAttribute),
			r.applyColor(match.SubInfo3, design.ColorAttribute),
		}
		data = append(data, row)
		*itemNumber++
	}

	// Add data to layout
	layout.Data = data

	dataTable := components.NewDataTableRowFromLayout(width, layout)

	// Render header
	headerRow := []string{
		"",
		r.applyColor("##", design.ColorTableHeader),
		r.applyColor("FULL NAME", design.ColorTableHeader),
		r.applyColor("LDAP", design.ColorTableHeader),
		r.applyColor("LEVEL", design.ColorTableHeader),
		r.applyColor("CREWS", design.ColorTableHeader),
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

// searchCrews performs fuzzy search on crew names and descriptions
func (r *SmartSearchRenderer) searchCrews(searchTerm string) ([]SearchResult, error) {
	var results []SearchResult

	// Cast dataStore to SimulationDataStore to access search methods
	simulationStore, ok := r.dataStore.(*data.SimulationDataStore)
	if !ok {
		return results, fmt.Errorf("dataStore is not a SimulationDataStore")
	}

	crews, err := simulationStore.SearchCrews(searchTerm)
	if err != nil {
		return nil, fmt.Errorf("crew search failed: %w", err)
	}

	for _, crew := range crews {
		// Count assets owned by this crew
		assetCount := len(crew.OwnedAssets)

		results = append(results, SearchResult{
			Type:        "crew",
			ID:          crew.ID,
			DisplayName: crew.VanityName,
			SubInfo1:    strings.TrimPrefix(crew.ID, "CREW-"), // Just the numeric part
			SubInfo2:    fmt.Sprintf("%d", len(crew.Members)),
			SubInfo3:    fmt.Sprintf("%d", assetCount),
			OriginalRef: crew,
		})
	}

	// Sort by relevance (by name for now)
	sort.Slice(results, func(i, j int) bool {
		return results[i].DisplayName < results[j].DisplayName
	})

	return results, nil
}

// searchAssets performs fuzzy search on asset names and descriptions
func (r *SmartSearchRenderer) searchAssets(searchTerm string) ([]SearchResult, error) {
	var results []SearchResult

	// Cast dataStore to SimulationDataStore to access search methods
	simulationStore, ok := r.dataStore.(*data.SimulationDataStore)
	if !ok {
		return results, fmt.Errorf("dataStore is not a SimulationDataStore")
	}

	assets, err := simulationStore.SearchAssets(searchTerm)
	if err != nil {
		return nil, fmt.Errorf("asset search failed: %w", err)
	}

	for _, asset := range assets {
		results = append(results, SearchResult{
			Type:        "asset",
			ID:          asset.ID,
			DisplayName: asset.VanityName,
			SubInfo1:    asset.Version,
			SubInfo2:    asset.ID,
			SubInfo3:    asset.OwningCrewID,
			OriginalRef: asset,
		})
	}

	// Sort by relevance (by name for now)
	sort.Slice(results, func(i, j int) bool {
		return results[i].DisplayName < results[j].DisplayName
	})

	return results, nil
}

// searchUsers performs fuzzy search on user names and LDAP IDs
func (r *SmartSearchRenderer) searchUsers(searchTerm string) ([]SearchResult, error) {
	var results []SearchResult

	// Cast dataStore to SimulationDataStore to access search methods
	simulationStore, ok := r.dataStore.(*data.SimulationDataStore)
	if !ok {
		return results, fmt.Errorf("dataStore is not a SimulationDataStore")
	}

	users, err := simulationStore.SearchUsers(searchTerm)
	if err != nil {
		return nil, fmt.Errorf("user search failed: %w", err)
	}

	for _, user := range users {
		// Count total memberships across all crews
		membershipCount := 0
		allCrews, _ := simulationStore.GetAllCrews()
		for _, crew := range allCrews {
			for _, member := range crew.Members {
				if member.UserID == user.UserID {
					membershipCount++
					break // Don't count same user multiple times in same crew
				}
			}
		}

		results = append(results, SearchResult{
			Type:        "user",
			ID:          user.UserID,
			DisplayName: user.FullName,
			SubInfo1:    "@" + user.UserID,
			SubInfo2:    user.Level,
			SubInfo3:    fmt.Sprintf("%d", membershipCount),
			OriginalRef: user,
		})
	}

	// Sort by relevance (by name for now)
	sort.Slice(results, func(i, j int) bool {
		return results[i].DisplayName < results[j].DisplayName
	})

	return results, nil
}

// GetResultByNumber returns a search result by its display number
func (r *SmartSearchRenderer) GetResultByNumber(results *SmartSearchResults, number int) *SearchResult {
	allResults := []SearchResult{}
	allResults = append(allResults, results.CrewMatches...)
	allResults = append(allResults, results.AssetMatches...)
	allResults = append(allResults, results.UserMatches...)

	if number < 1 || number > len(allResults) {
		return nil
	}

	return &allResults[number-1]
}

// applyColor applies ANSI color codes with 256-color support
func (r *SmartSearchRenderer) applyColor(text, colorCode string) string {
	if colorCode == "" || text == "" {
		return text
	}
	// Handle 256-color codes (3+ digits) vs basic ANSI codes (1-2 digits)
	if len(colorCode) >= 3 {
		return fmt.Sprintf("\033[38;5;%sm%s\033[%sm", colorCode, text, design.ColorReset)
	}
	return fmt.Sprintf("\033[%sm%s\033[%sm", colorCode, text, design.ColorReset)
}