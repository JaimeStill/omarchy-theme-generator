package formats_test

import (
	"image/color"
	"math"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func TestRGBAToHSLA(t *testing.T) {
	testCases := []struct {
		name      string
		input     color.RGBA
		expected  formats.HSLA
		tolerance float64
	}{
		{
			name:      "Pure Red",
			input:     color.RGBA{R: 255, G: 0, B: 0, A: 255},
			expected:  formats.NewHSLA(0, 1.0, 0.5, 1.0),
			tolerance: 0.01,
		},
		{
			name:      "Pure Green",
			input:     color.RGBA{R: 0, G: 255, B: 0, A: 255},
			expected:  formats.NewHSLA(120, 1.0, 0.5, 1.0),
			tolerance: 0.01,
		},
		{
			name:      "Pure Blue",
			input:     color.RGBA{R: 0, G: 0, B: 255, A: 255},
			expected:  formats.NewHSLA(240, 1.0, 0.5, 1.0),
			tolerance: 0.01,
		},
		{
			name:      "White",
			input:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expected:  formats.NewHSLA(0, 0.0, 1.0, 1.0),
			tolerance: 0.01,
		},
		{
			name:      "Black",
			input:     color.RGBA{R: 0, G: 0, B: 0, A: 255},
			expected:  formats.NewHSLA(0, 0.0, 0.0, 1.0),
			tolerance: 0.01,
		},
		{
			name:      "Gray",
			input:     color.RGBA{R: 128, G: 128, B: 128, A: 255},
			expected:  formats.NewHSLA(0, 0.0, 0.502, 1.0),
			tolerance: 0.01,
		},
		{
			name:      "Orange",
			input:     color.RGBA{R: 255, G: 165, B: 0, A: 255},
			expected:  formats.NewHSLA(38.8, 1.0, 0.5, 1.0),
			tolerance: 0.5,
		},
		{
			name:      "Semi-transparent Purple",
			input:     color.RGBA{R: 128, G: 0, B: 128, A: 128},
			expected:  formats.NewHSLA(300, 1.0, 0.251, 0.502),
			tolerance: 0.01,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.RGBAToHSLA(tc.input)

			if !floatEquals(result.H, tc.expected.H, tc.tolerance) {
				t.Errorf("Hue mismatch: expected %.2f, got %.2f", tc.expected.H, result.H)
			}
			if !floatEquals(result.S, tc.expected.S, tc.tolerance) {
				t.Errorf("Saturation mismatch: expected %.3f, got %.3f", tc.expected.S, result.S)
			}
			if !floatEquals(result.L, tc.expected.L, tc.tolerance) {
				t.Errorf("Lightness mismatch: expected %.3f, got %.3f", tc.expected.L, result.L)
			}
			if !floatEquals(result.A, tc.expected.A, tc.tolerance) {
				t.Errorf("Alpha mismatch: expected %.3f, got %.3f", tc.expected.A, result.A)
			}
		})
	}
}

func TestHSLAToRGBA(t *testing.T) {
	testCases := []struct {
		name     string
		input    formats.HSLA
		expected color.RGBA
	}{
		{
			name:     "Pure Red",
			input:    formats.NewHSLA(0, 1.0, 0.5, 1.0),
			expected: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:     "Pure Green",
			input:    formats.NewHSLA(120, 1.0, 0.5, 1.0),
			expected: color.RGBA{R: 0, G: 255, B: 0, A: 255},
		},
		{
			name:     "Pure Blue",
			input:    formats.NewHSLA(240, 1.0, 0.5, 1.0),
			expected: color.RGBA{R: 0, G: 0, B: 255, A: 255},
		},
		{
			name:     "White",
			input:    formats.NewHSLA(0, 0.0, 1.0, 1.0),
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:     "Black",
			input:    formats.NewHSLA(0, 0.0, 0.0, 1.0),
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		},
		{
			name:     "Gray",
			input:    formats.NewHSLA(0, 0.0, 0.5, 1.0),
			expected: color.RGBA{R: 128, G: 128, B: 128, A: 255},
		},
		{
			name:     "Orange",
			input:    formats.NewHSLA(30, 1.0, 0.5, 1.0),
			expected: color.RGBA{R: 255, G: 128, B: 0, A: 255},
		},
		{
			name:     "Semi-transparent Cyan",
			input:    formats.NewHSLA(180, 1.0, 0.5, 0.5),
			expected: color.RGBA{R: 0, G: 255, B: 255, A: 128},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.HSLAToRGBA(tc.input)

			if !colorEquals(result, tc.expected, 1) {
				t.Errorf("Expected RGBA(%d, %d, %d, %d), got RGBA(%d, %d, %d, %d)",
					tc.expected.R, tc.expected.G, tc.expected.B, tc.expected.A,
					result.R, result.G, result.B, result.A)
			}
		})
	}
}

func TestRoundTripConversion(t *testing.T) {
	// Test that converting RGBA -> HSLA -> RGBA preserves the color
	testColors := []color.RGBA{
		{R: 255, G: 0, B: 0, A: 255},     // Red
		{R: 0, G: 255, B: 0, A: 255},     // Green
		{R: 0, G: 0, B: 255, A: 255},     // Blue
		{R: 255, G: 255, B: 0, A: 255},   // Yellow
		{R: 255, G: 0, B: 255, A: 255},   // Magenta
		{R: 0, G: 255, B: 255, A: 255},   // Cyan
		{R: 128, G: 128, B: 128, A: 255}, // Gray
		{R: 255, G: 165, B: 0, A: 255},   // Orange
		{R: 75, G: 0, B: 130, A: 255},    // Indigo
	}

	for _, original := range testColors {
		hsla := formats.RGBAToHSLA(original)
		result := formats.HSLAToRGBA(hsla)

		if !colorEquals(result, original, 1) {
			t.Errorf("Round trip failed for RGBA(%d, %d, %d, %d): got RGBA(%d, %d, %d, %d)",
				original.R, original.G, original.B, original.A,
				result.R, result.G, result.B, result.A)
		}
	}
}

func TestHSLAMethods(t *testing.T) {
	t.Run("WithAlpha", func(t *testing.T) {
		hsla := formats.NewHSLA(180, 0.5, 0.5, 1.0)
		modified := hsla.WithAlpha(0.5)

		if modified.A != 0.5 {
			t.Errorf("Expected alpha 0.5, got %f", modified.A)
		}
		// Original should be unchanged
		if hsla.A != 1.0 {
			t.Errorf("Original HSLA was modified")
		}
	})

	t.Run("ToRGBA", func(t *testing.T) {
		hsla := formats.NewHSLA(120, 1.0, 0.5, 1.0)
		rgba := hsla.ToRGBA()
		expected := color.RGBA{R: 0, G: 255, B: 0, A: 255}

		if !colorEquals(rgba, expected, 1) {
			t.Errorf("Expected %v, got %v", expected, rgba)
		}
	})

	t.Run("RGBA Interface", func(t *testing.T) {
		hsla := formats.NewHSLA(0, 1.0, 0.5, 1.0)
		r, g, b, a := hsla.RGBA()

		// color.Color interface returns 16-bit values
		if r != 0xffff || g != 0 || b != 0 || a != 0xffff {
			t.Errorf("RGBA() interface method incorrect: got (%d, %d, %d, %d)", r, g, b, a)
		}
	})
}

func TestColorHelpers(t *testing.T) {
	t.Run("WithAlpha on RGBA", func(t *testing.T) {
		original := color.RGBA{R: 255, G: 128, B: 0, A: 255}
		modified := formats.WithAlpha(original, 0.5)

		if modified.A != 128 {
			t.Errorf("Expected alpha 128, got %d", modified.A)
		}
		if modified.R != original.R || modified.G != original.G || modified.B != original.B {
			t.Errorf("RGB values changed unexpectedly")
		}
	})

	t.Run("GetAlpha", func(t *testing.T) {
		c := color.RGBA{R: 255, G: 128, B: 0, A: 128}
		alpha := formats.GetAlpha(c)

		if !floatEquals(alpha, 0.502, 0.01) {
			t.Errorf("Expected alpha ~0.502, got %f", alpha)
		}
	})
}

// Helper functions
func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func colorEquals(a, b color.RGBA, tolerance uint8) bool {
	return absDiff(a.R, b.R) <= tolerance &&
		absDiff(a.G, b.G) <= tolerance &&
		absDiff(a.B, b.B) <= tolerance &&
		absDiff(a.A, b.A) <= tolerance
}

func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}
