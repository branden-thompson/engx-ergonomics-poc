package design

// Foundation defines the core design system constants and tokens
// for terminal UI components. This forms the basis for a future
// "Terminal Design and Component Language" extraction.

// Color System - ANSI-based color tokens
const (
	// Primary semantic colors
	ColorDarkGray      = "90"  // \033[90m - dividers, secondary text
	ColorBrightWhite   = "97"  // \033[97m - primary text, headers
	ColorBrightMagenta = "95"  // \033[95m - backtick commands
	ColorBrightGreen   = "92"  // \033[92m - success indicators
	ColorBrightYellow  = "93"  // \033[93m - warnings, subcmd
	ColorBrightRed     = "91"  // \033[91m - errors
	ColorBrightBlue    = "94"  // \033[94m - option values, arguments
	ColorBrightCyan    = "96"  // \033[96m - highlights
	ColorOrange        = "208" // \033[38;5;208m - cmd
	ColorLightPurple   = "135" // \033[38;5;135m - table headers
	ColorEngxPink      = "198" // \033[38;5;198m - #F06 equivalent for 'engx'
	ColorFlagGreen     = "48"  // \033[38;5;48m - #0FD equivalent for flags

	// Semantic colors for data tables and status indicators
	ColorPrimaryOnCall   = "196" // \033[38;5;196m - deep red for primary on-call (#D62 approximation)
	ColorSecondaryOnCall = "202" // \033[38;5;202m - orange-red for secondary on-call (#DA2 approximation)
	ColorTableHeader     = "95"  // \033[95m - bright magenta/light purple for table headers
	ColorRowNumber       = "90"  // \033[90m - dark grey for row numbers (same as ColorDarkGray)
	ColorAttribute       = "90"  // \033[90m - dark grey for attributes (same as ColorDarkGray)
	ColorDescription     = "90"  // \033[90m - dark grey for description text (same as ColorDarkGray)
	ColorLabel           = "97"  // \033[97m - bright white for labels/names (same as ColorBrightWhite)

	// Username colors
	ColorCurrentUser     = "92"  // \033[92m - bright green for current user's username
	ColorOtherUser       = "96"  // \033[96m - bright cyan/light blue for other usernames

	// Role colors (approximating hex colors with ANSI equivalents)
	ColorRoleOwner       = "95"  // \033[95m - bright magenta for Owner (#A4C approximation)
	ColorRoleAdmin       = "93"  // \033[93m - bright yellow for Admin (#FA0 approximation)
	ColorRoleMember      = "92"  // \033[92m - bright green for Member (#3AD approximation)
	ColorRoleAuto        = "94"  // \033[94m - bright blue for Auto (#69A approximation)
	ColorRoleTemp        = "91"  // \033[91m - bright red for Temp (#C66 approximation)

	// IC Level colors (Individual Contributor levels - 256-color approximations)
	ColorIC1             = "247" // \033[38;5;247m - light grey for IC1 (#999 approximation)
	ColorIC2             = "159" // \033[38;5;159m - light blue for IC2 (#9CF approximation)
	ColorIC3             = "75"  // \033[38;5;75m - bright blue for IC3 (#59D approximation)
	ColorIC4             = "33"  // \033[38;5;33m - blue for IC4 (#38E approximation)
	ColorIC5             = "87"  // \033[38;5;87m - cyan for IC5 (#3BE approximation)
	ColorIC6             = "39"  // \033[38;5;39m - bright cyan for IC6 (#0BF approximation)
	ColorIC7             = "80"  // \033[38;5;80m - cyan-green for IC7 (#5DD approximation)
	ColorIC8             = "43"  // \033[38;5;43m - green-cyan for IC8 (#2DC approximation)
	ColorIC9             = "51"  // \033[38;5;51m - bright cyan for IC9 (#0FD approximation)

	// MR Level colors (Management levels - 256-color approximations)
	ColorMR2             = "143" // \033[38;5;143m - brown/tan for MR2 (#8B7 approximation)
	ColorMR3             = "107" // \033[38;5;107m - olive for MR3 (#6B3 approximation)
	ColorMR4             = "136" // \033[38;5;136m - dark yellow for MR4 (#5C1 approximation)
	ColorMR5             = "71"  // \033[38;5;71m - green for MR5 (#3D3 approximation)
	ColorMR6             = "73"  // \033[38;5;73m - teal for MR6 (#3D8 approximation)
	ColorMR7             = "30"  // \033[38;5;30m - dark cyan for MR7 (#1E8 approximation)
	ColorMR8             = "36"  // \033[38;5;36m - cyan for MR8 (#0F8 approximation)
	ColorMR9             = "51"  // \033[38;5;51m - bright cyan for MR9 (#0FD approximation)

	// Crew type colors
	ColorStandardCrew    = "94"  // \033[94m - bright blue for STANDARD CREW (#08F approximation)
	ColorVirtualCrew     = "177" // \033[38;5;177m - light purple for VIRTUAL CREW (#D0F approximation)

	// Reset
	ColorReset = "0" // \033[0m - reset to default
)

// Spacing constants for consistent layout
const (
	SpaceNone   = 0
	SpaceXS     = 1
	SpaceSM     = 2
	SpaceMD     = 4
	SpaceLG     = 8
	SpaceXL     = 12
	SpaceXXL    = 16
)

// Typography scale for consistent text hierarchy
type TextStyle struct {
	Color     string
	Bold      bool
	Underline bool
}

var (
	// Text style presets
	StyleTitle     = TextStyle{Color: ColorBrightWhite, Bold: true}
	StyleHeader    = TextStyle{Color: ColorBrightWhite, Bold: false}
	StyleSubheader = TextStyle{Color: ColorDarkGray, Bold: false}
	StyleBody      = TextStyle{Color: ColorBrightWhite, Bold: false}
	StyleMuted     = TextStyle{Color: ColorDarkGray, Bold: false}
	StyleCommand   = TextStyle{Color: ColorBrightMagenta, Bold: false}
	StyleSuccess   = TextStyle{Color: ColorBrightGreen, Bold: false}
	StyleWarning   = TextStyle{Color: ColorBrightYellow, Bold: false}
	StyleError     = TextStyle{Color: ColorBrightRed, Bold: true}
)

// Layout constants for consistent component sizing
const (
	// Standard component widths
	WidthNarrow  = 40
	WidthMedium  = 60
	WidthWide    = 80
	WidthFull    = 120

	// Table column sizing
	ColXS  = 3   // numbers, icons
	ColSM  = 8   // short labels
	ColMD  = 12  // medium text
	ColLG  = 18  // long text
	ColXL  = 24  // very long text
	ColXXL = 32  // extra long text
)

// Border characters for consistent visual elements
const (
	BorderHorizontal = "-"
	BorderVertical   = "|"
	BorderCornerTL   = "+"
	BorderCornerTR   = "+"
	BorderCornerBL   = "+"
	BorderCornerBR   = "+"
)

// Common patterns for visual elements
const (
	SeparatorShort  = "----"
	SeparatorMedium = "----------------------------------------"
	SeparatorLong   = "-------------------------------------------------------------------------------"
)