# Requirements - Crew Command Plugin

## MAJOR FEATURE ENHANCEMENT | LEVEL-1 SEV-0
**Feature**: `crew-command-plugin`
**Date**: 2025-09-23
**Status**: RCC/PLAN Phase Complete

## 🎯 Core Mission

Implement comprehensive crew management system focusing on people groups, their relationships, ownership, and responsibilities within the EngX ecosystem.

## 📋 Functional Requirements

### 1. Smart Parameter Detection System
The root command `engx crews [parameter]` must intelligently detect parameter types:

- **CrewID Detection**: `CREW-1234` → Route to crew details view
- **User LDAP/Email**: `bthompso@company.com` → Show user's crew memberships
- **Catalog Asset ID/URN**: `asset://web-app/my-app` → Show owning crew details
- **Fallback Handling**: Invalid/ambiguous parameters → Search/suggestion mode

### 2. Complete Subcommand Suite

#### Core Management Commands
```bash
engx crews create                    # Guided crew creation workflow
engx crews manage <crewID>          # Edit crew details (admin/owner only)
engx crews details <crewID>         # Detailed crew information
engx crews transfer <crewID> --to <user> # Transfer ownership
```

#### Membership Operations
```bash
engx crews membership <user>        # Show user's crew memberships
engx crews membership change <user> --to [admin|member|auto|temp|removed]
```

#### Ownership & Asset Management
```bash
engx crews owner <assetID>          # Show asset's owning crew
engx crews assets <crewID>          # List crew's owned assets
engx crews on-call <crewID>         # Current on-call members
```

### 3. Data Entity Requirements

#### Crew Entity Structure
- **CrewID**: Unique identifier (CREW-####)
- **Vanity Name**: Human-readable team name
- **Description**: Purpose and responsibilities
- **Membership Roster**: Users with role assignments
- **On-call Rotation**: Current and scheduled on-call members
- **Owned Assets**: Catalog items under crew ownership
- **Access Delegations**: Permissions granted to external users

#### Member Role Hierarchy
- **admin**: Full crew management privileges
- **member**: Standard crew member access
- **temp**: Temporary access with expiration
- **auto**: Inherited access from asset ownership
- **removed**: Explicitly revoked access

### 4. User Experience Requirements

#### Smart Discovery Workflows
1. **Asset Ownership Discovery**: "Who owns this service?" → `engx crews owner service-name`
2. **User Membership Lookup**: "What teams is Sarah on?" → `engx crews membership sarah@company.com`
3. **On-call Resolution**: "Who's on-call for platform team?" → `engx crews on-call CREW-1234`
4. **Team Creation**: "Create team for new project" → `engx crews create` (guided flow)

#### Terminal UX Standards
- **Width Awareness**: Support 40+ column terminals with responsive layout
- **Color Consistency**: Use established lipgloss color scheme
- **Interactive Elements**: Guided prompts for complex operations
- **Status Indicators**: Visual cues for membership roles, on-call status
- **Progressive Disclosure**: Detailed information available on demand

## 🔗 Integration Requirements

### Relationship with Catalog Command
- **Crews**: Manage "people groups" and their relationships
- **Catalog**: Manage "things/assets" and their metadata
- **Bidirectional Queries**: Find crew from asset, find assets from crew
- **Ownership Chain**: Clear asset → crew → member relationship tracking

### Plugin Architecture Compliance
- **Cobra Integration**: Standard command registration patterns
- **Plugin Discovery**: Automatic registration via common plugin system
- **Error Handling**: Consistent error reporting and fallback behavior
- **Testing Framework**: Unit and integration tests for all commands

## 🎨 Visual Design Requirements

### Terminal Output Standards
- **Consistent Formatting**: Align with existing engx command styling
- **Dynamic Layout**: Adjust column widths based on terminal size
- **Status Indicators**: Color-coded role badges and on-call indicators
- **Table Rendering**: Professional tabular data display
- **Detail Views**: Card-style layouts for comprehensive information

### Color Scheme Guidelines
- **Primary Actions**: Bright colors for actionable items
- **Status Indicators**: Green (active/admin), Blue (member), Yellow (temp), Red (removed)
- **Information Hierarchy**: Consistent typography and spacing
- **Brand Consistency**: Maintain EngX visual identity

## ⚡ Performance Requirements

### Response Time Expectations
- **Command Discovery**: < 100ms for parameter detection
- **Data Retrieval**: < 200ms for crew information queries
- **Interactive Flows**: Responsive input handling for guided workflows
- **Terminal Rendering**: Smooth output with no visible delays

### Scalability Considerations
- **Large Teams**: Support crews with 100+ members
- **Asset Catalogs**: Handle queries across thousands of assets
- **Search Performance**: Efficient filtering and lookup operations
- **Memory Usage**: Minimal footprint for CLI tool requirements

## 🔒 Security & Access Requirements

### Authentication Integration
- **LDAP Integration**: Resolve usernames to full identity information
- **Permission Validation**: Verify user authority for management operations
- **Audit Logging**: Track crew membership and ownership changes
- **Data Privacy**: Respect user information access policies

### Authorization Controls
- **Admin-Only Operations**: Crew creation, member role changes, transfers
- **Member Permissions**: View crew details, list memberships
- **Public Information**: Basic crew information for asset ownership queries
- **Delegation Tracking**: Clear chain of access permissions

## 📊 Success Criteria

### Functional Validation
- ✅ All subcommands execute successfully with proper error handling
- ✅ Smart parameter detection correctly routes to appropriate functions
- ✅ Interactive workflows provide clear guidance and validation
- ✅ Terminal output renders correctly across different window sizes

### User Experience Validation
- ✅ Intuitive command discovery and usage patterns
- ✅ Clear visual hierarchy and information presentation
- ✅ Responsive design across terminal width breakpoints
- ✅ Consistent behavior with existing EngX command suite

### Technical Validation
- ✅ Plugin architecture follows established patterns
- ✅ Data models support all required operations
- ✅ Error handling provides meaningful feedback
- ✅ Performance meets response time requirements

## 🚀 Deployment Readiness

### Quality Gates
- **Code Coverage**: > 80% test coverage for all functions
- **Documentation**: Complete user guide and technical documentation
- **Error Handling**: Graceful degradation for all failure scenarios
- **Backward Compatibility**: No breaking changes to existing commands

### Integration Testing
- **Command Registration**: Proper integration with engx plugin system
- **Cross-Command Compatibility**: Seamless operation with catalog commands
- **Terminal Compatibility**: Verified operation across different terminal types
- **Data Consistency**: Accurate relationship mapping between crews and assets

---

**Requirements Status**: ✅ Complete
**Next Phase**: Technical Architecture Design
**Authorization**: HUM LEAD approval required for implementation