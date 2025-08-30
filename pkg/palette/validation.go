package palette

import (
	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
)

// ValidatePalette checks a palette for WCAG compliance and other quality metrics.
type PaletteValidator struct {
	MinContrast   float64 // Minimum contrast ratio (4.5 for AA, 7.0 for AAA)
	Background    *color.Color
	TestForeground bool // Whether to test as foreground colors
}

// NewPaletteValidator creates a validator with default WCAG AA compliance.
func NewPaletteValidator(background *color.Color) *PaletteValidator {
	if background == nil {
		background = color.NewRGB(255, 255, 255) // Default white background
	}
	
	return &PaletteValidator{
		MinContrast:    4.5, // WCAG AA
		Background:     background,
		TestForeground: true,
	}
}

// ValidationResult contains the results of palette validation.
type ValidationResult struct {
	TotalColors      int
	PassingColors    int
	FailingColors    int
	AverageContrast  float64
	MinimumContrast  float64
	MaximumContrast  float64
	FailingIndices   []int
	ContrastRatios   []float64
}

// Validate checks each color in the palette for WCAG compliance.
func (pv *PaletteValidator) Validate(palette []*color.Color) *ValidationResult {
	result := &ValidationResult{
		TotalColors:     len(palette),
		FailingIndices:  make([]int, 0),
		ContrastRatios:  make([]float64, len(palette)),
		MinimumContrast: 999999.0,
		MaximumContrast: 0.0,
	}
	
	totalContrast := 0.0
	
	for i, c := range palette {
		var contrast float64
		
		if pv.TestForeground {
			contrast = c.ContrastRatio(pv.Background)
		} else {
			contrast = pv.Background.ContrastRatio(c)
		}
		
		result.ContrastRatios[i] = contrast
		totalContrast += contrast
		
		if contrast >= pv.MinContrast {
			result.PassingColors++
		} else {
			result.FailingColors++
			result.FailingIndices = append(result.FailingIndices, i)
		}
		
		if contrast < result.MinimumContrast {
			result.MinimumContrast = contrast
		}
		if contrast > result.MaximumContrast {
			result.MaximumContrast = contrast
		}
	}
	
	if result.TotalColors > 0 {
		result.AverageContrast = totalContrast / float64(result.TotalColors)
	}
	
	return result
}

// EnsureContrast adjusts colors in a palette to meet minimum contrast requirements.
// It modifies colors in-place, adjusting lightness to achieve the target contrast.
func (pv *PaletteValidator) EnsureContrast(palette []*color.Color) []*color.Color {
	adjusted := make([]*color.Color, len(palette))
	copy(adjusted, palette)
	
	for i, c := range adjusted {
		contrast := c.ContrastRatio(pv.Background)
		
		if contrast < pv.MinContrast {
			// Adjust lightness to meet contrast requirement
			adjusted[i] = pv.adjustForContrast(c)
		}
	}
	
	return adjusted
}

// adjustForContrast modifies a color's lightness to meet the minimum contrast ratio.
func (pv *PaletteValidator) adjustForContrast(c *color.Color) *color.Color {
	h, s, l := c.HSL()
	backgroundLuminance := pv.Background.RelativeLuminance()
	
	// Binary search for the right lightness value
	minL, maxL := 0.0, 1.0
	bestL := l
	bestContrast := 0.0
	
	for iterations := 0; iterations < 20; iterations++ {
		testColor := color.NewHSL(h, s, bestL)
		contrast := testColor.ContrastRatio(pv.Background)
		
		if contrast >= pv.MinContrast {
			return testColor
		}
		
		if contrast > bestContrast {
			bestContrast = contrast
		}
		
		// Determine which direction to adjust
		if testColor.RelativeLuminance() > backgroundLuminance {
			// Make lighter
			minL = bestL
			bestL = (bestL + maxL) / 2
		} else {
			// Make darker
			maxL = bestL
			bestL = (minL + bestL) / 2
		}
	}
	
	// Return the best we could achieve
	return color.NewHSL(h, s, bestL)
}

// FindAccessiblePairs finds all color pairs in a palette that meet contrast requirements.
func FindAccessiblePairs(palette []*color.Color, minContrast float64) [][2]int {
	pairs := make([][2]int, 0)
	
	for i := 0; i < len(palette); i++ {
		for j := i + 1; j < len(palette); j++ {
			contrast := palette[i].ContrastRatio(palette[j])
			if contrast >= minContrast {
				pairs = append(pairs, [2]int{i, j})
			}
		}
	}
	
	return pairs
}

// PaletteMetrics provides comprehensive analysis of a color palette.
type PaletteMetrics struct {
	HueVariance        float64 // Variance in hue values (0-1)
	SaturationVariance float64 // Variance in saturation (0-1)
	LightnessVariance  float64 // Variance in lightness (0-1)
	ColorHarmony       float64 // Measure of color relationships (0-1)
	Distinctiveness    float64 // How distinct colors are from each other (0-1)
}

// AnalyzePalette computes various metrics for a color palette.
func AnalyzePalette(palette []*color.Color) *PaletteMetrics {
	if len(palette) == 0 {
		return &PaletteMetrics{}
	}
	
	metrics := &PaletteMetrics{}
	n := float64(len(palette))
	
	// Calculate means
	meanH, meanS, meanL := 0.0, 0.0, 0.0
	for _, c := range palette {
		h, s, l := c.HSL()
		meanH += h
		meanS += s
		meanL += l
	}
	meanH /= n
	meanS /= n
	meanL /= n
	
	// Calculate variances
	varH, varS, varL := 0.0, 0.0, 0.0
	for _, c := range palette {
		h, s, l := c.HSL()
		varH += (h - meanH) * (h - meanH)
		varS += (s - meanS) * (s - meanS)
		varL += (l - meanL) * (l - meanL)
	}
	metrics.HueVariance = varH / n
	metrics.SaturationVariance = varS / n
	metrics.LightnessVariance = varL / n
	
	// Calculate distinctiveness (average distance between colors)
	totalDistance := 0.0
	comparisons := 0
	for i := 0; i < len(palette); i++ {
		for j := i + 1; j < len(palette); j++ {
			totalDistance += palette[i].DistanceHSL(palette[j])
			comparisons++
		}
	}
	if comparisons > 0 {
		metrics.Distinctiveness = totalDistance / float64(comparisons)
	}
	
	// Calculate harmony (inverse of hue variance for now)
	metrics.ColorHarmony = 1.0 - metrics.HueVariance
	
	return metrics
}