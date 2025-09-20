package common

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/design"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/styles"
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