package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PromptConfig represents a single configurable prompt
type PromptConfig struct {
	ID             string            `json:"id" yaml:"id"`
	Trigger        string            `json:"trigger" yaml:"trigger"`           // Condition when prompt should appear
	Question       string            `json:"question" yaml:"question"`         // The question to ask
	UserOptions    map[string]string `json:"user_options" yaml:"user_options"` // Valid responses and their meanings
	ResponseFormat string            `json:"response_format" yaml:"response_format"` // Template for response message
	ConfigKey      string            `json:"config_key" yaml:"config_key"`     // Which config field this affects
}

// PromptConfiguration holds all prompt configurations
type PromptConfiguration struct {
	Prompts []PromptConfig `json:"prompts" yaml:"prompts"`
}

// LoadPromptConfiguration loads prompt config from JSON file
func LoadPromptConfiguration() (*PromptConfiguration, error) {
	// Default prompts if no config file exists
	defaultConfig := &PromptConfiguration{
		Prompts: []PromptConfig{
			{
				ID:      "production_data",
				Trigger: "always", // Always ask this question
				Question: "Will you need access to production data? (y/n)",
				UserOptions: map[string]string{
					"y":   "true",
					"yes": "true",
					"n":   "false",
					"no":  "false",
				},
				ResponseFormat: "Will %s production data access setup...",
				ConfigKey:     "ProductionDataAccess",
			},
			{
				ID:      "deployment_target",
				Trigger: "not_dev_only", // Only ask if --dev-only is NOT specified
				Question: "Will this deploy to 1) LI Data-Centers, 2) Azure Cloud? (1/2)",
				UserOptions: map[string]string{
					"1": "docker",
					"2": "azure",
				},
				ResponseFormat: "Will configure application for %s deployment",
				ConfigKey:     "DeploymentTarget",
			},
		},
	}

	// Try to load from config file
	configPath := filepath.Join(".", ".engx", "prompts.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return defaultConfig, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return defaultConfig, nil // Fall back to default on error
	}

	config := &PromptConfiguration{}
	if err := json.Unmarshal(data, config); err != nil {
		return defaultConfig, nil // Fall back to default on error
	}

	return config, nil
}

// ShouldTrigger determines if a prompt should be shown based on conditions
func (p *PromptConfig) ShouldTrigger(devOnly bool, flags []string) bool {
	switch p.Trigger {
	case "always":
		return true
	case "not_dev_only":
		return !devOnly
	case "dev_only":
		return devOnly
	default:
		// Could add more complex trigger logic here
		return false
	}
}

// GetResponseMessage formats the response message based on user input
func (p *PromptConfig) GetResponseMessage(userInput string) string {
	// Normalize user input
	normalizedInput := strings.ToLower(strings.TrimSpace(userInput))

	// Find the config value for this input
	configValue, exists := p.UserOptions[normalizedInput]
	if !exists {
		return ""
	}

	// Format the response message
	var displayValue string
	switch p.ID {
	case "production_data":
		if configValue == "true" {
			displayValue = "include"
		} else {
			displayValue = "bypass"
		}
	case "deployment_target":
		if configValue == "docker" {
			displayValue = "Docker Container"
		} else {
			displayValue = "Azure Container"
		}
	default:
		displayValue = configValue
	}

	return fmt.Sprintf(p.ResponseFormat, displayValue)
}

// GetResponseLines returns formatted response lines for enhanced tree display
func (p *PromptConfig) GetResponseLines(userInput string) []string {
	// Normalize user input
	normalizedInput := strings.ToLower(strings.TrimSpace(userInput))

	// Find the config value for this input
	configValue, exists := p.UserOptions[normalizedInput]
	if !exists {
		return []string{}
	}

	// Return detailed response lines based on prompt type
	switch p.ID {
	case "production_data":
		if configValue == "true" {
			return []string{
				"Including production data-access setup...",
			}
		} else {
			return []string{
				"Bypassing production data-access setup...",
			}
		}
	case "deployment_target":
		if configValue == "docker" {
			return []string{
				"Configuring for Docker Container deployment",
				"Setting Application to use standard LCD pipelines (can be changed later...)",
				"Setting fabric to CORP (can be changed later...)",
			}
		} else {
			return []string{
				"Configuring for Azure Container deployment",
				"Setting Application to use Azure DevOps pipelines (can be changed later...)",
				"Setting fabric to CLOUD (can be changed later...)",
			}
		}
	default:
		// Fallback to single line response
		return []string{fmt.Sprintf(p.ResponseFormat, configValue)}
	}
}

// IsValidInput checks if user input is valid for this prompt
func (p *PromptConfig) IsValidInput(input string) bool {
	normalizedInput := strings.ToLower(strings.TrimSpace(input))
	_, exists := p.UserOptions[normalizedInput]
	return exists
}