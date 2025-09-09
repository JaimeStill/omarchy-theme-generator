package processor

import (
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func (p *Processor) analyzeColors(colors []WeightedColor) *ColorProfile {
	profile := &ColorProfile{
		Mode: p.calculateThemeMode(colors),
	}

	var grayscaleWeight float64
	var totalWeight float64
	var nonGrayscaleColors []WeightedColor
	var nonGrayscaleHSLAs []formats.HSLA
	var weightedSaturation float64
	var weightedLuminance float64

	for _, wc := range colors {
		hsla, isGray := p.isGrayscaleWeighted(wc)
		totalWeight += wc.Weight

		if isGray {
			grayscaleWeight += wc.Weight
		} else {
			nonGrayscaleColors = append(nonGrayscaleColors, wc)
			nonGrayscaleHSLAs = append(nonGrayscaleHSLAs, formats.RGBAToHSLA(wc.RGBA))
		}

		weightedSaturation += hsla.S * wc.Weight
		weightedLuminance += chromatic.Luminance(wc.RGBA) * wc.Weight
	}

	profile.IsGrayscale = grayscaleWeight/totalWeight > p.settings.GrayscaleImageThreshold

	if totalWeight > 0 {
		profile.AvgSaturation = weightedSaturation / totalWeight
		profile.AvgLuminance = weightedLuminance / totalWeight
	}

	if len(nonGrayscaleColors) > 0 {
		profile.DominantHue = p.findWeightedDominantHue(nonGrayscaleColors)
		profile.HueVariance = chromatic.CalculateHueVariance(nonGrayscaleHSLAs)
		profile.IsMonochromatic = p.isMonochromaticWeighted(colors)
	} else {
		profile.DominantHue = math.NaN()
		profile.HueVariance = 0.0
		profile.IsMonochromatic = false
	}

	return profile
}

func (p *Processor) calculateThemeMode(colors []WeightedColor) ThemeMode {
	if len(colors) == 0 {
		return Dark
	}

	var weightedLuminance float64
	var totalWeight float64

	for _, wc := range colors {
		luminance := chromatic.Luminance(wc.RGBA)
		weightedLuminance += luminance * wc.Weight
		totalWeight += wc.Weight
	}

	avgLuminance := weightedLuminance / totalWeight

	if avgLuminance >= p.settings.ThemeModeThreshold {
		return Light
	}

	return Dark
}

func (p *Processor) findWeightedDominantHue(colors []WeightedColor) float64 {
	if len(colors) == 0 {
		return 0
	}

	if len(colors) == 1 {
		hsla := formats.RGBAToHSLA(colors[0].RGBA)
		return hsla.H
	}

	var sinSum, cosSum, totalWeight float64

	for _, wc := range colors {
		hsla := formats.RGBAToHSLA(wc.RGBA)
		radians := hsla.H * math.Pi / 180
		sinSum += math.Sin(radians) * wc.Weight
		cosSum += math.Cos(radians) * wc.Weight
		totalWeight += wc.Weight
	}

	if totalWeight == 0 {
		hsla := formats.RGBAToHSLA(colors[0].RGBA)
		return hsla.H
	}

	sinSum /= totalWeight
	cosSum /= totalWeight

	meanRadians := math.Atan2(sinSum, cosSum)
	meanDegrees := meanRadians * 180 / math.Pi

	if meanDegrees < 0 {
		meanDegrees += 360
	}
	if meanDegrees >= 360 {
		meanDegrees -= 360
	}

	return meanDegrees
}

func (p *Processor) isGrayscaleWeighted(wc WeightedColor) (formats.HSLA, bool) {
	hsla := formats.RGBAToHSLA(wc.RGBA)
	return hsla, hsla.S < p.settings.GrayscaleThreshold
}

func (p *Processor) isMonochromaticWeighted(colors []WeightedColor) bool {
	if len(colors) < 2 {
		return true
	}

	var validHues []float64
	var validWeights []float64

	for _, wc := range colors {
		hsla, gray := p.isGrayscaleWeighted(wc)
		if !gray {
			validHues = append(validHues, hsla.H)
			validWeights = append(validWeights, wc.Weight)
		}
	}

	if len(validHues) < 2 {
		return false
	}

	maxWeightIndex := 0
	for i, weight := range validWeights {
		if weight > validWeights[maxWeightIndex] {
			maxWeightIndex = i
		}
	}

	baseHue := validHues[maxWeightIndex]

	significantWeightThreshold := validWeights[maxWeightIndex] * p.settings.MonochromaticWeightThreshold

	for i, hue := range validHues {
		if validWeights[i] >= significantWeightThreshold {
			if !p.chroma.HuesWithinTolerance(baseHue, hue) {
				return false
			}
		}
	}

	return true
}
