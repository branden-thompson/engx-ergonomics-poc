# Implementation Strategy - Create Command Enhancements
**🛩️ BRTOPS v1.1.001** | Feature: create-cmd-enhancements | SEV-0

## Strategic Implementation Plan

### Recommended Approach: Modular Workflow Engine

**Decision Rationale**: Balance between architectural cleanliness, extensibility, and implementation complexity. Provides future-proof foundation while maintaining reasonable development timeline.

## Architecture Overview

### Core Components Architecture

```
plugins/create/plugin.go (Enhanced)
├── WorkflowOrchestrator (New)
│   ├── GuidedModeHandler (New)
│   ├── ArchetypeSelector (New)
│   ├── ContextualPrompter (Enhanced)
│   └── TUILauncher (Enhanced)
└── Legacy Direct Mode (Preserved)

internal/workflows/ (New Package)
├── archetype_selector.go
├── contextual_prompter.go
├── workflow_orchestrator.go
└── types.go

internal/archetypes/ (New Package)
├── registry.go
├── prod_web.go
├── dev_web.go
├── cli.go
├── service.go
├── hackday.go
└── engx_cmd.go
```

### Workflow Orchestration Pattern

```go
// Workflow stages with clean interfaces
type WorkflowStage interface {
    Execute(ctx *WorkflowContext) (*StageResult, error)
    CanSkip(ctx *WorkflowContext) bool
    GetName() string
}

// Orchestrator manages stage progression
type WorkflowOrchestrator struct {
    stages []WorkflowStage
    context *WorkflowContext
}

// Clean stage implementations
type ArchetypeSelectionStage struct{}
type ContextualPromptingStage struct{}
type TUIExecutionStage struct{}
```

## Implementation Phases

### Phase 1: Foundation & Infrastructure
**Target**: Core workflow system and archetype registry

**Deliverables:**
1. **WorkflowOrchestrator** - Stage management system
2. **ArchetypeRegistry** - Centralized archetype definitions
3. **WorkflowContext** - Shared state management
4. **Enhanced command structure** - Args handling for guided mode

**Files Created/Modified:**
- `internal/workflows/workflow_orchestrator.go`
- `internal/workflows/types.go`
- `internal/archetypes/registry.go`
- `plugins/create/plugin.go` (enhanced)

### Phase 2: Archetype Selection
**Target**: Interactive archetype selection using existing template system

**Deliverables:**
1. **ArchetypeSelector** - TUI-based archetype selection
2. **Template integration** - Leverage `pkg/common/templates_ui.go`
3. **Selection UI components** - Using `internal/tui/components/`
4. **Default highlighting** - prod-web as recommended default

**Files Created/Modified:**
- `internal/workflows/archetype_selector.go`
- Integration with existing `pkg/common/templates_ui.go`
- New selection components in `internal/tui/components/`

### Phase 3: Contextual Prompting
**Target**: Archetype-aware prompting system

**Deliverables:**
1. **ContextualPrompter** - Archetype-specific prompt routing
2. **Archetype prompt definitions** - 6 different prompt sets
3. **Enhanced prompt configs** - Extension of existing system
4. **Prompt validation** - Consistent with current patterns

**Files Created/Modified:**
- `internal/workflows/contextual_prompter.go`
- `internal/archetypes/prod_web.go` (and others)
- `internal/config/prompt_configs.go` (enhanced)

### Phase 4: TUI Integration
**Target**: Archetype-aware TUI with customized progress

**Deliverables:**
1. **TUILauncher** - Dynamic TUI configuration
2. **Archetype-specific steps** - Customized progress tables
3. **Design component integration** - Enhanced styling
4. **Dynamic step generation** - Based on archetype and config

**Files Created/Modified:**
- `internal/workflows/tui_launcher.go`
- `internal/tui/models/app.go` (enhanced)
- Enhanced step generation logic

## Technical Implementation Details

### Command Structure Enhancement

```go
// Enhanced Create command
cmd := &cobra.Command{
    Use:   "create [APP_NAME]",
    Args:  cobra.MaxArgs(1), // Changed from ExactArgs(1)
    RunE: func(cmd *cobra.Command, args []string) error {
        if len(args) == 0 {
            // NEW: Guided mode
            return p.executeGuidedMode(cmd, dependencies, flags...)
        } else {
            // EXISTING: Direct mode (preserved)
            return p.executeDirectMode(cmd, args, dependencies, flags...)
        }
    },
}
```

### Workflow Orchestrator Pattern

```go
type WorkflowOrchestrator struct {
    stages []WorkflowStage
    context *WorkflowContext
}

func (wo *WorkflowOrchestrator) Execute() error {
    for _, stage := range wo.stages {
        if stage.CanSkip(wo.context) {
            continue
        }

        result, err := stage.Execute(wo.context)
        if err != nil {
            return err
        }

        wo.context.Merge(result)
    }
    return nil
}
```

### Archetype Registry System

```go
type ArchetypeDefinition struct {
    ID           string
    Name         string
    Description  string
    IsDefault    bool
    PromptSet    string
    TUISteps     []StepDefinition
    Dependencies []string
}

type ArchetypeRegistry struct {
    archetypes map[string]*ArchetypeDefinition
}

func (ar *ArchetypeRegistry) GetDefault() *ArchetypeDefinition {
    // Returns prod-web as default
}

func (ar *ArchetypeRegistry) GetAvailable() []*ArchetypeDefinition {
    // Returns all registered archetypes
}
```

### Integration with Existing Systems

#### Template System Integration
```go
// Leverage existing template mappings
func (as *ArchetypeSelector) getArchetypeOptions() []ArchetypeOption {
    // Use existing mappings from pkg/common/templates_ui.go
    return convertTemplateArchetypesToSelectionOptions()
}
```

#### TUI Component Integration
```go
// Use enhanced design components
func (as *ArchetypeSelector) renderSelectionUI() string {
    header := components.NewHeader("Select Application Archetype").
        WithSubtitle("prod-web is recommended for most applications")

    list := components.NewList(archetypeNames).
        WithBullet("→").
        AsNumbered()

    return header.Render() + list.Render()
}
```

## Backward Compatibility Strategy

### Preservation Guarantees
1. **Existing Commands**: All `engx create <AppName>` usage remains identical
2. **Flag Support**: All current flags work exactly as before
3. **TUI Output**: Same progress table format and behavior
4. **AAR Generation**: Consistent AAR output format

### Migration Path
1. **Phase 1**: Infrastructure with no behavior changes
2. **Phase 2**: Guided mode addition (opt-in via no args)
3. **Phase 3**: Enhanced prompting (backwards compatible)
4. **Phase 4**: Enhanced TUI (same interface, better internals)

## Testing Strategy

### Unit Testing
- **Workflow orchestrator** stage execution
- **Archetype registry** selection logic
- **Contextual prompter** routing logic
- **TUI launcher** configuration generation

### Integration Testing
- **Guided mode flow** end-to-end
- **Direct mode preservation** regression testing
- **Template system integration** archetype mapping
- **TUI component integration** rendering tests

### User Acceptance Testing
- **Existing user workflows** regression validation
- **New guided mode flows** usability testing
- **Archetype selection UX** interaction testing
- **Cross-platform compatibility** verification

## Risk Mitigation

### High Priority Risks
1. **Breaking Existing Workflows**
   - Mitigation: Comprehensive regression testing
   - Fallback: Feature flag to disable guided mode

2. **Performance Degradation**
   - Mitigation: Lazy loading of archetype definitions
   - Monitoring: Startup time benchmarks

3. **Complex State Management**
   - Mitigation: Clean workflow context pattern
   - Validation: Isolated stage testing

### Implementation Risks
1. **Template System Coupling**
   - Mitigation: Abstract interface layer
   - Testing: Mock template system for tests

2. **TUI Integration Complexity**
   - Mitigation: Incremental integration approach
   - Fallback: Existing TUI system preserved

## Success Metrics

### Technical Metrics
- **Backward Compatibility**: 100% existing command preservation
- **Performance**: <10% startup time increase
- **Test Coverage**: >90% for new workflow components
- **Code Quality**: Clean architecture with <5 cyclomatic complexity

### User Experience Metrics
- **Discoverability**: Users can find archetype options naturally
- **Contextual Relevance**: Prompts match selected archetype
- **Smooth Transitions**: Seamless flow between workflow stages
- **Familiar Patterns**: Consistent with existing engx UX

---
**Next Steps**: Create detailed architecture decision records and begin Phase 1 implementation
**Authority**: HUM LEAD approval required for implementation approach