package formats_test

import (
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

// TestNewLAB tests the NewLAB constructor function
func TestNewLAB(t *testing.T) {
	testCases := []struct {
		name string
		l, a, b float64
		expected formats.LAB
		description string
	}{
		{
			name: "Standard white point",
			l: 100.0, a: 0.0, b: 0.0,
			expected: formats.LAB{L: 100.0, A: 0.0, B: 0.0},
			description: "Constructor should create white point in LAB space",
		},
		{
			name: "Standard black point",
			l: 0.0, a: 0.0, b: 0.0,
			expected: formats.LAB{L: 0.0, A: 0.0, B: 0.0},
			description: "Constructor should create black point in LAB space",
		},
		{
			name: "Positive a and b values",
			l: 50.0, a: 25.5, b: -15.3,
			expected: formats.LAB{L: 50.0, A: 25.5, B: -15.3},
			description: "Constructor should handle positive and negative a/b values",
		},
		{
			name: "Maximum typical values",
			l: 100.0, a: 127.0, b: 127.0,
			expected: formats.LAB{L: 100.0, A: 127.0, B: 127.0},
			description: "Constructor should handle maximum typical LAB values",
		},
		{
			name: "Minimum typical values",
			l: 0.0, a: -128.0, b: -128.0,
			expected: formats.LAB{L: 0.0, A: -128.0, B: -128.0},
			description: "Constructor should handle minimum typical LAB values",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lab := formats.NewLAB(tc.l, tc.a, tc.b)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Input values: L=%.1f, A=%.1f, B=%.1f", tc.l, tc.a, tc.b)
			t.Logf("Created LAB: %+v", lab)
			t.Logf("Expected LAB: %+v", tc.expected)
			t.Logf("Description: %s", tc.description)

			// Verify the created LAB matches expected values
			if lab.L != tc.expected.L {
				t.Errorf("L component mismatch: expected %.1f, got %.1f", tc.expected.L, lab.L)
			}
			if lab.A != tc.expected.A {
				t.Errorf("A component mismatch: expected %.1f, got %.1f", tc.expected.A, lab.A)
			}
			if lab.B != tc.expected.B {
				t.Errorf("B component mismatch: expected %.1f, got %.1f", tc.expected.B, lab.B)
			}

			t.Logf("✅ LAB constructor working correctly")
		})
	}
}

// TestLAB_IsValid tests the IsValid method on LAB
func TestLAB_IsValid(t *testing.T) {
	testCases := []struct {
		name string
		lab formats.LAB
		expected bool
		description string
	}{
		{
			name: "Valid white point",
			lab: formats.LAB{L: 100.0, A: 0.0, B: 0.0},
			expected: true,
			description: "White point should be valid",
		},
		{
			name: "Valid black point",
			lab: formats.LAB{L: 0.0, A: 0.0, B: 0.0},
			expected: true,
			description: "Black point should be valid",
		},
		{
			name: "Valid typical color",
			lab: formats.LAB{L: 50.0, A: 25.0, B: -15.0},
			expected: true,
			description: "Typical color within valid ranges should be valid",
		},
		{
			name: "Invalid L too high",
			lab: formats.LAB{L: 150.0, A: 0.0, B: 0.0},
			expected: false,
			description: "L values above 100 should be invalid",
		},
		{
			name: "Invalid L negative",
			lab: formats.LAB{L: -10.0, A: 0.0, B: 0.0},
			expected: false,
			description: "Negative L values should be invalid",
		},
		{
			name: "Invalid A too high",
			lab: formats.LAB{L: 50.0, A: 200.0, B: 0.0},
			expected: false,
			description: "A values above 127 should be invalid",
		},
		{
			name: "Invalid A too low",
			lab: formats.LAB{L: 50.0, A: -200.0, B: 0.0},
			expected: false,
			description: "A values below -128 should be invalid",
		},
		{
			name: "Invalid B too high",
			lab: formats.LAB{L: 50.0, A: 0.0, B: 200.0},
			expected: false,
			description: "B values above 127 should be invalid",
		},
		{
			name: "Invalid B too low",
			lab: formats.LAB{L: 50.0, A: 0.0, B: -200.0},
			expected: false,
			description: "B values below -128 should be invalid",
		},
		{
			name: "Edge case L=100",
			lab: formats.LAB{L: 100.0, A: 127.0, B: 127.0},
			expected: true,
			description: "Maximum valid values should be valid",
		},
		{
			name: "Edge case L=0",
			lab: formats.LAB{L: 0.0, A: -128.0, B: -128.0},
			expected: true,
			description: "Minimum valid values should be valid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := tc.lab.IsValid()

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("LAB values: L=%.1f, A=%.1f, B=%.1f", tc.lab.L, tc.lab.A, tc.lab.B)
			t.Logf("IsValid result: %t", isValid)
			t.Logf("Expected result: %t", tc.expected)
			t.Logf("Description: %s", tc.description)

			// Verify the validation result
			if isValid != tc.expected {
				t.Errorf("IsValid mismatch: expected %t, got %t", tc.expected, isValid)
			}

			t.Logf("✅ LAB IsValid method working correctly")
		})
	}
}

// TestLAB_String tests the String method on LAB
func TestLAB_String(t *testing.T) {
	testCases := []struct {
		name string
		lab formats.LAB
		expected string
		description string
	}{
		{
			name: "White point",
			lab: formats.LAB{L: 100.0, A: 0.0, B: 0.0},
			expected: "LAB(100.00, 0.00, 0.00)",
			description: "White point should format correctly",
		},
		{
			name: "Black point",
			lab: formats.LAB{L: 0.0, A: 0.0, B: 0.0},
			expected: "LAB(0.00, 0.00, 0.00)",
			description: "Black point should format correctly",
		},
		{
			name: "Positive values",
			lab: formats.LAB{L: 75.25, A: 12.75, B: 28.50},
			expected: "LAB(75.25, 12.75, 28.50)",
			description: "Positive values should format with two decimal places",
		},
		{
			name: "Negative values",
			lab: formats.LAB{L: 25.33, A: -15.67, B: -22.11},
			expected: "LAB(25.33, -15.67, -22.11)",
			description: "Negative A and B values should format correctly",
		},
		{
			name: "Maximum values",
			lab: formats.LAB{L: 100.0, A: 127.0, B: 127.0},
			expected: "LAB(100.00, 127.00, 127.00)",
			description: "Maximum values should format correctly",
		},
		{
			name: "Minimum values",
			lab: formats.LAB{L: 0.0, A: -128.0, B: -128.0},
			expected: "LAB(0.00, -128.00, -128.00)",
			description: "Minimum values should format correctly",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.lab.String()

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("LAB values: L=%.2f, A=%.2f, B=%.2f", tc.lab.L, tc.lab.A, tc.lab.B)
			t.Logf("String result: %s", result)
			t.Logf("Expected result: %s", tc.expected)
			t.Logf("Description: %s", tc.description)

			// Verify the string formatting
			if result != tc.expected {
				t.Errorf("String format mismatch:\nExpected: %s\nGot: %s", tc.expected, result)
			}

			t.Logf("✅ LAB String method working correctly")
		})
	}
}

// TestLAB_RGBA tests the RGBA method on LAB (color.Color interface)
func TestLAB_RGBA(t *testing.T) {
	testCases := []struct {
		name string
		lab formats.LAB
		description string
	}{
		{
			name: "White point",
			lab: formats.LAB{L: 100.0, A: 0.0, B: 0.0},
			description: "White point should convert to high RGBA values",
		},
		{
			name: "Black point",
			lab: formats.LAB{L: 0.0, A: 0.0, B: 0.0},
			description: "Black point should convert to low RGBA values",
		},
		{
			name: "Mid gray",
			lab: formats.LAB{L: 50.0, A: 0.0, B: 0.0},
			description: "Mid gray should convert to medium RGBA values",
		},
		{
			name: "Red-ish color",
			lab: formats.LAB{L: 50.0, A: 50.0, B: 0.0},
			description: "Red-ish color should have higher red component",
		},
		{
			name: "Green-ish color",
			lab: formats.LAB{L: 50.0, A: -50.0, B: 0.0},
			description: "Green-ish color should have higher green component",
		},
		{
			name: "Blue-ish color",
			lab: formats.LAB{L: 50.0, A: 0.0, B: -50.0},
			description: "Blue-ish color should have higher blue component",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, g, b, a := tc.lab.RGBA()

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("LAB values: L=%.1f, A=%.1f, B=%.1f", tc.lab.L, tc.lab.A, tc.lab.B)
			t.Logf("RGBA result: R=%d, G=%d, B=%d, A=%d", r, g, b, a)
			t.Logf("RGBA normalized: R=%.3f, G=%.3f, B=%.3f, A=%.3f", 
				float64(r)/65535.0, float64(g)/65535.0, float64(b)/65535.0, float64(a)/65535.0)
			t.Logf("Description: %s", tc.description)

			// Verify that values are in valid range (0-65535 for color.Color interface)
			if r > 65535 {
				t.Errorf("R component out of range: %d > 65535", r)
			}
			if g > 65535 {
				t.Errorf("G component out of range: %d > 65535", g)
			}
			if b > 65535 {
				t.Errorf("B component out of range: %d > 65535", b)
			}
			if a > 65535 {
				t.Errorf("A component out of range: %d > 65535", a)
			}

			// Alpha should always be maximum (opaque) for LAB colors
			if a != 65535 {
				t.Errorf("Alpha should be maximum (65535), got %d", a)
			}

			t.Logf("✅ LAB RGBA method working correctly")
		})
	}
}