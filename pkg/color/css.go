package color

import "fmt"

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
