package aar

import (
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/config"
)

// AARSummary represents the complete after-action report
type AARSummary struct {
	ProjectInfo    ProjectInfo      `json:"project_info"`
	ExecutionInfo  ExecutionInfo    `json:"execution_info"`
	StepResults    []StepResult     `json:"step_results"`
	NextSteps      []NextStep       `json:"next_steps"`
	Troubleshooting *TroubleshootingInfo `json:"troubleshooting,omitempty"`
}

// ProjectInfo contains information about the created project
type ProjectInfo struct {
	Name         string                        `json:"name"`
	Template     string                        `json:"template"`
	DevOnly      bool                          `json:"dev_only"`
	Features     map[string]bool               `json:"features"`
	Directory    string                        `json:"directory"`
	Configuration *config.UserConfiguration    `json:"configuration,omitempty"`
}

// ExecutionInfo contains timing and performance data
type ExecutionInfo struct {
	StartTime     time.Time         `json:"start_time"`
	EndTime       time.Time         `json:"end_time"`
	Duration      time.Duration     `json:"duration"`
	TotalSteps    int               `json:"total_steps"`
	SuccessSteps  int               `json:"success_steps"`
	FailedSteps   int               `json:"failed_steps"`
	SkippedSteps  int               `json:"skipped_steps"`
	Performance   PerformanceMetrics `json:"performance"`
}

// PerformanceMetrics contains performance and timing data
type PerformanceMetrics struct {
	AverageStepTime   time.Duration `json:"average_step_time"`
	SlowestStep       string        `json:"slowest_step"`
	SlowestStepTime   time.Duration `json:"slowest_step_time"`
	FastestStep       string        `json:"fastest_step"`
	FastestStepTime   time.Duration `json:"fastest_step_time"`
	ConfigurableTargets map[string]time.Duration `json:"configurable_targets"`
}

// StepResult contains the result of an individual step
type StepResult struct {
	Name         string        `json:"name"`
	Status       StepStatus    `json:"status"`
	Duration     time.Duration `json:"duration"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	ErrorMessage string        `json:"error_message,omitempty"`
	Details      string        `json:"details,omitempty"`
	SubSteps     []SubStepResult `json:"sub_steps,omitempty"`
}

// SubStepResult contains detailed sub-step information for verbose/debug modes
type SubStepResult struct {
	Name      string        `json:"name"`
	Status    StepStatus    `json:"status"`
	Duration  time.Duration `json:"duration"`
	Details   string        `json:"details,omitempty"`
}

// StepStatus represents the status of a step
type StepStatus int

const (
	StepStatusPending StepStatus = iota
	StepStatusRunning
	StepStatusSuccess
	StepStatusFailed
	StepStatusSkipped
)

func (s StepStatus) String() string {
	switch s {
	case StepStatusPending:
		return "pending"
	case StepStatusRunning:
		return "running"
	case StepStatusSuccess:
		return "success"
	case StepStatusFailed:
		return "failed"
	case StepStatusSkipped:
		return "skipped"
	default:
		return "unknown"
	}
}

// NextStep represents a recommended action for the user
type NextStep struct {
	Action       string        `json:"action"`
	Description  string        `json:"description"`
	Command      string        `json:"command,omitempty"`
	Priority     StepPriority  `json:"priority"`
	Category     StepCategory  `json:"category"`
	WorkingDir   string        `json:"working_dir,omitempty"`
}

// StepPriority represents the priority of a next step
type StepPriority int

const (
	PriorityLow StepPriority = iota
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

func (p StepPriority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityMedium:
		return "medium"
	case PriorityHigh:
		return "high"
	case PriorityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// StepCategory represents the category of a next step
type StepCategory int

const (
	CategoryDevelopment StepCategory = iota
	CategoryTesting
	CategoryDeployment
	CategoryConfiguration
	CategoryDocumentation
	CategoryTroubleshooting
)

func (c StepCategory) String() string {
	switch c {
	case CategoryDevelopment:
		return "development"
	case CategoryTesting:
		return "testing"
	case CategoryDeployment:
		return "deployment"
	case CategoryConfiguration:
		return "configuration"
	case CategoryDocumentation:
		return "documentation"
	case CategoryTroubleshooting:
		return "troubleshooting"
	default:
		return "unknown"
	}
}

// TroubleshootingInfo contains failure analysis and recovery suggestions
type TroubleshootingInfo struct {
	FailedSteps    []FailedStepAnalysis `json:"failed_steps"`
	Suggestions    []string             `json:"suggestions"`
	CommonIssues   []string             `json:"common_issues,omitempty"`
	SupportLinks   []SupportLink        `json:"support_links,omitempty"`
}

// FailedStepAnalysis contains analysis of a failed step
type FailedStepAnalysis struct {
	StepName       string   `json:"step_name"`
	ErrorMessage   string   `json:"error_message"`
	Suggestions    []string `json:"suggestions"`
	RecoverySteps  []string `json:"recovery_steps"`
	RelatedIssues  []string `json:"related_issues,omitempty"`
}

// SupportLink represents a link to documentation or support
type SupportLink struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}