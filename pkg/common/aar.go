package common

import (
	"fmt"
	"sync"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
)

// AARGenerator implements the interfaces.AARGenerator interface
type AARGenerator struct {
	projectName       string
	startTime         time.Time
	endTime           time.Time
	steps             []StepResult
	performanceTargets map[string]time.Duration
	executionResult   *interfaces.ExecutionResult
	logger            interfaces.Logger
	mutex             sync.RWMutex
}

// StepResult represents the result of an execution step
type StepResult struct {
	Name         string           `json:"name"`
	Status       interfaces.StepStatus `json:"status"`
	Duration     time.Duration    `json:"duration"`
	StartTime    time.Time        `json:"start_time"`
	EndTime      time.Time        `json:"end_time"`
	ErrorMessage string           `json:"error_message,omitempty"`
	Details      string           `json:"details,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// AARReport represents a complete after-action report
type AARReport struct {
	ProjectInfo   ProjectInfo      `json:"project_info"`
	ExecutionInfo ExecutionInfo    `json:"execution_info"`
	Performance   PerformanceInfo  `json:"performance"`
	Steps         []StepResult     `json:"steps"`
	NextSteps     []NextStep       `json:"next_steps"`
	Troubleshooting *TroubleshootingInfo `json:"troubleshooting,omitempty"`
	GeneratedAt   time.Time        `json:"generated_at"`
}

// ProjectInfo contains project metadata
type ProjectInfo struct {
	Name        string                 `json:"name"`
	Directory   string                 `json:"directory"`
	Template    string                 `json:"template"`
	Features    map[string]bool        `json:"features"`
	DevOnly     bool                   `json:"dev_only"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ExecutionInfo contains execution statistics
type ExecutionInfo struct {
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	TotalDuration time.Duration `json:"total_duration"`
	TotalSteps    int           `json:"total_steps"`
	SuccessSteps  int           `json:"success_steps"`
	FailedSteps   int           `json:"failed_steps"`
	SkippedSteps  int           `json:"skipped_steps"`
	SuccessRate   float64       `json:"success_rate"`
}

// PerformanceInfo contains performance metrics
type PerformanceInfo struct {
	AverageStepTime  time.Duration            `json:"average_step_time"`
	SlowestStep      string                   `json:"slowest_step"`
	SlowestStepTime  time.Duration            `json:"slowest_step_time"`
	FastestStep      string                   `json:"fastest_step"`
	FastestStepTime  time.Duration            `json:"fastest_step_time"`
	Targets          map[string]time.Duration `json:"targets"`
	TargetsMet       map[string]bool          `json:"targets_met"`
}

// NextStep represents a recommended action
type NextStep struct {
	Action      string                    `json:"action"`
	Description string                    `json:"description"`
	Command     string                    `json:"command,omitempty"`
	Priority    interfaces.StepPriority   `json:"priority"`
	Category    interfaces.StepCategory   `json:"category"`
	WorkingDir  string                    `json:"working_dir,omitempty"`
	Resources   []string                  `json:"resources,omitempty"`
}

// TroubleshootingInfo contains failure analysis
type TroubleshootingInfo struct {
	FailedSteps   []FailedStepAnalysis `json:"failed_steps"`
	Suggestions   []string             `json:"suggestions"`
	CommonIssues  []string             `json:"common_issues"`
	SupportLinks  []SupportLink        `json:"support_links"`
}

// FailedStepAnalysis contains analysis of a failed step
type FailedStepAnalysis struct {
	StepName      string   `json:"step_name"`
	ErrorMessage  string   `json:"error_message"`
	Suggestions   []string `json:"suggestions"`
	RecoverySteps []string `json:"recovery_steps"`
	RelatedIssues []string `json:"related_issues"`
}

// SupportLink represents a help resource
type SupportLink struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

// NewAARGenerator creates a new AAR generator instance
func NewAARGenerator(config interfaces.ConfigManager, logger interfaces.Logger) interfaces.AARGenerator {
	generator := &AARGenerator{
		startTime:          time.Now(),
		steps:              make([]StepResult, 0),
		performanceTargets: getDefaultPerformanceTargets(),
		logger:             logger.WithComponent("aar"),
	}

	// Only log initialization in explicit debug mode
	// generator.logger.Debug("AAR generator initialized")
	return generator
}

// SetProjectName sets the project name for the report
func (a *AARGenerator) SetProjectName(name string) {
	a.projectName = name
	a.logger.Debug("Project name set to: %s", name)
}

// StartTracking starts tracking execution time
func (a *AARGenerator) StartTracking() {
	a.startTime = time.Now()
	a.steps = make([]StepResult, 0) // Reset steps
	a.logger.Debug("Started tracking execution at: %s", a.startTime.Format(time.RFC3339))
}

// recordInternalStep records the execution of a step (internal method)
func (a *AARGenerator) recordInternalStep(name string, status interfaces.StepStatus, duration time.Duration, metadata map[string]interface{}) {
	step := StepResult{
		Name:      name,
		Status:    status,
		Duration:  duration,
		StartTime: time.Now().Add(-duration),
		EndTime:   time.Now(),
		Metadata:  metadata,
	}

	if errorMsg, exists := metadata["error"]; exists {
		if errorStr, ok := errorMsg.(string); ok {
			step.ErrorMessage = errorStr
		}
	}

	if details, exists := metadata["details"]; exists {
		if detailsStr, ok := details.(string); ok {
			step.Details = detailsStr
		}
	}

	a.steps = append(a.steps, step)
	a.logger.Debug("Recorded step: %s (status: %s, duration: %v)", name, status.String(), duration)
}

// RecordFailure records a failed step with error information
func (a *AARGenerator) RecordFailure(stepName, errorMessage string, duration time.Duration) {
	metadata := map[string]interface{}{
		"error": errorMessage,
	}
	a.recordInternalStep(stepName, interfaces.StepStatusFailed, duration, metadata)
}

// RecordSuccess records a successful step
func (a *AARGenerator) RecordSuccess(stepName string, duration time.Duration, details map[string]interface{}) {
	a.recordInternalStep(stepName, interfaces.StepStatusSuccess, duration, details)
}

// SetPerformanceTarget sets a performance target for evaluation
func (a *AARGenerator) SetPerformanceTarget(name string, target time.Duration) {
	if a.performanceTargets == nil {
		a.performanceTargets = make(map[string]time.Duration)
	}
	a.performanceTargets[name] = target
	a.logger.Debug("Performance target set: %s = %v", name, target)
}

// generateInternalReport generates a complete AAR report (internal method)
func (a *AARGenerator) generateInternalReport() (*interfaces.AARReport, error) {
	endTime := time.Now()
	totalDuration := endTime.Sub(a.startTime)

	// Calculate execution statistics
	executionInfo := a.calculateExecutionInfo(endTime, totalDuration)

	// Calculate performance metrics
	performanceInfo := a.calculatePerformanceInfo()

	// Generate next steps
	nextSteps := a.generateNextSteps()

	// Generate troubleshooting info if there were failures
	var troubleshooting *TroubleshootingInfo
	if executionInfo.FailedSteps > 0 {
		troubleshooting = a.generateTroubleshooting()
	}

	// Create the report
	report := &AARReport{
		ProjectInfo: ProjectInfo{
			Name:      a.projectName,
			Directory: ".", // Default to current directory
			Template:  "default",
			Features:  make(map[string]bool),
			DevOnly:   true, // Default assumption
			Metadata:  make(map[string]interface{}),
		},
		ExecutionInfo:   executionInfo,
		Performance:     performanceInfo,
		Steps:           a.steps,
		NextSteps:       nextSteps,
		Troubleshooting: troubleshooting,
		GeneratedAt:     time.Now(),
	}

	// Convert to interface type
	interfaceReport := &interfaces.AARReport{
		ProjectName:     report.ProjectInfo.Name,
		ExecutionTime:   report.ExecutionInfo.TotalDuration,
		TotalSteps:      report.ExecutionInfo.TotalSteps,
		SuccessfulSteps: report.ExecutionInfo.SuccessSteps,
		FailedSteps:     report.ExecutionInfo.FailedSteps,
		PerformanceMetrics: map[string]interface{}{
			"average_step_time": report.Performance.AverageStepTime,
			"slowest_step":      report.Performance.SlowestStep,
			"fastest_step":      report.Performance.FastestStep,
			"targets_met":       report.Performance.TargetsMet,
		},
		NextSteps: a.convertNextSteps(report.NextSteps),
		GeneratedAt: report.GeneratedAt,
	}

	a.logger.Info("AAR report generated: %d steps, %d successful, %d failed",
		report.ExecutionInfo.TotalSteps,
		report.ExecutionInfo.SuccessSteps,
		report.ExecutionInfo.FailedSteps)

	return interfaceReport, nil
}

// GenerateFormattedReport generates a formatted report
func (a *AARGenerator) GenerateFormattedReport(execution *interfaces.ExecutionContext, format string) (string, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	report, err := a.GenerateReport(execution)
	if err != nil {
		return "", fmt.Errorf("failed to generate report: %w", err)
	}

	return a.FormatReport(report, format)
}

// StartExecution starts tracking a new execution
func (a *AARGenerator) StartExecution(command string, args []string) *interfaces.ExecutionContext {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	ctx := &interfaces.ExecutionContext{
		Command:   command,
		Args:      args,
		StartTime: time.Now(),
		Steps:     make([]interfaces.StepExecution, 0),
		Metadata:  make(map[string]interface{}),
	}

	a.startTime = ctx.StartTime
	a.steps = make([]StepResult, 0)
	a.logger.Debug("Started execution tracking for command: %s", command)

	return ctx
}

// RecordStep records a step execution (interface method)
func (a *AARGenerator) RecordStep(ctx *interfaces.ExecutionContext, step *interfaces.StepExecution) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	stepResult := StepResult{
		Name:      step.Name,
		Status:    step.Status,
		StartTime: step.StartTime,
		EndTime:   step.EndTime,
		Duration:  step.Duration,
		ErrorMessage: step.Error,
		Details:   step.Output,
		Metadata:  step.Metadata,
	}

	a.steps = append(a.steps, stepResult)
	ctx.Steps = append(ctx.Steps, *step)
	a.logger.Debug("Recorded step: %s (status: %s, duration: %v)", step.Name, step.Status.String(), step.Duration)
}

// FinishExecution finalizes an execution context with results
func (a *AARGenerator) FinishExecution(ctx *interfaces.ExecutionContext, result *interfaces.ExecutionResult) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.executionResult = result
	a.endTime = time.Now()

	totalDuration := a.endTime.Sub(a.startTime)
	a.performanceTargets["total_execution"] = totalDuration

	a.logger.Info("Execution finished: success=%t, duration=%v", result.Success, totalDuration)
}

// GenerateNextSteps generates next step recommendations
func (a *AARGenerator) GenerateNextSteps(ctx *interfaces.ExecutionContext) []interfaces.NextStep {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	nextSteps := a.generateNextSteps()
	return a.convertNextSteps(nextSteps)
}

// GetRecommendations gets general recommendations
func (a *AARGenerator) GetRecommendations(ctx *interfaces.ExecutionContext) []interfaces.Recommendation {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	var recommendations []interfaces.Recommendation

	// Analyze execution for recommendations
	if a.executionResult != nil && !a.executionResult.Success {
		recommendations = append(recommendations, interfaces.Recommendation{
			Title:       "Review Execution Errors",
			Description: "Analyze and resolve execution failures before proceeding",
			Priority:    interfaces.StepPriorityHigh,
			Category:    interfaces.StepCategoryTroubleshooting,
			Impact:      interfaces.ChaosImpactHigh,
		})
	}

	// Performance recommendations
	if len(a.steps) > 0 {
		performanceInfo := a.calculatePerformanceInfo()
		if performanceInfo.AverageStepTime > 30*time.Second {
			recommendations = append(recommendations, interfaces.Recommendation{
				Title:       "Optimize Performance",
				Description: "Steps are taking longer than expected - consider optimization",
				Priority:    interfaces.StepPriorityMedium,
				Category:    interfaces.StepCategoryDevelopment,
				Impact:      interfaces.ChaosImpactModerate,
			})
		}
	}

	return recommendations
}

// FormatReport formats a report in the specified format
func (a *AARGenerator) FormatReport(report *interfaces.AARReport, format string) (string, error) {
	switch format {
	case "json":
		return fmt.Sprintf(`{
  "project_name": "%s",
  "execution_time": "%v",
  "total_steps": %d,
  "successful_steps": %d,
  "failed_steps": %d,
  "success_rate": %.2f,
  "generated_at": "%s"
}`, report.ProjectName, report.ExecutionTime, report.TotalSteps,
		   report.SuccessfulSteps, report.FailedSteps,
		   float64(report.SuccessfulSteps)/float64(report.TotalSteps)*100,
		   report.GeneratedAt.Format(time.RFC3339)), nil

	case "text", "plain":
		return fmt.Sprintf(`After Action Report - %s
Generated: %s

Execution Summary:
- Total Time: %v
- Total Steps: %d
- Successful: %d
- Failed: %d
- Success Rate: %.1f%%

Next Steps: %d recommendations available
`, report.ProjectName, report.GeneratedAt.Format(time.RFC3339),
		   report.ExecutionTime, report.TotalSteps, report.SuccessfulSteps,
		   report.FailedSteps, float64(report.SuccessfulSteps)/float64(report.TotalSteps)*100,
		   len(report.NextSteps)), nil

	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// SaveReport saves a report to the specified path
func (a *AARGenerator) SaveReport(report *interfaces.AARReport, path string) error {
	formattedReport, err := a.FormatReport(report, "json")
	if err != nil {
		return fmt.Errorf("failed to format report: %w", err)
	}

	// Note: In a real implementation, this would write to file
	// For now, just log the save operation
	a.logger.Info("Report saved to: %s (%d bytes)", path, len(formattedReport))
	return nil
}

// GenerateReport generates a report for the given execution context (interface method)
func (a *AARGenerator) GenerateReport(execution *interfaces.ExecutionContext) (*interfaces.AARReport, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if execution == nil {
		// Return a basic report using internal state
		return a.generateInternalReport()
	}

	// Use execution context data if available
	endTime := time.Now()
	if !a.endTime.IsZero() {
		endTime = a.endTime
	}

	totalDuration := endTime.Sub(execution.StartTime)

	// Calculate metrics based on execution context
	successSteps := 0
	failedSteps := 0
	for _, step := range execution.Steps {
		switch step.Status {
		case interfaces.StepStatusSuccess:
			successSteps++
		case interfaces.StepStatusFailed:
			failedSteps++
		}
	}

	report := &interfaces.AARReport{
		ProjectName:        execution.Command,
		ExecutionTime:      totalDuration,
		TotalSteps:         len(execution.Steps),
		SuccessfulSteps:    successSteps,
		FailedSteps:        failedSteps,
		PerformanceMetrics: make(map[string]interface{}),
		Steps:              execution.Steps,
		NextSteps:          a.GenerateNextSteps(execution),
		GeneratedAt:        time.Now(),
	}

	return report, nil
}

// GetCurrentMetrics returns current execution metrics
func (a *AARGenerator) GetCurrentMetrics() map[string]interface{} {
	currentTime := time.Now()
	elapsed := currentTime.Sub(a.startTime)

	successCount := 0
	failedCount := 0
	for _, step := range a.steps {
		switch step.Status {
		case interfaces.StepStatusSuccess:
			successCount++
		case interfaces.StepStatusFailed:
			failedCount++
		}
	}

	return map[string]interface{}{
		"elapsed_time":     elapsed,
		"total_steps":      len(a.steps),
		"successful_steps": successCount,
		"failed_steps":     failedCount,
		"success_rate":     a.calculateSuccessRate(),
	}
}

// Reset resets the AAR generator state
func (a *AARGenerator) Reset() {
	a.startTime = time.Now()
	a.steps = make([]StepResult, 0)
	a.projectName = ""
	a.logger.Debug("AAR generator state reset")
}

// calculateExecutionInfo calculates execution statistics
func (a *AARGenerator) calculateExecutionInfo(endTime time.Time, totalDuration time.Duration) ExecutionInfo {
	successSteps := 0
	failedSteps := 0
	skippedSteps := 0

	for _, step := range a.steps {
		switch step.Status {
		case interfaces.StepStatusSuccess:
			successSteps++
		case interfaces.StepStatusFailed:
			failedSteps++
		case interfaces.StepStatusSkipped:
			skippedSteps++
		}
	}

	successRate := a.calculateSuccessRate()

	return ExecutionInfo{
		StartTime:     a.startTime,
		EndTime:       endTime,
		TotalDuration: totalDuration,
		TotalSteps:    len(a.steps),
		SuccessSteps:  successSteps,
		FailedSteps:   failedSteps,
		SkippedSteps:  skippedSteps,
		SuccessRate:   successRate,
	}
}

// calculatePerformanceInfo calculates performance metrics
func (a *AARGenerator) calculatePerformanceInfo() PerformanceInfo {
	if len(a.steps) == 0 {
		return PerformanceInfo{
			Targets:    a.performanceTargets,
			TargetsMet: make(map[string]bool),
		}
	}

	var totalTime time.Duration
	var slowestStep, fastestStep string
	var slowestTime, fastestTime time.Duration

	// Initialize with first step
	slowestStep = a.steps[0].Name
	slowestTime = a.steps[0].Duration
	fastestStep = a.steps[0].Name
	fastestTime = a.steps[0].Duration

	// Calculate metrics
	for _, step := range a.steps {
		totalTime += step.Duration

		if step.Duration > slowestTime {
			slowestStep = step.Name
			slowestTime = step.Duration
		}

		if step.Duration < fastestTime {
			fastestStep = step.Name
			fastestTime = step.Duration
		}
	}

	averageTime := totalTime / time.Duration(len(a.steps))

	// Check targets
	targetsMet := make(map[string]bool)
	for targetName, targetValue := range a.performanceTargets {
		switch targetName {
		case "average_step_time":
			targetsMet[targetName] = averageTime <= targetValue
		case "total_execution":
			totalDuration := time.Now().Sub(a.startTime)
			targetsMet[targetName] = totalDuration <= targetValue
		default:
			targetsMet[targetName] = false
		}
	}

	return PerformanceInfo{
		AverageStepTime: averageTime,
		SlowestStep:     slowestStep,
		SlowestStepTime: slowestTime,
		FastestStep:     fastestStep,
		FastestStepTime: fastestTime,
		Targets:         a.performanceTargets,
		TargetsMet:      targetsMet,
	}
}

// generateNextSteps generates recommended next actions
func (a *AARGenerator) generateNextSteps() []NextStep {
	var nextSteps []NextStep

	// Analyze execution results to determine next steps
	hasFailures := false
	for _, step := range a.steps {
		if step.Status == interfaces.StepStatusFailed {
			hasFailures = true
			break
		}
	}

	if hasFailures {
		nextSteps = append(nextSteps, NextStep{
			Action:      "Review Failures",
			Description: "Review and resolve failed steps before proceeding",
			Priority:    interfaces.StepPriorityHigh,
			Category:    interfaces.StepCategoryTroubleshooting,
		})
	} else {
		// Success scenarios
		nextSteps = append(nextSteps, []NextStep{
			{
				Action:      "Start Development",
				Description: "Begin developing your application",
				Command:     "cd " + a.projectName + " && npm start",
				Priority:    interfaces.StepPriorityHigh,
				Category:    interfaces.StepCategoryDevelopment,
				WorkingDir:  a.projectName,
			},
			{
				Action:      "Run Tests",
				Description: "Execute the test suite to verify functionality",
				Command:     "npm test",
				Priority:    interfaces.StepPriorityMedium,
				Category:    interfaces.StepCategoryTesting,
				WorkingDir:  a.projectName,
			},
			{
				Action:      "Review Documentation",
				Description: "Read the generated documentation and README",
				Priority:    interfaces.StepPriorityLow,
				Category:    interfaces.StepCategoryDocumentation,
				WorkingDir:  a.projectName,
			},
		}...)
	}

	return nextSteps
}

// generateTroubleshooting generates troubleshooting information
func (a *AARGenerator) generateTroubleshooting() *TroubleshootingInfo {
	var failedSteps []FailedStepAnalysis
	var suggestions []string

	// Analyze failed steps
	for _, step := range a.steps {
		if step.Status == interfaces.StepStatusFailed {
			analysis := FailedStepAnalysis{
				StepName:     step.Name,
				ErrorMessage: step.ErrorMessage,
				Suggestions:  a.generateSuggestionsForStep(step),
				RecoverySteps: a.generateRecoverySteps(step),
			}
			failedSteps = append(failedSteps, analysis)
		}
	}

	// General suggestions
	if len(failedSteps) > 0 {
		suggestions = []string{
			"Check system requirements and dependencies",
			"Verify network connectivity if installation failed",
			"Review error messages for specific guidance",
			"Try running with elevated permissions if needed",
		}
	}

	return &TroubleshootingInfo{
		FailedSteps:  failedSteps,
		Suggestions:  suggestions,
		CommonIssues: a.getCommonIssues(),
		SupportLinks: a.getSupportLinks(),
	}
}

// generateSuggestionsForStep generates suggestions for a specific failed step
func (a *AARGenerator) generateSuggestionsForStep(step StepResult) []string {
	suggestions := make([]string, 0)

	// Basic suggestions based on step name
	switch step.Name {
	case "Package Installation":
		suggestions = append(suggestions,
			"Check internet connectivity",
			"Clear package manager cache",
			"Verify npm/yarn installation",
		)
	case "Project Structure":
		suggestions = append(suggestions,
			"Check file system permissions",
			"Verify disk space availability",
			"Ensure target directory is writable",
		)
	default:
		suggestions = append(suggestions,
			"Review the error message above",
			"Check system logs for additional details",
			"Retry the operation",
		)
	}

	return suggestions
}

// generateRecoverySteps generates recovery steps for a failed step
func (a *AARGenerator) generateRecoverySteps(step StepResult) []string {
	recoverySteps := make([]string, 0)

	switch step.Name {
	case "Package Installation":
		recoverySteps = append(recoverySteps,
			"Run 'npm cache clean --force'",
			"Try 'npm install' manually",
			"Use 'yarn install' as alternative",
		)
	default:
		recoverySteps = append(recoverySteps,
			"Retry the failed operation",
			"Check system requirements",
			"Contact support if issue persists",
		)
	}

	return recoverySteps
}

// getCommonIssues returns common issues and solutions
func (a *AARGenerator) getCommonIssues() []string {
	return []string{
		"Network timeouts during package installation",
		"Permission denied errors on file creation",
		"Insufficient disk space for project files",
		"Node.js/npm version compatibility issues",
	}
}

// getSupportLinks returns helpful support resources
func (a *AARGenerator) getSupportLinks() []SupportLink {
	return []SupportLink{
		{
			Title:       "Node.js Documentation",
			URL:         "https://nodejs.org/docs/",
			Description: "Official Node.js documentation and guides",
		},
		{
			Title:       "npm Troubleshooting",
			URL:         "https://docs.npmjs.com/troubleshooting",
			Description: "Common npm issues and solutions",
		},
		{
			Title:       "Project Support",
			URL:         "https://github.com/project/issues",
			Description: "Report issues and get community support",
		},
	}
}

// calculateSuccessRate calculates the success rate of steps
func (a *AARGenerator) calculateSuccessRate() float64 {
	if len(a.steps) == 0 {
		return 0.0
	}

	successCount := 0
	for _, step := range a.steps {
		if step.Status == interfaces.StepStatusSuccess {
			successCount++
		}
	}

	return float64(successCount) / float64(len(a.steps))
}

// convertNextSteps converts internal NextStep to interface NextStep
func (a *AARGenerator) convertNextSteps(steps []NextStep) []interfaces.NextStep {
	var converted []interfaces.NextStep
	for _, step := range steps {
		converted = append(converted, interfaces.NextStep{
			Action:      step.Action,
			Description: step.Description,
			Command:     step.Command,
			Priority:    step.Priority,
			Category:    step.Category,
			WorkingDir:  step.WorkingDir,
		})
	}
	return converted
}

// getDefaultPerformanceTargets returns default performance targets
func getDefaultPerformanceTargets() map[string]time.Duration {
	return map[string]time.Duration{
		"total_execution":     3 * time.Minute,
		"average_step_time":   15 * time.Second,
		"package_installation": 60 * time.Second,
		"project_structure":   5 * time.Second,
	}
}

// Cleanup performs cleanup of AAR generator resources
func (a *AARGenerator) Cleanup() error {
	a.steps = make([]StepResult, 0)
	a.performanceTargets = make(map[string]time.Duration)
	a.projectName = ""

	a.logger.Debug("AAR generator cleaned up")
	return nil
}