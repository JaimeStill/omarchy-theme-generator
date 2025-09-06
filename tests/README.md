# Test Suite Documentation

Comprehensive unit test coverage for the omarchy-theme-generator with diagnostic logging and real-world validation.

## Test Structure

The test suite is organized by package with comprehensive coverage of all implemented functionality. Each test package validates its corresponding pkg/ package with diagnostic output for debugging and verification.

### tests/formats/ - Color Space and Format Tests

Tests for color space conversions, hex parsing, and format utilities.

```bash
go test ./tests/formats -v
```

**Test Coverage:**
- **TestRGBAToHSLA**: Validates RGB to HSLA conversions with known color values
- **TestHSLAToRGBA**: Validates HSLA to RGB conversions with round-trip accuracy  
- **TestParseHex**: Tests hex string parsing (#RRGGBB format) with error handling
- **TestToHex**: Tests RGBA to hex string conversion with proper formatting
- **TestContrastRatio**: Validates WCAG contrast calculations with known ratios
- **TestLABConversions**: Tests LAB color space conversions for perceptual accuracy
- **TestXYZConversions**: Tests XYZ color space conversions as intermediate step

### tests/chromatic/ - Color Theory Algorithm Tests

Tests for color theory calculations, harmony detection, and chromatic analysis.

```bash
go test ./tests/chromatic -v
```

**Test Coverage:**
- **TestLuminance**: Validates perceptual luminance calculations for theme mode detection
- **TestHueDistance**: Tests circular hue distance calculations (0-360°)
- **TestColorDistance**: Validates perceptual color distance metrics
- **TestColorSchemeDetection**: Tests automatic color scheme classification
- **TestContrastRatio**: Validates contrast calculations with WCAG compliance
- **TestHueAnalysis**: Tests hue clustering and dominant hue detection
- **TestSaturationAnalysis**: Validates saturation-based grayscale detection

### tests/settings/ - Configuration Management Tests

Tests for settings loading, validation, and fallback handling.

```bash
go test ./tests/settings -v
```

**Test Coverage:**
- **TestDefaultSettings**: Validates default configuration values and thresholds
- **TestSettingsValidation**: Tests parameter validation and bounds checking
- **TestFallbackColorParsing**: Validates hex string parsing for fallback colors
- **TestViperIntegration**: Tests configuration loading and management
- **TestSettingsAsMethodsPattern**: Validates architectural pattern enforcement
- **TestThresholdValidation**: Tests empirical threshold ranges and defaults

### tests/loader/ - Image I/O and Validation Tests

Tests for image loading, format validation, and metadata extraction.

```bash
go test ./tests/loader -v
```

**Test Coverage:**
- **TestLoadImage**: Validates JPEG and PNG image loading with error handling
- **TestGetImageInfo**: Tests metadata extraction (dimensions, pixel count)
- **TestFormatValidation**: Validates supported format checking and conversion
- **TestImageValidation**: Tests image validation and error handling
- **TestMemoryEfficiency**: Validates memory-efficient image processing
- **TestFileNotFound**: Tests proper error handling for missing files

### tests/processor/ - End-to-End Processing Tests

Comprehensive tests using real wallpaper images from tests/images/ directory.

```bash
go test ./tests/processor -v
```

**Test Coverage:**
- **TestProcessor_New**: Validates processor initialization and configuration
- **TestProcessor_ProcessImage_SimpleImage**: Tests basic processing with simple test image
- **TestProcessor_ProcessImage_GrayscaleImage**: Validates grayscale detection and processing
- **TestProcessor_ProcessImage_MonochromeImage**: Tests monochromatic image handling
- **TestProcessor_ProcessImage_LowFrequencyFiltering**: Tests frequency threshold handling
- **TestProcessor_ThemeMode_Detection**: Validates theme mode detection with known images
- **TestProcessor_FallbackColors**: Tests fallback color application and parsing
- **TestProcessor_InvalidFallbackColors**: Tests graceful handling of invalid hex colors
- **TestProcessor_ColorExtractionQuality**: Validates color extraction quality with complex images

**Diagnostic Output Example:**
```
=== RUN   TestProcessor_ThemeMode_Detection/Dark_night_city_suggests_Dark_theme
    processor_test.go:176: Image: ../../tests/images/night-city.jpeg
    processor_test.go:177: Background luminance: 0.036
    processor_test.go:178: Foreground luminance: 0.847
    processor_test.go:179: Primary luminance: 0.445
    processor_test.go:180: Theme mode threshold: 0.5
    processor_test.go:188: Expected mode: Dark, Actual mode: Dark
--- PASS: TestProcessor_ThemeMode_Detection/Dark_night_city_suggests_Dark_theme (0.21s)
```

## Test Images

The `tests/images/` directory contains real wallpaper samples for validation:

- **abstract.jpeg** (2880x1800, 5.2MP) - Complex multi-color abstract
- **bokeh.jpeg** (4102x2735, 11.2MP) - Bokeh effects, large image
- **coast.jpeg** (1920x1080, 2.1MP) - Bright coastal scene  
- **grayscale.jpeg** (1920x1080, 2.1MP) - Grayscale validation
- **monochrome.jpeg** (1920x1080, 2.1MP) - Single hue validation
- **mountains.jpeg** (1920x1080, 2.1MP) - Natural landscape
- **nebula.jpeg** (3840x2160, 8.3MP) - 4K space scene
- **night-city.jpeg** (2559x1599, 4.1MP) - Dark urban scene
- **portal.jpeg** (1920x1080, 2.1MP) - Distinctive color scheme
- **simple.png** (5120x2880, 14.7MP) - Largest test image
- **warm.jpeg** (2048x1365, 2.8MP) - Warm color palette

## Diagnostic Logging Standards

All tests follow established diagnostic logging patterns:

```go
t.Logf("Image: %s", imagePath)
t.Logf("Background luminance: %v", bgLuminance)
t.Logf("Contrast ratio: %v", contrastRatio)
t.Logf("Expected: %v, Actual: %v", expected, actual)
```

This provides comprehensive debugging information and validates calculation accuracy.

## Running Tests

```bash
# Run all tests with verbose output
go test ./tests/... -v

# Run specific test packages
go test ./tests/formats -v
go test ./tests/chromatic -v  
go test ./tests/settings -v
go test ./tests/loader -v
go test ./tests/processor -v

# Run specific test functions
go test ./tests/processor -run TestProcessor_ThemeMode -v
go test ./tests/formats -run TestParseHex -v

# Run tests with coverage
go test ./tests/... -v -cover
```

## Performance Validation

The processor tests validate performance targets:

- **Processing Time**: <2 seconds for 4K images
- **Memory Usage**: <100MB peak allocation  
- **Success Rate**: 100% processing success across all test images
- **WCAG Compliance**: 4.5:1 minimum contrast ratio enforcement

All tests pass with verified performance compliance and comprehensive diagnostic output.