# After-Action Report (AAR) System Architecture

## System Overview

The AAR system generates professional, build-tool-style summaries with actionable next steps after project creation completion.

### Core Components

```
StateComplete Trigger → AAR Generator → Output Formatter → Stream Writer
                             ↓
                      Data Collectors → Summary Builder → Next Steps Engine
```

## Core Components

### 1. AAR Data Model
**File**: `internal/aar/models.go`
```go
type AARSummary struct {
    ProjectInfo    ProjectInfo    `json:"project_info"`
    ExecutionInfo  ExecutionInfo  `json:"execution_info"`
    StepResults    []StepResult   `json:"step_results"`
    NextSteps      []NextStep     `json:"next_steps"`
    Troubleshooting *TroubleshootingInfo `json:"troubleshooting,omitempty"`
}

type ProjectInfo struct {
    Name         string            `json:"name"`
    Template     string            `json:"template"`
    DevOnly      bool              `json:"dev_only"`
    Features     map[string]bool   `json:"features"`
    Directory    string            `json:"directory"`
}

type ExecutionInfo struct {
    StartTime    time.Time         `json:"start_time"`
    EndTime      time.Time         `json:"end_time"`
    Duration     time.Duration     `json:"duration"`
    TotalSteps   int               `json:"total_steps"`
    SuccessSteps int               `json:"success_steps"`
    FailedSteps  int               `json:"failed_steps"`
    Performance  PerformanceMetrics `json:"performance"`
}

type StepResult struct {
    Name         string        `json:"name"`
    Status       StepStatus    `json:"status"`
    Duration     time.Duration `json:"duration"`
    StartTime    time.Time     `json:"start_time"`
    ErrorMessage string        `json:"error_message,omitempty"`
}

type NextStep struct {
    Action       string        `json:"action"`
    Description  string        `json:"description"`
    Command      string        `json:"command,omitempty"`
    Priority     StepPriority  `json:"priority"`
    Category     StepCategory  `json:"category"`
}
```

### 2. AAR Generator
**File**: `internal/aar/generator.go`
```go
type AARGenerator struct {
    tracker     *progresssim.Tracker
    config      *config.UserConfiguration
    verbosity   *VerbosityConfig
    startTime   time.Time
    projectPath string
}

func NewAARGenerator(tracker *progresssim.Tracker, config *config.UserConfiguration) *AARGenerator

func (g *AARGenerator) Generate() (*AARSummary, error) {
    summary := &AARSummary{
        ProjectInfo:   g.buildProjectInfo(),
        ExecutionInfo: g.buildExecutionInfo(),
        StepResults:   g.buildStepResults(),
    }

    summary.NextSteps = g.generateNextSteps(summary)
    summary.Troubleshooting = g.generateTroubleshooting(summary)

    return summary, nil
}
```

### 3. Next Steps Engine
**File**: `internal/aar/nextsteps.go`
```go
type NextStepsEngine struct {
    templates map[string][]NextStepTemplate
    rules     []NextStepRule
}

type NextStepTemplate struct {
    Condition   func(*AARSummary) bool
    Action      string
    Description string
    Command     string
    Priority    StepPriority
}

func (e *NextStepsEngine) Generate(summary *AARSummary) []NextStep {
    var steps []NextStep

    // Immediate actions
    steps = append(steps, e.getImmediateSteps(summary)...)

    // Development workflow
    steps = append(steps, e.getDevelopmentSteps(summary)...)

    // Configuration and setup
    steps = append(steps, e.getSetupSteps(summary)...)

    return e.prioritizeSteps(steps)
}
```

### 4. Output Formatters
**File**: `internal/aar/formatters.go`
```go
type OutputFormatter interface {
    Format(summary *AARSummary, verbosity VerbosityLevel) string
}

type StandardFormatter struct{}
type JSONFormatter struct{}
type CompactFormatter struct{}

func (f *StandardFormatter) Format(summary *AARSummary, verbosity VerbosityLevel) string {
    var output strings.Builder

    f.writeHeader(&output, summary)
    f.writeExecutionSummary(&output, summary, verbosity)
    f.writeNextSteps(&output, summary, verbosity)

    if verbosity >= VerbosityVerbose {
        f.writeDetailedResults(&output, summary)
    }

    if summary.Troubleshooting != nil {
        f.writeTroubleshooting(&output, summary.Troubleshooting)
    }

    return output.String()
}
```

## Integration Architecture

### 1. App Model Integration
**File**: `internal/tui/models/app.go` (modifications)
```go
type AppModel struct {
    // ... existing fields
    aarGenerator *aar.AARGenerator
}

// In Update method StateComplete case:
case StateComplete:
    if m.shouldShowAAR() {
        summary, err := m.aarGenerator.Generate()
        if err != nil {
            return m, tea.Quit
        }

        return m, tea.Sequence(
            m.displayAAR(summary),
            tea.Quit,
        )
    }
    return m, tea.Quit
```

### 2. Data Collection Points
```go
// During execution, collect AAR data
func (m *AppModel) collectAARData(step int, result StepResult) {
    if m.aarGenerator != nil {
        m.aarGenerator.RecordStep(step, result)
    }
}
```

### 3. Timing Integration
```go
// Leverage existing timing infrastructure
func (g *AARGenerator) buildExecutionInfo() ExecutionInfo {
    return ExecutionInfo{
        StartTime:    g.startTime,
        EndTime:      time.Now(),
        Duration:     time.Since(g.startTime),
        TotalSteps:   g.tracker.GetTotalSteps(),
        SuccessSteps: g.tracker.GetCompletedSteps(),
        FailedSteps:  g.tracker.GetFailedSteps(),
    }
}
```

## Next Steps Generation Logic

### Template-Based Generation
```go
var nextStepTemplates = map[string][]NextStepTemplate{
    "typescript": {
        {
            Condition: func(s *AARSummary) bool { return s.ExecutionInfo.FailedSteps == 0 },
            Action: "Start Development Server",
            Description: "Launch the development server to begin coding",
            Command: "cd " + s.ProjectInfo.Name + " && npm run dev",
            Priority: PriorityHigh,
        },
        {
            Action: "Open in Editor",
            Description: "Open project in your preferred code editor",
            Command: "code " + s.ProjectInfo.Name,
            Priority: PriorityHigh,
        },
    },
}
```

### Dynamic Rule Engine
```go
type NextStepRule struct {
    Name      string
    Condition func(*AARSummary) bool
    Generator func(*AARSummary) []NextStep
}

var dynamicRules = []NextStepRule{
    {
        Name: "Failed Steps Recovery",
        Condition: func(s *AARSummary) bool { return s.ExecutionInfo.FailedSteps > 0 },
        Generator: func(s *AARSummary) []NextStep {
            return generateRecoverySteps(s.StepResults)
        },
    },
}
```

## AAR Output Formats

### Standard Format (Normal Verbosity)
```
╭─────────────────────────────────────────────────────────────────╮
│                    ✨ Project Creation Complete                 │
├─────────────────────────────────────────────────────────────────┤
│ Project: MyApp (TypeScript)                                     │
│ Created: /Users/dev/projects/MyApp                              │
│ Duration: 2m 34s                                                │
│ Steps: 12/12 completed successfully                             │
╰─────────────────────────────────────────────────────────────────╯

🚀 Next Steps:
  1. cd MyApp && npm run dev     # Start development server
  2. code MyApp                  # Open in VS Code
  3. git init && git add .       # Initialize version control

📚 Quick Commands:
  npm run dev      Start development server
  npm run build    Build for production
  npm run test     Run test suite
  npm run lint     Check code quality

💡 Learn More:
  • React Documentation: https://react.dev
  • TypeScript Guide: https://typescriptlang.org/docs
  • Project README: ./MyApp/README.md
```

### Verbose Format (--verbose)
```
╭─────────────────────────────────────────────────────────────────╮
│                    ✨ Project Creation Complete                 │
├─────────────────────────────────────────────────────────────────┤
│ Project: MyApp (TypeScript Template)                           │
│ Features: Hot Reload, ESLint+Prettier, Testing                 │
│ Created: /Users/dev/projects/MyApp                              │
│ Duration: 2m 34s (average: 2m 45s)                             │
│ Steps: 12/12 completed successfully                             │
╰─────────────────────────────────────────────────────────────────╯

📊 Execution Summary:
  ✅ Project Structure        0.234s
  ✅ Package Installation     45.123s
  ✅ TypeScript Config        0.156s
  ✅ ESLint Configuration     0.234s
  ✅ Test Framework Setup     12.456s
  ✅ Development Server       8.234s
  ✅ Build Configuration      2.123s
  ✅ Hot Reload Setup         1.456s
  ✅ VS Code Settings         0.123s
  ✅ Git Configuration        0.345s
  ✅ README Generation        0.234s
  ✅ Final Validation         0.456s

🚀 Next Steps:
  [... detailed next steps as above ...]
```

### Quiet Format (--quiet)
```
✨ MyApp created successfully (2m 34s)

Next: cd MyApp && npm run dev
```

## Error Handling and Troubleshooting

### Failure Analysis
```go
func (g *AARGenerator) generateTroubleshooting(summary *AARSummary) *TroubleshootingInfo {
    if summary.ExecutionInfo.FailedSteps == 0 {
        return nil
    }

    info := &TroubleshootingInfo{
        FailedSteps: []FailedStepAnalysis{},
        Suggestions: []string{},
    }

    for _, step := range summary.StepResults {
        if step.Status == StepStatusFailed {
            analysis := g.analyzeFailure(step)
            info.FailedSteps = append(info.FailedSteps, analysis)
        }
    }

    return info
}
```

### Recovery Suggestions
```go
func (g *AARGenerator) analyzeFailure(step StepResult) FailedStepAnalysis {
    return FailedStepAnalysis{
        StepName:     step.Name,
        ErrorMessage: step.ErrorMessage,
        Suggestions:  g.getSuggestions(step),
        Recovery:     g.getRecoverySteps(step),
    }
}
```

## Performance Optimization

### Lazy Data Collection
- Step results collected during execution (minimal overhead)
- Next steps generated only when AAR displayed
- Troubleshooting analysis deferred to failure cases

### Caching Strategy
```go
type AARCache struct {
    templates map[string][]NextStepTemplate
    rules     []NextStepRule
    mutex     sync.RWMutex
}
```

### Memory Management
- Bounded data collection during execution
- Cleanup after AAR display
- Configurable detail levels to manage memory usage