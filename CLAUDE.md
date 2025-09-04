# Omarchy Theme Generator

## Project Context
Go-based CLI tool that generates Omarchy themes from images using color extraction and palette generation based on color theory principles. Optional TUI interface planned for future enhancement.

## Development Philosophy
- **User-driven**: All code modifications require explicit user direction
- **Explanatory mode**: Provide educational insights while coding
- **Execution tests**: Validate immediately with minimal tests, no frameworks
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
2. Only modify code when explicitly directed
3. Use `go test ./tests/formats ./tests/extractor -v` for validation
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
- ✅ Foundation layer refactoring complete (pkg/formats, pkg/chromatic, pkg/settings, pkg/loader)
- ✅ LAB color space implementation with mathematical accuracy
- ✅ Color analysis and profile detection (pkg/analysis with Analyzer pattern)
- ✅ Flat settings architecture with Viper integration
- ✅ Color theory foundation (pkg/chromatic with harmony detection)
- ✅ Settings-as-methods architectural pattern established
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

# Run tests
go test ./tests/formats ./tests/extractor -v

# Run specific test suites
go test ./tests/formats -run TestParseHex -v
go test ./tests/extractor -run TestStrategySelection -v

# Generate image analysis documentation
go run tests/analyze-images/main.go

# Format code
go fmt ./...
```

## Package Structure (Refactored Architecture)

### Foundation Layer
- `pkg/formats/` - Color space representations and conversions (RGBA, HSLA, LAB, XYZ)
- `pkg/chromatic/` - Color theory, calculations, and analysis (foundational color science)
- `pkg/settings/` - Flat configuration structure with Viper integration
- `pkg/loader/` - Image I/O with validation and format support

### Analysis Layer  
- `pkg/analysis/` - High-level color profiles and comprehensive analysis (uses chromatic)

### Processing Layer
- `pkg/strategies/` - Extraction strategies with dependency injection (pending)
- `pkg/extractor/` - Extraction orchestration (simplified, pending)

### Generation Layer
- `pkg/theme/` - Theme file generation (pending)

### Application Layer
- `cmd/omarchy-theme-gen/` - CLI application (pending)

### Testing
- `tests/internal/` - Centralized test utilities and benchmarks
- `tests/samples/` - Reusable test images
- `tests/` - Standard Go test files

## Performance Targets
- 4K image: < 2 seconds
- Memory: < 100MB peak
- Contrast: WCAG AA (4.5:1)

## Current Development Focus
Architecture Refactoring

### Foundation Refactoring (Complete)
- ✅ Transform pkg/color → pkg/formats with standard library types
- ✅ Implement HSLA color space with alpha channel support
- ✅ Add WCAG accessibility calculations with proper types
- ✅ Create comprehensive color analysis utilities
- ✅ Reorganize tests into package-specific structure

### Purpose-Driven Extraction (Next)
- Extract pkg/analysis and pkg/strategies from extractor
- Implement role-based color organization
- Add profile detection and synthesis capabilities
- Build settings-driven configuration system

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

## Remember
- Start from fundamental understanding
- Build toward learnable challenges
- Document with technical precision
- Test immediately, adapt quickly
- Keep context clean and focused
