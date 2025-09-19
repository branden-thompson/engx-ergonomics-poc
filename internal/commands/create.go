package commands

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/models"
	"github.com/bthompso/engx-ergonomics-poc/internal/prompts"
	"github.com/bthompso/engx-ergonomics-poc/internal/config"
	"github.com/bthompso/engx-ergonomics-poc/internal/chaos"
	"github.com/spf13/cobra"
)

// NewCreateCommand creates the 'create' command
func NewCreateCommand() *cobra.Command {
	var devOnly bool
	var template string
	var chaosMarine bool
	var chaosLevel string
	var chaosSeed int64
	var chaosConfig string

	cmd := &cobra.Command{
		Use:   "create [APP_NAME]",
		Short: "Create a new React application",
		Long: `Create a new React application with modern development tooling.

This command simulates the creation of a new React application with:
- Modern React setup with TypeScript
- Development server configuration
- Testing framework setup
- Build pipeline configuration
- Deployment preparation (unless --dev-only)

Examples:
  engx create MyApp
  engx create MyApp --dev-only
  engx create MyApp --template=typescript
  engx create MyApp --verbose
  engx create MyApp --chaos-marine --chaos-level=scout
  engx create MyApp --chaos-marine --chaos-level=aggressive --chaos-seed=12345`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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
				chaosConfig, err := chaos.LoadChaosConfig(chaosLevel, chaosSeed, chaosConfig)
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
			if cmd.Flags().Changed("chaos-config") && chaosConfig != "" {
				flags = append(flags, fmt.Sprintf("--chaos-config=%s", chaosConfig))
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
		},
	}

	// Add command-specific flags
	cmd.Flags().BoolVar(&devOnly, "dev-only", false, "Create app for development only (skip production setup)")
	cmd.Flags().StringVar(&template, "template", "", "Template to use (typescript, javascript, minimal)")

	// Add chaos marine flags
	cmd.Flags().BoolVar(&chaosMarine, "chaos-marine", false, "Enable chaos injection for failure simulation")
	cmd.Flags().StringVar(&chaosLevel, "chaos-level", "default", "Chaos aggressiveness level (off, default, scout, aggressive, invasive, apocalyptic)")
	cmd.Flags().Int64Var(&chaosSeed, "chaos-seed", 0, "Random seed for deterministic chaos (0 = random)")
	cmd.Flags().StringVar(&chaosConfig, "chaos-config", "", "Path to chaos configuration file")

	return cmd
}