package chromatic

import (
	"image/color"
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
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

// ColorsSimilar determines if two colors should be clustered together using
// perceptual distance metrics with special handling for neutral colors.
// Uses LAB color space for accurate perceptual similarity assessment.
func (c *Chroma) ColorsSimilar(c1, c2 color.RGBA) bool {
	// Special handling for neutrals
	h1 := formats.RGBAToHSLA(c1)
	h2 := formats.RGBAToHSLA(c2)

	// If both are neutral, use configurable threshold based on lightness difference
	if h1.S < c.settings.Chromatic.NeutralThreshold && h2.S < c.settings.Chromatic.NeutralThreshold {
		return math.Abs(h1.L-h2.L) < c.settings.Chromatic.NeutralLightnessThreshold
	}

	// Use LAB distance for perceptual similarity
	distance := DistanceLAB(c1, c2)
	return distance <= c.settings.Chromatic.ColorMergeThreshold
}
