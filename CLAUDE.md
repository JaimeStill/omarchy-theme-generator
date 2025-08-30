# Omarchy Theme Generator

## Project Context
Go-based TUI application that generates Omarchy themes from images using color extraction and palette generation based on color theory principles.

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
- **Progress tracking**: `PROJECT.md` - Roadmap and session logs
- **Public overview**: `README.md` - User-facing documentation

## Development Rules
1. Operate in Explanatory mode (`/output-style explanatory`)
2. Only modify code when explicitly directed
3. Use `go run tests/test-*/main.go` for validation
4. Use `go vet ./...` for type checking
5. Reference existing implementations: "See pkg/color/hsl.go, pkg/color/lab.go"
6. Keep explanations technically precise

## Current Implementation Status
- ✅ Project structure established
- ✅ Core color types complete
- ✅ Color space conversions complete (RGB↔HSL, manipulation, WCAG, LAB)
- ✅ Image extraction with synthesis fallback complete
- ✅ Color synthesis strategies complete (6 strategies: monochromatic, analogous, complementary, triadic, tetradic, split-complementary)
- ✅ Palette generation pipeline complete (extraction → hybrid → synthesis)
- ✅ Computational generative system complete (material simulation, mathematical precision)
- ⏳ Theme orchestration integration pending
- ⏳ Config generation pending
- ⏳ TUI interface pending

## Key Technical Decisions
- RGBA with cached HSLA for performance
- AccessibilityLevel enum with automatic ratio lookup
- LAB color space with D65 illuminant for color science accuracy
- HSL distance weighting: lightness(2.0) > saturation(1.0) > hue(0.5)
- Extraction → Hybrid → Synthesis pipeline for edge case handling
- Color synthesis strategies for low-diversity images
- Octree quantization over k-means
- Template-based config generation
- 64x64 pixel regions for concurrency

## Commands
```bash
# Validate code
go vet ./...

# Run execution test
go run tests/test-color/main.go
go run tests/test-conversions/main.go

# Run with arguments (current image tests)
go run tests/test-load-image/main.go image.jpg

# Format code
go fmt ./...
```

## Package Structure
- `pkg/color/` - Color types and conversions
- `pkg/errors/` - Structured error handling
- `pkg/extractor/` - Image processing
- `pkg/generative/` - Computational image generation
- `pkg/palette/` - Color theory strategies
- `pkg/theme/` - Theme orchestration (Session 5)
- `pkg/template/` - Config generators (future)
- `internal/tui/` - UI components (future)
- `tests/` - Execution tests
- `cmd/omarchy-theme-gen/` - Main application (future)

## Performance Targets
- 4K image: < 2 seconds
- Memory: < 100MB peak
- Contrast: WCAG AA (4.5:1)

## Next Session Focus
Session 5: Theme Orchestration Integration (Current)
- Integrate extraction + synthesis pipeline with theme generation
- Implement light/dark mode detection with WCAG-accurate luminance analysis
- Add user color overrides with synthesis compatibility
- Complete all palette strategies with accessibility compliance
- Test with `tests/test-palette-strategies/main.go` execution test

Session 6: First Template Generator
- Create template interface with synthesis-compatible color mapping
- Implement alacritty.toml generator with synthesized color support
- Add color formatting functions for all synthesis strategies
- Test with `tests/test-generate-alacritty/main.go`

## Remember
- Start from fundamental understanding
- Build toward learnable challenges
- Document with technical precision
- Test immediately, adapt quickly
- Keep context clean and focused
