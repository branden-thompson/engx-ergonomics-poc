package main

import (
	"fmt"
	"time"

	progresssim "github.com/bthompso/engx-ergonomics-poc/internal/simulation/progress"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

func main() {
	fmt.Println("🔢 Test Step Numbering Logic")
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
	renderer := components.NewNPMStyleRenderer("Step Numbering Test", stepNames)

	// Start the tracker
	tracker.Start()
	totalSteps := tracker.TotalSteps()

	fmt.Printf("Testing step numbering for %d steps\n", totalSteps)
	fmt.Println()

	// Manually advance through steps to test the numbering logic
	for stepCycle := 0; stepCycle <= totalSteps; stepCycle++ {
		fmt.Printf("🔍 STEP CYCLE %d:\n", stepCycle)

		currentStepIndex := tracker.CurrentStep()
		stepInfo := tracker.CurrentStepInfo()

		fmt.Printf("   Before NextStep() - CurrentStep(): %d\n", currentStepIndex)
		if stepInfo != nil {
			fmt.Printf("   Before NextStep() - Step name: %s\n", stepInfo.Name)
		} else {
			fmt.Printf("   Before NextStep() - Step info: nil (completed)\n")
		}

		if stepCycle < totalSteps {
			// Advance to next step
			canAdvance := tracker.NextStep()
			fmt.Printf("   NextStep() returned: %v\n", canAdvance)

			if canAdvance {
				// Get step info after advancement
				newCurrentStep := tracker.CurrentStep()
				newStepInfo := tracker.CurrentStepInfo()

				fmt.Printf("   After NextStep() - CurrentStep(): %d\n", newCurrentStep)
				if newStepInfo != nil {
					fmt.Printf("   After NextStep() - Step name: %s\n", newStepInfo.Name)
				}

				// OLD LOGIC (WRONG): Step: m.tracker.CurrentStep() + 1
				oldProgressMsgStep := tracker.CurrentStep() + 1
				fmt.Printf("   OLD logic - ProgressMsg.Step would be: %d (WRONG)\n", oldProgressMsgStep)

				// NEW LOGIC (FIXED): Step: m.tracker.CurrentStep()
				newProgressMsgStep := tracker.CurrentStep()
				fmt.Printf("   NEW logic - ProgressMsg.Step would be: %d (CORRECT)\n", newProgressMsgStep)

				// Test what happens in ProgressMsg handler
				fmt.Printf("   ProgressMsg handler would:\n")
				fmt.Printf("     - Set m.currentStep = %d\n", newProgressMsgStep)

				// Mark previous step as complete (if step > 0)
				if newProgressMsgStep > 0 {
					prevStep := newProgressMsgStep - 1
					fmt.Printf("     - CompleteStep(%d) - marking step %d complete\n", prevStep, prevStep)
					renderer.CompleteStep(prevStep, time.Millisecond*1500)
				}

				// Set current step in renderer
				if newProgressMsgStep < totalSteps {
					fmt.Printf("     - SetCurrentStep(%d) - setting renderer current to %d\n", newProgressMsgStep, newProgressMsgStep)
					renderer.SetCurrentStep(newProgressMsgStep)
				}

				// Check completion condition
				if newProgressMsgStep >= totalSteps {
					fmt.Printf("     - COMPLETION: step %d >= %d (would set StateComplete)\n", newProgressMsgStep, totalSteps)
					renderer.CompleteStep(newProgressMsgStep-1, time.Millisecond*1500)
				}

				overallProgress := renderer.GetOverallProgress()
				fmt.Printf("   Current overall progress: %.1f%%\n", overallProgress*100)
			} else {
				fmt.Printf("   ✅ ALL STEPS COMPLETE!\n")

				// Simulate the completion ProgressMsg
				completionStep := totalSteps
				fmt.Printf("   Completion ProgressMsg.Step would be: %d\n", completionStep)
				fmt.Printf("   ProgressMsg handler would:\n")
				fmt.Printf("     - Set m.currentStep = %d\n", completionStep)
				fmt.Printf("     - Check: %d >= %d? %v (completion condition)\n", completionStep, totalSteps, completionStep >= totalSteps)

				if completionStep >= totalSteps {
					finalStepIndex := completionStep - 1
					fmt.Printf("     - CompleteStep(%d) - marking final step complete\n", finalStepIndex)
					renderer.CompleteStep(finalStepIndex, time.Millisecond*1500)
					fmt.Printf("     - Set state to StateComplete\n")
				}

				break
			}
		}

		fmt.Println()
	}

	fmt.Printf("🎯 FINAL RESULTS:\n")
	fmt.Printf("   Tracker completed: %v\n", tracker.IsCompleted())
	fmt.Printf("   Final progress: %.1f%%\n", renderer.GetOverallProgress()*100)

	// Validate step states
	fmt.Printf("   Step states in renderer:\n")
	for i := 0; i < renderer.GetStepCount(); i++ {
		step := renderer.GetStepAtIndex(i)
		if step != nil {
			fmt.Printf("     Step %d: %s - Progress: %.1f%% - Status: %s\n",
				i, step.Name, step.Progress*100, step.Status)
		}
	}
}