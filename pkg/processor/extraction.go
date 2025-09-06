package processor

import (
	"image/color"
	"math"
	"sort"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

type colorScore struct {
	color color.RGBA
	score float64
	freq  uint32
}

func (p *Processor) extractByRole(colorFreq map[color.RGBA]uint32, profile *ColorProfile) *ImageColors {
	scored := make([]colorScore, 0, len(colorFreq))

	for c, freq := range colorFreq {
		score := p.calculateVisualImportance(c, freq)
		scored = append(scored, colorScore{
			color: c,
			score: score,
			freq:  freq,
		})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	return p.assignRoles(scored, profile)
}

func (p *Processor) calculateVisualImportance(c color.RGBA, freq uint32) float64 {
	hsla := formats.RGBAToHSLA(c)

	freqScore := float64(freq)
	satScore := hsla.S
	lightScore := 1.0 - math.Abs(hsla.L-0.5)*2
	contrastScore := p.calculateContrastScore(c)

	weights := [4]float64{0.3, 0.3, 0.2, 0.2}

	score := weights[0]*freqScore + weights[1]*satScore + weights[2]*lightScore + weights[3]*contrastScore

	if hsla.L < p.settings.DarkLightThreshold || hsla.L > p.settings.BrightLightThreshold {
		score *= p.settings.ExtremeLightnessPenalty
	} else if hsla.L >= 0.2 && hsla.L <= 0.8 {
		score *= p.settings.OptimalLightnessBonus
	}

	if hsla.S < p.settings.MinSaturationForBonus {
		score *= 0.8
	}

	return score
}

func (p *Processor) calculateContrastScore(c color.RGBA) float64 {
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	whiteContrast := chromatic.ContrastRatio(c, white)
	blackContrast := chromatic.ContrastRatio(c, black)

	return math.Max(whiteContrast, blackContrast) / 21.0
}

func (p *Processor) assignRoles(scored []colorScore, profile *ColorProfile) *ImageColors {
	result := &ImageColors{}

	if len(scored) == 0 {
		return result
	}

	result.MostFrequent = scored[0].color
	result.Background = p.selectBackground(scored, profile)
	result.Foreground = p.selectForeground(scored, result.Background, profile)
	result.Primary = p.selectPrimary(scored, result.Background, result.Foreground)
	result.Secondary = p.selectSecondary(scored, result.Primary, result.Background)
	result.Accent = p.selectAccent(scored, result.Primary, result.Secondary)

	return result
}

func (p *Processor) selectBackground(scored []colorScore, profile *ColorProfile) color.RGBA {
	for _, cs := range scored {
		hsla := formats.RGBAToHSLA(cs.color)

		if profile.Mode == Light && hsla.L > p.settings.LightBackgroundThreshold {
			return cs.color
		}

		if profile.Mode == Dark && hsla.L < p.settings.DarkBackgroundThreshold {
			return cs.color
		}
	}

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

func (p *Processor) selectForeground(scored []colorScore, bg color.RGBA, profile *ColorProfile) color.RGBA {
	bestContrast := 0.0
	var bestColor color.RGBA

	for _, cs := range scored {
		if cs.color == bg {
			continue
		}

		contrast := chromatic.ContrastRatio(bg, cs.color)

		if contrast > bestContrast && contrast >= p.settings.MinContrastRatio {
			bestContrast = contrast
			bestColor = cs.color
		}
	}

	if bestContrast == 0 {
		if profile.Mode == Light {
			if rgba, err := formats.ParseHex(p.settings.LightForegroundFallback); err == nil {
				return rgba
			}

			return color.RGBA{R: 32, G: 32, B: 32, A: 255}
		}

		if rgba, err := formats.ParseHex(p.settings.DarkForegroundFallback); err == nil {
			return rgba
		}

		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}

	return bestColor
}

func (p *Processor) selectPrimary(scored []colorScore, bg, fg color.RGBA) color.RGBA {
	for _, cs := range scored {
		if cs.color == bg || cs.color == fg {
			continue
		}

		hsla := formats.RGBAToHSLA(cs.color)
		if hsla.S > p.settings.MinPrimarySaturation {
			return cs.color
		}
	}

	if len(scored) > 2 {
		return scored[2].color
	}

	if rgba, err := formats.ParseHex(p.settings.PrimaryFallback); err == nil {
		return rgba
	}

	return color.RGBA{R: 100, G: 150, B: 200, A: 255}
}

func (p *Processor) selectSecondary(scored []colorScore, primary, bg color.RGBA) color.RGBA {
	for _, cs := range scored {
		if cs.color == primary || cs.color == bg {
			continue
		}
		return cs.color
	}

	primaryHSLA := formats.RGBAToHSLA(primary)
	secondaryHue := math.Mod(primaryHSLA.H+120, 360) // Triadic relationship

	return formats.HSLAToRGBA(formats.HSLA{
		H: secondaryHue,
		S: primaryHSLA.S * 0.8,
		L: primaryHSLA.L,
		A: 1.0,
	})
}

func (p *Processor) selectAccent(scored []colorScore, primary, secondary color.RGBA) color.RGBA {
	for _, cs := range scored {
		if cs.color == primary || cs.color == secondary {
			continue
		}

		hsla := formats.RGBAToHSLA(cs.color)
		if hsla.S > p.settings.MinAccentSaturation &&
			(hsla.L > p.settings.MinAccentLightness && hsla.L < p.settings.MaxAccentLightness) {
			return cs.color
		}
	}

	primaryHSLA := formats.RGBAToHSLA(primary)
	accentHue := math.Mod(primaryHSLA.H+180, 360) // Complementary

	return formats.HSLAToRGBA(formats.HSLA{
		H: accentHue,
		S: math.Min(primaryHSLA.S*1.2, 1.0),
		L: primaryHSLA.L,
		A: 1.0,
	})
}
