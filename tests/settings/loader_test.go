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
	
	if s.Chromatic.NeutralThreshold != defaults.Chromatic.NeutralThreshold {
		t.Errorf("Without config file, should use default NeutralThreshold")
	}

	if s.Processor.MinFrequency != defaults.Processor.MinFrequency {
		t.Errorf("Without config file, should use default MinFrequency")
	}

	if s.Loader.MaxWidth != defaults.Loader.MaxWidth {
		t.Errorf("Without config file, should use default LoaderMaxWidth")
	}

	t.Logf("Successfully loaded default settings when no config file exists")
}

func TestSettings_Load_WithConfigFile(t *testing.T) {
	// Create temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-settings.yaml")
	
	configContent := `
chromatic:
  neutral_threshold: 0.08
  color_merge_threshold: 25.0
  dark_lightness_max: 0.4
processor:
  light_theme_threshold: 0.6
  min_frequency: 0.002
  max_ui_colors: 15
loader:
  max_width: 4096
  max_height: 4096
  allowed_formats: ["jpeg", "png"]
formats:
  quantization_bits: 6
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
	t.Logf("  Neutral threshold: %.3f", s.Chromatic.NeutralThreshold)
	t.Logf("  Light theme threshold: %.3f", s.Processor.LightThemeThreshold)
	t.Logf("  Loader max width: %d", s.Loader.MaxWidth)
	t.Logf("  Quantization bits: %d", s.Formats.QuantizationBits)

	// Verify values were loaded from config
	if s.Chromatic.NeutralThreshold != 0.08 {
		t.Errorf("Expected neutral threshold 0.08, got %.3f", s.Chromatic.NeutralThreshold)
	}

	if s.Processor.LightThemeThreshold != 0.6 {
		t.Errorf("Expected light theme threshold 0.6, got %.3f", s.Processor.LightThemeThreshold)
	}

	if s.Loader.MaxWidth != 4096 {
		t.Errorf("Expected loader max width 4096, got %d", s.Loader.MaxWidth)
	}

	if s.Formats.QuantizationBits != 6 {
		t.Errorf("Expected quantization bits 6, got %d", s.Formats.QuantizationBits)
	}

	if s.DefaultDark != "#121212" {
		t.Errorf("Expected default dark color #121212, got %s", s.DefaultDark)
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
	v.SetDefault("chromatic.neutral_threshold", 0.1)
	v.SetDefault("processor.min_frequency", 0.005)
	v.SetDefault("loader.max_width", 2048)
	
	var testSettings settings.Settings
	err := v.Unmarshal(&testSettings)
	if err != nil {
		t.Fatalf("Failed to unmarshal with Viper: %v", err)
	}
	
	t.Logf("Viper integration test:")
	t.Logf("  Neutral threshold: %.3f", testSettings.Chromatic.NeutralThreshold)
	t.Logf("  Min frequency: %.6f", testSettings.Processor.MinFrequency)
	t.Logf("  Loader max width: %d", testSettings.Loader.MaxWidth)

	if testSettings.Chromatic.NeutralThreshold != 0.1 {
		t.Errorf("Expected neutral threshold 0.1, got %.3f", testSettings.Chromatic.NeutralThreshold)
	}

	if testSettings.Processor.MinFrequency != 0.005 {
		t.Errorf("Expected min frequency 0.005, got %.6f", testSettings.Processor.MinFrequency)
	}

	if testSettings.Loader.MaxWidth != 2048 {
		t.Errorf("Expected loader max width 2048, got %d", testSettings.Loader.MaxWidth)
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
	t.Logf("  Original settings neutral threshold: %.3f", s.Chromatic.NeutralThreshold)
	t.Logf("  Retrieved settings neutral threshold: %.3f", retrievedSettings.Chromatic.NeutralThreshold)
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
		"chromatic": map[string]interface{}{
			"neutral_threshold":              0.07,
			"color_merge_threshold":          18.0,
			"dark_lightness_max":             0.3,
			"light_lightness_min":            0.8,
			"muted_saturation_max":           0.25,
			"vibrant_saturation_min":         0.75,
		},
		"processor": map[string]interface{}{
			"min_frequency":                  0.0005,
			"light_theme_threshold":          0.55,
			"max_ui_colors":                  25,
			"min_cluster_weight":             0.01,
		},
		"loader": map[string]interface{}{
			"max_width":                      6144,
			"max_height":                     6144,
			"allowed_formats":                []string{"jpeg", "png", "webp"},
		},
		"formats": map[string]interface{}{
			"quantization_bits":              6,
		},
		"default_dark":  "#181818",
		"default_light": "#f0f0f0",
		"default_gray":  "#909090",
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
	t.Logf("  Neutral threshold: %.3f (expected 0.07)", s.Chromatic.NeutralThreshold)
	t.Logf("  Loader max width: %d (expected 6144)", s.Loader.MaxWidth)
	t.Logf("  Processor max UI colors: %d (expected 25)", s.Processor.MaxUIColors)
	t.Logf("  Default dark: %s (expected #181818)", s.DefaultDark)

	// Verify unmarshal worked correctly
	if s.Chromatic.NeutralThreshold != 0.07 {
		t.Errorf("Expected neutral threshold 0.07, got %.3f", s.Chromatic.NeutralThreshold)
	}

	if s.Loader.MaxWidth != 6144 {
		t.Errorf("Expected loader max width 6144, got %d", s.Loader.MaxWidth)
	}

	if s.Processor.MaxUIColors != 25 {
		t.Errorf("Expected processor max UI colors 25, got %d", s.Processor.MaxUIColors)
	}

	if s.DefaultDark != "#181818" {
		t.Errorf("Expected default dark #181818, got %s", s.DefaultDark)
	}

	if s.Formats.QuantizationBits != 6 {
		t.Errorf("Expected quantization bits 6, got %d", s.Formats.QuantizationBits)
	}
}