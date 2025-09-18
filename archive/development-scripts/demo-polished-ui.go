package main

import (
	"fmt"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
)

func main() {
	fmt.Println("🎨 Polished UI Demo")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Create renderer with your example steps
	stepNames := []string{
		"Validating Configuration",
		"Setting up Git Repo",
		"Setting up Environment",
		"Installing Dependencies",
		"Generating Project Structure",
		"Verifying SSO Authentication",
		"Wiring Standard API Access",
		"Installing StoryBook",
		"Configuring Production Setup",
		"Finalizing Setup",
	}

	renderer := components.NewPolishedRenderer(
		"SampleApp",
		"./SampleApp",
		"TypeScript",
		stepNames,
	)

	fmt.Println("🔍 TESTING: Initial state (all pending)")
	fmt.Println()
	output := renderer.Render(80)
	fmt.Print(output)
	fmt.Println()
	fmt.Println()

	// Complete some steps
	for i := 0; i < 6; i++ {
		renderer.SetCurrentStep(i)
		renderer.UpdateStep(i, 1.0, "Completed", []string{})
		renderer.CompleteStep(i, time.Millisecond*1500)
	}

	// Set current step with partial progress
	renderer.SetCurrentStep(8) // Production setup
	renderer.UpdateStep(8, 0.765, "Configuring build optimization...", []string{})

	fmt.Println("🔍 TESTING: Progress state (6 complete, 1 running at 76.5%)")
	fmt.Println()
	output = renderer.Render(80)
	fmt.Print(output)
	fmt.Println()
	fmt.Println()

	fmt.Println("🎯 POLISHED UI FEATURES:")
	fmt.Printf("   ✅ Aligned progress bars with fixed width (%d chars)\n", 28)
	fmt.Printf("   ✅ Header with app name and total progress\n")
	fmt.Printf("   ✅ Current step info with spinner animation\n")
	fmt.Printf("   ✅ Consistent step formatting: [icon] name [progress] %%\n")
	fmt.Printf("   ✅ Footer with target directory and template\n")
	fmt.Printf("   ✅ Clean separator lines\n")
	fmt.Printf("   ✅ Right-aligned status information\n")

	// Test different templates
	fmt.Println()
	fmt.Println("🔍 TESTING: Different template display")
	jsRenderer := components.NewPolishedRenderer(
		"JSApp",
		"./JSApp",
		"JavaScript",
		[]string{"Setup", "Configure"},
	)
	jsOutput := jsRenderer.Render(80)

	// Extract just the footer to show template formatting
	lines := splitLines(jsOutput)
	if len(lines) > 0 {
		footer := lines[len(lines)-1]
		fmt.Printf("JavaScript template footer: %s\n", footer)
	}
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