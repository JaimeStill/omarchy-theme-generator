package formats

import (
	"fmt"
	"image/color"
	"strings"
)

// ToHex converts a color.RGBA to a hex color string in the format #RRGGBB.
// Alpha channel is ignored. All hex digits are uppercase.
func ToHex(c color.RGBA) string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

// ToHexA converts a color.RGBA to a hex color string with alpha in the format #RRGGBBAA.
// All hex digits are uppercase.
func ToHexA(c color.RGBA) string {
	return fmt.Sprintf("#%02X%02X%02X%02X", c.R, c.G, c.B, c.A)
}

// ParseHex parses a hex color string and returns a color.RGBA.
// Supports multiple hex formats:
//   - 3 digits: #RGB (expanded to #RRGGBB)
//   - 4 digits: #RGBA (expanded to #RRGGBBAA)
//   - 6 digits: #RRGGBB (alpha defaults to 255)
//   - 8 digits: #RRGGBBAA
//
// The leading # is optional. Case-insensitive hex digits are supported.
func ParseHex(hex string) (color.RGBA, error) {
	hex = strings.TrimPrefix(hex, "#")

	for _, c := range hex {
		if !isHexChar(c) {
			return color.RGBA{}, fmt.Errorf("invalid hex character: %c", c)
		}
	}

	var r, g, b, a uint8
	a = 255

	switch len(hex) {
	case 3:
		_, err := fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
		if err != nil {
			return color.RGBA{}, fmt.Errorf("invalid hex color: %s", hex)
		}

		r = r*16 + r
		g = g*16 + g
		b = b*16 + b
	case 4:
		_, err := fmt.Sscanf(hex, "%1x%1x%1x%1x", &r, &g, &b, &a)
		if err != nil {
			return color.RGBA{}, fmt.Errorf("invalid hex color: %s", hex)
		}

		r = r*16 + r
		g = g*16 + g
		b = b*16 + b
		a = a*16 + a
	case 6:
		_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
		if err != nil {
			return color.RGBA{}, fmt.Errorf("invalid hex color: %s", hex)
		}
	case 8:
		_, err := fmt.Sscanf(hex, "%02x%02x%02x%02x", &r, &g, &b, &a)
		if err != nil {
			return color.RGBA{}, fmt.Errorf("invalid hex color: %s", hex)
		}
	default:
		return color.RGBA{}, fmt.Errorf("invalid hex color length: %d", len(hex))
	}

	return color.RGBA{R: r, G: g, B: b, A: a}, nil
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

// ParseHexToHSLA parses a hex color string and returns an HSLA color.
// Accepts the same hex formats as ParseHex and converts the result to HSLA color space.
func ParseHexToHSLA(hex string) (HSLA, error) {
	rgba, err := ParseHex(hex)
	if err != nil {
		return HSLA{}, err
	}
	return RGBAToHSLA(rgba), nil
}

// isHexChar returns true if the rune is a valid hexadecimal digit (0-9, a-f, A-F).
func isHexChar(c rune) bool {
	return (c >= '0' && c <= '9') ||
		(c >= 'a' && c <= 'f') ||
		(c >= 'A' && c <= 'F')
}
