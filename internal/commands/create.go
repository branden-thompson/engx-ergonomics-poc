package commands

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/models"
	"github.com/bthompso/engx-ergonomics-poc/internal/prompts"
	"github.com/spf13/cobra"
)

// NewCreateCommand creates the 'create' command
func NewCreateCommand() *cobra.Command {
	var devOnly bool
	var template string

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
  engx create MyApp --verbose`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]

			// Collect only explicitly set flags
			var flags []string
			if cmd.Flags().Changed("dev-only") && devOnly {
				flags = append(flags, "--dev-only")
			}
			if cmd.Flags().Changed("template") && template != "" {
				flags = append(flags, fmt.Sprintf("--template=%s", template))
			}

			verbose, _ := cmd.Flags().GetBool("verbose")
			if verbose {
				flags = append(flags, "--verbose")
			}

			quiet, _ := cmd.Flags().GetBool("quiet")
			if quiet {
				flags = append(flags, "--quiet")
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
			model := models.NewAppModelWithConfig("create", appName, flags, userConfig)

			// Configure for inline mode with proper input/output handling
			program := tea.NewProgram(
				model,
				tea.WithInput(os.Stdin),
				tea.WithOutput(os.Stderr),
			)

			if _, err := program.Run(); err != nil {
				return fmt.Errorf("failed to run application: %w", err)
			}

			return nil
		},
	}

	// Add command-specific flags
	cmd.Flags().BoolVar(&devOnly, "dev-only", false, "Create app for development only (skip production setup)")
	cmd.Flags().StringVar(&template, "template", "", "Template to use (typescript, javascript, minimal)")

	return cmd
}