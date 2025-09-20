package interfaces

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfigManager defines the interface for configuration management services.
// This interface abstracts configuration loading, validation, and access
// across the plugin system.
type ConfigManager interface {
	// Configuration loading and management
	LoadConfig(path string) (*Config, error)
	GetConfig() *Config
	ReloadConfig() error

	// Verbosity and global settings
	GetVerbosity() VerbosityLevel
	SetVerbosity(level VerbosityLevel) error
	GetGlobalFlags() *GlobalFlags

	// Plugin-specific configuration
	GetPluginConfig(pluginName string) (map[string]interface{}, error)
	SetPluginConfig(pluginName string, config map[string]interface{}) error

	// Validation and schema
	ValidateConfig() error
	GetConfigSchema() map[string]interface{}
}

// TUIRegistry defines the interface for TUI component management.
// This interface provides access to reusable TUI components and models
// for consistent user interface across plugins.
type TUIRegistry interface {
	// Model creation
	CreateProgressModel(steps []Step) tea.Model
	CreateConfirmModel(prompt string) tea.Model
	CreateInputModel(prompt, placeholder string) tea.Model
	CreateSelectModel(prompt string, options []string) tea.Model

	// Component access
	GetTheme() *Theme
	GetComponents() *ComponentLibrary
	GetStyles() *StyleRegistry

	// Model execution
	RunModel(model tea.Model) (tea.Model, error)
	RunProgram(model tea.Model) error
}

// ChaosInjector defines the interface for chaos engineering functionality.
// This interface abstracts the chaos marine system for use across plugins
// that want to support educational failure injection.
type ChaosInjector interface {
	// Chaos injection decision
	ShouldInject(context string) bool
	ShouldInjectForStep(stepName string, stepIndex int) bool

	// Error template management
	GetErrorTemplate(scenario string) *ErrorTemplate
	GenerateErrorTemplate(stepName string, context string) *ErrorTemplate

	// Configuration access
	GetConfig() *ChaosConfig
	GetAggressivenessLevel() AggressivenessLevel
	IsEnabled() bool

	// Tracking and behavior
	TrackBehavior(action string, metadata map[string]interface{})
	GetBehaviorData() map[string]interface{}
}

// AARGenerator defines the interface for After Action Report generation.
// This interface provides capabilities for generating structured reports
// about command execution, outcomes, and lessons learned.
type AARGenerator interface {
	// Report generation
	GenerateReport(execution *ExecutionContext) (*AARReport, error)
	GenerateFormattedReport(execution *ExecutionContext, format string) (string, error)

	// Execution tracking
	StartExecution(command string, args []string) *ExecutionContext
	RecordStep(ctx *ExecutionContext, step *StepExecution)
	FinishExecution(ctx *ExecutionContext, result *ExecutionResult)

	// Next steps and recommendations
	GenerateNextSteps(ctx *ExecutionContext) []NextStep
	GetRecommendations(ctx *ExecutionContext) []Recommendation

	// Output and formatting
	FormatReport(report *AARReport, format string) (string, error)
	SaveReport(report *AARReport, path string) error
}

// Logger defines the interface for logging functionality across the plugin system.
// This interface provides structured logging capabilities with different levels
// and contexts for debugging and monitoring.
type Logger interface {
	// Standard logging levels
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})

	// Structured logging with context
	WithContext(ctx map[string]interface{}) Logger
	WithPlugin(pluginName string) Logger
	WithComponent(componentName string) Logger

	// Configuration
	SetLevel(level LogLevel)
	GetLevel() LogLevel
}

// FilesystemManager defines the interface for filesystem operations.
// This interface abstracts file system interactions for testing and
// provides common filesystem utilities for plugins.
type FilesystemManager interface {
	// File operations
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm uint32) error
	AppendToFile(path string, data []byte) error

	// Directory operations
	CreateDir(path string, perm uint32) error
	RemoveDir(path string) error
	ListDir(path string) ([]string, error)

	// Path operations
	Exists(path string) bool
	IsDir(path string) bool
	IsFile(path string) bool
	GetWorkingDir() (string, error)
	ChangeDir(path string) error

	// Utilities
	CopyFile(src, dst string) error
	MoveFile(src, dst string) error
	GetTempDir() (string, error)
}

// Supporting types and structures

// VerbosityLevel represents different levels of output verbosity
type VerbosityLevel int

const (
	VerbosityQuiet VerbosityLevel = iota
	VerbosityConcise
	VerbosityNormal
	VerbosityVerbose
	VerbosityDebug
)

func (vl VerbosityLevel) String() string {
	switch vl {
	case VerbosityQuiet:
		return "quiet"
	case VerbosityConcise:
		return "concise"
	case VerbosityNormal:
		return "normal"
	case VerbosityVerbose:
		return "verbose"
	case VerbosityDebug:
		return "debug"
	default:
		return "unknown"
	}
}

// LogLevel represents different logging levels
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func (ll LogLevel) String() string {
	switch ll {
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	default:
		return "unknown"
	}
}

// AggressivenessLevel represents chaos engineering aggressiveness levels
type AggressivenessLevel int

const (
	AggressivenessOff AggressivenessLevel = iota
	AggressivenessDefault
	AggressivenessScout
	AggressivenessAggressive
	AggressivenessInvasive
	AggressivenessApocalyptic
)

func (al AggressivenessLevel) String() string {
	switch al {
	case AggressivenessOff:
		return "off"
	case AggressivenessDefault:
		return "default"
	case AggressivenessScout:
		return "scout"
	case AggressivenessAggressive:
		return "aggressive"
	case AggressivenessInvasive:
		return "invasive"
	case AggressivenessApocalyptic:
		return "apocalyptic"
	default:
		return "unknown"
	}
}

// ChaosImpact defines the impact level of chaos scenarios
type ChaosImpact int

const (
	ChaosImpactLow ChaosImpact = iota
	ChaosImpactModerate
	ChaosImpactHigh
	ChaosImpactCritical
)

func (c ChaosImpact) String() string {
	switch c {
	case ChaosImpactLow:
		return "low"
	case ChaosImpactModerate:
		return "moderate"
	case ChaosImpactHigh:
		return "high"
	case ChaosImpactCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// StepStatus represents the status of an execution step
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

// StepPriority represents the priority of a step
type StepPriority int

const (
	StepPriorityLow StepPriority = iota
	StepPriorityMedium
	StepPriorityHigh
	StepPriorityCritical
)

func (p StepPriority) String() string {
	switch p {
	case StepPriorityLow:
		return "low"
	case StepPriorityMedium:
		return "medium"
	case StepPriorityHigh:
		return "high"
	case StepPriorityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// StepCategory represents the category of a step
type StepCategory int

const (
	StepCategoryDevelopment StepCategory = iota
	StepCategoryTesting
	StepCategoryDeployment
	StepCategoryConfiguration
	StepCategoryDocumentation
	StepCategoryTroubleshooting
)

func (c StepCategory) String() string {
	switch c {
	case StepCategoryDevelopment:
		return "development"
	case StepCategoryTesting:
		return "testing"
	case StepCategoryDeployment:
		return "deployment"
	case StepCategoryConfiguration:
		return "configuration"
	case StepCategoryDocumentation:
		return "documentation"
	case StepCategoryTroubleshooting:
		return "troubleshooting"
	default:
		return "unknown"
	}
}

// Placeholder types - these will be defined as we migrate existing code
type Config struct {
	Project      *ProjectConfig              `json:"project,omitempty"`
	Defaults     *DefaultsConfig             `json:"defaults,omitempty"`
	Environments map[string]*EnvConfig       `json:"environments,omitempty"`
	Commands     map[string]*CmdConfig       `json:"commands,omitempty"`
}

type ProjectConfig struct {
	Name        string `json:"name,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
}

type DefaultsConfig struct {
	Verbosity   VerbosityLevel `json:"verbosity,omitempty"`
	DevMode     bool          `json:"dev_mode,omitempty"`
	Debug       bool          `json:"debug,omitempty"`
}

type EnvConfig struct {
	Variables map[string]string `json:"variables,omitempty"`
}

type CmdConfig struct {
	Aliases []string `json:"aliases,omitempty"`
}

type GlobalFlags struct{}
type Step struct{}
type Theme struct{}
type ComponentLibrary struct{}
type StyleRegistry struct{}
type ErrorTemplate struct{}
type ChaosConfig struct{}

// ExecutionContext represents the context of a command execution
type ExecutionContext struct {
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	StartTime   time.Time         `json:"start_time"`
	Steps       []StepExecution   `json:"steps"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// AARReport represents an After Action Report
type AARReport struct {
	ProjectName        string                 `json:"project_name"`
	ExecutionTime      time.Duration          `json:"execution_time"`
	TotalSteps         int                    `json:"total_steps"`
	SuccessfulSteps    int                    `json:"successful_steps"`
	FailedSteps        int                    `json:"failed_steps"`
	PerformanceMetrics map[string]interface{} `json:"performance_metrics"`
	Steps              []StepExecution        `json:"steps"`
	NextSteps          []NextStep             `json:"next_steps"`
	GeneratedAt        time.Time              `json:"generated_at"`
}

// StepExecution represents the execution of a single step
type StepExecution struct {
	Name        string            `json:"name"`
	Status      StepStatus        `json:"status"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
	Duration    time.Duration     `json:"duration"`
	Output      string            `json:"output"`
	Error       string            `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
	Priority    StepPriority      `json:"priority"`
	Category    StepCategory      `json:"category"`
}

// ExecutionResult represents the final result of an execution
type ExecutionResult struct {
	Success       bool                   `json:"success"`
	Error         string                 `json:"error,omitempty"`
	Output        string                 `json:"output"`
	Duration      time.Duration          `json:"duration"`
	ExitCode      int                    `json:"exit_code"`
	Metrics       map[string]interface{} `json:"metrics"`
}

// NextStep represents a recommended next action
type NextStep struct {
	Action      string       `json:"action"`
	Description string       `json:"description"`
	Command     string       `json:"command,omitempty"`
	Priority    StepPriority `json:"priority"`
	Category    StepCategory `json:"category"`
	WorkingDir  string       `json:"working_dir,omitempty"`
	Automated   bool         `json:"automated"`
}

// Recommendation represents a general recommendation
type Recommendation struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Priority    StepPriority `json:"priority"`
	Category    StepCategory `json:"category"`
	Impact      ChaosImpact  `json:"impact"`
}