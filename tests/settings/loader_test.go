package settings_test

import (
	"context"
	"strings"
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
	
	if s.MinContrastRatio != defaults.MinContrastRatio {
		t.Errorf("Without config file, should use default MinContrastRatio")
	}
}

func TestSettings_LoadWithViper_ValidConfig(t *testing.T) {
	// Test loading custom config using Viper
	v := viper.New()
	
	// Set custom values
	v.Set("grayscale_threshold", 0.08)
	v.Set("monochromatic_tolerance", 20.0)
	v.Set("theme_mode_threshold", 0.6)
	v.Set("min_frequency", 0.002)
	v.Set("min_contrast_ratio", 7.0)
	v.Set("light_background_fallback", "#f0f0f0")
	v.Set("primary_fallback", "#2266aa")
	
	// Load settings from Viper
	s, err := settings.LoadWithViper(v)
	if err != nil {
		t.Fatalf("LoadWithViper() failed: %v", err)
	}
	
	// Verify custom values are loaded
	if s.GrayscaleThreshold != 0.08 {
		t.Errorf("Expected GrayscaleThreshold 0.08, got %v", s.GrayscaleThreshold)
	}
	
	if s.MonochromaticTolerance != 20.0 {
		t.Errorf("Expected MonochromaticTolerance 20.0, got %v", s.MonochromaticTolerance)
	}
	
	if s.MinContrastRatio != 7.0 {
		t.Errorf("Expected MinContrastRatio 7.0 (AAA), got %v", s.MinContrastRatio)
	}
	
	if s.LightBackgroundFallback != "#f0f0f0" {
		t.Errorf("Expected custom fallback color, got %s", s.LightBackgroundFallback)
	}
}

func TestSettings_LoadWithViper_InvalidConfig(t *testing.T) {
	// Test that LoadWithViper properly handles type mismatches
	v := viper.New()
	
	// Set some invalid values 
	v.Set("grayscale_threshold", "not_a_number") // Invalid type - string instead of float
	v.Set("min_contrast_ratio", -1.0)            // Invalid range (but valid type)
	
	// LoadWithViper should return error for type mismatches
	s, err := settings.LoadWithViper(v)
	if err == nil {
		t.Error("LoadWithViper() should return error for invalid type conversions")
	}
	
	if s != nil {
		t.Error("LoadWithViper() should return nil settings on error")
	}
	
	t.Logf("Expected error for invalid config: %v", err)
}

func TestSettings_LoadWithViper_PartialConfig(t *testing.T) {
	// Test partial configuration with defaults
	v := viper.New()
	
	// Apply defaults first (we can't access the private setDefaults function,
	// so we'll set some defaults manually for testing)
	v.SetDefault("grayscale_threshold", 0.05)
	v.SetDefault("min_contrast_ratio", 4.5)
	
	// Override only some values
	v.Set("min_contrast_ratio", 5.5)
	v.Set("light_background_fallback", "#fafafa")
	v.Set("loader_max_width", 4096)
	
	s, err := settings.LoadWithViper(v)
	if err != nil {
		t.Fatalf("LoadWithViper() should handle partial config: %v", err)
	}
	
	// Custom values should be loaded
	if s.MinContrastRatio != 5.5 {
		t.Errorf("Expected custom MinContrastRatio 5.5, got %v", s.MinContrastRatio)
	}
	
	if s.LightBackgroundFallback != "#fafafa" {
		t.Errorf("Expected custom fallback, got %s", s.LightBackgroundFallback)
	}
	
	if s.LoaderMaxWidth != 4096 {
		t.Errorf("Expected custom LoaderMaxWidth 4096, got %v", s.LoaderMaxWidth)
	}
	
	// Missing values should use defaults
	defaults := settings.DefaultSettings()
	if s.GrayscaleThreshold != defaults.GrayscaleThreshold {
		t.Errorf("Missing values should use defaults")
	}
}

func TestSettings_Context(t *testing.T) {
	s := settings.DefaultSettings()
	s.MinContrastRatio = 6.0 // Custom value for testing
	
	// Test setting context
	ctx := context.Background()
	ctxWithSettings := settings.WithSettings(ctx, s)
	
	// Test retrieving from context
	retrieved := settings.FromContext(ctxWithSettings)
	if retrieved == nil {
		t.Fatal("FromContext() should return settings")
	}
	
	if retrieved.MinContrastRatio != 6.0 {
		t.Errorf("Context should preserve custom settings, got %v", retrieved.MinContrastRatio)
	}
	
	// Test retrieving from context without settings (should load defaults)
	emptyCtx := context.Background()
	defaultRetrieved := settings.FromContext(emptyCtx)
	if defaultRetrieved == nil {
		t.Fatal("FromContext() should return defaults when no settings in context")
	}
	
	// Should be default value, not our custom 6.0
	if defaultRetrieved.MinContrastRatio == 6.0 {
		t.Error("Empty context should not return custom settings")
	}
}

func TestSettings_GetConfigPaths(t *testing.T) {
	// Test config path functions
	userPath := settings.GetUserConfigPath()
	if userPath == "" {
		t.Error("GetUserConfigPath() should return non-empty path")
	}
	
	systemPath := settings.GetSystemConfigPath()
	if systemPath == "" {
		t.Error("GetSystemConfigPath() should return non-empty path")
	}
	
	// Paths should be different
	if userPath == systemPath {
		t.Error("User and system config paths should be different")
	}
	
	// Should contain expected components
	if !strings.Contains(systemPath, "/etc") {
		t.Error("System path should contain /etc")
	}
}