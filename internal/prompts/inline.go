package prompts

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/internal/config"
	"golang.org/x/term"
)

// ANSI escape codes for styling
const (
	// Colors
	colorReset  = "\033[0m"
	colorGrey   = "\033[90m"   // Bright black (grey)

	// Styles
	styleItalic = "\033[3m"
	styleReset  = "\033[23m"   // Reset italic

	// Combined style for responses
	responseStyle = colorGrey + styleItalic
	resetStyle    = styleReset + colorReset
)

// InlinePrompter handles traditional CLI-style prompting
type InlinePrompter struct {
	config     *config.PromptConfiguration
	userConfig *config.UserConfiguration
	reader     *bufio.Reader
}

// NewInlinePrompter creates a new inline prompter
func NewInlinePrompter() (*InlinePrompter, error) {
	promptConfig, err := config.LoadPromptConfiguration()
	if err != nil {
		return nil, err
	}

	return &InlinePrompter{
		config:     promptConfig,
		userConfig: &config.UserConfiguration{},
		reader:     bufio.NewReader(os.Stdin),
	}, nil
}

// RunPrompts executes all applicable prompts based on conditions with new styling
func (ip *InlinePrompter) RunPrompts(devOnly bool, flags []string) (*config.UserConfiguration, error) {
	return ip.RunPromptsWithContext(devOnly, flags, "", "")
}

// RunPromptsWithContext executes all applicable prompts with archetype context and styling
func (ip *InlinePrompter) RunPromptsWithContext(devOnly bool, flags []string, archetypeName string, appName string) (*config.UserConfiguration, error) {
	// Import components for styling
	// Note: Will need to import the components package

	// Don't display styled header - remove the configuration header section
	// Just proceed directly to prompts after archetype selection

	// Initialize user config with defaults
	ip.userConfig = &config.UserConfiguration{
		ProjectName: "", // Will be set by caller
		Template: config.TemplateConfig{
			Type: config.TypeScript, // Default
		},
		DevFeatures: config.DevFeatureConfig{
			HotReload:    true,
			Linting:      true,
			Prettier:     true,
			VSCodeConfig: true,
		},
		ProductionSetup: config.ProductionConfig{
			Docker:     false,
			CI_CD:      false,
			Monitoring: false,
			Analytics:  false,
		},
		Testing: config.TestingConfig{
			UnitTesting: true,
			E2ETesting:  false,
			Coverage:    true,
		},
	}

	// Process each prompt
	for _, promptConfig := range ip.config.Prompts {
		if promptConfig.ShouldTrigger(devOnly, flags) {
			err := ip.askPrompt(&promptConfig)
			if err != nil {
				return nil, err
			}
		}
	}

	// No footer needed - clean prompt flow

	return ip.userConfig, nil
}

// askPrompt handles a single prompt interaction with enhanced formatting
func (ip *InlinePrompter) askPrompt(prompt *config.PromptConfig) error {
	for {
		// Show the question with single space formatting to align with checkmarks
		fmt.Printf(" ? %s ", prompt.Question)

		// Read user input
		input, err := ip.reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		// Clean up input
		input = strings.TrimSpace(input)

		// Validate input
		if !prompt.IsValidInput(input) {
			// Show valid options with proper indentation and styling
			validOptions := make([]string, 0, len(prompt.UserOptions))
			for option := range prompt.UserOptions {
				validOptions = append(validOptions, option)
			}
			errorMsg := fmt.Sprintf("Invalid input. Valid options: %s", strings.Join(validOptions, ", "))
			fmt.Printf("  └ %s%s%s\n", responseStyle, errorMsg, resetStyle)
			continue
		}

		// Apply the configuration
		err = ip.applyPromptResult(prompt, input)
		if err != nil {
			return err
		}

		// Show response message with enhanced formatting (italic grey)
		responseLines := prompt.GetResponseLines(input)
		if len(responseLines) > 0 {
			for i, line := range responseLines {
				if i == 0 {
					fmt.Printf("  └ %s%s%s\n", responseStyle, line, resetStyle)
				} else if i == len(responseLines)-1 {
					fmt.Printf("  └ %s%s%s\n", responseStyle, line, resetStyle)
				} else {
					fmt.Printf("  ├ %s%s%s\n", responseStyle, line, resetStyle)
				}
			}
		}

		// Add breathing room between prompts
		fmt.Println()

		break
	}

	return nil
}

// applyPromptResult applies the user's response to the configuration
func (ip *InlinePrompter) applyPromptResult(prompt *config.PromptConfig, input string) error {
	normalizedInput := strings.ToLower(strings.TrimSpace(input))
	configValue := prompt.UserOptions[normalizedInput]

	switch prompt.ConfigKey {
	case "ProductionDataAccess":
		needsProduction := configValue == "true"
		if needsProduction {
			// Enable production data access features
			ip.userConfig.ProductionSetup.TrustBridge = true
			ip.userConfig.ProductionSetup.GRPC = true
			ip.userConfig.ProductionSetup.GridHDFS = true
		}

	case "DeploymentTarget":
		if configValue == "docker" {
			ip.userConfig.ProductionSetup.Docker = true
			ip.userConfig.ProductionSetup.CI_CD = true
		} else if configValue == "azure" {
			ip.userConfig.ProductionSetup.Azure = true
			ip.userConfig.ProductionSetup.CI_CD = true
		}

	default:
		// Could add more configuration mappings here
	}

	return nil
}

// GetUserConfiguration returns the final user configuration
func (ip *InlinePrompter) GetUserConfiguration() *config.UserConfiguration {
	return ip.userConfig
}

// getTerminalWidth gets the current terminal width for responsive headers
func getTerminalWidth() int {
	// Try to get terminal width from stdout
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && width > 0 {
		return width
	}
	// Fallback to a reasonable default
	return 80
}