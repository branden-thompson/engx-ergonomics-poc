package common

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/design"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/styles"
	"golang.org/x/term"
)

// TemplateUI provides CLI interface for template discovery and selection
// Focused on ergonomic command patterns for React template selection
type TemplateUI struct {
	manager *TemplateManager
}

// NewTemplateUI creates a new template UI
func NewTemplateUI() *TemplateUI {
	return &TemplateUI{
		manager: NewTemplateManager(),
	}
}

// ShowTemplateList displays all available application archetypes with professional styling
func (tui *TemplateUI) ShowTemplateList() error {
	templates := tui.manager.ListTemplates()

	if len(templates) == 0 {
		fmt.Print(components.Warning("No archetypes available.\n"))
		return nil
	}

	// APPLICATION ARCHETYPES header using new header component
	header := components.NewHeader("APPLICATION ARCHETYPES")
	fmt.Print(header.Render())

	// Create flexible archetype table for terminal width awareness
	formatter := components.NewFlexibleArchetypeTable()

	// Archetype data
	archetypes := []struct {
		num        string
		name       string
		framework  string
		language   string
		appType    string
		isDefault  bool
	}{
		{"01.", "PROD Web App", "React Router 7", "TypeScript", "'prod-web'", true},
		{"02.", "DEV ONLY Web App", "React Router 7", "TypeScript", "'dev-web'", false},
		{"03.", "Hackathon Web App", "React Router 7", "TypeScript", "'hackday'", false},
		{"04.", "`engx` command", "Bubble", "GOLANG", "'engx-cmd'", false},
		{"05.", "Standalone CLI", "Bubble", "GOLANG", "'cli'", false},
		{"06.", "Headless Service", "Bubble", "GOLANG", "'service'", false},
		{"07.", "Agent / Sub-Agent", "Variable", "Python", "'agent'", false},
	}

	// Prepare data for width calculation
	data := make([][]string, len(archetypes))
	for i, archetype := range archetypes {
		data[i] = []string{
			archetype.num,
			archetype.name,
			archetype.framework,
			archetype.language,
			archetype.appType,
		}
	}

	// Calculate flexible widths, reserving space for the (Default) badge
	badgeWidth := 11 // "  (Default)" badge width
	columnWidths := formatter.CalculateFlexibleWidths(data, badgeWidth)

	// Print header with calculated widths
	fmt.Printf("%s\n", formatter.FormatHeaderWithWidths(columnWidths))

	// Print data rows using design system colors and calculated widths
	for _, archetype := range archetypes {
		values := []string{
			archetype.num,
			archetype.name,
			archetype.framework,
			archetype.language,
			archetype.appType,
		}

		colors := []string{
			design.ColorDarkGray,    // num
			design.ColorBrightWhite, // name
			design.ColorDarkGray,    // framework
			design.ColorDarkGray,    // language
			design.ColorBrightBlue,  // apptype (option values use light blue)
		}

		rowText := formatter.FormatRowWithWidths(values, colors, columnWidths)

		// Add default indicator using badge component
		if archetype.isDefault {
			badge := components.NewBadge("  (Default)").AsSuccess()
			rowText += badge.Render()
		}

		fmt.Printf("%s\n", rowText)
	}

	// Separator using new component
	separator := components.NewSeparator()
	fmt.Print("\n" + separator.Render() + "\n")

	// Commands section header
	commandHeader := components.NewHeader("Common Commands and Usage")
	fmt.Print(commandHeader.Render())

	// Command examples using smart command formatter
	commandFormatter := components.NewCommandTable()
	cmdFormatterSmart := components.NewCommandFormatter()

	commands := []struct {
		command string
		desc    string
	}{
		{"engx create <AppName> --app-type <app-type>", "Use Archetype for new App"},
		{"engx templates details <app-type>", "Get Archetype"},
		{"engx templates search <query>", "Find specific templates"},
	}

	for _, cmd := range commands {
		// Use smart formatter for the command, regular formatting for description
		formattedCmd := cmdFormatterSmart.FormatCommandInBackticks(cmd.command)
		row := []string{formattedCmd, cmd.desc}
		colors := []string{"", design.ColorDarkGray} // Command already formatted, just color the description
		fmt.Printf("%s\n", commandFormatter.FormatRow(row, colors))
	}

	// Final separator
	fmt.Print("\n" + separator.Render())

	return nil
}

// SearchTemplates searches and displays matching templates
func (tui *TemplateUI) SearchTemplates(query string) error {
	result := tui.manager.SearchTemplates(query)

	fmt.Printf("🔍 Template Search: \"%s\"\n", query)
	fmt.Println(strings.Repeat("=", 50))

	if len(result.Results) == 0 {
		fmt.Println("No templates found matching your query.")
		fmt.Println("")
		fmt.Println("💡 Try broader terms like 'typescript', 'vite', or 'nextjs'")
		fmt.Println("💡 Use 'engx templates list' to see all available templates")
		return nil
	}

	fmt.Printf("Found %d template(s):\n\n", result.Total)

	for i, template := range result.Results {
		tui.displayTemplateDetailed(i+1, template)
		fmt.Println("")
	}

	fmt.Println("💡 Use 'engx templates info <template-id>' for complete details")
	fmt.Println("💡 Use 'engx create MyApp --template <template-id>' to create project")

	return nil
}

// ShowTemplateInfo displays detailed information about a specific template
func (tui *TemplateUI) ShowTemplateInfo(templateID string) error {
	template, err := tui.manager.GetTemplate(templateID)
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	fmt.Printf("📦 Template: %s\n", template.Name)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("ID: %s\n", template.ID)
	fmt.Printf("Description: %s\n", template.Description)
	fmt.Printf("Framework: %s\n", template.Framework)
	fmt.Printf("Language: %s\n", template.Language)
	fmt.Printf("Bundler: %s\n", template.Bundler)
	fmt.Printf("React Version: %s\n", template.Version)
	fmt.Printf("Category: %s\n", template.Category)
	fmt.Printf("Complexity: %s\n", template.Complexity.String())
	fmt.Printf("Setup Time: %s\n", template.SetupTime)

	if template.Recommended {
		fmt.Printf("Status: ⭐ Recommended\n")
	} else if template.Popular {
		fmt.Printf("Status: 🔥 Popular\n")
	}

	if len(template.Features) > 0 {
		fmt.Printf("\n🛠️  Features:\n")
		for _, feature := range template.Features {
			fmt.Printf("   • %s\n", feature)
		}
	}

	if len(template.Tags) > 0 {
		fmt.Printf("\n🏷️  Tags: %s\n", strings.Join(template.Tags, ", "))
	}

	fmt.Printf("\n💡 Create project: engx create MyApp --template %s\n", template.ID)

	return nil
}

// ShowRecommended displays only recommended templates for quick selection
func (tui *TemplateUI) ShowRecommended() error {
	recommended := tui.manager.GetRecommended()

	fmt.Println("⭐ Recommended React Templates")
	fmt.Println(strings.Repeat("=", 50))

	if len(recommended) == 0 {
		fmt.Println("No recommended templates available.")
		return nil
	}

	fmt.Print("Perfect for new React projects:\n\n")

	for i, template := range recommended {
		fmt.Printf("%d. %s\n", i+1, template.Name)
		fmt.Printf("   📝 %s\n", template.Description)
		fmt.Printf("   🔧 %s + %s (%s)\n", template.Framework, template.Language, template.Bundler)
		fmt.Printf("   ⏱️  %s setup • %s complexity\n", template.SetupTime, template.Complexity.String())
		fmt.Printf("   💡 engx create MyApp --template %s\n", template.ID)
		fmt.Println("")
	}

	fmt.Println("💡 Use 'engx templates list' to see all available templates")

	return nil
}

// ShowByComplexity filters and displays templates by complexity level
func (tui *TemplateUI) ShowByComplexity(complexityStr string) error {
	var level ComplexityLevel
	switch strings.ToLower(complexityStr) {
	case "beginner":
		level = ComplexityBeginner
	case "intermediate":
		level = ComplexityIntermediate
	case "advanced":
		level = ComplexityAdvanced
	case "enterprise":
		level = ComplexityEnterprise
	default:
		return fmt.Errorf("invalid complexity level: %s (use: beginner, intermediate, advanced, enterprise)", complexityStr)
	}

	templates := tui.manager.FilterByComplexity(level)

	fmt.Printf("🎯 %s Templates\n", level.String())
	fmt.Println(strings.Repeat("=", 50))

	if len(templates) == 0 {
		fmt.Printf("No %s templates available.\n", strings.ToLower(level.String()))
		return nil
	}

	for i, template := range templates {
		tui.displayTemplateDetailed(i+1, template)
		fmt.Println("")
	}

	return nil
}

// ShowStats displays template statistics for CLI overview
func (tui *TemplateUI) ShowStats() error {
	stats := tui.manager.GetStats()

	fmt.Println("📊 Template Statistics")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("Total Templates: %v\n", stats["total_templates"])
	fmt.Printf("Recommended: %v\n", stats["recommended_count"])
	fmt.Printf("Popular: %v\n", stats["popular_count"])

	if frameworks, ok := stats["frameworks"].(map[string]int); ok && len(frameworks) > 0 {
		fmt.Printf("\n🚀 Frameworks:\n")
		for framework, count := range frameworks {
			fmt.Printf("   %s: %d\n", framework, count)
		}
	}

	if languages, ok := stats["languages"].(map[string]int); ok && len(languages) > 0 {
		fmt.Printf("\n💬 Languages:\n")
		for language, count := range languages {
			fmt.Printf("   %s: %d\n", language, count)
		}
	}

	if complexity, ok := stats["complexity_levels"].(map[string]int); ok && len(complexity) > 0 {
		fmt.Printf("\n🎯 Complexity Levels:\n")
		for level, count := range complexity {
			fmt.Printf("   %s: %d\n", level, count)
		}
	}

	return nil
}

// Helper methods for display formatting

func (tui *TemplateUI) displayTemplateCompact(index int, template *ReactTemplate) {
	icon := ""
	if template.Recommended {
		icon = "⭐"
	} else if template.Popular {
		icon = "🔥"
	}

	fmt.Printf("   %d. %s %s\n", index, icon, template.Name)
	fmt.Printf("      %s + %s (%s) • %s\n",
		template.Framework, template.Language, template.Bundler, template.SetupTime)
}

func (tui *TemplateUI) displayTemplateCompactStyled(index int, template *ReactTemplate) {
	icon := ""
	iconColor := styles.Gray500
	if template.Recommended {
		icon = "⭐"
		iconColor = styles.Success
	} else if template.Popular {
		icon = "🔥"
		iconColor = styles.Warning
	}

	// Main template name line
	indexStyle := lipgloss.NewStyle().
		Foreground(styles.Gray600).
		Bold(true)

	nameStyle := lipgloss.NewStyle().
		Foreground(styles.Gray800).
		Bold(true)

	iconStyle := lipgloss.NewStyle().
		Foreground(iconColor)

	fmt.Printf("   %s %s %s\n",
		indexStyle.Render(fmt.Sprintf("%d.", index)),
		iconStyle.Render(icon),
		nameStyle.Render(template.Name))

	// Details line
	detailStyle := lipgloss.NewStyle().
		Foreground(styles.Gray600).
		MarginLeft(6)

	fmt.Print(detailStyle.Render(fmt.Sprintf("%s + %s (%s) • %s",
		template.Framework, template.Language, template.Bundler, template.SetupTime)))
	fmt.Println()
}

func (tui *TemplateUI) displayTemplateDetailed(index int, template *ReactTemplate) {
	icon := ""
	if template.Recommended {
		icon = "⭐"
	} else if template.Popular {
		icon = "🔥"
	}

	fmt.Printf("%d. %s %s\n", index, icon, template.Name)
	fmt.Printf("   📝 %s\n", template.Description)
	fmt.Printf("   🔧 %s + %s (%s) • %s complexity • %s setup\n",
		template.Framework, template.Language, template.Bundler,
		template.Complexity.String(), template.SetupTime)

	if len(template.Features) > 0 {
		fmt.Printf("   ✨ %s\n", strings.Join(template.Features[:min(3, len(template.Features))], ", "))
	}

	fmt.Printf("   💡 engx create MyApp --template %s\n", template.ID)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// getTerminalWidth returns the current terminal width
func getTerminalWidth() int {
	// Try to get terminal width from stdout
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && width > 0 {
		return width
	}

	// Try COLUMNS environment variable
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width > 0 {
			return width
		}
	}

	// Fallback to reasonable default
	return 80
}

// formatBulletItems creates responsive column layout for bullet items
func formatBulletItems(items []string, terminalWidth int) string {
	if len(items) == 0 {
		return ""
	}

	const indent = "    "
	const bulletPrefix = "- "
	const minColumnWidth = 32
	const columnSpacing = 3

	// Calculate available width for content (minus indent)
	availableWidth := terminalWidth - len(indent)

	// Check if we should use single column (terminal too narrow or few items)
	useSingleColumn := availableWidth < (minColumnWidth*2 + columnSpacing) || len(items) <= 2

	var result strings.Builder

	if useSingleColumn {
		// Single column layout
		for _, item := range items {
			result.WriteString(fmt.Sprintf("%s\033[%sm%s%s\033[0m\n", indent, design.ColorBrightGreen, bulletPrefix, item))
		}
	} else {
		// Two column layout
		columnWidth := (availableWidth - columnSpacing) / 2

		for i := 0; i < len(items); i += 2 {
			leftItem := fmt.Sprintf("\033[%sm%s%s\033[0m", design.ColorBrightGreen, bulletPrefix, items[i])

			if i+1 < len(items) {
				rightItem := fmt.Sprintf("\033[%sm%s%s\033[0m", design.ColorBrightGreen, bulletPrefix, items[i+1])
				result.WriteString(fmt.Sprintf("%s%-*s   %s\n", indent, columnWidth, leftItem, rightItem))
			} else {
				result.WriteString(fmt.Sprintf("%s%s\n", indent, leftItem))
			}
		}
	}

	return result.String()
}

// stripANSI removes ANSI escape sequences from a string to calculate visual length
func stripANSI(str string) string {
	// Regular expression to match ANSI escape codes
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(str, "")
}

// formatCommandTable creates a 3-column table layout for commands
func formatCommandTable(label, command, status, statusColor string, terminalWidth int) string {
	const minSpacing = 2
	const statusMaxWidth = 18 // Width of "(Explicitly Set)" which is the longest
	const labelColumnWidth = 9 // Width of "Command:" to ensure alignment

	// Calculate column widths
	statusWidth := statusMaxWidth

	// Calculate available width for command column (flex)
	usedWidth := labelColumnWidth + statusWidth + (minSpacing * 2) // spacing between columns
	commandWidth := terminalWidth - usedWidth

	// If terminal is too narrow, fall back to simple format
	if commandWidth < 30 {
		return fmt.Sprintf("%s %s \033[%sm%s\033[0m\n", label, command, statusColor, status)
	}

	// Calculate the visual length of the command (without ANSI codes)
	// We need to strip ANSI codes to get the actual visual length
	visualCommand := stripANSI(command)

	// Calculate proper spacing based on actual command length
	spacing := commandWidth - len(visualCommand)
	if spacing < 1 {
		spacing = 1
	}

	// Create the formatted line with proper spacing, ensuring fixed-width label column
	return fmt.Sprintf("%-*s %s%*s \033[%sm%s\033[0m\n",
		labelColumnWidth, label,
		command,
		spacing, "",
		statusColor, status)
}

// ShowArchetypeDetails displays detailed information about a specific archetype by app-type
func (tui *TemplateUI) ShowArchetypeDetails(appType string) error {
	// Map app-types to specific information
	archetypeInfo := map[string]struct {
		name        string
		description string
		framework   string
		language    string
		features    []string
		useCase     string
		setupTime   string
		includes    []string
	}{
		"prod-web": {
			name:        "Production Web Application",
			description: "The default archetype for new internal tool applications, this archtype is set up for production deployment with full CI/CD pipelines and includes all necessary core data services and integrations to connect with critical services such as CREWS, CATALOG, SSO Authentication and gRPC web connectivity. If you're unsure about the archetype to use, this is the option that will prevent your app from having to be promoted later.",
			framework:   "React Router 7",
			language:    "TypeScript",
			features:    []string{"TypeScript", "React Router 7", "Tailwind CSS", "Radix UI", "ShadCN-based UI Design System (SUDS)"},
			useCase:     "Default archetype for new internal tool applications",
			setupTime:   "8-12 minutes",
			includes:    []string{"TrustBridge SSO", "gRPC Web", "CREWS API", "CATALOG API", "GRID/HDFS Access", "GitHub Actions", "ESLint + Prettier", "Vitest Testing", "Storybook", "GitHub Pages", "Microsoft Clarity", "LIX Integration", "Observe (Service Monitoring)", "LCD (LinkedIn Continuous Deployment)", "Docker Ready", "Azure Cloud"},
		},
		"dev-web": {
			name:        "Development Web Application",
			description: "Streamlined React application for rapid development and prototyping. Includes all the functionality of the production web archetype except deployment capabilities. Perfect for local development, testing new features, and rapid iteration without the overhead of production infrastructure setup.",
			framework:   "React Router 7",
			language:    "TypeScript",
			features:    []string{"TypeScript", "React Router 7", "Tailwind CSS", "Radix UI", "ShadCN-based UI Design System (SUDS)"},
			useCase:     "Quick prototypes and development-focused applications",
			setupTime:   "3-5 minutes",
			includes:    []string{"TrustBridge SSO", "gRPC Web", "CREWS API", "CATALOG API", "GRID/HDFS Access", "GitHub Actions", "ESLint + Prettier", "Vitest Testing", "Storybook", "GitHub Pages", "Microsoft Clarity", "LIX Integration", "Observe (Service Monitoring)"},
		},
		"hackday": {
			name:        "Hackathon Web Application",
			description: "Minimal React setup optimized for rapid iteration and experimentation",
			framework:   "React Router 7",
			language:    "TypeScript",
			features:    []string{"Router", "Minimal Setup", "Fast Build"},
			useCase:     "24-48 hour hackathons and rapid proof-of-concepts",
			setupTime:   "1-2 minutes",
			includes:    []string{"Minimal Dependencies", "Quick Start Scripts", "Pre-configured Styling", "Example Components"},
		},
		"engx-cmd": {
			name:        "ENGX Command Tool",
			description: "CLI application using Bubble Tea framework for terminal interfaces",
			framework:   "Bubble Tea (Go)",
			language:    "Go",
			features:    []string{"Terminal UI", "Interactive Forms", "Command Structure", "Plugin System"},
			useCase:     "Command-line tools and terminal-based applications",
			setupTime:   "5-8 minutes",
			includes:    []string{"Cobra CLI Framework", "Bubble Tea TUI", "Configuration Management", "Plugin Architecture"},
		},
		"cli": {
			name:        "Standalone CLI Tool",
			description: "Simple command-line interface application",
			framework:   "Cobra (Go)",
			language:    "Go",
			features:    []string{"Command Structure", "Flag Parsing", "Config Files"},
			useCase:     "Utility tools and command-line applications",
			setupTime:   "3-5 minutes",
			includes:    []string{"Command Framework", "Configuration", "Help System", "Build Scripts"},
		},
		"service": {
			name:        "Headless Service",
			description: "Background service or API without user interface",
			framework:   "Fiber/Gin (Go)",
			language:    "Go",
			features:    []string{"HTTP API", "Background Processing", "Database Integration"},
			useCase:     "APIs, microservices, and background processing services",
			setupTime:   "6-10 minutes",
			includes:    []string{"Web Framework", "Database Integration", "Logging", "Health Checks", "Docker Support"},
		},
		"agent": {
			name:        "AI Agent / Sub-Agent",
			description: "Python-based AI agent with variable framework support",
			framework:   "Variable (LangChain, CrewAI, etc.)",
			language:    "Python",
			features:    []string{"AI Integration", "Agent Framework", "Tool Integration", "Memory Management"},
			useCase:     "AI agents, automation tools, and intelligent systems",
			setupTime:   "4-7 minutes",
			includes:    []string{"Agent Framework", "LLM Integration", "Tool System", "Configuration Management"},
		},
	}

	info, exists := archetypeInfo[appType]
	if !exists {
		return fmt.Errorf("unknown app-type: %s. Available: prod-web, dev-web, hackday, engx-cmd, cli, service, agent", appType)
	}

	// Use the new component-based header
	header := components.NewHeader(fmt.Sprintf("ARCHETYPE: %s", strings.ToUpper(info.name)))
	fmt.Print(header.Render())

	// Get terminal width for responsive layout
	terminalWidth := getTerminalWidth()

	// Create command examples with proper colors and right-aligned status
	var content strings.Builder
	cmdFormatterSmart := components.NewCommandFormatter()

	// Format commands with 3-column table layout
	// Only prod-web shows both default and explicit commands since it's the assumed default
	if appType == "prod-web" {
		defaultCmd := cmdFormatterSmart.FormatCommandInBackticks("engx create <AppName>")
		explicitCmd := cmdFormatterSmart.FormatCommandInBackticks(fmt.Sprintf("engx create <AppName> --app-type %s", appType))

		content.WriteString(formatCommandTable("Command:", defaultCmd, "(Assumed Default)", design.ColorBrightGreen, terminalWidth))
		content.WriteString(formatCommandTable("", explicitCmd, "(Explicitly Set)", design.ColorBrightYellow, terminalWidth))
	} else {
		// All other archetypes only show the explicit command
		explicitCmd := cmdFormatterSmart.FormatCommandInBackticks(fmt.Sprintf("engx create <AppName> --app-type %s", appType))
		content.WriteString(formatCommandTable("Command:", explicitCmd, "(Explicitly Set)", design.ColorBrightYellow, terminalWidth))
	}
	content.WriteString("\n")

	// BLUF section with colors
	content.WriteString(fmt.Sprintf("\033[%smBottom-Line Up Front | BLUF:\033[0m\n", design.ColorBrightMagenta))

	// Format BLUF with command highlighting for dev-web
	var blufText string
	if appType == "dev-web" {
		// Special BLUF text with command highlighting for dev-web
		blufText = fmt.Sprintf("\033[%smStreamlined web-app stack to get going without having to worry about setting up additional configurations and deployment pipelines for publishing. This is a good choice if you need to get going quickly or are actually unsure if you're going to publish this application to a target environment.\n\nIf you determine you need to publish this app, you will need to invoke \033[0m%s\033[%sm to set up the required hooks and CI/CD configurations. Most of the time it is better to start off with a \033[0m%s\033[%sm app-type.\033[0m",
			design.ColorDarkGray,
			cmdFormatterSmart.FormatCommandInBackticks("engx promote <AppName> --production"),
			design.ColorDarkGray,
			cmdFormatterSmart.FormatCommandInBackticks("prod-web"),
			design.ColorDarkGray)
	} else {
		blufText = fmt.Sprintf("\033[%sm%s\033[0m", design.ColorDarkGray, info.description)
	}
	content.WriteString(blufText + "\n\n")

	// Core Technologies
	content.WriteString("• Core Technologies:\n")
	content.WriteString(formatBulletItems(info.features, terminalWidth))
	content.WriteString("\n")

	// Categorize includes into sections
	engxIntegrations := []string{}
	qualityTesting := []string{}
	analyticsMonitoring := []string{}
	deployment := []string{}

	for _, include := range info.includes {
		switch {
		case strings.Contains(include, "SSO") || strings.Contains(include, "gRPC") ||
			 strings.Contains(include, "CREWS") || strings.Contains(include, "CATALOG") ||
			 strings.Contains(include, "GRID") || strings.Contains(include, "GitHub Actions"):
			engxIntegrations = append(engxIntegrations, include)
		case strings.Contains(include, "ESLint") || strings.Contains(include, "Testing") ||
			 strings.Contains(include, "Storybook") || strings.Contains(include, "Pages"):
			qualityTesting = append(qualityTesting, include)
		case strings.Contains(include, "Clarity") || strings.Contains(include, "LIX") ||
			 strings.Contains(include, "Observe"):
			analyticsMonitoring = append(analyticsMonitoring, include)
		case strings.Contains(include, "LCD") || strings.Contains(include, "Docker") ||
			 strings.Contains(include, "Azure"):
			deployment = append(deployment, include)
		}
	}

	// EngX Integrations
	if len(engxIntegrations) > 0 {
		content.WriteString("• EngX Integrations: \n")
		content.WriteString(formatBulletItems(engxIntegrations, terminalWidth))
		content.WriteString("\n")
	}

	// Quality & Testing
	if len(qualityTesting) > 0 {
		content.WriteString("• Quality & Testing:\n")
		content.WriteString(formatBulletItems(qualityTesting, terminalWidth))
		content.WriteString("\n")
	}

	// User Analytics & Testing
	if len(analyticsMonitoring) > 0 {
		content.WriteString("• User Analytics & Testing:\n")
		content.WriteString(formatBulletItems(analyticsMonitoring, terminalWidth))
		content.WriteString("\n")
	}

	// Deployment Capabilities
	if len(deployment) > 0 {
		content.WriteString("• Deployment Capabilities:\n")
		content.WriteString(formatBulletItems(deployment, terminalWidth))
		content.WriteString("\n")
	} else if appType == "dev-web" {
		// Show deployment section for dev-web but indicate none are included
		content.WriteString("• Deployment Capabilities:\n")
		content.WriteString(fmt.Sprintf("    \033[%sm- None included by default (use --app-type prod-web for deployment features)\033[0m\n", design.ColorDarkGray))
		content.WriteString("\n")
	}

	// Print content directly for clean left alignment
	fmt.Print(content.String())

	// Add bottom separator
	separator := components.NewSeparator()
	fmt.Print(separator.Render())

	return nil
}