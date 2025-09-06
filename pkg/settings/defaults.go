package settings

import (
	"fmt"

	"github.com/spf13/viper"
	_ "golang.org/x/image/webp"
)

func setDefaults(v *viper.Viper) {
	v.SetDefault("grayscale_threshold", 0.05)
	v.SetDefault("monochromatic_tolerance", 15.0)
	v.SetDefault("theme_mode_threshold", 0.5)

	v.SetDefault("min_frequency", 0.001)

	v.SetDefault("light_background_threshold", 0.7)
	v.SetDefault("dark_background_threshold", 0.3)

	v.SetDefault("min_contrast_ratio", 4.5)

	v.SetDefault("min_primary_saturation", 0.3)
	v.SetDefault("min_accent_saturation", 0.5)
	v.SetDefault("min_accent_lightness", 0.3)
	v.SetDefault("max_accent_lightness", 0.7)

	v.SetDefault("dark_light_threshold", 0.1)
	v.SetDefault("bright_light_threshold", 0.9)
	v.SetDefault("extreme_lightness_penalty", 0.3)
	v.SetDefault("optimal_lightness_bonus", 1.2)
	v.SetDefault("min_saturation_for_bonus", 0.05)

	v.SetDefault("loader_max_width", 8192)
	v.SetDefault("loader_max_height", 8192)
	v.SetDefault("loader_allowed_formats", []string{
		"jpeg",
		"jpg",
		"png",
		"webp",
	})

	v.SetDefault("light_background_fallback", "#ffffff")
	v.SetDefault("dark_background_fallback", "#202020")
	v.SetDefault("light_foreground_fallback", "#202020")
	v.SetDefault("dark_foreground_fallback", "#ffffff")
	v.SetDefault("primary_fallback", "#6496c8")
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
