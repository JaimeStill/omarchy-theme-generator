package settings

import "context"

type contextKey string

const settingsKey contextKey = "settings"

type Settings struct {
	// Core extraction settings
	GrayscaleThreshold     float64 `mapstructure:"grayscale_threshold"`
	MonochromaticTolerance float64 `mapstructure:"monochromatic_tolerance"`
	ThemeModeThreshold     float64 `mapstructure:"theme_mode_threshold"`
	MinFrequency           float64 `mapstructure:"min_frequency"`

	// Loader settings
	LoaderMaxWidth       int      `mapstructure:"loader_max_width"`
	LoaderMaxHeight      int      `mapstructure:"loader_max_height"`
	LoaderAllowedFormats []string `mapstructure:"loader_allowed_formats"`

	// Fallback colors
	LightBackgroundFallback string `mapstructure:"light_background_fallback"`
	DarkBackgroundFallback  string `mapstructure:"dark_background_fallback"`
	LightForegroundFallback string `mapstructure:"light_foreground_fallback"`
	DarkForegroundFallback  string `mapstructure:"dark_foreground_fallback"`
	PrimaryFallback         string `mapstructure:"primary_fallback"`

	// Category-based extraction settings
	Categories      CategorySettings       `mapstructure:"categories"`
	CategoryScoring CategoryScoringWeights `mapstructure:"category_scoring"`
	Extraction      ExtractionSettings     `mapstructure:"extraction"`
}

type CategorySettings struct {
	Dark  map[string]CategoryCharacteristics `mapstructure:"dark"`
	Light map[string]CategoryCharacteristics `mapstructure:"light"`
}

type CategoryCharacteristics struct {
	MinLightness  float64  `mapstructure:"min_lightness"`
	MaxLightness  float64  `mapstructure:"max_lightness"`
	MinSaturation float64  `mapstructure:"min_saturation"`
	MaxSaturation float64  `mapstructure:"max_saturation"`
	MinContrast   float64  `mapstructure:"min_contrast"`
	HueCenter     *float64 `mapstructure:"hue_center"`
	HueTolerance  *float64 `mapstructure:"hue_tolerance"`
}

type CategoryScoringWeights struct {
	Frequency    float64 `mapstructure:"frequency"`
	Contrast     float64 `mapstructure:"contrast"`
	Saturation   float64 `mapstructure:"saturation"`
	HueAlignment float64 `mapstructure:"hue_alignment"`
	Lightness    float64 `mapstructure:"lightness"`
}

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
