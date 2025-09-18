package prompts

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/bthompso/engx-ergonomics-poc/internal/config"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/styles"
)

// TemplateChoice represents a template option
type TemplateChoice struct {
	template config.TemplateType
	selected bool
}

// Implement list.Item interface
func (tc TemplateChoice) FilterValue() string {
	return tc.template.DisplayName()
}

func (tc TemplateChoice) Title() string {
	title := tc.template.DisplayName()
	if tc.template == config.TypeScript {
		title += " (Recommended)"
	}
	return title
}

func (tc TemplateChoice) Description() string {
	return tc.template.Description()
}

// TemplateSelector implements template selection prompt
type TemplateSelector struct {
	BasePrompt
	list     list.Model
	choices  []TemplateChoice
	selected int
}

// NewTemplateSelector creates a new template selector
func NewTemplateSelector() *TemplateSelector {
	choices := []TemplateChoice{
		{template: config.TypeScript},
		{template: config.JavaScript},
		{template: config.Minimal},
	}

	// Convert to list items
	items := make([]list.Item, len(choices))
	for i, choice := range choices {
		items[i] = choice
	}

	// Create list
	l := list.New(items, NewTemplateDelegate(), 50, 8)
	l.Title = "Choose Your Template"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false) // We'll handle help ourselves

	// Set TypeScript as default selection
	l.Select(0)

	return &TemplateSelector{
		BasePrompt: BasePrompt{
			title:    "Template Selection",
			helpText: "Choose the project template that best fits your needs. TypeScript is recommended for most projects as it provides better IDE support and catches errors at compile time.",
			required: true,
		},
		list:     l,
		choices:  choices,
		selected: 0,
	}
}

// Init implements PromptComponent
func (ts *TemplateSelector) Init() tea.Cmd {
	return nil
}

// Update implements PromptComponent
func (ts *TemplateSelector) Update(msg tea.Msg) (PromptComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			ts.selected = ts.list.Index()
			ts.SetCompleted(true)
			return ts, tea.Cmd(func() tea.Msg {
				return CompletePromptMsg{}
			})
		case "h":
			ts.SetShowHelp(!ts.IsShowingHelp())
			return ts, nil
		case "q", "ctrl+c":
			return ts, tea.Quit
		}
	}

	var cmd tea.Cmd
	ts.list, cmd = ts.list.Update(msg)
	return ts, cmd
}

// View implements PromptComponent
func (ts *TemplateSelector) View() string {
	var view strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		MarginBottom(1)

	view.WriteString(headerStyle.Render("🎯 " + ts.list.Title))
	view.WriteString("\n")

	// List
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Padding(1).
		MarginBottom(1)

	view.WriteString(listStyle.Render(ts.list.View()))

	// Help text if showing
	if ts.IsShowingHelp() {
		helpStyle := lipgloss.NewStyle().
			Foreground(styles.Muted).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Muted).
			Padding(1).
			MarginTop(1)

		view.WriteString("\n")
		view.WriteString(helpStyle.Render("💡 " + ts.GetHelp()))
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
func (ts *TemplateSelector) GetValue() interface{} {
	if ts.selected >= 0 && ts.selected < len(ts.choices) {
		return ts.choices[ts.selected].template
	}
	return config.TypeScript // Default fallback
}

// SetValue implements PromptComponent
func (ts *TemplateSelector) SetValue(value interface{}) {
	if template, ok := value.(config.TemplateType); ok {
		for i, choice := range ts.choices {
			if choice.template == template {
				ts.selected = i
				ts.list.Select(i)
				break
			}
		}
	}
}

// Validate implements PromptComponent
func (ts *TemplateSelector) Validate() error {
	return ValidateRequired(ts.required, ts.GetValue())
}

// NewTemplateDelegate creates a custom list delegate for template items
func NewTemplateDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(styles.Primary).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(styles.Primary).
		PaddingLeft(1)

	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(styles.Muted).
		PaddingLeft(2)

	d.Styles.NormalTitle = d.Styles.NormalTitle.
		Foreground(styles.Gray700).
		PaddingLeft(2)

	d.Styles.NormalDesc = d.Styles.NormalDesc.
		Foreground(styles.Muted).
		PaddingLeft(2)

	return d
}