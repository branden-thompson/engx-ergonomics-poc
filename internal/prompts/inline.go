package prompts

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/internal/config"
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

// RunPrompts executes all applicable prompts based on conditions
func (ip *InlinePrompter) RunPrompts(devOnly bool, flags []string) (*config.UserConfiguration, error) {
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

	return ip.userConfig, nil
}

// askPrompt handles a single prompt interaction with enhanced formatting
func (ip *InlinePrompter) askPrompt(prompt *config.PromptConfig) error {
	for {
		// Show the question with Claude Code-style formatting
		fmt.Printf("? %s ", prompt.Question)

		// Read user input
		input, err := ip.reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		// Clean up input
		input = strings.TrimSpace(input)

		// Validate input
		if !prompt.IsValidInput(input) {
			// Show valid options with proper indentation
			validOptions := make([]string, 0, len(prompt.UserOptions))
			for option := range prompt.UserOptions {
				validOptions = append(validOptions, option)
			}
			fmt.Printf("  └ Invalid input. Valid options: %s\n", strings.Join(validOptions, ", "))
			continue
		}

		// Apply the configuration
		err = ip.applyPromptResult(prompt, input)
		if err != nil {
			return err
		}

		// Show response message with enhanced formatting
		responseLines := prompt.GetResponseLines(input)
		if len(responseLines) > 0 {
			for i, line := range responseLines {
				if i == 0 {
					fmt.Printf("  └ %s\n", line)
				} else if i == len(responseLines)-1 {
					fmt.Printf("  └ %s\n", line)
				} else {
					fmt.Printf("  ├ %s\n", line)
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