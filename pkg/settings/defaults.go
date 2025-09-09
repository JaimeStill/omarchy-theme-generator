package settings

import (
	"fmt"

	"github.com/spf13/viper"
	_ "golang.org/x/image/webp"
)

func setDefaults(v *viper.Viper) {
	v.SetDefault("grayscale_threshold", 0.05)
	v.SetDefault("grayscale_image_threshold", 0.8)
	v.SetDefault("monochromatic_tolerance", 15.0)
	v.SetDefault("monochromatic_weight_threshold", 0.1)
	v.SetDefault("theme_mode_threshold", 0.5)
	v.SetDefault("min_frequency", 0.0001)

	v.SetDefault("loader_max_width", 8192)
	v.SetDefault("loader_max_height", 8192)
	v.SetDefault("loader_allowed_formats", []string{
		"jpeg",
		"jpg",
		"png",
		"webp",
	})

	v.SetDefault("lightness_dark_max", 0.25)
	v.SetDefault("lightness_light_min", 0.75)

	v.SetDefault("saturation_gray_max", 0.05)
	v.SetDefault("saturation_muted_max", 0.25)
	v.SetDefault("saturation_normal_max", 0.70)

	v.SetDefault("hue_sector_count", 12)
	v.SetDefault("hue_sector_size", 30.0)

	v.SetDefault("extraction.max_colors_to_extract", 100000)
	v.SetDefault("extraction.dominant_color_count", 10)
	v.SetDefault("extraction.min_color_diversity", 0.1)
	v.SetDefault("extraction.adaptive_grouping", true)
	v.SetDefault("extraction.preserve_natural_clusters", true)

	v.SetDefault("fallbacks.default_dark", "#1a1a1a")
	v.SetDefault("fallbacks.default_light", "#f0f0f0")
	v.SetDefault("fallbacks.default_gray", "#808080")
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
