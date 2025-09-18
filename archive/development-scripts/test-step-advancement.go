package main

import (
	"fmt"
	"time"

	progresssim "github.com/bthompso/engx-ergonomics-poc/internal/simulation/progress"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

func main() {
	fmt.Println("🔧 Test Step Advancement Logic")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Create tracker for development flow (5 steps)
	tracker := progresssim.NewCreateTracker(true)
	stepNames := []string{
		"Validating configuration",
		"Setting up environment",
		"Installing dependencies",
		"Generating project structure",
		"Finalizing setup",
	}
	renderer := components.NewNPMStyleRenderer("Step Advancement Test", stepNames)

	// Start the tracker
	tracker.Start()
	totalSteps := tracker.TotalSteps()

	fmt.Printf("Testing step advancement for %d steps\n", totalSteps)
	fmt.Println()

	// Simulate the nextStep() logic like the TUI does
	for step := 0; step < totalSteps+2; step++ { // +2 to test completion
		fmt.Printf("🔍 CYCLE %d:\n", step+1)

		currentStepIndex := tracker.CurrentStep()
		stepInfo := tracker.CurrentStepInfo()

		fmt.Printf("   Tracker current step: %d\n", currentStepIndex)
		if stepInfo != nil {
			fmt.Printf("   Step info: %s\n", stepInfo.Name)
		} else {
			fmt.Printf("   Step info: nil (completed)\n")
		}

		// Check if step is ready (simulate IsStepReady)
		isReady := tracker.IsStepReady()
		fmt.Printf("   Is step ready: %v\n", isReady)

		if isReady {
			// Advance to next step (simulate NextStep)
			canAdvance := tracker.NextStep()
			fmt.Printf("   NextStep() returned: %v\n", canAdvance)

			if !canAdvance {
				// All steps complete
				fmt.Printf("   ✅ ALL STEPS COMPLETE!\n")
				fmt.Printf("   Progress: %.1f%%\n", renderer.GetOverallProgress()*100)
				break
			}

			// Get new current step info after advancement
			newCurrentStep := tracker.CurrentStep()
			newStepInfo := tracker.CurrentStepInfo()

			fmt.Printf("   After NextStep - Tracker current step: %d\n", newCurrentStep)
			if newStepInfo != nil {
				fmt.Printf("   After NextStep - Step info: %s\n", newStepInfo.Name)
			}

			// This is what the fixed code should send as ProgressMsg.Step
			progressMsgStep := tracker.CurrentStep()
			fmt.Printf("   ProgressMsg.Step would be: %d\n", progressMsgStep)

			// Simulate ProgressMsg handler
			// Mark previous step as complete
			if progressMsgStep > 0 {
				renderer.CompleteStep(progressMsgStep-1, time.Millisecond*1500)
				fmt.Printf("   Completed step %d in renderer\n", progressMsgStep-1)
			}

			// Set current step in renderer
			if progressMsgStep < totalSteps {
				renderer.SetCurrentStep(progressMsgStep)
				fmt.Printf("   Set renderer current step to: %d\n", progressMsgStep)
			}

			// Check if we're at completion
			if progressMsgStep >= totalSteps {
				fmt.Printf("   ✅ STATE: Complete (step %d >= %d)\n", progressMsgStep, totalSteps)
				// Mark final step as complete
				renderer.CompleteStep(progressMsgStep-1, time.Millisecond*1500)
			}

			overallProgress := renderer.GetOverallProgress()
			fmt.Printf("   Overall progress: %.1f%%\n", overallProgress*100)
		} else {
			// Simulate time passing to make step ready
			if stepInfo != nil {
				fmt.Printf("   Waiting for step duration (%v)...\n", stepInfo.Duration)
				time.Sleep(stepInfo.Duration)
			} else {
				fmt.Printf("   Simulating time passage...\n")
				time.Sleep(time.Millisecond * 10)
			}
		}

		fmt.Println()

		// Safety check
		if tracker.IsCompleted() {
			fmt.Printf("🎉 TRACKER COMPLETED!\n")
			break
		}
	}

	fmt.Printf("🎯 FINAL RESULTS:\n")
	fmt.Printf("   Tracker completed: %v\n", tracker.IsCompleted())
	fmt.Printf("   Final progress: %.1f%%\n", renderer.GetOverallProgress()*100)
}