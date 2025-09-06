package settings_test

import (
	"math"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

func TestCategoryDefaults_ExtractionSettings(t *testing.T) {
	s := settings.DefaultSettings()
	
	t.Logf("Testing extraction settings configuration:")
	t.Logf("  MaxCandidatesPerCategory: %d", s.Extraction.MaxCandidatesPerCategory)
	t.Logf("  AllowColoredBackgrounds: %t", s.Extraction.AllowColoredBackgrounds)
	t.Logf("  PreferVibrantAccents: %t", s.Extraction.PreferVibrantAccents)
	t.Logf("  MaintainHueConsistency: %t", s.Extraction.MaintainHueConsistency)
	t.Logf("  GrayscaleHueTemperature: %.1f", s.Extraction.GrayscaleHueTemperature)
	t.Logf("  MinimumColorFrequency: %.4f", s.Extraction.MinimumColorFrequency)

	testCases := []struct {
		name     string
		actual   interface{}
		expected interface{}
	}{
		{
			name:     "Max candidates per category",
			actual:   s.Extraction.MaxCandidatesPerCategory,
			expected: 5,
		},
		{
			name:     "Allow colored backgrounds",
			actual:   s.Extraction.AllowColoredBackgrounds,
			expected: false,
		},
		{
			name:     "Prefer vibrant accents",
			actual:   s.Extraction.PreferVibrantAccents,
			expected: true,
		},
		{
			name:     "Maintain hue consistency",
			actual:   s.Extraction.MaintainHueConsistency,
			expected: true,
		},
		{
			name:     "Grayscale hue temperature",
			actual:   s.Extraction.GrayscaleHueTemperature,
			expected: 220.0,
		},
		{
			name:     "Minimum color frequency",
			actual:   s.Extraction.MinimumColorFrequency,
			expected: 0.0001,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Extraction setting '%s': expected %v (%T), got %v (%T)",
				tc.name, tc.expected, tc.expected, tc.actual, tc.actual)

			if tc.actual != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, tc.actual)
			}
		})
	}
}

func TestCategoryDefaults_ScoringWeights(t *testing.T) {
	s := settings.DefaultSettings()

	// Calculate total weight sum for validation
	totalWeights := s.CategoryScoring.Frequency +
		s.CategoryScoring.Contrast +
		s.CategoryScoring.Saturation +
		s.CategoryScoring.HueAlignment +
		s.CategoryScoring.Lightness

	t.Logf("Total category scoring weights: %f (should sum to 1.0)", totalWeights)

	if math.Abs(totalWeights-1.0) > 0.001 {
		t.Errorf("Category scoring weights should sum to 1.0, got %f", totalWeights)
	}

	testCases := []struct {
		name     string
		actual   float64
		expected float64
	}{
		{
			name:     "Frequency weight",
			actual:   s.CategoryScoring.Frequency,
			expected: 0.25,
		},
		{
			name:     "Contrast weight",
			actual:   s.CategoryScoring.Contrast,
			expected: 0.25,
		},
		{
			name:     "Saturation weight",
			actual:   s.CategoryScoring.Saturation,
			expected: 0.20,
		},
		{
			name:     "Hue alignment weight",
			actual:   s.CategoryScoring.HueAlignment,
			expected: 0.15,
		},
		{
			name:     "Lightness weight",
			actual:   s.CategoryScoring.Lightness,
			expected: 0.15,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Scoring weight '%s': expected %.3f, got %.3f",
				tc.name, tc.expected, tc.actual)

			if tc.actual != tc.expected {
				t.Errorf("Expected %.3f, got %.3f", tc.expected, tc.actual)
			}
		})
	}
}

func TestCategoryDefaults_DarkModeCategories(t *testing.T) {
	s := settings.DefaultSettings()

	// Test core dark mode categories
	testCases := []struct {
		name                  string
		category              string
		expectedMinLightness  float64
		expectedMaxLightness  float64
		expectedMinSaturation float64
		expectedMaxSaturation float64
		expectedMinContrast   float64
	}{
		{
			name:                  "Background",
			category:              "background",
			expectedMinLightness:  0.0,
			expectedMaxLightness:  0.15,
			expectedMinSaturation: 0.0,
			expectedMaxSaturation: 0.2,
			expectedMinContrast:   0.0,
		},
		{
			name:                  "Foreground",
			category:              "foreground",
			expectedMinLightness:  0.85,
			expectedMaxLightness:  1.0,
			expectedMinSaturation: 0.0,
			expectedMaxSaturation: 0.1,
			expectedMinContrast:   4.5,
		},
		{
			name:                  "Cursor",
			category:              "cursor",
			expectedMinLightness:  0.7,
			expectedMaxLightness:  1.0,
			expectedMinSaturation: 0.0,
			expectedMaxSaturation: 0.5,
			expectedMinContrast:   7.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			characteristics, exists := s.Categories.Dark[tc.category]
			if !exists {
				t.Fatalf("Dark mode category '%s' not found", tc.category)
			}

			t.Logf("Dark mode %s category characteristics:", tc.name)
			t.Logf("  Lightness: [%.3f - %.3f], expected [%.3f - %.3f]",
				characteristics.MinLightness, characteristics.MaxLightness,
				tc.expectedMinLightness, tc.expectedMaxLightness)
			t.Logf("  Saturation: [%.3f - %.3f], expected [%.3f - %.3f]",
				characteristics.MinSaturation, characteristics.MaxSaturation,
				tc.expectedMinSaturation, tc.expectedMaxSaturation)
			t.Logf("  Min contrast: %.1f, expected %.1f",
				characteristics.MinContrast, tc.expectedMinContrast)

			// Validate lightness range
			if characteristics.MinLightness != tc.expectedMinLightness {
				t.Errorf("Min lightness: expected %.3f, got %.3f", tc.expectedMinLightness, characteristics.MinLightness)
			}
			if characteristics.MaxLightness != tc.expectedMaxLightness {
				t.Errorf("Max lightness: expected %.3f, got %.3f", tc.expectedMaxLightness, characteristics.MaxLightness)
			}

			// Validate saturation range
			if characteristics.MinSaturation != tc.expectedMinSaturation {
				t.Errorf("Min saturation: expected %.3f, got %.3f", tc.expectedMinSaturation, characteristics.MinSaturation)
			}
			if characteristics.MaxSaturation != tc.expectedMaxSaturation {
				t.Errorf("Max saturation: expected %.3f, got %.3f", tc.expectedMaxSaturation, characteristics.MaxSaturation)
			}

			// Validate contrast requirement
			if characteristics.MinContrast != tc.expectedMinContrast {
				t.Errorf("Min contrast: expected %.1f, got %.1f", tc.expectedMinContrast, characteristics.MinContrast)
			}
		})
	}
}

func TestCategoryDefaults_LightModeCategories(t *testing.T) {
	s := settings.DefaultSettings()

	// Test core light mode categories
	testCases := []struct {
		name                  string
		category              string
		expectedMinLightness  float64
		expectedMaxLightness  float64
		expectedMinSaturation float64
		expectedMaxSaturation float64
		expectedMinContrast   float64
	}{
		{
			name:                  "Background",
			category:              "background",
			expectedMinLightness:  0.9,
			expectedMaxLightness:  1.0,
			expectedMinSaturation: 0.0,
			expectedMaxSaturation: 0.15,
			expectedMinContrast:   0.0,
		},
		{
			name:                  "Foreground",
			category:              "foreground",
			expectedMinLightness:  0.0,
			expectedMaxLightness:  0.2,
			expectedMinSaturation: 0.0,
			expectedMaxSaturation: 0.1,
			expectedMinContrast:   4.5,
		},
		{
			name:                  "Cursor",
			category:              "cursor",
			expectedMinLightness:  0.0,
			expectedMaxLightness:  0.3,
			expectedMinSaturation: 0.0,
			expectedMaxSaturation: 0.5,
			expectedMinContrast:   7.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			characteristics, exists := s.Categories.Light[tc.category]
			if !exists {
				t.Fatalf("Light mode category '%s' not found", tc.category)
			}

			t.Logf("Light mode %s category characteristics:", tc.name)
			t.Logf("  Lightness: [%.3f - %.3f], expected [%.3f - %.3f]",
				characteristics.MinLightness, characteristics.MaxLightness,
				tc.expectedMinLightness, tc.expectedMaxLightness)
			t.Logf("  Saturation: [%.3f - %.3f], expected [%.3f - %.3f]",
				characteristics.MinSaturation, characteristics.MaxSaturation,
				tc.expectedMinSaturation, tc.expectedMaxSaturation)
			t.Logf("  Min contrast: %.1f, expected %.1f",
				characteristics.MinContrast, tc.expectedMinContrast)

			// Validate lightness range
			if characteristics.MinLightness != tc.expectedMinLightness {
				t.Errorf("Min lightness: expected %.3f, got %.3f", tc.expectedMinLightness, characteristics.MinLightness)
			}
			if characteristics.MaxLightness != tc.expectedMaxLightness {
				t.Errorf("Max lightness: expected %.3f, got %.3f", tc.expectedMaxLightness, characteristics.MaxLightness)
			}

			// Validate saturation range
			if characteristics.MinSaturation != tc.expectedMinSaturation {
				t.Errorf("Min saturation: expected %.3f, got %.3f", tc.expectedMinSaturation, characteristics.MinSaturation)
			}
			if characteristics.MaxSaturation != tc.expectedMaxSaturation {
				t.Errorf("Max saturation: expected %.3f, got %.3f", tc.expectedMaxSaturation, characteristics.MaxSaturation)
			}

			// Validate contrast requirement
			if characteristics.MinContrast != tc.expectedMinContrast {
				t.Errorf("Min contrast: expected %.1f, got %.1f", tc.expectedMinContrast, characteristics.MinContrast)
			}
		})
	}
}

func TestCategoryDefaults_TerminalColors(t *testing.T) {
	s := settings.DefaultSettings()

	// Test terminal colors with hue constraints
	testCases := []struct {
		name            string
		darkCategory    string
		lightCategory   string
		expectedHue     float64
		expectedTolerance float64
	}{
		{
			name:            "Red colors",
			darkCategory:    "normal_red",
			lightCategory:   "normal_red",
			expectedHue:     0.0,
			expectedTolerance: 25.0,
		},
		{
			name:            "Green colors",
			darkCategory:    "normal_green",
			lightCategory:   "normal_green",
			expectedHue:     120.0,
			expectedTolerance: 40.0,
		},
		{
			name:            "Blue colors",
			darkCategory:    "normal_blue",
			lightCategory:   "normal_blue",
			expectedHue:     240.0,
			expectedTolerance: 30.0,
		},
		{
			name:            "Yellow colors",
			darkCategory:    "normal_yellow",
			lightCategory:   "normal_yellow",
			expectedHue:     60.0,
			expectedTolerance: 20.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test dark mode
			darkChars, darkExists := s.Categories.Dark[tc.darkCategory]
			if !darkExists {
				t.Fatalf("Dark mode terminal category '%s' not found", tc.darkCategory)
			}

			// Test light mode
			lightChars, lightExists := s.Categories.Light[tc.lightCategory]
			if !lightExists {
				t.Fatalf("Light mode terminal category '%s' not found", tc.lightCategory)
			}

			t.Logf("%s terminal color settings:", tc.name)

			// Check dark mode hue settings
			if darkChars.HueCenter != nil {
				t.Logf("  Dark mode hue: %.1f° ± %.1f°", *darkChars.HueCenter,
					getToleranceValue(darkChars.HueTolerance))

				if *darkChars.HueCenter != tc.expectedHue {
					t.Errorf("Dark mode hue center: expected %.1f°, got %.1f°",
						tc.expectedHue, *darkChars.HueCenter)
				}

				tolerance := getToleranceValue(darkChars.HueTolerance)
				if tolerance != tc.expectedTolerance {
					t.Errorf("Dark mode hue tolerance: expected %.1f°, got %.1f°",
						tc.expectedTolerance, tolerance)
				}

				if tolerance <= 0 {
					t.Errorf("Dark mode tolerance %.1f° should be positive", tolerance)
				}
			} else {
				t.Errorf("Dark mode hue center should be set for terminal colors")
			}

			// Check light mode hue settings
			if lightChars.HueCenter != nil {
				t.Logf("  Light mode hue: %.1f° ± %.1f°", *lightChars.HueCenter,
					getToleranceValue(lightChars.HueTolerance))

				if *lightChars.HueCenter != tc.expectedHue {
					t.Errorf("Light mode hue center: expected %.1f°, got %.1f°",
						tc.expectedHue, *lightChars.HueCenter)
				}

				tolerance := getToleranceValue(lightChars.HueTolerance)
				if tolerance != tc.expectedTolerance {
					t.Errorf("Light mode hue tolerance: expected %.1f°, got %.1f°",
						tc.expectedTolerance, tolerance)
				}

				if tolerance <= 0 {
					t.Errorf("Light mode tolerance %.1f° should be positive", tolerance)
				}
			} else {
				t.Errorf("Light mode hue center should be set for terminal colors")
			}
		})
	}
}

func TestCategoryDefaults_SemanticColors(t *testing.T) {
	s := settings.DefaultSettings()

	// Test semantic colors with specific hue requirements
	testCases := []struct {
		name                string
		category           string
		expectedHue        float64
		expectedTolerance  float64
		description        string
	}{
		{
			name:              "Error colors",
			category:          "error",
			expectedHue:       0.0,
			expectedTolerance: 20.0,
			description:       "Red-based error indication",
		},
		{
			name:              "Warning colors",
			category:          "warning",
			expectedHue:       45.0,
			expectedTolerance: 15.0,
			description:       "Orange/yellow-based warning indication",
		},
		{
			name:              "Success colors",
			category:          "success",
			expectedHue:       120.0,
			expectedTolerance: 30.0,
			description:       "Green-based success indication",
		},
		{
			name:              "Info colors",
			category:          "info",
			expectedHue:       210.0,
			expectedTolerance: 30.0,
			description:       "Blue-based information indication",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			darkChars, darkExists := s.Categories.Dark[tc.category]
			if !darkExists {
				t.Fatalf("Dark mode semantic category '%s' not found", tc.category)
			}

			lightChars, lightExists := s.Categories.Light[tc.category]
			if !lightExists {
				t.Fatalf("Light mode semantic category '%s' not found", tc.category)
			}

			t.Logf("%s semantic color (%s):", tc.name, tc.description)

			// Validate hue settings
			if darkChars.HueCenter != nil {
				t.Logf("  Hue center: %.1f° (expected %.1f°)", *darkChars.HueCenter, tc.expectedHue)
				if *darkChars.HueCenter != tc.expectedHue {
					t.Errorf("Dark mode hue center: expected %.1f°, got %.1f°",
						tc.expectedHue, *darkChars.HueCenter)
				}
			} else {
				t.Errorf("Dark mode hue center should be set for semantic colors")
			}

			darkTolerance := getToleranceValue(darkChars.HueTolerance)
			t.Logf("  Hue tolerance: %.1f° (expected %.1f°)", darkTolerance, tc.expectedTolerance)

			if darkTolerance != tc.expectedTolerance {
				t.Errorf("Dark mode hue tolerance: expected %.1f°, got %.1f°",
					tc.expectedTolerance, darkTolerance)
			}

			// Test light mode matches
			if lightChars.HueCenter != nil && darkChars.HueCenter != nil {
				if *lightChars.HueCenter != *darkChars.HueCenter {
					t.Errorf("Light mode hue center %.1f° should match dark mode %.1f°",
						*lightChars.HueCenter, *darkChars.HueCenter)
				}
			}

			// Validate minimum contrast for visibility
			t.Logf("  Dark mode contrast requirement: %.1f", darkChars.MinContrast)
			t.Logf("  Light mode contrast requirement: %.1f", lightChars.MinContrast)

			if darkChars.MinContrast < 2.0 {
				t.Logf("Dark mode minimum contrast %.1f is below recommended 2.0", darkChars.MinContrast)
			}

			if lightChars.MinContrast < 3.0 {
				t.Logf("Light mode minimum contrast %.1f is below recommended 3.0", lightChars.MinContrast)
			}
		})
	}
}

// Helper function to safely get tolerance value
func getToleranceValue(tolerance *float64) float64 {
	if tolerance == nil {
		return 0.0
	}
	return *tolerance
}