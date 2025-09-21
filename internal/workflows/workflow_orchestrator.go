package workflows

import (
	"fmt"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
)

// WorkflowOrchestrator manages the execution of workflow stages
type WorkflowOrchestrator struct {
	stages  []WorkflowStage
	context *WorkflowContext
	logger  interfaces.Logger
}

// NewWorkflowOrchestrator creates a new workflow orchestrator
func NewWorkflowOrchestrator(context *WorkflowContext, logger interfaces.Logger) *WorkflowOrchestrator {
	return &WorkflowOrchestrator{
		stages:  make([]WorkflowStage, 0),
		context: context,
		logger:  logger,
	}
}

// AddStage adds a workflow stage to the orchestrator
func (wo *WorkflowOrchestrator) AddStage(stage WorkflowStage) {
	wo.stages = append(wo.stages, stage)
}

// Execute runs all workflow stages in sequence
func (wo *WorkflowOrchestrator) Execute() error {
	for _, stage := range wo.stages {
		stageName := stage.GetName()

		// Check if stage should be skipped
		if stage.CanSkip(wo.context) {
			continue
		}

		// Execute stage
		result, err := stage.Execute(wo.context)
		if err != nil {
			if wo.logger != nil {
				wo.logger.Error("Stage %s failed: %v", stageName, err)
			}
			return fmt.Errorf("stage %s failed: %w", stageName, err)
		}

		// Merge result into context
		if result != nil {
			err = wo.context.MergeResult(result)
			if err != nil {
				if wo.logger != nil {
					wo.logger.Error("Failed to merge result from stage %s: %v", stageName, err)
				}
				return fmt.Errorf("failed to merge stage result: %w", err)
			}
		}
	}

	return nil
}

// GetContext returns the current workflow context
func (wo *WorkflowOrchestrator) GetContext() *WorkflowContext {
	return wo.context
}

// GetStageCount returns the number of stages in the workflow
func (wo *WorkflowOrchestrator) GetStageCount() int {
	return len(wo.stages)
}

// GetStages returns all workflow stages
func (wo *WorkflowOrchestrator) GetStages() []WorkflowStage {
	// Return a copy to prevent external modification
	stages := make([]WorkflowStage, len(wo.stages))
	copy(stages, wo.stages)
	return stages
}

// ValidateWorkflow validates that the workflow is properly configured
func (wo *WorkflowOrchestrator) ValidateWorkflow() error {
	if len(wo.stages) == 0 {
		return fmt.Errorf("workflow must have at least one stage")
	}

	if wo.context == nil {
		return fmt.Errorf("workflow context cannot be nil")
	}

	// Validate each stage
	for i, stage := range wo.stages {
		if stage == nil {
			return fmt.Errorf("stage %d is nil", i)
		}

		stageName := stage.GetName()
		if stageName == "" {
			return fmt.Errorf("stage %d has empty name", i)
		}
	}

	return nil
}

// Reset clears all stages and resets the workflow
func (wo *WorkflowOrchestrator) Reset() {
	wo.stages = make([]WorkflowStage, 0)
}

// SetContext updates the workflow context
func (wo *WorkflowOrchestrator) SetContext(context *WorkflowContext) {
	wo.context = context
}