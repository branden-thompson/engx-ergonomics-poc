package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common"
	"github.com/bthompso/engx-ergonomics-poc/plugins/create"
	testerror "github.com/bthompso/engx-ergonomics-poc/plugins/test-error"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Initialize dependencies
	deps := common.NewDependencies()

	// Initialize plugin registry
	registry := common.NewPluginRegistry(deps)

	// Register plugins
	if err := registerPlugins(registry, deps); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to register plugins: %v\n", err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "engx",
		Short: "ENGX - Engineering Productivity CLI (POC)",
		Long: `ENGX POC - A terminal-based simulation of engineering productivity tools and workflows.
This tool demonstrates human-computer interaction patterns for developer tooling.

Focus: Terminal UI/UX patterns, not actual application scaffolding.`,
		Version: fmt.Sprintf("%s (%s) built on %s", version, commit, date),
	}

	// Add global verbosity flags (mutually exclusive)
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Show only essential information")
	rootCmd.PersistentFlags().Bool("concise", false, "Show less detail with components hidden")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show enhanced details with progress bars for multi-step")
	rootCmd.PersistentFlags().Bool("debug", false, "Show maximum verbosity with all system outputs")

	// Other global flags
	rootCmd.PersistentFlags().String("config", "", "Config file (default searches for .engx/config.yaml)")
	rootCmd.PersistentFlags().Bool("show-dev-tools", false, "Show developer tools section in interactive interface")

	// Mark verbosity flags as mutually exclusive
	rootCmd.MarkFlagsMutuallyExclusive("quiet", "concise", "verbose", "debug")

	// Add plugin-based commands
	if err := addPluginCommands(rootCmd, registry); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add plugin commands: %v\n", err)
		os.Exit(1)
	}

	// Add plugin discovery command for testing
	pluginDiscoveryCmd := &cobra.Command{
		Use:   "plugins",
		Short: "List discovered plugins",
		Long:  "Discover and list all available plugins in the system",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPluginDiscovery(registry, deps)
		},
	}
	rootCmd.AddCommand(pluginDiscoveryCmd)

	// Add plugin validation commands
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate plugins",
		Long:  "Validate registered plugins for compliance and dependency satisfaction",
	}

	validatePluginsCmd := &cobra.Command{
		Use:   "plugins",
		Short: "Validate all registered plugins",
		Long:  "Validate all registered plugins for compliance, dependencies, and proper implementation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPluginValidation(registry, deps)
		},
	}

	validatePluginCmd := &cobra.Command{
		Use:   "plugin [PLUGIN_NAME]",
		Short: "Validate a specific plugin",
		Long:  "Validate a specific plugin for compliance, dependencies, and proper implementation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSinglePluginValidation(registry, deps, args[0])
		},
	}

	validateServicesCmd := &cobra.Command{
		Use:   "services",
		Short: "List available services",
		Long:  "List all available services for plugin dependency checking",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServiceValidation(registry, deps)
		},
	}

	validateCmd.AddCommand(validatePluginsCmd)
	validateCmd.AddCommand(validatePluginCmd)
	validateCmd.AddCommand(validateServicesCmd)
	rootCmd.AddCommand(validateCmd)

	// Add plugin development commands
	devCmd := &cobra.Command{
		Use:   "dev",
		Short: "Plugin development tools",
		Long:  "Tools for developing, generating, and managing plugins",
	}

	devGenerateCmd := &cobra.Command{
		Use:   "generate [PLUGIN_NAME]",
		Short: "Generate a new plugin",
		Long:  "Generate a new plugin with scaffolding, templates, and documentation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPluginGeneration(registry, deps, args[0])
		},
	}

	devTemplatesCmd := &cobra.Command{
		Use:   "templates",
		Short: "List available plugin templates",
		Long:  "List all available plugin templates for generation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplateList(registry, deps)
		},
	}

	devExampleCmd := &cobra.Command{
		Use:   "example [TEMPLATE_TYPE]",
		Short: "Generate example plugin from template",
		Long:  "Generate an example plugin from a predefined template (basic, advanced, tui)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExampleGeneration(registry, deps, args[0])
		},
	}

	// Add hot-reload commands
	devHotReloadCmd := &cobra.Command{
		Use:   "hotreload",
		Short: "Manage hot-reload system",
		Long:  "Enable, disable, and manage the hot-reload system for plugin development",
	}

	devHotReloadEnableCmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable hot-reload",
		Long:  "Enable the hot-reload system to automatically reload plugins when files change",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHotReloadEnable(registry, deps)
		},
	}

	devHotReloadDisableCmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable hot-reload",
		Long:  "Disable the hot-reload system",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHotReloadDisable(registry, deps)
		},
	}

	devHotReloadStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show hot-reload status",
		Long:  "Show the current status of the hot-reload system and plugin states",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHotReloadStatus(registry, deps)
		},
	}

	devPluginCmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage plugins at runtime",
		Long:  "Enable, disable, and manage individual plugins at runtime",
	}

	devPluginEnableCmd := &cobra.Command{
		Use:   "enable [PLUGIN_NAME]",
		Short: "Enable a plugin",
		Long:  "Enable a plugin at runtime",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPluginEnable(registry, deps, args[0])
		},
	}

	devPluginDisableCmd := &cobra.Command{
		Use:   "disable [PLUGIN_NAME]",
		Short: "Disable a plugin",
		Long:  "Disable a plugin at runtime",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPluginDisable(registry, deps, args[0])
		},
	}

	devPluginListCmd := &cobra.Command{
		Use:   "list",
		Short: "List plugin states",
		Long:  "List all plugins and their current runtime states",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPluginList(registry, deps)
		},
	}

	devHotReloadCmd.AddCommand(devHotReloadEnableCmd)
	devHotReloadCmd.AddCommand(devHotReloadDisableCmd)
	devHotReloadCmd.AddCommand(devHotReloadStatusCmd)

	devPluginCmd.AddCommand(devPluginEnableCmd)
	devPluginCmd.AddCommand(devPluginDisableCmd)
	devPluginCmd.AddCommand(devPluginListCmd)

	// Add plugin configuration commands
	devConfigCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage plugin configurations",
		Long:  "Interactive configuration management for plugins",
	}

	devConfigListCmd := &cobra.Command{
		Use:   "list",
		Short: "List plugin configurations",
		Long:  "List all plugin configurations with their current status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigList(registry, deps)
		},
	}

	devConfigEditCmd := &cobra.Command{
		Use:   "edit [PLUGIN_NAME]",
		Short: "Edit plugin configuration",
		Long:  "Open an interactive editor for plugin configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigEdit(registry, deps, args[0])
		},
	}

	devConfigWizardCmd := &cobra.Command{
		Use:   "wizard [PLUGIN_NAME]",
		Short: "Run configuration wizard",
		Long:  "Run an interactive configuration wizard for a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigWizard(registry, deps, args[0])
		},
	}

	devConfigGetCmd := &cobra.Command{
		Use:   "get [PLUGIN_NAME] [KEY]",
		Short: "Get configuration value",
		Long:  "Get a specific configuration value for a plugin",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigGet(registry, deps, args[0], args[1])
		},
	}

	devConfigSetCmd := &cobra.Command{
		Use:   "set [PLUGIN_NAME] [KEY] [VALUE]",
		Short: "Set configuration value",
		Long:  "Set a specific configuration value for a plugin",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigSet(registry, deps, args[0], args[1], args[2])
		},
	}

	devConfigResetCmd := &cobra.Command{
		Use:   "reset [PLUGIN_NAME]",
		Short: "Reset plugin configuration",
		Long:  "Reset a plugin configuration to its defaults",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigReset(registry, deps, args[0])
		},
	}

	devConfigExportCmd := &cobra.Command{
		Use:   "export [PLUGIN_NAME] [FORMAT]",
		Short: "Export plugin configuration",
		Long:  "Export a plugin configuration to a file (yaml/json)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigExport(registry, deps, args[0], args[1])
		},
	}

	devConfigImportCmd := &cobra.Command{
		Use:   "import [PLUGIN_NAME] [FILE]",
		Short: "Import plugin configuration",
		Long:  "Import a plugin configuration from a file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigImport(registry, deps, args[0], args[1])
		},
	}

	devConfigUICmd := &cobra.Command{
		Use:   "ui",
		Short: "Open configuration UI",
		Long:  "Open an interactive terminal-based configuration interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigUI(registry, deps)
		},
	}

	devConfigStatsCmd := &cobra.Command{
		Use:   "stats",
		Short: "Show configuration statistics",
		Long:  "Show statistics about plugin configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigStats(registry, deps)
		},
	}

	devConfigCreateCmd := &cobra.Command{
		Use:   "create [PLUGIN_NAME]",
		Short: "Create plugin configuration",
		Long:  "Create a new plugin configuration with basic setup",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigCreate(registry, deps, args[0])
		},
	}

	devConfigCmd.AddCommand(devConfigCreateCmd)
	devConfigCmd.AddCommand(devConfigListCmd)
	devConfigCmd.AddCommand(devConfigEditCmd)
	devConfigCmd.AddCommand(devConfigWizardCmd)
	devConfigCmd.AddCommand(devConfigGetCmd)
	devConfigCmd.AddCommand(devConfigSetCmd)
	devConfigCmd.AddCommand(devConfigResetCmd)
	devConfigCmd.AddCommand(devConfigExportCmd)
	devConfigCmd.AddCommand(devConfigImportCmd)
	devConfigCmd.AddCommand(devConfigUICmd)
	devConfigCmd.AddCommand(devConfigStatsCmd)

	// Template discovery commands (replaces marketplace for simulation)
	templatesCmd := &cobra.Command{
		Use:   "templates",
		Short: "React template discovery and selection",
		Long:  "Browse, search, and explore React application templates for CLI ergonomics simulation",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default to list when no subcommand is provided
			return runTemplatesList(registry, deps)
		},
	}

	templatesListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all available React templates",
		Long:  "Display all React templates with recommendations and popularity indicators",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplatesList(registry, deps)
		},
	}

	templatesSearchCmd := &cobra.Command{
		Use:   "search [QUERY]",
		Short: "Search React templates",
		Long:  "Search for React templates by name, framework, features, or tags",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplatesSearch(registry, deps, args[0])
		},
	}

	templatesInfoCmd := &cobra.Command{
		Use:   "info [TEMPLATE_ID]",
		Short: "Show template information",
		Long:  "Display detailed information about a specific React template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplatesInfo(registry, deps, args[0])
		},
	}

	templatesRecommendedCmd := &cobra.Command{
		Use:   "recommended",
		Short: "Show recommended templates",
		Long:  "Display only recommended React templates for quick selection",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplatesRecommended(registry, deps)
		},
	}

	templatesComplexityCmd := &cobra.Command{
		Use:   "complexity [LEVEL]",
		Short: "Filter templates by complexity",
		Long:  "Show templates filtered by complexity level (beginner, intermediate, advanced, enterprise)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplatesComplexity(registry, deps, args[0])
		},
	}

	templatesStatsCmd := &cobra.Command{
		Use:   "stats",
		Short: "Show template statistics",
		Long:  "Display statistics about available React templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplatesStats(registry, deps)
		},
	}

	templatesDetailsCmd := &cobra.Command{
		Use:   "details [APP-TYPE]",
		Short: "Show detailed archetype information",
		Long:  "Display detailed information about a specific application archetype by app-type",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplatesDetails(registry, deps, args[0])
		},
	}

	templatesCmd.AddCommand(templatesListCmd)
	templatesCmd.AddCommand(templatesSearchCmd)
	templatesCmd.AddCommand(templatesInfoCmd)
	templatesCmd.AddCommand(templatesDetailsCmd)
	templatesCmd.AddCommand(templatesRecommendedCmd)
	templatesCmd.AddCommand(templatesComplexityCmd)
	templatesCmd.AddCommand(templatesStatsCmd)

	devCmd.AddCommand(devGenerateCmd)
	devCmd.AddCommand(devTemplatesCmd)
	devCmd.AddCommand(devExampleCmd)
	devCmd.AddCommand(devHotReloadCmd)
	devCmd.AddCommand(devPluginCmd)
	devCmd.AddCommand(devConfigCmd)
	// Analytics commands for CLI interaction pattern analysis
	analyticsCmd := &cobra.Command{
		Use:   "analytics",
		Short: "CLI interaction analytics and usage patterns",
		Long:  "Analyze CLI usage patterns, workflow sequences, and command ergonomics for simulation insights",
	}

	analyticsSummaryCmd := &cobra.Command{
		Use:   "summary",
		Short: "Show session analytics summary",
		Long:  "Display a summary of current CLI session analytics and usage patterns",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalyticsSummary(registry, deps)
		},
	}

	analyticsDetailsCmd := &cobra.Command{
		Use:   "details",
		Short: "Show detailed analytics",
		Long:  "Display comprehensive CLI analytics including command stats and patterns",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalyticsDetails(registry, deps)
		},
	}

	analyticsPatternsCmd := &cobra.Command{
		Use:   "patterns",
		Short: "Show workflow patterns",
		Long:  "Display detected workflow patterns and command sequences",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalyticsPatterns(registry, deps)
		},
	}

	analyticsStatsCmd := &cobra.Command{
		Use:   "stats",
		Short: "Show usage statistics",
		Long:  "Display command usage statistics and performance metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalyticsStats(registry, deps)
		},
	}

	analyticsExportCmd := &cobra.Command{
		Use:   "export [OUTPUT_FILE]",
		Short: "Export analytics data",
		Long:  "Export analytics data to JSON file for external analysis",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFile := ""
			if len(args) > 0 {
				outputFile = args[0]
			}
			return runAnalyticsExport(registry, deps, outputFile)
		},
	}

	analyticsStatusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show analytics system status",
		Long:  "Display analytics system status and configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalyticsStatus(registry, deps)
		},
	}

	analyticsCmd.AddCommand(analyticsSummaryCmd)
	analyticsCmd.AddCommand(analyticsDetailsCmd)
	analyticsCmd.AddCommand(analyticsPatternsCmd)
	analyticsCmd.AddCommand(analyticsStatsCmd)
	analyticsCmd.AddCommand(analyticsExportCmd)
	analyticsCmd.AddCommand(analyticsStatusCmd)

	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(templatesCmd)
	rootCmd.AddCommand(analyticsCmd)

	// Initialize root command router for interactive interface
	rootRouter := NewRootCommandRouter(rootCmd, registry, deps)
	rootRouter.ModifyRootCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// registerPlugins registers all available plugins with the registry
func registerPlugins(registry *common.PluginRegistry, deps *common.Dependencies) error {
	// Register create plugin
	createPlugin := create.NewPlugin(deps)
	if err := registry.Register(createPlugin); err != nil {
		return fmt.Errorf("failed to register create plugin: %w", err)
	}

	// Register test-error plugin
	testErrorPlugin := testerror.NewPlugin()
	if err := registry.Register(testErrorPlugin); err != nil {
		return fmt.Errorf("failed to register test-error plugin: %w", err)
	}

	// Initialize plugins
	if err := registry.InitializeAll(); err != nil {
		return fmt.Errorf("failed to initialize plugins: %w", err)
	}

	return nil
}

// addPluginCommands adds cobra commands for all registered plugins
func addPluginCommands(rootCmd *cobra.Command, registry *common.PluginRegistry) error {
	commands, err := registry.CreateCommands()
	if err != nil {
		return fmt.Errorf("failed to create plugin commands: %w", err)
	}

	for _, cmd := range commands {
		rootCmd.AddCommand(cmd)
	}

	return nil
}

// runPluginDiscovery demonstrates the plugin discovery system
func runPluginDiscovery(registry *common.PluginRegistry, deps *common.Dependencies) error {
	fmt.Println("🔍 Discovering plugins...")

	// Create plugin discovery system
	discovery := common.NewPluginDiscovery(registry, deps)

	// Discover plugins
	plugins, err := discovery.DiscoverPlugins()
	if err != nil {
		return fmt.Errorf("plugin discovery failed: %w", err)
	}

	// Display results
	fmt.Printf("\n📋 Found %d plugins:\n", len(plugins))
	for i, plugin := range plugins {
		fmt.Printf("\n%d. %s\n", i+1, plugin.Name)
		fmt.Printf("   📁 Path: %s\n", plugin.Path)
		fmt.Printf("   📦 Package: %s\n", plugin.PackageName)
		fmt.Printf("   🔗 Entry Point: %s\n", plugin.EntryPoint)
		fmt.Printf("   📝 Description: %s\n", plugin.Description)

		// Validate the plugin
		if err := discovery.ValidatePlugin(plugin); err != nil {
			fmt.Printf("   ❌ Validation: %v\n", err)
		} else {
			fmt.Printf("   ✅ Validation: OK\n")
		}
	}

	// Show currently registered plugins for comparison
	fmt.Printf("\n🔌 Currently registered plugins:\n")
	registeredPlugins := registry.ListPlugins()
	for _, pluginInfo := range registeredPlugins {
		fmt.Printf("   • %s - %s (v%s)\n", pluginInfo.Metadata.Name, pluginInfo.Metadata.Description, pluginInfo.Metadata.Version)
	}

	return nil
}

// runPluginValidation validates all registered plugins
func runPluginValidation(registry *common.PluginRegistry, deps *common.Dependencies) error {
	fmt.Println("🔍 Validating all registered plugins...")

	validator := registry.GetValidator()
	results, err := validator.ValidateAllPlugins()
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Printf("\n📋 Validation Results (%d plugins):\n", len(results))

	validCount := 0
	for pluginName, result := range results {
		fmt.Printf("\n🔧 Plugin: %s\n", pluginName)

		if result.IsValid {
			fmt.Printf("   ✅ Status: VALID\n")
			validCount++
		} else {
			fmt.Printf("   ❌ Status: INVALID\n")
		}

		// Show dependency status
		deps := result.Dependencies
		if deps.RequiredCount > 0 {
			fmt.Printf("   📦 Dependencies: %d/%d satisfied", deps.SatisfiedCount, deps.RequiredCount)
			if deps.Satisfied {
				fmt.Printf(" ✅\n")
			} else {
				fmt.Printf(" ❌\n")
				if len(deps.Missing) > 0 {
					fmt.Printf("      Missing: %s\n", fmt.Sprintf("%v", deps.Missing))
				}
			}
		}

		// Show errors
		if len(result.Errors) > 0 {
			fmt.Printf("   ❌ Errors:\n")
			for _, err := range result.Errors {
				fmt.Printf("      • %s\n", err)
			}
		}

		// Show warnings
		if len(result.Warnings) > 0 {
			fmt.Printf("   ⚠️  Warnings:\n")
			for _, warning := range result.Warnings {
				fmt.Printf("      • %s\n", warning)
			}
		}
	}

	fmt.Printf("\n📊 Summary: %d/%d plugins valid\n", validCount, len(results))

	if validCount == len(results) {
		fmt.Println("🎉 All plugins passed validation!")
	} else {
		fmt.Printf("⚠️  %d plugins failed validation\n", len(results)-validCount)
	}

	return nil
}

// runSinglePluginValidation validates a specific plugin
func runSinglePluginValidation(registry *common.PluginRegistry, deps *common.Dependencies, pluginName string) error {
	fmt.Printf("🔍 Validating plugin: %s\n", pluginName)

	plugin := registry.GetPlugin(pluginName)
	if plugin == nil {
		return fmt.Errorf("plugin '%s' not found", pluginName)
	}

	validator := registry.GetValidator()
	result, err := validator.ValidatePlugin(plugin)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Printf("\n📋 Validation Result for '%s':\n", pluginName)

	if result.IsValid {
		fmt.Printf("✅ Status: VALID\n")
	} else {
		fmt.Printf("❌ Status: INVALID\n")
	}

	// Show metadata validation
	metadata := result.Metadata
	fmt.Printf("\n📝 Metadata Validation:\n")
	fmt.Printf("   Name: %s %s\n", pluginName, getStatusIcon(metadata.NameValid))
	fmt.Printf("   Version: %s %s\n", plugin.Version(), getStatusIcon(metadata.VersionValid))
	fmt.Printf("   Description: %s %s\n", getStatusIcon(metadata.DescriptionValid), plugin.Description())

	if len(metadata.Issues) > 0 {
		fmt.Printf("   Issues:\n")
		for _, issue := range metadata.Issues {
			fmt.Printf("      • %s\n", issue)
		}
	}

	// Show detailed dependency status
	deps_result := result.Dependencies
	fmt.Printf("\n📦 Dependencies:\n")
	fmt.Printf("   Required: %d (satisfied: %d)\n", deps_result.RequiredCount, deps_result.SatisfiedCount)

	if len(deps_result.Required) > 0 {
		fmt.Printf("   Required Services:\n")
		for _, dep := range deps_result.Required {
			status := "❌"
			if dep.Available {
				status = "✅"
			}
			fmt.Printf("      %s %s - %s\n", status, dep.Name, dep.Description)
		}
	}

	if len(deps_result.Optional) > 0 {
		fmt.Printf("   Optional Services:\n")
		for _, dep := range deps_result.Optional {
			status := "❌"
			if dep.Available {
				status = "✅"
			}
			fmt.Printf("      %s %s - %s\n", status, dep.Name, dep.Description)
		}
	}

	// Show errors and warnings
	if len(result.Errors) > 0 {
		fmt.Printf("\n❌ Errors:\n")
		for _, err := range result.Errors {
			fmt.Printf("   • %s\n", err)
		}
	}

	if len(result.Warnings) > 0 {
		fmt.Printf("\n⚠️  Warnings:\n")
		for _, warning := range result.Warnings {
			fmt.Printf("   • %s\n", warning)
		}
	}

	return nil
}

// runServiceValidation shows available services
func runServiceValidation(registry *common.PluginRegistry, deps *common.Dependencies) error {
	fmt.Println("🔍 Available Services for Plugin Dependencies...")

	validator := registry.GetValidator()
	services := validator.GetKnownServices()

	fmt.Printf("\n📋 Known Services (%d):\n", len(services))

	for name, service := range services {
		status := "❌"
		if service.Available {
			status = "✅"
		}

		fmt.Printf("\n🔧 %s %s\n", status, name)
		fmt.Printf("   Version: %s\n", service.Version)
		fmt.Printf("   Description: %s\n", service.Description)
	}

	return nil
}

// getStatusIcon returns an appropriate status icon
func getStatusIcon(valid bool) string {
	if valid {
		return "✅"
	}
	return "❌"
}

// runPluginGeneration runs the interactive plugin generation process
func runPluginGeneration(registry *common.PluginRegistry, deps *common.Dependencies, pluginName string) error {
	fmt.Printf("🔧 Generating new plugin: %s\n", pluginName)

	generator := common.NewPluginGenerator(deps, registry)

	// Interactive template selection
	fmt.Println("\n📋 Available templates:")
	templates := generator.GetAvailableTemplates()
	for i, template := range templates {
		fmt.Printf("  %d. %s - %s\n", i+1, template.Name, template.Description)
	}

	fmt.Print("\nSelect template (1-3): ")
	var choice int
	if _, err := fmt.Scanln(&choice); err != nil || choice < 1 || choice > len(templates) {
		return fmt.Errorf("invalid template selection")
	}

	selectedTemplate := templates[choice-1]

	// Customize template
	pluginTemplate := common.PluginTemplate{
		Name:             selectedTemplate.Name,
		Description:      fmt.Sprintf("%s functionality", pluginName),
		Version:          "1.0.0",
		Author:           "Developer",
		PackageName:      pluginName,
		CommandName:      pluginName,
		RequiredServices: selectedTemplate.RequiredServices,
		OptionalServices: selectedTemplate.OptionalServices,
		Advanced:         selectedTemplate.Advanced,
		Flags:            selectedTemplate.Flags,
		Examples:         []string{fmt.Sprintf("engx %s", pluginName)},
	}

	// Get custom description
	fmt.Printf("\nEnter plugin description (default: %s): ", pluginTemplate.Description)
	var description string
	fmt.Scanln(&description)
	if description != "" {
		pluginTemplate.Description = description
	}

	// Get author
	fmt.Print("Enter author name (default: Developer): ")
	var author string
	fmt.Scanln(&author)
	if author != "" {
		pluginTemplate.Author = author
	}

	// Generate plugin
	result, err := generator.GeneratePlugin(pluginTemplate)
	if err != nil {
		return fmt.Errorf("plugin generation failed: %w", err)
	}

	// Display results
	if result.Success {
		fmt.Printf("\n✅ %s\n", result.Message)
		fmt.Printf("\n📁 Plugin created at: %s\n", result.PluginPath)
		fmt.Printf("📄 Files created:\n")
		for _, file := range result.FilesCreated {
			fmt.Printf("   • %s\n", file)
		}

		fmt.Printf("\n🔧 Next steps:\n")
		fmt.Printf("   1. Edit the generated plugin code in %s/plugin.go\n", result.PluginPath)
		fmt.Printf("   2. Build the project: go build -o dist/engx ./cmd/engx\n")
		fmt.Printf("   3. Test your plugin: ./dist/engx %s --help\n", pluginName)
		fmt.Printf("   4. Validate your plugin: ./dist/engx validate plugin %s\n", pluginName)
	} else {
		fmt.Printf("❌ %s\n", result.Message)
	}

	return nil
}

// runTemplateList lists available plugin templates
func runTemplateList(registry *common.PluginRegistry, deps *common.Dependencies) error {
	generator := common.NewPluginGenerator(deps, registry)
	fmt.Print(generator.ListPluginTemplates())
	return nil
}

// runExampleGeneration generates an example plugin from a template
func runExampleGeneration(registry *common.PluginRegistry, deps *common.Dependencies, templateType string) error {
	generator := common.NewPluginGenerator(deps, registry)
	templates := generator.GetAvailableTemplates()

	var selectedTemplate *common.PluginTemplate
	switch templateType {
	case "basic":
		selectedTemplate = &templates[0]
	case "advanced":
		selectedTemplate = &templates[1]
	case "tui":
		selectedTemplate = &templates[2]
	default:
		return fmt.Errorf("invalid template type '%s'. Available: basic, advanced, tui", templateType)
	}

	// Create example plugin with predefined settings
	pluginTemplate := common.PluginTemplate{
		Name:             selectedTemplate.Name,
		Description:      fmt.Sprintf("Example %s plugin", templateType),
		Version:          "1.0.0",
		Author:           "ENGX Plugin Generator",
		PackageName:      fmt.Sprintf("example_%s", templateType),
		CommandName:      fmt.Sprintf("example-%s", templateType),
		RequiredServices: selectedTemplate.RequiredServices,
		OptionalServices: selectedTemplate.OptionalServices,
		Advanced:         selectedTemplate.Advanced,
		Flags:            selectedTemplate.Flags,
		Examples:         []string{fmt.Sprintf("engx example-%s", templateType)},
	}

	fmt.Printf("🔧 Generating example %s plugin...\n", templateType)

	result, err := generator.GeneratePlugin(pluginTemplate)
	if err != nil {
		return fmt.Errorf("example plugin generation failed: %w", err)
	}

	if result.Success {
		fmt.Printf("✅ %s\n", result.Message)
		fmt.Printf("\n📁 Example plugin created at: %s\n", result.PluginPath)
		fmt.Printf("📄 Files created:\n")
		for _, file := range result.FilesCreated {
			fmt.Printf("   • %s\n", file)
		}

		fmt.Printf("\n📖 Study the generated code to understand:\n")
		if templateType == "basic" {
			fmt.Printf("   • Basic plugin structure and interface implementation\n")
			fmt.Printf("   • Simple command creation and execution\n")
			fmt.Printf("   • Dependency injection patterns\n")
		} else if templateType == "advanced" {
			fmt.Printf("   • Advanced plugin interface implementation\n")
			fmt.Printf("   • Configuration management and validation\n")
			fmt.Printf("   • Health checks and metadata handling\n")
			fmt.Printf("   • Event lifecycle hooks\n")
		} else if templateType == "tui" {
			fmt.Printf("   • Terminal User Interface integration\n")
			fmt.Printf("   • Interactive command patterns\n")
			fmt.Printf("   • TUI service dependency usage\n")
		}

		fmt.Printf("\n🔧 To integrate this example:\n")
		fmt.Printf("   1. Study the code in %s/\n", result.PluginPath)
		fmt.Printf("   2. Copy patterns to your own plugins\n")
		fmt.Printf("   3. Modify the registration in cmd/engx/main.go to include the example\n")
	} else {
		fmt.Printf("❌ %s\n", result.Message)
	}

	return nil
}
// Hot-reload command implementations

// runHotReloadEnable enables the hot-reload system
func runHotReloadEnable(registry *common.PluginRegistry, deps *common.Dependencies) error {
	fmt.Printf("🔥 Enabling hot-reload system...\n")

	if err := registry.EnableHotReload(); err != nil {
		return fmt.Errorf("failed to enable hot-reload: %w", err)
	}

	fmt.Printf("✅ Hot-reload system enabled\n")
	fmt.Printf("📁 Watching directories for plugin changes:\n")
	fmt.Printf("   • ./plugins\n")
	fmt.Printf("   • ./internal/plugins\n")
	fmt.Printf("\n💡 The system will automatically detect and reload plugins when .go files change\n")

	return nil
}

// runHotReloadDisable disables the hot-reload system
func runHotReloadDisable(registry *common.PluginRegistry, deps *common.Dependencies) error {
	fmt.Printf("🔥 Disabling hot-reload system...\n")

	if err := registry.DisableHotReload(); err != nil {
		return fmt.Errorf("failed to disable hot-reload: %w", err)
	}

	fmt.Printf("✅ Hot-reload system disabled\n")

	return nil
}

// runHotReloadStatus shows the hot-reload system status
func runHotReloadStatus(registry *common.PluginRegistry, deps *common.Dependencies) error {
	hotReload := registry.GetHotReloadManager()
	if hotReload == nil {
		fmt.Printf("❌ Hot-reload manager not initialized\n")
		return nil
	}

	// Refresh plugin states to ensure they're up to date
	hotReload.RefreshPluginStates()

	enabled := registry.IsHotReloadEnabled()
	fmt.Printf("🔥 Hot-Reload System Status\n\n")

	if enabled {
		fmt.Printf("Status: ✅ ENABLED\n")
	} else {
		fmt.Printf("Status: ❌ DISABLED\n")
	}

	// Get and display statistics
	stats := hotReload.GetStats()
	fmt.Printf("\n📊 Statistics:\n")
	fmt.Printf("   Total plugins: %v\n", stats["total_plugins"])
	fmt.Printf("   Enabled plugins: %v\n", stats["enabled_plugins"])
	fmt.Printf("   Loaded plugins: %v\n", stats["loaded_plugins"])
	fmt.Printf("   Total reloads: %v\n", stats["total_reloads"])
	fmt.Printf("   Callback count: %v\n", stats["callback_count"])

	if watchPaths, ok := stats["watch_paths"].([]string); ok {
		fmt.Printf("\n📁 Watch paths:\n")
		for _, path := range watchPaths {
			fmt.Printf("   • %s\n", path)
		}
	}

	// List plugin states
	states := hotReload.ListPluginStates()
	if len(states) > 0 {
		fmt.Printf("\n🔧 Plugin States:\n")
		for name, state := range states {
			status := "❌"
			if state.Enabled {
				status = "✅"
			}
			fmt.Printf("   %s %s - reloads: %d, last: %s\n",
				status, name, state.ReloadCount, state.LastReload.Format("15:04:05"))
		}
	}

	return nil
}

// runPluginEnable enables a plugin at runtime
func runPluginEnable(registry *common.PluginRegistry, deps *common.Dependencies, pluginName string) error {
	hotReload := registry.GetHotReloadManager()
	if hotReload == nil {
		return fmt.Errorf("hot-reload manager not initialized")
	}

	fmt.Printf("🔌 Enabling plugin: %s\n", pluginName)

	if err := hotReload.EnablePlugin(pluginName); err != nil {
		return fmt.Errorf("failed to enable plugin %s: %w", pluginName, err)
	}

	fmt.Printf("✅ Plugin %s enabled\n", pluginName)

	return nil
}

// runPluginDisable disables a plugin at runtime
func runPluginDisable(registry *common.PluginRegistry, deps *common.Dependencies, pluginName string) error {
	hotReload := registry.GetHotReloadManager()
	if hotReload == nil {
		return fmt.Errorf("hot-reload manager not initialized")
	}

	fmt.Printf("🔌 Disabling plugin: %s\n", pluginName)

	if err := hotReload.DisablePlugin(pluginName); err != nil {
		return fmt.Errorf("failed to disable plugin %s: %w", pluginName, err)
	}

	fmt.Printf("✅ Plugin %s disabled\n", pluginName)

	return nil
}

// runPluginList lists all plugin states
func runPluginList(registry *common.PluginRegistry, deps *common.Dependencies) error {
	hotReload := registry.GetHotReloadManager()
	if hotReload == nil {
		return fmt.Errorf("hot-reload manager not initialized")
	}

	// Refresh plugin states to ensure they're up to date
	hotReload.RefreshPluginStates()

	states := hotReload.ListPluginStates()

	fmt.Printf("🔧 Plugin Runtime States\n\n")

	if len(states) == 0 {
		fmt.Printf("No plugins found\n")
		return nil
	}

	fmt.Printf("%-15s %-8s %-8s %-10s %-20s\n", "PLUGIN", "ENABLED", "LOADED", "RELOADS", "LAST RELOAD")
	fmt.Printf("%-15s %-8s %-8s %-10s %-20s\n", "------", "-------", "------", "-------", "-----------")

	for name, state := range states {
		enabled := "❌"
		if state.Enabled {
			enabled = "✅"
		}

		loaded := "❌"
		if state.Loaded {
			loaded = "✅"
		}

		lastReload := "never"
		if !state.LastReload.IsZero() {
			lastReload = state.LastReload.Format("15:04:05")
		}

		fmt.Printf("%-15s %-8s %-8s %-10d %-20s\n",
			name, enabled, loaded, state.ReloadCount, lastReload)

		if state.Error != "" {
			fmt.Printf("   Error: %s\n", state.Error)
		}
	}

	fmt.Printf("\n💡 Use 'engx dev plugin enable/disable [name]' to manage plugin states\n")

	return nil
}

// Plugin configuration command handlers
func runConfigList(registry *common.PluginRegistry, deps *common.Dependencies) error {
	configManager := common.NewPluginConfigManager(deps)
	if err := configManager.LoadAll(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to load existing configurations: %v\n", err)
	}

	configUI := common.NewSimpleConfigUI(configManager, deps)
	return configUI.ShowPluginList()
}

func runConfigEdit(registry *common.PluginRegistry, deps *common.Dependencies, pluginName string) error {
	configManager := common.NewPluginConfigManager(deps)
	if err := configManager.LoadAll(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to load existing configurations: %v\n", err)
	}

	configUI := common.NewSimpleConfigUI(configManager, deps)
	return configUI.ShowPluginConfig(pluginName)
}

func runConfigWizard(registry *common.PluginRegistry, deps *common.Dependencies, pluginName string) error {
	configManager := common.NewPluginConfigManager(deps)
	if err := configManager.LoadAll(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to load existing configurations: %v\n", err)
	}

	configUI := common.NewSimpleConfigUI(configManager, deps)
	return configUI.InteractiveConfigWizard(pluginName)
}

func runConfigGet(registry *common.PluginRegistry, deps *common.Dependencies, pluginName, key string) error {
	configManager := common.NewPluginConfigManager(deps)
	if err := configManager.LoadAll(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to load existing configurations: %v\n", err)
	}

	value, err := configManager.GetConfigValue(pluginName, key)
	if err != nil {
		return fmt.Errorf("failed to get config value: %w", err)
	}

	fmt.Printf("%s.%s: %v\n", pluginName, key, value)
	return nil
}

func runConfigSet(registry *common.PluginRegistry, deps *common.Dependencies, pluginName, key, value string) error {
	configManager := common.NewPluginConfigManager(deps)
	if err := configManager.LoadAll(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to load existing configurations: %v\n", err)
	}

	// Try to parse value as different types
	var parsedValue interface{} = value
	if intVal, err := strconv.Atoi(value); err == nil {
		parsedValue = intVal
	} else if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		parsedValue = floatVal
	} else if boolVal, err := strconv.ParseBool(value); err == nil {
		parsedValue = boolVal
	}

	if err := configManager.SetConfigValue(pluginName, key, parsedValue); err != nil {
		return fmt.Errorf("failed to set config value: %w", err)
	}

	fmt.Printf("✅ Set %s.%s = %v\n", pluginName, key, parsedValue)
	return nil
}

func runConfigReset(registry *common.PluginRegistry, deps *common.Dependencies, pluginName string) error {
	configManager := common.NewPluginConfigManager(deps)
	if err := configManager.LoadAll(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to load existing configurations: %v\n", err)
	}

	if err := configManager.ResetConfig(pluginName); err != nil {
		return fmt.Errorf("failed to reset config: %w", err)
	}

	fmt.Printf("✅ Reset configuration for %s to defaults\n", pluginName)
	return nil
}

func runConfigExport(registry *common.PluginRegistry, deps *common.Dependencies, pluginName, format string) error {
	configManager := common.NewPluginConfigManager(deps)
	if err := configManager.LoadAll(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to load existing configurations: %v\n", err)
	}

	filename := fmt.Sprintf("%s-config.%s", pluginName, format)
	if err := configManager.ExportConfig(pluginName, format, filename); err != nil {
		return fmt.Errorf("failed to export config: %w", err)
	}

	fmt.Printf("✅ Exported %s configuration to %s\n", pluginName, filename)
	return nil
}

func runConfigImport(registry *common.PluginRegistry, deps *common.Dependencies, pluginName, filename string) error {
	configManager := common.NewPluginConfigManager(deps)
	if err := configManager.LoadAll(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to load existing configurations: %v\n", err)
	}

	if err := configManager.ImportConfig(pluginName, filename); err != nil {
		return fmt.Errorf("failed to import config: %w", err)
	}

	fmt.Printf("✅ Imported configuration for %s from %s\n", pluginName, filename)
	return nil
}

func runConfigUI(registry *common.PluginRegistry, deps *common.Dependencies) error {
	configManager := common.NewPluginConfigManager(deps)
	if err := configManager.LoadAll(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to load existing configurations: %v\n", err)
	}

	configUI := common.NewSimpleConfigUI(configManager, deps)

	fmt.Println("📦 Plugin Configuration Manager")
	fmt.Println("Use the following commands to manage plugin configurations:")
	fmt.Println("")

	return configUI.ShowPluginList()
}

func runConfigStats(registry *common.PluginRegistry, deps *common.Dependencies) error {
	configManager := common.NewPluginConfigManager(deps)
	if err := configManager.LoadAll(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to load existing configurations: %v\n", err)
	}

	configUI := common.NewSimpleConfigUI(configManager, deps)
	return configUI.ShowStats()
}

func runConfigCreate(registry *common.PluginRegistry, deps *common.Dependencies, pluginName string) error {
	configManager := common.NewPluginConfigManager(deps)
	if err := configManager.LoadAll(); err != nil {
		fmt.Printf("⚠️  Warning: Failed to load existing configurations: %v\n", err)
	}

	configUI := common.NewSimpleConfigUI(configManager, deps)
	return configUI.CreatePluginConfig(pluginName)
}

// Template discovery command handlers (simulation-focused)
func runTemplatesList(registry *common.PluginRegistry, deps *common.Dependencies) error {
	templateUI := common.NewTemplateUI()
	return templateUI.ShowTemplateList()
}

func runTemplatesSearch(registry *common.PluginRegistry, deps *common.Dependencies, query string) error {
	templateUI := common.NewTemplateUI()
	return templateUI.SearchTemplates(query)
}

func runTemplatesInfo(registry *common.PluginRegistry, deps *common.Dependencies, templateID string) error {
	templateUI := common.NewTemplateUI()
	return templateUI.ShowTemplateInfo(templateID)
}

func runTemplatesRecommended(registry *common.PluginRegistry, deps *common.Dependencies) error {
	templateUI := common.NewTemplateUI()
	return templateUI.ShowRecommended()
}

func runTemplatesComplexity(registry *common.PluginRegistry, deps *common.Dependencies, complexityStr string) error {
	templateUI := common.NewTemplateUI()
	return templateUI.ShowByComplexity(complexityStr)
}

func runTemplatesStats(registry *common.PluginRegistry, deps *common.Dependencies) error {
	templateUI := common.NewTemplateUI()
	return templateUI.ShowStats()
}

func runTemplatesDetails(registry *common.PluginRegistry, deps *common.Dependencies, appType string) error {
	templateUI := common.NewTemplateUI()
	return templateUI.ShowArchetypeDetails(appType)
}

// Analytics command handlers for CLI interaction pattern analysis
func runAnalyticsSummary(registry *common.PluginRegistry, deps *common.Dependencies) error {
	analytics := common.NewInteractionAnalytics(deps)
	analyticsUI := common.NewAnalyticsUI(analytics, deps)
	return analyticsUI.ShowSessionSummary()
}

func runAnalyticsDetails(registry *common.PluginRegistry, deps *common.Dependencies) error {
	analytics := common.NewInteractionAnalytics(deps)
	analyticsUI := common.NewAnalyticsUI(analytics, deps)
	return analyticsUI.ShowDetailedAnalytics()
}

func runAnalyticsPatterns(registry *common.PluginRegistry, deps *common.Dependencies) error {
	analytics := common.NewInteractionAnalytics(deps)
	analyticsUI := common.NewAnalyticsUI(analytics, deps)
	return analyticsUI.ShowWorkflowPatterns()
}

func runAnalyticsStats(registry *common.PluginRegistry, deps *common.Dependencies) error {
	analytics := common.NewInteractionAnalytics(deps)
	analyticsUI := common.NewAnalyticsUI(analytics, deps)
	return analyticsUI.ShowUsageStats()
}

func runAnalyticsExport(registry *common.PluginRegistry, deps *common.Dependencies, outputFile string) error {
	analytics := common.NewInteractionAnalytics(deps)
	analyticsUI := common.NewAnalyticsUI(analytics, deps)
	return analyticsUI.ExportAnalytics(outputFile)
}

func runAnalyticsStatus(registry *common.PluginRegistry, deps *common.Dependencies) error {
	analytics := common.NewInteractionAnalytics(deps)
	analyticsUI := common.NewAnalyticsUI(analytics, deps)
	return analyticsUI.ShowAnalyticsStatus()
}
