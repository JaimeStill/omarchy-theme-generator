package formats

import (
	"image/color"
	"math"
)

// AccessibilityLevel represents WCAG contrast ratio requirements for different
// accessibility conformance levels. Each level defines minimum contrast ratios
// for text readability.
type AccessibilityLevel string

const (
	// AA requires 4.5:1 contrast ratio for normal text (WCAG 2.1 Level AA)
	AA AccessibilityLevel = "AA"
	// AAA requires 7.0:1 contrast ratio for normal text (WCAG 2.1 Level AAA)
	AAA AccessibilityLevel = "AAA"
	// AALarge requires 3.0:1 contrast ratio for large text (WCAG 2.1 Level AA)
	AALarge AccessibilityLevel = "AA-large"
	// AAALarge requires 4.5:1 contrast ratio for large text (WCAG 2.1 Level AAA)
	AAALarge AccessibilityLevel = "AAA-large"
)

// Ratio returns the minimum contrast ratio required for the accessibility level.
// Returns 4.5 for AA (default), 7.0 for AAA, and 3.0 for AA-large.
func (a AccessibilityLevel) Ratio() float64 {
	switch a {
	case AAA:
		return 7.0
	case AALarge:
		return 3.0
	default:
		return 4.5
	}
}

// IsAccessible checks if two colors meet the specified WCAG accessibility level
// contrast ratio requirement. Returns true if the contrast ratio between c1 and c2
// meets or exceeds the level's minimum ratio.
func IsAccessible(c1, c2 color.RGBA, level AccessibilityLevel) bool {
	ratio := ContrastRatio(c1, c2)
	return ratio >= level.Ratio()
}

// Luminance calculates the relative luminance of a color according to WCAG 2.1.
// Returns a value between 0 (black) and 1 (white).
//
// Formula: L = 0.2126 * R + 0.7152 * G + 0.0722 * B
// where R, G, B are linearized sRGB values.
func Luminance(c color.RGBA) float64 {
	r := linearize(float64(c.R) / 255.0)
	g := linearize(float64(c.G) / 255.0)
	b := linearize(float64(c.B) / 255.0)

	return 0.2126*r + 0.7152*g + 0.0722*b
}

// linearize converts sRGB color values to linear RGB for luminance calculations.
// Applies the sRGB gamma correction inverse function as specified in WCAG 2.1.
// Values <= 0.03928 use linear scaling, values > 0.03928 use power function.
func linearize(v float64) float64 {
	if v <= 0.03928 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// ContrastRatio calculates the contrast ratio between two colors according to
// WCAG 2.1 guidelines. Returns a value from 1:1 (no contrast) to 21:1 (maximum contrast).
// The formula is: (L1 + 0.05) / (L2 + 0.05) where L1 is the lighter color's luminance.
func ContrastRatio(c1, c2 color.RGBA) float64 {
	l1 := Luminance(c1)
	l2 := Luminance(c2)

	if l1 < l2 {
		l1, l2 = l2, l1
	}

	return (l1 + 0.05) / (l2 + 0.05)
}
