# Omarchy Theme Generator

## Project Context
Go-based CLI tool that generates Omarchy themes from images using color extraction and palette generation based on color theory principles. Optional TUI interface planned for future enhancement.

## Development Philosophy
- **User-driven**: AI provides implementation guides, user develops code
- **Explanatory mode**: Provide educational insights while coding
- **Test-driven**: Create comprehensive unit tests in tests/ subdirectories
- **Precise language**: Use correct technical terminology always
- **Reference, don't repeat**: Link to existing code and docs

## Core Documents
- **Architecture design**: `docs/architecture.md` - Layered architecture and technical decisions
- **Development process**: `docs/development-methodology.md` - Intelligent Development principles
- **Testing approach**: `docs/testing-strategy.md` - Unit test patterns
- **Omarchy integration**: `OMARCHY.md` - Theme format standards and requirements
- **Progress tracking**: `PROJECT.md` - Roadmap and session logs
- **Public overview**: `README.md` - User-facing documentation

## Development Rules
1. Operate in Explanatory mode (`/output-style explanatory`)
2. Provide comprehensive implementation guides for user
3. Use `go test ./tests/... -v` for full test validation
4. Use `go vet ./...` for type checking
5. Reference existing implementations: "See pkg/formats/hsla.go"
6. Keep explanations technically precise

## Architectural Patterns

### Settings-as-Methods Pattern
- **Public functions requiring settings MUST be methods** on package configuration structures
- Private functions MAY accept settings parameters from calling methods
- This enforces explicit configuration and prevents hidden dependencies

Example:
```go
// ✅ Correct: Public method on configuration structure
func (a *Analyzer) AnalyzeColors(colors []color.RGBA) ColorProfile

// ❌ Incorrect: Public function with settings parameter
func AnalyzeColors(colors []color.RGBA, settings *Settings) ColorProfile
```

### Algorithmic Constants Only
- **Hard-coded values allowed ONLY for algorithmic constants** (mathematical formulas, ratios)
- **Any tunable parameter MUST be a setting** (thresholds, tolerances, limits)
- If a value can be varied without breaking the algorithm, it belongs in settings

Example:
```go
// ✅ Correct: Mathematical constant
luminance := 0.2126*r + 0.7152*g + 0.0722*b

// ❌ Incorrect: Tunable threshold
if saturation < 0.05 { // Should be: if saturation < a.grayscaleThreshold {
```

### Foundation Layer Responsibility
- **pkg/formats**: Color space representations and conversions
- **pkg/chromatic**: Color theory, calculations, and analysis (foundation)
- **pkg/analysis**: High-level color profiles using chromatic
- Use descriptive package names that reflect actual responsibility

## Current Implementation Status
- ✅ Multi-strategy image extraction complete (frequency/saliency algorithms)
- ✅ Strategy selection with empirical thresholds implemented
- ✅ Foundation layer structure complete (pkg/formats, pkg/chromatic, pkg/settings, pkg/loader)
- ✅ LAB and XYZ color space implementations in pkg/formats
- ✅ Settings-as-methods architectural pattern established
- ⏳ Unit tests needed for all foundation packages
- ⏳ Color derivation algorithms in pkg/chromatic in development
- ⏳ pkg/analysis partially extracted from extractor
- ⏳ Strategy extraction pending (pkg/strategies with dependency injection)
- ⏳ Extractor simplification pending (pure orchestration)
- ⏳ Theme generation pending (pkg/theme/)
- ⏳ CLI interface pending (cmd/omarchy-theme-gen/)

## Key Technical Decisions
- **Standard Types**: Use `color.RGBA` from standard library, not custom types
- **Purpose-Driven**: Colors organized by role (background, foreground, accent) not frequency
- **Settings vs Config**: System settings (HOW tool operates) separate from user config (WHAT user wants)
- **Layered Architecture**: Clear dependency layers with no circular dependencies
- **Profile Detection**: Grayscale, Monotone, Monochromatic, Duotone/Tritone for edge cases
- **Multi-Strategy**: Frequency for simple images, saliency for complex images
- **Vocabulary precision**: IsGrayscale (saturation < 0.05) vs IsMonochromatic (±15° hue tolerance)
- **Early termination algorithm**: Monochromatic detection with 80% threshold
- **HEXA color format**: For theme-gen.json metadata
- **CLI sub-commands**: For theme refinement

## Commands
```bash
# Validate code
go vet ./...

# Run all tests
go test ./tests/... -v

# Run package-specific tests
go test ./tests/formats -v
go test ./tests/extractor -v

# Run specific test functions
go test ./tests/formats -run TestParseHex -v
go test ./tests/extractor -run TestStrategySelection -v

# Run tests with coverage
go test ./tests/... -v -cover

# Generate image analysis documentation
go run tests/analyze-images/main.go

# Format code
go fmt ./...
```

## Package Structure (Current State)

### Foundation Layer (Structure Complete, Tests Needed)
- `pkg/formats/` - Color space representations and conversions (RGBA, HSLA, LAB, XYZ)
- `pkg/chromatic/` - Color theory foundation (algorithms in development)
- `pkg/settings/` - Flat configuration structure with Viper integration
- `pkg/loader/` - Image I/O with validation and format support
- `pkg/config/` - User preferences per theme (future, for theme-gen.json)

### Analysis Layer (Partial Implementation)
- `pkg/analysis/` - Color profiles and analysis (partially extracted from extractor)

### Processing Layer (Refactoring Needed)
- `pkg/extractor/` - Currently contains extraction + embedded strategies
- `pkg/strategies/` - Pending extraction from extractor

### Generation Layer (Not Implemented)
- `pkg/theme/` - Theme file generation (future)

### Application Layer (Not Implemented)
- `cmd/omarchy-theme-gen/` - CLI application (future)

### Testing (Structure Established)
- `tests/formats/` - Unit tests for pkg/formats (in development)
- `tests/extractor/` - Tests for extraction strategies
- `tests/images/` - Real-world test wallpapers
- `tests/analyze-images/` - Image analysis utility

## Performance Targets
- 4K image: < 2 seconds
- Memory: < 100MB peak
- Contrast: WCAG AA (4.5:1)

## Current Development Focus
Unit Testing and Algorithm Implementation

### Foundation Testing (Priority)
- ⏳ Create unit tests for pkg/formats
- ⏳ Create unit tests for pkg/chromatic
- ⏳ Create unit tests for pkg/settings
- ⏳ Create unit tests for pkg/loader
- ⏳ Create unit tests for pkg/analysis

### Algorithm Development (Active)
- ⏳ Implement color derivation algorithms in pkg/chromatic
- ⏳ Complete harmony generation functions
- ⏳ Add perceptual distance calculations

### Architecture Refactoring Phase 2 (Next)
- Extract pkg/strategies from extractor
- Complete pkg/analysis extraction
- Implement role-based color organization
- Simplify pkg/extractor to pure orchestration

## CLI Architecture
Commands planned:
```bash
omarchy-theme-gen generate --image photo.[jpg|png] [options]
omarchy-theme-gen set-scheme <theme-name> --scheme [monochromatic|analogous|complementary|split-complementary|triadic|tetradic|square]
omarchy-theme-gen set-mode <theme-name> --mode [light|dark]
omarchy-theme-gen clone <theme-name> <new-name> [options]
```

Command Options:

| Command | Option | Description | Default |
|---------|--------|-------------|---------|
| **generate** | `background` | background color | derived from image |
| | `foreground` | foreground color | derived from image |
| | `primary` | primary theme color, used as basis for color scheme operations | derived from image |
| | `mode` | light vs. dark mode | derived from image luminescence |
| | `scheme` | color scheme to apply | derived from image analysis |
| **clone** | `background` | background color | inherited from source theme |
| | `foreground` | foreground color | inherited from source theme |
| | `primary` | primary theme color, used as basis for color scheme operations | inherited from source theme |
| | `mode` | light vs. dark mode | inherited from source theme |
| | `scheme` | color scheme to apply | inherited from source theme |

Note: No `apply` command needed - themes integrate directly with Omarchy's system theme selection.

## AI Responsibilities
- Provide comprehensive implementation guides
- Create unit tests for all packages
- Maintain documentation accuracy
- Review code for best practices
- Ensure cross-references are valid

## User Responsibilities
- Architecture and design decisions
- Source code implementation
- Review and refine AI outputs
- Project direction and priorities
