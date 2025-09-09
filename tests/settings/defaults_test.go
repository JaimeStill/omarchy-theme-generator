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
	s := settings.DefaultSettings()

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
			name:     "Grayscale image threshold",
			getValue: func(s *settings.Settings) interface{} { return s.GrayscaleImageThreshold },
			expected: 0.8,
		},
		{
			name:     "Monochromatic tolerance",
			getValue: func(s *settings.Settings) interface{} { return s.MonochromaticTolerance },
			expected: 15.0,
		},
		{
			name:     "Monochromatic weight threshold",
			getValue: func(s *settings.Settings) interface{} { return s.MonochromaticWeightThreshold },
			expected: 0.1,
		},
		{
			name:     "Theme mode threshold",
			getValue: func(s *settings.Settings) interface{} { return s.ThemeModeThreshold },
			expected: 0.5,
		},
		{
			name:     "Min frequency",
			getValue: func(s *settings.Settings) interface{} { return s.MinFrequency },
			expected: 0.0001,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.getValue(s)
			
			t.Logf("Default setting '%s': expected %v (%T), got %v (%T)",
				tc.name, tc.expected, tc.expected, result, result)

			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestDefaultSettings_LoaderSettings(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Loader settings:")
	t.Logf("  Max width: %d pixels", s.LoaderMaxWidth)
	t.Logf("  Max height: %d pixels", s.LoaderMaxHeight)
	t.Logf("  Allowed formats: %v", s.LoaderAllowedFormats)

	// Test dimensions
	expectedMaxWidth := 8192
	expectedMaxHeight := 8192
	
	if s.LoaderMaxWidth != expectedMaxWidth {
		t.Errorf("Expected max width %d, got %d", expectedMaxWidth, s.LoaderMaxWidth)
	}
	if s.LoaderMaxHeight != expectedMaxHeight {
		t.Errorf("Expected max height %d, got %d", expectedMaxHeight, s.LoaderMaxHeight)
	}

	// Test allowed formats
	expectedFormats := []string{"jpeg", "jpg", "png", "webp"}
	if len(s.LoaderAllowedFormats) != len(expectedFormats) {
		t.Errorf("Expected %d formats, got %d", len(expectedFormats), len(s.LoaderAllowedFormats))
	}

	formatMap := make(map[string]bool)
	for _, format := range s.LoaderAllowedFormats {
		formatMap[format] = true
	}

	for _, expected := range expectedFormats {
		if !formatMap[expected] {
			t.Errorf("Expected format '%s' not found in allowed formats", expected)
		}
	}
}

func TestDefaultSettings_GroupingThresholds(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Grouping threshold settings:")
	t.Logf("  Lightness dark max: %.3f", s.LightnessDarkMax)
	t.Logf("  Lightness light min: %.3f", s.LightnessLightMin)
	t.Logf("  Saturation gray max: %.3f", s.SaturationGrayMax)
	t.Logf("  Saturation muted max: %.3f", s.SaturationMutedMax)
	t.Logf("  Saturation normal max: %.3f", s.SaturationNormalMax)

	// Expected values
	expectedValues := map[string]float64{
		"LightnessDarkMax":     0.25,
		"LightnessLightMin":    0.75,
		"SaturationGrayMax":    0.05,
		"SaturationMutedMax":   0.25,
		"SaturationNormalMax":  0.7,
	}

	actualValues := map[string]float64{
		"LightnessDarkMax":     s.LightnessDarkMax,
		"LightnessLightMin":    s.LightnessLightMin,
		"SaturationGrayMax":    s.SaturationGrayMax,
		"SaturationMutedMax":   s.SaturationMutedMax,
		"SaturationNormalMax":  s.SaturationNormalMax,
	}

	for name, expected := range expectedValues {
		actual := actualValues[name]
		if math.Abs(actual-expected) > 0.0001 {
			t.Errorf("%s: expected %.3f, got %.3f", name, expected, actual)
		}
	}
}

func TestDefaultSettings_HueSettings(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Hue organization settings:")
	t.Logf("  Hue sector count: %d", s.HueSectorCount)
	t.Logf("  Hue sector size: %.1f°", s.HueSectorSize)

	expectedSectorCount := 12
	expectedSectorSize := 30.0

	if s.HueSectorCount != expectedSectorCount {
		t.Errorf("Expected hue sector count %d, got %d", expectedSectorCount, s.HueSectorCount)
	}

	if math.Abs(s.HueSectorSize-expectedSectorSize) > 0.0001 {
		t.Errorf("Expected hue sector size %.1f°, got %.1f°", expectedSectorSize, s.HueSectorSize)
	}

	// Verify sector math
	calculatedSize := 360.0 / float64(s.HueSectorCount)
	if math.Abs(s.HueSectorSize-calculatedSize) > 0.0001 {
		t.Errorf("Sector size %.1f° doesn't match calculated %.1f° for %d sectors",
			s.HueSectorSize, calculatedSize, s.HueSectorCount)
	}
}

func TestDefaultSettings_ExtractionSettings(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Extraction settings:")
	t.Logf("  Max colors to extract: %d", s.Extraction.MaxColorsToExtract)
	t.Logf("  Dominant color count: %d", s.Extraction.DominantColorCount)
	t.Logf("  Min color diversity: %.3f", s.Extraction.MinColorDiversity)
	t.Logf("  Adaptive grouping: %t", s.Extraction.AdaptiveGrouping)
	t.Logf("  Preserve natural clusters: %t", s.Extraction.PreserveNaturalClusters)

	expectedValues := map[string]interface{}{
		"MaxColorsToExtract":        100,
		"DominantColorCount":        10,
		"MinColorDiversity":         0.1,
		"AdaptiveGrouping":          true,
		"PreserveNaturalClusters":   true,
	}

	actualValues := map[string]interface{}{
		"MaxColorsToExtract":        s.Extraction.MaxColorsToExtract,
		"DominantColorCount":        s.Extraction.DominantColorCount,
		"MinColorDiversity":         s.Extraction.MinColorDiversity,
		"AdaptiveGrouping":          s.Extraction.AdaptiveGrouping,
		"PreserveNaturalClusters":   s.Extraction.PreserveNaturalClusters,
	}

	for name, expected := range expectedValues {
		actual := actualValues[name]
		switch expected := expected.(type) {
		case float64:
			if actual, ok := actual.(float64); ok {
				if math.Abs(actual-expected) > 0.0001 {
					t.Errorf("%s: expected %.3f, got %.3f", name, expected, actual)
				}
			}
		case int:
			if actual != expected {
				t.Errorf("%s: expected %d, got %v", name, expected, actual)
			}
		case bool:
			if actual != expected {
				t.Errorf("%s: expected %t, got %v", name, expected, actual)
			}
		}
	}
}

func TestDefaultSettings_FallbackColors(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Fallback color settings:")
	t.Logf("  Default dark: %s", s.Fallbacks.DefaultDark)
	t.Logf("  Default light: %s", s.Fallbacks.DefaultLight)
	t.Logf("  Default gray: %s", s.Fallbacks.DefaultGray)

	expectedFallbacks := map[string]string{
		"DefaultDark":  "#1a1a1a",
		"DefaultLight": "#f0f0f0",
		"DefaultGray":  "#808080",
	}

	actualFallbacks := map[string]string{
		"DefaultDark":  s.Fallbacks.DefaultDark,
		"DefaultLight": s.Fallbacks.DefaultLight,
		"DefaultGray":  s.Fallbacks.DefaultGray,
	}

	for name, expected := range expectedFallbacks {
		actual := actualFallbacks[name]
		if actual != expected {
			t.Errorf("%s: expected %s, got %s", name, expected, actual)
		}

		// Verify hex color format
		if len(actual) != 7 || actual[0] != '#' {
			t.Errorf("%s color %s should be in 7-character hex format (#RRGGBB)", name, actual)
		}

		// Verify hex digits
		for i, char := range actual[1:] {
			if !isHexDigit(byte(char)) {
				t.Errorf("%s contains invalid hex character at position %d: '%c'",
					name, i+1, char)
				break
			}
		}
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
grayscale_image_threshold: 0.9
monochromatic_tolerance: 20.0
monochromatic_weight_threshold: 0.2
theme_mode_threshold: 0.6
min_frequency: 0.002
loader_max_width: 4096
loader_max_height: 4096
hue_sector_count: 8
hue_sector_size: 45.0
extraction:
  max_colors_to_extract: 50
  dominant_color_count: 5
fallbacks:
  default_dark: "#101010"
  default_light: "#f8f8f8"
  default_gray: "#888888"
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
	t.Logf("  Hue sector count: %d", loadedSettings.HueSectorCount)
	t.Logf("  Extraction max colors: %d", loadedSettings.Extraction.MaxColorsToExtract)

	// Verify some loaded values
	if math.Abs(loadedSettings.GrayscaleThreshold-0.08) > 0.0001 {
		t.Errorf("Expected grayscale threshold 0.08, got %.3f", loadedSettings.GrayscaleThreshold)
	}

	if math.Abs(loadedSettings.MonochromaticTolerance-20.0) > 0.0001 {
		t.Errorf("Expected monochromatic tolerance 20.0, got %.1f", loadedSettings.MonochromaticTolerance)
	}

	if loadedSettings.HueSectorCount != 8 {
		t.Errorf("Expected hue sector count 8, got %d", loadedSettings.HueSectorCount)
	}

	if loadedSettings.Extraction.MaxColorsToExtract != 50 {
		t.Errorf("Expected max colors to extract 50, got %d", loadedSettings.Extraction.MaxColorsToExtract)
	}

	if loadedSettings.Fallbacks.DefaultDark != "#101010" {
		t.Errorf("Expected fallback dark color #101010, got %s", loadedSettings.Fallbacks.DefaultDark)
	}
}

// Helper function to check if a byte is a valid hex digit
func isHexDigit(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'A' && b <= 'F') || (b >= 'a' && b <= 'f')
}