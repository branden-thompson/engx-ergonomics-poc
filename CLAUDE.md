# ENGX Ergonomics POC - Project Instructions
**🛩️ BRTOPS v1.1.001 PRODUCTION** - Close Air Support for your Agentic-enabled Teams

## PROJECT CLASSIFICATION
- **Type**: LEVEL-1 SEV-0 POC
- **Name**: engx-ergonomics-poc
- **Focus**: Terminal-based simulation of engineering productivity tools and workflows
- **Status**: PHASE 6 COMPLETE - Verbosity-Aware TUI Rendering System Implemented

## BRTOPS FRAMEWORK ACTIVE
**CRITICAL**: This project operates under BRTOPS protocols with SEV-0 quality requirements.

### CORE MISSION
Create a proof-of-concept terminal application that SIMULATES engineering productivity tools and workflows. Focus on human-computer interaction patterns, TUI design, and command structure - NOT actual functionality implementation.

#### MISSION SCOPE EVOLUTION
**Original Focus**: React application creation simulation
**Expanded Scope**: Comprehensive engineering productivity tool simulation including:
- **Application Archetypes**: React Apps, CLI tools, headless APIs/services, etc.
- **Asset Management**: Ownership determination via external 'crews' services
- **Workflow Management**: Approval/rejection flows, ownership transfers, exception handling
- **Environment Operations**: Production deployment status, environment management
- **Template Systems**: Extensible archetype and template discovery patterns

React application creation remains a **core use case** demonstrating the interaction patterns, with the plugin architecture enabling simulation of broader engineering tool categories.

### CURRENT IMPLEMENTATION STATUS

#### ✅ COMPLETED PHASES
1. **Phase 1**: Core project creation with TUI and interactive prompts
2. **Phase 2**: Enhanced user experience with AAR and workflow analytics
3. **Phase 3**: Template discovery system and CLI interaction analytics
4. **Phase 4**: TUI component system and footer styling enhancements
5. **Phase 5**: Terminal-width aware progress bar system
6. **Phase 6**: Verbosity-aware TUI rendering system

#### 🎯 KEY ACHIEVEMENTS
- **Template Discovery**: Complete React template browsing system (`engx templates`)
- **Analytics Framework**: CLI interaction tracking and pattern detection (`engx analytics`)
- **Professional Styling**: Consistent lipgloss-based terminal aesthetics
- **Modular Architecture**: Template system is completely removable as requested
- **Workflow Simulation**: Realistic engineering productivity tool interaction patterns
- **Plugin Architecture**: Extensible foundation for diverse engineering tool simulations
- **TUI Component System**: Archetype-aware application component display and animation
- **Dynamic Footer Styling**: Language-specific color schemes with proper alignment
- **Terminal-Width Aware Progress Bars**: Adaptive progress bars that fill available terminal width
- **Verbosity-Aware TUI Rendering**: Comprehensive verbosity control system with 5 distinct output levels

### ARCHITECTURE OVERVIEW

#### Core Systems
1. **Create Command** (`plugins/create/`) - Interactive project creation simulation
2. **Template Discovery** (`pkg/common/templates*.go`) - Lightweight template browsing
3. **Analytics System** (`pkg/common/analytics*.go`) - CLI usage pattern tracking
4. **Plugin Infrastructure** (`pkg/common/`) - Extensible command architecture

#### Command Structure
```bash
# Core archetype creation (React apps as primary example)
engx create MyApp [--dev-only] [--template <id>]

# Template/Archetype discovery (Phase 3)
engx templates {list|search|info|recommended|complexity|stats}

# Analytics & insights (Phase 3)
engx analytics {status|summary|details|patterns|stats|export}

# Development tools
engx dev {plugin|config|hotreload|generate}

# Future engineering productivity simulations:
# engx ownership {check|transfer|crews} <asset>
# engx approval {request|review|approve|reject} <workflow>
# engx deploy {status|logs|health} <environment>
# engx archetype {api|cli|service|webapp} <name>
```

### KEY REQUIREMENTS ✅
1. **Terminal Application**: ✅ Rich TUI with lipgloss styling
2. **Simulation Focus**: ✅ UX patterns over actual scaffolding
3. **Performance Critical**: ✅ Fast, responsive interface with in-memory data
4. **npx-like Experience**: ✅ Intuitive workflow with modern CLI patterns
5. **Human Ergonomics**: ✅ Optimized developer experience and interaction flows

### PHASE 3 REORIENTATION LEARNINGS
- **Original Plan**: Complex plugin marketplace infrastructure
- **Reoriented To**: Lightweight template discovery focused on CLI simulation
- **Key Insight**: Simulation fidelity > technical complexity for educational goals
- **Result**: Valuable UX patterns for React development tooling reference

### CURRENT CAPABILITIES

#### Template Discovery System
- Browse all React templates with visual hierarchy
- Search templates by framework, language, features
- Detailed template information and recommendations
- Complexity-based filtering (Beginner → Enterprise)
- Professional terminal styling with consistent UX

#### Analytics & Pattern Detection
- Real-time CLI interaction tracking
- Workflow pattern detection and analysis
- Session metrics and usage statistics
- Data export for external analysis
- Command sequence analysis

#### Development Experience
- Hot-reload plugin development
- Interactive configuration management
- Plugin generation and scaffolding
- Comprehensive testing and validation

### QUALITY GATES (SEV-0) ✅
- ✅ ALL folders + 80% optional files completed
- ✅ Comprehensive documentation at all phases
- ✅ Full validation and testing protocols
- ✅ Professional polish and consistent styling
- ✅ Modular architecture with removable components

### TECHNICAL LEARNINGS (Phase 4)

#### TUI Component System Implementation
- **Archetype Detection**: Component names used to infer archetype instead of passing through constructors
- **Dynamic Component Manager**: `NewComponentManagerForArchetype()` creates installation plans based on actual components
- **Animation State Management**: Components transition through "queued" → "installing" → "installed" states
- **6-Section Architecture**: Core Technologies, EngX Integrations, Quality & Testing, User Analytics, Agentic Capabilities, Deployment Capabilities

#### Footer Alignment Solution
- **Color Code Length Issue**: ANSI color codes caused misalignment in right-justified footer text
- **Dual Return Pattern**: Functions return both colored display string and plain text for length calculations
- **Padding Calculation**: Use plain text length for spacing, display colored version for output
- **Language-Specific Colors**: Framework (grey) + Language (color-coded: TypeScript=Bright Blue, GOLANG=Purple, etc.)

#### Terminal-Width Aware Progress Bars (Phase 5)
- **Dynamic Width Calculation**: Progress bars automatically fill available terminal width
- **Minimum Width Safety**: 10-character minimum bar width prevents terminal width panics
- **FillMode System**: Existing modular progress bar infrastructure utilized with fill mode
- **Right-Alignment Fix**: Manual percentage padding (PercentagePad: 0) for precise terminal alignment
- **Status Text Alignment**: Proper right-alignment of "✓ Done" and "⠋ Running..." with divider lines
- **Space Calculation**: Precise calculation of used space for both progress bars and status text
- **Cross-Terminal Support**: Tested across 40, 60, 80, and 120-column terminal widths

#### Verbosity-Aware TUI Rendering (Phase 6)
- **5-Level Verbosity System**: Quiet, Concise, Default, Verbose, Debug levels with distinct outputs
- **Conditional Section Rendering**: Helper methods control visibility of steps list, footer info, and application components
- **Quiet Mode Optimization**: Minimal output with timing information and single separator line
- **Modular Timing Display**: `renderTimingInfo()` method provides timing data for quiet mode without footer overhead
- **Separator Logic Fix**: Restructured conditional separator placement to prevent duplicate lines
- **Verbosity Config Integration**: `VerbosityConfig` passed through all renderer constructors for consistent behavior
- **Section Control Methods**: `shouldShowStepsList()`, `shouldShowFooterInfo()`, `shouldShowApplicationComponents()` for granular control

#### Architecture Patterns Discovered
- **Component Inference**: Detect archetype from component presence rather than explicit passing
- **Color-Aware Formatting**: Separate display formatting from layout calculations
- **Archetype-Agnostic Design**: Single renderer handles all archetype types dynamically
- **Terminal-Responsive UI**: Dynamic progress bar sizing based on available terminal width

### REFERENCE VALUE
This POC successfully demonstrates:
- **CLI Ergonomics**: Natural command discovery and interaction patterns for engineering tools
- **Template/Archetype Selection UX**: Intuitive browsing workflows extensible to multiple domains
- **Terminal Aesthetics**: Professional styling with lipgloss framework
- **Analytics Foundation**: Usage tracking for data-driven UX improvements
- **Modular Plugin Architecture**: Easy to extend for diverse engineering productivity simulations
- **Workflow Simulation Patterns**: Interaction models applicable to ownership, approval, deployment workflows
- **Dynamic TUI Systems**: Archetype-aware component rendering with proper color and alignment handling

---
**BRTOPS Version**: 1.1.001
**Project Status**: PHASE 6 COMPLETE - Verbosity-Aware TUI Rendering System with Comprehensive Output Control Implemented
**Next Phase**: Expand simulation scope with additional engineering productivity command categories and advanced TUI interactions