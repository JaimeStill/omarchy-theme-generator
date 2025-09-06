package settings

import "context"

type contextKey string

const settingsKey contextKey = "settings"

type Settings struct {
	GrayscaleThreshold       float64  `mapstructure:"grayscale_threshold"`
	MonochromaticTolerance   float64  `mapstructure:"monochromatic_tolerance"`
	ThemeModeThreshold       float64  `mapstructure:"theme_mode_threshold"`
	MinFrequency             float64  `mapstructure:"min_frequency"`
	LightBackgroundThreshold float64  `mapstructure:"light_background_threshold"`
	DarkBackgroundThreshold  float64  `mapstructure:"dark_background_threshold"`
	MinContrastRatio         float64  `mapstructure:"min_contrast_ratio"`
	MinPrimarySaturation     float64  `mapstructure:"min_primary_saturation"`
	MinAccentSaturation      float64  `mapstructure:"min_accent_saturation"`
	MinAccentLightness       float64  `mapstructure:"min_accent_lightness"`
	MaxAccentLightness       float64  `mapstructure:"max_accent_lightness"`
	DarkLightThreshold       float64  `mapstructure:"dark_light_threshold"`
	BrightLightThreshold     float64  `mapstructure:"bright_light_threshold"`
	ExtremeLightnessPenalty  float64  `mapstructure:"extreme_lightness_penalty"`
	OptimalLightnessBonus    float64  `mapstructure:"optimal_lightness_bonus"`
	MinSaturationForBonus    float64  `mapstructure:"min_saturation_for_bonus"`
	LoaderMaxWidth           int      `mapstructure:"loader_max_width"`
	LoaderMaxHeight          int      `mapstructure:"loader_max_height"`
	LoaderAllowedFormats     []string `mapstructure:"loader_allowed_formats"`
	LightBackgroundFallback  string   `mapstructure:"light_background_fallback"`
	DarkBackgroundFallback   string   `mapstructure:"dark_background_fallback"`
	LightForegroundFallback  string   `mapstructure:"light_foreground_fallback"`
	DarkForegroundFallback   string   `mapstructure:"dark_foreground_fallback"`
	PrimaryFallback          string   `mapstructure:"primary_fallback"`
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
