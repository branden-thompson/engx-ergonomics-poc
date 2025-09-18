package aar

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// OutputFormatter interface for different AAR output formats
type OutputFormatter interface {
	Format(summary *AARSummary) string
}

// StandardFormatter provides the default AAR output format
type StandardFormatter struct {
	width int
}

// NewStandardFormatter creates a new standard formatter
func NewStandardFormatter(width int) *StandardFormatter {
	if width <= 0 {
		width = 80 // Default width
	}
	return &StandardFormatter{width: width}
}

// Format generates the standard AAR output using the specified template with terminal-width awareness and styling
func (f *StandardFormatter) Format(summary *AARSummary) string {
	var output strings.Builder

	// Use terminal width, with sensible defaults
	width := f.width
	if width <= 0 || width > 120 {
		width = 80 // Reasonable default
	}

	// Calculate total duration
	duration := summary.ExecutionInfo.EndTime.Sub(summary.ExecutionInfo.StartTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	durationStr := fmt.Sprintf("%02dh %02dm %02ds", hours, minutes, seconds)

	// Determine dev-only vs production ready status
	isDevOnly := summary.ProjectInfo.DevOnly
	setupType := "PRODUCTION READY"
	setupNote := ""
	if isDevOnly {
		setupType = "DEV ONLY"
		setupNote = fmt.Sprintf("   └ run `engx promote %s --production` if you need to \n      deploy this application to a production environment.", summary.ProjectInfo.Name)
	} else {
		setupNote = fmt.Sprintf("   └ run `engx deploy %s --production` when ready to \n      deploy this application to a production environment.", summary.ProjectInfo.Name)
	}

	// Get dev server command and port
	devCommand := "yarn dev"
	port := "3000" // Default port

	// Look for dev server command in next steps
	for _, step := range summary.NextSteps {
		if step.Category == CategoryDevelopment && strings.Contains(step.Command, "npm run dev") {
			devCommand = "npm run dev"
			break
		}
	}

	// Create text components (NO STYLING YET - SPACING ONLY)
	headerText := "AFTER ACTION SUMMARY"
	successText := "OPERATION SUCCESS"
	footerSteps := fmt.Sprintf("%d/%d Steps Completed", len(summary.StepResults), len(summary.StepResults))
	footerTime := fmt.Sprintf("Total Elapsed time: %s", durationStr)

	// Calculate spacing for full-width layout
	headerPadding := width - len(headerText) - len(successText) - 12 // Account for dashes, spaces, and margin
	if headerPadding < 1 {
		headerPadding = 1
	}

	footerPadding := width - len(footerSteps) - len(footerTime) - 12 // Account for dashes, spaces, and margin
	if footerPadding < 1 {
		footerPadding = 1
	}

	// ANSI color codes - exact same as progress table enhanced_renderer.go:10-23
	const (
		colorReset         = "\033[0m"
		colorWhite         = "\033[97m"  // Bright white
		colorGreen         = "\033[92m"  // Green for OPERATION SUCCESS (matches [installed])
		colorLightGrey     = "\033[90m"  // Darker grey for dashes
		colorBrightOrange  = "\033[38;5;208m"  // Bright orange for PRODUCTION READY
		colorBrightMagenta = "\033[95m"  // Bright magenta for terminal commands
	)

	// formatCommandsInBackticks applies magenta color to text wrapped in backticks
	// while preserving the surrounding color context
	formatCommandsInBackticks := func(text string) string {
		// Replace `command` with colored version, returning to grey after
		result := text
		for {
			start := strings.Index(result, "`")
			if start == -1 {
				break
			}
			end := strings.Index(result[start+1:], "`")
			if end == -1 {
				break
			}
			end = start + 1 + end

			// Extract command without backticks
			command := result[start+1:end]
			// Replace with colored version that returns to grey
			coloredCommand := fmt.Sprintf("%s%s%s%s", colorBrightMagenta, command, colorReset, colorLightGrey)
			result = result[:start] + coloredCommand + result[end+1:]
		}
		return result
	}

	// Build the AAR with exact spacing AND COLORS
	output.WriteString("\n")

	// Header line - exact template format with colors
	output.WriteString(fmt.Sprintf("%s----%s %s%s%s %s%s%s %s%s%s %s----%s\n",
		colorLightGrey, colorReset,  // Grey dashes
		colorWhite, headerText, colorReset,  // White title
		colorLightGrey, strings.Repeat("-", headerPadding), colorReset,  // Grey middle dashes
		colorGreen, successText, colorReset,  // Green OPERATION SUCCESS
		colorLightGrey, colorReset))  // Grey end dashes

	output.WriteString("  \n") // Empty line with leading spaces

	// Action items section - exact spacing from template with colors
	output.WriteString(fmt.Sprintf("  %sYou can now:%s\n", colorWhite, colorReset))
	output.WriteString(fmt.Sprintf("   %sLaunch your DEV server:%s      %s%s%s\n",
		colorLightGrey, colorReset, colorWhite, devCommand, colorReset))
	output.WriteString(fmt.Sprintf("   %sOpen in your editor:%s         %scode %s%s\n",
		colorLightGrey, colorReset, colorWhite, summary.ProjectInfo.Name, colorReset))
	output.WriteString("\n")

	// Setup status section with colors
	setupTypeColored := setupType
	if isDevOnly {
		setupTypeColored = fmt.Sprintf("%s%s%s", colorBrightOrange, setupType, colorReset)
	} else {
		setupTypeColored = fmt.Sprintf("%s%s%s", colorBrightOrange, setupType, colorReset)
	}
	output.WriteString(fmt.Sprintf("  %sThis application is set up to be%s %s.\n",
		colorWhite, colorReset, setupTypeColored))

	if setupNote != "" {
		// Apply magenta formatting to commands in backticks
		formattedSetupNote := formatCommandsInBackticks(setupNote)
		output.WriteString(fmt.Sprintf("%s%s%s\n", colorLightGrey, formattedSetupNote, colorReset))
	}

	// Local server info section with colors
	output.WriteString(fmt.Sprintf("\n  %sOnce running, your development server will be available at:%s\n",
		colorLightGrey, colorReset))
	output.WriteString(fmt.Sprintf("   %s└%s %shttp://localhost:%s%s\n\n",
		colorLightGrey, colorReset, colorWhite, port, colorReset))

	// Footer line - exact template format with colors
	output.WriteString(fmt.Sprintf("%s----%s %s%s%s %s%s%s %s%s%s %s----%s\n",
		colorLightGrey, colorReset,  // Grey dashes
		colorWhite, footerSteps, colorReset,  // White steps
		colorLightGrey, strings.Repeat("-", footerPadding), colorReset,  // Grey middle dashes
		colorWhite, footerTime, colorReset,  // White time
		colorLightGrey, colorReset))  // Grey end dashes

	return output.String()
}

// writeHeader writes the header section
func (f *StandardFormatter) writeHeader(output *strings.Builder, summary *AARSummary) {
	// Create header box
	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(f.width - 4).
		Align(lipgloss.Center).
		Bold(true)

	headerContent := fmt.Sprintf("✨ Project Creation Complete")

	separatorStyle := lipgloss.NewStyle().
		Border(lipgloss.Border{Top: "─"}).
		BorderForeground(lipgloss.Color("62")).
		Width(f.width - 4)

	projectInfoStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Width(f.width - 4)

	// Build project info
	projectInfo := fmt.Sprintf("Project: %s (%s)",
		summary.ProjectInfo.Name,
		strings.Title(summary.ProjectInfo.Template))

	if summary.ProjectInfo.Directory != "" {
		projectInfo += fmt.Sprintf("\nCreated: %s", summary.ProjectInfo.Directory)
	}

	projectInfo += fmt.Sprintf("\nDuration: %s", f.formatDuration(summary.ExecutionInfo.Duration))

	projectInfo += fmt.Sprintf("\nSteps: %d/%d completed successfully",
		summary.ExecutionInfo.SuccessSteps,
		summary.ExecutionInfo.TotalSteps)

	// Write header
	output.WriteString(headerStyle.Render(headerContent))
	output.WriteString("\n")
	output.WriteString(separatorStyle.Render(""))
	output.WriteString("\n")
	output.WriteString(projectInfoStyle.Render(projectInfo))
	output.WriteString("\n\n")
}

// writeExecutionSummary writes execution summary (if there are interesting details)
func (f *StandardFormatter) writeExecutionSummary(output *strings.Builder, summary *AARSummary) {
	// Only show execution summary if there were failures or performance issues
	if summary.ExecutionInfo.FailedSteps > 0 || f.hasPerformanceIssues(summary) {
		titleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Bold(true)

		output.WriteString(titleStyle.Render("📊 Execution Summary:"))
		output.WriteString("\n")

		if summary.ExecutionInfo.FailedSteps > 0 {
			errorStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("196"))

			output.WriteString(errorStyle.Render(fmt.Sprintf("  ❌ %d step(s) failed", summary.ExecutionInfo.FailedSteps)))
			output.WriteString("\n")
		}

		if summary.ExecutionInfo.SkippedSteps > 0 {
			skipStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("214"))

			output.WriteString(skipStyle.Render(fmt.Sprintf("  ⏭️  %d step(s) skipped", summary.ExecutionInfo.SkippedSteps)))
			output.WriteString("\n")
		}

		// Performance warnings
		if f.hasPerformanceIssues(summary) {
			warnStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("214"))

			if summary.ExecutionInfo.Duration > summary.ExecutionInfo.Performance.ConfigurableTargets["total_execution"] {
				output.WriteString(warnStyle.Render("  ⚠️  Execution took longer than expected"))
				output.WriteString("\n")
			}
		}

		output.WriteString("\n")
	}
}

// writeNextSteps writes the next steps section
func (f *StandardFormatter) writeNextSteps(output *strings.Builder, summary *AARSummary) {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("82")).
		Bold(true)

	output.WriteString(titleStyle.Render("🚀 Next Steps:"))
	output.WriteString("\n")

	// Show high priority steps first
	highPrioritySteps := f.filterStepsByPriority(summary.NextSteps, PriorityHigh, PriorityCritical)

	for i, step := range highPrioritySteps {
		if i >= 3 { // Limit to 3 high priority steps
			break
		}

		stepStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			PaddingLeft(2)

		stepText := fmt.Sprintf("%d. %s", i+1, step.Description)
		if step.Command != "" {
			commandStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("244")).
				Italic(true)
			stepText += "\n   " + commandStyle.Render(step.Command)
		}

		output.WriteString(stepStyle.Render(stepText))
		output.WriteString("\n")
	}

	output.WriteString("\n")
}

// writeQuickCommands writes a quick commands reference
func (f *StandardFormatter) writeQuickCommands(output *strings.Builder, summary *AARSummary) {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("117")).
		Bold(true)

	output.WriteString(titleStyle.Render("📚 Quick Commands:"))
	output.WriteString("\n")

	// Extract commands from next steps
	commands := f.extractCommands(summary)

	commandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		PaddingLeft(2)

	for _, cmd := range commands {
		if cmd.description != "" && cmd.command != "" {
			line := fmt.Sprintf("%-20s %s", cmd.command, cmd.description)
			output.WriteString(commandStyle.Render(line))
			output.WriteString("\n")
		}
	}

	output.WriteString("\n")
}

// writeResources writes additional resources section
func (f *StandardFormatter) writeResources(output *strings.Builder, summary *AARSummary) {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("177")).
		Bold(true)

	output.WriteString(titleStyle.Render("💡 Learn More:"))
	output.WriteString("\n")

	resourceStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		PaddingLeft(2)

	// Standard resources based on template
	resources := f.getStandardResources(summary.ProjectInfo.Template)
	for _, resource := range resources {
		output.WriteString(resourceStyle.Render(resource))
		output.WriteString("\n")
	}

	// Project-specific resources
	if summary.ProjectInfo.Directory != "" {
		readmePath := fmt.Sprintf("• Project README: %s/README.md", summary.ProjectInfo.Directory)
		output.WriteString(resourceStyle.Render(readmePath))
		output.WriteString("\n")
	}

	output.WriteString("\n")
}

// writeTroubleshooting writes troubleshooting information
func (f *StandardFormatter) writeTroubleshooting(output *strings.Builder, troubleshooting *TroubleshootingInfo) {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	output.WriteString(titleStyle.Render("🔧 Troubleshooting:"))
	output.WriteString("\n")

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		PaddingLeft(2)

	for _, failedStep := range troubleshooting.FailedSteps {
		output.WriteString(errorStyle.Render(fmt.Sprintf("❌ %s:", failedStep.StepName)))
		output.WriteString("\n")

		if failedStep.ErrorMessage != "" {
			msgStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("244")).
				PaddingLeft(4)
			output.WriteString(msgStyle.Render(fmt.Sprintf("Error: %s", failedStep.ErrorMessage)))
			output.WriteString("\n")
		}

		for _, suggestion := range failedStep.Suggestions {
			suggestionStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				PaddingLeft(4)
			output.WriteString(suggestionStyle.Render(fmt.Sprintf("• %s", suggestion)))
			output.WriteString("\n")
		}

		output.WriteString("\n")
	}

	// General suggestions
	if len(troubleshooting.Suggestions) > 0 {
		generalStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			PaddingLeft(2)

		output.WriteString(generalStyle.Render("💡 General Suggestions:"))
		output.WriteString("\n")

		for _, suggestion := range troubleshooting.Suggestions {
			suggestionStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				PaddingLeft(4)
			output.WriteString(suggestionStyle.Render(fmt.Sprintf("• %s", suggestion)))
			output.WriteString("\n")
		}
	}
}

// Helper methods

func (f *StandardFormatter) formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
}

func (f *StandardFormatter) hasPerformanceIssues(summary *AARSummary) bool {
	targets := summary.ExecutionInfo.Performance.ConfigurableTargets

	if totalTarget, exists := targets["total_execution"]; exists {
		return summary.ExecutionInfo.Duration > totalTarget
	}

	return false
}

func (f *StandardFormatter) filterStepsByPriority(steps []NextStep, priorities ...StepPriority) []NextStep {
	var filtered []NextStep

	for _, step := range steps {
		for _, priority := range priorities {
			if step.Priority == priority {
				filtered = append(filtered, step)
				break
			}
		}
	}

	return filtered
}

type commandRef struct {
	command     string
	description string
}

func (f *StandardFormatter) extractCommands(summary *AARSummary) []commandRef {
	var commands []commandRef

	// Extract from high priority next steps
	highPrioritySteps := f.filterStepsByPriority(summary.NextSteps, PriorityHigh, PriorityCritical)

	for _, step := range highPrioritySteps {
		if step.Command != "" {
			// Extract just the command part (after cd part if present)
			command := step.Command
			if strings.Contains(command, " && ") {
				parts := strings.Split(command, " && ")
				if len(parts) > 1 {
					command = parts[len(parts)-1] // Get the last command
				}
			}

			commands = append(commands, commandRef{
				command:     command,
				description: step.Action,
			})
		}
	}

	// Add common commands based on template
	template := summary.ProjectInfo.Template
	if template == "typescript" || template == "javascript" {
		if !f.hasCommand(commands, "npm run build") {
			commands = append(commands, commandRef{
				command:     "npm run build",
				description: "Build for production",
			})
		}
		if summary.ProjectInfo.Features["unit_testing"] && !f.hasCommand(commands, "npm test") {
			commands = append(commands, commandRef{
				command:     "npm test",
				description: "Run test suite",
			})
		}
		if summary.ProjectInfo.Features["linting"] && !f.hasCommand(commands, "npm run lint") {
			commands = append(commands, commandRef{
				command:     "npm run lint",
				description: "Check code quality",
			})
		}
	}

	return commands[:min(len(commands), 6)] // Limit to 6 commands
}

func (f *StandardFormatter) hasCommand(commands []commandRef, command string) bool {
	for _, cmd := range commands {
		if cmd.command == command {
			return true
		}
	}
	return false
}

func (f *StandardFormatter) getStandardResources(template string) []string {
	switch template {
	case "typescript":
		return []string{
			"• React Documentation: https://react.dev",
			"• TypeScript Handbook: https://typescriptlang.org/docs",
			"• Modern React Patterns: https://react.dev/learn",
		}
	case "javascript":
		return []string{
			"• React Documentation: https://react.dev",
			"• JavaScript Guide: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide",
			"• Modern React Patterns: https://react.dev/learn",
		}
	default:
		return []string{
			"• React Documentation: https://react.dev",
			"• Getting Started Guide: https://react.dev/learn",
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}