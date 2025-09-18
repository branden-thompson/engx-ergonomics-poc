package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/styles"
	progresssim "github.com/bthompso/engx-ergonomics-poc/internal/simulation/progress"
)

// AppState represents the current state of the application
type AppState int

const (
	StateIdle AppState = iota
	StatePrompt
	StateExecuting
	StateComplete
	StateError
)

// AppModel represents the main application model
type AppModel struct {
	state      AppState
	command    string
	target     string
	flags      []string

	// UI components
	spinner    spinner.Model
	renderer   *components.EnhancedRenderer

	// Progress tracking
	tracker    *progresssim.Tracker
	startTime  time.Time

	// Execution state
	currentStep   int
	totalSteps    int
	stepName      string
	logs          []string
	error         error

	// Window dimensions
	width  int
	height int

	// Completion state
	completed bool
}

// getTemplateFromFlags extracts template from flags or returns default
func getTemplateFromFlags(flags []string) string {
	for _, flag := range flags {
		if strings.HasPrefix(flag, "--template=") {
			return strings.TrimPrefix(flag, "--template=")
		}
	}
	return "typescript"
}

// NewAppModel creates a new application model
func NewAppModel(command, target string, flags []string) *AppModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.InfoStyle

	// Create progress tracker based on command
	var tracker *progresssim.Tracker
	var stepNames []string
	var devOnly bool = false

	if command == "create" {
		for _, flag := range flags {
			if flag == "--dev-only" {
				devOnly = true
				break
			}
		}
		tracker = progresssim.NewCreateTracker(devOnly)

		// Define step names for npm-style renderer
		stepNames = []string{
			"Validating configuration",
			"Setting up environment",
			"Installing dependencies",
			"Generating project structure",
		}

		if !devOnly {
			stepNames = append(stepNames, "Configuring production setup")
		}

		stepNames = append(stepNames, "Finalizing setup")
	}

	// Create enhanced renderer
	appName := fmt.Sprintf("React application '%s'", target)
	targetDir := fmt.Sprintf("./%s", target)
	template := getTemplateFromFlags(flags)
	renderer := components.NewEnhancedRenderer(appName, targetDir, template, stepNames, devOnly)

	// Initialize with empty logs - all info is shown in the template
	initialLogs := []string{}

	return &AppModel{
		state:     StateIdle,
		command:   command,
		target:    target,
		flags:     flags,
		spinner:   s,
		renderer:  renderer,
		tracker:   tracker,
		startTime: time.Now(),
		logs:      initialLogs,
		completed: false,
	}
}

// Init implements tea.Model
func (m *AppModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.startExecution(),
		m.progressTicker(),
	)
}

// Update implements tea.Model
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.state == StateComplete || m.state == StateError {
				return m, tea.Quit
			}
		}

	case ProgressMsg:
		m.currentStep = msg.Step
		m.stepName = msg.StepName
		// Skip adding logs - all information is shown in the template

		// Mark previous step as complete if we advanced
		if msg.Step > 0 && m.renderer != nil {
			m.renderer.CompleteStep(msg.Step-1, time.Since(m.startTime))
		}

		// Set current step
		if m.renderer != nil && msg.Step < m.totalSteps {
			m.renderer.SetCurrentStep(msg.Step)
		}

		if msg.Step >= m.totalSteps {
			m.state = StateComplete
			m.completed = true
			// Mark final step as complete
			if m.renderer != nil {
				m.renderer.CompleteStep(msg.Step-1, time.Since(m.startTime))
			}
		} else {
			m.state = StateExecuting
		}

		cmds = append(cmds, m.nextStep())

	case ErrorMsg:
		m.state = StateError
		m.error = msg.Error
		// Skip adding error logs - errors will be shown in footer

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	// Removed old progress.FrameMsg handling for npm-style renderer

	case ProgressTickMsg:
		// Update npm-style renderer during step execution
		if m.state == StateExecuting && m.tracker != nil && !m.tracker.IsCompleted() {
			currentStep := m.tracker.CurrentStep()
			stepInfo := m.tracker.CurrentStepInfo()

			if stepInfo != nil && currentStep >= 0 {
				// Calculate individual step progress using tracker's step timing
				stepProgress := 0.0
				if currentStep < len(m.tracker.GetSteps()) {
					// Use the tracker's internal step timing
					stepStart := m.tracker.GetStepStart()
					elapsed := time.Since(stepStart)
					stepProgress = float64(elapsed) / float64(stepInfo.Duration)

					if stepProgress > 1.0 {
						stepProgress = 1.0
					}
					if stepProgress < 0 {
						stepProgress = 0.0
					}
				}

				// Update the renderer
				m.renderer.SetCurrentStep(currentStep)
				// Only update progress for steps that haven't been completed yet
				m.renderer.UpdateStep(currentStep, stepProgress, stepInfo.Message, m.getSubSteps(stepInfo.Name))
				// Update component statuses based on step progress
				m.renderer.UpdateComponentStatuses(stepInfo.Name, stepProgress)
			}
		}
		// Continue ticking for smooth animation
		cmds = append(cmds, m.progressTicker())

	case StepCheckMsg:
		// Continue checking for step advancement
		if m.state == StateExecuting && m.tracker != nil && !m.tracker.IsCompleted() {
			cmds = append(cmds, m.nextStep())
		}
	}

	// Progress bar updates are now handled in updateProgressBar method

	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m *AppModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var output strings.Builder

	// Use npm-style renderer for main progress display
	if m.renderer != nil {
		npmOutput := m.renderer.Render(m.width)
		output.WriteString(npmOutput)
		output.WriteString("\n")
	}

	// Skip logs - all information is shown in the template

	// Footer
	footer := m.renderFooter()
	if footer != "" {
		output.WriteString("\n")
		output.WriteString(footer)
	}

	return output.String()
}

// getSubSteps returns sub-steps for the current step
func (m *AppModel) getSubSteps(stepName string) []string {
	switch stepName {
	case "Validating configuration":
		return []string{
			"Checking project name validity",
			"Verifying directory permissions",
			"Scanning for naming conflicts",
			"Validating Node.js version",
		}
	case "Setting up environment":
		return []string{
			"Creating project directory structure",
			"Initializing Git repository",
			"Setting up .gitignore",
			"Configuring development environment",
		}
	case "Installing dependencies":
		return []string{
			"Installing React v18.2.0",
			"Installing TypeScript v5.1.6",
			"Installing Vite v4.4.5",
			"Installing testing dependencies",
		}
	case "Generating project structure":
		return []string{
			"Creating src/ directory",
			"Generating App.tsx component",
			"Setting up configuration files",
			"Adding package.json scripts",
		}
	case "Configuring production setup":
		return []string{
			"Setting up build optimization",
			"Configuring environment variables",
			"Preparing deployment scripts",
		}
	case "Finalizing setup":
		return []string{
			"Running initial build test",
			"Validating project structure",
			"Generating README.md",
		}
	default:
		return []string{}
	}
}

// Removed old render methods - using npm-style renderer instead

func (m *AppModel) renderFooter() string {
	switch m.state {
	case StateComplete:
		return styles.SuccessStyle.Render("✨ Success! Press Enter or 'q' to exit")
	case StateError:
		return styles.ErrorStyle.Render("💡 Check output above for troubleshooting. Press Enter or 'q' to exit")
	default:
		return styles.MutedStyle.Render("Press Ctrl+C or 'q' to quit")
	}
}

// Custom messages
type ProgressMsg struct {
	Step     int
	StepName string
	Message  string
}

type ErrorMsg struct {
	Error error
}

type ProgressTickMsg struct{}

type StepCheckMsg struct{}

// progressTicker provides continuous updates for smooth progress animation
func (m *AppModel) progressTicker() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return ProgressTickMsg{}
	})
}

// updateProgressBar is no longer needed with npm-style renderer
func (m *AppModel) updateProgressBar() tea.Cmd {
	return nil
}

// Command methods
func (m *AppModel) startExecution() tea.Cmd {
	if m.tracker != nil {
		m.tracker.Start()
		m.totalSteps = m.tracker.TotalSteps()
		m.state = StateExecuting
		return m.nextStep()
	}
	return nil
}

// Removed addStepLogs - using npm-style renderer sub-steps instead

func (m *AppModel) nextStep() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		if m.tracker == nil {
			return StepCheckMsg{}
		}

		// Check if current step should advance
		if m.tracker.IsStepReady() {
			// Advance to next step
			if !m.tracker.NextStep() {
				// All steps complete
				return ProgressMsg{
					Step:     m.tracker.TotalSteps(),
					StepName: "Complete",
					Message:  "✨ All steps completed successfully!",
				}
			}

			// Get current step info
			if stepInfo := m.tracker.CurrentStepInfo(); stepInfo != nil {
				return ProgressMsg{
					Step:     m.tracker.CurrentStep(),
					StepName: stepInfo.Name,
					Message:  stepInfo.Message,
				}
			}
		}

		// Continue checking for step completion
		return StepCheckMsg{}
	})
}