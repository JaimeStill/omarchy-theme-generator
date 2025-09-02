package extractor

import (
	"image"
	"image/color"
	"math"
)

// ImageType represents different categories of images based on visual characteristics.
// These classifications drive strategy selection for optimal color extraction.
type ImageType int

const (
	HighDetail ImageType = iota // Images with high edge density and fine details
	LowDetail                   // Images with minimal edge information and simple structure
	Smooth                      // Images with low edges but high color complexity
	Complex                     // Images with both high edges and high color complexity
)

// ImageCharacteristics contains computed metrics about an image's visual properties.
// These characteristics guide extraction strategy selection and parameter tuning
// for optimal color extraction results.
type ImageCharacteristics struct {
	Type               ImageType // Classified image category
	Width              int       // Image width in pixels
	Height             int       // Image height in pixels
	EdgeDensity        float64   // Proportion of pixels with significant edges (0.0-1.0)
	ColorComplexity    int       // Number of unique colors found during sampling
	ContrastLevel      float64   // Standard deviation of brightness values (0.0-1.0)
	HasDistinctRegions bool      // Whether image has identifiable separate regions
	AverageSaturation  float64   // Mean saturation across sampled pixels (0.0-1.0)
	DominancePattern   float64   // Dominance ratio of most frequent color (0.0-1.0)
}

// AnalyzeImageCharacteristics performs comprehensive image analysis using default settings.
// Returns characteristics including edge density, color complexity, and image type classification.
func AnalyzeImageCharacteristics(img image.Image) *ImageCharacteristics {
	return AnalyzeImageCharacteristicsWithSettings(img, CurrentSettings)
}

// AnalyzeImageCharacteristicsWithSettings performs image analysis with custom settings.
// Allows fine-tuning of analysis parameters for specialized use cases or testing.
func AnalyzeImageCharacteristicsWithSettings(img image.Image, settings *Settings) *ImageCharacteristics {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	characteristics := &ImageCharacteristics{
		Width:  width,
		Height: height,
	}

	characteristics.EdgeDensity = calculateEdgeDensity(img, settings)

	colorStats := analyzeColorDistribution(img, settings)
	characteristics.ColorComplexity = colorStats.UniqueColors
	characteristics.AverageSaturation = colorStats.AverageSaturation
	characteristics.DominancePattern = colorStats.DominanceRatio
	characteristics.ContrastLevel = colorStats.ContrastLevel

	characteristics.Type = classifyImageType(characteristics, settings)

	characteristics.HasDistinctRegions = detectRegionSeparation(img, characteristics.EdgeDensity, settings)

	return characteristics
}

// ColorStatistics provides aggregate metrics about color distribution in an image.
// Used internally during image analysis to compute characteristics for strategy selection.
type ColorStatistics struct {
	UniqueColors      int     // Count of distinct colors in sample
	AverageSaturation float64 // Mean saturation value (0.0-1.0)
	DominanceRatio    float64 // Proportion of pixels in dominant color (0.0-1.0)
	ContrastLevel     float64 // Brightness variation level (0.0-1.0)
}

// calculateEdgeDensity computes the proportion of pixels with significant edge gradients.
// Uses Sobel-like edge detection with configurable sampling rate for performance.
func calculateEdgeDensity(img image.Image, settings *Settings) float64 {
	bounds := img.Bounds()
	totalEdgeStrength := 0.0
	edgePixels := 0

	sampleRate := settings.Analysis.EdgeDetectionSampleRate
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y += sampleRate {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x += sampleRate {
			gx := getGrayscaleValue(img.At(x+1, y)) - getGrayscaleValue(img.At(x-1, y))
			gy := getGrayscaleValue(img.At(x, y+1)) - getGrayscaleValue(img.At(x, y-1))

			edgeStrength := math.Sqrt(float64(gx*gx + gy*gy))

			if edgeStrength > settings.Analysis.EdgeDetectionMinStrength {
				totalEdgeStrength += edgeStrength
				edgePixels++
			}
		}
	}

	sampledPixels := ((bounds.Dx() - 2) / sampleRate) * ((bounds.Dy() - 2) / sampleRate)
	if sampledPixels == 0 {
		return 0.0
	}

	return float64(edgePixels) / float64(sampledPixels)
}

func analyzeColorDistribution(img image.Image, settings *Settings) *ColorStatistics {
	bounds := img.Bounds()
	colorCount := make(map[uint32]int)
	totalSaturation := 0.0
	totalPixels := 0

	sampleRate := settings.Extraction.ColorDistSampleRate
	for y := bounds.Min.Y; y < bounds.Max.Y; y += sampleRate {
		for x := bounds.Min.X; x < bounds.Max.X; x += sampleRate {
			rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			packed := uint32(rgba.R)<<16 | uint32(rgba.G)<<8 | uint32(rgba.B)

			colorCount[packed]++

			r, g, b := float64(rgba.R)/255.0, float64(rgba.G)/255.0, float64(rgba.B)/255.0

			maxVal := math.Max(r, math.Max(g, b))
			minVal := math.Min(r, math.Min(g, b))

			saturation := 0.0
			if maxVal != 0 {
				saturation = (maxVal - minVal) / maxVal
			}
			totalSaturation += saturation
			totalPixels++
		}
	}

	maxCount := 0
	for _, count := range colorCount {
		if count > maxCount {
			maxCount = count
		}
	}

	dominanceRatio := 0.0
	if totalPixels > 0 {
		dominanceRatio = float64(maxCount) / float64(totalPixels)
	}

	avgSaturation := 0.0
	if totalPixels > 0 {
		avgSaturation = totalSaturation / float64(totalPixels)
	}

	return &ColorStatistics{
		UniqueColors:      len(colorCount),
		AverageSaturation: avgSaturation,
		DominanceRatio:    dominanceRatio,
		ContrastLevel:     calculateContrastLevel(img),
	}
}

func calculateContrastLevel(img image.Image) float64 {
	bounds := img.Bounds()
	values := make([]float64, 0, bounds.Dx()*bounds.Dy()/16)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 4 {
		for x := bounds.Min.X; x < bounds.Max.Y; x += 4 {
			brightness := getGrayscaleValue(img.At(x, y))
			values = append(values, float64(brightness))
		}
	}

	if len(values) < 2 {
		return 0.0
	}

	mean := 0.0
	for _, v := range values {
		mean += v
	}

	mean /= float64(len(values))

	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}

	variance /= float64(len(values))

	return math.Sqrt(variance) / 255.0
}

// classifyImageType determines the image category based on computed characteristics.
// Uses empirically-derived thresholds to classify images for strategy selection.
func classifyImageType(chars *ImageCharacteristics, settings *Settings) ImageType {
	analysis := settings.Analysis

	if chars.EdgeDensity > analysis.HighDetailEdgeThreshold {
		return HighDetail
	}

	if chars.EdgeDensity < analysis.SmoothEdgeThreshold && chars.ColorComplexity > analysis.SmoothColorThreshold {
		return Smooth
	}

	if chars.ColorComplexity < analysis.LowDetailColorThreshold {
		return LowDetail
	}

	if chars.ColorComplexity > analysis.ComplexColorThreshold && chars.EdgeDensity > analysis.ComplexEdgeThreshold {
		return Complex
	}

	return HighDetail
}

func detectRegionSeparation(img image.Image, edgeDensity float64, settings *Settings) bool {
	analysis := settings.Analysis
	return edgeDensity > analysis.RegionMinEdgeDensity && edgeDensity < analysis.RegionMaxEdgeDensity
}

// getGrayscaleValue converts a color to grayscale using luminance weighting.
// Uses standard RGB luminance formula: 0.299*R + 0.587*G + 0.114*B
func getGrayscaleValue(c color.Color) int {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return int(0.299*float64(rgba.R) + 0.587*float64(rgba.G) + 0.114*float64(rgba.B))
}
