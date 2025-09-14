# Architecture Documentation

## System Architecture

The omarchy-theme-generator uses a layered architecture with clear dependencies and separation of concerns, implementing a characteristic-based color extraction system with frequency weighting and cluster-based organization.

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
- QuantizeColor function for color clustering and similarity detection
- Pure functions with no external dependencies
- Optimized for repeated color space conversions

**pkg/chromatic** - Color theory algorithms and analysis
- Color harmony detection and scheme classification
- Contrast ratio calculations with WCAG compliance validation
- Hue analysis, chroma manipulation, and perceptual distance calculations
- Color derivation algorithms for palette generation
- Color similarity detection using Delta-E and perceptual distance
- Dependencies: pkg/formats

**pkg/settings** - Hierarchical configuration management
- Settings structure with extraction, clustering, and theme parameters
- Threshold-based configuration (grayscale detection, clustering tolerances)
- Extraction settings controlling color frequency and filtering
- UI-specific thresholds (lightness, saturation, vibrancy boundaries)
- Fallback color configurations for edge cases
- Dependencies: Viper configuration library

**pkg/loader** - Image I/O with validation and optimization
- JPEG, PNG, and WebP image loading with format validation
- Memory-efficient processing with configurable size limits
- Image metadata extraction (dimensions, pixel count, format)
- Error handling and validation for all supported formats
- Dependencies: Standard library image packages

### Processing Layer

**pkg/processor** - Characteristic-based color extraction and clustering
- **ColorCluster organization**: Representative colors with pre-calculated characteristics
- **Frequency weighting**: Colors weighted by pixel frequency and perceptual importance
- **Clustering algorithm**: Groups similar colors into distinct visual clusters
- **UI-optimized filtering**: Removes colors unsuitable for theme generation
- **Theme mode detection**: Light/dark classification based on luminance analysis
- **Grayscale detection**: Identifies images with insufficient color saturation
- **Performance optimized**: Concurrent processing with <2s/100MB targets
- Dependencies: pkg/formats, pkg/chromatic, pkg/settings, pkg/loader

## ColorCluster-Based Extraction System

### Color Clustering Approach

The processor extracts colors through a clustering algorithm that groups perceptually similar colors:

#### ColorCluster Structure
```go
type ColorCluster struct {
    color.RGBA                   // The representative color
    Weight      float64          // Combined weight (0.0-1.0)
    Lightness   float64          // Pre-calculated HSL lightness for efficiency
    Saturation  float64          // Pre-calculated HSL saturation for efficiency
    Hue         float64          // Hue in degrees (0-360)
    IsNeutral   bool            // Grayscale or very low saturation
    IsDark      bool            // L < 0.3
    IsLight     bool            // L > 0.7
    IsMuted     bool            // S < 0.3
    IsVibrant   bool            // S > 0.7
}
```

#### Characteristic Flags
- **IsNeutral**: Colors with minimal saturation suitable for backgrounds and text
- **IsDark/IsLight**: Lightness-based classification for theme mode compatibility
- **IsMuted/IsVibrant**: Saturation-based classification for accent vs primary usage

### Frequency-Weighted Analysis

Colors are evaluated based on perceptual importance and visual weight:

```go
type WeightedColor struct {
    color.RGBA           // Embedded RGBA for direct access
    Frequency  uint32    // Pixel count in source image
    Weight     float64   // Normalized importance (frequency/total)
}
```

### ColorProfile Output

The system returns a minimal, focused structure for theme generation:

```go
type ColorProfile struct {
    Mode       ThemeMode      // Light or Dark theme base
    Colors     []ColorCluster // Distinct colors, sorted by weight
    HasColor   bool          // False if image is essentially grayscale
    ColorCount int           // Number of distinct colors found
}
```

## Processing Pipeline

The processor implements a streamlined clustering-based pipeline:

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
│  Color Extraction   │ ← Frequency analysis with sampling
│  Minimum Threshold │   Filter colors below frequency threshold
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Color Clustering   │ ← Group similar colors into clusters
│  Similarity-Based  │   Use Delta-E and perceptual distance
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  UI Filtering       │ ← Remove unsuitable colors
│  Characteristic     │   Filter by lightness/saturation criteria
│  Flag Calculation   │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ Theme Mode & Stats  │ ← Light/dark mode determination
│ Profile Generation  │   Grayscale detection and final sorting
└─────────┬───────────┘
          │
          ▼
    ColorProfile
   (Clusters Only)
```

## Architectural Patterns

### Settings-as-Methods Pattern

All public functions requiring configuration are methods on package configuration structures:

```go
// ✅ Correct: Method on configuration structure
func (p *Processor) ProcessImage(img image.Image) (*ColorProfile, error)

// ✅ Correct: Private helper with settings from calling method
func (p *Processor) clusterColors(weighted []WeightedColor) []ColorCluster
```

### Characteristic-First Design

The architecture prioritizes natural color characteristics over premature semantic assignment:

- **Representative clustering**: Groups similar colors into distinct visual clusters
- **Pre-calculated characteristics**: Lightness, saturation, hue computed once
- **Boolean flags**: Efficient UI-relevant categorization (neutral, dark, light, muted, vibrant)
- **No semantic roles**: Avoids premature assignment to specific UI components
- **Weight-based sorting**: Colors ordered by visual importance

### Separation of Concerns

Clear distinction between extraction, semantic mapping, and theme generation:

- **pkg/processor**: Extracts and clusters colors by visual characteristics
- **pkg/palette (future)**: Maps ColorClusters to semantic roles based on requirements
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

The ColorCluster-based system maintains exceptional performance:

| Metric | Target | Achieved | Status |
|--------|--------|----------|--------|
| 4K Processing | <2s | ~500ms avg | ✅ 75% faster than target |
| Memory Usage | <100MB | ~12MB avg | ✅ 88% under limit |
| Peak Memory | <100MB | ~45MB max | ✅ 55% under limit |
| Color Extraction | Variable | 2-50 clusters | ✅ Suitable for theme generation |

### Algorithmic Complexity

- **Color extraction**: O(n) single-pass pixel analysis with sampling optimization
- **Color clustering**: O(n²) pairwise distance calculations with early termination
- **UI filtering**: O(n) linear pass through clusters
- **Characteristic calculation**: O(n) pre-computation of HSL values and flags
- **Total complexity**: O(n²) dominated by clustering algorithm

### Memory Efficiency

- **Bounded clustering**: Practical limit of ~50 clusters prevents excessive memory usage
- **Pre-calculated characteristics**: Avoids repeated HSL conversions
- **Embedded RGBA**: ColorCluster embeds color.RGBA to reduce pointer indirection
- **Sampling optimization**: Large images processed with pixel sampling

## Quality Assurance

### Theme Suitability

The system ensures extracted colors are appropriate for UI themes:

- **Lightness filtering**: Removes colors too similar in lightness for contrast
- **Saturation thresholds**: Maintains distinction between neutral and colorful elements
- **Weight-based ranking**: Prioritizes visually important colors from the source image
- **Grayscale detection**: Identifies images unsuitable for colorful themes

### Profile Detection

Automatic classification for edge cases and special images:

- **Grayscale Detection**: HasColor flag based on overall saturation analysis
- **Theme Mode Classification**: Light/Dark based on weighted average lightness
- **Color Count Tracking**: Number of distinct clusters for downstream processing

## Design Principles

1. **Clustering-Based Architecture**: Group similar colors into distinct visual clusters
2. **Frequency-Weighted Selection**: Visual importance drives color prioritization
3. **Pre-calculated Characteristics**: Efficient access to lightness, saturation, hue properties
4. **UI-Optimized Filtering**: Remove colors unsuitable for theme generation
5. **Minimal Output Structure**: Focus on essential data for downstream processing
6. **Performance Preservation**: Maintain <2s/100MB constraints through optimized algorithms
7. **Clear Separation**: Distinct phases for extraction, clustering, and filtering
8. **Future-Proof Design**: ColorCluster structure supports various theme generation strategies

This architecture successfully delivers focused color extraction capabilities through visual clustering while maintaining exceptional performance and providing a clean interface for downstream theme generation processes.