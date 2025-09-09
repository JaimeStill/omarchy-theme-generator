package processor

import (
	"fmt"
	"image"
	"image/color"
	"runtime"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

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
	colorFreq, totalPixels := p.extractColorFrequencies(img)

	filtered := p.filterByFrequency(colorFreq, totalPixels)

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no significant colors found")
	}

	weighted := p.createWeightedColors(filtered, totalPixels)

	pool := p.buildColorPool(weighted, totalPixels)

	profile := p.analyzeColors(weighted)

	profile.Pool = pool

	return profile, nil
}

func (p *Processor) createWeightedColors(colorFreq map[color.RGBA]uint32, totalPixels uint32) []WeightedColor {
	weighted := make([]WeightedColor, 0, len(colorFreq))

	for c, freq := range colorFreq {
		weighted = append(weighted, NewWeightedColor(c, freq, totalPixels))
	}

	sortByWeight(weighted)

	maxColors := p.settings.Extraction.MaxColorsToExtract
	if maxColors > 0 && len(weighted) > maxColors {
		weighted = weighted[:maxColors]
	}

	return weighted
}

func (p *Processor) extractColorFrequencies(img image.Image) (map[color.RGBA]uint32, uint32) {
	bounds := img.Bounds()
	totalPixels := uint32(bounds.Dx() * bounds.Dy())

	if totalPixels < 10000 {
		return p.extractSequential(img), totalPixels
	}

	return p.extractConcurrent(img, totalPixels), totalPixels
}

func (p *Processor) extractConcurrent(img image.Image, totalPixels uint32) map[color.RGBA]uint32 {
	bounds := img.Bounds()
	numWorkers := runtime.GOMAXPROCS(0)
	rowsPerWorker := bounds.Dy() / numWorkers

	if rowsPerWorker == 0 {
		rowsPerWorker = 1
		numWorkers = bounds.Dy()
	}

	type result struct {
		colors map[color.RGBA]uint32
	}

	results := make(chan result, numWorkers)

	for i := 0; i < numWorkers; i++ {
		startY := bounds.Min.Y + i*rowsPerWorker
		endY := startY + rowsPerWorker

		if i == numWorkers-1 {
			endY = bounds.Max.Y
		}

		go func(startY, endY int) {
			colors := make(map[color.RGBA]uint32, 1024)
			for y := startY; y < endY; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
					colors[rgba]++
				}
			}
			results <- result{colors: colors}
		}(startY, endY)
	}

	finalColors := make(map[color.RGBA]uint32)
	for i := 0; i < numWorkers; i++ {
		res := <-results
		for c, count := range res.colors {
			finalColors[c] += count
		}
	}

	return finalColors
}

func (p *Processor) extractSequential(img image.Image) map[color.RGBA]uint32 {
	bounds := img.Bounds()
	colorFreq := make(map[color.RGBA]uint32)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			colorFreq[rgba]++
		}
	}

	return colorFreq
}

func (p *Processor) filterByFrequency(colorFreq map[color.RGBA]uint32, totalPixels uint32) map[color.RGBA]uint32 {
	minCount := uint32(float64(totalPixels) * p.settings.MinFrequency)
	filtered := make(map[color.RGBA]uint32)

	for c, count := range colorFreq {
		if count >= minCount {
			filtered[c] = count
		}
	}

	return filtered
}
