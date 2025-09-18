package progress

import (
	"time"
)

// Step represents a single step in the simulation
type Step struct {
	Name        string
	Message     string
	Duration    time.Duration
	ErrorRate   float64 // 0.0 to 1.0, probability of this step failing
	CanRetry    bool
	Description string
}

// Tracker manages the progress simulation
type Tracker struct {
	steps       []Step
	currentStep int
	startTime   time.Time
	stepStart   time.Time
	completed   bool
	failed      bool
	lastError   error
}

// NewTracker creates a new progress tracker with predefined steps
func NewTracker(steps []Step) *Tracker {
	return &Tracker{
		steps:       steps,
		currentStep: 0,
		startTime:   time.Now(),
		stepStart:   time.Now(),
		completed:   false,
		failed:      false,
	}
}

// NewCreateTracker creates a tracker for the 'create' command simulation
func NewCreateTracker(devOnly bool) *Tracker {
	steps := []Step{
		{
			Name:        "Validating configuration",
			Message:     "🔍 Checking project configuration and dependencies...",
			Duration:    time.Millisecond * 1200,
			ErrorRate:   0.05, // 5% chance of config error
			CanRetry:    true,
			Description: "Validates project name, checks for conflicts, verifies system requirements",
		},
		{
			Name:        "Setting up environment",
			Message:     "⚙️ Preparing development environment...",
			Duration:    time.Millisecond * 1800,
			ErrorRate:   0.10, // 10% chance of environment error
			CanRetry:    true,
			Description: "Creates project directory, sets up git repository, configures development tools",
		},
		{
			Name:        "Installing dependencies",
			Message:     "📦 Installing React and core dependencies...",
			Duration:    time.Millisecond * 3000,
			ErrorRate:   0.15, // 15% chance of network/install error
			CanRetry:    true,
			Description: "Downloads and installs React, TypeScript, testing libraries, and build tools",
		},
		{
			Name:        "Generating project structure",
			Message:     "🏗️ Creating project files and folder structure...",
			Duration:    time.Millisecond * 2200,
			ErrorRate:   0.02, // 2% chance of filesystem error
			CanRetry:    true,
			Description: "Generates source code structure, configuration files, and example components",
		},
	}

	// Add production setup step if not dev-only
	if !devOnly {
		steps = append(steps, Step{
			Name:        "Configuring production setup",
			Message:     "🚀 Setting up build pipeline and deployment configuration...",
			Duration:    time.Millisecond * 2500,
			ErrorRate:   0.08, // 8% chance of deployment config error
			CanRetry:    true,
			Description: "Configures build scripts, environment variables, and deployment targets",
		})
	}

	// Add testing frameworks step
	steps = append(steps, Step{
		Name:        "Installing Testing Frameworks",
		Message:     "🧪 Setting up testing infrastructure...",
		Duration:    time.Millisecond * 1800,
		ErrorRate:   0.05, // 5% chance of testing setup error
		CanRetry:    true,
		Description: "Installs and configures Vitest, testing utilities, and coverage tools",
	})

	// Add documentation generation step
	steps = append(steps, Step{
		Name:        "Generating Documentation",
		Message:     "📚 Creating project documentation...",
		Duration:    time.Millisecond * 1200,
		ErrorRate:   0.02, // 2% chance of documentation error
		CanRetry:    true,
		Description: "Generates README, API docs, and component documentation",
	})

	// Final step
	steps = append(steps, Step{
		Name:        "Finalizing Setup",
		Message:     "✨ Project ready for development!",
		Duration:    time.Millisecond * 800,
		ErrorRate:   0.0, // No errors on final step
		CanRetry:    false,
		Description: "Completes setup, runs initial health checks, and prepares development server",
	})

	return NewTracker(steps)
}

// Start begins the progress simulation
func (t *Tracker) Start() {
	t.startTime = time.Now()
	t.stepStart = time.Now()
	t.currentStep = 0
}

// CurrentStep returns the current step number (0-based)
func (t *Tracker) CurrentStep() int {
	return t.currentStep
}

// TotalSteps returns the total number of steps
func (t *Tracker) TotalSteps() int {
	return len(t.steps)
}

// GetStep returns the step at the given index
func (t *Tracker) GetStep(index int) *Step {
	if index < 0 || index >= len(t.steps) {
		return nil
	}
	return &t.steps[index]
}

// CurrentStepInfo returns information about the current step
func (t *Tracker) CurrentStepInfo() *Step {
	if t.currentStep >= len(t.steps) {
		return nil
	}
	return &t.steps[t.currentStep]
}

// Progress returns the current progress as a percentage (0.0 to 1.0)
func (t *Tracker) Progress() float64 {
	if len(t.steps) == 0 {
		return 1.0
	}

	stepProgress := float64(t.currentStep) / float64(len(t.steps))

	// Add partial progress for current step based on elapsed time
	if t.currentStep < len(t.steps) && !t.completed && !t.failed {
		currentStep := &t.steps[t.currentStep]
		elapsed := time.Since(t.stepStart)
		stepPartial := float64(elapsed) / float64(currentStep.Duration)
		if stepPartial > 1.0 {
			stepPartial = 1.0
		}
		stepProgress += stepPartial / float64(len(t.steps))
	}

	if stepProgress > 1.0 {
		stepProgress = 1.0
	}

	return stepProgress
}

// IsStepReady returns true if the current step should complete
func (t *Tracker) IsStepReady() bool {
	if t.currentStep >= len(t.steps) || t.completed || t.failed {
		return false
	}

	currentStep := &t.steps[t.currentStep]
	return time.Since(t.stepStart) >= currentStep.Duration
}

// NextStep advances to the next step
func (t *Tracker) NextStep() bool {
	if t.currentStep >= len(t.steps) {
		t.completed = true
		return false
	}

	t.currentStep++
	t.stepStart = time.Now()

	if t.currentStep >= len(t.steps) {
		t.completed = true
		return false
	}

	return true
}

// IsCompleted returns true if all steps are finished
func (t *Tracker) IsCompleted() bool {
	return t.completed
}

// IsFailed returns true if the simulation failed
func (t *Tracker) IsFailed() bool {
	return t.failed
}

// GetError returns the last error that occurred
func (t *Tracker) GetError() error {
	return t.lastError
}

// EstimatedTimeRemaining calculates the estimated time to completion
func (t *Tracker) EstimatedTimeRemaining() time.Duration {
	if t.completed || t.failed {
		return 0
	}

	var remaining time.Duration
	for i := t.currentStep; i < len(t.steps); i++ {
		remaining += t.steps[i].Duration
	}

	// Subtract elapsed time from current step
	if t.currentStep < len(t.steps) {
		elapsed := time.Since(t.stepStart)
		stepDuration := t.steps[t.currentStep].Duration
		if elapsed < stepDuration {
			remaining -= elapsed
		}
	}

	return remaining
}

// TotalElapsed returns the total time elapsed since start
func (t *Tracker) TotalElapsed() time.Duration {
	return time.Since(t.startTime)
}

// GetStepStart returns the start time of the current step
func (t *Tracker) GetStepStart() time.Time {
	return t.stepStart
}

// GetSteps returns a copy of all steps
func (t *Tracker) GetSteps() []Step {
	return append([]Step(nil), t.steps...)
}

// Reset resets the tracker to the beginning
func (t *Tracker) Reset() {
	t.currentStep = 0
	t.startTime = time.Now()
	t.stepStart = time.Now()
	t.completed = false
	t.failed = false
	t.lastError = nil
}