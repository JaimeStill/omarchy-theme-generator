# Architecture Documentation

## System Architecture

The omarchy-theme-generator uses a layered architecture with clear dependencies and separation of concerns, implementing a characteristic-based color extraction system with frequency weighting and flexible color pool organization.

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

**pkg/settings** - Flat configuration management
- Simplified settings structure with ~20 core parameters
- Threshold-based configuration (grayscale, monochromatic, lightness, saturation)
- Extraction settings controlling color pool behavior
- Hue organization parameters (sector count, sector size)
- Fallback color configurations for edge cases
- Dependencies: Viper configuration library

**pkg/loader** - Image I/O with validation and optimization
- JPEG, PNG, and WebP image loading with format validation
- Memory-efficient processing with configurable size limits
- Image metadata extraction (dimensions, pixel count, format)
- Error handling and validation for all supported formats
- Dependencies: Standard library image packages

### Processing Layer

**pkg/processor** - Characteristic-based color extraction system
- **ColorPool organization**: Lightness, saturation, and hue-based grouping
- **Frequency weighting**: Colors weighted by pixel frequency and perceptual importance
- **Statistical analysis**: Chromatic diversity, contrast range, hue variance calculations
- **ColorProfile composition**: Comprehensive metadata with embedded ColorPool
- **Theme mode detection**: Light/dark classification based on luminance analysis
- **Flexible extraction**: Supports 2-30+ color requirements for diverse themes
- **Profile detection**: Grayscale, monochromatic, and color scheme identification
- **Performance optimized**: Concurrent processing with <2s/100MB targets
- Dependencies: pkg/formats, pkg/chromatic, pkg/settings, pkg/loader

## Characteristic-Based Extraction System

### Three-Dimensional Color Organization

The processor organizes colors by natural characteristics rather than semantic categories:

#### Lightness Groups
- **Dark** (L < 0.25): Suitable for dark theme backgrounds and deep accents
- **Mid** (0.25 ≤ L < 0.75): Primary colors for foregrounds and main accents  
- **Light** (L ≥ 0.75): Light theme backgrounds and highlight colors

#### Saturation Groups
- **Gray** (S < 0.05): Neutral colors for backgrounds and subtle elements
- **Muted** (0.05 ≤ S < 0.25): Subdued colors for secondary elements
- **Normal** (0.25 ≤ S < 0.70): Standard saturation for most theme elements
- **Vibrant** (S ≥ 0.70): High-impact colors for accents and highlights

#### Hue Families  
- **Sectored organization**: Configurable hue sectors (default: 12 sectors × 30°)
- **Natural clustering**: Colors grouped by hue proximity
- **Relationship preservation**: Maintains harmony and contrast relationships
- **Flexible mapping**: Supports various color scheme strategies

### Frequency-Weighted Analysis

Colors are evaluated based on perceptual importance rather than arbitrary scoring:

```go
type WeightedColor struct {
    color.RGBA           // Embedded RGBA for direct access
    Frequency  uint32    // Pixel count in source image
    Weight     float64   // Normalized importance (frequency/total)
}
```

### Statistical Metrics

The system computes comprehensive color statistics for theme analysis:

```go
type ColorStatistics struct {
    HueHistogram       []float64           // Distribution across hue sectors
    LightnessHistogram []float64           // Distribution across lightness ranges
    SaturationGroups   map[string]float64  // Ratios for each saturation group
    
    PrimaryHue         float64             // Most dominant hue direction
    SecondaryHue       float64             // Second most dominant hue
    TertiaryHue        float64             // Third most dominant hue
    ChromaticDiversity float64             // Color entropy (0-1)
    ContrastRange      float64             // Luminance range (0-1)
    
    HueVariance        float64             // Hue spread measurement
    LightnessSpread    float64             // Balance across lightness groups
    SaturationSpread   float64             // Distribution across saturation groups
}
```

## Processing Pipeline

The processor implements a streamlined three-stage pipeline with characteristic-based organization:

```
Image Input
    │
    ▼
┌─────────────────────┐
│    Image Loading    │ ← pkg/loader (format validation)
│   Memory Efficient │   Size limits and optimization
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Color Extraction   │ ← Concurrent frequency analysis
│  Frequency Weighting│   Minimum threshold filtering
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│Characteristic Group │ ← Lightness/Saturation/Hue grouping
│  ColorPool Build    │   Relationship preservation
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ Statistical Analysis│ ← Diversity, contrast, variance metrics
│  Profile Detection  │   Grayscale/monochromatic detection
└─────────┬───────────┘
          │
          ▼
    ColorProfile
 (Pool + Statistics)
```

## Data Structures

### ColorProfile Composition

The processor returns comprehensive analysis with embedded ColorPool:

```go
type ColorProfile struct {
    Mode            ThemeMode             // Light or Dark theme classification
    IsGrayscale     bool                 // Saturation-based classification (< 0.05)
    IsMonochromatic bool                 // Hue variance analysis (±15° tolerance)
    DominantHue     float64              // Primary hue direction (0-360°)
    HueVariance     float64              // Color spread measurement
    AvgLuminance    float64              // Overall brightness (0.0-1.0)
    AvgSaturation   float64              // Overall color intensity (0.0-1.0)
    Pool            ColorPool            // Characteristic-based organization
}
```

### ColorPool Organization

Colors are organized by natural characteristics with comprehensive statistics:

```go
type ColorPool struct {
    AllColors      []WeightedColor    // All extracted colors sorted by weight
    DominantColors []WeightedColor    // Top dominant colors by frequency
    
    ByLightness  LightnessGroups     // Dark, Mid, Light groupings
    BySaturation SaturationGroups    // Gray, Muted, Normal, Vibrant groupings
    ByHue        HueFamilies         // Hue sector organization
    
    TotalPixels  uint32              // Source image pixel count
    UniqueColors int                 // Distinct colors extracted
    
    Statistics   ColorStatistics     // Computed analysis metrics
}

type LightnessGroups struct {
    Dark  []WeightedColor    // L < 0.25
    Mid   []WeightedColor    // 0.25 ≤ L < 0.75  
    Light []WeightedColor    // L ≥ 0.75
}

type SaturationGroups struct {
    Gray    []WeightedColor  // S < 0.05
    Muted   []WeightedColor  // 0.05 ≤ S < 0.25
    Normal  []WeightedColor  // 0.25 ≤ S < 0.70
    Vibrant []WeightedColor  // S ≥ 0.70
}

type HueFamilies map[int][]WeightedColor  // Sector → Colors
```

## Architectural Patterns

### Settings-as-Methods Pattern

All public functions requiring configuration are methods on package configuration structures:

```go
// ✅ Correct: Method on configuration structure
func (p *Processor) ProcessImage(img image.Image) (*ColorProfile, error)

// ✅ Correct: Private helper with settings from calling method  
func (p *Processor) calculateStatistics(pool ColorPool) ColorStatistics
```

### Characteristic-First Design

The architecture prioritizes natural color organization over premature semantic assignment:

- **Natural grouping**: Colors grouped by perceptual characteristics (L, S, H)
- **Relationship preservation**: Maintains harmony and contrast relationships
- **Flexible mapping**: Supports various theme generation strategies
- **No premature semantics**: Avoids early assignment to specific UI roles
- **Statistical richness**: Comprehensive metrics for downstream processing

### Separation of Concerns

Clear distinction between extraction, semantic mapping, and theme generation:

- **pkg/processor**: Extracts and organizes colors by characteristics
- **pkg/palette (future)**: Maps colors to semantic roles based on requirements
- **pkg/theme (future)**: Generates component-specific configurations

### Dependency Management

Clear dependency layers prevent circular dependencies:

```
Foundation Layer: Standard library + Viper (settings only)
Processing Layer: Foundation packages only (formats, chromatic, settings, loader)
Generation Layer: Foundation + Processing packages (future)
Application Layer: All packages (future)
```

## Performance Characteristics

The characteristic-based system maintains exceptional performance:

| Metric | Target | Achieved | Status |
|--------|--------|----------|--------|
| 4K Processing | <2s | ~500ms avg | ✅ 75% faster than target |
| Memory Usage | <100MB | ~12MB avg | ✅ 88% under limit |
| Peak Memory | <100MB | ~45MB max | ✅ 55% under limit |
| Color Extraction | Variable | 2-100+ colors | ✅ Flexible by requirements |

### Algorithmic Complexity

- **Color extraction**: O(n) single-pass pixel analysis with concurrent processing
- **Characteristic grouping**: O(n) linear sorting into lightness/saturation/hue groups
- **Statistical analysis**: O(n) for most metrics, O(n log n) for dominant color ranking
- **Total complexity**: O(n log n) dominated by dominant color selection

### Memory Efficiency

- **Bounded extraction**: Configurable maximum colors extracted (default: 100)
- **WeightedColor optimization**: Embedded RGBA reduces pointer indirection
- **Concurrent processing**: Worker pools prevent excessive goroutine creation
- **Progressive filtering**: Early termination when minimum thresholds not met

## Quality Assurance

### Statistical Metrics

The system provides comprehensive quality assessment:

- **Chromatic Diversity**: Entropy-based color distribution measurement (0-1)
- **Contrast Range**: Luminance span across all extracted colors (0-1)  
- **Lightness Spread**: Balance across dark, mid, and light groupings
- **Saturation Spread**: Distribution across gray, muted, normal, vibrant groups
- **Hue Variance**: Angular spread of non-grayscale colors

### Profile Detection

Automatic classification for edge cases and special images:

- **Grayscale Detection**: Average saturation below configurable threshold
- **Monochromatic Detection**: Hue variance within tolerance range
- **Theme Mode Classification**: Based on weighted average luminance
- **Color Scheme Detection**: Integration with pkg/chromatic harmony analysis

## Design Principles

1. **Characteristic-First Architecture**: Natural color organization over premature semantic assignment
2. **Frequency-Weighted Analysis**: Perceptual importance drives color selection
3. **Statistical Richness**: Comprehensive metrics enable sophisticated downstream processing  
4. **Performance Preservation**: Streamlined pipeline within <2s/100MB constraints
5. **Flexible Extraction**: Supports diverse color requirements (2-30+ colors)
6. **Separation of Concerns**: Clear distinction between extraction and semantic mapping
7. **Relationship Preservation**: Maintains color harmony and contrast relationships
8. **Future-Proof Design**: Characteristic organization supports various theme strategies

This architecture successfully delivers flexible color extraction capabilities while maintaining exceptional performance and enabling sophisticated theme generation through downstream semantic mapping.