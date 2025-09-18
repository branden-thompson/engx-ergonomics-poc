package aar

import (
	"fmt"
)

// NextStepsEngine generates contextual next steps based on the project configuration
type NextStepsEngine struct {
	templates map[string][]NextStepTemplate
	rules     []NextStepRule
}

// NextStepTemplate defines a template for generating next steps
type NextStepTemplate struct {
	Condition   func(*AARSummary) bool
	Action      string
	Description string
	Command     string
	Priority    StepPriority
	Category    StepCategory
	WorkingDir  string
}

// NextStepRule defines a rule for dynamic next step generation
type NextStepRule struct {
	Name      string
	Condition func(*AARSummary) bool
	Generator func(*AARSummary) []NextStep
}

// NewNextStepsEngine creates a new next steps engine with default templates and rules
func NewNextStepsEngine() *NextStepsEngine {
	engine := &NextStepsEngine{
		templates: make(map[string][]NextStepTemplate),
		rules:     make([]NextStepRule, 0),
	}

	engine.initializeDefaultTemplates()
	engine.initializeDefaultRules()

	return engine
}

// Generate generates next steps based on the AAR summary
func (e *NextStepsEngine) Generate(summary *AARSummary) []NextStep {
	var steps []NextStep

	// Generate steps from templates
	templateSteps := e.generateFromTemplates(summary)
	steps = append(steps, templateSteps...)

	// Generate steps from dynamic rules
	ruleSteps := e.generateFromRules(summary)
	steps = append(steps, ruleSteps...)

	// Prioritize and limit steps
	return e.prioritizeAndLimit(steps, 8) // Limit to 8 most important steps
}

// generateFromTemplates generates steps using template matching
func (e *NextStepsEngine) generateFromTemplates(summary *AARSummary) []NextStep {
	var steps []NextStep

	// Get templates for the project template type
	templates, exists := e.templates[summary.ProjectInfo.Template]
	if !exists {
		// Fall back to default templates
		templates = e.templates["default"]
	}

	for _, template := range templates {
		if template.Condition == nil || template.Condition(summary) {
			step := NextStep{
				Action:      template.Action,
				Description: template.Description,
				Command:     e.processCommand(template.Command, summary),
				Priority:    template.Priority,
				Category:    template.Category,
				WorkingDir:  e.processWorkingDir(template.WorkingDir, summary),
			}
			steps = append(steps, step)
		}
	}

	return steps
}

// generateFromRules generates steps using dynamic rules
func (e *NextStepsEngine) generateFromRules(summary *AARSummary) []NextStep {
	var steps []NextStep

	for _, rule := range e.rules {
		if rule.Condition(summary) {
			ruleSteps := rule.Generator(summary)
			steps = append(steps, ruleSteps...)
		}
	}

	return steps
}

// processCommand processes command templates with variable substitution
func (e *NextStepsEngine) processCommand(command string, summary *AARSummary) string {
	if command == "" {
		return ""
	}

	// Simple variable substitution
	// Replace {{PROJECT_NAME}} with actual project name
	result := command
	if summary.ProjectInfo.Name != "" {
		result = fmt.Sprintf(result, summary.ProjectInfo.Name)
	}

	return result
}

// processWorkingDir processes working directory with variable substitution
func (e *NextStepsEngine) processWorkingDir(workingDir string, summary *AARSummary) string {
	if workingDir == "" || workingDir == "{{PROJECT_DIR}}" {
		return summary.ProjectInfo.Directory
	}
	return workingDir
}

// prioritizeAndLimit sorts steps by priority and limits the number returned
func (e *NextStepsEngine) prioritizeAndLimit(steps []NextStep, limit int) []NextStep {
	if len(steps) <= limit {
		return steps
	}

	// Sort by priority (higher priority first)
	priorityOrder := []StepPriority{PriorityCritical, PriorityHigh, PriorityMedium, PriorityLow}
	var sortedSteps []NextStep

	for _, priority := range priorityOrder {
		for _, step := range steps {
			if step.Priority == priority && len(sortedSteps) < limit {
				sortedSteps = append(sortedSteps, step)
			}
		}
	}

	return sortedSteps
}

// initializeDefaultTemplates sets up default next step templates
func (e *NextStepsEngine) initializeDefaultTemplates() {
	// TypeScript template steps
	e.templates["typescript"] = []NextStepTemplate{
		{
			Condition:   func(s *AARSummary) bool { return s.ExecutionInfo.FailedSteps == 0 },
			Action:      "Start Development Server",
			Description: "Launch the development server to begin coding",
			Command:     "cd %s && npm run dev",
			Priority:    PriorityHigh,
			Category:    CategoryDevelopment,
			WorkingDir:  "{{PROJECT_DIR}}",
		},
		{
			Condition:   func(s *AARSummary) bool { return s.ExecutionInfo.FailedSteps == 0 },
			Action:      "Open in Editor",
			Description: "Open project in your preferred code editor",
			Command:     "code %s",
			Priority:    PriorityHigh,
			Category:    CategoryDevelopment,
			WorkingDir:  "",
		},
		{
			Condition:   func(s *AARSummary) bool { return s.ProjectInfo.Features["unit_testing"] },
			Action:      "Run Tests",
			Description: "Execute the test suite to verify everything works",
			Command:     "cd %s && npm test",
			Priority:    PriorityMedium,
			Category:    CategoryTesting,
			WorkingDir:  "{{PROJECT_DIR}}",
		},
		{
			Action:      "Initialize Git Repository",
			Description: "Set up version control for your project",
			Command:     "cd %s && git init && git add . && git commit -m 'Initial commit'",
			Priority:    PriorityMedium,
			Category:    CategoryConfiguration,
			WorkingDir:  "{{PROJECT_DIR}}",
		},
		{
			Condition:   func(s *AARSummary) bool { return s.ProjectInfo.Features["linting"] },
			Action:      "Check Code Quality",
			Description: "Run linting to ensure code quality standards",
			Command:     "cd %s && npm run lint",
			Priority:    PriorityLow,
			Category:    CategoryDevelopment,
			WorkingDir:  "{{PROJECT_DIR}}",
		},
	}

	// JavaScript template steps (similar but no TypeScript-specific items)
	e.templates["javascript"] = []NextStepTemplate{
		{
			Condition:   func(s *AARSummary) bool { return s.ExecutionInfo.FailedSteps == 0 },
			Action:      "Start Development Server",
			Description: "Launch the development server to begin coding",
			Command:     "cd %s && npm run dev",
			Priority:    PriorityHigh,
			Category:    CategoryDevelopment,
			WorkingDir:  "{{PROJECT_DIR}}",
		},
		{
			Condition:   func(s *AARSummary) bool { return s.ExecutionInfo.FailedSteps == 0 },
			Action:      "Open in Editor",
			Description: "Open project in your preferred code editor",
			Command:     "code %s",
			Priority:    PriorityHigh,
			Category:    CategoryDevelopment,
			WorkingDir:  "",
		},
		{
			Action:      "Initialize Git Repository",
			Description: "Set up version control for your project",
			Command:     "cd %s && git init && git add . && git commit -m 'Initial commit'",
			Priority:    PriorityMedium,
			Category:    CategoryConfiguration,
			WorkingDir:  "{{PROJECT_DIR}}",
		},
	}

	// Minimal template steps
	e.templates["minimal"] = []NextStepTemplate{
		{
			Condition:   func(s *AARSummary) bool { return s.ExecutionInfo.FailedSteps == 0 },
			Action:      "Start Development Server",
			Description: "Launch the development server",
			Command:     "cd %s && npm run dev",
			Priority:    PriorityHigh,
			Category:    CategoryDevelopment,
			WorkingDir:  "{{PROJECT_DIR}}",
		},
		{
			Action:      "Open in Editor",
			Description: "Open project for editing",
			Command:     "code %s",
			Priority:    PriorityHigh,
			Category:    CategoryDevelopment,
			WorkingDir:  "",
		},
	}

	// Default fallback templates
	e.templates["default"] = e.templates["typescript"]
}

// initializeDefaultRules sets up dynamic rules for special cases
func (e *NextStepsEngine) initializeDefaultRules() {
	// Failed steps recovery rule
	e.rules = append(e.rules, NextStepRule{
		Name: "Failed Steps Recovery",
		Condition: func(s *AARSummary) bool {
			return s.ExecutionInfo.FailedSteps > 0
		},
		Generator: func(s *AARSummary) []NextStep {
			return []NextStep{
				{
					Action:      "Review Error Messages",
					Description: "Check the troubleshooting section above for specific guidance",
					Priority:    PriorityCritical,
					Category:    CategoryTroubleshooting,
				},
				{
					Action:      "Retry Project Creation",
					Description: "After addressing issues, try creating the project again",
					Command:     "engx create %s",
					Priority:    PriorityHigh,
					Category:    CategoryTroubleshooting,
				},
			}
		},
	})

	// Production setup rule
	e.rules = append(e.rules, NextStepRule{
		Name: "Production Features",
		Condition: func(s *AARSummary) bool {
			return !s.ProjectInfo.DevOnly &&
				   (s.ProjectInfo.Features["docker"] ||
				    s.ProjectInfo.Features["cicd"] ||
				    s.ProjectInfo.Features["monitoring"])
		},
		Generator: func(s *AARSummary) []NextStep {
			var steps []NextStep

			if s.ProjectInfo.Features["docker"] {
				steps = append(steps, NextStep{
					Action:      "Build Docker Image",
					Description: "Create Docker image for deployment",
					Command:     "cd %s && docker build -t %s .",
					Priority:    PriorityMedium,
					Category:    CategoryDeployment,
					WorkingDir:  "{{PROJECT_DIR}}",
				})
			}

			if s.ProjectInfo.Features["cicd"] {
				steps = append(steps, NextStep{
					Action:      "Configure CI/CD",
					Description: "Set up continuous integration and deployment",
					Priority:    PriorityMedium,
					Category:    CategoryDeployment,
				})
			}

			return steps
		},
	})

	// Development workflow rule
	e.rules = append(e.rules, NextStepRule{
		Name: "Development Workflow",
		Condition: func(s *AARSummary) bool {
			return s.ExecutionInfo.FailedSteps == 0 && s.ProjectInfo.Features["hot_reload"]
		},
		Generator: func(s *AARSummary) []NextStep {
			return []NextStep{
				{
					Action:      "Learn About Hot Reload",
					Description: "Your project includes hot reload for fast development iterations",
					Priority:    PriorityLow,
					Category:    CategoryDocumentation,
				},
			}
		},
	})
}