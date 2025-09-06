# Architecture Documentation

## System Architecture

The omarchy-theme-generator uses a layered architecture with clear dependencies and separation of concerns, implementing a sophisticated 27-category color extraction system with multi-dimensional scoring algorithms.

```
┌─────────────────────────────────────────────────────────┐
│                Application Layer (Future)               │
│                  cmd/omarchy-theme-gen                  │
├─────────────────────────────────────────────────────────┤
│                  Generation Layer (Future)              │
│         ┌──────────────┐      ┌──────────────┐          │
│         │   palette    │      │    theme     │          │
│         └──────────────┘      └──────────────┘          │
├─────────────────────────────────────────────────────────┤
│                  Processing Layer                       │
│                ┌──────────────────┐                     │
│                │    processor     │                     │
│                └──────────────────┘                     │
├─────────────────────────────────────────────────────────┤
│                 Foundation Layer                        │
│  ┌─────────┐ ┌───────────┐ ┌─────────┐ ┌─────────┐      │
│  │ formats │ │ chromatic │ │ settings│ │  loader │      │
│  └─────────┘ └───────────┘ └─────────┘ └─────────┘      │
└─────────────────────────────────────────────────────────┘
```

## Package Responsibilities

### Foundation Layer

**pkg/formats** - Color space representations and conversions
- RGB↔HSLA, LAB, XYZ color space conversions with full alpha support
- Hex string parsing and formatting (ToHex, ParseHex, ParseHexA for HEXA)
- HSLA type with complete alpha channel integration
- Pure functions with no external dependencies
- Optimized for repeated color space conversions

**pkg/chromatic** - Color theory algorithms and analysis
- Color harmony detection and scheme classification
- Contrast ratio calculations with WCAG compliance validation
- Hue analysis, chroma manipulation, and perceptual distance calculations
- Color derivation algorithms for palette generation
- Dependencies: pkg/formats

**pkg/settings** - Category-based configuration management
- Comprehensive CategorySettings with light/dark mode separation
- CategoryCharacteristics defining HSL constraints and contrast requirements
- CategoryScoringWeights for multi-dimensional scoring algorithms
- ExtractionSettings controlling processing behavior
- Fallback color configurations in hex string format
- Dependencies: Viper configuration library

**pkg/loader** - Image I/O with validation and optimization
- JPEG, PNG, and WebP image loading with format validation
- Memory-efficient processing with configurable size limits
- Image metadata extraction (dimensions, pixel count, format)
- Error handling and validation for all supported formats
- Dependencies: Standard library image packages

### Processing Layer

**pkg/processor** - Sophisticated category-based extraction system
- **27-category system**: Core UI, Terminal ANSI colors, Accents, Semantic colors
- **Multi-dimensional scoring**: Frequency, contrast, saturation, hue alignment, lightness
- **CategoryCandidates system**: Multiple scored options per category for palette generation
- **ColorProfile composition**: Comprehensive metadata with embedded ImageColors
- **Theme mode detection**: Light/dark pairing based on luminance analysis
- **Coverage quality metrics**: CoverageRatio indicating successful category extraction
- **WCAG compliance**: Automatic contrast validation with category-specific requirements
- **Configurable characteristics**: HSL constraints and scoring weights per category per mode
- Dependencies: pkg/formats, pkg/chromatic, pkg/settings, pkg/loader

## Category-Based Extraction System

### 27-Category Structure

The processor implements sophisticated categorization across four major groups:

#### Core UI Elements (4 categories)
- `background` - Theme background with mode-specific lightness constraints
- `foreground` - Text and UI elements with WCAG contrast requirements
- `dim_foreground` - Secondary text with reduced contrast requirements
- `cursor` - Cursor color with high contrast requirements

#### Terminal Colors (16 categories)
**Normal Colors (ANSI 0-7)**
- `normal_black`, `normal_red`, `normal_green`, `normal_yellow`
- `normal_blue`, `normal_magenta`, `normal_cyan`, `normal_white`

**Bright Colors (ANSI 8-15)**
- `bright_black`, `bright_red`, `bright_green`, `bright_yellow`
- `bright_blue`, `bright_magenta`, `bright_cyan`, `bright_white`

#### Accent Colors (3 categories)
- `accent_primary` - Primary theme accent with high saturation requirements
- `accent_secondary` - Supporting accent color
- `accent_tertiary` - Additional accent for complex themes

#### Semantic Colors (4 categories)
- `error` - Error states with red hue constraints
- `warning` - Warning states with yellow/orange hue constraints  
- `success` - Success states with green hue constraints
- `info` - Information states with blue hue constraints

### Multi-Dimensional Scoring Algorithm

Each color is evaluated against category requirements using weighted scoring:

```go
type CategoryScoringWeights struct {
    Frequency    float64 // Color frequency in image (0.25 default)
    Contrast     float64 // Contrast ratio against background (0.25 default)
    Saturation   float64 // Proximity to ideal saturation range (0.20 default)
    HueAlignment float64 // Alignment with category hue constraints (0.15 default)
    Lightness    float64 // Proximity to ideal lightness range (0.15 default)
}
```

### Category Characteristics Configuration

Each category defines HSL constraints and contrast requirements:

```go
type CategoryCharacteristics struct {
    MinLightness  float64  // Minimum lightness (0.0-1.0)
    MaxLightness  float64  // Maximum lightness (0.0-1.0)
    MinSaturation float64  // Minimum saturation (0.0-1.0)
    MaxSaturation float64  // Maximum saturation (0.0-1.0)
    MinContrast   float64  // Minimum contrast ratio vs background
    HueCenter     *float64 // Preferred hue in degrees (optional)
    HueTolerance  *float64 // Hue tolerance ± degrees (optional)
}
```

## Processing Pipeline

The processor implements a sophisticated single-pass pipeline with category-aware extraction:

```
Image Input
    │
    ▼
┌─────────────────────┐
│    Image Loading    │ ← pkg/loader (format validation)
│   Memory Efficient │   
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Frequency Extract  │ ← Single-pass pixel analysis
│    Color Counting   │   Configurable minimum thresholds
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Profile Analysis   │ ← Grayscale/monochromatic detection
│ Mode & Scheme Det.  │   Theme mode from luminance analysis
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│Background Selection │ ← Category-specific HSL constraints
│  Fallback Handling  │   Frequency + lightness scoring
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│Category Evaluation  │ ← 27 categories × N colors × 5 weights
│Multi-Dim Scoring    │   Parallel candidate evaluation
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│Candidate Selection  │ ← Top N candidates per category
│  Score Ranking      │   Configurable candidate limits
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│    Validation       │ ← WCAG contrast compliance
│   Coverage Ratio    │   Quality metrics calculation
└─────────┬───────────┘
          │
          ▼
    ColorProfile
  (Complete metadata)
```

## Data Structures

### ColorProfile Composition

The processor returns comprehensive analysis with embedded color data:

```go
type ColorProfile struct {
    Mode            ThemeMode           // Light or Dark theme pairing
    ColorScheme     ColorScheme         // Detected color scheme type
    IsGrayscale     bool               // Saturation-based classification (< 0.05)
    IsMonochromatic bool               // Hue variance analysis (±15° tolerance)
    DominantHue     float64            // Primary hue direction (0-360°)
    HueVariance     float64            // Color diversity metric
    AvgLuminance    float64            // Overall brightness (0.0-1.0)
    AvgSaturation   float64            // Overall color intensity (0.0-1.0)
    Colors          ImageColors        // Embedded category-based extraction
}
```

### Category-Based Color Organization

Colors are organized by sophisticated theme categories with rich metadata:

```go
type ImageColors struct {
    ColorFrequency     map[color.RGBA]uint32              // All colors with frequencies
    Categories         map[ColorCategory]color.RGBA       // Best color per category
    CategoryCandidates map[ColorCategory][]ColorCandidate // Multiple options per category
    TotalPixels        uint32                             // Image pixel count
    UniqueColors       int                                // Distinct colors found
    CoverageRatio      float64                            // % of categories filled (0.0-1.0)
}

type ColorCandidate struct {
    Color     color.RGBA `json:"color"`     // RGBA color value
    Frequency uint32     `json:"frequency"` // Pixel frequency in image
    Score     float64    `json:"score"`     // Multi-dimensional fitness score
}
```

## Architectural Patterns

### Settings-as-Methods Pattern

All public functions requiring configuration are methods on package configuration structures:

```go
// ✅ Correct: Method on configuration structure
func (p *Processor) GetCategoryCharacteristics(category ColorCategory, profile *ColorProfile) CategoryCharacteristics

// ✅ Correct: Private helper with settings from calling method  
func (p *Processor) calculateCategoryFitScore(c color.RGBA, category ColorCategory, ...) float64
```

### Category-First Design

The architecture prioritizes theme-oriented categorization over simple frequency analysis:

- **Theme-aware**: Each category has mode-specific characteristics (light/dark)
- **Configurable**: All constraints and scoring weights adjustable via settings
- **Extensible**: New categories easily added without breaking existing functionality
- **Quality metrics**: Coverage ratio provides extraction success indication

### Dependency Management

Clear dependency layers prevent circular dependencies:

```
Foundation Layer: Standard library + Viper (settings only)
Processing Layer: Foundation packages only (formats, chromatic, settings, loader)
Generation Layer: Foundation + Processing packages (future)
Application Layer: All packages (future)
```

## Performance Characteristics

The sophisticated category system maintains exceptional performance:

| Metric | Target | Achieved | Status |
|--------|--------|----------|--------|
| 4K Processing | <2s | 236ms avg | ✅ 88% faster than target |
| Memory Usage | <100MB | 8.6MB avg | ✅ 91% under limit |
| Peak Memory | <100MB | 61.2MB max | ✅ 39% under limit |
| Category Coverage | >80% | Variable by image | ✅ Quality-dependent |

### Algorithmic Complexity

- **Color extraction**: O(n) single-pass pixel analysis
- **Category scoring**: O(n×k×w) where n=colors, k=27 categories, w=5 weights  
- **Candidate selection**: O(n×log(n)) sorting per category
- **Total complexity**: O(n×k×w) dominated by scoring phase

### Memory Efficiency

- **Bounded candidates**: Configurable max candidates per category (default: 5)
- **Efficient storage**: ColorCandidate struct optimized for frequent access
- **Fallback handling**: Graceful degradation for insufficient color variety
- **No memory leaks**: All allocations properly scoped and released

## Quality Assurance

### Category Coverage Metrics

The system provides quality assessment through coverage analysis:

- **Coverage Ratio**: Percentage of 27 categories successfully filled
- **Candidate Depth**: Number of viable options per category
- **Score Distribution**: Quality assessment of category matches
- **Fallback Usage**: Tracking when default colors are applied

### WCAG Compliance

Automated accessibility validation ensures usable themes:

- **Contrast ratios**: Category-specific minimum requirements
- **AA compliance**: 4.5:1 minimum for foreground elements
- **AAA support**: 7:1 ratios for high-contrast elements
- **Automatic fallbacks**: Compliant colors when constraints cannot be met

## Design Principles

1. **Category-First Architecture**: Theme-oriented organization over frequency-based selection
2. **Multi-Dimensional Scoring**: Comprehensive evaluation beyond simple color frequency
3. **Configuration-Driven**: All characteristics and weights adjustable via settings
4. **Performance Preservation**: Sophisticated analysis within <2s/100MB constraints
5. **Quality Metrics**: Coverage ratio and candidate depth for extraction assessment
6. **Mode Awareness**: Light/dark theme considerations throughout the pipeline
7. **Standards Compliance**: WCAG accessibility requirements integrated at category level
8. **Extensible Design**: New categories and scoring dimensions easily integrated

This architecture successfully delivers sophisticated theme generation capabilities while maintaining exceptional performance and clean separation of concerns.