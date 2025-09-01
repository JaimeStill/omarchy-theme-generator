# Technical Specification

## Overview

Command-line tool in Go that generates Omarchy themes from images using intelligent, purpose-driven color extraction and color theory principles. The system uses a layered architecture with standard library types for optimal performance and maintainability.

## Input Parameters

**Required:**
- Source image file (JPEG, PNG)

**Optional Overrides:**
- Theme mode: `light` | `dark` | `auto`
- Primary color: Base for palette generation (hex format)
- Background color: Main background (hex format)
- Foreground color: Main text color (hex format)

When not provided, colors are intelligently derived from the source image.

## Output Structure

```
theme-name/
├── alacritty.toml       # Terminal emulator
├── btop.theme           # System monitor
├── hyprland.conf        # Window manager
├── hyprlock.conf        # Lock screen
├── mako.ini             # Notifications
├── neovim.lua           # Editor colorscheme
├── waybar.css           # Status bar
├── walker.css           # App launcher
├── swayosd.css          # On-screen display
├── backgrounds/         # Contains source image
│   └── wallpaper.jpg
└── light.mode           # Present only for light themes
```

## Core Technical Concepts

### 1. Purpose-Driven Color Extraction

The system organizes colors by their intended role in the theme rather than just their frequency:

**Color Roles:**
- **Background**: Colors suitable for window/terminal backgrounds based on mode and lightness
- **Foreground**: Colors suitable for text with proper contrast ratios
- **Accents**: Saturated colors for highlights and UI elements
- **Terminal Colors**: ANSI color palette mapping

**Mode-Aware Processing:**
```go
// Role assignment adapts based on detected theme mode
func AssignColorRoles(colors []color.RGBA, mode ThemeMode) map[ColorRole][]color.RGBA {
    roles := make(map[ColorRole][]color.RGBA)
    
    for _, c := range colors {
        _, _, l := formats.RGBToHSL(c)
        
        // Dark mode: darker colors → backgrounds
        if mode == ModeDark && l < 0.35 {
            roles[RoleBackground] = append(roles[RoleBackground], c)
        }
        
        // Light mode: lighter colors → backgrounds  
        if mode == ModeLight && l > 0.75 {
            roles[RoleBackground] = append(roles[RoleBackground], c)
        }
    }
    
    return roles
}
```

### 2. Standard Library Color Types

Uses `color.RGBA` from Go standard library instead of custom types:

```go
// pkg/formats/color.go
import "image/color"

// RGBToHSL converts standard color to HSL values
func RGBToHSL(c color.Color) (h, s, l float64) {
    r, g, b, _ := c.RGBA()
    // Convert uint32 [0, 0xffff] to float64 [0, 1]
    rf := float64(r) / 0xffff
    gf := float64(g) / 0xffff
    bf := float64(b) / 0xffff
    
    max := math.Max(math.Max(rf, gf), bf)
    min := math.Min(math.Min(rf, gf), bf)
    l = (max + min) / 2
    
    // Standard HSL conversion logic...
    return h, s, l
}

// ToHex converts color to hex string format (#RRGGBB)
func ToHex(c color.Color) string {
    rgba := color.RGBAModel.Convert(c).(color.RGBA)
    return fmt.Sprintf("#%02x%02x%02x", rgba.R, rgba.G, rgba.B)
}
```

### 2. Color Space Conversion
Bidirectional RGB ↔ HSL conversion for intuitive manipulation. HSL provides better lightness control for theme generation than HSV.

```go
func RGBToHSL(r, g, b uint8) (h, s, l float64) {
    rf, gf, bf := float64(r)/255.0, float64(g)/255.0, float64(b)/255.0
    max := math.Max(math.Max(rf, gf), bf)
    min := math.Min(math.Min(rf, gf), bf)
    l = (max + min) / 2
    
    if max == min {
        h, s = 0, 0 // achromatic
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
            if gf < bf { h += 6 }
        case gf:
            h = (bf - rf) / d + 2
        case bf:
            h = (rf - gf) / d + 4
        }
        h /= 6
    }
    return
}
```

### 3. Octree Color Quantization
Tree-based color reduction with O(n) insertion time and deterministic results. Maximum depth 8 allows up to 256 leaf nodes (color palette entries).

```go
type OctreeNode struct {
    red, green, blue uint32
    pixelCount       uint32
    paletteIndex     int
    children         [8]*OctreeNode
}

func (node *OctreeNode) Insert(r, g, b uint8, level int) {
    if level >= maxDepth {
        node.red += uint32(r)
        node.green += uint32(g)
        node.blue += uint32(b)
        node.pixelCount++
        return
    }
    
    // Determine octant based on RGB bit values
    index := 0
    mask := uint8(0x80 >> level)
    if r&mask != 0 { index |= 4 }
    if g&mask != 0 { index |= 2 }
    if b&mask != 0 { index |= 1 }
    
    if node.children[index] == nil {
        node.children[index] = &OctreeNode{}
    }
    node.children[index].Insert(r, g, b, level+1)
}
```

### 4. Color Theory Palette Generation
Applies established harmony principles to generate aesthetically pleasing palettes from a base color.

```go
type PaletteStrategy interface {
    Generate(baseColor Color, count int) []Color
    Name() string
}

// Monochromatic: Single hue with varied saturation/lightness
type MonochromaticStrategy struct{}

func (m MonochromaticStrategy) Generate(base Color, count int) []Color {
    h, s, l := base.HSL()
    palette := make([]Color, count)
    
    for i := 0; i < count; i++ {
        factor := float64(i) / float64(count-1)
        newL := l + (factor-0.5)*0.4   // ±40% lightness
        newS := s * (0.5 + factor*0.5) // 50-100% saturation
        palette[i] = NewHSL(h, clamp(newS, 0, 1), clamp(newL, 0, 1))
    }
    return palette
}

// Triadic: Colors 120° apart on color wheel
type TriadicStrategy struct{}

func (t TriadicStrategy) Generate(base Color, count int) []Color {
    h, s, l := base.HSL()
    palette := []Color{base}
    perColor := (count - 1) / 2
    
    // First triad at +120°
    for i := 0; i < perColor; i++ {
        factor := float64(i) / float64(perColor)
        newL := l + (factor-0.5)*0.3
        palette = append(palette, 
            NewHSL(math.Mod(h+1.0/3.0, 1), s, clamp(newL, 0, 1)))
    }
    
    // Second triad at +240°
    for i := 0; i < count-1-perColor; i++ {
        factor := float64(i) / float64(count-1-perColor)
        newL := l + (factor-0.5)*0.3
        palette = append(palette,
            NewHSL(math.Mod(h+2.0/3.0, 1), s, clamp(newL, 0, 1)))
    }
    
    return palette
}
```

### 5. Image Processing with User Overrides
Extracts colors while respecting optional user preferences for theme mode, primary, background, and foreground colors.

```go
type ThemeConfig struct {
    SourceImage     image.Image
    Mode            ThemeMode     // Light, Dark, or Auto
    PrimaryColor    *Color        // User override or nil
    BackgroundColor *Color        // User override or nil
    ForegroundColor *Color        // User override or nil
}

func ExtractThemeColors(config ThemeConfig) *Theme {
    theme := &Theme{SourceImage: config.SourceImage}
    
    // Determine primary color
    if config.PrimaryColor != nil {
        theme.Primary = *config.PrimaryColor
    } else {
        theme.Primary = findDominantColor(config.SourceImage)
    }
    
    // Auto-detect theme mode if needed
    mode := config.Mode
    if mode == ModeAuto {
        avgLuminance := calculateAverageLuminance(config.SourceImage)
        _, _, primaryL := theme.Primary.HSL()
        mode = ModeDark
        if avgLuminance > 0.5 && primaryL > 0.5 {
            mode = ModeLight
        }
    }
    theme.IsLight = (mode == ModeLight)
    
    // Derive or use provided colors
    if config.BackgroundColor != nil {
        theme.Background = *config.BackgroundColor
    } else {
        h, s, _ := theme.Primary.HSL()
        if theme.IsLight {
            theme.Background = NewHSL(h, s*0.1, 0.97)  // Light bg
        } else {
            theme.Background = NewHSL(h, s*0.3, 0.08)  // Dark bg
        }
    }
    
    return theme
}
```

### 6. Concurrent Processing
Divides images into 64x64 pixel regions for parallel color extraction using goroutines.

```go
func ExtractPaletteConcurrent(img image.Image, numColors int) []Color {
    bounds := img.Bounds()
    regionSize := 64
    results := make(chan map[Color]int)
    var wg sync.WaitGroup
    
    // Process regions in parallel
    for y := bounds.Min.Y; y < bounds.Max.Y; y += regionSize {
        for x := bounds.Min.X; x < bounds.Max.X; x += regionSize {
            wg.Add(1)
            go func(x0, y0 int) {
                defer wg.Done()
                colors := make(map[Color]int)
                
                maxY := min(y0+regionSize, bounds.Max.Y)
                maxX := min(x0+regionSize, bounds.Max.X)
                
                for y := y0; y < maxY; y++ {
                    for x := x0; x < maxX; x++ {
                        r, g, b, _ := img.At(x, y).RGBA()
                        c := NewRGB(uint8(r>>8), uint8(g>>8), uint8(b>>8))
                        colors[c]++
                    }
                }
                results <- colors
            }(x, y)
        }
    }
    
    // Aggregate results
    go func() {
        wg.Wait()
        close(results)
    }()
    
    allColors := make(map[Color]int)
    for regionColors := range results {
        for c, count := range regionColors {
            allColors[c] += count
        }
    }
    
    return quantizeColors(allColors, numColors)
}
```

### 7. WCAG Contrast Validation
Ensures text readability with AA compliance (4.5:1 minimum contrast ratio).

```go
func ContrastRatio(c1, c2 Color) float64 {
    l1 := relativeLuminance(c1)
    l2 := relativeLuminance(c2)
    
    lighter := math.Max(l1, l2)
    darker := math.Min(l1, l2)
    
    return (lighter + 0.05) / (darker + 0.05)
}

func relativeLuminance(c Color) float64 {
    // Convert to linear RGB with gamma correction
    r := toLinear(float64(c.R) / 255.0)
    g := toLinear(float64(c.G) / 255.0)
    b := toLinear(float64(c.B) / 255.0)
    
    // WCAG luminance formula
    return 0.2126*r + 0.7152*g + 0.0722*b
}

func toLinear(channel float64) float64 {
    if channel <= 0.03928 {
        return channel / 12.92
    }
    return math.Pow((channel+0.055)/1.055, 2.4)
}

func IsAccessible(fg, bg Color) bool {
    return ContrastRatio(fg, bg) >= 4.5
}
```

### 8. Template-Based Configuration Generation
Uses Go's text/template to generate all Omarchy config formats from a unified theme structure.

```go
type Theme struct {
    Name        string
    SourceImage image.Image
    IsLight     bool
    Primary     Color
    Background  Color
    Foreground  Color
    Colors      [16]Color  // Terminal palette
}

const alacrittyTemplate = `[colors.primary]
background = "{{ .Background.Hex }}"
foreground = "{{ .Foreground.Hex }}"

[colors.normal]
black   = "{{ (index .Colors 0).Hex }}"
red     = "{{ (index .Colors 1).Hex }}"
green   = "{{ index .Colors 2 }}"
# ... remaining colors
`

func (t *Theme) Export(basePath string) error {
    themePath := filepath.Join(basePath, t.Name)
    
    // Create structure with backgrounds/ directory
    os.MkdirAll(filepath.Join(themePath, "backgrounds"), 0755)
    
    // Save source image
    imgPath := filepath.Join(themePath, "backgrounds", "wallpaper.jpg")
    saveImage(t.SourceImage, imgPath)
    
    // Generate all configs
    configs := map[string]func() ([]byte, error){
        "alacritty.toml": t.GenerateAlacritty,
        "btop.theme":     t.GenerateBtop,
        "hyprland.conf":  t.GenerateHyprland,
        // ... other generators
    }
    
    for filename, generator := range configs {
        data, err := generator()
        if err != nil {
            return err
        }
        os.WriteFile(filepath.Join(themePath, filename), data, 0644)
    }
    
    // Add light.mode marker if needed
    if t.IsLight {
        os.WriteFile(filepath.Join(themePath, "light.mode"), []byte{}, 0644)
    }
    
    return nil
}
```

## Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| 4K Image Processing | < 2 seconds | Full extraction pipeline |
| Memory Usage | < 100MB | Peak during processing |
| Palette Generation | O(1) | After extraction |
| Color Conversion | 15ns | RGB to HSL |
| Concurrent Regions | 64x64 pixels | Goroutine chunk size |

## Architecture Layers

The refactored architecture uses clear dependency layers:

1. **Foundation Layer**: 
   - `pkg/formats` - Color conversions and formatting (refactored from pkg/color)
   - `pkg/settings` - System configuration and tool behavior
   - `pkg/config` - User preferences and theme-specific overrides

2. **Analysis Layer**:
   - `pkg/analysis` - Image analysis and profile detection (extracted from extractor)

3. **Processing Layer**:
   - `pkg/strategies` - Extraction strategies (extracted from extractor) 
   - `pkg/extractor` - Extraction orchestration (simplified)

4. **Generation Layer**:
   - `pkg/schemes` - Color theory scheme generation
   - `pkg/theme` - Theme file generation

5. **Application Layer**:
   - `cmd/omarchy-theme-gen` - CLI application

### Settings vs Configuration Architecture

**Settings** (`pkg/settings`) - HOW the tool operates:
- Extraction thresholds and parameters
- Algorithm behavior configuration  
- Multi-layer composition (defaults → system → user → workspace → env)

**Configuration** (`pkg/config`) - WHAT the user wants:
- Theme-specific color overrides
- User preferences and customizations
- Stored with generated themes

## Technical Decisions

| Decision | Rationale | Trade-off |
|----------|-----------|-----------|
| **Standard color.RGBA** | Better interoperability, proven implementation | Less control than custom types, but more stable |
| **Purpose-driven extraction** | Colors organized by theme role, not frequency | More complex than frequency-based, but better themes |
| **Layered architecture** | Clear dependencies, maintainable code | More packages than monolith, but easier to test |
| **Multi-strategy extraction** | Adapts to image characteristics | More complex selection logic, but handles edge cases |
| **Settings vs config separation** | Clear distinction between system and user preferences | Two configuration systems, but better organization |
| **Profile-based processing** | Handles edge cases (grayscale, monotone) gracefully | Additional complexity, but robust theme generation |
| **Standard Go testing** | Conventional, widely understood | More setup than execution tests, but standard practice |
| **CLI-first architecture** | Simpler than TUI, easier to implement | Less interactive than TUI, but more reliable |

## Dependencies

- **Go 1.25**: Latest language features
- **Standard Library**: image, image/color, text/template
- **No external dependencies**: Pure Go implementation for reliability and control

## Success Criteria

- [ ] Process 4K images under 2 seconds
- [ ] WCAG AA contrast compliance  
- [ ] Purpose-driven color role assignment
- [ ] Multi-strategy extraction (frequency, saliency)
- [ ] Profile detection (Grayscale, Monotone, Monochromatic, Duotone/Tritone)
- [ ] Color theory scheme generation
- [ ] All Omarchy config formats supported
- [ ] Settings vs config separation
- [ ] Standard library types throughout

## References

### Omarchy Documentation
- [Omarchy Overview](https://learn.omacom.io/2/the-omarchy-manual/91/welcome-to-omarchy) - Desktop environment introduction
- [Themes Documentation](https://learn.omacom.io/2/the-omarchy-manual/52/themes) - Theme system and components
- [Making Your Own Theme](https://learn.omacom.io/2/the-omarchy-manual/92/making-your-own-theme) - Theme creation guide
- [Theme Source Code](https://github.com/basecamp/omarchy/tree/master/themes) - Reference implementations

### Technical Foundations
- CSS Color Module Level 3 - HSL conversion specifications
- WCAG 2.1 Guidelines - Contrast ratio requirements
- CIE LAB Color Space - Perceptual color distance
- Octree Color Quantization (Gervautz & Purgathofer, 1988)

### Project Documentation
- [Development Methodology](development-methodology.md) - Intelligent Development principles
- [Testing Strategy](testing-strategy.md) - Execution test patterns
- [Project Roadmap](../PROJECT.md) - Session plan and progress
