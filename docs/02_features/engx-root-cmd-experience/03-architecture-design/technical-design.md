# Technical Design - EngX Root Command Experience

## Architecture Overview

### Component Architecture
```
cmd/engx/main.go
├── RootCommandRouter (NEW)
│   ├── DetectNoArgs() → launches InteractiveInterface
│   └── PreserveHelpFlag() → traditional help via --help
│
internal/tui/models/
├── root_command_interface.go (NEW)
│   ├── InteractiveCommandModel (bubbletea.Model)
│   ├── CommandTableRenderer
│   └── KeyboardNavigationHandler
│
internal/config/
├── roadmap_config.go (NEW)
│   ├── RoadmapLoader
│   ├── CommandStatusResolver
│   └── CrewConfigManager
│
pkg/common/
├── command_discovery.go (NEW)
│   ├── DynamicCommandDiscovery
│   ├── CommandFilteringSystem
│   └── PluginCommandIntegration
```

## Data Models

### RoadmapConfig Structure
```yaml
metadata:
  version: "1.0"
  cli_version: "v0.8.0"

user_context:
  default_crew: "CREW-1234"
  ldap_source: "system"

command_categories:
  - name: "Create & Manage Applications"
    commands:
      - cmd: "create"
        description: "Create new App from Scratch"
        owner_crew: "CREW-1234"
        status: "available"
```

### CommandEntry Go Structure
```go
type CommandEntry struct {
    Number      int                    `json:"number"`
    Command     string                 `json:"command"`
    UseCase     string                 `json:"use_case"`
    OwnerCrew   string                 `json:"owner_crew"`
    Status      config.CommandStatus   `json:"status"`
    Category    string                 `json:"category"`
    Available   bool                   `json:"available"`
    Selectable  bool                   `json:"selectable"`
}
```

## Responsive Layout Design

### Terminal Width Breakpoints
- **< 40 cols**: Emergency mode (command names only)
- **40-60 cols**: Compact table mode
- **60-80 cols**: Standard table mode
- **80+ cols**: Full table with expanded descriptions

### Column Layout Strategy
```
Standard Layout (80+ cols):
┌────┬────────────────────┬──────────────────────────────────────┬──────────────┐
│NUM │ COMMAND            │ USE TO                               │ OWNER CREW ID│
├────┼────────────────────┼──────────────────────────────────────┼──────────────┤
│ 1. │ `engx create`      │ Create new App from Scratch          │ CREW-1234    │
│ 2. │ `engx templates`   │ Discover Supported App Types, and    │ CREW-1234    │
│    │                    │ this is an example of a long use-to  │              │
│    │                    │ description that wraps within it's   │              │
│    │                    │ data-cell on the CLI                 │              │
└────┴────────────────────┴──────────────────────────────────────┴──────────────┘

Compact Layout (40-60 cols):
┌────┬──────────────┬──────────────────┬──────────┐
│NUM │ COMMAND      │ USE TO           │ CREW     │
├────┼──────────────┼──────────────────┼──────────┤
│ 1. │ create       │ Create new App   │ CREW-1234│
│ 2. │ templates    │ Discover Apps    │ CREW-1234│
└────┴──────────────┴──────────────────┴──────────┘
```

## Keyboard Interaction Design

### Navigation Patterns
```go
// Key Mappings
"↑"/"k"     → Previous item
"↓"/"j"     → Next item
"1"-"9"     → Select by number (if <= 9 items)
"0"         → Select item 10 (if exists)
"Enter"     → Execute selected command
"h"/"?"     → Toggle help
"q"/"Ctrl+C" → Quit
```

### Selection State Management
- Visual highlight with lipgloss styling
- Skip non-selectable items during navigation
- Number-based quick selection with validation
- Error feedback for invalid selections

## Command Execution Flow

### Root Command Detection
```go
func (r *RootCommandRouter) ShouldShowInteractiveInterface(args []string) bool {
    // Show interactive interface only when:
    // 1. No arguments provided (just 'engx')
    // 2. Not requesting help (--help, -h)
    // 3. Not requesting version (--version)
    return len(args) == 0 && !hasSpecialFlags(args)
}
```

### Interface Preservation Strategy
```go
// After selection, display final state before execution
func (m *InteractiveCommandModel) ExecuteCommand(cmd CommandEntry) tea.Cmd {
    // 1. Command selected and interface state preserved
    // 2. Execute via tea.ExecProcess for seamless transition
    return tea.ExecProcess(
        exec.Command("engx", formatCommand(cmd)...),
        func(err error) tea.Msg { return tea.Quit() },
    )
}
```

## Integration Points

### Cobra Command System
- Intercept root command execution before cobra processing
- Preserve all existing global flags and help functionality
- Access registered cobra commands for discovery validation
- Maintain plugin system integration

### Existing TUI Infrastructure
- Reuse lipgloss styles from `internal/tui/styles/`
- Follow bubbletea patterns from `internal/tui/models/`
- Integrate with verbosity system (--quiet, --verbose, etc.)
- Use established component architecture

### Configuration System
- YAML-based roadmap configuration in `.engx/roadmap.yaml`
- Environment variable integration for user context
- Dynamic command status resolution
- Extensible categorization system

## Error Handling & Fallbacks

### Configuration Loading Failures
```go
// If roadmap config fails to load, fall back to help
if err := loadRoadmapConfig(); err != nil {
    fmt.Fprintf(os.Stderr, "Warning: Could not load roadmap configuration: %v\n", err)
    return r.rootCmd.Help()
}
```

### Terminal Compatibility
- TTY detection for interactive interface eligibility
- Graceful fallback to help screen for non-terminal environments
- Alt-screen buffer usage to preserve terminal state

### Command Discovery Failures
- Validate cobra command registry against roadmap
- Mark commands as "coming soon" if not implemented
- Filter excluded commands (dev tools, internal commands)

## Security Considerations

### Input Validation
- Sanitize roadmap configuration inputs
- Validate command names against whitelist
- Prevent command injection through configuration

### Environment Access
- Safe environment variable access for user context
- No sensitive data exposure in debug output
- Proper file permission handling for configuration files

## Performance Characteristics

### Startup Performance
- Lazy loading of roadmap configuration
- Minimal overhead for normal command execution
- Fast command discovery with caching

### Interactive Performance
- Smooth 60fps rendering with minimal state updates
- Efficient text wrapping and layout calculations
- Responsive keyboard handling without input lag

### Memory Usage
- Lightweight command structures
- Minimal bubbletea model state
- Efficient string operations for display rendering