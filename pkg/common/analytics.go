package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
)

// InteractionAnalytics tracks CLI usage patterns for simulation analysis
// Focus: Understanding developer workflow patterns and command ergonomics
type InteractionAnalytics struct {
	sessionID    string
	startTime    time.Time
	events       []InteractionEvent
	commandStats map[string]*CommandStats
	sessionStats *SessionStats
	logger       interfaces.Logger
	mu           sync.RWMutex
	enabled      bool
}

// InteractionEvent represents a single CLI interaction
type InteractionEvent struct {
	Timestamp   time.Time         `json:"timestamp"`
	SessionID   string            `json:"session_id"`
	Command     string            `json:"command"`
	Subcommand  string            `json:"subcommand,omitempty"`
	Args        []string          `json:"args,omitempty"`
	Flags       map[string]string `json:"flags,omitempty"`
	Duration    time.Duration     `json:"duration"`
	Success     bool              `json:"success"`
	ErrorMsg    string            `json:"error_msg,omitempty"`
	UserPath    string            `json:"user_path"` // Command sequence leading to this event
	Context     map[string]string `json:"context,omitempty"`
}

// CommandStats tracks statistics for individual commands
type CommandStats struct {
	Command      string        `json:"command"`
	TotalUses    int           `json:"total_uses"`
	SuccessCount int           `json:"success_count"`
	FailureCount int           `json:"failure_count"`
	AvgDuration  time.Duration `json:"avg_duration"`
	TotalTime    time.Duration `json:"total_time"`
	FirstUsed    time.Time     `json:"first_used"`
	LastUsed     time.Time     `json:"last_used"`
	CommonArgs   []string      `json:"common_args"`
	CommonFlags  []string      `json:"common_flags"`
	ErrorTypes   map[string]int `json:"error_types"`
}

// SessionStats tracks overall session metrics
type SessionStats struct {
	SessionID       string            `json:"session_id"`
	StartTime       time.Time         `json:"start_time"`
	Duration        time.Duration     `json:"duration"`
	TotalCommands   int               `json:"total_commands"`
	UniqueCommands  int               `json:"unique_commands"`
	SuccessRate     float64           `json:"success_rate"`
	MostUsedCommand string            `json:"most_used_command"`
	CommandFlow     []string          `json:"command_flow"`
	Patterns        []WorkflowPattern `json:"patterns"`
}

// WorkflowPattern represents common command sequences
type WorkflowPattern struct {
	Name        string    `json:"name"`
	Sequence    []string  `json:"sequence"`
	Frequency   int       `json:"frequency"`
	AvgDuration time.Duration `json:"avg_duration"`
	Description string    `json:"description"`
}

// NewInteractionAnalytics creates a new analytics tracker
func NewInteractionAnalytics(deps *Dependencies) *InteractionAnalytics {
	sessionID := fmt.Sprintf("session_%d", time.Now().Unix())

	return &InteractionAnalytics{
		sessionID:    sessionID,
		startTime:    time.Now(),
		events:       make([]InteractionEvent, 0),
		commandStats: make(map[string]*CommandStats),
		sessionStats: &SessionStats{
			SessionID:      sessionID,
			StartTime:      time.Now(),
			CommandFlow:    make([]string, 0),
			Patterns:       make([]WorkflowPattern, 0),
		},
		logger:  deps.Logger.WithComponent("analytics"),
		enabled: true,
	}
}

// TrackCommand records a command execution for analysis
func (ia *InteractionAnalytics) TrackCommand(command, subcommand string, args []string, flags map[string]string, duration time.Duration, success bool, errorMsg string) {
	if !ia.enabled {
		return
	}

	ia.mu.Lock()
	defer ia.mu.Unlock()

	// Create interaction event
	event := InteractionEvent{
		Timestamp:  time.Now(),
		SessionID:  ia.sessionID,
		Command:    command,
		Subcommand: subcommand,
		Args:       args,
		Flags:      flags,
		Duration:   duration,
		Success:    success,
		ErrorMsg:   errorMsg,
		UserPath:   ia.buildUserPath(),
		Context:    make(map[string]string),
	}

	// Add contextual information
	if subcommand != "" {
		event.Context["full_command"] = fmt.Sprintf("%s %s", command, subcommand)
	} else {
		event.Context["full_command"] = command
	}

	ia.events = append(ia.events, event)

	// Update command statistics
	ia.updateCommandStats(event)

	// Update session statistics
	ia.updateSessionStats(event)

	// Detect workflow patterns
	ia.detectWorkflowPatterns()

	ia.logger.Debug("Tracked command interaction", map[string]interface{}{
		"command":     command,
		"subcommand":  subcommand,
		"duration":    duration,
		"success":     success,
		"session_id":  ia.sessionID,
	})
}

// TrackTemplateSelection tracks template discovery patterns
func (ia *InteractionAnalytics) TrackTemplateSelection(templateID string, discoveryMethod string, searchQuery string, timeToSelect time.Duration) {
	if !ia.enabled {
		return
	}

	context := map[string]string{
		"template_id":       templateID,
		"discovery_method":  discoveryMethod, // list, search, recommended, etc.
		"search_query":     searchQuery,
		"time_to_select":   timeToSelect.String(),
	}

	ia.TrackCommand("templates", "select", []string{templateID}, map[string]string{
		"method": discoveryMethod,
		"query":  searchQuery,
	}, timeToSelect, true, "")

	ia.logger.Info("Template selection tracked", context)
}

// GetSessionAnalytics returns current session analytics
func (ia *InteractionAnalytics) GetSessionAnalytics() *SessionStats {
	ia.mu.RLock()
	defer ia.mu.RUnlock()

	// Update duration
	ia.sessionStats.Duration = time.Since(ia.startTime)

	// Calculate success rate
	if ia.sessionStats.TotalCommands > 0 {
		successCount := 0
		for _, event := range ia.events {
			if event.Success {
				successCount++
			}
		}
		ia.sessionStats.SuccessRate = float64(successCount) / float64(ia.sessionStats.TotalCommands) * 100
	}

	return ia.sessionStats
}

// GetCommandStats returns statistics for all commands
func (ia *InteractionAnalytics) GetCommandStats() map[string]*CommandStats {
	ia.mu.RLock()
	defer ia.mu.RUnlock()

	// Return a copy to prevent modification
	result := make(map[string]*CommandStats)
	for k, v := range ia.commandStats {
		result[k] = v
	}
	return result
}

// GetWorkflowPatterns returns detected workflow patterns
func (ia *InteractionAnalytics) GetWorkflowPatterns() []WorkflowPattern {
	ia.mu.RLock()
	defer ia.mu.RUnlock()

	return ia.sessionStats.Patterns
}

// ExportAnalytics exports analytics data to JSON file
func (ia *InteractionAnalytics) ExportAnalytics(outputPath string) error {
	ia.mu.RLock()
	defer ia.mu.RUnlock()

	// Prepare export data
	exportData := struct {
		SessionStats *SessionStats              `json:"session_stats"`
		CommandStats map[string]*CommandStats   `json:"command_stats"`
		Events       []InteractionEvent         `json:"events"`
		ExportTime   time.Time                  `json:"export_time"`
	}{
		SessionStats: ia.GetSessionAnalytics(),
		CommandStats: ia.commandStats,
		Events:       ia.events,
		ExportTime:   time.Now(),
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write JSON file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(exportData); err != nil {
		return fmt.Errorf("failed to encode analytics data: %w", err)
	}

	ia.logger.Info("Analytics exported", map[string]interface{}{
		"output_path":   outputPath,
		"event_count":   len(ia.events),
		"command_count": len(ia.commandStats),
		"session_id":    ia.sessionID,
	})

	return nil
}

// Internal helper methods

func (ia *InteractionAnalytics) buildUserPath() string {
	if len(ia.sessionStats.CommandFlow) == 0 {
		return ""
	}

	// Return last 3 commands as context
	start := len(ia.sessionStats.CommandFlow) - 3
	if start < 0 {
		start = 0
	}

	path := ""
	for i := start; i < len(ia.sessionStats.CommandFlow); i++ {
		if i > start {
			path += " → "
		}
		path += ia.sessionStats.CommandFlow[i]
	}
	return path
}

func (ia *InteractionAnalytics) updateCommandStats(event InteractionEvent) {
	cmdKey := event.Command
	if event.Subcommand != "" {
		cmdKey = fmt.Sprintf("%s %s", event.Command, event.Subcommand)
	}

	stats, exists := ia.commandStats[cmdKey]
	if !exists {
		stats = &CommandStats{
			Command:     cmdKey,
			FirstUsed:   event.Timestamp,
			ErrorTypes:  make(map[string]int),
			CommonArgs:  make([]string, 0),
			CommonFlags: make([]string, 0),
		}
		ia.commandStats[cmdKey] = stats
	}

	// Update basic counters
	stats.TotalUses++
	stats.LastUsed = event.Timestamp
	stats.TotalTime += event.Duration

	if event.Success {
		stats.SuccessCount++
	} else {
		stats.FailureCount++
		if event.ErrorMsg != "" {
			stats.ErrorTypes[event.ErrorMsg]++
		}
	}

	// Calculate average duration
	stats.AvgDuration = stats.TotalTime / time.Duration(stats.TotalUses)

	// Track common arguments and flags
	ia.updateCommonUsage(stats, event.Args, event.Flags)
}

func (ia *InteractionAnalytics) updateCommonUsage(stats *CommandStats, args []string, flags map[string]string) {
	// Update common args (simplified - just track first arg)
	if len(args) > 0 && !contains(stats.CommonArgs, args[0]) {
		stats.CommonArgs = append(stats.CommonArgs, args[0])
	}

	// Update common flags
	for flag := range flags {
		if !contains(stats.CommonFlags, flag) {
			stats.CommonFlags = append(stats.CommonFlags, flag)
		}
	}
}

func (ia *InteractionAnalytics) updateSessionStats(event InteractionEvent) {
	cmdKey := event.Command
	if event.Subcommand != "" {
		cmdKey = fmt.Sprintf("%s %s", event.Command, event.Subcommand)
	}

	ia.sessionStats.TotalCommands++
	ia.sessionStats.CommandFlow = append(ia.sessionStats.CommandFlow, cmdKey)

	// Update unique commands count
	unique := make(map[string]bool)
	for _, cmd := range ia.sessionStats.CommandFlow {
		unique[cmd] = true
	}
	ia.sessionStats.UniqueCommands = len(unique)

	// Find most used command
	maxUses := 0
	for _, stats := range ia.commandStats {
		if stats.TotalUses > maxUses {
			maxUses = stats.TotalUses
			ia.sessionStats.MostUsedCommand = stats.Command
		}
	}
}

func (ia *InteractionAnalytics) detectWorkflowPatterns() {
	// Simple pattern detection - look for common 2-3 command sequences
	if len(ia.sessionStats.CommandFlow) < 2 {
		return
	}

	patterns := make(map[string]*WorkflowPattern)

	// Detect 2-command patterns
	for i := 0; i < len(ia.sessionStats.CommandFlow)-1; i++ {
		seq := []string{ia.sessionStats.CommandFlow[i], ia.sessionStats.CommandFlow[i+1]}
		key := fmt.Sprintf("%s → %s", seq[0], seq[1])

		if pattern, exists := patterns[key]; exists {
			pattern.Frequency++
		} else {
			patterns[key] = &WorkflowPattern{
				Name:        ia.classifyPattern(seq),
				Sequence:    seq,
				Frequency:   1,
				Description: ia.describePattern(seq),
			}
		}
	}

	// Convert to slice and update session stats
	ia.sessionStats.Patterns = make([]WorkflowPattern, 0, len(patterns))
	for _, pattern := range patterns {
		ia.sessionStats.Patterns = append(ia.sessionStats.Patterns, *pattern)
	}
}

func (ia *InteractionAnalytics) classifyPattern(sequence []string) string {
	if len(sequence) >= 2 {
		first, second := sequence[0], sequence[1]

		if first == "templates list" && second == "templates search" {
			return "Discovery → Search"
		}
		if first == "templates search" && second == "templates info" {
			return "Search → Details"
		}
		if first == "templates info" && second == "create" {
			return "Details → Create"
		}
		if first == "templates recommended" && second == "create" {
			return "Quick Start"
		}
	}

	return "Custom Sequence"
}

func (ia *InteractionAnalytics) describePattern(sequence []string) string {
	if len(sequence) >= 2 {
		first, second := sequence[0], sequence[1]

		if first == "templates list" && second == "templates search" {
			return "User browsed all templates then searched for specific criteria"
		}
		if first == "templates search" && second == "templates info" {
			return "User searched templates then requested detailed information"
		}
		if first == "templates info" && second == "create" {
			return "User reviewed template details then created project"
		}
		if first == "templates recommended" && second == "create" {
			return "User chose recommended template and created project immediately"
		}
	}

	return fmt.Sprintf("Custom workflow: %s", sequence)
}

// Utility functions
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Enable/Disable analytics
func (ia *InteractionAnalytics) Enable() {
	ia.mu.Lock()
	defer ia.mu.Unlock()
	ia.enabled = true
}

func (ia *InteractionAnalytics) Disable() {
	ia.mu.Lock()
	defer ia.mu.Unlock()
	ia.enabled = false
}

func (ia *InteractionAnalytics) IsEnabled() bool {
	ia.mu.RLock()
	defer ia.mu.RUnlock()
	return ia.enabled
}