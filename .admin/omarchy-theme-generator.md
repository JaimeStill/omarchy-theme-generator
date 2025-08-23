# Omarchy Theme Generator - Technical Specification

## Project Overview

A Terminal User Interface (TUI) application written in Go that generates Omarchy themes from images. The app extracts dominant colors from user-provided images, allows real-time tweaking of the generated color palette through an interactive preview, and outputs a complete Omarchy theme package ready for installation.

### Core Features
- Image-based dominant color extraction with automatic palette derivation
- Optional user overrides for theme colors (primary, background, foreground)
- Light/dark theme mode selection or automatic detection
- Color theory-driven palette generation (monochromatic, complementary, triadic, analogous, etc.)
- Interactive color adjustment with live preview
- Generation of all required Omarchy configuration files
- Theme validation and export

### Input Parameters
- **Required**: Source image file
- **Optional Overrides** (if not provided, derived from image):
  - Theme mode: light or dark
  - Primary color: base color for palette generation
  - Background color: main background color
  - Foreground color: main text color

### Target Output
A complete theme directory containing:
- `alacritty.toml` - Terminal emulator configuration
- `btop.theme` - System monitor theme
- `hyprland.conf` - Window manager configuration
- `hyprlock.conf` - Lock screen configuration
- `mako.ini` - Notification daemon settings
- `neovim.lua` - Editor colorscheme
- `waybar.css` - Status bar styling
- `walker.css` - Application launcher styling
- `swayosd.css` - On-screen display styling
- `backgrounds/` - Directory containing the source image
- Optional: `light.mode` marker file (for light themes)
- Optional: `icons.theme` file (icon set selection)

## Development Methodology

### Project Memory Configuration

**Concept**: CLAUDE.md files provide persistent project context for Claude Code, automatically loading when sessions begin. This eliminates repetitive explanations and maintains consistent development practices across sessions.

The project's CLAUDE.md file establishes development rules, project structure, execution test patterns, and tracks progress. It serves as both documentation and active configuration, ensuring Claude Code understands the project's unique requirements and methodology.

#### CLAUDE.md File Structure

```markdown
# Omarchy Theme Generator

## Project Overview
A Go-based TUI application that generates Omarchy themes from images using color extraction and palette generation algorithms based on color theory principles.

## Development Philosophy
- **User-driven development**: All code modifications require explicit user direction
- **Claude Code output mode**: Use Explanatory mode to provide educational insights while coding
- **No formal testing during development**: Use minimal execution tests to validate functionality
- **Iterative refinement**: Adapt architecture based on execution test results
- **Progressive enhancement**: Start with core functionality, add features incrementally
- **Color theory focus**: Apply established color harmony principles for aesthetic palettes

## Development Rules
1. Claude Code operates in **Explanatory mode** (`/output-style explanatory`) to provide insights about implementation choices
2. Code modifications are only made when explicitly directed by the user
3. Formal unit tests are deferred until project completion
4. Execution tests should be:
   - Extremely small and focused
   - Free of complicated boilerplate
   - Designed to validate specific functionality
   - Run as soon as realistically possible
5. Architecture evolves based on execution test results
6. Each development session builds upon previous work with continuous validation
```

#### Project Structure Definition

```markdown
## Project Structure
omarchy-theme-gen/
├── CLAUDE.md              # Project memory and configuration
├── cmd/
│   ├── examples/          # Execution test programs
│   │   ├── test_color.go     # Color operations validation
│   │   ├── test_extract.go   # Extraction algorithm testing
│   │   ├── test_strategies.go # Palette strategy comparison
│   │   └── ...
│   └── omarchy-theme/     # Main TUI application
├── pkg/                   # Core packages (importable)
│   ├── color/            # Color types and conversions
│   ├── quantizer/        # Quantization algorithms
│   ├── extractor/        # Image processing
│   ├── palette/          # Color theory strategies
│   ├── template/         # Config generators
│   ├── theme/           # Theme orchestration
│   └── preview/         # Preview rendering
└── internal/
    └── tui/             # TUI components (not importable)
```

#### Execution Test Pattern Documentation

```markdown
## Execution Test Pattern
Standard pattern for all execution tests:

// cmd/examples/test_name.go
package main

import (
    "fmt"
    "omarchy-theme-gen/pkg/[package]"
)

func main() {
    // Minimal setup - no test framework overhead
    // Direct function call - immediate feedback
    // Clear output - visual validation
    fmt.Printf("Result: %v\n", result)
}

Run with: go run cmd/examples/test_name.go [args]

Example tests:
- Color operations: go run cmd/examples/test_color.go
- Image extraction: go run cmd/examples/test_extract.go image.jpg
- Palette strategies: go run cmd/examples/test_strategies.go image.jpg
```

#### Development Commands Reference

```markdown
## Commands
- Run example: go run cmd/examples/[test].go
- Validate code: go vet ./...
- Format code: go fmt ./...
- Module init: go mod init omarchy-theme-gen
- Add dependency: go get [package]
- Tidy modules: go mod tidy
- Run main app: go run cmd/omarchy-theme/main.go
```

#### Progress Tracking

```markdown
## Current Phase
Phase 1: Foundation - Session 1 of 5
Status: Setting up project structure and core color types

## Session Log
### Session 1: [Date]
- Initialized Go module
- Created project structure
- Implemented basic RGB color type
- **Insight**: [Document key learnings]
- **Decision**: [Document architectural choices]
- **Next**: [What to tackle in next session]
```

## Core Technical Concepts

### 1. Color Representation and Types
**Concept**: Define a robust color type system that efficiently handles multiple color spaces. Native RGBA storage with cached HSL values provides the best performance balance.

**Performance Consideration**: Images provide pixels in RGB format natively. Converting every pixel to HSL during extraction would add O(n) overhead. The optimal approach is RGBA native storage with lazy HSL calculation and caching for colors that need manipulation.

```go
// Color represents a color with RGBA native storage and cached HSL
type Color struct {
    R, G, B, A uint8
    
    // Cached HSL values (calculated on first access)
    hsl     *hslCache
    hslOnce sync.Once
}

type hslCache struct {
    H, S, L float64 // Hue [0,1], Saturation [0,1], Lightness [0,1]
}

// NewRGB creates a new opaque color
func NewRGB(r, g, b uint8) Color {
    return Color{R: r, G: g, B: b, A: 255}
}

// NewRGBA creates a new color with alpha
func NewRGBA(r, g, b, a uint8) Color {
    return Color{R: r, G: g, B: b, A: a}
}

// NewHSL creates a color from HSL values
func NewHSL(h, s, l float64) Color {
    r, g, b := hslToRGB(h, s, l)
    c := Color{R: r, G: g, B: b, A: 255}
    // Pre-cache the HSL values since we have them
    c.hsl = &hslCache{H: h, S: s, L: l}
    return c
}

// HSL returns the color in HSL space (cached after first calculation)
func (c *Color) HSL() (h, s, l float64) {
    c.hslOnce.Do(func() {
        h, s, l := rgbToHSL(c.R, c.G, c.B)
        c.hsl = &hslCache{H: h, S: s, L: l}
    })
    return c.hsl.H, c.hsl.S, c.hsl.L
}

// Hex returns the color as a hex string
func (c Color) Hex() string {
    if c.A == 255 {
        return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
    }
    return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}

// CSS returns the color in CSS format
func (c Color) CSS() string {
    if c.A == 255 {
        return fmt.Sprintf("rgb(%d, %d, %d)", c.R, c.G, c.B)
    }
    return fmt.Sprintf("rgba(%d, %d, %d, %.2f)", 
                       c.R, c.G, c.B, float64(c.A)/255.0)
}

// HSLA returns the color in HSLA format
func (c Color) HSLA() string {
    h, s, l := c.HSL()
    if c.A == 255 {
        return fmt.Sprintf("hsl(%.0f, %.0f%%, %.0f%%)", 
                          h*360, s*100, l*100)
    }
    return fmt.Sprintf("hsla(%.0f, %.0f%%, %.0f%%, %.2f)", 
                       h*360, s*100, l*100, float64(c.A)/255.0)
}

// RGBA implements the color.Color interface
func (c Color) RGBA() (r, g, b, a uint32) {
    r = uint32(c.R) << 8 | uint32(c.R)
    g = uint32(c.G) << 8 | uint32(c.G)
    b = uint32(c.B) << 8 | uint32(c.B)
    a = uint32(c.A) << 8 | uint32(c.A)
    return
}

// WithAlpha returns a new color with the specified alpha
func (c Color) WithAlpha(a uint8) Color {
    return Color{R: c.R, G: c.G, B: c.B, A: a}
}

// AdjustHSL returns a new color with adjusted HSL values
func (c Color) AdjustHSL(hDelta, sDelta, lDelta float64) Color {
    h, s, l := c.HSL()
    h = math.Mod(h+hDelta+1, 1)
    s = clamp(s+sDelta, 0, 1)
    l = clamp(l+lDelta, 0, 1)
    return NewHSL(h, s, l).WithAlpha(c.A)
}
```

### 2. Color Space Conversion
**Concept**: Convert between RGB and perceptually uniform color spaces (HSL, LAB) to enable intuitive color manipulation and accurate color relationships.

```go
// RGB to HSL conversion for intuitive adjustments
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

// HSL to RGB conversion
func HSLToRGB(h, s, l float64) Color {
    var r, g, b float64
    
    if s == 0 {
        r, g, b = l, l, l // achromatic
    } else {
        var q float64
        if l < 0.5 {
            q = l * (1 + s)
        } else {
            q = l + s - l*s
        }
        p := 2*l - q
        r = hueToRGB(p, q, h+1.0/3.0)
        g = hueToRGB(p, q, h)
        b = hueToRGB(p, q, h-1.0/3.0)
    }
    
    return NewRGB(uint8(r*255), uint8(g*255), uint8(b*255))
}
```

### 3. Image Color Extraction with Optional Overrides
**Concept**: Extract color information from images while respecting user preferences. When colors are provided by the user, use them; otherwise derive them intelligently from the image.

```go
// ThemeConfig holds user preferences and derived values
type ThemeConfig struct {
    SourceImage    image.Image
    Mode           ThemeMode      // Light, Dark, or Auto
    PrimaryColor   *Color         // User override or nil
    BackgroundColor *Color        // User override or nil
    ForegroundColor *Color        // User override or nil
}

type ThemeMode int
const (
    ModeAuto ThemeMode = iota
    ModeLight
    ModeDark
)

// ExtractThemeColors derives theme colors respecting user overrides
func ExtractThemeColors(config ThemeConfig) *Theme {
    theme := &Theme{
        SourceImage: config.SourceImage,
    }
    
    // Determine primary color
    if config.PrimaryColor != nil {
        theme.Primary = *config.PrimaryColor
    } else {
        theme.Primary = findDominantColor(config.SourceImage)
    }
    
    // Determine theme mode if auto
    mode := config.Mode
    if mode == ModeAuto {
        mode = detectThemeMode(config.SourceImage, theme.Primary)
    }
    theme.IsLight = (mode == ModeLight)
    
    // Determine background color
    if config.BackgroundColor != nil {
        theme.Background = *config.BackgroundColor
    } else {
        theme.Background = deriveBackground(theme.Primary, theme.IsLight)
    }
    
    // Determine foreground color
    if config.ForegroundColor != nil {
        theme.Foreground = *config.ForegroundColor
    } else {
        theme.Foreground = deriveForeground(theme.Background, theme.IsLight)
    }
    
    return theme
}

// detectThemeMode analyzes image brightness to determine light/dark
func detectThemeMode(img image.Image, primary Color) ThemeMode {
    avgLuminance := calculateAverageLuminance(img)
    _, _, primaryL := primary.HSL()
    
    // If image and primary are both bright, suggest light theme
    if avgLuminance > 0.5 && primaryL > 0.5 {
        return ModeLight
    }
    return ModeDark
}

// deriveBackground generates appropriate background from primary
func deriveBackground(primary Color, isLight bool) Color {
    h, s, _ := primary.HSL()
    
    if isLight {
        // Light theme: very light, desaturated version
        return NewHSL(h, s*0.1, 0.97)
    } else {
        // Dark theme: very dark, slightly desaturated version
        return NewHSL(h, s*0.3, 0.08)
    }
}

// deriveForeground generates readable text color for background
func deriveForeground(bg Color, isLight bool) Color {
    if isLight {
        // Dark text on light background
        h, s, _ := bg.HSL()
        return NewHSL(h, s*0.2, 0.15)
    } else {
        // Light text on dark background
        h, s, _ := bg.HSL()
        return NewHSL(h, s*0.1, 0.92)
    }
}
```

### 4. Color Quantization Algorithms
**Concept**: Reduce millions of colors in an image to a manageable palette while preserving visual essence using tree-based or clustering algorithms.

```go
// Octree node for color quantization
type OctreeNode struct {
    red, green, blue uint32
    pixelCount       uint32
    paletteIndex     int
    children         [8]*OctreeNode
}

// Insert color into octree
func (node *OctreeNode) Insert(r, g, b uint8, level int) {
    if level >= maxDepth {
        node.red += uint32(r)
        node.green += uint32(g)
        node.blue += uint32(b)
        node.pixelCount++
        return
    }
    
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

// GetPalette extracts the color palette from the octree
func (node *OctreeNode) GetPalette() []Color {
    if node.pixelCount > 0 {
        r := uint8(node.red / node.pixelCount)
        g := uint8(node.green / node.pixelCount)
        b := uint8(node.blue / node.pixelCount)
        return []Color{NewRGB(r, g, b)}
    }
    
    var palette []Color
    for _, child := range node.children {
        if child != nil {
            palette = append(palette, child.GetPalette()...)
        }
    }
    return palette
}
```

### 5. Color Theory-Based Palette Generation
**Concept**: Apply established color harmony principles (monochromatic, complementary, triadic, analogous) to generate aesthetically pleasing palettes from a base color.

```go
// PaletteStrategy defines how to generate a color palette
type PaletteStrategy interface {
    Generate(baseColor Color, count int) []Color
    Name() string
}

// MonochromaticStrategy generates variations of a single hue
type MonochromaticStrategy struct{}

func (m MonochromaticStrategy) Name() string { return "Monochromatic" }

func (m MonochromaticStrategy) Generate(base Color, count int) []Color {
    h, s, l := RGBToHSL(base.R, base.G, base.B)
    palette := make([]Color, count)
    
    for i := 0; i < count; i++ {
        // Vary lightness and saturation while keeping hue constant
        factor := float64(i) / float64(count-1)
        newL := l + (factor-0.5)*0.4 // Vary lightness ±40%
        newS := s * (0.5 + factor*0.5) // Vary saturation 50-100%
        
        palette[i] = HSLToRGB(h, clamp(newS, 0, 1), clamp(newL, 0, 1))
    }
    return palette
}

// TriadicStrategy generates colors 120° apart on the color wheel
type TriadicStrategy struct{}

func (t TriadicStrategy) Name() string { return "Triadic" }

func (t TriadicStrategy) Generate(base Color, count int) []Color {
    h, s, l := RGBToHSL(base.R, base.G, base.B)
    
    // Primary triadic colors
    palette := []Color{base}
    perColor := (count - 1) / 2
    
    // Add variations of first triad (120° rotation)
    for i := 0; i < perColor; i++ {
        factor := float64(i) / float64(perColor)
        newL := l + (factor-0.5)*0.3
        palette = append(palette, HSLToRGB(
            math.Mod(h+1.0/3.0, 1), s, clamp(newL, 0, 1)))
    }
    
    // Add variations of second triad (240° rotation)
    for i := 0; i < count-1-perColor; i++ {
        factor := float64(i) / float64(count-1-perColor)
        newL := l + (factor-0.5)*0.3
        palette = append(palette, HSLToRGB(
            math.Mod(h+2.0/3.0, 1), s, clamp(newL, 0, 1)))
    }
    
    return palette
}
```

### 6. Template-Based Configuration Generation
**Concept**: Use Go's text/template package to generate various configuration file formats (TOML, INI, CSS, Lua) from a unified theme structure.

```go
// Theme represents the complete color scheme and assets
type Theme struct {
    Name        string
    SourceImage image.Image    // Original image for backgrounds/
    IsLight     bool           // Light or dark theme
    Primary     Color          // Primary/accent color
    Background  Color          // Main background
    Foreground  Color          // Main text color
    Colors      [16]Color      // Full terminal palette
}

// Export creates the complete theme directory structure
func (t *Theme) Export(basePath string) error {
    themePath := filepath.Join(basePath, t.Name)
    
    // Create directory structure
    dirs := []string{
        themePath,
        filepath.Join(themePath, "backgrounds"),
    }
    
    for _, dir := range dirs {
        if err := os.MkdirAll(dir, 0755); err != nil {
            return err
        }
    }
    
    // Save source image to backgrounds/
    imgPath := filepath.Join(themePath, "backgrounds", "wallpaper.jpg")
    if err := saveImage(t.SourceImage, imgPath); err != nil {
        return err
    }
    
    // Generate all config files
    configs := map[string]func() ([]byte, error){
        "alacritty.toml": t.GenerateAlacritty,
        "btop.theme":     t.GenerateBtop,
        "hyprland.conf":  t.GenerateHyprland,
        "hyprlock.conf":  t.GenerateHyprlock,
        "mako.ini":       t.GenerateMako,
        "neovim.lua":     t.GenerateNeovim,
        "waybar.css":     t.GenerateWaybar,
        "walker.css":     t.GenerateWalker,
        "swayosd.css":    t.GenerateSwayOSD,
    }
    
    for filename, generator := range configs {
        data, err := generator()
        if err != nil {
            return err
        }
        
        path := filepath.Join(themePath, filename)
        if err := os.WriteFile(path, data, 0644); err != nil {
            return err
        }
    }
    
    // Add light.mode marker if needed
    if t.IsLight {
        marker := filepath.Join(themePath, "light.mode")
        if err := os.WriteFile(marker, []byte{}, 0644); err != nil {
            return err
        }
    }
    
    return nil
}

// Alacritty TOML template
const alacrittyTemplate = `[colors.primary]
background = "{{ .Background.Hex }}"
foreground = "{{ .Foreground.Hex }}"

[colors.normal]
black   = "{{ index .Colors 0 }}"
red     = "{{ index .Colors 1 }}"
green   = "{{ index .Colors 2 }}"
yellow  = "{{ index .Colors 3 }}"
blue    = "{{ index .Colors 4 }}"
magenta = "{{ index .Colors 5 }}"
cyan    = "{{ index .Colors 6 }}"
white   = "{{ index .Colors 7 }}"
`

// GenerateAlacritty creates the alacritty.toml configuration
func (t *Theme) GenerateAlacritty() ([]byte, error) {
    tmpl := template.New("alacritty")
    tmpl, err := tmpl.Parse(alacrittyTemplate)
    if err != nil {
        return nil, err
    }
    
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, t); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}
```

### 7. Contrast Ratio Validation
**Concept**: Calculate WCAG contrast ratios to ensure generated themes meet accessibility standards for readability.

```go
// ContrastRatio calculates WCAG contrast ratio between two colors
func ContrastRatio(c1, c2 Color) float64 {
    l1 := relativeLuminance(c1)
    l2 := relativeLuminance(c2)
    
    lighter := math.Max(l1, l2)
    darker := math.Min(l1, l2)
    
    return (lighter + 0.05) / (darker + 0.05)
}

// relativeLuminance calculates the relative luminance of a color
func relativeLuminance(c Color) float64 {
    // Convert to linear RGB
    r := toLinear(float64(c.R) / 255.0)
    g := toLinear(float64(c.G) / 255.0)
    b := toLinear(float64(c.B) / 255.0)
    
    // Apply WCAG luminance formula
    return 0.2126*r + 0.7152*g + 0.0722*b
}

// toLinear applies gamma correction
func toLinear(channel float64) float64 {
    if channel <= 0.03928 {
        return channel / 12.92
    }
    return math.Pow((channel+0.055)/1.055, 2.4)
}

// IsAccessible checks if contrast meets WCAG AA standard (4.5:1)
func IsAccessible(fg, bg Color) bool {
    return ContrastRatio(fg, bg) >= 4.5
}
```

### 8. Concurrent Image Processing
**Concept**: Process image regions in parallel using goroutines and channels for faster color extraction on multi-core systems.

```go
// ExtractPaletteConcurrent extracts colors using parallel processing
func ExtractPaletteConcurrent(img image.Image, numColors int) []Color {
    bounds := img.Bounds()
    regionSize := 64 // Process in 64x64 chunks
    
    type regionResult struct {
        colors map[Color]int
    }
    
    // Channel for collecting results
    results := make(chan regionResult)
    var wg sync.WaitGroup
    
    // Process regions concurrently
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
                results <- regionResult{colors}
            }(x, y)
        }
    }
    
    // Close results channel when done
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Merge results
    allColors := make(map[Color]int)
    for result := range results {
        for c, count := range result.colors {
            allColors[c] += count
        }
    }
    
    // Apply quantization to get final palette
    return quantizeColors(allColors, numColors)
}
```

### 9. Bubble Tea TUI Components
**Concept**: Build composable, reactive terminal UI components using the Bubble Tea framework for user interaction and real-time feedback.

```go
// ColorAdjuster provides interactive HSL adjustment
type ColorAdjuster struct {
    color    Color
    hue      float64
    sat      float64
    light    float64
    selected int // 0: hue, 1: saturation, 2: lightness
}

// Init implements tea.Model
func (c ColorAdjuster) Init() tea.Cmd {
    return nil
}

// Update handles user input
func (c ColorAdjuster) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up", "k":
            c.selected = (c.selected - 1 + 3) % 3
        case "down", "j":
            c.selected = (c.selected + 1) % 3
        case "left", "h":
            c.adjustValue(-0.01)
        case "right", "l":
            c.adjustValue(0.01)
        case "q", "ctrl+c":
            return c, tea.Quit
        }
    }
    return c, nil
}

// View renders the component
func (c ColorAdjuster) View() string {
    items := []string{
        fmt.Sprintf("Hue:        %.2f", c.hue),
        fmt.Sprintf("Saturation: %.2f", c.sat),
        fmt.Sprintf("Lightness:  %.2f", c.light),
    }
    
    // Add selection indicator
    items[c.selected] = "> " + items[c.selected]
    
    return strings.Join(items, "\n") + 
           fmt.Sprintf("\n\nColor: %s", c.color)
}

// adjustValue modifies the selected HSL component
func (c *ColorAdjuster) adjustValue(delta float64) {
    switch c.selected {
    case 0: // hue
        c.hue = math.Mod(c.hue+delta+1, 1)
    case 1: // saturation
        c.sat = clamp(c.sat+delta, 0, 1)
    case 2: // lightness
        c.light = clamp(c.light+delta, 0, 1)
    }
    c.color = HSLToRGB(c.hue, c.sat, c.light)
}
```

## Architecture Layers

### Layer 1: Core Domain (No Dependencies)
```
pkg/
├── color/
│   ├── color.go         // Color type definitions
│   ├── space.go          // Color space conversions
│   └── palette.go        // Palette type and operations
```

### Layer 2: Algorithms (Depends on Core)
```
pkg/
├── quantizer/
│   ├── octree.go         // Octree quantization
│   ├── kmeans.go         // K-means clustering
│   └── median_cut.go     // Median cut algorithm
├── extractor/
│   ├── extractor.go      // Image color extraction interface
│   ├── dominant.go       // Dominant color detection
│   └── concurrent.go     // Parallel extraction
├── palette/
│   ├── strategy.go       // Palette generation interface
│   ├── monochromatic.go  // Single hue variations
│   ├── complementary.go  // Opposite colors
│   ├── triadic.go        // 120° separated colors
│   ├── analogous.go      // Adjacent colors
│   └── tetradic.go       // Four-color schemes
```

### Layer 3: Infrastructure (Depends on Core)
```
pkg/
├── template/
│   ├── engine.go         // Template engine
│   ├── alacritty.go      // Alacritty generator
│   ├── hyprland.go       // Hyprland generator
│   ├── mako.go           // Mako generator
│   └── ... (other generators)
├── validator/
│   ├── contrast.go       // Accessibility validation
│   └── syntax.go         // Config syntax validation
```

### Layer 4: Application (Depends on Core, Algorithms, Infrastructure)
```
pkg/
├── theme/
│   ├── theme.go          // Theme aggregate
│   ├── generator.go      // Theme generation orchestration
│   └── exporter.go       // Theme export functionality
├── preview/
│   ├── renderer.go       // Preview rendering logic
│   └── mockup.go         // UI mockup generation
```

### Layer 5: Presentation (Depends on All)
```
cmd/
├── examples/             // Execution test programs
│   ├── test_color.go    // Test color operations
│   ├── test_extract.go  // Test extraction
│   └── ...
└── omarchy-theme/
    └── main.go           // Entry point
internal/
└── tui/
    ├── app.go            // Main TUI application
    ├── components/
    │   ├── picker.go     // Color picker component
    │   ├── preview.go    // Preview component
    │   ├── adjuster.go   // Adjustment controls
    │   └── exporter.go   // Export dialog
    └── styles/
        └── theme.go      // TUI styling
```

## User-Driven Development Roadmap

### Phase 1: Foundation (Sessions 1-5)

**Session 1: Project Setup & Core Types**
- Initialize Go module (go mod init omarchy-theme-gen)
- Create project structure and CLAUDE.md
- Define basic color type with RGB representation
- **Execution Test**: `cmd/examples/test_color.go` - Create colors and print values
- Adapt based on test results

**Session 2: Color Space Conversions**
- Implement RGB to HSL conversion
- Implement HSL to RGB conversion
- **Execution Test**: `cmd/examples/test_conversion.go` - Convert known values and verify
- Refine algorithms based on accuracy

**Session 3: Basic Image Loading**
- Set up image loading from file path
- Iterate through pixels and count colors
- **Execution Test**: `cmd/examples/test_load_image.go` - Load image and print dimensions/pixel count
- Adjust approach based on performance

**Session 4: Color Theory-Based Extraction**
- Build color frequency map from image
- Identify dominant color (or use provided primary)
- Implement palette strategies (monochromatic, complementary, triadic, analogous)
- Handle light/dark mode detection
- **Execution Test**: `cmd/examples/test_extract_strategies.go` - Extract palettes with optional overrides
- Compare aesthetic results

**Session 5: First Template Generator**
- Create template interface
- Implement basic alacritty.toml generator
- **Execution Test**: `cmd/examples/test_generate_alacritty.go` - Generate config and validate syntax
- Refine template based on output

### Phase 2: Algorithms (Sessions 6-10)

**Session 6: Octree Implementation**
- Build octree data structure
- Implement color insertion
- **Execution Test**: `cmd/examples/test_octree.go` - Build octree from image and extract palette
- Compare with simple extraction

**Session 7: Dominant Color Detection**
- Implement color clustering for dominant color
- Add perceptual color distance metrics
- Test different dominant color algorithms
- **Execution Test**: `cmd/examples/test_dominant.go` - Compare dominant color detection methods
- Choose most reliable approach

**Session 8: Concurrent Processing**
- Divide image into regions
- Process regions in parallel
- **Execution Test**: `cmd/examples/test_concurrent.go` - Benchmark serial vs parallel
- Optimize based on results

**Session 9: Advanced Palette Strategies**
- Implement tetradic (square) color scheme
- Add split-complementary strategy
- Create custom weighted strategies
- **Execution Test**: `cmd/examples/test_advanced_harmony.go` - Generate complex color schemes
- Compare with professional palettes

**Session 10: Accessibility Validation**
- Implement contrast ratio calculation
- Add automatic adjustment
- **Execution Test**: `cmd/examples/test_contrast.go` - Check various color pairs
- Ensure WCAG compliance

### Phase 3: Configuration Generation (Sessions 11-15)

**Session 11: Multiple Config Generators**
- Implement mako.ini generator
- Add btop.theme generator
- **Execution Test**: `cmd/examples/test_generate_configs.go` - Generate multiple configs
- Validate each format

**Session 12: CSS Generation**
- Create waybar.css generator
- Add walker.css generator
- **Execution Test**: `cmd/examples/test_generate_css.go` - Generate and validate CSS
- Check for syntax errors

**Session 13: Lua Generation**
- Implement neovim.lua generator
- Map syntax highlighting groups
- **Execution Test**: `cmd/examples/test_generate_lua.go` - Generate Lua config
- Test in actual Neovim if available

**Session 14: Hyprland Configuration**
- Create hyprland.conf generator
- Add hyprlock.conf generator
- **Execution Test**: `cmd/examples/test_generate_hypr.go` - Generate window manager configs
- Validate configuration syntax

**Session 15: Complete Theme Package**
- Assemble all generators
- Create directory structure
- Copy source image to backgrounds/ directory
- Add light.mode file if applicable
- **Execution Test**: `cmd/examples/test_full_theme.go` - Generate complete theme
- Verify all files present and valid

### Phase 4: TUI Development (Sessions 16-22)

**Session 16: Bubble Tea Setup**
- Initialize Bubble Tea app
- Create basic model
- **Execution Test**: `cmd/examples/test_tui_basic.go` - Launch TUI and handle quit
- Verify basic interaction

**Session 17: File Selection & Theme Options**
- Create file browser view
- Add theme mode selector (light/dark/auto)
- Add optional color override inputs
- Implement navigation
- **Execution Test**: Run TUI and test file selection with options
- Ensure smooth navigation and input

**Session 18: Palette Display & Strategy Selection**
- Show extracted dominant color
- Add palette strategy selector (monochromatic, triadic, etc.)
- Display generated palette based on selected strategy
- Show color values and relationships
- **Execution Test**: Load image, select strategies, and see different palettes
- Verify color theory application

**Session 19: Color Adjustment**
- Create HSL sliders
- Implement real-time updates
- **Execution Test**: Adjust colors and see changes
- Ensure responsive controls

**Session 20: Preview Component**
- Design terminal mockup
- Show color applications
- **Execution Test**: Display preview with current palette
- Validate preview accuracy

**Session 21: Export Functionality**
- Add export dialog
- Implement file writing
- **Execution Test**: Export theme and verify files
- Check generated theme structure

**Session 22: Full Integration**
- Wire all components
- Add state management
- **Execution Test**: Complete workflow from image to export
- Ensure smooth user experience

### Phase 5: Polish & Features (Sessions 23-28)

**Session 23: History & Undo**
- Implement undo/redo
- Add command history
- **Execution Test**: Make changes and undo them
- Verify state consistency

**Session 24: Theme Variations**
- Generate light/dark modes
- Create color variations
- **Execution Test**: Generate multiple variations
- Compare output quality

**Session 25: Batch Processing**
- Support multiple images
- Add batch export
- **Execution Test**: Process image folder
- Verify batch operation

**Session 26: Settings & Persistence**
- Add configuration file
- Remember preferences
- **Execution Test**: Save and load settings
- Ensure persistence works

**Session 27: Error Handling**
- Add comprehensive error handling
- Implement recovery
- **Execution Test**: Test with invalid inputs
- Verify graceful handling

**Session 28: Performance Optimization**
- Profile application
- Optimize hot paths
- **Execution Test**: Process large images
- Measure improvements

### Phase 6: Finalization (Sessions 29-30)

**Session 29: Documentation**
- Write user documentation
- Add help system
- Create example themes
- Document API if needed

**Session 30: Testing & Release**
- Write formal test suite
- Add integration tests
- Build release binaries
- Create installation guide

## Execution Test Examples

### Example: Color Operations Test
```go
// cmd/examples/test_color.go
package main

import (
    "fmt"
    "omarchy-theme-gen/pkg/color"
)

func main() {
    // Create a color
    c := color.NewRGB(255, 128, 64)
    fmt.Printf("RGB: %v\n", c)
    
    // Convert to HSL
    h, s, l := c.ToHSL()
    fmt.Printf("HSL: %.2f, %.2f, %.2f\n", h, s, l)
    
    // Convert back
    c2 := color.FromHSL(h, s, l)
    fmt.Printf("RGB (converted back): %v\n", c2)
}
```

### Example: Palette Strategy Test
```go
// cmd/examples/test_extract_strategies.go
package main

import (
    "fmt"
    "image"
    _ "image/jpeg"
    _ "image/png"
    "os"
    "omarchy-theme-gen/pkg/extractor"
    "omarchy-theme-gen/pkg/palette"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run test_extract_strategies.go <image-path>")
        return
    }
    
    file, err := os.Open(os.Args[1])
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    img, _, err := image.Decode(file)
    if err != nil {
        panic(err)
    }
    
    // Find dominant color
    dominant := extractor.FindDominantColor(img)
    fmt.Printf("Dominant color: %v\n\n", dominant)
    
    // Test different strategies
    strategies := map[string]palette.Strategy{
        "Monochromatic": &palette.MonochromaticStrategy{},
        "Complementary": &palette.ComplementaryStrategy{},
        "Triadic":       &palette.TriadicStrategy{},
        "Analogous":     &palette.AnalogousStrategy{},
    }
    
    for name, strategy := range strategies {
        fmt.Printf("%s Palette:\n", name)
        colors := strategy.Generate(dominant, 8)
        for i, c := range colors {
            fmt.Printf("  %d: %v\n", i, c)
        }
        fmt.Println()
    }
}
```

### Example: Theme Generation with Overrides
```go
// cmd/examples/test_theme_overrides.go
package main

import (
    "flag"
    "fmt"
    "image"
    _ "image/jpeg"
    _ "image/png"
    "os"
    "omarchy-theme-gen/pkg/theme"
)

func main() {
    // Parse command line flags
    imagePath := flag.String("image", "", "Source image (required)")
    mode := flag.String("mode", "auto", "Theme mode: light|dark|auto")
    primary := flag.String("primary", "", "Primary color (hex)")
    bg := flag.String("bg", "", "Background color (hex)")
    fg := flag.String("fg", "", "Foreground color (hex)")
    flag.Parse()
    
    if *imagePath == "" {
        fmt.Println("Usage: go run test_theme_overrides.go -image=path.jpg [options]")
        return
    }
    
    // Load image
    file, err := os.Open(*imagePath)
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    img, _, err := image.Decode(file)
    if err != nil {
        panic(err)
    }
    
    // Build config with optional overrides
    config := theme.Config{
        SourceImage: img,
        Mode:       parseMode(*mode),
    }
    
    if *primary != "" {
        config.PrimaryColor = parseHexColor(*primary)
    }
    if *bg != "" {
        config.BackgroundColor = parseHexColor(*bg)
    }
    if *fg != "" {
        config.ForegroundColor = parseHexColor(*fg)
    }
    
    // Generate theme
    t := theme.Generate(config)
    
    // Display results
    fmt.Printf("Theme: %s\n", t.Name)
    fmt.Printf("Mode: %s\n", t.Mode())
    fmt.Printf("Primary: %s\n", t.Primary.Hex())
    fmt.Printf("Background: %s\n", t.Background.Hex())
    fmt.Printf("Foreground: %s\n", t.Foreground.Hex())
    
    // Export theme
    err = t.Export("./output")
    if err != nil {
        panic(err)
    }
    fmt.Println("Theme exported to ./output")
}
```

## Success Metrics

- **Performance**: Process a 4K image in under 2 seconds
- **Quality**: Generated themes pass WCAG AA contrast requirements
- **Aesthetics**: Palettes follow established color theory principles
- **Flexibility**: Support at least 5 different palette generation strategies
- **Usability**: Complete theme generation in under 10 user interactions
- **Compatibility**: All generated configs validated by their respective tools
- **Completeness**: Every theme includes source image in backgrounds/ directory
- **Reliability**: Zero crashes during normal operation

## Technical Decisions

- **Go 1.25**: Latest language features and performance improvements
- **Bubble Tea**: Modern, well-maintained TUI framework
- **RGBA with Cached HSL**: Native RGBA for efficient extraction, cached HSL for manipulation
- **Dominant Color + Strategy**: Extract key color then apply color theory
- **Multiple Quantization Options**: Octree, k-means, median-cut for flexibility
- **Template-based Generation**: Maintainable and extensible
- **Concurrent Processing**: Utilize modern multi-core systems
- **User Override Support**: Allow manual color selection while maintaining intelligent defaults
- **User-driven Development**: Claude Code in Explanatory mode with execution tests
