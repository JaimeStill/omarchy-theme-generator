# Omarchy Theme Generator

> [!IMPORTANT]
> My day job took over and I lost bandwidth to engage with this project. These were some fun concepts to work through and learn from. I learned a lot working on this.

Color extraction and analysis system for generating Omarchy themes from images using characteristic-based color organization.

## What This Tool Is

The Omarchy Theme Generator is a Go-based color analysis system that extracts and organizes colors from images by their intrinsic properties (lightness, saturation, hue). It provides a foundation for generating themes ranging from minimal to extended configurations for the Omarchy desktop environment.

## Current Implementation Status

### ✅ Complete Foundation
The core processing pipeline is fully implemented and tested:

- **Color Extraction**: ColorCluster-based system with frequency weighting
- **Image Processing**: JPEG, PNG support with validation and format detection
- **Color Analysis**: HSLA, LAB, XYZ color space conversions and calculations
- **Performance**: <500ms average processing time, <50MB memory usage
- **Color Theory**: Similarity detection, contrast calculations, accessibility metrics

### 🔄 In Development
- **Theme Generation**: Semantic color mapping and component generation
- **CLI Interface**: Command-line application for end users
- **Palette Strategies**: Color scheme application and theme optimization

## What You Can Currently Do

### Process Images Through the Analysis Pipeline
Test the color extraction system directly:

```bash
# Run comprehensive unit tests
go test ./tests/... -v

# Run performance benchmarks
go test ./tests/benchmarks -bench=. -benchmem

# Generate image analysis documentation
go run tools/analyze-images/main.go

# Run performance validation
go run tools/performance-test/main.go
```

### Analyze Color Processing Results
The processing pipeline currently extracts:

- **ColorCluster structures** with UI-relevant metadata (weight, lightness, characteristics)
- **Theme mode detection** (Light/Dark based on weighted luminance)
- **Color characteristics** (neutral, vibrant, muted classifications)
- **Processing statistics** and performance metrics

See **[tests/images/README.md](tests/images/README.md)** for examples of current processing output.

## Architecture Status

### Foundation Layer (Complete ✅)
- **pkg/formats** - Color space conversions (RGBA, HSLA, LAB, XYZ)
- **pkg/chromatic** - Color theory algorithms and calculations
- **pkg/settings** - Configuration management with Viper integration
- **pkg/loader** - Image loading with format validation
- **pkg/errors** - Comprehensive error handling with sentinel errors

### Processing Layer (Complete ✅)
- **pkg/processor** - ColorCluster-based extraction with characteristic analysis

### Generation Layer (Planned 🔄)
- **pkg/palette** - Semantic color mapping (not yet implemented)
- **pkg/theme** - Component-specific configuration generation (not yet implemented)

### Application Layer (Planned 🔄)
- **cmd/omarchy-theme-gen** - CLI interface (not yet implemented)

## Performance Characteristics

Validated performance on real-world images:
- **Small images** (~2MP): ~84ms processing, ~17MB memory
- **Large images** (~15MP): ~533ms processing, ~118MB memory
- **Success rate**: 100% across diverse image types and formats
- **Meets targets**: <2s processing, <100MB memory for all tested images

## Development and Testing

### Requirements
- Go 1.25+
- No external dependencies (pure Go implementation)
- Linux/Unix environment (developed and tested on Linux)

### Development Commands
```bash
# Run all tests with verbose output
go test ./tests/... -v

# Run performance benchmarks
go test ./tests/benchmarks -bench=. -benchmem

# Validate code quality
go vet ./...
go fmt ./...

# Generate image analysis documentation
go run tools/analyze-images/main.go

# Run comprehensive performance tests
go run tools/performance-test/main.go
```

### Test Infrastructure
- **[tests/](tests/)** - Comprehensive unit test suite with diagnostic logging
- **[tests/benchmarks/](tests/benchmarks/)** - Performance validation against targets
- **[tools/](tools/)** - Analysis and validation tools for development
- **[tests/images/](tests/images/)** - Test image collection with analysis results

## Current ColorProfile Output

The processor generates ColorProfile structures containing:

```go
type ColorProfile struct {
    Mode       ThemeMode      // Light or Dark theme recommendation
    Colors     []ColorCluster // Distinct colors sorted by weight
    HasColor   bool          // Whether image contains significant color
    ColorCount int           // Number of distinct colors found
}

type ColorCluster struct {
    color.RGBA                // Representative color
    Weight      float64       // Combined frequency weight (0.0-1.0)
    Lightness   float64       // HSL lightness for UI decisions
    Saturation  float64       // HSL saturation
    Hue         float64       // Hue in degrees (0-360)
    IsNeutral   bool         // Grayscale or very low saturation
    IsDark      bool         // Low lightness (< 0.3)
    IsLight     bool         // High lightness (> 0.7)
    IsMuted     bool         // Low saturation
    IsVibrant   bool         // High saturation
}
```

## Future CLI Design (Planned)

Once generation layers are implemented:

```bash
# Generate Omarchy theme from image
omarchy-theme-gen generate --image wallpaper.jpg --name "my-theme"

# Clone and modify existing themes
omarchy-theme-gen clone my-theme my-variant --scheme complementary

# Set theme properties
omarchy-theme-gen set-mode my-theme --mode dark
omarchy-theme-gen set-scheme my-theme --scheme triadic
```

## Development Approach

This project follows Intelligent Development principles:
- **[docs/development-methodology.md](docs/development-methodology.md)** - Development practices
- **[docs/testing-strategy.md](docs/testing-strategy.md)** - Testing standards
- **[CLAUDE.md](CLAUDE.md)** - Project context and AI-assisted development guidelines
- **[PROJECT.md](PROJECT.md)** - Detailed roadmap and session logs

## Project Status and Roadmap

For detailed development status, architectural decisions, and future phases:
- **[PROJECT.md](PROJECT.md)** - Complete project roadmap and status tracking
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - Architectural design and technical decisions
- **[THEMES.md](THEMES.md)** - Analysis of Omarchy theme requirements and structure

## Acknowledgments

Built for the [Omarchy](https://omarchy.org) desktop environment.
