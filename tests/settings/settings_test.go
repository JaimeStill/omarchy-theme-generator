package settings_test

import (
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

func TestDefaultSettings(t *testing.T) {
	s := settings.DefaultSettings()
	
	if s == nil {
		t.Fatal("DefaultSettings() should return non-nil settings")
	}
	
	// Test critical thresholds
	if s.GrayscaleThreshold <= 0 || s.GrayscaleThreshold >= 1 {
		t.Errorf("GrayscaleThreshold %v should be between 0 and 1", s.GrayscaleThreshold)
	}
	
	if s.MonochromaticTolerance <= 0 || s.MonochromaticTolerance > 180 {
		t.Errorf("MonochromaticTolerance %v should be between 0 and 180 degrees", s.MonochromaticTolerance)
	}
	
	if s.ThemeModeThreshold <= 0 || s.ThemeModeThreshold >= 1 {
		t.Errorf("ThemeModeThreshold %v should be between 0 and 1", s.ThemeModeThreshold)
	}
	
	// Test frequency settings
	if s.MinFrequency <= 0 || s.MinFrequency > 1 {
		t.Errorf("MinFrequency %v should be between 0 and 1", s.MinFrequency)
	}
	
	// Test contrast requirements
	if s.MinContrastRatio < 3.0 || s.MinContrastRatio > 21.0 {
		t.Errorf("MinContrastRatio %v should be reasonable (3.0-21.0)", s.MinContrastRatio)
	}
	
	// Verify WCAG AA compliance default
	expectedWCAGAA := 4.5
	if s.MinContrastRatio != expectedWCAGAA {
		t.Errorf("Expected WCAG AA contrast ratio %v, got %v", expectedWCAGAA, s.MinContrastRatio)
	}
	
	// Test background thresholds
	if s.LightBackgroundThreshold <= 0.5 || s.LightBackgroundThreshold > 1 {
		t.Errorf("LightBackgroundThreshold %v should be > 0.5 for light colors", s.LightBackgroundThreshold)
	}
	
	if s.DarkBackgroundThreshold <= 0 || s.DarkBackgroundThreshold >= 0.5 {
		t.Errorf("DarkBackgroundThreshold %v should be < 0.5 for dark colors", s.DarkBackgroundThreshold)
	}
	
	// Test saturation thresholds
	if s.MinPrimarySaturation < 0 || s.MinPrimarySaturation > 1 {
		t.Errorf("MinPrimarySaturation %v should be between 0 and 1", s.MinPrimarySaturation)
	}
	
	if s.MinAccentSaturation < 0 || s.MinAccentSaturation > 1 {
		t.Errorf("MinAccentSaturation %v should be between 0 and 1", s.MinAccentSaturation)
	}
	
	// Test lightness bounds
	if s.MinAccentLightness >= s.MaxAccentLightness {
		t.Errorf("MinAccentLightness %v should be < MaxAccentLightness %v", 
			s.MinAccentLightness, s.MaxAccentLightness)
	}
	
	if s.MinAccentLightness < 0 || s.MaxAccentLightness > 1 {
		t.Error("Accent lightness bounds should be within [0, 1]")
	}
}

func TestSettings_ScoringParameters(t *testing.T) {
	s := settings.DefaultSettings()
	
	// Test lightness penalty thresholds
	if s.DarkLightThreshold >= s.BrightLightThreshold {
		t.Errorf("DarkLightThreshold %v should be < BrightLightThreshold %v",
			s.DarkLightThreshold, s.BrightLightThreshold)
	}
	
	// Test penalty/bonus multipliers
	if s.ExtremeLightnessPenalty <= 0 || s.ExtremeLightnessPenalty > 1 {
		t.Errorf("ExtremeLightnessPenalty %v should be between 0 and 1 (penalty)", s.ExtremeLightnessPenalty)
	}
	
	if s.OptimalLightnessBonus <= 1 {
		t.Errorf("OptimalLightnessBonus %v should be > 1 (bonus multiplier)", s.OptimalLightnessBonus)
	}
	
	if s.MinSaturationForBonus < 0 || s.MinSaturationForBonus > 1 {
		t.Errorf("MinSaturationForBonus %v should be between 0 and 1", s.MinSaturationForBonus)
	}
}

func TestSettings_LoaderConfiguration(t *testing.T) {
	s := settings.DefaultSettings()
	
	// Test image size limits
	if s.LoaderMaxWidth <= 0 || s.LoaderMaxHeight <= 0 {
		t.Error("Loader max dimensions should be positive")
	}
	
	if s.LoaderMaxWidth < 1920 || s.LoaderMaxHeight < 1080 {
		t.Error("Loader should support at least Full HD dimensions")
	}
	
	// Test allowed formats
	if len(s.LoaderAllowedFormats) == 0 {
		t.Error("Should have at least one allowed format")
	}
	
	expectedFormats := []string{"jpeg", "jpg", "png", "webp"}
	for _, expected := range expectedFormats {
		found := false
		for _, format := range s.LoaderAllowedFormats {
			if format == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected format %s not found in allowed formats", expected)
		}
	}
}

func TestSettings_FallbackColors(t *testing.T) {
	s := settings.DefaultSettings()
	
	// Test that fallback colors are valid hex strings
	fallbackColors := map[string]string{
		"LightBackgroundFallback": s.LightBackgroundFallback,
		"DarkBackgroundFallback":  s.DarkBackgroundFallback,
		"LightForegroundFallback": s.LightForegroundFallback,
		"DarkForegroundFallback":  s.DarkForegroundFallback,
		"PrimaryFallback":         s.PrimaryFallback,
	}
	
	for name, hexColor := range fallbackColors {
		if hexColor == "" {
			t.Errorf("%s should not be empty", name)
			continue
		}
		
		if hexColor[0] != '#' {
			t.Errorf("%s should start with # (got: %s)", name, hexColor)
			continue
		}
		
		if len(hexColor) != 7 && len(hexColor) != 9 {
			t.Errorf("%s should be 7 or 9 characters (#RRGGBB or #RRGGBBAA) (got: %s)", name, hexColor)
		}
		
		// Test that hex digits are valid
		for i, r := range hexColor[1:] {
			if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
				t.Errorf("%s contains invalid hex character at position %d: %c", name, i+1, r)
			}
		}
	}
}

func TestSettings_ColorTheoryCompliance(t *testing.T) {
	s := settings.DefaultSettings()
	
	// Test that settings align with color theory principles
	
	// Accent colors should be more saturated than primary colors
	if s.MinAccentSaturation <= s.MinPrimarySaturation {
		t.Errorf("MinAccentSaturation %v should be > MinPrimarySaturation %v for visual hierarchy",
			s.MinAccentSaturation, s.MinPrimarySaturation)
	}
	
	// Grayscale threshold should be low (most colors have some saturation)
	if s.GrayscaleThreshold > 0.1 {
		t.Errorf("GrayscaleThreshold %v should be low (< 0.1) to detect truly gray colors", s.GrayscaleThreshold)
	}
	
	// Monochromatic tolerance should be reasonable (not too strict, not too loose)
	if s.MonochromaticTolerance < 10 || s.MonochromaticTolerance > 30 {
		t.Errorf("MonochromaticTolerance %v should be between 10-30 degrees for practical use",
			s.MonochromaticTolerance)
	}
	
	// Theme mode threshold should be around middle luminance
	if s.ThemeModeThreshold < 0.4 || s.ThemeModeThreshold > 0.6 {
		t.Errorf("ThemeModeThreshold %v should be around 0.5 for balanced theme detection",
			s.ThemeModeThreshold)
	}
}

func TestSettings_AccessibilityCompliance(t *testing.T) {
	s := settings.DefaultSettings()
	
	// Test WCAG compliance levels
	testCases := []struct {
		name        string
		ratio       float64
		level       string
		description string
	}{
		{"WCAG AA Normal", 4.5, "AA", "Normal text"},
		{"WCAG AAA Normal", 7.0, "AAA", "Enhanced accessibility"},
		{"WCAG AA Large", 3.0, "AA Large", "Large text"},
	}
	
	// Verify our default meets at least AA
	if s.MinContrastRatio < 4.5 {
		t.Error("Default contrast ratio should meet WCAG AA minimum (4.5:1)")
	}
	
	// Log which level we meet
	for _, tc := range testCases {
		if s.MinContrastRatio >= tc.ratio {
			t.Logf("Settings meet %s standard (%v:1) for %s", tc.level, tc.ratio, tc.description)
		}
	}
	
	// Background thresholds should create accessible combinations
	if s.LightBackgroundThreshold < 0.7 {
		t.Error("LightBackgroundThreshold should be high enough for dark text readability")
	}
	
	if s.DarkBackgroundThreshold > 0.3 {
		t.Error("DarkBackgroundThreshold should be low enough for light text readability")
	}
}