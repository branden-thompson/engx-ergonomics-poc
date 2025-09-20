# ENGX Ergonomics POC

> **Terminal-based simulation of React web application creation**
> Focus: Human-computer interaction patterns, TUI design, and command structure

## Overview

ENGX POC is a proof-of-concept terminal application that **simulates** the creation of React-based web applications. The primary goal is to demonstrate and test terminal UX patterns, CLI ergonomics, and human-computer interaction flows that could be used in real development tooling.

**Key Point**: This is a simulation tool, not actual application scaffolding. It focuses on CLI interaction patterns and terminal user experience rather than generating real React projects.

## Architecture

### Core Systems

1. **Create Command** (`plugins/create/`)
   - Interactive React project creation simulation
   - TUI-based prompts and configuration
   - Multi-step progress visualization
   - After Action Reports (AAR) for workflow analysis

2. **Template Discovery** (`pkg/common/templates*.go`)
   - Lightweight React template browsing system
   - Search and filtering capabilities
   - Ergonomic CLI patterns for template selection
   - **Completely modular and removable**

3. **Analytics System** (`pkg/common/analytics*.go`)
   - CLI interaction tracking and pattern detection
   - Workflow analysis and session metrics
   - Export capabilities for external analysis
   - Real-time usage statistics

4. **Plugin Infrastructure** (`pkg/common/`)
   - Hot-reload development system
   - Configuration management
   - Extensible command architecture

### Design System

- **Styling**: Consistent lipgloss-based terminal styling
- **Colors**: Professional color palette optimized for terminal use
- **Layout**: Structured information hierarchy with visual polish
- **Typography**: Clear visual hierarchy with emojis and formatting

## Available Commands

### Core Application Creation
```bash
# Interactive project creation with TUI
./dist/engx create MyApp

# Development-only mode (faster, no deployment simulation)
./dist/engx create MyApp --dev-only

# Template-based creation
./dist/engx create MyApp --template react-typescript-vite
```

### Template Discovery
```bash
# Browse all available React templates
./dist/engx templates list

# Search templates by criteria
./dist/engx templates search typescript

# Get detailed template information
./dist/engx templates info react-typescript-vite

# Show only recommended templates
./dist/engx templates recommended

# Filter by complexity level
./dist/engx templates complexity beginner

# View template statistics
./dist/engx templates stats
```

### Analytics & Usage Patterns
```bash
# View session analytics summary
./dist/engx analytics summary

# Comprehensive usage analysis
./dist/engx analytics details

# Workflow pattern detection
./dist/engx analytics patterns

# Command usage statistics
./dist/engx analytics stats

# Export analytics data
./dist/engx analytics export session-data.json

# System status
./dist/engx analytics status
```

### Development Tools
```bash
# Plugin management
./dist/engx dev plugin list
./dist/engx dev hotreload enable

# Configuration management
./dist/engx dev config list
./dist/engx dev config edit create

# Plugin generation
./dist/engx dev generate MyPlugin
```

## Getting Started

### Build and Run
```bash
# Build the application
go build -o dist/engx ./cmd/engx

# Run basic help
./dist/engx --help

# Try template discovery
./dist/engx templates list

# Create a simulated project
./dist/engx create TestApp --dev-only
```

### Development Mode
```bash
# Enable hot-reload for plugin development
./dist/engx dev hotreload enable

# Generate example plugins
./dist/engx dev example basic
./dist/engx dev example advanced
```

## Key Features Demonstrated

### CLI Ergonomics
- **Natural Workflows**: Intuitive command discovery paths
- **Progressive Disclosure**: Information revealed as needed
- **Contextual Help**: Relevant suggestions throughout interface
- **Error Guidance**: Clear error messages with next steps

### Template Selection UX
- **Visual Hierarchy**: Recommended (⭐) and popular (🔥) indicators
- **Multiple Views**: List, search, detailed info, filtered views
- **Smart Search**: Relevance-based ranking and suggestions
- **Quick Access**: Shortcuts for common operations

### Analytics & Insights
- **Workflow Patterns**: Detection of common command sequences
- **Usage Metrics**: Command frequency and performance tracking
- **Session Analysis**: Real-time session statistics
- **Export Capabilities**: Data export for external analysis

### Development Experience
- **Hot Reload**: Real-time plugin development
- **Configuration Management**: Flexible plugin configuration
- **Extensible Architecture**: Easy to add new commands
- **Testing Support**: Built-in testing and validation tools

## Architecture Decisions

### Phase 3 Reorientation
Originally planned as complex plugin infrastructure, Phase 3 was reoriented toward simulation goals:

- **Template Discovery**: Replaced marketplace with lightweight React template system
- **Analytics**: Added CLI interaction tracking instead of plugin performance metrics
- **Modular Design**: Template system is completely removable as requested
- **Styling Consistency**: Unified lipgloss styling across all components

See `docs/phase3-learnings.md` for detailed architectural decisions and insights.

### Simulation vs. Reality
This tool prioritizes **simulation fidelity** (realistic workflows) over **technical complexity** (actual functionality):

- Template system uses in-memory data for instant response
- Analytics track real CLI usage patterns
- TUI demonstrates modern terminal interface patterns
- Commands show realistic output without actual file generation

## File Structure

```
engx-ergonomics-poc/
├── cmd/engx/               # Main CLI application
├── pkg/common/             # Core shared functionality
│   ├── templates*.go       # Template discovery system
│   ├── analytics*.go       # CLI usage analytics
│   └── *.go               # Plugin infrastructure
├── plugins/                # Plugin implementations
│   ├── create/            # Project creation plugin
│   └── */                 # Other plugins
├── internal/               # Internal packages
│   ├── tui/               # Terminal UI components
│   ├── prompts/           # Interactive prompts
│   └── config/            # Configuration management
└── docs/                  # Documentation and learnings
```

## Development Philosophy

1. **Simulation First**: Focus on UX patterns over real functionality
2. **CLI Ergonomics**: Optimize for developer experience and interaction patterns
3. **Visual Polish**: Professional terminal aesthetics with consistent styling
4. **Modular Design**: Components should be easily removable/replaceable
5. **Analytics Driven**: Capture interaction data to improve UX
6. **Documentation Heavy**: Record decisions and learnings for reference

## Success Metrics

✅ **Template Discovery**: Natural, intuitive template selection workflows
✅ **CLI Consistency**: Unified styling and interaction patterns
✅ **Analytics Foundation**: Comprehensive usage tracking and pattern detection
✅ **Modular Architecture**: Easy to modify, extend, or remove components
✅ **Professional Polish**: Terminal interface that feels modern and responsive

## Future Considerations

- Interactive template selection TUI
- Advanced workflow recording and playback
- Performance metrics and optimization
- Multi-language template support
- Enhanced analytics visualization

---

**Purpose**: Reference implementation for terminal-based React development tooling
**Status**: Phase 3 Complete - Template discovery and analytics systems integrated
**Next**: Visual consistency improvements and enhanced interaction patterns