# Technical Design - Crew Command Plugin

## MAJOR FEATURE ENHANCEMENT | LEVEL-1 SEV-0
**Feature**: `crew-command-plugin`
**Date**: 2025-09-23
**Phase**: Architecture Design
**Status**: Technical Specification Complete

## 🏗️ System Architecture Overview

### Plugin Integration Strategy

Following established EngX patterns for consistent command experience:

```
EngX CLI Core
├── cmd/engx/main.go (cobra root)
├── pkg/common/ (shared utilities)
└── plugins/
    ├── create/ (existing)
    ├── templates/ (existing)
    ├── analytics/ (existing)
    └── crews/ (NEW)
        ├── cmd.go (command registration)
        ├── smart_resolver.go (parameter detection)
        ├── subcommands/ (operation handlers)
        ├── models/ (data structures)
        ├── renderers/ (terminal output)
        └── data/ (simulation data)
```

### Component Architecture

```mermaid
graph TD
    A[crews command] --> B[Smart Parameter Resolver]
    B --> C{Parameter Type?}
    C -->|CREW-####| D[Details Handler]
    C -->|user@domain| E[Membership Handler]
    C -->|asset://| F[Owner Handler]
    C -->|unknown| G[Search Handler]

    D --> H[Crew Details Renderer]
    E --> I[Membership List Renderer]
    F --> J[Owner Details Renderer]
    G --> K[Search Results Renderer]

    H --> L[Terminal Output]
    I --> L
    J --> L
    K --> L
```

## 🧠 Smart Parameter Detection System

### Detection Pipeline Architecture

```go
// Parameter detection with priority-based routing
type ParameterResolver struct {
    patterns map[string]*regexp.Regexp
    handlers map[string]SubcommandHandler
    fallback SubcommandHandler
}

// Detection patterns (evaluated in order)
var DetectionPatterns = map[string]string{
    "crewID":    `^CREW-\d{4}$`,
    "email":     `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
    "assetURN":  `^(asset://|urn:asset:)`,
    "ldapUser":  `^[a-zA-Z][a-zA-Z0-9._-]*$`,
}

// Routing logic
func (r *ParameterResolver) Resolve(param string) (SubcommandHandler, error) {
    switch {
    case r.patterns["crewID"].MatchString(param):
        return r.handlers["details"], nil
    case r.patterns["email"].MatchString(param):
        return r.handlers["membership"], nil
    case r.patterns["assetURN"].MatchString(param):
        return r.handlers["owner"], nil
    case r.patterns["ldapUser"].MatchString(param):
        return r.handlers["membership"], nil
    default:
        return r.fallback, nil
    }
}
```

### Parameter Examples and Routing

| Input Example | Pattern Match | Route to Function | Output |
|---|---|---|---|
| `CREW-1234` | crewID | `Details()` | Crew information card |
| `bthompso@company.com` | email | `Membership()` | User's crew list |
| `asset://web-app/dashboard` | assetURN | `Owner()` | Owning crew details |
| `sarah.jones` | ldapUser | `Membership()` | User's crew list |
| `invalid-input` | none | `Search()` | Search suggestions |

## 📊 Data Model Architecture

### Core Entity Definitions

```go
// Primary crew entity
type Crew struct {
    ID              string              `json:"id"`              // CREW-1234
    VanityName      string              `json:"vanity_name"`     // "Web Platform Team"
    Description     string              `json:"description"`     // Team purpose
    CreatedAt       time.Time           `json:"created_at"`
    CreatedBy       string              `json:"created_by"`      // Creator LDAP
    Members         []Member            `json:"members"`
    OnCallSchedule  OnCallSchedule      `json:"oncall_schedule"`
    OwnedAssets     []string            `json:"owned_assets"`    // Asset URNs
    AccessGrants    []AccessGrant       `json:"access_grants"`   // External permissions
    Status          CrewStatus          `json:"status"`          // active, archived, etc.
}

// Member relationship with role-based access
type Member struct {
    UserID          string              `json:"user_id"`         // LDAP username
    Email           string              `json:"email"`           // Full email address
    FullName        string              `json:"full_name"`       // Display name
    Role            MemberRole          `json:"role"`            // Access level
    JoinedAt        time.Time           `json:"joined_at"`
    AddedBy         string              `json:"added_by"`        // Who added them
    IsOnCall        bool                `json:"is_oncall"`       // Current on-call status
    OnCallSchedule  []OnCallPeriod      `json:"oncall_periods"`  // Scheduled rotations
    LastActive      time.Time           `json:"last_active"`     // Activity tracking
}

// Role-based access control
type MemberRole string
const (
    RoleOwner   MemberRole = "owner"   // Full control, can transfer ownership
    RoleAdmin   MemberRole = "admin"   // Management privileges
    RoleMember  MemberRole = "member"  // Standard access
    RoleTemp    MemberRole = "temp"    // Temporary access with expiration
    RoleAuto    MemberRole = "auto"    // Inherited from asset ownership
    RoleRemoved MemberRole = "removed" // Explicitly revoked
)

// On-call scheduling system
type OnCallSchedule struct {
    CurrentOnCall   []string            `json:"current_oncall"`  // Current rotation
    Schedule        []OnCallPeriod      `json:"schedule"`        // Upcoming rotations
    RotationType    string              `json:"rotation_type"`   // weekly, daily, etc.
    EscalationPath  []string            `json:"escalation"`      // Backup contacts
}

type OnCallPeriod struct {
    UserID          string              `json:"user_id"`
    StartTime       time.Time           `json:"start_time"`
    EndTime         time.Time           `json:"end_time"`
    Type            string              `json:"type"`            // primary, backup
}

// External access delegation
type AccessGrant struct {
    GrantedTo       string              `json:"granted_to"`      // User or crew
    GrantedBy       string              `json:"granted_by"`      // Granting user
    GrantedAt       time.Time           `json:"granted_at"`
    AccessLevel     string              `json:"access_level"`    // read, write, admin
    AssetScope      []string            `json:"asset_scope"`     // Specific assets
    ExpiresAt       *time.Time          `json:"expires_at"`      // Optional expiration
}
```

### Asset Relationship Model

```go
// Asset ownership mapping
type AssetOwnership struct {
    AssetURN        string              `json:"asset_urn"`       // Unique asset identifier
    AssetName       string              `json:"asset_name"`      // Human-readable name
    AssetType       string              `json:"asset_type"`      // web-app, service, etc.
    OwnerCrewID     string              `json:"owner_crew_id"`   // Primary owner
    Delegates       []string            `json:"delegates"`       // Delegated crews
    CreatedAt       time.Time           `json:"created_at"`
    LastModified    time.Time           `json:"last_modified"`
    Metadata        map[string]string   `json:"metadata"`        // Additional properties
}
```

## 🎨 Terminal Rendering Architecture

### Renderer Component Design

```go
// Base renderer interface for consistent output
type Renderer interface {
    Render(data interface{}, width int) (string, error)
    SetColorScheme(scheme ColorScheme)
    SetVerbosity(level VerbosityLevel)
}

// Specialized renderers for different views
type CrewDetailsRenderer struct {
    colorScheme ColorScheme
    verbosity   VerbosityLevel
    components  *lipgloss.Style
}

type MembershipListRenderer struct {
    colorScheme ColorScheme
    tableConfig TableConfig
    statusIcons map[MemberRole]string
}

type AssetListRenderer struct {
    colorScheme ColorScheme
    groupBy     string // crew, type, etc.
    showCounts  bool
}
```

### Visual Design Templates

#### Crew Details View
```
┌─ CREW-1234: Web Platform Team ────────────────────────────────┐
│ 📋 Description: Manages web application platform infrastructure│
│ 📅 Created: 2024-01-15 by bthompso | Members: 12 | Assets: 23 │
├─ 👥 MEMBERSHIP ───────────────────────────────────────────────┤
│ 🟢 bthompso (owner, on-call)     │ 🔵 sarah.dev (admin)        │
│ ⚪ john.ops (member)              │ 🟡 temp.user (temp)         │
├─ 🚨 ON-CALL ROTATION ────────────────────────────────────────┤
│ Current: bthompso (until Fri 6PM) | Next: sarah.dev (Fri-Sun) │
├─ 📦 OWNED ASSETS (showing 5 of 23) ──────────────────────────┤
│ web-app/dashboard    │ service/api-gateway  │ db/user-store    │
└────────────────────────────────────────────────────────────────┘
```

#### Membership List View
```
┌─ CREW MEMBERSHIPS: bthompso@company.com ──────────────────────┐
│ CREW ID    │ CREW NAME              │ ROLE   │ ON-CALL │ SINCE │
├────────────┼────────────────────────┼────────┼─────────┼───────┤
│ CREW-1234  │ Web Platform Team      │ owner  │   🟢    │ 2024  │
│ CREW-2345  │ Security Response      │ member │   ⚪    │ 2024  │
│ CREW-3456  │ Database Operations    │ admin  │   ⚪    │ 2023  │
└────────────┴────────────────────────┴────────┴─────────┴───────┘
```

### Color Scheme Definition

```go
type ColorScheme struct {
    // Role-based colors
    OwnerColor   lipgloss.Color // Bright green
    AdminColor   lipgloss.Color // Blue
    MemberColor  lipgloss.Color // Light gray
    TempColor    lipgloss.Color // Yellow
    RemovedColor lipgloss.Color // Red

    // Status indicators
    OnCallColor  lipgloss.Color // Bright green
    OffCallColor lipgloss.Color // Light gray

    // UI elements
    BorderColor  lipgloss.Color // Dark gray
    HeaderColor  lipgloss.Color // Purple
    TextColor    lipgloss.Color // White
    AccentColor  lipgloss.Color // Bright blue
}
```

## 🔄 Command Flow Architecture

### Subcommand Handler Pattern

```go
type SubcommandHandler func(cmd *cobra.Command, args []string) error

// Standard handler registration
var SubcommandMap = map[string]SubcommandHandler{
    "create":     handleCreate,
    "details":    handleDetails,
    "membership": handleMembership,
    "assets":     handleAssets,
    "oncall":     handleOnCall,
    "transfer":   handleTransfer,
    "manage":     handleManage,
    "owner":      handleOwner,
}

// Main command router
func (c *CrewsCommand) Execute(cmd *cobra.Command, args []string) error {
    if len(args) == 0 {
        return c.showHelp(cmd)
    }

    // Check for explicit subcommand
    if handler, exists := SubcommandMap[args[0]]; exists {
        return handler(cmd, args[1:])
    }

    // Use smart parameter detection
    resolver := NewParameterResolver()
    handler, err := resolver.Resolve(args[0])
    if err != nil {
        return fmt.Errorf("could not resolve parameter: %w", err)
    }

    return handler(cmd, args)
}
```

### Interactive Workflow Pattern

```go
// Guided creation workflow (similar to engx create)
func handleCreate(cmd *cobra.Command, args []string) error {
    workflow := NewCrewCreationWorkflow()

    // Step 1: Basic information
    crewInfo, err := workflow.CollectBasicInfo()
    if err != nil {
        return err
    }

    // Step 2: Initial membership
    members, err := workflow.CollectInitialMembers()
    if err != nil {
        return err
    }

    // Step 3: On-call setup (optional)
    onCallConfig, err := workflow.CollectOnCallConfig()
    if err != nil {
        return err
    }

    // Step 4: Confirmation and creation
    crew, err := workflow.CreateCrew(crewInfo, members, onCallConfig)
    if err != nil {
        return err
    }

    // Step 5: Display results
    renderer := NewCrewDetailsRenderer()
    output, err := renderer.Render(crew, getTerminalWidth())
    if err != nil {
        return err
    }

    fmt.Print(output)
    return nil
}
```

## 🔍 Data Access Layer

### Simulation Data Strategy

For POC implementation, use in-memory data structures that simulate real backend:

```go
type SimulationDataStore struct {
    crews       map[string]*Crew
    assets      map[string]*AssetOwnership
    userIndex   map[string][]string // user -> crew IDs
    assetIndex  map[string]string   // asset -> crew ID
}

// CRUD operations interface
type CrewDataStore interface {
    GetCrew(id string) (*Crew, error)
    GetCrewsByUser(userID string) ([]*Crew, error)
    GetCrewByAsset(assetURN string) (*Crew, error)
    CreateCrew(crew *Crew) error
    UpdateCrew(crew *Crew) error
    DeleteCrew(id string) error

    // Membership operations
    AddMember(crewID string, member *Member) error
    UpdateMemberRole(crewID, userID string, role MemberRole) error
    RemoveMember(crewID, userID string) error

    // Asset operations
    GetAssetsByCrewID(crewID string) ([]*AssetOwnership, error)
    TransferAssetOwnership(assetURN, fromCrewID, toCrewID string) error
}
```

### Sample Data Population

```go
func PopulateSimulationData() *SimulationDataStore {
    store := NewSimulationDataStore()

    // Sample crews
    store.crews["CREW-1234"] = &Crew{
        ID:          "CREW-1234",
        VanityName:  "Web Platform Team",
        Description: "Manages web application platform infrastructure",
        Members: []Member{
            {UserID: "bthompso", Email: "bthompso@company.com", Role: RoleOwner, IsOnCall: true},
            {UserID: "sarah.dev", Email: "sarah.dev@company.com", Role: RoleAdmin},
            {UserID: "john.ops", Email: "john.ops@company.com", Role: RoleMember},
        },
        OwnedAssets: []string{
            "asset://web-app/dashboard",
            "asset://service/api-gateway",
            "asset://database/user-store",
        },
    }

    // Additional crews for realistic scenarios...

    return store
}
```

## 🧪 Testing Strategy

### Unit Test Coverage

```go
// Parameter resolution testing
func TestParameterResolver(t *testing.T) {
    resolver := NewParameterResolver()

    tests := []struct {
        input    string
        expected string
    }{
        {"CREW-1234", "details"},
        {"user@domain.com", "membership"},
        {"asset://web-app/test", "owner"},
        {"invalid", "search"},
    }

    for _, test := range tests {
        handler, _ := resolver.Resolve(test.input)
        assert.Equal(t, test.expected, handler.Name())
    }
}

// Renderer testing
func TestCrewDetailsRenderer(t *testing.T) {
    renderer := NewCrewDetailsRenderer()
    crew := &Crew{ID: "CREW-1234", VanityName: "Test Team"}

    output, err := renderer.Render(crew, 80)
    assert.NoError(t, err)
    assert.Contains(t, output, "CREW-1234")
    assert.Contains(t, output, "Test Team")
}
```

### Integration Test Scenarios

```go
func TestEndToEndWorkflows(t *testing.T) {
    // Test complete command execution
    cmd := NewCrewsCommand()

    // Test details lookup
    err := cmd.Execute(nil, []string{"CREW-1234"})
    assert.NoError(t, err)

    // Test membership lookup
    err = cmd.Execute(nil, []string{"bthompso@company.com"})
    assert.NoError(t, err)

    // Test asset ownership lookup
    err = cmd.Execute(nil, []string{"asset://web-app/dashboard"})
    assert.NoError(t, err)
}
```

## 🚀 Implementation Phases

### Phase 1: Core Infrastructure
1. **Plugin Registration**: Cobra command setup and registration
2. **Parameter Resolution**: Smart detection system implementation
3. **Basic Data Models**: Core entity definitions
4. **Simulation Data**: Sample data for testing

### Phase 2: Primary Commands
1. **Details Command**: Crew information display
2. **Membership Command**: User crew listing
3. **Owner Command**: Asset ownership queries
4. **Basic Rendering**: Terminal output formatting

### Phase 3: Management Operations
1. **Create Command**: Guided crew creation workflow
2. **Manage Command**: Crew editing capabilities
3. **Transfer Command**: Ownership transfer operations
4. **Advanced Rendering**: Interactive elements and improved layouts

### Phase 4: Advanced Features
1. **On-call Management**: Rotation scheduling and status
2. **Access Delegation**: Permission granting system
3. **Asset Management**: Bulk asset operations
4. **Search and Discovery**: Enhanced parameter handling

---

**Technical Design Status**: ✅ Complete
**Next Phase**: Implementation
**Architecture Review**: HUM LEAD approval required