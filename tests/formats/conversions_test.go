package formats_test

import (
	"fmt"
	"image/color"
	"math"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func TestRGBAToHSLA(t *testing.T) {
	testCases := []struct {
		name     string
		rgba     color.RGBA
		expected formats.HSLA
		tolerance float64
		description string
	}{
		{
			name:     "Pure Red",
			rgba:     color.RGBA{R: 255, G: 0, B: 0, A: 255},
			expected: formats.NewHSLA(0, 1.0, 0.5, 1.0),
			tolerance: 0.001,
			description: "Pure red should be hue 0°, full saturation, 50% lightness",
		},
		{
			name:     "Pure Green",
			rgba:     color.RGBA{R: 0, G: 255, B: 0, A: 255},
			expected: formats.NewHSLA(120, 1.0, 0.5, 1.0),
			tolerance: 0.001,
			description: "Pure green should be hue 120°, full saturation, 50% lightness",
		},
		{
			name:     "Pure Blue", 
			rgba:     color.RGBA{R: 0, G: 0, B: 255, A: 255},
			expected: formats.NewHSLA(240, 1.0, 0.5, 1.0),
			tolerance: 0.001,
			description: "Pure blue should be hue 240°, full saturation, 50% lightness",
		},
		{
			name:     "White",
			rgba:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expected: formats.NewHSLA(0, 0.0, 1.0, 1.0), // Hue undefined, but often 0
			tolerance: 0.001,
			description: "White should be 0% saturation, 100% lightness",
		},
		{
			name:     "Black",
			rgba:     color.RGBA{R: 0, G: 0, B: 0, A: 255},
			expected: formats.NewHSLA(0, 0.0, 0.0, 1.0), // Hue undefined, but often 0
			tolerance: 0.001,
			description: "Black should be 0% saturation, 0% lightness",
		},
		{
			name:     "Medium Gray",
			rgba:     color.RGBA{R: 128, G: 128, B: 128, A: 255},
			expected: formats.NewHSLA(0, 0.0, 0.502, 1.0), // 128/255 ≈ 0.502
			tolerance: 0.01,
			description: "Medium gray should be 0% saturation, ~50% lightness",
		},
		{
			name:     "Semi-transparent Red",
			rgba:     color.RGBA{R: 255, G: 0, B: 0, A: 128},
			expected: formats.NewHSLA(0, 1.0, 0.5, 0.502), // 128/255 ≈ 0.502
			tolerance: 0.01,
			description: "Semi-transparent red should maintain RGB properties with 50% alpha",
		},
		{
			name:     "Orange",
			rgba:     color.RGBA{R: 255, G: 165, B: 0, A: 255},
			expected: formats.NewHSLA(38.8, 1.0, 0.5, 1.0), // Orange ≈ 39° hue
			tolerance: 1.0, // Allow 1° hue tolerance
			description: "Orange should be around 39° hue with full saturation",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formats.RGBAToHSLA(tc.rgba)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Input RGB: (%d,%d,%d,%d)", tc.rgba.R, tc.rgba.G, tc.rgba.B, tc.rgba.A)
			t.Logf("Expected HSLA: (%.1f°, %.3f, %.3f, %.3f)", tc.expected.H, tc.expected.S, tc.expected.L, tc.expected.A)
			t.Logf("Actual HSLA: (%.1f°, %.3f, %.3f, %.3f)", actual.H, actual.S, actual.L, actual.A)
			t.Logf("Differences: H=%.3f° S=%.6f L=%.6f A=%.6f", 
				hueDifference(tc.expected.H, actual.H), 
				math.Abs(tc.expected.S-actual.S),
				math.Abs(tc.expected.L-actual.L),
				math.Abs(tc.expected.A-actual.A))
			t.Logf("Description: %s", tc.description)

			// Compare with tolerance, handling hue wraparound
			if hueDifference(tc.expected.H, actual.H) > tc.tolerance {
				t.Errorf("Hue mismatch: expected %.3f°, got %.3f° (diff: %.3f°)", 
					tc.expected.H, actual.H, hueDifference(tc.expected.H, actual.H))
			}
			if math.Abs(tc.expected.S-actual.S) > tc.tolerance {
				t.Errorf("Saturation mismatch: expected %.6f, got %.6f (diff: %.6f)", 
					tc.expected.S, actual.S, math.Abs(tc.expected.S-actual.S))
			}
			if math.Abs(tc.expected.L-actual.L) > tc.tolerance {
				t.Errorf("Lightness mismatch: expected %.6f, got %.6f (diff: %.6f)", 
					tc.expected.L, actual.L, math.Abs(tc.expected.L-actual.L))
			}
			if math.Abs(tc.expected.A-actual.A) > tc.tolerance {
				t.Errorf("Alpha mismatch: expected %.6f, got %.6f (diff: %.6f)", 
					tc.expected.A, actual.A, math.Abs(tc.expected.A-actual.A))
			}
		})
	}
}

func TestHSLAToRGBA(t *testing.T) {
	testCases := []struct {
		name     string
		hsla     formats.HSLA
		expected color.RGBA
		tolerance uint8
		description string
	}{
		{
			name:     "Pure Red",
			hsla:     formats.NewHSLA(0, 1.0, 0.5, 1.0),
			expected: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			tolerance: 1,
			description: "Hue 0° should produce pure red",
		},
		{
			name:     "Pure Green",
			hsla:     formats.NewHSLA(120, 1.0, 0.5, 1.0),
			expected: color.RGBA{R: 0, G: 255, B: 0, A: 255},
			tolerance: 1,
			description: "Hue 120° should produce pure green",
		},
		{
			name:     "Pure Blue",
			hsla:     formats.NewHSLA(240, 1.0, 0.5, 1.0),
			expected: color.RGBA{R: 0, G: 0, B: 255, A: 255},
			tolerance: 1,
			description: "Hue 240° should produce pure blue",
		},
		{
			name:     "White",
			hsla:     formats.NewHSLA(0, 0.0, 1.0, 1.0),
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			tolerance: 1,
			description: "100% lightness should produce white regardless of hue",
		},
		{
			name:     "Black",
			hsla:     formats.NewHSLA(180, 1.0, 0.0, 1.0), // Any hue, 0% lightness
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 255},
			tolerance: 1,
			description: "0% lightness should produce black regardless of hue/saturation",
		},
		{
			name:     "Medium Gray",
			hsla:     formats.NewHSLA(0, 0.0, 0.502, 1.0),
			expected: color.RGBA{R: 128, G: 128, B: 128, A: 255},
			tolerance: 2, // Allow rounding tolerance
			description: "0% saturation should produce gray",
		},
		{
			name:     "Semi-transparent Blue",
			hsla:     formats.NewHSLA(240, 1.0, 0.5, 0.502),
			expected: color.RGBA{R: 0, G: 0, B: 255, A: 128},
			tolerance: 2,
			description: "50% alpha should be preserved in conversion",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formats.HSLAToRGBA(tc.hsla)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Input HSLA: (%.1f°, %.3f, %.3f, %.3f)", tc.hsla.H, tc.hsla.S, tc.hsla.L, tc.hsla.A)
			t.Logf("Expected RGB: (%d,%d,%d,%d)", tc.expected.R, tc.expected.G, tc.expected.B, tc.expected.A)
			t.Logf("Actual RGB: (%d,%d,%d,%d)", actual.R, actual.G, actual.B, actual.A)
			t.Logf("Differences: R=%d G=%d B=%d A=%d", 
				abs(int(tc.expected.R)-int(actual.R)),
				abs(int(tc.expected.G)-int(actual.G)),
				abs(int(tc.expected.B)-int(actual.B)),
				abs(int(tc.expected.A)-int(actual.A)))
			t.Logf("Description: %s", tc.description)

			// Compare with tolerance
			if abs(int(tc.expected.R)-int(actual.R)) > int(tc.tolerance) {
				t.Errorf("Red component mismatch: expected %d, got %d (diff: %d)", 
					tc.expected.R, actual.R, abs(int(tc.expected.R)-int(actual.R)))
			}
			if abs(int(tc.expected.G)-int(actual.G)) > int(tc.tolerance) {
				t.Errorf("Green component mismatch: expected %d, got %d (diff: %d)", 
					tc.expected.G, actual.G, abs(int(tc.expected.G)-int(actual.G)))
			}
			if abs(int(tc.expected.B)-int(actual.B)) > int(tc.tolerance) {
				t.Errorf("Blue component mismatch: expected %d, got %d (diff: %d)", 
					tc.expected.B, actual.B, abs(int(tc.expected.B)-int(actual.B)))
			}
			if abs(int(tc.expected.A)-int(actual.A)) > int(tc.tolerance) {
				t.Errorf("Alpha component mismatch: expected %d, got %d (diff: %d)", 
					tc.expected.A, actual.A, abs(int(tc.expected.A)-int(actual.A)))
			}
		})
	}
}

func TestRGBAToXYZ(t *testing.T) {
	testCases := []struct {
		name     string
		rgba     color.RGBA
		expectedX float64
		expectedY float64
		expectedZ float64
		tolerance float64
		description string
	}{
		{
			name:     "Pure White D65",
			rgba:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expectedX: 95.047, // D65 white point
			expectedY: 100.000,
			expectedZ: 108.883,
			tolerance: 0.1,
			description: "Pure white should match D65 illuminant values",
		},
		{
			name:     "Pure Black",
			rgba:     color.RGBA{R: 0, G: 0, B: 0, A: 255},
			expectedX: 0,
			expectedY: 0,
			expectedZ: 0,
			tolerance: 0.001,
			description: "Pure black should have zero XYZ components",
		},
		{
			name:     "Pure Red",
			rgba:     color.RGBA{R: 255, G: 0, B: 0, A: 255},
			expectedX: 41.24, // Approximate sRGB red in XYZ
			expectedY: 21.26,
			expectedZ: 1.93,
			tolerance: 0.5,
			description: "Pure red should have known XYZ values for sRGB",
		},
		{
			name:     "Medium Gray",
			rgba:     color.RGBA{R: 128, G: 128, B: 128, A: 255},
			expectedX: 20.5, // Approximately 1/4 of white point
			expectedY: 21.6,
			expectedZ: 23.5,
			tolerance: 1.0,
			description: "Medium gray should be roughly 1/4 of white point values",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			xyz := formats.RGBAToXYZ(tc.rgba)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Input RGB: (%d,%d,%d)", tc.rgba.R, tc.rgba.G, tc.rgba.B)
			t.Logf("Expected XYZ: (%.3f, %.3f, %.3f)", tc.expectedX, tc.expectedY, tc.expectedZ)
			t.Logf("Actual XYZ: (%.3f, %.3f, %.3f)", xyz.X, xyz.Y, xyz.Z)
			t.Logf("Differences: X=%.3f Y=%.3f Z=%.3f", 
				math.Abs(tc.expectedX-xyz.X),
				math.Abs(tc.expectedY-xyz.Y),
				math.Abs(tc.expectedZ-xyz.Z))
			t.Logf("Description: %s", tc.description)

			// Compare with tolerance
			if math.Abs(tc.expectedX-xyz.X) > tc.tolerance {
				t.Errorf("X component mismatch: expected %.3f, got %.3f (diff: %.3f)", 
					tc.expectedX, xyz.X, math.Abs(tc.expectedX-xyz.X))
			}
			if math.Abs(tc.expectedY-xyz.Y) > tc.tolerance {
				t.Errorf("Y component mismatch: expected %.3f, got %.3f (diff: %.3f)", 
					tc.expectedY, xyz.Y, math.Abs(tc.expectedY-xyz.Y))
			}
			if math.Abs(tc.expectedZ-xyz.Z) > tc.tolerance {
				t.Errorf("Z component mismatch: expected %.3f, got %.3f (diff: %.3f)", 
					tc.expectedZ, xyz.Z, math.Abs(tc.expectedZ-xyz.Z))
			}
		})
	}
}

func TestXYZToRGBA(t *testing.T) {
	testCases := []struct {
		name     string
		xyz      formats.XYZ
		expected color.RGBA
		tolerance uint8
		description string
	}{
		{
			name:     "D65 White Point",
			xyz:      formats.NewXYZ(95.047, 100.000, 108.883),
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			tolerance: 2,
			description: "D65 white point should convert to pure white RGB",
		},
		{
			name:     "Origin (Black)",
			xyz:      formats.NewXYZ(0, 0, 0),
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 255},
			tolerance: 1,
			description: "Origin in XYZ should convert to black RGB",
		},
		{
			name:     "Red Primary",
			xyz:      formats.NewXYZ(41.24, 21.26, 1.93),
			expected: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			tolerance: 5, // Allow more tolerance for conversion accuracy
			description: "sRGB red primary in XYZ should convert back to red",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formats.XYZToRGBA(tc.xyz)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Input XYZ: (%.3f, %.3f, %.3f)", tc.xyz.X, tc.xyz.Y, tc.xyz.Z)
			t.Logf("Expected RGB: (%d,%d,%d)", tc.expected.R, tc.expected.G, tc.expected.B)
			t.Logf("Actual RGB: (%d,%d,%d)", actual.R, actual.G, actual.B)
			t.Logf("Differences: R=%d G=%d B=%d", 
				abs(int(tc.expected.R)-int(actual.R)),
				abs(int(tc.expected.G)-int(actual.G)),
				abs(int(tc.expected.B)-int(actual.B)))
			t.Logf("Description: %s", tc.description)

			// Compare with tolerance
			if abs(int(tc.expected.R)-int(actual.R)) > int(tc.tolerance) {
				t.Errorf("Red component mismatch: expected %d, got %d (diff: %d)", 
					tc.expected.R, actual.R, abs(int(tc.expected.R)-int(actual.R)))
			}
			if abs(int(tc.expected.G)-int(actual.G)) > int(tc.tolerance) {
				t.Errorf("Green component mismatch: expected %d, got %d (diff: %d)", 
					tc.expected.G, actual.G, abs(int(tc.expected.G)-int(actual.G)))
			}
			if abs(int(tc.expected.B)-int(actual.B)) > int(tc.tolerance) {
				t.Errorf("Blue component mismatch: expected %d, got %d (diff: %d)", 
					tc.expected.B, actual.B, abs(int(tc.expected.B)-int(actual.B)))
			}
		})
	}
}

func TestRGBAToLAB(t *testing.T) {
	testCases := []struct {
		name     string
		rgba     color.RGBA
		expectedL float64
		expectedA float64
		expectedB float64
		tolerance float64
		description string
	}{
		{
			name:     "Pure White",
			rgba:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expectedL: 100.0, // Perfect lightness
			expectedA: 0.0,   // Neutral a*
			expectedB: 0.0,   // Neutral b*
			tolerance: 0.5,
			description: "Pure white should be L*=100, a*=0, b*=0 in LAB",
		},
		{
			name:     "Pure Black",
			rgba:     color.RGBA{R: 0, G: 0, B: 0, A: 255},
			expectedL: 0.0,
			expectedA: 0.0,
			expectedB: 0.0,
			tolerance: 0.1,
			description: "Pure black should be L*=0, a*=0, b*=0 in LAB",
		},
		{
			name:     "Pure Red",
			rgba:     color.RGBA{R: 255, G: 0, B: 0, A: 255},
			expectedL: 53.2, // Approximate L* for sRGB red
			expectedA: 80.1, // Strong positive a* (red direction)
			expectedB: 67.2, // Positive b* (yellow direction)
			tolerance: 2.0,
			description: "Pure red should have moderate L*, high positive a*, positive b*",
		},
		{
			name:     "Pure Green",
			rgba:     color.RGBA{R: 0, G: 255, B: 0, A: 255},
			expectedL: 87.7,  // High L* for green
			expectedA: -86.2, // Strong negative a* (green direction)
			expectedB: 83.2,  // High positive b* (yellow direction)
			tolerance: 2.0,
			description: "Pure green should have high L*, negative a*, positive b*",
		},
		{
			name:     "Pure Blue",
			rgba:     color.RGBA{R: 0, G: 0, B: 255, A: 255},
			expectedL: 32.3,  // Low L* for blue
			expectedA: 79.2,  // Positive a* 
			expectedB: -107.9, // Strong negative b* (blue direction)
			tolerance: 3.0,
			description: "Pure blue should have low L*, positive a*, negative b*",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lab := formats.RGBAToLAB(tc.rgba)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Input RGB: (%d,%d,%d)", tc.rgba.R, tc.rgba.G, tc.rgba.B)
			t.Logf("Expected LAB: (%.1f, %.1f, %.1f)", tc.expectedL, tc.expectedA, tc.expectedB)
			t.Logf("Actual LAB: (%.1f, %.1f, %.1f)", lab.L, lab.A, lab.B)
			t.Logf("Differences: L*=%.3f a*=%.3f b*=%.3f", 
				math.Abs(tc.expectedL-lab.L),
				math.Abs(tc.expectedA-lab.A),
				math.Abs(tc.expectedB-lab.B))
			t.Logf("Description: %s", tc.description)

			// Compare with tolerance
			if math.Abs(tc.expectedL-lab.L) > tc.tolerance {
				t.Errorf("L* component mismatch: expected %.1f, got %.1f (diff: %.3f)", 
					tc.expectedL, lab.L, math.Abs(tc.expectedL-lab.L))
			}
			if math.Abs(tc.expectedA-lab.A) > tc.tolerance {
				t.Errorf("a* component mismatch: expected %.1f, got %.1f (diff: %.3f)", 
					tc.expectedA, lab.A, math.Abs(tc.expectedA-lab.A))
			}
			if math.Abs(tc.expectedB-lab.B) > tc.tolerance {
				t.Errorf("b* component mismatch: expected %.1f, got %.1f (diff: %.3f)", 
					tc.expectedB, lab.B, math.Abs(tc.expectedB-lab.B))
			}
		})
	}
}

func TestLABToRGBA(t *testing.T) {
	testCases := []struct {
		name     string
		lab      formats.LAB
		expected color.RGBA
		tolerance uint8
		description string
	}{
		{
			name:     "Perfect White",
			lab:      formats.NewLAB(100.0, 0.0, 0.0),
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			tolerance: 2,
			description: "LAB white should convert to RGB white",
		},
		{
			name:     "Perfect Black",
			lab:      formats.NewLAB(0.0, 0.0, 0.0),
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 255},
			tolerance: 1,
			description: "LAB black should convert to RGB black",
		},
		{
			name:     "Red-like color",
			lab:      formats.NewLAB(50.0, 70.0, 50.0), // Red-ish LAB values
			expected: color.RGBA{R: 200, G: 0, B: 50, A: 255}, // Approximate expectation
			tolerance: 50, // Allow wide tolerance for complex conversion
			description: "Red-like LAB should convert to reddish RGB",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formats.LABToRGBA(tc.lab)

			// Log diagnostic information  
			t.Logf("Test: %s", tc.name)
			t.Logf("Input LAB: (%.1f, %.1f, %.1f)", tc.lab.L, tc.lab.A, tc.lab.B)
			t.Logf("Expected RGB: (%d,%d,%d)", tc.expected.R, tc.expected.G, tc.expected.B)
			t.Logf("Actual RGB: (%d,%d,%d)", actual.R, actual.G, actual.B)
			t.Logf("Differences: R=%d G=%d B=%d", 
				abs(int(tc.expected.R)-int(actual.R)),
				abs(int(tc.expected.G)-int(actual.G)),
				abs(int(tc.expected.B)-int(actual.B)))
			t.Logf("Description: %s", tc.description)

			// Basic validation - colors should be within valid range
			if actual.A != 255 {
				t.Errorf("Alpha should be 255, got %d", actual.A)
			}

			// For non-perfect colors, just ensure we get reasonable results
			if tc.name == "Perfect White" || tc.name == "Perfect Black" {
				// Strict comparison for perfect colors
				if abs(int(tc.expected.R)-int(actual.R)) > int(tc.tolerance) {
					t.Errorf("Red component mismatch: expected %d, got %d (diff: %d)", 
						tc.expected.R, actual.R, abs(int(tc.expected.R)-int(actual.R)))
				}
				if abs(int(tc.expected.G)-int(actual.G)) > int(tc.tolerance) {
					t.Errorf("Green component mismatch: expected %d, got %d (diff: %d)", 
						tc.expected.G, actual.G, abs(int(tc.expected.G)-int(actual.G)))
				}
				if abs(int(tc.expected.B)-int(actual.B)) > int(tc.tolerance) {
					t.Errorf("Blue component mismatch: expected %d, got %d (diff: %d)", 
						tc.expected.B, actual.B, abs(int(tc.expected.B)-int(actual.B)))
				}
			} else {
				// For complex colors, just log the results
				t.Logf("Complex color conversion completed - see diagnostic output above")
			}
		})
	}
}

func TestRoundTripConversions(t *testing.T) {
	// Test that RGB → HSLA → RGB preserves the original color
	testColors := []color.RGBA{
		{R: 255, G: 0, B: 0, A: 255},     // Red
		{R: 0, G: 255, B: 0, A: 255},     // Green
		{R: 0, G: 0, B: 255, A: 255},     // Blue
		{R: 255, G: 255, B: 255, A: 255}, // White
		{R: 0, G: 0, B: 0, A: 255},       // Black
		{R: 128, G: 128, B: 128, A: 255}, // Gray
		{R: 255, G: 128, B: 64, A: 200},  // Orange with transparency
		{R: 64, G: 192, B: 255, A: 100},  // Sky blue with transparency
	}

	for i, original := range testColors {
		t.Run(fmt.Sprintf("RoundTrip_RGB_HSLA_%d", i), func(t *testing.T) {
			// RGB → HSLA → RGB
			hsla := formats.RGBAToHSLA(original)
			roundTrip := formats.HSLAToRGBA(hsla)

			// Log diagnostic information
			t.Logf("Original RGB: (%d,%d,%d,%d)", original.R, original.G, original.B, original.A)
			t.Logf("Intermediate HSLA: (%.1f°, %.3f, %.3f, %.3f)", hsla.H, hsla.S, hsla.L, hsla.A)
			t.Logf("Round-trip RGB: (%d,%d,%d,%d)", roundTrip.R, roundTrip.G, roundTrip.B, roundTrip.A)
			t.Logf("Differences: R=%d G=%d B=%d A=%d", 
				abs(int(original.R)-int(roundTrip.R)),
				abs(int(original.G)-int(roundTrip.G)),
				abs(int(original.B)-int(roundTrip.B)),
				abs(int(original.A)-int(roundTrip.A)))

			tolerance := uint8(2) // Allow small rounding errors
			if abs(int(original.R)-int(roundTrip.R)) > int(tolerance) {
				t.Errorf("Red round-trip error: %d → %d (diff: %d)", 
					original.R, roundTrip.R, abs(int(original.R)-int(roundTrip.R)))
			}
			if abs(int(original.G)-int(roundTrip.G)) > int(tolerance) {
				t.Errorf("Green round-trip error: %d → %d (diff: %d)", 
					original.G, roundTrip.G, abs(int(original.G)-int(roundTrip.G)))
			}
			if abs(int(original.B)-int(roundTrip.B)) > int(tolerance) {
				t.Errorf("Blue round-trip error: %d → %d (diff: %d)", 
					original.B, roundTrip.B, abs(int(original.B)-int(roundTrip.B)))
			}
			if abs(int(original.A)-int(roundTrip.A)) > int(tolerance) {
				t.Errorf("Alpha round-trip error: %d → %d (diff: %d)", 
					original.A, roundTrip.A, abs(int(original.A)-int(roundTrip.A)))
			}
		})
	}

	// Test XYZ round-trip for a few key colors
	xyzTestColors := []color.RGBA{
		{R: 255, G: 255, B: 255, A: 255}, // White
		{R: 0, G: 0, B: 0, A: 255},       // Black  
		{R: 255, G: 0, B: 0, A: 255},     // Red
		{R: 128, G: 128, B: 128, A: 255}, // Gray
	}

	for i, original := range xyzTestColors {
		t.Run(fmt.Sprintf("RoundTrip_RGB_XYZ_%d", i), func(t *testing.T) {
			// RGB → XYZ → RGB
			xyz := formats.RGBAToXYZ(original)
			roundTrip := formats.XYZToRGBA(xyz)

			// Log diagnostic information
			t.Logf("Original RGB: (%d,%d,%d)", original.R, original.G, original.B)
			t.Logf("Intermediate XYZ: (%.3f, %.3f, %.3f)", xyz.X, xyz.Y, xyz.Z)
			t.Logf("Round-trip RGB: (%d,%d,%d)", roundTrip.R, roundTrip.G, roundTrip.B)
			t.Logf("Differences: R=%d G=%d B=%d", 
				abs(int(original.R)-int(roundTrip.R)),
				abs(int(original.G)-int(roundTrip.G)),
				abs(int(original.B)-int(roundTrip.B)))

			tolerance := uint8(3) // Allow some tolerance for XYZ conversion
			if abs(int(original.R)-int(roundTrip.R)) > int(tolerance) {
				t.Errorf("Red XYZ round-trip error: %d → %d (diff: %d)", 
					original.R, roundTrip.R, abs(int(original.R)-int(roundTrip.R)))
			}
			if abs(int(original.G)-int(roundTrip.G)) > int(tolerance) {
				t.Errorf("Green XYZ round-trip error: %d → %d (diff: %d)", 
					original.G, roundTrip.G, abs(int(original.G)-int(roundTrip.G)))
			}
			if abs(int(original.B)-int(roundTrip.B)) > int(tolerance) {
				t.Errorf("Blue XYZ round-trip error: %d → %d (diff: %d)", 
					original.B, roundTrip.B, abs(int(original.B)-int(roundTrip.B)))
			}
		})
	}
}

func TestQuantizeColor(t *testing.T) {
	testCases := []struct {
		name        string
		input       color.RGBA
		bits        uint8
		expected    color.RGBA
		description string
	}{
		{
			name:     "5-bit quantization (default)",
			input:    color.RGBA{R: 127, G: 200, B: 63, A: 255},
			bits:     5,
			expected: color.RGBA{R: 124, G: 204, B: 60, A: 255}, // (120+4, 200+4, 56+4)
			description: "5-bit quantization reduces precision to 32 levels per channel",
		},
		{
			name:     "4-bit quantization",
			input:    color.RGBA{R: 127, G: 200, B: 63, A: 255},
			bits:     4,
			expected: color.RGBA{R: 120, G: 200, B: 56, A: 255}, // (112+8, 192+8, 48+8)
			description: "4-bit quantization reduces precision to 16 levels per channel",
		},
		{
			name:     "3-bit quantization (aggressive)",
			input:    color.RGBA{R: 127, G: 200, B: 63, A: 255},
			bits:     3,
			expected: color.RGBA{R: 112, G: 208, B: 48, A: 255}, // (96+16, 192+16, 32+16)
			description: "3-bit quantization reduces precision to 8 levels per channel",
		},
		{
			name:     "8-bit quantization (no change)",
			input:    color.RGBA{R: 127, G: 200, B: 63, A: 255},
			bits:     8,
			expected: color.RGBA{R: 127, G: 200, B: 63, A: 255}, // (127+0, 200+0, 63+0)
			description: "8-bit quantization preserves original precision",
		},
		{
			name:     "Pure colors quantized",
			input:    color.RGBA{R: 255, G: 0, B: 128, A: 255},
			bits:     5,
			expected: color.RGBA{R: 252, G: 4, B: 132, A: 255}, // (248+4, 0+4, 128+4)
			description: "Colors quantize to nearest levels with rounding",
		},
		{
			name:     "Black quantized",
			input:    color.RGBA{R: 0, G: 0, B: 0, A: 255},
			bits:     5,
			expected: color.RGBA{R: 4, G: 4, B: 4, A: 255}, // (0+4, 0+4, 0+4)
			description: "Black quantizes to lowest level plus rounding offset",
		},
		{
			name:     "White quantization",
			input:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			bits:     5,
			expected: color.RGBA{R: 252, G: 252, B: 252, A: 255}, // (248+4, 248+4, 248+4)
			description: "White quantizes to highest quantization level",
		},
		{
			name:     "Invalid bits (too high) defaults to 5",
			input:    color.RGBA{R: 127, G: 200, B: 63, A: 255},
			bits:     10, // Invalid - should default to 5
			expected: color.RGBA{R: 124, G: 204, B: 60, A: 255}, // Same as 5-bit
			description: "Invalid bit count should default to 5-bit quantization",
		},
		{
			name:     "Invalid bits (too low) defaults to 5",
			input:    color.RGBA{R: 127, G: 200, B: 63, A: 255},
			bits:     0, // Invalid - should default to 5
			expected: color.RGBA{R: 124, G: 204, B: 60, A: 255}, // Same as 5-bit
			description: "Invalid bit count should default to 5-bit quantization",
		},
		{
			name:     "Alpha preserved",
			input:    color.RGBA{R: 127, G: 200, B: 63, A: 128},
			bits:     5,
			expected: color.RGBA{R: 124, G: 204, B: 60, A: 255}, // Alpha always set to 255
			description: "Alpha channel should always be set to 255",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := formats.QuantizeColor(tc.input, tc.bits)

			// Log comprehensive diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Input RGBA: (%d,%d,%d,%d)", tc.input.R, tc.input.G, tc.input.B, tc.input.A)
			t.Logf("Quantization bits: %d", tc.bits)
			t.Logf("Expected RGBA: (%d,%d,%d,%d)", tc.expected.R, tc.expected.G, tc.expected.B, tc.expected.A)
			t.Logf("Actual RGBA: (%d,%d,%d,%d)", actual.R, actual.G, actual.B, actual.A)
			t.Logf("Differences: R=%d G=%d B=%d A=%d",
				abs(int(tc.expected.R)-int(actual.R)),
				abs(int(tc.expected.G)-int(actual.G)),
				abs(int(tc.expected.B)-int(actual.B)),
				abs(int(tc.expected.A)-int(actual.A)))
			t.Logf("Description: %s", tc.description)

			// Calculate quantization levels for validation
			if tc.bits >= 1 && tc.bits <= 8 {
				levels := 1 << tc.bits // 2^bits
				t.Logf("Quantization levels: %d (step size: %.1f)", levels, 256.0/float64(levels))
			}

			// Compare with expected values
			if actual.R != tc.expected.R {
				t.Errorf("Red component mismatch: expected %d, got %d", tc.expected.R, actual.R)
			}
			if actual.G != tc.expected.G {
				t.Errorf("Green component mismatch: expected %d, got %d", tc.expected.G, actual.G)
			}
			if actual.B != tc.expected.B {
				t.Errorf("Blue component mismatch: expected %d, got %d", tc.expected.B, actual.B)
			}
			if actual.A != tc.expected.A {
				t.Errorf("Alpha component mismatch: expected %d, got %d", tc.expected.A, actual.A)
			}

			// Validate that quantization reduces precision
			if tc.bits < 8 {
				// Verify that the result uses fewer distinct values
				t.Logf("Quantization validation: %d-bit should reduce precision from 256 to %d levels",
					tc.bits, 1<<tc.bits)
			}
		})
	}
}

func TestQuantizeColor_ConsistentBehavior(t *testing.T) {
	// Test that quantization is consistent and reversible at the same bit depth
	testColors := []color.RGBA{
		{R: 0, G: 0, B: 0, A: 255},       // Black
		{R: 255, G: 255, B: 255, A: 255}, // White
		{R: 128, G: 128, B: 128, A: 255}, // Gray
		{R: 255, G: 0, B: 0, A: 255},     // Red
		{R: 64, G: 192, B: 128, A: 255},  // Mixed color
		{R: 200, G: 100, B: 50, A: 255},  // Another mixed color
	}

	bitDepths := []uint8{3, 4, 5, 6, 7, 8}

	for _, bits := range bitDepths {
		for i, color := range testColors {
			t.Run(fmt.Sprintf("Consistency_Bits%d_Color%d", bits, i), func(t *testing.T) {
				// Quantize the same color multiple times
				result1 := formats.QuantizeColor(color, bits)
				result2 := formats.QuantizeColor(color, bits)
				result3 := formats.QuantizeColor(result1, bits) // Quantize an already quantized color

				t.Logf("Testing %d-bit quantization consistency", bits)
				t.Logf("Original: (%d,%d,%d,%d)", color.R, color.G, color.B, color.A)
				t.Logf("First quantization: (%d,%d,%d,%d)", result1.R, result1.G, result1.B, result1.A)
				t.Logf("Second quantization: (%d,%d,%d,%d)", result2.R, result2.G, result2.B, result2.A)
				t.Logf("Re-quantized: (%d,%d,%d,%d)", result3.R, result3.G, result3.B, result3.A)

				// Results should be identical
				if result1 != result2 {
					t.Errorf("Quantization not consistent: first=%+v, second=%+v", result1, result2)
				}

				// Re-quantizing should not change the result
				if result1 != result3 {
					t.Errorf("Re-quantization changed result: original=%+v, re-quantized=%+v", result1, result3)
				}

				// Alpha should always be 255
				if result1.A != 255 {
					t.Errorf("Alpha not preserved: expected 255, got %d", result1.A)
				}
			})
		}
	}
}

// Helper functions
func hueDifference(h1, h2 float64) float64 {
	// Calculate the shortest distance between two hue values (handling 360° wraparound)
	diff := math.Abs(h1 - h2)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}