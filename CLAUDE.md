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
3. Use `go run cmd/examples/test_*.go` for validation
4. Use `go vet ./...` for type checking
5. Reference existing implementations: "See pkg/color/space.go"
6. Keep explanations technically precise

## Current Implementation Status
- ✅ Project structure established
- ✅ Core color types complete
- ⏳ Color space conversions pending (partial - RGB↔HSL done)
- ⏳ Image extraction pending
- ⏳ Palette strategies pending
- ⏳ Config generation pending
- ⏳ TUI interface pending

## Key Technical Decisions
- RGBA with cached HSLA for performance
- Octree quantization over k-means
- Template-based config generation
- 64x64 pixel regions for concurrency

## Commands
```bash
# Validate code
go vet ./...

# Run execution test
go run cmd/examples/test_[name].go

# Run with arguments
go run cmd/examples/test_extract.go image.jpg

# Format code
go fmt ./...
```

## Package Structure
- `pkg/color/` - Color types and conversions
- `pkg/quantizer/` - Quantization algorithms
- `pkg/extractor/` - Image processing
- `pkg/palette/` - Color theory strategies
- `pkg/template/` - Config generators
- `pkg/theme/` - Theme orchestration
- `internal/tui/` - UI components
- `cmd/examples/` - Execution tests
- `cmd/omarchy-theme-gen/` - Main application

## Performance Targets
- 4K image: < 2 seconds
- Memory: < 100MB peak
- Contrast: WCAG AA (4.5:1)

## Next Session Focus
Session 2: Color Space Conversions
- Note: RGB↔HSL conversion already implemented in Session 1
- Add additional color manipulation methods if needed
- Create comprehensive test_conversion.go execution test
- Verify all conversions against CSS Color Module Level 3
- Begin Session 3: Basic Image Loading

## Remember
- Start from fundamental understanding
- Build toward learnable challenges
- Document with technical precision
- Test immediately, adapt quickly
- Keep context clean and focused
