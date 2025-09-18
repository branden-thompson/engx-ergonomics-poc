package main

import (
	"fmt"
	"time"

	progresssim "github.com/bthompso/engx-ergonomics-poc/internal/simulation/progress"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

// Simulate the TUI message types
type MsgType string

const (
	ProgressTickMsg MsgType = "ProgressTickMsg"
	StepCheckMsg    MsgType = "StepCheckMsg"
	ProgressMsg     MsgType = "ProgressMsg"
)

type SimulatedMessage struct {
	Type     MsgType
	Step     int
	StepName string
	Message  string
}

func main() {
	fmt.Println("🔍 Diagnose TUI Flow - Exact Simulation")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Create tracker and renderer (dev-only for faster testing)
	tracker := progresssim.NewCreateTracker(true)
	stepNames := []string{
		"Validating configuration",
		"Setting up environment",
		"Installing dependencies",
		"Generating project structure",
		"Finalizing setup",
	}
	renderer := components.NewNPMStyleRenderer("TUI Flow Diagnosis", stepNames)

	// Start tracker
	tracker.Start()
	totalSteps := tracker.TotalSteps()

	fmt.Printf("Total steps: %d\n", totalSteps)
	fmt.Printf("Simulating TUI message flow with exact timing...\n")
	fmt.Println()

	// Simulate TUI state
	state := "StateExecuting"

	// Track message timing like the real TUI
	lastProgressTick := time.Now()
	lastStepCheck := time.Now()

	// Simulation loop
	maxCycles := 2000 // Safety limit - increased to see final step completion
	for cycle := 0; cycle < maxCycles; cycle++ {
		now := time.Now()

		// Simulate messages based on timing (like real TUI)
		var msg SimulatedMessage

		// ProgressTickMsg every 50ms
		if now.Sub(lastProgressTick) >= time.Millisecond*50 {
			msg = SimulatedMessage{Type: ProgressTickMsg}
			lastProgressTick = now
		} else if now.Sub(lastStepCheck) >= time.Millisecond*200 {
			// StepCheckMsg every 200ms
			msg = SimulatedMessage{Type: StepCheckMsg}
			lastStepCheck = now
		} else {
			// Small delay to prevent busy loop
			time.Sleep(time.Millisecond * 10)
			continue
		}

		// Process message like the real TUI Update() method
		switch msg.Type {
		case ProgressTickMsg:
			// Like ProgressTickMsg case in app.go
			if state == "StateExecuting" && tracker != nil && !tracker.IsCompleted() {
				currentStepIndex := tracker.CurrentStep()
				stepInfo := tracker.CurrentStepInfo()

				if stepInfo != nil && currentStepIndex >= 0 {
					// Calculate step progress like real TUI
					stepStart := tracker.GetStepStart()
					elapsed := time.Since(stepStart)
					stepProgress := float64(elapsed) / float64(stepInfo.Duration)

					if stepProgress > 1.0 {
						stepProgress = 1.0
					}
					if stepProgress < 0 {
						stepProgress = 0.0
					}

					// Update renderer like real TUI
					renderer.SetCurrentStep(currentStepIndex)
					renderer.UpdateStep(currentStepIndex, stepProgress, stepInfo.Message, []string{
						fmt.Sprintf("Processing %s...", stepInfo.Name),
					})

					// Log progress periodically
					if cycle%40 == 0 { // Every ~2 seconds at 20fps
						overallProgress := renderer.GetOverallProgress()
						fmt.Printf("[Cycle %4d] ProgressTick - Step %d: %.1f%% | Overall: %.1f%%\n",
							cycle, currentStepIndex, stepProgress*100, overallProgress*100)
					}
				}
			}

		case StepCheckMsg:
			// Like StepCheckMsg case in app.go
			if state == "StateExecuting" && tracker != nil && !tracker.IsCompleted() {
				// This triggers nextStep() which checks IsStepReady()
				if tracker.IsStepReady() {
					// Like nextStep() logic in app.go
					if !tracker.NextStep() {
						// All steps complete
						msg = SimulatedMessage{
							Type:     ProgressMsg,
							Step:     tracker.TotalSteps(),
							StepName: "Complete",
							Message:  "✨ All steps completed successfully!",
						}
						fmt.Printf("[Cycle %4d] StepCheck -> COMPLETION ProgressMsg (Step %d)\n",
							cycle, msg.Step)
					} else {
						// Step advanced
						stepInfo := tracker.CurrentStepInfo()
						if stepInfo != nil {
							msg = SimulatedMessage{
								Type:     ProgressMsg,
								Step:     tracker.CurrentStep(), // FIXED: no +1
								StepName: stepInfo.Name,
								Message:  stepInfo.Message,
							}
							fmt.Printf("[Cycle %4d] StepCheck -> ADVANCE ProgressMsg (Step %d: %s)\n",
								cycle, msg.Step, stepInfo.Name)
						}
					}

					// Process the generated ProgressMsg immediately
					// Like ProgressMsg case in app.go
					if msg.Step > 0 && renderer != nil {
						renderer.CompleteStep(msg.Step-1, time.Since(tracker.GetStepStart()))
						fmt.Printf("               -> Completed step %d\n", msg.Step-1)
					}

					if renderer != nil && msg.Step < totalSteps {
						renderer.SetCurrentStep(msg.Step)
						fmt.Printf("               -> Set current step to %d\n", msg.Step)
					}

					if msg.Step >= totalSteps {
						state = "StateComplete"
						if renderer != nil {
							renderer.CompleteStep(msg.Step-1, time.Since(tracker.GetStepStart()))
							fmt.Printf("               -> FINAL step %d completed\n", msg.Step-1)
						}
						fmt.Printf("               -> State set to COMPLETE\n")
						break
					} else {
						state = "StateExecuting"
					}

					overallProgress := renderer.GetOverallProgress()
					fmt.Printf("               -> Overall progress: %.1f%%\n", overallProgress*100)
				}
			}
		}

		// Break if completed
		if tracker.IsCompleted() && state == "StateComplete" {
			fmt.Printf("[Cycle %4d] SIMULATION COMPLETE!\n", cycle)
			break
		}

		// Small delay to prevent overwhelming output
		if cycle%10 == 0 {
			time.Sleep(time.Millisecond * 1)
		}
	}

	// Final results
	fmt.Println()
	fmt.Printf("🎯 FINAL DIAGNOSIS:\n")
	fmt.Printf("   Tracker completed: %v\n", tracker.IsCompleted())
	fmt.Printf("   State: %s\n", state)
	fmt.Printf("   Final progress: %.1f%%\n", renderer.GetOverallProgress()*100)

	if renderer.GetOverallProgress() >= 0.99 {
		fmt.Printf("   ✅ SUCCESS: Progress reaches 100%%\n")
	} else {
		fmt.Printf("   ❌ ISSUE: Progress stuck at %.1f%%\n", renderer.GetOverallProgress()*100)
		fmt.Printf("   Step details:\n")
		for i := 0; i < renderer.GetStepCount(); i++ {
			step := renderer.GetStepAtIndex(i)
			if step != nil {
				fmt.Printf("     Step %d: %.1f%% - %s\n", i, step.Progress*100, step.Status)
			}
		}
	}
}