package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

func main() {
	fmt.Println("🛩️ NPM/Yarn Style UI Redesign Demo")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Create the npm-style renderer
	stepNames := []string{
		"Validating configuration",
		"Setting up environment",
		"Installing dependencies",
		"Generating project structure",
		"Configuring production setup",
		"Finalizing setup",
	}

	renderer := components.NewNPMStyleRenderer("Creating React application 'SampleApp'", stepNames)

	// Simulate progress through all steps
	fmt.Println("📊 SIMULATING NPM/YARN STYLE PROGRESS:")
	fmt.Println()

	// Step 1: Show initial state
	output := renderer.Render(80)
	fmt.Println(output)
	fmt.Println(
		"Key Features:\n" +
		"✅ Fixed progress bar/percentage synchronization - no more mismatch!\n" +
		"✅ Individual step progress bars for each task\n" +
		"✅ Overall command progress tracking\n" +
		"✅ Clean npm/yarn-style output without separate boxes\n" +
		"✅ Sub-steps displayed under current running step\n" +
		"✅ Step completion with duration timing\n" +
		"✅ Professional status indicators (⏳ pending, ⚡ running, ✅ complete)\n")

	// Simulate some progress
	renderer.SetCurrentStep(0)
	renderer.UpdateStep(0, 0.5, "Checking project name validity", []string{
		"Checking project name validity",
		"Verifying directory permissions",
	})

	fmt.Println("\n" + strings.Repeat("─", 80))
	fmt.Println("Progress Example (Step 1 at 50%):")
	fmt.Println(strings.Repeat("─", 80))
	output = renderer.Render(80)
	fmt.Println(output)

	// Complete step 1, move to step 2
	renderer.CompleteStep(0, time.Millisecond*1200)
	renderer.SetCurrentStep(1)
	renderer.UpdateStep(1, 0.3, "Creating project directory structure", []string{
		"Creating project directory structure",
		"Initializing Git repository",
		"Setting up .gitignore",
	})

	fmt.Println("\n" + strings.Repeat("─", 80))
	fmt.Println("Progress Example (Step 1 complete, Step 2 at 30%):")
	fmt.Println(strings.Repeat("─", 80))
	output = renderer.Render(80)
	fmt.Println(output)

	fmt.Println()
	fmt.Printf("🎯 IMPROVEMENTS ACHIEVED:\n")
	fmt.Printf("   • No more progress bar/percentage mismatch\n")
	fmt.Printf("   • Removed confusing 4-box layout\n")
	fmt.Printf("   • Added individual step progress bars\n")
	fmt.Printf("   • Overall progress calculation works correctly\n")
	fmt.Printf("   • Sub-steps shown for current running step\n")
	fmt.Printf("   • Clean, professional npm/yarn-style output\n")
	fmt.Printf("   • Step completion timing displayed\n")
}