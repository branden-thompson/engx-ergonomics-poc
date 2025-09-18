package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/bthompso/engx-ergonomics-poc/internal/config"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components/prompts"
)

// PromptOrchestrator manages the flow of interactive prompts
type PromptOrchestrator struct {
	prompts       []prompts.PromptStep
	currentIndex  int
	config        *config.UserConfiguration
	navigation    prompts.NavigationState
	completed     bool
	projectName   string
}

// NewPromptOrchestrator creates a new prompt orchestrator
func NewPromptOrchestrator(projectName string) PromptOrchestrator {
	// Start with smart defaults
	config := config.GetSmartDefaults(projectName)

	// Define the prompt flow (confirmation will be added dynamically)
	promptSteps := []prompts.PromptStep{
		{
			ID:        "template",
			Title:     "Template Selection",
			Component: prompts.NewTemplateSelector(),
			Required:  true,
			HelpText:  "Choose the project template that best fits your needs",
			Type:      prompts.PromptTypeTemplate,
		},
		{
			ID:        "dev-features",
			Title:     "Development Features",
			Component: prompts.NewDevFeatureSelector(),
			Required:  false,
			HelpText:  "Select development tools to enhance your workflow",
			Type:      prompts.PromptTypeDevFeatures,
		},
		{
			ID:        "production-setup",
			Title:     "Production Setup",
			Component: prompts.NewProductionFeatureSelector(),
			Required:  false,
			HelpText:  "Configure production deployment and monitoring",
			Type:      prompts.PromptTypeProductionSetup,
		},
		{
			ID:        "testing",
			Title:     "Testing Configuration",
			Component: prompts.NewTestingFeatureSelector(),
			Required:  false,
			HelpText:  "Set up testing frameworks and tools",
			Type:      prompts.PromptTypeTesting,
		},
		{
			ID:        "navigation",
			Title:     "Navigation Configuration",
			Component: prompts.NewNavigationSelector(),
			Required:  true,
			HelpText:  "Choose between federated global navigation or standalone app navigation",
			Type:      prompts.PromptTypeNavigation,
		},
	}

	return PromptOrchestrator{
		prompts:     promptSteps,
		currentIndex: 0,
		config:      &config,
		navigation: prompts.NavigationState{
			CanGoBack:    false,
			CanGoForward: true,
			CanSkip:      false,
			ShowHelp:     false,
			CurrentStep:  1,
			TotalSteps:   len(promptSteps),
		},
		completed:   false,
		projectName: projectName,
	}
}

// Init initializes the orchestrator
func (po *PromptOrchestrator) Init() tea.Cmd {
	if len(po.prompts) > 0 {
		return po.prompts[po.currentIndex].Component.Init()
	}
	return nil
}

// Update handles orchestrator updates
func (po *PromptOrchestrator) Update(msg tea.Msg) (PromptOrchestrator, tea.Cmd) {
	if po.completed || len(po.prompts) == 0 {
		return *po, nil
	}

	currentPrompt := &po.prompts[po.currentIndex]

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return *po, tea.Quit
		case "b", "left":
			if po.navigation.CanGoBack {
				return po.goToPreviousPrompt()
			}
		case "s":
			if po.navigation.CanSkip && !currentPrompt.Required {
				return po.skipCurrentPrompt()
			}
		}

	case prompts.CompletePromptMsg:
		// Save the current prompt's value
		po.savePromptValue(currentPrompt)

		// Check if we need to add confirmation step
		if po.currentIndex == len(po.prompts)-1 && currentPrompt.Type != prompts.PromptTypeConfirmation {
			// Add confirmation step dynamically
			confirmStep := prompts.PromptStep{
				ID:        "confirmation",
				Title:     "Configuration Summary",
				Component: prompts.NewConfigurationSummary(po.config),
				Required:  true,
				HelpText:  "Review your configuration and choose how to proceed",
				Type:      prompts.PromptTypeConfirmation,
			}
			po.prompts = append(po.prompts, confirmStep)
			po.navigation.TotalSteps = len(po.prompts)
			return po.goToNextPrompt()
		}

		// Move to next prompt or complete
		if po.currentIndex < len(po.prompts)-1 {
			return po.goToNextPrompt()
		} else {
			po.completed = true
			return *po, nil
		}

	case prompts.NextPromptMsg:
		return po.goToNextPrompt()

	case prompts.PrevPromptMsg:
		return po.goToPreviousPrompt()

	case prompts.SkipPromptMsg:
		return po.skipCurrentPrompt()
	}

	// Update current prompt
	updatedComponent, cmd := currentPrompt.Component.Update(msg)
	currentPrompt.Component = updatedComponent

	return *po, cmd
}

// View renders the current prompt
func (po *PromptOrchestrator) View() string {
	if po.completed || len(po.prompts) == 0 {
		return "Configuration complete!"
	}

	return po.prompts[po.currentIndex].Component.View()
}

// IsComplete returns whether all prompts are completed
func (po *PromptOrchestrator) IsComplete() bool {
	return po.completed
}

// GetConfiguration returns the final user configuration
func (po *PromptOrchestrator) GetConfiguration() *config.UserConfiguration {
	return po.config
}

// GetCurrentStep returns current step information for progress display
func (po *PromptOrchestrator) GetCurrentStep() (int, int, string) {
	if po.completed {
		return len(po.prompts), len(po.prompts), "Configuration Complete"
	}

	step := po.currentIndex + 1
	total := len(po.prompts)
	title := po.prompts[po.currentIndex].Title

	return step, total, title
}

// Private methods

func (po *PromptOrchestrator) goToNextPrompt() (PromptOrchestrator, tea.Cmd) {
	if po.currentIndex < len(po.prompts)-1 {
		po.currentIndex++
		po.updateNavigation()

		// Initialize the new prompt
		cmd := po.prompts[po.currentIndex].Component.Init()
		return *po, cmd
	}
	return *po, nil
}

func (po *PromptOrchestrator) goToPreviousPrompt() (PromptOrchestrator, tea.Cmd) {
	if po.currentIndex > 0 {
		po.currentIndex--
		po.updateNavigation()

		// Re-initialize the previous prompt
		cmd := po.prompts[po.currentIndex].Component.Init()
		return *po, cmd
	}
	return *po, nil
}

func (po *PromptOrchestrator) skipCurrentPrompt() (PromptOrchestrator, tea.Cmd) {
	currentPrompt := &po.prompts[po.currentIndex]
	if !currentPrompt.Required {
		// Don't save value when skipping
		if po.currentIndex < len(po.prompts)-1 {
			return po.goToNextPrompt()
		} else {
			po.completed = true
			return *po, nil
		}
	}
	return *po, nil
}

func (po *PromptOrchestrator) updateNavigation() {
	po.navigation.CanGoBack = po.currentIndex > 0
	po.navigation.CanGoForward = po.currentIndex < len(po.prompts)-1
	po.navigation.CanSkip = !po.prompts[po.currentIndex].Required
	po.navigation.CurrentStep = po.currentIndex + 1
}

func (po *PromptOrchestrator) savePromptValue(promptStep *prompts.PromptStep) {
	value := promptStep.Component.GetValue()

	switch promptStep.Type {
	case prompts.PromptTypeTemplate:
		if template, ok := value.(config.TemplateType); ok {
			po.config.Template.Type = template
		}

	case prompts.PromptTypeDevFeatures:
		if devFeatures, ok := value.(config.DevFeatureConfig); ok {
			po.config.DevFeatures = devFeatures
		}

	case prompts.PromptTypeProductionSetup:
		if prodSetup, ok := value.(config.ProductionConfig); ok {
			po.config.ProductionSetup = prodSetup
		}

	case prompts.PromptTypeTesting:
		if testing, ok := value.(config.TestingConfig); ok {
			po.config.Testing = testing
		}

	case prompts.PromptTypeNavigation:
		if navigation, ok := value.(config.NavigationConfig); ok {
			po.config.Navigation = navigation
		}
	}
}