// Package color provides color manipulation methods for HSL-based transformations.
// All methods preserve immutability by returning new Color instances rather than
// modifying the receiver.
package color

import "math"

// Lighten increases the lightness of the color by the specified amount (0.0-1.0).
// The amount is clamped to valid range and applied in HSL color space.
// Returns a new Color instance, preserving the original color's immutability.
func (c *Color) Lighten(amount float64) *Color {
	h, s, l := c.HSL()
	return NewHSLA(h, s, clamp(l+amount, 0, 1), c.Alpha())
}

// Darken decreases the lightness of the color by the specified amount (0.0-1.0).
// The amount is clamped to valid range and applied in HSL color space.
func (c *Color) Darken(amount float64) *Color {
	h, s, l := c.HSL()
	return NewHSLA(h, s, clamp(l-amount, 0, 1), c.Alpha())
}

// Saturate increases the saturation of the color by the specified amount (0.0-1.0).
// The amount is clamped to valid range and applied in HSL color space.
func (c *Color) Saturate(amount float64) *Color {
	h, s, l := c.HSL()
	return NewHSLA(h, clamp(s+amount, 0, 1), l, c.Alpha())
}

// Desaturate decreases the saturation of the color by the specified amount (0.0-1.0).
// The amount is clamped to valid range and applied in HSL color space.
func (c *Color) Desaturate(amount float64) *Color {
	h, s, l := c.HSL()
	return NewHSLA(h, clamp(s-amount, 0, 1), l, c.Alpha())
}

// RotateHue rotates the hue by the specified number of degrees.
// Positive values rotate clockwise, negative values rotate counterclockwise.
// The rotation wraps around the full 360-degree hue circle.
func (c *Color) RotateHue(degrees float64) *Color {
	h, s, l := c.HSL()
	return NewHSLA(math.Mod(h+(degrees/360.0)+1, 1), s, l, c.Alpha())
}

// ToGrayscale converts the color to grayscale by setting saturation to 0.
// Preserves the original hue and lightness values for potential future use.
func (c *Color) ToGrayscale() *Color {
	h, _, l := c.HSL()
	return NewHSLA(h, 0, l, c.Alpha())
}

// AdjustLightness multiplies the current lightness by the specified factor.
// Factor > 1.0 makes the color lighter, factor < 1.0 makes it darker.
// The result is clamped to the valid lightness range (0.0-1.0).
func (c *Color) AdjustLightness(factor float64) *Color {
	h, s, l := c.HSL()
	return NewHSLA(h, s, clamp(l*factor, 0, 1), c.Alpha())
}

// AdjustSaturation multiplies the current saturation by the specified factor.
// Factor > 1.0 makes the color more saturated, factor < 1.0 makes it less saturated.
// The result is clamped to the valid saturation range (0.0-1.0).
func (c *Color) AdjustSaturation(factor float64) *Color {
	h, s, l := c.HSL()
	return NewHSLA(h, clamp(s*factor, 0, 1), l, c.Alpha())
}

// Mix blends this color with another color by the specified amount (0.0-1.0).
// Amount 0.0 returns the original color, 1.0 returns the mixed color.
// Mixing is performed in RGB color space with linear interpolation.
func (c *Color) Mix(a *Color, amount float64) *Color {
	amount = clamp(amount, 0, 1)

	rC, gC, bC, aC := float64(c.R), float64(c.G), float64(c.B), float64(c.A)
	rA, gA, bA, aA := float64(a.R), float64(a.G), float64(a.B), float64(a.A)

	rR := uint8(rC*(1-amount) + rA*amount)
	gR := uint8(gC*(1-amount) + gA*amount)
	bR := uint8(bC*(1-amount) + bA*amount)
	aR := uint8(aC*(1-amount) + aA*amount)

	return &Color{R: rR, G: gR, B: bR, A: aR}
}

// Invert returns the color with inverted RGB values (255 - original).
// This creates the photographic negative effect while preserving alpha.
func (c *Color) Invert() *Color {
	return &Color{
		R: 255 - c.R,
		G: 255 - c.G,
		B: 255 - c.B,
		A: c.A,
	}
}

// Complement returns the complementary color by rotating hue 180 degrees.
// This creates the color opposite on the color wheel, useful for contrast.
func (c *Color) Complement() *Color {
	return c.RotateHue(180)
}
