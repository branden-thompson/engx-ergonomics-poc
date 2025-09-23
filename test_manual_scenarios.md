# Manual Testing Scenarios - EngX Root Command Experience

## MANUAL TESTING PROTOCOL
These tests must be run in a real terminal environment to validate interactive functionality.

## 🔧 Prerequisites
1. Build the binary: `go build -o engx ./cmd/engx`
2. Ensure you're in a real terminal (not CI/automated environment)
3. Terminal should be at least 40 columns wide

---

## 🧪 TEST SCENARIO 1: Basic Interactive Interface
**Objective**: Verify interactive interface launches and displays correctly

### Steps:
1. Run: `./engx`
2. Verify interface appears with:
   - Header: "---- EngX CLI v.v0.8.0 ------------------------------ OCT 2025 Release ----"
   - Command categories displayed
   - Available commands marked with ✓
   - Coming soon commands marked with ⏳
   - Footer with keyboard shortcuts

### Expected Results:
- ✅ Interface launches without errors
- ✅ Commands are numbered 1-15
- ✅ Categories include "Create & Manage Applications", "Analytics & Development Tools", etc.
- ✅ @{your-username} appears in "Specific to you" section
- ✅ Available commands: create, templates, analytics, dev
- ✅ Coming soon commands: promote, deploy, my apps, etc.

---

## 🧪 TEST SCENARIO 2: Keyboard Navigation
**Objective**: Verify all keyboard navigation patterns work

### Steps:
1. Launch: `./engx`
2. Test arrow key navigation:
   - Press ↓ multiple times → selection should move to next selectable command
   - Press ↑ multiple times → selection should move to previous selectable command
3. Test vi-style navigation:
   - Press `j` → should move down
   - Press `k` → should move up
4. Test number selection:
   - Press `1` → should select command #1 (create)
   - Press `2` → should select command #2 (templates)
5. Test help toggle:
   - Press `h` → detailed help should appear
   - Press `h` again → help should disappear

### Expected Results:
- ✅ Arrow keys navigate only to selectable commands (skip coming soon)
- ✅ Vi keys (j/k) work identically to arrow keys
- ✅ Number keys 1-9 work for quick selection
- ✅ Selected command is visually highlighted
- ✅ Help toggle works correctly

---

## 🧪 TEST SCENARIO 3: Command Execution
**Objective**: Verify command execution preserves interface and works correctly

### Steps:
1. Launch: `./engx`
2. Navigate to command #1 (create)
3. Press Enter
4. Verify the interface output remains visible
5. Verify `engx create` help screen appears

### Expected Results:
- ✅ Interface table remains visible before command execution
- ✅ Smooth transition to `engx create` command
- ✅ No errors or interface corruption

---

## 🧪 TEST SCENARIO 4: Terminal Width Responsiveness
**Objective**: Verify responsive design at different terminal widths

### Steps:
1. Set terminal to exactly 40 columns wide
2. Run: `./engx`
3. Verify interface fits and is readable
4. Expand terminal to 80+ columns
5. Run: `./engx` again
6. Verify descriptions expand to use available space

### Expected Results:
- ✅ 40-column mode: Compact but readable
- ✅ 80+ column mode: Full descriptions visible
- ✅ No text overflow or truncation issues
- ✅ Proper text wrapping for long descriptions

---

## 🧪 TEST SCENARIO 5: Error Handling
**Objective**: Verify graceful error handling and fallback behavior

### Steps:
1. Test invalid number selection:
   - Launch: `./engx`
   - Press `9` (should be invalid if < 9 commands available)
   - Verify error message appears
2. Test coming soon command selection:
   - Navigate to a "⏳ Coming Soon" command
   - Press Enter
   - Verify appropriate error message

### Expected Results:
- ✅ Invalid selections show clear error messages
- ✅ Coming soon commands display "not yet available" message
- ✅ Interface remains stable after errors

---

## 🧪 TEST SCENARIO 6: Backward Compatibility Verification
**Objective**: Ensure existing workflows are unchanged

### Steps:
1. Test help preservation: `./engx --help`
2. Test create command: `./engx create --help`
3. Test templates command: `./engx templates --help`
4. Test analytics command: `./engx analytics --help`
5. Test version flag: `./engx --version`
6. Test global flags: `./engx create MyApp --verbose`

### Expected Results:
- ✅ All help screens work exactly as before
- ✅ Version information displays correctly
- ✅ Global verbosity flags work with existing commands
- ✅ No breaking changes to any existing functionality

---

## 🧪 TEST SCENARIO 7: Configuration System
**Objective**: Verify configuration loading and user context

### Steps:
1. Check current user context:
   - Launch: `./engx`
   - Verify your username appears in "Specific to you (@{username})" section
2. Verify crew assignments:
   - Check that commands show "CREW-1234" in the crew column
3. Test configuration fallback:
   - Temporarily rename `.engx/roadmap.yaml` to `.engx/roadmap.yaml.backup`
   - Run: `./engx`
   - Verify graceful fallback to help screen
   - Restore the config file

### Expected Results:
- ✅ Your system username appears correctly
- ✅ Crew assignments display properly
- ✅ Missing config gracefully falls back to help

---

## 🧪 TEST SCENARIO 8: Exit and Quit Behavior
**Objective**: Verify clean exit behavior

### Steps:
1. Launch: `./engx`
2. Press `q` → should quit cleanly
3. Launch: `./engx`
4. Press `Ctrl+C` → should quit cleanly
5. Verify terminal state is restored properly

### Expected Results:
- ✅ `q` key exits immediately
- ✅ `Ctrl+C` exits immediately
- ✅ Terminal cursor and screen state restored
- ✅ No hanging processes or artifacts

---

## 📋 MANUAL TEST CHECKLIST

Copy this checklist and mark off each item as you test:

```
BASIC FUNCTIONALITY:
[ ] Interactive interface launches without errors
[ ] Header displays with correct version
[ ] All command categories appear
[ ] Available vs coming soon status correct
[ ] Username appears in personalized section

KEYBOARD NAVIGATION:
[ ] Arrow keys navigate correctly
[ ] Vi-style j/k keys work
[ ] Number selection (1-9) works
[ ] Help toggle (h) works
[ ] Selection highlighting visible

COMMAND EXECUTION:
[ ] Enter key executes selected command
[ ] Interface output preserved
[ ] Smooth transition to command help

RESPONSIVE DESIGN:
[ ] Works at 40-column width
[ ] Expands properly at 80+ columns
[ ] Text wrapping handles long descriptions
[ ] No overflow or layout issues

ERROR HANDLING:
[ ] Invalid number selections handled gracefully
[ ] Coming soon commands show appropriate message
[ ] Interface remains stable after errors

BACKWARD COMPATIBILITY:
[ ] `./engx --help` works unchanged
[ ] `./engx create --help` works unchanged
[ ] `./engx templates --help` works unchanged
[ ] `./engx --version` works unchanged
[ ] Global flags preserved

CONFIGURATION:
[ ] Username resolved from system
[ ] Crew assignments display correctly
[ ] Missing config falls back gracefully

EXIT BEHAVIOR:
[ ] `q` key exits cleanly
[ ] `Ctrl+C` exits cleanly
[ ] Terminal state restored properly
```

---

## 🚨 KNOWN LIMITATIONS

1. **TTY Requirement**: Interactive mode only works in real terminals, not in CI/automated environments
2. **Terminal Width**: Minimum 40 columns required for proper display
3. **Configuration Dependency**: Requires `.engx/roadmap.yaml` for full functionality

These are expected behaviors and not bugs. The system gracefully handles these limitations with appropriate fallbacks.