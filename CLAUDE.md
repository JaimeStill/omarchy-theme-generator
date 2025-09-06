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
- **Architecture design**: `ARCHITECTURE.md` - Layered architecture and technical decisions
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
7. **All unit tests MUST include diagnostic logging** - Use `t.Logf()` to output calculation metrics, expected vs actual values, thresholds, and intermediate results for debugging

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
- **pkg/processor**: Unified processing pipeline combining extraction and analysis
- Use descriptive package names that reflect actual responsibility

## Current Implementation Status

### ✅ Completed Architecture (Foundation + Processing Layers)
- **pkg/formats**: Complete with comprehensive unit tests (RGBA, HSLA, LAB, XYZ, hex parsing)
- **pkg/chromatic**: Complete with comprehensive unit tests (color theory, harmony, contrast)
- **pkg/settings**: Complete with comprehensive unit tests (Viper integration, fallback configs)
- **pkg/loader**: Complete with comprehensive unit tests (JPEG/PNG loading, validation)
- **pkg/processor**: Complete unified processing with comprehensive unit tests using real images
- **Settings-as-methods pattern**: Enforced across all packages
- **ColorProfile composition**: Embedded ImageColors with complete metadata
- **Performance validated**: 100% compliance with <2s/100MB targets (88% faster than target)
- **Documentation**: Complete overhaul reflecting unified architecture

### ✅ Eliminated Packages (40-60% Performance Improvement)
- ❌ **pkg/analysis** → Merged into pkg/processor
- ❌ **pkg/extractor** → Merged into pkg/processor  
- ❌ **pkg/strategies** → Eliminated (frequency-only approach)

### ✅ Testing & Tools Infrastructure Complete
- **tests/**: Package-specific unit tests with diagnostic logging standards
- **tools/analyze-images/**: Image analysis documentation generator
- **tools/performance-test/**: Statistical performance validation across test dataset

### 🔄 Next Development Phase: Theme Generation
- ⏳ **pkg/palette**: Theme color derivation from ColorProfile metadata (will use pkg/chromatic algorithms)
- ⏳ **pkg/theme**: Omarchy configuration file generation
- ⏳ **cmd/omarchy-theme-gen/**: CLI application interface

## Key Technical Decisions
- **Standard Types**: Use `color.RGBA` from standard library, not custom types
- **Category-Based**: Colors organized by 27 theme categories with configurable characteristics
- **Settings vs Config**: System settings (HOW tool operates) separate from user config (WHAT user wants)
- **Layered Architecture**: Clear dependency layers with no circular dependencies
- **Profile Detection**: Grayscale, Monotone, Monochromatic, Duotone/Tritone for edge cases
- **Single-Strategy**: Optimized frequency-based extraction for all image types
- **Vocabulary precision**: IsGrayscale (saturation < 0.05) vs IsMonochromatic (±15° hue tolerance)
- **Early termination algorithm**: Monochromatic detection with 80% threshold
- **Color Storage**: HEXA format (#RRGGBBAA) in theme-gen.json for human readability
- **Color Bridge**: ParseHexA function to convert HEXA → color.RGBA
- **pkg/chromatic vs pkg/palette**: Foundation algorithms vs complete palette generation
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

## Package Structure

### Foundation Layer (Structure Complete)
- `pkg/formats/` - Color space representations and conversions (RGBA, HSLA, LAB, XYZ)
  - Needs: ParseHexA function for HEXA (#RRGGBBAA) parsing
- `pkg/chromatic/` - Color theory algorithms and calculations
- `pkg/settings/` - Flat configuration structure with Viper integration
- `pkg/loader/` - Image I/O with validation and format support

### Generation Layer (Not Implemented)
- `pkg/palette/` - Theme palette generation from processor output
- `pkg/theme/` - Theme file generation from templates

### Processing Layer (Complete)
- `pkg/processor/` - Unified processing pipeline with frequency-based extraction
- `pkg/palette/` - Pending creation for palette generation engine

### Generation Layer (Not Implemented)
- `pkg/config/` - User preferences per theme (theme-gen.json)
- `pkg/theme/` - Theme file generation with templates

### Application Layer (Not Implemented)
- `cmd/omarchy-theme-gen/` - CLI application

### Testing (Structure Established)
- `tests/*/` - Unit tests per package (comprehensive coverage pending)

## Performance Targets
- 4K image: < 2 seconds
- Memory: < 100MB peak
- Contrast: WCAG AA (4.5:1)

## Current Development Focus
Phase 1: Complete Refactoring & Testing (See PROJECT.md for detailed roadmap)

### Active Tasks
1. Documentation updates (in progress)
2. Strategy extraction from pkg/extractor
3. Category-based extraction and scoring implementation
4. pkg/palette package creation
5. Comprehensive test coverage for all packages

### Development Phases
- **Phase 1**: Complete refactoring & testing (current)
- **Phase 2**: Theme generation implementation
- **Phase 3**: CLI application

Refer to PROJECT.md for complete roadmap with task breakdowns.

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
