# Architecture Documentation

## System Architecture

The omarchy-theme-generator uses a layered architecture with clear dependencies and separation of concerns. Each layer depends only on layers below it, preventing circular dependencies and ensuring maintainable code.

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
- RGB↔HSL, LAB, XYZ color space conversions
- Hex string parsing and formatting (ToHex, ParseHex)
- HSLA type with full alpha channel support 
- WCAG accessibility calculations (ContrastRatio)
- Pure functions with no external dependencies

**pkg/chromatic** - Color theory foundation
- Color harmony detection and scheme classification
- Contrast ratio and perceptual distance calculations
- Hue analysis and chroma manipulation utilities
- Color derivation algorithms for theme generation
- Dependencies: pkg/formats

**pkg/settings** - System configuration management
- Flat settings structure with Viper integration
- Empirical thresholds and performance parameters
- Settings-as-methods architectural pattern enforcement
- Fallback color configurations (hex string format)
- Dependencies: Standard library + Viper

**pkg/loader** - Image I/O and validation
- JPEG and PNG image loading with format validation
- Memory-efficient image processing and error handling
- Image metadata extraction (dimensions, pixel count)
- Format support validation and conversion
- Dependencies: Standard library image packages

### Processing Layer

**pkg/processor** - Unified image processing and analysis
- Single-pass processing pipeline: Load → Extract → Analyze → Assign roles
- ColorProfile composition with comprehensive metadata including embedded ImageColors
- Frequency-based color extraction optimized for all image types
- Role-based color assignment to background/foreground/primary/secondary/accent
- Integrated analysis: Grayscale, monochromatic, and color scheme detection
- Theme mode detection based on luminance analysis (light/dark pairing)
- WCAG compliance validation with automatic fallback handling
- Dependencies: pkg/formats, pkg/chromatic, pkg/settings, pkg/loader

## Processing Pipeline

The processor implements a single-pass processing pipeline:

```
Image Input
    │
    ▼
┌─────────────────────┐
│    Image Loading    │ ← pkg/loader
│   (Format Check)    │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Color Extraction   │ ← Frequency-based analysis
│   (Single-pass)     │   
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Profile Analysis   │ ← Grayscale/monochromatic detection
│ (Theme Mode, etc.)  │   Color scheme classification
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Role Assignment    │ ← Background/foreground selection
│  (Purpose-driven)   │   Primary/secondary/accent mapping
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│    Validation       │ ← WCAG contrast compliance
│   (WCAG + Fallback) │   Fallback color application
└─────────┬───────────┘
          │
          ▼
    ColorProfile
   (Complete metadata)
```

## Data Structures

### ColorProfile Composition

The processor returns a comprehensive ColorProfile that embeds ImageColors and includes all analysis metadata:

```go
type ColorProfile struct {
    Mode            ThemeMode           // Light or Dark
    ColorScheme     ColorScheme         // Detected scheme type
    IsGrayscale     bool               // Saturation analysis
    IsMonochromatic bool               // Hue variance analysis  
    DominantHue     float64            // Primary hue direction
    HueVariance     float64            // Color diversity metric
    AvgLuminance    float64            // Overall brightness
    AvgSaturation   float64            // Overall color intensity
    Colors          ImageColors        // Embedded role-based colors
}
```

### Role-Based Color Organization

Colors are organized by purpose rather than frequency:

```go
type ImageColors struct {
    Background   color.RGBA  // Theme background color
    Foreground   color.RGBA  // Text and UI elements
    Primary      color.RGBA  // Brand/accent color
    Secondary    color.RGBA  // Supporting elements  
    Accent       color.RGBA  // Highlights and emphasis
    MostFrequent color.RGBA  // Statistical reference
}
```

## Architectural Patterns

### Settings-as-Methods Pattern

All public functions requiring configuration are methods on package configuration structures:

```go
// Correct: Method on configuration structure
func (p *Processor) ProcessImage(img image.Image) (*ColorProfile, error)

// Incorrect: Function with settings parameter  
func ProcessImage(img image.Image, settings *Settings) (*ColorProfile, error)
```

### Dependency Management

Clear dependency layers prevent circular dependencies:

```
Foundation Layer: No external dependencies (except Viper for settings)
Processing Layer: Foundation packages only
Generation Layer: Foundation + Processing packages
Application Layer: All packages
```

## Performance Characteristics

The architecture achieves the following performance metrics:

| Metric | Target | Achieved | Status |
|--------|--------|----------|--------|
| 4K Processing | <2s | 236ms avg | ✅ 88% faster |
| Memory Usage | <100MB | 8.6MB avg | ✅ 91% under limit |
| Peak Memory | <100MB | 61.2MB max | ✅ 39% under limit |
| Target Compliance | 100% | 100% (15/15) | ✅ Perfect |

### Performance by Image Size

- **Medium (2-8MP)**: 12 images, 147ms average
- **Large (>8MP)**: 3 images, 593ms average  
- **All images**: Sub-second processing with full analysis

## Design Principles

1. **Single Responsibility**: Each package has one clear purpose
2. **Performance First**: All decisions favor speed and memory efficiency
3. **Standards Compliance**: WCAG accessibility and Go best practices
4. **Simplicity**: Eliminate abstraction layers that don't add value
5. **Testability**: Comprehensive unit test coverage with diagnostic logging