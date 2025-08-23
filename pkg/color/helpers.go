package color

import "math"

// clamp ensures value is within the specified range [min, max].
func clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

// rgbToHSL converts RGB [0,255] to HSL [0,1] following CSS Color Module Level 3 specification.
func rgbToHSL(r, g, b uint8) (h, s, l float64) {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	max := math.Max(math.Max(rf, gf), bf)
	min := math.Min(math.Min(rf, gf), bf)

	l = (max + min) / 2

	if max == min {
		h, s = 0, 0
	} else {
		delta := max - min

		if l > 0.5 {
			s = delta / (2 - max - min)
		} else {
			s = delta / (max + min)
		}

		switch max {
		case rf:
			h = (gf - bf) / delta
			if gf < bf {
				h += 6
			}
		case gf:
			h = (bf-rf)/delta + 2
		case bf:
			h = (rf-gf)/delta + 4
		}

		h /= 6
	}

	return h, s, l
}

// hslToRGB converts HSL [0,1] to RGB [0,255] following CSS Color Module Level 3 specification.
func hslToRGB(h, s, l float64) (r, g, b uint8) {
	if s == 0 {
		gray := uint8(l * 255)
		return gray, gray, gray
	}

	var c2 float64
	if l < 0.5 {
		c2 = l * (1 + s)
	} else {
		c2 = l + s - l*s
	}
	c1 := 2*l - c2

	r = uint8(clamp(hueToRGB(c1, c2, h+1.0/3.0), 0, 1) * 255)
	g = uint8(clamp(hueToRGB(c1, c2, h), 0, 1) * 255)
	b = uint8(clamp(hueToRGB(c1, c2, h-1.0/3.0), 0, 1) * 255)

	return r, g, b
}

// hueToRGB is a helper function for HSL to RGB conversion.
func hueToRGB(c1, c2, hue float64) float64 {
	if hue < 0 {
		hue += 1
	}
	if hue > 1 {
		hue -= 1
	}

	if hue < 1.0/6.0 {
		return c1 + (c2-c1)*6*hue
	}
	if hue < 1.0/2.0 {
		return c2
	}
	if hue < 2.0/3.0 {
		return c1 + (c2-c1)*(2.0/3.0-hue)*6
	}
	return c1
}

// toAlpha converts byte (0-255) to alpha (0.0-1.0) with rounding for clean output.
func toAlpha(b uint8) float64 {
	return roundAlpha(float64(b) / 255.0)
}

// alphaToByte converts alpha (0.0-1.0) to byte (0-255) with clamping.
func alphaToByte(a float64) uint8 {
	return uint8(clamp(a, 0, 1) * 255)
}

// roundAlpha rounds alpha values to 3 decimal places for consistent display.
func roundAlpha(a float64) float64 {
	return math.Round(a * 1000) / 1000
}
