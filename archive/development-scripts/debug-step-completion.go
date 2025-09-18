package main

import (
	"fmt"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

func main() {
	fmt.Println("🔍 Debug Step Completion Process")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Create renderer
	stepNames := []string{
		"Validating configuration",
		"Setting up environment",
		"Installing dependencies",
	}
	renderer := components.NewNPMStyleRenderer("Debug Test", stepNames)

	fmt.Println("📊 SIMULATING STEP COMPLETION SEQUENCE:")
	fmt.Println()

	// Step 1: Start with step 0
	fmt.Println("1. Set current step to 0, update progress to 45%")
	renderer.SetCurrentStep(0)
	renderer.UpdateStep(0, 0.45, "Checking configuration...", []string{})

	step0 := renderer.GetStepAtIndex(0)
	fmt.Printf("   Step 0 state: Status=%s, Progress=%.1f%%\n", step0.Status, step0.Progress*100)

	output := renderer.Render(80)
	fmt.Printf("   Visual: %s\n", extractStepLine(output, 0))
	fmt.Println()

	// Step 2: Complete step 0
	fmt.Println("2. Complete step 0 (this should set it to 100%)")
	renderer.CompleteStep(0, time.Millisecond*1200)

	step0After := renderer.GetStepAtIndex(0)
	fmt.Printf("   Step 0 state after CompleteStep: Status=%s, Progress=%.1f%%\n", step0After.Status, step0After.Progress*100)

	output = renderer.Render(80)
	fmt.Printf("   Visual: %s\n", extractStepLine(output, 0))
	fmt.Println()

	// Step 3: Move to step 1
	fmt.Println("3. Set current step to 1, update progress to 30%")
	renderer.SetCurrentStep(1)
	renderer.UpdateStep(1, 0.30, "Setting up environment...", []string{})

	step1 := renderer.GetStepAtIndex(1)
	fmt.Printf("   Step 1 state: Status=%s, Progress=%.1f%%\n", step1.Status, step1.Progress*100)

	// Check if step 0 is still showing correctly
	step0Check := renderer.GetStepAtIndex(0)
	fmt.Printf("   Step 0 state (should still be complete): Status=%s, Progress=%.1f%%\n", step0Check.Status, step0Check.Progress*100)

	output = renderer.Render(80)
	fmt.Printf("   Step 0 visual: %s\n", extractStepLine(output, 0))
	fmt.Printf("   Step 1 visual: %s\n", extractStepLine(output, 1))
	fmt.Println()

	// Step 4: Complete step 1
	fmt.Println("4. Complete step 1")
	renderer.CompleteStep(1, time.Millisecond*1800)

	step1After := renderer.GetStepAtIndex(1)
	fmt.Printf("   Step 1 state after CompleteStep: Status=%s, Progress=%.1f%%\n", step1After.Status, step1After.Progress*100)

	output = renderer.Render(80)
	fmt.Printf("   Step 0 visual: %s\n", extractStepLine(output, 0))
	fmt.Printf("   Step 1 visual: %s\n", extractStepLine(output, 1))
	fmt.Println()

	// Final check
	fmt.Println("🎯 FINAL STATE CHECK:")
	for i := 0; i < 3; i++ {
		step := renderer.GetStepAtIndex(i)
		if step != nil {
			fmt.Printf("   Step %d: Status=%s, Progress=%.1f%%\n", i, step.Status, step.Progress*100)
		}
	}

	fmt.Println()
	fmt.Printf("🔍 DEBUGGING CONCLUSION:\n")
	fmt.Printf("   If completed steps show Status=Complete but Progress<100%%,\n")
	fmt.Printf("   then CompleteStep() is not properly setting Progress=1.0\n")
	fmt.Printf("   \n")
	fmt.Printf("   If completed steps show Status=Complete and Progress=100%%\n")
	fmt.Printf("   but visual doesn't show full bars, then renderStep() has an issue\n")
}

func extractStepLine(output string, stepIndex int) string {
	lines := splitLines(output)
	stepLineIndex := 3 + stepIndex // Header is 3 lines, then steps start
	if stepLineIndex < len(lines) {
		return lines[stepLineIndex]
	}
	return "Line not found"
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}