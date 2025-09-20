package common

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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

// ShowTemplateList displays all available React templates with professional styling
func (tui *TemplateUI) ShowTemplateList() error {
	templates := tui.manager.ListTemplates()

	// Header with styled title
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		MarginBottom(1)

	fmt.Print(headerStyle.Render("🚀 React Templates"))
	fmt.Println()

	dividerStyle := lipgloss.NewStyle().
		Foreground(styles.Gray400)
	fmt.Print(dividerStyle.Render(strings.Repeat("─", 70)))
	fmt.Println()

	if len(templates) == 0 {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.Warning).
			MarginTop(1)
		fmt.Print(errorStyle.Render("No templates available."))
		fmt.Println()
		return nil
	}

	// Show recommended templates first
	recommended := tui.manager.GetRecommended()
	if len(recommended) > 0 {
		sectionStyle := lipgloss.NewStyle().
			Foreground(styles.Success).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)
		fmt.Print(sectionStyle.Render("\n⭐ Recommended for new projects:"))
		fmt.Println()

		for i, template := range recommended {
			tui.displayTemplateCompactStyled(i+1, template)
		}
	}

	// Show popular templates
	popular := tui.manager.GetPopular()
	if len(popular) > 0 {
		sectionStyle := lipgloss.NewStyle().
			Foreground(styles.Warning).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)
		fmt.Print(sectionStyle.Render("\n🔥 Popular choices:"))
		fmt.Println()

		popularCount := 0
		for _, template := range popular {
			if !template.Recommended { // Don't duplicate recommended ones
				popularCount++
				tui.displayTemplateCompactStyled(popularCount, template)
			}
		}
	}

	// Show all templates in a styled table
	tableHeaderStyle := lipgloss.NewStyle().
		Foreground(styles.Secondary).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)
	fmt.Print(tableHeaderStyle.Render("\n📋 All Templates:"))
	fmt.Println()

	// Table header with styling
	headerRowStyle := lipgloss.NewStyle().
		Foreground(styles.Gray600).
		Bold(true)

	fmt.Print(headerRowStyle.Render(fmt.Sprintf("%-3s %-25s %-15s %-12s %-10s %-8s",
		"#", "NAME", "FRAMEWORK", "LANGUAGE", "COMPLEXITY", "SETUP")))
	fmt.Println()

	fmt.Print(dividerStyle.Render(strings.Repeat("─", 70)))
	fmt.Println()

	// Table rows with alternating styling
	for i, template := range templates {
		complexity := template.Complexity.String()
		if len(complexity) > 10 {
			complexity = complexity[:10]
		}

		rowStyle := lipgloss.NewStyle().
			Foreground(styles.Gray700)
		if i%2 == 1 {
			rowStyle = rowStyle.Foreground(styles.Gray600)
		}

		fmt.Print(rowStyle.Render(fmt.Sprintf("%-3d %-25s %-15s %-12s %-10s %-8s",
			i+1,
			truncateString(template.Name, 25),
			template.Framework,
			template.Language,
			complexity,
			template.SetupTime)))
		fmt.Println()
	}

	// Footer with helpful tips
	tipStyle := lipgloss.NewStyle().
		Foreground(styles.Info).
		MarginTop(1)

	fmt.Println()
	fmt.Print(tipStyle.Render("💡 Use 'engx templates info <template-id>' for details"))
	fmt.Println()
	fmt.Print(tipStyle.Render("💡 Use 'engx templates search <query>' to find specific templates"))
	fmt.Println()
	fmt.Print(tipStyle.Render("💡 Use 'engx create MyApp --template <template-id>' to use a template"))
	fmt.Println()

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

	fmt.Println("Perfect for new React projects:\n")

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