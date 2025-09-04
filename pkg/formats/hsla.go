package formats

import (
	"image/color"
	"math"
)

// HSLA represents a color in HSLA color space.
// H is hue in degrees [0-360)
// S is saturation [0-1]
// L is lightness [0-1], A is alpha [0-1]
type HSLA struct {
	H float64
	S float64
	L float64
	A float64
}

// RGBA converts HSLA to the color.Color interface.
// This allows HSLA to satisfy the color.Color interface from the standard library.
func (c HSLA) RGBA() (r, g, b, a uint32) {
	rgba := HSLAToRGBA(c)
	r = uint32(rgba.R) * 0x101
	g = uint32(rgba.G) * 0x101
	b = uint32(rgba.B) * 0x101
	a = uint32(rgba.A) * 0x101
	return
}

// NewHSLA creates a new HSLA color with normalized values.
// Hue is normalized to [0-360), saturation/lightness/alpha are clamped to [0-1].
func NewHSLA(h, s, l, a float64) HSLA {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	return HSLA{
		H: h,
		S: clamp(s, 0, 1),
		L: clamp(l, 0, 1),
		A: clamp(a, 0, 1),
	}
}

// NewHSL creates a new HSLA color with full opacity (alpha = 1.0).
// Convenience function for creating opaque HSLA colors.
func NewHSL(h, s, l float64) HSLA {
	return NewHSLA(h, s, l, 1.0)
}

// ToRGBA converts the HSLA color to color.RGBA.
// Convenience method that calls HSLAToRGBA.
func (c HSLA) ToRGBA() color.RGBA {
	return HSLAToRGBA(c)
}

// WithAlpha returns a new HSLA color with the specified alpha value.
// The alpha value is clamped to [0-1]. Other color components remain unchanged.
func (c HSLA) WithAlpha(alpha float64) HSLA {
	return HSLA{
		H: c.H,
		S: c.S,
		L: c.L,
		A: clamp(alpha, 0, 1),
	}
}

// WithAlpha returns a new color.RGBA with the specified alpha value.
// The alpha value is clamped to [0-1] and converted to [0-255] range.
// RGB components remain unchanged.
func WithAlpha(c color.RGBA, alpha float64) color.RGBA {
	return color.RGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: uint8(math.Round(clamp(alpha, 0, 1) * 255)),
	}
}

// GetAlpha extracts the alpha channel from a color.RGBA as a float64 in range [0-1].
// Converts from [0-255] uint8 range to [0-1] float64 range.
func GetAlpha(c color.RGBA) float64 {
	return float64(c.A) / 255.0
}
