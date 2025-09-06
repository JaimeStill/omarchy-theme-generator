package formats_test

import (
	"image/color"
	"math"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

// TestNewHSLA tests the NewHSLA constructor function
func TestNewHSLA(t *testing.T) {
	testCases := []struct {
		name     string
		h, s, l, a float64
		expected formats.HSLA
		description string
	}{
		{
			name:     "Standard white",
			h: 0, s: 0, l: 1, a: 1,
			expected: formats.HSLA{H: 0, S: 0, L: 1, A: 1},
			description: "White should be H=0, S=0, L=1, A=1",
		},
		{
			name:     "Standard black",
			h: 0, s: 0, l: 0, a: 1,
			expected: formats.HSLA{H: 0, S: 0, L: 0, A: 1},
			description: "Black should be H=0, S=0, L=0, A=1",
		},
		{
			name:     "Pure red",
			h: 0, s: 1, l: 0.5, a: 1,
			expected: formats.HSLA{H: 0, S: 1, L: 0.5, A: 1},
			description: "Pure red should be H=0, S=1, L=0.5, A=1",
		},
		{
			name:     "Hue normalization positive",
			h: 390, s: 0.5, l: 0.5, a: 1,
			expected: formats.HSLA{H: 30, S: 0.5, L: 0.5, A: 1},
			description: "H=390 should normalize to H=30",
		},
		{
			name:     "Hue normalization negative",
			h: -30, s: 0.5, l: 0.5, a: 1,
			expected: formats.HSLA{H: 330, S: 0.5, L: 0.5, A: 1},
			description: "H=-30 should normalize to H=330",
		},
		{
			name:     "Saturation clamping high",
			h: 120, s: 1.5, l: 0.5, a: 1,
			expected: formats.HSLA{H: 120, S: 1, L: 0.5, A: 1},
			description: "S=1.5 should clamp to S=1",
		},
		{
			name:     "Saturation clamping low",
			h: 120, s: -0.2, l: 0.5, a: 1,
			expected: formats.HSLA{H: 120, S: 0, L: 0.5, A: 1},
			description: "S=-0.2 should clamp to S=0",
		},
		{
			name:     "Lightness clamping high",
			h: 240, s: 0.5, l: 1.2, a: 1,
			expected: formats.HSLA{H: 240, S: 0.5, L: 1, A: 1},
			description: "L=1.2 should clamp to L=1",
		},
		{
			name:     "Alpha clamping low",
			h: 180, s: 0.5, l: 0.5, a: -0.1,
			expected: formats.HSLA{H: 180, S: 0.5, L: 0.5, A: 0},
			description: "A=-0.1 should clamp to A=0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.NewHSLA(tc.h, tc.s, tc.l, tc.a)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Input HSLA: H=%.1f°, S=%.2f, L=%.2f, A=%.2f", tc.h, tc.s, tc.l, tc.a)
			t.Logf("Result HSLA: H=%.1f°, S=%.2f, L=%.2f, A=%.2f", result.H, result.S, result.L, result.A)
			t.Logf("Expected HSLA: H=%.1f°, S=%.2f, L=%.2f, A=%.2f", tc.expected.H, tc.expected.S, tc.expected.L, tc.expected.A)
			t.Logf("Description: %s", tc.description)

			// Check each component with appropriate tolerance
			tolerance := 0.0001
			if math.Abs(result.H-tc.expected.H) > tolerance {
				t.Errorf("H component mismatch: expected %.3f, got %.3f", tc.expected.H, result.H)
			}
			if math.Abs(result.S-tc.expected.S) > tolerance {
				t.Errorf("S component mismatch: expected %.3f, got %.3f", tc.expected.S, result.S)
			}
			if math.Abs(result.L-tc.expected.L) > tolerance {
				t.Errorf("L component mismatch: expected %.3f, got %.3f", tc.expected.L, result.L)
			}
			if math.Abs(result.A-tc.expected.A) > tolerance {
				t.Errorf("A component mismatch: expected %.3f, got %.3f", tc.expected.A, result.A)
			}

			t.Logf("✅ NewHSLA working correctly")
		})
	}
}

// TestNewHSL tests the NewHSL convenience constructor
func TestNewHSL(t *testing.T) {
	testCases := []struct {
		name     string
		h, s, l  float64
		expected formats.HSLA
		description string
	}{
		{
			name:     "Red HSL",
			h: 0, s: 1, l: 0.5,
			expected: formats.HSLA{H: 0, S: 1, L: 0.5, A: 1},
			description: "NewHSL should default alpha to 1",
		},
		{
			name:     "Green HSL",
			h: 120, s: 1, l: 0.5,
			expected: formats.HSLA{H: 120, S: 1, L: 0.5, A: 1},
			description: "Green HSL should have full alpha",
		},
		{
			name:     "Blue HSL",
			h: 240, s: 1, l: 0.5,
			expected: formats.HSLA{H: 240, S: 1, L: 0.5, A: 1},
			description: "Blue HSL should have full alpha",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.NewHSL(tc.h, tc.s, tc.l)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Input HSL: H=%.1f°, S=%.2f, L=%.2f", tc.h, tc.s, tc.l)
			t.Logf("Result HSLA: H=%.1f°, S=%.2f, L=%.2f, A=%.2f", result.H, result.S, result.L, result.A)
			t.Logf("Expected HSLA: H=%.1f°, S=%.2f, L=%.2f, A=%.2f", tc.expected.H, tc.expected.S, tc.expected.L, tc.expected.A)
			t.Logf("Description: %s", tc.description)

			// Verify the result
			if result != tc.expected {
				t.Errorf("NewHSL result mismatch:\nExpected: %+v\nGot: %+v", tc.expected, result)
			}

			t.Logf("✅ NewHSL working correctly")
		})
	}
}

// TestHSLA_RGBA tests the RGBA method (color.Color interface)
func TestHSLA_RGBA(t *testing.T) {
	testCases := []struct {
		name     string
		hsla     formats.HSLA
		description string
	}{
		{
			name:     "White",
			hsla:     formats.HSLA{H: 0, S: 0, L: 1, A: 1},
			description: "White should convert to maximum RGB values",
		},
		{
			name:     "Black",
			hsla:     formats.HSLA{H: 0, S: 0, L: 0, A: 1},
			description: "Black should convert to minimum RGB values",
		},
		{
			name:     "Pure red",
			hsla:     formats.HSLA{H: 0, S: 1, L: 0.5, A: 1},
			description: "Pure red should have high R component",
		},
		{
			name:     "Pure green",
			hsla:     formats.HSLA{H: 120, S: 1, L: 0.5, A: 1},
			description: "Pure green should have high G component",
		},
		{
			name:     "Pure blue",
			hsla:     formats.HSLA{H: 240, S: 1, L: 0.5, A: 1},
			description: "Pure blue should have high B component",
		},
		{
			name:     "Semi-transparent",
			hsla:     formats.HSLA{H: 180, S: 0.5, L: 0.5, A: 0.5},
			description: "Semi-transparent color should have proper alpha",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, g, b, a := tc.hsla.RGBA()

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("HSLA input: H=%.1f°, S=%.2f, L=%.2f, A=%.2f", tc.hsla.H, tc.hsla.S, tc.hsla.L, tc.hsla.A)
			t.Logf("RGBA result: R=%d, G=%d, B=%d, A=%d", r, g, b, a)
			t.Logf("RGBA normalized: R=%.3f, G=%.3f, B=%.3f, A=%.3f", 
				float64(r)/65535.0, float64(g)/65535.0, float64(b)/65535.0, float64(a)/65535.0)
			t.Logf("Description: %s", tc.description)

			// Verify values are in valid range (0-65535 for color.Color interface)
			if r > 65535 || g > 65535 || b > 65535 || a > 65535 {
				t.Errorf("RGBA components out of range: R=%d, G=%d, B=%d, A=%d", r, g, b, a)
			}

			t.Logf("✅ HSLA RGBA method working correctly")
		})
	}
}

// TestHSLA_ToRGBA tests the ToRGBA convenience method
func TestHSLA_ToRGBA(t *testing.T) {
	testCases := []struct {
		name     string
		hsla     formats.HSLA
		description string
	}{
		{
			name:     "White to RGBA",
			hsla:     formats.HSLA{H: 0, S: 0, L: 1, A: 1},
			description: "White HSLA should convert to RGBA(255,255,255,255)",
		},
		{
			name:     "Black to RGBA",
			hsla:     formats.HSLA{H: 0, S: 0, L: 0, A: 1},
			description: "Black HSLA should convert to RGBA(0,0,0,255)",
		},
		{
			name:     "Red to RGBA",
			hsla:     formats.HSLA{H: 0, S: 1, L: 0.5, A: 1},
			description: "Red HSLA should convert to RGBA(255,0,0,255)",
		},
		{
			name:     "Semi-transparent to RGBA",
			hsla:     formats.HSLA{H: 120, S: 1, L: 0.5, A: 0.5},
			description: "Semi-transparent green should have alpha=128",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.hsla.ToRGBA()

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("HSLA input: H=%.1f°, S=%.2f, L=%.2f, A=%.2f", tc.hsla.H, tc.hsla.S, tc.hsla.L, tc.hsla.A)
			t.Logf("RGBA result: R=%d, G=%d, B=%d, A=%d", result.R, result.G, result.B, result.A)
			t.Logf("Description: %s", tc.description)

			// Verify values are in valid range
			if result.R > 255 || result.G > 255 || result.B > 255 || result.A > 255 {
				t.Errorf("RGBA components out of range: %+v", result)
			}

			t.Logf("✅ HSLA ToRGBA method working correctly")
		})
	}
}

// TestHSLA_WithAlpha tests the WithAlpha method
func TestHSLA_WithAlpha(t *testing.T) {
	testCases := []struct {
		name     string
		original formats.HSLA
		newAlpha float64
		expected formats.HSLA
		description string
	}{
		{
			name:     "Set full alpha",
			original: formats.HSLA{H: 120, S: 0.8, L: 0.6, A: 0.5},
			newAlpha: 1.0,
			expected: formats.HSLA{H: 120, S: 0.8, L: 0.6, A: 1.0},
			description: "WithAlpha should only change alpha component",
		},
		{
			name:     "Set half alpha",
			original: formats.HSLA{H: 240, S: 1.0, L: 0.5, A: 1.0},
			newAlpha: 0.5,
			expected: formats.HSLA{H: 240, S: 1.0, L: 0.5, A: 0.5},
			description: "WithAlpha should preserve HSL components",
		},
		{
			name:     "Alpha clamping high",
			original: formats.HSLA{H: 60, S: 0.7, L: 0.3, A: 0.2},
			newAlpha: 1.5,
			expected: formats.HSLA{H: 60, S: 0.7, L: 0.3, A: 1.0},
			description: "WithAlpha should clamp high alpha to 1.0",
		},
		{
			name:     "Alpha clamping low",
			original: formats.HSLA{H: 300, S: 0.9, L: 0.7, A: 0.8},
			newAlpha: -0.2,
			expected: formats.HSLA{H: 300, S: 0.9, L: 0.7, A: 0.0},
			description: "WithAlpha should clamp low alpha to 0.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.original.WithAlpha(tc.newAlpha)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Original HSLA: H=%.1f°, S=%.2f, L=%.2f, A=%.2f", tc.original.H, tc.original.S, tc.original.L, tc.original.A)
			t.Logf("New alpha: %.2f", tc.newAlpha)
			t.Logf("Result HSLA: H=%.1f°, S=%.2f, L=%.2f, A=%.2f", result.H, result.S, result.L, result.A)
			t.Logf("Expected HSLA: H=%.1f°, S=%.2f, L=%.2f, A=%.2f", tc.expected.H, tc.expected.S, tc.expected.L, tc.expected.A)
			t.Logf("Description: %s", tc.description)

			// Verify the result
			tolerance := 0.0001
			if math.Abs(result.H-tc.expected.H) > tolerance ||
			   math.Abs(result.S-tc.expected.S) > tolerance ||
			   math.Abs(result.L-tc.expected.L) > tolerance ||
			   math.Abs(result.A-tc.expected.A) > tolerance {
				t.Errorf("WithAlpha result mismatch:\nExpected: %+v\nGot: %+v", tc.expected, result)
			}

			t.Logf("✅ HSLA WithAlpha method working correctly")
		})
	}
}

// TestWithAlpha tests the standalone WithAlpha function
func TestWithAlpha(t *testing.T) {
	testCases := []struct {
		name     string
		rgba     color.RGBA
		newAlpha float64
		expected color.RGBA
		description string
	}{
		{
			name:     "Set full alpha",
			rgba:     color.RGBA{R: 255, G: 128, B: 64, A: 128},
			newAlpha: 1.0,
			expected: color.RGBA{R: 255, G: 128, B: 64, A: 255},
			description: "WithAlpha should only change alpha component",
		},
		{
			name:     "Set half alpha",
			rgba:     color.RGBA{R: 100, G: 150, B: 200, A: 255},
			newAlpha: 0.5,
			expected: color.RGBA{R: 100, G: 150, B: 200, A: 128},
			description: "WithAlpha should preserve RGB components",
		},
		{
			name:     "Alpha clamping high",
			rgba:     color.RGBA{R: 50, G: 75, B: 25, A: 100},
			newAlpha: 1.5,
			expected: color.RGBA{R: 50, G: 75, B: 25, A: 255},
			description: "WithAlpha should clamp high alpha to 255",
		},
		{
			name:     "Alpha clamping low",
			rgba:     color.RGBA{R: 200, G: 100, B: 50, A: 200},
			newAlpha: -0.1,
			expected: color.RGBA{R: 200, G: 100, B: 50, A: 0},
			description: "WithAlpha should clamp low alpha to 0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.WithAlpha(tc.rgba, tc.newAlpha)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Original RGBA: R=%d, G=%d, B=%d, A=%d", tc.rgba.R, tc.rgba.G, tc.rgba.B, tc.rgba.A)
			t.Logf("New alpha: %.2f", tc.newAlpha)
			t.Logf("Result RGBA: R=%d, G=%d, B=%d, A=%d", result.R, result.G, result.B, result.A)
			t.Logf("Expected RGBA: R=%d, G=%d, B=%d, A=%d", tc.expected.R, tc.expected.G, tc.expected.B, tc.expected.A)
			t.Logf("Description: %s", tc.description)

			// Verify the result
			if result != tc.expected {
				t.Errorf("WithAlpha result mismatch:\nExpected: %+v\nGot: %+v", tc.expected, result)
			}

			t.Logf("✅ WithAlpha function working correctly")
		})
	}
}

// TestGetAlpha tests the GetAlpha function
func TestGetAlpha(t *testing.T) {
	testCases := []struct {
		name     string
		rgba     color.RGBA
		expected float64
		description string
	}{
		{
			name:     "Full alpha",
			rgba:     color.RGBA{R: 255, G: 128, B: 64, A: 255},
			expected: 1.0,
			description: "Alpha=255 should return 1.0",
		},
		{
			name:     "No alpha",
			rgba:     color.RGBA{R: 100, G: 150, B: 200, A: 0},
			expected: 0.0,
			description: "Alpha=0 should return 0.0",
		},
		{
			name:     "Half alpha",
			rgba:     color.RGBA{R: 50, G: 75, B: 25, A: 128},
			expected: 128.0/255.0,
			description: "Alpha=128 should return ~0.502",
		},
		{
			name:     "Quarter alpha",
			rgba:     color.RGBA{R: 200, G: 100, B: 50, A: 64},
			expected: 64.0/255.0,
			description: "Alpha=64 should return ~0.251",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.GetAlpha(tc.rgba)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("RGBA input: R=%d, G=%d, B=%d, A=%d", tc.rgba.R, tc.rgba.G, tc.rgba.B, tc.rgba.A)
			t.Logf("Alpha result: %.6f", result)
			t.Logf("Expected alpha: %.6f", tc.expected)
			t.Logf("Description: %s", tc.description)

			// Verify the result with small tolerance for floating point
			tolerance := 0.001
			if math.Abs(result-tc.expected) > tolerance {
				t.Errorf("GetAlpha result mismatch: expected %.6f, got %.6f", tc.expected, result)
			}

			// Verify result is in valid range
			if result < 0 || result > 1 {
				t.Errorf("GetAlpha result out of range [0,1]: %.6f", result)
			}

			t.Logf("✅ GetAlpha function working correctly")
		})
	}
}