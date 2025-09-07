# Themes

This repository captures all of the details associated with Omarchy theming: the components that are used and the values established in the baseline Omarchy themes.

## Components

### Alacritty

File: `alacritty.toml`

Standardized configuration. Variables set varies between themes.

References:
- repo: https://github.com/alacritty/alacritty
- site: https://alacritty.org/
- config: https://alacritty.org/config-alacritty.html

### BTOP

File: `btop.theme`

Standardized configuration. Variables set slightly varies between themes. 

References:
- repo: https://github.com/aristocratos/btop
- samples: https://github.com/aristocratos/btop/tree/main/themes

### Chromium

File: `chromium.theme`

Simple RGB color value in plaintext. Varies based on theme palette.

Format:

```
239,241,245
```

### Hyprland

File `hyprland.conf`

Standardized configuration. Variables set is the same between themes. Exception for Kanagawa, which modifies Alacritty window opacity.

References:
- repo: https://github.com/hyprwm/Hyprland
- site: https://wiki.hypr.land/
- config: https://wiki.hypr.land/Configuring/Variables/

### Hyprlock

File: `hyprlock.conf`

Standardized configuration. Variables set is the same between themes.

References:
- repo: https://github.com/hyprwm/hyprlock
- config: https://wiki.hypr.land/Hypr-Ecosystem/hyprlock/

### Yaru

File: `icons.theme`

Pre-defined icon themes. Set based on theme color profile.

Format:

```
Yaru-[theme]
```

Themes:

Theme | Base Color
------|-----------
Yaru-bark | `#757557`
Yaru-blue | `#0171df`
Yaru-magenta | `#af4baf`
Yaru-mate | `#84a154`
Yaru-olive | `#418202`
Yaru-prussiangreen | `#307f7d`
Yaru-purple | `#7462d2`
Yaru-red | `#d4334f`
Yaru-sage | `#637867`
Yaru-viridian | `#048459`
Yaru-wartybrown | `#8f6f49`
Yaru-yellow | `#9b6a01`

References:
- repo: https://github.com/ubuntu/yaru
- gtk: https://github.com/ubuntu/yaru/tree/master/gtk/src
- icons: https://github.com/ubuntu/yaru/tree/master/icons

### Mako

File `maki.ini`

Standardized configuration. Variables set is the same between themes.

References:
- repo: https://github.com/emersion/mako
- sample: https://github.com/emersion/mako/wiki/Example-configuration
- config: https://man.archlinux.org/man/mako.5.en

### Neovim

File: `neovim.lua`

Variable Neovim theme configurations. Need to work out a UX for setting this. Ideas:
- Establish constants for popular Neovim theme configurations and their associated color profiles. Select based on the ColorProfile image metadata processed from the image.
- Allow explicit theme constant to be specified (i.e. - `--neovim-theme kanagawa`)
- Allow explicit `neovim.lua` to be specified (i.e. - `--neovim-theme-file [path-to]/neovim.lua`)

References:
- topic: https://github.com/topics/neovim-colorscheme
- dotfyle: https://dotfyle.com/neovim/colorscheme/trending

### SwayOSD

File: `swayosd.css`

Standardized configuration. Variables set is the same between themes.

References:
- repo: https://github.com/ErikReider/SwayOSD
- sample: https://github.com/ErikReider/SwayOSD/blob/main/data/style/style.scss
- details: https://github.com/ErikReider/SwayOSD/issues/36

### Walker

File `walker.css`

Standardized configuration. Variables set is the same between themes.

References:
- repo: https://github.com/abenz1267/walker

### Waybar

File: `waybar.css`

Standardized configuration. Variables set is the same between themes. Exception for catppuccin-latte, which adds `border` and `accent`.

References:
- repo: https://github.com/Alexays/Waybar
- config: https://github.com/Alexays/Waybar/wiki/Configuration
- styling: https://github.com/Alexays/Waybar/wiki/Styling

## Themes

Default Omarchy themes: https://github.com/basecamp/omarchy/tree/master/themes

### Theme Analysis Overview

#### Theme Consistency Patterns

All Omarchy themes follow a universal structure with consistent file organization and color application patterns:

1. **Core Palette Structure (10-16 colors)**:
   - Background & Foreground base colors (2)
   - 8 Terminal normal colors (ANSI 0-7: black, red, green, yellow, blue, magenta, cyan, white)
   - 8 Terminal bright colors (ANSI 8-15: bright variants of normal colors)
   - Optional: cursor, selection, dim variants

2. **Component Color Requirements**:
   - **Minimal (2-4 colors)**: waybar, hyprland, mako, swayosd, walker
     - Typically: background, foreground, accent/border
   - **Standard (10-16 colors)**: alacritty terminal
     - Full ANSI palette plus UI colors
   - **Extended (20-30+ colors)**: btop system monitor
     - Gradients for metrics (2-3 colors per metric type)
     - Box outlines for different system components

3. **Color Reuse Patterns**:
   - Same colors appear across multiple components
   - Terminal blue frequently becomes UI accent/border color
   - Background/foreground consistently applied across all components
   - Bright colors used for emphasis, alerts, and active states

#### Color Organization Standards

Themes organize their color palettes using several strategies:

1. **Lightness-based Organization**:
   - Dark themes: backgrounds < 0.25 L, foregrounds > 0.70 L
   - Light themes: backgrounds > 0.85 L, foregrounds < 0.30 L
   - Mid-tones used for secondary elements and dim states

2. **Saturation Grouping**:
   - Vibrant colors (S > 0.7) for accents and alerts
   - Normal saturation (0.3-0.7) for standard UI elements
   - Muted colors (S < 0.3) for subtle distinctions
   - Grayscale (S < 0.05) for neutral elements

3. **Hue Distribution**:
   - Terminal colors cover spectrum: red (0°), yellow (60°), green (120°), cyan (180°), blue (240°), magenta (300°)
   - Themes maintain hue relationships for color harmony
   - Some themes use analogous hues, others use complementary

#### Theme Personality Types

1. **Vibrant Themes** (Tokyo Night, Catppuccin):
   - Distinct bright colors with higher saturation/lightness
   - Clear separation between normal and bright variants
   - Rich, saturated accent colors
   - High contrast between elements

2. **Muted Themes** (Nord, Gruvbox, Everforest):
   - Identical or very similar normal/bright colors
   - Earth tones and pastels predominant
   - Lower overall saturation
   - Softer contrast relationships

3. **Minimal Themes** (Matte Black):
   - Limited color variation
   - Heavy use of grayscale
   - Strategic color placement for emphasis
   - Focus on contrast over color diversity

4. **Artistic Themes** (Osaka-Jade, Ristretto, Rose-Pine):
   - Unique, non-traditional color combinations
   - Creative interpretations of standard colors
   - Distinctive hue choices
   - Personality-driven palettes

#### Palette Flexibility

Themes demonstrate remarkable flexibility in color application:

- **Color Multiplexing**: Single colors serve multiple roles across components
- **Adaptive Scaling**: 2-color schemes expand to 30+ through interpolation and variation
- **Context Sensitivity**: Same color appears different based on surrounding colors
- **Gradient Generation**: Base colors extended through lightness/saturation shifts

### [catppuccin-latte](https://github.com/basecamp/omarchy/tree/master/themes/catppuccin-latte)

**Mode**: Light  
**Palette Style**: Pastel/Soft  
**Personality**: Clean, modern, gentle contrast

#### Color Palette

**Core Colors**:
- Background: `#eff1f5` (very light blue-gray)
- Foreground: `#4c4f69` (dark blue-gray)
- Dim Foreground: `#8c8fa1` (medium gray-blue)
- Bright Foreground: `#4c4f69` (same as foreground)
- Cursor: `#dc8a78` (warm pink)
- Cursor Text: `#eff1f5` (matches background)

**Terminal Normal Colors**:
- Black: `#bcc0cc` (light gray)
- Red: `#d20f39` (bright red)
- Green: `#40a02b` (forest green)
- Yellow: `#df8e1d` (amber)
- Blue: `#1e66f5` (vibrant blue)
- Magenta: `#ea76cb` (pink)
- Cyan: `#179299` (teal)
- White: `#5c5f77` (dark gray)

**Terminal Bright Colors**:
- Bright Black: `#acb0be` (medium gray)
- Bright Red: `#d20f39` (same as normal)
- Bright Green: `#40a02b` (same as normal)
- Bright Yellow: `#df8e1d` (same as normal)
- Bright Blue: `#1e66f5` (same as normal)
- Bright Magenta: `#ea76cb` (same as normal)
- Bright Cyan: `#179299` (same as normal)
- Bright White: `#6c6f85` (darker gray)

**Additional Colors**:
- Indexed 16: `#fe640b` (orange)
- Indexed 17: `#dc8a78` (matches cursor)

#### Component Configurations

**Waybar**:
- Foreground: `#4c4f69`
- Background: `#eff1f5`
- Border: `#dce0e8` (crust)
- Accent: `#1e66f5` (blue)

**Hyprland**:
- Active Border: `rgb(1e66f5)` (blue)

**Mako** (notifications):
- Text: `#4c4f69`
- Border: `#1e66f5`
- Background: `#eff1f5`

**SwayOSD**:
- Background: `#eff1f5`
- Border: `#1e66f5`
- Label/Progress: `#4c4f69`

**Walker** (app launcher):
- Selected Text: `#1e66f5`
- Text: `#4c4f69`
- Base: `#eff1f5`
- Border: `#dce0e8`

**Hyprlock**:
- Background: `rgba(239,241,245,1.0)`
- Inner: `rgba(239,241,245,0.8)`
- Outer: `rgba(30,102,245,1.0)`
- Font: `rgba(76,79,105,1.0)`
- Check: `rgba(4,165,229,1.0)`

**Chromium**: `239,241,245`

**Icon Theme**: `Yaru-blue`

**Neovim**: `catppuccin-latte` colorscheme

#### Theme Characteristics

- **Unique Features**: Light theme with soft pastel colors, includes border and accent variables in waybar
- **Color Count**: ~20 unique colors across all components
- **Contrast Strategy**: Gentle contrast suitable for daylight use
- **Color Reuse**: Blue (`#1e66f5`) used consistently as accent across components

### [catppuccin](https://github.com/basecamp/omarchy/tree/master/themes/catppuccin)

**Mode**: Dark  
**Palette Style**: Pastel/Soft  
**Personality**: Modern, gentle, soothing

#### Color Palette

**Core Colors**:
- Background: `#24273a` (deep blue-black)
- Foreground: `#cad3f5` (light blue-gray)
- Dim Foreground: `#8087a2` (muted blue-gray)
- Bright Foreground: `#cad3f5` (same as foreground)
- Cursor: `#f4dbd6` (warm beige)
- Cursor Text: `#24273a` (matches background)
- Vi Mode Cursor: `#b7bdf8` (lavender)

**Terminal Normal Colors**:
- Black: `#494d64` (dark gray)
- Red: `#ed8796` (coral)
- Green: `#a6da95` (mint)
- Yellow: `#eed49f` (cream)
- Blue: `#8aadf4` (sky blue)
- Magenta: `#f5bde6` (pink)
- Cyan: `#8bd5ca` (aqua)
- White: `#b8c0e0` (light gray)

**Terminal Bright Colors**:
- Bright Black: `#5b6078` (medium gray)
- Bright Red: `#ed8796` (same as normal)
- Bright Green: `#a6da95` (same as normal)
- Bright Yellow: `#eed49f` (same as normal)
- Bright Blue: `#8aadf4` (same as normal)
- Bright Magenta: `#f5bde6` (same as normal)
- Bright Cyan: `#8bd5ca` (same as normal)
- Bright White: `#a5adcb` (lighter gray)

**Additional Colors**:
- Indexed 16: `#f5a97f` (peach)
- Indexed 17: `#f4dbd6` (matches cursor)

#### Theme Characteristics

- **Unique Features**: Dark counterpart to latte, soft pastel palette throughout
- **Color Count**: ~18 unique colors
- **Contrast Strategy**: Moderate contrast with soothing colors
- **Color Reuse**: Consistent use of pastel tones across components

### [everforest](https://github.com/basecamp/omarchy/tree/master/themes/everforest)

**Mode**: Dark  
**Palette Style**: Earth-tone/Natural  
**Personality**: Calm, organic, forest-inspired

#### Color Palette

**Core Colors**:
- Background: `#2d353b` (dark slate)
- Foreground: `#d3c6aa` (warm beige)

**Terminal Colors** (Normal and Bright are identical):
- Black: `#475258` (dark gray-green)
- Red: `#e67e80` (coral)
- Green: `#a7c080` (sage)
- Yellow: `#dbbc7f` (wheat)
- Blue: `#7fbbb3` (teal)
- Magenta: `#d699b6` (mauve)
- Cyan: `#83c092` (mint)
- White: `#d3c6aa` (warm beige)

#### Theme Characteristics

- **Unique Features**: Normal and bright colors are identical, nature-inspired palette
- **Color Count**: ~9 unique colors (minimal variation)
- **Contrast Strategy**: Soft, comfortable contrast
- **Color Reuse**: Earth tones create cohesive natural aesthetic

### [gruvbox](https://github.com/basecamp/omarchy/tree/master/themes/gruvbox)

**Mode**: Dark  
**Palette Style**: Retro/Warm  
**Personality**: Vintage, warm, comfortable

#### Color Palette

**Core Colors**:
- Background: `#282828` (dark brown-gray)
- Foreground: `#d4be98` (warm tan)

**Terminal Colors** (Normal and Bright are identical):
- Black: `#3c3836` (dark brown)
- Red: `#ea6962` (warm red)
- Green: `#a9b665` (olive)
- Yellow: `#d8a657` (gold)
- Blue: `#7daea3` (muted teal)
- Magenta: `#d3869b` (dusty rose)
- Cyan: `#89b482` (sage)
- White: `#d4be98` (warm tan)

**Component Configurations**:

**Waybar**:
- Foreground: `#d4be98`
- Background: `#282828`

#### Theme Characteristics

- **Unique Features**: Identical normal/bright colors, retro aesthetic
- **Color Count**: ~9 unique colors
- **Contrast Strategy**: Warm, comfortable contrast
- **Color Reuse**: Minimal palette with strategic reuse

### [kanagawa](https://github.com/basecamp/omarchy/tree/master/themes/kanagawa)

**Mode**: Dark  
**Palette Style**: Japanese-inspired/Muted  
**Personality**: Elegant, subdued, artistic

#### Color Palette

**Core Colors**:
- Background: `#1f1f28` (deep charcoal)
- Foreground: `#dcd7ba` (warm white)
- Selection Background: `#2d4f67` (dark blue)
- Selection Foreground: `#c8c093` (beige)

**Terminal Normal Colors**:
- Black: `#090618` (near black)
- Red: `#c34043` (crimson)
- Green: `#76946a` (moss)
- Yellow: `#c0a36e` (amber)
- Blue: `#7e9cd8` (periwinkle)
- Magenta: `#957fb8` (lavender)
- Cyan: `#6a9589` (jade)
- White: `#c8c093` (beige)

**Terminal Bright Colors**:
- Bright Black: `#727169` (gray)
- Bright Red: `#e82424` (bright red)
- Bright Green: `#98bb6c` (lime)
- Bright Yellow: `#e6c384` (light amber)
- Bright Blue: `#7fb4ca` (sky)
- Bright Magenta: `#938aa9` (muted purple)
- Bright Cyan: `#7aa89f` (seafoam)
- Bright White: `#dcd7ba` (warm white)

**Additional Colors**:
- Indexed 16: `#ffa066` (orange)
- Indexed 17: `#ff5d62` (coral)

#### Theme Characteristics

- **Unique Features**: Japanese aesthetic, distinct bright variants, modifies Alacritty opacity
- **Color Count**: ~20 unique colors
- **Contrast Strategy**: Subtle, elegant contrast
- **Color Reuse**: Carefully curated palette with Japanese influence

### [matte-black](https://github.com/basecamp/omarchy/tree/master/themes/matte-black)

**Mode**: Dark  
**Palette Style**: Minimal/Monochrome  
**Personality**: Ultra-minimal, stark, focused

#### Color Palette

**Core Colors**:
- Background: `#121212` (deep black)
- Foreground: `#bebebe` (light gray)
- Dim Foreground: `#8a8a8d` (muted gray)
- Cursor: `#eaeaea` (off-white)
- Cursor Text: `#121212` (deep black)

**Terminal Normal Colors**:
- Black: `#333333` (dark gray)
- Red: `#D35F5F` (muted red)
- Green: `#FFC107` (amber)
- Yellow: `#b91c1c` (dark red)
- Blue: `#e68e0d` (orange)
- Magenta: `#D35F5F` (same as red)
- Cyan: `#bebebe` (light gray)
- White: `#bebebe` (light gray)

**Terminal Bright Colors**:
- Bright Black: `#8a8a8d` (gray)
- Bright Red: `#B91C1C` (crimson)
- Bright Green: `#FFC107` (amber)
- Bright Yellow: `#b90a0a` (dark red)
- Bright Blue: `#f59e0b` (bright orange)
- Bright Magenta: `#B91C1C` (crimson)
- Bright Cyan: `#eaeaea` (off-white)
- Bright White: `#ffffff` (pure white)

#### Theme Characteristics

- **Unique Features**: Extremely limited color palette, heavy grayscale usage
- **Color Count**: ~12 unique colors
- **Contrast Strategy**: High contrast monochrome with selective color
- **Color Reuse**: Strategic use of grays and minimal accent colors

### [nord](https://github.com/basecamp/omarchy/tree/master/themes/nord)

**Mode**: Dark  
**Palette Style**: Arctic/Cool  
**Personality**: Cool, professional, nordic

#### Color Palette

**Core Colors**:
- Background: `#2e3440` (polar night)
- Foreground: `#d8dee9` (snow storm)
- Dim Foreground: `#a5abb6` (muted gray)
- Cursor: `#d8dee9` (snow storm)
- Cursor Text: `#2e3440` (polar night)
- Selection Background: `#4c566a` (dark gray)

**Terminal Normal Colors**:
- Black: `#3b4252` (dark gray)
- Red: `#bf616a` (aurora red)
- Green: `#a3be8c` (aurora green)
- Yellow: `#ebcb8b` (aurora yellow)
- Blue: `#81a1c1` (frost blue)
- Magenta: `#b48ead` (aurora purple)
- Cyan: `#88c0d0` (frost cyan)
- White: `#e5e9f0` (snow)

**Terminal Bright Colors**:
- Bright Black: `#4c566a` (gray)
- Bright Red: `#bf616a` (same as normal)
- Bright Green: `#a3be8c` (same as normal)
- Bright Yellow: `#ebcb8b` (same as normal)
- Bright Blue: `#81a1c1` (same as normal)
- Bright Magenta: `#b48ead` (same as normal)
- Bright Cyan: `#8fbcbb` (lighter frost)
- Bright White: `#eceff4` (bright snow)

**Terminal Dim Colors**:
- Dim Black: `#373e4d`
- Dim Red: `#94545d`
- Dim Green: `#809575`
- Dim Yellow: `#b29e75`
- Dim Blue: `#68809a`
- Dim Magenta: `#8c738c`
- Dim Cyan: `#6d96a5`
- Dim White: `#aeb3bb`

#### Theme Characteristics

- **Unique Features**: Includes third "dim" palette, arctic-inspired colors
- **Color Count**: ~24 unique colors (with dim variants)
- **Contrast Strategy**: Moderate contrast with cool tones
- **Color Reuse**: Consistent frost and aurora color families

### [osaka-jade](https://github.com/basecamp/omarchy/tree/master/themes/osaka-jade)

**Mode**: Dark  
**Palette Style**: Jade/Oriental  
**Personality**: Unique, artistic, jade-inspired

#### Color Palette

**Core Colors**:
- Background: `#111c18` (deep jade black)
- Foreground: `#C1C497` (pale gold)
- Cursor: `#D7C995` (gold)
- Cursor Text: `#000000` (pure black)

**Terminal Normal Colors**:
- Black: `#23372B` (dark jade)
- Red: `#FF5345` (coral red)
- Green: `#549e6a` (jade green)
- Yellow: `#459451` (forest)
- Blue: `#509475` (teal)
- Magenta: `#D2689C` (rose)
- Cyan: `#2DD5B7` (bright jade)
- White: `#F6F5DD` (cream)

**Terminal Bright Colors**:
- Bright Black: `#53685B` (gray-green)
- Bright Red: `#db9f9c` (pale coral)
- Bright Green: `#143614` (deep forest)
- Bright Yellow: `#E5C736` (gold)
- Bright Blue: `#ACD4CF` (pale teal)
- Bright Magenta: `#75bbb3` (seafoam)
- Bright Cyan: `#8CD3CB` (mint)
- Bright White: `#9eebb3` (pale jade)

#### Theme Characteristics

- **Unique Features**: Jade and gold color scheme, unusual bright color choices
- **Color Count**: ~18 unique colors
- **Contrast Strategy**: Artistic contrast with jade emphasis
- **Color Reuse**: Jade greens and golds throughout

### [ristretto](https://github.com/basecamp/omarchy/tree/master/themes/ristretto)

**Mode**: Dark  
**Palette Style**: Coffee/Warm  
**Personality**: Warm, cozy, coffee-inspired

#### Color Palette

**Core Colors**:
- Background: `#2c2525` (dark coffee)
- Foreground: `#e6d9db` (cream)
- Cursor: `#c3b7b8` (light brown)
- Selection Background: `#403e41` (medium brown)

**Terminal Normal Colors**:
- Black: `#72696a` (gray-brown)
- Red: `#fd6883` (pink-red)
- Green: `#adda78` (lime)
- Yellow: `#f9cc6c` (butter)
- Blue: `#f38d70` (salmon)
- Magenta: `#a8a9eb` (lavender)
- Cyan: `#85dacc` (mint)
- White: `#e6d9db` (cream)

**Terminal Bright Colors**:
- Bright Black: `#948a8b` (lighter gray-brown)
- Bright Red: `#ff8297` (bright pink)
- Bright Green: `#c8e292` (pale lime)
- Bright Yellow: `#fcd675` (pale butter)
- Bright Blue: `#f8a788` (peach)
- Bright Magenta: `#bebffd` (bright lavender)
- Bright Cyan: `#9bf1e1` (bright mint)
- Bright White: `#f1e5e7` (light cream)

#### Theme Characteristics

- **Unique Features**: Coffee-inspired warm palette, unique color assignments
- **Color Count**: ~17 unique colors
- **Contrast Strategy**: Warm, comfortable contrast
- **Color Reuse**: Warm browns and creams with pastel accents

### [rose-pine](https://github.com/basecamp/omarchy/tree/master/themes/rose-pine)

**Mode**: Light  
**Palette Style**: Floral/Soft  
**Personality**: Gentle, romantic, floral

#### Color Palette

**Core Colors**:
- Background: `#faf4ed` (soft cream)
- Foreground: `#575279` (muted purple)
- Dim Text: `#797593` (lighter purple)
- Bright Text: `#575279` (same as foreground)
- Cursor: `#cecacd` (light gray)
- Selection Background: `#dfdad9` (pale pink)
- Search Background: `#f2e9e1` (lighter cream)

**Terminal Colors** (simplified palette):
- Red: `#b4637a` (rose)
- Green: `#286983` (deep teal)
- Yellow: `#ea9d34` (amber)
- Blue: `#56949f` (slate blue)
- Magenta: `#907aa9` (lavender)
- Cyan: `#d7827e` (coral)

**Component Configurations**:

**Waybar**:
- Foreground: `#575279`
- Background: `#faf4ed`

#### Theme Characteristics

- **Unique Features**: Light theme with floral inspiration, simplified color set
- **Color Count**: ~12 unique colors
- **Contrast Strategy**: Soft, romantic contrast
- **Color Reuse**: Muted purples and warm creams

### [tokyo-night](https://github.com/basecamp/omarchy/tree/master/themes/tokyo-night)

**Mode**: Dark  
**Palette Style**: Neon/Vibrant  
**Personality**: Modern, vibrant, city-night inspired

#### Color Palette

**Core Colors**:
- Background: `#1a1b26` (deep navy)
- Foreground: `#a9b1d6` (pale purple)
- Selection Background: `#7aa2f7` (blue)

**Terminal Normal Colors**:
- Black: `#32344a` (dark purple)
- Red: `#f7768e` (pink-red)
- Green: `#9ece6a` (lime)
- Yellow: `#e0af68` (gold)
- Blue: `#7aa2f7` (sky blue)
- Magenta: `#ad8ee6` (purple)
- Cyan: `#449dab` (teal)
- White: `#787c99` (gray)

**Terminal Bright Colors**:
- Bright Black: `#444b6a` (lighter purple)
- Bright Red: `#ff7a93` (bright pink)
- Bright Green: `#b9f27c` (bright lime)
- Bright Yellow: `#ff9e64` (orange)
- Bright Blue: `#7da6ff` (bright sky)
- Bright Magenta: `#bb9af7` (bright purple)
- Bright Cyan: `#0db9d7` (bright cyan)
- Bright White: `#acb0d0` (light purple)

#### Theme Characteristics

- **Unique Features**: Vibrant neon colors, distinct bright variants, city-night aesthetic
- **Color Count**: ~17 unique colors
- **Contrast Strategy**: High contrast with vibrant accents
- **Color Reuse**: Blues and purples create cohesive night theme

## pkg/processor Optimization Analysis

Based on the comprehensive analysis of Omarchy themes, the current pkg/processor package requires significant refactoring to better support flexible theme generation. The existing approach of pre-categorizing colors into 27 specific theme roles (background, foreground, normal_red, etc.) is too rigid and forces premature decisions about color usage.

### Current Limitations

1. **Premature Categorization**: The processor assigns colors to specific roles before understanding component requirements
2. **Rigid Category System**: 27 fixed categories don't match the flexible nature of theme components
3. **Missing Tiered Requirements**: Components need different color counts (2-30+) which isn't reflected
4. **Lost Color Relationships**: Focus on categories loses important color relationships and harmonies

### Recommended Refactoring

#### 1. Focus on Color Extraction and Organization

Instead of semantic categories, organize colors by intrinsic properties:

```go
type ColorProfile struct {
    // Image analysis metadata
    Mode            ThemeMode
    ColorScheme     chromatic.ColorScheme  
    IsGrayscale     bool
    IsMonochromatic bool
    
    // Organized color data
    ColorPool       ColorPool
    Statistics      ColorStatistics
}

type ColorPool struct {
    // Frequency-based extraction
    DominantColors  []WeightedColor  // Top N by frequency
    
    // Characteristic-based organization
    ByLightness     LightnessGroups  // dark/mid/light buckets
    BySaturation    SaturationGroups // vibrant/normal/muted/gray
    ByHue           HueFamilies      // 12 hue sectors (30° each)
    
    // Relationship tracking
    ContrastPairs   []ColorPair      // Pre-calculated high contrast pairs
    HarmonyGroups   []ColorGroup     // Colors that work well together
}
```

#### 2. Rich Metadata Collection

```go
type LightnessGroups struct {
    Dark   []color.RGBA  // L < 0.33
    Mid    []color.RGBA  // 0.33 <= L < 0.66  
    Light  []color.RGBA  // L >= 0.66
}

type SaturationGroups struct {
    Vibrant []color.RGBA  // S >= 0.7
    Normal  []color.RGBA  // 0.3 <= S < 0.7
    Muted   []color.RGBA  // 0.05 <= S < 0.3
    Gray    []color.RGBA  // S < 0.05
}

type HueFamilies map[float64][]color.RGBA // Hue sector -> colors
```

#### 3. Enhanced Statistics

```go
type ColorStatistics struct {
    // Distribution metrics
    HueHistogram       []float64  // 360 buckets
    LightnessHistogram []float64  // 100 buckets
    SaturationDistribution map[string]float64
    
    // Dominant characteristics
    PrimaryHue         float64
    SecondaryHue       float64
    AverageLightness   float64
    AverageSaturation  float64
    
    // Coverage metrics
    HueVariance        float64
    ChromaticDiversity float64
    ContrastRange      float64
}
```

### Separation of Concerns

The refactored architecture should clearly separate responsibilities:

1. **pkg/processor**: Extract and organize colors by characteristics
   - Frequency-based extraction
   - Group by lightness, saturation, hue
   - Calculate relationships and statistics
   - No semantic role assignment

2. **pkg/palette**: Map colors to theme components
   - Take ColorPool from processor
   - Apply theme generation strategies
   - Handle component-specific requirements
   - Create semantic mappings

3. **pkg/theme**: Generate configuration files
   - Use palette selections
   - Apply component templates
   - Handle format-specific requirements
   - Write theme files

### Benefits of This Approach

1. **Flexibility**: Colors aren't forced into roles prematurely
2. **Rich Selection**: Multiple candidates available for each component need
3. **Preserves Relationships**: Color harmonies and contrasts maintained
4. **Adaptive**: Can generate minimal (2 colors) to extended (30+ colors) schemes
5. **Component-Aware**: Each component can have custom selection logic
6. **Theme Style Support**: Can generate vibrant, muted, minimal, or artistic themes

### Implementation Priority

1. Refactor ColorProfile to remove categories
2. Implement characteristic-based grouping
3. Add relationship tracking (contrast, harmony)
4. Create rich statistics collection
5. Move semantic mapping to pkg/palette

This refactoring aligns with the discovered theme patterns and provides the flexibility needed to generate diverse theme styles from image colors.
