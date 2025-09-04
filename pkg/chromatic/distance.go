package chromatic

import (
	"image/color"
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

// DistanceRGB calculates the Euclidean distance between two colors in RGB space.
// Returns a value representing color similarity, where 0 means identical colors
// and larger values indicate greater difference. Max distance is ~441 for RGB.
func DistanceRGB(c1, c2 color.RGBA) float64 {
	dr := float64(c1.R) - float64(c2.R)
	dg := float64(c1.G) - float64(c2.G)
	db := float64(c1.B) - float64(c2.B)

	return math.Sqrt(dr*dr + dg*dg + db*db)
}

// DistanceHSL calculates the distance between two colors in HSL space.
// Provides perceptually more uniform color difference than RGB distance.
// Special handling for grayscale colors (low saturation) to focus on lightness.
// Returns normalized distance where components are weighted appropriately.
func DistanceHSL(c1, c2 color.RGBA) float64 {
	h1 := formats.RGBAToHSLA(c1)
	h2 := formats.RGBAToHSLA(c2)

	if h1.S < 0.01 && h2.S < 0.01 {
		dl := h1.L - h2.L
		return math.Abs(dl)
	}

	if h1.S < 0.01 || h2.S < 0.01 {
		ds := h1.S - h2.S
		dl := h1.L - h2.L
		return math.Sqrt(ds*ds + dl*dl)
	}

	dh := hueDistance(h1.H, h2.H) / 180.0
	ds := h1.S - h2.S
	dl := h1.L - h2.L

	return math.Sqrt(dh*dh + ds*ds + dl*dl)
}

// DistanceLAB calculates the Euclidean distance in CIE LAB color space.
// LAB space is designed to be perceptually uniform, making it ideal for
// measuring actual color differences as perceived by humans.
// This is the most accurate method for color similarity comparison.
func DistanceLAB(c1, c2 color.RGBA) float64 {
	lab1 := formats.RGBAToLAB(c1)
	lab2 := formats.RGBAToLAB(c2)

	dl := lab1.L - lab2.L
	da := lab1.A - lab2.A
	db := lab1.B - lab2.B

	return math.Sqrt(dl*dl + da*da + db*db)
}

// hueDistance calculates the shortest angular distance between two hues in degrees.
// Accounts for the circular nature of hue (0° = 360°) by taking the minimum
// of clockwise and counterclockwise distances. Returns value in range [0-180].
func hueDistance(h1, h2 float64) float64 {
	diff := math.Abs(h1 - h2)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}
