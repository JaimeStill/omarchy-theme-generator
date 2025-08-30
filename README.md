# Omarchy Theme Generator

Generate beautiful, cohesive terminal themes from any image using intelligent color extraction and color theory principles.

## Features

- 🎨 **Smart Color Extraction** - Automatically identifies dominant colors using octree quantization
- 🎭 **Light/Dark Mode** - Auto-detects or manually specify theme brightness
- 🎯 **Color Theory** - Generates harmonious palettes (monochromatic, complementary, triadic, analogous)
- ♿ **WCAG Compliant** - Ensures text readability with AA contrast ratios
- 🖼️ **Wallpaper Included** - Source image automatically included in theme package
- ⚡ **Fast** - Processes 4K images in under 2 seconds

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

### Interactive TUI

```bash
omarchy-theme-gen
```

Navigate the interactive interface to:

1. Select your source image
2. Choose theme mode (light/dark/auto)
3. Optionally override colors
4. Select palette strategy
5. Preview and adjust
6. Export complete theme

### Command Line

```bash
# Auto-generate from image
omarchy-theme-gen sunset.jpg

# Specify theme mode
omarchy-theme-gen sunset.jpg --mode dark

# Override primary color
omarchy-theme-gen sunset.jpg --primary "#ff6b35"

# Full control
omarchy-theme-gen sunset.jpg \
  --mode dark \
  --primary "#ff6b35" \
  --background "#1a1a1a" \
  --foreground "#e0e0e0"
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
├── backgrounds/        # Wallpapers
│   └── wallpaper.jpg
└── light.mode          # (if light theme)
```

## Installing Generated Themes

```bash
# Copy to Omarchy themes directory
cp -r my-theme ~/.config/omarchy/themes/

# Or use Omarchy's theme installer
omarchy-theme-install ./my-theme
```

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
go vet ./...

# Run examples
go run tests/test-load-image/main.go sample.jpg

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
