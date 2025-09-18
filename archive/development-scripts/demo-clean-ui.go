package main

import (
	"fmt"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

func main() {
	fmt.Println("🧹 Clean UI Demo - Redundant Info Removed")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Create renderer with simplified constructor
	stepNames := []string{
		"Validating Configuration",
		"Setting up Git Repo",
		"Setting up Environment",
		"Installing Dependencies",
		"Generating Project Structure",
		"Configuring Production Setup",
		"Finalizing Setup",
	}

	renderer := components.NewPolishedRenderer("SampleApp", stepNames)

	fmt.Println("🔍 CLEAN UI LAYOUT:")
	fmt.Println()

	// Complete some steps and show running state
	for i := 0; i < 4; i++ {
		renderer.SetCurrentStep(i)
		renderer.UpdateStep(i, 1.0, "Completed", []string{})
		renderer.CompleteStep(i, time.Millisecond*1500)
	}

	// Set current step with partial progress
	renderer.SetCurrentStep(5) // Production setup
	renderer.UpdateStep(5, 0.763, "Optimizing build configuration...", []string{})

	output := renderer.Render(80)
	fmt.Print(output)
	fmt.Println()
	fmt.Println("🔍 COMPLETION STATE:")

	// Complete all steps
	for i := 0; i < len(stepNames); i++ {
		renderer.SetCurrentStep(i)
		renderer.UpdateStep(i, 1.0, "Completed", []string{})
		renderer.CompleteStep(i, time.Millisecond*1500)
	}

	renderer.SetCurrentStep(len(stepNames) - 1)
	completedOutput := renderer.Render(80)
	fmt.Print(completedOutput)

}