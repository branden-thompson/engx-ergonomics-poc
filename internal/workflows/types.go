package workflows

import (
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/chaos"
	"github.com/bthompso/engx-ergonomics-poc/internal/config"
	"github.com/bthompso/engx-ergonomics-poc/pkg/common"
)

// StageType represents the type of workflow stage
type StageType string

const (
	StageTypeArchetypeSelection   StageType = "archetype_selection"
	StageTypeContextualPrompting  StageType = "contextual_prompting"
	StageTypeTUIExecution        StageType = "tui_execution"
)

// WorkflowStage represents a single step in the workflow
type WorkflowStage interface {
	Execute(ctx *WorkflowContext) (*StageResult, error)
	CanSkip(ctx *WorkflowContext) bool
	GetName() string
	GetDescription() string
}

// WorkflowContext carries state between workflow stages
type WorkflowContext struct {
	// User inputs
	AppName string
	Flags   []string

	// Stage results
	SelectedArchetype *ArchetypeDefinition
	UserConfiguration *config.UserConfiguration
	TUIConfiguration  *TUIConfig

	// System context
	Dependencies    *common.Dependencies
	VerbosityConfig *config.VerbosityConfig
	ChaosInjector   chaos.ChaosInjector

	// Mode tracking
	IsGuidedMode bool
}

// StageResult contains the output of a workflow stage
type StageResult struct {
	StageType       StageType
	Data           map[string]interface{}
	NextStageHints []string
	Errors         []error
}

// MergeResult merges stage result data into the workflow context
func (wc *WorkflowContext) MergeResult(result *StageResult) error {
	if result == nil {
		return nil
	}

	// Merge specific data based on stage type
	switch result.StageType {
	case StageTypeArchetypeSelection:
		if archetype, ok := result.Data["selectedArchetype"]; ok {
			if archetypeDef, ok := archetype.(*ArchetypeDefinition); ok {
				wc.SelectedArchetype = archetypeDef
			}
		}
	case StageTypeContextualPrompting:
		if userConfig, ok := result.Data["userConfiguration"]; ok {
			if config, ok := userConfig.(*config.UserConfiguration); ok {
				wc.UserConfiguration = config
			}
		}
	case StageTypeTUIExecution:
		if tuiConfig, ok := result.Data["tuiConfig"]; ok {
			if config, ok := tuiConfig.(*TUIConfig); ok {
				wc.TUIConfiguration = config
			}
		}
	}

	return nil
}

// ArchetypeDefinition represents an application archetype
type ArchetypeDefinition struct {
	ID          string
	Name        string
	Description string
	IsDefault   bool
	Category    ArchetypeCategory
	PromptSet   string
	TUISteps    []StepDefinition
	Features    []string
}

// ArchetypeCategory represents the category of an archetype
type ArchetypeCategory string

// StepDefinition represents a single step in a workflow
type StepDefinition struct {
	ID          string
	Name        string
	Description string
	Type        StepType
	Duration    time.Duration
}

// StepType represents the type of workflow step
type StepType string

// Archetype categories
const (
	CategoryWebApplication  ArchetypeCategory = "web_application"
	CategoryCLITool         ArchetypeCategory = "cli_tool"
	CategoryBackendService  ArchetypeCategory = "backend_service"
	CategoryPrototype       ArchetypeCategory = "prototype"
	CategoryPlugin          ArchetypeCategory = "plugin"
)

// Step types
const (
	StepTypeInitialization StepType = "initialization"
	StepTypeDependencies   StepType = "dependencies"
	StepTypeConfiguration  StepType = "configuration"
	StepTypeBuild          StepType = "build"
	StepTypeDeployment     StepType = "deployment"
	StepTypeFinalization   StepType = "finalization"
)

// TUIConfig contains configuration for TUI execution
type TUIConfig struct {
	Steps           []StepDefinition
	Theme           string
	ProgressStyle   string
	CompletionHooks []string
	ErrorHandlers   []string
}

// ArchetypeOption represents an archetype option for selection UI
type ArchetypeOption struct {
	ID          string
	Name        string
	Description string
	IsDefault   bool
	Category    ArchetypeCategory
}