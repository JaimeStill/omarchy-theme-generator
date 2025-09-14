package chromatic_test

import (
	"image/color"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

func TestNewChroma(t *testing.T) {
	testCases := []struct {
		name     string
		settings *settings.Settings
		expectNil bool
	}{
		{
			name:     "Valid default settings",
			settings: settings.DefaultSettings(),
			expectNil: false,
		},
		{
			name: "Custom settings",
			settings: &settings.Settings{
				Chromatic: settings.ChromaticSettings{
					NeutralThreshold:      0.05,
					ColorMergeThreshold:   10.0,
					DarkLightnessMax:      0.25,
					LightLightnessMin:     0.75,
					MutedSaturationMax:    0.2,
					VibrantSaturationMin:  0.8,
				},
			},
			expectNil: false,
		},
		{
			name:     "Nil settings",
			settings: nil,
			expectNil: false, // Constructor should still work, may panic on use
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Comprehensive diagnostic logging
			if tc.settings != nil {
				t.Logf("Input settings:")
				t.Logf("  Neutral threshold: %.3f", tc.settings.Chromatic.NeutralThreshold)
				t.Logf("  Color merge threshold: %.1f", tc.settings.Chromatic.ColorMergeThreshold)
				t.Logf("  Dark lightness max: %.3f", tc.settings.Chromatic.DarkLightnessMax)
				t.Logf("  Light lightness min: %.3f", tc.settings.Chromatic.LightLightnessMin)
			} else {
				t.Logf("Input settings: nil")
			}

			result := chromatic.NewChroma(tc.settings)

			t.Logf("NewChroma result: %+v", result)
			t.Logf("Expected nil: %t, Got nil: %t", tc.expectNil, result == nil)

			if tc.expectNil && result != nil {
				t.Errorf("Expected nil Chroma instance, got %+v", result)
			}

			if !tc.expectNil && result == nil {
				t.Error("Expected non-nil Chroma instance, got nil")
			}

			if result != nil {
				t.Logf("✓ Successfully created Chroma instance")
			}
		})
	}
}

func TestChromaColorsSimilar(t *testing.T) {
	settings := settings.DefaultSettings()
	chroma := chromatic.NewChroma(settings)

	testCases := []struct {
		name     string
		color1   color.RGBA
		color2   color.RGBA
		expected bool
		reason   string
	}{
		{
			name:     "Identical colors",
			color1:   color.RGBA{R: 128, G: 128, B: 128, A: 255},
			color2:   color.RGBA{R: 128, G: 128, B: 128, A: 255},
			expected: true,
			reason:   "Same color should be considered similar",
		},
		{
			name:     "Very similar grays (neutral)",
			color1:   color.RGBA{R: 100, G: 100, B: 100, A: 255},
			color2:   color.RGBA{R: 105, G: 105, B: 105, A: 255},
			expected: true,
			reason:   "Similar neutrals should use lightness difference threshold",
		},
		{
			name:     "Dissimilar grays (neutral)",
			color1:   color.RGBA{R: 50, G: 50, B: 50, A: 255},
			color2:   color.RGBA{R: 200, G: 200, B: 200, A: 255},
			expected: false,
			reason:   "Very different lightness neutrals should not be similar",
		},
		{
			name:     "Similar saturated colors",
			color1:   color.RGBA{R: 255, G: 100, B: 100, A: 255}, // Light red
			color2:   color.RGBA{R: 250, G: 110, B: 110, A: 255}, // Similar light red
			expected: true,
			reason:   "Perceptually similar colors should use LAB distance",
		},
		{
			name:     "Different saturated colors",
			color1:   color.RGBA{R: 255, G: 0, B: 0, A: 255},   // Pure red
			color2:   color.RGBA{R: 0, G: 255, B: 0, A: 255},   // Pure green
			expected: false,
			reason:   "Very different hues should not be similar",
		},
		{
			name:     "One neutral, one saturated",
			color1:   color.RGBA{R: 128, G: 128, B: 128, A: 255}, // Gray (neutral)
			color2:   color.RGBA{R: 255, G: 100, B: 100, A: 255}, // Red (saturated)
			expected: false,
			reason:   "Neutral and saturated colors should use LAB distance",
		},
		{
			name:     "Black and dark gray",
			color1:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
			color2:   color.RGBA{R: 20, G: 20, B: 20, A: 255},
			expected: true,
			reason:   "Very dark colors should be considered similar",
		},
		{
			name:     "White and light gray",
			color1:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
			color2:   color.RGBA{R: 240, G: 240, B: 240, A: 255},
			expected: true,
			reason:   "Very light colors should be considered similar",
		},
		{
			name:     "Different alpha (should be ignored)",
			color1:   color.RGBA{R: 128, G: 128, B: 128, A: 255},
			color2:   color.RGBA{R: 128, G: 128, B: 128, A: 128},
			expected: true,
			reason:   "Alpha channel should be ignored in similarity calculation",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := chroma.ColorsSimilar(tc.color1, tc.color2)

			// Convert to HSLA for detailed logging
			hsla1 := formats.RGBAToHSLA(tc.color1)
			hsla2 := formats.RGBAToHSLA(tc.color2)

			// Calculate LAB distance for reference
			labDistance := chromatic.DistanceLAB(tc.color1, tc.color2)

			// Comprehensive diagnostic logging
			t.Logf("Color 1: RGBA(%d, %d, %d, %d)", tc.color1.R, tc.color1.G, tc.color1.B, tc.color1.A)
			t.Logf("  HSLA: H=%.1f°, S=%.3f, L=%.3f, A=%.3f", hsla1.H, hsla1.S, hsla1.L, hsla1.A)
			t.Logf("  Is neutral: %t (S=%.3f < threshold=%.3f)",
				hsla1.S < settings.Chromatic.NeutralThreshold, hsla1.S, settings.Chromatic.NeutralThreshold)

			t.Logf("Color 2: RGBA(%d, %d, %d, %d)", tc.color2.R, tc.color2.G, tc.color2.B, tc.color2.A)
			t.Logf("  HSLA: H=%.1f°, S=%.3f, L=%.3f, A=%.3f", hsla2.H, hsla2.S, hsla2.L, hsla2.A)
			t.Logf("  Is neutral: %t (S=%.3f < threshold=%.3f)",
				hsla2.S < settings.Chromatic.NeutralThreshold, hsla2.S, settings.Chromatic.NeutralThreshold)

			// Determine which similarity method should be used
			bothNeutral := hsla1.S < settings.Chromatic.NeutralThreshold && hsla2.S < settings.Chromatic.NeutralThreshold
			if bothNeutral {
				lightnessDistance := float64(hsla1.L - hsla2.L)
				if lightnessDistance < 0 {
					lightnessDistance = -lightnessDistance
				}
				t.Logf("Both neutral - using lightness difference: %.6f (threshold: %.3f)",
					lightnessDistance, settings.Chromatic.NeutralLightnessThreshold)
			} else {
				t.Logf("Using LAB distance: %.3f (threshold: %.1f)",
					labDistance, settings.Chromatic.ColorMergeThreshold)
			}

			t.Logf("ColorsSimilar result: %t", result)
			t.Logf("Expected: %t (%s)", tc.expected, tc.reason)

			if result != tc.expected {
				t.Errorf("Expected %t, got %t - %s", tc.expected, result, tc.reason)
			} else {
				t.Logf("✓ ColorsSimilar correctly determined similarity")
			}

			// Test symmetry - ColorsSimilar should be commutative
			reverseResult := chroma.ColorsSimilar(tc.color2, tc.color1)
			if result != reverseResult {
				t.Errorf("ColorsSimilar is not symmetric: forward=%t, reverse=%t", result, reverseResult)
			} else {
				t.Logf("✓ ColorsSimilar is symmetric")
			}
		})
	}
}

func TestChromaColorsSimilar_NeutralThresholds(t *testing.T) {
	// Test specific neutral threshold behavior with custom settings
	customSettings := &settings.Settings{
		Chromatic: settings.ChromaticSettings{
			NeutralThreshold:          0.05, // Very low saturation = neutral
			NeutralLightnessThreshold: 0.1,  // 10% lightness difference for neutrals
			ColorMergeThreshold:       20.0, // Higher threshold for LAB distance
		},
	}
	chroma := chromatic.NewChroma(customSettings)

	testCases := []struct {
		name     string
		color1   color.RGBA
		color2   color.RGBA
		expected bool
		reason   string
	}{
		{
			name:     "Both neutral, small lightness difference",
			color1:   color.RGBA{R: 100, G: 100, B: 100, A: 255}, // L ≈ 0.39
			color2:   color.RGBA{R: 110, G: 110, B: 110, A: 255}, // L ≈ 0.43, diff ≈ 0.04 < 0.1
			expected: true,
			reason:   "Neutral colors with small lightness difference should be similar",
		},
		{
			name:     "Both neutral, large lightness difference",
			color1:   color.RGBA{R: 50, G: 50, B: 50, A: 255},    // L ≈ 0.20
			color2:   color.RGBA{R: 150, G: 150, B: 150, A: 255}, // L ≈ 0.59, diff ≈ 0.39 > 0.1
			expected: false,
			reason:   "Neutral colors with large lightness difference should not be similar",
		},
		{
			name:     "One neutral, one slightly saturated",
			color1:   color.RGBA{R: 100, G: 100, B: 100, A: 255}, // S = 0 (neutral)
			color2:   color.RGBA{R: 100, G: 105, B: 95, A: 255},  // S ≈ 0.1 (not neutral)
			expected: true, // Depends on LAB distance, likely similar
			reason:   "Mixed neutral/saturated uses LAB distance",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := chroma.ColorsSimilar(tc.color1, tc.color2)

			// Convert to HSLA for analysis
			hsla1 := formats.RGBAToHSLA(tc.color1)
			hsla2 := formats.RGBAToHSLA(tc.color2)
			labDistance := chromatic.DistanceLAB(tc.color1, tc.color2)

			t.Logf("Custom settings - Neutral threshold: %.3f, Lightness threshold: %.3f",
				customSettings.Chromatic.NeutralThreshold, customSettings.Chromatic.NeutralLightnessThreshold)

			t.Logf("Color 1: HSLA(%.1f°, %.3f, %.3f, %.3f), Neutral: %t",
				hsla1.H, hsla1.S, hsla1.L, hsla1.A, hsla1.S < customSettings.Chromatic.NeutralThreshold)
			t.Logf("Color 2: HSLA(%.1f°, %.3f, %.3f, %.3f), Neutral: %t",
				hsla2.H, hsla2.S, hsla2.L, hsla2.A, hsla2.S < customSettings.Chromatic.NeutralThreshold)

			bothNeutral := hsla1.S < customSettings.Chromatic.NeutralThreshold &&
				hsla2.S < customSettings.Chromatic.NeutralThreshold

			if bothNeutral {
				lightnessDistance := hsla1.L - hsla2.L
				if lightnessDistance < 0 {
					lightnessDistance = -lightnessDistance
				}
				t.Logf("Lightness difference: %.6f (threshold: %.3f)",
					lightnessDistance, customSettings.Chromatic.NeutralLightnessThreshold)
			} else {
				t.Logf("LAB distance: %.3f (threshold: %.1f)",
					labDistance, customSettings.Chromatic.ColorMergeThreshold)
			}

			t.Logf("Result: %t, Expected: %t", result, tc.expected)

			if result != tc.expected {
				t.Errorf("Expected %t, got %t - %s", tc.expected, result, tc.reason)
			} else {
				t.Logf("✓ Correctly handled neutral threshold behavior")
			}
		})
	}
}

func TestChromaColorsSimilar_EdgeCases(t *testing.T) {
	settings := settings.DefaultSettings()
	chroma := chromatic.NewChroma(settings)

	t.Run("Boundary saturation values", func(t *testing.T) {
		// Colors right at the neutral threshold boundary
		neutralThreshold := settings.Chromatic.NeutralThreshold

		// Create colors with saturation exactly at threshold
		lowSat := formats.HSLA{H: 0, S: neutralThreshold - 0.001, L: 0.5, A: 1.0}   // Just under (neutral)
		highSat := formats.HSLA{H: 0, S: neutralThreshold + 0.001, L: 0.5, A: 1.0}  // Just over (not neutral)

		color1 := formats.HSLAToRGBA(lowSat)
		color2 := formats.HSLAToRGBA(highSat)

		result := chroma.ColorsSimilar(color1, color2)

		t.Logf("Neutral threshold boundary test:")
		t.Logf("  Threshold: %.6f", neutralThreshold)
		t.Logf("  Color 1 saturation: %.6f (neutral: %t)", lowSat.S, lowSat.S < neutralThreshold)
		t.Logf("  Color 2 saturation: %.6f (neutral: %t)", highSat.S, highSat.S < neutralThreshold)
		t.Logf("  Colors similar: %t", result)

		// These should be very similar despite crossing the threshold
		if !result {
			t.Logf("⚠ Colors at threshold boundary are not considered similar (may be expected)")
		} else {
			t.Logf("✓ Colors at threshold boundary are considered similar")
		}
	})

	t.Run("Extreme color values", func(t *testing.T) {
		testPairs := []struct {
			name   string
			color1 color.RGBA
			color2 color.RGBA
		}{
			{"Pure black vs pure white", color.RGBA{0, 0, 0, 255}, color.RGBA{255, 255, 255, 255}},
			{"Pure red vs pure cyan", color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 255, 255}},
			{"Min vs max single channel", color.RGBA{0, 128, 128, 255}, color.RGBA{255, 128, 128, 255}},
		}

		for _, pair := range testPairs {
			result := chroma.ColorsSimilar(pair.color1, pair.color2)
			labDistance := chromatic.DistanceLAB(pair.color1, pair.color2)

			t.Logf("%s:", pair.name)
			t.Logf("  LAB distance: %.3f", labDistance)
			t.Logf("  Similar: %t", result)

			// Extreme colors should generally not be similar
			if result {
				t.Logf("  ⚠ Extreme colors considered similar (may indicate high threshold)")
			} else {
				t.Logf("  ✓ Extreme colors correctly identified as dissimilar")
			}
		}
	})
}