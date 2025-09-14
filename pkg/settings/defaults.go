package settings

import (
	"fmt"

	"github.com/spf13/viper"
	_ "golang.org/x/image/webp"
)

func setDefaults(v *viper.Viper) {
	// Foundation layer settings - anchored to lowest-level package ownership

	// Loader settings
	v.SetDefault("loader.max_width", 8192)
	v.SetDefault("loader.max_height", 8192)
	v.SetDefault("loader.allowed_formats", []string{
		"jpeg",
		"jpg",
		"png",
		"webp",
	})

	// Formats settings
	v.SetDefault("formats.quantization_bits", 5) // 32 levels per channel

	// Chromatic settings
	v.SetDefault("chromatic.color_merge_threshold", 15.0)       // Delta-E threshold for color similarity
	v.SetDefault("chromatic.neutral_threshold", 0.1)           // 10% saturation threshold for neutrals
	v.SetDefault("chromatic.neutral_lightness_threshold", 0.08) // 8% lightness difference for neutral clustering
	v.SetDefault("chromatic.dark_lightness_max", 0.3)           // 30% maximum lightness for dark classification
	v.SetDefault("chromatic.light_lightness_min", 0.7)          // 70% minimum lightness for light classification
	v.SetDefault("chromatic.muted_saturation_max", 0.3)         // 30% maximum saturation for muted classification
	v.SetDefault("chromatic.vibrant_saturation_min", 0.7)       // 70% minimum saturation for vibrant classification

	// Processing layer settings
	v.SetDefault("processor.min_frequency", 0.0001)              // 0.01% minimum frequency
	v.SetDefault("processor.min_cluster_weight", 0.005)          // 0.5% minimum cluster weight
	v.SetDefault("processor.min_ui_color_weight", 0.01)          // 1% minimum for UI inclusion
	v.SetDefault("processor.max_ui_colors", 20)                  // Maximum colors for UI palette
	v.SetDefault("processor.pure_black_threshold", 0.01)         // 1% lightness threshold for pure black
	v.SetDefault("processor.pure_white_threshold", 0.99)         // 99% lightness threshold for pure white
	v.SetDefault("processor.light_theme_threshold", 0.5)         // 50% lightness threshold for light theme
	v.SetDefault("processor.theme_mode_max_clusters", 5)         // Maximum clusters to consider for theme mode
	v.SetDefault("processor.significant_color_threshold", 0.1)   // 10% weight threshold for significant color content

	// Global settings
	v.SetDefault("default_dark", "#1a1a1a")
	v.SetDefault("default_light", "#f0f0f0")
	v.SetDefault("default_gray", "#808080")
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
