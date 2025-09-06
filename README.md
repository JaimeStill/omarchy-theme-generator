# Omarchy Theme Generator

> [!IMPORTANT]
> This project is still early in development and the README describes the intended API for this tool. Stay tuned, it's coming together!

Generate beautiful, cohesive terminal themes from any image using intelligent color extraction and color theory principles.

## Features

- 🎨 **Intelligent Color Extraction** - Optimized single-pass frequency analysis
- 🎯 **Purpose-Driven Colors** - Organizes colors by their role (background, foreground, accents)
- ♿ **WCAG Compliant** - Ensures text readability with proper contrast ratios  
- 🎭 **Light/Dark Mode** - Automatic detection with manual override options
- 🖼️ **Edge Case Handling** - Gracefully handles grayscale, monotone, and monochromatic images
- ⚡ **High Performance** - Processes 4K images in under 2s

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

Following comprehensive refactoring, the generator uses a simplified layered architecture optimized for performance:

### Core Packages

#### Foundation Layer (✅ Complete)
- **pkg/formats** - Color space representations and conversions (RGBA, HSLA, LAB, XYZ)
- **pkg/chromatic** - Color theory foundation and harmony calculations  
- **pkg/settings** - System configuration with Viper integration
- **pkg/loader** - Image I/O with validation and format support

#### Processing Layer (✅ Complete) 
- **pkg/processor** - Unified image processing and analysis with single-pass pipeline

#### Generation Layer (🔄 Future)
- **pkg/palette** - Complete theme palette generation using color theory algorithms
- **pkg/theme** - Theme file generation from templates

### Processing Pipeline

**Unified Single-Pass Processing:**
1. **Image Loading** - Load and validate image format
2. **Color Extraction** - Frequency-based analysis for optimal performance  
3. **Profile Analysis** - Detect grayscale, monochromatic, and color schemes
4. **Role Assignment** - Map colors to background/foreground/primary/secondary/accent
5. **Validation** - Ensure WCAG compliance with automatic fallbacks
6. **Result** - Complete ColorProfile with embedded ImageColors and metadata

## Supported Formats

- **Input**: JPEG, PNG images
- **Output**: All Omarchy configuration formats
- **Color Spaces**: RGB, HSL with automatic conversion
- **Color Schemes**: Automatic detection and classification
  - Monochromatic (single hue variations)
  - Complementary (opposite colors)  
  - Triadic (three-color harmony)
  - Analogous (adjacent colors)
  - Complex multi-color schemes

## Requirements

- Go 1.25+
- No external dependencies (pure Go implementation)

## Development

```bash
# Run all tests
go test ./tests/... -v

# Run package-specific tests  
go test ./tests/formats -v
go test ./tests/processor -v
go test ./tests/chromatic -v
go test ./tests/settings -v
go test ./tests/loader -v

# Validate code
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
