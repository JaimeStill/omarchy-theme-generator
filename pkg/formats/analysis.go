package formats

import (
	"image/color"
	"math"
)

// ThemeMode represents the theme brightness mode based on color analysis.
// Used to determine whether extracted colors suggest a light or dark theme.
type ThemeMode string

const (
	// Light indicates colors suggest a light theme (dark colors on light background)
	Light ThemeMode = "Light"
	// Dark indicates colors suggest a dark theme (light colors on dark background)
	Dark ThemeMode = "Dark"
)

// IsGrayscale checks if a color is grayscale (has very low saturation).
// Returns the color's HSLA representation and true if saturation < 0.05.
// This threshold distinguishes truly grayscale colors from slightly desaturated colors.
func IsGrayscale(c color.RGBA) (HSLA, bool) {
	hsla := RGBAToHSLA(c)
	return hsla, hsla.S < 0.05
}

// IsMonochromatic determines if a set of colors are monochromatic (similar hues).
// Colors are considered monochromatic if their hues fall within the tolerance (degrees).
// Grayscale colors are excluded from hue analysis since they have no meaningful hue.
// Returns false if fewer than 2 non-grayscale colors are provided.
func IsMonochromatic(colors []color.RGBA, tolerance float64) bool {
	if len(colors) < 2 {
		return true
	}

	var validHues []float64

	for _, c := range colors {
		hsla, gray := IsGrayscale(c)
		if !gray {
			validHues = append(validHues, hsla.H)
		}
	}

	if len(validHues) < 2 {
		return false
	}

	baseHue := validHues[0]
	for _, hue := range validHues[1:] {
		if !huesWithinTolerance(baseHue, hue, tolerance) {
			return false
		}
	}

	return true
}

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
	h1 := RGBAToHSLA(c1)
	h2 := RGBAToHSLA(c2)

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
	l1, a1, b1 := rgbaToLAB(c1)
	l2, a2, b2 := rgbaToLAB(c2)

	dl := l1 - l2
	da := a1 - a2
	db := b1 - b2

	return math.Sqrt(dl*dl + da*da + db*db)
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
		totalLuminance += Luminance(c)
	}

	avgLuminance := totalLuminance / float64(len(colors))

	if avgLuminance < 0.5 {
		return Light
	}

	return Dark
}

// huesWithinTolerance checks if two hues are within the specified tolerance in degrees.
// Handles hue wraparound (e.g., 350° and 10° are 20° apart, not 340°).
// Used internally by IsMonochromatic to determine color similarity.
func huesWithinTolerance(h1, h2, tolerance float64) bool {
	diff := math.Abs(h1 - h2)

	if diff > 180 {
		diff = 360 - diff
	}

	return diff <= tolerance
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

// rgbaToLAB converts color.RGBA to CIE LAB color space coordinates.
// LAB is designed to approximate human vision and provides uniform color differences.
// Returns L* (lightness 0-100), a* (green-red axis), and b* (blue-yellow axis).
//
// The conversion process:
//  1. Apply gamma correction to linearize RGB values
//  2. Convert to CIE XYZ color space using sRGB matrix
//  3. Normalize by standard illuminant D65 white point
//  4. Apply LAB transformation with cube root for perceptual uniformity
func rgbaToLAB(c color.RGBA) (l, a, b float64) {
	r := gammaCorrect(float64(c.R) / 255.0)
	g := gammaCorrect(float64(c.G) / 255.0)
	bl := gammaCorrect(float64(c.B) / 255.0)

	x := r*0.4124564 + g*0.3575761 + bl*0.1804375
	y := r*0.2126729 + g*0.7151522 + bl*0.0721750
	z := r*0.0193339 + g*0.1191920 + bl*0.9503041

	x /= 0.95047
	y /= 1.00000
	z /= 1.08883

	x = labTransform(x)
	y = labTransform(y)
	z = labTransform(z)

	l = 116*y - 16
	a = 500 * (x - y)
	b = 200 * (y - z)

	return l, a, b
}

// gammaCorrect applies sRGB gamma correction to linearize color values.
// Converts from sRGB (gamma-corrected) to linear RGB for accurate calculations.
// Uses the standard sRGB specification with threshold at 0.04045.
func gammaCorrect(c float64) float64 {
	if c > 0.04045 {
		return math.Pow((c+0.055)/1.055, 2.4)
	}
	return c / 12.92
}

// labTransform applies the CIE LAB cube root transformation.
// Converts linear tristimulus values to perceptually uniform LAB space.
// Uses threshold at 0.008856 to avoid numerical instability near zero.
func labTransform(t float64) float64 {
	if t > 0.008856 {
		return math.Pow(t, 1.0/3.0)
	}
	return 7.787*t + 16.0/116.0
}
