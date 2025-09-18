package main

import (
	"fmt"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

func main() {
	fmt.Println("🏁 Test Race Condition Fix")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Create renderer
	stepNames := []string{
		"Validating configuration",
		"Setting up environment",
		"Installing dependencies",
	}
	renderer := components.NewNPMStyleRenderer("Race Condition Test", stepNames)

	fmt.Println("🔍 SIMULATING TUI RACE CONDITION:")
	fmt.Println()

	// Step 1: Start step 0 with partial progress
	fmt.Println("1. UpdateStep(0, 45%) - like ProgressTickMsg")
	renderer.SetCurrentStep(0)
	renderer.UpdateStep(0, 0.45, "Checking configuration...", []string{})

	step0 := renderer.GetStepAtIndex(0)
	fmt.Printf("   Step 0: Status=%s, Progress=%.1f%%\n", step0.Status, step0.Progress*100)

	// Step 2: Complete step 0 - like ProgressMsg
	fmt.Println("2. CompleteStep(0) - like ProgressMsg when step advances")
	renderer.CompleteStep(0, time.Millisecond*1200)

	step0After := renderer.GetStepAtIndex(0)
	fmt.Printf("   Step 0: Status=%s, Progress=%.1f%%\n", step0After.Status, step0After.Progress*100)

	// Step 3: Try to update step 0 again - like ProgressTickMsg continuing to run
	fmt.Println("3. UpdateStep(0, 60%) - RACE CONDITION: ProgressTickMsg after completion")
	renderer.UpdateStep(0, 0.60, "Should NOT overwrite", []string{})

	step0Final := renderer.GetStepAtIndex(0)
	fmt.Printf("   Step 0: Status=%s, Progress=%.1f%%\n", step0Final.Status, step0Final.Progress*100)

	if step0Final.Progress >= 1.0 && step0Final.Status.String() == "Complete" {
		fmt.Printf("   ✅ SUCCESS: Completed step protected from overwrites\n")
	} else {
		fmt.Printf("   ❌ FAILED: Completed step was overwritten!\n")
	}

	fmt.Println()

	// Test visual output
	fmt.Println("4. Visual output:")
	output := renderer.Render(80)
	fmt.Printf("%s\n", output)

	fmt.Println("🎯 FIX VERIFICATION:")
	fmt.Printf("   • CompleteStep() sets Status=Complete and Progress=100%%\n")
	fmt.Printf("   • UpdateStep() checks Status=Complete and returns early\n")
	fmt.Printf("   • Race condition between ProgressTickMsg and ProgressMsg resolved\n")
	fmt.Printf("   • Completed steps maintain 100%% progress bars\n")
}