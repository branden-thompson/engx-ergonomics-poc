# Requirements & Context Collection - Create Command Enhancements
**🛩️ BRTOPS v1.1.001** | Feature: create-cmd-enhancements | SEV-0

## Mission Statement
Transform `engx create` from a direct-to-TUI command into a guided archetype-driven workflow that provides contextual selection and customized experiences based on application type.

## Context Analysis

### Current Architecture Assessment
**Analyzed**: `plugins/create/plugin.go` and `internal/prompts/inline.go`

**Current Flow:**
```
engx create <AppName> → InlinePrompter (3 prompts) → TUI Progress Table
```

**Current State Characteristics:**
- **Hard-coded React focus**: Command explicitly states "Create a new React application"
- **Fixed prompting**: 3 prompts focused on prod-web archetype (ProductionDataAccess, DeploymentTarget)
- **Single TUI flow**: Same progress table regardless of archetype selected
- **Args requirement**: `cobra.ExactArgs(1)` - requires app name
- **Template awareness**: Limited to typescript/javascript templates

### Target State Vision

**Enhanced Flow:**
```
engx create → Archetype Selection → Contextual Prompts → Customized TUI Progress
```

**New Flow Components:**
1. **Guided Mode** (`engx create` with no args)
2. **Archetype Selection** (prod-web, dev-web, cli, service, etc.)
3. **Contextual Prompting** (different questions per archetype)
4. **Customized TUI** (archetype-specific steps and details)

## Requirements Specification

### R1: Guided Mode Entry Point
- **R1.1**: `engx create` (no arguments) enters guided mode
- **R1.2**: `engx create <AppName>` preserves legacy direct mode
- **R1.3**: Cobra command args change from `ExactArgs(1)` to `MaxArgs(1)`

### R2: Archetype Selection System
- **R2.1**: Present available archetypes using template system integration
- **R2.2**: Display archetype descriptions with "prod-web is default" guidance
- **R2.3**: Support keyboard navigation for archetype selection
- **R2.4**: Integrate with existing `pkg/common/templates_ui.go` archetype mappings

### R3: Contextual Prompting Framework
- **R3.1**: Archetype-aware prompt configuration system
- **R3.2**: Different prompt sets per archetype:
  - **prod-web**: Current 3 prompts (ProductionDataAccess, DeploymentTarget, etc.)
  - **dev-web**: Development-focused prompts (no production deployment)
  - **cli**: CLI-specific prompts (packaging, distribution)
  - **service**: API/service prompts (protocols, monitoring)
- **R3.3**: Maintain existing prompt validation and styling

### R4: Customized TUI Integration
- **R4.1**: Archetype-specific progress steps in TUI table
- **R4.2**: Use enhanced design components from `internal/tui/components/`
- **R4.3**: Dynamic step descriptions based on selected archetype
- **R4.4**: Preserve existing chaos marine and verbosity support

### R5: Modular Architecture Requirements
- **R5.1**: Composable workflow stages (selection → prompting → TUI)
- **R5.2**: Clean separation between archetype logic and execution
- **R5.3**: Maintain backward compatibility with existing flags
- **R5.4**: Plugin architecture compliance

## Archetype Mapping Specification

### Available Archetypes (from templates system)
1. **prod-web** - Production React applications (default)
2. **dev-web** - Development/prototyping React apps
3. **cli** - Command-line tools and utilities
4. **service** - Backend services and APIs
5. **hackday** - Rapid prototyping projects
6. **engx-cmd** - EngX plugin development

### Archetype-Specific Prompt Sets

#### prod-web (Current)
- ProductionDataAccess: "Do you need access to production data sources?"
- DeploymentTarget: "What's your primary deployment target?"
- Additional: TBD based on analysis

#### dev-web
- DevelopmentMode: "Is this for rapid prototyping or learning?"
- SharedAccess: "Will this be shared with other developers?"
- TestingScope: "What level of testing do you need?"

#### cli
- DistributionMethod: "How will users install this CLI?"
- PackageManager: "Which package manager support?"
- CrossPlatform: "Do you need cross-platform binaries?"

#### service
- ProtocolType: "What communication protocol?" (REST, gRPC, GraphQL)
- DataPersistence: "Do you need database integration?"
- MonitoringLevel: "What monitoring/observability level?"

## Technical Constraints

### Backward Compatibility
- **C1**: Existing `engx create <AppName>` usage must continue working
- **C2**: All current flags must remain functional
- **C3**: TUI output format should remain consistent
- **C4**: AAR generation must work with new flow

### Integration Requirements
- **C5**: Must use existing `pkg/common/templates_ui.go` archetype mappings
- **C6**: Must leverage `internal/tui/components/` design system
- **C7**: Must integrate with current chaos marine and verbosity systems
- **C8**: Must support all existing flags (--dev-only, --template, etc.)

### Performance Requirements
- **C9**: Guided mode should not add significant startup delay
- **C10**: Archetype selection should be responsive
- **C11**: TUI performance must remain unchanged

## Success Criteria

### User Experience Goals
1. **Intuitive Discovery**: Engineers can discover archetype options naturally
2. **Contextual Guidance**: Prompts are relevant to selected archetype
3. **Smooth Transition**: Seamless flow from selection to TUI
4. **Backward Compatible**: Existing users experience no disruption

### Technical Achievement Goals
1. **Modular Design**: Clean separation of concerns between workflow stages
2. **Extensible Framework**: Easy to add new archetypes and prompts
3. **Component Integration**: Proper use of enhanced design components
4. **Maintainable Code**: Clear architecture for future enhancements

## Risk Assessment

### High Priority Risks
- **R-H1**: Breaking existing user workflows during transition
- **R-H2**: Complexity explosion in prompt configuration system
- **R-H3**: TUI integration challenges with dynamic content

### Medium Priority Risks
- **R-M1**: Performance degradation in guided mode
- **R-M2**: Inconsistent archetype behavior across different flows
- **R-M3**: Template system integration coupling issues

### Mitigation Strategies
- **M1**: Extensive testing of legacy command patterns
- **M2**: Phased implementation with feature flags
- **M3**: Comprehensive documentation of new workflows

## Documentation Requirements

### SEV-0 Documentation Standards
- **D1**: Complete 7-folder structure documentation
- **D2**: Architecture decision records for major design choices
- **D3**: User workflow documentation for all archetype paths
- **D4**: Technical implementation details for maintainers
- **D5**: Migration guide for existing users

---
**Status**: RCC Phase Complete - Ready for Planning Phase
**Next Phase**: Strategic planning and architecture design
**Authority**: HUM LEAD collaboration required for planning approval