package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/bthompso/engx-ergonomics-poc/internal/tui/models"
)

func main() {
	fmt.Println("🎯 User Configuration Prompts Demo")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("🔍 TESTING: Interactive prompting system")
	fmt.Println("This demo shows the new user configuration prompts that determine")
	fmt.Println("which components get installed during project creation.")
	fmt.Println()

	// Create app model without any flags (triggers prompting)
	model := models.NewAppModel("create", "MyTestApp", []string{})

	// Run the TUI
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}

	fmt.Println()
	fmt.Println("✨ DEMO COMPLETE")
	fmt.Println("The interactive prompting system allows users to:")
	fmt.Println("• Select project template (TypeScript, JavaScript, Minimal)")
	fmt.Println("• Choose development features (linting, testing, etc.)")
	fmt.Println("• Configure production setup (Docker, CI/CD, monitoring)")
	fmt.Println("• Set up testing frameworks")
	fmt.Println("• Review and confirm configuration before proceeding")
	fmt.Println()
	fmt.Println("This demonstrates the interaction patterns that the engineering")
	fmt.Println("team should implement in the real version of this tool.")
}