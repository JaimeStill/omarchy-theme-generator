package processor_test

import (
	"image/color"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/processor"
)

func TestNewWeightedColor(t *testing.T) {
	testCases := []struct {
		name           string
		color          color.RGBA
		frequency      uint32
		total          uint32
		expectedWeight float64
	}{
		{
			name:           "High frequency color",
			color:          color.RGBA{R: 255, G: 0, B: 0, A: 255}, // Red
			frequency:      500,
			total:          1000,
			expectedWeight: 0.5,
		},
		{
			name:           "Low frequency color",
			color:          color.RGBA{R: 0, G: 255, B: 0, A: 255}, // Green
			frequency:      10,
			total:          1000,
			expectedWeight: 0.01,
		},
		{
			name:           "Maximum frequency (100%)",
			color:          color.RGBA{R: 0, G: 0, B: 255, A: 255}, // Blue
			frequency:      1000,
			total:          1000,
			expectedWeight: 1.0,
		},
		{
			name:           "Minimum frequency",
			color:          color.RGBA{R: 128, G: 128, B: 128, A: 255}, // Gray
			frequency:      1,
			total:          10000,
			expectedWeight: 0.0001,
		},
		{
			name:           "Zero frequency",
			color:          color.RGBA{R: 255, G: 255, B: 0, A: 255}, // Yellow
			frequency:      0,
			total:          1000,
			expectedWeight: 0.0,
		},
		{
			name:           "Single pixel image",
			color:          color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White
			frequency:      1,
			total:          1,
			expectedWeight: 1.0,
		},
		{
			name:           "High precision calculation",
			color:          color.RGBA{R: 100, G: 50, B: 200, A: 255},
			frequency:      12345,
			total:          1000000,
			expectedWeight: 0.012345,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := processor.NewWeightedColor(tc.color, tc.frequency, tc.total)

			// Comprehensive diagnostic logging
			t.Logf("Input color: RGBA(%d, %d, %d, %d)", tc.color.R, tc.color.G, tc.color.B, tc.color.A)
			t.Logf("Frequency: %d pixels", tc.frequency)
			t.Logf("Total pixels: %d", tc.total)
			t.Logf("Expected weight: %.6f", tc.expectedWeight)
			t.Logf("Calculated weight: %.6f", result.Weight)
			t.Logf("Weight difference: %.8f", abs64(result.Weight-tc.expectedWeight))

			// Verify the color was set correctly
			if result.RGBA != tc.color {
				t.Errorf("Expected color RGBA(%d, %d, %d, %d), got RGBA(%d, %d, %d, %d)",
					tc.color.R, tc.color.G, tc.color.B, tc.color.A,
					result.RGBA.R, result.RGBA.G, result.RGBA.B, result.RGBA.A)
			} else {
				t.Logf("✓ Color correctly set: RGBA(%d, %d, %d, %d)",
					result.RGBA.R, result.RGBA.G, result.RGBA.B, result.RGBA.A)
			}

			// Verify the weight calculation (allowing for small floating point errors)
			tolerance := 0.0000001
			if abs64(result.Weight-tc.expectedWeight) > tolerance {
				t.Errorf("Expected weight %.8f, got %.8f (difference: %.8f, tolerance: %.8f)",
					tc.expectedWeight, result.Weight,
					abs64(result.Weight-tc.expectedWeight), tolerance)
			} else {
				t.Logf("✓ Weight correctly calculated within tolerance")
			}

			// Verify weight is in valid range [0, 1]
			if result.Weight < 0.0 || result.Weight > 1.0 {
				t.Errorf("Weight %.6f is outside valid range [0.0, 1.0]", result.Weight)
			} else {
				t.Logf("✓ Weight is within valid range [0.0, 1.0]")
			}

			// Additional validation for percentage representation
			percentage := result.Weight * 100
			t.Logf("Weight as percentage: %.4f%%", percentage)

			// Verify consistency with input parameters
			if tc.total > 0 {
				expectedWeightFromRatio := float64(tc.frequency) / float64(tc.total)
				if abs64(result.Weight-expectedWeightFromRatio) > tolerance {
					t.Errorf("Weight calculation inconsistent: expected %.8f from ratio, got %.8f",
						expectedWeightFromRatio, result.Weight)
				}
			}
		})
	}
}

func TestWeightedColor_EdgeCases(t *testing.T) {
	t.Run("Zero total (division by zero protection)", func(t *testing.T) {
		color := color.RGBA{R: 255, G: 128, B: 64, A: 255}
		frequency := uint32(100)
		total := uint32(0)

		// This should handle division by zero gracefully
		result := processor.NewWeightedColor(color, frequency, total)

		t.Logf("Edge case: zero total")
		t.Logf("Color: RGBA(%d, %d, %d, %d)", color.R, color.G, color.B, color.A)
		t.Logf("Frequency: %d", frequency)
		t.Logf("Total: %d (zero - division by zero case)", total)
		t.Logf("Resulting weight: %.6f", result.Weight)

		// The behavior here depends on implementation - could be 0, NaN, or Inf
		// We mainly want to ensure it doesn't crash
		t.Logf("✓ Function handled zero total without crashing")

		// Check if the result is reasonable
		if result.Weight != 0.0 && !isNaN64(result.Weight) && !isInf64(result.Weight) {
			t.Logf("⚠ Unexpected weight value for zero total: %.6f", result.Weight)
		}
	})

	t.Run("Maximum possible values", func(t *testing.T) {
		color := color.RGBA{R: 255, G: 255, B: 255, A: 255}
		maxUint32 := ^uint32(0) // Maximum uint32 value
		frequency := maxUint32
		total := maxUint32

		result := processor.NewWeightedColor(color, frequency, total)

		t.Logf("Edge case: maximum uint32 values")
		t.Logf("Frequency: %d (max uint32)", frequency)
		t.Logf("Total: %d (max uint32)", total)
		t.Logf("Resulting weight: %.6f", result.Weight)

		// Should result in weight = 1.0
		expectedWeight := 1.0
		tolerance := 0.0000001

		if abs64(result.Weight-expectedWeight) > tolerance {
			t.Errorf("Expected weight %.6f for equal max values, got %.6f",
				expectedWeight, result.Weight)
		} else {
			t.Logf("✓ Correctly handled maximum uint32 values")
		}
	})

	t.Run("Frequency greater than total", func(t *testing.T) {
		color := color.RGBA{R: 128, G: 64, B: 32, A: 255}
		frequency := uint32(1500)
		total := uint32(1000)

		result := processor.NewWeightedColor(color, frequency, total)

		t.Logf("Edge case: frequency > total")
		t.Logf("Color: RGBA(%d, %d, %d, %d)", color.R, color.G, color.B, color.A)
		t.Logf("Frequency: %d", frequency)
		t.Logf("Total: %d", total)
		t.Logf("Resulting weight: %.6f", result.Weight)

		// This is mathematically invalid but we test the behavior
		// The weight will be > 1.0, which may or may not be intended
		if result.Weight > 1.0 {
			t.Logf("⚠ Weight %.6f > 1.0 when frequency > total", result.Weight)
		} else {
			t.Logf("✓ Implementation handles frequency > total case")
		}
	})

	t.Run("Very small weights", func(t *testing.T) {
		color := color.RGBA{R: 1, G: 1, B: 1, A: 255}
		frequency := uint32(1)
		total := uint32(10000000) // 10 million

		result := processor.NewWeightedColor(color, frequency, total)

		expectedWeight := 0.0000001 // 1 in 10 million
		t.Logf("Very small weight test:")
		t.Logf("Frequency: %d", frequency)
		t.Logf("Total: %d", total)
		t.Logf("Expected weight: %.8f", expectedWeight)
		t.Logf("Calculated weight: %.8f", result.Weight)

		tolerance := 0.00000001
		if abs64(result.Weight-expectedWeight) > tolerance {
			t.Errorf("Expected weight %.8f, got %.8f", expectedWeight, result.Weight)
		} else {
			t.Logf("✓ Correctly calculated very small weight")
		}

		// Verify precision is maintained
		if result.Weight == 0.0 {
			t.Logf("⚠ Weight rounded to zero - may indicate precision loss")
		}
	})
}

func TestWeightedColor_Consistency(t *testing.T) {
	t.Run("Multiple calls with same parameters", func(t *testing.T) {
		color := color.RGBA{R: 200, G: 100, B: 50, A: 255}
		frequency := uint32(750)
		total := uint32(2000)

		// Call the function multiple times with the same parameters
		results := make([]processor.WeightedColor, 5)
		for i := 0; i < len(results); i++ {
			results[i] = processor.NewWeightedColor(color, frequency, total)
		}

		t.Logf("Consistency test: 5 calls with same parameters")
		t.Logf("Parameters: RGBA(%d, %d, %d, %d), freq=%d, total=%d",
			color.R, color.G, color.B, color.A, frequency, total)

		// All results should be identical
		for i := 1; i < len(results); i++ {
			if results[i].RGBA != results[0].RGBA {
				t.Errorf("Color inconsistency at call %d: expected RGBA(%d, %d, %d, %d), got RGBA(%d, %d, %d, %d)",
					i, results[0].RGBA.R, results[0].RGBA.G, results[0].RGBA.B, results[0].RGBA.A,
					results[i].RGBA.R, results[i].RGBA.G, results[i].RGBA.B, results[i].RGBA.A)
			}

			if results[i].Weight != results[0].Weight {
				t.Errorf("Weight inconsistency at call %d: expected %.8f, got %.8f",
					i, results[0].Weight, results[i].Weight)
			}
		}

		t.Logf("All results: weight=%.6f", results[0].Weight)
		t.Logf("✓ All calls produced identical results")
	})

	t.Run("Weight proportionality", func(t *testing.T) {
		// Test that doubling frequency doubles the weight (for fixed total)
		color := color.RGBA{R: 64, G: 128, B: 192, A: 255}
		total := uint32(1000)

		wc1 := processor.NewWeightedColor(color, 100, total) // 10%
		wc2 := processor.NewWeightedColor(color, 200, total) // 20%

		t.Logf("Proportionality test:")
		t.Logf("Weight 1 (freq=100): %.6f", wc1.Weight)
		t.Logf("Weight 2 (freq=200): %.6f", wc2.Weight)
		t.Logf("Ratio: %.6f (expected: 2.0)", wc2.Weight/wc1.Weight)

		expectedRatio := 2.0
		actualRatio := wc2.Weight / wc1.Weight
		tolerance := 0.0001

		if abs64(actualRatio-expectedRatio) > tolerance {
			t.Errorf("Expected ratio %.3f, got %.6f", expectedRatio, actualRatio)
		} else {
			t.Logf("✓ Weights are correctly proportional")
		}
	})
}

// Helper functions for floating point operations
func abs64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func isNaN64(x float64) bool {
	return x != x
}

func isInf64(x float64) bool {
	return x > 1e308 || x < -1e308
}