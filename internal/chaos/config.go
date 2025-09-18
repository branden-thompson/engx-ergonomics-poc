package chaos

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AggressivenessLevel defines the chaos failure rate levels
type AggressivenessLevel int

const (
	Off AggressivenessLevel = iota // 0% failure rate
	Default                        // 0.1% failure rate
	Scout                          // 0.5% failure rate
	Aggressive                     // 1% failure rate
	Invasive                       // 5% failure rate
	Apocalyptic                    // 10% failure rate
)

// String returns the string representation of the aggressiveness level
func (a AggressivenessLevel) String() string {
	switch a {
	case Off:
		return "off"
	case Default:
		return "default"
	case Scout:
		return "scout"
	case Aggressive:
		return "aggressive"
	case Invasive:
		return "invasive"
	case Apocalyptic:
		return "apocalyptic"
	default:
		return "unknown"
	}
}

// FailureRate returns the base failure rate for this aggressiveness level
func (a AggressivenessLevel) FailureRate() float64 {
	switch a {
	case Off:
		return 0.0
	case Default:
		return 0.001 // 0.1%
	case Scout:
		return 0.005 // 0.5%
	case Aggressive:
		return 0.01 // 1%
	case Invasive:
		return 0.05 // 5%
	case Apocalyptic:
		return 0.10 // 10%
	default:
		return 0.0
	}
}

// ParseAggressivenessLevel parses a string into an AggressivenessLevel
func ParseAggressivenessLevel(s string) (AggressivenessLevel, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "off", "false", "0":
		return Off, nil
	case "default", "":
		return Default, nil
	case "scout":
		return Scout, nil
	case "aggressive":
		return Aggressive, nil
	case "invasive":
		return Invasive, nil
	case "apocalyptic":
		return Apocalyptic, nil
	default:
		return Off, fmt.Errorf("unknown aggressiveness level: %s", s)
	}
}

// ChaosConfig defines all chaos marine configuration
type ChaosConfig struct {
	// Basic activation
	Enabled             bool                `json:"enabled" yaml:"enabled"`
	AggressivenessLevel AggressivenessLevel `json:"aggressiveness_level" yaml:"aggressiveness_level"`
	RandomSeed          int64               `json:"random_seed,omitempty" yaml:"random_seed,omitempty"`

	// Safety configuration
	SafetyMode          bool     `json:"safety_mode" yaml:"safety_mode"`
	MaxInjectionCount   int64    `json:"max_injection_count" yaml:"max_injection_count"`
	AllowedOperations   []string `json:"allowed_operations" yaml:"allowed_operations"`
	ProhibitedPaths     []string `json:"prohibited_paths" yaml:"prohibited_paths"`

	// User experience
	AdaptiveDifficulty  bool `json:"adaptive_difficulty" yaml:"adaptive_difficulty"`
	EducationalMode     bool `json:"educational_mode" yaml:"educational_mode"`
	RecoveryValidation  bool `json:"recovery_validation" yaml:"recovery_validation"`
	ProgressiveHints    bool `json:"progressive_hints" yaml:"progressive_hints"`

	// Performance limits
	MaxMemoryUsageMB    int64   `json:"max_memory_usage_mb" yaml:"max_memory_usage_mb"`
	MaxCPUUsagePercent  float64 `json:"max_cpu_usage_percent" yaml:"max_cpu_usage_percent"`
	OperationTimeoutSec int64   `json:"operation_timeout_sec" yaml:"operation_timeout_sec"`

	// Analytics and telemetry
	TelemetryEnabled     bool `json:"telemetry_enabled" yaml:"telemetry_enabled"`
	AnonymousReporting   bool `json:"anonymous_reporting" yaml:"anonymous_reporting"`
	MetricsRetentionDays int  `json:"metrics_retention_days" yaml:"metrics_retention_days"`

	// Advanced features (Phase 3)
	FailureChaining bool `json:"failure_chaining" yaml:"failure_chaining"`
	CascadePrevent  bool `json:"cascade_prevent" yaml:"cascade_prevent"`
}

// NewDefaultConfig creates a default chaos configuration
func NewDefaultConfig() *ChaosConfig {
	return &ChaosConfig{
		Enabled:             false, // Must be explicitly enabled
		AggressivenessLevel: Default,
		RandomSeed:          0, // 0 means use random seed

		// Safety defaults (maximum safety)
		SafetyMode:        true,
		MaxInjectionCount: 1000,
		AllowedOperations: []string{}, // Empty means all operations allowed
		ProhibitedPaths: []string{
			"/", "/usr", "/bin", "/sbin", "/etc", "/var", "/opt", "/home",
			"C:\\", "C:\\Windows", "C:\\Program Files", "C:\\Users",
		},

		// User experience defaults
		AdaptiveDifficulty: true,
		EducationalMode:    true,
		RecoveryValidation: true,
		ProgressiveHints:   true,

		// Performance limits (conservative defaults)
		MaxMemoryUsageMB:    10,    // 10MB max memory usage
		MaxCPUUsagePercent:  5.0,   // 5% max CPU usage
		OperationTimeoutSec: 30,    // 30 second timeout

		// Telemetry defaults (privacy-first)
		TelemetryEnabled:     false, // Must be explicitly enabled
		AnonymousReporting:   false,
		MetricsRetentionDays: 7,

		// Advanced features disabled by default
		FailureChaining: false,
		CascadePrevent:  true,
	}
}

// Validate validates the chaos configuration
func (c *ChaosConfig) Validate() error {
	// Memory validation
	if c.MaxMemoryUsageMB < 1 || c.MaxMemoryUsageMB > 1024 {
		return errors.New("max_memory_usage_mb must be between 1 and 1024")
	}

	// CPU validation
	if c.MaxCPUUsagePercent < 0.1 || c.MaxCPUUsagePercent > 50.0 {
		return errors.New("max_cpu_usage_percent must be between 0.1 and 50.0")
	}

	// Injection count validation
	if c.MaxInjectionCount < 1 || c.MaxInjectionCount > 10000 {
		return errors.New("max_injection_count must be between 1 and 10000")
	}

	// Timeout validation
	if c.OperationTimeoutSec < 1 || c.OperationTimeoutSec > 3600 {
		return errors.New("operation_timeout_sec must be between 1 and 3600")
	}

	// Metrics retention validation
	if c.MetricsRetentionDays < 1 || c.MetricsRetentionDays > 90 {
		return errors.New("metrics_retention_days must be between 1 and 90")
	}

	return nil
}

// LoadChaosConfig loads chaos configuration from various sources
func LoadChaosConfig(level string, seed int64, configPath string) (*ChaosConfig, error) {
	// Start with defaults
	config := NewDefaultConfig()

	// Load from configuration file if provided
	if configPath != "" {
		if err := loadConfigFromFile(config, configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from file %s: %w", configPath, err)
		}
	}

	// Override with command line parameters
	if level != "" {
		aggressiveness, err := ParseAggressivenessLevel(level)
		if err != nil {
			return nil, fmt.Errorf("invalid aggressiveness level: %w", err)
		}
		config.AggressivenessLevel = aggressiveness
		config.Enabled = aggressiveness != Off
	}

	// Set random seed if provided
	if seed != 0 {
		config.RandomSeed = seed
	}

	// Validate final configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// loadConfigFromFile loads configuration from a JSON file
func loadConfigFromFile(config *ChaosConfig, path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("configuration file does not exist: %s", path)
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse configuration JSON: %w", err)
	}

	return nil
}

// SaveConfigToFile saves configuration to a JSON file
func (c *ChaosConfig) SaveConfigToFile(path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	return nil
}

// IsOperationAllowed checks if an operation is allowed to be chaos-injected
func (c *ChaosConfig) IsOperationAllowed(operation string) bool {
	// If no allowed operations specified, all are allowed
	if len(c.AllowedOperations) == 0 {
		return true
	}

	// Check if operation is in allowed list
	for _, allowed := range c.AllowedOperations {
		if strings.EqualFold(operation, allowed) {
			return true
		}
	}

	return false
}

// IsPathProhibited checks if a path is prohibited for chaos operations
func (c *ChaosConfig) IsPathProhibited(path string) bool {
	for _, prohibited := range c.ProhibitedPaths {
		if strings.HasPrefix(path, prohibited) {
			return true
		}
	}
	return false
}

// GetOperationTimeout returns the operation timeout as a duration
func (c *ChaosConfig) GetOperationTimeout() time.Duration {
	return time.Duration(c.OperationTimeoutSec) * time.Second
}

// String returns a string representation of the configuration
func (c *ChaosConfig) String() string {
	status := "disabled"
	if c.Enabled {
		status = fmt.Sprintf("enabled (%s)", c.AggressivenessLevel.String())
	}

	return fmt.Sprintf("ChaosConfig{status=%s, safety=%t, adaptive=%t, educational=%t}",
		status, c.SafetyMode, c.AdaptiveDifficulty, c.EducationalMode)
}