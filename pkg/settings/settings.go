package settings

import "context"

type contextKey string

const settingsKey contextKey = "settings"

type Settings struct {
	// Foundation layer settings (anchored to lowest-level package ownership)
	Loader    LoaderSettings    `mapstructure:"loader"`
	Formats   FormatsSettings   `mapstructure:"formats"`
	Chromatic ChromaticSettings `mapstructure:"chromatic"`

	// Processing layer settings
	Processor ProcessorSettings `mapstructure:"processor"`

	// Global settings
	DefaultDark  string `mapstructure:"default_dark"`  // Fallback dark color
	DefaultLight string `mapstructure:"default_light"` // Fallback light color
	DefaultGray  string `mapstructure:"default_gray"`  // Fallback gray color
}

type LoaderSettings struct {
	MaxWidth       int      `mapstructure:"max_width"`       // Maximum image width
	MaxHeight      int      `mapstructure:"max_height"`      // Maximum image height
	AllowedFormats []string `mapstructure:"allowed_formats"` // Supported image formats
}

type FormatsSettings struct {
	QuantizationBits int `mapstructure:"quantization_bits"` // Color precision (1-8 bits per channel)
}

type ChromaticSettings struct {
	ColorMergeThreshold       float64 `mapstructure:"color_merge_threshold"`       // LAB distance for color similarity
	NeutralThreshold          float64 `mapstructure:"neutral_threshold"`          // Saturation threshold for neutral colors
	NeutralLightnessThreshold float64 `mapstructure:"neutral_lightness_threshold"` // Lightness difference threshold for neutral clustering
	DarkLightnessMax          float64 `mapstructure:"dark_lightness_max"`          // Maximum lightness for dark classification
	LightLightnessMin         float64 `mapstructure:"light_lightness_min"`         // Minimum lightness for light classification
	MutedSaturationMax        float64 `mapstructure:"muted_saturation_max"`        // Maximum saturation for muted classification
	VibrantSaturationMin      float64 `mapstructure:"vibrant_saturation_min"`      // Minimum saturation for vibrant classification
}

type ProcessorSettings struct {
	// Color extraction
	MinFrequency float64 `mapstructure:"min_frequency"` // Minimum frequency to consider

	// Clustering
	MinClusterWeight float64 `mapstructure:"min_cluster_weight"` // Minimum weight to keep cluster

	// UI filtering
	MinUIColorWeight         float64 `mapstructure:"min_ui_color_weight"`         // Minimum weight for UI inclusion
	MaxUIColors              int     `mapstructure:"max_ui_colors"`               // Maximum colors for UI palette
	PureBlackThreshold       float64 `mapstructure:"pure_black_threshold"`        // Lightness threshold for pure black
	PureWhiteThreshold       float64 `mapstructure:"pure_white_threshold"`        // Lightness threshold for pure white

	// Theme analysis
	LightThemeThreshold       float64 `mapstructure:"light_theme_threshold"`       // Lightness threshold for light theme
	ThemeModeMaxClusters      int     `mapstructure:"theme_mode_max_clusters"`     // Maximum clusters to consider for theme mode
	SignificantColorThreshold float64 `mapstructure:"significant_color_threshold"` // Weight threshold for significant color content
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
