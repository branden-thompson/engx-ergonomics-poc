package main

import (
	"fmt"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

func main() {
	fmt.Println("🎬 Full Component Animation Demo")
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

	// Complete first two steps
	renderer.SetCurrentStep(0)
	renderer.UpdateStep(0, 1.0, "Completed", []string{})
	renderer.CompleteStep(0, time.Millisecond*1500)
	renderer.UpdateComponentStatuses("Validating configuration", 1.0)

	renderer.SetCurrentStep(1)
	renderer.UpdateStep(1, 1.0, "Completed", []string{})
	renderer.CompleteStep(1, time.Millisecond*1500)
	renderer.UpdateComponentStatuses("Setting up environment", 1.0)

	fmt.Println("📦 Step 3: Installing Dependencies")
	fmt.Println("---")
	renderer.SetCurrentStep(2)
	renderer.UpdateStep(2, 1.0, "Installing dependencies...", []string{})
	renderer.UpdateComponentStatuses("Installing dependencies", 1.0)
	output := renderer.Render(89)
	fmt.Print(output)
	fmt.Println()

	fmt.Println("📦 Step 4: Generating Project Structure")
	fmt.Println("---")
	renderer.SetCurrentStep(3)
	renderer.UpdateStep(3, 1.0, "Generating project structure...", []string{})
	renderer.UpdateComponentStatuses("Generating project structure", 1.0)
	output = renderer.Render(89)
	fmt.Print(output)
	fmt.Println()

	fmt.Println("📦 Step 5: Finalizing Setup")
	fmt.Println("---")
	renderer.SetCurrentStep(4)
	renderer.UpdateStep(4, 1.0, "Finalizing setup...", []string{})
	renderer.UpdateComponentStatuses("Finalizing setup", 1.0)
	output = renderer.Render(89)
	fmt.Print(output)
	fmt.Println()

	fmt.Println("✨ FULL COMPONENT ANIMATION TEST COMPLETE")
}