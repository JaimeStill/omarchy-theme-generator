# Omarchy Theme Style Guide

## Overview
This guide documents the format, structure, and conventions for generating themes compatible with the Omarchy desktop environment. Based on analysis of native Omarchy themes (catppuccin, catppuccin-latte, everforest, gruvbox, kanagawa, etc.).

## Theme Directory Structure

### Required Structure
```
theme-name/
├── alacritty.toml          # Terminal emulator colors
├── btop.theme              # System monitor theme  
├── hyprland.conf           # Window manager colors
├── hyprlock.conf           # Lock screen appearance
├── mako.ini                # Notification daemon
├── neovim.lua              # Editor colorscheme
├── swayosd.css             # On-screen display
├── walker.css              # Application launcher
├── waybar.css              # Status bar styling
├── theme-gen.json          # Theme metadata (our addition)
├── backgrounds/            # Wallpaper directory
│   └── [image-file]        # Source image
└── [optional files]       # See below
```

### Optional Files
- `chromium.theme` - Browser theme reference (minimal content)
- `icons.theme` - Icon theme reference (minimal content)  
- `light.mode` - Present only for light themes (indicates prefer-light mode)

## Color Format Standards

### Format by Configuration Type

| File Type | Color Format | Example | Notes |
|-----------|-------------|---------|-------|
| `alacritty.toml` | Hex strings | `background = "#24273a"` | No alpha |
| `btop.theme` | Hex strings | `theme[main_fg]="#cad3f5"` | No alpha |
| `hyprland.conf` | RGB function | `col.active_border = rgb(c6d0f5)` | No # prefix |
| `hyprlock.conf` | RGBA function | `color = rgba(36, 39, 58, 1.0)` | Decimal RGBA |
| `mako.ini` | Hex strings | `text-color=#cad3f5` | No alpha |
| `neovim.lua` | Hex strings | `bg = "#24273a"` | Lua table format |
| `CSS files` | Hex strings | `background-color: #24273a;` | Standard CSS |

### Color Conversion Requirements
From our HEXA format in `theme-gen.json`:
- **HEXA → Hex**: Strip alpha component for most formats
- **HEXA → RGB**: Extract RGB components for hyprland.conf  
- **HEXA → RGBA**: Convert to decimal for hyprlock.conf
- **HEXA → CSS**: Standard hex format for all CSS files

## Theme Metadata (theme-gen.json)

### Required Structure
```json
{
  "version": "1.0.0",
  "source_image": "backgrounds/photo.jpg",
  "extracted_colors": {
    "dominant": "#2e3440ff",
    "palette": [
      "#2e3440ff",
      "#88c0d0ff", 
      "#81a1c1ff"
    ],
    "unique_count": 1247,
    "coverage_map": {
      "#2e3440ff": {"pixels": 28470, "percentage": 34.2}
    }
  },
  "analysis": {
    "is_grayscale": false,
    "is_monochromatic": false,
    "average_luminance": 0.42,
    "perceptual_diversity": 0.67
  },
  "generation": {
    "mode": "dark",
    "scheme": "complementary",
    "primary": "#88c0d0ff",
    "background": "#2e3440ff",
    "foreground": "#eceff4ff",
    "accent1": "#bf616aff",
    "accent2": "#d08770ff",
    "accent3": "#ebcb8bff",
    "timestamp": "2025-08-31T10:30:00Z"
  },
  "overrides": {
    "background": null,
    "foreground": null,
    "primary": null
  }
}
```

### Color Storage
- **All colors in HEXA format** (#RRGGBBAA) preserving full precision
- **No conversion losses** - parse once, convert per output format
- **Enables refinement** - cached extraction data allows scheme changes

## Configuration File Standards

### alacritty.toml Structure
```toml
[colors.primary]
background = "#24273a"
foreground = "#cad3f5"
dim_foreground = "#8087a2"

[colors.cursor]
text = "#24273a"
cursor = "#f4dbd6"

[colors.normal]
black = "#494d64"
red = "#ed8796"
# ... (8 colors)

[colors.bright]  
black = "#5b6078"
red = "#ed8796"
# ... (8 colors)
```

### btop.theme Structure
```ini
# Main colors
theme[main_bg]="#24273a"
theme[main_fg]="#cad3f5"
theme[title]="#cad3f5"

# Graph colors  
theme[inactive_fg]="#8087a2"
theme[proc_misc]="#f4dbd6"

# CPU/Memory specific
theme[cpu_box]="#c6d0f5"
theme[mem_box]="#a6da95"
# ... (50+ theme properties)
```

### CSS Files (waybar.css, walker.css, swayosd.css)
```css
/* Standard CSS with hex colors */
* {
    background-color: #24273a;
    color: #cad3f5;
    border-color: #c6d0f5;
}

/* Component-specific styling */
.module {
    padding: 4px 8px;
    margin: 2px;
}
```

## Light vs Dark Mode Detection

### Light Mode Characteristics
- **light.mode file present** - Key indicator
- **Background colors** - High lightness (L > 0.8 in HSL)
- **Foreground colors** - Low lightness (L < 0.3 in HSL)
- **Example themes**: catppuccin-latte, rose-pine

### Dark Mode Characteristics  
- **No light.mode file** - Default assumption
- **Background colors** - Low lightness (L < 0.2 in HSL)
- **Foreground colors** - High lightness (L > 0.7 in HSL)
- **Example themes**: catppuccin, gruvbox, kanagawa

## Color Theory Scheme Implementation

### Supported Schemes
1. **Monochromatic** - Single hue with lightness/saturation variations
2. **Analogous** - Adjacent hues (±30° on color wheel)  
3. **Complementary** - Opposite hues (180° separation)
4. **Split-Complementary** - Base hue + two adjacent to complement
5. **Triadic** - Three hues equally spaced (120° separation)
6. **Tetradic** - Two complementary pairs (90° separation)
7. **Square** - Four hues equally spaced (90° separation)

### Scheme-to-Config Mapping
- **Primary**: Most prominent accent color
- **Background**: Dominant image color or scheme-derived
- **Foreground**: High contrast to background (WCAG AA compliance)
- **Accent colors**: Scheme-derived harmonious colors
- **Terminal colors**: Map to 16-color ANSI palette

## File Generation Priority

### Critical Files (Must Generate)
1. `alacritty.toml` - Terminal colors
2. `theme-gen.json` - Metadata for refinement
3. `backgrounds/` - Source image

### Standard Files (Should Generate)  
4. `btop.theme` - System monitor
5. `mako.ini` - Notifications
6. `waybar.css` - Status bar
7. `hyprland.conf` - Window manager
8. `hyprlock.conf` - Lock screen

### Optional Files (Can Generate)
9. `neovim.lua` - Editor (if applicable)
10. `walker.css` - App launcher
11. `swayosd.css` - OSD styling
12. `light.mode` - Light theme marker

## Validation Requirements

### WCAG Compliance
- **Minimum contrast**: 4.5:1 (AA standard)
- **Text on backgrounds**: Validate all fg/bg combinations  
- **UI elements**: 3.0:1 minimum for non-text

### Format Validation
- **Hex colors**: Valid #RRGGBB format
- **TOML syntax**: Valid TOML structure
- **CSS syntax**: Valid CSS properties
- **JSON structure**: Valid theme-gen.json schema

### Integration Testing
- **Omarchy compatibility**: Theme loads without errors
- **File permissions**: 644 for files, 755 for directories
- **Image copying**: Source image properly copied to backgrounds/

## Development Notes

### Color Precision
- **Store in HEXA**: Preserve full precision including alpha
- **Convert per format**: Avoid cumulative conversion losses
- **Round carefully**: Consistent rounding for display values

### Performance Considerations
- **Template-based generation**: Pre-defined templates per format
- **Batch color conversion**: Convert all colors once per format
- **Efficient file I/O**: Write all files in single pass

### Error Handling
- **Missing source image**: Clear error message
- **Invalid colors**: Fallback to safe defaults
- **File permissions**: Handle read-only directories gracefully
- **Malformed metadata**: Validate theme-gen.json structure

This style guide ensures consistent, high-quality theme generation that integrates seamlessly with the Omarchy desktop environment.

## References

- [Omarchy Overview](https://learn.omacom.io/2/the-omarchy-manual/91/welcome-to-omarchy)
- [Omarchy Themes](https://learn.omacom.io/2/the-omarchy-manual/52/themes)
- [Making your own Theme](https://learn.omacom.io/2/the-omarchy-manual/92/making-your-own-theme)
- [Themes Repo Source](https://github.com/basecamp/omarchy/tree/master/themes)
