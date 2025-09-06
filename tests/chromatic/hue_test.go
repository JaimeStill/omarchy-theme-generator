package chromatic_test

import (
	"math"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

func TestFindDominantHue(t *testing.T) {
	testCases := []struct {
		name     string
		colors   []formats.HSLA
		expected float64
		tolerance float64
	}{
		{
			name: "All same hue",
			colors: []formats.HSLA{
				{H: 120, S: 0.5, L: 0.5, A: 1.0},
				{H: 120, S: 0.8, L: 0.3, A: 1.0},
				{H: 120, S: 0.3, L: 0.7, A: 1.0},
			},
			expected: 120.0,
			tolerance: 0.1,
		},
		{
			name: "Close hues",
			colors: []formats.HSLA{
				{H: 100, S: 0.5, L: 0.5, A: 1.0},
				{H: 110, S: 0.5, L: 0.5, A: 1.0},
				{H: 120, S: 0.5, L: 0.5, A: 1.0},
			},
			expected: 110.0, // Average
			tolerance: 5.0,
		},
		{
			name: "Hues around zero crossing",
			colors: []formats.HSLA{
				{H: 350, S: 0.5, L: 0.5, A: 1.0},
				{H: 10, S: 0.5, L: 0.5, A: 1.0},
			},
			expected: 0.0, // Should handle wraparound
			tolerance: 10.0,
		},
		{
			name: "Single color",
			colors: []formats.HSLA{
				{H: 240, S: 0.5, L: 0.5, A: 1.0},
			},
			expected: 240.0,
			tolerance: 0.1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := chromatic.FindDominantHue(tc.colors)
			
			// Handle circular hue comparison
			diff := math.Abs(result - tc.expected)
			if diff > 180 {
				diff = 360 - diff
			}
			
			if diff > tc.tolerance {
				t.Errorf("Expected dominant hue %v ± %v, got %v",
					tc.expected, tc.tolerance, result)
			}
		})
	}
}

func TestCalculateHueVariance(t *testing.T) {
	testCases := []struct {
		name        string
		colors      []formats.HSLA
		minVariance float64
		maxVariance float64
	}{
		{
			name: "No variance - same hue",
			colors: []formats.HSLA{
				{H: 120, S: 0.5, L: 0.5, A: 1.0},
				{H: 120, S: 0.8, L: 0.3, A: 1.0},
			},
			minVariance: 0.0,
			maxVariance: 0.1,
		},
		{
			name: "Small variance",
			colors: []formats.HSLA{
				{H: 100, S: 0.5, L: 0.5, A: 1.0},
				{H: 110, S: 0.5, L: 0.5, A: 1.0},
				{H: 120, S: 0.5, L: 0.5, A: 1.0},
			},
			minVariance: 5.0,
			maxVariance: 15.0,
		},
		{
			name: "Large variance",
			colors: []formats.HSLA{
				{H: 0, S: 0.5, L: 0.5, A: 1.0},
				{H: 120, S: 0.5, L: 0.5, A: 1.0},
				{H: 240, S: 0.5, L: 0.5, A: 1.0},
			},
			minVariance: 80.0,
			maxVariance: 120.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := chromatic.CalculateHueVariance(tc.colors)
			
			if result < tc.minVariance || result > tc.maxVariance {
				t.Errorf("Expected variance between %v and %v, got %v",
					tc.minVariance, tc.maxVariance, result)
			}
		})
	}
}

func TestChroma_HuesWithinTolerance(t *testing.T) {
	s := settings.DefaultSettings()
	s.MonochromaticTolerance = 15.0
	c := chromatic.NewChroma(s)

	testCases := []struct {
		name     string
		hue1     float64
		hue2     float64
		expected bool
	}{
		{
			name:     "Same hue",
			hue1:     120.0,
			hue2:     120.0,
			expected: true,
		},
		{
			name:     "Within tolerance",
			hue1:     120.0,
			hue2:     130.0,
			expected: true,
		},
		{
			name:     "Outside tolerance",
			hue1:     120.0,
			hue2:     140.0,
			expected: false,
		},
		{
			name:     "Wraparound - within tolerance",
			hue1:     5.0,
			hue2:     355.0,
			expected: true, // 10 degrees apart
		},
		{
			name:     "Wraparound - outside tolerance",
			hue1:     20.0,
			hue2:     340.0,
			expected: false, // 40 degrees apart
		},
		{
			name:     "Opposite hues",
			hue1:     0.0,
			hue2:     180.0,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := c.HuesWithinTolerance(tc.hue1, tc.hue2)
			if result != tc.expected {
				t.Errorf("Expected HuesWithinTolerance(%v, %v) = %v, got %v",
					tc.hue1, tc.hue2, tc.expected, result)
			}
		})
	}
}

func TestChroma_IdentifyColorScheme(t *testing.T) {
	s := settings.DefaultSettings()
	c := chromatic.NewChroma(s)

	testCases := []struct {
		name       string
		variance   float64
		colorCount int
		hues       []float64
		expected   chromatic.ColorScheme
	}{
		{
			name:       "Grayscale - no colors",
			variance:   0,
			colorCount: 0,
			hues:       []float64{},
			expected:   chromatic.Grayscale,
		},
		{
			name:       "Monochromatic - single color",
			variance:   0,
			colorCount: 1,
			hues:       []float64{120},
			expected:   chromatic.Monochromatic,
		},
		{
			name:       "Complementary - two opposite hues",
			variance:   180,
			colorCount: 2,
			hues:       []float64{0, 180},
			expected:   chromatic.Complementary,
		},
		{
			name:       "Analogous - close hues",
			variance:   25,
			colorCount: 3,
			hues:       []float64{100, 120, 125},
			expected:   chromatic.Analogous,
		},
		{
			name:       "Triadic - three evenly spaced",
			variance:   120,
			colorCount: 3,
			hues:       []float64{0, 120, 240},
			expected:   chromatic.Triadic,
		},
		{
			name:       "Custom - no pattern",
			variance:   90,
			colorCount: 5,
			hues:       []float64{0, 45, 90, 180, 270},
			expected:   chromatic.Custom,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := c.IdentifyColorScheme(tc.variance, tc.colorCount, tc.hues)
			if result != tc.expected {
				t.Errorf("Expected color scheme %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestHueDistance(t *testing.T) {
	testCases := []struct {
		name     string
		h1       float64
		h2       float64
		expected float64
	}{
		{
			name:     "Same hue",
			h1:       120,
			h2:       120,
			expected: 0,
		},
		{
			name:     "Simple difference",
			h1:       100,
			h2:       150,
			expected: 50,
		},
		{
			name:     "Wraparound shorter",
			h1:       10,
			h2:       350,
			expected: 20, // Going backwards is shorter
		},
		{
			name:     "Half circle",
			h1:       0,
			h2:       180,
			expected: 180,
		},
		{
			name:     "Reverse order",
			h1:       200,
			h2:       100,
			expected: 100, // Should be positive
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Since hueDistance is likely private, test through HuesWithinTolerance
			// or other public methods that use it
			s := settings.DefaultSettings()
			s.MonochromaticTolerance = tc.expected + 0.1
			c := chromatic.NewChroma(s)
			
			// If distance is less than expected, should be within tolerance
			result := c.HuesWithinTolerance(tc.h1, tc.h2)
			if !result && tc.expected < s.MonochromaticTolerance {
				t.Errorf("Hues %v and %v should be within tolerance %v",
					tc.h1, tc.h2, s.MonochromaticTolerance)
			}
		})
	}
}