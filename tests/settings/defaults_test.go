package settings_test

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

func TestDefaultSettings_CoreValues(t *testing.T) {
	testCases := []struct {
		name     string
		getValue func(*settings.Settings) interface{}
		expected interface{}
	}{
		{
			name:     "Grayscale threshold",
			getValue: func(s *settings.Settings) interface{} { return s.GrayscaleThreshold },
			expected: 0.05,
		},
		{
			name:     "Monochromatic tolerance",
			getValue: func(s *settings.Settings) interface{} { return s.MonochromaticTolerance },
			expected: 15.0,
		},
		{
			name:     "Theme mode threshold",
			getValue: func(s *settings.Settings) interface{} { return s.ThemeModeThreshold },
			expected: 0.5,
		},
		{
			name:     "Min frequency",
			getValue: func(s *settings.Settings) interface{} { return s.MinFrequency },
			expected: 0.001,
		},
	}

	settings := settings.DefaultSettings()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.getValue(settings)
			
			t.Logf("Default setting '%s': expected %v (%T), got %v (%T)",
				tc.name, tc.expected, tc.expected, result, result)

			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestDefaultSettings_LoaderSettings(t *testing.T) {
	settings := settings.DefaultSettings()

	t.Logf("Loader settings:")
	t.Logf("  Max width: %d pixels", settings.LoaderMaxWidth)
	t.Logf("  Max height: %d pixels", settings.LoaderMaxHeight)
	t.Logf("  Allowed formats: %v", settings.LoaderAllowedFormats)

	// Test dimensions
	expectedMaxWidth := 8192
	expectedMaxHeight := 8192
	
	if settings.LoaderMaxWidth != expectedMaxWidth {
		t.Errorf("Expected max width %d, got %d", expectedMaxWidth, settings.LoaderMaxWidth)
	}
	if settings.LoaderMaxHeight != expectedMaxHeight {
		t.Errorf("Expected max height %d, got %d", expectedMaxHeight, settings.LoaderMaxHeight)
	}

	// Test allowed formats
	expectedFormats := []string{"jpeg", "jpg", "png", "webp"}
	if len(settings.LoaderAllowedFormats) != len(expectedFormats) {
		t.Errorf("Expected %d formats, got %d", len(expectedFormats), len(settings.LoaderAllowedFormats))
	}

	formatMap := make(map[string]bool)
	for _, format := range settings.LoaderAllowedFormats {
		formatMap[format] = true
	}

	for _, expected := range expectedFormats {
		if !formatMap[expected] {
			t.Errorf("Expected format '%s' not found in allowed formats", expected)
		}
	}

	// Verify reasonable dimension limits
	if settings.LoaderMaxWidth <= 0 || settings.LoaderMaxWidth > 32768 {
		t.Errorf("Max width %d should be between 1 and 32768", settings.LoaderMaxWidth)
	}
	if settings.LoaderMaxHeight <= 0 || settings.LoaderMaxHeight > 32768 {
		t.Errorf("Max height %d should be between 1 and 32768", settings.LoaderMaxHeight)
	}
}

func TestDefaultSettings_FallbackColors(t *testing.T) {
	testCases := []struct {
		name     string
		getValue func(*settings.Settings) string
		expected string
		description string
	}{
		{
			name:        "Light background fallback",
			getValue:    func(s *settings.Settings) string { return s.LightBackgroundFallback },
			expected:    "#ffffff",
			description: "Pure white for light themes",
		},
		{
			name:        "Dark background fallback",
			getValue:    func(s *settings.Settings) string { return s.DarkBackgroundFallback },
			expected:    "#202020",
			description: "Dark gray for dark themes",
		},
		{
			name:        "Light foreground fallback",
			getValue:    func(s *settings.Settings) string { return s.LightForegroundFallback },
			expected:    "#202020",
			description: "Dark gray text for light themes",
		},
		{
			name:        "Dark foreground fallback",
			getValue:    func(s *settings.Settings) string { return s.DarkForegroundFallback },
			expected:    "#ffffff",
			description: "White text for dark themes",
		},
		{
			name:        "Primary fallback",
			getValue:    func(s *settings.Settings) string { return s.PrimaryFallback },
			expected:    "#6496c8",
			description: "Blue accent color",
		},
	}

	settings := settings.DefaultSettings()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.getValue(settings)
			
			t.Logf("Fallback color '%s' (%s): expected %s, got %s",
				tc.name, tc.description, tc.expected, result)

			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}

			// Verify hex color format
			if len(result) != 7 || result[0] != '#' {
				t.Errorf("Color %s should be in 7-character hex format (#RRGGBB)", result)
			}

			// Verify hex digits
			for i, char := range result[1:] {
				if !isHexDigit(byte(char)) {
					t.Errorf("Character at position %d ('%c') in color %s is not a valid hex digit",
						i+1, char, result)
				}
			}
		})
	}
}

func TestDefaultSettings_ExtractionSettings(t *testing.T) {
	settings := settings.DefaultSettings()

	t.Logf("Extraction settings:")
	t.Logf("  Max candidates per category: %d", settings.Extraction.MaxCandidatesPerCategory)
	t.Logf("  Allow colored backgrounds: %t", settings.Extraction.AllowColoredBackgrounds)
	t.Logf("  Prefer vibrant accents: %t", settings.Extraction.PreferVibrantAccents)
	t.Logf("  Maintain hue consistency: %t", settings.Extraction.MaintainHueConsistency)
	t.Logf("  Grayscale hue temperature: %.1f°", settings.Extraction.GrayscaleHueTemperature)
	t.Logf("  Minimum color frequency: %.6f", settings.Extraction.MinimumColorFrequency)

	// Test reasonable values
	expectedMaxCandidates := 5
	if settings.Extraction.MaxCandidatesPerCategory != expectedMaxCandidates {
		t.Errorf("Expected max candidates %d, got %d",
			expectedMaxCandidates, settings.Extraction.MaxCandidatesPerCategory)
	}

	expectedAllowColored := false
	if settings.Extraction.AllowColoredBackgrounds != expectedAllowColored {
		t.Errorf("Expected allow colored backgrounds %t, got %t",
			expectedAllowColored, settings.Extraction.AllowColoredBackgrounds)
	}

	expectedPreferVibrant := true
	if settings.Extraction.PreferVibrantAccents != expectedPreferVibrant {
		t.Errorf("Expected prefer vibrant accents %t, got %t",
			expectedPreferVibrant, settings.Extraction.PreferVibrantAccents)
	}

	expectedMaintainHue := true
	if settings.Extraction.MaintainHueConsistency != expectedMaintainHue {
		t.Errorf("Expected maintain hue consistency %t, got %t",
			expectedMaintainHue, settings.Extraction.MaintainHueConsistency)
	}

	expectedHueTemp := 220.0
	if math.Abs(settings.Extraction.GrayscaleHueTemperature-expectedHueTemp) > 0.0001 {
		t.Errorf("Expected grayscale hue temperature %.1f°, got %.1f°",
			expectedHueTemp, settings.Extraction.GrayscaleHueTemperature)
	}

	expectedMinFreq := 0.0001
	if math.Abs(settings.Extraction.MinimumColorFrequency-expectedMinFreq) > 0.000001 {
		t.Errorf("Expected minimum color frequency %.6f, got %.6f",
			expectedMinFreq, settings.Extraction.MinimumColorFrequency)
	}

	// Verify reasonable ranges
	if settings.Extraction.MaxCandidatesPerCategory <= 0 {
		t.Errorf("Max candidates per category should be positive, got %d",
			settings.Extraction.MaxCandidatesPerCategory)
	}

	if settings.Extraction.GrayscaleHueTemperature < 0 || settings.Extraction.GrayscaleHueTemperature > 360 {
		t.Errorf("Grayscale hue temperature %.1f° should be in range [0-360°]",
			settings.Extraction.GrayscaleHueTemperature)
	}

	if settings.Extraction.MinimumColorFrequency < 0 || settings.Extraction.MinimumColorFrequency > 1 {
		t.Errorf("Minimum color frequency %.6f should be in range [0-1]",
			settings.Extraction.MinimumColorFrequency)
	}
}

func TestDefaultSettings_CategoryScoringWeights(t *testing.T) {
	settings := settings.DefaultSettings()

	weights := settings.CategoryScoring
	
	t.Logf("Category scoring weights:")
	t.Logf("  Frequency: %.3f", weights.Frequency)
	t.Logf("  Contrast: %.3f", weights.Contrast)
	t.Logf("  Saturation: %.3f", weights.Saturation)
	t.Logf("  Hue alignment: %.3f", weights.HueAlignment)
	t.Logf("  Lightness: %.3f", weights.Lightness)

	// Test individual weights
	expectedWeights := map[string]float64{
		"frequency":     0.25,
		"contrast":      0.25,
		"saturation":    0.20,
		"hue_alignment": 0.15,
		"lightness":     0.15,
	}

	actualWeights := map[string]float64{
		"frequency":     weights.Frequency,
		"contrast":      weights.Contrast,
		"saturation":    weights.Saturation,
		"hue_alignment": weights.HueAlignment,
		"lightness":     weights.Lightness,
	}

	for name, expected := range expectedWeights {
		actual := actualWeights[name]
		if math.Abs(actual-expected) > 0.0001 {
			t.Errorf("Weight '%s': expected %.3f, got %.3f", name, expected, actual)
		}
	}

	// Test that weights sum to 1.0
	totalWeight := weights.Frequency + weights.Contrast + weights.Saturation + 
		weights.HueAlignment + weights.Lightness
	
	t.Logf("Total weight sum: %.6f (should be 1.0)", totalWeight)
	
	if math.Abs(totalWeight-1.0) > 0.0001 {
		t.Errorf("Weights should sum to 1.0, got %.6f", totalWeight)
	}

	// Test that all weights are positive
	if weights.Frequency <= 0 {
		t.Errorf("Frequency weight %.3f should be positive", weights.Frequency)
	}
	if weights.Contrast <= 0 {
		t.Errorf("Contrast weight %.3f should be positive", weights.Contrast)
	}
	if weights.Saturation <= 0 {
		t.Errorf("Saturation weight %.3f should be positive", weights.Saturation)
	}
	if weights.HueAlignment <= 0 {
		t.Errorf("Hue alignment weight %.3f should be positive", weights.HueAlignment)
	}
	if weights.Lightness <= 0 {
		t.Errorf("Lightness weight %.3f should be positive", weights.Lightness)
	}
}

func TestDefaultSettings_CategoryStructure(t *testing.T) {
	settings := settings.DefaultSettings()

	// Verify category structure exists
	if settings.Categories.Dark == nil {
		t.Error("Dark category settings should not be nil")
	}
	if settings.Categories.Light == nil {
		t.Error("Light category settings should not be nil")
	}

	t.Logf("Dark mode categories: %d", len(settings.Categories.Dark))
	t.Logf("Light mode categories: %d", len(settings.Categories.Light))

	// Expected categories based on category_defaults.go
	expectedCategories := []string{
		"background", "foreground", "dim_foreground", "cursor",
		"normal_black", "normal_red", "normal_green", "normal_yellow",
		"normal_blue", "normal_magenta", "normal_cyan", "normal_white",
		"bright_black", "bright_red", "bright_green", "bright_yellow",
		"bright_blue", "bright_magenta", "bright_cyan", "bright_white",
		"accent_primary", "accent_secondary", "accent_tertiary",
		"error", "warning", "success", "info",
	}

	// Check that expected categories exist in both modes
	for _, category := range expectedCategories {
		darkChar, darkExists := settings.Categories.Dark[category]
		lightChar, lightExists := settings.Categories.Light[category]
		
		t.Logf("Category '%s': dark=%t, light=%t", category, darkExists, lightExists)
		
		if !darkExists {
			t.Errorf("Category '%s' missing from dark mode settings", category)
		}
		if !lightExists {
			t.Errorf("Category '%s' missing from light mode settings", category)
		}

		// Verify characteristics structure for existing categories
		if darkExists {
			validateCategoryCharacteristics(t, category+" (dark)", darkChar)
		}
		if lightExists {
			validateCategoryCharacteristics(t, category+" (light)", lightChar)
		}
	}

	expectedCount := len(expectedCategories)
	if len(settings.Categories.Dark) != expectedCount {
		t.Errorf("Expected %d dark categories, got %d", expectedCount, len(settings.Categories.Dark))
	}
	if len(settings.Categories.Light) != expectedCount {
		t.Errorf("Expected %d light categories, got %d", expectedCount, len(settings.Categories.Light))
	}
}

func TestSettings_WithContext(t *testing.T) {
	s := settings.DefaultSettings()
	ctx := context.Background()

	// Test setting context
	ctxWithSettings := settings.WithSettings(ctx, s)
	
	// Test retrieving from context
	retrievedSettings := settings.FromContext(ctxWithSettings)
	
	t.Logf("Original settings grayscale threshold: %.3f", s.GrayscaleThreshold)
	t.Logf("Retrieved settings grayscale threshold: %.3f", retrievedSettings.GrayscaleThreshold)

	if retrievedSettings == nil {
		t.Error("Retrieved settings should not be nil")
	}

	if retrievedSettings.GrayscaleThreshold != s.GrayscaleThreshold {
		t.Errorf("Expected grayscale threshold %.3f, got %.3f",
			s.GrayscaleThreshold, retrievedSettings.GrayscaleThreshold)
	}

	// Test context without settings - should return default
	emptyCtx := context.Background()
	defaultSettings := settings.FromContext(emptyCtx)
	
	if defaultSettings == nil {
		t.Error("Default settings should not be nil")
	}

	t.Logf("Default settings from empty context: grayscale threshold %.3f",
		defaultSettings.GrayscaleThreshold)
}

func TestSettings_LoadFromFile(t *testing.T) {
	// Create temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")
	
	configContent := `
grayscale_threshold: 0.08
monochromatic_tolerance: 20.0
theme_mode_threshold: 0.6
min_frequency: 0.002
loader_max_width: 4096
loader_max_height: 4096
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Test loading from file
	v := viper.New()
	v.SetConfigFile(configPath)
	
	err = v.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var loadedSettings settings.Settings
	err = v.Unmarshal(&loadedSettings)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	t.Logf("Loaded settings from file:")
	t.Logf("  Grayscale threshold: %.3f", loadedSettings.GrayscaleThreshold)
	t.Logf("  Monochromatic tolerance: %.1f", loadedSettings.MonochromaticTolerance)
	t.Logf("  Theme mode threshold: %.1f", loadedSettings.ThemeModeThreshold)
	t.Logf("  Min frequency: %.3f", loadedSettings.MinFrequency)
	t.Logf("  Loader max width: %d", loadedSettings.LoaderMaxWidth)

	// Verify loaded values
	expectedValues := map[string]interface{}{
		"GrayscaleThreshold":     0.08,
		"MonochromaticTolerance": 20.0,
		"ThemeModeThreshold":     0.6,
		"MinFrequency":          0.002,
		"LoaderMaxWidth":        4096,
		"LoaderMaxHeight":       4096,
	}

	actualValues := map[string]interface{}{
		"GrayscaleThreshold":     loadedSettings.GrayscaleThreshold,
		"MonochromaticTolerance": loadedSettings.MonochromaticTolerance,
		"ThemeModeThreshold":     loadedSettings.ThemeModeThreshold,
		"MinFrequency":          loadedSettings.MinFrequency,
		"LoaderMaxWidth":        loadedSettings.LoaderMaxWidth,
		"LoaderMaxHeight":       loadedSettings.LoaderMaxHeight,
	}

	for field, expected := range expectedValues {
		actual := actualValues[field]
		switch expected := expected.(type) {
		case float64:
			if actual, ok := actual.(float64); ok {
				if math.Abs(actual-expected) > 0.0001 {
					t.Errorf("Field %s: expected %.6f, got %.6f", field, expected, actual)
				}
			} else {
				t.Errorf("Field %s: expected float64, got %T", field, actual)
			}
		case int:
			if actual, ok := actual.(int); ok {
				if actual != expected {
					t.Errorf("Field %s: expected %d, got %d", field, expected, actual)
				}
			} else {
				t.Errorf("Field %s: expected int, got %T", field, actual)
			}
		}
	}
}

// Helper function to validate category characteristics
func validateCategoryCharacteristics(t *testing.T, categoryName string, char settings.CategoryCharacteristics) {
	t.Helper()
	
	// Verify lightness range is valid
	if char.MinLightness < 0 || char.MinLightness > 1 {
		t.Errorf("%s: min lightness %.3f should be in range [0-1]", categoryName, char.MinLightness)
	}
	if char.MaxLightness < 0 || char.MaxLightness > 1 {
		t.Errorf("%s: max lightness %.3f should be in range [0-1]", categoryName, char.MaxLightness)
	}
	if char.MinLightness > char.MaxLightness {
		t.Errorf("%s: min lightness %.3f should not exceed max lightness %.3f",
			categoryName, char.MinLightness, char.MaxLightness)
	}

	// Verify saturation range is valid
	if char.MinSaturation < 0 || char.MinSaturation > 1 {
		t.Errorf("%s: min saturation %.3f should be in range [0-1]", categoryName, char.MinSaturation)
	}
	if char.MaxSaturation < 0 || char.MaxSaturation > 1 {
		t.Errorf("%s: max saturation %.3f should be in range [0-1]", categoryName, char.MaxSaturation)
	}
	if char.MinSaturation > char.MaxSaturation {
		t.Errorf("%s: min saturation %.3f should not exceed max saturation %.3f",
			categoryName, char.MinSaturation, char.MaxSaturation)
	}

	// Verify contrast is reasonable
	if char.MinContrast < 0 {
		t.Errorf("%s: min contrast %.1f should not be negative", categoryName, char.MinContrast)
	}
	if char.MinContrast > 21 { // Maximum theoretical contrast ratio
		t.Errorf("%s: min contrast %.1f exceeds maximum possible (21)", categoryName, char.MinContrast)
	}

	// Verify hue settings if present
	if char.HueCenter != nil {
		if *char.HueCenter < 0 || *char.HueCenter > 360 {
			t.Errorf("%s: hue center %.1f° should be in range [0-360°]", categoryName, *char.HueCenter)
		}
	}
	if char.HueTolerance != nil {
		if *char.HueTolerance <= 0 || *char.HueTolerance > 180 {
			t.Errorf("%s: hue tolerance %.1f° should be in range (0-180°]", categoryName, *char.HueTolerance)
		}
	}

	t.Logf("%s characteristics: L[%.3f-%.3f], S[%.3f-%.3f], C≥%.1f",
		categoryName, char.MinLightness, char.MaxLightness,
		char.MinSaturation, char.MaxSaturation, char.MinContrast)
	
	if char.HueCenter != nil && char.HueTolerance != nil {
		t.Logf("  Hue: %.1f° ± %.1f°", *char.HueCenter, *char.HueTolerance)
	}
}

// Helper function to check if a byte is a valid hex digit
func isHexDigit(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'A' && b <= 'F') || (b >= 'a' && b <= 'f')
}

