package chaos

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ChaosInjector defines the interface for chaos injection operations
type ChaosInjector interface {
	// Configuration
	LoadConfig(config *ChaosConfig) error
	IsEnabled() bool
	GetAggressivenessLevel() AggressivenessLevel

	// Injection decisions
	ShouldInject(operation string) bool
	SelectScenario(operation string) *ChaosScenario
	CalculateEnhancedErrorRate(operation string, baseRate float64) float64

	// Failure execution
	InjectFailure(operation string, scenario *ChaosScenario) error
	ExecuteScenario(operation string) error

	// User behavior tracking
	RecordUserAction(action UserAction) error
	AnalyzeBehaviorPattern() *BehaviorPattern
	AdjustDifficulty(pattern *BehaviorPattern) AggressivenessLevel

	// Safety and monitoring
	ValidateSafetyBoundaries() error
	GetOperationHistory() []InjectionEvent
	ResetState() error
}

// SafeChaosInjector implements the ChaosInjector interface with safety-first design
type SafeChaosInjector struct {
	config        *ChaosConfig
	scenarios     map[string]*ChaosScenario
	userBehavior  *BehaviorTracker
	safetyMonitor *SafetyMonitor
	operationLog  []InjectionEvent
	random        *rand.Rand
	metrics       *InjectionMetrics
	mutex         sync.RWMutex
	startTime     time.Time
}

// ErrorScenario represents a basic error scenario (placeholder for future error system)
type ErrorScenario struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// ChaosScenario extends the basic error scenario framework
type ChaosScenario struct {
	*ErrorScenario

	// Chaos-specific enhancements
	TriggerProbability    float64                    `json:"trigger_probability"`
	UserSkillModifier     map[SkillLevel]float64     `json:"user_skill_modifier"`
	ChainableFailures     []string                   `json:"chainable_failures"`

	// Educational features
	LearningObjectives    []string                   `json:"learning_objectives"`
	SkillAssessment      *SkillTest                 `json:"skill_assessment"`
	ProgressiveHints     []string                   `json:"progressive_hints"`

	// Simulation parameters
	MinDuration          time.Duration              `json:"min_duration"`
	MaxDuration          time.Duration              `json:"max_duration"`
	ResourceRequirements []ResourceType             `json:"resource_requirements"`

	// Recovery validation
	ExpectedActions      []ExpectedAction           `json:"expected_actions"`
	ValidationScript     string                     `json:"validation_script"`
	AlternativeApproaches []RecoveryPath            `json:"alternative_approaches"`
}

// UserAction represents a user's action during chaos scenarios
type UserAction struct {
	Timestamp   time.Time     `json:"timestamp"`
	ActionType  ActionType    `json:"action_type"`
	Command     string        `json:"command"`
	Context     string        `json:"context"`
	Success     bool          `json:"success"`
	Duration    time.Duration `json:"duration"`
}

// ActionType defines the types of user actions we track
type ActionType int

const (
	CommandExecution ActionType = iota
	FileOperation
	ConfigurationChange
	HelpRequest
	RetryAttempt
	RecoveryAction
)

// BehaviorPattern represents analyzed user behavior patterns
type BehaviorPattern struct {
	SkillLevel           SkillLevel    `json:"skill_level"`
	RecentSuccessRate    float64       `json:"recent_success_rate"`
	AverageResolutionTime time.Duration `json:"average_resolution_time"`
	HelpRequestFrequency float64       `json:"help_request_frequency"`
	RetryPatterns        []RetryPattern `json:"retry_patterns"`
	ShowsFrustration     bool          `json:"shows_frustration"`
	ConfidenceLevel      float64       `json:"confidence_level"`
}

// InjectionEvent represents a chaos injection event
type InjectionEvent struct {
	Timestamp   time.Time      `json:"timestamp"`
	Operation   string         `json:"operation"`
	Scenario    string         `json:"scenario"`
	Success     bool           `json:"success"`
	UserResponse *UserResponse `json:"user_response,omitempty"`
	Duration    time.Duration  `json:"duration"`
}

// InjectionMetrics tracks chaos injection performance and effectiveness
type InjectionMetrics struct {
	TotalInjections     int64         `json:"total_injections"`
	SuccessfulInjections int64        `json:"successful_injections"`
	UserResolutions     int64         `json:"user_resolutions"`
	AverageResolutionTime time.Duration `json:"average_resolution_time"`
	EducationalValue    float64       `json:"educational_value"`
}

// NewSafeChaosInjector creates a new chaos injector with safety guarantees
func NewSafeChaosInjector(config *ChaosConfig) (*SafeChaosInjector, error) {
	if config == nil {
		return nil, errors.New("chaos configuration cannot be nil")
	}

	// Initialize safety monitor
	safetyMonitor, err := NewSafetyMonitor(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize safety monitor: %w", err)
	}

	// Initialize random number generator
	var rng *rand.Rand
	if config.RandomSeed != 0 {
		// Deterministic mode for testing
		rng = rand.New(rand.NewSource(config.RandomSeed))
	} else {
		// Random seed for normal operation
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	// Initialize behavior tracker
	behaviorTracker := NewBehaviorTracker()

	// Initialize scenarios
	scenarios := initializeDefaultScenarios()

	injector := &SafeChaosInjector{
		config:        config,
		scenarios:     scenarios,
		userBehavior:  behaviorTracker,
		safetyMonitor: safetyMonitor,
		operationLog:  make([]InjectionEvent, 0),
		random:        rng,
		metrics:       &InjectionMetrics{},
		startTime:     time.Now(),
	}

	return injector, nil
}

// LoadConfig loads a new configuration
func (injector *SafeChaosInjector) LoadConfig(config *ChaosConfig) error {
	injector.mutex.Lock()
	defer injector.mutex.Unlock()

	if config == nil {
		return errors.New("configuration cannot be nil")
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	injector.config = config
	return nil
}

// IsEnabled returns whether chaos injection is currently enabled
func (injector *SafeChaosInjector) IsEnabled() bool {
	injector.mutex.RLock()
	defer injector.mutex.RUnlock()
	return injector.config.Enabled
}

// GetAggressivenessLevel returns the current aggressiveness level
func (injector *SafeChaosInjector) GetAggressivenessLevel() AggressivenessLevel {
	injector.mutex.RLock()
	defer injector.mutex.RUnlock()
	return injector.config.AggressivenessLevel
}

// ShouldInject determines if chaos should be injected for the given operation
func (injector *SafeChaosInjector) ShouldInject(operation string) bool {
	injector.mutex.RLock()
	defer injector.mutex.RUnlock()

	// Quick exit if disabled
	if !injector.config.Enabled {
		return false
	}

	// Safety check
	if err := injector.safetyMonitor.IsOperationSafe(operation); err != nil {
		return false
	}

	// Check operation is allowed
	if !injector.config.IsOperationAllowed(operation) {
		return false
	}

	// Get base failure rate
	baseRate := injector.config.AggressivenessLevel.FailureRate()

	// Apply adaptive difficulty if enabled
	if injector.config.AdaptiveDifficulty {
		pattern := injector.userBehavior.GetCurrentPattern()
		baseRate = injector.applyAdaptiveDifficulty(baseRate, pattern)
	}

	// Make injection decision based on probability
	return injector.random.Float64() < baseRate
}

// applyAdaptiveDifficulty adjusts the base failure rate based on user behavior
func (injector *SafeChaosInjector) applyAdaptiveDifficulty(baseRate float64, pattern *BehaviorPattern) float64 {
	if pattern == nil {
		return baseRate
	}

	adjustment := 1.0

	// Skill-based adjustment
	switch pattern.SkillLevel {
	case Novice:
		adjustment *= 0.5 // Reduce difficulty for beginners
	case Intermediate:
		adjustment *= 0.8 // Slightly reduce for intermediate
	case Advanced:
		adjustment *= 1.2 // Increase for advanced users
	case Expert:
		adjustment *= 1.5 // Significant increase for experts
	}

	// Recent performance adjustment
	if pattern.RecentSuccessRate < 0.3 {
		adjustment *= 0.6 // Ease up on struggling users
	} else if pattern.RecentSuccessRate > 0.9 {
		adjustment *= 1.3 // Challenge successful users
	}

	// Frustration indicators
	if pattern.ShowsFrustration {
		adjustment *= 0.4 // Significantly reduce difficulty
	}

	// Apply adjustment with reasonable bounds
	adjustedRate := baseRate * adjustment
	return clamp(adjustedRate, 0.0, 0.5) // Cap at 50% maximum
}

// clamp ensures a value is within the specified range
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// SelectScenario selects an appropriate chaos scenario for the operation
func (injector *SafeChaosInjector) SelectScenario(operation string) *ChaosScenario {
	injector.mutex.RLock()
	defer injector.mutex.RUnlock()

	// Get scenarios applicable to this operation
	candidates := injector.getApplicableScenarios(operation)
	if len(candidates) == 0 {
		return nil
	}

	// Apply user skill-based filtering if adaptive difficulty is enabled
	if injector.config.AdaptiveDifficulty {
		pattern := injector.userBehavior.GetCurrentPattern()
		candidates = injector.filterScenariosBySkill(candidates, pattern)
	}

	if len(candidates) == 0 {
		return nil
	}

	// Select scenario based on weighted probability
	return injector.selectWeightedScenario(candidates)
}

// getApplicableScenarios returns scenarios that apply to the given operation
func (injector *SafeChaosInjector) getApplicableScenarios(operation string) []*ChaosScenario {
	candidates := make([]*ChaosScenario, 0)

	for _, scenario := range injector.scenarios {
		if injector.isScenarioApplicable(scenario, operation) {
			candidates = append(candidates, scenario)
		}
	}

	return candidates
}

// isScenarioApplicable checks if a scenario applies to the given operation
func (injector *SafeChaosInjector) isScenarioApplicable(scenario *ChaosScenario, operation string) bool {
	// For now, implement basic matching logic
	// In a full implementation, this would be more sophisticated
	return scenario.ErrorScenario != nil
}

// filterScenariosBySkill filters scenarios based on user skill level
func (injector *SafeChaosInjector) filterScenariosBySkill(scenarios []*ChaosScenario, pattern *BehaviorPattern) []*ChaosScenario {
	if pattern == nil {
		return scenarios
	}

	filtered := make([]*ChaosScenario, 0)
	for _, scenario := range scenarios {
		if injector.isScenarioAppropriateForSkill(scenario, pattern.SkillLevel) {
			filtered = append(filtered, scenario)
		}
	}

	return filtered
}

// isScenarioAppropriateForSkill checks if a scenario is appropriate for the user's skill level
func (injector *SafeChaosInjector) isScenarioAppropriateForSkill(scenario *ChaosScenario, skillLevel SkillLevel) bool {
	// Check if scenario has skill level modifiers
	if modifier, exists := scenario.UserSkillModifier[skillLevel]; exists {
		// If modifier is very low, skip this scenario for this skill level
		return modifier > 0.1
	}

	// Default: allow all scenarios for all skill levels
	return true
}

// selectWeightedScenario selects a scenario using weighted random selection
func (injector *SafeChaosInjector) selectWeightedScenario(scenarios []*ChaosScenario) *ChaosScenario {
	if len(scenarios) == 0 {
		return nil
	}

	// Calculate total weight
	totalWeight := 0.0
	for _, scenario := range scenarios {
		totalWeight += scenario.TriggerProbability
	}

	if totalWeight == 0 {
		// If no weights, select randomly
		return scenarios[injector.random.Intn(len(scenarios))]
	}

	// Weighted random selection
	target := injector.random.Float64() * totalWeight
	currentWeight := 0.0

	for _, scenario := range scenarios {
		currentWeight += scenario.TriggerProbability
		if currentWeight >= target {
			return scenario
		}
	}

	// Fallback to last scenario
	return scenarios[len(scenarios)-1]
}

// CalculateEnhancedErrorRate calculates an enhanced error rate for existing error scenarios
func (injector *SafeChaosInjector) CalculateEnhancedErrorRate(operation string, baseRate float64) float64 {
	injector.mutex.RLock()
	defer injector.mutex.RUnlock()

	if !injector.config.Enabled {
		return baseRate
	}

	// Apply aggressiveness multiplier
	aggressivenessRate := injector.config.AggressivenessLevel.FailureRate()
	enhancedRate := baseRate + aggressivenessRate

	// Apply adaptive difficulty if enabled
	if injector.config.AdaptiveDifficulty {
		pattern := injector.userBehavior.GetCurrentPattern()
		enhancedRate = injector.applyAdaptiveDifficulty(enhancedRate, pattern)
	}

	return clamp(enhancedRate, 0.0, 0.95) // Cap at 95% to ensure some operations succeed
}

// InjectFailure injects a specific failure scenario
func (injector *SafeChaosInjector) InjectFailure(operation string, scenario *ChaosScenario) error {
	startTime := time.Now()

	// Safety verification
	if err := injector.safetyMonitor.IsOperationSafe(operation); err != nil {
		return fmt.Errorf("safety check failed: %w", err)
	}

	// Record injection attempt
	if err := injector.safetyMonitor.RecordInjection(operation); err != nil {
		return fmt.Errorf("injection recording failed: %w", err)
	}

	// Execute the scenario (simulation only)
	err := injector.simulateFailure(scenario)

	// Record injection event
	event := InjectionEvent{
		Timestamp: startTime,
		Operation: operation,
		Scenario:  scenario.ErrorScenario.Type,
		Success:   err == nil,
		Duration:  time.Since(startTime),
	}

	injector.mutex.Lock()
	injector.operationLog = append(injector.operationLog, event)
	injector.metrics.TotalInjections++
	if err == nil {
		injector.metrics.SuccessfulInjections++
	}
	injector.mutex.Unlock()

	return err
}

// ExecuteScenario executes a chaos scenario for the given operation
func (injector *SafeChaosInjector) ExecuteScenario(operation string) error {
	scenario := injector.SelectScenario(operation)
	if scenario == nil {
		return errors.New("no applicable scenario found")
	}

	return injector.InjectFailure(operation, scenario)
}

// simulateFailure simulates a failure scenario without actual system impact
func (injector *SafeChaosInjector) simulateFailure(scenario *ChaosScenario) error {
	// This is where the actual chaos simulation would happen
	// For now, we'll implement basic simulation logic

	if scenario.ErrorScenario == nil {
		return errors.New("invalid scenario: missing error scenario")
	}

	// Simulate processing time
	minDuration := scenario.MinDuration
	if minDuration == 0 {
		minDuration = 100 * time.Millisecond
	}

	maxDuration := scenario.MaxDuration
	if maxDuration == 0 {
		maxDuration = 2 * time.Second
	}

	// Random duration within bounds
	duration := minDuration + time.Duration(injector.random.Float64()*float64(maxDuration-minDuration))
	time.Sleep(duration)

	// Return the simulated error (this would integrate with existing error scenarios)
	return fmt.Errorf("CHAOS INJECTION: %s - %s",
		scenario.ErrorScenario.Type,
		scenario.ErrorScenario.Message)
}

// RecordUserAction records a user action for behavior analysis
func (injector *SafeChaosInjector) RecordUserAction(action UserAction) error {
	return injector.userBehavior.RecordAction(action)
}

// AnalyzeBehaviorPattern analyzes current user behavior patterns
func (injector *SafeChaosInjector) AnalyzeBehaviorPattern() *BehaviorPattern {
	return injector.userBehavior.GetCurrentPattern()
}

// AdjustDifficulty adjusts aggressiveness level based on user behavior
func (injector *SafeChaosInjector) AdjustDifficulty(pattern *BehaviorPattern) AggressivenessLevel {
	if pattern == nil || !injector.config.AdaptiveDifficulty {
		return injector.config.AggressivenessLevel
	}

	currentLevel := injector.config.AggressivenessLevel

	// Increase difficulty for expert users with high success rates
	if pattern.SkillLevel >= Advanced && pattern.RecentSuccessRate > 0.85 {
		if currentLevel < Apocalyptic {
			return currentLevel + 1
		}
	}

	// Decrease difficulty for struggling users
	if pattern.RecentSuccessRate < 0.4 || pattern.ShowsFrustration {
		if currentLevel > Off {
			return currentLevel - 1
		}
	}

	return currentLevel
}

// ValidateSafetyBoundaries validates all safety boundaries are intact
func (injector *SafeChaosInjector) ValidateSafetyBoundaries() error {
	return injector.safetyMonitor.PerformHealthCheck()
}

// GetOperationHistory returns the history of chaos injection operations
func (injector *SafeChaosInjector) GetOperationHistory() []InjectionEvent {
	injector.mutex.RLock()
	defer injector.mutex.RUnlock()

	// Return a copy to prevent external modification
	history := make([]InjectionEvent, len(injector.operationLog))
	copy(history, injector.operationLog)
	return history
}

// ResetState resets the injector state while preserving configuration
func (injector *SafeChaosInjector) ResetState() error {
	injector.mutex.Lock()
	defer injector.mutex.Unlock()

	// Reset safety monitor
	if err := injector.safetyMonitor.Reset(); err != nil {
		return fmt.Errorf("failed to reset safety monitor: %w", err)
	}

	// Reset behavior tracker
	injector.userBehavior.Reset()

	// Clear operation log
	injector.operationLog = make([]InjectionEvent, 0)

	// Reset metrics
	injector.metrics = &InjectionMetrics{}

	// Reset start time
	injector.startTime = time.Now()

	return nil
}

// initializeDefaultScenarios creates the default set of chaos scenarios
func initializeDefaultScenarios() map[string]*ChaosScenario {
	scenarios := make(map[string]*ChaosScenario)

	// Network failure scenario
	networkScenario := &ChaosScenario{
		ErrorScenario: &ErrorScenario{
			Type:    "network_failure",
			Message: "Network connection timed out",
		},
		TriggerProbability: 0.3,
		UserSkillModifier: map[SkillLevel]float64{
			Novice:       0.5,
			Intermediate: 0.7,
			Advanced:     1.0,
			Expert:       1.2,
		},
		LearningObjectives: []string{
			"Understanding network troubleshooting",
			"Learning timeout handling",
		},
		MinDuration: 200 * time.Millisecond,
		MaxDuration: 1 * time.Second,
	}
	scenarios["network_failure"] = networkScenario

	// Permission denied scenario
	permissionScenario := &ChaosScenario{
		ErrorScenario: &ErrorScenario{
			Type:    "permission_denied",
			Message: "Permission denied: insufficient privileges",
		},
		TriggerProbability: 0.25,
		UserSkillModifier: map[SkillLevel]float64{
			Novice:       0.8,
			Intermediate: 1.0,
			Advanced:     1.0,
			Expert:       0.8,
		},
		LearningObjectives: []string{
			"Understanding file permissions",
			"Learning privilege escalation",
		},
		MinDuration: 100 * time.Millisecond,
		MaxDuration: 500 * time.Millisecond,
	}
	scenarios["permission_denied"] = permissionScenario

	// Resource exhaustion scenario
	resourceScenario := &ChaosScenario{
		ErrorScenario: &ErrorScenario{
			Type:    "resource_exhausted",
			Message: "Insufficient disk space or memory",
		},
		TriggerProbability: 0.2,
		UserSkillModifier: map[SkillLevel]float64{
			Novice:       0.3,
			Intermediate: 0.6,
			Advanced:     1.0,
			Expert:       1.3,
		},
		LearningObjectives: []string{
			"Understanding resource management",
			"Learning capacity planning",
		},
		MinDuration: 300 * time.Millisecond,
		MaxDuration: 2 * time.Second,
	}
	scenarios["resource_exhausted"] = resourceScenario

	return scenarios
}

// String method for ActionType
func (at ActionType) String() string {
	switch at {
	case CommandExecution:
		return "COMMAND_EXECUTION"
	case FileOperation:
		return "FILE_OPERATION"
	case ConfigurationChange:
		return "CONFIGURATION_CHANGE"
	case HelpRequest:
		return "HELP_REQUEST"
	case RetryAttempt:
		return "RETRY_ATTEMPT"
	case RecoveryAction:
		return "RECOVERY_ACTION"
	default:
		return "UNKNOWN_ACTION"
	}
}