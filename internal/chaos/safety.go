package chaos

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// SafetyMonitor implements comprehensive safety boundaries for chaos operations
type SafetyMonitor struct {
	config           *ChaosConfig
	startTime        time.Time
	injectionCount   int64
	emergencyStop    bool
	mutex            sync.RWMutex
	systemSnapshot   *SystemSnapshot
	resourceMonitor  *ResourceMonitor
	healthCheck      *HealthChecker
}

// SystemSnapshot captures the current system state for integrity verification
type SystemSnapshot struct {
	workingDirectory string
	environmentVars  map[string]string
	timestamp        time.Time
	processID        int
}

// ResourceMonitor tracks resource usage to enforce limits
type ResourceMonitor struct {
	startMemory      int64
	currentMemory    int64
	maxMemoryMB      int64
	maxCPUPercent    float64
	operationTimeout time.Duration
	startTime        time.Time
}

// HealthChecker provides continuous monitoring of system integrity
type HealthChecker struct {
	lastCheck        time.Time
	checkInterval    time.Duration
	anomalyDetected  bool
	violations       []SafetyViolation
}

// SafetyViolation represents a detected safety boundary violation
type SafetyViolation struct {
	Type        ViolationType
	Description string
	Severity    ViolationSeverity
	Timestamp   time.Time
	Action      string
}

// ViolationType defines the types of safety violations
type ViolationType int

const (
	PathViolation ViolationType = iota
	ResourceViolation
	PermissionViolation
	IntegrityViolation
	TimeoutViolation
)

// ViolationSeverity defines the severity levels of safety violations
type ViolationSeverity int

const (
	Warning ViolationSeverity = iota
	Critical
	Emergency
)

// NewSafetyMonitor creates a new safety monitor with the given configuration
func NewSafetyMonitor(config *ChaosConfig) (*SafetyMonitor, error) {
	if config == nil {
		return nil, errors.New("chaos configuration cannot be nil")
	}

	// Capture initial system state
	snapshot, err := captureSystemSnapshot()
	if err != nil {
		return nil, fmt.Errorf("failed to capture initial system snapshot: %w", err)
	}

	// Initialize resource monitor
	resourceMonitor := &ResourceMonitor{
		maxMemoryMB:      config.MaxMemoryUsageMB,
		maxCPUPercent:    config.MaxCPUUsagePercent,
		operationTimeout: config.GetOperationTimeout(),
		startTime:        time.Now(),
	}

	// Initialize health checker
	healthChecker := &HealthChecker{
		lastCheck:     time.Now(),
		checkInterval: 30 * time.Second, // Check every 30 seconds
		violations:    make([]SafetyViolation, 0),
	}

	monitor := &SafetyMonitor{
		config:          config,
		startTime:       time.Now(),
		injectionCount:  0,
		emergencyStop:   false,
		systemSnapshot:  snapshot,
		resourceMonitor: resourceMonitor,
		healthCheck:     healthChecker,
	}

	return monitor, nil
}

// captureSystemSnapshot captures the current system state
func captureSystemSnapshot() (*SystemSnapshot, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	envVars := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envVars[parts[0]] = parts[1]
		}
	}

	return &SystemSnapshot{
		workingDirectory: wd,
		environmentVars:  envVars,
		timestamp:        time.Now(),
		processID:        os.Getpid(),
	}, nil
}

// IsPathSafe checks if a given path is safe for chaos operations
func (sm *SafetyMonitor) IsPathSafe(path string) error {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Emergency stop check
	if sm.emergencyStop {
		return errors.New("SAFETY VIOLATION: Emergency stop active")
	}

	// Clean and normalize the path
	cleanPath := filepath.Clean(path)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Check prohibited paths
	if sm.config.IsPathProhibited(absPath) {
		violation := SafetyViolation{
			Type:        PathViolation,
			Description: fmt.Sprintf("Attempted access to prohibited path: %s", absPath),
			Severity:    Critical,
			Timestamp:   time.Now(),
			Action:      "Path access blocked",
		}
		sm.recordViolation(violation)
		return fmt.Errorf("SAFETY VIOLATION: Path '%s' is prohibited", absPath)
	}

	// Additional platform-specific safety checks
	if err := sm.platformSpecificPathCheck(absPath); err != nil {
		return err
	}

	return nil
}

// platformSpecificPathCheck performs OS-specific path safety checks
func (sm *SafetyMonitor) platformSpecificPathCheck(path string) error {
	switch runtime.GOOS {
	case "windows":
		return sm.windowsPathCheck(path)
	case "darwin", "linux":
		return sm.unixPathCheck(path)
	default:
		// Conservative approach for unknown platforms
		return fmt.Errorf("SAFETY VIOLATION: Platform %s not supported", runtime.GOOS)
	}
}

// windowsPathCheck performs Windows-specific safety checks
func (sm *SafetyMonitor) windowsPathCheck(path string) error {
	// Critical Windows system paths
	dangerousPaths := []string{
		"C:\\Windows",
		"C:\\Program Files",
		"C:\\Program Files (x86)",
		"C:\\Users\\Default",
		"C:\\ProgramData",
	}

	upperPath := strings.ToUpper(path)
	for _, dangerous := range dangerousPaths {
		if strings.HasPrefix(upperPath, strings.ToUpper(dangerous)) {
			return fmt.Errorf("SAFETY VIOLATION: Windows system path access blocked: %s", path)
		}
	}

	return nil
}

// unixPathCheck performs Unix/Linux/macOS specific safety checks
func (sm *SafetyMonitor) unixPathCheck(path string) error {
	// Critical Unix system paths
	dangerousPaths := []string{
		"/bin", "/sbin", "/usr/bin", "/usr/sbin",
		"/etc", "/var", "/opt", "/boot",
		"/sys", "/proc", "/dev",
		"/Library/System", "/System",
	}

	for _, dangerous := range dangerousPaths {
		if strings.HasPrefix(path, dangerous) {
			return fmt.Errorf("SAFETY VIOLATION: Unix system path access blocked: %s", path)
		}
	}

	return nil
}

// IsOperationSafe checks if a chaos operation can be safely executed
func (sm *SafetyMonitor) IsOperationSafe(operation string) error {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Emergency stop check
	if sm.emergencyStop {
		return errors.New("SAFETY VIOLATION: Emergency stop active")
	}

	// Check injection count limits
	if sm.injectionCount >= sm.config.MaxInjectionCount {
		violation := SafetyViolation{
			Type:        ResourceViolation,
			Description: fmt.Sprintf("Maximum injection count exceeded: %d", sm.injectionCount),
			Severity:    Critical,
			Timestamp:   time.Now(),
			Action:      "Operation blocked",
		}
		sm.recordViolation(violation)
		return fmt.Errorf("SAFETY VIOLATION: Maximum injection count (%d) exceeded", sm.config.MaxInjectionCount)
	}

	// Check if operation is allowed
	if !sm.config.IsOperationAllowed(operation) {
		violation := SafetyViolation{
			Type:        PermissionViolation,
			Description: fmt.Sprintf("Operation not in allowed list: %s", operation),
			Severity:    Warning,
			Timestamp:   time.Now(),
			Action:      "Operation blocked",
		}
		sm.recordViolation(violation)
		return fmt.Errorf("SAFETY VIOLATION: Operation '%s' not allowed", operation)
	}

	// Check resource limits
	if err := sm.checkResourceLimits(); err != nil {
		return err
	}

	// Check operation timeout
	if time.Since(sm.startTime) > sm.config.GetOperationTimeout() {
		violation := SafetyViolation{
			Type:        TimeoutViolation,
			Description: fmt.Sprintf("Operation timeout exceeded: %v", sm.config.GetOperationTimeout()),
			Severity:    Critical,
			Timestamp:   time.Now(),
			Action:      "Operation blocked",
		}
		sm.recordViolation(violation)
		return fmt.Errorf("SAFETY VIOLATION: Operation timeout exceeded")
	}

	return nil
}

// checkResourceLimits verifies resource usage is within configured limits
func (sm *SafetyMonitor) checkResourceLimits() error {
	// Get current memory usage
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	currentMemoryMB := int64(m.Alloc) / 1024 / 1024

	// Check memory limit
	if currentMemoryMB > sm.config.MaxMemoryUsageMB {
		violation := SafetyViolation{
			Type:        ResourceViolation,
			Description: fmt.Sprintf("Memory usage exceeded: %dMB > %dMB", currentMemoryMB, sm.config.MaxMemoryUsageMB),
			Severity:    Critical,
			Timestamp:   time.Now(),
			Action:      "Operation blocked",
		}
		sm.recordViolation(violation)
		return fmt.Errorf("SAFETY VIOLATION: Memory usage exceeded limit (%dMB > %dMB)",
			currentMemoryMB, sm.config.MaxMemoryUsageMB)
	}

	// CPU usage check would require more complex monitoring
	// For now, we implement a simple goroutine count check as a proxy
	numGoroutines := runtime.NumGoroutine()
	if numGoroutines > 100 { // Arbitrary limit for chaos operations
		violation := SafetyViolation{
			Type:        ResourceViolation,
			Description: fmt.Sprintf("Goroutine count too high: %d", numGoroutines),
			Severity:    Warning,
			Timestamp:   time.Now(),
			Action:      "Warning issued",
		}
		sm.recordViolation(violation)
	}

	return nil
}

// RecordInjection records a chaos injection for tracking and limits
func (sm *SafetyMonitor) RecordInjection(operation string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if sm.emergencyStop {
		return errors.New("SAFETY VIOLATION: Emergency stop active")
	}

	sm.injectionCount++

	// Log injection for audit trail
	if sm.config.TelemetryEnabled {
		// In a real implementation, this would go to a proper logging system
		fmt.Printf("[CHAOS AUDIT] Injection #%d: %s at %v\n",
			sm.injectionCount, operation, time.Now().Format(time.RFC3339))
	}

	return nil
}

// recordViolation records a safety violation
func (sm *SafetyMonitor) recordViolation(violation SafetyViolation) {
	sm.healthCheck.violations = append(sm.healthCheck.violations, violation)
	sm.healthCheck.anomalyDetected = true

	// For critical and emergency violations, trigger emergency stop
	if violation.Severity >= Critical {
		sm.emergencyStop = true
	}

	// Log violation for audit trail
	fmt.Printf("[SAFETY VIOLATION] %s: %s (%s)\n",
		violation.Type, violation.Description, violation.Severity)
}

// PerformHealthCheck performs a comprehensive system health check
func (sm *SafetyMonitor) PerformHealthCheck() error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	now := time.Now()
	sm.healthCheck.lastCheck = now

	// Verify system integrity hasn't changed
	currentSnapshot, err := captureSystemSnapshot()
	if err != nil {
		return fmt.Errorf("failed to capture current system snapshot: %w", err)
	}

	// Check working directory hasn't changed unexpectedly
	if currentSnapshot.workingDirectory != sm.systemSnapshot.workingDirectory {
		violation := SafetyViolation{
			Type:        IntegrityViolation,
			Description: fmt.Sprintf("Working directory changed: %s -> %s",
				sm.systemSnapshot.workingDirectory, currentSnapshot.workingDirectory),
			Severity:    Warning,
			Timestamp:   now,
			Action:      "Directory change detected",
		}
		sm.recordViolation(violation)
	}

	// Check resource usage
	if err := sm.checkResourceLimits(); err != nil {
		return err
	}

	// Clear old violations (older than 1 hour)
	sm.cleanupOldViolations(1 * time.Hour)

	return nil
}

// cleanupOldViolations removes violations older than the specified duration
func (sm *SafetyMonitor) cleanupOldViolations(maxAge time.Duration) {
	cutoff := time.Now().Add(-maxAge)
	violations := make([]SafetyViolation, 0)

	for _, violation := range sm.healthCheck.violations {
		if violation.Timestamp.After(cutoff) {
			violations = append(violations, violation)
		}
	}

	sm.healthCheck.violations = violations
	sm.healthCheck.anomalyDetected = len(violations) > 0
}

// EmergencyStop immediately stops all chaos operations
func (sm *SafetyMonitor) EmergencyStop() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.emergencyStop = true

	violation := SafetyViolation{
		Type:        IntegrityViolation,
		Description: "Emergency stop activated",
		Severity:    Emergency,
		Timestamp:   time.Now(),
		Action:      "All operations halted",
	}
	sm.recordViolation(violation)
}

// Reset resets the safety monitor to allow operations again
func (sm *SafetyMonitor) Reset() error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Perform final health check before reset
	if err := sm.checkResourceLimits(); err != nil {
		return fmt.Errorf("cannot reset with resource violations: %w", err)
	}

	sm.emergencyStop = false
	sm.injectionCount = 0
	sm.startTime = time.Now()
	sm.healthCheck.violations = make([]SafetyViolation, 0)
	sm.healthCheck.anomalyDetected = false

	// Capture fresh system snapshot
	snapshot, err := captureSystemSnapshot()
	if err != nil {
		return fmt.Errorf("failed to capture fresh system snapshot: %w", err)
	}
	sm.systemSnapshot = snapshot

	return nil
}

// GetSafetyStatus returns the current safety status
func (sm *SafetyMonitor) GetSafetyStatus() SafetyStatus {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	status := SafetyStatus{
		EmergencyStop:     sm.emergencyStop,
		InjectionCount:    sm.injectionCount,
		MaxInjections:     sm.config.MaxInjectionCount,
		AnomalyDetected:   sm.healthCheck.anomalyDetected,
		ViolationCount:    len(sm.healthCheck.violations),
		OperationUptime:   time.Since(sm.startTime),
		LastHealthCheck:   sm.healthCheck.lastCheck,
		ResourceUsage:     sm.getResourceUsage(),
	}

	return status
}

// SafetyStatus represents the current safety monitoring status
type SafetyStatus struct {
	EmergencyStop     bool
	InjectionCount    int64
	MaxInjections     int64
	AnomalyDetected   bool
	ViolationCount    int
	OperationUptime   time.Duration
	LastHealthCheck   time.Time
	ResourceUsage     ResourceUsage
}

// ResourceUsage represents current resource utilization
type ResourceUsage struct {
	MemoryMB      int64
	MaxMemoryMB   int64
	GoroutineCount int
	CPUPercent    float64 // Placeholder for future implementation
}

// getResourceUsage returns current resource usage statistics
func (sm *SafetyMonitor) getResourceUsage() ResourceUsage {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	currentMemoryMB := int64(m.Alloc) / 1024 / 1024

	return ResourceUsage{
		MemoryMB:       currentMemoryMB,
		MaxMemoryMB:    sm.config.MaxMemoryUsageMB,
		GoroutineCount: runtime.NumGoroutine(),
		CPUPercent:     0.0, // Placeholder - would need platform-specific implementation
	}
}

// String methods for enums
func (vt ViolationType) String() string {
	switch vt {
	case PathViolation:
		return "PATH_VIOLATION"
	case ResourceViolation:
		return "RESOURCE_VIOLATION"
	case PermissionViolation:
		return "PERMISSION_VIOLATION"
	case IntegrityViolation:
		return "INTEGRITY_VIOLATION"
	case TimeoutViolation:
		return "TIMEOUT_VIOLATION"
	default:
		return "UNKNOWN_VIOLATION"
	}
}

func (vs ViolationSeverity) String() string {
	switch vs {
	case Warning:
		return "WARNING"
	case Critical:
		return "CRITICAL"
	case Emergency:
		return "EMERGENCY"
	default:
		return "UNKNOWN"
	}
}