package config

import (
	"fmt"
	"strings"
)

// TemplateType defines the available project templates
type TemplateType string

const (
	TypeScript TemplateType = "typescript"
	JavaScript TemplateType = "javascript"
	Minimal    TemplateType = "minimal"
)

func (t TemplateType) String() string {
	return string(t)
}

func (t TemplateType) DisplayName() string {
	switch t {
	case TypeScript:
		return "TypeScript"
	case JavaScript:
		return "JavaScript"
	case Minimal:
		return "Minimal"
	default:
		return string(t)
	}
}

func (t TemplateType) Description() string {
	switch t {
	case TypeScript:
		return "Full type safety with modern tooling (Recommended)"
	case JavaScript:
		return "Fast setup with traditional JavaScript"
	case Minimal:
		return "Bare minimum setup for custom configuration"
	default:
		return ""
	}
}

// UserConfiguration represents the complete user configuration from prompts
type UserConfiguration struct {
	ProjectName string         `json:"projectName"`
	Template    TemplateConfig `json:"template"`
	DevFeatures DevFeatureConfig `json:"devFeatures"`
	ProductionSetup ProductionConfig `json:"productionSetup"`
	Testing     TestingConfig  `json:"testing"`
}

// TemplateConfig contains template-specific configuration
type TemplateConfig struct {
	Type TemplateType `json:"type"`
}

// DevFeatureConfig contains development feature selections
type DevFeatureConfig struct {
	HotReload    bool `json:"hotReload"`
	Linting      bool `json:"linting"`
	Prettier     bool `json:"prettier"`
	Husky        bool `json:"husky"`
	VSCodeConfig bool `json:"vscodeConfig"`
	DevTools     bool `json:"devTools"`
}

// ProductionConfig contains production setup selections
type ProductionConfig struct {
	Docker     bool `json:"docker"`
	Azure      bool `json:"azure"`
	CI_CD      bool `json:"cicd"`
	Monitoring bool `json:"monitoring"`
	Analytics  bool `json:"analytics"`

	// Production data access features
	TrustBridge bool `json:"trustBridge"`
	GRPC        bool `json:"grpc"`
	GridHDFS    bool `json:"gridHdfs"`
}

// TestingConfig contains testing framework selections
type TestingConfig struct {
	UnitTesting bool `json:"unitTesting"`
	E2ETesting  bool `json:"e2eTesting"`
	Coverage    bool `json:"coverage"`
}

// GetSelectedDevFeatures returns a list of selected development features
func (d DevFeatureConfig) GetSelected() []string {
	var selected []string
	if d.HotReload {
		selected = append(selected, "Hot Reload")
	}
	if d.Linting {
		selected = append(selected, "ESLint")
	}
	if d.Prettier {
		selected = append(selected, "Prettier")
	}
	if d.Husky {
		selected = append(selected, "Husky Pre-commit")
	}
	if d.VSCodeConfig {
		selected = append(selected, "VS Code Config")
	}
	if d.DevTools {
		selected = append(selected, "React DevTools")
	}
	return selected
}

// GetSelectedProductionFeatures returns a list of selected production features
func (p ProductionConfig) GetSelected() []string {
	var selected []string
	if p.Docker {
		selected = append(selected, "Docker")
	}
	if p.CI_CD {
		selected = append(selected, "CI/CD Pipeline")
	}
	if p.Monitoring {
		selected = append(selected, "Monitoring")
	}
	if p.Analytics {
		selected = append(selected, "Analytics")
	}
	return selected
}

// GetSelectedTestingFeatures returns a list of selected testing features
func (t TestingConfig) GetSelected() []string {
	var selected []string
	if t.UnitTesting {
		selected = append(selected, "Unit Testing")
	}
	if t.E2ETesting {
		selected = append(selected, "E2E Testing")
	}
	if t.Coverage {
		selected = append(selected, "Coverage Reports")
	}
	return selected
}

// GetSummary returns a formatted summary of the configuration
func (c UserConfiguration) GetSummary() string {
	var summary strings.Builder

	summary.WriteString(fmt.Sprintf("Project: %s\n", c.ProjectName))
	summary.WriteString(fmt.Sprintf("Template: %s\n\n", c.Template.Type.DisplayName()))

	devFeatures := c.DevFeatures.GetSelected()
	if len(devFeatures) > 0 {
		summary.WriteString("Development Features:\n")
		for _, feature := range devFeatures {
			summary.WriteString(fmt.Sprintf("• %s\n", feature))
		}
		summary.WriteString("\n")
	}

	prodFeatures := c.ProductionSetup.GetSelected()
	if len(prodFeatures) > 0 {
		summary.WriteString("Production Setup:\n")
		for _, feature := range prodFeatures {
			summary.WriteString(fmt.Sprintf("• %s\n", feature))
		}
		summary.WriteString("\n")
	}

	testFeatures := c.Testing.GetSelected()
	if len(testFeatures) > 0 {
		summary.WriteString("Testing:\n")
		for _, feature := range testFeatures {
			summary.WriteString(fmt.Sprintf("• %s\n", feature))
		}
		summary.WriteString("\n")
	}

	return summary.String()
}

// GetSmartDefaults returns sensible default configuration
func GetSmartDefaults(projectName string) UserConfiguration {
	return UserConfiguration{
		ProjectName: projectName,
		Template: TemplateConfig{
			Type: TypeScript, // Most popular choice
		},
		DevFeatures: DevFeatureConfig{
			HotReload: true,  // Essential for development
			Linting:   true,  // Code quality
			Prettier:  true,  // Code formatting
			DevTools:  true,  // Debugging
			// Leave advanced features unselected by default
		},
		ProductionSetup: ProductionConfig{
			// Leave production features unselected by default
			// User should consciously choose these
		},
		Testing: TestingConfig{
			UnitTesting: true, // Basic testing is good practice
			// Leave advanced testing unselected
		},
	}
}

// ValidateConfiguration checks for configuration conflicts
func ValidateConfiguration(config UserConfiguration) []string {
	var warnings []string

	// Check template vs feature compatibility
	if config.Template.Type == Minimal {
		devCount := len(config.DevFeatures.GetSelected())
		if devCount > 2 {
			warnings = append(warnings,
				"Minimal template with many dev features may be complex to set up")
		}
	}

	// Check production setup dependencies
	if config.ProductionSetup.Docker && !config.ProductionSetup.CI_CD {
		warnings = append(warnings,
			"Docker setup is recommended with CI/CD pipeline for automated deployments")
	}

	// Check testing setup
	if config.Testing.E2ETesting && !config.Testing.UnitTesting {
		warnings = append(warnings,
			"E2E testing is typically combined with unit testing for comprehensive coverage")
	}

	return warnings
}

// EstimateSetupTime calculates estimated setup time based on configuration
func EstimateSetupTime(config UserConfiguration) int {
	baseTime := 180 // 3 minutes base

	// Template complexity
	switch config.Template.Type {
	case TypeScript:
		baseTime += 60 // Additional TypeScript setup
	case Minimal:
		baseTime -= 30 // Simpler setup
	}

	// Development features
	devFeatures := len(config.DevFeatures.GetSelected())
	baseTime += devFeatures * 30 // 30 seconds per dev feature

	// Production features
	prodFeatures := len(config.ProductionSetup.GetSelected())
	baseTime += prodFeatures * 60 // 1 minute per production feature

	// Testing features
	testFeatures := len(config.Testing.GetSelected())
	baseTime += testFeatures * 45 // 45 seconds per testing feature

	return baseTime
}