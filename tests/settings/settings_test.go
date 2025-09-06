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
	
	// Test category scoring weights
	weights := s.CategoryScoring
	if weights.Frequency <= 0 || weights.Frequency > 1 {
		t.Errorf("CategoryScoring.Frequency %v should be between 0 and 1", weights.Frequency)
	}
	if weights.Contrast <= 0 || weights.Contrast > 1 {
		t.Errorf("CategoryScoring.Contrast %v should be between 0 and 1", weights.Contrast)
	}
	if weights.Saturation <= 0 || weights.Saturation > 1 {
		t.Errorf("CategoryScoring.Saturation %v should be between 0 and 1", weights.Saturation)
	}
	if weights.HueAlignment < 0 || weights.HueAlignment > 1 {
		t.Errorf("CategoryScoring.HueAlignment %v should be between 0 and 1", weights.HueAlignment)
	}
	if weights.Lightness < 0 || weights.Lightness > 1 {
		t.Errorf("CategoryScoring.Lightness %v should be between 0 and 1", weights.Lightness)
	}
	
	// Verify scoring weights sum to a reasonable value (should be around 1.0)
	totalWeight := weights.Frequency + weights.Contrast + weights.Saturation + weights.HueAlignment + weights.Lightness
	if totalWeight < 0.8 || totalWeight > 1.2 {
		t.Errorf("Total category scoring weights %v should be near 1.0 for balanced scoring", totalWeight)
	}
	t.Logf("Category scoring weights sum: %.3f", totalWeight)
	
	// Test extraction settings
	if s.Extraction.MaxCandidatesPerCategory <= 0 {
		t.Error("MaxCandidatesPerCategory should be positive")
	}
	if s.Extraction.MinimumColorFrequency < 0 || s.Extraction.MinimumColorFrequency > 1 {
		t.Errorf("MinimumColorFrequency %v should be between 0 and 1", s.Extraction.MinimumColorFrequency)
	}
	if s.Extraction.GrayscaleHueTemperature < 0 || s.Extraction.GrayscaleHueTemperature >= 360 {
		t.Errorf("GrayscaleHueTemperature %v should be between 0 and 360 degrees", s.Extraction.GrayscaleHueTemperature)
	}
}

func TestSettings_CategoryCharacteristics(t *testing.T) {
	s := settings.DefaultSettings()
	
	// Test that we have category definitions for both modes
	if len(s.Categories.Light) == 0 {
		t.Error("Should have light mode category characteristics defined")
	}
	if len(s.Categories.Dark) == 0 {
		t.Error("Should have dark mode category characteristics defined")
	}
	
	// Test category characteristic validity for light mode
	for category, chars := range s.Categories.Light {
		if chars.MinLightness < 0 || chars.MinLightness > 1 {
			t.Errorf("Category %s light mode: MinLightness %v should be between 0 and 1", category, chars.MinLightness)
		}
		if chars.MaxLightness < 0 || chars.MaxLightness > 1 {
			t.Errorf("Category %s light mode: MaxLightness %v should be between 0 and 1", category, chars.MaxLightness)
		}
		if chars.MinLightness >= chars.MaxLightness {
			t.Errorf("Category %s light mode: MinLightness %v should be < MaxLightness %v", 
				category, chars.MinLightness, chars.MaxLightness)
		}
		if chars.MinSaturation < 0 || chars.MinSaturation > 1 {
			t.Errorf("Category %s light mode: MinSaturation %v should be between 0 and 1", category, chars.MinSaturation)
		}
		if chars.MaxSaturation < 0 || chars.MaxSaturation > 1 {
			t.Errorf("Category %s light mode: MaxSaturation %v should be between 0 and 1", category, chars.MaxSaturation)
		}
		if chars.MinContrast < 0.0 || chars.MinContrast > 21.0 {
			t.Errorf("Category %s light mode: MinContrast %v should be between 0.0 and 21.0", category, chars.MinContrast)
		}
		// Background categories can have MinContrast of 0 since they don't need contrast with themselves
		if category == "background" && chars.MinContrast != 0.0 {
			t.Logf("Category %s light mode: MinContrast %v (background can be 0)", category, chars.MinContrast)
		}
	}
	
	// Test category characteristic validity for dark mode
	for category, chars := range s.Categories.Dark {
		if chars.MinLightness < 0 || chars.MinLightness > 1 {
			t.Errorf("Category %s dark mode: MinLightness %v should be between 0 and 1", category, chars.MinLightness)
		}
		if chars.MaxLightness < 0 || chars.MaxLightness > 1 {
			t.Errorf("Category %s dark mode: MaxLightness %v should be between 0 and 1", category, chars.MaxLightness)
		}
		if chars.MinLightness >= chars.MaxLightness {
			t.Errorf("Category %s dark mode: MinLightness %v should be < MaxLightness %v", 
				category, chars.MinLightness, chars.MaxLightness)
		}
	}
	t.Logf("Found %d light mode categories and %d dark mode categories", 
		len(s.Categories.Light), len(s.Categories.Dark))
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
	
	// Check category-based saturation expectations
	// Accent categories should have higher minimum saturation than background categories
	if accentChars, hasAccent := s.Categories.Light["accent_primary"]; hasAccent {
		if bgChars, hasBg := s.Categories.Light["background"]; hasBg {
			if accentChars.MinSaturation <= bgChars.MinSaturation {
				t.Logf("Note: Accent primary min saturation (%.3f) should typically be higher than background (%.3f)", 
					accentChars.MinSaturation, bgChars.MinSaturation)
			}
		}
	}
	
	// Extraction settings should support vibrant accents
	if !s.Extraction.PreferVibrantAccents {
		t.Log("Note: PreferVibrantAccents is disabled - may result in muted accent colors")
	}
}

func TestSettings_AccessibilityCompliance(t *testing.T) {
	s := settings.DefaultSettings()
	
	// Test WCAG compliance levels in category characteristics
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
	
	// Check foreground category contrast requirements
	for mode, categories := range map[string]map[string]settings.CategoryCharacteristics{
		"light": s.Categories.Light,
		"dark":  s.Categories.Dark,
	} {
		if fgChars, hasFg := categories["foreground"]; hasFg {
			if fgChars.MinContrast < 4.5 {
				t.Errorf("%s mode foreground contrast %v should meet WCAG AA minimum (4.5:1)", 
					mode, fgChars.MinContrast)
			} else {
				// Log which level we meet
				for _, tc := range testCases {
					if fgChars.MinContrast >= tc.ratio {
						t.Logf("%s mode foreground meets %s standard (%.1f:1) for %s", 
							mode, tc.level, tc.ratio, tc.description)
						break
					}
				}
			}
		}
		
		// Check that background categories have appropriate lightness ranges
		if bgChars, hasBg := categories["background"]; hasBg {
			if mode == "light" && bgChars.MinLightness < 0.7 {
				t.Errorf("Light mode background min lightness %v should be high (>= 0.7) for readability", 
					bgChars.MinLightness)
			}
			if mode == "dark" && bgChars.MaxLightness > 0.3 {
				t.Errorf("Dark mode background max lightness %v should be low (<= 0.3) for readability", 
					bgChars.MaxLightness)
			}
		}
	}
}