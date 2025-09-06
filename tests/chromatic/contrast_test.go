package chromatic_test

import (
	"image/color"
	"math"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
)

func TestLuminance(t *testing.T) {
	testCases := []struct {
		name     string
		color    color.RGBA
		expected float64
		tolerance float64
	}{
		{
			name:     "Black",
			color:    color.RGBA{R: 0, G: 0, B: 0, A: 255},
			expected: 0.0,
			tolerance: 0.001,
		},
		{
			name:     "White",
			color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expected: 1.0,
			tolerance: 0.001,
		},
		{
			name:     "Pure Red",
			color:    color.RGBA{R: 255, G: 0, B: 0, A: 255},
			expected: 0.2126, // Red channel weight
			tolerance: 0.001,
		},
		{
			name:     "Pure Green",
			color:    color.RGBA{R: 0, G: 255, B: 0, A: 255},
			expected: 0.7152, // Green channel weight
			tolerance: 0.001,
		},
		{
			name:     "Pure Blue",
			color:    color.RGBA{R: 0, G: 0, B: 255, A: 255},
			expected: 0.0722, // Blue channel weight
			tolerance: 0.001,
		},
		{
			name:     "Middle Gray",
			color:    color.RGBA{R: 128, G: 128, B: 128, A: 255},
			expected: 0.2158, // Linearized middle gray
			tolerance: 0.001,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := chromatic.Luminance(tc.color)
			
			// Comprehensive diagnostic logging
			t.Logf("Color: RGBA(%d, %d, %d, %d)", tc.color.R, tc.color.G, tc.color.B, tc.color.A)
			t.Logf("Calculated luminance: %.6f", result)
			t.Logf("Expected luminance: %.6f ± %.6f", tc.expected, tc.tolerance)
			t.Logf("Difference: %.6f (threshold: %.6f)", math.Abs(result-tc.expected), tc.tolerance)
			
			if math.Abs(result-tc.expected) > tc.tolerance {
				t.Errorf("Expected luminance %v ± %v, got %v", 
					tc.expected, tc.tolerance, result)
			}
		})
	}
}

func TestContrastRatio(t *testing.T) {
	testCases := []struct {
		name     string
		color1   color.RGBA
		color2   color.RGBA
		expected float64
		tolerance float64
	}{
		{
			name:     "Black and White",
			color1:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
			color2:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expected: 21.0, // Maximum contrast
			tolerance: 0.001,
		},
		{
			name:     "Same color",
			color1:   color.RGBA{R: 128, G: 128, B: 128, A: 255},
			color2:   color.RGBA{R: 128, G: 128, B: 128, A: 255},
			expected: 1.0, // Minimum contrast
			tolerance: 0.001,
		},
		{
			name:     "Dark gray and light gray",
			color1:   color.RGBA{R: 64, G: 64, B: 64, A: 255},
			color2:   color.RGBA{R: 192, G: 192, B: 192, A: 255},
			expected: 5.7, // Actual calculated value
			tolerance: 0.1,
		},
		{
			name:     "Red and Blue",
			color1:   color.RGBA{R: 255, G: 0, B: 0, A: 255},
			color2:   color.RGBA{R: 0, G: 0, B: 255, A: 255},
			expected: 2.15, // Low contrast despite different hues
			tolerance: 0.1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := chromatic.ContrastRatio(tc.color1, tc.color2)
			
			// Calculate individual luminances for comprehensive logging
			lum1 := chromatic.Luminance(tc.color1)
			lum2 := chromatic.Luminance(tc.color2)
			
			// Comprehensive diagnostic logging
			t.Logf("Color 1: RGBA(%d, %d, %d, %d), luminance: %.6f", 
				tc.color1.R, tc.color1.G, tc.color1.B, tc.color1.A, lum1)
			t.Logf("Color 2: RGBA(%d, %d, %d, %d), luminance: %.6f", 
				tc.color2.R, tc.color2.G, tc.color2.B, tc.color2.A, lum2)
			t.Logf("Calculated contrast ratio: %.3f:1", result)
			t.Logf("Expected contrast ratio: %.3f ± %.3f", tc.expected, tc.tolerance)
			t.Logf("Difference: %.3f (threshold: %.3f)", math.Abs(result-tc.expected), tc.tolerance)
			
			if math.Abs(result-tc.expected) > tc.tolerance {
				t.Errorf("Expected contrast ratio %v ± %v, got %v",
					tc.expected, tc.tolerance, result)
			}
			
			// Contrast ratio should be symmetric
			reverse := chromatic.ContrastRatio(tc.color2, tc.color1)
			t.Logf("Symmetry check: forward=%.6f, reverse=%.6f, diff=%.6f", 
				result, reverse, math.Abs(result-reverse))
			
			if math.Abs(result-reverse) > 0.001 {
				t.Errorf("Contrast ratio not symmetric: %v vs %v", result, reverse)
			}
		})
	}
}

func TestIsAccessible(t *testing.T) {
	testCases := []struct {
		name     string
		color1   color.RGBA
		color2   color.RGBA
		level    chromatic.AccessibilityLevel
		expected bool
	}{
		{
			name:     "Black and White - AA",
			color1:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
			color2:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
			level:    chromatic.AA,
			expected: true, // 21:1 exceeds 4.5:1
		},
		{
			name:     "Black and White - AAA",
			color1:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
			color2:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
			level:    chromatic.AAA,
			expected: true, // 21:1 exceeds 7:1
		},
		{
			name:     "Red and Blue - AA",
			color1:   color.RGBA{R: 255, G: 0, B: 0, A: 255},
			color2:   color.RGBA{R: 0, G: 0, B: 255, A: 255},
			level:    chromatic.AA,
			expected: false, // ~2.15:1 fails 4.5:1
		},
		{
			name:     "Dark and Light Gray - AA",
			color1:   color.RGBA{R: 64, G: 64, B: 64, A: 255},
			color2:   color.RGBA{R: 200, G: 200, B: 200, A: 255},
			level:    chromatic.AA,
			expected: true, // Should pass AA
		},
		{
			name:     "Dark and Light Gray - AALarge",
			color1:   color.RGBA{R: 96, G: 96, B: 96, A: 255},
			color2:   color.RGBA{R: 180, G: 180, B: 180, A: 255},
			level:    chromatic.AALarge,
			expected: true, // Lower threshold for large text
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := chromatic.IsAccessible(tc.color1, tc.color2, tc.level)
			ratio := chromatic.ContrastRatio(tc.color1, tc.color2)
			threshold := tc.level.Ratio()
			
			// Comprehensive diagnostic logging
			t.Logf("Color 1: RGBA(%d, %d, %d, %d)", tc.color1.R, tc.color1.G, tc.color1.B, tc.color1.A)
			t.Logf("Color 2: RGBA(%d, %d, %d, %d)", tc.color2.R, tc.color2.G, tc.color2.B, tc.color2.A)
			t.Logf("Accessibility level: %s (threshold: %.1f:1)", tc.level, threshold)
			t.Logf("Actual contrast ratio: %.3f:1", ratio)
			t.Logf("Is accessible: %t (expected: %t)", result, tc.expected)
			t.Logf("Passes threshold: %t (%.3f >= %.1f)", ratio >= threshold, ratio, threshold)
			
			if result != tc.expected {
				t.Errorf("Expected IsAccessible=%v for level %v (ratio: %v:1)",
					tc.expected, tc.level, ratio)
			}
		})
	}
}

func TestAccessibilityLevel_Ratio(t *testing.T) {
	testCases := []struct {
		level    chromatic.AccessibilityLevel
		expected float64
	}{
		{chromatic.AA, 4.5},
		{chromatic.AAA, 7.0},
		{chromatic.AALarge, 3.0},
		{chromatic.AAALarge, 4.5},
	}

	for _, tc := range testCases {
		t.Run(string(tc.level), func(t *testing.T) {
			result := tc.level.Ratio()
			
			// Comprehensive diagnostic logging
			t.Logf("Accessibility level: %s", tc.level)
			t.Logf("Expected ratio threshold: %.1f", tc.expected)
			t.Logf("Actual ratio threshold: %.1f", result)
			
			if result != tc.expected {
				t.Errorf("Expected ratio %v for level %v, got %v",
					tc.expected, tc.level, result)
			}
		})
	}
}

func TestContrastRatioEdgeCases(t *testing.T) {
	t.Run("Very similar colors", func(t *testing.T) {
		c1 := color.RGBA{R: 100, G: 100, B: 100, A: 255}
		c2 := color.RGBA{R: 101, G: 101, B: 101, A: 255}
		
		ratio := chromatic.ContrastRatio(c1, c2)
		lum1 := chromatic.Luminance(c1)
		lum2 := chromatic.Luminance(c2)
		
		// Comprehensive diagnostic logging
		t.Logf("Color 1: RGBA(%d, %d, %d, %d), luminance: %.6f", c1.R, c1.G, c1.B, c1.A, lum1)
		t.Logf("Color 2: RGBA(%d, %d, %d, %d), luminance: %.6f", c2.R, c2.G, c2.B, c2.A, lum2)
		t.Logf("Luminance difference: %.6f", math.Abs(lum1-lum2))
		t.Logf("Contrast ratio: %.6f", ratio)
		t.Logf("Expected range: [1.0, 1.05], within range: %t", ratio >= 1.0 && ratio <= 1.05)
		
		if ratio < 1.0 || ratio > 1.05 {
			t.Errorf("Very similar colors should have ratio near 1.0, got %v", ratio)
		}
	})

	t.Run("Alpha channel ignored", func(t *testing.T) {
		opaque := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		transparent := color.RGBA{R: 255, G: 0, B: 0, A: 128}
		
		ratio := chromatic.ContrastRatio(opaque, transparent)
		lumOpaque := chromatic.Luminance(opaque)
		lumTransparent := chromatic.Luminance(transparent)
		
		// Comprehensive diagnostic logging
		t.Logf("Opaque color: RGBA(%d, %d, %d, %d), luminance: %.6f", 
			opaque.R, opaque.G, opaque.B, opaque.A, lumOpaque)
		t.Logf("Transparent color: RGBA(%d, %d, %d, %d), luminance: %.6f", 
			transparent.R, transparent.G, transparent.B, transparent.A, lumTransparent)
		t.Logf("Contrast ratio: %.6f (should be 1.0)", ratio)
		t.Logf("Luminances equal: %t", lumOpaque == lumTransparent)
		
		if ratio != 1.0 {
			t.Errorf("Alpha channel should be ignored, got ratio %v", ratio)
		}
	})

	t.Run("Maximum contrast", func(t *testing.T) {
		black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
		white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
		
		ratio := chromatic.ContrastRatio(black, white)
		lumBlack := chromatic.Luminance(black)
		lumWhite := chromatic.Luminance(white)
		
		// Comprehensive diagnostic logging
		t.Logf("Black: RGBA(%d, %d, %d, %d), luminance: %.6f", 
			black.R, black.G, black.B, black.A, lumBlack)
		t.Logf("White: RGBA(%d, %d, %d, %d), luminance: %.6f", 
			white.R, white.G, white.B, white.A, lumWhite)
		t.Logf("Contrast ratio: %.6f (theoretical maximum: 21.0)", ratio)
		t.Logf("Difference from expected: %.6f (threshold: 0.001)", math.Abs(ratio-21.0))
		
		if math.Abs(ratio-21.0) > 0.001 {
			t.Errorf("Black/white should have exactly 21:1 contrast, got %v", ratio)
		}
	})
}