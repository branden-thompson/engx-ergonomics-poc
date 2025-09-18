# Implementation Insights - User Configuration Prompts

## Key Technical Learnings

### 1. Bubble Tea Full-Screen vs Inline Mode
**Discovery**: `tea.WithAltScreen()` creates jarring user experience when transitioning from CLI prompts
**Solution**: Remove alt-screen mode and use inline rendering for seamless terminal history preservation
**Impact**: Natural build tool workflow that engineers expect

### 2. Configuration State Management
**Challenge**: Bridging CLI prompts with TUI execution phase
**Solution**: `NewAppModelWithConfig()` function to accept pre-configured state
**Pattern**: Two-phase workflow: prompts → execution with shared configuration

### 3. Prompt System Architecture
**Design**: Configurable JSON/YAML-driven prompt system with conditional triggers
**Benefits**:
- Extensible for future prompts
- Environment-specific prompt logic
- Clean separation of concerns

### 4. Progress Display Synchronization
**Issue**: Step count misalignment between tracker and renderer
**Root Cause**: Different step enumeration methods
**Fix**: Extract step names directly from tracker to guarantee alignment

### 5. Terminal I/O Handling
**Learning**: Bubble Tea requires explicit I/O configuration for non-interactive environments
**Solution**: `tea.WithInput(os.Stdin)` and `tea.WithOutput(os.Stderr)` for proper stream handling

## Best Practices Established

### CLI Workflow Design
1. **Traditional Prompts First** - Use familiar CLI interaction patterns
2. **Seamless Transition** - No jarring screen changes
3. **History Preservation** - Keep everything in terminal scrollback
4. **Auto-completion** - No manual exit required

### Code Organization
1. **Configuration Layer** - Clean separation of prompt config from logic
2. **Conditional Logic** - Trigger-based prompt system for flexibility
3. **State Transfer** - Explicit configuration passing between phases

### Error Handling
1. **Input Validation** - Robust prompt validation with helpful error messages
2. **Graceful Degradation** - Fallback to defaults when configuration fails
3. **User Feedback** - Clear response messages for user actions

## Performance Insights

### Memory Usage
- Inline rendering uses less memory than full-screen TUI
- Configuration objects are lightweight and fast to serialize

### User Experience
- Reduced cognitive load with familiar CLI patterns
- Faster workflow without screen transitions
- Better debugging with persistent terminal output

## Reusable Components

### 1. Inline Prompter (`internal/prompts/inline.go`)
- Reusable for any CLI prompt needs
- Configurable validation and response handling

### 2. Configuration System (`internal/config/prompts.go`)
- JSON/YAML driven prompt definitions
- Conditional trigger system
- Extensible for new prompt types

### 3. State Transfer Pattern
- `NewAppModelWithConfig()` pattern for passing state between phases
- Applicable to other multi-phase CLI tools

## Future Considerations

### Scalability
- Prompt system easily extends to support new question types
- Configuration format supports complex conditional logic
- TUI rendering can be enhanced with colors and animations

### Maintainability
- Clear separation between prompt logic and TUI rendering
- Configuration-driven approach reduces code changes for new prompts
- Documented patterns for extending functionality

### User Experience Enhancements
- Terminal history preservation enables better debugging
- Inline mode supports future color and styling improvements
- Natural workflow matches developer expectations

## Recent Enhancement Learnings (Session 2)

### 6. Visual Consistency and Color Systems
**Challenge**: Maintaining consistent visual styling across different UI components
**Solution**: Implemented comprehensive color system with state-based styling
**Impact**: Professional appearance with clear visual hierarchy

### 7. Progress Bar Architecture
**Initial Issue**: Multiple percentage formatting implementations leading to inconsistencies
**Root Cause**: Different rendering paths for total progress vs step progress
**Solution**: Modular progress bar system with configurable width and percentage display

#### Modular Progress Bar System
```go
type ProgressBarConfig struct {
    Width           int    // Fixed width or 0 for fill mode
    ShowPercentage  bool   // Toggle percentage display
    PercentagePad   int    // Padding for alignment
    FillMode        bool   // Dynamic width calculation
}
```

**Benefits**:
- Single source of truth for percentage formatting
- Consistent "100.0%" display across all progress indicators
- Configurable width for different use cases
- Eliminates layout calculation errors

### 8. Layout Calculation Edge Cases
**Discovery**: Dynamic width calculations can result in negative values
**Issue**: `slice bounds out of range [:-4]` panic when terminal width constraints create impossible layouts
**Solution**: Comprehensive bounds checking with minimum width enforcement
**Pattern**: Always validate calculated widths before slice operations

### 9. Component Status Visualization
**Enhancement**: Modular component installation status system
**Implementation**: Color-coded status indicators with consistent styling
- `[queued]` (White)
- `[installing...]` (Blue)
- `[skipped]` (Grey Italic)
- `[installed]` (Green)
- `[failed]` (Red)

### 10. Floating Point Precision in Percentage Display
**Subtle Bug**: `fmt.Sprintf("%.1f%%", 100.0)` can produce "100%" instead of "100.0%"
**Root Cause**: Floating point representation and format string behavior
**Robust Solution**: Post-processing string correction to ensure consistent format
```go
percentText = fmt.Sprintf("%.1f%%", percentage)
if percentText == "100%" {
    percentText = "100.0%"
}
```

## Enhanced Best Practices

### Modular Component Design
1. **Configuration Structs** - Use typed configuration for complex components
2. **Result Types** - Return structured data for layout calculations
3. **Bounds Validation** - Always validate calculated dimensions
4. **Color Consistency** - Centralized color system for state-based styling

### Progress Visualization
1. **Fixed vs Fill Modes** - Choose appropriate sizing strategy based on context
2. **Percentage Consistency** - Use single formatting function across all displays
3. **Layout Robustness** - Handle edge cases gracefully with minimums and maximums

### Error Prevention
1. **Slice Bounds** - Validate all calculated indices before slice operations
2. **Width Calculations** - Ensure positive values with minimum constraints
3. **String Formatting** - Post-process format results for consistency

## Performance and UX Impact

### Improved Visual Experience
- Consistent percentage formatting eliminates user confusion
- Proper progress bar sizing preserves readability of step names
- Color-coded status system provides immediate visual feedback

### Code Maintainability
- Modular progress bar system reduces duplication
- Centralized styling makes future enhancements easier
- Robust error handling prevents runtime panics

### Development Velocity
- Configuration-driven approach allows rapid iteration
- Reusable components accelerate feature development
- Comprehensive error handling reduces debugging time