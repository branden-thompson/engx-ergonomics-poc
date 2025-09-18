# Verbosity System Architecture

## System Overview

### Verbosity Level Hierarchy
```
DEBUG (--debug)     │ Maximum verbosity, all system outputs, technical details
  ↓                 │
VERBOSE (--verbose) │ Enhanced details, progress bars for multi-step processes
  ↓                 │
DEFAULT (default)   │ Current view (chosen when no verbosity option specified)
  ↓                 │
CONCISE (--concise) │ Less detail, components info and granular items hidden
  ↓                 │
QUIET (--quiet)     │ Essential info only: total progress, current step, footer
```

## Core Components

### 1. Verbosity Configuration
**File**: `internal/config/verbosity.go`
```go
type VerbosityLevel int

const (
    VerbosityQuiet VerbosityLevel = iota
    VerbosityConcise
    VerbosityDefault
    VerbosityVerbose
    VerbosityDebug
)

type VerbosityConfig struct {
    Level           VerbosityLevel
    ShowTimestamps  bool
    ShowStepDetails bool
    ShowDebugInfo   bool
    ShowAAR         bool
}
```

### 2. Output Controller
**File**: `internal/output/controller.go`
```go
type OutputController struct {
    config    *VerbosityConfig
    writer    io.Writer
    startTime time.Time
}

func (oc *OutputController) ShouldShow(level VerbosityLevel) bool
func (oc *OutputController) WriteProgress(step string, details string)
func (oc *OutputController) WriteDebug(message string)
func (oc *OutputController) WriteAAR(summary AARSummary)
```

### 3. Enhanced Renderer Integration
**File**: `internal/tui/components/enhanced_renderer.go` (modifications)
```go
type EnhancedRenderer struct {
    // ... existing fields
    verbosity *VerbosityConfig
    output    *OutputController
}

func (er *EnhancedRenderer) RenderWithVerbosity() string
func (er *EnhancedRenderer) ShouldShowDetail(component string) bool
```

## Flag System Architecture

### CLI Flag Integration
**File**: `cmd/engx/main.go` (modifications)
```go
// Enhanced flag system
rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Essential info only")
rootCmd.PersistentFlags().Bool("concise", false, "Less detail, hide components info")
rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enhanced details with progress bars")
rootCmd.PersistentFlags().BoolP("debug", "d", false, "Maximum verbosity, all system outputs")
```

### Flag Processing
**File**: `internal/commands/create.go` (modifications)
```go
func determineVerbosityLevel(cmd *cobra.Command) VerbosityLevel {
    // Highest precedence wins (debug > verbose > default > concise > quiet)
    if debug, _ := cmd.Flags().GetBool("debug"); debug {
        return VerbosityDebug
    }
    if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
        return VerbosityVerbose
    }
    if concise, _ := cmd.Flags().GetBool("concise"); concise {
        return VerbosityConcise
    }
    if quiet, _ := cmd.Flags().GetBool("quiet"); quiet {
        return VerbosityQuiet
    }
    return VerbosityDefault
}
```

## Verbosity Level Specifications

### Quiet Mode (--quiet)
**Behavior**: Essential info only: total progress, current step, footer
- **Progress**: Single line progress indicator only
- **Steps**: Current major step name only
- **Components**: Hidden completely
- **Timing**: Basic footer timing only
- **AAR**: One-line completion status
- **Errors**: Error messages only

### Concise Mode (--concise)
**Behavior**: Less detail, components info and granular items hidden
- **Progress**: Main progress sections only
- **Steps**: Major steps without sub-details
- **Components**: Component info hidden
- **Timing**: Essential timing information
- **AAR**: Brief summary with essential next steps
- **Errors**: Standard error handling

### Default Mode (no flag)
**Behavior**: Current user experience maintained exactly
- **Progress**: Current progress bar and percentages
- **Steps**: Current step name display
- **Components**: Current component status display
- **Timing**: Current timing display
- **AAR**: Standard summary with next steps
- **Errors**: Current error handling

### Verbose Mode (--verbose)
**Behavior**: Enhanced details, progress bars for multi-step processes
- **Progress**: Enhanced progress bars for all multi-step processes
- **Steps**: Individual step timing and detailed status
- **Components**: Detailed component status with descriptions
- **Timing**: Individual step timing displayed
- **AAR**: Comprehensive summary with detailed metrics
- **Errors**: Enhanced error context

### Debug Mode (--debug)
**Behavior**: Maximum verbosity, all system outputs, technical details
- **Progress**: All internal state changes and transitions
- **Steps**: Technical implementation details and diagnostics
- **Components**: Component lifecycle and technical information
- **Timing**: Microsecond precision timing with performance metrics
- **AAR**: Full technical report with comprehensive diagnostics
- **Errors**: Stack traces and complete debug context

## Integration Points

### 1. App Model Integration
```go
type AppModel struct {
    // ... existing fields
    verbosity *VerbosityConfig
    output    *OutputController
}

func (m *AppModel) UpdateWithVerbosity(msg tea.Msg) (tea.Model, tea.Cmd)
```

### 2. Progress Tracker Integration
```go
type Tracker struct {
    // ... existing fields
    verbosity *VerbosityConfig
}

func (t *Tracker) LogWithVerbosity(level VerbosityLevel, message string)
```

### 3. State Transition Hooks
```go
// In app.go Update method
case StateComplete:
    if m.verbosity.ShowAAR {
        return m, m.generateAAR()
    }
```

## Output Stream Management

### Stream Separation
- **stdout**: AAR and final output
- **stderr**: Progress and status (current TUI behavior)
- **Debug**: Separate debug stream option

### Buffer Management
```go
type BufferedOutput struct {
    stdout bytes.Buffer
    stderr bytes.Buffer
    debug  bytes.Buffer
}

func (bo *BufferedOutput) FlushAppropriate(level VerbosityLevel)
```

## Performance Considerations

### Lazy Evaluation
- Debug information only computed when debug mode active
- AAR data collection happens during execution, generation on demand
- String formatting deferred until output time

### Memory Management
```go
type VerbosityBuffer struct {
    messages []VerbosityMessage
    maxSize  int
    level    VerbosityLevel
}

func (vb *VerbosityBuffer) AddMessage(level VerbosityLevel, msg string)
func (vb *VerbosityBuffer) Flush() []string
```

## Configuration Management

### Environment Variables
```bash
ENGX_VERBOSITY=normal|quiet|verbose|debug
ENGX_SHOW_TIMESTAMPS=true|false
ENGX_AAR_ENABLED=true|false
```

### Config File Support
```yaml
verbosity:
  level: normal
  show_timestamps: false
  show_aar: true
  debug_to_file: false
```

## Error Handling

### Verbosity-Aware Error Display
```go
func (oc *OutputController) WriteError(err error, context string) {
    switch oc.config.Level {
    case VerbosityQuiet:
        fmt.Fprintf(oc.writer, "Error: %s\n", err.Error())
    case VerbosityNormal:
        // Current behavior
    case VerbosityVerbose:
        // Enhanced context
    case VerbosityDebug:
        // Stack trace and debug info
    }
}
```

## Testing Strategy

### Unit Testing
- Test each verbosity level output
- Validate flag parsing logic
- Buffer management testing

### Integration Testing
- End-to-end verbosity behavior
- Performance impact measurement
- Cross-platform output verification