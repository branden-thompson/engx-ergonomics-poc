# Plugin System Architecture Design

## 🎯 **Architectural Proposals**

### **Proposal A: Minimal Refactor (Conservative)**

#### **Architecture Overview**
```go
// cmd/engx/main.go - Stays mostly the same
func main() {
    rootCmd := &cobra.Command{...}

    // Registry-based command loading
    registry := commands.NewRegistry()
    for _, cmd := range registry.GetCommands() {
        rootCmd.AddCommand(cmd)
    }

    rootCmd.Execute()
}

// internal/commands/registry.go - New file
type Registry struct {
    commands []*cobra.Command
}

func NewRegistry() *Registry {
    return &Registry{
        commands: []*cobra.Command{
            NewCreateCommand(),    // Existing function
            NewTestErrorCommand(), // Existing function
        },
    }
}
```

**Directory Structure**:
```
engx-ergonomics-poc/
├── cmd/engx/main.go                 # Minimal changes
├── internal/commands/
│   ├── registry.go                  # New: Command registration
│   ├── create.go                    # Existing: Unchanged
│   ├── test_error.go               # Existing: Unchanged
│   └── update.go                   # New: Future commands
└── internal/[existing packages]     # No changes
```

**Pros**:
- ✅ Minimal risk - small changes
- ✅ Quick implementation (1-2 days)
- ✅ Zero breaking changes
- ✅ Easy to understand and review

**Cons**:
- ❌ Limited modularity improvement
- ❌ Still requires main.go changes for new commands
- ❌ No shared service architecture
- ❌ Minimal collaboration benefits

---

### **Proposal B: Package-Based Plugin System (Recommended)**

#### **Architecture Overview**
```go
// cmd/engx/main.go - Minimal bootstrap
func main() {
    deps := common.NewDependencies()
    registry := registry.NewManager(deps)

    // Auto-discover and register plugins
    plugins := registry.DiscoverPlugins()
    for _, plugin := range plugins {
        registry.Register(plugin)
    }

    rootCmd := createRootCommand()
    for _, cmd := range registry.GetCommands() {
        rootCmd.AddCommand(cmd)
    }

    rootCmd.Execute()
}

// pkg/common/interfaces/plugin.go
type CommandPlugin interface {
    Name() string
    Description() string
    Version() string
    Create(deps *Dependencies) *cobra.Command
    Initialize() error
    Cleanup() error
    RequiredServices() []string
}

// pkg/commands/create/plugin.go
type Plugin struct {
    initialized bool
}

func (p *Plugin) Name() string { return "create" }
func (p *Plugin) Create(deps *Dependencies) *cobra.Command {
    return createCommand(deps)
}
```

**Directory Structure**:
```
engx-ergonomics-poc/
├── cmd/engx/main.go                 # Minimal bootstrap (~40 lines)
├── pkg/
│   ├── commands/                    # Plugin commands
│   │   ├── create/
│   │   │   ├── plugin.go           # Plugin interface impl
│   │   │   ├── command.go          # Cobra command logic
│   │   │   ├── handlers.go         # Business logic
│   │   │   └── types.go            # Command types
│   │   ├── test/
│   │   │   └── [similar structure]
│   │   └── update/                 # Future commands
│   ├── common/                      # Shared utilities
│   │   ├── interfaces/             # Plugin interfaces
│   │   ├── config/                 # Shared config
│   │   ├── tui/                    # Shared TUI components
│   │   ├── chaos/                  # Moved from internal
│   │   ├── aar/                    # Moved from internal
│   │   └── dependencies.go         # Dependency injection
│   └── registry/                    # Plugin management
│       ├── manager.go              # Plugin lifecycle
│       ├── discovery.go            # Auto-discovery
│       └── validation.go           # Plugin validation
└── internal/ → pkg/common/          # Migration path
```

**Pros**:
- ✅ True modularity - commands completely isolated
- ✅ Team collaboration - independent development
- ✅ Shared service architecture
- ✅ Auto-discovery of plugins
- ✅ Dependency injection for testing
- ✅ Follows Go standard project layout

**Cons**:
- ⚠️ Medium complexity - requires careful migration
- ⚠️ More files and abstractions
- ⚠️ Learning curve for team

---

### **Proposal C: Advanced Plugin System with Hot-Reload (Aggressive)**

#### **Architecture Overview**
```go
// Advanced plugin loading with validation and hot-reload
type AdvancedPlugin interface {
    CommandPlugin
    Validate() error
    Dependencies() PluginDependencies
    Configuration() PluginConfig
    HealthCheck() error
}

type PluginManager struct {
    plugins     map[string]AdvancedPlugin
    watcher     *fsnotify.Watcher
    hotReload   bool
}

func (pm *PluginManager) WatchForChanges() {
    // File system watching for development
    // Hot-reload plugins on change
}
```

**Features**:
- Plugin validation and health checks
- Configuration management per plugin
- Hot-reload during development
- Plugin dependency graphs
- Advanced error handling and recovery

**Pros**:
- ✅ Maximum flexibility and features
- ✅ Hot-reload for rapid development
- ✅ Advanced plugin validation
- ✅ Configuration management

**Cons**:
- ❌ High complexity
- ❌ Long development time (3-4 weeks)
- ❌ Over-engineering for current needs
- ❌ Higher maintenance burden

## 🎯 **Dependency Injection Architecture**

### **Core Dependencies System**
```go
// pkg/common/dependencies.go
type Dependencies struct {
    Config     ConfigManager
    TUI        TUIRegistry
    Chaos      ChaosInjector
    AAR        AARGenerator
    Logger     Logger
    Filesystem FilesystemManager
}

func NewDependencies() *Dependencies {
    return &Dependencies{
        Config:     config.NewManager(),
        TUI:        tui.NewRegistry(),
        Chaos:      chaos.NewInjector(),
        AAR:        aar.NewGenerator(),
        Logger:     logging.NewLogger(),
        Filesystem: fs.NewManager(),
    }
}

// Service interfaces for clean testing
type ConfigManager interface {
    LoadConfig(path string) (*Config, error)
    GetVerbosity() VerbosityLevel
    GetChaosConfig() *ChaosConfig
}

type TUIRegistry interface {
    CreateProgressModel(steps []Step) tea.Model
    CreateConfirmModel(prompt string) tea.Model
    GetTheme() *Theme
}
```

### **Plugin Implementation Pattern**
```go
// pkg/commands/create/plugin.go
type Plugin struct {
    deps *common.Dependencies
}

func (p *Plugin) Create(deps *common.Dependencies) *cobra.Command {
    p.deps = deps

    cmd := &cobra.Command{
        Use:   "create [APP_NAME]",
        Short: "Create a new React application",
        RunE:  p.runCreate,
    }

    // Add command-specific flags
    cmd.Flags().BoolVar(&devOnly, "dev-only", false, "Skip deployment preparation")
    cmd.Flags().BoolVar(&chaosMarine, "chaos-marine", false, "Enable chaos engineering")

    return cmd
}

func (p *Plugin) runCreate(cmd *cobra.Command, args []string) error {
    // Access all shared services through p.deps
    config := p.deps.Config.LoadConfig("")
    tuiModel := p.deps.TUI.CreateProgressModel(steps)

    // Business logic using injected dependencies
    return p.executeCreate(args[0], config, tuiModel)
}
```

## 📋 **Migration Strategy**

### **Phase 1: Infrastructure Setup (Week 1)**
```
Day 1-2: Create pkg/ structure and interfaces
Day 3-4: Implement plugin registry and dependency injection
Day 5:   Move common utilities to pkg/common/
```

### **Phase 2: Command Migration (Week 2)**
```
Day 1-2: Convert create command to plugin
Day 3-4: Convert test-error command to plugin
Day 5:   Update main.go and integration testing
```

### **Phase 3: Enhancement & Polish (Week 3)**
```
Day 1-2: Plugin validation and error handling
Day 3-4: Documentation and developer guides
Day 5:   Performance testing and optimization
```

### **Migration Safety Measures**
1. **Feature Flags**: Toggle between old/new system
2. **Parallel Implementation**: Keep old code until new code proven
3. **Comprehensive Testing**: Test every command after migration
4. **Rollback Plan**: Quick revert to pre-migration state

## ⚖️ **Risk Assessment & Mitigation**

### **High Priority Risks**

**Risk 1: Breaking Existing Functionality**
- **Probability**: Medium
- **Impact**: High
- **Mitigation**:
  - Comprehensive integration tests before migration
  - Feature flag for old/new system switching
  - Gradual migration with immediate rollback capability

**Risk 2: Over-Engineering**
- **Probability**: Medium
- **Impact**: Medium
- **Mitigation**:
  - Start with Proposal B (balanced approach)
  - YAGNI principle - implement only current needs
  - Regular review with team for complexity check

**Risk 3: Team Learning Curve**
- **Probability**: Low
- **Impact**: Medium
- **Mitigation**:
  - Comprehensive documentation and examples
  - Plugin development template
  - Pair programming for first plugin

### **Medium Priority Risks**

**Risk 4: Performance Degradation**
- **Probability**: Low
- **Impact**: Low
- **Mitigation**:
  - Benchmark tests before/after migration
  - Profile memory usage and startup time
  - Optimize plugin discovery if needed

## 🎯 **Resource Requirements**

### **Development Time Estimates**

**Proposal A (Minimal)**: 2-3 days
- 1 day: Registry implementation
- 1 day: Testing and integration
- 0.5 day: Documentation

**Proposal B (Recommended)**: 2-3 weeks
- Week 1: Infrastructure and interfaces
- Week 2: Command migration
- Week 3: Testing, docs, and polish

**Proposal C (Advanced)**: 4-5 weeks
- Additional complexity requires extended development

### **Team Resource Requirements**
- **Primary Developer**: 1 person full-time
- **Code Reviews**: 2-3 team members for architecture review
- **Testing**: QA support for comprehensive testing
- **Documentation**: Technical writing support

---

**📊 Planning Status**: Multiple Proposals Ready for Decision
**🎯 Recommendation**: Proposal B (Package-Based Plugin System)
**👥 Decision Required**: HUM LEAD architectural approach selection
**⏱️ Timeline**: 2-3 weeks for Proposal B implementation