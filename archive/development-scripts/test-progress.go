package main

import (
	"fmt"
	"time"

	progresssim "github.com/bthompso/engx-ergonomics-poc/internal/simulation/progress"
)

func main() {
	fmt.Println("🛩️ Testing Progress Simulation Engine")

	// Create a tracker for the create command
	tracker := progresssim.NewCreateTracker(false)
	tracker.Start()

	fmt.Printf("Total steps: %d\n", tracker.TotalSteps())

	// Simulate the progress over time
	for !tracker.IsCompleted() {
		currentStep := tracker.CurrentStep()
		stepInfo := tracker.CurrentStepInfo()
		progress := tracker.Progress()
		eta := tracker.EstimatedTimeRemaining()

		if stepInfo != nil {
			fmt.Printf("Step %d/%d: %s (%.1f%% complete, ETA: %v)\n",
				currentStep+1, tracker.TotalSteps(), stepInfo.Name, progress*100, eta.Round(time.Second))
		}

		// Check if step should advance
		if tracker.IsStepReady() {
			fmt.Printf("  ✅ Step completed: %s\n", stepInfo.Message)
			if !tracker.NextStep() {
				break
			}
		}

		time.Sleep(time.Millisecond * 100) // Simulate UI refresh rate
	}

	fmt.Println("✨ All steps completed!")
	fmt.Printf("Total time: %v\n", tracker.TotalElapsed())
}