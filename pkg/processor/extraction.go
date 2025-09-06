package processor

import (
	"image/color"
	"math"
	"sort"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func (p *Processor) extractColors(
	colorFreq map[color.RGBA]uint32,
	profile *ColorProfile,
	totalPixels uint32,
) *ImageColors {
	result := &ImageColors{
		ColorFrequency:     colorFreq,
		Categories:         make(map[ColorCategory]color.RGBA),
		CategoryCandidates: make(map[ColorCategory][]ColorCandidate),
		TotalPixels:        totalPixels,
		UniqueColors:       len(colorFreq),
	}

	background := p.selectBackground(colorFreq, profile)
	result.Categories[CategoryBackground] = background

	p.categorizeColors(colorFreq, profile, background, result)
	p.selectBestCategoryColors(result)

	allCategories := GetAllCategories()
	filledCount := len(result.Categories)
	result.CoverageRatio = float64(filledCount) / float64(len(allCategories))

	return result
}

func (p *Processor) selectBackground(
	colorFreq map[color.RGBA]uint32,
	profile *ColorProfile,
) color.RGBA {
	chars := p.GetCategoryCharacteristics(CategoryBackground, profile)

	var bestColor color.RGBA
	bestScore := -1.0

	for c, freq := range colorFreq {
		hsla := formats.RGBAToHSLA(c)

		if hsla.L >= chars.MinLightness && hsla.L <= chars.MaxLightness &&
			hsla.S >= chars.MinSaturation && hsla.S <= chars.MaxSaturation {
			freqScore := float64(freq)
			idealLight := (chars.MinLightness + chars.MaxLightness) / 2
			lightScore := 1.0 - math.Abs(hsla.L-idealLight)

			score := freqScore * lightScore

			if score > bestScore {
				bestScore = score
				bestColor = c
			}
		}
	}

	if bestScore < 0 {
		if profile.Mode == Light {
			if rgba, err := formats.ParseHex(p.settings.LightBackgroundFallback); err == nil {
				return rgba
			}
			return color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}

		if rgba, err := formats.ParseHex(p.settings.DarkBackgroundFallback); err == nil {
			return rgba
		}
		return color.RGBA{R: 32, G: 32, B: 32, A: 255}
	}

	return bestColor
}

func (p *Processor) categorizeColors(
	colorFreq map[color.RGBA]uint32,
	profile *ColorProfile,
	background color.RGBA,
	result *ImageColors,
) {
	categoryOrder := p.GetCategoryPriorityOrder(profile)
	maxCandidates := p.settings.Extraction.MaxCandidatesPerCategory

	for _, category := range categoryOrder {
		if category == CategoryBackground {
			continue // Already handled
		}

		candidates := []ColorCandidate{}

		for c, freq := range colorFreq {
			if p.fitsCategory(c, category, profile, background) {
				score := p.calculateCategoryFitScore(
					c, category, profile, background, freq, result.TotalPixels,
				)

				if score > 0 {
					candidates = append(candidates, ColorCandidate{
						Color:     c,
						Frequency: freq,
						Score:     score,
					})
				}
			}
		}

		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Score > candidates[j].Score
		})

		if len(candidates) > maxCandidates {
			candidates = candidates[:maxCandidates]
		}

		if len(candidates) > 0 {
			result.CategoryCandidates[category] = candidates
		}
	}
}

func (p *Processor) selectBestCategoryColors(result *ImageColors) {
	for category, candidates := range result.CategoryCandidates {
		if len(candidates) > 0 {
			result.Categories[category] = candidates[0].Color
		}
	}
}
