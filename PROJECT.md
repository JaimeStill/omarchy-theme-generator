# Omarchy Theme Generator - Project Status

## Current Implementation

### Infrastructure
- **Go module**: `github.com/JaimeStill/omarchy-theme-generator`
- **Go version**: 1.25.0
- **Binary name**: `omarchy-theme-gen`

### Packages Implemented

#### pkg/extractor
- Multi-strategy extraction (frequency and saliency)
- Automatic strategy selection based on image analysis
- Image characteristic analysis (edge detection, complexity)
- Settings-driven configuration with empirical thresholds
- **Status**: Working but needs decomposition for refactoring

#### pkg/color
- Color type with RGBA storage and HSL conversion
- Contrast calculation for WCAG compliance
- Distance metrics (RGB, HSL, LAB)
- Hex color parsing and formatting
- **Status**: Over-engineered, needs simplification to pkg/formats

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

### Architecture Refactoring (Active)
**Goal**: Transform from frequency-based to purpose-driven extraction with layered architecture

**Tasks**:
- [ ] Complete documentation cleanup (Session 1 - Current)
- [ ] Refactor pkg/color → pkg/formats with standard library types
- [ ] Extract pkg/analysis and pkg/strategies from extractor
- [ ] Implement role-based color organization
- [ ] Create pkg/settings and pkg/config packages
- [ ] Add profile detection (Grayscale, Monotone, Monochromatic, Duotone/Tritone)

---

## Components & Features (Ordered by Dependency)

### Layer 1: Foundation

#### **pkg/formats** - Data formatting and conversion
*Purpose*: Handle color transformations using standard library types

**Features**:
- `RGBToHSL()` - Convert colors to HSL for analysis
- `ContrastRatio()` - WCAG accessibility calculations
- `ToHex()`, `ToHexA()` - Color formatting
- `ParseHex()` - Hex string to color parsing
- Theme type definitions (ColorRole, ThemeMode, etc.)

*Dependencies*: Standard library only
*Status*: Needs creation (refactor from pkg/color)

#### **pkg/settings** - System configuration
*Purpose*: Tool behavior and operational thresholds

**Features**:
- Settings structure with layered composition
- Default values and empirical thresholds
- JSON loading from multiple sources
- Multi-layer override system (defaults → system → user → workspace → env)

*Dependencies*: Standard library only
*Status*: Not implemented

#### **pkg/config** - User preferences
*Purpose*: Theme-specific user overrides

**Features**:
- User preference structure for themes
- Color overrides and extraction hints
- Theme-gen.json integration
- Per-theme storage and retrieval

*Dependencies*: pkg/formats
*Status*: Not implemented

---

### Layer 2: Analysis

#### **pkg/analysis** - Image and color analysis
*Purpose*: Analyze images to determine extraction approach and profile

**Features**:
- Image profile detection (Grayscale, Monotone, Monochromatic, Duotone/Tritone)
- Theme mode detection (light/dark based on luminance)
- Color clustering and perceptual grouping
- Edge detection and complexity analysis
- Role assignment logic for purpose-driven extraction

*Dependencies*: pkg/formats, pkg/settings
*Status*: Not yet extracted from pkg/extractor

---

### Layer 3: Processing

#### **pkg/strategies** - Extraction strategies
*Purpose*: Pluggable extraction algorithms

**Features**:
- Strategy interface for extensibility
- Frequency strategy for simple images
- Saliency strategy for complex images
- Strategy selector based on image characteristics
- Configurable thresholds and parameters

*Dependencies*: pkg/formats, pkg/analysis, pkg/settings
*Status*: Not yet extracted from pkg/extractor

#### **pkg/extractor** - Extraction orchestration
*Purpose*: Coordinate the extraction pipeline

**Features**:
- Pipeline coordination and orchestration
- Result aggregation and validation
- Profile-specific processing workflows
- Integration with analysis and strategies

*Dependencies*: pkg/formats, pkg/analysis, pkg/strategies, pkg/settings
*Status*: Needs simplification after decomposition

---

### Layer 4: Generation

#### **pkg/schemes** - Color scheme generation
*Purpose*: Apply color theory to generate complete palettes

**Features**:
- Color theory schemes (complementary, triadic, etc.)
- Synthesis for minimal-color images
- WCAG compliance validation
- Role-based scheme application

*Dependencies*: pkg/formats, pkg/analysis, pkg/config
*Status*: Not implemented

#### **pkg/theme** - Theme configuration generation
*Purpose*: Generate theme configuration files

**Features**:
- Template-based generation for supported formats
- Role → configuration mapping
- Format-specific color conversion
- Metadata generation (theme-gen.json)

*Dependencies*: pkg/formats, pkg/schemes, pkg/config
*Status*: Not implemented

---

### Layer 5: Application

#### **cmd/omarchy-theme-gen** - CLI application
*Purpose*: Command-line interface for theme generation

**Features**:
- `generate` - Create theme from image
- `set-scheme` - Apply color theory scheme
- `set-mode` - Switch light/dark mode
- `clone` - Duplicate and modify theme
- Settings and preferences management

*Dependencies*: All packages
*Status*: Not implemented

---

## Implementation Notes

### Simplification Principles
1. **Use standard library types** where possible (`color.RGBA`)
2. **Only build what's needed** - no speculative features
3. **Settings over hardcoding** - all thresholds configurable
4. **Clear dependency layers** - each package has specific purpose
5. **Purpose-driven organization** - role-based over frequency-based

### File Structure (Target State)
```
omarchy-theme-generator/
├── pkg/
│   ├── formats/        # Color utilities (refactor from pkg/color)
│   ├── settings/       # System configuration (new)
│   ├── config/         # User preferences (new)  
│   ├── analysis/       # Image analysis (extract from extractor)
│   ├── strategies/     # Extraction strategies (extract from extractor)
│   ├── extractor/      # Orchestration (simplify)
│   ├── schemes/        # Color theory schemes (new)
│   └── theme/          # Theme generation (new)
├── cmd/
│   └── omarchy-theme-gen/
├── tests/
│   ├── internal/       # Test utilities
│   ├── samples/        # Reusable test images
│   └── *_test.go       # Standard Go tests
└── docs/
    ├── architecture.md
    ├── glossary.md
    └── ...
```

### Next Implementation Steps
1. **Complete documentation cleanup** (Current work)
2. **Refactor pkg/color → pkg/formats**
3. **Decompose extractor package**
4. **Implement purpose-driven extraction**
5. **Add scheme generation**
6. **Create theme generators**
7. **Build CLI interface**

---

## Completed Features

### Foundation Work (Sessions 1-4)
- ✅ Project structure and Go module setup
- ✅ Color type with RGBA storage and cached HSLA conversion
- ✅ Color space conversions (RGB↔HSL, manipulation, WCAG, LAB)
- ✅ Image loading and validation infrastructure
- ✅ Multi-strategy extraction system (frequency vs saliency)
- ✅ Strategy selection based on image characteristics
- ✅ Settings-driven configuration with empirical thresholds
- ✅ Grayscale vs monochromatic classification with proper vocabulary
- ✅ Comprehensive test suite with real wallpaper validation
- ✅ Performance optimization (<2s for 4K images)

### Architectural Decisions Made
- **RGBA with cached HSLA**: Chosen over pure HSL for performance
- **Multi-strategy extraction**: Frequency for simple, saliency for complex images
- **Settings-driven configuration**: All thresholds configurable, no hardcoded values
- **Vocabulary precision**: IsGrayscale vs IsMonochromatic with proper definitions
- **Early termination algorithms**: 80% threshold for monochromatic detection
- **CLI-first architecture**: Optional TUI enhancement in future
- **HEXA color format**: For theme-gen.json metadata preservation

---

## Testing Strategy

- Unit tests for each package using standard Go test files
- Integration tests for complete pipeline validation  
- Benchmark tests for performance monitoring
- Real image validation with diverse wallpaper samples
- Test utilities centralized in `tests/internal/`

---

## Performance Metrics

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| 4K Processing | <2s | 241ms | ✅ |
| Memory Usage | <100MB | 72MB | ✅ |
| Extraction Strategies | 2+ | 2 | ✅ |
| WCAG Compliance | AA (4.5:1) | Infrastructure ready | ⏳ |
| Color Schemes | 7+ | 0 | ⏳ |
| Config Formats | 9 | 0 | ⏳ |

---

## References

- [Architecture Documentation](docs/architecture.md)
- [Development Methodology](docs/development-methodology.md)  
- [Testing Strategy](docs/testing-strategy.md)
- [Omarchy Integration](OMARCHY.md)
- [Glossary](docs/glossary.md)
