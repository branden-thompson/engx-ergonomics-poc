package main

import (
	"fmt"
	"time"

	progresssim "github.com/bthompso/engx-ergonomics-poc/internal/simulation/progress"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

func main() {
	fmt.Println("🔍 Debug Progress Timing Issues")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Create tracker and renderer
	tracker := progresssim.NewCreateTracker(true) // dev-only for faster testing
	stepNames := []string{
		"Validating configuration",
		"Setting up environment",
		"Installing dependencies",
		"Generating project structure",
		"Finalizing setup",
	}
	renderer := components.NewNPMStyleRenderer("Debug Test", stepNames)

	// Start the tracker
	tracker.Start()

	fmt.Printf("Total steps: %d\n", tracker.TotalSteps())
	fmt.Println()

	// Simulate what happens in the actual TUI
	for currentStepIndex := 0; currentStepIndex < tracker.TotalSteps(); currentStepIndex++ {
		stepInfo := tracker.GetStep(currentStepIndex)
		if stepInfo == nil {
			break
		}

		fmt.Printf("🔍 DEBUGGING STEP %d: %s\n", currentStepIndex+1, stepInfo.Name)
		fmt.Printf("   Expected duration: %v\n", stepInfo.Duration)

		// Simulate the step start
		tracker.NextStep() // This sets stepStart time in tracker

		// Set current step in renderer
		renderer.SetCurrentStep(currentStepIndex)

		// Simulate progress over time (like ProgressTickMsg does)
		simulationDuration := stepInfo.Duration + time.Millisecond*100 // Add 100ms buffer

		for elapsed := time.Duration(0); elapsed <= simulationDuration; elapsed += time.Millisecond * 100 {
			// Calculate step progress like the actual app does
			stepProgress := float64(elapsed) / float64(stepInfo.Duration)

			if stepProgress > 1.0 {
				stepProgress = 1.0
			}
			if stepProgress < 0 {
				stepProgress = 0.0
			}

			// Update renderer
			renderer.UpdateStep(currentStepIndex, stepProgress, stepInfo.Message, []string{})

			// Check overall progress
			overallProgress := renderer.GetOverallProgress()

			fmt.Printf("   Elapsed: %6v | Step Progress: %5.1f%% | Overall: %5.1f%%\n",
				elapsed, stepProgress*100, overallProgress*100)

			// If step is complete, break
			if stepProgress >= 1.0 {
				break
			}

			time.Sleep(time.Millisecond * 100)
		}

		// Mark step complete (like ProgressMsg handler does)
		renderer.CompleteStep(currentStepIndex, stepInfo.Duration)

		finalProgress := renderer.GetOverallProgress()
		fmt.Printf("   ✅ Step completed - Final overall progress: %.1f%%\n", finalProgress*100)
		fmt.Println()
	}

	// Final check
	finalOverallProgress := renderer.GetOverallProgress()
	fmt.Printf("🎯 FINAL RESULTS:\n")
	fmt.Printf("   Final overall progress: %.1f%%\n", finalOverallProgress*100)

	if finalOverallProgress >= 0.99 {
		fmt.Printf("   ✅ Progress calculation working correctly\n")
	} else {
		fmt.Printf("   ❌ Progress calculation issue detected\n")
		fmt.Printf("   Individual step progress:\n")
		for i := 0; i < renderer.GetStepCount(); i++ {
			step := renderer.GetStepAtIndex(i)
			if step != nil {
				fmt.Printf("     Step %d: %.1f%% (%s)\n", i+1, step.Progress*100, step.Status)
			}
		}
	}
}