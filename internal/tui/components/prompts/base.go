package prompts

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

// PromptComponent defines the interface for all prompt components
type PromptComponent interface {
	Init() tea.Cmd
	Update(tea.Msg) (PromptComponent, tea.Cmd)
	View() string
	GetValue() interface{}
	SetValue(interface{})
	Validate() error
	IsComplete() bool
	GetTitle() string
	GetHelp() string
}

// PromptType defines the different types of prompts available
type PromptType int

const (
	PromptTypeTemplate PromptType = iota
	PromptTypeDevFeatures
	PromptTypeProductionSetup
	PromptTypeTesting
	PromptTypeConfirmation
)

// PromptStep represents a single step in the prompting flow
type PromptStep struct {
	ID          string
	Title       string
	Component   PromptComponent
	Required    bool
	HelpText    string
	Type        PromptType
}

// NavigationState tracks the current navigation state
type NavigationState struct {
	CanGoBack    bool
	CanGoForward bool
	CanSkip      bool
	ShowHelp     bool
	CurrentStep  int
	TotalSteps   int
}

// PromptResult represents the result of a completed prompt
type PromptResult struct {
	StepID    string
	Value     interface{}
	Skipped   bool
	Completed bool
}

// Common prompt messages
type (
	// Navigation messages
	NextPromptMsg    struct{}
	PrevPromptMsg    struct{}
	SkipPromptMsg    struct{}
	CompletePromptMsg struct{}

	// Help messages
	ShowHelpMsg      struct{}
	HideHelpMsg      struct{}

	// Validation messages
	ValidationErrorMsg struct {
		Error error
	}

	// Value change messages
	ValueChangedMsg struct {
		StepID string
		Value  interface{}
	}
)

// BasePrompt provides common functionality for all prompts
type BasePrompt struct {
	title       string
	helpText    string
	required    bool
	completed   bool
	showHelp    bool
	error       error
}

// GetTitle returns the prompt title
func (bp *BasePrompt) GetTitle() string {
	return bp.title
}

// GetHelp returns the help text
func (bp *BasePrompt) GetHelp() string {
	return bp.helpText
}

// IsComplete returns whether the prompt is completed
func (bp *BasePrompt) IsComplete() bool {
	return bp.completed
}

// SetCompleted marks the prompt as completed
func (bp *BasePrompt) SetCompleted(completed bool) {
	bp.completed = completed
}

// SetError sets an error for the prompt
func (bp *BasePrompt) SetError(err error) {
	bp.error = err
}

// GetError returns the current error
func (bp *BasePrompt) GetError() error {
	return bp.error
}

// SetShowHelp toggles help display
func (bp *BasePrompt) SetShowHelp(show bool) {
	bp.showHelp = show
}

// IsShowingHelp returns whether help is currently shown
func (bp *BasePrompt) IsShowingHelp() bool {
	return bp.showHelp
}

// ValidateRequired checks if a required prompt has a value
func ValidateRequired(required bool, value interface{}) error {
	if !required {
		return nil
	}

	if value == nil {
		return ErrRequiredField
	}

	// Check for empty strings
	if str, ok := value.(string); ok && str == "" {
		return ErrRequiredField
	}

	return nil
}

// Common errors
var (
	ErrRequiredField = fmt.Errorf("this field is required")
	ErrInvalidValue  = fmt.Errorf("invalid value provided")
)