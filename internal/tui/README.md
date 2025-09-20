# Terminal Design System

**Status**: Ready for extraction as standalone "Terminal Design and Component Language"

This directory contains a modular terminal UI design system that provides consistent styling, components, and layout utilities for terminal applications. It's designed to be easily extractable as a standalone library similar to bubble-tea and lipgloss.

## Architecture

### 📐 Foundation (`design/foundation.go`)
Core design tokens and constants:
- **Color System**: ANSI-based semantic color tokens
- **Spacing Scale**: Consistent spacing values (XS, SM, MD, LG, etc.)
- **Typography**: Text style presets (Title, Header, Body, etc.)
- **Layout Constants**: Standard component dimensions
- **Visual Elements**: Border characters and separator patterns

### 🧱 Components (`components/`)
Reusable UI components built on the foundation:

#### `table.go` - Data Tables
- **TableFormatter**: Flexible table rendering with exact spacing
- **TableColumn**: Column configuration (width, color, alignment)
- **Preset Tables**: Ready-to-use configurations
  - `NewArchetypeTable()` - Application archetype listings
  - `NewCommandTable()` - Command reference tables
  - `NewStatusTable()` - Status and monitoring displays

#### `layout.go` - Layout Components
- **Header**: Section headers with dividers and subtitles
- **Separator**: Horizontal separators with customizable styling
- **Panel**: Content panels with optional borders and padding
- **List**: Numbered and bulleted lists with consistent formatting

#### `text.go` - Text Styling
- **Text**: Styled text component with fluent API
- **Badge**: Status indicators and labels
- **Convenience Functions**: Quick styling utilities
  - `Title()`, `HeaderText()`, `Subheader()`, `Muted()`
  - `Command()`, `Success()`, `Warning()`, `Error()`
  - `Highlight()`, `Code()`, `Link()`

### 🎨 Legacy Styles (`styles/`)
Legacy lipgloss-based styles for backward compatibility:
- `colors.go` - Original color definitions
- `table.go` - Original table implementation (deprecated)

## Design Principles

### 1. **Modular Architecture**
Each component is self-contained and can be used independently. Components follow a consistent API pattern with fluent interfaces.

### 2. **Design Token System**
All styling uses centralized design tokens from `foundation.go`. This ensures consistency and makes theming straightforward.

### 3. **ANSI-First Approach**
Direct ANSI escape sequence usage for maximum compatibility and performance across terminal environments.

### 4. **Preset Configurations**
Common use cases have preset configurations to reduce boilerplate while maintaining customization options.

### 5. **Future-Ready**
Structured for easy extraction as a standalone library with minimal dependencies.

## Usage Examples

### Basic Table
```go
import "github.com/bthompso/engx-ergonomics-poc/internal/tui/components"

// Using preset
formatter := components.NewArchetypeTable()
fmt.Println(formatter.FormatHeader())

// Custom table
columns := []components.TableColumn{
    {Header: "Name", Width: 20, Color: design.ColorBrightWhite, Alignment: "left"},
    {Header: "Status", Width: 10, Color: design.ColorDarkGray, Alignment: "center"},
}
formatter := components.NewTableFormatter(columns)
```

### Styled Text
```go
// Using components
title := components.NewText("My Application").AsTitle()
fmt.Println(title.Render())

// Using convenience functions
fmt.Println(components.Success("✅ Operation completed"))
fmt.Println(components.Command("`engx create MyApp`"))
```

### Layout Components
```go
// Header with dividers
header := components.NewHeader("Configuration").WithSubtitle("Current settings")
fmt.Print(header.Render())

// Separator
separator := components.NewSeparator().WithColor(design.ColorBrightBlue)
fmt.Print(separator.Render())

// List
items := []string{"Item 1", "Item 2", "Item 3"}
list := components.NewList(items).AsNumbered()
fmt.Print(list.Render())
```

## Extraction Strategy

### Phase 1: Standalone Package
- Move `design/` and `components/` to standalone module
- Add comprehensive test suite
- Create example applications
- Document API thoroughly

### Phase 2: Enhanced Features
- Terminal width detection and responsive layouts
- Theme system with custom color schemes
- Animation support for progress indicators
- Advanced table features (sorting, filtering)

### Phase 3: Ecosystem
- Plugin system for custom components
- Integration helpers for popular CLI frameworks
- Performance optimizations for large datasets
- Cross-platform terminal compatibility testing

## Dependencies

### Current
- None (uses only Go standard library)
- Legacy components depend on `lipgloss` (to be removed)

### Future Extraction
- Minimal dependencies (possibly just Go standard library)
- Optional integrations with `bubbletea` for interactive components
- Plugin system for extended functionality

## Benefits Over Existing Solutions

### vs. lipgloss
- **Direct ANSI Control**: No abstraction layer, better performance
- **Design Token System**: Centralized design constants
- **Component Presets**: Ready-to-use configurations
- **Table System**: Advanced table formatting capabilities

### vs. bubble-tea
- **Complementary**: Focuses on styling and layout, not interaction
- **Lightweight**: No event handling overhead for static displays
- **Immediate Mode**: Direct rendering without state management

### vs. Custom Terminal Code
- **Consistency**: Standardized design tokens and patterns
- **Productivity**: Pre-built components and presets
- **Maintainability**: Modular architecture and clear API
- **Professional Polish**: Carefully designed visual hierarchy

## Current Integration

The design system is currently integrated into:
- **Template Listings** (`pkg/common/templates_ui.go`)
- **Analytics Displays** (`pkg/common/analytics_ui.go`)
- **Command Output** (various CLI commands)

This provides real-world testing and validation of the component API before extraction.