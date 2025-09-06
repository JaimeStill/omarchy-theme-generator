package processor

import (
	"fmt"
	"image"
	"image/color"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

type ThemeMode string

const (
	Light ThemeMode = "Light"
	Dark  ThemeMode = "Dark"
)

type ColorProfile struct {
	Mode            ThemeMode
	ColorScheme     chromatic.ColorScheme
	IsGrayscale     bool
	IsMonochromatic bool
	DominantHue     float64
	HueVariance     float64
	AvgLuminance    float64
	AvgSaturation   float64
	Colors          ImageColors
}

type ImageColors struct {
	ColorFrequency     map[color.RGBA]uint32              `json:"color_frequency"`
	Categories         map[ColorCategory]color.RGBA       `json:"categories"`
	CategoryCandidates map[ColorCategory][]ColorCandidate `json:"category_candidates"`
	TotalPixels        uint32                             `json:"total_pixels"`
	UniqueColors       int                                `json:"unique_colors"`
	CoverageRatio      float64                            `json:"coverage_ratio"`
}

type ColorCandidate struct {
	Color     color.RGBA `json:"color"`
	Frequency uint32     `json:"frequency"`
	Score     float64    `json:"score"`
}

type Processor struct {
	settings *settings.Settings
	chroma   *chromatic.Chroma
}

func New(s *settings.Settings) *Processor {
	return &Processor{
		settings: s,
		chroma:   chromatic.NewChroma(s),
	}
}

func (p *Processor) ProcessImage(img image.Image) (*ColorProfile, error) {
	bounds := img.Bounds()
	colorFreq := make(map[color.RGBA]uint32)
	totalPixels := uint32(bounds.Dx() * bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			colorFreq[rgba]++
		}
	}

	minCount := uint32(float64(totalPixels) * p.settings.MinFrequency)
	filtered := make(map[color.RGBA]uint32)
	for c, count := range colorFreq {
		if count >= minCount {
			filtered[c] = count
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no significant colors found")
	}

	profile := p.analyzeColors(filtered)
	profile.Colors = *p.extractColors(filtered, profile, totalPixels)

	return profile, nil
}

func (p *Processor) calculateTotalPixels(colorFreq map[color.RGBA]uint32) uint32 {
	total := uint32(0)
	for _, freq := range colorFreq {
		total += freq
	}

	return total
}
