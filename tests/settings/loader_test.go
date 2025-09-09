package settings_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

func TestSettings_Load_NoConfigFile(t *testing.T) {
	// Test loading settings when no config file exists (should use defaults)
	s, err := settings.Load()
	if err != nil {
		t.Fatalf("Load() should not error when no config file exists: %v", err)
	}
	
	if s == nil {
		t.Fatal("Load() should return non-nil settings even without config file")
	}
	
	// Should be equivalent to default settings
	defaults := settings.DefaultSettings()
	
	if s.GrayscaleThreshold != defaults.GrayscaleThreshold {
		t.Errorf("Without config file, should use default GrayscaleThreshold")
	}
	
	if s.MinFrequency != defaults.MinFrequency {
		t.Errorf("Without config file, should use default MinFrequency")
	}

	if s.LoaderMaxWidth != defaults.LoaderMaxWidth {
		t.Errorf("Without config file, should use default LoaderMaxWidth")
	}

	t.Logf("Successfully loaded default settings when no config file exists")
}

func TestSettings_Load_WithConfigFile(t *testing.T) {
	// Create temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-settings.yaml")
	
	configContent := `
grayscale_threshold: 0.08
grayscale_image_threshold: 0.9
monochromatic_tolerance: 25.0
theme_mode_threshold: 0.6
min_frequency: 0.002
loader_max_width: 4096
loader_max_height: 4096
extraction:
  max_colors_to_extract: 75
  dominant_color_count: 8
fallbacks:
  default_dark: "#121212"
  default_light: "#fafafa"
  default_gray: "#888888"
`
	
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}
	
	// Set config file path
	originalConfigPath := os.Getenv("OMARCHY_CONFIG")
	defer func() {
		if originalConfigPath == "" {
			os.Unsetenv("OMARCHY_CONFIG")
		} else {
			os.Setenv("OMARCHY_CONFIG", originalConfigPath)
		}
	}()
	
	os.Setenv("OMARCHY_CONFIG", configPath)
	
	// Test loading from config file
	s, err := settings.Load()
	if err != nil {
		t.Fatalf("Load() should succeed with valid config file: %v", err)
	}
	
	if s == nil {
		t.Fatal("Load() should return non-nil settings with config file")
	}
	
	t.Logf("Loaded settings from config file:")
	t.Logf("  Grayscale threshold: %.3f", s.GrayscaleThreshold)
	t.Logf("  Theme mode threshold: %.3f", s.ThemeModeThreshold)
	t.Logf("  Loader max width: %d", s.LoaderMaxWidth)
	t.Logf("  Extraction max colors: %d", s.Extraction.MaxColorsToExtract)
	
	// Verify values were loaded from config
	if s.GrayscaleThreshold != 0.08 {
		t.Errorf("Expected grayscale threshold 0.08, got %.3f", s.GrayscaleThreshold)
	}
	
	if s.ThemeModeThreshold != 0.6 {
		t.Errorf("Expected theme mode threshold 0.6, got %.3f", s.ThemeModeThreshold)
	}
	
	if s.LoaderMaxWidth != 4096 {
		t.Errorf("Expected loader max width 4096, got %d", s.LoaderMaxWidth)
	}

	if s.Extraction.MaxColorsToExtract != 75 {
		t.Errorf("Expected extraction max colors 75, got %d", s.Extraction.MaxColorsToExtract)
	}

	if s.Fallbacks.DefaultDark != "#121212" {
		t.Errorf("Expected fallback dark color #121212, got %s", s.Fallbacks.DefaultDark)
	}
}

func TestSettings_Load_InvalidConfigFile(t *testing.T) {
	// Create invalid config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid-settings.yaml")
	
	invalidContent := `
invalid_yaml: [
  unclosed_bracket: true
  missing_bracket
`
	
	err := os.WriteFile(configPath, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}
	
	// Set config file path
	originalConfigPath := os.Getenv("OMARCHY_CONFIG")
	defer func() {
		if originalConfigPath == "" {
			os.Unsetenv("OMARCHY_CONFIG")
		} else {
			os.Setenv("OMARCHY_CONFIG", originalConfigPath)
		}
	}()
	
	os.Setenv("OMARCHY_CONFIG", configPath)
	
	// Test loading invalid config should return error when explicit path is set
	s, err := settings.Load()
	if err == nil {
		t.Error("Load() should return error for invalid explicit config file")
	} else {
		t.Logf("Load() correctly returned error for invalid config: %v", err)
	}
	
	if s != nil {
		t.Error("Load() should return nil settings when explicit config file is invalid")
	}
	
	t.Logf("Successfully handled invalid config file")
}

func TestSettings_ViperIntegration(t *testing.T) {
	// Test direct Viper integration
	v := viper.New()
	v.SetDefault("grayscale_threshold", 0.1)
	v.SetDefault("min_frequency", 0.005)
	v.SetDefault("loader_max_width", 2048)
	
	var testSettings settings.Settings
	err := v.Unmarshal(&testSettings)
	if err != nil {
		t.Fatalf("Failed to unmarshal with Viper: %v", err)
	}
	
	t.Logf("Viper integration test:")
	t.Logf("  Grayscale threshold: %.3f", testSettings.GrayscaleThreshold)
	t.Logf("  Min frequency: %.6f", testSettings.MinFrequency)
	t.Logf("  Loader max width: %d", testSettings.LoaderMaxWidth)
	
	if testSettings.GrayscaleThreshold != 0.1 {
		t.Errorf("Expected grayscale threshold 0.1, got %.3f", testSettings.GrayscaleThreshold)
	}
	
	if testSettings.MinFrequency != 0.005 {
		t.Errorf("Expected min frequency 0.005, got %.6f", testSettings.MinFrequency)
	}
	
	if testSettings.LoaderMaxWidth != 2048 {
		t.Errorf("Expected loader max width 2048, got %d", testSettings.LoaderMaxWidth)
	}
}

func TestSettings_ContextIntegration(t *testing.T) {
	// Test loading and context integration
	s, err := settings.Load()
	if err != nil {
		t.Fatalf("Failed to load settings: %v", err)
	}
	
	ctx := context.Background()
	ctxWithSettings := settings.WithSettings(ctx, s)
	
	retrievedSettings := settings.FromContext(ctxWithSettings)
	if retrievedSettings == nil {
		t.Fatal("Settings should be retrievable from context")
	}
	
	// Test that settings are the same instance
	if retrievedSettings != s {
		t.Error("Retrieved settings should be the same instance as stored")
	}
	
	t.Logf("Context integration successful")
	t.Logf("  Original settings grayscale threshold: %.3f", s.GrayscaleThreshold)
	t.Logf("  Retrieved settings grayscale threshold: %.3f", retrievedSettings.GrayscaleThreshold)
}

func TestSettings_ConfigPaths(t *testing.T) {
	// Test various config path scenarios
	testCases := []struct {
		name        string
		envValue    string
		shouldExist bool
		description string
	}{
		{
			name:        "Empty environment",
			envValue:    "",
			shouldExist: false,
			description: "No OMARCHY_CONFIG set - should use defaults",
		},
		{
			name:        "Non-existent file",
			envValue:    "/non/existent/path/config.yaml", 
			shouldExist: false,
			description: "Non-existent config file - should return error",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment
			originalConfig := os.Getenv("OMARCHY_CONFIG")
			defer func() {
				if originalConfig == "" {
					os.Unsetenv("OMARCHY_CONFIG")
				} else {
					os.Setenv("OMARCHY_CONFIG", originalConfig)
				}
			}()
			
			if tc.envValue == "" {
				os.Unsetenv("OMARCHY_CONFIG")
			} else {
				os.Setenv("OMARCHY_CONFIG", tc.envValue)
			}
			
			t.Logf("Testing %s: %s", tc.name, tc.description)
			
			s, err := settings.Load()
			
			// Handle different expectations based on the test case
			if tc.name == "Non-existent file" {
				// Explicit non-existent config file should return error and nil settings
				if err == nil {
					t.Error("Expected error for non-existent explicit config file")
				}
				if s != nil {
					t.Error("Expected nil settings for non-existent explicit config file")
				}
				t.Logf("Correctly handled non-existent explicit config file: %v", err)
			} else {
				// Other cases should succeed
				if err != nil && tc.shouldExist {
					t.Errorf("Expected successful load but got error: %v", err)
				}
				
				if s == nil {
					t.Fatal("Settings should never be nil for default config behavior")
				}
				
				t.Logf("Successfully loaded settings for case: %s", tc.name)
			}
		})
	}
}

func TestSettings_StructTags(t *testing.T) {
	// Test that mapstructure tags work correctly
	configData := map[string]interface{}{
		"grayscale_threshold":              0.07,
		"grayscale_image_threshold":        0.85,
		"monochromatic_tolerance":          18.0,
		"monochromatic_weight_threshold":   0.15,
		"theme_mode_threshold":             0.55,
		"min_frequency":                    0.0005,
		"loader_max_width":                 6144,
		"loader_max_height":                6144,
		"lightness_dark_max":               0.3,
		"lightness_light_min":              0.8,
		"saturation_gray_max":              0.04,
		"hue_sector_count":                 16,
		"hue_sector_size":                  22.5,
		"extraction": map[string]interface{}{
			"max_colors_to_extract":      80,
			"dominant_color_count":       12,
			"min_color_diversity":        0.4,
			"adaptive_grouping":          false,
			"preserve_natural_clusters":  false,
		},
		"fallbacks": map[string]interface{}{
			"default_dark":  "#181818",
			"default_light": "#f0f0f0",
			"default_gray":  "#909090",
		},
	}
	
	v := viper.New()
	for key, value := range configData {
		v.Set(key, value)
	}
	
	var s settings.Settings
	err := v.Unmarshal(&s)
	if err != nil {
		t.Fatalf("Failed to unmarshal settings: %v", err)
	}
	
	t.Logf("Struct tag validation:")
	t.Logf("  Grayscale threshold: %.3f (expected 0.07)", s.GrayscaleThreshold)
	t.Logf("  Hue sector count: %d (expected 16)", s.HueSectorCount)
	t.Logf("  Extraction max colors: %d (expected 80)", s.Extraction.MaxColorsToExtract)
	t.Logf("  Fallback dark: %s (expected #181818)", s.Fallbacks.DefaultDark)
	
	// Verify unmarshal worked correctly
	if s.GrayscaleThreshold != 0.07 {
		t.Errorf("Expected grayscale threshold 0.07, got %.3f", s.GrayscaleThreshold)
	}
	
	if s.HueSectorCount != 16 {
		t.Errorf("Expected hue sector count 16, got %d", s.HueSectorCount)
	}
	
	if s.Extraction.MaxColorsToExtract != 80 {
		t.Errorf("Expected extraction max colors 80, got %d", s.Extraction.MaxColorsToExtract)
	}
	
	if s.Fallbacks.DefaultDark != "#181818" {
		t.Errorf("Expected fallback dark #181818, got %s", s.Fallbacks.DefaultDark)
	}
	
	if s.Extraction.AdaptiveGrouping {
		t.Errorf("Expected adaptive grouping false, got %t", s.Extraction.AdaptiveGrouping)
	}
}