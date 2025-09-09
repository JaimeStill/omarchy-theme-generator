package processor

import "sort"

func (p *Processor) buildColorPool(colors []WeightedColor, totalPixels uint32) ColorPool {
	pool := ColorPool{
		AllColors:    colors,
		TotalPixels:  totalPixels,
		UniqueColors: len(colors),
	}

	pool.DominantColors = p.selectDominant(colors, p.settings.Extraction.DominantColorCount)

	pool.ByLightness = p.groupByLightness(colors)
	pool.BySaturation = p.groupBySaturation(colors)
	pool.ByHue = p.groupByHue(colors)

	pool.Statistics = p.calculateStatistics(pool)

	return pool
}

func (p *Processor) selectDominant(colors []WeightedColor, count int) []WeightedColor {
	if count <= 0 || count >= len(colors) {
		return colors
	}

	sorted := make([]WeightedColor, len(colors))
	copy(sorted, colors)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Weight > sorted[j].Weight
	})

	return sorted[:count]
}
