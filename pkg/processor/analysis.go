package processor

import (
	"image/color"
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func (p *Processor) analyzeColors(colorFreq map[color.RGBA]uint32) *ColorProfile {
	colors := make([]color.RGBA, 0, len(colorFreq))
	for c := range colorFreq {
		colors = append(colors, c)
	}

	profile := &ColorProfile{
		Mode: p.calculateThemeMode(colors),
	}

	grayscaleCount := 0
	var nonGrayscaleHSLAs []formats.HSLA
	totalSaturation := 0.0
	totalLuminance := 0.0

	for _, c := range colors {
		hsla, isGray := p.isGrayscale(c)
		if isGray {
			grayscaleCount++
		} else {
			nonGrayscaleHSLAs = append(nonGrayscaleHSLAs, hsla)
		}
		totalSaturation += hsla.S
		totalLuminance += chromatic.Luminance(c)
	}

	profile.IsGrayscale = grayscaleCount == len(colors)
	profile.AvgSaturation = totalSaturation / float64(len(colors))
	profile.AvgLuminance = totalLuminance / float64(len(colors))

	if len(nonGrayscaleHSLAs) > 0 {
		profile.DominantHue = chromatic.FindDominantHue(nonGrayscaleHSLAs)
		profile.HueVariance = chromatic.CalculateHueVariance(nonGrayscaleHSLAs)
		profile.IsMonochromatic = p.isMonochromatic(colors)
	} else {
		profile.DominantHue = math.NaN()
		profile.HueVariance = 0.0
		profile.IsMonochromatic = false
	}

	return profile
}

func (p *Processor) isGrayscale(c color.RGBA) (formats.HSLA, bool) {
	hsla := formats.RGBAToHSLA(c)
	return hsla, hsla.S < p.settings.GrayscaleThreshold
}

func (p *Processor) isMonochromatic(colors []color.RGBA) bool {
	if len(colors) < 2 {
		return true
	}

	var validHues []float64

	for _, c := range colors {
		hsla, gray := p.isGrayscale(c)
		if !gray {
			validHues = append(validHues, hsla.H)
		}
	}

	if len(validHues) < 2 {
		return false
	}

	baseHue := validHues[0]
	for _, hue := range validHues[1:] {
		if !p.chroma.HuesWithinTolerance(baseHue, hue) {
			return false
		}
	}

	return true
}

func (p *Processor) calculateThemeMode(colors []color.RGBA) ThemeMode {
	if len(colors) == 0 {
		return Dark
	}

	totalLuminance := 0.0
	for _, c := range colors {
		totalLuminance += chromatic.Luminance(c)
	}

	avgLuminance := totalLuminance / float64(len(colors))

	if avgLuminance >= p.settings.ThemeModeThreshold {
		return Light
	}

	return Dark
}
