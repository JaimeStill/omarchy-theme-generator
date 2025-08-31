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
- **Testing approach**: `docs/testing-strategy.md` - Execution test patterns
- **Omarchy integration**: `OMARCHY.md` - Theme format standards and requirements
- **Progress tracking**: `PROJECT.md` - Roadmap and session logs
- **Public overview**: `README.md` - User-facing documentation

## Development Rules
1. Operate in Explanatory mode (`/output-style explanatory`)
2. Only modify code when explicitly directed
3. Use `go run tests/test-*/main.go` for validation
4. Use `go vet ./...` for type checking
5. Reference existing implementations: "See pkg/color/space.go"
6. Keep explanations technically precise

## Current Implementation Status
- ✅ Project structure established
- ✅ Core color types complete
- ✅ Color space conversions complete (RGB↔HSL, manipulation, WCAG, LAB)
- ✅ Image extraction with vocabulary corrections complete
- ✅ Grayscale vs monochromatic classification implemented
- ⏳ Color theory schemes pending (pkg/palette/ package)
- ⏳ Palette generation pipeline pending
- ⏳ Config generation pending
- ⏳ CLI interface pending
- ⏳ TUI interface (optional future enhancement)

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

# Run execution test
go run tests/test-color/main.go
go run tests/test-conversions/main.go

# Run with arguments (for future image tests)
go run tests/test-extract/main.go image.jpg

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
- `tests/` - Execution tests
- `cmd/omarchy-theme-gen/` - CLI application
- `internal/tui/` - UI components (future enhancement)

## Performance Targets
- 4K image: < 2 seconds
- Memory: < 100MB peak
- Contrast: WCAG AA (4.5:1)

## Next Session Focus
Session 5: Color Theory Schemes Implementation
- Create pkg/palette/ package with color theory scheme generators
- Implement monochromatic, analogous, complementary, and triadic schemes
- Build SynthesisOptions configuration for fallback scenarios
- Integrate schemes with extraction→hybrid→synthesis pipeline
- Test color theory schemes with grayscale and monochromatic images

## CLI Architecture
Commands planned:
```bash
omarchy-theme-gen generate --image photo.[jpg|png] [options]
omarchy-theme-gen set-scheme <theme-name> --scheme [monochromatic|analogous|complementary|split-complementary|triadic|tetradic|square]
omarchy-theme-gen set-mode <theme-name> --mode [light|dark]
omarchy-theme-gen clone <theme-name> <new-name>
```

Note: No `apply` command needed - themes integrate directly with Omarchy's system theme selection.

## Remember
- Start from fundamental understanding
- Build toward learnable challenges
- Document with technical precision
- Test immediately, adapt quickly
- Keep context clean and focused
