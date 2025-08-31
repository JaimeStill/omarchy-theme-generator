package color

import (
	"fmt"
	"strconv"
	"strings"
)

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

// ParseHEXA parses an 8-digit hexadecimal color string including alpha channel.
// Accepts formats: #RRGGBBAA, RRGGBBAA
// Returns error if the string is not a valid HEXA color format.
func ParseHEXA(hexa string) (*Color, error) {
	// Remove # prefix if present
	hexa = strings.TrimPrefix(hexa, "#")

	// Validate length
	if len(hexa) != 8 {
		return nil, fmt.Errorf("invalid HEXA format: expected 8 characters, got %d", len(hexa))
	}

	// Parse each component
	r, err := strconv.ParseUint(hexa[0:2], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid red component in HEXA: %w", err)
	}

	g, err := strconv.ParseUint(hexa[2:4], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid green component in HEXA: %w", err)
	}

	b, err := strconv.ParseUint(hexa[4:6], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid blue component in HEXA: %w", err)
	}

	a, err := strconv.ParseUint(hexa[6:8], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid alpha component in HEXA: %w", err)
	}

	return &Color{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}, nil
}

// ParseHEX parses a 6-digit hexadecimal color string with full alpha (255).
// Accepts formats: #RRGGBB, RRGGBB
// Returns error if the string is not a valid HEX color format.
func ParseHEX(hex string) (*Color, error) {
	// Remove # prefix if present
	hex = strings.TrimPrefix(hex, "#")

	// Validate length
	if len(hex) != 6 {
		return nil, fmt.Errorf("invalid HEX format: expected 6 characters, got %d", len(hex))
	}

	// Parse each component
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid red component in HEX: %w", err)
	}

	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid green component in HEX: %w", err)
	}

	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid blue component in HEX: %w", err)
	}

	return &Color{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255, // Full opacity for 6-digit hex
	}, nil
}

// ParseHexString automatically parses either HEX (6-digit) or HEXA (8-digit) format.
// This is a convenience function that determines the format based on string length.
func ParseHexString(hexStr string) (*Color, error) {
	hexStr = strings.TrimPrefix(hexStr, "#")

	switch len(hexStr) {
	case 6:
		return ParseHEX("#" + hexStr)
	case 8:
		return ParseHEXA("#" + hexStr)
	default:
		return nil, fmt.Errorf("invalid hex format: expected 6 or 8 characters, got %d", len(hexStr))
	}
}
