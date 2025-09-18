package prompts

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/bthompso/engx-ergonomics-poc/internal/config"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/styles"
)

// ConfigurationSummary displays final configuration for user confirmation
type ConfigurationSummary struct {
	BasePrompt
	config       *config.UserConfiguration
	confirmed    bool
	selectedOption int
	options      []string
}

// NewConfigurationSummary creates a new configuration summary prompt
func NewConfigurationSummary(config *config.UserConfiguration) *ConfigurationSummary {
	options := []string{
		"✅ Create Project",
		"✏️  Edit Configuration",
		"💾 Save Configuration",
		"❌ Cancel",
	}

	return &ConfigurationSummary{
		BasePrompt: BasePrompt{
			title:    "Configuration Summary",
			helpText: "Review your configuration and choose how to proceed",
			required: true,
		},
		config:         config,
		confirmed:      false,
		selectedOption: 0, // Default to "Create Project"
		options:        options,
	}
}

// Init implements PromptComponent
func (cs *ConfigurationSummary) Init() tea.Cmd {
	return nil
}

// Update implements PromptComponent
func (cs *ConfigurationSummary) Update(msg tea.Msg) (PromptComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if cs.selectedOption > 0 {
				cs.selectedOption--
			}
			return cs, nil
		case "down", "j":
			if cs.selectedOption < len(cs.options)-1 {
				cs.selectedOption++
			}
			return cs, nil
		case "enter":
			switch cs.selectedOption {
			case 0: // Create Project
				cs.confirmed = true
				cs.SetCompleted(true)
				return cs, tea.Cmd(func() tea.Msg {
					return CompletePromptMsg{}
				})
			case 1: // Edit Configuration
				return cs, tea.Cmd(func() tea.Msg {
					return PrevPromptMsg{}
				})
			case 2: // Save Configuration
				// TODO: Implement save functionality
				return cs, nil
			case 3: // Cancel
				return cs, tea.Quit
			}
		case "h":
			cs.SetShowHelp(!cs.IsShowingHelp())
			return cs, nil
		case "q", "ctrl+c":
			return cs, tea.Quit
		}
	}

	return cs, nil
}

// View implements PromptComponent
func (cs *ConfigurationSummary) View() string {
	var view strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		MarginBottom(1)

	view.WriteString(headerStyle.Render("📋 Configuration Summary"))
	view.WriteString("\n\n")

	// Configuration details
	summaryStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Padding(1, 2).
		MarginBottom(2)

	view.WriteString(summaryStyle.Render(cs.config.GetSummary()))

	// Estimated setup time
	estimatedTime := config.EstimateSetupTime(*cs.config)
	timeStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		MarginBottom(2)

	view.WriteString(timeStyle.Render(
		fmt.Sprintf("⏱️  Estimated setup time: %s", formatDuration(estimatedTime))))
	view.WriteString("\n\n")

	// Configuration warnings
	warnings := config.ValidateConfiguration(*cs.config)
	if len(warnings) > 0 {
		warningStyle := lipgloss.NewStyle().
			Foreground(styles.Warning).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Warning).
			Padding(1).
			MarginBottom(1)

		warningText := "⚠️  Configuration Warnings:\n"
		for _, warning := range warnings {
			warningText += fmt.Sprintf("• %s\n", warning)
		}

		view.WriteString(warningStyle.Render(warningText))
		view.WriteString("\n")
	}

	// Options
	optionsStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Padding(1)

	optionsText := cs.renderOptions()
	view.WriteString(optionsStyle.Render(optionsText))

	// Help text if showing
	if cs.IsShowingHelp() {
		helpStyle := lipgloss.NewStyle().
			Foreground(styles.Muted).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Muted).
			Padding(1).
			MarginTop(1)

		view.WriteString("\n")
		view.WriteString(helpStyle.Render("💡 " + cs.GetHelp()))
	}

	// Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.Muted).
		MarginTop(1)

	footer := "[Enter] Select • [↑↓] Navigate • [h] Help • [q] Quit"
	view.WriteString("\n")
	view.WriteString(footerStyle.Render(footer))

	return view.String()
}

// GetValue implements PromptComponent
func (cs *ConfigurationSummary) GetValue() interface{} {
	return cs.confirmed
}

// SetValue implements PromptComponent
func (cs *ConfigurationSummary) SetValue(value interface{}) {
	if confirmed, ok := value.(bool); ok {
		cs.confirmed = confirmed
	}
}

// Validate implements PromptComponent
func (cs *ConfigurationSummary) Validate() error {
	// No validation needed for confirmation
	return nil
}

// Helper methods

func (cs *ConfigurationSummary) renderOptions() string {
	var options strings.Builder

	for i, option := range cs.options {
		if i == cs.selectedOption {
			selectedStyle := lipgloss.NewStyle().
				Foreground(styles.Primary).
				Bold(true).
				PaddingLeft(2)

			options.WriteString(selectedStyle.Render("▶ " + option))
		} else {
			normalStyle := lipgloss.NewStyle().
				Foreground(styles.Gray700).
				PaddingLeft(4)

			options.WriteString(normalStyle.Render(option))
		}

		if i < len(cs.options)-1 {
			options.WriteString("\n")
		}
	}

	return options.String()
}

func formatDuration(seconds int) string {
	duration := time.Duration(seconds) * time.Second

	if duration >= time.Minute {
		minutes := int(duration.Minutes())
		remainingSeconds := int(duration.Seconds()) % 60

		if remainingSeconds == 0 {
			return fmt.Sprintf("%d minute%s", minutes, pluralize(minutes))
		}
		return fmt.Sprintf("%d minute%s %d second%s",
			minutes, pluralize(minutes),
			remainingSeconds, pluralize(remainingSeconds))
	}

	return fmt.Sprintf("%d second%s", seconds, pluralize(seconds))
}

func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}