package errors

import (
	"fmt"
)

// ErrorScenario represents a simulated error with recovery information
type ErrorScenario struct {
	Code        string
	Title       string
	Description string
	Causes      []string
	Actions     []Action
	QuickFix    *QuickFix
	HelpCommand string
	Severity    Severity
}

// Action represents a suggested action for error recovery
type Action struct {
	Description string
	Command     string
	AutoFix     bool
}

// QuickFix represents an automated fix option
type QuickFix struct {
	Description string
	Command     string
	Confirmation bool // Whether to ask for user confirmation
}

// Severity represents the severity level of an error
type Severity int

const (
	SeverityLow Severity = iota
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

// Common error scenarios for the create command
var CreateErrorScenarios = map[string]*ErrorScenario{
	"CONFIG_INVALID": {
		Code:        "CONFIG_INVALID",
		Title:       "Invalid Configuration",
		Description: "Project configuration contains invalid or conflicting settings",
		Causes: []string{
			"Project name contains invalid characters",
			"Target directory already exists",
			"Missing required configuration values",
		},
		Actions: []Action{
			{
				Description: "Validate project name follows naming conventions",
				Command:     "dpx-web validate name",
				AutoFix:     false,
			},
			{
				Description: "Check if directory exists and choose different name",
				Command:     "ls -la | grep [project-name]",
				AutoFix:     false,
			},
			{
				Description: "Use configuration wizard to set up valid config",
				Command:     "dpx-web config init",
				AutoFix:     true,
			},
		},
		QuickFix: &QuickFix{
			Description:  "Run configuration validator and auto-fix common issues",
			Command:      "dpx-web fix config",
			Confirmation: true,
		},
		HelpCommand: "dpx-web help create troubleshooting",
		Severity:    SeverityMedium,
	},

	"NETWORK_ERROR": {
		Code:        "NETWORK_ERROR",
		Title:       "Network Connection Failed",
		Description: "Unable to download dependencies due to network connectivity issues",
		Causes: []string{
			"No internet connection available",
			"Corporate firewall blocking npm registry",
			"NPM registry is temporarily unavailable",
			"Proxy configuration required",
		},
		Actions: []Action{
			{
				Description: "Check internet connectivity",
				Command:     "ping google.com",
				AutoFix:     false,
			},
			{
				Description: "Test NPM registry access",
				Command:     "npm ping",
				AutoFix:     false,
			},
			{
				Description: "Configure corporate proxy settings",
				Command:     "npm config set proxy http://your-proxy:port",
				AutoFix:     false,
			},
			{
				Description: "Use offline mode with cached dependencies",
				Command:     "dpx-web create --offline",
				AutoFix:     true,
			},
		},
		QuickFix: &QuickFix{
			Description:  "Run network diagnostics and attempt automatic fixes",
			Command:      "dpx-web fix network",
			Confirmation: true,
		},
		HelpCommand: "dpx-web help network troubleshooting",
		Severity:    SeverityHigh,
	},

	"PERMISSION_DENIED": {
		Code:        "PERMISSION_DENIED",
		Title:       "Insufficient Permissions",
		Description: "Unable to create files or directories due to permission restrictions",
		Causes: []string{
			"Target directory requires elevated permissions",
			"File system is read-only",
			"User account lacks write permissions",
		},
		Actions: []Action{
			{
				Description: "Check directory permissions",
				Command:     "ls -la .",
				AutoFix:     false,
			},
			{
				Description: "Change to a directory you have write access to",
				Command:     "cd ~/Documents && dpx-web create [app-name]",
				AutoFix:     false,
			},
			{
				Description: "Run with elevated permissions (if appropriate)",
				Command:     "sudo dpx-web create [app-name]",
				AutoFix:     false,
			},
		},
		QuickFix: &QuickFix{
			Description:  "Suggest alternative directory with proper permissions",
			Command:      "dpx-web fix permissions",
			Confirmation: true,
		},
		HelpCommand: "dpx-web help permissions",
		Severity:    SeverityMedium,
	},

	"DEPENDENCY_CONFLICT": {
		Code:        "DEPENDENCY_CONFLICT",
		Title:       "Dependency Version Conflict",
		Description: "Conflicting package versions detected in dependency resolution",
		Causes: []string{
			"Global packages conflict with project requirements",
			"Peer dependency version mismatch",
			"Outdated package manager version",
		},
		Actions: []Action{
			{
				Description: "Clear NPM cache and retry",
				Command:     "npm cache clean --force",
				AutoFix:     true,
			},
			{
				Description: "Update package manager to latest version",
				Command:     "npm install -g npm@latest",
				AutoFix:     false,
			},
			{
				Description: "Use fresh dependency resolution",
				Command:     "rm -rf node_modules package-lock.json && npm install",
				AutoFix:     true,
			},
		},
		QuickFix: &QuickFix{
			Description:  "Clean dependency cache and retry with fresh resolution",
			Command:      "dpx-web fix dependencies",
			Confirmation: true,
		},
		HelpCommand: "dpx-web help dependencies troubleshooting",
		Severity:    SeverityMedium,
	},

	"DISK_SPACE": {
		Code:        "DISK_SPACE",
		Title:       "Insufficient Disk Space",
		Description: "Not enough free disk space to complete project creation",
		Causes: []string{
			"Target drive is full or nearly full",
			"Temporary directory lacks space",
			"Node modules require significant space",
		},
		Actions: []Action{
			{
				Description: "Check available disk space",
				Command:     "df -h .",
				AutoFix:     false,
			},
			{
				Description: "Clean temporary files and caches",
				Command:     "npm cache clean --force && rm -rf /tmp/*",
				AutoFix:     true,
			},
			{
				Description: "Choose different target directory with more space",
				Command:     "dpx-web create [app-name] --target /path/with/space",
				AutoFix:     false,
			},
		},
		QuickFix: &QuickFix{
			Description:  "Clean caches and temporary files to free up space",
			Command:      "dpx-web fix disk-space",
			Confirmation: true,
		},
		HelpCommand: "dpx-web help disk-space",
		Severity:    SeverityHigh,
	},
}

// GetErrorScenario returns a specific error scenario by code
func GetErrorScenario(code string) *ErrorScenario {
	return CreateErrorScenarios[code]
}

// GetRandomErrorScenario returns a random error scenario for simulation
func GetRandomErrorScenario() *ErrorScenario {
	scenarioKeys := []string{
		"CONFIG_INVALID",
		"NETWORK_ERROR",
		"PERMISSION_DENIED",
		"DEPENDENCY_CONFLICT",
		"DISK_SPACE",
	}

	// In a real implementation, this would use proper randomization
	// For demo purposes, return the network error as it's common and relatable
	_ = scenarioKeys // Acknowledge the variable for future randomization implementation
	return CreateErrorScenarios["NETWORK_ERROR"]
}

// FormatErrorMessage formats an error scenario into a user-friendly message
func FormatErrorMessage(scenario *ErrorScenario) string {
	msg := fmt.Sprintf("❌ Error: %s\n\n", scenario.Title)
	msg += fmt.Sprintf("🔍 Problem: %s\n", scenario.Description)

	if len(scenario.Causes) > 0 {
		msg += "📋 Likely Causes:\n"
		for _, cause := range scenario.Causes {
			msg += fmt.Sprintf("   • %s\n", cause)
		}
		msg += "\n"
	}

	if len(scenario.Actions) > 0 {
		msg += "🛠️  Suggested Actions:\n"
		for i, action := range scenario.Actions {
			msg += fmt.Sprintf("   %d. %s: %s\n", i+1, action.Description, action.Command)
		}
		msg += "\n"
	}

	if scenario.QuickFix != nil {
		msg += fmt.Sprintf("⚡ Quick Fix Available:\n   Run '%s' to %s\n\n",
			scenario.QuickFix.Command, scenario.QuickFix.Description)
	}

	if scenario.HelpCommand != "" {
		msg += fmt.Sprintf("❓ More Help: %s\n", scenario.HelpCommand)
	}

	return msg
}

// ShouldSimulateError determines if an error should be simulated based on error rate
func ShouldSimulateError(errorRate float64) bool {
	// For demo purposes, simulate errors at a lower rate than specified
	// In a real implementation, this would use proper random number generation
	return errorRate > 0.10 // Only simulate if error rate is above 10%
}