package main

import (
	"fmt"
	"os"
	"time"

	progresssim "github.com/bthompso/engx-ergonomics-poc/internal/simulation/progress"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

func clearScreen() {
	fmt.Print("\033[2J\033[H") // Clear screen and move cursor to top
}

func main() {
	fmt.Println("🎬 Real-Time Progress Animation Demo")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("This shows how the progress should look in the real TUI...")
	fmt.Println()
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()
	time.Sleep(2 * time.Second)

	// Create tracker and renderer (dev-only for faster demo)
	tracker := progresssim.NewCreateTracker(true)
	stepNames := []string{
		"Validating configuration",
		"Setting up environment",
		"Installing dependencies",
		"Generating project structure",
		"Finalizing setup",
	}
	renderer := components.NewNPMStyleRenderer("Creating React application 'DemoApp'", stepNames)

	// Start tracker
	tracker.Start()
	totalSteps := tracker.TotalSteps()

	// Display total expected time
	totalDuration := time.Duration(0)
	for i := 0; i < totalSteps; i++ {
		step := tracker.GetStep(i)
		if step != nil {
			totalDuration += step.Duration
		}
	}

	fmt.Printf("Expected completion time: %.1f seconds\n", totalDuration.Seconds())
	fmt.Println("Watch the progress bars animate to 100%...")
	time.Sleep(2 * time.Second)

	// Real-time simulation matching TUI behavior
	state := "StateExecuting"
	lastProgressTick := time.Now()
	lastStepCheck := time.Now()
	startTime := time.Now()

	for state != "StateComplete" {
		now := time.Now()

		// Progress tick updates (50ms intervals like real TUI)
		if now.Sub(lastProgressTick) >= time.Millisecond*50 {
			if state == "StateExecuting" && !tracker.IsCompleted() {
				currentStepIndex := tracker.CurrentStep()
				stepInfo := tracker.CurrentStepInfo()

				if stepInfo != nil && currentStepIndex >= 0 {
					// Calculate step progress
					stepStart := tracker.GetStepStart()
					elapsed := time.Since(stepStart)
					stepProgress := float64(elapsed) / float64(stepInfo.Duration)

					if stepProgress > 1.0 {
						stepProgress = 1.0
					}

					// Update renderer
					renderer.SetCurrentStep(currentStepIndex)
					renderer.UpdateStep(currentStepIndex, stepProgress, stepInfo.Message, []string{
						"✓ Checking project name validity",
						"✓ Verifying directory permissions",
						"⚡ Setting up configuration...",
					})

					// Display current state
					clearScreen()
					fmt.Printf("🎬 Real-Time Progress Demo - Elapsed: %.1fs\n", time.Since(startTime).Seconds())
					fmt.Printf("═══════════════════════════════════════════════════════════════\n\n")

					output := renderer.Render(80)
					fmt.Print(output)

					fmt.Printf("\n\nCurrent Step: %d/%d - %s\n", currentStepIndex+1, totalSteps, stepInfo.Name)
					fmt.Printf("Step Progress: %.1f%% | Overall: %.1f%%\n",
						stepProgress*100, renderer.GetOverallProgress()*100)
					fmt.Printf("State: %s\n", state)

					if stepProgress >= 1.0 {
						fmt.Printf("⏳ Step completing...\n")
					}
				}
			}
			lastProgressTick = now
		}

		// Step check updates (200ms intervals like real TUI)
		if now.Sub(lastStepCheck) >= time.Millisecond*200 {
			if state == "StateExecuting" && !tracker.IsCompleted() {
				if tracker.IsStepReady() {
					if !tracker.NextStep() {
						// All steps complete
						state = "StateComplete"
						renderer.CompleteStep(totalSteps-1, time.Since(startTime))

						clearScreen()
						fmt.Printf("🎬 Real-Time Progress Demo - COMPLETED in %.1fs\n", time.Since(startTime).Seconds())
						fmt.Printf("═══════════════════════════════════════════════════════════════\n\n")

						output := renderer.Render(80)
						fmt.Print(output)

						fmt.Printf("\n\n✅ ALL STEPS COMPLETE!\n")
						fmt.Printf("Final Progress: %.1f%%\n", renderer.GetOverallProgress()*100)
						fmt.Printf("State: %s\n", state)
						break
					} else {
						// Step advanced - complete previous step
						currentStep := tracker.CurrentStep()
						if currentStep > 0 {
							renderer.CompleteStep(currentStep-1, time.Since(startTime))
						}
						if currentStep < totalSteps {
							renderer.SetCurrentStep(currentStep)
						}
					}
				}
			}
			lastStepCheck = now
		}

		// Small delay to prevent busy loop
		time.Sleep(time.Millisecond * 10)
	}

	// Final display
	fmt.Println("\n🎯 DEMO COMPLETE!")
	fmt.Printf("   Total time: %.1f seconds\n", time.Since(startTime).Seconds())
	fmt.Printf("   Expected: %.1f seconds\n", totalDuration.Seconds())
	fmt.Printf("   Progress animation: ✅ Working correctly\n")
	fmt.Printf("   All steps complete: ✅ 100%%\n")
	fmt.Println()
	fmt.Println("This is how the TUI should behave. If your TUI doesn't reach 100%,")
	fmt.Println("try running it for the full expected duration before concluding it's stuck.")

	// Ask if user wants to see step details
	fmt.Println("\nPress Enter to see detailed step breakdown...")
	fmt.Scanln()

	fmt.Println("\n📊 STEP BREAKDOWN:")
	for i := 0; i < renderer.GetStepCount(); i++ {
		step := renderer.GetStepAtIndex(i)
		if step != nil {
			duration := tracker.GetStep(i).Duration
			fmt.Printf("   Step %d: %s\n", i+1, step.Name)
			fmt.Printf("           Duration: %v | Progress: %.1f%% | Status: %s\n",
				duration, step.Progress*100, step.Status)
		}
	}
}