package color

import "fmt"

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
