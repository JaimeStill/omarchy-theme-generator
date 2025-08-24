// Package color provides CIE LAB color space conversion and Delta-E calculations.
// LAB color space is designed to be perceptually uniform, making it ideal for
// accurate color difference measurements in color science applications.
package color

import "math"

// LABColor represents a color in the CIE LAB color space.
// L* represents lightness (0-100), A* represents green-red axis (-128 to +127),
// B* represents blue-yellow axis (-128 to +127).
type LABColor struct {
	L, A, B float64
}

// ToLAB converts the color to CIE LAB color space using D65 illuminant.
// Conversion: sRGB → linear RGB → XYZ → LAB with proper gamma correction.
// Uses standard CIE transformation matrices for accurate color science results.
func (c *Color) ToLAB() LABColor {
	r := toLinearRGB(float64(c.R) / 255.0)
	g := toLinearRGB(float64(c.G) / 255.0)
	b := toLinearRGB(float64(c.B) / 255.0)

	x := r*0.4124564 + g*0.3575761 + b*0.1804375
	y := r*0.2126729 + g*0.7151522 + b*0.0721750
	z := r*0.0193339 + g*0.1191920 + b*0.9503041

	return xyzToLAB(x, y, z)
}

// DeltaE76 calculates CIE76 Delta-E color difference (simple Euclidean distance).
// Values: 0-1 (identical), 1-2.3 (barely perceptible), 2.3-5 (noticeable).
// Fast calculation but less accurate than CIE94 for perceptual uniformity.
func (c *Color) DeltaE76(a *Color) float64 {
	lab1 := c.ToLAB()
	lab2 := a.ToLAB()

	dl := lab1.L - lab2.L
	da := lab1.A - lab2.A
	db := lab1.B - lab2.B

	return math.Sqrt(dl*dl + da*da + db*db)
}

// DeltaE94 calculates CIE94 Delta-E color difference with perceptual weighting.
// More accurate than CIE76 for human color perception, accounting for chroma and hue.
// Industry standard for color matching and quality control applications.
func (c *Color) DeltaE94(a *Color) float64 {
	lab1 := c.ToLAB()
	lab2 := a.ToLAB()

	dl := lab1.L - lab2.L
	da := lab1.A - lab2.A
	db := lab1.B - lab2.B

	c1 := math.Sqrt(lab1.A*lab1.A + lab1.B*lab1.B)
	c2 := math.Sqrt(lab2.A*lab2.A + lab2.B*lab2.B)
	dc := c1 - c2

	dh := math.Sqrt(da*da + db*db - dc*dc)

	kL, kC, kH := 1.0, 1.0, 1.0
	k1, k2 := 0.045, 0.015

	sL := 1.0
	sC := 1 + k1*c1
	sH := 1 + k2*c1

	term1 := dl / (kL * sL)
	term2 := dc / (kC * sC)
	term3 := dh / (kH * sH)

	return math.Sqrt(term1*term1 + term2*term2 + term3*term3)
}

// IsPerceptuallyIdentical checks if colors are indistinguishable (Delta-E ≤ 1.0).
// Based on CIE76 Delta-E threshold for human color discrimination limits.
func (c *Color) IsPerceptuallyIdentical(a *Color) bool {
	return c.DeltaE76(a) <= 1.0
}

// IsPerceptuallySimilar checks if colors are barely distinguishable (Delta-E ≤ 2.3).
// Uses established threshold for just-noticeable color differences in controlled conditions.
func (c *Color) IsPerceptuallySimilar(a *Color) bool {
	return c.DeltaE76(a) <= 2.3
}

// FindPerceptuallyClosest finds the most perceptually similar color using Delta-E CIE76.
// Returns the index of the closest color and its Delta-E value.
// Returns (-1, 0) for empty candidate sets.
func (c *Color) FindPerceptuallyClosest(candidates []*Color) (index int, deltaE float64) {
	if len(candidates) == 0 {
		return -1, 0
	}

	minDeltaE := math.Inf(1)
	closestIndex := 0

	for i, a := range candidates {
		de := c.DeltaE76(a)
		if de < minDeltaE {
			minDeltaE = de
			closestIndex = i
		}
	}

	return closestIndex, minDeltaE
}

// IsPerceptuallyDistinct checks if color is sufficiently different from all others.
// Uses Delta-E CIE76 with custom threshold for palette generation and color selection.
// Typical thresholds: 1.0 (identical), 2.3 (just noticeable), 5.0 (clearly different).
func (c *Color) IsPerceptuallyDistinct(colors []*Color, threshold float64) bool {
	for _, a := range colors {
		if c.DeltaE76(a) <= threshold {
			return false
		}
	}

	return true
}

// xyzToLAB converts CIE XYZ coordinates to LAB color space.
// Uses D65 illuminant white point (0.95047, 1.0, 1.08883) for standard daylight conditions.
// Applies LAB transformation function for perceptually uniform color space.
func xyzToLAB(x, y, z float64) LABColor {
	x = x / 0.95047
	y = y / 1.0
	z = z / 1.08883

	fx := labTransform(x)
	fy := labTransform(y)
	fz := labTransform(z)

	l := 116*fy - 16
	a := 500 * (fx - fy)
	b := 200 * (fy - fz)

	return LABColor{L: l, A: a, B: b}
}

// labTransform applies the CIE LAB transformation function.
// Handles the cube root for t > 0.008856, linear approximation otherwise.
// Essential for maintaining perceptual uniformity in LAB color space.
func labTransform(t float64) float64 {
	if t > 0.008856 {
		return math.Pow(t, 1.0/3.0)
	}

	return (7.787*t + 16.0/116.0)
}
