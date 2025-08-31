package extractor

import (
	"fmt"
	"image"
	"math"

	tcolor "github.com/JaimeStill/omarchy-theme-generator/pkg/color"
)

// ColorFrequency represents a color and its occurrence frequency in an image.
// It includes both raw count and percentage for convenient display and analysis.
type ColorFrequency struct {
	Color      *tcolor.Color // The color value
	Count      uint32        // Number of pixels with this color
	Percentage float64       // Percentage of total pixels (0.0-100.0)
}

type ExtractionResult struct {
	Image            image.Image
	FrequencyMap     *FrequencyMap
	DominantColor    *tcolor.Color
	TopColors        []*ColorFrequency
	UniqueColors     int
	TotalPixels      uint32
	SelectedStrategy string
}

// ExtractionOptions configures the color extraction process.
// All fields are optional with sensible defaults provided by DefaultOptions().
type ExtractionOptions struct {
	TopColorCount     int     // Number of top colors to return (default: 10)
	MinThreshold      float64 // Minimum percentage for significant colors (default: 0.1%)
	MaxImageDimension int     // Maximum width or height, 0 = unlimited (default: 8192)
}

// DefaultOptions returns sensible default values for extraction configuration.
// Top 10 colors, 0.1% threshold for significance, 8K maximum dimension.
func DefaultOptions() *ExtractionOptions {
	return &ExtractionOptions{
		TopColorCount:     10,
		MinThreshold:      0.1,
		MaxImageDimension: 8192,
	}
}

var defaultSelector *Selector

func InitializeStrategies() {
	defaultSelector = NewSelector()

	defaultSelector.Register(&SaliencyStrategy{
		Settings: DefaultSettings(),
	})

	defaultSelector.SetFallback(&FrequencyStrategy{
		Settings: DefaultSettings(),
	})
}

// ExtractColors loads an image from the file path and extracts its color palette.
//
// The function handles the complete pipeline: loading, validation, resizing (if needed),
// strategy selection, and color extraction. It automatically selects the optimal
// extraction strategy based on image characteristics.
//
// Parameters:
//   - imagePath: File path to image (supports JPEG, PNG, GIF, WebP)
//   - options: Extraction configuration, use nil for defaults
//
// Returns comprehensive extraction results including:
//   - Color frequency analysis with percentages
//   - Dominant color identification  
//   - Strategy recommendation for theme generation
//   - Image characteristics analysis
//
// Error conditions:
//   - File not found or permission denied: wrapped os.PathError
//   - Unsupported image format: ErrUnsupportedFormat
//   - Image too large (>MaxImageDimension): ErrImageTooLarge  
//   - Corrupted image data: image decoding errors
//   - No colors extracted: ErrNoColors (rare, typically corrupted images)
//
// Performance: 4K JPEG images typically process in 1-2 seconds.
func ExtractColors(imagePath string, options *ExtractionOptions) (*ExtractionResult, error) {
	if options == nil {
		options = DefaultOptions()
	}

	var img image.Image
	var err error

	if options.MaxImageDimension > 0 {
		img, err = LoadImageWithValidation(imagePath, options.MaxImageDimension, options.MaxImageDimension)
	} else {
		img, err = LoadImage(imagePath)
	}

	if err != nil {
		return nil, fmt.Errorf("extraction failed: %w", err)
	}

	return ExtractFromLoadedImage(img, options)
}

// ExtractFromLoadedImage extracts colors from an already-loaded image.
//
// This function is useful when the image source is not a file path,
// such as images from memory, network streams, or embedded resources.
// It performs the complete extraction pipeline including frequency analysis
// and returns comprehensive results with strategy recommendations.
//
// If options is nil, DefaultOptions() will be used. The returned
// ExtractionResult contains analysis data to guide synthesis strategies:
// "extract" for sufficient color diversity, "hybrid" for partial extraction,
// or "synthesize" for color theory generation.
//
// Performance: Processes 4K images in <2 seconds with <100MB memory usage.
func ExtractFromLoadedImage(img image.Image, options *ExtractionOptions) (*ExtractionResult, error) {
	if options == nil {
		options = DefaultOptions()
	}

	if defaultSelector == nil {
		InitializeStrategies()
	}

	return defaultSelector.Extract(img, options)
}

func (r *ExtractionResult) GetSignificantColors(minThreshold float64) []*ColorFrequency {
	if r.FrequencyMap == nil {
		return nil
	}

	return r.FrequencyMap.FilterByThreshold(minThreshold)
}

func (r *ExtractionResult) GetColorPalette(requestedSize int) ([]*tcolor.Color, error) {
	if requestedSize <= 0 {
		return nil, fmt.Errorf("palette size must be positive, got %d", requestedSize)
	}

	var colors []*ColorFrequency

	if requestedSize > r.UniqueColors {
		colors = r.FrequencyMap.GetTopColors(0)
	} else {
		colors = r.FrequencyMap.GetTopColors(requestedSize)
	}

	palette := make([]*tcolor.Color, len(colors))
	for i, cf := range colors {
		palette[i] = cf.Color
	}

	return palette, nil
}

type ColorDistribution struct {
	TotalPixels      uint32
	UniqueColors     int
	DominantCoverage float64
	Top10Coverage    float64
	DiversityScore   float64
}

func (r *ExtractionResult) AnalyzeColorDistribution() *ColorDistribution {
	if r.FrequencyMap == nil || len(r.TopColors) == 0 {
		return nil
	}

	dist := &ColorDistribution{
		TotalPixels:      r.TotalPixels,
		UniqueColors:     r.UniqueColors,
		DominantCoverage: 0,
		Top10Coverage:    0,
	}

	if r.DominantColor != nil && len(r.TopColors) > 0 {
		dist.DominantCoverage = r.TopColors[0].Percentage
	}

	for i, cf := range r.TopColors {
		if i >= 10 {
			break
		}
		dist.Top10Coverage += cf.Percentage
	}

	if len(r.TopColors) >= 2 {
		diversityRatio := float64(r.TopColors[1].Count) / float64(r.TopColors[0].Count)
		dist.DiversityScore = diversityRatio
	}

	return dist
}

// ThemeGenerationAnalysis provides detailed analysis for determining theme generation strategy.
// It replaces traditional pass/fail validation with actionable guidance for synthesis systems.
type ThemeGenerationAnalysis struct {
	CanExtract        bool    // True if sufficient colors exist for direct extraction
	NeedsSynthesis    bool    // True if color theory synthesis is recommended
	IsGrayscale       bool    // True if image has no color information (saturation ≈ 0)
	IsMonochromatic   bool    // True if image uses single hue with variations (within hue tolerance)
	DominantCoverage  float64 // Percentage of pixels in the most frequent color
	UniqueColors      int     // Total unique colors available for extraction
	SuggestedStrategy string  // "extract", "synthesize", or "hybrid"
	AverageSaturation float64 // Average saturation across top colors (0.0-1.0)
	DominantHue       float64 // Primary hue for monochromatic images (0.0-360, NaN if grayscale)
	HueTolerance      float64 // Tolerance in degrees for monochromatic detection (typically 15°)
}

// AnalyzeForThemeGeneration analyzes extraction results to determine theme generation strategy.
// Unlike traditional validation that fails, this always succeeds and provides synthesis guidance.
// The analysis considers color diversity, dominance patterns, and saturation characteristics.
func (r *ExtractionResult) AnalyzeForThemeGeneration() *ThemeGenerationAnalysis {
	analysis := &ThemeGenerationAnalysis{
		UniqueColors: r.UniqueColors,
	}

	if len(r.TopColors) > 0 {
		analysis.DominantCoverage = r.TopColors[0].Percentage
	}

	totalSaturation := 0.0
	sampleSize := min(len(r.TopColors), 100)

	for i := 0; i < sampleSize; i++ {
		_, s, _ := r.TopColors[i].Color.HSL()
		totalSaturation += s
	}

	if sampleSize > 0 {
		analysis.AverageSaturation = totalSaturation / float64(sampleSize)
	}

	analysis.HueTolerance = 15.0
	analysis.DominantHue = math.NaN()

	analysis.IsGrayscale = analysis.AverageSaturation < 0.05

	if !analysis.IsGrayscale {
		var dominantHue float64
		foundDominant := false
		foundConflicting := false

		for _, cf := range r.TopColors {
			h, s, _ := cf.Color.HSL()
			if s > 0.05 {
				hueInDegrees := h * 360.0

				if !foundDominant {
					dominantHue = hueInDegrees
					analysis.DominantHue = dominantHue
					foundDominant = true
				} else if !isWithinHueTolerance(dominantHue, hueInDegrees, analysis.HueTolerance) {
					foundConflicting = true
					analysis.IsMonochromatic = false
					break
				}
			}
		}

		if !foundConflicting && foundDominant {
			analysis.IsMonochromatic = true
		}
	}

	minColorsForPureExtraction := 8
	maxDominanceForPureExtraction := 80.0

	if r.UniqueColors >= minColorsForPureExtraction && analysis.DominantCoverage < maxDominanceForPureExtraction {
		analysis.CanExtract = true
		analysis.NeedsSynthesis = false
		analysis.SuggestedStrategy = "extract"
	} else if r.UniqueColors >= 3 && !analysis.IsGrayscale {
		analysis.CanExtract = false
		analysis.NeedsSynthesis = true
		analysis.SuggestedStrategy = "hybrid"
	} else {
		analysis.CanExtract = false
		analysis.NeedsSynthesis = true
		analysis.SuggestedStrategy = "synthesize"
	}

	return analysis
}

// GetPrimaryNonGrayscale returns the first non-grayscale color found in the top colors.
// It searches through colors in frequency order and returns the first with saturation >= threshold.
// Returns nil if all top colors are grayscale. Useful for finding synthesis seed colors.
func (r *ExtractionResult) GetPrimaryNonGrayscale(saturationThreshold float64) *tcolor.Color {
	for _, cf := range r.TopColors {
		_, s, _ := cf.Color.HSL()
		if s >= saturationThreshold {
			return cf.Color
		}
	}
	return nil
}

// String provides a human-readable summary of extraction results including analysis strategy.
// The output includes dimensions, color counts, dominant color, and recommended strategy.
func (r *ExtractionResult) String() string {
	bounds := r.Image.Bounds()
	analysis := r.AnalyzeForThemeGeneration()
	return fmt.Sprintf(
		"ExtractionResult: %dx%d image, %d unique colors from %d pixels, dominant: %s (%.2f%%), strategy: %s",
		bounds.Dx(), bounds.Dy(),
		r.UniqueColors, r.TotalPixels,
		r.DominantColor.HEX(), r.TopColors[0].Percentage,
		analysis.SuggestedStrategy,
	)
}

func isWithinHueTolerance(hue1, hue2, tolerance float64) bool {
	diff := math.Abs(hue1 - hue2)
	if diff > 180.0 {
		diff = 360.0 - diff
	}
	return diff <= tolerance
}
