# PROJECT.md Restructure

## Current Implementation

### Core Infrastructure
**pkg/formats/** - Data formatting and color analysis
- `color.go` - HSL conversion, contrast ratio, hex formatting
- `types.go` - Theme-specific type definitions (ColorRole, ThemeMode)

**pkg/extractor/** - Image processing and color extraction
- `loader.go` - Image loading with format validation
- `strategy_frequency.go` - Frequency-based color extraction
- `strategy_saliency.go` - Saliency-based extraction for complex images
- `strategies.go` - Strategy selection based on image characteristics
- `image_analysis.go` - Edge detection, color complexity analysis
- `settings.go` - Configurable thresholds (empirically derived)
- `extractor.go` - Main extraction pipeline

**pkg/errors/** - Centralized error handling
- `extractor.go` - Image processing errors
- `docs.go` - Package documentation

**tests/** - Validation and analysis
- `strategies_test.go` - Multi-strategy validation with real images
- `analyze-images/` - Image characteristic analysis tool
- `images/` - 15 test wallpapers with documented characteristics
- `internal/` - Shared test utilities

### Capabilities
- ✅ Multi-strategy extraction (frequency vs saliency)
- ✅ Automatic strategy selection based on image analysis
- ✅ Grayscale vs monochromatic detection
- ✅ Performance: <2s for 4K images
- ✅ Settings-driven configuration

---

## Current Work

### Refactoring: pkg/color → pkg/formats
**Goal**: Simplify to only what's used, leverage standard library

**Tasks**:
- [ ] Remove unused color space conversions (LAB, distance calculations)
- [ ] Replace custom Color type with standard `color.RGBA`
- [ ] Extract only needed functions: `RGBToHSL()`, `ContrastRatio()`, `ToHex()`
- [ ] Update extractor to use `color.RGBA` directly
- [ ] Update FrequencyMap to use `map[color.RGBA]uint32`

---

## Components & Features (Ordered by Dependency)

### Layer 1: Foundation

#### **pkg/formats** - Data formatting and conversion
*Purpose*: Handle all data transformations between different formats

**Features**:
- `RGBToHSL()` - Convert colors to HSL for analysis
- `ContrastRatio()` - WCAG accessibility calculations
- `ToHex()` - Color to hex string conversion
- `ParseHex()` - Hex string to color parsing
- Theme type definitions (ColorRole, ThemeMode, etc.)

*Dependencies*: Standard library only

---

### Layer 2: Analysis

#### **pkg/analysis** - Image and color analysis
*Purpose*: Analyze images to determine extraction strategy and theme mode

**Features**:
- Image profile detection (grayscale, monochromatic, full-color)
- Theme mode detection (light/dark based on luminance)
- Color clustering (group similar colors by perceptual distance)
- Edge detection for complexity analysis
- Saturation and temperature analysis

*Dependencies*: pkg/formats

---

### Layer 3: Extraction

#### **pkg/extractor** - Purpose-driven color extraction
*Purpose*: Extract colors organized by their role in the theme

**Features**:
- Role-based extraction (backgrounds, foregrounds, accents)
- Mode-aware role assignment
- Profile-specific strategies (grayscale → synthesis, etc.)
- Settings-driven configuration
- User preferences support

*Dependencies*: pkg/formats, pkg/analysis

---

### Layer 4: Generation

#### **pkg/palette** - Color scheme generation
*Purpose*: Apply color theory to generate complete palettes

**Features**:
- Color theory schemes (complementary, triadic, etc.)
- Synthesis for minimal-color images
- WCAG compliance validation
- Scheme application to role-based colors

*Dependencies*: pkg/formats, pkg/analysis

#### **pkg/theme** - Theme configuration generation
*Purpose*: Generate theme configuration files for all supported formats

**Features**:
- Template-based generation for each format
- Role → configuration mapping
- Format-specific color conversion
- Metadata generation (theme-gen.json)

*Dependencies*: pkg/formats, pkg/palette

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

---

## Implementation Notes

### Simplification Principles
1. **Use standard library types** where possible (`color.RGBA`)
2. **Only build what's needed** - no speculative features
3. **Settings over hardcoding** - all thresholds configurable
4. **Clear dependency layers** - each package has specific purpose

### File Structure
```
omarchy-theme-generator/
├── pkg/
│   ├── formats/        # Data formatting (was pkg/color)
│   │   ├── color.go    # Color conversions and formatting
│   │   └── types.go    # Theme types and constants
│   ├── analysis/       # Image and color analysis (NEW)
│   ├── extractor/      # Color extraction (REFACTOR)
│   ├── palette/        # Color schemes (NEW)
│   └── theme/          # Theme generation (NEW)
├── cmd/
│   └── omarchy-theme-gen/
└── tests/
```

### Next Steps
1. Complete pkg/formats refactor
2. Build pkg/analysis for profile detection
3. Enhance extractor with role-based extraction
4. Implement palette generation
5. Create theme generators

---

## Completed Features

### Session 1-3 (Foundation)
- ✅ Project structure
- ✅ Basic color types with HSL conversion
- ✅ Image loading and validation
- ✅ Frequency-based extraction
- ✅ Performance optimization

### Session 4 (Multi-Strategy System)
- ✅ Saliency strategy implementation
- ✅ Strategy selection based on image analysis
- ✅ Edge detection and complexity analysis
- ✅ Settings-driven configuration
- ✅ Comprehensive testing with real images

---

# pkg/formats Package Design

## Overview
Simplified package for data formatting, replacing the over-engineered pkg/color.

## Core Files

### formats/color.go
```go
package formats

import (
    "fmt"
    "image/color"
    "math"
)

// RGBToHSL converts any color to HSL values (0-1 range)
func RGBToHSL(c color.Color) (h, s, l float64) {
    r, g, b, _ := c.RGBA()
    // Note: RGBA() returns uint32 in range [0, 0xffff]
    rf := float64(r) / 0xffff
    gf := float64(g) / 0xffff
    bf := float64(b) / 0xffff
    
    max := math.Max(math.Max(rf, gf), bf)
    min := math.Min(math.Min(rf, gf), bf)
    l = (max + min) / 2
    
    if max == min {
        h, s = 0, 0
    } else {
        d := max - min
        if l > 0.5 {
            s = d / (2 - max - min)
        } else {
            s = d / (max + min)
        }
        
        switch max {
        case rf:
            h = (gf - bf) / d
            if gf < bf {
                h += 6
            }
        case gf:
            h = (bf-rf)/d + 2
        case bf:
            h = (rf-gf)/d + 4
        }
        h /= 6
    }
    
    return h, s, l
}

// ContrastRatio calculates WCAG contrast ratio between two colors
func ContrastRatio(c1, c2 color.Color) float64 {
    l1 := relativeLuminance(c1)
    l2 := relativeLuminance(c2)
    
    if l1 < l2 {
        l1, l2 = l2, l1
    }
    
    return (l1 + 0.05) / (l2 + 0.05)
}

// ToHex converts a color to hex string format (#RRGGBB)
func ToHex(c color.Color) string {
    rgba := color.RGBAModel.Convert(c).(color.RGBA)
    return fmt.Sprintf("#%02x%02x%02x", rgba.R, rgba.G, rgba.B)
}

// ToHexA converts a color to hex string with alpha (#RRGGBBAA)
func ToHexA(c color.Color) string {
    rgba := color.RGBAModel.Convert(c).(color.RGBA)
    return fmt.Sprintf("#%02x%02x%02x%02x", rgba.R, rgba.G, rgba.B, rgba.A)
}

// ParseHex parses a hex color string to color.RGBA
func ParseHex(hex string) (color.RGBA, error) {
    // Implementation
}

// Helper functions
func relativeLuminance(c color.Color) float64 {
    r, g, b, _ := c.RGBA()
    rf := toLinearRGB(float64(r) / 0xffff)
    gf := toLinearRGB(float64(g) / 0xffff)
    bf := toLinearRGB(float64(b) / 0xffff)
    
    return 0.2126*rf + 0.7152*gf + 0.0722*bf
}

func toLinearRGB(channel float64) float64 {
    if channel <= 0.03928 {
        return channel / 12.92
    }
    return math.Pow((channel+0.055)/1.055, 2.4)
}
```

### formats/types.go
```go
package formats

// ColorRole represents the purpose of a color in the theme
type ColorRole string

const (
    // Core UI roles
    RoleBackground    ColorRole = "background"
    RoleForeground    ColorRole = "foreground"
    RolePrimary       ColorRole = "primary"
    RoleSecondary     ColorRole = "secondary"
    
    // Terminal colors
    RoleTerminalBlack ColorRole = "terminal-black"
    RoleTerminalRed   ColorRole = "terminal-red"
    // ... etc
)

// ThemeMode represents light or dark theme preference
type ThemeMode string

const (
    ModeDark  ThemeMode = "dark"
    ModeLight ThemeMode = "light"
    ModeAuto  ThemeMode = "auto"
)

// ImageProfile represents the color characteristics of an image
type ImageProfile int

const (
    ProfileFullColor ImageProfile = iota
    ProfileGrayscale
    ProfileMonochromatic
    ProfileDuotone
)
```

## Migration Examples

### Before (pkg/color)
```go
c := color.NewRGB(255, 128, 64)
hex := c.HEX()
h, s, l := c.HSL()
contrast := c.ContrastRatio(other)
```

### After (pkg/formats)
```go
c := color.RGBA{255, 128, 64, 255}
hex := formats.ToHex(c)
h, s, l := formats.RGBToHSL(c)
contrast := formats.ContrastRatio(c, other)
```

## Benefits
1. **90% less code** - Only what we actually use
2. **Standard types** - Better interoperability
3. **Clear purpose** - Formatting and conversion only
4. **No premature optimization** - No caching, no unused features
5. **Extensible** - Easy to add formats as needed
