package chaos

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/simulation/progress"
)

// ChaosAwareTracker wraps the existing progress.Tracker with chaos injection capabilities
type ChaosAwareTracker struct {
	*progress.Tracker                // Embed existing tracker for full compatibility

	chaosInjector     ChaosInjector   // Chaos injection engine
	userBehavior      *BehaviorTracker // User behavior analysis
	injectionHistory  []InjectionEvent // History of chaos injections
	adaptationLog     []AdaptationEvent // History of difficulty adaptations

	// Configuration
	enabled          bool
	currentSession   string

	// State tracking
	stepFailures     map[int]bool      // Track which steps have failed due to chaos
	recoveryAttempts map[int]int       // Track recovery attempts per step

	// Thread safety
	mutex sync.RWMutex
}

// StepExecutionResult represents the result of executing a step with potential chaos
type StepExecutionResult struct {
	StepIndex        int           `json:"step_index"`
	StepName         string        `json:"step_name"`
	Success          bool          `json:"success"`
	ChaosInjected    bool          `json:"chaos_injected"`
	InjectedScenario string        `json:"injected_scenario,omitempty"`
	ExecutionTime    time.Duration `json:"execution_time"`
	ErrorMessage     string        `json:"error_message,omitempty"`
	RecoveryRequired bool          `json:"recovery_required"`
}

// NewChaosAwareTracker creates a new chaos-aware tracker wrapping an existing tracker
func NewChaosAwareTracker(baseTracker *progress.Tracker, injector ChaosInjector) *ChaosAwareTracker {
	tracker := &ChaosAwareTracker{
		Tracker:          baseTracker,
		chaosInjector:    injector,
		userBehavior:     NewBehaviorTracker(),
		injectionHistory: make([]InjectionEvent, 0),
		adaptationLog:    make([]AdaptationEvent, 0),
		enabled:          injector != nil && injector.IsEnabled(),
		stepFailures:     make(map[int]bool),
		recoveryAttempts: make(map[int]int),
	}

	// Start behavior tracking session
	if tracker.enabled {
		tracker.currentSession = tracker.userBehavior.StartSession()
	}

	return tracker
}

// ExecuteStep executes a step with potential chaos injection
func (cat *ChaosAwareTracker) ExecuteStep(stepIndex int) *StepExecutionResult {
	startTime := time.Now()

	step := cat.GetStep(stepIndex)
	if step == nil {
		return &StepExecutionResult{
			StepIndex:     stepIndex,
			Success:       false,
			ChaosInjected: false,
			ExecutionTime: time.Since(startTime),
			ErrorMessage:  "Invalid step index",
		}
	}

	result := &StepExecutionResult{
		StepIndex:     stepIndex,
		StepName:      step.Name,
		Success:       true,
		ChaosInjected: false,
		ExecutionTime: 0,
	}

	// Check if chaos should be injected for this step
	if cat.enabled && cat.shouldInjectChaosForStep(step) {
		result.ChaosInjected = true
		chaosResult := cat.executeChaosScenario(step)

		if chaosResult.Error != nil {
			result.Success = false
			result.ErrorMessage = chaosResult.Error.Error()
			result.InjectedScenario = chaosResult.ScenarioType
			result.RecoveryRequired = true

			// Mark this step as failed due to chaos
			cat.mutex.Lock()
			cat.stepFailures[stepIndex] = true
			cat.mutex.Unlock()

			// Record injection event
			cat.recordInjectionEvent(stepIndex, step.Name, chaosResult)
		}
	} else {
		// Execute normal step logic with existing error rate
		if cat.shouldStepFail(step) {
			result.Success = false
			result.ErrorMessage = fmt.Sprintf("Step failed: %s", step.Name)
			result.RecoveryRequired = step.CanRetry
		}
	}

	result.ExecutionTime = time.Since(startTime)

	// Record user action for behavior analysis
	if cat.enabled {
		action := UserAction{
			Timestamp:  startTime,
			ActionType: CommandExecution,
			Command:    step.Name,
			Context:    "step_execution",
			Success:    result.Success,
			Duration:   result.ExecutionTime,
		}
		cat.userBehavior.RecordAction(action)
	}

	return result
}

// shouldInjectChaosForStep determines if chaos should be injected for a specific step
func (cat *ChaosAwareTracker) shouldInjectChaosForStep(step *progress.Step) bool {
	if !cat.enabled || cat.chaosInjector == nil {
		return false
	}

	// Use the step name as the operation identifier
	return cat.chaosInjector.ShouldInject(step.Name)
}

// executeChaosScenario executes a chaos scenario for the given step
func (cat *ChaosAwareTracker) executeChaosScenario(step *progress.Step) *ChaosExecutionResult {
	result := &ChaosExecutionResult{
		StepName:     step.Name,
		StartTime:    time.Now(),
		Success:      true,
	}

	// Select and execute scenario
	scenario := cat.chaosInjector.SelectScenario(step.Name)
	if scenario != nil {
		result.ScenarioType = scenario.ErrorScenario.Type
		err := cat.chaosInjector.InjectFailure(step.Name, scenario)
		if err != nil {
			result.Error = err
			result.Success = false
		}
	} else {
		// No scenario found, but still inject a failure based on enhanced error rate
		enhancedRate := cat.chaosInjector.CalculateEnhancedErrorRate(step.Name, step.ErrorRate)
		random := rand.New(rand.NewSource(time.Now().UnixNano()))

		if random.Float64() < enhancedRate {
			result.Error = fmt.Errorf("CHAOS INJECTION: Enhanced failure for %s", step.Name)
			result.Success = false
			result.ScenarioType = "enhanced_failure"
		}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result
}

// shouldStepFail determines if a step should fail based on its natural error rate
func (cat *ChaosAwareTracker) shouldStepFail(step *progress.Step) bool {
	if step.ErrorRate <= 0 {
		return false
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return random.Float64() < step.ErrorRate
}

// recordInjectionEvent records a chaos injection event
func (cat *ChaosAwareTracker) recordInjectionEvent(stepIndex int, stepName string, chaosResult *ChaosExecutionResult) {
	cat.mutex.Lock()
	defer cat.mutex.Unlock()

	event := InjectionEvent{
		Timestamp: chaosResult.StartTime,
		Operation: stepName,
		Scenario:  chaosResult.ScenarioType,
		Success:   chaosResult.Success,
		Duration:  chaosResult.Duration,
	}

	cat.injectionHistory = append(cat.injectionHistory, event)
}

// AttemptStepRecovery attempts to recover from a failed step
func (cat *ChaosAwareTracker) AttemptStepRecovery(stepIndex int) (*RecoveryResult, error) {
	cat.mutex.Lock()
	defer cat.mutex.Unlock()

	step := cat.GetStep(stepIndex)
	if step == nil {
		return nil, fmt.Errorf("invalid step index: %d", stepIndex)
	}

	// Check if this step failed due to chaos
	failedDueToChaos, exists := cat.stepFailures[stepIndex]
	if !exists || !failedDueToChaos {
		return nil, fmt.Errorf("step %d did not fail due to chaos", stepIndex)
	}

	// Increment recovery attempts
	cat.recoveryAttempts[stepIndex]++
	attempts := cat.recoveryAttempts[stepIndex]

	result := &RecoveryResult{
		StepIndex:       stepIndex,
		StepName:        step.Name,
		AttemptNumber:   attempts,
		StartTime:       time.Now(),
	}

	// Recovery logic depends on the number of attempts and user behavior
	if cat.enabled && cat.chaosInjector != nil {
		pattern := cat.userBehavior.GetCurrentPattern()

		// Provide progressive assistance based on attempt number and user skill
		if attempts == 1 {
			// First attempt - minimal help
			result.AssistanceLevel = MinimalAssistance
			result.Success = cat.attemptBasicRecovery(step, pattern)
		} else if attempts == 2 {
			// Second attempt - provide hints
			result.AssistanceLevel = HintProvided
			result.Success = cat.attemptGuidedRecovery(step, pattern)
			result.Hint = cat.generateRecoveryHint(step, pattern)
		} else {
			// Third+ attempt - provide solution
			result.AssistanceLevel = SolutionProvided
			result.Success = true // Always succeed with solution
			result.Solution = cat.generateRecoverySolution(step, pattern)
		}
	} else {
		// Basic recovery without chaos intelligence
		result.Success = step.CanRetry
		result.AssistanceLevel = NoAssistance
	}

	// Clear failure state if recovery was successful
	if result.Success {
		delete(cat.stepFailures, stepIndex)
		delete(cat.recoveryAttempts, stepIndex)
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// Record recovery action for behavior analysis
	if cat.enabled {
		action := UserAction{
			Timestamp:  result.StartTime,
			ActionType: RecoveryAction,
			Command:    fmt.Sprintf("recover_%s", step.Name),
			Context:    fmt.Sprintf("attempt_%d", attempts),
			Success:    result.Success,
			Duration:   result.Duration,
		}
		cat.userBehavior.RecordAction(action)
	}

	return result, nil
}

// attemptBasicRecovery attempts basic recovery without assistance
func (cat *ChaosAwareTracker) attemptBasicRecovery(step *progress.Step, pattern *BehaviorPattern) bool {
	// Success rate depends on user skill level and step complexity
	baseSuccessRate := 0.3 // 30% base success rate for unassisted recovery

	if pattern != nil {
		switch pattern.SkillLevel {
		case Expert:
			baseSuccessRate = 0.8
		case Advanced:
			baseSuccessRate = 0.6
		case Intermediate:
			baseSuccessRate = 0.4
		case Novice:
			baseSuccessRate = 0.2
		}
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return random.Float64() < baseSuccessRate
}

// attemptGuidedRecovery attempts recovery with hints
func (cat *ChaosAwareTracker) attemptGuidedRecovery(step *progress.Step, pattern *BehaviorPattern) bool {
	// Higher success rate with hints
	baseSuccessRate := 0.7

	if pattern != nil {
		switch pattern.SkillLevel {
		case Expert:
			baseSuccessRate = 0.95
		case Advanced:
			baseSuccessRate = 0.85
		case Intermediate:
			baseSuccessRate = 0.75
		case Novice:
			baseSuccessRate = 0.6
		}
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return random.Float64() < baseSuccessRate
}

// generateRecoveryHint generates a helpful hint for step recovery
func (cat *ChaosAwareTracker) generateRecoveryHint(step *progress.Step, pattern *BehaviorPattern) string {
	// This would be more sophisticated in a full implementation
	hints := map[string]string{
		"Validating configuration":       "💡 Check your project name for invalid characters or existing conflicts",
		"Setting up environment":         "💡 Ensure you have write permissions to the current directory",
		"Installing dependencies":        "💡 Check your network connection and try clearing npm cache",
		"Generating project structure":   "💡 Verify available disk space and file system permissions",
		"Configuring production setup":   "💡 Check environment variables and deployment configuration",
		"Installing Testing Frameworks":  "💡 Ensure compatible Node.js version for testing tools",
		"Generating Documentation":       "💡 Check for documentation template conflicts or missing files",
	}

	hint, exists := hints[step.Name]
	if !exists {
		hint = fmt.Sprintf("💡 Try rerunning the '%s' step after checking system resources", step.Name)
	}

	return hint
}

// generateRecoverySolution generates a complete solution for step recovery
func (cat *ChaosAwareTracker) generateRecoverySolution(step *progress.Step, pattern *BehaviorPattern) string {
	solutions := map[string]string{
		"Validating configuration":       "🔧 Run: npm init --force && npm config set registry https://registry.npmjs.org/",
		"Setting up environment":         "🔧 Run: mkdir -p temp_project && cd temp_project",
		"Installing dependencies":        "🔧 Run: npm cache clean --force && npm install --no-optional",
		"Generating project structure":   "🔧 Run: rm -rf node_modules && npm install",
		"Configuring production setup":   "🔧 Run: npm run build --if-present",
		"Installing Testing Frameworks":  "🔧 Run: npm install --save-dev vitest @testing-library/react",
		"Generating Documentation":       "🔧 Run: npm run docs:generate --if-present",
	}

	solution, exists := solutions[step.Name]
	if !exists {
		solution = fmt.Sprintf("🔧 Retry the '%s' step with elevated permissions", step.Name)
	}

	return solution
}

// GenerateErrorTemplate creates an error template for a failed step
func (cat *ChaosAwareTracker) GenerateErrorTemplate(stepIndex int, result *StepExecutionResult) *ErrorTemplate {
	if !result.ChaosInjected || result.Success {
		return nil
	}

	step := cat.GetStep(stepIndex)
	if step == nil {
		return nil
	}

	config := cat.chaosInjector.GetConfig()

	// Create error template
	template := NewChaosErrorTemplate(
		result.InjectedScenario,
		"create",
		step.Name,
		config.AggressivenessLevel,
	)

	// Use predefined template if available, otherwise create custom
	if predefined, exists := DefaultErrorTemplates[result.InjectedScenario]; exists {
		template.BottomLineMessage = predefined.BottomLineMessage
		template.FirstAction = predefined.FirstAction
		template.SecondAction = predefined.SecondAction
		template.Summary = predefined.Summary
		template.AdditionalContext = predefined.AdditionalContext
		template.OnCallCrew = predefined.OnCallCrew
		template.StackTrace = cat.generateStepSpecificStackTrace(step.Name, result.InjectedScenario)
	} else {
		// Generate custom template for unknown scenario
		template.BottomLineMessage = fmt.Sprintf("%s failed: %s", step.Name, result.ErrorMessage)
		template.FirstAction = "Retry the operation"
		template.SecondAction = "Check system resources and try again"
		template.Summary = fmt.Sprintf("The '%s' step encountered an error during execution.", step.Name)
		template.AdditionalContext = "This error was generated by the Chaos Marine for educational purposes."
		template.OnCallCrew = "DevOps Team"
		template.StackTrace = cat.generateStepSpecificStackTrace(step.Name, result.InjectedScenario)
	}

	return template
}

// generateStepSpecificStackTrace creates a realistic stack trace for the specific step and error type
func (cat *ChaosAwareTracker) generateStepSpecificStackTrace(stepName, scenarioType string) string {
	switch stepName {
	case "Installing dependencies":
		if scenarioType == "network_failure" {
			return `  Stack Trace:
   at npmInstall (/usr/local/lib/node_modules/engx/lib/install.js:142:15)
   at fetchPackage (/usr/local/lib/node_modules/engx/lib/fetch.js:87:12)
   at Registry.request (/usr/local/lib/node_modules/engx/lib/registry.js:234:9)
   at ClientRequest.onError (/usr/local/lib/node_modules/engx/lib/request.js:156:21)
   Error: ENOTFOUND registry.npmjs.org
   Code: ENOTFOUND
   Errno: -3008`
		}

	case "Setting up environment":
		if scenarioType == "permission_denied" {
			return `  Stack Trace:
   at mkdir (/usr/local/lib/node_modules/engx/lib/filesystem.js:45:18)
   at createProjectStructure (/usr/local/lib/node_modules/engx/lib/creator.js:123:7)
   at executeStep (/usr/local/lib/node_modules/engx/lib/runner.js:89:12)
   at processStep (/usr/local/lib/node_modules/engx/lib/progress.js:67:5)
   Error: EACCES: permission denied, mkdir '/restricted/project'
   Code: EACCES
   Errno: -13`
		}

	case "Generating project structure":
		if scenarioType == "resource_exhausted" {
			return `  Stack Trace:
   at writeFile (/usr/local/lib/node_modules/engx/lib/filesystem.js:78:14)
   at generateFiles (/usr/local/lib/node_modules/engx/lib/generator.js:156:9)
   at createProject (/usr/local/lib/node_modules/engx/lib/creator.js:201:11)
   at executeCommand (/usr/local/lib/node_modules/engx/lib/runner.js:134:8)
   Error: ENOSPC: no space left on device, write
   Code: ENOSPC
   Errno: -28`
		}
	}

	// Default stack trace
	return `  Stack Trace:
   at executeOperation (/usr/local/lib/node_modules/engx/lib/runner.js:234:18)
   at processCommand (/usr/local/lib/node_modules/engx/lib/processor.js:78:11)
   at main (/usr/local/lib/node_modules/engx/bin/engx:145:7)
   Error: ` + scenarioType + ` error occurred in ` + stepName
}

// AdaptDifficulty adapts the chaos difficulty based on user behavior
func (cat *ChaosAwareTracker) AdaptDifficulty() {
	if !cat.enabled || cat.chaosInjector == nil {
		return
	}

	cat.mutex.Lock()
	defer cat.mutex.Unlock()

	pattern := cat.userBehavior.GetCurrentPattern()
	if pattern == nil {
		return
	}

	previousLevel := cat.chaosInjector.GetAggressivenessLevel()
	newLevel := cat.chaosInjector.AdjustDifficulty(pattern)

	if newLevel != previousLevel {
		event := AdaptationEvent{
			Timestamp:     time.Now(),
			PreviousLevel: previousLevel,
			NewLevel:      newLevel,
			Reason:        cat.generateAdaptationReason(pattern),
			UserMetrics:   cat.userBehavior.competenceMetrics,
		}

		cat.adaptationLog = append(cat.adaptationLog, event)
	}
}

// generateAdaptationReason generates a human-readable reason for difficulty adaptation
func (cat *ChaosAwareTracker) generateAdaptationReason(pattern *BehaviorPattern) string {
	if pattern.ShowsFrustration {
		return "User showing frustration patterns - reducing difficulty"
	}

	if pattern.RecentSuccessRate > 0.9 && pattern.SkillLevel >= Advanced {
		return "High success rate with advanced skill level - increasing challenge"
	}

	if pattern.RecentSuccessRate < 0.4 {
		return "Low success rate detected - reducing difficulty to improve learning experience"
	}

	if pattern.SkillLevel == Expert && pattern.ConfidenceLevel > 0.8 {
		return "Expert user with high confidence - providing maximum challenge"
	}

	return "Adjusting difficulty based on user behavior patterns"
}

// GetChaosMetrics returns comprehensive chaos injection metrics
func (cat *ChaosAwareTracker) GetChaosMetrics() *ChaosMetrics {
	cat.mutex.RLock()
	defer cat.mutex.RUnlock()

	totalSteps := cat.TotalSteps()
	failedSteps := len(cat.stepFailures)
	totalInjections := len(cat.injectionHistory)

	successfulInjections := 0
	for _, event := range cat.injectionHistory {
		if event.Success {
			successfulInjections++
		}
	}

	totalRecoveryAttempts := 0
	for _, attempts := range cat.recoveryAttempts {
		totalRecoveryAttempts += attempts
	}

	pattern := cat.userBehavior.GetCurrentPattern()

	return &ChaosMetrics{
		Enabled:               cat.enabled,
		TotalSteps:            totalSteps,
		FailedSteps:           failedSteps,
		TotalInjections:       totalInjections,
		SuccessfulInjections:  successfulInjections,
		TotalRecoveryAttempts: totalRecoveryAttempts,
		CurrentSession:        cat.currentSession,
		UserBehaviorPattern:   pattern,
		InjectionHistory:      cat.injectionHistory,
		AdaptationHistory:     cat.adaptationLog,
	}
}

// GetChaosInjector returns the chaos injector
func (cat *ChaosAwareTracker) GetChaosInjector() ChaosInjector {
	return cat.chaosInjector
}

// Reset resets the chaos-aware tracker state
func (cat *ChaosAwareTracker) Reset() {
	cat.mutex.Lock()
	defer cat.mutex.Unlock()

	// Reset base tracker
	cat.Tracker.Reset()

	// Reset chaos state
	cat.stepFailures = make(map[int]bool)
	cat.recoveryAttempts = make(map[int]int)
	cat.injectionHistory = make([]InjectionEvent, 0)
	cat.adaptationLog = make([]AdaptationEvent, 0)

	// Reset behavior tracking
	if cat.enabled && cat.userBehavior != nil {
		cat.userBehavior.Reset()
		cat.currentSession = cat.userBehavior.StartSession()
	}

	// Reset chaos injector
	if cat.enabled && cat.chaosInjector != nil {
		cat.chaosInjector.ResetState()
	}
}

// Supporting types and structures

// ChaosExecutionResult represents the result of executing a chaos scenario
type ChaosExecutionResult struct {
	StepName     string        `json:"step_name"`
	ScenarioType string        `json:"scenario_type"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	Duration     time.Duration `json:"duration"`
	Success      bool          `json:"success"`
	Error        error         `json:"error,omitempty"`
}

// RecoveryResult represents the result of a recovery attempt
type RecoveryResult struct {
	StepIndex       int               `json:"step_index"`
	StepName        string            `json:"step_name"`
	AttemptNumber   int               `json:"attempt_number"`
	StartTime       time.Time         `json:"start_time"`
	EndTime         time.Time         `json:"end_time"`
	Duration        time.Duration     `json:"duration"`
	Success         bool              `json:"success"`
	AssistanceLevel AssistanceLevel   `json:"assistance_level"`
	Hint            string            `json:"hint,omitempty"`
	Solution        string            `json:"solution,omitempty"`
}

// AssistanceLevel defines the level of assistance provided during recovery
type AssistanceLevel int

const (
	NoAssistance AssistanceLevel = iota
	MinimalAssistance
	HintProvided
	SolutionProvided
)

// ChaosMetrics represents comprehensive chaos injection metrics
type ChaosMetrics struct {
	Enabled               bool                  `json:"enabled"`
	TotalSteps            int                   `json:"total_steps"`
	FailedSteps           int                   `json:"failed_steps"`
	TotalInjections       int                   `json:"total_injections"`
	SuccessfulInjections  int                   `json:"successful_injections"`
	TotalRecoveryAttempts int                   `json:"total_recovery_attempts"`
	CurrentSession        string                `json:"current_session"`
	UserBehaviorPattern   *BehaviorPattern      `json:"user_behavior_pattern"`
	InjectionHistory      []InjectionEvent      `json:"injection_history"`
	AdaptationHistory     []AdaptationEvent     `json:"adaptation_history"`
}

// String methods for enums
func (al AssistanceLevel) String() string {
	switch al {
	case NoAssistance:
		return "none"
	case MinimalAssistance:
		return "minimal"
	case HintProvided:
		return "hint"
	case SolutionProvided:
		return "solution"
	default:
		return "unknown"
	}
}