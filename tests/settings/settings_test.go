package settings_test

import (
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

func TestDefaultSettings(t *testing.T) {
	s := settings.DefaultSettings()

	if s == nil {
		t.Fatal("DefaultSettings() returned nil")
	}

	t.Logf("Testing hierarchical settings structure:")

	// Foundation layer settings
	t.Logf("  Loader.MaxWidth: %d", s.Loader.MaxWidth)
	t.Logf("  Loader.MaxHeight: %d", s.Loader.MaxHeight)
	t.Logf("  Loader.AllowedFormats: %v", s.Loader.AllowedFormats)

	t.Logf("  Formats.QuantizationBits: %d", s.Formats.QuantizationBits)

	t.Logf("  Chromatic.NeutralThreshold: %.3f", s.Chromatic.NeutralThreshold)
	t.Logf("  Chromatic.ColorMergeThreshold: %.1f", s.Chromatic.ColorMergeThreshold)
	t.Logf("  Chromatic.DarkLightnessMax: %.3f", s.Chromatic.DarkLightnessMax)
	t.Logf("  Chromatic.LightLightnessMin: %.3f", s.Chromatic.LightLightnessMin)
	t.Logf("  Chromatic.MutedSaturationMax: %.3f", s.Chromatic.MutedSaturationMax)
	t.Logf("  Chromatic.VibrantSaturationMin: %.3f", s.Chromatic.VibrantSaturationMin)

	// Processing layer settings
	t.Logf("  Processor.MinFrequency: %.6f", s.Processor.MinFrequency)
	t.Logf("  Processor.MinClusterWeight: %.6f", s.Processor.MinClusterWeight)
	t.Logf("  Processor.LightThemeThreshold: %.3f", s.Processor.LightThemeThreshold)
	t.Logf("  Processor.MaxUIColors: %d", s.Processor.MaxUIColors)

	// Global settings
	t.Logf("  DefaultDark: %s", s.DefaultDark)
	t.Logf("  DefaultLight: %s", s.DefaultLight)
	t.Logf("  DefaultGray: %s", s.DefaultGray)
}

func TestCoreSettingsValidation(t *testing.T) {
	s := settings.DefaultSettings()

	// Test core threshold values are sensible
	testCases := []struct {
		name     string
		value    float64
		min      float64
		max      float64
	}{
		{"Chromatic.NeutralThreshold", s.Chromatic.NeutralThreshold, 0.05, 0.2},
		{"Chromatic.ColorMergeThreshold", s.Chromatic.ColorMergeThreshold, 10.0, 25.0},
		{"Chromatic.DarkLightnessMax", s.Chromatic.DarkLightnessMax, 0.2, 0.4},
		{"Chromatic.LightLightnessMin", s.Chromatic.LightLightnessMin, 0.6, 0.8},
		{"Chromatic.MutedSaturationMax", s.Chromatic.MutedSaturationMax, 0.2, 0.4},
		{"Chromatic.VibrantSaturationMin", s.Chromatic.VibrantSaturationMin, 0.6, 0.8},
		{"Processor.MinFrequency", s.Processor.MinFrequency, 0.0001, 0.01},
		{"Processor.LightThemeThreshold", s.Processor.LightThemeThreshold, 0.3, 0.7},
		{"Processor.MinClusterWeight", s.Processor.MinClusterWeight, 0.001, 0.02},
		{"Processor.MinUIColorWeight", s.Processor.MinUIColorWeight, 0.005, 0.05},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Validating %s: %.4f (range: %.4f - %.4f)", 
				tc.name, tc.value, tc.min, tc.max)

			if tc.value < tc.min || tc.value > tc.max {
				t.Errorf("%s value %.4f outside expected range [%.4f, %.4f]",
					tc.name, tc.value, tc.min, tc.max)
			}
		})
	}
}

func TestProcessorSettings(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Testing processor settings:")
	t.Logf("  MinFrequency: %.6f", s.Processor.MinFrequency)
	t.Logf("  MinClusterWeight: %.6f", s.Processor.MinClusterWeight)
	t.Logf("  MinUIColorWeight: %.6f", s.Processor.MinUIColorWeight)
	t.Logf("  MaxUIColors: %d", s.Processor.MaxUIColors)
	t.Logf("  LightThemeThreshold: %.3f", s.Processor.LightThemeThreshold)
	t.Logf("  ThemeModeMaxClusters: %d", s.Processor.ThemeModeMaxClusters)

	// Validate processor constraints
	if s.Processor.MinFrequency <= 0 {
		t.Error("MinFrequency must be positive")
	}

	if s.Processor.MinClusterWeight <= 0 {
		t.Error("MinClusterWeight must be positive")
	}

	if s.Processor.MaxUIColors <= 0 {
		t.Error("MaxUIColors must be positive")
	}

	if s.Processor.LightThemeThreshold < 0 || s.Processor.LightThemeThreshold > 1 {
		t.Errorf("LightThemeThreshold %.3f must be in range [0, 1]", s.Processor.LightThemeThreshold)
	}

	// Weight hierarchy validation
	if s.Processor.MinUIColorWeight <= s.Processor.MinClusterWeight {
		t.Errorf("MinUIColorWeight %.6f should be > MinClusterWeight %.6f",
			s.Processor.MinUIColorWeight, s.Processor.MinClusterWeight)
	}
}

func TestFormatsSettings(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Testing formats settings:")
	t.Logf("  QuantizationBits: %d", s.Formats.QuantizationBits)

	// Validate quantization constraints
	if s.Formats.QuantizationBits < 1 || s.Formats.QuantizationBits > 8 {
		t.Errorf("QuantizationBits %d must be in range [1, 8]", s.Formats.QuantizationBits)
	}

	// Calculate expected quantization levels
	expectedLevels := 1 << s.Formats.QuantizationBits
	t.Logf("  Expected quantization levels per channel: %d", expectedLevels)

	// Validate reasonable quantization (not too coarse, not too fine)
	if expectedLevels < 8 {
		t.Errorf("QuantizationBits %d results in too few levels (%d), may lose color detail",
			s.Formats.QuantizationBits, expectedLevels)
	}

	if expectedLevels > 256 {
		t.Errorf("QuantizationBits %d results in too many levels (%d), defeats quantization purpose",
			s.Formats.QuantizationBits, expectedLevels)
	}
}

func TestLoaderSettings(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Testing loader settings:")
	t.Logf("  MaxWidth: %d", s.Loader.MaxWidth)
	t.Logf("  MaxHeight: %d", s.Loader.MaxHeight)
	t.Logf("  AllowedFormats: %v", s.Loader.AllowedFormats)

	// Validate loader constraints
	if s.Loader.MaxWidth <= 0 {
		t.Error("MaxWidth must be positive")
	}

	if s.Loader.MaxHeight <= 0 {
		t.Error("MaxHeight must be positive")
	}

	if len(s.Loader.AllowedFormats) == 0 {
		t.Error("AllowedFormats cannot be empty")
	}

	// Should at least support common formats
	hasJPEG := false
	hasPNG := false
	for _, format := range s.Loader.AllowedFormats {
		if format == "jpeg" || format == "jpg" {
			hasJPEG = true
		}
		if format == "png" {
			hasPNG = true
		}
	}

	if !hasJPEG {
		t.Error("AllowedFormats should include JPEG support")
	}

	if !hasPNG {
		t.Error("AllowedFormats should include PNG support")
	}

	// Validate reasonable size limits
	if s.Loader.MaxWidth < 1024 || s.Loader.MaxHeight < 1024 {
		t.Errorf("MaxWidth/MaxHeight (%dx%d) should support at least 1024x1024 images",
			s.Loader.MaxWidth, s.Loader.MaxHeight)
	}
}

func TestGlobalSettings(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Testing global fallback settings:")
	t.Logf("  DefaultDark: %s", s.DefaultDark)
	t.Logf("  DefaultLight: %s", s.DefaultLight)
	t.Logf("  DefaultGray: %s", s.DefaultGray)

	// Validate fallback colors are valid hex
	fallbacks := map[string]string{
		"DefaultDark":  s.DefaultDark,
		"DefaultLight": s.DefaultLight,
		"DefaultGray":  s.DefaultGray,
	}

	for name, value := range fallbacks {
		if value == "" {
			t.Errorf("%s fallback cannot be empty", name)
			continue
		}

		if value[0] != '#' {
			t.Errorf("%s fallback '%s' should start with #", name, value)
			continue
		}

		expectedLength := 7 // #RRGGBB
		if len(value) != expectedLength {
			t.Errorf("%s fallback '%s' should be %d characters (including #)",
				name, value, expectedLength)
		}

		// Basic hex validation
		for i, c := range value[1:] {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				t.Errorf("%s fallback '%s' contains invalid hex character '%c' at position %d",
					name, value, c, i+1)
				break
			}
		}
	}
}

func TestSettingsConsistency(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Testing hierarchical settings consistency:")

	// Lightness thresholds should be ordered
	if s.Chromatic.DarkLightnessMax >= s.Chromatic.LightLightnessMin {
		t.Errorf("Lightness thresholds overlap: DarkMax=%.3f >= LightMin=%.3f",
			s.Chromatic.DarkLightnessMax, s.Chromatic.LightLightnessMin)
	}

	// Saturation thresholds should be ordered
	if s.Chromatic.MutedSaturationMax >= s.Chromatic.VibrantSaturationMin {
		t.Errorf("Saturation thresholds overlap: MutedMax=%.3f >= VibrantMin=%.3f",
			s.Chromatic.MutedSaturationMax, s.Chromatic.VibrantSaturationMin)
	}

	// Processor weight hierarchy should be ordered
	if s.Processor.MinFrequency >= s.Processor.MinClusterWeight {
		t.Errorf("MinFrequency %.6f should be < MinClusterWeight %.6f",
			s.Processor.MinFrequency, s.Processor.MinClusterWeight)
	}

	if s.Processor.MinClusterWeight >= s.Processor.MinUIColorWeight {
		t.Errorf("MinClusterWeight %.6f should be < MinUIColorWeight %.6f",
			s.Processor.MinClusterWeight, s.Processor.MinUIColorWeight)
	}

	// Neutral threshold should be reasonable for color similarity
	if s.Chromatic.NeutralThreshold > s.Chromatic.MutedSaturationMax {
		t.Errorf("NeutralThreshold %.3f > MutedSaturationMax %.3f (inconsistent neutral classification)",
			s.Chromatic.NeutralThreshold, s.Chromatic.MutedSaturationMax)
	}

	t.Logf("Hierarchical settings consistency validation passed")
}