package config

import (
	"fmt"
	"os"
	"strings"
)

// VerbosityLevel represents the different verbosity levels available
type VerbosityLevel int

const (
	// VerbosityQuiet shows only essential information
	VerbosityQuiet VerbosityLevel = iota
	// VerbosityConcise shows less detail with components hidden
	VerbosityConcise
	// VerbosityDefault shows current view (no flag needed)
	VerbosityDefault
	// VerbosityVerbose shows enhanced details with progress bars for multi-step
	VerbosityVerbose
	// VerbosityDebug shows maximum verbosity with all system outputs
	VerbosityDebug
)

// String returns the string representation of the verbosity level
func (v VerbosityLevel) String() string {
	switch v {
	case VerbosityQuiet:
		return "quiet"
	case VerbosityConcise:
		return "concise"
	case VerbosityDefault:
		return "default"
	case VerbosityVerbose:
		return "verbose"
	case VerbosityDebug:
		return "debug"
	default:
		return "unknown"
	}
}

// ParseVerbosityLevel parses a string into a VerbosityLevel
func ParseVerbosityLevel(s string) (VerbosityLevel, error) {
	switch strings.ToLower(s) {
	case "quiet", "q":
		return VerbosityQuiet, nil
	case "concise", "c":
		return VerbosityConcise, nil
	case "default", "d", "":
		return VerbosityDefault, nil
	case "verbose", "v":
		return VerbosityVerbose, nil
	case "debug":
		return VerbosityDebug, nil
	default:
		return VerbosityDefault, fmt.Errorf("invalid verbosity level: %s", s)
	}
}

// VerbosityConfig holds the verbosity configuration
type VerbosityConfig struct {
	Level           VerbosityLevel `json:"level"`
	ShowProgress    bool           `json:"show_progress"`
	ShowTimings     bool           `json:"show_timings"`
	ShowComponents  bool           `json:"show_components"`
	ShowDebugInfo   bool           `json:"show_debug_info"`
	ShowSystemInfo  bool           `json:"show_system_info"`
	DetailLevel     int            `json:"detail_level"`     // 1-5 scale
	ProgressFormat  string         `json:"progress_format"`  // "simple", "detailed", "spinner"
	OutputFormat    string         `json:"output_format"`    // "standard", "compact", "expanded"
}

// NewVerbosityConfig creates a new verbosity configuration based on the level
func NewVerbosityConfig(level VerbosityLevel) *VerbosityConfig {
	config := &VerbosityConfig{
		Level: level,
	}

	// Configure settings based on verbosity level
	switch level {
	case VerbosityQuiet:
		config.ShowProgress = false
		config.ShowTimings = false
		config.ShowComponents = false
		config.ShowDebugInfo = false
		config.ShowSystemInfo = false
		config.DetailLevel = 1
		config.ProgressFormat = "simple"
		config.OutputFormat = "compact"

	case VerbosityConcise:
		config.ShowProgress = true
		config.ShowTimings = false
		config.ShowComponents = false
		config.ShowDebugInfo = false
		config.ShowSystemInfo = false
		config.DetailLevel = 2
		config.ProgressFormat = "simple"
		config.OutputFormat = "standard"

	case VerbosityDefault:
		config.ShowProgress = true
		config.ShowTimings = true
		config.ShowComponents = true
		config.ShowDebugInfo = false
		config.ShowSystemInfo = false
		config.DetailLevel = 3
		config.ProgressFormat = "detailed"
		config.OutputFormat = "standard"

	case VerbosityVerbose:
		config.ShowProgress = true
		config.ShowTimings = true
		config.ShowComponents = true
		config.ShowDebugInfo = false
		config.ShowSystemInfo = true
		config.DetailLevel = 4
		config.ProgressFormat = "detailed"
		config.OutputFormat = "expanded"

	case VerbosityDebug:
		config.ShowProgress = true
		config.ShowTimings = true
		config.ShowComponents = true
		config.ShowDebugInfo = true
		config.ShowSystemInfo = true
		config.DetailLevel = 5
		config.ProgressFormat = "detailed"
		config.OutputFormat = "expanded"
	}

	return config
}

// DetermineVerbosityLevel determines the verbosity level based on CLI flags and environment
// Precedence: CLI flags > Environment variables > Default
func DetermineVerbosityLevel(quiet, concise, verbose, debug bool) VerbosityLevel {
	// CLI flags take highest precedence
	if debug {
		return VerbosityDebug
	}
	if verbose {
		return VerbosityVerbose
	}
	if concise {
		return VerbosityConcise
	}
	if quiet {
		return VerbosityQuiet
	}

	// Check environment variable
	if envLevel := os.Getenv("ENGX_VERBOSITY"); envLevel != "" {
		if level, err := ParseVerbosityLevel(envLevel); err == nil {
			return level
		}
	}

	// Default to normal verbosity
	return VerbosityDefault
}

// ShouldShow determines if content should be shown based on verbosity configuration
func (vc *VerbosityConfig) ShouldShow(contentType string) bool {
	switch contentType {
	case "progress":
		return vc.ShowProgress
	case "timings":
		return vc.ShowTimings
	case "components":
		return vc.ShowComponents
	case "debug":
		return vc.ShowDebugInfo
	case "system":
		return vc.ShowSystemInfo
	default:
		return true // Show unknown content types by default
	}
}

// ShouldShowDetailLevel determines if content should be shown based on detail level
func (vc *VerbosityConfig) ShouldShowDetailLevel(requiredLevel int) bool {
	return vc.DetailLevel >= requiredLevel
}

// GetProgressFormat returns the appropriate progress format for this verbosity level
func (vc *VerbosityConfig) GetProgressFormat() string {
	return vc.ProgressFormat
}

// GetOutputFormat returns the appropriate output format for this verbosity level
func (vc *VerbosityConfig) GetOutputFormat() string {
	return vc.OutputFormat
}

// IsQuiet returns true if the verbosity level is quiet
func (vc *VerbosityConfig) IsQuiet() bool {
	return vc.Level == VerbosityQuiet
}

// IsConcise returns true if the verbosity level is concise
func (vc *VerbosityConfig) IsConcise() bool {
	return vc.Level == VerbosityConcise
}

// IsDefault returns true if the verbosity level is default
func (vc *VerbosityConfig) IsDefault() bool {
	return vc.Level == VerbosityDefault
}

// IsVerbose returns true if the verbosity level is verbose
func (vc *VerbosityConfig) IsVerbose() bool {
	return vc.Level == VerbosityVerbose
}

// IsDebug returns true if the verbosity level is debug
func (vc *VerbosityConfig) IsDebug() bool {
	return vc.Level == VerbosityDebug
}

// DebugPrint prints debug information if debug mode is enabled
func (vc *VerbosityConfig) DebugPrint(format string, args ...interface{}) {
	if vc.IsDebug() {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
	}
}

// VerbosePrint prints verbose information if verbose mode or higher is enabled
func (vc *VerbosityConfig) VerbosePrint(format string, args ...interface{}) {
	if vc.Level >= VerbosityVerbose {
		fmt.Fprintf(os.Stderr, "[VERBOSE] "+format+"\n", args...)
	}
}