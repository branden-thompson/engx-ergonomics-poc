package prompts

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/bthompso/engx-ergonomics-poc/internal/config"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/styles"
)

// FeatureChoice represents a selectable feature
type FeatureChoice struct {
	name        string
	description string
	selected    bool
	recommended bool
}

// Implement list.Item interface
func (fc FeatureChoice) FilterValue() string {
	return fc.name
}

func (fc FeatureChoice) Title() string {
	checkbox := "[ ]"
	if fc.selected {
		checkbox = "[✓]"
	}

	title := fmt.Sprintf("%s %s", checkbox, fc.name)
	if fc.recommended {
		title += " (Recommended)"
	}
	return title
}

func (fc FeatureChoice) Description() string {
	return fc.description
}

// FeatureSelector implements multi-select feature selection
type FeatureSelector struct {
	BasePrompt
	list        list.Model
	choices     []FeatureChoice
	category    string
	minSelected int
}

// NewDevFeatureSelector creates a development features selector
func NewDevFeatureSelector() *FeatureSelector {
	choices := []FeatureChoice{
		{name: "Hot Reload", description: "Fast development with auto-refresh", selected: true, recommended: true},
		{name: "ESLint + Prettier", description: "Code quality and formatting", selected: true, recommended: true},
		{name: "Husky + Lint-staged", description: "Pre-commit hooks for code quality", selected: false},
		{name: "VS Code Configuration", description: "Workspace settings and extensions", selected: false},
		{name: "React DevTools", description: "Browser debugging extensions", selected: true, recommended: true},
		{name: "Storybook", description: "Component development environment", selected: false},
	}

	return newFeatureSelector("Development Features", choices, 0)
}

// NewProductionFeatureSelector creates a production features selector
func NewProductionFeatureSelector() *FeatureSelector {
	choices := []FeatureChoice{
		{name: "Docker", description: "Containerization for consistent deployments", selected: false},
		{name: "CI/CD Pipeline", description: "Automated testing and deployment", selected: false},
		{name: "Monitoring", description: "Error tracking and performance monitoring", selected: false},
		{name: "Analytics", description: "User behavior and performance analytics", selected: false},
	}

	return newFeatureSelector("Production Setup", choices, 0)
}

// NewTestingFeatureSelector creates a testing features selector
func NewTestingFeatureSelector() *FeatureSelector {
	choices := []FeatureChoice{
		{name: "Unit Testing", description: "Jest + React Testing Library", selected: true, recommended: true},
		{name: "E2E Testing", description: "Cypress or Playwright", selected: false},
		{name: "Coverage Reports", description: "Code coverage tracking", selected: false},
	}

	return newFeatureSelector("Testing Setup", choices, 0)
}

// NewNavigationSelector creates a navigation configuration selector
func NewNavigationSelector() *FeatureSelector {
	choices := []FeatureChoice{
		{name: "Federated Global Nav & Chrome", description: "Use shared navigation with global chrome templates", selected: false},
		{name: "Standalone App Header & Chrome", description: "Use standalone app header and chrome templates", selected: true, recommended: true},
	}

	return newFeatureSelector("Navigation Configuration", choices, 1) // Require exactly one selection
}

// newFeatureSelector creates a new feature selector with the given configuration
func newFeatureSelector(category string, choices []FeatureChoice, minSelected int) *FeatureSelector {
	// Convert to list items
	items := make([]list.Item, len(choices))
	for i, choice := range choices {
		items[i] = choice
	}

	// Create list
	l := list.New(items, NewFeatureDelegate(), 60, len(choices)+4)
	l.Title = category
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	return &FeatureSelector{
		BasePrompt: BasePrompt{
			title:    category,
			helpText: fmt.Sprintf("Select the %s you want to include in your project. Use Space to toggle selection, Enter to continue.", strings.ToLower(category)),
			required: minSelected > 0,
		},
		list:        l,
		choices:     choices,
		category:    category,
		minSelected: minSelected,
	}
}

// Init implements PromptComponent
func (fs *FeatureSelector) Init() tea.Cmd {
	return nil
}

// Update implements PromptComponent
func (fs *FeatureSelector) Update(msg tea.Msg) (PromptComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ", "spacebar":
			// Toggle selection
			if fs.list.Index() >= 0 && fs.list.Index() < len(fs.choices) {
				fs.choices[fs.list.Index()].selected = !fs.choices[fs.list.Index()].selected
				// Update the list item
				fs.list.SetItem(fs.list.Index(), fs.choices[fs.list.Index()])
			}
			return fs, nil
		case "enter":
			if err := fs.Validate(); err != nil {
				fs.SetError(err)
				return fs, nil
			}
			fs.SetCompleted(true)
			return fs, tea.Cmd(func() tea.Msg {
				return CompletePromptMsg{}
			})
		case "a":
			// Select all
			for i := range fs.choices {
				fs.choices[i].selected = true
				fs.list.SetItem(i, fs.choices[i])
			}
			return fs, nil
		case "n":
			// Select none
			for i := range fs.choices {
				fs.choices[i].selected = false
				fs.list.SetItem(i, fs.choices[i])
			}
			return fs, nil
		case "h":
			fs.SetShowHelp(!fs.IsShowingHelp())
			return fs, nil
		case "q", "ctrl+c":
			return fs, tea.Quit
		}
	}

	var cmd tea.Cmd
	fs.list, cmd = fs.list.Update(msg)
	return fs, cmd
}

// View implements PromptComponent
func (fs *FeatureSelector) View() string {
	var view strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		MarginBottom(1)

	icon := "🛠️"
	if fs.category == "Production Setup" {
		icon = "🚀"
	} else if fs.category == "Testing Setup" {
		icon = "🧪"
	} else if fs.category == "Navigation Configuration" {
		icon = "🧭"
	}

	view.WriteString(headerStyle.Render(icon + " " + fs.list.Title))
	view.WriteString("\n")

	// List
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Padding(1).
		MarginBottom(1)

	view.WriteString(listStyle.Render(fs.list.View()))

	// Selection count
	selectedCount := fs.getSelectedCount()
	countStyle := lipgloss.NewStyle().
		Foreground(styles.Primary).
		MarginBottom(1)

	view.WriteString(countStyle.Render(fmt.Sprintf("💡 Selected: %d/%d features", selectedCount, len(fs.choices))))
	view.WriteString("\n")

	// Error display
	if fs.GetError() != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.Error).
			MarginBottom(1)

		view.WriteString(errorStyle.Render("❌ " + fs.GetError().Error()))
		view.WriteString("\n")
	}

	// Help text if showing
	if fs.IsShowingHelp() {
		helpStyle := lipgloss.NewStyle().
			Foreground(styles.Muted).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Muted).
			Padding(1).
			MarginTop(1)

		view.WriteString(helpStyle.Render("💡 " + fs.GetHelp()))
		view.WriteString("\n")
	}

	// Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.Muted).
		MarginTop(1)

	footer := "[Space] Toggle • [Enter] Continue • [↑↓] Navigate • [a] All • [n] None • [h] Help"
	view.WriteString(footerStyle.Render(footer))

	return view.String()
}

// GetValue implements PromptComponent - returns the appropriate config struct
func (fs *FeatureSelector) GetValue() interface{} {
	switch fs.category {
	case "Development Features":
		return config.DevFeatureConfig{
			HotReload:    fs.isSelected("Hot Reload"),
			Linting:      fs.isSelected("ESLint + Prettier"),
			Prettier:     fs.isSelected("ESLint + Prettier"), // Same as linting
			Husky:        fs.isSelected("Husky + Lint-staged"),
			VSCodeConfig: fs.isSelected("VS Code Configuration"),
			DevTools:     fs.isSelected("React DevTools"),
		}
	case "Production Setup":
		return config.ProductionConfig{
			Docker:     fs.isSelected("Docker"),
			CI_CD:      fs.isSelected("CI/CD Pipeline"),
			Monitoring: fs.isSelected("Monitoring"),
			Analytics:  fs.isSelected("Analytics"),
		}
	case "Testing Setup":
		return config.TestingConfig{
			UnitTesting: fs.isSelected("Unit Testing"),
			E2ETesting:  fs.isSelected("E2E Testing"),
			Coverage:    fs.isSelected("Coverage Reports"),
		}
	case "Navigation Configuration":
		return config.NavigationConfig{
			UseFederatedNav: fs.isSelected("Federated Global Nav & Chrome"),
		}
	}
	return nil
}

// SetValue implements PromptComponent
func (fs *FeatureSelector) SetValue(value interface{}) {
	switch v := value.(type) {
	case config.DevFeatureConfig:
		fs.setSelected("Hot Reload", v.HotReload)
		fs.setSelected("ESLint + Prettier", v.Linting)
		fs.setSelected("Husky + Lint-staged", v.Husky)
		fs.setSelected("VS Code Configuration", v.VSCodeConfig)
		fs.setSelected("React DevTools", v.DevTools)
	case config.ProductionConfig:
		fs.setSelected("Docker", v.Docker)
		fs.setSelected("CI/CD Pipeline", v.CI_CD)
		fs.setSelected("Monitoring", v.Monitoring)
		fs.setSelected("Analytics", v.Analytics)
	case config.TestingConfig:
		fs.setSelected("Unit Testing", v.UnitTesting)
		fs.setSelected("E2E Testing", v.E2ETesting)
		fs.setSelected("Coverage Reports", v.Coverage)
	case config.NavigationConfig:
		fs.setSelected("Federated Global Nav & Chrome", v.UseFederatedNav)
		fs.setSelected("Standalone App Header & Chrome", !v.UseFederatedNav)
	}
}

// Validate implements PromptComponent
func (fs *FeatureSelector) Validate() error {
	selectedCount := fs.getSelectedCount()
	if selectedCount < fs.minSelected {
		return fmt.Errorf("please select at least %d feature(s)", fs.minSelected)
	}
	return nil
}

// Helper methods
func (fs *FeatureSelector) isSelected(name string) bool {
	for _, choice := range fs.choices {
		if choice.name == name {
			return choice.selected
		}
	}
	return false
}

func (fs *FeatureSelector) setSelected(name string, selected bool) {
	for i, choice := range fs.choices {
		if choice.name == name {
			fs.choices[i].selected = selected
			fs.list.SetItem(i, fs.choices[i])
			break
		}
	}
}

func (fs *FeatureSelector) getSelectedCount() int {
	count := 0
	for _, choice := range fs.choices {
		if choice.selected {
			count++
		}
	}
	return count
}

// NewFeatureDelegate creates a custom delegate for feature items
func NewFeatureDelegate() list.DefaultDelegate {
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