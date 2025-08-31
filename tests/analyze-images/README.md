# Image Analysis Utility

A utility tool that analyzes wallpaper images and generates comprehensive characteristic documentation for the omarchy-theme-generator color extraction system.

## Purpose

This utility analyzes all images in `tests/images/` and generates detailed metadata about their visual characteristics, color extraction behavior, and strategy selection patterns.

## Usage

```bash
# Generate or update tests/images/README.md
go run tests/analyze-images/main.go
```

## Output

The tool generates `tests/images/README.md` containing detailed analysis for each image:

### Image Characteristics
- **Dimensions** - Width x height in pixels
- **Edge Density** - Measure of fine detail and texture (0.0-1.0)
- **Color Complexity** - Number of unique colors in the image
- **Contrast Level** - Brightness variation across the image
- **Average Saturation** - Overall color intensity (0.0-1.0)
- **Dominance Pattern** - How much the most frequent color dominates
- **Distinct Regions** - Whether the image has clear visual boundaries

### Strategy Selection Analysis
- **Selected Strategy** - Which extraction strategy (frequency or saliency) was chosen
- **Selection Logic** - Explanation of why each strategy can/cannot handle the image
- **Priority Scoring** - Relative priority values for strategy selection

### Color Extraction Results
- **Unique Color Count** - Total distinct colors extracted
- **Dominant Color** - Most frequent color with percentage
- **Top 5 Color Palette** - Most significant colors with hex codes and percentages

### Theme Generation Guidance
- **Suggested Strategy** - Recommended approach (extract, synthesize, hybrid)
- **Color Classification** - Whether the image is grayscale, monochromatic, or diverse
- **Synthesis Recommendations** - Guidance for cases where direct extraction isn't optimal

## Integration

The utility uses the same extraction algorithms and settings as the production system in `pkg/extractor/`, ensuring that the documentation accurately reflects real extraction behavior and strategy selection patterns.