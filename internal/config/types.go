package config

import "time"

// Config represents the complete configuration structure
type Config struct {
	Project      *ProjectConfig      `yaml:"project,omitempty"`
	Defaults     *DefaultsConfig     `yaml:"defaults,omitempty"`
	Environments map[string]*EnvConfig `yaml:"environments,omitempty"`
	Commands     map[string]*CmdConfig `yaml:"custom_commands,omitempty"`
}

// ProjectConfig contains project-specific settings
type ProjectConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description,omitempty"`
	Repository  string `yaml:"repository,omitempty"`
}

// DefaultsConfig contains default behavior settings
type DefaultsConfig struct {
	Verbosity        string        `yaml:"verbosity"`         // normal, verbose, quiet
	DeploymentTarget string        `yaml:"deployment_target"` // development, staging, production
	Timeout          time.Duration `yaml:"timeout"`
	Theme            string        `yaml:"theme"`             // auto, dark, light
	Template         string        `yaml:"template,omitempty"` // typescript, javascript, minimal
}

// EnvConfig contains environment-specific settings
type EnvConfig struct {
	APIEndpoint string            `yaml:"api_endpoint"`
	Variables   map[string]string `yaml:"variables,omitempty"`
	Timeout     time.Duration     `yaml:"timeout,omitempty"`
}

// CmdConfig contains custom command configurations
type CmdConfig struct {
	Flags []string          `yaml:"flags,omitempty"`
	Env   map[string]string `yaml:"env,omitempty"`
}

// NewDefaultConfig returns a configuration with sensible defaults
func NewDefaultConfig() *Config {
	return &Config{
		Defaults: &DefaultsConfig{
			Verbosity:        "normal",
			DeploymentTarget: "production",
			Timeout:          5 * time.Minute,
			Theme:            "auto",
			Template:         "typescript",
		},
		Environments: map[string]*EnvConfig{
			"development": {
				APIEndpoint: "http://localhost:3000",
				Timeout:     30 * time.Second,
			},
			"staging": {
				APIEndpoint: "https://staging.company.com",
				Timeout:     2 * time.Minute,
			},
			"production": {
				APIEndpoint: "https://api.company.com",
				Timeout:     5 * time.Minute,
			},
		},
		Commands: make(map[string]*CmdConfig),
	}
}

// Merge merges another config into this one, with the other config taking precedence
func (c *Config) Merge(other *Config) {
	if other == nil {
		return
	}

	// Merge project config
	if other.Project != nil {
		c.Project = other.Project
	}

	// Merge defaults
	if other.Defaults != nil {
		if c.Defaults == nil {
			c.Defaults = &DefaultsConfig{}
		}
		c.mergeDefaults(other.Defaults)
	}

	// Merge environments
	if other.Environments != nil {
		if c.Environments == nil {
			c.Environments = make(map[string]*EnvConfig)
		}
		for name, env := range other.Environments {
			c.Environments[name] = env
		}
	}

	// Merge commands
	if other.Commands != nil {
		if c.Commands == nil {
			c.Commands = make(map[string]*CmdConfig)
		}
		for name, cmd := range other.Commands {
			c.Commands[name] = cmd
		}
	}
}

func (c *Config) mergeDefaults(other *DefaultsConfig) {
	if other.Verbosity != "" {
		c.Defaults.Verbosity = other.Verbosity
	}
	if other.DeploymentTarget != "" {
		c.Defaults.DeploymentTarget = other.DeploymentTarget
	}
	if other.Timeout != 0 {
		c.Defaults.Timeout = other.Timeout
	}
	if other.Theme != "" {
		c.Defaults.Theme = other.Theme
	}
	if other.Template != "" {
		c.Defaults.Template = other.Template
	}
}

// GetEnvironment returns the configuration for a specific environment
func (c *Config) GetEnvironment(name string) *EnvConfig {
	if c.Environments == nil {
		return nil
	}
	return c.Environments[name]
}

// GetCommand returns the configuration for a specific command
func (c *Config) GetCommand(name string) *CmdConfig {
	if c.Commands == nil {
		return nil
	}
	return c.Commands[name]
}