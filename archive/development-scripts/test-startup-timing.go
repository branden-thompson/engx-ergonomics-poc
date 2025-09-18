package main

import (
	"fmt"
	"time"

	progresssim "github.com/bthompso/engx-ergonomics-poc/internal/simulation/progress"
)

func main() {
	fmt.Println("🕒 Test Startup Timing Issues")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Create tracker
	tracker := progresssim.NewCreateTracker(true)

	fmt.Println("📋 STEP ANALYSIS:")
	totalDuration := time.Duration(0)
	for i := 0; i < tracker.TotalSteps(); i++ {
		step := tracker.GetStep(i)
		if step != nil {
			totalDuration += step.Duration
			fmt.Printf("  Step %d: %s - %v\n", i, step.Name, step.Duration)
		}
	}
	fmt.Printf("  Total expected: %v\n", totalDuration)
	fmt.Println()

	// Test the startup sequence like the real TUI
	fmt.Println("🚀 SIMULATING TUI STARTUP SEQUENCE:")

	// 1. tracker.Start() - like startExecution()
	fmt.Println("1. tracker.Start() called")
	tracker.Start()
	startTime := time.Now()

	// 2. Check initial state
	fmt.Printf("   Current step: %d\n", tracker.CurrentStep())
	fmt.Printf("   Step info: %s\n", tracker.CurrentStepInfo().Name)
	fmt.Printf("   Is step ready: %v\n", tracker.IsStepReady())

	// 3. Simulate the first nextStep() call (200ms delay)
	fmt.Println("2. nextStep() scheduled (200ms delay)")
	time.Sleep(time.Millisecond * 200)

	fmt.Printf("   After 200ms - Is step ready: %v\n", tracker.IsStepReady())
	stepInfo := tracker.CurrentStepInfo()
	if stepInfo != nil {
		elapsed := time.Since(tracker.GetStepStart())
		stepProgress := float64(elapsed) / float64(stepInfo.Duration)
		fmt.Printf("   Step progress: %.1f%% (elapsed: %v / %v)\n",
			stepProgress*100, elapsed, stepInfo.Duration)
	}

	// 4. Wait for first step to become ready
	fmt.Println("3. Waiting for first step to become ready...")
	stepReady := false
	for !stepReady && time.Since(startTime) < time.Second*5 {
		stepReady = tracker.IsStepReady()
		if !stepReady {
			time.Sleep(time.Millisecond * 50)
		}
	}

	if stepReady {
		elapsed := time.Since(startTime)
		fmt.Printf("   ✅ First step ready after: %v\n", elapsed)

		// Try to advance
		canAdvance := tracker.NextStep()
		fmt.Printf("   NextStep() returned: %v\n", canAdvance)

		if canAdvance {
			newStepInfo := tracker.CurrentStepInfo()
			fmt.Printf("   Advanced to: %s\n", newStepInfo.Name)
		}
	} else {
		fmt.Printf("   ❌ First step never became ready within 5 seconds\n")
	}

	fmt.Println()
	fmt.Println("🔍 TIMING ANALYSIS:")

	// Analyze each step's expected timing
	currentTime := time.Duration(0)
	for i := 0; i < tracker.TotalSteps(); i++ {
		step := tracker.GetStep(i)
		if step != nil {
			startAt := currentTime
			endAt := currentTime + step.Duration
			fmt.Printf("  Step %d: %8v - %8v (%v) - %s\n",
				i, startAt, endAt, step.Duration, step.Name)
			currentTime = endAt
		}
	}

	fmt.Println()
	fmt.Printf("🎯 EXPECTED BEHAVIOR:\n")
	fmt.Printf("  • Progress ticks every 50ms should show smooth animation\n")
	fmt.Printf("  • Step checks every 200ms should advance steps when ready\n")
	fmt.Printf("  • Total time should be ~%v\n", totalDuration)
	fmt.Printf("  • Each step should reach 100%% before advancing\n")
	fmt.Printf("  • Final step should trigger completion at 100%%\n")
}