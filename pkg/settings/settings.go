package settings

import "context"

type contextKey string

const settingsKey contextKey = "settings"

type Settings struct {
	// Core extraction settings
	GrayscaleThreshold           float64 `mapstructure:"grayscale_threshold"`
	GrayscaleImageThreshold      float64 `mapstructure:"grayscale_image_threshold"`
	MonochromaticTolerance       float64 `mapstructure:"monochromatic_tolerance"`
	MonochromaticWeightThreshold float64 `mapstructure:"monochromatic_weight_threshold"`
	ThemeModeThreshold           float64 `mapstructure:"theme_mode_threshold"`
	MinFrequency                 float64 `mapstructure:"min_frequency"`

	// Loader settings
	LoaderMaxWidth       int      `mapstructure:"loader_max_width"`
	LoaderMaxHeight      int      `mapstructure:"loader_max_height"`
	LoaderAllowedFormats []string `mapstructure:"loader_allowed_formats"`

	LightnessDarkMax  float64 `mapstructure:"lightness_dark_max"`
	LightnessLightMin float64 `mapstructure:"lightness_light_min"`

	SaturationGrayMax   float64 `mapstructure:"saturation_gray_max"`
	SaturationMutedMax  float64 `mapstructure:"saturation_muted_max"`
	SaturationNormalMax float64 `mapstructure:"saturation_normal_max"`

	HueSectorCount int     `mapstructure:"hue_sector_count"`
	HueSectorSize  float64 `mapstructure:"hue_sector_size"`

	Extraction ExtractionSettings `mapstructure:"extraction"`
	Fallbacks  FallbackSettings   `mapstructure:"fallbacks"`
}

type ExtractionSettings struct {
	MaxColorsToExtract      int     `mapstructure:"max_colors_to_extract"`
	DominantColorCount      int     `mapstructure:"dominant_color_count"`
	MinColorDiversity       float64 `mapstructure:"min_color_diversity"`
	AdaptiveGrouping        bool    `mapstructure:"adaptive_grouping"`
	PreserveNaturalClusters bool    `mapstructure:"preserve_natural_clusters"`
}

type FallbackSettings struct {
	DefaultDark  string `mapstructure:"default_dark"`
	DefaultLight string `mapstructure:"default_light"`
	DefaultGray  string `mapstructure:"default_gray"`
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
