package common

import (
	"fmt"
	"sync"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
)

// ChaosInjector implements the interfaces.ChaosInjector interface
type ChaosInjector struct {
	config       *ChaosConfig
	enabled      bool
	level        interfaces.AggressivenessLevel
	scenarios    map[string]*ChaosScenario
	injectionLog []InjectionEvent
	metrics      *ChaosMetrics
	logger       interfaces.Logger
	mutex        sync.RWMutex
}

// ChaosConfig represents chaos injection configuration
type ChaosConfig struct {
	Enabled             bool                            `json:"enabled"`
	AggressivenessLevel interfaces.AggressivenessLevel `json:"aggressiveness_level"`
	SafetyMode          bool                            `json:"safety_mode"`
	AdaptiveDifficulty  bool                            `json:"adaptive_difficulty"`
	MaxInjectionCount   int64                           `json:"max_injection_count"`
	AllowedOperations   []string                        `json:"allowed_operations"`
	ProhibitedPaths     []string                        `json:"prohibited_paths"`
}

// ChaosScenario represents a failure scenario
type ChaosScenario struct {
	Name               string                      `json:"name"`
	Type               string                      `json:"type"`
	Description        string                      `json:"description"`
	TriggerProbability float64                     `json:"trigger_probability"`
	Impact             interfaces.ChaosImpact      `json:"impact"`
	Duration           time.Duration               `json:"duration"`
	RecoveryHints      []string                    `json:"recovery_hints"`
	LearningObjectives []string                    `json:"learning_objectives"`
}

// InjectionEvent represents a chaos injection event
type InjectionEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Scenario  string    `json:"scenario"`
	Success   bool      `json:"success"`
	Duration  time.Duration `json:"duration"`
	Impact    interfaces.ChaosImpact `json:"impact"`
}

// ChaosMetrics tracks chaos injection statistics
type ChaosMetrics struct {
	TotalInjections      int64         `json:"total_injections"`
	SuccessfulInjections int64         `json:"successful_injections"`
	FailedInjections     int64         `json:"failed_injections"`
	AverageImpact        float64       `json:"average_impact"`
	TotalUptime          time.Duration `json:"total_uptime"`
	LastInjection        time.Time     `json:"last_injection"`
}

// NewChaosInjector creates a new chaos injector instance
func NewChaosInjector(config interfaces.ConfigManager, logger interfaces.Logger) interfaces.ChaosInjector {
	chaosConfig := &ChaosConfig{
		Enabled:             false, // Default disabled
		AggressivenessLevel: interfaces.AggressivenessOff,
		SafetyMode:          true,
		AdaptiveDifficulty:  true,
		MaxInjectionCount:   1000,
		AllowedOperations:   []string{},
		ProhibitedPaths: []string{
			"/", "/usr", "/bin", "/sbin", "/etc", "/var", "/opt", "/home",
			"C:\\", "C:\\Windows", "C:\\Program Files", "C:\\Users",
		},
	}

	injector := &ChaosInjector{
		config:       chaosConfig,
		enabled:      false,
		level:        interfaces.AggressivenessOff,
		scenarios:    make(map[string]*ChaosScenario),
		injectionLog: make([]InjectionEvent, 0),
		metrics:      &ChaosMetrics{},
		logger:       logger.WithComponent("chaos"),
		mutex:        sync.RWMutex{},
	}

	// Initialize default scenarios
	injector.initializeDefaultScenarios()

	// Only log initialization in explicit debug mode
	// injector.logger.Debug("Chaos injector initialized (disabled by default)")
	return injector
}

// Enable enables chaos injection with the specified level
func (ci *ChaosInjector) Enable(level interfaces.AggressivenessLevel) error {
	ci.mutex.Lock()
	defer ci.mutex.Unlock()

	ci.enabled = true
	ci.level = level
	ci.config.Enabled = true
	ci.config.AggressivenessLevel = level

	ci.logger.Info("Chaos injection enabled at level: %s", level.String())
	return nil
}

// Disable disables chaos injection
func (ci *ChaosInjector) Disable() error {
	ci.mutex.Lock()
	defer ci.mutex.Unlock()

	ci.enabled = false
	ci.level = interfaces.AggressivenessOff
	ci.config.Enabled = false
	ci.config.AggressivenessLevel = interfaces.AggressivenessOff

	ci.logger.Info("Chaos injection disabled")
	return nil
}

// IsEnabled returns whether chaos injection is currently enabled
func (ci *ChaosInjector) IsEnabled() bool {
	ci.mutex.RLock()
	defer ci.mutex.RUnlock()
	return ci.enabled
}

// ShouldInjectForStep determines if chaos should be injected for a specific step
func (ci *ChaosInjector) ShouldInjectForStep(stepName string, stepIndex int) bool {
	ci.mutex.RLock()
	defer ci.mutex.RUnlock()

	if !ci.enabled {
		return false
	}

	// More specific logic based on step name and index
	if !ci.isOperationAllowed(stepName) {
		return false
	}

	// Safety checks with step context
	if ci.config.SafetyMode && ci.metrics.TotalInjections >= ci.config.MaxInjectionCount {
		ci.logger.Warn("Maximum injection count reached for step %s[%d]", stepName, stepIndex)
		return false
	}

	// Calculate injection probability with step context
	baseRate := ci.getFailureRate()

	// Adjust rate based on step index (later steps might have different rates)
	adjustedRate := baseRate
	if stepIndex > 3 {
		adjustedRate *= 0.5 // Reduce probability for later steps
	}

	return adjustedRate > 0.0
}

// GetErrorTemplate returns an error template for a scenario
func (ci *ChaosInjector) GetErrorTemplate(scenario string) *interfaces.ErrorTemplate {
	ci.mutex.RLock()
	defer ci.mutex.RUnlock()

	if _, exists := ci.scenarios[scenario]; exists {
		return &interfaces.ErrorTemplate{
			// These would be properly defined based on the actual ErrorTemplate structure
		}
	}

	return nil
}

// GenerateErrorTemplate generates a new error template
func (ci *ChaosInjector) GenerateErrorTemplate(stepName string, context string) *interfaces.ErrorTemplate {
	ci.mutex.RLock()
	defer ci.mutex.RUnlock()

	// Generate a template based on step name and context
	return &interfaces.ErrorTemplate{
		// These would be properly defined based on the actual ErrorTemplate structure
	}
}

// GetConfig returns the chaos configuration
func (ci *ChaosInjector) GetConfig() *interfaces.ChaosConfig {
	ci.mutex.RLock()
	defer ci.mutex.RUnlock()

	return &interfaces.ChaosConfig{
		// These would be properly defined based on the actual ChaosConfig structure
	}
}

// GetAggressivenessLevel returns the current aggressiveness level
func (ci *ChaosInjector) GetAggressivenessLevel() interfaces.AggressivenessLevel {
	ci.mutex.RLock()
	defer ci.mutex.RUnlock()
	return ci.level
}

// TrackBehavior tracks behavior for analytics
func (ci *ChaosInjector) TrackBehavior(action string, metadata map[string]interface{}) {
	ci.mutex.Lock()
	defer ci.mutex.Unlock()

	// Record behavior tracking
	ci.logger.Debug("Tracking behavior: %s with metadata: %v", action, metadata)

	// In a real implementation, this would store behavior data
}

// GetBehaviorData returns behavior analytics data
func (ci *ChaosInjector) GetBehaviorData() map[string]interface{} {
	ci.mutex.RLock()
	defer ci.mutex.RUnlock()

	return map[string]interface{}{
		"total_injections": ci.metrics.TotalInjections,
		"success_rate":     ci.calculateSuccessRate(),
		"level":           ci.level.String(),
		"enabled":         ci.enabled,
	}
}

// GetLevel returns the current aggressiveness level
func (ci *ChaosInjector) GetLevel() interfaces.AggressivenessLevel {
	ci.mutex.RLock()
	defer ci.mutex.RUnlock()
	return ci.level
}

// SetLevel sets the aggressiveness level
func (ci *ChaosInjector) SetLevel(level interfaces.AggressivenessLevel) error {
	ci.mutex.Lock()
	defer ci.mutex.Unlock()

	ci.level = level
	ci.config.AggressivenessLevel = level

	if level == interfaces.AggressivenessOff {
		ci.enabled = false
		ci.config.Enabled = false
	} else if !ci.enabled {
		ci.enabled = true
		ci.config.Enabled = true
	}

	ci.logger.Debug("Chaos injection level set to: %s", level.String())
	return nil
}

// ShouldInject determines if chaos should be injected for the given operation
func (ci *ChaosInjector) ShouldInject(operation string) bool {
	ci.mutex.RLock()
	defer ci.mutex.RUnlock()

	if !ci.enabled {
		return false
	}

	// Check if operation is allowed
	if !ci.isOperationAllowed(operation) {
		return false
	}

	// Safety checks
	if ci.config.SafetyMode && ci.metrics.TotalInjections >= ci.config.MaxInjectionCount {
		ci.logger.Warn("Maximum injection count reached, skipping injection")
		return false
	}

	// Calculate injection probability based on aggressiveness level
	baseRate := ci.getFailureRate()

	// For demonstration purposes, use simplified probability
	// In a real implementation, this would be more sophisticated
	return baseRate > 0.0
}

// InjectFailure injects a failure for the specified operation
func (ci *ChaosInjector) InjectFailure(operation string) error {
	startTime := time.Now()

	ci.mutex.Lock()
	defer ci.mutex.Unlock()

	if !ci.enabled {
		return fmt.Errorf("chaos injection is disabled")
	}

	// Select a scenario for this operation
	scenario := ci.selectScenario(operation)
	if scenario == nil {
		return fmt.Errorf("no applicable scenario found for operation: %s", operation)
	}

	// Simulate the failure
	err := ci.simulateFailure(scenario)

	// Record the injection event
	event := InjectionEvent{
		Timestamp: startTime,
		Operation: operation,
		Scenario:  scenario.Name,
		Success:   err == nil,
		Duration:  time.Since(startTime),
		Impact:    scenario.Impact,
	}

	ci.injectionLog = append(ci.injectionLog, event)
	ci.metrics.TotalInjections++

	if err == nil {
		ci.metrics.SuccessfulInjections++
	} else {
		ci.metrics.FailedInjections++
	}
	ci.metrics.LastInjection = startTime

	ci.logger.Debug("Chaos injection executed for %s: scenario=%s, success=%t",
		operation, scenario.Name, err == nil)

	if err == nil {
		// Return a chaos-specific error to indicate injection
		return fmt.Errorf("CHAOS INJECTION: %s - %s", scenario.Type, scenario.Description)
	}

	return err
}

// GetScenarios returns all available chaos scenarios
func (ci *ChaosInjector) GetScenarios() []map[string]interface{} {
	ci.mutex.RLock()
	defer ci.mutex.RUnlock()

	var scenarios []map[string]interface{}
	for _, scenario := range ci.scenarios {
		scenarios = append(scenarios, map[string]interface{}{
			"name":                scenario.Name,
			"type":                scenario.Type,
			"description":         scenario.Description,
			"trigger_probability": scenario.TriggerProbability,
			"impact":              scenario.Impact.String(),
			"learning_objectives": scenario.LearningObjectives,
		})
	}

	return scenarios
}

// GetMetrics returns chaos injection metrics
func (ci *ChaosInjector) GetMetrics() map[string]interface{} {
	ci.mutex.RLock()
	defer ci.mutex.RUnlock()

	return map[string]interface{}{
		"total_injections":      ci.metrics.TotalInjections,
		"successful_injections": ci.metrics.SuccessfulInjections,
		"failed_injections":     ci.metrics.FailedInjections,
		"success_rate":          ci.calculateSuccessRate(),
		"last_injection":        ci.metrics.LastInjection,
		"uptime":                time.Since(ci.metrics.LastInjection),
	}
}

// GetHistory returns the injection event history
func (ci *ChaosInjector) GetHistory() []map[string]interface{} {
	ci.mutex.RLock()
	defer ci.mutex.RUnlock()

	var history []map[string]interface{}
	for _, event := range ci.injectionLog {
		history = append(history, map[string]interface{}{
			"timestamp": event.Timestamp,
			"operation": event.Operation,
			"scenario":  event.Scenario,
			"success":   event.Success,
			"duration":  event.Duration,
			"impact":    event.Impact.String(),
		})
	}

	return history
}

// Reset resets chaos injection state
func (ci *ChaosInjector) Reset() error {
	ci.mutex.Lock()
	defer ci.mutex.Unlock()

	ci.injectionLog = make([]InjectionEvent, 0)
	ci.metrics = &ChaosMetrics{}

	ci.logger.Debug("Chaos injection state reset")
	return nil
}

// isOperationAllowed checks if an operation is allowed for chaos injection
func (ci *ChaosInjector) isOperationAllowed(operation string) bool {
	// If no allowed operations specified, all are allowed
	if len(ci.config.AllowedOperations) == 0 {
		return true
	}

	// Check if operation is in allowed list
	for _, allowed := range ci.config.AllowedOperations {
		if operation == allowed {
			return true
		}
	}

	return false
}

// getFailureRate returns the failure rate for the current aggressiveness level
func (ci *ChaosInjector) getFailureRate() float64 {
	switch ci.level {
	case interfaces.AggressivenessOff:
		return 0.0
	case interfaces.AggressivenessDefault:
		return 0.001 // 0.1%
	case interfaces.AggressivenessScout:
		return 0.005 // 0.5%
	case interfaces.AggressivenessAggressive:
		return 0.01 // 1%
	case interfaces.AggressivenessInvasive:
		return 0.05 // 5%
	case interfaces.AggressivenessApocalyptic:
		return 0.10 // 10%
	default:
		return 0.0
	}
}

// selectScenario selects an appropriate scenario for the operation
func (ci *ChaosInjector) selectScenario(operation string) *ChaosScenario {
	// For now, return a simple scenario based on operation type
	// In a real implementation, this would be more sophisticated

	for _, scenario := range ci.scenarios {
		// Simple matching - could be more sophisticated
		if scenario.Type == "general" {
			return scenario
		}
	}

	// Fallback to first available scenario
	for _, scenario := range ci.scenarios {
		return scenario
	}

	return nil
}

// simulateFailure simulates a failure scenario
func (ci *ChaosInjector) simulateFailure(scenario *ChaosScenario) error {
	// Simulate processing time
	if scenario.Duration > 0 {
		time.Sleep(scenario.Duration)
	} else {
		// Default short delay
		time.Sleep(100 * time.Millisecond)
	}

	// Simulation successful (no actual failure injection in this POC)
	return nil
}

// calculateSuccessRate calculates the success rate of injections
func (ci *ChaosInjector) calculateSuccessRate() float64 {
	if ci.metrics.TotalInjections == 0 {
		return 0.0
	}
	return float64(ci.metrics.SuccessfulInjections) / float64(ci.metrics.TotalInjections)
}

// initializeDefaultScenarios creates default chaos scenarios
func (ci *ChaosInjector) initializeDefaultScenarios() {
	scenarios := map[string]*ChaosScenario{
		"network-timeout": {
			Name:               "network-timeout",
			Type:               "network",
			Description:        "Network connection timeout during operation",
			TriggerProbability: 0.3,
			Impact:             interfaces.ChaosImpactModerate,
			Duration:           200 * time.Millisecond,
			RecoveryHints:      []string{"Check network connectivity", "Retry the operation", "Use offline mode"},
			LearningObjectives: []string{"Network troubleshooting", "Timeout handling"},
		},
		"permission-denied": {
			Name:               "permission-denied",
			Type:               "security",
			Description:        "Permission denied when accessing resources",
			TriggerProbability: 0.25,
			Impact:             interfaces.ChaosImpactLow,
			Duration:           100 * time.Millisecond,
			RecoveryHints:      []string{"Check file permissions", "Run with elevated privileges", "Use different directory"},
			LearningObjectives: []string{"File system permissions", "Security models"},
		},
		"resource-exhaustion": {
			Name:               "resource-exhaustion",
			Type:               "system",
			Description:        "System resources (memory/disk) exhausted",
			TriggerProbability: 0.2,
			Impact:             interfaces.ChaosImpactHigh,
			Duration:           500 * time.Millisecond,
			RecoveryHints:      []string{"Free up disk space", "Close other applications", "Use smaller dataset"},
			LearningObjectives: []string{"Resource management", "System monitoring"},
		},
		"general-failure": {
			Name:               "general-failure",
			Type:               "general",
			Description:        "General operation failure for testing resilience",
			TriggerProbability: 0.1,
			Impact:             interfaces.ChaosImpactModerate,
			Duration:           150 * time.Millisecond,
			RecoveryHints:      []string{"Retry the operation", "Check system logs", "Verify configuration"},
			LearningObjectives: []string{"Error handling", "Resilience patterns"},
		},
	}

	for name, scenario := range scenarios {
		ci.scenarios[name] = scenario
	}
}

// Cleanup performs cleanup of chaos injector resources
func (ci *ChaosInjector) Cleanup() error {
	ci.mutex.Lock()
	defer ci.mutex.Unlock()

	// Disable chaos injection
	ci.enabled = false
	ci.level = interfaces.AggressivenessOff

	// Clear scenarios and logs
	ci.scenarios = make(map[string]*ChaosScenario)
	ci.injectionLog = make([]InjectionEvent, 0)
	ci.metrics = &ChaosMetrics{}

	ci.logger.Debug("Chaos injector cleaned up")
	return nil
}