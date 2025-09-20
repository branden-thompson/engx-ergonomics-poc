package common

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

// AnalyticsUI provides CLI interface for interaction analytics
type AnalyticsUI struct {
	analytics *InteractionAnalytics
	logger    interfaces.Logger
}

// NewAnalyticsUI creates a new analytics UI
func NewAnalyticsUI(analytics *InteractionAnalytics, deps *Dependencies) *AnalyticsUI {
	return &AnalyticsUI{
		analytics: analytics,
		logger:    deps.Logger.WithComponent("analytics-ui"),
	}
}

// ShowSessionSummary displays current session analytics summary
func (aui *AnalyticsUI) ShowSessionSummary() error {
	sessionStats := aui.analytics.GetSessionAnalytics()

	fmt.Println("📊 CLI Session Analytics Summary")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Session ID: %s\n", sessionStats.SessionID)
	fmt.Printf("Duration: %s\n", formatDuration(sessionStats.Duration))
	fmt.Printf("Start Time: %s\n\n", sessionStats.StartTime.Format("2006-01-02 15:04:05"))

	fmt.Printf("📋 Command Activity:\n")
	fmt.Printf("   Total Commands: %d\n", sessionStats.TotalCommands)
	fmt.Printf("   Unique Commands: %d\n", sessionStats.UniqueCommands)
	fmt.Printf("   Success Rate: %.1f%%\n", sessionStats.SuccessRate)
	if sessionStats.MostUsedCommand != "" {
		fmt.Printf("   Most Used: %s\n", sessionStats.MostUsedCommand)
	}

	if len(sessionStats.CommandFlow) > 0 {
		fmt.Printf("\n🔄 Recent Command Flow:\n")
		recentCommands := sessionStats.CommandFlow
		if len(recentCommands) > 5 {
			recentCommands = recentCommands[len(recentCommands)-5:]
		}
		fmt.Printf("   %s\n", strings.Join(recentCommands, " → "))
	}

	if len(sessionStats.Patterns) > 0 {
		fmt.Printf("\n🎯 Detected Workflow Patterns:\n")
		for i, pattern := range sessionStats.Patterns {
			if i >= 3 { // Show top 3 patterns
				break
			}
			fmt.Printf("   %s (%dx)\n", pattern.Name, pattern.Frequency)
			fmt.Printf("      %s\n", strings.Join(pattern.Sequence, " → "))
		}
	}

	cmdFormatter := components.NewCommandFormatter()
	fmt.Printf("\n💡 Use %s for comprehensive analysis\n", cmdFormatter.FormatCommandInBackticks("engx analytics details"))
	fmt.Printf("💡 Use %s to save data for review\n", cmdFormatter.FormatCommandInBackticks("engx analytics export"))

	return nil
}

// ShowDetailedAnalytics displays comprehensive analytics data
func (aui *AnalyticsUI) ShowDetailedAnalytics() error {
	sessionStats := aui.analytics.GetSessionAnalytics()
	commandStats := aui.analytics.GetCommandStats()

	fmt.Println("📊 Detailed CLI Analytics")
	fmt.Println(strings.Repeat("=", 60))

	// Session Overview
	fmt.Printf("📋 Session Overview:\n")
	fmt.Printf("   Session ID: %s\n", sessionStats.SessionID)
	fmt.Printf("   Duration: %s\n", formatDuration(sessionStats.Duration))
	fmt.Printf("   Commands Executed: %d\n", sessionStats.TotalCommands)
	fmt.Printf("   Unique Commands: %d\n", sessionStats.UniqueCommands)
	fmt.Printf("   Success Rate: %.1f%%\n\n", sessionStats.SuccessRate)

	// Command Statistics
	if len(commandStats) > 0 {
		fmt.Printf("🔧 Command Usage Statistics:\n")
		fmt.Printf("%-25s %-6s %-6s %-8s %-12s %-10s\n",
			"COMMAND", "USES", "RATE", "SUCCESS", "AVG TIME", "LAST USED")
		fmt.Println(strings.Repeat("-", 75))

		// Sort commands by usage
		sortedCommands := make([]*CommandStats, 0, len(commandStats))
		for _, stats := range commandStats {
			sortedCommands = append(sortedCommands, stats)
		}
		sort.Slice(sortedCommands, func(i, j int) bool {
			return sortedCommands[i].TotalUses > sortedCommands[j].TotalUses
		})

		for _, stats := range sortedCommands {
			successRate := float64(stats.SuccessCount) / float64(stats.TotalUses) * 100
			lastUsed := "never"
			if !stats.LastUsed.IsZero() {
				lastUsed = stats.LastUsed.Format("15:04:05")
			}

			fmt.Printf("%-25s %-6d %-6.1f%% %-8d %-12s %-10s\n",
				truncateString(stats.Command, 25),
				stats.TotalUses,
				successRate,
				stats.SuccessCount,
				formatDuration(stats.AvgDuration),
				lastUsed)
		}
	}

	// Workflow Patterns
	patterns := aui.analytics.GetWorkflowPatterns()
	if len(patterns) > 0 {
		fmt.Printf("\n🎯 Workflow Pattern Analysis:\n")

		// Sort patterns by frequency
		sort.Slice(patterns, func(i, j int) bool {
			return patterns[i].Frequency > patterns[j].Frequency
		})

		for i, pattern := range patterns {
			if i >= 5 { // Show top 5 patterns
				break
			}
			fmt.Printf("\n%d. %s (occurred %d time(s))\n", i+1, pattern.Name, pattern.Frequency)
			fmt.Printf("   Sequence: %s\n", strings.Join(pattern.Sequence, " → "))
			fmt.Printf("   Description: %s\n", pattern.Description)
		}
	}

	// Error Analysis
	if len(commandStats) > 0 {
		fmt.Printf("\n❌ Error Analysis:\n")
		hasErrors := false
		for _, stats := range commandStats {
			if stats.FailureCount > 0 {
				hasErrors = true
				fmt.Printf("   %s: %d failures\n", stats.Command, stats.FailureCount)
				if len(stats.ErrorTypes) > 0 {
					for errorMsg, count := range stats.ErrorTypes {
						fmt.Printf("      • %s (%dx)\n", errorMsg, count)
					}
				}
			}
		}
		if !hasErrors {
			fmt.Printf("   ✅ No errors detected in this session\n")
		}
	}

	cmdFormatter := components.NewCommandFormatter()
	fmt.Printf("\n💡 Use %s to save this data\n", cmdFormatter.FormatCommandInBackticks("engx analytics export"))

	return nil
}

// ShowWorkflowPatterns displays discovered workflow patterns
func (aui *AnalyticsUI) ShowWorkflowPatterns() error {
	patterns := aui.analytics.GetWorkflowPatterns()

	fmt.Println("🎯 CLI Workflow Pattern Analysis")
	fmt.Println(strings.Repeat("=", 50))

	if len(patterns) == 0 {
		fmt.Println("No workflow patterns detected yet.")
		fmt.Println("Try using multiple commands to build up pattern data.")
		return nil
	}

	// Sort patterns by frequency
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Frequency > patterns[j].Frequency
	})

	fmt.Printf("Detected %d workflow pattern(s):\n\n", len(patterns))

	for i, pattern := range patterns {
		fmt.Printf("%d. %s\n", i+1, pattern.Name)
		fmt.Printf("   🔄 Sequence: %s\n", strings.Join(pattern.Sequence, " → "))
		fmt.Printf("   📊 Frequency: %d occurrence(s)\n", pattern.Frequency)
		fmt.Printf("   📝 %s\n", pattern.Description)

		if i < len(patterns)-1 {
			fmt.Println()
		}
	}

	fmt.Printf("\n💡 Patterns help identify common developer workflows\n")
	cmdFormatter := components.NewCommandFormatter()
	fmt.Printf("💡 Use %s for complete analysis\n", cmdFormatter.FormatCommandInBackticks("engx analytics details"))

	return nil
}

// ShowUsageStats displays command usage statistics
func (aui *AnalyticsUI) ShowUsageStats() error {
	commandStats := aui.analytics.GetCommandStats()
	sessionStats := aui.analytics.GetSessionAnalytics()

	fmt.Println("📈 CLI Usage Statistics")
	fmt.Println(strings.Repeat("=", 45))

	if len(commandStats) == 0 {
		fmt.Println("No usage data available yet.")
		fmt.Println("Execute some commands to build up statistics.")
		return nil
	}

	fmt.Printf("Session Duration: %s\n", formatDuration(sessionStats.Duration))
	fmt.Printf("Total Commands: %d\n", sessionStats.TotalCommands)
	fmt.Printf("Overall Success Rate: %.1f%%\n\n", sessionStats.SuccessRate)

	// Most/least used commands
	sortedCommands := make([]*CommandStats, 0, len(commandStats))
	for _, stats := range commandStats {
		sortedCommands = append(sortedCommands, stats)
	}
	sort.Slice(sortedCommands, func(i, j int) bool {
		return sortedCommands[i].TotalUses > sortedCommands[j].TotalUses
	})

	fmt.Printf("🔥 Most Used Commands:\n")
	for i, stats := range sortedCommands {
		if i >= 3 { // Top 3
			break
		}
		successRate := float64(stats.SuccessCount) / float64(stats.TotalUses) * 100
		fmt.Printf("   %d. %s (%d uses, %.1f%% success)\n",
			i+1, stats.Command, stats.TotalUses, successRate)
	}

	// Command performance
	fmt.Printf("\n⚡ Command Performance:\n")
	for i, stats := range sortedCommands {
		if i >= 5 { // Top 5
			break
		}
		fmt.Printf("   %s: avg %s\n",
			truncateString(stats.Command, 20), formatDuration(stats.AvgDuration))
	}

	// Recent activity
	if len(sessionStats.CommandFlow) > 0 {
		fmt.Printf("\n🕒 Recent Activity:\n")
		recentCommands := sessionStats.CommandFlow
		if len(recentCommands) > 8 {
			recentCommands = recentCommands[len(recentCommands)-8:]
		}

		for i, cmd := range recentCommands {
			if i > 0 {
				fmt.Printf(" → ")
			}
			fmt.Printf("%s", cmd)
		}
		fmt.Println()
	}

	cmdFormatter := components.NewCommandFormatter()
	fmt.Printf("\n💡 Use %s to see workflow analysis\n", cmdFormatter.FormatCommandInBackticks("engx analytics patterns"))

	return nil
}

// ExportAnalytics exports analytics data to file
func (aui *AnalyticsUI) ExportAnalytics(outputPath string) error {
	if outputPath == "" {
		// Generate default filename
		sessionStats := aui.analytics.GetSessionAnalytics()
		timestamp := time.Now().Format("20060102_150405")
		outputPath = fmt.Sprintf("engx_analytics_%s_%s.json",
			sessionStats.SessionID, timestamp)
	}

	fmt.Printf("📤 Exporting analytics data...\n")

	if err := aui.analytics.ExportAnalytics(outputPath); err != nil {
		fmt.Printf("❌ Export failed: %v\n", err)
		return err
	}

	fmt.Printf("✅ Analytics exported to: %s\n", outputPath)

	// Show export summary
	sessionStats := aui.analytics.GetSessionAnalytics()
	commandStats := aui.analytics.GetCommandStats()

	fmt.Printf("\n📋 Export Summary:\n")
	fmt.Printf("   Session ID: %s\n", sessionStats.SessionID)
	fmt.Printf("   Commands: %d\n", sessionStats.TotalCommands)
	fmt.Printf("   Unique Commands: %d\n", len(commandStats))
	fmt.Printf("   Patterns: %d\n", len(sessionStats.Patterns))
	fmt.Printf("   Export Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}

// ShowAnalyticsStatus displays analytics system status
func (aui *AnalyticsUI) ShowAnalyticsStatus() error {
	fmt.Println("🔍 Analytics System Status")
	fmt.Println(strings.Repeat("=", 40))

	if aui.analytics.IsEnabled() {
		fmt.Printf("Status: ✅ ENABLED\n")
	} else {
		fmt.Printf("Status: ❌ DISABLED\n")
	}

	sessionStats := aui.analytics.GetSessionAnalytics()
	fmt.Printf("Session ID: %s\n", sessionStats.SessionID)
	fmt.Printf("Uptime: %s\n", formatDuration(sessionStats.Duration))
	fmt.Printf("Events Tracked: %d\n", sessionStats.TotalCommands)

	commandStats := aui.analytics.GetCommandStats()
	fmt.Printf("Commands Analyzed: %d\n", len(commandStats))

	patterns := aui.analytics.GetWorkflowPatterns()
	fmt.Printf("Patterns Detected: %d\n", len(patterns))

	cmdFormatter := components.NewCommandFormatter()
	fmt.Printf("\n💡 Use %s for session overview\n", cmdFormatter.FormatCommandInBackticks("engx analytics summary"))
	fmt.Printf("💡 Use %s for full analysis\n", cmdFormatter.FormatCommandInBackticks("engx analytics details"))

	return nil
}

// Helper functions

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

