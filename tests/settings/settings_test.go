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

	t.Logf("Testing default settings structure:")
	
	// Core extraction settings
	t.Logf("  GrayscaleThreshold: %.3f", s.GrayscaleThreshold)
	t.Logf("  GrayscaleImageThreshold: %.3f", s.GrayscaleImageThreshold)
	t.Logf("  MonochromaticTolerance: %.1f°", s.MonochromaticTolerance)
	t.Logf("  MonochromaticWeightThreshold: %.3f", s.MonochromaticWeightThreshold)
	t.Logf("  ThemeModeThreshold: %.3f", s.ThemeModeThreshold)
	t.Logf("  MinFrequency: %.4f", s.MinFrequency)

	// Loader settings
	t.Logf("  LoaderMaxWidth: %d", s.LoaderMaxWidth)
	t.Logf("  LoaderMaxHeight: %d", s.LoaderMaxHeight)
	t.Logf("  LoaderAllowedFormats: %v", s.LoaderAllowedFormats)

	// Grouping thresholds
	t.Logf("  LightnessDarkMax: %.3f", s.LightnessDarkMax)
	t.Logf("  LightnessLightMin: %.3f", s.LightnessLightMin)
	t.Logf("  SaturationGrayMax: %.3f", s.SaturationGrayMax)
	t.Logf("  SaturationMutedMax: %.3f", s.SaturationMutedMax)
	t.Logf("  SaturationNormalMax: %.3f", s.SaturationNormalMax)

	// Hue organization
	t.Logf("  HueSectorCount: %d", s.HueSectorCount)
	t.Logf("  HueSectorSize: %.1f°", s.HueSectorSize)

	// Extraction settings
	t.Logf("  Extraction.MaxColorsToExtract: %d", s.Extraction.MaxColorsToExtract)
	t.Logf("  Extraction.DominantColorCount: %d", s.Extraction.DominantColorCount)
	t.Logf("  Extraction.MinColorDiversity: %.3f", s.Extraction.MinColorDiversity)
	t.Logf("  Extraction.AdaptiveGrouping: %t", s.Extraction.AdaptiveGrouping)
	t.Logf("  Extraction.PreserveNaturalClusters: %t", s.Extraction.PreserveNaturalClusters)

	// Fallback settings
	t.Logf("  Fallbacks.DefaultDark: %s", s.Fallbacks.DefaultDark)
	t.Logf("  Fallbacks.DefaultLight: %s", s.Fallbacks.DefaultLight)
	t.Logf("  Fallbacks.DefaultGray: %s", s.Fallbacks.DefaultGray)
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
		{"GrayscaleThreshold", s.GrayscaleThreshold, 0.0, 0.2},
		{"GrayscaleImageThreshold", s.GrayscaleImageThreshold, 0.5, 1.0},
		{"MonochromaticTolerance", s.MonochromaticTolerance, 5.0, 30.0},
		{"MonochromaticWeightThreshold", s.MonochromaticWeightThreshold, 0.1, 0.9},
		{"ThemeModeThreshold", s.ThemeModeThreshold, 0.3, 0.7},
		{"MinFrequency", s.MinFrequency, 0.0001, 0.01},
		{"LightnessDarkMax", s.LightnessDarkMax, 0.2, 0.4},
		{"LightnessLightMin", s.LightnessLightMin, 0.6, 0.8},
		{"SaturationGrayMax", s.SaturationGrayMax, 0.03, 0.08},
		{"SaturationMutedMax", s.SaturationMutedMax, 0.2, 0.4},
		{"SaturationNormalMax", s.SaturationNormalMax, 0.6, 0.8},
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

func TestHueSettings(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Testing hue sector configuration:")
	t.Logf("  HueSectorCount: %d", s.HueSectorCount)
	t.Logf("  HueSectorSize: %.1f°", s.HueSectorSize)

	// Validate hue sector math
	expectedSectorSize := 360.0 / float64(s.HueSectorCount)
	if s.HueSectorSize != expectedSectorSize {
		t.Errorf("HueSectorSize %.1f° doesn't match calculated value %.1f° for %d sectors",
			s.HueSectorSize, expectedSectorSize, s.HueSectorCount)
	}

	// Validate reasonable sector counts
	if s.HueSectorCount < 6 || s.HueSectorCount > 24 {
		t.Errorf("HueSectorCount %d outside reasonable range [6, 24]", s.HueSectorCount)
	}
}

func TestExtractionSettings(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Testing extraction settings:")
	t.Logf("  MaxColorsToExtract: %d", s.Extraction.MaxColorsToExtract)
	t.Logf("  DominantColorCount: %d", s.Extraction.DominantColorCount)
	t.Logf("  MinColorDiversity: %.3f", s.Extraction.MinColorDiversity)

	// Validate extraction limits
	if s.Extraction.MaxColorsToExtract <= 0 {
		t.Error("MaxColorsToExtract must be positive")
	}

	if s.Extraction.DominantColorCount <= 0 {
		t.Error("DominantColorCount must be positive")
	}

	if s.Extraction.DominantColorCount > s.Extraction.MaxColorsToExtract {
		t.Error("DominantColorCount cannot exceed MaxColorsToExtract")
	}

	if s.Extraction.MinColorDiversity < 0 || s.Extraction.MinColorDiversity > 1 {
		t.Errorf("MinColorDiversity %.3f must be in range [0, 1]", s.Extraction.MinColorDiversity)
	}

	t.Logf("  AdaptiveGrouping: %t", s.Extraction.AdaptiveGrouping)
	t.Logf("  PreserveNaturalClusters: %t", s.Extraction.PreserveNaturalClusters)
}

func TestLoaderSettings(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Testing loader settings:")
	t.Logf("  MaxWidth: %d", s.LoaderMaxWidth)
	t.Logf("  MaxHeight: %d", s.LoaderMaxHeight)
	t.Logf("  AllowedFormats: %v", s.LoaderAllowedFormats)

	// Validate loader constraints
	if s.LoaderMaxWidth <= 0 {
		t.Error("LoaderMaxWidth must be positive")
	}

	if s.LoaderMaxHeight <= 0 {
		t.Error("LoaderMaxHeight must be positive")
	}

	if len(s.LoaderAllowedFormats) == 0 {
		t.Error("LoaderAllowedFormats cannot be empty")
	}

	// Should at least support common formats
	hasJPEG := false
	hasPNG := false
	for _, format := range s.LoaderAllowedFormats {
		if format == "jpeg" || format == "jpg" {
			hasJPEG = true
		}
		if format == "png" {
			hasPNG = true
		}
	}

	if !hasJPEG {
		t.Error("LoaderAllowedFormats should include JPEG support")
	}

	if !hasPNG {
		t.Error("LoaderAllowedFormats should include PNG support")
	}
}

func TestFallbackSettings(t *testing.T) {
	s := settings.DefaultSettings()

	t.Logf("Testing fallback settings:")
	t.Logf("  DefaultDark: %s", s.Fallbacks.DefaultDark)
	t.Logf("  DefaultLight: %s", s.Fallbacks.DefaultLight)
	t.Logf("  DefaultGray: %s", s.Fallbacks.DefaultGray)

	// Validate fallback colors are valid hex
	fallbacks := map[string]string{
		"DefaultDark":  s.Fallbacks.DefaultDark,
		"DefaultLight": s.Fallbacks.DefaultLight,
		"DefaultGray":  s.Fallbacks.DefaultGray,
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

	t.Logf("Testing settings consistency:")

	// Lightness thresholds should be ordered
	if s.LightnessDarkMax >= s.LightnessLightMin {
		t.Errorf("Lightness thresholds overlap: DarkMax=%.3f >= LightMin=%.3f",
			s.LightnessDarkMax, s.LightnessLightMin)
	}

	// Saturation thresholds should be ordered
	thresholds := []float64{s.SaturationGrayMax, s.SaturationMutedMax, s.SaturationNormalMax}
	for i := 0; i < len(thresholds)-1; i++ {
		if thresholds[i] >= thresholds[i+1] {
			t.Errorf("Saturation thresholds not ordered: %.3f >= %.3f at positions %d, %d",
				thresholds[i], thresholds[i+1], i, i+1)
		}
	}

	// Grayscale thresholds should be consistent
	if s.GrayscaleThreshold > s.SaturationGrayMax {
		t.Errorf("GrayscaleThreshold %.3f > SaturationGrayMax %.3f (inconsistent)",
			s.GrayscaleThreshold, s.SaturationGrayMax)
	}

	t.Logf("Settings consistency validation passed")
}