package components

import (
	"fmt"
	"strings"
	"time"
)

// ANSI color codes for styling
const (
	// Basic colors
	colorReset         = "\033[0m"
	colorWhite         = "\033[97m"  // Bright white
	colorBlue          = "\033[94m"  // Blue for running/installing
	colorYellow        = "\033[93m"  // Yellow for stalling/paused/JavaScript
	colorRed           = "\033[91m"  // Red for failed
	colorGreen         = "\033[92m"  // Green for done/installed
	colorLightGrey     = "\033[90m"  // Darker grey for lines
	colorGrey          = "\033[37m"  // Light grey for skipped items

	// Header-specific colors
	colorBrightMagenta = "\033[95m"  // Bright magenta for app name/paths
	colorBrightOrange  = "\033[38;5;208m"  // Bright orange for PRODUCTION READY

	// Text styles
	styleBold   = "\033[1m"
	styleItalic = "\033[3m"
	styleReset  = "\033[22m"
)

// Helper function to create colored separator lines
func (r *EnhancedRenderer) renderSeparatorLine() string {
	return fmt.Sprintf("%s%s%s", colorLightGrey, strings.Repeat("-", r.totalWidth), colorReset)
}

// EnhancedRenderer renders progress in the comprehensive template format
type EnhancedRenderer struct {
	steps         []Step
	currentStep   int
	appName       string
	targetDir     string
	template      string
	startTime     time.Time
	isDevOnly     bool

	// Layout configuration
	stepNameWidth    int
	progressBarWidth int
	totalWidth       int

	// Component sections
	coreTechnologies []Component
	engxIntegrations []Component
	qualityComponents []QualityComponent

	// Component management
	componentManager *ComponentManager
}

// Component represents any technology component with status
type Component struct {
	Name   string
	Status string // "installed", "installing", "queued"
	Icon   string // "[✓]", "[✓ ]", "[ ]"
}

// QualityComponent represents a quality/testing tool
type QualityComponent struct {
	Name   string
	Status string // "queued", "installing", "complete"
}

// NewEnhancedRenderer creates a new enhanced renderer with comprehensive layout
func NewEnhancedRenderer(appName, targetDir, template string, stepNames []string, isDevOnly bool) *EnhancedRenderer {
	steps := make([]Step, len(stepNames))
	for i, name := range stepNames {
		steps[i] = Step{
			Name:     name,
			Status:   StepPending,
			Progress: 0.0,
			SubSteps: make([]string, 0),
		}
	}

	// Calculate optimal layout dimensions
	maxStepNameLength := 0
	for _, name := range stepNames {
		if len(name) > maxStepNameLength {
			maxStepNameLength = len(name)
		}
	}

	// Initialize component sections (all start as queued)
	coreTechnologies := []Component{
		{"TypeScript", "queued", "[ ]"},
		{"React", "queued", "[ ]"},
		{"React Router 7", "queued", "[ ]"},
		{"Tailwind CSS", "queued", "[ ]"},
		{"Radix UI", "queued", "[ ]"},
		{"ShadCN-based UI Design System (SUDS)", "queued", "[ ]"},
	}

	engxIntegrations := []Component{
		{"TrustBridge SSO", "queued", "[ ]"},
		{"gRPC Web", "queued", "[ ]"},
		{"GRID/HDFS Access", "queued", "[ ]"},
		{"CREWS API", "queued", "[ ]"},
		{"LI CATALOG API", "queued", "[ ]"},
		{"GitHub Actions", "queued", "[ ]"},
	}

	qualityComponents := []QualityComponent{
		{"Vitest", "queued"},
		{"EngX TypeScript Linters", "queued"},
		{"GitHub Pages", "queued"},
		{"StoryBook (UI Components & Documentation)", "queued"},
	}

	return &EnhancedRenderer{
		steps:             steps,
		currentStep:       0,
		appName:           appName,
		targetDir:         targetDir,
		template:          template,
		startTime:         time.Now(),
		isDevOnly:         isDevOnly,
		stepNameWidth:     42, // Fixed width for alignment
		progressBarWidth:  30, // Fixed width like in template
		totalWidth:        89, // Match template width
		coreTechnologies:  coreTechnologies,
		engxIntegrations:  engxIntegrations,
		qualityComponents: qualityComponents,
		componentManager:  NewComponentManager(),
	}
}

// UpdateStep updates the current step's progress and status
func (r *EnhancedRenderer) UpdateStep(stepIndex int, progress float64, message string, subSteps []string) {
	if stepIndex >= 0 && stepIndex < len(r.steps) {
		// Don't overwrite completed steps
		if r.steps[stepIndex].Status == StepComplete {
			return
		}

		r.steps[stepIndex].Progress = progress
		r.steps[stepIndex].Message = message
		r.steps[stepIndex].SubSteps = subSteps

		if progress >= 1.0 {
			r.steps[stepIndex].Status = StepComplete
			r.steps[stepIndex].Duration = time.Since(r.startTime)
		} else if progress > 0 {
			r.steps[stepIndex].Status = StepRunning
		}
	}
}

// SetCurrentStep sets which step is currently active
func (r *EnhancedRenderer) SetCurrentStep(stepIndex int) {
	r.currentStep = stepIndex
	if stepIndex >= 0 && stepIndex < len(r.steps) {
		r.steps[stepIndex].Status = StepRunning
	}
}

// CompleteStep marks a step as complete
func (r *EnhancedRenderer) CompleteStep(stepIndex int, duration time.Duration) {
	if stepIndex >= 0 && stepIndex < len(r.steps) {
		r.steps[stepIndex].Status = StepComplete
		r.steps[stepIndex].Progress = 1.0
		r.steps[stepIndex].Duration = duration
	}
}

// GetOverallProgress calculates overall progress (0.0 to 1.0)
func (r *EnhancedRenderer) GetOverallProgress() float64 {
	if len(r.steps) == 0 {
		return 0.0
	}

	total := 0.0
	for _, step := range r.steps {
		total += step.Progress
	}
	return total / float64(len(r.steps))
}

// Render generates the comprehensive enhanced output
func (r *EnhancedRenderer) Render(width int) string {
	var output strings.Builder

	// Store the width for consistent formatting
	r.totalWidth = width

	// Header section (includes empty line between header and progress)
	output.WriteString(r.renderHeader())
	output.WriteString("\n")

	// Current step info
	allComplete := r.GetOverallProgress() >= 1.0
	if r.currentStep >= 0 && r.currentStep < len(r.steps) {
		currentStepInfo := r.steps[r.currentStep]
		if currentStepInfo.Status == StepRunning || allComplete {
			output.WriteString(r.renderCurrentStepInfo(currentStepInfo))
			output.WriteString("\n")
		}
	}

	// Empty line before separator
	output.WriteString("\n")

	// Top separator
	output.WriteString(r.renderSeparatorLine())
	output.WriteString("\n")

	// Main steps section
	for i, step := range r.steps {
		stepLine := r.renderStepLine(i, step)
		output.WriteString(stepLine)
		output.WriteString("\n")
	}

	// Middle separator
	output.WriteString(r.renderSeparatorLine())
	output.WriteString("\n")

	// Footer info section
	output.WriteString(r.renderFooterInfo())
	output.WriteString("\n")

	// Bottom separator
	output.WriteString(r.renderSeparatorLine())
	output.WriteString("\n")

	// Empty line for breathing space (as shown in template)
	output.WriteString("\n")

	// Application Components section
	output.WriteString(r.renderApplicationComponents())

	// Final separator
	output.WriteString(r.renderSeparatorLine())

	return output.String()
}

// renderHeader creates the header with app name and total progress
func (r *EnhancedRenderer) renderHeader() string {
	// Determine setup type and colors based on devOnly flag
	var setupType, setupColor string
	if r.isDevOnly {
		setupType = "DEV SETUP"
		setupColor = colorBlue
	} else {
		setupType = "PRODUCTION READY SETUP"
		setupColor = colorBrightOrange
	}

	// Create colored header components
	dashPrefix := fmt.Sprintf("%s---%s", colorLightGrey, colorReset)
	coloredAppName := fmt.Sprintf("%s'%s'%s", colorBrightMagenta, r.appName, colorReset)
	creatingText := fmt.Sprintf(" Creating %s ", coloredAppName)

	// Calculate padding for setup type (accounting for color codes in length calculation)
	plainHeader := fmt.Sprintf("--- Creating '%s' ", r.appName)
	plainSetupPadding := fmt.Sprintf(" %s ", setupType)
	plainTotalLength := len(plainHeader) + len(plainSetupPadding) + 4 // 4 dashes at end

	var headerText string
	if plainTotalLength < r.totalWidth {
		middlePadding := r.totalWidth - len(plainHeader) - len(plainSetupPadding) - 4
		middleDashes := fmt.Sprintf("%s%s%s", colorLightGrey, strings.Repeat("-", middlePadding), colorReset)
		coloredSetupType := fmt.Sprintf(" %s%s%s ", setupColor, setupType, colorReset)
		endDashes := fmt.Sprintf("%s----%s", colorLightGrey, colorReset)

		headerText = dashPrefix + creatingText + middleDashes + coloredSetupType + endDashes
	} else {
		coloredSetupType := fmt.Sprintf(" %s%s%s ", setupColor, setupType, colorReset)
		endDashes := fmt.Sprintf("%s----%s", colorLightGrey, colorReset)
		headerText = dashPrefix + creatingText + coloredSetupType + endDashes
	}

	// Total progress line with colored progress bar and percentage
	overallProgress := r.GetOverallProgress()

	// Determine overall progress state
	var progressState ProgressState
	if overallProgress >= 1.0 {
		progressState = StateDone
	} else if overallProgress > 0 {
		progressState = StateRunning
	} else {
		progressState = StateQueued
	}

	totalProgressBar := r.renderProgressBarWithState(overallProgress, 47, progressState)
	overallProgressPercent := r.renderColoredPercentage(overallProgress, progressState)

	progressText := fmt.Sprintf("Total Progress: %s %6s", totalProgressBar, overallProgressPercent)

	return fmt.Sprintf("%s\n\n%s", headerText, progressText)
}

// renderCurrentStepInfo shows the current running step with colored spinner or completion
func (r *EnhancedRenderer) renderCurrentStepInfo(step Step) string {
	// Check if all steps are complete
	allComplete := r.GetOverallProgress() >= 1.0

	var message, statusText string

	if allComplete {
		// Show completion state with green color
		message = "Completed Successfully"
		statusText = fmt.Sprintf("%s✓ Done%s", colorGreen, colorReset)
	} else {
		// Show running state with colored spinner
		spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		spinnerIndex := int(time.Since(r.startTime)/time.Millisecond/100) % len(spinnerChars)
		spinner := spinnerChars[spinnerIndex]

		message = step.Message
		if message == "" {
			message = step.Name
		}

		// Determine spinner color based on step status
		var spinnerColor string
		switch step.Status {
		case StepPending:
			spinnerColor = colorWhite // Queued
		case StepRunning:
			spinnerColor = colorBlue // Running
		case StepComplete:
			spinnerColor = colorGreen // Done
		case StepError:
			spinnerColor = colorRed // Failed
		default:
			// Determine based on progress
			if step.Progress >= 1.0 {
				spinnerColor = colorGreen
			} else if step.Progress > 0 {
				spinnerColor = colorBlue
			} else {
				spinnerColor = colorWhite
			}
		}

		coloredSpinner := fmt.Sprintf("%s%s%s", spinnerColor, spinner, colorReset)
		statusText = fmt.Sprintf("%s Running...", coloredSpinner)
	}

	// Calculate spacing to right-align the status (use plain text lengths for calculation)
	plainStatusText := "✓ Done"   // Sample status for length calculation
	maxMessageLength := r.totalWidth - len("Current Step: ") - len(plainStatusText) - 3

	if len(message) > maxMessageLength {
		message = message[:maxMessageLength-3] + "..."
	}

	padding := maxMessageLength - len(message)
	return fmt.Sprintf("Current Step: %s%s %s", message, strings.Repeat(" ", padding), statusText)
}

// renderStepLine creates a single aligned step line with dynamic width
func (r *EnhancedRenderer) renderStepLine(index int, step Step) string {
	// Use modular status icon system
	iconType := statusIconTypeFromStepStatus(step.Status)
	icon := r.renderStatusIcon(iconType)

	// Use modular progress state system
	progressState := progressStateFromStepStatus(step.Status, step.Progress)

	// Calculate dynamic progress bar width (account for plain icon text for spacing)
	plainIcon := "[ ]" // Use plain icon for width calculation
	plainProgressPercent := " 0.0%" // Use sample percentage for width calculation
	remainingWidth := r.totalWidth - len(plainIcon) - 1 - 1 - len(plainProgressPercent)

	// Use 30 chars for progress bar, rest for step name
	progressBarWidth := 30
	stepNameWidth := remainingWidth - progressBarWidth - 1

	// Step name with dynamic width
	stepName := step.Name
	if len(stepName) > stepNameWidth {
		stepName = stepName[:stepNameWidth-3] + "..."
	}

	// Pad step name to calculated width
	stepNamePadded := fmt.Sprintf("%-*s", stepNameWidth, stepName)

	// Colored progress bar and percentage
	progressBar := r.renderProgressBarWithState(step.Progress, progressBarWidth, progressState)
	progressPercent := r.renderColoredPercentage(step.Progress, progressState)

	// Combine with proper spacing
	return fmt.Sprintf("%s %s %s %s", icon, stepNamePadded, progressBar, progressPercent)
}

// renderFooterInfo creates the footer with timing and directory info
func (r *EnhancedRenderer) renderFooterInfo() string {
	// First line: Target Directory and Template with colors
	var templateDisplay, templateColor string
	switch strings.ToLower(r.template) {
	case "typescript":
		templateDisplay = "TypeScript"
		templateColor = colorBlue // Blue for TypeScript
	case "javascript":
		templateDisplay = "JavaScript"
		templateColor = colorYellow // Yellow for JavaScript
	default:
		templateDisplay = r.template
		templateColor = colorWhite // Default white
	}

	// Color the directory path bright magenta
	coloredTargetDir := fmt.Sprintf("%s%s%s", colorBrightMagenta, r.targetDir, colorReset)
	coloredTemplate := fmt.Sprintf("%s%s%s", templateColor, templateDisplay, colorReset)

	// Calculate padding accounting for the original text length (without color codes)
	plainLine1 := fmt.Sprintf("Target Directory: %s", r.targetDir)
	plainTemplatePart := templateDisplay
	padding1 := r.totalWidth - len(plainLine1) - len(plainTemplatePart)

	var line1 string
	if padding1 > 0 {
		line1 = fmt.Sprintf("Target Directory: %s%s%s", coloredTargetDir, strings.Repeat(" ", padding1), coloredTemplate)
	} else {
		line1 = fmt.Sprintf("Target Directory: %s %s", coloredTargetDir, coloredTemplate)
	}

	// Second line: Timing information
	elapsed := time.Since(r.startTime)
	elapsedFormatted := formatDuration(elapsed)

	// Estimate remaining time based on progress
	var estimatedRemaining string
	if progress := r.GetOverallProgress(); progress > 0 && progress < 1.0 {
		totalEstimated := time.Duration(float64(elapsed) / progress)
		remaining := totalEstimated - elapsed
		estimatedRemaining = formatDuration(remaining)
	} else {
		estimatedRemaining = "00h 00m 00s"
	}

	line2Left := fmt.Sprintf("Estimated Time Remaining: %s", estimatedRemaining)
	line2Right := fmt.Sprintf("Running Time: %s", elapsedFormatted)
	padding2 := r.totalWidth - len(line2Left) - len(line2Right)
	if padding2 > 0 {
		line2 := line2Left + strings.Repeat(" ", padding2) + line2Right
		return fmt.Sprintf("%s\n%s", line1, line2)
	}

	return line1
}

// renderApplicationComponents creates the Application Components section
func (r *EnhancedRenderer) renderApplicationComponents() string {
	var output strings.Builder

	// Section header - white title with grey dashes
	headerText := "---- APPLICATION COMPONENTS "
	padding := r.totalWidth - len(headerText)
	var fullHeaderText string
	if padding > 0 {
		paddingDashes := fmt.Sprintf("%s%s%s", colorLightGrey, strings.Repeat("-", padding), colorReset)
		// Use grey for dashes but white for title
		dashPrefix := fmt.Sprintf("%s----%s", colorLightGrey, colorReset)
		whiteTitle := fmt.Sprintf("%s APPLICATION COMPONENTS %s", colorWhite, colorReset)
		fullHeaderText = dashPrefix + whiteTitle + paddingDashes
	} else {
		dashPrefix := fmt.Sprintf("%s----%s", colorLightGrey, colorReset)
		whiteTitle := fmt.Sprintf("%s APPLICATION COMPONENTS %s", colorWhite, colorReset)
		fullHeaderText = dashPrefix + whiteTitle
	}
	output.WriteString(fullHeaderText + "\n")

	// Core Technologies section
	output.WriteString("• Core Technologies:\n")
	for _, component := range r.coreTechnologies {
		output.WriteString(r.renderComponentLine(component))
	}
	output.WriteString("\n")

	// EngX Integrations section
	output.WriteString("• EngX Integrations:\n")
	for _, component := range r.engxIntegrations {
		output.WriteString(r.renderComponentLine(component))
	}
	output.WriteString("\n")

	// Quality & Testing section
	output.WriteString("• Quality & Testing:\n")
	for _, component := range r.qualityComponents {
		// Convert QualityComponent to Component for consistent rendering
		var icon string
		if component.Status == "installed" {
			icon = "[✓]"
		} else if component.Status == "installing" {
			icon = "[✓ ]"
		} else {
			icon = "[ ]"
		}
		comp := Component{Name: component.Name, Status: component.Status, Icon: icon}
		output.WriteString(r.renderComponentLine(comp))
	}

	return output.String()
}

// renderComponentLine creates a properly aligned component line
func (r *EnhancedRenderer) renderComponentLine(component Component) string {
	// Status display mapping
	var statusDisplay string
	if component.Status == "installed" {
		statusDisplay = "[installed]"
	} else if component.Status == "installing" {
		statusDisplay = "[installing...]"
	} else {
		statusDisplay = "[ queued ]"
	}

	// Calculate dynamic width for proper right alignment
	// Total width - indent - icon - spaces - status = remaining for component name
	remainingWidth := r.totalWidth - 2 - len(component.Icon) - 1 - len(statusDisplay) - 1

	// Component name with dynamic width
	componentName := component.Name
	if len(componentName) > remainingWidth {
		componentName = componentName[:remainingWidth-3] + "..."
	}

	// Pad component name to calculated width for right alignment
	componentNamePadded := fmt.Sprintf("%-*s", remainingWidth, componentName)

	// Combine with proper spacing
	return fmt.Sprintf("  %s %s %s\n", component.Icon, componentNamePadded, statusDisplay)
}

// ProgressState represents different states for progress visualization
type ProgressState int

const (
	StateQueued ProgressState = iota
	StateRunning
	StateStalling
	StateFailed
	StateDone
)

// renderProgressBar creates a colored progress bar with the specified width
func (r *EnhancedRenderer) renderProgressBar(progress float64, width int) string {
	return r.renderProgressBarWithState(progress, width, StateRunning)
}

// renderProgressBarWithState creates a colored progress bar with specific state
func (r *EnhancedRenderer) renderProgressBarWithState(progress float64, width int, state ProgressState) string {
	filled := int(progress * float64(width))
	empty := width - filled

	if filled < 0 {
		filled = 0
	}
	if empty < 0 {
		empty = 0
	}

	// Determine color based on state
	var barColor string
	switch state {
	case StateQueued:
		barColor = colorWhite
	case StateRunning:
		barColor = colorBlue
	case StateStalling:
		barColor = colorYellow
	case StateFailed:
		barColor = colorRed
	case StateDone:
		barColor = colorGreen
	default:
		// If progress is 100%, use green regardless of state
		if progress >= 1.0 {
			barColor = colorGreen
		} else if progress > 0 {
			barColor = colorBlue
		} else {
			barColor = colorWhite
		}
	}

	// Create colored progress bar
	coloredFilled := fmt.Sprintf("%s%s%s", barColor, strings.Repeat("#", filled), colorReset)
	emptySpace := strings.Repeat(" ", empty)

	return fmt.Sprintf("[%s%s]", coloredFilled, emptySpace)
}

// renderColoredPercentage colors percentage text to match progress bar state
func (r *EnhancedRenderer) renderColoredPercentage(progress float64, state ProgressState) string {
	// Calculate percentage
	percentage := progress * 100
	var percentText string
	if percentage == 0.0 {
		percentText = "0.0%"
	} else if percentage >= 100.0 {
		percentText = "100.0%"  // Ensure 100% always shows as "100.0%"
	} else {
		percentText = fmt.Sprintf("%.1f%%", percentage)
	}

	// Determine color based on state (same logic as progress bars)
	var percentColor string
	switch state {
	case StateQueued:
		percentColor = colorWhite
	case StateRunning:
		percentColor = colorBlue
	case StateStalling:
		percentColor = colorYellow
	case StateFailed:
		percentColor = colorRed
	case StateDone:
		percentColor = colorGreen
	default:
		if progress >= 1.0 {
			percentColor = colorGreen
		} else if progress > 0 {
			percentColor = colorBlue
		} else {
			percentColor = colorWhite
		}
	}

	return fmt.Sprintf("%s%s%s", percentColor, percentText, colorReset)
}

// StatusIconType represents different types of status indicators
type StatusIconType int

const (
	IconPending StatusIconType = iota
	IconRunning
	IconComplete
	IconError
	IconQueued  // Alias for pending
)

// renderStatusIcon creates a colored status icon based on type
func (r *EnhancedRenderer) renderStatusIcon(iconType StatusIconType) string {
	var icon, color string

	switch iconType {
	case IconPending, IconQueued:
		icon = "[ ]"
		color = colorWhite
	case IconRunning:
		icon = "[✓ ]" // Space after checkmark for running steps
		color = colorBlue
	case IconComplete:
		icon = "[✓]"
		color = colorGreen
	case IconError:
		icon = "[✗]"
		color = colorRed
	default:
		icon = "[ ]"
		color = colorWhite
	}

	return fmt.Sprintf("%s%s%s", color, icon, colorReset)
}

// statusIconTypeFromStepStatus converts step status to icon type
func statusIconTypeFromStepStatus(status StepStatus) StatusIconType {
	switch status {
	case StepPending:
		return IconPending
	case StepRunning:
		return IconRunning
	case StepComplete:
		return IconComplete
	case StepError:
		return IconError
	default:
		return IconPending
	}
}

// progressStateFromStepStatus converts step status to progress state
func progressStateFromStepStatus(status StepStatus, progress float64) ProgressState {
	switch status {
	case StepPending:
		return StateQueued
	case StepRunning:
		return StateRunning
	case StepComplete:
		return StateDone
	case StepError:
		return StateFailed
	default:
		// Fallback based on progress value
		if progress >= 1.0 {
			return StateDone
		} else if progress > 0 {
			return StateRunning
		} else {
			return StateQueued
		}
	}
}

// formatDuration formats a duration into HH:MM:SS format
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02dh %02dm %02ds", hours, minutes, seconds)
}

// GetStepCount returns the number of steps
func (r *EnhancedRenderer) GetStepCount() int {
	return len(r.steps)
}

// GetStepAtIndex returns the step at the given index for debugging
func (r *EnhancedRenderer) GetStepAtIndex(index int) *Step {
	if index >= 0 && index < len(r.steps) {
		return &r.steps[index]
	}
	return nil
}

// UpdateComponentStatuses updates component statuses based on step progress using the modular system
func (r *EnhancedRenderer) UpdateComponentStatuses(stepName string, stepProgress float64) {
	// Map step name to installation phase
	phase := MapStepNameToPhase(stepName)

	// Get all component updates up to this phase and progress
	// This ensures previous phase components remain installed and don't re-animate
	updates := r.componentManager.GetAllComponentsUpToPhase(phase, stepProgress)

	// Apply updates to components (with protection against re-animation)
	for _, update := range updates {
		r.applyComponentUpdateSafe(update)
	}

	// Special case: Finalizing step - ensure everything is installed
	if (stepName == "Finalizing Setup" || stepName == "Finalizing setup") && stepProgress >= 1.0 {
		allUpdates := r.componentManager.GetAllComponentsForPhase(PhaseFinalizing)
		for _, update := range allUpdates {
			r.applyComponentUpdateSafe(update)
		}
	}
}

// applyComponentUpdate applies a component update to the appropriate section
func (r *EnhancedRenderer) applyComponentUpdate(update ComponentUpdate) {
	// Try to update in core technologies
	for i := range r.coreTechnologies {
		if r.coreTechnologies[i].Name == update.ComponentName {
			r.coreTechnologies[i].Status = update.NewStatus
			r.coreTechnologies[i].Icon = update.NewIcon
			return
		}
	}

	// Try to update in engx integrations
	for i := range r.engxIntegrations {
		if r.engxIntegrations[i].Name == update.ComponentName {
			r.engxIntegrations[i].Status = update.NewStatus
			r.engxIntegrations[i].Icon = update.NewIcon
			return
		}
	}

	// Try to update in quality components
	for i := range r.qualityComponents {
		if r.qualityComponents[i].Name == update.ComponentName {
			r.qualityComponents[i].Status = update.NewStatus
			return
		}
	}
}

// applyComponentUpdateSafe applies a component update but prevents downgrading from installed to installing
func (r *EnhancedRenderer) applyComponentUpdateSafe(update ComponentUpdate) {
	// Try to update in core technologies
	for i := range r.coreTechnologies {
		if r.coreTechnologies[i].Name == update.ComponentName {
			// Don't downgrade from installed to installing
			if r.coreTechnologies[i].Status == "installed" && update.NewStatus == "installing" {
				return
			}
			r.coreTechnologies[i].Status = update.NewStatus
			r.coreTechnologies[i].Icon = update.NewIcon
			return
		}
	}

	// Try to update in engx integrations
	for i := range r.engxIntegrations {
		if r.engxIntegrations[i].Name == update.ComponentName {
			// Don't downgrade from installed to installing
			if r.engxIntegrations[i].Status == "installed" && update.NewStatus == "installing" {
				return
			}
			r.engxIntegrations[i].Status = update.NewStatus
			r.engxIntegrations[i].Icon = update.NewIcon
			return
		}
	}

	// Try to update in quality components
	for i := range r.qualityComponents {
		if r.qualityComponents[i].Name == update.ComponentName {
			// Don't downgrade from installed to installing
			if r.qualityComponents[i].Status == "installed" && update.NewStatus == "installing" {
				return
			}
			r.qualityComponents[i].Status = update.NewStatus
			return
		}
	}
}

// updateComponentStatus updates a specific component's status
func (r *EnhancedRenderer) updateComponentStatus(section, name, status, icon string) {
	if section == "coreTech" {
		for i := range r.coreTechnologies {
			if r.coreTechnologies[i].Name == name {
				r.coreTechnologies[i].Status = status
				r.coreTechnologies[i].Icon = icon
				break
			}
		}
	} else if section == "engx" {
		for i := range r.engxIntegrations {
			if r.engxIntegrations[i].Name == name {
				r.engxIntegrations[i].Status = status
				r.engxIntegrations[i].Icon = icon
				break
			}
		}
	}
}

// updateQualityComponentStatus updates a quality component's status
func (r *EnhancedRenderer) updateQualityComponentStatus(name, status string) {
	for i := range r.qualityComponents {
		if r.qualityComponents[i].Name == name {
			r.qualityComponents[i].Status = status
			break
		}
	}
}