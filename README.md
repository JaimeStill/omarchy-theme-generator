# Omarchy Theme Generator

Generate beautiful, cohesive terminal themes from any image using intelligent color extraction and color theory principles.

## Features

- 🎨 **Smart Color Extraction** - Automatically identifies dominant colors using octree quantization
- 🎭 **Light/Dark Mode** - Auto-detects or manually specify theme brightness
- 🎯 **Color Theory** - Generates harmonious color schemes (monochromatic, analogous, complementary, split-complementary, triadic, tetradic, square)
- ♿ **WCAG Compliant** - Ensures text readability with AA contrast ratios
- 🖼️ **Wallpaper Included** - Source image automatically included in theme package
- ⚡ **Reliable** - Clean implementation with solid error handling

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

### Advanced Usage

```bash
# Generate with specific options
omarchy-theme-gen generate --image sunset.jpg --scheme monochromatic --mode dark

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
├── theme-gen.json      # Theme metadata for refinement
├── backgrounds/        # Wallpapers
│   └── wallpaper.jpg
└── light.mode          # (if light theme)
```

## Generated Theme Integration

Generated themes integrate directly with Omarchy's theme system and include a `theme-gen.json` file for refinement:

```json
{
  "version": "1.0.0",
  "source_image": "backgrounds/sunset.jpg",
  "generation": {
    "mode": "dark",
    "scheme": "complementary",
    "primary": "#88c0d0ff",
    "background": "#2e3440ff",
    "foreground": "#eceff4ff"
  }
}
```

Themes are automatically available in Omarchy after generation.

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
go run tests/test-extract/main.go sample.jpg

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
