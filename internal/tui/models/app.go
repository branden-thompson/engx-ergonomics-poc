package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/bthompso/engx-ergonomics-poc/internal/config"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/styles"
	progresssim "github.com/bthompso/engx-ergonomics-poc/internal/simulation/progress"
)

// AppState represents the current state of the application
type AppState int

const (
	StateIdle AppState = iota
	StatePrompting    // NEW: Interactive configuration prompts
	StateValidating   // NEW: Validate user configuration
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

	// NEW: Prompting and configuration
	userConfig         *config.UserConfiguration
	promptOrchestrator PromptOrchestrator
	skipPrompts        bool

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

	// Check if we should skip prompts (flags provided)
	skipPrompts := hasConfigurationFlags(flags)

	// Initialize prompt orchestrator
	promptOrchestrator := NewPromptOrchestrator(target)

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

		// Define step names for npm-style renderer (will be updated based on config)
		stepNames = []string{
			"Validating configuration",
			"Setting up environment",
			"Installing dependencies",
			"Generating project structure",
		}

		if !devOnly {
			stepNames = append(stepNames, "Configuring production setup")
		}

		stepNames = append(stepNames, "Finalizing Setup")
	}

	// Create enhanced renderer (will be updated with user config later)
	appName := fmt.Sprintf("React application '%s'", target)
	targetDir := fmt.Sprintf("./%s", target)
	template := getTemplateFromFlags(flags)
	renderer := components.NewEnhancedRenderer(appName, targetDir, template, stepNames, devOnly)

	// Initialize with empty logs - all info is shown in the template
	initialLogs := []string{}

	return &AppModel{
		state:              StateIdle,
		command:            command,
		target:             target,
		flags:              flags,
		skipPrompts:        skipPrompts,
		promptOrchestrator: promptOrchestrator,
		spinner:            s,
		renderer:           renderer,
		tracker:            tracker,
		startTime:          time.Now(),
		logs:               initialLogs,
		completed:          false,
	}
}

// NewAppModelWithConfig creates a new app model with pre-configured settings from inline prompts
func NewAppModelWithConfig(command, target string, flags []string, userConfig *config.UserConfiguration) *AppModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.InfoStyle

	// Extract devOnly flag from flags
	var devOnly bool = false
	for _, flag := range flags {
		if flag == "--dev-only" {
			devOnly = true
			break
		}
	}

	// Create progress tracker based on command and configuration
	var tracker *progresssim.Tracker
	var stepNames []string

	if command == "create" {
		tracker = progresssim.NewCreateTracker(devOnly)

		// Extract step names directly from tracker to ensure perfect alignment
		tempTracker := progresssim.NewCreateTracker(devOnly)
		tempTracker.Start()

		stepNames = make([]string, tempTracker.TotalSteps())
		for i := 0; i < tempTracker.TotalSteps(); i++ {
			// Advance to step i
			for j := 0; j < i; j++ {
				tempTracker.NextStep()
			}
			stepInfo := tempTracker.CurrentStepInfo()
			if stepInfo != nil {
				stepNames[i] = stepInfo.Name
			}
			// Reset for next iteration
			tempTracker = progresssim.NewCreateTracker(devOnly)
			tempTracker.Start()
		}
	}

	// Create enhanced renderer with user configuration
	appName := fmt.Sprintf("React application '%s'", target)
	targetDir := fmt.Sprintf("./%s", target)
	template := userConfig.Template.Type.String()
	renderer := components.NewEnhancedRenderer(appName, targetDir, template, stepNames, devOnly)

	// Initialize with empty logs - all info is shown in the template
	initialLogs := []string{}

	return &AppModel{
		state:              StateIdle,
		command:            command,
		target:             target,
		flags:              flags,
		userConfig:         userConfig,    // Pre-configured from inline prompts
		skipPrompts:        true,          // Skip TUI prompts since we already have config
		promptOrchestrator: PromptOrchestrator{}, // Empty orchestrator since we skip prompts
		spinner:            s,
		renderer:           renderer,
		tracker:            tracker,
		startTime:          time.Now(),
		logs:               initialLogs,
		completed:          false,
	}
}

// hasConfigurationFlags checks if configuration flags are provided
func hasConfigurationFlags(flags []string) bool {
	for _, flag := range flags {
		if strings.HasPrefix(flag, "--template=") ||
			strings.HasPrefix(flag, "--dev-only") ||
			strings.HasPrefix(flag, "--production") {
			return true
		}
	}
	return false
}

// Init implements tea.Model
func (m *AppModel) Init() tea.Cmd {
	if m.skipPrompts {
		// Skip prompting and go straight to execution
		return tea.Batch(
			m.spinner.Tick,
			m.startExecution(),
			m.progressTicker(),
		)
	} else {
		// Start with prompting
		m.state = StatePrompting
		return tea.Batch(
			m.spinner.Tick,
			m.promptOrchestrator.Init(),
		)
	}
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
		// Handle global key messages first
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		// Remove 'q' key handling for inline mode
		}

		// Handle state-specific key messages
		switch m.state {
		case StatePrompting, StateValidating:
			// Let the prompt orchestrator handle the keys
			orchestrator, cmd := m.promptOrchestrator.Update(msg)
			m.promptOrchestrator = orchestrator

			// Check if prompting is complete
			if m.promptOrchestrator.IsComplete() {
				m.userConfig = m.promptOrchestrator.GetConfiguration()
				m.state = StateValidating
				return m, m.validateAndStartExecution()
			}

			return m, cmd
		}

	// Handle prompting messages (like CompletePromptMsg)
	default:
		if m.state == StatePrompting || m.state == StateValidating {
			// Let the prompt orchestrator handle all messages
			orchestrator, cmd := m.promptOrchestrator.Update(msg)
			m.promptOrchestrator = orchestrator

			// Check if prompting is complete
			if m.promptOrchestrator.IsComplete() {
				m.userConfig = m.promptOrchestrator.GetConfiguration()
				m.state = StateValidating
				return m, m.validateAndStartExecution()
			}

			return m, cmd
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
			// Auto-quit in inline mode after a brief pause to show completion
			cmds = append(cmds, tea.Sequence(
				tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
					return tea.Quit()
				}),
			))
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

	// Handle prompting state
	if m.state == StatePrompting || m.state == StateValidating {
		return m.promptOrchestrator.View()
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
	case "Finalizing Setup":
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
		return styles.SuccessStyle.Render("✨ Success! Application created successfully")
	case StateError:
		return styles.ErrorStyle.Render("💡 Check output above for troubleshooting. Press Ctrl+C to exit")
	default:
		return styles.MutedStyle.Render("Press Ctrl+C to quit")
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

// validateAndStartExecution validates user configuration and starts execution
func (m *AppModel) validateAndStartExecution() tea.Cmd {
	if m.userConfig == nil {
		return tea.Cmd(func() tea.Msg {
			return ErrorMsg{Error: fmt.Errorf("no user configuration available")}
		})
	}

	// Update the renderer and tracker based on user configuration
	m.updateComponentsFromConfig()

	// Start execution
	m.state = StateExecuting
	return tea.Batch(
		m.startExecution(),
		m.progressTicker(),
	)
}

// updateComponentsFromConfig updates the renderer and tracker based on user configuration
func (m *AppModel) updateComponentsFromConfig() {
	if m.userConfig == nil {
		return
	}

	// Create tracker first to get the exact step names
	devOnly := !m.userConfig.ProductionSetup.Docker &&
		!m.userConfig.ProductionSetup.CI_CD &&
		!m.userConfig.ProductionSetup.Monitoring &&
		!m.userConfig.ProductionSetup.Analytics

	// Create new tracker with updated configuration
	tempTracker := progresssim.NewCreateTracker(devOnly)
	tempTracker.Start()

	// Extract step names directly from tracker to ensure perfect alignment
	stepNames := make([]string, tempTracker.TotalSteps())
	for i := 0; i < tempTracker.TotalSteps(); i++ {
		// Advance to step i
		for j := 0; j < i; j++ {
			tempTracker.NextStep()
		}
		stepInfo := tempTracker.CurrentStepInfo()
		if stepInfo != nil {
			stepNames[i] = stepInfo.Name
		}
		// Reset for next iteration
		tempTracker = progresssim.NewCreateTracker(devOnly)
		tempTracker.Start()
	}

	// Create new tracker with updated configuration
	m.tracker = progresssim.NewCreateTracker(devOnly)
	m.totalSteps = m.tracker.TotalSteps()

	// Create new renderer with user configuration
	appName := fmt.Sprintf("React application '%s'", m.target)
	targetDir := fmt.Sprintf("./%s", m.target)
	template := m.userConfig.Template.Type.String()
	m.renderer = components.NewEnhancedRenderer(appName, targetDir, template, stepNames, devOnly)
}