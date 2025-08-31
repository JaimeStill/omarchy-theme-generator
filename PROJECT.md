# Project Roadmap

## Development Phases

### Phase 1: Foundation (5 sessions)
Core types and basic functionality

### Phase 2: Color Theory Schemes (5 sessions)  
Advanced color extraction and color theory scheme generation

### Phase 3: Configuration Generation (5 sessions)
All Omarchy file formats with theme-gen.json metadata

### Phase 4: CLI Development (4 sessions)
Command-line interface with sub-commands

### Phase 5: Polish & Integration (4 sessions)
Advanced capabilities and Omarchy integration

### Phase 6: Optional TUI Enhancement (4 sessions)
Interactive interface with Bubble Tea (future enhancement)

### Phase 7: Finalization (2 sessions)
Documentation and release

---

## Session Plan

### Phase 1: Foundation

#### Session 1: Project Setup & Core Types
- [x] Initialize Go module: `go mod init omarchy-theme-gen`
- [x] Create project structure
- [x] Implement Color type with RGBA storage
- [x] Add HSLA caching mechanism
- [x] **Test**: `tests/test-color/main.go`
- [x] Document decisions

#### Session 2: Color Space Conversions
- [x] Implement RGB to HSL conversion
- [x] Implement HSL to RGB conversion
- [x] Add output format methods (Hex, CSS, HSLA)
- [x] **Test**: `tests/test-conversions/main.go`
- [x] Verify against CSS Color Module Level 3

#### Session 3: Basic Image Loading
- [x] Image loading from file path
- [x] Pixel iteration and color counting
- [x] **Test**: `tests/test-load-image/main.go`
- [x] Benchmark performance

#### Session 4: Color Synthesis & Palette Generation
- [x] **Vocabulary Correction**: Replace IsMonochrome with proper IsGrayscale and IsMonochromatic detection
- [x] Hue tolerance algorithm implemented - 15-degree tolerance with wraparound handling
- [x] Strategy decision logic updated - Grayscale images properly use color-based strategy
- [x] All references corrected in tests and documentation
- [x] **Test**: Empirical validation complete via `tests/test-load-image/main.go`

#### Session 5: Color Theory Schemes Implementation
- [ ] Create pkg/palette/ package with extensible SchemeGenerator interface
- [ ] Implement SchemeOptions configuration for flexible palette generation
- [ ] Build all 7 color theory scheme generators:
  - [ ] Monochromatic - Single hue with saturation/lightness variations
  - [ ] Analogous - Adjacent hues (±30° on color wheel)
  - [ ] Complementary - Opposite hues (180° separation)
  - [ ] Split-Complementary - Base hue + two adjacent to complement
  - [ ] Triadic - Three hues equally spaced (120° separation)
  - [ ] Tetradic - Rectangle pattern (60°, 120°, 180°)
  - [ ] Square - Four hues equally spaced (90° separation)
- [ ] Create scheme registry for dynamic selection
- [ ] Build extraction → hybrid → scheme pipeline integration
- [ ] Implement light/dark mode palette variations
- [ ] Add WCAG compliance validation for all schemes
- [ ] **Test**: `tests/test-color-schemes/main.go`

### Phase 2: Algorithms

#### Session 6: First Template Generator
- [ ] Create template interface with color theory scheme-compatible color mapping
- [ ] Implement alacritty.toml generator with synthesized color support
- [ ] Add color formatting functions for all color theory schemes
- [ ] **Test**: `tests/test-generate-alacritty/main.go`

#### Session 7: Octree Implementation (Optimization)
- [ ] Build octree data structure for efficient color quantization
- [ ] Implement color insertion and tree reduction algorithms
- [ ] Optimize memory usage and processing speed
- [ ] **Test**: `tests/test-octree/main.go`

#### Session 8: Dominant Color Detection (Optimization)
- [ ] Implement advanced color clustering with color theory scheme integration
- [ ] Add perceptual distance metrics for better color selection
- [ ] Compare extraction vs color theory scheme quality metrics
- [ ] **Test**: `tests/test-dominant/main.go`

#### Session 9: Concurrent Processing
- [ ] Divide image into 64x64 regions for parallel processing
- [ ] Implement parallel extraction with color theory scheme fallback coordination
- [ ] Add result aggregation and performance optimization
- [ ] **Test**: `tests/test-concurrent/main.go`

#### Session 10: Advanced Schemes & Custom Generators
- [ ] Implement custom scheme builders for edge cases
- [ ] Add weighted color schemes with user preferences
- [ ] Create double-complementary and other advanced schemes
- [ ] Build accessibility compliance reports for all generation modes
- [ ] Add scheme mutation and refinement capabilities
- [ ] **Test**: `tests/test-advanced-schemes/main.go`

### Phase 3: Configuration Generation

#### Session 11: Multiple Config Generators
- [ ] Implement mako.ini generator
- [ ] Add btop.theme generator
- [ ] **Test**: `tests/test-generate-configs/main.go`

#### Session 12: CSS Generation
- [ ] Create waybar.css generator
- [ ] Add walker.css generator
- [ ] Add swayosd.css generator
- [ ] **Test**: `tests/test-generate-css/main.go`

#### Session 13: Lua Generation
- [ ] Implement neovim.lua generator
- [ ] Map syntax highlighting groups
- [ ] **Test**: `tests/test-generate-lua/main.go`

#### Session 14: Hyprland Configuration
- [ ] Create hyprland.conf generator
- [ ] Add hyprlock.conf generator
- [ ] **Test**: `tests/test-generate-hypr/main.go`

#### Session 15: Complete Theme Package
- [ ] Assemble all generators
- [ ] Create directory structure
- [ ] Copy image to backgrounds/
- [ ] Add light.mode marker
- [ ] **Test**: `tests/test-full-theme/main.go`

### Phase 4: TUI Development

#### Session 16: Bubble Tea Setup
- [ ] Initialize Bubble Tea application
- [ ] Create basic model structure
- [ ] Implement key handling
- [ ] **Test**: Launch and quit

#### Session 17: File Selection & Theme Options
- [ ] Create file browser
- [ ] Add theme mode selector
- [ ] Add color override inputs
- [ ] **Test**: Navigation and input

#### Session 18: Palette Display & Strategy Selection
- [ ] Show dominant color
- [ ] Add strategy selector
- [ ] Display generated palette
- [ ] **Test**: Strategy switching

#### Session 19: Color Adjustment
- [ ] Create HSL sliders
- [ ] Implement real-time updates
- [ ] Add keyboard controls
- [ ] **Test**: Color manipulation

#### Session 20: Preview Component
- [ ] Design terminal mockup
- [ ] Show color applications
- [ ] Update with changes
- [ ] **Test**: Preview accuracy

#### Session 21: Export Dialog
- [ ] Create export view
- [ ] Add theme naming
- [ ] Handle file writing
- [ ] Copy image to backgrounds/
- [ ] **Test**: Export verification

#### Session 22: Full Integration
- [ ] Wire all components
- [ ] Add state management
- [ ] Implement navigation
- [ ] **Test**: Complete workflow

### Phase 5: Polish & Features

#### Session 23: History & Undo
- [ ] Implement command pattern
- [ ] Add history stack
- [ ] Create undo/redo handlers
- [ ] **Test**: State consistency

#### Session 24: Theme Variations
- [ ] Generate light/dark pairs
- [ ] Create color variations
- [ ] Add batch generation
- [ ] **Test**: Variation quality

#### Session 25: Batch Processing
- [ ] Support multiple images
- [ ] Add batch export
- [ ] Create comparison view
- [ ] **Test**: Batch operations

#### Session 26: Settings & Persistence
- [ ] Add configuration file
- [ ] Remember preferences
- [ ] Save recent files
- [ ] **Test**: Settings persistence

#### Session 27: Error Handling
- [ ] Add comprehensive error handling
- [ ] Implement recovery
- [ ] Create error messages
- [ ] **Test**: Edge cases

#### Session 28: Performance Optimization
- [ ] Profile application
- [ ] Optimize hot paths
- [ ] Reduce allocations
- [ ] **Test**: Performance benchmarks

### Phase 6: Finalization

#### Session 29: Documentation
- [ ] Write user guide
- [ ] Add help system
- [ ] Create example themes
- [ ] Document API

#### Session 30: Testing & Release
- [ ] Convert to formal tests
- [ ] Add integration suite
- [ ] Build binaries
- [ ] Create installation guide

---

## Progress Log

### Session Template
```markdown
### Session N: [Date]
**Completed:**
- ✅ Task description - reference to code/test
- ✅ Task description - reference to code/test

**Insights:**
- Key learning or discovery

**Decision:**
- Architectural choice made - link to docs/decisions/

**Next:**
- What to tackle in next session
```

### Sessions Completed

### Session 1: 2025-08-23
**Completed:**
- ✅ Go module initialized - `go.mod` created with module name `omarchy-theme-gen`
- ✅ Project structure established - `pkg/color/` and `tests/` directories
- ✅ Color type implemented - RGBA storage with lazy-cached HSLA conversion (`pkg/color/color.go`)
- ✅ Thread-safe caching added - `sync.Once` for HSLA computation
- ✅ Comprehensive testing - `test_color.go` validates all functionality including concurrency
- ✅ Documentation complete - Full godoc coverage for all public and private functions

**Insights:**
- Alpha standardized to 0.0-1.0 range throughout, removed opacity concept for API simplicity
- HSLA caching provides significant performance improvement over repeated conversion
- Value semantics for `WithAlpha()` ensures immutability of original colors

**Decision:**
- RGBA with cached HSLA chosen over pure HSL storage for native image processing performance
- Pointer return from constructors, value return from `WithAlpha()` for correct Go semantics
- Added `roundAlpha()` for consistent 3-decimal display in CSS output

**Next:**
- Session 3: Basic image loading functionality (completed)

### Session 2: 2025-08-24
**Completed:**
- ✅ Complete color manipulation infrastructure - 11 methods in `pkg/color/manipulation.go`
- ✅ WCAG contrast calculations with gamma correction - `pkg/color/contrast.go`
- ✅ Color distance metrics in multiple spaces - `pkg/color/distance.go` 
- ✅ LAB color space with Delta-E calculations - `pkg/color/lab.go`
- ✅ Transparent testing methodology - `tests/test-conversions/main.go`
- ✅ Comprehensive documentation - Full godoc coverage for all new functions

**Insights:**
- Gamma correction essential for accurate WCAG luminance calculation
- Delta-E CIE76/CIE94 provide perceptually-uniform color differences
- Transparent testing with detailed explanations improves understanding
- LAB color space crucial for professional color science applications

**Decision:**  
- AccessibilityLevel enum with automatic ratio lookup for type safety
- LAB conversion uses D65 illuminant for standard daylight conditions
- HSL distance weighting: lightness(2.0) > saturation(1.0) > hue(0.5)

**Next:**
- Session 4: Color theory schemes with vocabulary corrections

### Architectural Decision: Color Synthesis Pipeline (Session 3)
**Context:**
Images may lack sufficient color diversity for theme generation (grayscale, noir, monochrome cases).

**Decision:**
Implement extraction → hybrid → color theory scheme pipeline with automatic failover:
1. **Extraction**: Traditional image-based color extraction
2. **Hybrid**: Combine extracted colors with synthesized ones when insufficient diversity
3. **Synthesis**: Pure color theory-based generation when extraction fails

**Impact:**
- Sessions 4-5 restructured to prioritize color theory scheme architecture
- Sessions 6-10 reordered with color theory scheme integration
- All color theory schemes must support color generation modes
- Template generators must handle synthesized color palettes
- WCAG compliance required for both extracted and synthesized colors

**Technical Implementation:**
- `pkg/palette/` package for color theory schemes
- `SchemeOptions` configuration for flexible palette generation
- Color theory schemes: all 7 standard schemes (monochromatic through square)
- Extensible `SchemeGenerator` interface for future additions
- Edge case testing for low-diversity images

**Vocabulary Correction Required:**
- Current `IsMonochrome` actually detects grayscale images (saturation ≈ 0)
- Proper terminology: `IsGrayscale` (no hue) vs `IsMonochromatic` (single hue variations)
- Strategy implications: grayscale → synthesize, monochromatic → extract/hybrid
- Session 4 must implement proper color classification

### Session 3: 2025-08-25
**Completed:**
- ✅ Image loading infrastructure - `pkg/extractor/loader.go` with JPEG/PNG support
- ✅ Structured error handling - `pkg/errors/extractor.go` with comprehensive error types
- ✅ Color frequency mapping - `pkg/extractor/frequency.go` with optimized pixel access
- ✅ Main extraction pipeline - `pkg/extractor/extractor.go` with analysis-based validation
- ✅ Performance benchmarking - `pkg/extractor/performance.go` with 4K testing capabilities
- ✅ Comprehensive execution test - `tests/test-load-image/main.go` with visual samples

**Performance Achievements:**
- 4K image processing: 241ms (target: <2s) - **6x faster than target**
- Memory usage: 72MB (target: <100MB) - **28% under target**
- Processing rate: 34M pixels/second
- Edge case handling: Proper detection of grayscale and high-contrast scenarios

**Insights:**
- Analysis-based validation eliminates hard failures while providing color theory scheme guidance
- Type-specific pixel access (RGBA vs generic) provides significant performance improvements
- Visual test documentation dramatically improves comprehension of edge cases
- Vocabulary precision (monochromatic vs grayscale) is critical for accurate classification

**Architectural Decision:**
- Replaced strict validation with intelligent analysis that guides color theory schemes
- Extraction → Hybrid → Synthesis pipeline architecture established
- Sessions 4-10 restructured to prioritize color theory scheme integration

**Vocabulary Correction Identified:**
- Current `IsMonochrome` detects grayscale (no color information)
- Need separate `IsGrayscale` vs `IsMonochromatic` (single hue with variations)
- Strategy logic must distinguish: grayscale → synthesize, monochromatic → extract/hybrid

**Next:**
- Session 4: Implement vocabulary corrections and color theory schemes

### Maintenance Session: 2025-08-31
**Completed:**
- ✅ Testing infrastructure reorganization - Moved `pkg/extractor/performance.go` utilities to `tests/internal/` with proper separation (`tests/internal/benchmark.go`, `tests/internal/generators.go`, `tests/internal/suite.go`)
- ✅ Centralized test samples - Created `tests/samples/` directory with 6 reusable test images for consistent benchmarking across all tests
- ✅ HEXA color format support - Implemented `ParseHEXA()`, `ParseHEX()`, and `ParseHexString()` with comprehensive test coverage (Tests 12-13 in `test-color`)
- ✅ CLI-first architecture transition - Updated all documentation from TUI-first to CLI-first approach with optional future TUI enhancement
- ✅ Color theory terminology standardization - Fixed 21+ instances of "synthesis/strategies" → "color theory schemes" across all documentation
- ✅ CLI command structure refinement - Established proper command architecture with color theory scheme terminology
- ✅ Omarchy integration documentation - Created comprehensive `OMARCHY.md` style guide with theme format standards, color conversion requirements, and validation criteria
- ✅ Documentation consistency review - Engaged docs-consistency-checker to identify and fix inconsistencies across all project documents

**Architectural Decisions:**
- **CLI Commands**: Finalized structure with `generate`, `set-scheme`, `set-mode`, `clone` sub-commands
- **Color Storage**: HEXA format (#RRGGBBAA) for `theme-gen.json` metadata with format-specific conversion
- **Testing Strategy**: Centralized utilities in `tests/internal/` with reusable samples in `tests/samples/`
- **Terminology**: Proper color theory academic terminology (schemes vs strategies) throughout project
- **Integration Approach**: Direct Omarchy theme system compatibility without separate apply command

**Scope Refinements:**
- **Primary Focus**: Clean, reliable CLI implementation over premature performance optimization
- **Theme Refinement**: Cached extraction data in `theme-gen.json` enables iterative scheme adjustments
- **Phase Restructure**: Updated PROJECT.md phases to reflect CLI-first with optional TUI in Phase 6
- **Documentation Structure**: Moved style guide to `OMARCHY.md` as core reference document alongside `CLAUDE.md`

**Technical Improvements:**
- **Color Precision**: Full HEXA support preserves alpha channel information without conversion losses
- **Test Infrastructure**: Benchmark utilities properly separated from production code with importable test functions
- **Error Handling**: Enhanced validation with proper color parsing error messages and format verification
- **Cross-references**: Updated all internal documentation references to reflect new structure and terminology

**Impact on Development:**
- **Simplified Scope**: CLI-first approach reduces complexity while maintaining extensibility for future TUI
- **Consistent Terminology**: Proper color theory language aligns with academic standards and improves communication
- **Enhanced Maintainability**: Centralized testing infrastructure and clear documentation structure
- **Foundation Strengthened**: Solid architectural base established for Session 4 color theory scheme implementation

**Completed:**
- ✅ Testing infrastructure reorganization - Moved `pkg/extractor/performance.go` utilities to `tests/internal/` with proper separation (`tests/internal/benchmark.go`, `tests/internal/generators.go`, `tests/internal/suite.go`)
- ✅ Centralized test samples - Created `tests/samples/` directory with 6 reusable test images for consistent benchmarking across all tests
- ✅ HEXA color format support - Implemented `ParseHEXA()`, `ParseHEX()`, and `ParseHexString()` with comprehensive test coverage (Tests 12-13 in `test-color`)
- ✅ CLI-first architecture transition - Updated all documentation from TUI-first to CLI-first approach with optional future TUI enhancement
- ✅ Color theory terminology standardization - Fixed 21+ instances of "synthesis/strategies" → "color theory schemes" across all documentation
- ✅ CLI command structure refinement - Established proper command architecture with color theory scheme terminology
- ✅ Omarchy integration documentation - Created comprehensive `OMARCHY.md` style guide with theme format standards, color conversion requirements, and validation criteria
- ✅ Documentation consistency review - Engaged docs-consistency-checker to identify and fix inconsistencies across all project documents

**Architectural Decisions:**
- **CLI Commands**: Finalized structure with `generate`, `set-scheme`, `set-mode`, `clone` sub-commands
- **Color Storage**: HEXA format (#RRGGBBAA) for `theme-gen.json` metadata with format-specific conversion
- **Testing Strategy**: Centralized utilities in `tests/internal/` with reusable samples in `tests/samples/`
- **Terminology**: Proper color theory academic terminology (schemes vs strategies) throughout project
- **Integration Approach**: Direct Omarchy theme system compatibility without separate apply command

**Scope Refinements:**
- **Primary Focus**: Clean, reliable CLI implementation over premature performance optimization
- **Theme Refinement**: Cached extraction data in `theme-gen.json` enables iterative scheme adjustments
- **Phase Restructure**: Updated PROJECT.md phases to reflect CLI-first with optional TUI in Phase 6
- **Documentation Structure**: Moved style guide to `OMARCHY.md` as core reference document alongside `CLAUDE.md`

**Technical Improvements:**
- **Color Precision**: Full HEXA support preserves alpha channel information without conversion losses
- **Test Infrastructure**: Benchmark utilities properly separated from production code with importable test functions
- **Error Handling**: Enhanced validation with proper color parsing error messages and format verification
- **Cross-references**: Updated all internal documentation references to reflect new structure and terminology

**Impact on Development:**
- **Simplified Scope**: CLI-first approach reduces complexity while maintaining extensibility for future TUI
- **Consistent Terminology**: Proper color theory language aligns with academic standards and improves communication
- **Enhanced Maintainability**: Centralized testing infrastructure and clear documentation structure
- **Foundation Strengthened**: Solid architectural base established for Session 4 color theory scheme implementation

**Next:**
- Session 4: Color theory schemes with vocabulary corrections and CLI foundation

### Session 4: 2025-08-31
**Completed:**
- ✅ Vocabulary correction complete - `IsMonochrome` → `IsGrayscale` + `IsMonochromatic` with proper color science terminology (See `pkg/extractor/extractor.go:169-180`)
- ✅ Hue tolerance algorithm implemented - 15-degree tolerance with wraparound handling for monochromatic detection (`pkg/extractor/extractor.go:254-260`)
- ✅ Strategy decision logic updated - Grayscale images properly use color-based strategy logic (`pkg/extractor/extractor.go:244`)
- ✅ All references corrected - Tests, documentation, and internal suite updated with precise vocabulary (`tests/internal/suite.go:54`, `tests/test-load-image/main.go:111-115`)
- ✅ Empirical validation complete - Test execution confirms correct grayscale/monochromatic classification (`tests/test-load-image/README.md`)

**Not Completed (moved to Session 5):**
- pkg/palette/ package creation with color theory schemes
- SchemeOptions configuration implementation
- Extraction → hybrid → scheme pipeline integration
- Color scheme validation and testing

**Insights:**
- Early termination algorithm more elegant than accumulation for monochromatic detection
- Proper color science vocabulary critical for extraction→hybrid→scheme pipeline accuracy
- Session focused on foundational vocabulary corrections rather than full implementation

**Decision:**
- Monochromatic detection uses ±15° hue tolerance with 80% threshold requirement
- Grayscale threshold tightened from 0.1 to 0.05 saturation for precision
- Strategy logic distinguishes grayscale (scheme candidate) from monochromatic (extraction viable)
- Session 5 expanded to include all 7 standard color theory schemes with extensible architecture

**Next:**
- Session 5: Complete color theory schemes implementation with all 7 standard schemes

---

## Metrics Tracking

| Metric | Target | Current | Status |
|--------|--------|---------|---------|
| 4K Processing | < 2s | 241ms | ✅ |
| Memory Usage | < 100MB | 72MB | ✅ |
| WCAG Compliance | AA (4.5:1) | Infrastructure ready | ⏳ |
| Color Theory Schemes | 7+ | 0 (Session 4) | ⏳ |
| Extraction Strategies | 3+ | 3 (frequency, type-specific, generic) | ✅ |
| Config Formats | 9 | 0 (Sessions 6+) | ⏳ |
| Edge Case Support | 100% | Analysis ready (needs color theory schemes) | ⏳ |
| Test Coverage | 80% | Extraction pipeline covered | ⏳ |

---

## Links

- [Technical Specification](docs/technical-specification.md)
- [Development Methodology](docs/development-methodology.md)
- [Testing Strategy](docs/testing-strategy.md)
- [Memory File](CLAUDE.md)
- [Public README](README.md)
