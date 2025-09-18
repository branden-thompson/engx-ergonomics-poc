package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/styles"
)

// StepStatus represents the status of a step
type StepStatus int

const (
	StepPending StepStatus = iota
	StepRunning
	StepComplete
	StepError
)

// String returns string representation of StepStatus
func (s StepStatus) String() string {
	switch s {
	case StepPending:
		return "Pending"
	case StepRunning:
		return "Running"
	case StepComplete:
		return "Complete"
	case StepError:
		return "Error"
	default:
		return "Unknown"
	}
}

// Step represents a single step in the process
type Step struct {
	Name         string
	Status       StepStatus
	Progress     float64 // 0.0 to 1.0
	Message      string
	Duration     time.Duration
	SubSteps     []string
	Error        error
}

// NPMStyleRenderer renders progress in npm/yarn style
type NPMStyleRenderer struct {
	steps         []Step
	currentStep   int
	overallTitle  string
	startTime     time.Time
}

// NewNPMStyleRenderer creates a new npm-style renderer
func NewNPMStyleRenderer(title string, stepNames []string) *NPMStyleRenderer {
	steps := make([]Step, len(stepNames))
	for i, name := range stepNames {
		steps[i] = Step{
			Name:     name,
			Status:   StepPending,
			Progress: 0.0,
			SubSteps: make([]string, 0),
		}
	}

	return &NPMStyleRenderer{
		steps:        steps,
		currentStep:  0,
		overallTitle: title,
		startTime:    time.Now(),
	}
}

// UpdateStep updates the current step's progress and status
// Does not overwrite steps that are already marked as complete
func (r *NPMStyleRenderer) UpdateStep(stepIndex int, progress float64, message string, subSteps []string) {
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
func (r *NPMStyleRenderer) SetCurrentStep(stepIndex int) {
	r.currentStep = stepIndex
	if stepIndex >= 0 && stepIndex < len(r.steps) {
		r.steps[stepIndex].Status = StepRunning
	}
}

// CompleteStep marks a step as complete
func (r *NPMStyleRenderer) CompleteStep(stepIndex int, duration time.Duration) {
	if stepIndex >= 0 && stepIndex < len(r.steps) {
		r.steps[stepIndex].Status = StepComplete
		r.steps[stepIndex].Progress = 1.0
		r.steps[stepIndex].Duration = duration
	}
}

// GetOverallProgress calculates overall progress (0.0 to 1.0)
func (r *NPMStyleRenderer) GetOverallProgress() float64 {
	if len(r.steps) == 0 {
		return 0.0
	}

	total := 0.0
	for _, step := range r.steps {
		total += step.Progress
	}
	return total / float64(len(r.steps))
}

// Render generates the npm/yarn-style output
func (r *NPMStyleRenderer) Render(width int) string {
	var output strings.Builder

	// Header with overall progress
	overallProgress := r.GetOverallProgress()
	progressBar := r.renderProgressBar(overallProgress, 40)

	headerStyle := styles.HeaderStyle
	titleLine := headerStyle.Render(fmt.Sprintf("🛩️ %s", r.overallTitle))
	progressLine := fmt.Sprintf("%s %.1f%%", progressBar, overallProgress*100)

	output.WriteString(titleLine + "\n")
	output.WriteString(progressLine + "\n")
	output.WriteString(strings.Repeat("─", width) + "\n")

	// Render each step
	for i, step := range r.steps {
		stepLine := r.renderStep(i, step, i == r.currentStep)
		output.WriteString(stepLine + "\n")

		// Show substeps for current step
		if i == r.currentStep && len(step.SubSteps) > 0 {
			for _, subStep := range step.SubSteps {
				subStepLine := styles.MutedStyle.Render(fmt.Sprintf("  %s", subStep))
				output.WriteString(subStepLine + "\n")
			}
		}
	}

	return output.String()
}

// renderStep renders a single step line
func (r *NPMStyleRenderer) renderStep(index int, step Step, isCurrent bool) string {
	var icon string
	var style lipgloss.Style

	switch step.Status {
	case StepPending:
		icon = "⏳"
		style = styles.MutedStyle
	case StepRunning:
		icon = "⚡"
		style = styles.InfoStyle
	case StepComplete:
		icon = "✅"
		style = styles.SuccessStyle
	case StepError:
		icon = "❌"
		style = styles.ErrorStyle
	}

	stepText := fmt.Sprintf("%s %s", icon, step.Name)

	// Add step progress bar for running and completed steps
	if (step.Status == StepRunning || step.Status == StepComplete) && step.Progress > 0 {
		stepProgressBar := r.renderProgressBar(step.Progress, 20)
		stepText = fmt.Sprintf("%s %s %s %.1f%%", icon, step.Name, stepProgressBar, step.Progress*100)
	} else if step.Status == StepRunning {
		// Show 0.0% for running steps that haven't started yet
		stepProgressBar := r.renderProgressBar(0.0, 20)
		stepText = fmt.Sprintf("%s %s %s 0.0%%", icon, step.Name, stepProgressBar)
	}

	if isCurrent && step.Status == StepRunning && step.Message != "" {
		stepText += fmt.Sprintf(" - %s", step.Message)
	}

	return style.Render(stepText)
}

// renderProgressBar creates a simple progress bar
func (r *NPMStyleRenderer) renderProgressBar(progress float64, width int) string {
	filled := int(progress * float64(width))
	empty := width - filled

	if filled < 0 {
		filled = 0
	}
	if empty < 0 {
		empty = 0
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	return fmt.Sprintf("[%s]", bar)
}

// GetStepCount returns the number of steps
func (r *NPMStyleRenderer) GetStepCount() int {
	return len(r.steps)
}

// GetStepAtIndex returns the step at the given index for debugging
func (r *NPMStyleRenderer) GetStepAtIndex(index int) *Step {
	if index >= 0 && index < len(r.steps) {
		return &r.steps[index]
	}
	return nil
}