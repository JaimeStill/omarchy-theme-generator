package color

// hslaCache stores computed HSLA values in normalized [0,1] range.
// This struct is used internally by Color to cache expensive HSL conversion results.
type hslaCache struct {
	H, S, L, A float64
}

// NewHSL creates a new Color from HSL values with full opacity (alpha = 1.0).
// All parameters should be in range [0.0, 1.0].
func NewHSL(h, s, l float64) *Color {
	r, g, b := hslToRGB(h, s, l)
	return NewRGB(r, g, b)
}

// NewHSLA creates a new Color from HSLA values.
// All parameters should be in range [0.0, 1.0] and will be clamped if outside this range.
func NewHSLA(h, s, l, a float64) *Color {
	r, g, b := hslToRGB(h, s, l)
	return NewRGBA(r, g, b, a)
}

// HSLA returns the color's hue, saturation, lightness, and alpha values.
// All returned values are in the range [0.0, 1.0].
// Values are computed once per Color instance and cached for performance.
func (c *Color) HSLA() (h, s, l, a float64) {
	c.hslaOnce.Do(func() {
		h, s, l := rgbToHSL(c.R, c.G, c.B)
		c.hsla = &hslaCache{H: h, S: s, L: l, A: toAlpha(c.A)}
	})

	return c.hsla.H, c.hsla.S, c.hsla.L, c.hsla.A
}

// HSL returns the color's hue, saturation, and lightness values.
// This is a convenience method that calls HSLA() and discards the alpha value.
func (c *Color) HSL() (h, s, l float64) {
	h, s, l, _ = c.HSLA()
	return h, s, l
}
