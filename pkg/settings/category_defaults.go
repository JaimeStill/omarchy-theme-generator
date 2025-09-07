package settings

import "github.com/spf13/viper"

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

	// Core UI elements - relaxed ranges for better coverage
	v.SetDefault(prefix+"background.min_lightness", 0.0)
	v.SetDefault(prefix+"background.max_lightness", 0.25) // Increased from 0.15
	v.SetDefault(prefix+"background.min_saturation", 0.0)
	v.SetDefault(prefix+"background.max_saturation", 0.4) // Increased from 0.2
	v.SetDefault(prefix+"background.min_contrast", 0.0)

	v.SetDefault(prefix+"foreground.min_lightness", 0.70) // Decreased from 0.85
	v.SetDefault(prefix+"foreground.max_lightness", 1.0)
	v.SetDefault(prefix+"foreground.min_saturation", 0.0)
	v.SetDefault(prefix+"foreground.max_saturation", 0.3) // Increased from 0.1
	v.SetDefault(prefix+"foreground.min_contrast", 3.0)   // Decreased from 4.5

	v.SetDefault(prefix+"dim_foreground.min_lightness", 0.35) // Decreased from 0.4
	v.SetDefault(prefix+"dim_foreground.max_lightness", 0.65) // Increased from 0.6
	v.SetDefault(prefix+"dim_foreground.min_saturation", 0.0)
	v.SetDefault(prefix+"dim_foreground.max_saturation", 0.5) // Increased from 0.3
	v.SetDefault(prefix+"dim_foreground.min_contrast", 2.0)   // Decreased from 3.0

	v.SetDefault(prefix+"cursor.min_lightness", 0.6) // Decreased from 0.7
	v.SetDefault(prefix+"cursor.max_lightness", 1.0)
	v.SetDefault(prefix+"cursor.min_saturation", 0.0)
	v.SetDefault(prefix+"cursor.max_saturation", 0.7) // Increased from 0.5
	v.SetDefault(prefix+"cursor.min_contrast", 4.5)   // Decreased from 7.0

	// Terminal normal colors - wider ranges and lower requirements
	v.SetDefault(prefix+"normal_black.min_lightness", 0.10) // Decreased from 0.15
	v.SetDefault(prefix+"normal_black.max_lightness", 0.30) // Increased from 0.25
	v.SetDefault(prefix+"normal_black.min_saturation", 0.0)
	v.SetDefault(prefix+"normal_black.max_saturation", 0.2) // Increased from 0.1
	v.SetDefault(prefix+"normal_black.min_contrast", 1.2)   // Decreased from 1.5

	v.SetDefault(prefix+"normal_red.min_lightness", 0.25) // Decreased from 0.35
	v.SetDefault(prefix+"normal_red.max_lightness", 0.65) // Increased from 0.55
	v.SetDefault(prefix+"normal_red.min_saturation", 0.3) // Decreased from 0.6
	v.SetDefault(prefix+"normal_red.max_saturation", 1.0)
	v.SetDefault(prefix+"normal_red.min_contrast", 1.5) // Decreased from 2.0
	v.SetDefault(prefix+"normal_red.hue_center", 0.0)
	v.SetDefault(prefix+"normal_red.hue_tolerance", 35.0) // Increased from 25.0

	v.SetDefault(prefix+"normal_green.min_lightness", 0.25) // Decreased from 0.35
	v.SetDefault(prefix+"normal_green.max_lightness", 0.65) // Increased from 0.55
	v.SetDefault(prefix+"normal_green.min_saturation", 0.3) // Decreased from 0.5
	v.SetDefault(prefix+"normal_green.max_saturation", 1.0)
	v.SetDefault(prefix+"normal_green.min_contrast", 1.5) // Decreased from 2.0
	v.SetDefault(prefix+"normal_green.hue_center", 120.0)
	v.SetDefault(prefix+"normal_green.hue_tolerance", 50.0) // Increased from 40.0

	v.SetDefault(prefix+"normal_yellow.min_lightness", 0.35) // Decreased from 0.45
	v.SetDefault(prefix+"normal_yellow.max_lightness", 0.75) // Increased from 0.65
	v.SetDefault(prefix+"normal_yellow.min_saturation", 0.4) // Decreased from 0.7
	v.SetDefault(prefix+"normal_yellow.max_saturation", 1.0)
	v.SetDefault(prefix+"normal_yellow.min_contrast", 2.0) // Decreased from 2.5
	v.SetDefault(prefix+"normal_yellow.hue_center", 60.0)
	v.SetDefault(prefix+"normal_yellow.hue_tolerance", 30.0) // Increased from 20.0

	v.SetDefault(prefix+"normal_blue.min_lightness", 0.25) // Decreased from 0.35
	v.SetDefault(prefix+"normal_blue.max_lightness", 0.65) // Increased from 0.55
	v.SetDefault(prefix+"normal_blue.min_saturation", 0.3) // Decreased from 0.6
	v.SetDefault(prefix+"normal_blue.max_saturation", 1.0)
	v.SetDefault(prefix+"normal_blue.min_contrast", 1.5) // Decreased from 2.0
	v.SetDefault(prefix+"normal_blue.hue_center", 240.0)
	v.SetDefault(prefix+"normal_blue.hue_tolerance", 40.0) // Increased from 30.0

	v.SetDefault(prefix+"normal_magenta.min_lightness", 0.25) // Decreased from 0.35
	v.SetDefault(prefix+"normal_magenta.max_lightness", 0.65) // Increased from 0.55
	v.SetDefault(prefix+"normal_magenta.min_saturation", 0.3) // Decreased from 0.5
	v.SetDefault(prefix+"normal_magenta.max_saturation", 1.0)
	v.SetDefault(prefix+"normal_magenta.min_contrast", 1.5) // Decreased from 2.0
	v.SetDefault(prefix+"normal_magenta.hue_center", 300.0)
	v.SetDefault(prefix+"normal_magenta.hue_tolerance", 40.0) // Increased from 30.0

	v.SetDefault(prefix+"normal_cyan.min_lightness", 0.25) // Decreased from 0.35
	v.SetDefault(prefix+"normal_cyan.max_lightness", 0.65) // Increased from 0.55
	v.SetDefault(prefix+"normal_cyan.min_saturation", 0.3) // Decreased from 0.5
	v.SetDefault(prefix+"normal_cyan.max_saturation", 1.0)
	v.SetDefault(prefix+"normal_cyan.min_contrast", 1.5) // Decreased from 2.0
	v.SetDefault(prefix+"normal_cyan.hue_center", 180.0)
	v.SetDefault(prefix+"normal_cyan.hue_tolerance", 40.0) // Increased from 30.0

	v.SetDefault(prefix+"normal_white.min_lightness", 0.5)  // Decreased from 0.6
	v.SetDefault(prefix+"normal_white.max_lightness", 0.85) // Increased from 0.8
	v.SetDefault(prefix+"normal_white.min_saturation", 0.0)
	v.SetDefault(prefix+"normal_white.max_saturation", 0.2) // Increased from 0.1
	v.SetDefault(prefix+"normal_white.min_contrast", 3.0)   // Decreased from 3.5

	// Terminal bright colors - wider ranges for better capture
	v.SetDefault(prefix+"bright_black.min_lightness", 0.20) // Decreased from 0.25
	v.SetDefault(prefix+"bright_black.max_lightness", 0.40) // Increased from 0.35
	v.SetDefault(prefix+"bright_black.min_saturation", 0.0)
	v.SetDefault(prefix+"bright_black.max_saturation", 0.2) // Increased from 0.1
	v.SetDefault(prefix+"bright_black.min_contrast", 1.5)   // Decreased from 2.0

	v.SetDefault(prefix+"bright_red.min_lightness", 0.4)  // Decreased from 0.5
	v.SetDefault(prefix+"bright_red.max_lightness", 0.8)  // Increased from 0.7
	v.SetDefault(prefix+"bright_red.min_saturation", 0.5) // Decreased from 0.8
	v.SetDefault(prefix+"bright_red.max_saturation", 1.0)
	v.SetDefault(prefix+"bright_red.min_contrast", 2.5) // Decreased from 3.0
	v.SetDefault(prefix+"bright_red.hue_center", 0.0)
	v.SetDefault(prefix+"bright_red.hue_tolerance", 35.0) // Increased from 25.0

	v.SetDefault(prefix+"bright_green.min_lightness", 0.4)  // Decreased from 0.5
	v.SetDefault(prefix+"bright_green.max_lightness", 0.8)  // Increased from 0.7
	v.SetDefault(prefix+"bright_green.min_saturation", 0.4) // Decreased from 0.7
	v.SetDefault(prefix+"bright_green.max_saturation", 1.0)
	v.SetDefault(prefix+"bright_green.min_contrast", 2.5) // Decreased from 3.0
	v.SetDefault(prefix+"bright_green.hue_center", 120.0)
	v.SetDefault(prefix+"bright_green.hue_tolerance", 50.0) // Increased from 40.0

	v.SetDefault(prefix+"bright_yellow.min_lightness", 0.5)  // Decreased from 0.6
	v.SetDefault(prefix+"bright_yellow.max_lightness", 0.9)  // Increased from 0.8
	v.SetDefault(prefix+"bright_yellow.min_saturation", 0.5) // Decreased from 0.8
	v.SetDefault(prefix+"bright_yellow.max_saturation", 1.0)
	v.SetDefault(prefix+"bright_yellow.min_contrast", 3.0) // Decreased from 3.5
	v.SetDefault(prefix+"bright_yellow.hue_center", 60.0)
	v.SetDefault(prefix+"bright_yellow.hue_tolerance", 30.0) // Increased from 20.0

	v.SetDefault(prefix+"bright_blue.min_lightness", 0.4)  // Decreased from 0.5
	v.SetDefault(prefix+"bright_blue.max_lightness", 0.8)  // Increased from 0.7
	v.SetDefault(prefix+"bright_blue.min_saturation", 0.5) // Decreased from 0.8
	v.SetDefault(prefix+"bright_blue.max_saturation", 1.0)
	v.SetDefault(prefix+"bright_blue.min_contrast", 2.5) // Decreased from 3.0
	v.SetDefault(prefix+"bright_blue.hue_center", 240.0)
	v.SetDefault(prefix+"bright_blue.hue_tolerance", 40.0) // Increased from 30.0

	v.SetDefault(prefix+"bright_magenta.min_lightness", 0.4)  // Decreased from 0.5
	v.SetDefault(prefix+"bright_magenta.max_lightness", 0.8)  // Increased from 0.7
	v.SetDefault(prefix+"bright_magenta.min_saturation", 0.4) // Decreased from 0.7
	v.SetDefault(prefix+"bright_magenta.max_saturation", 1.0)
	v.SetDefault(prefix+"bright_magenta.min_contrast", 2.5) // Decreased from 3.0
	v.SetDefault(prefix+"bright_magenta.hue_center", 300.0)
	v.SetDefault(prefix+"bright_magenta.hue_tolerance", 40.0) // Increased from 30.0

	v.SetDefault(prefix+"bright_cyan.min_lightness", 0.4)  // Decreased from 0.5
	v.SetDefault(prefix+"bright_cyan.max_lightness", 0.8)  // Increased from 0.7
	v.SetDefault(prefix+"bright_cyan.min_saturation", 0.4) // Decreased from 0.7
	v.SetDefault(prefix+"bright_cyan.max_saturation", 1.0)
	v.SetDefault(prefix+"bright_cyan.min_contrast", 2.5) // Decreased from 3.0
	v.SetDefault(prefix+"bright_cyan.hue_center", 180.0)
	v.SetDefault(prefix+"bright_cyan.hue_tolerance", 40.0) // Increased from 30.0

	v.SetDefault(prefix+"bright_white.min_lightness", 0.7) // Decreased from 0.8
	v.SetDefault(prefix+"bright_white.max_lightness", 1.0)
	v.SetDefault(prefix+"bright_white.min_saturation", 0.0)
	v.SetDefault(prefix+"bright_white.max_saturation", 0.2) // Increased from 0.1
	v.SetDefault(prefix+"bright_white.min_contrast", 3.5)   // Decreased from 4.0

	// Accent colors - more inclusive ranges
	v.SetDefault(prefix+"accent_primary.min_lightness", 0.3)  // Decreased from 0.4
	v.SetDefault(prefix+"accent_primary.max_lightness", 0.8)  // Increased from 0.7
	v.SetDefault(prefix+"accent_primary.min_saturation", 0.4) // Decreased from 0.6
	v.SetDefault(prefix+"accent_primary.max_saturation", 1.0)
	v.SetDefault(prefix+"accent_primary.min_contrast", 2.5) // Decreased from 3.0

	v.SetDefault(prefix+"accent_secondary.min_lightness", 0.25)  // Decreased from 0.35
	v.SetDefault(prefix+"accent_secondary.max_lightness", 0.75)  // Increased from 0.65
	v.SetDefault(prefix+"accent_secondary.min_saturation", 0.3)  // Decreased from 0.5
	v.SetDefault(prefix+"accent_secondary.max_saturation", 0.95) // Increased from 0.9
	v.SetDefault(prefix+"accent_secondary.min_contrast", 2.0)    // Decreased from 2.5

	v.SetDefault(prefix+"accent_tertiary.min_lightness", 0.2)   // Decreased from 0.3
	v.SetDefault(prefix+"accent_tertiary.max_lightness", 0.7)   // Increased from 0.6
	v.SetDefault(prefix+"accent_tertiary.min_saturation", 0.25) // Decreased from 0.4
	v.SetDefault(prefix+"accent_tertiary.max_saturation", 0.9)  // Increased from 0.8
	v.SetDefault(prefix+"accent_tertiary.min_contrast", 1.5)    // Decreased from 2.0

	// Semantic colors - wider tolerances
	v.SetDefault(prefix+"error.min_lightness", 0.3)  // Decreased from 0.4
	v.SetDefault(prefix+"error.max_lightness", 0.7)  // Increased from 0.6
	v.SetDefault(prefix+"error.min_saturation", 0.5) // Decreased from 0.7
	v.SetDefault(prefix+"error.max_saturation", 1.0)
	v.SetDefault(prefix+"error.min_contrast", 2.5) // Decreased from 3.0
	v.SetDefault(prefix+"error.hue_center", 0.0)
	v.SetDefault(prefix+"error.hue_tolerance", 30.0) // Increased from 20.0

	v.SetDefault(prefix+"warning.min_lightness", 0.35) // Decreased from 0.45
	v.SetDefault(prefix+"warning.max_lightness", 0.75) // Increased from 0.65
	v.SetDefault(prefix+"warning.min_saturation", 0.5) // Decreased from 0.7
	v.SetDefault(prefix+"warning.max_saturation", 1.0)
	v.SetDefault(prefix+"warning.min_contrast", 2.5) // Decreased from 3.0
	v.SetDefault(prefix+"warning.hue_center", 45.0)
	v.SetDefault(prefix+"warning.hue_tolerance", 25.0) // Increased from 15.0

	v.SetDefault(prefix+"success.min_lightness", 0.25) // Decreased from 0.35
	v.SetDefault(prefix+"success.max_lightness", 0.65) // Increased from 0.55
	v.SetDefault(prefix+"success.min_saturation", 0.3) // Decreased from 0.5
	v.SetDefault(prefix+"success.max_saturation", 1.0)
	v.SetDefault(prefix+"success.min_contrast", 2.5) // Decreased from 3.0
	v.SetDefault(prefix+"success.hue_center", 120.0)
	v.SetDefault(prefix+"success.hue_tolerance", 40.0) // Increased from 30.0

	v.SetDefault(prefix+"info.min_lightness", 0.25) // Decreased from 0.35
	v.SetDefault(prefix+"info.max_lightness", 0.65) // Increased from 0.55
	v.SetDefault(prefix+"info.min_saturation", 0.3) // Decreased from 0.5
	v.SetDefault(prefix+"info.max_saturation", 1.0)
	v.SetDefault(prefix+"info.min_contrast", 2.5) // Decreased from 3.0
	v.SetDefault(prefix+"info.hue_center", 210.0)
	v.SetDefault(prefix+"info.hue_tolerance", 40.0) // Increased from 30.0
}

func setLightCategoryDefaults(v *viper.Viper) {
	prefix := "categories.light."

	// Core UI elements - relaxed for light mode
	v.SetDefault(prefix+"background.min_lightness", 0.85) // Decreased from 0.90
	v.SetDefault(prefix+"background.max_lightness", 1.0)
	v.SetDefault(prefix+"background.min_saturation", 0.0)
	v.SetDefault(prefix+"background.max_saturation", 0.25) // Increased from 0.15
	v.SetDefault(prefix+"background.min_contrast", 0.0)

	v.SetDefault(prefix+"foreground.min_lightness", 0.0)
	v.SetDefault(prefix+"foreground.max_lightness", 0.30) // Increased from 0.20
	v.SetDefault(prefix+"foreground.min_saturation", 0.0)
	v.SetDefault(prefix+"foreground.max_saturation", 0.2) // Increased from 0.1
	v.SetDefault(prefix+"foreground.min_contrast", 3.0)   // Decreased from 4.5

	v.SetDefault(prefix+"dim_foreground.min_lightness", 0.25) // Decreased from 0.3
	v.SetDefault(prefix+"dim_foreground.max_lightness", 0.55) // Increased from 0.5
	v.SetDefault(prefix+"dim_foreground.min_saturation", 0.0)
	v.SetDefault(prefix+"dim_foreground.max_saturation", 0.4) // Increased from 0.3
	v.SetDefault(prefix+"dim_foreground.min_contrast", 2.0)   // Decreased from 3.0

	v.SetDefault(prefix+"cursor.min_lightness", 0.0)
	v.SetDefault(prefix+"cursor.max_lightness", 0.4) // Increased from 0.3
	v.SetDefault(prefix+"cursor.min_saturation", 0.0)
	v.SetDefault(prefix+"cursor.max_saturation", 0.7) // Increased from 0.5
	v.SetDefault(prefix+"cursor.min_contrast", 4.5)   // Decreased from 7.0

	// Terminal normal colors - wider ranges for light mode
	v.SetDefault(prefix+"normal_black.min_lightness", 0.0)
	v.SetDefault(prefix+"normal_black.max_lightness", 0.20) // Increased from 0.15
	v.SetDefault(prefix+"normal_black.min_saturation", 0.0)
	v.SetDefault(prefix+"normal_black.max_saturation", 0.2) // Increased from 0.1
	v.SetDefault(prefix+"normal_black.min_contrast", 4.5)   // Decreased from 7.0

	v.SetDefault(prefix+"normal_red.min_lightness", 0.20) // Decreased from 0.25
	v.SetDefault(prefix+"normal_red.max_lightness", 0.50) // Increased from 0.45
	v.SetDefault(prefix+"normal_red.min_saturation", 0.4) // Decreased from 0.6
	v.SetDefault(prefix+"normal_red.max_saturation", 1.0)
	v.SetDefault(prefix+"normal_red.min_contrast", 3.0) // Decreased from 4.5
	v.SetDefault(prefix+"normal_red.hue_center", 0.0)
	v.SetDefault(prefix+"normal_red.hue_tolerance", 35.0) // Increased from 25.0

	v.SetDefault(prefix+"normal_green.min_lightness", 0.20) // Decreased from 0.25
	v.SetDefault(prefix+"normal_green.max_lightness", 0.50) // Increased from 0.45
	v.SetDefault(prefix+"normal_green.min_saturation", 0.3) // Decreased from 0.5
	v.SetDefault(prefix+"normal_green.max_saturation", 1.0)
	v.SetDefault(prefix+"normal_green.min_contrast", 3.0) // Decreased from 4.5
	v.SetDefault(prefix+"normal_green.hue_center", 120.0)
	v.SetDefault(prefix+"normal_green.hue_tolerance", 50.0) // Increased from 40.0

	v.SetDefault(prefix+"normal_yellow.min_lightness", 0.25) // Decreased from 0.3
	v.SetDefault(prefix+"normal_yellow.max_lightness", 0.55) // Increased from 0.5
	v.SetDefault(prefix+"normal_yellow.min_saturation", 0.5) // Decreased from 0.7
	v.SetDefault(prefix+"normal_yellow.max_saturation", 1.0)
	v.SetDefault(prefix+"normal_yellow.min_contrast", 3.0) // Decreased from 4.5
	v.SetDefault(prefix+"normal_yellow.hue_center", 60.0)
	v.SetDefault(prefix+"normal_yellow.hue_tolerance", 30.0) // Increased from 20.0

	v.SetDefault(prefix+"normal_blue.min_lightness", 0.20) // Decreased from 0.25
	v.SetDefault(prefix+"normal_blue.max_lightness", 0.50) // Increased from 0.45
	v.SetDefault(prefix+"normal_blue.min_saturation", 0.4) // Decreased from 0.6
	v.SetDefault(prefix+"normal_blue.max_saturation", 1.0)
	v.SetDefault(prefix+"normal_blue.min_contrast", 3.0) // Decreased from 4.5
	v.SetDefault(prefix+"normal_blue.hue_center", 240.0)
	v.SetDefault(prefix+"normal_blue.hue_tolerance", 40.0) // Increased from 30.0

	v.SetDefault(prefix+"normal_magenta.min_lightness", 0.20) // Decreased from 0.25
	v.SetDefault(prefix+"normal_magenta.max_lightness", 0.50) // Increased from 0.45
	v.SetDefault(prefix+"normal_magenta.min_saturation", 0.3) // Decreased from 0.5
	v.SetDefault(prefix+"normal_magenta.max_saturation", 1.0)
	v.SetDefault(prefix+"normal_magenta.min_contrast", 3.0) // Decreased from 4.5
	v.SetDefault(prefix+"normal_magenta.hue_center", 300.0)
	v.SetDefault(prefix+"normal_magenta.hue_tolerance", 40.0) // Increased from 30.0

	v.SetDefault(prefix+"normal_cyan.min_lightness", 0.20) // Decreased from 0.25
	v.SetDefault(prefix+"normal_cyan.max_lightness", 0.50) // Increased from 0.45
	v.SetDefault(prefix+"normal_cyan.min_saturation", 0.3) // Decreased from 0.5
	v.SetDefault(prefix+"normal_cyan.max_saturation", 1.0)
	v.SetDefault(prefix+"normal_cyan.min_contrast", 3.0) // Decreased from 4.5
	v.SetDefault(prefix+"normal_cyan.hue_center", 180.0)
	v.SetDefault(prefix+"normal_cyan.hue_tolerance", 40.0) // Increased from 30.0

	v.SetDefault(prefix+"normal_white.min_lightness", 0.35) // Decreased from 0.4
	v.SetDefault(prefix+"normal_white.max_lightness", 0.65) // Increased from 0.6
	v.SetDefault(prefix+"normal_white.min_saturation", 0.0)
	v.SetDefault(prefix+"normal_white.max_saturation", 0.2) // Increased from 0.1
	v.SetDefault(prefix+"normal_white.min_contrast", 3.0)   // Decreased from 3.5

	// Terminal bright colors - relaxed for light mode
	v.SetDefault(prefix+"bright_black.min_lightness", 0.10) // Decreased from 0.15
	v.SetDefault(prefix+"bright_black.max_lightness", 0.30) // Increased from 0.25
	v.SetDefault(prefix+"bright_black.min_saturation", 0.0)
	v.SetDefault(prefix+"bright_black.max_saturation", 0.2) // Increased from 0.1
	v.SetDefault(prefix+"bright_black.min_contrast", 3.5)   // Decreased from 5.0

	v.SetDefault(prefix+"bright_red.min_lightness", 0.10) // Decreased from 0.15
	v.SetDefault(prefix+"bright_red.max_lightness", 0.40) // Increased from 0.35
	v.SetDefault(prefix+"bright_red.min_saturation", 0.6) // Decreased from 0.8
	v.SetDefault(prefix+"bright_red.max_saturation", 1.0)
	v.SetDefault(prefix+"bright_red.min_contrast", 5.0) // Decreased from 7.0
	v.SetDefault(prefix+"bright_red.hue_center", 0.0)
	v.SetDefault(prefix+"bright_red.hue_tolerance", 35.0) // Increased from 25.0

	v.SetDefault(prefix+"bright_green.min_lightness", 0.10) // Decreased from 0.15
	v.SetDefault(prefix+"bright_green.max_lightness", 0.40) // Increased from 0.35
	v.SetDefault(prefix+"bright_green.min_saturation", 0.5) // Decreased from 0.7
	v.SetDefault(prefix+"bright_green.max_saturation", 1.0)
	v.SetDefault(prefix+"bright_green.min_contrast", 5.0) // Decreased from 7.0
	v.SetDefault(prefix+"bright_green.hue_center", 120.0)
	v.SetDefault(prefix+"bright_green.hue_tolerance", 50.0) // Increased from 40.0

	v.SetDefault(prefix+"bright_yellow.min_lightness", 0.15) // Decreased from 0.2
	v.SetDefault(prefix+"bright_yellow.max_lightness", 0.45) // Increased from 0.4
	v.SetDefault(prefix+"bright_yellow.min_saturation", 0.6) // Decreased from 0.8
	v.SetDefault(prefix+"bright_yellow.max_saturation", 1.0)
	v.SetDefault(prefix+"bright_yellow.min_contrast", 5.0) // Decreased from 7.0
	v.SetDefault(prefix+"bright_yellow.hue_center", 60.0)
	v.SetDefault(prefix+"bright_yellow.hue_tolerance", 30.0) // Increased from 20.0

	v.SetDefault(prefix+"bright_blue.min_lightness", 0.10) // Decreased from 0.15
	v.SetDefault(prefix+"bright_blue.max_lightness", 0.40) // Increased from 0.35
	v.SetDefault(prefix+"bright_blue.min_saturation", 0.6) // Decreased from 0.8
	v.SetDefault(prefix+"bright_blue.max_saturation", 1.0)
	v.SetDefault(prefix+"bright_blue.min_contrast", 5.0) // Decreased from 7.0
	v.SetDefault(prefix+"bright_blue.hue_center", 240.0)
	v.SetDefault(prefix+"bright_blue.hue_tolerance", 40.0) // Increased from 30.0

	v.SetDefault(prefix+"bright_magenta.min_lightness", 0.10) // Decreased from 0.15
	v.SetDefault(prefix+"bright_magenta.max_lightness", 0.40) // Increased from 0.35
	v.SetDefault(prefix+"bright_magenta.min_saturation", 0.5) // Decreased from 0.7
	v.SetDefault(prefix+"bright_magenta.max_saturation", 1.0)
	v.SetDefault(prefix+"bright_magenta.min_contrast", 5.0) // Decreased from 7.0
	v.SetDefault(prefix+"bright_magenta.hue_center", 300.0)
	v.SetDefault(prefix+"bright_magenta.hue_tolerance", 40.0) // Increased from 30.0

	v.SetDefault(prefix+"bright_cyan.min_lightness", 0.10) // Decreased from 0.15
	v.SetDefault(prefix+"bright_cyan.max_lightness", 0.40) // Increased from 0.35
	v.SetDefault(prefix+"bright_cyan.min_saturation", 0.5) // Decreased from 0.7
	v.SetDefault(prefix+"bright_cyan.max_saturation", 1.0)
	v.SetDefault(prefix+"bright_cyan.min_contrast", 5.0) // Decreased from 7.0
	v.SetDefault(prefix+"bright_cyan.hue_center", 180.0)
	v.SetDefault(prefix+"bright_cyan.hue_tolerance", 40.0) // Increased from 30.0

	v.SetDefault(prefix+"bright_white.min_lightness", 0.15) // Decreased from 0.2
	v.SetDefault(prefix+"bright_white.max_lightness", 0.45) // Increased from 0.4
	v.SetDefault(prefix+"bright_white.min_saturation", 0.0)
	v.SetDefault(prefix+"bright_white.max_saturation", 0.2) // Increased from 0.1
	v.SetDefault(prefix+"bright_white.min_contrast", 5.0)   // Decreased from 7.0

	// Accent colors - wider ranges for light mode
	v.SetDefault(prefix+"accent_primary.min_lightness", 0.25) // Decreased from 0.3
	v.SetDefault(prefix+"accent_primary.max_lightness", 0.65) // Increased from 0.6
	v.SetDefault(prefix+"accent_primary.min_saturation", 0.4) // Decreased from 0.6
	v.SetDefault(prefix+"accent_primary.max_saturation", 1.0)
	v.SetDefault(prefix+"accent_primary.min_contrast", 2.5) // Decreased from 3.0

	v.SetDefault(prefix+"accent_secondary.min_lightness", 0.30)  // Decreased from 0.35
	v.SetDefault(prefix+"accent_secondary.max_lightness", 0.70)  // Increased from 0.65
	v.SetDefault(prefix+"accent_secondary.min_saturation", 0.3)  // Decreased from 0.5
	v.SetDefault(prefix+"accent_secondary.max_saturation", 0.95) // Increased from 0.9
	v.SetDefault(prefix+"accent_secondary.min_contrast", 2.0)    // Decreased from 2.5

	v.SetDefault(prefix+"accent_tertiary.min_lightness", 0.35)  // Decreased from 0.4
	v.SetDefault(prefix+"accent_tertiary.max_lightness", 0.75)  // Increased from 0.7
	v.SetDefault(prefix+"accent_tertiary.min_saturation", 0.25) // Decreased from 0.4
	v.SetDefault(prefix+"accent_tertiary.max_saturation", 0.9)  // Increased from 0.8
	v.SetDefault(prefix+"accent_tertiary.min_contrast", 1.5)    // Decreased from 2.0

	// Semantic colors - relaxed for light mode
	v.SetDefault(prefix+"error.min_lightness", 0.25) // Decreased from 0.3
	v.SetDefault(prefix+"error.max_lightness", 0.55) // Increased from 0.5
	v.SetDefault(prefix+"error.min_saturation", 0.5) // Decreased from 0.7
	v.SetDefault(prefix+"error.max_saturation", 1.0)
	v.SetDefault(prefix+"error.min_contrast", 3.0) // Decreased from 4.5
	v.SetDefault(prefix+"error.hue_center", 0.0)
	v.SetDefault(prefix+"error.hue_tolerance", 30.0) // Increased from 20.0

	v.SetDefault(prefix+"warning.min_lightness", 0.30) // Decreased from 0.35
	v.SetDefault(prefix+"warning.max_lightness", 0.60) // Increased from 0.55
	v.SetDefault(prefix+"warning.min_saturation", 0.5) // Decreased from 0.7
	v.SetDefault(prefix+"warning.max_saturation", 1.0)
	v.SetDefault(prefix+"warning.min_contrast", 3.0) // Decreased from 4.5
	v.SetDefault(prefix+"warning.hue_center", 45.0)
	v.SetDefault(prefix+"warning.hue_tolerance", 25.0) // Increased from 15.0

	v.SetDefault(prefix+"success.min_lightness", 0.20) // Decreased from 0.25
	v.SetDefault(prefix+"success.max_lightness", 0.50) // Increased from 0.45
	v.SetDefault(prefix+"success.min_saturation", 0.3) // Decreased from 0.5
	v.SetDefault(prefix+"success.max_saturation", 1.0)
	v.SetDefault(prefix+"success.min_contrast", 3.0) // Decreased from 4.5
	v.SetDefault(prefix+"success.hue_center", 120.0)
	v.SetDefault(prefix+"success.hue_tolerance", 40.0) // Increased from 30.0

	v.SetDefault(prefix+"info.min_lightness", 0.20) // Decreased from 0.25
	v.SetDefault(prefix+"info.max_lightness", 0.50) // Increased from 0.45
	v.SetDefault(prefix+"info.min_saturation", 0.3) // Decreased from 0.5
	v.SetDefault(prefix+"info.max_saturation", 1.0)
	v.SetDefault(prefix+"info.min_contrast", 3.0) // Decreased from 4.5
	v.SetDefault(prefix+"info.hue_center", 210.0)
	v.SetDefault(prefix+"info.hue_tolerance", 40.0) // Increased from 30.0
}
