package archetypes

import (
	"fmt"
	"sync"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/internal/workflows"
)

// Registry manages application archetype definitions
type Registry struct {
	archetypes map[string]*workflows.ArchetypeDefinition
	mutex      sync.RWMutex
}

// NewRegistry creates a new archetype registry
func NewRegistry() *Registry {
	registry := &Registry{
		archetypes: make(map[string]*workflows.ArchetypeDefinition),
	}

	// Register built-in archetypes
	registry.registerBuiltInArchetypes()

	return registry
}

// Register adds an archetype definition to the registry
func (r *Registry) Register(archetype *workflows.ArchetypeDefinition) error {
	if archetype == nil {
		return fmt.Errorf("archetype cannot be nil")
	}

	if archetype.ID == "" {
		return fmt.Errorf("archetype ID cannot be empty")
	}

	if archetype.Name == "" {
		return fmt.Errorf("archetype name cannot be empty")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.archetypes[archetype.ID] = archetype
	return nil
}

// GetByID retrieves an archetype by its ID
func (r *Registry) GetByID(id string) (*workflows.ArchetypeDefinition, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	archetype, exists := r.archetypes[id]
	if !exists {
		return nil, fmt.Errorf("archetype with ID '%s' not found", id)
	}

	return archetype, nil
}

// GetDefault returns the default archetype (prod-web)
func (r *Registry) GetDefault() *workflows.ArchetypeDefinition {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, archetype := range r.archetypes {
		if archetype.IsDefault {
			return archetype
		}
	}

	// Fallback to prod-web if no default is marked
	if archetype, exists := r.archetypes["prod-web"]; exists {
		return archetype
	}

	return nil
}

// GetAvailable returns all available archetypes
func (r *Registry) GetAvailable() []*workflows.ArchetypeDefinition {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	archetypes := make([]*workflows.ArchetypeDefinition, 0, len(r.archetypes))
	for _, archetype := range r.archetypes {
		archetypes = append(archetypes, archetype)
	}

	return archetypes
}

// GetByCategory returns archetypes filtered by category
func (r *Registry) GetByCategory(category workflows.ArchetypeCategory) []*workflows.ArchetypeDefinition {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var filtered []*workflows.ArchetypeDefinition
	for _, archetype := range r.archetypes {
		if archetype.Category == category {
			filtered = append(filtered, archetype)
		}
	}

	return filtered
}

// GetIDs returns all registered archetype IDs
func (r *Registry) GetIDs() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	ids := make([]string, 0, len(r.archetypes))
	for id := range r.archetypes {
		ids = append(ids, id)
	}

	return ids
}

// Count returns the number of registered archetypes
func (r *Registry) Count() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.archetypes)
}

// registerBuiltInArchetypes registers the built-in archetypes
func (r *Registry) registerBuiltInArchetypes() {
	// Production React Application
	prodWebArchetype := &workflows.ArchetypeDefinition{
		ID:          "prod-web",
		Name:        "Production React App",
		Description: "Full-featured React application with production optimizations",
		IsDefault:   true,
		Category:    workflows.CategoryWebApplication,
		PromptSet:   "prod-web-prompts",
		TUISteps: []workflows.StepDefinition{
			{ID: "init", Name: "Project Initialization", Description: "Setting up project structure", Type: workflows.StepTypeInitialization, Duration: 2 * time.Second},
			{ID: "deps", Name: "Installing Dependencies", Description: "npm install with production packages", Type: workflows.StepTypeDependencies, Duration: 5 * time.Second},
			{ID: "config", Name: "Configuration Setup", Description: "ESLint, Prettier, TypeScript config", Type: workflows.StepTypeConfiguration, Duration: 3 * time.Second},
			{ID: "build", Name: "Production Build Setup", Description: "Webpack optimization, CI/CD prep", Type: workflows.StepTypeBuild, Duration: 4 * time.Second},
			{ID: "deploy", Name: "Deployment Preparation", Description: "Docker, environment configuration", Type: workflows.StepTypeDeployment, Duration: 3 * time.Second},
			{ID: "finalize", Name: "Finalization", Description: "Final checks and documentation", Type: workflows.StepTypeFinalization, Duration: 2 * time.Second},
		},
		Features: []string{"typescript", "eslint", "prettier", "jest", "docker", "ci-cd"},
	}

	// Development React Application
	devWebArchetype := &workflows.ArchetypeDefinition{
		ID:          "dev-web",
		Name:        "Development React App",
		Description: "Lightweight React setup for rapid development and prototyping",
		IsDefault:   false,
		Category:    workflows.CategoryWebApplication,
		PromptSet:   "dev-web-prompts",
		TUISteps: []workflows.StepDefinition{
			{ID: "init", Name: "Quick Project Setup", Description: "Minimal project structure", Type: workflows.StepTypeInitialization, Duration: 1 * time.Second},
			{ID: "deps", Name: "Essential Dependencies", Description: "Core React packages only", Type: workflows.StepTypeDependencies, Duration: 3 * time.Second},
			{ID: "config", Name: "Dev Configuration", Description: "Hot reload, basic linting", Type: workflows.StepTypeConfiguration, Duration: 2 * time.Second},
			{ID: "finalize", Name: "Ready to Code", Description: "Development server ready", Type: workflows.StepTypeFinalization, Duration: 1 * time.Second},
		},
		Features: []string{"typescript", "hot-reload", "basic-linting"},
	}

	// CLI Tool
	cliArchetype := &workflows.ArchetypeDefinition{
		ID:          "cli",
		Name:        "Command Line Tool",
		Description: "Go-based CLI application with modern tooling",
		IsDefault:   false,
		Category:    workflows.CategoryCLITool,
		PromptSet:   "cli-prompts",
		TUISteps: []workflows.StepDefinition{
			{ID: "init", Name: "CLI Project Setup", Description: "Go module and project structure", Type: workflows.StepTypeInitialization, Duration: 1 * time.Second},
			{ID: "deps", Name: "CLI Dependencies", Description: "Cobra, Viper, and essential packages", Type: workflows.StepTypeDependencies, Duration: 2 * time.Second},
			{ID: "config", Name: "CLI Configuration", Description: "Command structure and config files", Type: workflows.StepTypeConfiguration, Duration: 3 * time.Second},
			{ID: "build", Name: "Build System", Description: "Cross-platform build configuration", Type: workflows.StepTypeBuild, Duration: 2 * time.Second},
			{ID: "finalize", Name: "CLI Ready", Description: "Command registration complete", Type: workflows.StepTypeFinalization, Duration: 1 * time.Second},
		},
		Features: []string{"cobra", "viper", "cross-platform", "packaging"},
	}

	// Backend Service
	serviceArchetype := &workflows.ArchetypeDefinition{
		ID:          "service",
		Name:        "Backend Service",
		Description: "REST API or microservice with database integration",
		IsDefault:   false,
		Category:    workflows.CategoryBackendService,
		PromptSet:   "service-prompts",
		TUISteps: []workflows.StepDefinition{
			{ID: "init", Name: "Service Setup", Description: "API project structure", Type: workflows.StepTypeInitialization, Duration: 2 * time.Second},
			{ID: "deps", Name: "Service Dependencies", Description: "HTTP framework and database drivers", Type: workflows.StepTypeDependencies, Duration: 4 * time.Second},
			{ID: "config", Name: "Service Configuration", Description: "API routes and middleware", Type: workflows.StepTypeConfiguration, Duration: 3 * time.Second},
			{ID: "build", Name: "Service Build", Description: "Container and deployment setup", Type: workflows.StepTypeBuild, Duration: 3 * time.Second},
			{ID: "deploy", Name: "Service Deployment", Description: "Production deployment configuration", Type: workflows.StepTypeDeployment, Duration: 2 * time.Second},
			{ID: "finalize", Name: "Service Ready", Description: "API endpoints configured", Type: workflows.StepTypeFinalization, Duration: 1 * time.Second},
		},
		Features: []string{"rest-api", "database", "docker", "monitoring"},
	}

	// Hackday Prototype
	hackdayArchetype := &workflows.ArchetypeDefinition{
		ID:          "hackday",
		Name:        "Hackday Prototype",
		Description: "Rapid prototyping setup for hackathons and experiments",
		IsDefault:   false,
		Category:    workflows.CategoryPrototype,
		PromptSet:   "hackday-prompts",
		TUISteps: []workflows.StepDefinition{
			{ID: "init", Name: "Rapid Setup", Description: "Minimal project structure", Type: workflows.StepTypeInitialization, Duration: 1 * time.Second},
			{ID: "deps", Name: "Essential Tools", Description: "Core packages for rapid development", Type: workflows.StepTypeDependencies, Duration: 2 * time.Second},
			{ID: "finalize", Name: "Hack Ready", Description: "Ready for rapid development", Type: workflows.StepTypeFinalization, Duration: 1 * time.Second},
		},
		Features: []string{"minimal", "rapid-development", "experimental"},
	}

	// EngX Command Plugin
	engxCmdArchetype := &workflows.ArchetypeDefinition{
		ID:          "engx-cmd",
		Name:        "EngX Command Plugin",
		Description: "Plugin for extending EngX functionality",
		IsDefault:   false,
		Category:    workflows.CategoryPlugin,
		PromptSet:   "engx-cmd-prompts",
		TUISteps: []workflows.StepDefinition{
			{ID: "init", Name: "Plugin Setup", Description: "EngX plugin structure", Type: workflows.StepTypeInitialization, Duration: 1 * time.Second},
			{ID: "deps", Name: "Plugin Dependencies", Description: "EngX interfaces and common packages", Type: workflows.StepTypeDependencies, Duration: 2 * time.Second},
			{ID: "config", Name: "Plugin Configuration", Description: "Command interface implementation", Type: workflows.StepTypeConfiguration, Duration: 3 * time.Second},
			{ID: "finalize", Name: "Plugin Ready", Description: "EngX plugin registration complete", Type: workflows.StepTypeFinalization, Duration: 1 * time.Second},
		},
		Features: []string{"engx-plugin", "command-interface", "cobra"},
	}

	// Register all built-in archetypes
	archetypes := []*workflows.ArchetypeDefinition{
		prodWebArchetype,
		devWebArchetype,
		cliArchetype,
		serviceArchetype,
		hackdayArchetype,
		engxCmdArchetype,
	}

	for _, archetype := range archetypes {
		r.archetypes[archetype.ID] = archetype
	}
}

// GetDefaultRegistry returns a singleton instance of the registry
var defaultRegistry *Registry
var registryOnce sync.Once

func GetDefaultRegistry() *Registry {
	registryOnce.Do(func() {
		defaultRegistry = NewRegistry()
	})
	return defaultRegistry
}