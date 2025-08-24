// Package color provides WCAG-compliant contrast calculation methods.
// All calculations follow WCAG 2.1 guidelines with proper gamma correction
// for accessibility compliance validation.
package color

import "math"

// luminanceCache stores the computed relative luminance value to avoid
// expensive recalculation on repeated contrast ratio computations.
type luminanceCache struct {
	value float64
}

// ContrastRatio calculates the WCAG-compliant contrast ratio between two colors.
// Returns a value from 1:1 (no contrast) to 21:1 (maximum contrast).
// The calculation ensures the lighter color is always the numerator.
func (c *Color) ContrastRatio(a *Color) float64 {
	l1 := c.RelativeLuminance()
	l2 := a.RelativeLuminance()

	if l1 < l2 {
		l1, l2 = l2, l1
	}

	return (l1 + 0.05) / (l2 + 0.05)
}

// IsAccessible checks if the color pair meets the specified WCAG accessibility level.
// Returns true if the contrast ratio meets or exceeds the required threshold.
func (c *Color) IsAccessible(a *Color, level AccessibilityLevel) bool {
	ratio := c.ContrastRatio(a)
	return ratio >= level.Ratio()
}

// MeetsWCAG is a convenience method that checks AA compliance (4.5:1 ratio).
// Equivalent to IsAccessible(other, AA) for standard accessibility requirements.
func (c *Color) MeetsWCAG(a *Color) bool {
	return c.IsAccessible(a, AA)
}

// RelativeLuminance calculates the relative luminance according to WCAG 2.1.
// Uses cached computation with thread-safe sync.Once for performance.
// Formula: 0.2126*R + 0.7152*G + 0.0722*B with gamma correction.
func (c *Color) RelativeLuminance() float64 {
	c.luminanceOnce.Do(func() {
		r := toLinearRGB(float64(c.R) / 255.0)
		g := toLinearRGB(float64(c.G) / 255.0)
		b := toLinearRGB(float64(c.B) / 255.0)

		c.luminance = &luminanceCache{
			value: 0.2126*r + 0.7152*g + 0.0722*b,
		}
	})

	return c.luminance.value
}

// toLinearRGB applies gamma correction to convert sRGB to linear RGB.
// Essential for accurate luminance calculation per WCAG 2.1 specification.
func toLinearRGB(channel float64) float64 {
	if channel <= 0.03928 {
		return channel / 12.92
	}
	return math.Pow((channel+0.055)/1.055, 2.4)
}

// AccessibilityLevel represents WCAG compliance levels with associated contrast ratios.
// Provides type safety and automatic ratio lookup for accessibility validation.
type AccessibilityLevel string

const (
	AA       AccessibilityLevel = "AA"
	AAA      AccessibilityLevel = "AAA"
	AALarge  AccessibilityLevel = "AA-large"
	AAALarge AccessibilityLevel = "AAA-large"
)

// Ratio returns the minimum contrast ratio required for the accessibility level.
// Based on WCAG 2.1 guidelines: AA(4.5:1), AAA(7.0:1), with large text exceptions.
func (a AccessibilityLevel) Ratio() float64 {
	switch a {
	case AA:
		return 4.5
	case AAA:
		return 7.0
	case AALarge:
		return 3.0
	case AAALarge:
		return 4.5
	default:
		return 4.5
	}
}
