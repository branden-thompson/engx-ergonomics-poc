package styles

import "github.com/charmbracelet/lipgloss"

// Color palette for DPX Web theme
//
// STYLING CONSISTENCY RULES (based on create command AAR):
// 1. Dividers/separators: DarkGray (ANSI 90)
// 2. Main section headers: BrightWhite (ANSI 97) - NOT BOLD
// 3. Backticked commands: BrightMagenta (ANSI 95)
// 4. Table headers: LightPurple - NOT BOLD
// 5. Main content labels: BrightWhite (ANSI 97)
// 6. Secondary content: Gray500 (lighter gray)
// 7. Numbers/indexes: DarkGray (ANSI 90)
// 8. Success indicators: Success green
// 9. Special highlights: Primary colors
var (
	// Primary colors
	Primary   = lipgloss.Color("#06B6D4") // Cyan-500
	Secondary = lipgloss.Color("#8B5CF6") // Violet-500
	Success   = lipgloss.Color("#10B981") // Emerald-500
	Warning   = lipgloss.Color("#F59E0B") // Amber-500
	Error     = lipgloss.Color("#EF4444") // Red-500
	Info      = lipgloss.Color("#3B82F6") // Blue-500

	// Neutral colors
	White     = lipgloss.Color("#FFFFFF")
	Gray100   = lipgloss.Color("#F3F4F6")
	Gray200   = lipgloss.Color("#E5E7EB")
	Gray300   = lipgloss.Color("#D1D5DB")
	Gray400   = lipgloss.Color("#9CA3AF")
	Gray500   = lipgloss.Color("#6B7280")
	Gray600   = lipgloss.Color("#4B5563")
	Gray700   = lipgloss.Color("#374151")
	Gray800   = lipgloss.Color("#1F2937")
	Gray900   = lipgloss.Color("#111827")

	// Special colors
	Highlight = lipgloss.Color("#FBBF24") // Amber-400
	Muted     = Gray500
	Border    = Gray300
	Background = Gray100

	// Terminal ANSI colors (matching create command styling exactly)
	ANSI_BrightWhite   = lipgloss.Color("97")  // ^[[97m - bright white
	ANSI_BrightMagenta = lipgloss.Color("95")  // ^[[95m - bright magenta
	ANSI_DarkGray      = lipgloss.Color("90")  // ^[[90m - dark gray
	ANSI_BrightGreen   = lipgloss.Color("92")  // ^[[92m - bright green
	ANSI_Orange        = lipgloss.Color("208") // ^[[38;5;208m - orange
)

// Common style presets
var (
	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			MarginBottom(1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	InfoStyle = lipgloss.NewStyle().
			Foreground(Info)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	// Layout styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(1, 2)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Border).
			Padding(0, 1)

	HeaderBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Background(Background).
			Padding(1, 2).
			MarginBottom(1)

	ProgressBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Padding(1, 2).
			MarginBottom(1)

	LogsBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Gray400).
			Height(8).
			Padding(1, 2).
			MarginBottom(1)

	FooterStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(Gray300).
			Padding(1, 2).
			Foreground(Muted)
)