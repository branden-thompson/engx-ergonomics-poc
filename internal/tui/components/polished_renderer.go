package components

import (
	"fmt"
	"strings"
	"time"
)

// PolishedRenderer renders progress in a clean, aligned format
type PolishedRenderer struct {
	steps         []Step
	currentStep   int
	appName       string
	startTime     time.Time

	// Layout configuration
	stepNameWidth    int
	progressBarWidth int
	totalWidth       int
}

// NewPolishedRenderer creates a new polished renderer
func NewPolishedRenderer(appName string, stepNames []string) *PolishedRenderer {
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

	return &PolishedRenderer{
		steps:            steps,
		currentStep:      0,
		appName:          appName,
		startTime:        time.Now(),
		stepNameWidth:    maxStepNameLength + 2, // Add padding
		progressBarWidth: 28, // Fixed width like in your example
		totalWidth:      80,
	}
}

// UpdateStep updates the current step's progress and status
// Does not overwrite steps that are already marked as complete
func (r *PolishedRenderer) UpdateStep(stepIndex int, progress float64, message string, subSteps []string) {
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
func (r *PolishedRenderer) SetCurrentStep(stepIndex int) {
	r.currentStep = stepIndex
	if stepIndex >= 0 && stepIndex < len(r.steps) {
		r.steps[stepIndex].Status = StepRunning
	}
}

// CompleteStep marks a step as complete
func (r *PolishedRenderer) CompleteStep(stepIndex int, duration time.Duration) {
	if stepIndex >= 0 && stepIndex < len(r.steps) {
		r.steps[stepIndex].Status = StepComplete
		r.steps[stepIndex].Progress = 1.0
		r.steps[stepIndex].Duration = duration
	}
}

// GetOverallProgress calculates overall progress (0.0 to 1.0)
func (r *PolishedRenderer) GetOverallProgress() float64 {
	if len(r.steps) == 0 {
		return 0.0
	}

	total := 0.0
	for _, step := range r.steps {
		total += step.Progress
	}
	return total / float64(len(r.steps))
}

// Render generates the polished aligned output
func (r *PolishedRenderer) Render(width int) string {
	var output strings.Builder

	// Store the width for consistent formatting
	r.totalWidth = width

	// Header section
	output.WriteString(r.renderHeader())
	output.WriteString("\n")

	// Current step info (if running or if all complete)
	allComplete := r.GetOverallProgress() >= 1.0
	if r.currentStep >= 0 && r.currentStep < len(r.steps) {
		currentStepInfo := r.steps[r.currentStep]
		if currentStepInfo.Status == StepRunning || allComplete {
			output.WriteString(r.renderCurrentStepInfo(currentStepInfo))
			output.WriteString("\n")
		}
	}

	// Separator
	output.WriteString(strings.Repeat("-", r.totalWidth))
	output.WriteString("\n")

	// Steps section
	for i, step := range r.steps {
		stepLine := r.renderStepLine(i, step)
		output.WriteString(stepLine)
		output.WriteString("\n")
	}

	// Final separator
	output.WriteString(strings.Repeat("-", r.totalWidth))

	return output.String()
}

// renderHeader creates the header with app name and total progress
func (r *PolishedRenderer) renderHeader() string {
	headerText := fmt.Sprintf("---- Creating %s ", r.appName)
	padding := r.totalWidth - len(headerText)
	if padding > 0 {
		headerText += strings.Repeat("-", padding)
	}

	// Total progress line
	overallProgress := r.GetOverallProgress()
	totalProgressBar := r.renderProgressBar(overallProgress, 40) // Wider bar for total progress
	progressText := fmt.Sprintf("Total Progress: %s %5.1f%%", totalProgressBar, overallProgress*100)

	return fmt.Sprintf("%s\n%s", headerText, progressText)
}

// renderCurrentStepInfo shows the current running step with spinner or completion
func (r *PolishedRenderer) renderCurrentStepInfo(step Step) string {
	// Check if all steps are complete
	allComplete := r.GetOverallProgress() >= 1.0

	var message, statusText string

	if allComplete {
		// Show completion state
		message = "Completed Successfully"
		statusText = "✓ Done"
	} else {
		// Show running state with spinner
		spinnerChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		spinnerIndex := int(time.Since(r.startTime)/time.Millisecond/100) % len(spinnerChars)
		spinner := spinnerChars[spinnerIndex]

		message = step.Message
		if message == "" {
			message = fmt.Sprintf("Processing %s...", step.Name)
		}
		statusText = fmt.Sprintf("%s Running...", spinner)
	}

	// Calculate spacing to right-align the status
	maxMessageLength := r.totalWidth - len("Current Step: ") - len(statusText) - 3

	if len(message) > maxMessageLength {
		message = message[:maxMessageLength-3] + "..."
	}

	padding := maxMessageLength - len(message)
	return fmt.Sprintf("Current Step: %s%s %s", message, strings.Repeat(" ", padding), statusText)
}

// renderStepLine creates a single aligned step line
func (r *PolishedRenderer) renderStepLine(index int, step Step) string {
	// Status icon
	var icon string
	switch step.Status {
	case StepPending:
		icon = "[ ]"
	case StepRunning:
		icon = "[⚡]"
	case StepComplete:
		icon = "[✓]"
	case StepError:
		icon = "[✗]"
	}

	// Step name with fixed width
	stepName := step.Name
	if len(stepName) > r.stepNameWidth {
		stepName = stepName[:r.stepNameWidth-3] + "..."
	}

	// Pad step name to fixed width
	stepNamePadded := fmt.Sprintf("%-*s", r.stepNameWidth, stepName)

	// Progress bar
	progressBar := r.renderProgressBar(step.Progress, r.progressBarWidth)

	// Progress percentage
	progressPercent := fmt.Sprintf("%5.1f%%", step.Progress*100)

	// Combine with proper spacing
	return fmt.Sprintf("%s %s %s %s", icon, stepNamePadded, progressBar, progressPercent)
}


// renderProgressBar creates a progress bar with the specified width
func (r *PolishedRenderer) renderProgressBar(progress float64, width int) string {
	filled := int(progress * float64(width))
	empty := width - filled

	if filled < 0 {
		filled = 0
	}
	if empty < 0 {
		empty = 0
	}

	bar := strings.Repeat("#", filled) + strings.Repeat(" ", empty)
	return fmt.Sprintf("[%s]", bar)
}

// GetStepCount returns the number of steps
func (r *PolishedRenderer) GetStepCount() int {
	return len(r.steps)
}

// GetStepAtIndex returns the step at the given index for debugging
func (r *PolishedRenderer) GetStepAtIndex(index int) *Step {
	if index >= 0 && index < len(r.steps) {
		return &r.steps[index]
	}
	return nil
}