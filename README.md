# Omarchy Theme Generator

> [!IMPORTANT]
> This project is still early in development. The core processing pipeline is complete, but the CLI application and theme generation features are not yet implemented.

Color extraction and analysis engine for generating Omarchy themes from images using characteristic-based color organization and flexible palette mapping.

## What This Tool Is

The Omarchy Theme Generator is a Go-based color analysis system that extracts and organizes colors from images by their intrinsic properties (lightness, saturation, hue). It provides a rich pool of colors that can be flexibly mapped to generate themes ranging from minimal (2-4 colors) to extended (30+ colors) configurations for all Omarchy components.

## What You Can Currently Do

### Test and Validate Color Processing
The test infrastructure provides comprehensive coverage of the processing pipeline:
- **[tests/](tests/)** - Unit test suite with diagnostic logging for all packages
- **[tests/integration/](tests/integration/)** - End-to-end pipeline validation
- **[tests/benchmarks/](tests/benchmarks/)** - Performance benchmarking suite

### Analyze Images with Development Tools
Analysis and testing tools for development and validation:
- **[tools/analyze-images/](tools/analyze-images/)** - Image analysis and documentation generation
- **[tools/performance-test/](tools/performance-test/)** - Performance testing and statistical analysis
- **[tests/images/README.md](tests/images/README.md)** - Demonstration of current processing pipeline output

### Process Images Through the Pipeline
The core processing engine can be tested directly:

```bash
# Run image analysis
go run tools/analyze-images/main.go

# Run performance tests
go run tools/performance-test/main.go

# Run unit tests
go test ./tests/... -v

# Run benchmarks
go test ./tests/benchmarks -bench=. -benchmem
```

**Current Capabilities:**
- Characteristic-based color extraction (lightness, saturation, hue organization)
- Multi-dimensional color analysis (frequency, relationships, harmonies)
- Color scheme identification and profile detection
- Performance targets met: <2s processing, <100MB memory for 4K images
- Contrast relationship tracking for accessibility

## What You Will Be Able To Do

### Complete Theme Generation (In Development)
Future CLI application with theme generation:

```bash
# Generate Omarchy theme from image
omarchy-theme-gen generate --image wallpaper.jpg --name "my-theme"

# Clone and modify existing themes
omarchy-theme-gen clone my-theme my-variant --scheme complementary
```

Generated themes will integrate directly with Omarchy's theme selection system.

See **[PROJECT.md](PROJECT.md)** for the development roadmap and **[docs/architecture.md](docs/architecture.md)** for architectural design details.

## How It's Being Built

This project follows Intelligent Development principles with AI-assisted implementation:

### Development Methodology
- **[docs/development-methodology.md](docs/development-methodology.md)** - Development principles and practices
- **[docs/testing-strategy.md](docs/testing-strategy.md)** - Testing standards and requirements
- **[CLAUDE.md](CLAUDE.md)** - Project context for AI-assisted development

### AI-Assisted Development
- **[.claude/agents/](.claude/agents/)** - Specialized agents for development tasks
- User-driven development with AI implementation assistance
- Test-driven development with comprehensive coverage requirements

### Architecture & Design
- Layered architecture with separation of concerns
- Settings-as-Methods pattern for configuration management
- Public Infrastructure Testing standard for API coverage
- Comprehensive diagnostic logging in all tests

## Current Architecture Status

### Foundation Layer (Complete)
- **pkg/formats** - Color space conversions (RGBA, HSLA, LAB, XYZ)
- **pkg/chromatic** - Color theory algorithms and calculations
- **pkg/settings** - Configuration management with category defaults
- **pkg/loader** - Image loading and format validation

### Processing Layer (Refactoring Required)
- **pkg/processor** - Color extraction and characteristic-based organization
- **pkg/errors** - Error handling with sentinel errors

### Generation Layer (In Development)
- **pkg/palette** - Semantic color mapping and theme strategy application
- **pkg/theme** - Component-specific configuration generation

### Application Layer (Planned)
- **cmd/omarchy-theme-gen** - CLI application interface

## Performance Characteristics

Benchmark results on Intel i7-9700K:
- Small images (2.1MP): ~84ms processing, 16.6MB memory
- Large images (14.7MP): ~533ms processing, 118MB memory
- All images meet <2s processing and <100MB memory targets
- Typical category coverage: 25-30% of 27 categories

## Development Commands

```bash
# Run test suite
go test ./tests/... -v

# Run benchmarks
go test ./tests/benchmarks -bench=. -benchmem

# Run integration tests
go test ./tests/integration -v

# Code validation
go vet ./...
go fmt ./...

# Generate image analysis
go run tools/analyze-images/main.go

# Run performance tests
go run tools/performance-test/main.go
```

## Requirements

- Go 1.25+
- No external dependencies (pure Go implementation)
- Linux/Unix environment (developed on Linux)

## Example Processing Output

See **[tests/images/README.md](tests/images/README.md)** for examples of the processing pipeline output:

- Category-based color extraction results
- Color analysis with theme mode detection
- Performance metrics and processing times
- Color scheme identification

## Acknowledgments

Built for the [Omarchy](https://omarchy.org) desktop environment.