package chromatic

import (
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

type Chroma struct {
	settings *settings.Settings
}

func NewChroma(s *settings.Settings) *Chroma {
	return &Chroma{
		settings: s,
	}
}

func (c *Chroma) IdentifyColorScheme(variance float64, colorCount int, hues []float64) ColorScheme {
	tolerance := c.settings.MonochromaticTolerance

	if colorCount == 0 {
		return Grayscale
	}

	if colorCount == 1 {
		return Monochromatic
	}

	if colorCount == 2 {
		hueDiff := hueDistance(hues[0], hues[1])
		if hueDiff >= 170 && hueDiff <= 190 {
			return Complementary
		}
		if variance <= 30 {
			return Analogous
		}
		return Custom
	}

	if colorCount == 3 {
		if isTriadic(hues, tolerance) {
			return Triadic
		}
		if isSplitComplementary(hues, tolerance) {
			return SplitComplementary
		}
		if variance <= 30 {
			return Analogous
		}
		return Custom
	}

	if colorCount == 4 {
		if isSquare(hues, tolerance) {
			return Square
		}
		if isTetradic(hues, tolerance) {
			return Tetradic
		}
		if variance <= 30 {
			return Analogous
		}
		return Custom
	}

	if variance <= tolerance {
		return Monochromatic
	}
	if variance <= 30 {
		return Analogous
	}

	return Custom
}

// HuesWithinTolerance checks if two hues are within the specified tolerance in degrees.
// Handles hue wraparound (e.g., 350° and 10° are 20° apart, not 340°).
// Used internally by IsMonochromatic to determine color similarity.
func (c *Chroma) HuesWithinTolerance(h1, h2 float64) bool {
	tolerance := c.settings.MonochromaticTolerance
	diff := math.Abs(h1 - h2)

	if diff > 180 {
		diff = 360 - diff
	}

	return diff <= tolerance
}
