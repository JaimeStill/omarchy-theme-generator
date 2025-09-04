# Omarchy Theme Generator

Generate beautiful, cohesive terminal themes from any image using intelligent color extraction and color theory principles.

## Features

- 🎨 **Intelligent Color Extraction** - Multi-strategy system adapts to image characteristics
- 🎯 **Purpose-Driven Colors** - Organizes colors by their role (background, foreground, accents)
- ♿ **WCAG Compliant** - Ensures text readability with proper contrast ratios  
- 🎭 **Light/Dark Mode** - Automatic detection with manual override options
- 🖼️ **Edge Case Handling** - Gracefully handles grayscale, monotone, and monochromatic images
- ⚡ **High Performance** - Processes 4K images in under 2 seconds

## Installation

```bash
go install github.com/JaimeStill/omarchy-theme-generator/cmd/omarchy-theme-gen@latest
```

Or clone and run directly:

```bash
git clone https://github.com/JaimeStill/omarchy-theme-generator
cd omarchy-theme-generator
go run cmd/omarchy-theme-gen/main.go
```

## Usage

### CLI Usage

```bash
omarchy-theme-gen generate --image photo.jpg [options]
```

Generate themes directly from command line:

1. **Generate theme from image**
   ```bash
   omarchy-theme-gen generate --image sunset.jpg --mode dark --name "sunset-dark"
   ```

2. **Adjust color scheme after generation**
   ```bash
   omarchy-theme-gen set-scheme sunset-dark --scheme complementary
   ```

3. **Switch between light and dark modes**
   ```bash
   omarchy-theme-gen set-mode sunset-dark --mode light
   ```

4. **Clone existing theme**
   ```bash
   omarchy-theme-gen clone sunset-dark sunset-variant
   ```

### Customization Options

```bash
# Force dark mode
omarchy-theme-gen generate --image wallpaper.jpg --mode dark

# Override primary color
omarchy-theme-gen generate --image wallpaper.jpg --primary "#e94560"

# Apply color scheme
omarchy-theme-gen generate --image wallpaper.jpg --scheme complementary

# Refine existing theme
omarchy-theme-gen set-scheme my-theme --scheme triadic
omarchy-theme-gen set-mode my-theme --mode light
```

## Generated Theme Structure

```
my-theme/
├── alacritty.toml      # Terminal emulator
├── btop.theme          # System monitor
├── hyprland.conf       # Window manager
├── hyprlock.conf       # Lock screen
├── mako.ini            # Notifications
├── neovim.lua          # Editor colorscheme
├── waybar.css          # Status bar
├── walker.css          # App launcher
├── swayosd.css         # On-screen display
├── theme-gen.json      # Theme metadata for refinement
├── backgrounds/        # Wallpapers
│   └── wallpaper.jpg
└── light.mode          # (if light theme)
```

## Generated Theme Integration

Generated themes integrate directly with Omarchy's theme system and include a `theme-gen.json` file containing extraction metadata and user preferences for easy refinement.

## Architecture

The generator uses a layered architecture with clear dependencies:

### Core Packages

#### Foundation Layer
- **pkg/formats** - Color space representations and conversions (RGBA, HSLA, LAB, XYZ)
- **pkg/chromatic** - Foundational color theory and mathematical calculations  
- **pkg/settings** - Flat configuration structure with Viper integration
- **pkg/loader** - Image I/O with validation and format support
- **pkg/errors** - Error handling utilities

#### Analysis Layer  
- **pkg/analysis** - High-level color profiles and comprehensive analysis

#### Processing Layer
- **pkg/extractor** - Color extraction orchestration (contains embedded strategies)
- **pkg/strategies** - Pluggable extraction strategies (pending extraction)

#### Generation Layer
- **pkg/theme** - Theme file generation from templates (pending implementation)

### Extraction Pipeline

1. **Analysis** - Detect image characteristics and color profile
2. **Strategy Selection** - Choose optimal extraction strategy
3. **Extraction** - Extract colors using selected strategy
4. **Role Assignment** - Categorize colors by purpose
5. **Calculation** - Calculate missing colors using color theory
6. **Validation** - Ensure WCAG compliance
7. **Generation** - Create theme configuration files

## Supported Formats

- **Input**: JPEG, PNG images
- **Output**: All Omarchy configuration formats
- **Color Spaces**: RGB, HSL with automatic conversion
- **Palette Strategies**:
  - Monochromatic (single hue variations)
  - Complementary (opposite colors)
  - Triadic (three-color harmony)
  - Analogous (adjacent colors)
  - Tetradic (four-color schemes)

## Requirements

- Go 1.25+
- No external dependencies (pure Go implementation)

## Development

```bash
# Run tests
go test ./tests -v
go vet ./...

# Format code
go fmt ./...
```

See [docs/](docs/) for technical documentation and development guidelines.

## Philosophy

This project follows [Intelligent Development](docs/development-methodology.md) principles:

- Precise technical language
- Immediate validation through execution tests
- User-driven development with AI assistance
- Knowledge transfer as primary output

## License

MIT

## Acknowledgments

Built for the [Omarchy](https://omarchy.org) desktop environment by [DHH](https://github.com/DHH).
