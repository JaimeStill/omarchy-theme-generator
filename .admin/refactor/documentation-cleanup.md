# Documentation Infrastructure Cleanup

## Phase 0: Documentation Alignment (Execute First!)

This phase must be completed before any code changes to ensure development sessions have accurate references.

---

## 1. PROJECT.md Update

```markdown
# Omarchy Theme Generator - Project Status

## Current Implementation

### Infrastructure
- **Go module**: `github.com/JaimeStill/omarchy-theme-generator`
- **Go version**: 1.25.0
- **Binary name**: `omarchy-theme-gen`

### Packages Implemented

#### pkg/color (to be refactored → pkg/formats)
- Color type with RGBA storage and HSL conversion
- Contrast calculation for WCAG compliance
- Distance metrics (RGB, HSL, LAB)
- Hex color parsing and formatting
- **Status**: Over-engineered, needs simplification

#### pkg/extractor
- Multi-strategy extraction (frequency and saliency)
- Automatic strategy selection based on image analysis
- Image characteristic analysis (edge detection, complexity)
- Settings-driven configuration with empirical thresholds
- **Status**: Working but needs decomposition

#### pkg/errors
- Centralized error handling
- Domain-specific error types
- **Status**: Complete and adequate

#### tests/
- Strategy validation with 15 test images
- Image analysis utility
- Benchmark suite
- **Status**: Comprehensive coverage

### Capabilities
- ✅ Process 4K images in <2 seconds
- ✅ Memory usage <100MB
- ✅ Multi-strategy extraction
- ✅ Grayscale vs monochromatic detection
- ✅ Settings-driven configuration

---

## Current Work

### Documentation Cleanup (Active)
- [ ] Update PROJECT.md structure
- [ ] Align README.md with new architecture
- [ ] Update technical documentation
- [ ] Clean up CLAUDE.md
- [ ] Update development prompts

### Architecture Refactoring (Planned)
- [ ] pkg/color → pkg/formats (simplification)
- [ ] Extract pkg/analysis from extractor
- [ ] Extract pkg/strategies from extractor
- [ ] Create pkg/settings and pkg/config
- [ ] Implement purpose-driven extraction

---

## Architecture Overview

### Dependency Layers

```
Layer 1: Foundation
├── pkg/formats         # Color conversion, types (no dependencies)
├── pkg/settings        # System configuration
└── pkg/config          # User preferences

Layer 2: Analysis
└── pkg/analysis        # Image analysis (depends on: formats)

Layer 3: Processing
├── pkg/strategies      # Extraction strategies (depends on: formats, analysis)
└── pkg/extractor       # Orchestration (depends on: formats, analysis, strategies)

Layer 4: Generation
├── pkg/palette         # Color schemes (depends on: formats, analysis)
└── pkg/theme           # Theme files (depends on: formats, palette)

Layer 5: Application
└── cmd/omarchy-theme-gen  # CLI (depends on: all packages)
```

---

## Components & Features

### pkg/formats (Refactor of pkg/color)
**Purpose**: Minimal color conversion and formatting utilities

**Required Features**:
- `RGBToHSL()` - HSL conversion for analysis
- `ContrastRatio()` - WCAG compliance checking
- `ToHex()`, `ToHexA()` - Color formatting
- `ParseHex()` - Hex string parsing
- Type definitions (ColorRole, ThemeMode, etc.)

**Status**: Needs refactoring (remove 90% of current code)

### pkg/analysis (Extract from extractor)
**Purpose**: Image and color analysis

**Required Features**:
- Image characteristic analysis
- Color profile detection (grayscale, monochromatic, etc.)
- Theme mode detection (light/dark)
- Role assignment logic

**Status**: Not yet extracted

### pkg/strategies (Extract from extractor)
**Purpose**: Pluggable extraction strategies

**Required Features**:
- Strategy interface
- Frequency strategy
- Saliency strategy
- Strategy selector

**Status**: Not yet extracted

### pkg/extractor (Simplify)
**Purpose**: Extraction orchestration only

**Required Features**:
- Pipeline coordination
- Result aggregation

**Status**: Needs simplification

### pkg/palette (New)
**Purpose**: Color scheme generation

**Required Features**:
- Color theory schemes (complementary, triadic, etc.)
- Edge case synthesis
- WCAG validation

**Status**: Not implemented

### pkg/theme (New)
**Purpose**: Theme file generation

**Required Features**:
- Template processing
- Format-specific generation
- Metadata creation

**Status**: Not implemented

### pkg/settings (New)
**Purpose**: System configuration

**Required Features**:
- Settings structure
- Default values
- JSON loading

**Status**: Not implemented

### pkg/config (New)
**Purpose**: User preferences

**Required Features**:
- Preference structure
- Override system
- Per-theme storage

**Status**: Not implemented

---

## Next Steps

1. **Complete documentation cleanup** (Current)
2. **Refactor pkg/color → pkg/formats**
3. **Decompose extractor package**
4. **Implement purpose-driven extraction**
5. **Add palette generation**
6. **Create theme generators**
7. **Build CLI interface**

---

## Testing Strategy

- Unit tests for each package
- Integration tests for pipeline
- Benchmark tests for performance
- Real image validation

---

## Performance Metrics

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| 4K Processing | <2s | 241ms | ✅ |
| Memory Usage | <100MB | 72MB | ✅ |
| Extraction Strategies | 2+ | 2 | ✅ |
| Color Schemes | 7 | 0 | ⏳ |
| Output Formats | 9 | 0 | ⏳ |
```

---

## 2. README.md Update

```markdown
# Omarchy Theme Generator

Generate beautiful, accessible terminal themes from any image using intelligent color extraction and color theory principles.

## Features

- 🎨 **Intelligent Color Extraction** - Multi-strategy system that adapts to image characteristics
- 🎯 **Purpose-Driven Colors** - Organizes colors by their role (background, foreground, accents)
- ♿ **WCAG Compliant** - Ensures text readability with proper contrast ratios
- 🎭 **Light/Dark Mode** - Automatic detection with manual override options
- 🖼️ **Edge Case Handling** - Gracefully handles grayscale, duotone, and monochromatic images
- ⚡ **High Performance** - Processes 4K images in under 2 seconds

## Installation

```bash
go install github.com/JaimeStill/omarchy-theme-generator/cmd/omarchy-theme-gen@latest
```

Or build from source:

```bash
git clone https://github.com/JaimeStill/omarchy-theme-generator
cd omarchy-theme-generator
go build -o omarchy-theme-gen cmd/omarchy-theme-gen/main.go
```

## Usage

### Generate Theme

```bash
omarchy-theme-gen generate --image wallpaper.jpg
```

### Customize Generation

```bash
# Force dark mode
omarchy-theme-gen generate --image wallpaper.jpg --mode dark

# Override primary color
omarchy-theme-gen generate --image wallpaper.jpg --primary "#e94560"

# Apply color scheme
omarchy-theme-gen generate --image wallpaper.jpg --scheme complementary
```

### Refine Existing Theme

```bash
# Change color scheme
omarchy-theme-gen set-scheme my-theme --scheme triadic

# Switch mode
omarchy-theme-gen set-mode my-theme --mode light

# Clone and modify
omarchy-theme-gen clone my-theme my-variant --primary "#3498db"
```

## Architecture

The generator uses a layered architecture with clear dependencies:

### Core Packages

- **pkg/formats** - Color conversion and formatting utilities
- **pkg/analysis** - Image characteristic and profile detection
- **pkg/extractor** - Color extraction orchestration
- **pkg/strategies** - Pluggable extraction strategies (frequency, saliency)
- **pkg/palette** - Color theory and scheme generation
- **pkg/theme** - Theme file generation from templates
- **pkg/settings** - System configuration and thresholds
- **pkg/config** - User preferences and overrides

### Extraction Pipeline

1. **Analysis** - Detect image characteristics and color profile
2. **Strategy Selection** - Choose optimal extraction strategy
3. **Extraction** - Extract colors using selected strategy
4. **Role Assignment** - Categorize colors by purpose
5. **Synthesis** - Generate missing colors if needed
6. **Validation** - Ensure WCAG compliance
7. **Generation** - Create theme configuration files

## Configuration

### System Settings

Global extraction behavior in `config/settings.json`:

```json
{
  "extraction": {
    "saliency_edge_threshold": 0.036,
    "frequency_sample_rate": 4
  },
  "accessibility": {
    "min_contrast_ratio": 4.5
  }
}
```

### User Preferences

Per-theme overrides in `theme-name/theme-preferences.json`:

```json
{
  "mode": "dark",
  "color_overrides": {
    "background": "#1a1a2e",
    "primary": "#e94560"
  }
}
```

## Development

```bash
# Run tests
go test ./...

# Run benchmarks
go test -bench=. ./tests

# Analyze test images
go run tests/analyze-images/main.go

# Validate code
go vet ./...
```

## Requirements

- Go 1.25 or later
- No external dependencies (pure Go implementation)

## License

MIT

## Acknowledgments

Built for the [Omarchy](https://omarchy.org) desktop environment.
```

---

## 3. docs/ Directory Updates

### docs/technical-specification.md

```markdown
# Technical Specification

## Overview

Omarchy Theme Generator extracts colors from images and generates accessible, aesthetically pleasing themes using color theory principles.

## Architecture

### Layered Design

The system uses a strict layered architecture where each layer depends only on layers below it:

1. **Foundation Layer** - Basic types and utilities
2. **Analysis Layer** - Image and color analysis
3. **Processing Layer** - Color extraction
4. **Generation Layer** - Theme creation
5. **Application Layer** - User interface

### Core Concepts

#### Purpose-Driven Extraction

Instead of extracting colors by frequency, the system categorizes colors by their intended role:

- **Backgrounds** - Colors suitable for window/terminal backgrounds
- **Foregrounds** - Colors suitable for text
- **Accents** - Saturated colors for highlights
- **Terminal Colors** - ANSI color palette mapping

#### Mode-Aware Processing

Role assignment adapts based on detected theme mode:

```go
// Dark mode: darker colors become backgrounds
if mode == ModeDark && lightness < 0.35 {
    role = RoleBackground
}

// Light mode: lighter colors become backgrounds
if mode == ModeLight && lightness > 0.75 {
    role = RoleBackground
}
```

#### Settings vs Configuration

- **Settings** (`config/settings.json`) - HOW the tool operates
- **Configuration** (`theme/preferences.json`) - WHAT the user wants

### Performance Requirements

| Metric | Target | Notes |
|--------|--------|-------|
| 4K Image Processing | <2 seconds | Full pipeline |
| Memory Usage | <100MB peak | During extraction |
| Minimum Contrast | 4.5:1 | WCAG AA compliance |

### Color Space Operations

Uses standard Go `image/color` types with minimal extensions:

- RGB to HSL conversion for saturation analysis
- Contrast ratio calculation for accessibility
- Hex formatting for configuration output

## Implementation Details

[Rest of current technical content...]
```

### docs/architecture.md (NEW)

```markdown
# Architecture Documentation

## System Architecture

```
┌─────────────────────────────────────┐
│      Application Layer (CLI)         │
├─────────────────────────────────────┤
│      Generation Layer                │
│   ┌──────────┐    ┌──────────┐     │
│   │ palette  │    │  theme   │     │
│   └──────────┘    └──────────┘     │
├─────────────────────────────────────┤
│      Processing Layer                │
│   ┌──────────┐    ┌──────────┐     │
│   │extractor │    │strategies│     │
│   └──────────┘    └──────────┘     │
├─────────────────────────────────────┤
│      Analysis Layer                  │
│   ┌──────────────────────────┐      │
│   │       analysis           │      │
│   └──────────────────────────┘      │
├─────────────────────────────────────┤
│      Foundation Layer                │
│   ┌─────────┐ ┌────────┐ ┌────────┐│
│   │ formats │ │settings│ │ config ││
│   └─────────┘ └────────┘ └────────┘│
└─────────────────────────────────────┘
```

## Package Responsibilities

### Foundation Layer

**pkg/formats**
- Color space conversions
- Format utilities
- Type definitions
- No dependencies

**pkg/settings**
- System configuration
- Default values
- Threshold management

**pkg/config**
- User preferences
- Override system
- Theme-specific settings

### Analysis Layer

**pkg/analysis**
- Image characteristic detection
- Color profile classification
- Mode detection
- Role assignment logic

### Processing Layer

**pkg/extractor**
- Orchestrates extraction pipeline
- Coordinates strategies
- Aggregates results

**pkg/strategies**
- Pluggable extraction algorithms
- Strategy selection
- Frequency and saliency strategies

### Generation Layer

**pkg/palette**
- Color theory schemes
- Edge case synthesis
- Color harmony validation

**pkg/theme**
- Template processing
- Configuration generation
- Format-specific output

### Application Layer

**cmd/omarchy-theme-gen**
- CLI interface
- Command handling
- User interaction

## Data Flow

1. **Input** → Image file
2. **Analysis** → Characteristics & profile
3. **Strategy Selection** → Optimal algorithm
4. **Extraction** → Raw colors
5. **Role Assignment** → Categorized colors
6. **Synthesis** → Missing colors generated
7. **Validation** → Accessibility checks
8. **Generation** → Theme files
9. **Output** → Configuration files

## Design Principles

1. **Separation of Concerns** - Each package has a single responsibility
2. **Dependency Injection** - Packages depend on interfaces, not implementations
3. **Settings-Driven** - No hardcoded values
4. **Standard Types** - Use Go standard library where possible
5. **Purpose-Driven** - Organize by intent, not implementation
```

---

## 4. CLAUDE.md Update

```markdown
# Omarchy Theme Generator

## Project Context
Go-based CLI tool that generates Omarchy themes from images using intelligent color extraction and color theory principles.

## Current Architecture

### Layered Design
- **Foundation**: formats, settings, config (no dependencies)
- **Analysis**: Image and color analysis (depends on foundation)
- **Processing**: Extraction and strategies (depends on analysis)
- **Generation**: Palette and theme creation (depends on processing)
- **Application**: CLI interface (depends on all)

### Key Technical Decisions
- **Standard Types**: Use `color.RGBA` from standard library, not custom types
- **Purpose-Driven**: Colors organized by role (background, foreground, accent)
- **Settings vs Config**: System settings (HOW) separate from user config (WHAT)
- **Multi-Strategy**: Frequency for simple images, saliency for complex
- **Edge Case Handling**: Synthesis for grayscale, duotone, monochromatic

## Current Implementation Status
- ✅ Multi-strategy extraction (frequency/saliency)
- ✅ Image analysis and characterization
- ✅ Settings-driven configuration
- ✅ Performance targets met (<2s for 4K)
- 🔄 Architecture refactoring in progress
- ⏳ Purpose-driven extraction pending
- ⏳ Color scheme generation pending
- ⏳ Theme file generation pending

## Development Rules
1. Use standard library types where possible
2. No hardcoded values - use settings
3. Clear dependency layers - no circular deps
4. Test with real images immediately
5. Document with technical precision

## Package Structure
```
pkg/
├── formats/      # Color utilities (refactor of pkg/color)
├── analysis/     # Image analysis (extract from extractor)
├── strategies/   # Extraction strategies (extract from extractor)
├── extractor/    # Orchestration (simplify)
├── palette/      # Color schemes (new)
├── theme/        # Theme generation (new)
├── settings/     # System config (new)
└── config/       # User preferences (new)
```

## Commands
```bash
# Validate
go vet ./...

# Test
go test ./tests -v

# Benchmark
go test -bench=. ./tests

# Format
go fmt ./...
```

## Performance Targets
- 4K image: <2 seconds ✅
- Memory: <100MB ✅
- Contrast: WCAG AA (4.5:1) ⏳

## Next Development Focus
1. Clean up documentation (current)
2. Refactor pkg/color → pkg/formats
3. Decompose extractor package
4. Implement purpose-driven extraction
5. Add color scheme generation

## References
- Technical details: `docs/technical-specification.md`
- Architecture: `docs/architecture.md`
- Progress: `PROJECT.md`
```

---

## 5. prompts/ Directory Updates

### prompts/session-initiate.md (Update)

Add section about architecture:

```markdown
## Architecture Context

The project uses a layered architecture with clear dependencies:

### Package Layers
1. **Foundation** (formats, settings, config) - No dependencies
2. **Analysis** (analysis) - Depends on foundation
3. **Processing** (strategies, extractor) - Depends on analysis
4. **Generation** (palette, theme) - Depends on processing
5. **Application** (cmd) - Depends on all

### Key Principles
- Use standard library types (`color.RGBA` not custom types)
- Settings-driven (no hardcoded values)
- Purpose-driven extraction (role-based, not frequency)
- Clear separation of concerns
```

### prompts/architecture-guide.md (NEW)

```markdown
# Architecture Guide for Development

## Package Dependencies

When implementing features, respect the dependency hierarchy:

```
formats, settings, config
    ↑
analysis
    ↑
strategies, extractor
    ↑
palette, theme
    ↑
cmd
```

## Implementation Guidelines

### Adding New Features

1. **Identify the correct layer** for your feature
2. **Only import from lower layers** (or same layer)
3. **Use interfaces** for cross-package communication
4. **Add settings** for any thresholds or parameters
5. **Write tests** at the package level

### Refactoring Existing Code

1. **Check current dependencies** before moving code
2. **Update imports** in all affected files
3. **Run tests** after each change
4. **Update documentation** to reflect new structure

## Common Patterns

### Settings Pattern
```go
// In pkg/settings/settings.go
type FeatureSettings struct {
    Threshold float64 `json:"threshold"`
    Enabled   bool    `json:"enabled"`
}

// Usage
settings := settings.Load()
if settings.Feature.Enabled {
    // Apply feature with threshold
}
```

### Strategy Pattern
```go
// In pkg/strategies/interface.go
type Strategy interface {
    Process(input Input) (Output, error)
    CanHandle(input Input) bool
}

// Implementation
type MyStrategy struct{}
func (s *MyStrategy) Process(input Input) (Output, error) { ... }
func (s *MyStrategy) CanHandle(input Input) bool { ... }
```

### Role-Based Organization
```go
// Instead of frequency-based
topColors := ExtractTopColors(img, 10)

// Use role-based
colorMap := ExtractColorsByRole(img)
background := colorMap[RoleBackground]
foreground := colorMap[RoleForeground]
```
```

---

## Execution Order

1. **Update PROJECT.md** with new structure ✅
2. **Update README.md** to reflect architecture ✅
3. **Create/Update docs/** ✅
   - Update technical-specification.md
   - Create architecture.md
4. **Update CLAUDE.md** with current state ✅
5. **Update prompts/** ✅
   - Update session-initiate.md
   - Create architecture-guide.md

## Notes

- This cleanup ensures all documentation reflects the new architecture
- Future development sessions will have accurate references
- No code changes yet - documentation first!
- After this cleanup, we can proceed with the actual refactoring
