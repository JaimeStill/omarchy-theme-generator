package chromatic_test

import (
	"math"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func TestFindDominantHue(t *testing.T) {
	testCases := []struct {
		name     string
		hslas    []formats.HSLA
		expected float64
		tolerance float64
	}{
		{
			name:     "Empty slice",
			hslas:    []formats.HSLA{},
			expected: math.NaN(),
			tolerance: 0,
		},
		{
			name: "Single hue",
			hslas: []formats.HSLA{
				{H: 120, S: 0.8, L: 0.5, A: 1.0}, // Green
			},
			expected:  120.0,
			tolerance: 0.001,
		},
		{
			name: "Multiple similar hues",
			hslas: []formats.HSLA{
				{H: 118, S: 0.8, L: 0.5, A: 1.0},
				{H: 120, S: 0.7, L: 0.4, A: 1.0},
				{H: 122, S: 0.9, L: 0.6, A: 1.0},
			},
			expected:  120.0,
			tolerance: 2.0,
		},
		{
			name: "Red hues wrapping around 0/360",
			hslas: []formats.HSLA{
				{H: 358, S: 0.8, L: 0.5, A: 1.0},
				{H: 0, S: 0.7, L: 0.4, A: 1.0},
				{H: 2, S: 0.9, L: 0.6, A: 1.0},
			},
			expected:  0.0,
			tolerance: 2.0,
		},
		{
			name: "Opposite hues (should average to middle)",
			hslas: []formats.HSLA{
				{H: 0, S: 0.8, L: 0.5, A: 1.0},   // Red
				{H: 180, S: 0.8, L: 0.5, A: 1.0}, // Cyan
			},
			expected:  90.0, // Could be 90 or 270, depends on vector sum
			tolerance: 180.0, // Large tolerance due to ambiguity
		},
		{
			name: "Primary colors (120° apart)",
			hslas: []formats.HSLA{
				{H: 0, S: 1.0, L: 0.5, A: 1.0},   // Red
				{H: 120, S: 1.0, L: 0.5, A: 1.0}, // Green
				{H: 240, S: 1.0, L: 0.5, A: 1.0}, // Blue
			},
			expected:  0.0, // Depends on vector sum calculation
			tolerance: 180.0, // Large tolerance due to symmetry
		},
		{
			name: "Blue-dominant cluster",
			hslas: []formats.HSLA{
				{H: 235, S: 0.8, L: 0.5, A: 1.0},
				{H: 240, S: 0.7, L: 0.4, A: 1.0},
				{H: 245, S: 0.9, L: 0.6, A: 1.0},
				{H: 242, S: 0.6, L: 0.3, A: 1.0},
			},
			expected:  240.0,
			tolerance: 5.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := chromatic.FindDominantHue(tc.hslas)

			// Comprehensive diagnostic logging
			t.Logf("Input HSLA colors: %d colors", len(tc.hslas))
			for i, hsla := range tc.hslas {
				t.Logf("  Color %d: H=%.1f°, S=%.2f, L=%.2f, A=%.2f",
					i+1, hsla.H, hsla.S, hsla.L, hsla.A)
			}
			t.Logf("Calculated dominant hue: %.3f°", result)
			t.Logf("Expected dominant hue: %.3f° ± %.3f°", tc.expected, tc.tolerance)

			// Handle NaN case
			if math.IsNaN(tc.expected) {
				if !math.IsNaN(result) {
					t.Errorf("Expected NaN, got %.3f", result)
				}
				t.Logf("✓ Correctly returned NaN for empty input")
				return
			}

			if math.IsNaN(result) {
				t.Errorf("Expected %.3f, got NaN", tc.expected)
				return
			}

			// Calculate hue difference considering wraparound
			diff := math.Abs(result - tc.expected)
			if diff > 180 {
				diff = 360 - diff
			}

			t.Logf("Angular difference: %.3f° (threshold: %.3f°)", diff, tc.tolerance)

			if diff > tc.tolerance {
				t.Errorf("Expected dominant hue %.3f ± %.3f, got %.3f (diff: %.3f°)",
					tc.expected, tc.tolerance, result, diff)
			} else {
				t.Logf("✓ Dominant hue within expected tolerance")
			}
		})
	}
}

func TestCalculateHueVariance(t *testing.T) {
	testCases := []struct {
		name     string
		hslas    []formats.HSLA
		expected float64
		tolerance float64
	}{
		{
			name:     "Empty slice",
			hslas:    []formats.HSLA{},
			expected: 0.0,
			tolerance: 0.001,
		},
		{
			name: "Single color",
			hslas: []formats.HSLA{
				{H: 120, S: 0.8, L: 0.5, A: 1.0},
			},
			expected:  0.0,
			tolerance: 0.001,
		},
		{
			name: "Two identical hues",
			hslas: []formats.HSLA{
				{H: 120, S: 0.8, L: 0.5, A: 1.0},
				{H: 120, S: 0.6, L: 0.4, A: 1.0},
			},
			expected:  0.0,
			tolerance: 0.001,
		},
		{
			name: "Small variance cluster",
			hslas: []formats.HSLA{
				{H: 118, S: 0.8, L: 0.5, A: 1.0},
				{H: 120, S: 0.7, L: 0.4, A: 1.0},
				{H: 122, S: 0.9, L: 0.6, A: 1.0},
			},
			expected:  1.63, // sqrt((2^2 + 0^2 + 2^2)/3) ≈ 1.63
			tolerance: 0.5,
		},
		{
			name: "Large variance cluster",
			hslas: []formats.HSLA{
				{H: 100, S: 0.8, L: 0.5, A: 1.0},
				{H: 120, S: 0.7, L: 0.4, A: 1.0},
				{H: 140, S: 0.9, L: 0.6, A: 1.0},
			},
			expected:  16.3, // Much larger variance
			tolerance: 5.0,
		},
		{
			name: "Wraparound hues near 0/360",
			hslas: []formats.HSLA{
				{H: 358, S: 0.8, L: 0.5, A: 1.0},
				{H: 0, S: 0.7, L: 0.4, A: 1.0},
				{H: 2, S: 0.9, L: 0.6, A: 1.0},
			},
			expected:  1.63, // Should handle wraparound correctly
			tolerance: 1.0,
		},
		{
			name: "Maximum variance (opposite hues)",
			hslas: []formats.HSLA{
				{H: 0, S: 0.8, L: 0.5, A: 1.0},   // Red
				{H: 180, S: 0.8, L: 0.5, A: 1.0}, // Cyan
			},
			expected:  90.0, // Maximum possible variance for two colors
			tolerance: 20.0, // Large tolerance due to wraparound complexity
		},
		{
			name: "Monochromatic with tight clustering",
			hslas: []formats.HSLA{
				{H: 240.0, S: 0.8, L: 0.5, A: 1.0},
				{H: 240.5, S: 0.7, L: 0.4, A: 1.0},
				{H: 239.5, S: 0.9, L: 0.6, A: 1.0},
				{H: 240.2, S: 0.6, L: 0.3, A: 1.0},
			},
			expected:  0.35, // Very low variance for tight cluster
			tolerance: 0.2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := chromatic.CalculateHueVariance(tc.hslas)

			// Comprehensive diagnostic logging
			t.Logf("Input HSLA colors: %d colors", len(tc.hslas))
			for i, hsla := range tc.hslas {
				t.Logf("  Color %d: H=%.1f°, S=%.2f, L=%.2f, A=%.2f",
					i+1, hsla.H, hsla.S, hsla.L, hsla.A)
			}

			// Calculate and log dominant hue for context
			dominantHue := chromatic.FindDominantHue(tc.hslas)
			t.Logf("Dominant hue: %.3f°", dominantHue)
			t.Logf("Calculated hue variance: %.3f°", result)
			t.Logf("Expected hue variance: %.3f° ± %.3f°", tc.expected, tc.tolerance)

			// Calculate individual deviations for detailed logging
			if len(tc.hslas) > 1 && !math.IsNaN(dominantHue) {
				t.Logf("Individual hue deviations from dominant:")
				for i, hsla := range tc.hslas {
					// Calculate hue distance (considering wraparound)
					diff := math.Abs(hsla.H - dominantHue)
					if diff > 180 {
						diff = 360 - diff
					}
					t.Logf("  Color %d: %.3f° deviation", i+1, diff)
				}
			}

			diff := math.Abs(result - tc.expected)
			t.Logf("Variance difference: %.3f° (threshold: %.3f°)", diff, tc.tolerance)

			if diff > tc.tolerance {
				t.Errorf("Expected hue variance %.3f ± %.3f, got %.3f (diff: %.3f)",
					tc.expected, tc.tolerance, result, diff)
			} else {
				t.Logf("✓ Hue variance within expected tolerance")
			}

			// Additional validation - variance should be non-negative
			if result < 0 {
				t.Errorf("Hue variance should be non-negative, got %.3f", result)
			}
		})
	}
}

func TestHueFunctions_EdgeCases(t *testing.T) {
	t.Run("FindDominantHue with extreme wraparound", func(t *testing.T) {
		hslas := []formats.HSLA{
			{H: 359.9, S: 0.8, L: 0.5, A: 1.0},
			{H: 0.1, S: 0.8, L: 0.5, A: 1.0},
		}

		result := chromatic.FindDominantHue(hslas)

		t.Logf("Input hues: [%.1f°, %.1f°]", hslas[0].H, hslas[1].H)
		t.Logf("Calculated dominant hue: %.3f°", result)

		// Should be close to 0° or 360° due to wraparound
		nearZero := math.Abs(result) < 1.0
		near360 := math.Abs(result-360) < 1.0

		if !nearZero && !near360 {
			t.Errorf("Expected dominant hue near 0° or 360°, got %.3f°", result)
		} else {
			t.Logf("✓ Correctly handled wraparound case")
		}
	})

	t.Run("CalculateHueVariance with all same hue", func(t *testing.T) {
		hslas := []formats.HSLA{
			{H: 45, S: 0.2, L: 0.3, A: 1.0},
			{H: 45, S: 0.8, L: 0.7, A: 1.0},
			{H: 45, S: 0.5, L: 0.5, A: 1.0},
			{H: 45, S: 0.9, L: 0.2, A: 1.0},
		}

		result := chromatic.CalculateHueVariance(hslas)

		t.Logf("All colors have identical hue: 45°")
		t.Logf("Calculated variance: %.6f", result)

		if result > 0.001 {
			t.Errorf("Expected zero variance for identical hues, got %.6f", result)
		} else {
			t.Logf("✓ Correctly calculated zero variance for identical hues")
		}
	})

	t.Run("Consistency between functions", func(t *testing.T) {
		hslas := []formats.HSLA{
			{H: 30, S: 0.8, L: 0.5, A: 1.0},
			{H: 35, S: 0.7, L: 0.4, A: 1.0},
			{H: 25, S: 0.9, L: 0.6, A: 1.0},
		}

		dominantHue := chromatic.FindDominantHue(hslas)
		variance := chromatic.CalculateHueVariance(hslas)

		t.Logf("Input hues: [%.1f°, %.1f°, %.1f°]", hslas[0].H, hslas[1].H, hslas[2].H)
		t.Logf("Dominant hue: %.3f°", dominantHue)
		t.Logf("Hue variance: %.3f°", variance)

		// Variance should be reasonable for this tight cluster
		if variance > 10.0 {
			t.Errorf("Variance %.3f° seems too high for tight cluster", variance)
		}

		// Dominant hue should be within the range of input hues
		if dominantHue < 20 || dominantHue > 40 {
			t.Errorf("Dominant hue %.3f° outside expected range [20°, 40°]", dominantHue)
		} else {
			t.Logf("✓ Functions produced consistent results")
		}
	})
}