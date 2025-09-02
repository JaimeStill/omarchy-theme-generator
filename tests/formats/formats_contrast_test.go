package formats_test

import (
	"image/color"
	"math"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func TestLuminance(t *testing.T) {
	testCases := []struct {
		name      string
		input     color.RGBA
		expected  float64
		tolerance float64
	}{
		{
			name:      "Black",
			input:     color.RGBA{R: 0, G: 0, B: 0, A: 255},
			expected:  0.0,
			tolerance: 0.001,
		},
		{
			name:      "White",
			input:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expected:  1.0,
			tolerance: 0.001,
		},
		{
			name:      "Red",
			input:     color.RGBA{R: 255, G: 0, B: 0, A: 255},
			expected:  0.2126,
			tolerance: 0.001,
		},
		{
			name:      "Green",
			input:     color.RGBA{R: 0, G: 255, B: 0, A: 255},
			expected:  0.7152,
			tolerance: 0.001,
		},
		{
			name:      "Blue",
			input:     color.RGBA{R: 0, G: 0, B: 255, A: 255},
			expected:  0.0722,
			tolerance: 0.001,
		},
		{
			name:      "Middle Gray",
			input:     color.RGBA{R: 128, G: 128, B: 128, A: 255},
			expected:  0.2159,
			tolerance: 0.001,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.Luminance(tc.input)

			if math.Abs(result-tc.expected) > tc.tolerance {
				t.Errorf("Expected luminance %.4f, got %.4f", tc.expected, result)
			}
		})
	}
}

func TestContrastRatio(t *testing.T) {
	testCases := []struct {
		name      string
		color1    color.RGBA
		color2    color.RGBA
		expected  float64
		tolerance float64
	}{
		{
			name:      "Black on White",
			color1:    color.RGBA{R: 0, G: 0, B: 0, A: 255},
			color2:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expected:  21.0,
			tolerance: 0.01,
		},
		{
			name:      "White on Black",
			color1:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			color2:    color.RGBA{R: 0, G: 0, B: 0, A: 255},
			expected:  21.0,
			tolerance: 0.01,
		},
		{
			name:      "Same Color",
			color1:    color.RGBA{R: 128, G: 128, B: 128, A: 255},
			color2:    color.RGBA{R: 128, G: 128, B: 128, A: 255},
			expected:  1.0,
			tolerance: 0.01,
		},
		{
			name:      "Dark Gray on Light Gray",
			color1:    color.RGBA{R: 64, G: 64, B: 64, A: 255},
			color2:    color.RGBA{R: 192, G: 192, B: 192, A: 255},
			expected:  5.70,
			tolerance: 0.1,
		},
		{
			name:      "Blue on Yellow (complementary)",
			color1:    color.RGBA{R: 0, G: 0, B: 255, A: 255},
			color2:    color.RGBA{R: 255, G: 255, B: 0, A: 255},
			expected:  8.00,
			tolerance: 0.1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.ContrastRatio(tc.color1, tc.color2)

			if math.Abs(result-tc.expected) > tc.tolerance {
				t.Errorf("Expected contrast ratio %.2f, got %.2f", tc.expected, result)
			}
		})
	}
}

func TestIsAccessible(t *testing.T) {
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	darkGray := color.RGBA{R: 64, G: 64, B: 64, A: 255}
	lightGray := color.RGBA{R: 192, G: 192, B: 192, A: 255}

	testCases := []struct {
		name     string
		color1   color.RGBA
		color2   color.RGBA
		level    formats.AccessibilityLevel
		expected bool
	}{
		// Black on white tests
		{
			name:     "Black on White - AA",
			color1:   black,
			color2:   white,
			level:    formats.AA,
			expected: true,
		},
		{
			name:     "Black on White - AAA",
			color1:   black,
			color2:   white,
			level:    formats.AAA,
			expected: true,
		},
		{
			name:     "Black on White - AA Large",
			color1:   black,
			color2:   white,
			level:    formats.AALarge,
			expected: true,
		},
		{
			name:     "Black on White - AAA Large",
			color1:   black,
			color2:   white,
			level:    formats.AAALarge,
			expected: true,
		},
		// Gray combinations
		{
			name:     "Dark Gray on Light Gray - AA",
			color1:   darkGray,
			color2:   lightGray,
			level:    formats.AA,
			expected: true, // 5.70 > 4.5
		},
		{
			name:     "Dark Gray on Light Gray - AA Large",
			color1:   darkGray,
			color2:   lightGray,
			level:    formats.AALarge,
			expected: true, // 5.70 > 3.0
		},
		{
			name:     "Dark Gray on Light Gray - AAA",
			color1:   darkGray,
			color2:   lightGray,
			level:    formats.AAA,
			expected: false, // 5.70 < 7.0
		},
		// Same color (no contrast)
		{
			name:     "Same Color - AA",
			color1:   lightGray,
			color2:   lightGray,
			level:    formats.AA,
			expected: false, // 1.0 < 4.5
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.IsAccessible(tc.color1, tc.color2, tc.level)

			if result != tc.expected {
				ratio := formats.ContrastRatio(tc.color1, tc.color2)
				t.Errorf("Expected %v for %s (ratio: %.2f, required: %.1f), got %v",
					tc.expected, tc.level, ratio, tc.level.Ratio(), result)
			}
		})
	}
}

func TestAccessibilityLevelRatios(t *testing.T) {
	testCases := []struct {
		level    formats.AccessibilityLevel
		expected float64
	}{
		{formats.AA, 4.5},
		{formats.AAA, 7.0},
		{formats.AALarge, 3.0},
		{formats.AAALarge, 4.5},
	}

	for _, tc := range testCases {
		t.Run(string(tc.level), func(t *testing.T) {
			ratio := tc.level.Ratio()
			if ratio != tc.expected {
				t.Errorf("Expected ratio %.1f for %s, got %.1f", tc.expected, tc.level, ratio)
			}
		})
	}
}
