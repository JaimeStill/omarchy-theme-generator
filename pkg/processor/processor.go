package processor

import (
	"fmt"
	"image"
	"image/color"
	"runtime"
	"sort"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
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
	colorFreq, totalSamples := p.extractColors(img)

	if len(colorFreq) == 0 {
		return nil, fmt.Errorf("no colors found in image")
	}

	weighted := p.createWeightedColors(colorFreq, totalSamples)
	clusters := p.clusterColors(weighted)
	clusters = p.filterForUI(clusters)

	if len(clusters) == 0 {
		return nil, fmt.Errorf("no suitable colors found for UI theme")
	}

	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Weight > clusters[j].Weight
	})

	mode := p.calculateThemeMode(clusters)
	hasColor := p.hasSignificantColor(clusters)

	return &ColorProfile{
		Mode:       mode,
		Colors:     clusters,
		HasColor:   hasColor,
		ColorCount: len(clusters),
	}, nil
}

func (p *Processor) extractColors(img image.Image) (map[color.RGBA]uint32, uint32) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	totalPixels := width * height

	sampleRate := p.calculateSampleRate(width, height)

	if totalPixels > 100000 && runtime.GOMAXPROCS(0) > 1 {
		return p.extractColorsConcurrent(img, sampleRate)
	}

	return p.extractColorsSequential(img, sampleRate)
}

func (p *Processor) extractColorsSequential(img image.Image, sampleRate int) (map[color.RGBA]uint32, uint32) {
	bounds := img.Bounds()
	colorFreq := make(map[color.RGBA]uint32)
	var totalSamples uint32

	for y := bounds.Min.Y; y < bounds.Max.Y; y += sampleRate {
		for x := bounds.Min.X; x < bounds.Max.X; x += sampleRate {
			rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			rgba = formats.QuantizeColor(rgba, uint8(p.settings.Formats.QuantizationBits))
			colorFreq[rgba]++
			totalSamples++
		}
	}

	return colorFreq, totalSamples
}

func (p *Processor) extractColorsConcurrent(img image.Image, sampleRate int) (map[color.RGBA]uint32, uint32) {
	bounds := img.Bounds()
	numWorkers := runtime.GOMAXPROCS(0)
	rowsPerWorker := bounds.Dy() / numWorkers

	if rowsPerWorker == 0 {
		rowsPerWorker = 1
		numWorkers = bounds.Dy()
	}

	type result struct {
		colors  map[color.RGBA]uint32
		samples uint32
	}

	results := make(chan result, numWorkers)

	for i := 0; i < numWorkers; i++ {
		startY := bounds.Min.Y + i*rowsPerWorker
		endY := startY + rowsPerWorker
		if i == numWorkers-1 {
			endY = bounds.Max.Y
		}

		go func(startY, endY int) {
			colors := make(map[color.RGBA]uint32)
			var samples uint32

			for y := startY; y < endY; y += sampleRate {
				for x := bounds.Min.X; x < bounds.Max.X; x += sampleRate {
					rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
					rgba = formats.QuantizeColor(rgba, uint8(p.settings.Formats.QuantizationBits))
					colors[rgba]++
					samples++
				}
			}

			results <- result{colors: colors, samples: samples}
		}(startY, endY)
	}

	finalColors := make(map[color.RGBA]uint32)
	var totalSamples uint32

	for i := 0; i < numWorkers; i++ {
		res := <-results
		for c, count := range res.colors {
			finalColors[c] += count
		}
		totalSamples += res.samples
	}

	return finalColors, totalSamples
}

func (p *Processor) createWeightedColors(colorFreq map[color.RGBA]uint32, totalSamples uint32) []WeightedColor {
	weighted := make([]WeightedColor, 0, len(colorFreq))
	minFreq := uint32(float64(totalSamples) * p.settings.Processor.MinFrequency)

	for c, freq := range colorFreq {
		if freq >= minFreq {
			weighted = append(weighted, NewWeightedColor(c, freq, totalSamples))
		}
	}

	return weighted
}

func (p *Processor) clusterColors(colors []WeightedColor) []ColorCluster {
	if len(colors) == 0 {
		return nil
	}

	sort.Slice(colors, func(i, j int) bool {
		return colors[i].Weight > colors[j].Weight
	})

	var clusters []ColorCluster
	used := make([]bool, len(colors))

	for i, color := range colors {
		if used[i] {
			continue
		}

		cluster := p.createCluster(color)
		used[i] = true

		for j := i + 1; j < len(colors); j++ {
			if used[j] {
				continue
			}

			if p.chroma.ColorsSimilar(color.RGBA, colors[j].RGBA) {
				cluster.Weight += colors[j].Weight
				used[j] = true
			}
		}

		if cluster.Weight >= p.settings.Processor.MinClusterWeight {
			clusters = append(clusters, cluster)
		}
	}

	return clusters
}

func (p *Processor) createCluster(wc WeightedColor) ColorCluster {
	hsla := formats.RGBAToHSLA(wc.RGBA)

	return ColorCluster{
		RGBA:       wc.RGBA,
		Weight:     wc.Weight,
		Lightness:  hsla.L,
		Saturation: hsla.S,
		Hue:        hsla.H,
		IsNeutral:  hsla.S < p.settings.Chromatic.NeutralThreshold,
		IsDark:     hsla.L < p.settings.Chromatic.DarkLightnessMax,
		IsLight:    hsla.L > p.settings.Chromatic.LightLightnessMin,
		IsMuted:    hsla.S < p.settings.Chromatic.MutedSaturationMax && hsla.S >= p.settings.Chromatic.NeutralThreshold,
		IsVibrant:  hsla.S > p.settings.Chromatic.VibrantSaturationMin,
	}
}

func (p *Processor) filterForUI(clusters []ColorCluster) []ColorCluster {
	filtered := make([]ColorCluster, 0, len(clusters))
	var hasPureBlack, hasPureWhite bool

	for _, cluster := range clusters {
		if cluster.Lightness < p.settings.Processor.PureBlackThreshold {
			if !hasPureBlack && cluster.Weight > 0.01 {
				hasPureBlack = true
				filtered = append(filtered, cluster)
			}
			continue
		}
		if cluster.Lightness > p.settings.Processor.PureWhiteThreshold {
			if !hasPureWhite && cluster.Weight > 0.01 {
				hasPureWhite = true
				filtered = append(filtered, cluster)
			}
			continue
		}

		if cluster.Weight < p.settings.Processor.MinUIColorWeight {
			continue
		}

		filtered = append(filtered, cluster)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Weight > filtered[j].Weight
	})

	if len(filtered) > p.settings.Processor.MaxUIColors {
		filtered = filtered[:p.settings.Processor.MaxUIColors]
	}

	return filtered
}

func (p *Processor) calculateThemeMode(clusters []ColorCluster) ThemeMode {
	if len(clusters) == 0 {
		return Dark
	}

	var weightedLightness float64
	var totalWeight float64

	maxClusters := p.settings.Processor.ThemeModeMaxClusters
	if len(clusters) < maxClusters {
		maxClusters = len(clusters)
	}

	for i := 0; i < maxClusters; i++ {
		cluster := clusters[i]
		weightedLightness += cluster.Lightness * cluster.Weight
		totalWeight += cluster.Weight
	}

	avgLightness := weightedLightness / totalWeight

	if avgLightness > p.settings.Processor.LightThemeThreshold {
		return Light
	}
	return Dark
}

func (p *Processor) hasSignificantColor(clusters []ColorCluster) bool {
	colorWeight := 0.0

	for _, cluster := range clusters {
		if !cluster.IsNeutral {
			colorWeight += cluster.Weight
		}
	}

	return colorWeight > p.settings.Processor.SignificantColorThreshold
}

func (p *Processor) calculateSampleRate(width, height int) int {
	pixels := width * height

	switch {
	case pixels > 8000000:
		return 4
	case pixels > 4000000:
		return 3
	case pixels > 2000000:
		return 2
	default:
		return 1
	}
}