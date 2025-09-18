package main

import (
	"fmt"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

func main() {
	fmt.Println("🎨 Enhanced Comprehensive Layout Demo")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Create enhanced renderer matching the template
	stepNames := []string{
		"Validating Configuration",
		"Setting Up Environment",
		"Installing Dependencies",
		"Generating Project Structure",
		"Installing Testing Frameworks",
		"Generating Documentation",
		"Finalizing Setup",
	}

	renderer := components.NewEnhancedRenderer(
		"SampleApp",
		"./my-react-app",
		"TypeScript",
		stepNames,
		false, // Production deployable
	)

	fmt.Println("🔍 TESTING: Enhanced Layout with Comprehensive Sections")
	fmt.Println()

	// Complete some steps to show the layout in action
	for i := 0; i < 2; i++ {
		renderer.SetCurrentStep(i)
		renderer.UpdateStep(i, 1.0, "Completed", []string{})
		renderer.CompleteStep(i, time.Millisecond*1500)
	}

	// Set current step with partial progress (Installing Dependencies)
	renderer.SetCurrentStep(2)
	renderer.UpdateStep(2, 0.655, "Installing React and TypeScript...", []string{})

	// Sleep for a moment to show some elapsed time
	time.Sleep(time.Millisecond * 100)

	output := renderer.Render(89) // Match template width
	fmt.Print(output)
	fmt.Println()
	fmt.Println("🔍 TESTING: Completion State")

	// Complete all steps and simulate component installation
	for i := 0; i < len(stepNames); i++ {
		renderer.SetCurrentStep(i)
		renderer.UpdateStep(i, 1.0, "Completed", []string{})
		renderer.CompleteStep(i, time.Millisecond*1500)
		// Manually trigger component status updates for demo
		renderer.UpdateComponentStatuses(stepNames[i], 1.0)
	}

	renderer.SetCurrentStep(len(stepNames) - 1)
	completedOutput := renderer.Render(89)
	fmt.Print(completedOutput)

}