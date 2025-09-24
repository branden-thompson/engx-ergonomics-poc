package crews

import (
	"github.com/bthompso/engx-ergonomics-poc/pkg/common"
	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
	"github.com/spf13/cobra"
)

// Plugin implements the CommandPlugin interface for the crews command
type Plugin struct {
	deps *common.Dependencies
}

// NewPlugin creates a new crews command plugin
func NewPlugin(deps *common.Dependencies) interfaces.CommandPlugin {
	return &Plugin{deps: deps}
}

// Name returns the plugin name
func (p *Plugin) Name() string {
	return "crews"
}

// Description returns the plugin description
func (p *Plugin) Description() string {
	return "Manage and query crew information, membership, and ownership"
}

// Version returns the plugin version
func (p *Plugin) Version() string {
	return "1.0.0"
}

// Create returns the cobra command for this plugin
func (p *Plugin) Create(deps interface{}) *cobra.Command {
	// Return the crews command
	return NewCrewsCommand()
}

// Initialize initializes the plugin
func (p *Plugin) Initialize() error {
	return nil
}

// Cleanup cleans up plugin resources
func (p *Plugin) Cleanup() error {
	return nil
}

// RequiredServices returns the list of required services
func (p *Plugin) RequiredServices() []string {
	return []string{}
}

// OptionalServices returns the list of optional services
func (p *Plugin) OptionalServices() []string {
	return []string{}
}