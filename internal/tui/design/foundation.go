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