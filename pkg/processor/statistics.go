package processor

import (
	"math"
	"sort"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func (p *Processor) calculateStatistics(pool ColorPool) ColorStatistics {
	stats := ColorStatistics{
		SaturationGroups: make(map[string]float64),
	}

	stats.HueHistogram = p.calculateHueHistogram(pool.AllColors)

	stats.LightnessHistogram = p.calculateLightnessHistorgram(pool.AllColors)

	total := float64(len(pool.AllColors))
	if total > 0 {
		stats.SaturationGroups["gray"] = float64(len(pool.BySaturation.Gray)) / total
		stats.SaturationGroups["muted"] = float64(len(pool.BySaturation.Muted)) / total
		stats.SaturationGroups["normal"] = float64(len(pool.BySaturation.Normal)) / total
		stats.SaturationGroups["vibrant"] = float64(len(pool.BySaturation.Vibrant)) / total
	}

	stats.PrimaryHue, stats.SecondaryHue, stats.TertiaryHue = p.findDominantHues(pool.ByHue)
	stats.ChromaticDiversity = p.calculateChromaticDiversity(stats.HueHistogram)
	stats.ContrastRange = p.calculateContrastRange(pool.AllColors)

	stats.HueVariance = p.calculateHueVariance(pool.AllColors)
	stats.LightnessSpread = p.calculateLightnessSpread(pool.ByLightness)
	stats.SaturationSpread = p.calculateSaturationSpread(pool.BySaturation)

	return stats
}

func (p *Processor) calculateChromaticDiversity(hueHistogram []float64) float64 {
	entropy := 0.0
	for _, prob := range hueHistogram {
		if prob > 0 {
			entropy -= prob * math.Log2(prob)
		}
	}

	maxEntropy := math.Log2(float64(len(hueHistogram)))
	if maxEntropy > 0 {
		return entropy / maxEntropy
	}

	return 0
}

func (p *Processor) calculateContrastRange(colors []WeightedColor) float64 {
	if len(colors) == 0 {
		return 0
	}

	// Initialize with extreme values to find actual range
	minLum := 1.0  // Start with maximum possible luminance
	maxLum := 0.0  // Start with minimum possible luminance

	for _, wc := range colors {
		lum := chromatic.Luminance(wc.RGBA)
		if lum < minLum {
			minLum = lum
		}
		if lum > maxLum {
			maxLum = lum
		}
	}

	return maxLum - minLum
}

func (p *Processor) calculateHueHistogram(colors []WeightedColor) []float64 {
	sectorCount := p.settings.HueSectorCount
	sectorSize := p.settings.HueSectorSize
	histogram := make([]float64, sectorCount)
	totalWeight := 0.0

	for _, wc := range colors {
		hsla := formats.RGBAToHSLA(wc.RGBA)
		if hsla.S >= p.settings.GrayscaleThreshold {
			sector := int(hsla.H / sectorSize)
			if sector >= 0 && sector < sectorCount {
				histogram[sector] += wc.Weight
				totalWeight += wc.Weight
			}
		}
	}

	if totalWeight > 0 {
		for i := range histogram {
			histogram[i] /= totalWeight
		}
	}

	return histogram
}

func (p *Processor) calculateHueVariance(colors []WeightedColor) float64 {
	var nonGrayscaleHSLAs []formats.HSLA
	for _, wc := range colors {
		hsla := formats.RGBAToHSLA(wc.RGBA)
		if hsla.S >= p.settings.GrayscaleThreshold {
			nonGrayscaleHSLAs = append(nonGrayscaleHSLAs, hsla)
		}
	}

	if len(nonGrayscaleHSLAs) > 0 {
		return chromatic.CalculateHueVariance(nonGrayscaleHSLAs)
	}

	return 0
}

func (p *Processor) calculateLightnessHistorgram(colors []WeightedColor) []float64 {
	histogram := make([]float64, 10)
	totalWeight := 0.0

	for _, wc := range colors {
		hsla := formats.RGBAToHSLA(wc.RGBA)
		bucket := int(hsla.L * 10)
		if bucket >= 10 {
			bucket = 9
		}
		if bucket < 0 {
			bucket = 0
		}
		histogram[bucket] += wc.Weight
		totalWeight += wc.Weight
	}

	if totalWeight > 0 {
		for i := range histogram {
			histogram[i] /= totalWeight
		}
	}

	return histogram
}

func (p *Processor) calculateLightnessSpread(groups LightnessGroups) float64 {
	total := len(groups.Dark) + len(groups.Mid) + len(groups.Light)
	if total == 0 {
		return 0
	}

	darkRatio := float64(len(groups.Dark)) / float64(total)
	midRatio := float64(len(groups.Mid)) / float64(total)
	lightRatio := float64(len(groups.Light)) / float64(total)

	idealRatio := 1.0 / 3.0
	deviation := math.Abs(darkRatio-idealRatio) + math.Abs(midRatio-idealRatio) + math.Abs(lightRatio-idealRatio)

	return 1.0 - (deviation / 2.0)
}

func (p *Processor) calculateSaturationSpread(groups SaturationGroups) float64 {
	groupCount := 0
	if len(groups.Gray) > 0 {
		groupCount++
	}
	if len(groups.Muted) > 0 {
		groupCount++
	}
	if len(groups.Normal) > 0 {
		groupCount++
	}
	if len(groups.Vibrant) > 0 {
		groupCount++
	}

	return float64(groupCount) / 4.0
}

func (p *Processor) findDominantHues(hueFamilies HueFamilies) (primary, secondary, tertiary float64) {
	type hueBucket struct {
		hue    float64
		weight float64
	}

	buckets := make([]hueBucket, 0, len(hueFamilies))
	sectorSize := p.settings.HueSectorSize

	for sector, colors := range hueFamilies {
		totalWeight := 0.0
		for _, wc := range colors {
			totalWeight += wc.Weight
		}
		buckets = append(buckets, hueBucket{
			hue:    float64(sector) * sectorSize,
			weight: totalWeight,
		})
	}

	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].weight > buckets[j].weight
	})

	if len(buckets) > 0 {
		primary = buckets[0].hue
	}
	if len(buckets) > 1 {
		secondary = buckets[1].hue
	}
	if len(buckets) > 2 {
		tertiary = buckets[2].hue
	}

	return primary, secondary, tertiary
}
