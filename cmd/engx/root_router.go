package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bthompso/engx-ergonomics-poc/internal/config"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/models"
	"github.com/bthompso/engx-ergonomics-poc/pkg/common"
)

// RootCommandRouter handles the logic for determining when to show the interactive interface
type RootCommandRouter struct {
	rootCmd        *cobra.Command
	pluginRegistry *common.PluginRegistry
	deps           *common.Dependencies
}

// NewRootCommandRouter creates a new root command router
func NewRootCommandRouter(rootCmd *cobra.Command, pluginRegistry *common.PluginRegistry, deps *common.Dependencies) *RootCommandRouter {
	return &RootCommandRouter{
		rootCmd:        rootCmd,
		pluginRegistry: pluginRegistry,
		deps:           deps,
	}
}

// ShouldShowInteractiveInterface determines if the interactive interface should be shown
func (r *RootCommandRouter) ShouldShowInteractiveInterface(args []string) bool {
	// Show interactive interface only when:
	// 1. No arguments provided (just 'engx')
	// 2. Not requesting help (--help, -h)
	// 3. Not requesting version (--version)
	// 4. Only --show-dev-tools flag is provided

	if len(args) == 0 {
		return true
	}

	// Allow --show-dev-tools flag alone to trigger interactive interface
	if len(args) == 1 && args[0] == "--show-dev-tools" {
		return true
	}

	// Any other arguments - use normal cobra behavior
	return false
}

// ShowInteractiveInterface launches the interactive command selection interface
func (r *RootCommandRouter) ShowInteractiveInterface() error {
	// Check for --show-dev-tools flag
	showDevTools, _ := r.rootCmd.Flags().GetBool("show-dev-tools")
	// Load roadmap configuration
	roadmapConfig, err := config.LoadRoadmapConfig()
	if err != nil {
		// If roadmap config fails to load, fall back to help
		fmt.Fprintf(os.Stderr, "Warning: Could not load roadmap configuration: %v\n", err)
		fmt.Fprintf(os.Stderr, "Falling back to help screen.\n\n")
		return r.rootCmd.Help()
	}

	// Validate configuration
	if err := roadmapConfig.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Invalid roadmap configuration: %v\n", err)
		fmt.Fprintf(os.Stderr, "Falling back to help screen.\n\n")
		return r.rootCmd.Help()
	}

	// Create command discovery
	commandDiscovery := common.NewCommandDiscovery(roadmapConfig, r.rootCmd, r.pluginRegistry)

	// Create interactive selector
	selector, err := models.NewInteractiveCommandSelector(roadmapConfig, commandDiscovery, showDevTools)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create interactive interface: %v\n", err)
		fmt.Fprintf(os.Stderr, "Falling back to help screen.\n\n")
		return r.rootCmd.Help()
	}

	// Run the interactive interface
	selectedCmd, err := selector.Run()
	if err != nil {
		// User quit or error occurred
		return nil // Don't treat user quit as error
	}

	// Execute the selected command
	if selectedCmd != nil {
		return selector.ExecuteCommand(*selectedCmd)
	}

	return nil
}

// HandleRootCommand determines the appropriate action for the root command
func (r *RootCommandRouter) HandleRootCommand(args []string) error {
	// Check for special flags that should be handled normally
	for _, arg := range args {
		if arg == "--help" || arg == "-h" || arg == "--version" {
			// Let cobra handle these flags normally
			return nil
		}
	}

	// If no arguments and not a special flag, show interactive interface
	if r.ShouldShowInteractiveInterface(args) {
		return r.ShowInteractiveInterface()
	}

	// Otherwise, proceed with normal cobra behavior
	return nil
}

// ModifyRootCommand modifies the root command to include interactive behavior
func (r *RootCommandRouter) ModifyRootCommand() {
	// Store the original RunE function if it exists
	originalRunE := r.rootCmd.RunE

	// Set new RunE that checks for interactive mode
	r.rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// Get the raw arguments passed to the command
		rawArgs := os.Args[1:]

		// Handle root command routing
		if err := r.HandleRootCommand(rawArgs); err != nil {
			return err
		}

		// If we get here and there are no arguments, it means we should show help
		// (interactive interface was not shown for some reason)
		if len(rawArgs) == 0 {
			return cmd.Help()
		}

		// If original RunE exists, call it
		if originalRunE != nil {
			return originalRunE(cmd, args)
		}

		// Default behavior - show help
		return cmd.Help()
	}

	// Modify the help text to mention the interactive interface
	originalLong := r.rootCmd.Long
	enhancedLong := originalLong + `

Interactive Mode:
  Run 'engx' without any arguments to launch the interactive command selection interface.
  This provides a guided way to discover and execute available commands.`

	r.rootCmd.Long = enhancedLong
}

// ValidateEnvironment checks if the environment is suitable for the interactive interface
func (r *RootCommandRouter) ValidateEnvironment() error {
	// Check if stdout is a terminal
	if !isTerminal() {
		return fmt.Errorf("interactive interface requires a terminal")
	}

	// Check if roadmap configuration exists
	if _, err := os.Stat(".engx/roadmap.yaml"); os.IsNotExist(err) {
		return fmt.Errorf("roadmap configuration not found at .engx/roadmap.yaml")
	}

	return nil
}

// isTerminal checks if stdout is connected to a terminal
func isTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// GetDebugInfo returns debug information for troubleshooting
func (r *RootCommandRouter) GetDebugInfo() map[string]interface{} {
	info := make(map[string]interface{})

	// Environment info
	info["is_terminal"] = isTerminal()
	info["args"] = os.Args

	// Configuration info
	if roadmapConfig, err := config.LoadRoadmapConfig(); err == nil {
		configInfo := map[string]interface{}{
			"version":           roadmapConfig.Metadata.Version,
			"categories":        len(roadmapConfig.CommandCategories),
			"excluded_commands": len(roadmapConfig.ExcludedCommands),
			"ldap_user":         roadmapConfig.UserContext.ResolvedLDAP,
			"interface_config": map[string]interface{}{
				"minimum_width":     roadmapConfig.Interface.MinimumWidth,
				"default_width":     roadmapConfig.Interface.DefaultWidth,
				"expandable_column": roadmapConfig.Interface.ExpandableColumn,
			},
		}
		info["roadmap_config"] = configInfo
	} else {
		info["roadmap_config_error"] = err.Error()
	}

	// Command info
	info["cobra_commands"] = len(r.rootCmd.Commands())

	return info
}