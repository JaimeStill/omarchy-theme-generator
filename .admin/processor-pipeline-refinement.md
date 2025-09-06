# Enhanced Color Categorization Implementation Guide

## Overview

This guide provides a complete implementation plan for enhancing the `pkg/processor` pipeline to categorize colors by their potential theme roles during extraction. This fundamental shift transforms the processor from a simple color extractor into an intelligent theme-color analyzer that understands not just what colors exist, but what role they could play in a cohesive theme.

## Key Architectural Changes

### Before: Role-Based Assignment
- Colors extracted by frequency
- Roles assigned after extraction based on simple rules
- Limited metadata for palette generation
- Hard-coded characteristics

### After: Category-Based Extraction
- Colors evaluated for theme potential during extraction
- Multiple candidates per category with scoring
- Rich metadata for palette generation
- Fully configurable characteristics via settings
- Statistical analysis of color distribution

## Implementation Steps

### Step 1: Expand Settings Structure

Create comprehensive settings for category characteristics and scoring weights.

**File: `pkg/settings/settings.go`**

```go
package settings

import "context"

type contextKey string

const settingsKey contextKey = "settings"

type Settings struct {
    // Core analysis settings (preserved)
    GrayscaleThreshold       float64  `mapstructure:"grayscale_threshold"`
    MonochromaticTolerance   float64  `mapstructure:"monochromatic_tolerance"`
    ThemeModeThreshold       float64  `mapstructure:"theme_mode_threshold"`
    MinFrequency             float64  `mapstructure:"min_frequency"`
    
    // Loader settings (preserved)
    LoaderMaxWidth           int      `mapstructure:"loader_max_width"`
    LoaderMaxHeight          int      `mapstructure:"loader_max_height"`
    LoaderAllowedFormats     []string `mapstructure:"loader_allowed_formats"`
    
    // Fallback colors (preserved)
    LightBackgroundFallback  string   `mapstructure:"light_background_fallback"`
    DarkBackgroundFallback   string   `mapstructure:"dark_background_fallback"`
    LightForegroundFallback  string   `mapstructure:"light_foreground_fallback"`
    DarkForegroundFallback   string   `mapstructure:"dark_foreground_fallback"`
    PrimaryFallback          string   `mapstructure:"primary_fallback"`
    
    // REMOVED: These settings are now handled by category characteristics
    // - LightBackgroundThreshold, DarkBackgroundThreshold (replaced by category lightness ranges)
    // - MinContrastRatio (now per-category min_contrast)
    // - MinPrimarySaturation, MinAccentSaturation (now category min_saturation)
    // - MinAccentLightness, MaxAccentLightness (now category lightness ranges)
    // - DarkLightThreshold, BrightLightThreshold (replaced by category ranges)
    // - ExtremeLightnessPenalty, OptimalLightnessBonus (replaced by scoring system)
    // - MinSaturationForBonus (replaced by scoring system)
    
    // NEW: Category-based extraction settings
    Categories      CategorySettings        `mapstructure:"categories"`
    CategoryScoring CategoryScoringWeights  `mapstructure:"category_scoring"`
    Extraction      ExtractionSettings      `mapstructure:"extraction"`
}

// CategorySettings holds characteristics for all categories by theme mode
type CategorySettings struct {
    Dark  map[string]CategoryCharacteristics `mapstructure:"dark"`
    Light map[string]CategoryCharacteristics `mapstructure:"light"`
}

// CategoryCharacteristics defines HSL and contrast requirements for a category
type CategoryCharacteristics struct {
    MinLightness   float64  `mapstructure:"min_lightness"`
    MaxLightness   float64  `mapstructure:"max_lightness"`
    MinSaturation  float64  `mapstructure:"min_saturation"`
    MaxSaturation  float64  `mapstructure:"max_saturation"`
    MinContrast    float64  `mapstructure:"min_contrast"`
    HueCenter      *float64 `mapstructure:"hue_center"`     // Optional: preferred hue
    HueTolerance   *float64 `mapstructure:"hue_tolerance"`  // Optional: ± degrees from center
}

// CategoryScoringWeights controls how colors are scored for categories
type CategoryScoringWeights struct {
    Frequency    float64 `mapstructure:"frequency"`     // Weight for color frequency
    Contrast     float64 `mapstructure:"contrast"`      // Weight for contrast ratio
    Saturation   float64 `mapstructure:"saturation"`    // Weight for saturation fit
    HueAlignment float64 `mapstructure:"hue_alignment"` // Weight for hue matching
    Lightness    float64 `mapstructure:"lightness"`     // Weight for lightness fit
}

// ExtractionSettings controls the extraction process
type ExtractionSettings struct {
    MaxCandidatesPerCategory int     `mapstructure:"max_candidates_per_category"`
    AllowColoredBackgrounds  bool    `mapstructure:"allow_colored_backgrounds"`
    PreferVibrantAccents     bool    `mapstructure:"prefer_vibrant_accents"`
    MaintainHueConsistency   bool    `mapstructure:"maintain_hue_consistency"`
    GrayscaleHueTemperature  float64 `mapstructure:"grayscale_hue_temperature"`
    MinimumColorFrequency    float64 `mapstructure:"minimum_color_frequency"`
}

func WithSettings(ctx context.Context, s *Settings) context.Context {
    return context.WithValue(ctx, settingsKey, s)
}

func FromContext(ctx context.Context) *Settings {
    if s, ok := ctx.Value(settingsKey).(*Settings); ok {
        return s
    }
    s, _ := Load()
    return s
}
```

### Step 2: Update Default Settings

Add comprehensive defaults for all category characteristics.

**File: `pkg/settings/category_defaults.go`**

```go
package settings

import "github.com/spf13/viper"

// setCategoryDefaults sets all category characteristic defaults
func setCategoryDefaults(v *viper.Viper) {
    // Extraction settings
    v.SetDefault("extraction.max_candidates_per_category", 5)
    v.SetDefault("extraction.allow_colored_backgrounds", false)
    v.SetDefault("extraction.prefer_vibrant_accents", true)
    v.SetDefault("extraction.maintain_hue_consistency", true)
    v.SetDefault("extraction.grayscale_hue_temperature", 220.0)
    v.SetDefault("extraction.minimum_color_frequency", 0.0001)
    
    // Category scoring weights (must sum to 1.0 for normalized scores)
    v.SetDefault("category_scoring.frequency", 0.25)
    v.SetDefault("category_scoring.contrast", 0.25)
    v.SetDefault("category_scoring.saturation", 0.20)
    v.SetDefault("category_scoring.hue_alignment", 0.15)
    v.SetDefault("category_scoring.lightness", 0.15)
    
    // Set dark and light mode defaults
    setDarkCategoryDefaults(v)
    setLightCategoryDefaults(v)
}

func setDarkCategoryDefaults(v *viper.Viper) {
    prefix := "categories.dark."
    
    // Core UI elements
    v.SetDefault(prefix+"background.min_lightness", 0.0)
    v.SetDefault(prefix+"background.max_lightness", 0.15)
    v.SetDefault(prefix+"background.min_saturation", 0.0)
    v.SetDefault(prefix+"background.max_saturation", 0.2)
    v.SetDefault(prefix+"background.min_contrast", 0.0)
    
    v.SetDefault(prefix+"foreground.min_lightness", 0.85)
    v.SetDefault(prefix+"foreground.max_lightness", 1.0)
    v.SetDefault(prefix+"foreground.min_saturation", 0.0)
    v.SetDefault(prefix+"foreground.max_saturation", 0.1)
    v.SetDefault(prefix+"foreground.min_contrast", 4.5)
    
    v.SetDefault(prefix+"dim_foreground.min_lightness", 0.4)
    v.SetDefault(prefix+"dim_foreground.max_lightness", 0.6)
    v.SetDefault(prefix+"dim_foreground.min_saturation", 0.0)
    v.SetDefault(prefix+"dim_foreground.max_saturation", 0.3)
    v.SetDefault(prefix+"dim_foreground.min_contrast", 3.0)
    
    v.SetDefault(prefix+"cursor.min_lightness", 0.7)
    v.SetDefault(prefix+"cursor.max_lightness", 1.0)
    v.SetDefault(prefix+"cursor.min_saturation", 0.0)
    v.SetDefault(prefix+"cursor.max_saturation", 0.5)
    v.SetDefault(prefix+"cursor.min_contrast", 7.0)
    
    // Terminal normal colors
    v.SetDefault(prefix+"normal_black.min_lightness", 0.15)
    v.SetDefault(prefix+"normal_black.max_lightness", 0.25)
    v.SetDefault(prefix+"normal_black.min_saturation", 0.0)
    v.SetDefault(prefix+"normal_black.max_saturation", 0.1)
    v.SetDefault(prefix+"normal_black.min_contrast", 1.5)
    
    v.SetDefault(prefix+"normal_red.min_lightness", 0.35)
    v.SetDefault(prefix+"normal_red.max_lightness", 0.55)
    v.SetDefault(prefix+"normal_red.min_saturation", 0.6)
    v.SetDefault(prefix+"normal_red.max_saturation", 1.0)
    v.SetDefault(prefix+"normal_red.min_contrast", 2.0)
    v.SetDefault(prefix+"normal_red.hue_center", 0.0)
    v.SetDefault(prefix+"normal_red.hue_tolerance", 25.0)
    
    v.SetDefault(prefix+"normal_green.min_lightness", 0.35)
    v.SetDefault(prefix+"normal_green.max_lightness", 0.55)
    v.SetDefault(prefix+"normal_green.min_saturation", 0.5)
    v.SetDefault(prefix+"normal_green.max_saturation", 1.0)
    v.SetDefault(prefix+"normal_green.min_contrast", 2.0)
    v.SetDefault(prefix+"normal_green.hue_center", 120.0)
    v.SetDefault(prefix+"normal_green.hue_tolerance", 40.0)
    
    v.SetDefault(prefix+"normal_yellow.min_lightness", 0.45)
    v.SetDefault(prefix+"normal_yellow.max_lightness", 0.65)
    v.SetDefault(prefix+"normal_yellow.min_saturation", 0.7)
    v.SetDefault(prefix+"normal_yellow.max_saturation", 1.0)
    v.SetDefault(prefix+"normal_yellow.min_contrast", 2.5)
    v.SetDefault(prefix+"normal_yellow.hue_center", 60.0)
    v.SetDefault(prefix+"normal_yellow.hue_tolerance", 20.0)
    
    v.SetDefault(prefix+"normal_blue.min_lightness", 0.35)
    v.SetDefault(prefix+"normal_blue.max_lightness", 0.55)
    v.SetDefault(prefix+"normal_blue.min_saturation", 0.6)
    v.SetDefault(prefix+"normal_blue.max_saturation", 1.0)
    v.SetDefault(prefix+"normal_blue.min_contrast", 2.0)
    v.SetDefault(prefix+"normal_blue.hue_center", 240.0)
    v.SetDefault(prefix+"normal_blue.hue_tolerance", 30.0)
    
    v.SetDefault(prefix+"normal_magenta.min_lightness", 0.35)
    v.SetDefault(prefix+"normal_magenta.max_lightness", 0.55)
    v.SetDefault(prefix+"normal_magenta.min_saturation", 0.5)
    v.SetDefault(prefix+"normal_magenta.max_saturation", 1.0)
    v.SetDefault(prefix+"normal_magenta.min_contrast", 2.0)
    v.SetDefault(prefix+"normal_magenta.hue_center", 300.0)
    v.SetDefault(prefix+"normal_magenta.hue_tolerance", 30.0)
    
    v.SetDefault(prefix+"normal_cyan.min_lightness", 0.35)
    v.SetDefault(prefix+"normal_cyan.max_lightness", 0.55)
    v.SetDefault(prefix+"normal_cyan.min_saturation", 0.5)
    v.SetDefault(prefix+"normal_cyan.max_saturation", 1.0)
    v.SetDefault(prefix+"normal_cyan.min_contrast", 2.0)
    v.SetDefault(prefix+"normal_cyan.hue_center", 180.0)
    v.SetDefault(prefix+"normal_cyan.hue_tolerance", 30.0)
    
    v.SetDefault(prefix+"normal_white.min_lightness", 0.6)
    v.SetDefault(prefix+"normal_white.max_lightness", 0.8)
    v.SetDefault(prefix+"normal_white.min_saturation", 0.0)
    v.SetDefault(prefix+"normal_white.max_saturation", 0.1)
    v.SetDefault(prefix+"normal_white.min_contrast", 3.5)
    
    // Terminal bright colors (higher lightness/saturation than normal)
    v.SetDefault(prefix+"bright_black.min_lightness", 0.25)
    v.SetDefault(prefix+"bright_black.max_lightness", 0.35)
    v.SetDefault(prefix+"bright_black.min_saturation", 0.0)
    v.SetDefault(prefix+"bright_black.max_saturation", 0.1)
    v.SetDefault(prefix+"bright_black.min_contrast", 2.0)
    
    v.SetDefault(prefix+"bright_red.min_lightness", 0.5)
    v.SetDefault(prefix+"bright_red.max_lightness", 0.7)
    v.SetDefault(prefix+"bright_red.min_saturation", 0.8)
    v.SetDefault(prefix+"bright_red.max_saturation", 1.0)
    v.SetDefault(prefix+"bright_red.min_contrast", 3.0)
    v.SetDefault(prefix+"bright_red.hue_center", 0.0)
    v.SetDefault(prefix+"bright_red.hue_tolerance", 25.0)
    
    v.SetDefault(prefix+"bright_green.min_lightness", 0.5)
    v.SetDefault(prefix+"bright_green.max_lightness", 0.7)
    v.SetDefault(prefix+"bright_green.min_saturation", 0.7)
    v.SetDefault(prefix+"bright_green.max_saturation", 1.0)
    v.SetDefault(prefix+"bright_green.min_contrast", 3.0)
    v.SetDefault(prefix+"bright_green.hue_center", 120.0)
    v.SetDefault(prefix+"bright_green.hue_tolerance", 40.0)
    
    v.SetDefault(prefix+"bright_yellow.min_lightness", 0.6)
    v.SetDefault(prefix+"bright_yellow.max_lightness", 0.8)
    v.SetDefault(prefix+"bright_yellow.min_saturation", 0.8)
    v.SetDefault(prefix+"bright_yellow.max_saturation", 1.0)
    v.SetDefault(prefix+"bright_yellow.min_contrast", 3.5)
    v.SetDefault(prefix+"bright_yellow.hue_center", 60.0)
    v.SetDefault(prefix+"bright_yellow.hue_tolerance", 20.0)
    
    v.SetDefault(prefix+"bright_blue.min_lightness", 0.5)
    v.SetDefault(prefix+"bright_blue.max_lightness", 0.7)
    v.SetDefault(prefix+"bright_blue.min_saturation", 0.8)
    v.SetDefault(prefix+"bright_blue.max_saturation", 1.0)
    v.SetDefault(prefix+"bright_blue.min_contrast", 3.0)
    v.SetDefault(prefix+"bright_blue.hue_center", 240.0)
    v.SetDefault(prefix+"bright_blue.hue_tolerance", 30.0)
    
    v.SetDefault(prefix+"bright_magenta.min_lightness", 0.5)
    v.SetDefault(prefix+"bright_magenta.max_lightness", 0.7)
    v.SetDefault(prefix+"bright_magenta.min_saturation", 0.7)
    v.SetDefault(prefix+"bright_magenta.max_saturation", 1.0)
    v.SetDefault(prefix+"bright_magenta.min_contrast", 3.0)
    v.SetDefault(prefix+"bright_magenta.hue_center", 300.0)
    v.SetDefault(prefix+"bright_magenta.hue_tolerance", 30.0)
    
    v.SetDefault(prefix+"bright_cyan.min_lightness", 0.5)
    v.SetDefault(prefix+"bright_cyan.max_lightness", 0.7)
    v.SetDefault(prefix+"bright_cyan.min_saturation", 0.7)
    v.SetDefault(prefix+"bright_cyan.max_saturation", 1.0)
    v.SetDefault(prefix+"bright_cyan.min_contrast", 3.0)
    v.SetDefault(prefix+"bright_cyan.hue_center", 180.0)
    v.SetDefault(prefix+"bright_cyan.hue_tolerance", 30.0)
    
    v.SetDefault(prefix+"bright_white.min_lightness", 0.8)
    v.SetDefault(prefix+"bright_white.max_lightness", 1.0)
    v.SetDefault(prefix+"bright_white.min_saturation", 0.0)
    v.SetDefault(prefix+"bright_white.max_saturation", 0.1)
    v.SetDefault(prefix+"bright_white.min_contrast", 4.0)
    
    // Accent colors
    v.SetDefault(prefix+"accent_primary.min_lightness", 0.4)
    v.SetDefault(prefix+"accent_primary.max_lightness", 0.7)
    v.SetDefault(prefix+"accent_primary.min_saturation", 0.6)
    v.SetDefault(prefix+"accent_primary.max_saturation", 1.0)
    v.SetDefault(prefix+"accent_primary.min_contrast", 3.0)
    
    v.SetDefault(prefix+"accent_secondary.min_lightness", 0.35)
    v.SetDefault(prefix+"accent_secondary.max_lightness", 0.65)
    v.SetDefault(prefix+"accent_secondary.min_saturation", 0.5)
    v.SetDefault(prefix+"accent_secondary.max_saturation", 0.9)
    v.SetDefault(prefix+"accent_secondary.min_contrast", 2.5)
    
    v.SetDefault(prefix+"accent_tertiary.min_lightness", 0.3)
    v.SetDefault(prefix+"accent_tertiary.max_lightness", 0.6)
    v.SetDefault(prefix+"accent_tertiary.min_saturation", 0.4)
    v.SetDefault(prefix+"accent_tertiary.max_saturation", 0.8)
    v.SetDefault(prefix+"accent_tertiary.min_contrast", 2.0)
    
    // Semantic colors
    v.SetDefault(prefix+"error.min_lightness", 0.4)
    v.SetDefault(prefix+"error.max_lightness", 0.6)
    v.SetDefault(prefix+"error.min_saturation", 0.7)
    v.SetDefault(prefix+"error.max_saturation", 1.0)
    v.SetDefault(prefix+"error.min_contrast", 3.0)
    v.SetDefault(prefix+"error.hue_center", 0.0)
    v.SetDefault(prefix+"error.hue_tolerance", 20.0)
    
    v.SetDefault(prefix+"warning.min_lightness", 0.45)
    v.SetDefault(prefix+"warning.max_lightness", 0.65)
    v.SetDefault(prefix+"warning.min_saturation", 0.7)
    v.SetDefault(prefix+"warning.max_saturation", 1.0)
    v.SetDefault(prefix+"warning.min_contrast", 3.0)
    v.SetDefault(prefix+"warning.hue_center", 45.0)
    v.SetDefault(prefix+"warning.hue_tolerance", 15.0)
    
    v.SetDefault(prefix+"success.min_lightness", 0.35)
    v.SetDefault(prefix+"success.max_lightness", 0.55)
    v.SetDefault(prefix+"success.min_saturation", 0.5)
    v.SetDefault(prefix+"success.max_saturation", 1.0)
    v.SetDefault(prefix+"success.min_contrast", 3.0)
    v.SetDefault(prefix+"success.hue_center", 120.0)
    v.SetDefault(prefix+"success.hue_tolerance", 30.0)
    
    v.SetDefault(prefix+"info.min_lightness", 0.35)
    v.SetDefault(prefix+"info.max_lightness", 0.55)
    v.SetDefault(prefix+"info.min_saturation", 0.5)
    v.SetDefault(prefix+"info.max_saturation", 1.0)
    v.SetDefault(prefix+"info.min_contrast", 3.0)
    v.SetDefault(prefix+"info.hue_center", 210.0)
    v.SetDefault(prefix+"info.hue_tolerance", 30.0)
}

func setLightCategoryDefaults(v *viper.Viper) {
    prefix := "categories.light."
    
    // Core UI elements - light mode inverts most lightness ranges
    v.SetDefault(prefix+"background.min_lightness", 0.90)
    v.SetDefault(prefix+"background.max_lightness", 1.0)
    v.SetDefault(prefix+"background.min_saturation", 0.0)
    v.SetDefault(prefix+"background.max_saturation", 0.15)
    v.SetDefault(prefix+"background.min_contrast", 0.0)
    
    v.SetDefault(prefix+"foreground.min_lightness", 0.0)
    v.SetDefault(prefix+"foreground.max_lightness", 0.20)
    v.SetDefault(prefix+"foreground.min_saturation", 0.0)
    v.SetDefault(prefix+"foreground.max_saturation", 0.1)
    v.SetDefault(prefix+"foreground.min_contrast", 4.5)
    
    v.SetDefault(prefix+"dim_foreground.min_lightness", 0.3)
    v.SetDefault(prefix+"dim_foreground.max_lightness", 0.5)
    v.SetDefault(prefix+"dim_foreground.min_saturation", 0.0)
    v.SetDefault(prefix+"dim_foreground.max_saturation", 0.3)
    v.SetDefault(prefix+"dim_foreground.min_contrast", 3.0)
    
    v.SetDefault(prefix+"cursor.min_lightness", 0.0)
    v.SetDefault(prefix+"cursor.max_lightness", 0.3)
    v.SetDefault(prefix+"cursor.min_saturation", 0.0)
    v.SetDefault(prefix+"cursor.max_saturation", 0.5)
    v.SetDefault(prefix+"cursor.min_contrast", 7.0)
    
    // Terminal normal colors - darker for light mode
    v.SetDefault(prefix+"normal_black.min_lightness", 0.0)
    v.SetDefault(prefix+"normal_black.max_lightness", 0.15)
    v.SetDefault(prefix+"normal_black.min_saturation", 0.0)
    v.SetDefault(prefix+"normal_black.max_saturation", 0.1)
    v.SetDefault(prefix+"normal_black.min_contrast", 7.0)
    
    v.SetDefault(prefix+"normal_red.min_lightness", 0.25)
    v.SetDefault(prefix+"normal_red.max_lightness", 0.45)
    v.SetDefault(prefix+"normal_red.min_saturation", 0.6)
    v.SetDefault(prefix+"normal_red.max_saturation", 1.0)
    v.SetDefault(prefix+"normal_red.min_contrast", 4.5)
    v.SetDefault(prefix+"normal_red.hue_center", 0.0)
    v.SetDefault(prefix+"normal_red.hue_tolerance", 25.0)
    
    v.SetDefault(prefix+"normal_green.min_lightness", 0.25)
    v.SetDefault(prefix+"normal_green.max_lightness", 0.45)
    v.SetDefault(prefix+"normal_green.min_saturation", 0.5)
    v.SetDefault(prefix+"normal_green.max_saturation", 1.0)
    v.SetDefault(prefix+"normal_green.min_contrast", 4.5)
    v.SetDefault(prefix+"normal_green.hue_center", 120.0)
    v.SetDefault(prefix+"normal_green.hue_tolerance", 40.0)
    
    v.SetDefault(prefix+"normal_yellow.min_lightness", 0.3)
    v.SetDefault(prefix+"normal_yellow.max_lightness", 0.5)
    v.SetDefault(prefix+"normal_yellow.min_saturation", 0.7)
    v.SetDefault(prefix+"normal_yellow.max_saturation", 1.0)
    v.SetDefault(prefix+"normal_yellow.min_contrast", 4.5)
    v.SetDefault(prefix+"normal_yellow.hue_center", 60.0)
    v.SetDefault(prefix+"normal_yellow.hue_tolerance", 20.0)
    
    v.SetDefault(prefix+"normal_blue.min_lightness", 0.25)
    v.SetDefault(prefix+"normal_blue.max_lightness", 0.45)
    v.SetDefault(prefix+"normal_blue.min_saturation", 0.6)
    v.SetDefault(prefix+"normal_blue.max_saturation", 1.0)
    v.SetDefault(prefix+"normal_blue.min_contrast", 4.5)
    v.SetDefault(prefix+"normal_blue.hue_center", 240.0)
    v.SetDefault(prefix+"normal_blue.hue_tolerance", 30.0)
    
    v.SetDefault(prefix+"normal_magenta.min_lightness", 0.25)
    v.SetDefault(prefix+"normal_magenta.max_lightness", 0.45)
    v.SetDefault(prefix+"normal_magenta.min_saturation", 0.5)
    v.SetDefault(prefix+"normal_magenta.max_saturation", 1.0)
    v.SetDefault(prefix+"normal_magenta.min_contrast", 4.5)
    v.SetDefault(prefix+"normal_magenta.hue_center", 300.0)
    v.SetDefault(prefix+"normal_magenta.hue_tolerance", 30.0)
    
    v.SetDefault(prefix+"normal_cyan.min_lightness", 0.25)
    v.SetDefault(prefix+"normal_cyan.max_lightness", 0.45)
    v.SetDefault(prefix+"normal_cyan.min_saturation", 0.5)
    v.SetDefault(prefix+"normal_cyan.max_saturation", 1.0)
    v.SetDefault(prefix+"normal_cyan.min_contrast", 4.5)
    v.SetDefault(prefix+"normal_cyan.hue_center", 180.0)
    v.SetDefault(prefix+"normal_cyan.hue_tolerance", 30.0)
    
    v.SetDefault(prefix+"normal_white.min_lightness", 0.4)
    v.SetDefault(prefix+"normal_white.max_lightness", 0.6)
    v.SetDefault(prefix+"normal_white.min_saturation", 0.0)
    v.SetDefault(prefix+"normal_white.max_saturation", 0.1)
    v.SetDefault(prefix+"normal_white.min_contrast", 3.5)
    
    // Terminal bright colors
    v.SetDefault(prefix+"bright_black.min_lightness", 0.15)
    v.SetDefault(prefix+"bright_black.max_lightness", 0.25)
    v.SetDefault(prefix+"bright_black.min_saturation", 0.0)
    v.SetDefault(prefix+"bright_black.max_saturation", 0.1)
    v.SetDefault(prefix+"bright_black.min_contrast", 5.0)
    
    v.SetDefault(prefix+"bright_red.min_lightness", 0.15)
    v.SetDefault(prefix+"bright_red.max_lightness", 0.35)
    v.SetDefault(prefix+"bright_red.min_saturation", 0.8)
    v.SetDefault(prefix+"bright_red.max_saturation", 1.0)
    v.SetDefault(prefix+"bright_red.min_contrast", 7.0)
    v.SetDefault(prefix+"bright_red.hue_center", 0.0)
    v.SetDefault(prefix+"bright_red.hue_tolerance", 25.0)
    
    v.SetDefault(prefix+"bright_green.min_lightness", 0.15)
    v.SetDefault(prefix+"bright_green.max_lightness", 0.35)
    v.SetDefault(prefix+"bright_green.min_saturation", 0.7)
    v.SetDefault(prefix+"bright_green.max_saturation", 1.0)
    v.SetDefault(prefix+"bright_green.min_contrast", 7.0)
    v.SetDefault(prefix+"bright_green.hue_center", 120.0)
    v.SetDefault(prefix+"bright_green.hue_tolerance", 40.0)
    
    v.SetDefault(prefix+"bright_yellow.min_lightness", 0.2)
    v.SetDefault(prefix+"bright_yellow.max_lightness", 0.4)
    v.SetDefault(prefix+"bright_yellow.min_saturation", 0.8)
    v.SetDefault(prefix+"bright_yellow.max_saturation", 1.0)
    v.SetDefault(prefix+"bright_yellow.min_contrast", 7.0)
    v.SetDefault(prefix+"bright_yellow.hue_center", 60.0)
    v.SetDefault(prefix+"bright_yellow.hue_tolerance", 20.0)
    
    v.SetDefault(prefix+"bright_blue.min_lightness", 0.15)
    v.SetDefault(prefix+"bright_blue.max_lightness", 0.35)
    v.SetDefault(prefix+"bright_blue.min_saturation", 0.8)
    v.SetDefault(prefix+"bright_blue.max_saturation", 1.0)
    v.SetDefault(prefix+"bright_blue.min_contrast", 7.0)
    v.SetDefault(prefix+"bright_blue.hue_center", 240.0)
    v.SetDefault(prefix+"bright_blue.hue_tolerance", 30.0)
    
    v.SetDefault(prefix+"bright_magenta.min_lightness", 0.15)
    v.SetDefault(prefix+"bright_magenta.max_lightness", 0.35)
    v.SetDefault(prefix+"bright_magenta.min_saturation", 0.7)
    v.SetDefault(prefix+"bright_magenta.max_saturation", 1.0)
    v.SetDefault(prefix+"bright_magenta.min_contrast", 7.0)
    v.SetDefault(prefix+"bright_magenta.hue_center", 300.0)
    v.SetDefault(prefix+"bright_magenta.hue_tolerance", 30.0)
    
    v.SetDefault(prefix+"bright_cyan.min_lightness", 0.15)
    v.SetDefault(prefix+"bright_cyan.max_lightness", 0.35)
    v.SetDefault(prefix+"bright_cyan.min_saturation", 0.7)
    v.SetDefault(prefix+"bright_cyan.max_saturation", 1.0)
    v.SetDefault(prefix+"bright_cyan.min_contrast", 7.0)
    v.SetDefault(prefix+"bright_cyan.hue_center", 180.0)
    v.SetDefault(prefix+"bright_cyan.hue_tolerance", 30.0)
    
    v.SetDefault(prefix+"bright_white.min_lightness", 0.2)
    v.SetDefault(prefix+"bright_white.max_lightness", 0.4)
    v.SetDefault(prefix+"bright_white.min_saturation", 0.0)
    v.SetDefault(prefix+"bright_white.max_saturation", 0.1)
    v.SetDefault(prefix+"bright_white.min_contrast", 7.0)
    
    // Accent colors
    v.SetDefault(prefix+"accent_primary.min_lightness", 0.3)
    v.SetDefault(prefix+"accent_primary.max_lightness", 0.6)
    v.SetDefault(prefix+"accent_primary.min_saturation", 0.6)
    v.SetDefault(prefix+"accent_primary.max_saturation", 1.0)
    v.SetDefault(prefix+"accent_primary.min_contrast", 3.0)
    
    v.SetDefault(prefix+"accent_secondary.min_lightness", 0.35)
    v.SetDefault(prefix+"accent_secondary.max_lightness", 0.65)
    v.SetDefault(prefix+"accent_secondary.min_saturation", 0.5)
    v.SetDefault(prefix+"accent_secondary.max_saturation", 0.9)
    v.SetDefault(prefix+"accent_secondary.min_contrast", 2.5)
    
    v.SetDefault(prefix+"accent_tertiary.min_lightness", 0.4)
    v.SetDefault(prefix+"accent_tertiary.max_lightness", 0.7)
    v.SetDefault(prefix+"accent_tertiary.min_saturation", 0.4)
    v.SetDefault(prefix+"accent_tertiary.max_saturation", 0.8)
    v.SetDefault(prefix+"accent_tertiary.min_contrast", 2.0)
    
    // Semantic colors
    v.SetDefault(prefix+"error.min_lightness", 0.3)
    v.SetDefault(prefix+"error.max_lightness", 0.5)
    v.SetDefault(prefix+"error.min_saturation", 0.7)
    v.SetDefault(prefix+"error.max_saturation", 1.0)
    v.SetDefault(prefix+"error.min_contrast", 4.5)
    v.SetDefault(prefix+"error.hue_center", 0.0)
    v.SetDefault(prefix+"error.hue_tolerance", 20.0)
    
    v.SetDefault(prefix+"warning.min_lightness", 0.35)
    v.SetDefault(prefix+"warning.max_lightness", 0.55)
    v.SetDefault(prefix+"warning.min_saturation", 0.7)
    v.SetDefault(prefix+"warning.max_saturation", 1.0)
    v.SetDefault(prefix+"warning.min_contrast", 4.5)
    v.SetDefault(prefix+"warning.hue_center", 45.0)
    v.SetDefault(prefix+"warning.hue_tolerance", 15.0)
    
    v.SetDefault(prefix+"success.min_lightness", 0.25)
    v.SetDefault(prefix+"success.max_lightness", 0.45)
    v.SetDefault(prefix+"success.min_saturation", 0.5)
    v.SetDefault(prefix+"success.max_saturation", 1.0)
    v.SetDefault(prefix+"success.min_contrast", 4.5)
    v.SetDefault(prefix+"success.hue_center", 120.0)
    v.SetDefault(prefix+"success.hue_tolerance", 30.0)
    
    v.SetDefault(prefix+"info.min_lightness", 0.25)
    v.SetDefault(prefix+"info.max_lightness", 0.45)
    v.SetDefault(prefix+"info.min_saturation", 0.5)
    v.SetDefault(prefix+"info.max_saturation", 1.0)
    v.SetDefault(prefix+"info.min_contrast", 4.5)
    v.SetDefault(prefix+"info.hue_center", 210.0)
    v.SetDefault(prefix+"info.hue_tolerance", 30.0)
}
```

### Step 3: Update defaults.go to Remove Unnecessary Settings

**File: `pkg/settings/defaults.go` (updated)**

```go
package settings

import (
    "fmt"
    "github.com/spf13/viper"
    _ "golang.org/x/image/webp"
)

func setDefaults(v *viper.Viper) {
    // Core analysis settings (preserved)
    v.SetDefault("grayscale_threshold", 0.05)
    v.SetDefault("monochromatic_tolerance", 15.0)
    v.SetDefault("theme_mode_threshold", 0.5)
    v.SetDefault("min_frequency", 0.001)
    
    // Loader settings (preserved)
    v.SetDefault("loader_max_width", 8192)
    v.SetDefault("loader_max_height", 8192)
    v.SetDefault("loader_allowed_formats", []string{
        "jpeg",
        "jpg",
        "png",
        "webp",
    })
    
    // Fallback colors (preserved)
    v.SetDefault("light_background_fallback", "#ffffff")
    v.SetDefault("dark_background_fallback", "#202020")
    v.SetDefault("light_foreground_fallback", "#202020")
    v.SetDefault("dark_foreground_fallback", "#ffffff")
    v.SetDefault("primary_fallback", "#6496c8")
    
    // REMOVED: Legacy settings replaced by category system
    // - light_background_threshold, dark_background_threshold
    // - min_contrast_ratio, min_primary_saturation, min_accent_saturation
    // - min_accent_lightness, max_accent_lightness
    // - dark_light_threshold, bright_light_threshold
    // - extreme_lightness_penalty, optimal_lightness_bonus, min_saturation_for_bonus
    
    // NEW: Set category defaults
    setCategoryDefaults(v)
}

func DefaultSettings() *Settings {
    v := viper.New()
    setDefaults(v)
    
    var settings Settings
    if err := v.Unmarshal(&settings); err != nil {
        panic(fmt.Sprintf("failed to unmarshal default settings: %v", err))
    }
    
    return &settings
}
```

### Step 4: Create Category System

**File: `pkg/processor/categories.go`**

```go
package processor

import (
    "math"
    "github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

// ColorCategory represents a theme-oriented color role
type ColorCategory string

// Define all available categories
const (
    // Core UI elements
    CategoryBackground    ColorCategory = "background"
    CategoryForeground    ColorCategory = "foreground"
    CategoryDimForeground ColorCategory = "dim_foreground"
    CategoryCursor        ColorCategory = "cursor"
    
    // Terminal normal colors (ANSI 0-7)
    CategoryNormalBlack   ColorCategory = "normal_black"
    CategoryNormalRed     ColorCategory = "normal_red"
    CategoryNormalGreen   ColorCategory = "normal_green"
    CategoryNormalYellow  ColorCategory = "normal_yellow"
    CategoryNormalBlue    ColorCategory = "normal_blue"
    CategoryNormalMagenta ColorCategory = "normal_magenta"
    CategoryNormalCyan    ColorCategory = "normal_cyan"
    CategoryNormalWhite   ColorCategory = "normal_white"
    
    // Terminal bright colors (ANSI 8-15)
    CategoryBrightBlack   ColorCategory = "bright_black"
    CategoryBrightRed     ColorCategory = "bright_red"
    CategoryBrightGreen   ColorCategory = "bright_green"
    CategoryBrightYellow  ColorCategory = "bright_yellow"
    CategoryBrightBlue    ColorCategory = "bright_blue"
    CategoryBrightMagenta ColorCategory = "bright_magenta"
    CategoryBrightCyan    ColorCategory = "bright_cyan"
    CategoryBrightWhite   ColorCategory = "bright_white"
    
    // Accent colors
    CategoryAccentPrimary   ColorCategory = "accent_primary"
    CategoryAccentSecondary ColorCategory = "accent_secondary"
    CategoryAccentTertiary  ColorCategory = "accent_tertiary"
    
    // Semantic colors
    CategoryError   ColorCategory = "error"
    CategoryWarning ColorCategory = "warning"
    CategorySuccess ColorCategory = "success"
    CategoryInfo    ColorCategory = "info"
)

// GetAllCategories returns all available categories
func GetAllCategories() []ColorCategory {
    return []ColorCategory{
        CategoryBackground,
        CategoryForeground,
        CategoryDimForeground,
        CategoryCursor,
        CategoryNormalBlack,
        CategoryNormalRed,
        CategoryNormalGreen,
        CategoryNormalYellow,
        CategoryNormalBlue,
        CategoryNormalMagenta,
        CategoryNormalCyan,
        CategoryNormalWhite,
        CategoryBrightBlack,
        CategoryBrightRed,
        CategoryBrightGreen,
        CategoryBrightYellow,
        CategoryBrightBlue,
        CategoryBrightMagenta,
        CategoryBrightCyan,
        CategoryBrightWhite,
        CategoryAccentPrimary,
        CategoryAccentSecondary,
        CategoryAccentTertiary,
        CategoryError,
        CategoryWarning,
        CategorySuccess,
        CategoryInfo,
    }
}

// GetCategoryPriorityOrder returns categories in evaluation order
func (p *Processor) GetCategoryPriorityOrder(profile *ColorProfile) []ColorCategory {
    // Background must be first (needed for contrast calculations)
    // Then high-priority UI elements, then terminal colors, then accents
    return []ColorCategory{
        CategoryBackground,
        CategoryForeground,
        CategoryAccentPrimary,
        CategoryDimForeground,
        CategoryCursor,
        CategoryError,
        CategoryWarning,
        CategorySuccess,
        CategoryInfo,
        CategoryAccentSecondary,
        CategoryAccentTertiary,
        CategoryNormalRed,
        CategoryNormalGreen,
        CategoryNormalBlue,
        CategoryNormalYellow,
        CategoryNormalMagenta,
        CategoryNormalCyan,
        CategoryNormalBlack,
        CategoryNormalWhite,
        CategoryBrightRed,
        CategoryBrightGreen,
        CategoryBrightBlue,
        CategoryBrightYellow,
        CategoryBrightMagenta,
        CategoryBrightCyan,
        CategoryBrightBlack,
        CategoryBrightWhite,
    }
}

// GetCategoryCharacteristics retrieves characteristics from settings
func (p *Processor) GetCategoryCharacteristics(
    category ColorCategory,
    profile *ColorProfile,
) settings.CategoryCharacteristics {
    
    var categoryMap map[string]settings.CategoryCharacteristics
    
    if profile.Mode == Light {
        categoryMap = p.settings.Categories.Light
    } else {
        categoryMap = p.settings.Categories.Dark
    }
    
    if chars, ok := categoryMap[string(category)]; ok {
        return chars
    }
    
    // Return permissive defaults if not configured
    return settings.CategoryCharacteristics{
        MinLightness:  0.0,
        MaxLightness:  1.0,
        MinSaturation: 0.0,
        MaxSaturation: 1.0,
        MinContrast:   2.0,
    }
}

// fitsCategory checks if a color meets category requirements
func (p *Processor) fitsCategory(
    c color.RGBA,
    category ColorCategory,
    profile *ColorProfile,
    background color.RGBA,
) bool {
    
    chars := p.GetCategoryCharacteristics(category, profile)
    hsla := formats.RGBAToHSLA(c)
    
    // Check lightness bounds
    if hsla.L < chars.MinLightness || hsla.L > chars.MaxLightness {
        return false
    }
    
    // Check saturation bounds
    if hsla.S < chars.MinSaturation || hsla.S > chars.MaxSaturation {
        return false
    }
    
    // Check contrast requirement (skip for background itself)
    if category != CategoryBackground && chars.MinContrast > 0 {
        contrast := chromatic.ContrastRatio(c, background)
        if contrast < chars.MinContrast {
            return false
        }
    }
    
    // Check hue constraints if specified
    if chars.HueCenter != nil && chars.HueTolerance != nil {
        hueDiff := math.Abs(hsla.H - *chars.HueCenter)
        if hueDiff > 180 {
            hueDiff = 360 - hueDiff
        }
        if hueDiff > *chars.HueTolerance {
            return false
        }
    }
    
    return true
}
```

### Step 5: Create Scoring System

**File: `pkg/processor/scoring.go`**

```go
package processor

import (
    "image/color"
    "math"
    "github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
    "github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

// calculateCategoryFitScore computes how well a color fits a category
func (p *Processor) calculateCategoryFitScore(
    c color.RGBA,
    category ColorCategory,
    profile *ColorProfile,
    background color.RGBA,
    frequency uint32,
    totalPixels uint32,
) float64 {
    
    chars := p.GetCategoryCharacteristics(category, profile)
    hsla := formats.RGBAToHSLA(c)
    weights := p.settings.CategoryScoring
    
    score := 0.0
    
    // Frequency score (normalized by total pixels)
    if weights.Frequency > 0 {
        freqRatio := float64(frequency) / float64(totalPixels)
        // Use logarithmic scale to prevent dominant colors from overwhelming
        freqScore := math.Min(1.0, math.Log10(freqRatio*10000 + 1) / 4)
        score += freqScore * weights.Frequency
    }
    
    // Contrast score (only if minimum contrast required)
    if weights.Contrast > 0 && chars.MinContrast > 0 {
        contrast := chromatic.ContrastRatio(c, background)
        if contrast >= chars.MinContrast {
            // Score increases with contrast above minimum
            contrastScore := math.Min(1.0, (contrast - chars.MinContrast) / 10.0)
            score += contrastScore * weights.Contrast
        } else {
            // Insufficient contrast disqualifies the color
            return 0
        }
    }
    
    // Saturation score (proximity to ideal range midpoint)
    if weights.Saturation > 0 {
        idealSat := (chars.MinSaturation + chars.MaxSaturation) / 2
        satRange := chars.MaxSaturation - chars.MinSaturation
        if satRange > 0 {
            satDiff := math.Abs(hsla.S - idealSat)
            satScore := 1.0 - (satDiff / satRange)
            score += satScore * weights.Saturation
        } else {
            score += weights.Saturation // Perfect fit if range is zero
        }
    }
    
    // Lightness score (proximity to ideal range midpoint)
    if weights.Lightness > 0 {
        idealLight := (chars.MinLightness + chars.MaxLightness) / 2
        lightRange := chars.MaxLightness - chars.MinLightness
        if lightRange > 0 {
            lightDiff := math.Abs(hsla.L - idealLight)
            lightScore := 1.0 - (lightDiff / lightRange)
            score += lightScore * weights.Lightness
        } else {
            score += weights.Lightness // Perfect fit if range is zero
        }
    }
    
    // Hue alignment score (if hue constraints specified)
    if weights.HueAlignment > 0 && chars.HueCenter != nil && chars.HueTolerance != nil {
        hueDiff := math.Abs(hsla.H - *chars.HueCenter)
        if hueDiff > 180 {
            hueDiff = 360 - hueDiff
        }
        if hueDiff <= *chars.HueTolerance {
            hueScore := 1.0 - (hueDiff / *chars.HueTolerance)
            score += hueScore * weights.HueAlignment
        } else {
            // Outside hue tolerance disqualifies for hue-specific categories
            return 0
        }
    }
    
    return score
}
```

### Step 6: Update Processor Structure

**File: `pkg/processor/processor.go` (updated)**

```go
package processor

import (
    "fmt"
    "image"
    "image/color"
    "github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
    "github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

type ThemeMode string

const (
    Light ThemeMode = "Light"
    Dark  ThemeMode = "Dark"
)

type ColorProfile struct {
    Mode            ThemeMode
    ColorScheme     chromatic.ColorScheme
    IsGrayscale     bool
    IsMonochromatic bool
    DominantHue     float64
    HueVariance     float64
    AvgLuminance    float64
    AvgSaturation   float64
    Colors          ImageColors  // Now contains category-based extraction
}

// ImageColors now focuses on category-based organization
type ImageColors struct {
    // All colors found in image with frequency
    ColorFrequency map[color.RGBA]uint32 `json:"color_frequency"`
    
    // Best color per category
    Categories map[ColorCategory]color.RGBA `json:"categories"`
    
    // Multiple candidates per category for palette generation flexibility
    CategoryCandidates map[ColorCategory][]ColorCandidate `json:"category_candidates"`
    
    // Statistical metadata
    TotalPixels   uint32  `json:"total_pixels"`
    UniqueColors  int     `json:"unique_colors"`
    CoverageRatio float64 `json:"coverage_ratio"` // % of categories filled
}

// ColorCandidate tracks a potential color for a category
// Note: HSL data removed - use formats.RGBAToHSLA() for conversion as needed
type ColorCandidate struct {
    Color     color.RGBA `json:"color"`
    Frequency uint32     `json:"frequency"`
    Score     float64    `json:"score"`
}

type Processor struct {
    settings *settings.Settings
    chroma   *chromatic.Chroma
}

func New(s *settings.Settings) *Processor {
    return &Processor{
        settings: s,
        chroma:   chromatic.NewChroma(s),
    }
}

func (p *Processor) ProcessImage(img image.Image) (*ColorProfile, error) {
    bounds := img.Bounds()
    colorFreq := make(map[color.RGBA]uint32)
    totalPixels := uint32(bounds.Dx() * bounds.Dy())
    
    // Extract color frequencies
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
            colorFreq[rgba]++
        }
    }
    
    // Filter by minimum frequency
    minCount := uint32(float64(totalPixels) * p.settings.MinFrequency)
    filtered := make(map[color.RGBA]uint32)
    for c, count := range colorFreq {
        if count >= minCount {
            filtered[c] = count
        }
    }
    
    if len(filtered) == 0 {
        return nil, fmt.Errorf("no significant colors found")
    }
    
    // Analyze color characteristics
    profile := p.analyzeColors(filtered)
    
    // NEW: Extract colors by category instead of by role
    profile.Colors = *p.extractByCategory(filtered, profile, totalPixels)
    
    return profile, nil
}

// calculateTotalPixels sums all frequency counts
func (p *Processor) calculateTotalPixels(colorFreq map[color.RGBA]uint32) uint32 {
    total := uint32(0)
    for _, freq := range colorFreq {
        total += freq
    }
    return total
}
```

### Step 7: Implement Category-Based Extraction

**File: `pkg/processor/extraction.go` (completely rewritten)**

```go
package processor

import (
    "image/color"
    "sort"
    "github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

// extractByCategory performs category-based color extraction
func (p *Processor) extractByCategory(
    colorFreq map[color.RGBA]uint32,
    profile *ColorProfile,
    totalPixels uint32,
) *ImageColors {
    
    result := &ImageColors{
        ColorFrequency:     colorFreq,
        Categories:         make(map[ColorCategory]color.RGBA),
        CategoryCandidates: make(map[ColorCategory][]ColorCandidate),
        TotalPixels:        totalPixels,
        UniqueColors:       len(colorFreq),
    }
    
    // First, establish background color (critical for contrast calculations)
    background := p.selectBackground(colorFreq, profile)
    result.Categories[CategoryBackground] = background
    
    // Categorize all colors based on their characteristics
    p.categorizeColors(colorFreq, profile, background, result)
    
    // Select best candidates for each category
    p.selectBestCategoryColors(result)
    
    // Calculate coverage ratio
    allCategories := GetAllCategories()
    filledCount := len(result.Categories)
    result.CoverageRatio = float64(filledCount) / float64(len(allCategories))
    
    return result
}

// selectBackground chooses the most appropriate background color
func (p *Processor) selectBackground(
    colorFreq map[color.RGBA]uint32,
    profile *ColorProfile,
) color.RGBA {
    
    chars := p.GetCategoryCharacteristics(CategoryBackground, profile)
    
    // Find the most frequent color that fits background requirements
    var bestColor color.RGBA
    bestScore := -1.0
    
    for c, freq := range colorFreq {
        hsla := formats.RGBAToHSLA(c)
        
        // Check if color fits background characteristics
        if hsla.L >= chars.MinLightness && hsla.L <= chars.MaxLightness &&
           hsla.S >= chars.MinSaturation && hsla.S <= chars.MaxSaturation {
            
            // Score based on frequency and proximity to ideal lightness
            freqScore := float64(freq)
            idealLight := (chars.MinLightness + chars.MaxLightness) / 2
            lightScore := 1.0 - math.Abs(hsla.L - idealLight)
            
            score := freqScore * lightScore
            
            if score > bestScore {
                bestScore = score
                bestColor = c
            }
        }
    }
    
    // Use fallback if no suitable background found
    if bestScore < 0 {
        if profile.Mode == Light {
            if rgba, err := formats.ParseHex(p.settings.LightBackgroundFallback); err == nil {
                return rgba
            }
            return color.RGBA{R: 255, G: 255, B: 255, A: 255}
        }
        
        if rgba, err := formats.ParseHex(p.settings.DarkBackgroundFallback); err == nil {
            return rgba
        }
        return color.RGBA{R: 32, G: 32, B: 32, A: 255}
    }
    
    return bestColor
}

// categorizeColors evaluates all colors for category fitness
func (p *Processor) categorizeColors(
    colorFreq map[color.RGBA]uint32,
    profile *ColorProfile,
    background color.RGBA,
    result *ImageColors,
) {
    
    categoryOrder := p.GetCategoryPriorityOrder(profile)
    maxCandidates := p.settings.Extraction.MaxCandidatesPerCategory
    
    // Track which colors have been assigned to primary categories
    assigned := make(map[color.RGBA]bool)
    assigned[background] = true
    
    for _, category := range categoryOrder {
        if category == CategoryBackground {
            continue // Already handled
        }
        
        candidates := []ColorCandidate{}
        
        for c, freq := range colorFreq {
            // Skip if already assigned to a higher-priority category
            if assigned[c] && p.isExclusiveCategory(category) {
                continue
            }
            
            // Check if color fits this category
            if p.fitsCategory(c, category, profile, background) {
                score := p.calculateCategoryFitScore(
                    c, category, profile, background, freq, result.TotalPixels,
                )
                
                if score > 0 {
                    candidates = append(candidates, ColorCandidate{
                        Color:     c,
                        Frequency: freq,
                        Score:     score,
                    })
                }
            }
        }
        
        // Sort candidates by score (highest first)
        sort.Slice(candidates, func(i, j int) bool {
            return candidates[i].Score > candidates[j].Score
        })
        
        // Keep top N candidates
        if len(candidates) > maxCandidates {
            candidates = candidates[:maxCandidates]
        }
        
        if len(candidates) > 0 {
            result.CategoryCandidates[category] = candidates
            // Mark the best candidate as assigned for exclusive categories
            if p.isExclusiveCategory(category) {
                assigned[candidates[0].Color] = true
            }
        }
    }
}

// selectBestCategoryColors picks the top candidate for each category
func (p *Processor) selectBestCategoryColors(result *ImageColors) {
    for category, candidates := range result.CategoryCandidates {
        if len(candidates) > 0 {
            result.Categories[category] = candidates[0].Color
        }
    }
}

// isExclusiveCategory returns true for categories that shouldn't share colors
func (p *Processor) isExclusiveCategory(category ColorCategory) bool {
    switch category {
    case CategoryForeground, CategoryAccentPrimary, CategoryAccentSecondary:
        return true
    default:
        return false
    }
}
```

## Testing Considerations

The enhanced processor requires comprehensive testing to validate:

1. **Category assignment accuracy**
2. **Scoring algorithm correctness**
3. **Settings integration**
4. **Edge cases (grayscale, monochromatic, etc.)**

## Usage Example

```go
// Initialize with custom settings
s := settings.DefaultSettings()
// Users can modify category characteristics
s.Categories.Dark["background"].MaxSaturation = 0.3 // Allow more colored backgrounds

p := processor.New(s)
profile, err := p.ProcessImage(img)

// Access categorized colors
bgColor := profile.Colors.Categories[processor.CategoryBackground]
fgColor := profile.Colors.Categories[processor.CategoryForeground]

// Check category coverage
fmt.Printf("Category coverage: %.1f%%\n", profile.Colors.CoverageRatio * 100)

// Access candidates for palette generation
if candidates, ok := profile.Colors.CategoryCandidates[processor.CategoryAccentPrimary]; ok {
    fmt.Printf("Found %d accent color candidates\n", len(candidates))
}
```

---

## Documentation and Testing Requirements for Claude Code

### Documentation Updates Required

1. **pkg/processor/docs.go** - Update package documentation to reflect category-based extraction
2. **docs/architecture.md** - Document the new category system and scoring algorithm
3. **CLAUDE.md** - Update with new extraction approach and remove backwards compatibility notes
4. **PROJECT.md** - Mark processor enhancement as complete
5. **README.md** - Update usage examples to show category-based extraction

### Test Files to Create/Update

1. **tests/processor/categories_test.go** - Test category system
   - Test `GetCategoryCharacteristics` with different modes
   - Test `fitsCategory` with various color/category combinations
   - Test priority ordering

2. **tests/processor/scoring_test.go** - Test scoring system
   - Test score calculation with different weights
   - Test edge cases (zero weights, perfect matches)
   - Test disqualification conditions

3. **tests/processor/extraction_test.go** - Update extraction tests
   - Remove role-based tests
   - Add category-based extraction tests
   - Test candidate selection

4. **tests/settings/categories_test.go** - Test category settings
   - Test default loading
   - Test custom settings override
   - Test validation of characteristics

### Comment Updates

All files should have updated comments reflecting:
- Category-based approach instead of role-based
- Settings-driven characteristics
- No hard-coded values
- Clear explanation of scoring algorithm

### Tool Updates

1. **tools/analyze-images/main.go** - Update to show category coverage
   ```go
   // Add category analysis section
   fmt.Printf("### Category Coverage: %.1f%%\n\n", profile.Colors.CoverageRatio * 100)
   ```

2. **tools/performance-test/main.go** - No changes needed (processing time should remain similar)

### Claude Artifacts to Update

The CLAUDE.md file should be updated to reflect:
- Remove "backwards compatibility" mentions
- Update ColorProfile composition to show new ImageColors structure
- Document category system as primary extraction method
- Remove references to "Primary", "Secondary", "Accent" fields in ImageColors

### Integration Points

For pkg/palette (future development):
- Document how palette generator will consume CategoryCandidates
- Explain how empty categories will be filled using color theory
- Note that frequency data is preserved for weighting decisions

This enhancement provides a solid foundation for sophisticated theme generation while maintaining clean architecture and full configurability.
