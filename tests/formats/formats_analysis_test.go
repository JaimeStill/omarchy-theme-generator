package formats_test

import (
	"image/color"
	"math"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func TestIsGrayscale(t *testing.T) {
	testCases := []struct {
		name        string
		input       color.RGBA
		expectGray  bool
		expectedHue float64
		expectedSat float64
	}{
		{
			name:        "Pure Black",
			input:       color.RGBA{R: 0, G: 0, B: 0, A: 255},
			expectGray:  true,
			expectedSat: 0.0,
		},
		{
			name:        "Pure White",
			input:       color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expectGray:  true,
			expectedSat: 0.0,
		},
		{
			name:        "Middle Gray",
			input:       color.RGBA{R: 128, G: 128, B: 128, A: 255},
			expectGray:  true,
			expectedSat: 0.0,
		},
		{
			name:        "Pure Red",
			input:       color.RGBA{R: 255, G: 0, B: 0, A: 255},
			expectGray:  false,
			expectedHue: 0.0,
			expectedSat: 1.0,
		},
		{
			name:        "Slightly Desaturated Blue",
			input:       color.RGBA{R: 120, G: 120, B: 130, A: 255},
			expectGray:  true, // Saturation < 0.05
			expectedSat: 0.04,
		},
		{
			name:        "Saturated Orange",
			input:       color.RGBA{R: 255, G: 165, B: 0, A: 255},
			expectGray:  false,
			expectedSat: 1.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hsla, isGray := formats.IsGrayscale(tc.input)

			if isGray != tc.expectGray {
				t.Errorf("Expected grayscale=%v, got %v (saturation=%.3f)", tc.expectGray, isGray, hsla.S)
			}

			// Verify the returned HSLA is correct
			if tc.expectedSat > 0 {
				if !floatEquals(hsla.S, tc.expectedSat, 0.1) {
					t.Errorf("Expected saturation %.2f, got %.2f", tc.expectedSat, hsla.S)
				}
			}
		})
	}
}

func TestIsMonochromatic(t *testing.T) {
	testCases := []struct {
		name      string
		colors    []color.RGBA
		tolerance float64
		expected  bool
	}{
		{
			name: "All Red Shades",
			colors: []color.RGBA{
				{R: 255, G: 0, B: 0, A: 255},   // Pure red
				{R: 200, G: 0, B: 0, A: 255},   // Dark red
				{R: 255, G: 50, B: 50, A: 255}, // Light red
			},
			tolerance: 15.0,
			expected:  true,
		},
		{
			name: "Mixed Hues",
			colors: []color.RGBA{
				{R: 255, G: 0, B: 0, A: 255}, // Red (0°)
				{R: 0, G: 255, B: 0, A: 255}, // Green (120°)
				{R: 0, G: 0, B: 255, A: 255}, // Blue (240°)
			},
			tolerance: 15.0,
			expected:  false,
		},
		{
			name: "All Grayscale",
			colors: []color.RGBA{
				{R: 0, G: 0, B: 0, A: 255},       // Black
				{R: 128, G: 128, B: 128, A: 255}, // Gray
				{R: 255, G: 255, B: 255, A: 255}, // White
			},
			tolerance: 15.0,
			expected:  false, // Grayscale colors have no valid hue
		},
		{
			name: "Blue to Cyan Range",
			colors: []color.RGBA{
				{R: 0, G: 0, B: 255, A: 255},   // Blue (240°)
				{R: 0, G: 128, B: 255, A: 255}, // Light blue (~210°)
				{R: 0, G: 200, B: 255, A: 255}, // Cyan-ish (~193°)
			},
			tolerance: 50.0, // Larger tolerance
			expected:  true,
		},
		{
			name:      "Single Color",
			colors:    []color.RGBA{{R: 255, G: 0, B: 0, A: 255}},
			tolerance: 15.0,
			expected:  true,
		},
		{
			name:      "Empty Array",
			colors:    []color.RGBA{},
			tolerance: 15.0,
			expected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.IsMonochromatic(tc.colors, tc.tolerance)

			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestDistanceRGB(t *testing.T) {
	testCases := []struct {
		name      string
		color1    color.RGBA
		color2    color.RGBA
		expected  float64
		tolerance float64
	}{
		{
			name:      "Same Color",
			color1:    color.RGBA{R: 128, G: 128, B: 128, A: 255},
			color2:    color.RGBA{R: 128, G: 128, B: 128, A: 255},
			expected:  0.0,
			tolerance: 0.01,
		},
		{
			name:      "Black to White",
			color1:    color.RGBA{R: 0, G: 0, B: 0, A: 255},
			color2:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expected:  441.67, // sqrt(255^2 * 3)
			tolerance: 0.1,
		},
		{
			name:      "Red to Blue",
			color1:    color.RGBA{R: 255, G: 0, B: 0, A: 255},
			color2:    color.RGBA{R: 0, G: 0, B: 255, A: 255},
			expected:  360.62, // sqrt(255^2 * 2)
			tolerance: 0.1,
		},
		{
			name:      "Small Difference",
			color1:    color.RGBA{R: 100, G: 100, B: 100, A: 255},
			color2:    color.RGBA{R: 110, G: 105, B: 95, A: 255},
			expected:  12.25, // sqrt(10^2 + 5^2 + 5^2)
			tolerance: 0.1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.DistanceRGB(tc.color1, tc.color2)

			if math.Abs(result-tc.expected) > tc.tolerance {
				t.Errorf("Expected distance %.2f, got %.2f", tc.expected, result)
			}
		})
	}
}

func TestDistanceHSL(t *testing.T) {
	testCases := []struct {
		name   string
		color1 color.RGBA
		color2 color.RGBA
	}{
		{
			name:   "Same Color",
			color1: color.RGBA{R: 128, G: 128, B: 128, A: 255},
			color2: color.RGBA{R: 128, G: 128, B: 128, A: 255},
		},
		{
			name:   "Both Grayscale",
			color1: color.RGBA{R: 0, G: 0, B: 0, A: 255},
			color2: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:   "One Grayscale",
			color1: color.RGBA{R: 128, G: 128, B: 128, A: 255},
			color2: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:   "Different Hues",
			color1: color.RGBA{R: 255, G: 0, B: 0, A: 255}, // Red
			color2: color.RGBA{R: 0, G: 0, B: 255, A: 255}, // Blue
		},
		{
			name:   "Opposite Hues",
			color1: color.RGBA{R: 255, G: 0, B: 0, A: 255},   // Red (0°)
			color2: color.RGBA{R: 0, G: 255, B: 255, A: 255}, // Cyan (180°)
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.DistanceHSL(tc.color1, tc.color2)

			// Basic sanity checks
			if result < 0 {
				t.Errorf("Distance should not be negative, got %.2f", result)
			}

			// Same color should have zero distance
			if tc.color1 == tc.color2 && result > 0.01 {
				t.Errorf("Same color should have zero distance, got %.2f", result)
			}
		})
	}
}

func TestDistanceLAB(t *testing.T) {
	testCases := []struct {
		name   string
		color1 color.RGBA
		color2 color.RGBA
	}{
		{
			name:   "Same Color",
			color1: color.RGBA{R: 128, G: 128, B: 128, A: 255},
			color2: color.RGBA{R: 128, G: 128, B: 128, A: 255},
		},
		{
			name:   "Black to White",
			color1: color.RGBA{R: 0, G: 0, B: 0, A: 255},
			color2: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:   "Red to Green",
			color1: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			color2: color.RGBA{R: 0, G: 255, B: 0, A: 255},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.DistanceLAB(tc.color1, tc.color2)

			// Basic sanity checks
			if result < 0 {
				t.Errorf("Distance should not be negative, got %.2f", result)
			}

			// Same color should have very small distance
			if tc.color1 == tc.color2 && result > 0.1 {
				t.Errorf("Same color should have near-zero distance, got %.2f", result)
			}

			// LAB distance should be more perceptually uniform
			// Red to Green should have significant distance
			if tc.name == "Red to Green" && result < 50 {
				t.Errorf("Red to Green should have significant LAB distance, got %.2f", result)
			}
		})
	}
}

func TestCalculateThemeMode(t *testing.T) {
	testCases := []struct {
		name     string
		colors   []color.RGBA
		expected formats.ThemeMode
	}{
		{
			name:     "Empty Array",
			colors:   []color.RGBA{},
			expected: formats.Dark,
		},
		{
			name: "Mostly Dark Colors",
			colors: []color.RGBA{
				{R: 0, G: 0, B: 0, A: 255},    // Black
				{R: 32, G: 32, B: 32, A: 255}, // Dark gray
				{R: 64, G: 0, B: 0, A: 255},   // Dark red
			},
			expected: formats.Light, // Dark images need light text
		},
		{
			name: "Mostly Light Colors",
			colors: []color.RGBA{
				{R: 255, G: 255, B: 255, A: 255}, // White
				{R: 200, G: 200, B: 200, A: 255}, // Light gray
				{R: 255, G: 255, B: 200, A: 255}, // Light yellow
			},
			expected: formats.Dark, // Light images need dark text
		},
		{
			name: "Mixed Colors",
			colors: []color.RGBA{
				{R: 0, G: 0, B: 0, A: 255},       // Black
				{R: 255, G: 255, B: 255, A: 255}, // White
				{R: 128, G: 128, B: 128, A: 255}, // Gray
			},
			expected: formats.Light, // Average luminance ~0.4
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.CalculateThemeMode(tc.colors)

			if result != tc.expected {
				t.Errorf("Expected %s theme mode, got %s", tc.expected, result)
			}
		})
	}
}
