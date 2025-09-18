package aar

import (
	"path/filepath"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/config"
	progresssim "github.com/bthompso/engx-ergonomics-poc/internal/simulation/progress"
)

// AARGenerator generates after-action reports
type AARGenerator struct {
	tracker       *progresssim.Tracker
	config        *config.UserConfiguration
	startTime     time.Time
	projectPath   string
	stepResults   []StepResult
	performanceTargets map[string]time.Duration
}

// NewAARGenerator creates a new AAR generator
func NewAARGenerator(tracker *progresssim.Tracker, config *config.UserConfiguration, startTime time.Time, projectPath string) *AARGenerator {
	return &AARGenerator{
		tracker:     tracker,
		config:      config,
		startTime:   startTime,
		projectPath: projectPath,
		stepResults: make([]StepResult, 0),
		performanceTargets: getDefaultPerformanceTargets(),
	}
}

// getDefaultPerformanceTargets returns configurable performance targets
func getDefaultPerformanceTargets() map[string]time.Duration {
	return map[string]time.Duration{
		"total_execution":     3 * time.Minute,    // Target: under 3 minutes total
		"package_installation": 60 * time.Second,  // Target: under 1 minute for packages
		"project_structure":   5 * time.Second,    // Target: under 5 seconds for structure
		"configuration":       10 * time.Second,   // Target: under 10 seconds for config
		"step_average":        15 * time.Second,   // Target: under 15 seconds per step
	}
}

// RecordStep records the result of a step for AAR generation
func (g *AARGenerator) RecordStep(stepName string, status StepStatus, duration time.Duration, errorMessage string) {
	result := StepResult{
		Name:         stepName,
		Status:       status,
		Duration:     duration,
		StartTime:    time.Now().Add(-duration),
		EndTime:      time.Now(),
		ErrorMessage: errorMessage,
	}

	g.stepResults = append(g.stepResults, result)
}

// RecordStepWithDetails records a step with additional details for verbose/debug modes
func (g *AARGenerator) RecordStepWithDetails(stepName string, status StepStatus, duration time.Duration, errorMessage, details string, subSteps []SubStepResult) {
	result := StepResult{
		Name:         stepName,
		Status:       status,
		Duration:     duration,
		StartTime:    time.Now().Add(-duration),
		EndTime:      time.Now(),
		ErrorMessage: errorMessage,
		Details:      details,
		SubSteps:     subSteps,
	}

	g.stepResults = append(g.stepResults, result)
}

// SetPerformanceTarget sets a configurable performance target
func (g *AARGenerator) SetPerformanceTarget(key string, target time.Duration) {
	if g.performanceTargets == nil {
		g.performanceTargets = make(map[string]time.Duration)
	}
	g.performanceTargets[key] = target
}

// Generate creates the complete AAR summary
func (g *AARGenerator) Generate() (*AARSummary, error) {
	endTime := time.Now()
	duration := endTime.Sub(g.startTime)

	summary := &AARSummary{
		ProjectInfo:   g.buildProjectInfo(),
		ExecutionInfo: g.buildExecutionInfo(endTime, duration),
		StepResults:   g.stepResults,
	}

	// Generate next steps
	nextStepsEngine := NewNextStepsEngine()
	summary.NextSteps = nextStepsEngine.Generate(summary)

	// Generate troubleshooting info if there were failures
	if summary.ExecutionInfo.FailedSteps > 0 {
		summary.Troubleshooting = g.generateTroubleshooting(summary)
	}

	return summary, nil
}

// buildProjectInfo creates the project information section
func (g *AARGenerator) buildProjectInfo() ProjectInfo {
	// Extract template type from config
	template := "typescript" // default
	if g.config != nil && g.config.Template.Type != "" {
		template = string(g.config.Template.Type)
	}

	// Build features map
	features := make(map[string]bool)
	if g.config != nil {
		features["hot_reload"] = g.config.DevFeatures.HotReload
		features["linting"] = g.config.DevFeatures.Linting
		features["prettier"] = g.config.DevFeatures.Prettier
		features["husky"] = g.config.DevFeatures.Husky
		features["vscode_config"] = g.config.DevFeatures.VSCodeConfig
		features["dev_tools"] = g.config.DevFeatures.DevTools
		features["unit_testing"] = g.config.Testing.UnitTesting
		features["e2e_testing"] = g.config.Testing.E2ETesting
		features["coverage"] = g.config.Testing.Coverage
		features["docker"] = g.config.ProductionSetup.Docker
		features["cicd"] = g.config.ProductionSetup.CI_CD
		features["monitoring"] = g.config.ProductionSetup.Monitoring
		features["analytics"] = g.config.ProductionSetup.Analytics
	}

	return ProjectInfo{
		Name:         g.config.ProjectName,
		Template:     template,
		DevOnly:      g.isDevOnly(),
		Features:     features,
		Directory:    g.projectPath,
		Configuration: g.config,
	}
}

// buildExecutionInfo creates the execution information section
func (g *AARGenerator) buildExecutionInfo(endTime time.Time, duration time.Duration) ExecutionInfo {
	totalSteps := len(g.stepResults)
	successSteps := 0
	failedSteps := 0
	skippedSteps := 0

	for _, step := range g.stepResults {
		switch step.Status {
		case StepStatusSuccess:
			successSteps++
		case StepStatusFailed:
			failedSteps++
		case StepStatusSkipped:
			skippedSteps++
		}
	}

	return ExecutionInfo{
		StartTime:    g.startTime,
		EndTime:      endTime,
		Duration:     duration,
		TotalSteps:   totalSteps,
		SuccessSteps: successSteps,
		FailedSteps:  failedSteps,
		SkippedSteps: skippedSteps,
		Performance:  g.buildPerformanceMetrics(duration),
	}
}

// buildPerformanceMetrics creates performance metrics
func (g *AARGenerator) buildPerformanceMetrics(totalDuration time.Duration) PerformanceMetrics {
	if len(g.stepResults) == 0 {
		return PerformanceMetrics{
			ConfigurableTargets: g.performanceTargets,
		}
	}

	var totalStepTime time.Duration
	var slowestStep, fastestStep string
	var slowestTime, fastestTime time.Duration

	// Initialize with first step
	if len(g.stepResults) > 0 {
		slowestStep = g.stepResults[0].Name
		slowestTime = g.stepResults[0].Duration
		fastestStep = g.stepResults[0].Name
		fastestTime = g.stepResults[0].Duration
	}

	// Calculate metrics
	for _, step := range g.stepResults {
		totalStepTime += step.Duration

		if step.Duration > slowestTime {
			slowestStep = step.Name
			slowestTime = step.Duration
		}

		if step.Duration < fastestTime {
			fastestStep = step.Name
			fastestTime = step.Duration
		}
	}

	averageStepTime := totalStepTime / time.Duration(len(g.stepResults))

	return PerformanceMetrics{
		AverageStepTime:    averageStepTime,
		SlowestStep:        slowestStep,
		SlowestStepTime:    slowestTime,
		FastestStep:        fastestStep,
		FastestStepTime:    fastestTime,
		ConfigurableTargets: g.performanceTargets,
	}
}

// generateTroubleshooting creates troubleshooting information for failed steps
func (g *AARGenerator) generateTroubleshooting(summary *AARSummary) *TroubleshootingInfo {
	info := &TroubleshootingInfo{
		FailedSteps:  make([]FailedStepAnalysis, 0),
		Suggestions:  make([]string, 0),
		CommonIssues: make([]string, 0),
		SupportLinks: g.getDefaultSupportLinks(),
	}

	// Analyze each failed step
	for _, step := range summary.StepResults {
		if step.Status == StepStatusFailed {
			analysis := g.analyzeFailedStep(step)
			info.FailedSteps = append(info.FailedSteps, analysis)
		}
	}

	// Add general suggestions based on failure patterns
	info.Suggestions = g.generateGeneralSuggestions(info.FailedSteps)

	return info
}

// analyzeFailedStep analyzes a failed step and provides suggestions
func (g *AARGenerator) analyzeFailedStep(step StepResult) FailedStepAnalysis {
	suggestions := make([]string, 0)
	recoverySteps := make([]string, 0)

	// Basic suggestions based on step name and error
	switch step.Name {
	case "Package Installation":
		suggestions = append(suggestions,
			"Check internet connectivity",
			"Verify npm/yarn is properly installed",
			"Clear package manager cache",
		)
		recoverySteps = append(recoverySteps,
			"Run 'npm cache clean --force' or 'yarn cache clean'",
			"Try running 'npm install' or 'yarn install' manually",
		)
	case "TypeScript Configuration":
		suggestions = append(suggestions,
			"Verify TypeScript is installed",
			"Check for conflicting tsconfig.json files",
		)
		recoverySteps = append(recoverySteps,
			"Run 'npm install typescript --save-dev'",
			"Remove existing tsconfig.json and regenerate",
		)
	default:
		suggestions = append(suggestions,
			"Check system permissions",
			"Verify disk space availability",
		)
		recoverySteps = append(recoverySteps,
			"Retry the operation",
			"Check system logs for more details",
		)
	}

	return FailedStepAnalysis{
		StepName:      step.Name,
		ErrorMessage:  step.ErrorMessage,
		Suggestions:   suggestions,
		RecoverySteps: recoverySteps,
	}
}

// generateGeneralSuggestions creates general suggestions based on failure patterns
func (g *AARGenerator) generateGeneralSuggestions(failedSteps []FailedStepAnalysis) []string {
	suggestions := make([]string, 0)

	if len(failedSteps) > 0 {
		suggestions = append(suggestions,
			"Review the error messages above for specific guidance",
			"Ensure all system requirements are met",
			"Try running the command again with elevated permissions if needed",
		)
	}

	return suggestions
}

// getDefaultSupportLinks returns default support links
func (g *AARGenerator) getDefaultSupportLinks() []SupportLink {
	return []SupportLink{
		{
			Title:       "React Documentation",
			URL:         "https://react.dev",
			Description: "Official React documentation and guides",
		},
		{
			Title:       "TypeScript Handbook",
			URL:         "https://www.typescriptlang.org/docs/",
			Description: "Complete TypeScript documentation",
		},
		{
			Title:       "Node.js Troubleshooting",
			URL:         "https://nodejs.org/en/docs/guides/debugging-getting-started/",
			Description: "Node.js debugging and troubleshooting guide",
		},
	}
}

// isDevOnly determines if this is a dev-only configuration
func (g *AARGenerator) isDevOnly() bool {
	if g.config == nil {
		return false
	}

	// Check if production features are disabled
	return !g.config.ProductionSetup.Docker &&
		   !g.config.ProductionSetup.CI_CD &&
		   !g.config.ProductionSetup.Monitoring &&
		   !g.config.ProductionSetup.Analytics
}

// GetProjectDirectory returns the full path to the created project
func (g *AARGenerator) GetProjectDirectory() string {
	if g.projectPath == "" && g.config != nil {
		// Build path from current directory + project name
		return filepath.Join(".", g.config.ProjectName)
	}
	return g.projectPath
}