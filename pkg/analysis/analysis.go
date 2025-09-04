package analysis

import (
	"image/color"
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

type Analyzer struct {
	grayscaleThreshold     float64
	monochromaticTolerance float64
	chroma                 *chromatic.Chroma
}

type ColorProfile struct {
	Mode            ThemeMode
	ColorScheme     chromatic.ColorScheme
	IsGrayscale     bool
	IsMonochromatic bool
	DominantHue     float64
	HueVariance     float64
	AvgLuminance    float64
	AvgSaturation   float64
}

// ThemeMode represents the theme brightness mode based on color analysis.
// Used to determine whether extracted colors suggest a light or dark theme.
type ThemeMode string

const (
	// Light indicates colors suggest a light theme (dark colors on light background)
	Light ThemeMode = "Light"
	// Dark indicates colors suggest a dark theme (light colors on dark background)
	Dark ThemeMode = "Dark"
)

func NewAnalyzer(s *settings.Settings) *Analyzer {
	return &Analyzer{
		grayscaleThreshold:     s.GrayscaleThreshold,
		monochromaticTolerance: s.MonochromaticTolerance,
		chroma:                 chromatic.NewChroma(s),
	}
}

func (a *Analyzer) AnalyzeColors(colors []color.RGBA) ColorProfile {
	profile := ColorProfile{
		Mode: CalculateThemeMode(colors),
	}

	grayscaleCount := 0
	var nonGrayscaleHSLAs []formats.HSLA
	totalSaturation := 0.0
	totalLuminance := 0.0

	for _, c := range colors {
		hsla, isGray := a.IsGrayscale(c)
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
		profile.IsMonochromatic = a.IsMonochromatic(colors)
	} else {
		profile.DominantHue = math.NaN()
		profile.HueVariance = 0.0
		profile.IsMonochromatic = false
	}

	return profile
}

// IsGrayscale checks if a color is grayscale (has very low saturation).
// Returns the color's HSLA representation and true if saturation < 0.05.
// This threshold distinguishes truly grayscale colors from slightly desaturated colors.
func (a *Analyzer) IsGrayscale(c color.RGBA) (formats.HSLA, bool) {
	hsla := formats.RGBAToHSLA(c)
	return hsla, hsla.S < a.grayscaleThreshold
}

// IsMonochromatic determines if a set of colors are monochromatic (similar hues).
// Colors are considered monochromatic if their hues fall within the tolerance (degrees).
// Grayscale colors are excluded from hue analysis since they have no meaningful hue.
// Returns false if fewer than 2 non-grayscale colors are provided.
func (a *Analyzer) IsMonochromatic(colors []color.RGBA) bool {
	if len(colors) < 2 {
		return true
	}

	var validHues []float64

	for _, c := range colors {
		hsla, gray := a.IsGrayscale(c)
		if !gray {
			validHues = append(validHues, hsla.H)
		}
	}

	if len(validHues) < 2 {
		return false
	}

	baseHue := validHues[0]
	for _, hue := range validHues[1:] {
		if !a.chroma.HuesWithinTolerance(baseHue, hue) {
			return false
		}
	}

	return true
}

// CalculateThemeMode determines the appropriate theme mode based on color luminance.
// Analyzes the average luminance of provided colors to suggest Light or Dark theme.
// Colors with average luminance < 0.5 suggest Light theme (dark colors need light background).
// Colors with average luminance >= 0.5 suggest Dark theme (light colors need dark background).
func CalculateThemeMode(colors []color.RGBA) ThemeMode {
	if len(colors) == 0 {
		return Dark
	}

	totalLuminance := 0.0
	for _, c := range colors {
		totalLuminance += chromatic.Luminance(c)
	}

	avgLuminance := totalLuminance / float64(len(colors))

	if avgLuminance < 0.5 {
		return Light
	}

	return Dark
}
