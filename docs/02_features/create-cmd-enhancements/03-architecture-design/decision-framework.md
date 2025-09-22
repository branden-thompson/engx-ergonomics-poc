# Architecture Decision Framework - Create Command Enhancements
**🛩️ BRTOPS v1.1.001** | Feature: create-cmd-enhancements | SEV-0

## Decision Summary

**Selected Approach**: Modular Workflow Engine (Approach 2)
**Authority**: Pending HUM LEAD approval
**Confidence**: High (8/10)

## Decision Matrix Analysis

| Criteria | Approach 1: Minimal | Approach 2: Modular | Approach 3: Plugin |
|----------|-------------------|--------------------|--------------------|
| **Implementation Speed** | ⭐⭐⭐⭐⭐ Fast | ⭐⭐⭐ Medium | ⭐ Slow |
| **Extensibility** | ⭐⭐ Limited | ⭐⭐⭐⭐⭐ Excellent | ⭐⭐⭐⭐⭐ Ultimate |
| **Maintainability** | ⭐⭐ Poor | ⭐⭐⭐⭐ Good | ⭐⭐⭐ Complex |
| **Risk Level** | ⭐⭐⭐⭐⭐ Low | ⭐⭐⭐ Medium | ⭐ High |
| **Future-Proof** | ⭐⭐ Limited | ⭐⭐⭐⭐ Strong | ⭐⭐⭐⭐⭐ Ultimate |
| **Testing Ease** | ⭐⭐⭐ Medium | ⭐⭐⭐⭐ Good | ⭐⭐ Complex |
| **Total Score** | **17/30** | **23/30** | **18/30** |

## Key Architecture Decisions

### ADR-001: Workflow Orchestration Pattern
**Decision**: Implement stage-based workflow orchestrator
**Rationale**:
- Clean separation of concerns (selection → prompting → TUI)
- Easy to test individual stages
- Extensible for future workflow modifications
- Familiar pattern from existing codebase

**Alternatives Considered**:
- Monolithic flow control (rejected: poor separation)
- Event-driven architecture (rejected: over-complex)

**Implementation**:
```go
type WorkflowStage interface {
    Execute(ctx *WorkflowContext) (*StageResult, error)
    CanSkip(ctx *WorkflowContext) bool
    GetName() string
}
```

### ADR-002: Archetype Registry Design
**Decision**: Centralized registry with static definitions
**Rationale**:
- Integration with existing template system
- Compile-time type safety
- Performance (no runtime loading)
- Simpler than configuration files

**Alternatives Considered**:
- JSON/YAML configuration (rejected: premature flexibility)
- Embedded archetype plugins (rejected: complexity)

**Implementation**:
```go
type ArchetypeRegistry struct {
    archetypes map[string]*ArchetypeDefinition
}
```

### ADR-003: Backward Compatibility Strategy
**Decision**: Preserve existing direct mode completely
**Rationale**:
- Zero risk to existing users
- Clean separation between old and new flows
- Easy rollback if issues discovered

**Implementation**:
```go
func (p *Plugin) Create(deps interface{}) *cobra.Command {
    cmd := &cobra.Command{
        Args: cobra.MaxArgs(1), // Changed from ExactArgs(1)
        RunE: func(cmd *cobra.Command, args []string) error {
            if len(args) == 0 {
                return p.executeGuidedMode(...)  // NEW
            } else {
                return p.executeDirectMode(...)  // EXISTING
            }
        },
    }
}
```

### ADR-004: TUI Integration Approach
**Decision**: Enhance existing TUI model with dynamic configuration
**Rationale**:
- Leverages proven TUI system
- Maintains consistent user experience
- Uses existing design components
- Minimal disruption to current flows

**Alternatives Considered**:
- Separate TUI for guided mode (rejected: code duplication)
- Complete TUI rewrite (rejected: high risk, unnecessary)

### ADR-005: Template System Integration
**Decision**: Abstract layer over existing template mappings
**Rationale**:
- Reuses existing archetype definitions
- Maintains consistency with `engx templates` command
- Avoids duplication of archetype information

**Implementation**:
```go
func convertTemplateArchetypesToWorkflowOptions() []ArchetypeOption {
    // Bridge existing template system to workflow system
}
```

## Component Dependencies

### Core Dependency Graph
```
WorkflowOrchestrator
├── ArchetypeSelector
│   └── pkg/common/templates_ui.go (existing)
├── ContextualPrompter
│   └── internal/prompts/inline.go (enhanced)
└── TUILauncher
    └── internal/tui/models/app.go (enhanced)
```

### New Package Structure
```
internal/workflows/
├── workflow_orchestrator.go    # Core orchestration logic
├── archetype_selector.go       # TUI-based archetype selection
├── contextual_prompter.go      # Archetype-aware prompting
├── tui_launcher.go            # Dynamic TUI configuration
└── types.go                   # Shared types and interfaces

internal/archetypes/
├── registry.go                # Centralized archetype definitions
├── prod_web.go               # Production React app archetype
├── dev_web.go                # Development React app archetype
├── cli.go                    # CLI tool archetype
├── service.go                # Backend service archetype
├── hackday.go                # Rapid prototyping archetype
└── engx_cmd.go               # EngX plugin archetype
```

## Integration Points

### With Existing Systems

#### Template System Integration
- **File**: `pkg/common/templates_ui.go`
- **Integration**: Read archetype mappings, convert to workflow options
- **Risk**: Low (read-only access)

#### Prompt System Enhancement
- **File**: `internal/prompts/inline.go`
- **Integration**: Extend with archetype-aware routing
- **Risk**: Medium (requires careful extension)

#### TUI Model Enhancement
- **File**: `internal/tui/models/app.go`
- **Integration**: Dynamic step configuration based on archetype
- **Risk**: Medium (core TUI logic changes)

#### Design Component Usage
- **Files**: `internal/tui/components/*.go`
- **Integration**: Use existing components for archetype selection UI
- **Risk**: Low (consumer relationship)

## Implementation Phases Detail

### Phase 1: Foundation (Week 1)
**Objective**: Core infrastructure without behavior changes

**Tasks**:
1. Create `internal/workflows/` package structure
2. Implement `WorkflowOrchestrator` base pattern
3. Create `WorkflowContext` and `StageResult` types
4. Enhance command structure with optional args
5. Add workflow routing logic (guided vs direct mode)

**Validation**: All existing functionality works identically

### Phase 2: Archetype Selection (Week 2)
**Objective**: Interactive archetype selection

**Tasks**:
1. Implement `ArchetypeSelector` stage
2. Create archetype registry with template integration
3. Build selection UI using existing components
4. Add archetype selection to guided workflow
5. Implement default highlighting (prod-web)

**Validation**: Guided mode shows archetype selection, direct mode unchanged

### Phase 3: Contextual Prompting (Week 3)
**Objective**: Archetype-aware prompting

**Tasks**:
1. Implement `ContextualPrompter` stage
2. Define archetype-specific prompt sets
3. Enhance prompt configuration system
4. Add prompt routing logic
5. Integrate with existing prompt validation

**Validation**: Different archetypes show different prompts

### Phase 4: TUI Integration (Week 4)
**Objective**: Customized TUI based on archetype

**Tasks**:
1. Implement `TUILauncher` stage
2. Add dynamic step generation logic
3. Enhance TUI model with archetype awareness
4. Integrate archetype-specific progress steps
5. Preserve all existing TUI features

**Validation**: TUI shows archetype-specific steps and descriptions

## Risk Assessment & Mitigation

### Technical Risks

#### High Impact Risks
1. **Breaking Existing Workflows**
   - Probability: Low
   - Impact: High
   - Mitigation: Comprehensive regression testing, feature flags

2. **TUI Integration Complexity**
   - Probability: Medium
   - Impact: High
   - Mitigation: Incremental integration, fallback preservation

#### Medium Impact Risks
1. **Template System Coupling**
   - Probability: Medium
   - Impact: Medium
   - Mitigation: Abstract interface layer, dependency injection

2. **Performance Degradation**
   - Probability: Low
   - Impact: Medium
   - Mitigation: Lazy loading, startup benchmarks

### Implementation Risks

1. **State Management Complexity**
   - Mitigation: Clean WorkflowContext pattern, immutable data structures

2. **Testing Challenge**
   - Mitigation: Stage isolation, comprehensive mocking

3. **Code Maintenance Burden**
   - Mitigation: Clear documentation, simple interfaces

## Success Criteria

### Technical Objectives
- ✅ **Backward Compatibility**: 100% preservation of existing functionality
- ✅ **Performance**: <10% startup time increase
- ✅ **Test Coverage**: >90% for new workflow components
- ✅ **Architecture Quality**: Clean separation of concerns, <5 cyclomatic complexity

### User Experience Objectives
- ✅ **Discoverability**: Users can naturally find archetype options
- ✅ **Contextual Relevance**: Prompts match selected archetype appropriately
- ✅ **Smooth Flow**: Seamless transitions between workflow stages
- ✅ **Consistency**: Familiar patterns consistent with existing engx UX

### Business Objectives
- ✅ **Extensibility**: Easy addition of new archetypes in future
- ✅ **Maintainability**: Clear codebase organization for long-term maintenance
- ✅ **User Adoption**: Existing users continue without disruption
- ✅ **Future Readiness**: Foundation for advanced workflow features

---
**Next Phase**: Begin implementation with Phase 1 foundation work
**Decision Authority**: Awaiting HUM LEAD final approval for architecture approach