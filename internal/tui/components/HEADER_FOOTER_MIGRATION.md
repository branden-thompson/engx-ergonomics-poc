# Header/Footer Component Migration Guide

## The Problem This Solves

Previously, every plugin was manually building header/footer lines with inconsistent colors and complex width calculations.

## Before (Manual Implementation)
```go
// Manual dash building - no colors, error-prone width calculation
func buildHeader(title, info string, width int) string {
    titleLen := len(title)
    infoLen := len(info)
    dashLen := width - titleLen - infoLen - 8 // 8 for spaces and lead/tail
    if dashLen < 4 {
        dashLen = 4
    }
    return fmt.Sprintf("---- %s %s %s ----\n",
        title,
        strings.Repeat("-", dashLen),
        info)
}

// Usage everywhere:
header := buildHeader("Web Platform Team", "CREW ID: CREW-1234", terminalWidth)
footer := buildHeader("STANDARD CREW", "Last Updated: 09/18/2025", terminalWidth)
```

**Problems:**
- ❌ No colors applied
- ❌ Width calculation bugs common
- ❌ Inconsistent formatting across plugins
- ❌ Duplicate code everywhere
- ❌ Hard to maintain

## After (Modular Component with Default Colors)
```go
// Import once
import "github.com/bthompso/engx-ergonomics-poc/internal/tui/components"

// Create once
headerFooter := components.NewHeaderFooterComponent(terminalWidth)

// Use everywhere with automatic colors
header := headerFooter.RenderHeader("Web Platform Team", "CREW ID: CREW-1234")
footer := headerFooter.RenderFooter("STANDARD CREW", "Last Updated: 09/18/2025")
separator := headerFooter.RenderSeparator()
```

**Benefits:**
- ✅ **Grey dashes, white labels automatically**
- ✅ **Perfect width calculation** (matches proven header/footer precision)
- ✅ **Consistent styling** across all plugins
- ✅ **One line of code** for complex headers
- ✅ **Easy to maintain** - change colors in one place

## Custom Colors When Needed
```go
// Override colors for special cases (like EngX CLI branding)
header := headerFooter.RenderHeaderWithColors(
    "EngX CLI", "v1.0",
    design.ColorEngxPink,    // Custom left label color
    design.ColorBrightGreen  // Custom right label color
)
// Dashes are still grey by default
```

## Migration Checklist

For each renderer file with manual header/footer building:

1. **Replace imports:**
   ```go
   // Add this import
   import "github.com/bthompso/engx-ergonomics-poc/internal/tui/components"
   ```

2. **Replace manual dash building:**
   ```go
   // OLD: Manual building
   header := fmt.Sprintf("---- %s %s %s ----", ...)

   // NEW: Component usage
   headerFooter := components.NewHeaderFooterComponent(width)
   header := headerFooter.RenderHeader(leftLabel, rightLabel)
   ```

3. **Remove manual width calculations:**
   - Delete `strings.Repeat("-", dashCount)` code
   - Delete manual spacing calculations
   - Let component handle everything

4. **Test output** to ensure colors appear correctly

## Files Ready for Migration

These files currently do manual header/footer building and could benefit:

- `plugins/crews/renderers/asset_owner.go`
- `plugins/crews/renderers/membership_list.go`
- `plugins/crews/renderers/search_results.go`
- `internal/aar/formatters.go`
- `pkg/common/analytics_ui.go`
- `pkg/common/marketplace_ui.go`

## Result

✅ **Consistent visual styling** across entire application
✅ **Reduced code duplication** and maintenance burden
✅ **Automatic color application** - no more forgetting colors
✅ **Proven width calculation** - no more alignment bugs