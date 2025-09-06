# Development Tools

Utility tools for analyzing, documenting, and validating the omarchy-theme-generator system.

## tools/analyze-images/ - Image Analysis Documentation Generator

Generates comprehensive image analysis documentation from test images using the complete processing pipeline.

### Usage

```bash
# Analyze images in default directory (tests/images)
go run tools/analyze-images/main.go

# Analyze images in custom directory
go run tools/analyze-images/main.go --images /path/to/images

# View help
go run tools/analyze-images/main.go --help
```

### Command-Line Options

- `--images` - Directory containing test images (default: "tests/images")

### Output

Creates `README.md` in the specified images directory with comprehensive analysis including:

- **Image metadata** - Dimensions, pixel count, file format
- **Color extraction results** - All role-based colors with hex values and visual swatches
- **Color profile analysis** - Theme mode, color scheme, grayscale/monochromatic detection
- **Visual color swatches** - HTML color squares for easy visualization

### Example Output

```markdown
## abstract.jpeg

![abstract.jpeg](./abstract.jpeg)

**Dimensions**: 2880 x 1800 px

### Image Colors

| Function | Value | Color |
|----------|-------|-------|
| Background | `#1a1a1a` | <span style="..."></span> |
| Foreground | `#f0f0f0` | <span style="..."></span> |
| Primary | `#ff6b35` | <span style="..."></span> |

### Color Profile

| Property | Value |
|----------|-------|
| Mode | Dark |
| Color Scheme | Complex |
| Dominant Hue | 15.2° |
```

## tools/performance-test/ - Comprehensive Performance Validation

Runs statistical performance analysis across all test images to validate system performance targets.

### Usage

```bash
# Test images in default directory (tests/images)
go run tools/performance-test/main.go

# Test images in custom directory  
go run tools/performance-test/main.go --images /path/to/images

# View help
go run tools/performance-test/main.go --help
```

### Command-Line Options

- `--images` - Directory containing test images (default: "tests/images")

### Performance Metrics

The tool validates against established performance targets:

- **Processing Time**: <2 seconds per image (target)
- **Memory Usage**: <100MB peak allocation (target)
- **Success Rate**: 100% processing success

### Output Format

```
Comprehensive Performance Test
Target: < 2 seconds for 4K images (4096x2160 = 8.8MP)
Target: < 100MB peak memory usage

Processing 15 images...

abstract.jpeg        2880x1800    5.2MP    269ms   21.4MB ✅ PASS
bokeh.jpeg           4102x2735   11.2MP    693ms   45.9MB ✅ PASS
nebula.jpeg          3840x2160    8.3MP    465ms    0.0MB ✅ PASS

======================================================================
PERFORMANCE ANALYSIS SUMMARY
======================================================================

Processing Time Statistics:
  Average:    236ms
  Median:     137ms
  Min:         99ms
  Max:        693ms

Memory Usage Statistics:
  Average:    8.6MB
  Median:     0.0MB
  Min:        0.0MB
  Max:       61.2MB

Performance by Image Size:
  Small (<2MP):   0 images, avg    0ms
  Medium (2-8MP): 12 images, avg  147ms
  Large (>8MP):   3 images, avg  593ms

Performance Target Analysis:
  Time Target (< 2s):     15/15 images (100.0%)
  Memory Target (< 100MB): 15/15 images (100.0%)
  🎉 ALL PERFORMANCE TARGETS MET!
```

### Statistical Analysis

- **Individual image metrics** - Processing time, memory usage, dimensions
- **Statistical summaries** - Min/max/median/average for time and memory
- **Performance categorization** - Analysis by image size (small/medium/large)
- **Target compliance reporting** - Success rates against established targets
- **Error summary** - Details of any processing failures

## Integration with Development Workflow

### Image Analysis Integration

The analyze-images tool is integrated into the development workflow for:

- **Test image documentation** - Automatic generation of visual analysis
- **Color extraction validation** - Verify extraction results across image types
- **Visual debugging** - HTML color swatches for manual verification
- **Documentation maintenance** - Keep test image documentation current

### Performance Validation Integration

The performance-test tool provides:

- **Continuous performance monitoring** - Regular validation against targets
- **Regression detection** - Identify performance degradation
- **Scalability analysis** - Performance characteristics across image sizes
- **Memory profiling** - Memory allocation patterns and peak usage

## Development Requirements

Both tools require:

- **Go 1.25+** - Compatible with project Go version
- **Test images** - Located in tests/images/ or specified directory
- **Package dependencies** - All foundation and processing packages

The tools use the same processing pipeline as the main application, ensuring consistency between development analysis and production behavior.