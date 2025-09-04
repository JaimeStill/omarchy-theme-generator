package settings

import (
	"fmt"

	"github.com/spf13/viper"
	_ "golang.org/x/image/webp"
)

func setDefaults(v *viper.Viper) {
	v.SetDefault("grayscale_threshold", 0.05)
	v.SetDefault("monochromatic_tolerance", 15.0)

	v.SetDefault("loader_max_width", 8192)
	v.SetDefault("loader_max_height", 8192)
	v.SetDefault("loader_allowed_formats", []string{
		"jpeg",
		"jpg",
		"png",
		"webp",
	})

	v.SetDefault("extractor_max_colors", 10)
	v.SetDefault("extractor_min_threshold", 0.1)

	v.SetDefault("extractor_edge_threshold", 0.036)
	v.SetDefault("extractor_color_complexity", 10000)
	v.SetDefault("extractor_saturation_threshold", 0.4)

	v.SetDefault("extractor_max_candidates", 20)
	v.SetDefault("extractor_dominance_threshold", 60.0)
	v.SetDefault("extractor_optimal_lightness_min", 0.2)
	v.SetDefault("extractor_optimal_lightness_max", 0.8)

	v.SetDefault("extractor_spread_divisor", 3.0)
	v.SetDefault("extractor_weight", 0.3)
	v.SetDefault("extractor_saliency_weight", 0.7)
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
