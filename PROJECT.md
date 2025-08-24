# Project Roadmap

## Development Phases

### Phase 1: Foundation (5 sessions)
Core types and basic functionality

### Phase 2: Algorithms (5 sessions)
Advanced color extraction and palette generation

### Phase 3: Configuration Generation (5 sessions)
All Omarchy file formats

### Phase 4: TUI Development (7 sessions)
Interactive interface with Bubble Tea

### Phase 5: Polish & Features (6 sessions)
Advanced capabilities and optimization

### Phase 6: Finalization (2 sessions)
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
- [ ] Image loading from file path
- [ ] Pixel iteration and color counting
- [ ] **Test**: `tests/test-load-image/main.go`
- [ ] Benchmark performance

#### Session 4: Color Theory-Based Extraction
- [ ] Build color frequency map
- [ ] Implement dominant color detection
- [ ] Add palette strategies (mono, complementary, triadic, analogous)
- [ ] Handle light/dark mode detection
- [ ] **Test**: `tests/test-extract-strategies/main.go`

#### Session 5: First Template Generator
- [ ] Create template interface
- [ ] Implement alacritty.toml generator
- [ ] Add color formatting functions
- [ ] **Test**: `tests/test-generate-alacritty/main.go`

### Phase 2: Algorithms

#### Session 6: Octree Implementation
- [ ] Build octree data structure
- [ ] Implement color insertion
- [ ] Add tree reduction
- [ ] **Test**: `tests/test-octree/main.go`

#### Session 7: Dominant Color Detection
- [ ] Implement color clustering
- [ ] Add perceptual distance metrics
- [ ] Compare detection methods
- [ ] **Test**: `tests/test-dominant/main.go`

#### Session 8: Concurrent Processing
- [ ] Divide image into 64x64 regions
- [ ] Implement parallel extraction
- [ ] Add result aggregation
- [ ] **Test**: `tests/test-concurrent/main.go`

#### Session 9: Advanced Palette Strategies
- [ ] Implement tetradic scheme
- [ ] Add split-complementary
- [ ] Create weighted strategies
- [ ] **Test**: `tests/test-advanced-harmony/main.go`

#### Session 10: Accessibility Validation
- [ ] Implement WCAG contrast calculation
- [ ] Add automatic adjustment
- [ ] Create validation reports
- [ ] **Test**: `tests/test-contrast/main.go`

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
- Session 3: Basic image loading with pixel iteration and performance validation

---

## Metrics Tracking

| Metric | Target | Current | Status |
|--------|--------|---------|---------|
| 4K Processing | < 2s | - | ⏳ |
| Memory Usage | < 100MB | - | ⏳ |
| WCAG Compliance | AA (4.5:1) | - | ⏳ |
| Palette Strategies | 5+ | 0 | ⏳ |
| Config Formats | 9 | 0 | ⏳ |
| Test Coverage | 80% | - | ⏳ |

---

## Links

- [Technical Specification](docs/technical-specification.md)
- [Development Methodology](docs/development-methodology.md)
- [Testing Strategy](docs/testing-strategy.md)
- [Memory File](CLAUDE.md)
- [Public README](README.md)
