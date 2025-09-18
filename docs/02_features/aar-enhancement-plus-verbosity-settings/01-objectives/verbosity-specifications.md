# Verbosity Level Specifications

## Updated Verbosity Hierarchy

```
DEBUG (--debug)     │ Maximum verbosity, all system outputs, technical details
  ↓                 │
VERBOSE (--verbose) │ Enhanced details, progress bars for multi-step processes
  ↓                 │
DEFAULT (default)   │ Current view (chosen when no verbosity option specified)
  ↓                 │
CONCISE (--concise) │ Less detail, components info and granular items hidden
  ↓                 │
QUIET (--quiet)     │ Essential info only: total progress, current step, footer
```

## Detailed Level Specifications

### 1. DEBUG Mode (--debug)
**Purpose**: Maximum verbosity for troubleshooting and development
**Target Users**: Developers debugging issues, contributors, technical support

**Display Elements**:
- ✅ All system outputs displayed
- ✅ Technical implementation details
- ✅ Timing information with microsecond precision
- ✅ Internal state changes and transitions
- ✅ Memory usage and performance metrics
- ✅ Stack traces for any errors
- ✅ Component lifecycle information
- ✅ Network requests and responses
- ✅ File system operations
- ✅ Configuration loading and parsing

**Output Example**:
```
[DEBUG 00:00:01.234567] StateManager: Transitioning StateExecuting → StateComplete
[DEBUG 00:00:01.234789] ProgressTracker: Step 12/12 completed in 45.123ms
[DEBUG 00:00:01.234890] Memory: Heap usage 15.2MB (+0.3MB), GC triggered
[DEBUG 00:00:01.235001] Renderer: Calculating layout for 89 components
[DEBUG 00:00:01.235123] TUI: Buffer flush scheduled, 4KB pending

🎯 PROJECT SETUP PROGRESS
├─ 🏗️  Project Structure           ✅ (234ms) [DEBUG: 47 files created]
├─ 📦 Package Installation         ✅ (45.1s) [DEBUG: npm install, 234 packages]
├─ ⚙️  TypeScript Configuration    ✅ (156ms) [DEBUG: tsconfig.json, 12 rules]
└─ ✅ Setup Complete              100.0% [DEBUG: Total memory: 15.2MB]
```

### 2. VERBOSE Mode (--verbose)
**Purpose**: Enhanced details for power users and detailed progress tracking
**Target Users**: Experienced developers, CI/CD systems, detailed monitoring

**Display Elements**:
- ✅ Enhanced progress bars for all multi-step processes
- ✅ Non-technical system outputs displayed
- ✅ Individual step timing and status
- ✅ Feature installation details
- ✅ Configuration choices and their implications
- ✅ Performance metrics (human-readable)
- ✅ Detailed component status
- ✅ Helpful context and explanations

**Output Example**:
```
🎯 PROJECT SETUP PROGRESS

📁 PROJECT STRUCTURE
├─ Creating directory structure     ✅ (234ms)
├─ Initializing package.json       ✅ (156ms)
├─ Setting up TypeScript config     ✅ (189ms)
└─ Creating source directories      ✅ (123ms)

📦 PACKAGE INSTALLATION
Installing dependencies [████████████████████████████████] 100% (45.1s)
├─ React & TypeScript              ✅ (12.3s) - Core framework
├─ Development tools               ✅ (15.2s) - ESLint, Prettier, etc.
├─ Testing framework               ✅ (8.9s)  - Jest, Testing Library
└─ Build tools                     ✅ (8.7s)  - Webpack, PostCSS

⚙️ APPLICATION COMPONENTS
├─ Hot Reload Setup                ✅ (1.4s)  - Fast development refresh
├─ ESLint Configuration            ✅ (0.8s)  - Code quality rules
├─ Testing Environment             ✅ (2.1s)  - Jest + React Testing Library
└─ VS Code Settings                ✅ (0.2s)  - Workspace configuration

✅ Setup Complete                  100.0% (2m 34s)
```

### 3. DEFAULT Mode (no flag)
**Purpose**: Current balanced view - unchanged from existing behavior
**Target Users**: General developers, standard usage

**Display Elements**:
- ✅ Current progress view maintained exactly
- ✅ Standard step progression
- ✅ Current timing display
- ✅ Current component status
- ✅ Current footer and progress bar

**Output Example**:
```
🎯 PROJECT SETUP PROGRESS

├─ 🏗️  Project Structure           ✅
├─ 📦 Package Installation         ✅
├─ ⚙️  TypeScript Configuration    ✅
├─ 🛠️  Development Features        ✅
├─ 🧪 Testing Setup               ✅
└─ ✅ Setup Complete              100.0%

📦 APPLICATION COMPONENTS
├─ [✓] Hot Reload                 ✅
├─ [✓] ESLint + Prettier          ✅
├─ [✓] React DevTools             ✅
└─ [✓] VS Code Configuration      ✅

✨ Running Time: 2m 34s
```

### 4. CONCISE Mode (--concise)
**Purpose**: Reduced detail, focus on essential progress only
**Target Users**: Quick development iterations, automated scripts, CI/CD (when some output needed)

**Display Elements**:
- ✅ Main progress sections only
- ❌ Component info hidden
- ❌ Granular line items hidden
- ✅ Overall progress percentage
- ✅ Current major step
- ✅ Essential timing information
- ✅ Footer information

**Output Example**:
```
🎯 PROJECT SETUP PROGRESS

├─ 🏗️  Project Structure           ✅
├─ 📦 Package Installation         ✅
├─ ⚙️  Configuration               ✅
└─ ✅ Setup Complete              100.0%

✨ Running Time: 2m 34s
```

### 5. QUIET Mode (--quiet)
**Purpose**: Minimal essential information only
**Target Users**: Automated systems, scripts, minimal output requirements

**Display Elements**:
- ✅ Total progress percentage
- ✅ Current major step name only
- ✅ Basic footer info (timing)
- ❌ Individual component status hidden
- ❌ Step details hidden
- ❌ Progress bars simplified to single line

**Output Example**:
```
🎯 Setup Complete 100.0% (2m 34s)
```

## Implementation Requirements

### Flag Mapping
- `--debug` → VerbosityDebug
- `--verbose` → VerbosityVerbose
- (no flag) → VerbosityDefault
- `--concise` → VerbosityConcise
- `--quiet` → VerbosityQuiet

### Precedence Rules
1. If multiple flags specified, highest verbosity wins
2. `--debug` overrides all others
3. `--quiet` is lowest precedence (overridden by any other)
4. No flag specified defaults to DEFAULT mode

### Output Element Control Matrix

| Element | Debug | Verbose | Default | Concise | Quiet |
|---------|-------|---------|---------|---------|-------|
| Technical details | ✅ | ❌ | ❌ | ❌ | ❌ |
| Component progress bars | ✅ | ✅ | ✅ | ❌ | ❌ |
| Individual step timing | ✅ | ✅ | ✅ | ❌ | ❌ |
| Component status list | ✅ | ✅ | ✅ | ❌ | ❌ |
| Main progress sections | ✅ | ✅ | ✅ | ✅ | ❌ |
| Total progress % | ✅ | ✅ | ✅ | ✅ | ✅ |
| Current step name | ✅ | ✅ | ✅ | ✅ | ✅ |
| Footer timing | ✅ | ✅ | ✅ | ✅ | ✅ |

### AAR Verbosity Scaling
- **DEBUG**: Full technical report with diagnostics
- **VERBOSE**: Comprehensive summary with detailed metrics
- **DEFAULT**: Standard AAR with next steps (current planned behavior)
- **CONCISE**: Brief summary with essential next steps
- **QUIET**: One-line completion status with primary next action