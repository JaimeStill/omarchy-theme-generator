package processor

import (
	"image/color"
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

type ColorCategory string

const (
	// Core UI elements
	CategoryBackground    ColorCategory = "background"
	CategoryForeground    ColorCategory = "foreground"
	CategoryDimForeground ColorCategory = "dim_foreground"
	CategoryCursor        ColorCategory = "cursor"

	// Terminal normal colors (ANSI 0-7)
	CategoryNormalBlack   ColorCategory = "normal_black"
	CategoryNormalRed     ColorCategory = "normal_red"
	CategoryNormalGreen   ColorCategory = "normal_green"
	CategoryNormalYellow  ColorCategory = "normal_yellow"
	CategoryNormalBlue    ColorCategory = "normal_blue"
	CategoryNormalMagenta ColorCategory = "normal_magenta"
	CategoryNormalCyan    ColorCategory = "normal_cyan"
	CategoryNormalWhite   ColorCategory = "normal_white"

	// Terminal bright colors (ANSI 8-15)
	CategoryBrightBlack   ColorCategory = "bright_black"
	CategoryBrightRed     ColorCategory = "bright_red"
	CategoryBrightGreen   ColorCategory = "bright_green"
	CategoryBrightYellow  ColorCategory = "bright_yellow"
	CategoryBrightBlue    ColorCategory = "bright_blue"
	CategoryBrightMagenta ColorCategory = "bright_magenta"
	CategoryBrightCyan    ColorCategory = "bright_cyan"
	CategoryBrightWhite   ColorCategory = "bright_white"

	// Accent colors
	CategoryAccentPrimary   ColorCategory = "accent_primary"
	CategoryAccentSecondary ColorCategory = "accent_secondary"
	CategoryAccentTertiary  ColorCategory = "accent_tertiary"

	// Semantic colors
	CategoryError   ColorCategory = "error"
	CategoryWarning ColorCategory = "warning"
	CategorySuccess ColorCategory = "success"
	CategoryInfo    ColorCategory = "info"
)

func GetAllCategories() []ColorCategory {
	return []ColorCategory{
		CategoryBackground,
		CategoryForeground,
		CategoryDimForeground,
		CategoryCursor,
		CategoryNormalBlack,
		CategoryNormalRed,
		CategoryNormalGreen,
		CategoryNormalYellow,
		CategoryNormalBlue,
		CategoryNormalMagenta,
		CategoryNormalCyan,
		CategoryNormalWhite,
		CategoryBrightBlack,
		CategoryBrightRed,
		CategoryBrightGreen,
		CategoryBrightYellow,
		CategoryBrightBlue,
		CategoryBrightMagenta,
		CategoryBrightCyan,
		CategoryBrightWhite,
		CategoryAccentPrimary,
		CategoryAccentSecondary,
		CategoryAccentTertiary,
		CategoryError,
		CategoryWarning,
		CategorySuccess,
		CategoryInfo,
	}
}

func (p *Processor) GetCategoryPriorityOrder(profile *ColorProfile) []ColorCategory {
	return []ColorCategory{
		CategoryBackground,
		CategoryForeground,
		CategoryAccentPrimary,
		CategoryDimForeground,
		CategoryCursor,
		CategoryError,
		CategoryWarning,
		CategorySuccess,
		CategoryInfo,
		CategoryAccentSecondary,
		CategoryAccentTertiary,
		CategoryNormalRed,
		CategoryNormalGreen,
		CategoryNormalBlue,
		CategoryNormalYellow,
		CategoryNormalMagenta,
		CategoryNormalCyan,
		CategoryNormalBlack,
		CategoryNormalWhite,
		CategoryBrightRed,
		CategoryBrightGreen,
		CategoryBrightBlue,
		CategoryBrightYellow,
		CategoryBrightMagenta,
		CategoryBrightCyan,
		CategoryBrightBlack,
		CategoryBrightWhite,
	}
}

func (p *Processor) GetCategoryCharacteristics(
	category ColorCategory,
	profile *ColorProfile,
) settings.CategoryCharacteristics {

	var categoryMap map[string]settings.CategoryCharacteristics

	if profile.Mode == Light {
		categoryMap = p.settings.Categories.Light
	} else {
		categoryMap = p.settings.Categories.Dark
	}

	if chars, ok := categoryMap[string(category)]; ok {
		return chars
	}

	// Return permissive defaults if not configured
	return settings.CategoryCharacteristics{
		MinLightness:  0.0,
		MaxLightness:  1.0,
		MinSaturation: 0.0,
		MaxSaturation: 1.0,
		MinContrast:   2.0,
	}
}

func (p *Processor) fitsCategory(
	c color.RGBA,
	category ColorCategory,
	profile *ColorProfile,
	background color.RGBA,
) bool {

	chars := p.GetCategoryCharacteristics(category, profile)
	hsla := formats.RGBAToHSLA(c)

	// Check lightness bounds
	if hsla.L < chars.MinLightness || hsla.L > chars.MaxLightness {
		return false
	}

	// Check saturation bounds
	if hsla.S < chars.MinSaturation || hsla.S > chars.MaxSaturation {
		return false
	}

	// Check contrast requirement (skip for background itself)
	if category != CategoryBackground && chars.MinContrast > 0 {
		contrast := chromatic.ContrastRatio(c, background)
		if contrast < chars.MinContrast {
			return false
		}
	}

	// Check hue constraints if specified
	if chars.HueCenter != nil && chars.HueTolerance != nil {
		hueDiff := math.Abs(hsla.H - *chars.HueCenter)
		if hueDiff > 180 {
			hueDiff = 360 - hueDiff
		}
		if hueDiff > *chars.HueTolerance {
			return false
		}
	}

	return true
}
