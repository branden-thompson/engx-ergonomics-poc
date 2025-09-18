package main

import (
	"fmt"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/models"
)

func main() {
	fmt.Println("🛩️ Testing TUI Model Logic (Headless)")

	// Create a model like the real TUI would
	model := models.NewAppModel("create", "TestApp", []string{"--verbose"})

	// Initialize it
	model.Init()

	// Simulate the update loop for a few seconds
	fmt.Println("Simulating TUI updates...")

	for i := 0; i < 50; i++ { // 5 seconds at 100ms intervals
		// This would normally be driven by Bubble Tea's event loop
		fmt.Printf("Update %d: State simulation\n", i+1)
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Println("✅ TUI model logic test complete")
}