package color

import (
	"sync"
)

// Color represents an RGBA color with performance-optimized caching.
//
// Storage and Performance:
//   - Native RGBA uint8 storage for zero-allocation image processing
//   - Thread-safe HSLA conversion caching via sync.Once (computed once)
//   - Thread-safe luminance caching for contrast calculations
//   - Memory overhead: 32 bytes base + 48 bytes when caches populated
//
// The Color type is safe for concurrent use across goroutines. Expensive
// conversions (RGB→HSL, luminance calculation) are computed once and cached
// automatically on first access.
//
// Alpha values are stored internally as uint8 (0-255) but exposed as
// float64 (0.0-1.0) for CSS compatibility and mathematical operations.
type Color struct {
	R, G, B, A    uint8
	hsla          *hslaCache
	hslaOnce      sync.Once
	luminance     *luminanceCache
	luminanceOnce sync.Once
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
