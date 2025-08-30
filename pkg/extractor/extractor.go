package extractor

import (
	"fmt"
	"image"

	otgcolor "github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/errors"
)

type ExtractionResult struct {
	Image         image.Image
	FrequencyMap  *FrequencyMap
	DominantColor *otgcolor.Color
	TopColors     []*ColorFrequency
	UniqueColors  int
	TotalPixels   uint32
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

// ExtractColors performs complete color extraction from an image file.
// This is the main entry point for file-based extraction with comprehensive analysis.
// If options is nil, DefaultOptions() will be used. Supports optional dimension validation.
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

	return ExtractColorsFromImage(img, options)
}

// ExtractColorsFromImage performs extraction on an already-loaded image.
// Useful when the image comes from memory, network, or other non-file sources.
// Returns complete extraction results with analysis data for synthesis guidance.
func ExtractColorsFromImage(img image.Image, options *ExtractionOptions) (*ExtractionResult, error) {
	if options == nil {
		options = DefaultOptions()
	}

	fm, err := ExtractFromImage(img)
	if err != nil {
		return nil, fmt.Errorf("failed to build frequency map: %w", err)
	}

	dominant := fm.GetDominantColor()
	if dominant == nil {
		return nil, &errors.ExtractionError{
			Stage:   "dominant",
			Details: "no dominant color found",
			Err:     errors.ErrNoColors,
		}
	}

	topColors := fm.GetTopColors(options.TopColorCount)

	result := &ExtractionResult{
		Image:         img,
		FrequencyMap:  fm,
		DominantColor: dominant.Color,
		TopColors:     topColors,
		UniqueColors:  fm.Size(),
		TotalPixels:   fm.Total(),
	}

	return result, nil
}

func (r *ExtractionResult) GetSignificantColors(minThreshold float64) []*ColorFrequency {
	if r.FrequencyMap == nil {
		return nil
	}

	return r.FrequencyMap.FilterByThreshold(minThreshold)
}

func (r *ExtractionResult) GetColorPalette(requestedSize int) ([]*otgcolor.Color, error) {
	if requestedSize <= 0 {
		return nil, fmt.Errorf("palette size must be positive, got %d", requestedSize)
	}

	var colors []*ColorFrequency

	if requestedSize > r.UniqueColors {
		colors = r.FrequencyMap.GetTopColors(0)
	} else {
		colors = r.FrequencyMap.GetTopColors(requestedSize)
	}

	palette := make([]*otgcolor.Color, len(colors))
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
// Distinguishes between grayscale (no hue) and monochromatic (single dominant hue).
type ThemeGenerationAnalysis struct {
	CanExtract        bool    // True if sufficient colors exist for direct extraction
	NeedsSynthesis    bool    // True if color theory synthesis is recommended
	IsGrayscale       bool    // True if image has no color information (saturation ≈ 0)
	IsMonochromatic   bool    // True if image has single dominant hue (±10°) with optional grays
	DominantHue       float64 // The dominant hue in degrees (0-360) if monochromatic
	DominantCoverage  float64 // Percentage of pixels in the most frequent color
	UniqueColors      int     // Total unique colors available for extraction
	SuggestedStrategy string  // "extract", "synthesize", or "hybrid"
	AverageSaturation float64 // Average saturation across top colors (0.0-1.0)
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

	// Analyze saturation and hue distribution
	totalSaturation := 0.0
	sampleSize := min(len(r.TopColors), 100)
	
	// Track hues for monochromatic detection (10-degree bins)
	hueBins := make(map[int]int) // bin index -> count
	grayscaleCount := 0
	
	for i := 0; i < sampleSize; i++ {
		h, s, _ := r.TopColors[i].Color.HSL()
		totalSaturation += s
		
		// Consider colors with saturation < 0.05 as grayscale
		if s < 0.05 {
			grayscaleCount++
		} else {
			// Bin hues into 10-degree segments
			bin := int(h * 360 / 10) % 36 // 36 bins of 10 degrees each
			hueBins[bin]++
		}
	}

	if sampleSize > 0 {
		analysis.AverageSaturation = totalSaturation / float64(sampleSize)
	}

	// Determine if image is grayscale (all colors have very low saturation)
	grayscaleRatio := float64(grayscaleCount) / float64(sampleSize)
	analysis.IsGrayscale = grayscaleRatio > 0.95 // 95% or more grayscale pixels
	
	// Determine if image is monochromatic (single dominant hue with ±10° tolerance)
	if !analysis.IsGrayscale && len(hueBins) > 0 {
		// Find the most populated hue bin and its neighbors
		maxBin := -1
		maxCount := 0
		for bin, count := range hueBins {
			if count > maxCount {
				maxCount = count
				maxBin = bin
			}
		}
		
		if maxBin >= 0 {
			// Count colors in the dominant bin and adjacent bins (±10° = ±1 bin)
			monochromaticCount := hueBins[maxBin]
			prevBin := (maxBin - 1 + 36) % 36
			nextBin := (maxBin + 1) % 36
			monochromaticCount += hueBins[prevBin]
			monochromaticCount += hueBins[nextBin]
			
			// Calculate what percentage of colored pixels fall in this hue range
			coloredPixels := sampleSize - grayscaleCount
			if coloredPixels > 0 {
				monochromaticRatio := float64(monochromaticCount) / float64(coloredPixels)
				analysis.IsMonochromatic = monochromaticRatio > 0.9 // 90% of colored pixels in same hue range
				
				if analysis.IsMonochromatic {
					// Set the dominant hue (center of the bin)
					analysis.DominantHue = float64(maxBin*10 + 5) // Center of the 10-degree bin
				}
			}
		}
	}

	// Determine strategy based on analysis
	minColorsForPureExtraction := 8
	maxDominanceForPureExtraction := 80.0

	if analysis.IsGrayscale {
		// Pure grayscale needs full synthesis
		analysis.CanExtract = false
		analysis.NeedsSynthesis = true
		analysis.SuggestedStrategy = "synthesize"
	} else if analysis.IsMonochromatic {
		// Monochromatic can use hybrid approach
		analysis.CanExtract = false
		analysis.NeedsSynthesis = true
		analysis.SuggestedStrategy = "hybrid"
	} else if r.UniqueColors >= minColorsForPureExtraction && analysis.DominantCoverage < maxDominanceForPureExtraction {
		// Sufficient diversity for pure extraction
		analysis.CanExtract = true
		analysis.NeedsSynthesis = false
		analysis.SuggestedStrategy = "extract"
	} else if r.UniqueColors >= 3 {
		// Some colors but not enough diversity
		analysis.CanExtract = false
		analysis.NeedsSynthesis = true
		analysis.SuggestedStrategy = "hybrid"
	} else {
		// Very few colors, need synthesis
		analysis.CanExtract = false
		analysis.NeedsSynthesis = true
		analysis.SuggestedStrategy = "synthesize"
	}

	return analysis
}

// GetPrimaryNonGrayscale returns the first non-grayscale color found in the top colors.
// It searches through colors in frequency order and returns the first with saturation >= threshold.
// Returns nil if all top colors are grayscale. Useful for finding synthesis seed colors.
func (r *ExtractionResult) GetPrimaryNonGrayscale(saturationThreshold float64) *otgcolor.Color {
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
