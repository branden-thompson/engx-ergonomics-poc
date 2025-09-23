# Requirements - EngX Root Command Experience

## Feature Classification
- **Type**: MAJOR FEATURE ENHANCEMENT
- **SEV Level**: SEV-0 (Critical)
- **Feature ID**: engx-root-cmd-experience

## Core Objective
Implement an interactive command selection interface when users invoke `engx` without arguments, following the principle of least surprise while maintaining all existing functionality.

## Primary Requirements

### 1. Interactive Interface Requirements
- **R1.1**: Interactive view with keyboard functionality (↑↓, j/k, 1-9, Enter)
- **R1.2**: Terminal width aware with 40-column minimum safe operation
- **R1.3**: No constant re-renders - smooth, stable interface
- **R1.4**: Preserve final selection state in terminal history before command execution
- **R1.5**: Safe fallback number-only selection mode for accessibility

### 2. Dynamic Command System Requirements
- **R2.1**: Not hard-coded - auto-discover available commands
- **R2.2**: Configurable command filtering and categorization
- **R2.3**: Integration with existing plugin system
- **R2.4**: Support for "Coming Soon" roadmap items with visual indicators

### 3. User Context Requirements
- **R3.1**: Pull @LDAP-username from system environment
- **R3.2**: Configurable CREW-ID assignments per command
- **R3.3**: Template-based interface layout with substitution variables

### 4. Compatibility Requirements
- **R4.1**: Preserve existing `engx --help` functionality
- **R4.2**: Maintain all existing command interfaces unchanged
- **R4.3**: Support existing global flags (--quiet, --verbose, etc.)
- **R4.4**: No breaking changes to current workflows

### 5. Responsive Design Requirements
- **R5.1**: Minimum width threshold: 40 columns
- **R5.2**: Graceful degradation for narrow terminals
- **R5.3**: Expandable description column for wider terminals
- **R5.4**: Different layouts for 40-60, 60-80, 80+ column widths

## Success Criteria
1. ✅ `engx` with no args shows interactive interface
2. ✅ `engx --help` shows traditional help unchanged
3. ✅ All existing commands work identically
4. ✅ Interface displays available vs coming soon commands
5. ✅ Keyboard navigation works reliably
6. ✅ Command execution preserves interface output
7. ✅ Terminal width responsiveness operates correctly
8. ✅ Configuration system allows roadmap expansion

## Non-Requirements
- Real implementation of "Coming Soon" commands
- Complex authentication or permission systems
- Multi-user collaboration features
- Network-based command discovery
- Mobile or web interface support

## Acceptance Test Cases
- **AC1**: User runs `engx` → sees interactive interface with current commands
- **AC2**: User selects command via number → command executes with interface preserved
- **AC3**: User navigates with arrow keys → selection updates smoothly
- **AC4**: Terminal resized → interface adjusts appropriately
- **AC5**: User runs `engx --help` → traditional help displayed unchanged
- **AC6**: User runs `engx create MyApp` → existing create workflow unchanged