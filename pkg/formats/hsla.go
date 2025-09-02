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

// RGBAToHSLA converts a color.RGBA to HSLA color space.
func RGBAToHSLA(c color.RGBA) HSLA {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0
	a := float64(c.A) / 255.0

	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	delta := max - min

	l := (max + min) / 2.0

	if delta == 0 {
		return NewHSLA(0, 0, l, a)
	}

	var s float64
	if l < 0.5 {
		s = delta / (max + min)
	} else {
		s = delta / (2.0 - max - min)
	}

	var h float64
	switch max {
	case r:
		h = (g - b) / delta
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/delta + 2
	case b:
		h = (r-g)/delta + 4
	}

	h *= 60

	return NewHSLA(h, s, l, a)
}

// HSLAtoRGBA converts HSLA color space to color.RGBA.
func HSLAToRGBA(c HSLA) color.RGBA {
	// Normalize hue to [0, 360)
	h := math.Mod(c.H, 360)
	if h < 0 {
		h += 360
	}

	s := clamp(c.S, 0, 1)
	l := clamp(c.L, 0, 1)
	a := clamp(c.A, 0, 1)

	if s == 0 {
		gray := uint8(math.Round(l * 255))
		return color.RGBA{
			R: gray,
			G: gray,
			B: gray,
			A: uint8(math.Round(a * 255)),
		}
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	hk := h / 360.0

	tr := hk + 1.0/3.0
	tg := hk
	tb := hk - 1.0/3.0

	tr = normalizeHue(tr)
	tg = normalizeHue(tg)
	tb = normalizeHue(tb)

	r := hueToRGB(p, q, tr)
	g := hueToRGB(p, q, tg)
	b := hueToRGB(p, q, tb)

	return color.RGBA{
		R: uint8(math.Round(r * 255)),
		G: uint8(math.Round(g * 255)),
		B: uint8(math.Round(b * 255)),
		A: uint8(math.Round(a * 255)),
	}
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

// ToRGBA converts the HSLA color to color.RGBA.
// Convenience method that calls HSLAToRGBA.
func (c HSLA) ToRGBA() color.RGBA {
	return HSLAToRGBA(c)
}

// hueToRGB converts a hue value to RGB using the HSL algorithm.
// Used internally by HSLAToRGBA for color space conversion.
// Parameters p, q are intermediate values, t is the normalized hue component.
func hueToRGB(p, q, t float64) float64 {
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 0.5 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

// normalizeHue ensures hue values stay within [0-1] range for internal calculations.
// Used during HSL to RGB conversion to handle hue wraparound.
func normalizeHue(h float64) float64 {
	if h < 0 {
		return h + 1
	}
	if h > 1 {
		return h - 1
	}
	return h
}

// clamp constrains a value to the specified range [min, max].
// Used throughout the package to ensure color component values stay within valid bounds.
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
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
