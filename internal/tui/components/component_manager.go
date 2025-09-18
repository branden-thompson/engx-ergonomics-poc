package components

// ComponentInstallationPhase defines when components get installed
type ComponentInstallationPhase int

const (
	PhaseDependencies ComponentInstallationPhase = iota
	PhaseProjectStructure
	PhaseTestingFrameworks
	PhaseDocumentation
	PhaseFinalizing
)

// ComponentInstallationStep defines a specific installation step within a phase
type ComponentInstallationStep struct {
	Phase           ComponentInstallationPhase
	ProgressStart   float64 // 0.0 to 1.0 - when this step starts within the phase
	ProgressEnd     float64 // 0.0 to 1.0 - when this step completes within the phase
	ComponentNames  []string
	SuccessRate     float64 // 0.0 to 1.0 - for future failure simulation
}

// ComponentManager handles the installation simulation of all components
type ComponentManager struct {
	installationPlan []ComponentInstallationStep
	currentPhase     ComponentInstallationPhase
	phaseProgress    float64
}

// NewComponentManager creates a new component manager with the default installation plan
func NewComponentManager() *ComponentManager {
	plan := []ComponentInstallationStep{
		// DEPENDENCIES PHASE: Core Technologies (10%-60%) + EngX Start (70%-80%)
		{
			Phase:          PhaseDependencies,
			ProgressStart:  0.1,
			ProgressEnd:    0.2,
			ComponentNames: []string{"TypeScript"},
			SuccessRate:    0.98,
		},
		{
			Phase:          PhaseDependencies,
			ProgressStart:  0.2,
			ProgressEnd:    0.3,
			ComponentNames: []string{"React"},
			SuccessRate:    0.97,
		},
		{
			Phase:          PhaseDependencies,
			ProgressStart:  0.3,
			ProgressEnd:    0.4,
			ComponentNames: []string{"React Router 7"},
			SuccessRate:    0.96,
		},
		{
			Phase:          PhaseDependencies,
			ProgressStart:  0.4,
			ProgressEnd:    0.5,
			ComponentNames: []string{"Tailwind CSS"},
			SuccessRate:    0.95,
		},
		{
			Phase:          PhaseDependencies,
			ProgressStart:  0.5,
			ProgressEnd:    0.6,
			ComponentNames: []string{"Radix UI"},
			SuccessRate:    0.94,
		},
		{
			Phase:          PhaseDependencies,
			ProgressStart:  0.6,
			ProgressEnd:    0.7,
			ComponentNames: []string{"ShadCN-based UI Design System (SUDS)"},
			SuccessRate:    0.93,
		},
		{
			Phase:          PhaseDependencies,
			ProgressStart:  0.7,
			ProgressEnd:    0.8,
			ComponentNames: []string{"TrustBridge SSO"},
			SuccessRate:    0.92,
		},
		{
			Phase:          PhaseDependencies,
			ProgressStart:  0.8,
			ProgressEnd:    0.9,
			ComponentNames: []string{"gRPC Web"},
			SuccessRate:    0.91,
		},
		{
			Phase:          PhaseDependencies,
			ProgressStart:  0.9,
			ProgressEnd:    1.0,
			ComponentNames: []string{"GRID/HDFS Access"},
			SuccessRate:    0.90,
		},

		// PROJECT STRUCTURE PHASE: Remaining EngX (20%-80%)
		{
			Phase:          PhaseProjectStructure,
			ProgressStart:  0.2,
			ProgressEnd:    0.5,
			ComponentNames: []string{"CREWS API"},
			SuccessRate:    0.89,
		},
		{
			Phase:          PhaseProjectStructure,
			ProgressStart:  0.5,
			ProgressEnd:    0.8,
			ComponentNames: []string{"LI CATALOG API"},
			SuccessRate:    0.88,
		},
		{
			Phase:          PhaseProjectStructure,
			ProgressStart:  0.8,
			ProgressEnd:    1.0,
			ComponentNames: []string{"GitHub Actions"},
			SuccessRate:    0.87,
		},

		// TESTING FRAMEWORKS PHASE: Quality & Testing (25%-100%)
		{
			Phase:          PhaseTestingFrameworks,
			ProgressStart:  0.25,
			ProgressEnd:    0.5,
			ComponentNames: []string{"Vitest"},
			SuccessRate:    0.95,
		},
		{
			Phase:          PhaseTestingFrameworks,
			ProgressStart:  0.5,
			ProgressEnd:    0.75,
			ComponentNames: []string{"EngX TypeScript Linters"},
			SuccessRate:    0.94,
		},
		{
			Phase:          PhaseTestingFrameworks,
			ProgressStart:  0.75,
			ProgressEnd:    0.9,
			ComponentNames: []string{"GitHub Pages"},
			SuccessRate:    0.93,
		},
		{
			Phase:          PhaseTestingFrameworks,
			ProgressStart:  0.9,
			ProgressEnd:    1.0,
			ComponentNames: []string{"StoryBook (UI Components & Documentation)"},
			SuccessRate:    0.92,
		},
	}

	return &ComponentManager{
		installationPlan: plan,
		currentPhase:     PhaseDependencies,
		phaseProgress:    0.0,
	}
}

// GetInstallationUpdates returns the components that should change status at the given phase and progress
func (cm *ComponentManager) GetInstallationUpdates(phase ComponentInstallationPhase, progress float64) []ComponentUpdate {
	var updates []ComponentUpdate

	// Find all steps that should be active at this phase and progress
	for _, step := range cm.installationPlan {
		if step.Phase != phase {
			continue
		}

		for _, componentName := range step.ComponentNames {
			if progress >= step.ProgressStart && progress < step.ProgressEnd {
				// Component should be installing
				updates = append(updates, ComponentUpdate{
					ComponentName: componentName,
					NewStatus:     "installing",
					NewIcon:       "[✓ ]",
				})
			} else if progress >= step.ProgressEnd {
				// Component should be installed
				updates = append(updates, ComponentUpdate{
					ComponentName: componentName,
					NewStatus:     "installed",
					NewIcon:       "[✓]",
				})
			}
		}
	}

	return updates
}

// GetAllComponentsUpToPhase returns all components that should be installed up to and including the given phase
func (cm *ComponentManager) GetAllComponentsUpToPhase(phase ComponentInstallationPhase, progress float64) []ComponentUpdate {
	var updates []ComponentUpdate
	processedComponents := make(map[string]bool)

	// Get all components from earlier phases (should be completed)
	for _, step := range cm.installationPlan {
		if step.Phase < phase {
			for _, componentName := range step.ComponentNames {
				if !processedComponents[componentName] {
					updates = append(updates, ComponentUpdate{
						ComponentName: componentName,
						NewStatus:     "installed",
						NewIcon:       "[✓]",
					})
					processedComponents[componentName] = true
				}
			}
		}
	}

	// Get components from current phase based on progress
	currentPhaseUpdates := cm.GetInstallationUpdates(phase, progress)
	for _, update := range currentPhaseUpdates {
		if !processedComponents[update.ComponentName] {
			updates = append(updates, update)
			processedComponents[update.ComponentName] = true
		}
	}

	return updates
}

// GetAllComponentsForPhase returns all components that should be installed during or before a phase
func (cm *ComponentManager) GetAllComponentsForPhase(phase ComponentInstallationPhase) []ComponentUpdate {
	var updates []ComponentUpdate
	processedComponents := make(map[string]bool)

	// Get all components from this phase and earlier phases
	for _, step := range cm.installationPlan {
		if step.Phase <= phase {
			for _, componentName := range step.ComponentNames {
				if !processedComponents[componentName] {
					updates = append(updates, ComponentUpdate{
						ComponentName: componentName,
						NewStatus:     "installed",
						NewIcon:       "[✓]",
					})
					processedComponents[componentName] = true
				}
			}
		}
	}

	return updates
}

// ComponentUpdate represents a status change for a component
type ComponentUpdate struct {
	ComponentName string
	NewStatus     string
	NewIcon       string
}

// MapStepNameToPhase converts step names to installation phases
func MapStepNameToPhase(stepName string) ComponentInstallationPhase {
	switch stepName {
	case "Installing Dependencies", "Installing dependencies":
		return PhaseDependencies
	case "Generating Project Structure", "Generating project structure":
		return PhaseProjectStructure
	case "Installing Testing Frameworks", "Installing testing frameworks":
		return PhaseTestingFrameworks
	case "Generating Documentation", "Generating documentation":
		return PhaseDocumentation
	case "Finalizing Setup", "Finalizing setup":
		return PhaseFinalizing
	default:
		return PhaseDependencies
	}
}