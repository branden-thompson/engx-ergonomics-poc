package commands

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/models"
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

			// Collect flags
			var flags []string
			if devOnly {
				flags = append(flags, "--dev-only")
			}
			if template != "" {
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

			// Initialize and run TUI
			model := models.NewAppModel("create", appName, flags)
			program := tea.NewProgram(model, tea.WithAltScreen())

			if _, err := program.Run(); err != nil {
				return fmt.Errorf("failed to run application: %w", err)
			}

			return nil
		},
	}

	// Add command-specific flags
	cmd.Flags().BoolVar(&devOnly, "dev-only", false, "Create app for development only (skip production setup)")
	cmd.Flags().StringVar(&template, "template", "typescript", "Template to use (typescript, javascript, minimal)")

	return cmd
}