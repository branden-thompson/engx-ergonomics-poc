# Technical Specifications - Chaos Marine Injection Engine

## SPECIFICATION CLASSIFICATION
- **Type**: MAJOR SEV-0 SYSTEM FEATURE
- **Technical Level**: COMPREHENSIVE implementation specifications
- **Integration Depth**: SYSTEM-WIDE with existing infrastructure leverage
- **Safety Requirements**: MAXIMUM with complete isolation guarantees

---

## CORE TECHNICAL SPECIFICATIONS

### 1. FILE STRUCTURE AND ORGANIZATION

```
internal/chaos/
├── config.go              # Core configuration structures and validation
├── injector.go             # Main injection engine implementation
├── safety.go               # Safety boundaries and validation
├── behavior.go             # User behavior tracking and analysis
├── scenarios/
│   ├── enhanced.go         # Enhanced scenario structures
│   ├── educational.go      # Educational features and assessment
│   └── chaining.go         # Failure cascade management
├── analytics/
│   ├── metrics.go          # Performance and usage metrics
│   ├── telemetry.go        # Data collection framework
│   └── reporting.go        # Report generation
└── testing/
    ├── deterministic.go    # Deterministic testing framework
    ├── safety_test.go      # Safety boundary validation
    └── performance_test.go # Performance regression testing
```

### 2. CORE INTERFACE SPECIFICATIONS

#### 2.1 Primary Chaos Injector Interface
```go
package chaos

import (
    "context"
    "time"
    "github.com/bthompso/engx-ergonomics-poc/internal/simulation/errors"
)

// ChaosInjector defines the primary interface for chaos injection
type ChaosInjector interface {
    // Core lifecycle
    Initialize(config *ChaosConfig) error
    Shutdown(ctx context.Context) error
    IsEnabled() bool

    // Injection decisions
    ShouldInject(operation string) (bool, error)
    SelectScenario(operation string, context *OperationContext) (*ChaosScenario, error)
    CalculateEnhancedErrorRate(operation string, baseRate float64) float64

    // Execution
    InjectFailure(operation string, scenario *ChaosScenario) error
    SimulateDelay(operation string, baseDelay time.Duration) time.Duration
    SimulateResourceConstraint(operation string, resourceType ResourceType) error

    // Behavior tracking
    RecordUserAction(action UserAction) error
    GetBehaviorPattern() (*BehaviorPattern, error)
    AdjustDifficulty() error

    // Safety and monitoring
    ValidateSafetyBoundaries() error
    GetInjectionHistory() ([]InjectionEvent, error)
    EmergencyShutdown() error
}

// Production implementation
type SafeChaosInjector struct {
    // Configuration
    config          *ChaosConfig
    initialized     bool
    shutdownChan    chan struct{}

    // Core components
    scenarioManager *ScenarioManager
    behaviorTracker *BehaviorTracker
    safetyMonitor   *SafetyMonitor
    metricsCollector *MetricsCollector

    // State management
    operationHistory []InjectionEvent
    userContext      *UserContext
    performanceStats *PerformanceStats

    // Safety systems
    emergencyStop   bool
    resourceMonitor *ResourceMonitor
    healthChecker   *HealthChecker

    // Concurrency control
    mutex           sync.RWMutex
    operationLimiter *rate.Limiter
}
```

#### 2.2 Configuration Specifications
```go
// ChaosConfig defines all chaos marine configuration
type ChaosConfig struct {
    // Basic activation
    Enabled             bool                `json:"enabled" yaml:"enabled"`
    AggressivenessLevel AggressivenessLevel `json:"aggressiveness_level" yaml:"aggressiveness_level"`
    RandomSeed          int64               `json:"random_seed,omitempty" yaml:"random_seed,omitempty"`

    // Safety configuration
    SafetyMode          bool                `json:"safety_mode" yaml:"safety_mode"`
    MaxInjectionCount   int64               `json:"max_injection_count" yaml:"max_injection_count"`
    AllowedOperations   []string            `json:"allowed_operations" yaml:"allowed_operations"`
    ProhibitedPaths     []string            `json:"prohibited_paths" yaml:"prohibited_paths"`

    // User experience
    AdaptiveDifficulty  bool                `json:"adaptive_difficulty" yaml:"adaptive_difficulty"`
    EducationalMode     bool                `json:"educational_mode" yaml:"educational_mode"`
    RecoveryValidation  bool                `json:"recovery_validation" yaml:"recovery_validation"`
    ProgressiveHints    bool                `json:"progressive_hints" yaml:"progressive_hints"`

    // Performance limits
    MaxMemoryUsageMB    int64               `json:"max_memory_usage_mb" yaml:"max_memory_usage_mb"`
    MaxCPUUsagePercent  float64             `json:"max_cpu_usage_percent" yaml:"max_cpu_usage_percent"`
    OperationTimeoutSec int64               `json:"operation_timeout_sec" yaml:"operation_timeout_sec"`

    // Analytics and telemetry
    TelemetryEnabled    bool                `json:"telemetry_enabled" yaml:"telemetry_enabled"`
    AnonymousReporting  bool                `json:"anonymous_reporting" yaml:"anonymous_reporting"`
    MetricsRetentionDays int                `json:"metrics_retention_days" yaml:"metrics_retention_days"`

    // Advanced features
    CustomScenarios     map[string]*CustomScenario `json:"custom_scenarios,omitempty" yaml:"custom_scenarios,omitempty"`
    FailureChaining     bool                       `json:"failure_chaining" yaml:"failure_chaining"`
    CascadePrevent      bool                       `json:"cascade_prevent" yaml:"cascade_prevent"`
}

// Configuration validation
func (c *ChaosConfig) Validate() error {
    if c.MaxMemoryUsageMB < 1 || c.MaxMemoryUsageMB > 1024 {
        return errors.New("max_memory_usage_mb must be between 1 and 1024")
    }

    if c.MaxCPUUsagePercent < 0.1 || c.MaxCPUUsagePercent > 50.0 {
        return errors.New("max_cpu_usage_percent must be between 0.1 and 50.0")
    }

    if c.MaxInjectionCount < 1 || c.MaxInjectionCount > 10000 {
        return errors.New("max_injection_count must be between 1 and 10000")
    }

    return nil
}
```

#### 2.3 Enhanced Scenario Specifications
```go
// ChaosScenario extends existing ErrorScenario with chaos features
type ChaosScenario struct {
    // Inherit existing error scenario structure
    *errors.ErrorScenario

    // Chaos-specific metadata
    ID                  string              `json:"id"`
    Version             string              `json:"version"`
    CreatedAt           time.Time           `json:"created_at"`
    LastModified        time.Time           `json:"last_modified"`

    // Probability and triggering
    BaseProbability     float64             `json:"base_probability"`
    SkillModifiers      map[SkillLevel]float64 `json:"skill_modifiers"`
    ContextRequirements []ContextRequirement `json:"context_requirements"`
    CooldownPeriod      time.Duration       `json:"cooldown_period"`

    // Educational features
    LearningObjectives  []string            `json:"learning_objectives"`
    DifficultyRating    int                 `json:"difficulty_rating"` // 1-10
    PrerequisiteSkills  []string            `json:"prerequisite_skills"`
    ProgressiveHints    []ProgressiveHint   `json:"progressive_hints"`

    // Execution parameters
    MinDuration         time.Duration       `json:"min_duration"`
    MaxDuration         time.Duration       `json:"max_duration"`
    ResourceImpact      []ResourceImpact    `json:"resource_impact"`
    EnvironmentReqs     []string            `json:"environment_requirements"`

    // Recovery and validation
    ExpectedActions     []ExpectedAction    `json:"expected_actions"`
    ValidationCommands  []ValidationCommand `json:"validation_commands"`
    RecoveryTimeouts    []time.Duration     `json:"recovery_timeouts"`
    AlternativePaths    []RecoveryPath      `json:"alternative_paths"`

    // Failure chaining
    ChainableTo         []string            `json:"chainable_to"`
    ChainProbability    float64             `json:"chain_probability"`
    ChainCooldown       time.Duration       `json:"chain_cooldown"`
    MaxChainLength      int                 `json:"max_chain_length"`

    // Metrics and analytics
    UsageStats          *ScenarioStats      `json:"usage_stats,omitempty"`
    EffectivenessScore  float64             `json:"effectiveness_score"`
    UserSatisfaction    float64             `json:"user_satisfaction"`
    LearningOutcomes    []LearningOutcome   `json:"learning_outcomes"`
}

// Progressive hint system
type ProgressiveHint struct {
    Level       int           `json:"level"`        // 1=subtle, 5=explicit
    Trigger     HintTrigger   `json:"trigger"`      // When to show this hint
    Content     string        `json:"content"`      // Hint text
    ActionType  string        `json:"action_type"`  // "guidance", "command", "explanation"
    Delay       time.Duration `json:"delay"`        // Delay before showing
}

type HintTrigger struct {
    TimeElapsed     time.Duration `json:"time_elapsed,omitempty"`
    FailedAttempts  int           `json:"failed_attempts,omitempty"`
    UserRequest     bool          `json:"user_request,omitempty"`
    FrustrationLevel float64      `json:"frustration_level,omitempty"`
}
```

### 3. INTEGRATION SPECIFICATIONS

#### 3.1 Command Line Integration
**File**: `internal/commands/create.go`

```go
// Enhanced create command with chaos integration
func NewCreateCommand() *cobra.Command {
    var devOnly bool
    var template string

    // Chaos-specific flags
    var chaosEnabled bool
    var chaosLevel string
    var chaosSeed int64
    var chaosConfig string

    cmd := &cobra.Command{
        Use:   "create [APP_NAME]",
        Short: "Create a new React application",
        Long: `Create a new React application with modern development tooling.

Chaos Marine Options:
  --chaos-marine           Enable chaos injection for error simulation
  --chaos-level string     Set aggressiveness level (default|scout|aggressive|invasive|apocalyptic)
  --chaos-seed int         Set random seed for deterministic chaos (testing)
  --chaos-config string    Path to custom chaos configuration file

Examples:
  engx create MyApp --chaos-marine --chaos-level=scout
  engx create MyApp --chaos-marine=false  # Explicit disable
  engx create MyApp --chaos-level=apocalyptic --chaos-seed=12345`,
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            appName := args[0]

            // ... existing verbosity and configuration logic ...

            // Initialize chaos injector if enabled
            var chaosInjector chaos.ChaosInjector
            if chaosEnabled {
                config, err := chaos.LoadChaosConfig(chaosLevel, chaosSeed, chaosConfig)
                if err != nil {
                    return fmt.Errorf("failed to load chaos configuration: %w", err)
                }

                chaosInjector, err = chaos.NewSafeChaosInjector(config)
                if err != nil {
                    return fmt.Errorf("failed to initialize chaos injector: %w", err)
                }
                defer chaosInjector.Shutdown(context.Background())
            }

            // ... existing prompt collection logic ...

            // Enhanced model with chaos integration
            model := models.NewAppModelWithChaos("create", appName, flags, userConfig, chaosInjector)

            // ... existing TUI execution logic ...
        },
    }

    // Add chaos flags
    cmd.Flags().BoolVar(&chaosEnabled, "chaos-marine", false, "Enable chaos injection")
    cmd.Flags().StringVar(&chaosLevel, "chaos-level", "default", "Chaos aggressiveness level")
    cmd.Flags().Int64Var(&chaosSeed, "chaos-seed", 0, "Random seed for deterministic chaos")
    cmd.Flags().StringVar(&chaosConfig, "chaos-config", "", "Custom chaos configuration file")

    // ... existing flags ...

    return cmd
}
```

#### 3.2 Progress Tracker Integration
**File**: `internal/simulation/progress/tracker.go`

```go
// Enhanced tracker with chaos capabilities
type ChaosAwareTracker struct {
    *Tracker  // Embed existing tracker for backward compatibility

    // Chaos components
    chaosInjector chaos.ChaosInjector
    userContext   *chaos.UserContext
    injectionLog  []chaos.InjectionEvent
}

// Factory function for chaos-aware tracker
func NewChaosAwareTracker(steps []Step, injector chaos.ChaosInjector) *ChaosAwareTracker {
    baseTracker := NewTracker(steps)

    return &ChaosAwareTracker{
        Tracker:       baseTracker,
        chaosInjector: injector,
        userContext:   chaos.NewUserContext(),
        injectionLog:  make([]chaos.InjectionEvent, 0),
    }
}

// Enhanced step execution with chaos integration
func (t *ChaosAwareTracker) ExecuteCurrentStep() error {
    step := t.CurrentStepInfo()
    if step == nil {
        return errors.New("no current step available")
    }

    // Check for chaos injection
    if t.chaosInjector != nil && t.chaosInjector.IsEnabled() {
        shouldInject, err := t.chaosInjector.ShouldInject(step.Name)
        if err != nil {
            return fmt.Errorf("chaos injection decision failed: %w", err)
        }

        if shouldInject {
            return t.executeStepWithChaos(step)
        }
    }

    // Normal execution (existing behavior preserved)
    return t.executeStepNormally(step)
}

func (t *ChaosAwareTracker) executeStepWithChaos(step *Step) error {
    // Create operation context
    opContext := &chaos.OperationContext{
        OperationName: step.Name,
        BaseDuration:  step.Duration,
        BaseErrorRate: step.ErrorRate,
        UserContext:   t.userContext,
        StepIndex:     t.currentStep,
        TotalSteps:    len(t.steps),
    }

    // Select and execute chaos scenario
    scenario, err := t.chaosInjector.SelectScenario(step.Name, opContext)
    if err != nil {
        return fmt.Errorf("scenario selection failed: %w", err)
    }

    // Log injection event
    event := chaos.InjectionEvent{
        Timestamp:     time.Now(),
        OperationName: step.Name,
        ScenarioID:    scenario.ID,
        UserContext:   *t.userContext,
    }
    t.injectionLog = append(t.injectionLog, event)

    // Execute chaos scenario
    err = t.chaosInjector.InjectFailure(step.Name, scenario)
    if err != nil {
        // Record failed injection
        event.Result = chaos.InjectionResultFailed
        event.Error = err.Error()

        // Fall back to normal execution
        return t.executeStepNormally(step)
    }

    // Record successful injection
    event.Result = chaos.InjectionResultSuccess
    return nil
}
```

#### 3.3 TUI Model Integration
**File**: `internal/tui/models/app.go`

```go
// Enhanced AppModel with chaos integration
type AppModel struct {
    // ... existing fields ...

    // Chaos integration
    chaosInjector    chaos.ChaosInjector
    chaosEnabled     bool
    injectionHistory []chaos.InjectionEvent
    userBehavior     *chaos.BehaviorTracker
}

// Factory function with chaos support
func NewAppModelWithChaos(
    command, projectName string,
    flags []string,
    config *config.UserConfiguration,
    chaosInjector chaos.ChaosInjector,
) *AppModel {
    // Create base model
    model := NewAppModel(command, projectName, flags, config)

    // Add chaos integration if provided
    if chaosInjector != nil {
        model.chaosInjector = chaosInjector
        model.chaosEnabled = true
        model.userBehavior = chaos.NewBehaviorTracker()
        model.injectionHistory = make([]chaos.InjectionEvent, 0)

        // Create chaos-aware progress tracker
        chaosTracker := progress.NewChaosAwareTracker(model.steps, chaosInjector)
        model.progressTracker = chaosTracker
    }

    return model
}

// Enhanced update method with chaos behavior tracking
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Track user input for behavior analysis
        if m.chaosEnabled {
            userAction := chaos.UserAction{
                Type:        chaos.ActionTypeKeyPress,
                Key:         msg.String(),
                Timestamp:   time.Now(),
                Context:     m.getCurrentContext(),
            }
            m.userBehavior.RecordAction(userAction)
        }

        // ... existing key handling ...

    case progress.StepCompleteMsg:
        // Handle step completion with chaos awareness
        if m.chaosEnabled && msg.Failed {
            // This might be a chaos-injected failure
            m.handleChaosFailure(msg)
        }

        // ... existing step handling ...

    // ... other existing message handling ...
    }

    // ... existing update logic ...
}

func (m *AppModel) handleChaosFailure(msg progress.StepCompleteMsg) {
    // Check if this was a chaos-injected failure
    if injectionEvent := m.findInjectionForStep(msg.StepName); injectionEvent != nil {
        // Start recovery validation process
        m.startRecoveryValidation(injectionEvent)
    }
}
```

### 4. SAFETY SYSTEM SPECIFICATIONS

#### 4.1 Safety Boundary Implementation
```go
// SafetyMonitor ensures chaos operations never affect real system
type SafetyMonitor struct {
    config           *SafetyConfig
    systemSnapshot   *SystemSnapshot
    resourceMonitor  *ResourceMonitor
    emergencyTrigger chan struct{}
    healthTicker     *time.Ticker
}

type SafetyConfig struct {
    // File system protection
    AllowedPaths      []string `json:"allowed_paths"`
    ProhibitedPaths   []string `json:"prohibited_paths"`
    VirtualFSOnly     bool     `json:"virtual_fs_only"`

    // Network protection
    BlockRealNetwork  bool     `json:"block_real_network"`
    MockNetworkOnly   bool     `json:"mock_network_only"`

    // Process protection
    BlockProcessOps   bool     `json:"block_process_ops"`
    ProcessWhitelist  []string `json:"process_whitelist"`

    // Resource limits
    MaxMemoryMB       int64    `json:"max_memory_mb"`
    MaxCPUPercent     float64  `json:"max_cpu_percent"`
    MaxDurationSec    int64    `json:"max_duration_sec"`

    // Emergency controls
    HealthCheckInterval time.Duration `json:"health_check_interval"`
    EmergencyTimeout    time.Duration `json:"emergency_timeout"`
}

func (sm *SafetyMonitor) ValidateOperation(op *ChaosOperation) error {
    // Check file system access
    if err := sm.validateFileSystemAccess(op); err != nil {
        return fmt.Errorf("file system safety violation: %w", err)
    }

    // Check network access
    if err := sm.validateNetworkAccess(op); err != nil {
        return fmt.Errorf("network safety violation: %w", err)
    }

    // Check process operations
    if err := sm.validateProcessOperations(op); err != nil {
        return fmt.Errorf("process safety violation: %w", err)
    }

    // Check resource usage
    if err := sm.validateResourceUsage(op); err != nil {
        return fmt.Errorf("resource safety violation: %w", err)
    }

    return nil
}

func (sm *SafetyMonitor) validateFileSystemAccess(op *ChaosOperation) error {
    for _, path := range op.FileSystemAccess {
        // Check against prohibited paths
        for _, prohibited := range sm.config.ProhibitedPaths {
            if strings.HasPrefix(path, prohibited) {
                return fmt.Errorf("access to prohibited path: %s", path)
            }
        }

        // Check allowed paths if whitelist is configured
        if len(sm.config.AllowedPaths) > 0 {
            allowed := false
            for _, allowedPath := range sm.config.AllowedPaths {
                if strings.HasPrefix(path, allowedPath) {
                    allowed = true
                    break
                }
            }
            if !allowed {
                return fmt.Errorf("access to non-whitelisted path: %s", path)
            }
        }
    }

    return nil
}
```

#### 4.2 Resource Monitoring Specifications
```go
// ResourceMonitor tracks resource usage and enforces limits
type ResourceMonitor struct {
    config       *ResourceConfig
    startTime    time.Time
    memoryUsage  int64
    cpuUsage     float64
    diskUsage    int64
    networkUsage int64

    // Monitoring goroutines
    stopChan     chan struct{}
    alertChan    chan ResourceAlert
    metricsChan  chan ResourceMetrics
}

type ResourceConfig struct {
    MaxMemoryMB      int64         `json:"max_memory_mb"`
    MaxCPUPercent    float64       `json:"max_cpu_percent"`
    MaxDiskUsageMB   int64         `json:"max_disk_usage_mb"`
    MaxNetworkKBps   int64         `json:"max_network_kbps"`
    MonitorInterval  time.Duration `json:"monitor_interval"`
    AlertThreshold   float64       `json:"alert_threshold"` // Percentage of limit
}

func (rm *ResourceMonitor) StartMonitoring() error {
    rm.startTime = time.Now()

    // Start monitoring goroutine
    go rm.monitorResources()

    return nil
}

func (rm *ResourceMonitor) monitorResources() {
    ticker := time.NewTicker(rm.config.MonitorInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            metrics := rm.collectMetrics()

            // Check limits
            if alert := rm.checkLimits(metrics); alert != nil {
                select {
                case rm.alertChan <- *alert:
                case <-rm.stopChan:
                    return
                }
            }

            // Send metrics
            select {
            case rm.metricsChan <- metrics:
            case <-rm.stopChan:
                return
            }

        case <-rm.stopChan:
            return
        }
    }
}

func (rm *ResourceMonitor) checkLimits(metrics ResourceMetrics) *ResourceAlert {
    // Memory check
    if metrics.MemoryUsageMB > rm.config.MaxMemoryMB {
        return &ResourceAlert{
            Type:        AlertTypeMemoryLimit,
            Severity:    SeverityCritical,
            Message:     fmt.Sprintf("Memory usage (%d MB) exceeds limit (%d MB)", metrics.MemoryUsageMB, rm.config.MaxMemoryMB),
            Metrics:     metrics,
            Timestamp:   time.Now(),
        }
    }

    // CPU check
    if metrics.CPUUsagePercent > rm.config.MaxCPUPercent {
        return &ResourceAlert{
            Type:        AlertTypeCPULimit,
            Severity:    SeverityCritical,
            Message:     fmt.Sprintf("CPU usage (%.2f%%) exceeds limit (%.2f%%)", metrics.CPUUsagePercent, rm.config.MaxCPUPercent),
            Metrics:     metrics,
            Timestamp:   time.Now(),
        }
    }

    return nil
}
```

---

## PERFORMANCE SPECIFICATIONS

### 1. Performance Requirements
- **Chaos Disabled**: 0% measurable overhead (compile-time optimization)
- **Chaos Enabled**: <2% overall performance impact
- **Memory Usage**: <5MB additional memory consumption
- **Startup Time**: <50ms additional initialization time
- **Injection Latency**: <10ms decision time per operation

### 2. Performance Testing Framework
```go
// Performance testing suite
func BenchmarkChaosInjection(b *testing.B) {
    injector := chaos.NewSafeChaosInjector(testConfig)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        shouldInject, _ := injector.ShouldInject("test_operation")
        _ = shouldInject
    }
}

func BenchmarkChaosDisabled(b *testing.B) {
    // Test with chaos completely disabled
    const ChaosEnabled = false

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        if ChaosEnabled {
            // This code should be optimized away
            _ = performChaosCheck()
        }
    }
}
```

---

## TESTING SPECIFICATIONS

### 1. Testing Framework Requirements
- **Unit Tests**: >95% code coverage for all chaos components
- **Integration Tests**: Full system testing with/without chaos
- **Safety Tests**: Comprehensive boundary validation
- **Performance Tests**: Regression testing for all benchmarks
- **Deterministic Tests**: Reproducible chaos scenarios for CI/CD

### 2. Test Structure
```
internal/chaos/testing/
├── unit_test.go           # Unit tests for core components
├── integration_test.go    # Full system integration tests
├── safety_test.go         # Safety boundary validation
├── performance_test.go    # Performance regression tests
├── deterministic_test.go  # Reproducible chaos scenarios
└── fixtures/
    ├── test_configs/      # Test configuration files
    ├── mock_scenarios/    # Mock scenario definitions
    └── expected_outputs/  # Expected test outputs
```

---

**Technical Specifications Status**: ✅ COMPREHENSIVE SPECIFICATIONS COMPLETE
**Implementation Ready**: YES - All interfaces and integrations defined
**Safety Validated**: YES - Complete safety framework specified
**Performance Characterized**: YES - Detailed performance requirements defined
**Next Phase**: CODE - Begin Phase 1 Implementation