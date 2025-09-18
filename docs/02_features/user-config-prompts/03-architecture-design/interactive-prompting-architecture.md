# Interactive Prompting Architecture
**🛩️ BRTOPS PLAN Phase - Strategic Planning & Architecture Design**

## ARCHITECTURE OVERVIEW

### State Machine Extension
Extend existing AppModel state machine to include prompting phases:

```
Current Flow:
┌─────────┐    ┌──────────┐    ┌─────────────┐    ┌──────────┐
│  IDLE   │───▶│  PROMPT  │───▶│  EXECUTING  │───▶│ COMPLETE │
└─────────┘    └──────────┘    └─────────────┘    └──────────┘

Enhanced Flow:
┌─────────┐    ┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌──────────┐
│  IDLE   │───▶│ PROMPTING   │───▶│ VALIDATING   │───▶│  EXECUTING  │───▶│ COMPLETE │
└─────────┘    └─────────────┘    └──────────────┘    └─────────────┘    └──────────┘
                       │                   │
                       ▼                   │
               ┌─────────────┐             │
               │ PROMPT_HELP │◀────────────┘
               └─────────────┘
```

### Core Components Architecture

#### 1. Configuration Management Layer
```go
// Core configuration structures
type UserConfiguration struct {
    ProjectName     string
    Template        TemplateConfig
    DevFeatures     DevFeatureConfig
    ProductionSetup ProductionConfig
    Integrations    IntegrationConfig
    Testing         TestingConfig
    Advanced        AdvancedConfig
}

type TemplateConfig struct {
    Type        TemplateType  // typescript, javascript, minimal
    StyleSystem StyleType     // styled-components, emotion, css-modules
    StateManagement StateType // redux, zustand, context
}

type DevFeatureConfig struct {
    HotReload     bool
    Linting       bool
    Prettier      bool
    Husky         bool
    VSCodeConfig  bool
    DevTools      bool
}

type ProductionConfig struct {
    Docker        bool
    CI_CD         CICDType     // github-actions, azure-devops, none
    Monitoring    MonitoringType // sentry, datadog, none
    Analytics     AnalyticsType  // google, mixpanel, none
    Optimization  bool
}
```

#### 2. Interactive Component System
```go
// Base prompting interface
type PromptComponent interface {
    Init() tea.Cmd
    Update(tea.Msg) (PromptComponent, tea.Cmd)
    View() string
    GetValue() interface{}
    SetValue(interface{})
    Validate() error
}

// Specific prompt implementations
type TemplateSelector struct {
    list        list.Model
    choices     []TemplateChoice
    selected    int
    description string
}

type FeatureMultiSelect struct {
    list     list.Model
    choices  []FeatureChoice
    selected map[string]bool
    category string
}

type TextInput struct {
    input       textinput.Model
    validator   func(string) error
    placeholder string
    required    bool
}

type ConfirmDialog struct {
    message string
    result  *bool
    focused bool
}
```

#### 3. Prompt Flow Orchestrator
```go
type PromptOrchestrator struct {
    currentPrompt int
    prompts       []PromptComponent
    config        *UserConfiguration
    navigation    NavigationState
    help          HelpSystem
}

type NavigationState struct {
    canGoBack    bool
    canGoForward bool
    canSkip      bool
    showHelp     bool
}

type HelpSystem struct {
    contextHelp map[string]string
    examples    map[string][]string
    tips        map[string]string
}
```

## DETAILED DESIGN SPECIFICATIONS

### 1. Template Selection Flow
```
┌─────────────────────────────────────────────────────────────────┐
│ 🎯 Choose Your Template                                         │
├─────────────────────────────────────────────────────────────────┤
│ ▶ TypeScript (Recommended)                                     │
│   └─ Full type safety with modern tooling                      │
│                                                                 │
│   JavaScript                                                   │
│   └─ Fast setup with traditional JS                            │
│                                                                 │
│   Minimal                                                      │
│   └─ Bare minimum setup for custom configuration               │
├─────────────────────────────────────────────────────────────────┤
│ 💡 TypeScript provides better IDE support and catches errors   │
│    at compile time. Recommended for most projects.             │
├─────────────────────────────────────────────────────────────────┤
│ [Enter] Select • [↑↓] Navigate • [h] Help • [q] Quit          │
└─────────────────────────────────────────────────────────────────┘
```

### 2. Development Features Multi-Select
```
┌─────────────────────────────────────────────────────────────────┐
│ 🛠️  Development Features                                        │
├─────────────────────────────────────────────────────────────────┤
│ [✓] Hot Reload                 Fast development with auto-refresh │
│ [✓] ESLint + Prettier         Code quality and formatting       │
│ [✓] Husky + Lint-staged       Pre-commit hooks                  │
│ [ ] VS Code Configuration     Workspace settings and extensions │
│ [✓] React DevTools           Browser debugging extensions        │
│ [ ] Storybook                 Component development environment  │
├─────────────────────────────────────────────────────────────────┤
│ 💡 Selected: 4/6 features                                      │
│    These tools improve development speed and code quality       │
├─────────────────────────────────────────────────────────────────┤
│ [Space] Toggle • [Enter] Continue • [↑↓] Navigate • [a] All    │
└─────────────────────────────────────────────────────────────────┘
```

### 3. Configuration Summary
```
┌─────────────────────────────────────────────────────────────────┐
│ 📋 Configuration Summary                                        │
├─────────────────────────────────────────────────────────────────┤
│ Project: MyReactApp                                             │
│ Template: TypeScript                                            │
│                                                                 │
│ Development Features:                                           │
│ • Hot Reload, ESLint + Prettier, Husky                        │
│ • React DevTools                                               │
│                                                                 │
│ Production Setup:                                               │
│ • Docker containerization                                      │
│ • GitHub Actions CI/CD                                         │
│                                                                 │
│ Testing:                                                        │
│ • Jest + React Testing Library                                │
│ • Cypress E2E testing                                         │
│                                                                 │
│ Estimated setup time: ~8 minutes                              │
├─────────────────────────────────────────────────────────────────┤
│ [Enter] Create Project • [b] Back • [e] Edit • [s] Save Config │
└─────────────────────────────────────────────────────────────────┘
```

## COMPONENT SELECTION LOGIC

### Template-Driven Base Components
```go
func GetBaseComponentsForTemplate(template TemplateType) ComponentSet {
    switch template {
    case TypeScript:
        return ComponentSet{
            CoreTech: []string{"React", "TypeScript", "Vite"},
            DevDeps:  []string{"@types/react", "@types/node", "typescript"},
            Scripts:  map[string]string{
                "dev": "vite",
                "build": "tsc && vite build",
                "typecheck": "tsc --noEmit",
            },
        }
    case JavaScript:
        return ComponentSet{
            CoreTech: []string{"React", "Vite"},
            DevDeps:  []string{"@vitejs/plugin-react"},
            Scripts:  map[string]string{
                "dev": "vite",
                "build": "vite build",
            },
        }
    case Minimal:
        return ComponentSet{
            CoreTech: []string{"React"},
            DevDeps:  []string{},
            Scripts:  map[string]string{
                "start": "react-scripts start",
            },
        }
    }
}
```

### Feature-Driven Component Addition
```go
func AddFeatureComponents(base ComponentSet, config UserConfiguration) ComponentSet {
    if config.DevFeatures.Linting {
        base.DevDeps = append(base.DevDeps, "eslint", "prettier")
        base.ConfigFiles = append(base.ConfigFiles, ".eslintrc.js", ".prettierrc")
    }

    if config.DevFeatures.Husky {
        base.DevDeps = append(base.DevDeps, "husky", "lint-staged")
        base.Scripts["prepare"] = "husky install"
    }

    if config.Testing.UnitTesting {
        base.DevDeps = append(base.DevDeps, "jest", "@testing-library/react")
        base.Scripts["test"] = "jest"
        base.Scripts["test:watch"] = "jest --watch"
    }

    return base
}
```

### Dynamic Installation Phase Generation
```go
func GenerateInstallationPhases(config UserConfiguration) []InstallationPhase {
    phases := []InstallationPhase{
        {Name: "Validating configuration", Duration: 1200 * time.Millisecond},
        {Name: "Setting up environment", Duration: 1800 * time.Millisecond},
    }

    // Dynamic dependency installation based on selections
    depDuration := calculateDependencyDuration(config)
    phases = append(phases, InstallationPhase{
        Name: "Installing dependencies",
        Duration: depDuration,
        Components: getSelectedComponents(config),
    })

    if config.Testing.UnitTesting || config.Testing.E2ETesting {
        phases = append(phases, InstallationPhase{
            Name: "Setting up testing frameworks",
            Duration: 2200 * time.Millisecond,
        })
    }

    if config.ProductionSetup.Docker || config.ProductionSetup.CI_CD != None {
        phases = append(phases, InstallationPhase{
            Name: "Configuring production setup",
            Duration: 2500 * time.Millisecond,
        })
    }

    phases = append(phases, InstallationPhase{
        Name: "Finalizing setup",
        Duration: 800 * time.Millisecond,
    })

    return phases
}
```

## BUBBLE TEA INTEGRATION DESIGN

### Enhanced AppModel Structure
```go
type AppModel struct {
    // Existing fields
    state      AppState
    command    string
    target     string
    flags      []string

    // New prompting fields
    promptOrchestrator *PromptOrchestrator
    userConfig         *UserConfiguration
    configValidation   []ValidationError

    // Enhanced execution fields
    dynamicSteps       []string
    dynamicComponents  ComponentSet

    // Existing TUI fields
    spinner    spinner.Model
    renderer   *components.EnhancedRenderer
    tracker    *progresssim.Tracker
    // ... rest of existing fields
}
```

### State Transition Logic
```go
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch m.state {
    case StateIdle:
        return m, m.startPrompting()

    case StatePrompting:
        return m.handlePromptingUpdate(msg)

    case StateValidating:
        return m.handleValidationUpdate(msg)

    case StateExecuting:
        // Use dynamic configuration for execution
        return m.handleExecutionUpdate(msg)

    // ... existing cases
    }
}

func (m *AppModel) handlePromptingUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
    orchestrator, cmd := m.promptOrchestrator.Update(msg)
    m.promptOrchestrator = orchestrator.(*PromptOrchestrator)

    if m.promptOrchestrator.IsComplete() {
        m.userConfig = m.promptOrchestrator.GetConfiguration()
        m.state = StateValidating
        return m, m.validateConfiguration()
    }

    return m, cmd
}
```

## VALIDATION AND ERROR HANDLING

### Configuration Validation Rules
```go
type ValidationRule struct {
    Name      string
    Validator func(UserConfiguration) error
    Severity  ValidationSeverity
}

var validationRules = []ValidationRule{
    {
        Name: "template-feature-compatibility",
        Validator: func(config UserConfiguration) error {
            if config.Template.Type == Minimal && len(config.DevFeatures.GetSelected()) > 2 {
                return errors.New("minimal template incompatible with many dev features")
            }
            return nil
        },
        Severity: Warning,
    },
    {
        Name: "production-dependencies",
        Validator: func(config UserConfiguration) error {
            if config.ProductionSetup.Docker && !config.ProductionSetup.CI_CD.IsSet() {
                return errors.New("docker setup recommended with CI/CD pipeline")
            }
            return nil
        },
        Severity: Warning,
    },
}
```

### Smart Defaults and Suggestions
```go
func GetSmartDefaults(context PromptContext) UserConfiguration {
    defaults := UserConfiguration{
        Template: TemplateConfig{Type: TypeScript}, // Most popular choice
    }

    // Context-aware defaults
    if context.IsFirstTime {
        defaults.DevFeatures = DevFeatureConfig{
            HotReload: true,
            Linting:   true,
            Prettier:  true,
        }
    }

    if context.ProjectSize == Large {
        defaults.Testing = TestingConfig{
            UnitTesting: true,
            E2ETesting:  true,
        }
    }

    return defaults
}
```

---
**PLAN Status**: ✅ ARCHITECTURE DESIGN COMPLETE
**Next Phase**: Implementation Strategy & Component Development
**Quality Gates**: SEV-0 Detailed Architecture Documentation Complete
**Ready for**: CODE Phase Implementation