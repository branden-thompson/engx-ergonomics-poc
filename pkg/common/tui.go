package common

import (
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
)

// TUIRegistry implements the interfaces.TUIRegistry interface
type TUIRegistry struct {
	components    map[string]*TUIComponent
	styles        map[string]*TUIStyle
	prompts       map[string]*TUIPrompt
	logger        interfaces.Logger
	mutex         sync.RWMutex
}

// TUIComponent represents a registered TUI component
type TUIComponent struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`         // "progress", "prompt", "renderer", "style"
	Category    string                 `json:"category"`     // "installation", "validation", "interaction"
	Status      string                 `json:"status"`       // "available", "active", "disabled"
	Config      map[string]interface{} `json:"config"`
	Factory     ComponentFactory       `json:"-"`
}

// TUIStyle represents a style configuration
type TUIStyle struct {
	Name        string                 `json:"name"`
	Theme       string                 `json:"theme"`        // "auto", "dark", "light"
	Colors      map[string]string      `json:"colors"`
	Attributes  map[string]interface{} `json:"attributes"`
}

// TUIPrompt represents a prompt configuration
type TUIPrompt struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`         // "confirmation", "selection", "input", "multi-select"
	Config      map[string]interface{} `json:"config"`
	Handler     PromptHandler          `json:"-"`
}

// ComponentFactory creates TUI components
type ComponentFactory func(config map[string]interface{}) (interface{}, error)

// PromptHandler handles prompt interactions
type PromptHandler func(config map[string]interface{}) (interface{}, error)

// NewTUIRegistry creates a new TUI registry instance
func NewTUIRegistry(logger interfaces.Logger) interfaces.TUIRegistry {
	registry := &TUIRegistry{
		components: make(map[string]*TUIComponent),
		styles:     make(map[string]*TUIStyle),
		prompts:    make(map[string]*TUIPrompt),
		logger:     logger.WithComponent("tui"),
		mutex:      sync.RWMutex{},
	}

	// Register default components
	registry.registerDefaultComponents()
	registry.registerDefaultStyles()
	registry.registerDefaultPrompts()

	return registry
}

// RegisterComponent registers a TUI component
func (tr *TUIRegistry) RegisterComponent(name, componentType string, config map[string]interface{}) error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	if _, exists := tr.components[name]; exists {
		return fmt.Errorf("component %s already registered", name)
	}

	component := &TUIComponent{
		Name:     name,
		Type:     componentType,
		Category: getComponentCategory(componentType),
		Status:   "available",
		Config:   config,
	}

	tr.components[name] = component
	tr.logger.Debug("Registered TUI component: %s (type: %s)", name, componentType)
	return nil
}

// GetComponent retrieves a TUI component by name
func (tr *TUIRegistry) GetComponent(name string) (map[string]interface{}, bool) {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	component, exists := tr.components[name]
	if !exists {
		return nil, false
	}

	result := map[string]interface{}{
		"name":     component.Name,
		"type":     component.Type,
		"category": component.Category,
		"status":   component.Status,
		"config":   component.Config,
	}

	return result, true
}

// ListComponents returns all registered components
func (tr *TUIRegistry) ListComponents() []map[string]interface{} {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	var components []map[string]interface{}
	for _, component := range tr.components {
		components = append(components, map[string]interface{}{
			"name":     component.Name,
			"type":     component.Type,
			"category": component.Category,
			"status":   component.Status,
		})
	}

	return components
}

// SetTheme sets the current TUI theme
func (tr *TUIRegistry) SetTheme(theme string) error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	validThemes := map[string]bool{
		"auto":  true,
		"light": true,
		"dark":  true,
	}

	if !validThemes[theme] {
		return fmt.Errorf("invalid theme: %s (valid options: auto, light, dark)", theme)
	}

	// Update all styles to use the new theme
	for _, style := range tr.styles {
		style.Theme = theme
	}

	tr.logger.Info("TUI theme set to: %s", theme)
	return nil
}

// GetTheme returns the current TUI theme
func (tr *TUIRegistry) GetTheme() *interfaces.Theme {
	return &interfaces.Theme{}
}

// CreateProgressIndicator creates a progress indicator component
func (tr *TUIRegistry) CreateProgressIndicator(progressType string, config map[string]interface{}) (interface{}, error) {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	// Look for a registered progress component
	for _, component := range tr.components {
		if component.Type == "progress" && component.Status == "available" {
			if component.Factory != nil {
				return component.Factory(config)
			}
		}
	}

	// Fallback to default progress indicator
	return tr.createDefaultProgressIndicator(progressType, config)
}

// CreatePrompt creates a prompt component
func (tr *TUIRegistry) CreatePrompt(promptType string, config map[string]interface{}) (interface{}, error) {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	prompt, exists := tr.prompts[promptType]
	if !exists {
		return nil, fmt.Errorf("unknown prompt type: %s", promptType)
	}

	if prompt.Handler != nil {
		return prompt.Handler(config)
	}

	// Fallback to default prompt
	return tr.createDefaultPrompt(promptType, config)
}

// GetAvailableStyles returns all available style configurations
func (tr *TUIRegistry) GetAvailableStyles() []string {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	var styles []string
	for name := range tr.styles {
		styles = append(styles, name)
	}
	return styles
}

// ApplyStyle applies a style configuration by name
func (tr *TUIRegistry) ApplyStyle(styleName string) error {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	style, exists := tr.styles[styleName]
	if !exists {
		return fmt.Errorf("style %s not found", styleName)
	}

	tr.logger.Debug("Applied TUI style: %s (theme: %s)", styleName, style.Theme)
	return nil
}

// registerDefaultComponents registers built-in TUI components
func (tr *TUIRegistry) registerDefaultComponents() {
	defaultComponents := map[string]*TUIComponent{
		"npm-style-progress": {
			Name:     "npm-style-progress",
			Type:     "progress",
			Category: "installation",
			Status:   "available",
			Config: map[string]interface{}{
				"format": "npm",
				"colors": true,
				"spinner": true,
			},
		},
		"enhanced-renderer": {
			Name:     "enhanced-renderer",
			Type:     "renderer",
			Category: "interaction",
			Status:   "available",
			Config: map[string]interface{}{
				"detailed_output": true,
				"color_support":   true,
			},
		},
		"component-manager": {
			Name:     "component-manager",
			Type:     "manager",
			Category: "installation",
			Status:   "available",
			Config: map[string]interface{}{
				"phases": []string{"dependencies", "structure", "testing", "documentation", "finalizing"},
			},
		},
	}

	for name, component := range defaultComponents {
		tr.components[name] = component
	}
}

// registerDefaultStyles registers built-in style configurations
func (tr *TUIRegistry) registerDefaultStyles() {
	defaultStyles := map[string]*TUIStyle{
		"default": {
			Name:  "default",
			Theme: "auto",
			Colors: map[string]string{
				"primary":   "#3b82f6",
				"secondary": "#6b7280",
				"success":   "#10b981",
				"warning":   "#f59e0b",
				"error":     "#ef4444",
			},
			Attributes: map[string]interface{}{
				"bold":      true,
				"underline": false,
				"italic":    false,
			},
		},
		"minimal": {
			Name:  "minimal",
			Theme: "auto",
			Colors: map[string]string{
				"primary":   "#000000",
				"secondary": "#666666",
				"success":   "#008000",
				"warning":   "#ff8800",
				"error":     "#ff0000",
			},
			Attributes: map[string]interface{}{
				"bold":      false,
				"underline": false,
				"italic":    false,
			},
		},
	}

	for name, style := range defaultStyles {
		tr.styles[name] = style
	}
}

// registerDefaultPrompts registers built-in prompt configurations
func (tr *TUIRegistry) registerDefaultPrompts() {
	defaultPrompts := map[string]*TUIPrompt{
		"confirmation": {
			Name: "confirmation",
			Type: "confirmation",
			Config: map[string]interface{}{
				"default_yes": true,
				"show_help":   false,
			},
		},
		"template-selector": {
			Name: "template-selector",
			Type: "selection",
			Config: map[string]interface{}{
				"allow_multiple": false,
				"show_preview":   true,
			},
		},
		"feature-selector": {
			Name: "feature-selector",
			Type: "multi-select",
			Config: map[string]interface{}{
				"allow_multiple": true,
				"show_description": true,
			},
		},
	}

	for name, prompt := range defaultPrompts {
		tr.prompts[name] = prompt
	}
}

// getComponentCategory determines the category for a component type
func getComponentCategory(componentType string) string {
	categoryMap := map[string]string{
		"progress":  "installation",
		"prompt":    "interaction",
		"renderer":  "interaction",
		"style":     "visual",
		"manager":   "installation",
		"validator": "validation",
	}

	if category, exists := categoryMap[componentType]; exists {
		return category
	}
	return "general"
}

// createDefaultProgressIndicator creates a fallback progress indicator
func (tr *TUIRegistry) createDefaultProgressIndicator(progressType string, config map[string]interface{}) (interface{}, error) {
	return map[string]interface{}{
		"type":   progressType,
		"config": config,
		"status": "created",
	}, nil
}

// createDefaultPrompt creates a fallback prompt
func (tr *TUIRegistry) createDefaultPrompt(promptType string, config map[string]interface{}) (interface{}, error) {
	return map[string]interface{}{
		"type":   promptType,
		"config": config,
		"status": "created",
	}, nil
}

// Cleanup performs cleanup of TUI resources
func (tr *TUIRegistry) Cleanup() error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	// Clear all registrations
	tr.components = make(map[string]*TUIComponent)
	tr.styles = make(map[string]*TUIStyle)
	tr.prompts = make(map[string]*TUIPrompt)

	tr.logger.Debug("TUI registry cleaned up")
	return nil
}

// CreateProgressModel creates a progress model
func (tr *TUIRegistry) CreateProgressModel(steps []interfaces.Step) tea.Model {
	// Return a simple mock model for now
	return &mockModel{content: "Progress Model"}
}

// CreateConfirmModel creates a confirmation model
func (tr *TUIRegistry) CreateConfirmModel(prompt string) tea.Model {
	// Return a simple mock model for now
	return &mockModel{content: prompt}
}

// CreateInputModel creates an input model
func (tr *TUIRegistry) CreateInputModel(prompt, placeholder string) tea.Model {
	// Return a simple mock model for now
	return &mockModel{content: prompt + " " + placeholder}
}

// CreateSelectModel creates a selection model
func (tr *TUIRegistry) CreateSelectModel(prompt string, options []string) tea.Model {
	// Return a simple mock model for now
	return &mockModel{content: prompt}
}


// GetComponents returns the component library
func (tr *TUIRegistry) GetComponents() *interfaces.ComponentLibrary {
	return &interfaces.ComponentLibrary{}
}

// GetStyles returns the style registry
func (tr *TUIRegistry) GetStyles() *interfaces.StyleRegistry {
	return &interfaces.StyleRegistry{}
}

// RunModel runs a model and returns the final state
func (tr *TUIRegistry) RunModel(model tea.Model) (tea.Model, error) {
	return model, nil
}

// RunProgram runs a bubbletea program
func (tr *TUIRegistry) RunProgram(model tea.Model) error {
	return nil
}

// mockModel is a simple mock model for testing
type mockModel struct {
	content string
}

func (m *mockModel) Init() tea.Cmd {
	return nil
}

func (m *mockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *mockModel) View() string {
	return m.content
}