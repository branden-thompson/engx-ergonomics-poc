# Modular Workflow Design Specification
**🛩️ BRTOPS v1.1.001** | Feature: create-cmd-enhancements | SEV-0

## Workflow Architecture Overview

### Design Philosophy
**Composable Stages**: Each workflow stage is an independent, testable component with clean interfaces and single responsibility.

**State Management**: Immutable workflow context passed between stages with additive state accumulation.

**Extensibility**: New archetypes and workflow stages can be added without modifying existing code.

## Core Workflow Pattern

### Workflow Orchestration Interface
```go
// WorkflowOrchestrator manages the execution of workflow stages
type WorkflowOrchestrator struct {
    stages  []WorkflowStage
    context *WorkflowContext
    logger  interfaces.Logger
}

// WorkflowStage represents a single step in the workflow
type WorkflowStage interface {
    Execute(ctx *WorkflowContext) (*StageResult, error)
    CanSkip(ctx *WorkflowContext) bool
    GetName() string
    GetDescription() string
}

// WorkflowContext carries state between stages
type WorkflowContext struct {
    // User inputs
    AppName     string
    Flags       []string

    // Stage results
    SelectedArchetype *ArchetypeDefinition
    UserConfiguration *config.UserConfiguration
    TUIConfiguration  *TUIConfig

    // System context
    Dependencies *common.Dependencies
    VerbosityConfig *config.VerbosityConfig
    ChaosInjector   chaos.ChaosInjector
}

// StageResult contains the output of a workflow stage
type StageResult struct {
    StageType       StageType
    Data           map[string]interface{}
    NextStageHints []string
    Errors         []error
}
```

### Workflow Execution Flow
```go
func (wo *WorkflowOrchestrator) Execute() error {
    for _, stage := range wo.stages {
        // Check if stage should be skipped
        if stage.CanSkip(wo.context) {
            wo.logger.Debug("Skipping stage: %s", stage.GetName())
            continue
        }

        // Execute stage
        wo.logger.Info("Executing stage: %s", stage.GetName())
        result, err := stage.Execute(wo.context)
        if err != nil {
            return fmt.Errorf("stage %s failed: %w", stage.GetName(), err)
        }

        // Merge result into context
        err = wo.context.MergeResult(result)
        if err != nil {
            return fmt.Errorf("failed to merge stage result: %w", err)
        }
    }

    return nil
}
```

## Stage Implementations

### Stage 1: Archetype Selection
```go
type ArchetypeSelectionStage struct {
    registry *ArchetypeRegistry
    selector *ArchetypeSelector
}

func (as *ArchetypeSelectionStage) Execute(ctx *WorkflowContext) (*StageResult, error) {
    // Get available archetypes from registry
    archetypes := as.registry.GetAvailable()

    // Present selection UI using existing components
    selected, err := as.selector.ShowSelectionUI(archetypes)
    if err != nil {
        return nil, err
    }

    // Return result
    return &StageResult{
        StageType: StageTypeArchetypeSelection,
        Data: map[string]interface{}{
            "selectedArchetype": selected,
            "availableArchetypes": archetypes,
        },
    }, nil
}

func (as *ArchetypeSelectionStage) CanSkip(ctx *WorkflowContext) bool {
    // Skip if archetype already selected (direct mode)
    return ctx.SelectedArchetype != nil
}
```

### Stage 2: Contextual Prompting
```go
type ContextualPromptingStage struct {
    prompter *ContextualPrompter
}

func (cp *ContextualPromptingStage) Execute(ctx *WorkflowContext) (*StageResult, error) {
    // Get prompts for selected archetype
    prompts := cp.prompter.GetPromptsForArchetype(ctx.SelectedArchetype)

    // Execute prompts using enhanced inline prompter
    userConfig, err := cp.prompter.RunPrompts(prompts, ctx.Flags)
    if err != nil {
        return nil, err
    }

    return &StageResult{
        StageType: StageTypeContextualPrompting,
        Data: map[string]interface{}{
            "userConfiguration": userConfig,
            "promptsExecuted": len(prompts),
        },
    }, nil
}

func (cp *ContextualPromptingStage) CanSkip(ctx *WorkflowContext) bool {
    // Skip if no prompts required for archetype
    return len(cp.prompter.GetPromptsForArchetype(ctx.SelectedArchetype)) == 0
}
```

### Stage 3: TUI Execution
```go
type TUIExecutionStage struct {
    launcher *TUILauncher
}

func (te *TUIExecutionStage) Execute(ctx *WorkflowContext) (*StageResult, error) {
    // Generate TUI configuration based on archetype and user config
    tuiConfig := te.launcher.GenerateConfiguration(ctx.SelectedArchetype, ctx.UserConfiguration)

    // Launch TUI with dynamic configuration
    aarOutput, err := te.launcher.LaunchTUI(ctx.AppName, ctx.Flags, tuiConfig, ctx.VerbosityConfig, ctx.ChaosInjector)
    if err != nil {
        return nil, err
    }

    return &StageResult{
        StageType: StageTypeTUIExecution,
        Data: map[string]interface{}{
            "aarOutput": aarOutput,
            "tuiConfig": tuiConfig,
        },
    }, nil
}

func (te *TUIExecutionStage) CanSkip(ctx *WorkflowContext) bool {
    // TUI stage should never be skipped
    return false
}
```

## Archetype System Design

### Archetype Registry
```go
type ArchetypeRegistry struct {
    archetypes map[string]*ArchetypeDefinition
    mutex      sync.RWMutex
}

type ArchetypeDefinition struct {
    ID           string                 `json:"id"`
    Name         string                 `json:"name"`
    Description  string                 `json:"description"`
    IsDefault    bool                   `json:"isDefault"`
    Category     ArchetypeCategory      `json:"category"`
    PromptSet    string                 `json:"promptSet"`
    TUISteps     []StepDefinition       `json:"tuiSteps"`
    Dependencies []string               `json:"dependencies"`
    Features     []string               `json:"features"`
    Metadata     map[string]interface{} `json:"metadata"`
}

type StepDefinition struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Type        StepType          `json:"type"`
    Duration    time.Duration     `json:"duration"`
    Metadata    map[string]string `json:"metadata"`
}

// Registry methods
func (ar *ArchetypeRegistry) Register(archetype *ArchetypeDefinition) error
func (ar *ArchetypeRegistry) GetByID(id string) (*ArchetypeDefinition, error)
func (ar *ArchetypeRegistry) GetDefault() *ArchetypeDefinition
func (ar *ArchetypeRegistry) GetAvailable() []*ArchetypeDefinition
func (ar *ArchetypeRegistry) GetByCategory(category ArchetypeCategory) []*ArchetypeDefinition
```

### Archetype Definitions
```go
// Production React Application
var ProdWebArchetype = &ArchetypeDefinition{
    ID:          "prod-web",
    Name:        "Production React App",
    Description: "Full-featured React application with production optimizations",
    IsDefault:   true,
    Category:    CategoryWebApplication,
    PromptSet:   "prod-web-prompts",
    TUISteps: []StepDefinition{
        {ID: "init", Name: "Project Initialization", Description: "Setting up project structure", Type: StepTypeInitialization, Duration: 2 * time.Second},
        {ID: "deps", Name: "Installing Dependencies", Description: "npm install with production packages", Type: StepTypeDependencies, Duration: 5 * time.Second},
        {ID: "config", Name: "Configuration Setup", Description: "ESLint, Prettier, TypeScript config", Type: StepTypeConfiguration, Duration: 3 * time.Second},
        {ID: "build", Name: "Production Build Setup", Description: "Webpack optimization, CI/CD prep", Type: StepTypeBuild, Duration: 4 * time.Second},
        {ID: "deploy", Name: "Deployment Preparation", Description: "Docker, environment configuration", Type: StepTypeDeployment, Duration: 3 * time.Second},
        {ID: "finalize", Name: "Finalization", Description: "Final checks and documentation", Type: StepTypeFinalization, Duration: 2 * time.Second},
    },
    Features: []string{"typescript", "eslint", "prettier", "jest", "docker", "ci-cd"},
}

// Development React Application
var DevWebArchetype = &ArchetypeDefinition{
    ID:          "dev-web",
    Name:        "Development React App",
    Description: "Lightweight React setup for rapid development and prototyping",
    IsDefault:   false,
    Category:    CategoryWebApplication,
    PromptSet:   "dev-web-prompts",
    TUISteps: []StepDefinition{
        {ID: "init", Name: "Quick Project Setup", Description: "Minimal project structure", Type: StepTypeInitialization, Duration: 1 * time.Second},
        {ID: "deps", Name: "Essential Dependencies", Description: "Core React packages only", Type: StepTypeDependencies, Duration: 3 * time.Second},
        {ID: "config", Name: "Dev Configuration", Description: "Hot reload, basic linting", Type: StepTypeConfiguration, Duration: 2 * time.Second},
        {ID: "finalize", Name: "Ready to Code", Description: "Development server ready", Type: StepTypeFinalization, Duration: 1 * time.Second},
    },
    Features: []string{"typescript", "hot-reload", "basic-linting"},
}
```

## Prompt System Enhancement

### Contextual Prompter
```go
type ContextualPrompter struct {
    promptSets map[string][]config.PromptConfig
    validator  *PromptValidator
}

func (cp *ContextualPrompter) GetPromptsForArchetype(archetype *ArchetypeDefinition) []config.PromptConfig {
    if prompts, exists := cp.promptSets[archetype.PromptSet]; exists {
        return prompts
    }
    return cp.promptSets["default"]
}

// Archetype-specific prompt sets
var PromptSets = map[string][]config.PromptConfig{
    "prod-web-prompts": {
        {
            ID:       "production-data-access",
            Question: "Do you need access to production data sources?",
            UserOptions: map[string]string{
                "yes": "true",
                "no":  "false",
            },
            ConfigKey: "ProductionDataAccess",
            ResponseTemplates: map[string][]string{
                "yes": {"Enabling production data access", "TrustBridge and GRPC will be configured"},
                "no":  {"Skipping production data setup", "Development mode data sources only"},
            },
        },
        {
            ID:       "deployment-target",
            Question: "What's your primary deployment target?",
            UserOptions: map[string]string{
                "docker": "docker",
                "azure":  "azure",
                "local":  "local",
            },
            ConfigKey: "DeploymentTarget",
        },
    },
    "dev-web-prompts": {
        {
            ID:       "prototype-mode",
            Question: "Is this for rapid prototyping or learning?",
            UserOptions: map[string]string{
                "prototype": "prototype",
                "learning":  "learning",
                "demo":      "demo",
            },
            ConfigKey: "DevelopmentMode",
        },
        {
            ID:       "shared-access",
            Question: "Will this be shared with other developers?",
            UserOptions: map[string]string{
                "yes": "true",
                "no":  "false",
            },
            ConfigKey: "SharedAccess",
        },
    },
    "cli-prompts": {
        {
            ID:       "distribution-method",
            Question: "How will users install this CLI?",
            UserOptions: map[string]string{
                "npm":     "npm",
                "brew":    "brew",
                "binary":  "binary",
                "source":  "source",
            },
            ConfigKey: "DistributionMethod",
        },
        {
            ID:       "cross-platform",
            Question: "Do you need cross-platform binaries?",
            UserOptions: map[string]string{
                "yes": "true",
                "no":  "false",
            },
            ConfigKey: "CrossPlatform",
        },
    },
}
```

## TUI Integration Design

### Dynamic TUI Configuration
```go
type TUILauncher struct {
    stepGenerator *StepGenerator
    modelFactory  *ModelFactory
}

type TUIConfig struct {
    Steps           []StepDefinition
    Theme           *TUITheme
    ProgressStyle   ProgressStyle
    CompletionHooks []CompletionHook
    ErrorHandlers   []ErrorHandler
}

func (tl *TUILauncher) GenerateConfiguration(archetype *ArchetypeDefinition, userConfig *config.UserConfiguration) *TUIConfig {
    return &TUIConfig{
        Steps:         tl.stepGenerator.GenerateSteps(archetype, userConfig),
        Theme:         tl.getThemeForArchetype(archetype),
        ProgressStyle: tl.getProgressStyleForArchetype(archetype),
    }
}

func (tl *TUILauncher) LaunchTUI(appName string, flags []string, tuiConfig *TUIConfig, verbosityConfig *config.VerbosityConfig, chaosInjector chaos.ChaosInjector) (string, error) {
    // Create TUI model with dynamic configuration
    model := tl.modelFactory.CreateModel(appName, flags, tuiConfig, verbosityConfig, chaosInjector)

    // Configure tea program
    program := tea.NewProgram(
        model,
        tea.WithInput(os.Stdin),
        tea.WithOutput(os.Stderr),
    )

    // Run TUI
    finalModel, err := program.Run()
    if err != nil {
        return "", err
    }

    // Extract AAR output
    if appModel, ok := finalModel.(*models.AppModel); ok {
        return appModel.GetAAROutput(), nil
    }

    return "", nil
}
```

### Step Generation Logic
```go
type StepGenerator struct {
    stepTemplates map[StepType]*StepTemplate
}

func (sg *StepGenerator) GenerateSteps(archetype *ArchetypeDefinition, userConfig *config.UserConfiguration) []StepDefinition {
    var steps []StepDefinition

    for _, stepDef := range archetype.TUISteps {
        // Customize step based on user configuration
        customizedStep := sg.customizeStep(stepDef, userConfig)
        steps = append(steps, customizedStep)
    }

    return steps
}

func (sg *StepGenerator) customizeStep(step StepDefinition, userConfig *config.UserConfiguration) StepDefinition {
    customized := step

    // Modify step based on user configuration
    switch step.ID {
    case "deploy":
        if userConfig.DevOnly {
            customized.Name = "Skipping Deployment Setup"
            customized.Description = "Development mode - deployment skipped"
            customized.Duration = 0.5 * time.Second
        }
    case "deps":
        if userConfig.Template.Type == config.JavaScript {
            customized.Description = "Installing JavaScript dependencies"
        } else {
            customized.Description = "Installing TypeScript dependencies"
        }
    }

    return customized
}
```

## Integration Points

### Template System Bridge
```go
// Bridge existing template system to workflow system
func convertTemplateArchetypesToWorkflowOptions() []ArchetypeOption {
    templateArchetypes := []string{"prod-web", "dev-web", "cli", "service", "hackday", "engx-cmd"}

    var options []ArchetypeOption
    for _, id := range templateArchetypes {
        if archetype, err := ArchetypeRegistry.GetByID(id); err == nil {
            options = append(options, ArchetypeOption{
                ID:          archetype.ID,
                Name:        archetype.Name,
                Description: archetype.Description,
                IsDefault:   archetype.IsDefault,
            })
        }
    }

    return options
}
```

### Design Component Integration
```go
func (as *ArchetypeSelector) renderSelectionUI(archetypes []*ArchetypeDefinition) string {
    // Use existing design components
    header := components.NewHeader("Select Application Archetype").
        WithSubtitle("prod-web is recommended for most applications").
        WithRightText(fmt.Sprintf("%d options", len(archetypes)))

    var items []string
    for _, archetype := range archetypes {
        label := archetype.Name
        if archetype.IsDefault {
            label += " (Recommended)"
        }
        items = append(items, label)
    }

    list := components.NewList(items).
        WithBullet("→").
        AsNumbered()

    panel := components.NewPanel(list.Render()).
        WithTitle("Available Archetypes")

    return header.Render() + panel.Render()
}
```

## Testing Strategy

### Unit Testing Framework
```go
// Stage testing with mocks
func TestArchetypeSelectionStage(t *testing.T) {
    registry := &MockArchetypeRegistry{}
    selector := &MockArchetypeSelector{}
    stage := &ArchetypeSelectionStage{registry: registry, selector: selector}

    ctx := &WorkflowContext{}
    result, err := stage.Execute(ctx)

    assert.NoError(t, err)
    assert.Equal(t, StageTypeArchetypeSelection, result.StageType)
}

// Integration testing with real components
func TestWorkflowOrchestrator(t *testing.T) {
    orchestrator := NewWorkflowOrchestrator()
    orchestrator.AddStage(&ArchetypeSelectionStage{})
    orchestrator.AddStage(&ContextualPromptingStage{})
    orchestrator.AddStage(&TUIExecutionStage{})

    ctx := &WorkflowContext{AppName: "TestApp"}
    err := orchestrator.Execute()

    assert.NoError(t, err)
    assert.NotNil(t, ctx.SelectedArchetype)
}
```

---
**Status**: PLAN Phase Complete - Ready for HUM LEAD Implementation Approval
**Next Phase**: CODE - Begin Phase 1 foundation implementation
**Authority**: HUM LEAD approval required to proceed with implementation