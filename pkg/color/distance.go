// Package color provides various color distance and similarity calculation methods.
// Different distance metrics serve different purposes: RGB for computational speed,
// HSL for perceptual similarity, LAB for color science accuracy.
package color

import "math"

// DistanceRGB calculates Euclidean distance in RGB color space.
// Fast calculation but doesn't account for human color perception.
// Useful for computational clustering and basic similarity checks.
func (c *Color) DistanceRGB(a *Color) float64 {
	dr := float64(c.R) - float64(a.R)
	dg := float64(c.G) - float64(a.G)
	db := float64(c.B) - float64(a.B)

	return math.Sqrt(dr*dr + dg*dg + db*db)
}

// DistanceHSL calculates perceptually-weighted distance in HSL color space.
// Weights lightness (2.0) higher than saturation (1.0) and hue (0.5) for better
// perceptual similarity matching. Handles hue wraparound correctly.
func (c *Color) DistanceHSL(a *Color) float64 {
	h1, s1, l1 := c.HSL()
	h2, s2, l2 := a.HSL()

	dh := hueDistance(h1, h2)
	ds := s1 - s2
	dl := l1 - l2

	return math.Sqrt(dh*dh*0.5 + ds*ds*1.0 + dl*dl*2.0)
}

// DistanceLAB calculates distance in CIE LAB color space (Delta-E CIE76).
// Returns both LAB representations and the perceptual distance.
// Most accurate for human color perception, used in color science applications.
func (c *Color) DistanceLAB(a *Color) (LABColor, LABColor, float64) {
	lab1 := c.ToLAB()
	lab2 := a.ToLAB()

	dl := lab1.L - lab2.L
	da := lab1.A - lab2.A
	db := lab1.B - lab2.B

	distance := math.Sqrt(dl*dl + da*da + db*db)

	return lab1, lab2, distance
}

// DistanceLuminance calculates the absolute difference in relative luminance.
// Useful for accessibility analysis and contrast-based color grouping.
func (c *Color) DistanceLuminance(a *Color) float64 {
	l1 := c.RelativeLuminance()
	l2 := a.RelativeLuminance()

	return math.Abs(l1 - l2)
}

// IsSimilar checks if two colors are within the specified HSL distance threshold.
// Uses perceptually-weighted HSL distance for natural similarity detection.
func (c *Color) IsSimilar(a *Color, threshold float64) bool {
	return c.DistanceHSL(a) <= threshold
}

// ClosestColor finds the nearest color from a set of candidates using HSL distance.
// Returns the index of the closest color and its distance value.
// Returns (-1, 0) for empty candidate sets.
func (c *Color) ClosestColor(candidates []*Color) (index int, distance float64) {
	if len(candidates) == 0 {
		return -1, 0
	}

	minDistance := math.Inf(1)
	closestIndex := 0

	for i, candidate := range candidates {
		dist := c.DistanceHSL(candidate)
		if dist < minDistance {
			minDistance = dist
			closestIndex = i
		}
	}

	return closestIndex, minDistance
}

// IsDistinct checks if this color is sufficiently different from all colors in the set.
// Returns true if the color maintains the threshold distance from all existing colors.
// Useful for palette generation to avoid similar colors.
func (c *Color) IsDistinct(colors []*Color, threshold float64) bool {
	for _, a := range colors {
		if c.IsSimilar(a, threshold) {
			return false
		}
	}

	return true
}

// WeightedDistance calculates HSL distance with custom component weights.
// Allows fine-tuning distance calculation for specific use cases.
// Standard weights: hue(0.5), saturation(1.0), lightness(2.0).
func (c *Color) WeightedDistance(a *Color, hueWeight, satWeight, lightWeight float64) float64 {
	h1, s1, l1 := c.HSL()
	h2, s2, l2 := a.HSL()

	dh := hueDistance(h1, h2)
	ds := s1 - s2
	dl := l1 - l2

	return math.Sqrt(dh*dh*hueWeight + ds*ds*satWeight + dl*dl*lightWeight)
}

// ColorTemperatureDistance estimates warmth/coolness difference between colors.
// Returns normalized distance (0.0-1.0) based on approximate color temperature.
// Useful for grouping colors by warm/cool characteristics.
func (c *Color) ColorTemperatureDistance(a *Color) float64 {
	t1 := c.approximateTemperature()
	t2 := a.approximateTemperature()

	return math.Abs(t1-t2) / 2.0
}

// approximateTemperature estimates color warmth on a 0-1 scale.
// Based on red/yellow content (warm) vs blue content (cool).
// Simplified approximation for basic temperature classification.
func (c *Color) approximateTemperature() float64 {
	r, g, b := float64(c.R), float64(c.G), float64(c.B)

	warmth := (r + g*0.7) / (r + g + b + 1)
	coolness := b / (r + g + b + 1)

	return clamp(warmth-coolness+0.5, 0, 1)
}
