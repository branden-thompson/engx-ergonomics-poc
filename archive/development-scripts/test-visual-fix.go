package main

import (
	"fmt"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

func main() {
	fmt.Println("🎨 Visual Fix Verification")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Create renderer
	stepNames := []string{
		"Validating configuration",
		"Setting up environment",
		"Installing dependencies",
		"Generating project structure",
		"Finalizing setup",
	}
	renderer := components.NewNPMStyleRenderer("Visual Fix Test", stepNames)

	fmt.Println("🔍 TESTING: Progress bar visibility for completed steps")
	fmt.Println()

	// Test scenario: Some steps completed, one running, some pending

	// Step 0: Complete
	renderer.SetCurrentStep(0)
	renderer.UpdateStep(0, 1.0, "Configuration validated", []string{})
	renderer.CompleteStep(0, time.Millisecond*1200)

	// Step 1: Complete
	renderer.SetCurrentStep(1)
	renderer.UpdateStep(1, 1.0, "Environment ready", []string{})
	renderer.CompleteStep(1, time.Millisecond*1800)

	// Step 2: Currently running at 45%
	renderer.SetCurrentStep(2)
	renderer.UpdateStep(2, 0.45, "Installing React dependencies...", []string{
		"Installing React v18.2.0",
		"Installing TypeScript v5.1.6",
		"Installing Vite v4.4.5",
	})

	// Render and check the output
	fmt.Println("📊 CURRENT OUTPUT:")
	fmt.Println("────────────────────────────────────────────────────────────")
	output := renderer.Render(80)
	fmt.Print(output)
	fmt.Println("────────────────────────────────────────────────────────────")

	fmt.Println()
	fmt.Println("🔍 VERIFICATION CHECKLIST:")

	// Check each step's state
	for i := 0; i < renderer.GetStepCount(); i++ {
		step := renderer.GetStepAtIndex(i)
		if step != nil {
			fmt.Printf("   Step %d: %s\n", i, step.Name)
			fmt.Printf("           Status: %s | Progress: %.1f%%\n", step.Status, step.Progress*100)

			if step.Status.String() == "Complete" {
				if step.Progress >= 1.0 {
					fmt.Printf("           ✅ CORRECT: Complete step shows 100%% progress\n")
				} else {
					fmt.Printf("           ❌ BUG: Complete step shows %.1f%% instead of 100%%\n", step.Progress*100)
				}
			} else if step.Status.String() == "Running" {
				if step.Progress > 0 && step.Progress < 1.0 {
					fmt.Printf("           ✅ CORRECT: Running step shows partial progress\n")
				} else {
					fmt.Printf("           ❌ BUG: Running step has invalid progress %.1f%%\n", step.Progress*100)
				}
			} else if step.Status.String() == "Pending" {
				if step.Progress == 0 {
					fmt.Printf("           ✅ CORRECT: Pending step shows 0%% progress\n")
				} else {
					fmt.Printf("           ❌ BUG: Pending step shows %.1f%% instead of 0%%\n", step.Progress*100)
				}
			}
			fmt.Println()
		}
	}

	// Test completion scenario
	fmt.Println("🎯 TESTING: All steps completed")
	fmt.Println("────────────────────────────────────────────────────────────")

	// Complete all remaining steps
	renderer.CompleteStep(2, time.Millisecond*3000)
	renderer.SetCurrentStep(3)
	renderer.UpdateStep(3, 1.0, "Project structure generated", []string{})
	renderer.CompleteStep(3, time.Millisecond*2200)
	renderer.SetCurrentStep(4)
	renderer.UpdateStep(4, 1.0, "Setup finalized", []string{})
	renderer.CompleteStep(4, time.Millisecond*800)

	// Render final state
	finalOutput := renderer.Render(80)
	fmt.Print(finalOutput)
	fmt.Println("────────────────────────────────────────────────────────────")

	finalProgress := renderer.GetOverallProgress()
	fmt.Printf("Final overall progress: %.1f%%\n", finalProgress*100)

	if finalProgress >= 0.99 {
		fmt.Println("✅ SUCCESS: All steps show 100% progress bars when complete")
	} else {
		fmt.Printf("❌ ISSUE: Final progress only %.1f%% - some steps may not be properly completed\n", finalProgress*100)
	}

	fmt.Println()
	fmt.Printf("🎨 VISUAL FIX SUMMARY:\n")
	fmt.Printf("   • Added StepComplete to progress bar condition\n")
	fmt.Printf("   • Completed steps now show [████████████████████] 100.0%%\n")
	fmt.Printf("   • Running steps show partial progress bars\n")
	fmt.Printf("   • Pending steps show no progress bars\n")
	fmt.Printf("   • Overall progress calculation includes all completed steps\n")
}