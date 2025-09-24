package renderers

// ColorScheme defines the color palette for crew command output
type ColorScheme struct {
	// ANSI Reset
	ResetColor string

	// Text colors
	HeaderColor    string
	TextColor      string
	AccentColor    string
	CrewIDColor    string

	// Role-based colors
	OwnerColor   string
	AdminColor   string
	MemberColor  string
	TempColor    string
	RemovedColor string

	// Status indicators
	OnCallActiveColor   string
	OnCallInactiveColor string
	ActiveColor         string
	InactiveColor       string

	// Icons and symbols
	OwnerIcon       string
	AdminIcon       string
	MemberIcon      string
	TempIcon        string

	// Section icons
	DescriptionIcon string
	CalendarIcon    string
	MembershipIcon  string
	OnCallIcon      string
	AssetsIcon      string
}

// DefaultColorScheme returns the default color scheme for crew commands
func DefaultColorScheme() ColorScheme {
	return ColorScheme{
		// ANSI Reset
		ResetColor: "\033[0m",

		// Text colors
		HeaderColor: "\033[97m",   // Bright white for headers
		TextColor:   "\033[37m",   // Light gray for regular text
		AccentColor: "\033[94m",   // Bright blue for accents
		CrewIDColor: "\033[94m",   // Bright blue for crew IDs

		// Role-based colors (matching existing engx patterns)
		OwnerColor:   "\033[92m",  // Bright green for owners
		AdminColor:   "\033[94m",  // Bright blue for admins
		MemberColor:  "\033[37m",  // Light gray for members
		TempColor:    "\033[93m",  // Bright yellow for temp
		RemovedColor: "\033[91m",  // Bright red for removed

		// Status indicators
		OnCallActiveColor:   "\033[92m", // Bright green for active on-call
		OnCallInactiveColor: "\033[90m", // Dark gray for inactive
		ActiveColor:         "\033[92m", // Bright green for active status
		InactiveColor:       "\033[90m", // Dark gray for inactive status

		// Icons and symbols (using Unicode symbols compatible with terminals)
		OwnerIcon:  "🟢", // Green circle for owners
		AdminIcon:  "🔵", // Blue circle for admins
		MemberIcon: "⚪", // White circle for members
		TempIcon:   "🟡", // Yellow circle for temp users

		// Section icons
		DescriptionIcon: "📋", // Clipboard for description
		CalendarIcon:    "📅", // Calendar for dates
		MembershipIcon:  "👥", // People for membership
		OnCallIcon:      "🚨", // Alert for on-call
		AssetsIcon:      "📦", // Package for assets
	}
}

// MonochromeColorScheme returns a monochrome color scheme for terminals without color support
func MonochromeColorScheme() ColorScheme {
	return ColorScheme{
		// ANSI Reset
		ResetColor: "\033[0m",

		// Text colors (all white/gray)
		HeaderColor: "\033[1m",   // Bold for headers
		TextColor:   "\033[0m",   // Normal text
		AccentColor: "\033[1m",   // Bold for accents
		CrewIDColor: "\033[1m",   // Bold for crew IDs

		// Role-based colors (using different text styles)
		OwnerColor:   "\033[1m",  // Bold for owners
		AdminColor:   "\033[4m",  // Underline for admins
		MemberColor:  "\033[0m",  // Normal for members
		TempColor:    "\033[3m",  // Italic for temp
		RemovedColor: "\033[9m",  // Strikethrough for removed

		// Status indicators
		OnCallActiveColor:   "\033[1m", // Bold for active
		OnCallInactiveColor: "\033[2m", // Dim for inactive
		ActiveColor:         "\033[1m", // Bold for active
		InactiveColor:       "\033[2m", // Dim for inactive

		// Icons and symbols (ASCII-safe)
		OwnerIcon:  "*", // Asterisk for owners
		AdminIcon:  "+", // Plus for admins
		MemberIcon: "o", // Circle for members
		TempIcon:   "~", // Tilde for temp users

		// Section icons (ASCII-safe)
		DescriptionIcon: "#", // Hash for description
		CalendarIcon:    "@", // At symbol for dates
		MembershipIcon:  "&", // Ampersand for membership
		OnCallIcon:      "!", // Exclamation for on-call
		AssetsIcon:      "$", // Dollar for assets
	}
}