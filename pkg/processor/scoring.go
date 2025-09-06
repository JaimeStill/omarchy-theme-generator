package processor

import (
	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
	"image/color"
	"math"
)

func (p *Processor) calculateCategoryFitScore(
	c color.RGBA,
	category ColorCategory,
	profile *ColorProfile,
	background color.RGBA,
	frequency uint32,
	totalPixels uint32,
) float64 {

	chars := p.GetCategoryCharacteristics(category, profile)
	hsla := formats.RGBAToHSLA(c)
	weights := p.settings.CategoryScoring

	score := 0.0

	// Frequency score (normalized by total pixels)
	if weights.Frequency > 0 {
		freqRatio := float64(frequency) / float64(totalPixels)
		// Use logarithmic scale to prevent dominant colors from overwhelming
		freqScore := math.Min(1.0, math.Log10(freqRatio*10000+1)/4)
		score += freqScore * weights.Frequency
	}

	// Contrast score (only if minimum contrast required)
	if weights.Contrast > 0 && chars.MinContrast > 0 {
		contrast := chromatic.ContrastRatio(c, background)
		if contrast >= chars.MinContrast {
			// Score increases with contrast above minimum
			contrastScore := math.Min(1.0, (contrast-chars.MinContrast)/10.0)
			score += contrastScore * weights.Contrast
		} else {
			// Insufficient contrast disqualifies the color
			return 0
		}
	}

	// Saturation score (proximity to ideal range midpoint)
	if weights.Saturation > 0 {
		idealSat := (chars.MinSaturation + chars.MaxSaturation) / 2
		satRange := chars.MaxSaturation - chars.MinSaturation
		if satRange > 0 {
			satDiff := math.Abs(hsla.S - idealSat)
			satScore := 1.0 - (satDiff / satRange)
			score += satScore * weights.Saturation
		} else {
			score += weights.Saturation // Perfect fit if range is zero
		}
	}

	// Lightness score (proximity to ideal range midpoint)
	if weights.Lightness > 0 {
		idealLight := (chars.MinLightness + chars.MaxLightness) / 2
		lightRange := chars.MaxLightness - chars.MinLightness
		if lightRange > 0 {
			lightDiff := math.Abs(hsla.L - idealLight)
			lightScore := 1.0 - (lightDiff / lightRange)
			score += lightScore * weights.Lightness
		} else {
			score += weights.Lightness // Perfect fit if range is zero
		}
	}

	// Hue alignment score (if hue constraints specified)
	if weights.HueAlignment > 0 && chars.HueCenter != nil && chars.HueTolerance != nil {
		hueDiff := math.Abs(hsla.H - *chars.HueCenter)
		if hueDiff > 180 {
			hueDiff = 360 - hueDiff
		}
		if hueDiff <= *chars.HueTolerance {
			hueScore := 1.0 - (hueDiff / *chars.HueTolerance)
			score += hueScore * weights.HueAlignment
		} else {
			// Outside hue tolerance disqualifies for hue-specific categories
			return 0
		}
	}

	return score
}
