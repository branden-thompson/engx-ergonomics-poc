package create

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/models"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
	"github.com/bthompso/engx-ergonomics-poc/internal/prompts"
	"github.com/bthompso/engx-ergonomics-poc/internal/config"
	"github.com/bthompso/engx-ergonomics-poc/internal/chaos"
	"github.com/bthompso/engx-ergonomics-poc/internal/workflows"
	"github.com/bthompso/engx-ergonomics-poc/internal/archetypes"
	"github.com/bthompso/engx-ergonomics-poc/pkg/common"
	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
	"github.com/spf13/cobra"
)

// Plugin implements the CommandPlugin interface for the create command
type Plugin struct {
	deps *common.Dependencies
}

// NewPlugin creates a new create command plugin
func NewPlugin(deps *common.Dependencies) interfaces.CommandPlugin {
	return &Plugin{deps: deps}
}

// Name returns the plugin name
func (p *Plugin) Name() string {
	return "create"
}

// Description returns the plugin description
func (p *Plugin) Description() string {
	return "Creates new application projects with interactive setup"
}

// Version returns the plugin version
func (p *Plugin) Version() string {
	return "1.0.0"
}

// Create returns the cobra command for this plugin
func (p *Plugin) Create(deps interface{}) *cobra.Command {
	// Cast deps back to our Dependencies type
	var dependencies *common.Dependencies
	if d, ok := deps.(*common.Dependencies); ok {
		dependencies = d
		p.deps = d
	} else {
		// Fallback to our stored deps
		dependencies = p.deps
	}

	var devOnly bool
	var template string
	var chaosMarine bool
	var chaosLevel string
	var chaosSeed int64
	var chaosConfig string

	cmd := &cobra.Command{
		Use:   "create [APP_NAME]",
		Short: "Create a new application project with guided archetype selection",
		Long: `Create a new application project with guided archetype selection.

This command can create various types of applications:
- Production React applications with full deployment pipeline
- Development React apps for rapid prototyping
- CLI tools and utilities
- Backend services and APIs
- Hackday prototypes
- EngX command plugins

Usage Modes:
  engx create              # Guided mode with archetype selection
  engx create <MyApp>      # Direct mode using default archetype

Examples:
  engx create                                    # Guided archetype selection
  engx create <MyApp>                           # Direct mode (prod-web archetype)
  engx create <MyApp> --dev-only                # Skip production setup
  engx create <MyApp> --template=typescript     # Specify template type
  engx create <MyApp> --chaos-marine            # Enable chaos testing`,
		Args: cobra.RangeArgs(0, 1), // Changed from ExactArgs(1) to support guided mode
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				// NEW: Guided mode
				return p.executeGuidedMode(cmd, dependencies, devOnly, template, chaosMarine, chaosLevel, chaosSeed, chaosConfig)
			} else {
				// EXISTING: Direct mode (preserved for backward compatibility)
				return p.executeDirectMode(cmd, args, dependencies, devOnly, template, chaosMarine, chaosLevel, chaosSeed, chaosConfig)
			}
		},
	}

	// Add flags
	cmd.Flags().BoolVar(&devOnly, "dev-only", false, "Skip production optimizations and deployment setup")
	cmd.Flags().StringVar(&template, "template", "typescript", "Project template (typescript, javascript)")
	cmd.Flags().BoolVar(&chaosMarine, "chaos-marine", false, "Enable chaos marine for testing error handling")
	cmd.Flags().StringVar(&chaosLevel, "chaos-level", "scout", "Chaos level (scout, marine, aggressive)")
	cmd.Flags().Int64Var(&chaosSeed, "chaos-seed", 0, "Random seed for chaos injection (0 for random)")
	cmd.Flags().StringVar(&chaosConfig, "chaos-config", "", "Path to chaos configuration file")

	// Add hidden test flag for guaranteed chaos injection
	var testGuaranteedChaos bool
	cmd.Flags().BoolVar(&testGuaranteedChaos, "test-guaranteed-chaos", false, "Test TUI with guaranteed chaos injection (hidden)")
	cmd.Flags().MarkHidden("test-guaranteed-chaos")

	return cmd
}

// Initialize initializes the plugin
func (p *Plugin) Initialize() error {
	if p.deps == nil {
		return fmt.Errorf("dependencies not provided")
	}
	return nil
}

// Cleanup performs any necessary cleanup
func (p *Plugin) Cleanup() error {
	return nil
}

// RequiredServices returns required service names
func (p *Plugin) RequiredServices() []string {
	return []string{"config", "filesystem", "logger", "aar"}
}

// OptionalServices returns optional service names
func (p *Plugin) OptionalServices() []string {
	return []string{"tui", "chaos"}
}

// executeDirectMode implements the original direct command logic (preserved for backward compatibility)
func (p *Plugin) executeDirectMode(cmd *cobra.Command, args []string, deps *common.Dependencies, devOnly bool, template string, chaosMarine bool, chaosLevel string, chaosSeed int64, chaosConfigPath string) error {
	appName := args[0]

	// Determine verbosity level from flags
	quiet, _ := cmd.Flags().GetBool("quiet")
	concise, _ := cmd.Flags().GetBool("concise")
	verbose, _ := cmd.Flags().GetBool("verbose")
	debug, _ := cmd.Flags().GetBool("debug")

	verbosityLevel := config.DetermineVerbosityLevel(quiet, concise, verbose, debug)
	verbosityConfig := config.NewVerbosityConfig(verbosityLevel)

	// Debug output for verbosity level determination
	verbosityConfig.DebugPrint("Verbosity level determined: %s", verbosityLevel.String())

	// Initialize chaos configuration if chaos marine is enabled
	var chaosInjector chaos.ChaosInjector
	if chaosMarine {
		chaosConfig, err := chaos.LoadChaosConfig(chaosLevel, chaosSeed, chaosConfigPath)
		if err != nil {
			return fmt.Errorf("failed to load chaos configuration: %w", err)
		}

		chaosInjector, err = chaos.NewSafeChaosInjector(chaosConfig)
		if err != nil {
			return fmt.Errorf("failed to initialize chaos injector: %w", err)
		}

		verbosityConfig.DebugPrint("Chaos Marine enabled: level=%s, seed=%d", chaosLevel, chaosSeed)
	}

	// Collect only explicitly set flags for display purposes
	var flags []string
	if cmd.Flags().Changed("dev-only") && devOnly {
		flags = append(flags, "--dev-only")
	}
	if cmd.Flags().Changed("template") && template != "" {
		flags = append(flags, fmt.Sprintf("--template=%s", template))
	}
	if cmd.Flags().Changed("chaos-marine") && chaosMarine {
		flags = append(flags, "--chaos-marine")
	}
	if cmd.Flags().Changed("chaos-level") && chaosLevel != "" {
		flags = append(flags, fmt.Sprintf("--chaos-level=%s", chaosLevel))
	}
	if cmd.Flags().Changed("chaos-seed") && chaosSeed != 0 {
		flags = append(flags, fmt.Sprintf("--chaos-seed=%d", chaosSeed))
	}
	if cmd.Flags().Changed("chaos-config") && chaosConfigPath != "" {
		flags = append(flags, fmt.Sprintf("--chaos-config=%s", chaosConfigPath))
	}

	// Add verbosity flags to display
	if quiet {
		flags = append(flags, "--quiet")
	}
	if concise {
		flags = append(flags, "--concise")
	}
	if verbose {
		flags = append(flags, "--verbose")
	}
	if debug {
		flags = append(flags, "--debug")
	}

	// Run inline prompts first (traditional CLI style)
	prompter, err := prompts.NewInlinePrompter()
	if err != nil {
		return fmt.Errorf("failed to initialize prompter: %w", err)
	}

	userConfig, err := prompter.RunPrompts(devOnly, flags)
	if err != nil {
		return fmt.Errorf("failed to run prompts: %w", err)
	}

	// Set the project name in config
	userConfig.ProjectName = appName

	// Initialize and run TUI with configuration already set (inline mode)
	var model *models.AppModel
	if chaosInjector != nil {
		model = models.NewAppModelWithChaos("create", appName, flags, userConfig, verbosityConfig, chaosInjector)
	} else {
		model = models.NewAppModelWithVerbosity("create", appName, flags, userConfig, verbosityConfig)
	}

	// Configure for inline mode with proper input/output handling
	program := tea.NewProgram(
		model,
		tea.WithInput(os.Stdin),
		tea.WithOutput(os.Stderr),
	)

	finalModel, err := program.Run()
	if err != nil {
		return fmt.Errorf("failed to run application: %w", err)
	}

	// Print AAR after TUI exits if available
	if appModel, ok := finalModel.(*models.AppModel); ok && appModel.GetAAROutput() != "" {
		fmt.Print(appModel.GetAAROutput())
	}

	return nil
}

// executeGuidedMode implements the new guided workflow with archetype selection
func (p *Plugin) executeGuidedMode(cmd *cobra.Command, deps *common.Dependencies, devOnly bool, template string, chaosMarine bool, chaosLevel string, chaosSeed int64, chaosConfigPath string) error {
	// Determine verbosity level from flags
	quiet, _ := cmd.Flags().GetBool("quiet")
	concise, _ := cmd.Flags().GetBool("concise")
	verbose, _ := cmd.Flags().GetBool("verbose")
	debug, _ := cmd.Flags().GetBool("debug")

	verbosityLevel := config.DetermineVerbosityLevel(quiet, concise, verbose, debug)
	verbosityConfig := config.NewVerbosityConfig(verbosityLevel)

	// Debug output for verbosity level determination
	verbosityConfig.DebugPrint("Guided mode verbosity level determined: %s", verbosityLevel.String())

	// Initialize chaos configuration if chaos marine is enabled
	var chaosInjector chaos.ChaosInjector
	if chaosMarine {
		chaosConfig, err := chaos.LoadChaosConfig(chaosLevel, chaosSeed, chaosConfigPath)
		if err != nil {
			return fmt.Errorf("failed to load chaos configuration: %w", err)
		}

		chaosInjector, err = chaos.NewSafeChaosInjector(chaosConfig)
		if err != nil {
			return fmt.Errorf("failed to initialize chaos injector: %w", err)
		}

		verbosityConfig.DebugPrint("Chaos Marine enabled in guided mode: level=%s, seed=%d", chaosLevel, chaosSeed)
	}

	// Collect explicitly set flags for display purposes
	var flags []string
	if cmd.Flags().Changed("dev-only") && devOnly {
		flags = append(flags, "--dev-only")
	}
	if cmd.Flags().Changed("template") && template != "" {
		flags = append(flags, fmt.Sprintf("--template=%s", template))
	}
	if cmd.Flags().Changed("chaos-marine") && chaosMarine {
		flags = append(flags, "--chaos-marine")
	}
	if cmd.Flags().Changed("chaos-level") && chaosLevel != "" {
		flags = append(flags, fmt.Sprintf("--chaos-level=%s", chaosLevel))
	}
	if cmd.Flags().Changed("chaos-seed") && chaosSeed != 0 {
		flags = append(flags, fmt.Sprintf("--chaos-seed=%d", chaosSeed))
	}
	if cmd.Flags().Changed("chaos-config") && chaosConfigPath != "" {
		flags = append(flags, fmt.Sprintf("--chaos-config=%s", chaosConfigPath))
	}

	// Add verbosity flags to display
	if quiet {
		flags = append(flags, "--quiet")
	}
	if concise {
		flags = append(flags, "--concise")
	}
	if verbose {
		flags = append(flags, "--verbose")
	}
	if debug {
		flags = append(flags, "--debug")
	}

	// Initialize workflow context
	workflowContext := &workflows.WorkflowContext{
		AppName:         "", // Will be set during archetype selection
		Flags:           flags,
		Dependencies:    deps,
		VerbosityConfig: verbosityConfig,
		ChaosInjector:   chaosInjector,
		IsGuidedMode:    true,
	}

	// Create workflow orchestrator
	orchestrator := workflows.NewWorkflowOrchestrator(workflowContext, deps.Logger)

	// Add archetype selection stage
	registry := archetypes.GetDefaultRegistry()
	archetypeStage := workflows.NewArchetypeSelectionStage(registry)
	orchestrator.AddStage(archetypeStage)

	// Professional wizard header using components
	wizardHeader := components.NewHeader("ENGX NEW APPLICATION WIZARD")
	fmt.Print(wizardHeader.Render())

	// Single paragraph with proper styling and smart syntax highlighting
	fmt.Printf("\033[90mLet's get your application configured and setup for development, this is the \033[93mGUIDED MODE\033[90m that is the default with \033[38;5;198mengx\033[0m \033[38;5;208mcreate\033[0m\033[90m. You can bypass this by specifying an \033[38;5;48m--app-type\033[0m\033[90m directly in the create command:\033[0m\n\n")

	// Smart command formatting with syntax highlighting
	cmdFormatter := components.NewCommandFormatter()
	fmt.Printf("%s\n\n", cmdFormatter.FormatCommandInBackticks("engx create <AppName> --app-type <app-type>"))

	// Execute workflow stages
	err := orchestrator.Execute()
	if err != nil {
		return fmt.Errorf("guided workflow failed: %w", err)
	}

	// Workflow complete - get the updated context
	updatedContext := orchestrator.GetContext()

	// Prompt for app name if not set during archetype selection
	if updatedContext.AppName == "" {
		appName, err := archetypeStage.PromptForAppName()
		if err != nil {
			return fmt.Errorf("failed to get app name: %w", err)
		}
		updatedContext.AppName = appName
	}

	// Initialize and run inline prompts (existing behavior preserved)
	prompter, err := prompts.NewInlinePrompter()
	if err != nil {
		return fmt.Errorf("failed to initialize prompter: %w", err)
	}

	userConfig, err := prompter.RunPromptsWithContext(devOnly, flags, updatedContext.SelectedArchetype.Name, updatedContext.AppName)
	if err != nil {
		return fmt.Errorf("failed to run prompts: %w", err)
	}

	// Set the project name in config
	userConfig.ProjectName = updatedContext.AppName
	updatedContext.UserConfiguration = userConfig

	// Initialize and run TUI with configuration already set and archetype information
	var model *models.AppModel
	archetypeID := ""
	if updatedContext.SelectedArchetype != nil {
		archetypeID = updatedContext.SelectedArchetype.ID
	}

	if chaosInjector != nil {
		model = models.NewAppModelWithArchetypeAndChaos("create", updatedContext.AppName, flags, userConfig, verbosityConfig, archetypeID, chaosInjector)
	} else {
		model = models.NewAppModelWithArchetype("create", updatedContext.AppName, flags, userConfig, verbosityConfig, archetypeID)
	}

	// Configure for inline mode with proper input/output handling
	program := tea.NewProgram(
		model,
		tea.WithInput(os.Stdin),
		tea.WithOutput(os.Stderr),
	)

	finalModel, err := program.Run()
	if err != nil {
		return fmt.Errorf("failed to run guided workflow: %w", err)
	}

	// Print AAR after TUI exits if available
	if appModel, ok := finalModel.(*models.AppModel); ok && appModel.GetAAROutput() != "" {
		fmt.Print(appModel.GetAAROutput())
	}

	return nil
}