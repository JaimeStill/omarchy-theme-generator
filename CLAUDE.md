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
- **Technical details**: `docs/technical-specification.md` - Algorithms, architecture, performance targets
- **Development process**: `docs/development-methodology.md` - Intelligent Development principles
- **Testing approach**: `docs/testing-strategy.md` - Unit test patterns
- **Omarchy integration**: `OMARCHY.md` - Theme format standards and requirements
- **Progress tracking**: `PROJECT.md` - Roadmap and session logs
- **Public overview**: `README.md` - User-facing documentation

## Development Rules
1. Operate in Explanatory mode (`/output-style explanatory`)
2. Only modify code when explicitly directed
3. Use `go test ./tests -v` for validation
4. Use `go vet ./...` for type checking
5. Reference existing implementations: "See pkg/color/hsl.go"
6. Keep explanations technically precise

## Current Implementation Status
- ✅ Project structure established
- ✅ Core color types complete (pkg/color/)
- ✅ Color space conversions complete (RGB↔HSL, manipulation, WCAG, LAB)
- ✅ Multi-strategy image extraction complete (frequency/saliency algorithms)
- ✅ Strategy selection with empirical thresholds implemented
- ✅ Grayscale vs monochromatic classification implemented
- ✅ Unit test suite with real wallpaper validation
- ⏳ Color theory schemes pending (pkg/palette/ package)
- ⏳ Template-based config generation pending (pkg/template/)
- ⏳ Theme orchestration pending (pkg/theme/)
- ⏳ CLI interface pending (cmd/omarchy-theme-gen/)
- 📋 TUI interface (optional future enhancement)

## Key Technical Decisions
- RGBA with cached HSLA for performance
- AccessibilityLevel enum with automatic ratio lookup
- LAB color space with D65 illuminant for color science accuracy
- HSL distance weighting: lightness(2.0) > saturation(1.0) > hue(0.5)
- Extraction → Hybrid → Scheme Generation pipeline for edge case handling
- Vocabulary precision: IsGrayscale (saturation < 0.05) vs IsMonochromatic (±15° hue tolerance)
- Early termination algorithm for monochromatic detection with 80% threshold
- Color theory schemes for low-diversity images (monochromatic, analogous, complementary, etc.)
- Octree quantization over k-means
- Template-based config generation
- 64x64 pixel regions for concurrency
- HEXA color format for theme-gen.json metadata
- CLI sub-commands for theme refinement

## Commands
```bash
# Validate code
go vet ./...

# Run tests
go test ./tests -v

# Run specific test suites
go test ./tests -run TestStrategySelection -v

# Generate image analysis documentation
go run tests/analyze-images/main.go

# Format code
go fmt ./...
```

## Package Structure
- `pkg/color/` - Color types, conversions, and HEXA parsing
- `pkg/quantizer/` - Quantization algorithms
- `pkg/extractor/` - Image processing and color extraction
- `pkg/palette/` - Color theory schemes
- `pkg/template/` - Config generators and theme-gen.json
- `pkg/theme/` - Theme orchestration and CLI commands
- `pkg/metadata/` - Theme metadata serialization
- `tests/internal/` - Centralized test utilities and benchmarks
- `tests/samples/` - Reusable test images
- `tests/` - Unit tests and validation
- `cmd/omarchy-theme-gen/` - CLI application
- `internal/tui/` - UI components (future enhancement)

## Performance Targets
- 4K image: < 2 seconds
- Memory: < 100MB peak
- Contrast: WCAG AA (4.5:1)

## Next Session Focus
Session 5: Color Theory Schemes Implementation
- Create pkg/palette/ package with color theory scheme generators
- Implement core schemes: monochromatic, analogous, complementary, split-complementary, triadic, tetradic, square
- Build SchemeOptions configuration for flexible palette generation  
- Design SchemeGenerator interface for extensibility
- Integrate schemes with existing multi-strategy extraction system
- Validate scheme generation with test suite

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
