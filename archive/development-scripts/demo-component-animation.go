package main

import (
	"fmt"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

func main() {
	fmt.Println("🎬 Component Animation Demo")
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

	renderer := components.NewEnhancedRenderer(
		"TestApp",
		"./my-test-app",
		"TypeScript",
		stepNames,
		false,
	)

	fmt.Println("🔍 TESTING: Component Animation During Dependencies Step")
	fmt.Println()

	// Set up initial state
	renderer.SetCurrentStep(0)
	renderer.UpdateStep(0, 1.0, "Completed", []string{})
	renderer.CompleteStep(0, time.Millisecond*1500)

	renderer.SetCurrentStep(1)
	renderer.UpdateStep(1, 1.0, "Completed", []string{})
	renderer.CompleteStep(1, time.Millisecond*1500)

	// Now start the dependencies step and show component animation
	renderer.SetCurrentStep(2)

	// Show progression at different points
	progressSteps := []float64{0.0, 0.3, 0.6, 0.9, 1.0}

	for _, progress := range progressSteps {
		fmt.Printf("📦 Installing Dependencies... %.1f%%\n", progress*100)
		fmt.Println("---")

		renderer.UpdateStep(2, progress, "Installing dependencies...", []string{})
		renderer.UpdateComponentStatuses("Installing dependencies", progress)

		output := renderer.Render(89)
		fmt.Print(output)
		fmt.Println()
	}

	fmt.Println("✨ COMPONENT ANIMATION TEST COMPLETE")
}