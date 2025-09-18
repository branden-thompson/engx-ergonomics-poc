package main

import (
	"fmt"
	"strings"
	"time"

	progresssim "github.com/bthompso/engx-ergonomics-poc/internal/simulation/progress"
)

func main() {
	fmt.Println("🛩️ DPX Web Progress Animation Demo")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Create tracker like the TUI would
	tracker := progresssim.NewCreateTracker(false)
	tracker.Start()

	// Simulate what the TUI would show with continuous updates
	lastProgress := -1.0
	updateCount := 0

	for !tracker.IsCompleted() {
		currentStep := tracker.CurrentStep()
		stepInfo := tracker.CurrentStepInfo()
		progress := tracker.Progress()
		eta := tracker.EstimatedTimeRemaining()

		// Only show updates when progress changes (like TUI redraws)
		progressPercent := int(progress * 100)
		if progressPercent != int(lastProgress*100) {
			updateCount++

			// Clear previous line (simulate TUI redraw)
			if updateCount > 1 {
				fmt.Print("\033[2K\r") // Clear line and return to start
			}

			// Show what the TUI progress bar would look like
			progressBar := renderProgressBar(progress, 40)

			if stepInfo != nil {
				fmt.Printf("Step %d/%d: %s\n",
					currentStep+1, tracker.TotalSteps(), stepInfo.Name)
				fmt.Printf("%s %.1f%%\n", progressBar, progress*100)
				fmt.Printf("ETA: %v\n", eta.Round(time.Second))
				fmt.Printf("%s\n", stepInfo.Message)
				fmt.Println(strings.Repeat("─", 50))
			}

			lastProgress = progress
		}

		// Check if step should advance
		if tracker.IsStepReady() {
			if !tracker.NextStep() {
				break
			}
		}

		time.Sleep(time.Millisecond * 50) // 20fps like TUI
	}

	fmt.Printf("\n✨ Complete! Total time: %v\n", tracker.TotalElapsed())
	fmt.Printf("📊 Updates shown: %d (simulating TUI redraws)\n", updateCount)
}

func renderProgressBar(progress float64, width int) string {
	filled := int(progress * float64(width))
	empty := width - filled

	bar := ""
	bar += "█" // Start
	bar += strings.Repeat("█", filled)
	bar += strings.Repeat("░", empty)
	bar += "█" // End

	return fmt.Sprintf("[%s]", bar[1:len(bar)-1]) // Remove the extra chars
}