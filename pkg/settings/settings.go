package settings

import "context"

type contextKey string

const settingsKey contextKey = "settings"

type Settings struct {
	GrayscaleThreshold           float64  `mapstructure:"grayscale_threshold"`
	MonochromaticTolerance       float64  `mapstructure:"monochromatic_tolerance"`
	LoaderMaxWidth               int      `mapstructure:"loader_max_width"`
	LoaderMaxHeight              int      `mapstructure:"loader_max_height"`
	LoaderAllowedFormats         []string `mapstructure:"loader_allowed_formats"`
	ExtractorMaxColors           int      `mapstructure:"extractor_max_colors"`
	ExtractorMinThreshold        float64  `mapstructure:"extractor_min_threshold"`
	ExtractorEdgeThreshold       float64  `mapstructure:"extractor_edge_threshold"`
	ExtractorColorComplexity     int      `mapstructure:"extractor_color_complexity"`
	ExtractorSaturationThreshold float64  `mapstructure:"extractor_saturation_threshold"`
	ExtractorMaxCandidates       int      `mapstructure:"extractor_max_candidates"`
	ExtractorDominanceThreshold  float64  `mapstructure:"extractor_dominance_threshold"`
	ExtractorOptimalLightnessMin float64  `mapstructure:"extractor_optimal_lightness_min"`
	ExtractorOptimalLightnessMax float64  `mapstructure:"extractor_optimal_lightness_max"`
	ExtractorSpreadDivisor       float64  `mapstructure:"extractor_spread_divisor"`
	ExtractorFrequencyWeight     float64  `mapstructure:"extractor_frequency_weight"`
	ExtractorSaliencyWeight      float64  `mapstructure:"extractor_saliency_weight"`
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
