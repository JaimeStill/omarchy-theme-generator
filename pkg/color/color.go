package color

import (
	"fmt"
	"sync"
)

// Color represents an RGBA color with lazy-cached HSLA conversion.
// Colors use native RGBA storage for optimal image processing performance
// while supporting efficient HSL manipulation through thread-safe caching.
//
// The Color type is safe for concurrent use. HSLA values are computed once
// per instance using sync.Once and cached for subsequent access.
//
// Alpha values are stored as uint8 (0-255) internally but exposed as float64
// (0.0-1.0) in all public methods for consistency with CSS specifications.
type Color struct {
	R, G, B, A uint8
	hsla       *hslaCache
	hslaOnce   sync.Once
}

// hslaCache stores computed HSLA values in normalized [0,1] range.
// This struct is used internally by Color to cache expensive HSL conversion results.
type hslaCache struct {
	H, S, L, A float64
}

// NewRGB creates a new Color from RGB values with full opacity (alpha = 1.0).
func NewRGB(r, g, b uint8) *Color {
	return &Color{R: r, G: g, B: b, A: 255}
}

// NewRGBA creates a new Color from RGBA values.
// Alpha should be in range [0.0, 1.0] and will be clamped if outside this range.
func NewRGBA(r, g, b uint8, a float64) *Color {
	return &Color{R: r, G: g, B: b, A: alphaToByte(a)}
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

// HEXA returns the color as an 8-digit hexadecimal string including alpha.
// Format: #RRGGBBAA where each component is two lowercase hexadecimal digits.
func (c *Color) HEXA() string {
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}

// HEX returns the color as a 6-digit hexadecimal string, ignoring alpha.
// Format: #RRGGBB where each component is two lowercase hexadecimal digits.
func (c *Color) HEX() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// RGBA returns the color's red, green, blue, and alpha components as uint8 values.
// Each component is in the range [0, 255].
func (c *Color) RGBA() (r, g, b, a uint8) {
	return c.R, c.G, c.B, c.A
}

// RGB returns the color's red, green, and blue components as uint8 values.
// This is a convenience method that calls RGBA() and discards the alpha value.
func (c *Color) RGB() (r, g, b uint8) {
	return c.R, c.G, c.B
}

// CSSRGBA returns the color as a CSS rgba() function string.
// Format: rgba(r, g, b, a) where RGB are integers [0-255] and alpha is [0.000-1.000].
func (c *Color) CSSRGBA() string {
	alpha := toAlpha(c.A)
	return fmt.Sprintf("rgba(%d, %d, %d, %.3f)", c.R, c.G, c.B, alpha)
}

// CSSRGB returns the color as a CSS rgb() function string, ignoring alpha.
// Format: rgb(r, g, b) where each component is an integer [0-255].
func (c *Color) CSSRGB() string {
	return fmt.Sprintf("rgb(%d, %d, %d)", c.R, c.G, c.B)
}

// CSSHSLA returns the color as a CSS hsla() function string.
// Format: hsla(h, s%, l%, a) where h is degrees [0.0-360.0], s and l are percentages [0.0-100.0], and a is alpha [0.000-1.000].
func (c *Color) CSSHSLA() string {
	h, s, l, a := c.HSLA()
	return fmt.Sprintf("hsla(%.1f, %.1f, %.1f, %.3f)", h*360, s*100, l*100, a)
}

// CSSHSL returns the color as a CSS hsl() function string, ignoring alpha.
// Format: hsl(h, s%, l%) where h is degrees [0.0-360.0] and s, l are percentages [0.0-100.0].
func (c *Color) CSSHSL() string {
	h, s, l := c.HSL()
	return fmt.Sprintf("hsl(%.1f, %.1f, %.1f)", h*360, s*100, l*100)
}

// WithAlpha returns a new Color with the specified alpha value.
// Alpha will be clamped to the range [0.0, 1.0] if outside this range.
// The original Color is unchanged (value semantics).
func (c *Color) WithAlpha(a float64) Color {
	alpha := alphaToByte(a)
	return Color{R: c.R, G: c.G, B: c.B, A: alpha}
}

// Alpha returns the color's alpha (opacity) value as a float64 in range [0.0, 1.0].
// 0.0 represents fully transparent, 1.0 represents fully opaque.
func (c *Color) Alpha() float64 {
	return toAlpha(c.A)
}

// IsOpaque returns true if the color is fully opaque (alpha = 1.0).
func (c *Color) IsOpaque() bool {
	return c.A == 255
}

// IsTransparent returns true if the color is fully transparent (alpha = 0.0).
func (c *Color) IsTransparent() bool {
	return c.A == 0
}
