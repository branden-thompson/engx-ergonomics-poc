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

### TUI COMPONENT ARCHITECTURE & SHARED COMPONENTS

#### 🚨 CRITICAL: ALWAYS USE SHARED COMPONENTS
**NEVER re-implement headers, footers, separators, data tables, progress bars, or other TUI components when creating new views. The project provides a comprehensive shared component library that MUST be used.**

#### Shared Component Library (`internal/tui/components/`)

1. **HeaderFooterComponent** (`header_footer.go`)
   - Terminal-width aware headers and footers with customizable labels
   - ANSI-aware width calculations (handles colored text properly)
   - Methods: `RenderHeader()`, `RenderFooter()`, `RenderHeaderWithColors()`, `RenderFooterWithColors()`, `RenderSeparator()`
   - **Usage**: `NewHeaderFooterComponent(terminalWidth)`

2. **DataTableRow** (`data_table_row.go`)
   - Professional data table rendering with column definitions
   - Automatic width calculation and alignment across terminal sizes
   - Badge support with right-alignment and color coordination
   - **Architecture**: Use `DataTableLayout` struct for clean data/presentation separation
   - **Usage**: `NewDataTableRowFromLayout(terminalWidth, layout)`

3. **BadgeAligner** (`badge_aligner.go`)
   - Right-aligned badges with proper spacing and gutter control
   - Prevents text wrapping in narrow terminals
   - **Usage**: `NewBadgeAligner(terminalWidth).AlignRight(content, badge, gutter)`

4. **ProgressBarRenderer** (`progress_bar_renderer.go`)
   - Terminal-width aware progress bars with multiple display modes
   - Dynamic sizing that fills available terminal width
   - **Usage**: Various constructors for different display patterns

#### Design System (`internal/tui/design/foundation.go`)

**ALWAYS use design system constants for consistent styling:**

- **Colors**: Semantic color constants (e.g., `ColorBrightWhite`, `ColorDarkGray`, `ColorPrimaryOnCall`)
- **Spacing**: Consistent spacing scale (`SpaceXS`, `SpaceSM`, `SpaceMD`, etc.)
- **Typography**: Text style presets (`StyleTitle`, `StyleHeader`, `StyleMuted`, etc.)
- **Layout**: Standard widths and column sizing constants

#### Color Implementation Best Practices

1. **ANSI Format**: Use basic ANSI format `\033[%sm` for maximum compatibility
2. **Width Calculations**: Always strip ANSI codes when calculating display width
3. **Color Separation**: Apply colors AFTER formatting to avoid width issues
4. **Semantic Colors**: Use design system constants for consistent color meaning

#### 🚨 CRITICAL DATA TABLE PATTERNS - PROVEN SOLUTIONS

**WARNING**: Multiple attempts at customizing table patterns have led to repeated alignment and wrapping issues. **ALWAYS USE THE EXACT PATTERN FROM WORKING COMMANDS.**

#### ✅ PROVEN PATTERN: Crews Detail Column Structure

**ALWAYS use this EXACT column structure for table views with badges/indicators:**

```go
// ✅ PROVEN: Separate prefix column for indicators (NEVER combine with other columns)
columnDefs := []components.ColumnDefinition{
    {Name: "prefix", Header: "", Width: 3, Alignment: "left", Color: ""},           // ⏺ indicator space
    {Name: "number", Header: "##", Width: 3, Alignment: "left", Color: ""},        // Row numbers
    {Name: "name", Header: "CREW NAME", Width: 0, MinWidth: 20, Alignment: "left", Color: ""}, // FILL column
    {Name: "id", Header: "ID", Width: 16, Alignment: "left", Color: ""},           // Fixed width
    {Name: "type", Header: "TYPE", Width: 12, Alignment: "left", Color: ""},       // Fixed width
    {Name: "role", Header: "ROLE", Width: 8, Alignment: "left", Color: ""},        // Fixed width
}
```

#### ✅ PROVEN PATTERN: Badge Alignment (NO table columns for badges)

```go
// ✅ PROVEN: Use BadgeAligner AFTER table formatting (not as table column)
r.dataTable = components.NewDataTableRowFromLayout(width, layout)
r.dataTable.SetBadge("(On-Call 1)", 4).SetBadgeRenderMode(false) // Reserve space only

// ✅ PROVEN: Manual badge alignment per row
rowLine := r.dataTable.FormatDataRow(rowData)
if needsBadge {
    coloredBadge := r.applyColor(badge, badgeColor)
    rowLine = r.badgeAligner.AlignRight(rowLine, coloredBadge, 4)
}
```

#### ✅ PROVEN PATTERN: Header Alignment

```go
// ✅ PROVEN: Treat headers exactly like data rows with same prefix structure
headerRow := []string{
    "   ", // EXACT same 3-space prefix as data rows
    r.applyColor("##", design.ColorTableHeader),
    r.applyColor("CREW NAME", design.ColorTableHeader),
    r.applyColor("ID", design.ColorTableHeader),
    r.applyColor("TYPE", design.ColorTableHeader),
    r.applyColor("ROLE", design.ColorTableHeader),
}
headerLine := r.dataTable.FormatDataRow(headerRow)
```

#### 🚫 ANTI-PATTERNS THAT CAUSE REPEATED ISSUES

1. **🚫 NEVER**: Customize the proven column structure
   - Adding badge as table column → Always causes wrapping
   - Combining indicator with number column → Breaks alignment
   - Custom width calculations → Always breaks on different terminal sizes

2. **🚫 NEVER**: Use `len()` for width calculations with ANSI codes
   ```go
   // 🚫 WRONG: Causes badge wrapping and misalignment
   badgeWidth := len(badge) // Counts ANSI codes as characters

   // ✅ CORRECT: Use getDisplayWidth() for ANSI-aware calculations
   badgeWidth := b.getDisplayWidth(badge) // Strips ANSI codes
   ```

3. **🚫 NEVER**: Apply colors before width calculations
   - Always calculate widths with plain text first
   - Apply colors AFTER formatting is complete

#### ✅ REQUIRED FIXES: ANSI Color Support

**1. BadgeAligner ANSI Width Calculation:**
```go
// ✅ FIXED: badge_aligner.go now uses getDisplayWidth() for colored badges
badgeDisplayWidth := b.getDisplayWidth(badge) // Badge may contain ANSI codes
```

**2. 256-Color ANSI Sequence Format (CRITICAL):**
```go
// 🚨 CRITICAL: applyColor() must handle 256-color codes properly
func applyColor(text, colorCode string) string {
    if colorCode == "" || text == "" {
        return text
    }
    // Handle 256-color codes (3+ digits) vs basic ANSI codes (1-2 digits)
    if len(colorCode) >= 3 {
        return fmt.Sprintf("\033[38;5;%sm%s\033[%sm", colorCode, text, design.ColorReset)
    }
    return fmt.Sprintf("\033[%sm%s\033[%sm", colorCode, text, design.ColorReset)
}
```

**ROOT CAUSE**: 256-color codes like `"196"` and `"202"` require `\033[38;5;CODEm` format, not basic `\033[CODEm`. This was the reason colors weren't displaying despite correct logic.

#### DataTableLayout Pattern Implementation

1. **DataTableLayout Pattern**: Separate data from presentation
   ```go
   layout := &components.DataTableLayout{
       Columns: columnDefinitions,
       Data:    formattedData,
   }
   ```

2. **Fill Column Pattern**: Always use ONE fill column (Width: 0, MinWidth: X)
   ```go
   {Name: "name", Header: "CREW NAME", Width: 0, MinWidth: 20, Alignment: "left", Color: ""}
   ```

3. **Badge Space Reservation**: Use SetBadge() with SetBadgeRenderMode(false)
   ```go
   r.dataTable.SetBadge("(On-Call 1)", 4).SetBadgeRenderMode(false)
   ```

4. **Strategic Color Hierarchy**:
   - Headers: Light Purple (`ColorTableHeader`)
   - Names: White or semantic colors (on-call status matching)
   - Attributes: Dark Grey (`ColorAttribute`) for secondary information
   - Badges: Match associated icon colors (`ColorPrimaryOnCall`, `ColorSecondaryOnCall`)

#### 📋 NEW TABLE VIEW IMPLEMENTATION CHECKLIST

**Before implementing any new table view, follow this checklist to prevent repeated issues:**

1. **✅ Column Structure**:
   - [ ] Use EXACT crews detail column pattern (prefix + number + fill name + fixed columns)
   - [ ] Include separate 3-char prefix column for indicators (NEVER combine)
   - [ ] Use ONE fill column with `Width: 0, MinWidth: X`
   - [ ] All other columns have fixed widths

2. **✅ Header Implementation**:
   - [ ] Treat headers as data rows using `FormatDataRow()`
   - [ ] Include same 3-space prefix as data rows (`"   "`)
   - [ ] Apply `ColorTableHeader` to all header text

3. **✅ Badge Implementation**:
   - [ ] Use BadgeAligner AFTER table formatting (NEVER as table column)
   - [ ] Call `SetBadge().SetBadgeRenderMode(false)` for space reservation
   - [ ] Apply colors to badge before passing to `AlignRight()`
   - [ ] Use 4-character minimum gutter for badge alignment

4. **✅ ANSI Code Handling (CRITICAL)**:
   - [ ] Use `getDisplayWidth()` for all width calculations involving colored text
   - [ ] Apply colors AFTER formatting and width calculations are complete
   - [ ] Ensure BadgeAligner uses `getDisplayWidth(badge)` not `len(badge)`
   - [ ] **🚨 CRITICAL**: Use proper 256-color ANSI format in `applyColor()` function

5. **✅ Color Scheme**:
   - [ ] Use semantic color hierarchy (headers purple, attributes dark grey)
   - [ ] Match badge colors to indicator colors for visual consistency
   - [ ] Keep color usage minimal and purposeful

6. **✅ Data-Driven Status (NEW)**:
   - [ ] Use data model's status fields (e.g., `member.IsCurrentlyOnCall()`, `member.OnCallType`)
   - [ ] NEVER implement view-specific status logic
   - [ ] Ensure status data exists in core models before implementing view

7. **✅ Testing**:
   - [ ] Test across multiple terminal widths (40, 60, 80, 120 columns)
   - [ ] Verify no badge wrapping at any terminal width
   - [ ] Confirm header-data column alignment is perfect
   - [ ] Check that indicators and badges align properly
   - [ ] **🚨 CRITICAL**: Verify colors are actually displaying (not just ANSI codes in output)

**If ANY step fails or produces unexpected results, STOP and use the exact crews detail pattern without modification.**

#### ✅ DATA-DRIVEN STATUS PATTERNS (MEMBERSHIP VIEW LEARNINGS)

**ALWAYS use data model status fields instead of view-specific logic. This ensures consistency across commands and supports future features like `engx my oncall`.**

**1. ✅ PROVEN PATTERN: On-Call Status from Data Model**
```go
// ✅ CORRECT: Use data model's on-call information
isOnCall := member.IsCurrentlyOnCall()
var onCallColor string
if isOnCall {
    // Determine color based on OnCallType from data model
    if member.OnCallType == models.OnCallPrimary {
        onCallColor = design.ColorPrimaryOnCall // Primary on-call (deep red)
    } else {
        onCallColor = design.ColorSecondaryOnCall // Secondary/backup on-call (orange-red)
    }
}

// Apply to icon, name, and badge consistently
prefix = " " + r.applyColor("⏺", onCallColor) + " "
nameColor = onCallColor
badge = "(On-Call 1)" or "(On-Call 2)"
badgeColor = onCallColor
```

**2. 🚫 ANTI-PATTERN: View-Specific Status Logic**
```go
// 🚫 WRONG: Custom logic per view
func (r *Renderer) isOnCall(member *models.Member, crew *models.Crew) bool {
    if member.UserID == "bthompso" {
        return strings.Contains(crew.VanityName, "Universal") // View-specific logic
    }
    return false
}
```

**3. ✅ REQUIRED DATA MODEL FIELDS**
```go
type Member struct {
    IsOnCall   bool       `json:"is_oncall"`    // Current on-call status
    OnCallType OnCallType `json:"oncall_type"`  // Primary, backup, temp
    // ... other fields
}

// With helper method
func (m *Member) IsCurrentlyOnCall() bool {
    return m.IsActive() && m.IsOnCall
}
```

**4. ✅ SEMANTIC COLOR COORDINATION**
```go
// All three elements use the SAME color for consistency
var onCallColor string // Determined by OnCallType
- Icon: r.applyColor("⏺", onCallColor)
- Name: r.applyColor(crew.VanityName, onCallColor)
- Badge: r.applyColor("(On-Call 1)", onCallColor)
```

#### ✅ COMPLETE WORKING REFERENCE: Membership View Implementation

**This is the PROVEN, working implementation pattern from `membership_list.go` - use as reference for future views:**

```go
// 1. Data-driven status logic
isOnCall := member.IsCurrentlyOnCall()
var onCallColor string
if isOnCall {
    if member.OnCallType == models.OnCallPrimary {
        onCallColor = design.ColorPrimaryOnCall // "196"
    } else {
        onCallColor = design.ColorSecondaryOnCall // "202"
    }
}

// 2. Prefix column with indicator
prefix := "   " // 3 spaces default
if isOnCall {
    prefix = " " + r.applyColor("⏺", onCallColor) + " "
}

// 3. Name color coordination
nameColor := design.ColorLabel // Default white
if isOnCall && onCallColor != "" {
    nameColor = onCallColor // Match on-call color
}

// 4. Badge implementation
if member != nil && member.IsCurrentlyOnCall() {
    var badge string
    var badgeColor string
    if member.OnCallType == models.OnCallPrimary {
        badge = "(On-Call 1)"
        badgeColor = design.ColorPrimaryOnCall
    } else {
        badge = "(On-Call 2)"
        badgeColor = design.ColorSecondaryOnCall
    }
    coloredBadge := r.applyColor(badge, badgeColor)
    rowLine = r.badgeAligner.AlignRight(rowLine, coloredBadge, 4)
}

// 5. Critical: 256-color ANSI support
func (r *Renderer) applyColor(text, colorCode string) string {
    if colorCode == "" || text == "" {
        return text
    }
    // Handle 256-color codes (3+ digits) vs basic ANSI codes (1-2 digits)
    if len(colorCode) >= 3 {
        return fmt.Sprintf("\033[38;5;%sm%s\033[%sm", colorCode, text, design.ColorReset)
    }
    return fmt.Sprintf("\033[%sm%s\033[%sm", colorCode, text, design.ColorReset)
}
```

**Result**: Perfect color coordination with deep red (#196) for primary on-call, orange-red (#202) for secondary on-call, across icon, name, and badge.

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

#### Professional Data Tables & Color Systems (Crews Command - Phase 7)
- **DataTableLayout Architecture**: Clean separation of data from presentation using structured layout objects
- **ANSI-Aware Width Calculations**: HeaderFooterComponent enhanced with `getDisplayWidth()` method to handle colored text in width calculations
- **Header-Data Alignment**: Treating headers as data rows ensures perfect column alignment across all terminal widths
- **Strategic Color Hierarchy**: Professional color scheme with semantic meaning:
  - **Headers**: Light purple for distinction without overwhelming
  - **Primary Information**: White for names and key data
  - **Secondary Information**: Dark grey for attributes, timestamps, metadata
  - **Status Indicators**: Semantic colors (red/orange for on-call status)
  - **User Context**: Current user highlighting with green, other users in cyan
- **Badge Space Management**: Including badge columns in layout prevents text wrapping in narrow terminals
- **Color-Width Separation**: Apply ANSI colors AFTER formatting to avoid width calculation interference
- **On-Call Status Integration**: Icon, name, and badge color coordination for clear visual status communication

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
- **Professional Data Tables**: ANSI-aware width calculations with strategic color hierarchies for complex information display
- **Shared Component Library**: Reusable TUI components preventing reimplementation and ensuring design consistency

---
**BRTOPS Version**: 1.1.001
**Project Status**: PHASE 7 COMPLETE - Professional Data Tables & Color Systems with Comprehensive Shared Component Library
**Next Phase**: Expand simulation scope with additional engineering productivity command categories leveraging established TUI component architecture