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
			name:     "Chromatic neutral threshold",
			getValue: func(s *settings.Settings) interface{} { return s.Chromatic.NeutralThreshold },
			expected: 0.1,
		},
		{
			name:     "Chromatic color merge threshold",
			getValue: func(s *settings.Settings) interface{} { return s.Chromatic.ColorMergeThreshold },
			expected: 15.0,
		},
		{
			name:     "Processor min cluster weight",
			getValue: func(s *settings.Settings) interface{} { return s.Processor.MinClusterWeight },
			expected: 0.005,
		},
		{
			name:     "Processor light theme threshold",
			getValue: func(s *settings.Settings) interface{} { return s.Processor.LightThemeThreshold },
			expected: 0.5,
		},
		{
			name:     "Processor min frequency",
			getValue: func(s *settings.Settings) interface{} { return s.Processor.MinFrequency },
			expected: 0.0001,
		},
		{
			name:     "Formats quantization bits",
			getValue: func(s *settings.Settings) interface{} { return s.Formats.QuantizationBits },
			expected: 5,
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
	t.Logf("  Max width: %d pixels", s.Loader.MaxWidth)
	t.Logf("  Max height: %d pixels", s.Loader.MaxHeight)
	t.Logf("  Allowed formats: %v", s.Loader.AllowedFormats)

	// Test dimensions
	expectedMaxWidth := 8192
	expectedMaxHeight := 8192

	if s.Loader.MaxWidth != expectedMaxWidth {
		t.Errorf("Expected max width %d, got %d", expectedMaxWidth, s.Loader.MaxWidth)
	}
	if s.Loader.MaxHeight != expectedMaxHeight {
		t.Errorf("Expected max height %d, got %d", expectedMaxHeight, s.Loader.MaxHeight)
	}

	// Test allowed formats
	expectedFormats := []string{"jpeg", "jpg", "png", "webp"}
	if len(s.Loader.AllowedFormats) != len(expectedFormats) {
		t.Errorf("Expected %d formats, got %d", len(expectedFormats), len(s.Loader.AllowedFormats))
	}

	formatMap := make(map[string]bool)
	for _, format := range s.Loader.AllowedFormats {
		formatMap[format] = true
	}

	for _, expected := range expectedFormats {
		if !formatMap[expected] {
			t.Errorf("Expected format '%s' not found in allowed formats", expected)
		}
	}
}

func TestDefaultSettings_ChromaticThresholds(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Chromatic threshold settings:")
	t.Logf("  Dark lightness max: %.3f", s.Chromatic.DarkLightnessMax)
	t.Logf("  Light lightness min: %.3f", s.Chromatic.LightLightnessMin)
	t.Logf("  Neutral threshold: %.3f", s.Chromatic.NeutralThreshold)
	t.Logf("  Muted saturation max: %.3f", s.Chromatic.MutedSaturationMax)
	t.Logf("  Vibrant saturation min: %.3f", s.Chromatic.VibrantSaturationMin)
	t.Logf("  Color merge threshold: %.1f", s.Chromatic.ColorMergeThreshold)

	// Expected values based on defaults.go
	expectedValues := map[string]float64{
		"DarkLightnessMax":      0.3,
		"LightLightnessMin":     0.7,
		"NeutralThreshold":      0.1,
		"MutedSaturationMax":    0.3,
		"VibrantSaturationMin":  0.7,
		"ColorMergeThreshold":   15.0,
		"NeutralLightnessThreshold": 0.08,
	}

	actualValues := map[string]float64{
		"DarkLightnessMax":      s.Chromatic.DarkLightnessMax,
		"LightLightnessMin":     s.Chromatic.LightLightnessMin,
		"NeutralThreshold":      s.Chromatic.NeutralThreshold,
		"MutedSaturationMax":    s.Chromatic.MutedSaturationMax,
		"VibrantSaturationMin":  s.Chromatic.VibrantSaturationMin,
		"ColorMergeThreshold":   s.Chromatic.ColorMergeThreshold,
		"NeutralLightnessThreshold": s.Chromatic.NeutralLightnessThreshold,
	}

	for name, expected := range expectedValues {
		actual := actualValues[name]
		if math.Abs(actual-expected) > 0.0001 {
			t.Errorf("%s: expected %.3f, got %.3f", name, expected, actual)
		}
	}
}

func TestDefaultSettings_ProcessorSettings(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Processor settings:")
	t.Logf("  Min frequency: %.6f", s.Processor.MinFrequency)
	t.Logf("  Min cluster weight: %.6f", s.Processor.MinClusterWeight)
	t.Logf("  Min UI color weight: %.6f", s.Processor.MinUIColorWeight)
	t.Logf("  Max UI colors: %d", s.Processor.MaxUIColors)
	t.Logf("  Light theme threshold: %.3f", s.Processor.LightThemeThreshold)
	t.Logf("  Theme mode max clusters: %d", s.Processor.ThemeModeMaxClusters)

	// Expected values based on defaults.go
	expectedValues := map[string]interface{}{
		"MinFrequency":              0.0001,
		"MinClusterWeight":          0.005,
		"MinUIColorWeight":          0.01,
		"MaxUIColors":               20,
		"LightThemeThreshold":       0.5,
		"ThemeModeMaxClusters":      5,
		"SignificantColorThreshold": 0.1,
		"PureBlackThreshold":        0.01,
		"PureWhiteThreshold":        0.99,
	}

	actualValues := map[string]interface{}{
		"MinFrequency":              s.Processor.MinFrequency,
		"MinClusterWeight":          s.Processor.MinClusterWeight,
		"MinUIColorWeight":          s.Processor.MinUIColorWeight,
		"MaxUIColors":               s.Processor.MaxUIColors,
		"LightThemeThreshold":       s.Processor.LightThemeThreshold,
		"ThemeModeMaxClusters":      s.Processor.ThemeModeMaxClusters,
		"SignificantColorThreshold": s.Processor.SignificantColorThreshold,
		"PureBlackThreshold":        s.Processor.PureBlackThreshold,
		"PureWhiteThreshold":        s.Processor.PureWhiteThreshold,
	}

	for name, expected := range expectedValues {
		actual := actualValues[name]
		switch expected := expected.(type) {
		case float64:
			if actual, ok := actual.(float64); ok {
				if math.Abs(actual-expected) > 0.0001 {
					t.Errorf("%s: expected %.6f, got %.6f", name, expected, actual)
				}
			}
		case int:
			if actual != expected {
				t.Errorf("%s: expected %d, got %v", name, expected, actual)
			}
		}
	}
}

func TestDefaultSettings_GlobalSettings(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Global fallback color settings:")
	t.Logf("  Default dark: %s", s.DefaultDark)
	t.Logf("  Default light: %s", s.DefaultLight)
	t.Logf("  Default gray: %s", s.DefaultGray)

	expectedValues := map[string]string{
		"DefaultDark":  "#1a1a1a",
		"DefaultLight": "#f0f0f0",
		"DefaultGray":  "#808080",
	}

	actualValues := map[string]string{
		"DefaultDark":  s.DefaultDark,
		"DefaultLight": s.DefaultLight,
		"DefaultGray":  s.DefaultGray,
	}

	for name, expected := range expectedValues {
		actual := actualValues[name]
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

func TestDefaultSettings_FormatsSettings(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Formats settings:")
	t.Logf("  Quantization bits: %d", s.Formats.QuantizationBits)

	expectedBits := 5
	if s.Formats.QuantizationBits != expectedBits {
		t.Errorf("Expected quantization bits %d, got %d", expectedBits, s.Formats.QuantizationBits)
	}

	// Verify quantization bits are within valid range (1-8)
	if s.Formats.QuantizationBits < 1 || s.Formats.QuantizationBits > 8 {
		t.Errorf("Quantization bits %d should be in range [1-8]", s.Formats.QuantizationBits)
	}

	// Calculate expected quantization levels
	expectedLevels := 1 << s.Formats.QuantizationBits
	t.Logf("  Quantization levels per channel: %d", expectedLevels)

	if expectedLevels != 32 {
		t.Errorf("Expected 32 quantization levels for 5 bits, got %d", expectedLevels)
	}
}

func TestSettings_WithContext(t *testing.T) {
	s := settings.DefaultSettings()
	ctx := context.Background()

	// Test setting context
	ctxWithSettings := settings.WithSettings(ctx, s)
	
	// Test retrieving from context
	retrievedSettings := settings.FromContext(ctxWithSettings)
	
	t.Logf("Original settings neutral threshold: %.3f", s.Chromatic.NeutralThreshold)
	t.Logf("Retrieved settings neutral threshold: %.3f", retrievedSettings.Chromatic.NeutralThreshold)

	if retrievedSettings == nil {
		t.Error("Retrieved settings should not be nil")
	}

	if retrievedSettings.Chromatic.NeutralThreshold != s.Chromatic.NeutralThreshold {
		t.Errorf("Expected neutral threshold %.3f, got %.3f",
			s.Chromatic.NeutralThreshold, retrievedSettings.Chromatic.NeutralThreshold)
	}

	// Test context without settings - should return default
	emptyCtx := context.Background()
	defaultSettings := settings.FromContext(emptyCtx)
	
	if defaultSettings == nil {
		t.Error("Default settings should not be nil")
	}

	t.Logf("Default settings from empty context: neutral threshold %.3f",
		defaultSettings.Chromatic.NeutralThreshold)
}

func TestSettings_LoadFromFile(t *testing.T) {
	// Create temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")
	
	configContent := `
chromatic:
  neutral_threshold: 0.08
  color_merge_threshold: 20.0
  dark_lightness_max: 0.35
  light_lightness_min: 0.75
processor:
  min_frequency: 0.002
  light_theme_threshold: 0.6
  min_cluster_weight: 0.01
loader:
  max_width: 4096
  max_height: 4096
  allowed_formats: ["jpeg", "png"]
formats:
  quantization_bits: 6
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
	t.Logf("  Neutral threshold: %.3f", loadedSettings.Chromatic.NeutralThreshold)
	t.Logf("  Color merge threshold: %.1f", loadedSettings.Chromatic.ColorMergeThreshold)
	t.Logf("  Quantization bits: %d", loadedSettings.Formats.QuantizationBits)
	t.Logf("  Min frequency: %.6f", loadedSettings.Processor.MinFrequency)
	t.Logf("  Loader max width: %d", loadedSettings.Loader.MaxWidth)

	// Verify some loaded values
	if math.Abs(loadedSettings.Chromatic.NeutralThreshold-0.08) > 0.0001 {
		t.Errorf("Expected neutral threshold 0.08, got %.3f", loadedSettings.Chromatic.NeutralThreshold)
	}

	if math.Abs(loadedSettings.Chromatic.ColorMergeThreshold-20.0) > 0.0001 {
		t.Errorf("Expected color merge threshold 20.0, got %.1f", loadedSettings.Chromatic.ColorMergeThreshold)
	}

	if loadedSettings.Formats.QuantizationBits != 6 {
		t.Errorf("Expected quantization bits 6, got %d", loadedSettings.Formats.QuantizationBits)
	}

	if math.Abs(loadedSettings.Processor.MinFrequency-0.002) > 0.0001 {
		t.Errorf("Expected min frequency 0.002, got %.6f", loadedSettings.Processor.MinFrequency)
	}

	if loadedSettings.DefaultDark != "#101010" {
		t.Errorf("Expected default dark color #101010, got %s", loadedSettings.DefaultDark)
	}
}

// Helper function to check if a byte is a valid hex digit
func isHexDigit(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'A' && b <= 'F') || (b >= 'a' && b <= 'f')
}