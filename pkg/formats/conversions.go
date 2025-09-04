package formats

import (
	"image/color"
	"math"
)

func GetIlluminant(illuminant string) XYZ {
	switch illuminant {
	case "D50":
		return D50Illuminant
	default:
		return D65Illuminant
	}
}

// HSLAToHex converts an HSLA color to hex format #RRGGBB.
// Alpha channel is ignored in the output.
func HSLAToHex(h HSLA) string {
	return ToHex(HSLAToRGBA(h))
}

// HSLAToHexA converts an HSLA color to hex format with alpha #RRGGBBAA.
// Includes the alpha channel in the output.
func HSLAToHexA(h HSLA) string {
	return ToHexA(HSLAToRGBA(h))
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

func LABToRGBA(lab LAB) color.RGBA {
	return LABToRGBAWithIlluminant(lab, D65Illuminant)
}

func LABToRGBAWithIlluminant(lab LAB, illuminant XYZ) color.RGBA {
	xyz := LABToXYZ(lab, illuminant)
	return XYZToRGBA(xyz)
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

func RGBAToLAB(c color.RGBA) LAB {
	return RGBAToLABWithIlluminant(c, D65Illuminant)
}

func RGBAToLABWithIlluminant(c color.RGBA, illuminant XYZ) LAB {
	xyz := RGBAToXYZ(c)
	return XYZToLAB(xyz, illuminant)
}

func RGBAToLABWithSettings(c color.RGBA, illuminant string) LAB {
	return RGBAToLABWithIlluminant(c, GetIlluminant(illuminant))
}

func RGBAToXYZ(c color.RGBA) XYZ {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	r = inverseSRGBGamma(r)
	g = inverseSRGBGamma(g)
	b = inverseSRGBGamma(b)

	x := r*0.4124564 + g*0.3575761 + b*0.1804375
	y := r*0.2126729 + g*0.7151522 + b*0.0721750
	z := r*0.0193339 + g*0.1191920 + b*0.9503041

	return XYZ{
		X: x * 100.0,
		Y: y * 100.0,
		Z: z * 100.0,
	}
}

func XYZToLAB(xyz XYZ, illuminant XYZ) LAB {
	xn := xyz.X / illuminant.X
	yn := xyz.Y / illuminant.Y
	zn := xyz.Z / illuminant.Z

	fx := labTransform(xn)
	fy := labTransform(yn)
	fz := labTransform(zn)

	l := 116*fy - 16
	a := 500 * (fx - fy)
	b := 200 * (fy - fz)

	return LAB{L: l, A: a, B: b}
}

func LABToXYZ(lab LAB, illuminant XYZ) XYZ {
	fy := (lab.L + 16) / 116
	fx := lab.A/500 + fy
	fz := fy - lab.B/200

	x := inverseLABTransform(fx) * illuminant.X
	y := inverseLABTransform(fy) * illuminant.Y
	z := inverseLABTransform(fz) * illuminant.Z

	return XYZ{X: x, Y: y, Z: z}
}

func XYZToRGBA(xyz XYZ) color.RGBA {
	x := xyz.X / 100.0
	y := xyz.Y / 100.0
	z := xyz.Z / 100.0

	r := x*3.2404542 + y*(-1.5371385) + z*(-0.4985314)
	g := x*(-0.9692660) + y*1.8760108 + z*0.0415560
	b := x*0.0556434 + y*(-0.2040259) + z*1.0572252

	r = sRGBGamma(r)
	g = sRGBGamma(g)
	b = sRGBGamma(b)

	return color.RGBA{
		R: uint8(clamp(r*255.0, 0, 255)),
		G: uint8(clamp(g*255.0, 0, 255)),
		B: uint8(clamp(b*255.0, 0, 255)),
		A: 255,
	}
}

// clamp constrains a value to the specified range [min, max].
// Used throughout the package to ensure color component values stay within valid bounds.
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
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

func inverseLABTransform(t float64) float64 {
	if t > 0.206897 {
		return math.Pow(t, 3.0)
	}

	return 3 * (6.0 / 29.0) * (6.0 / 29.0) * (t - 4.0/29.0)
}

func inverseSRGBGamma(value float64) float64 {
	if value <= 0.04045 {
		return value / 12.92
	}
	return math.Pow((value+0.055)/1.055, 2.4)
}

func labTransform(t float64) float64 {
	if t > 0.008856 {
		return math.Pow(t, 1.0/3.0)
	}
	return (903.3*t + 16) / 116
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

func sRGBGamma(value float64) float64 {
	if value <= 0.0031308 {
		return 12.92 * value
	}
	return 1.055*math.Pow(value, 1.0/2.4) - 0.055
}
