package chromatic

import (
	"math"
	"sort"
)

type ColorScheme string

const (
	Grayscale          ColorScheme = "Grayscale"
	Monochromatic      ColorScheme = "Monochromatic"
	Analogous          ColorScheme = "Analogous"
	Complementary      ColorScheme = "Complementary"
	SplitComplementary ColorScheme = "SplitComplementary"
	Triadic            ColorScheme = "Triadic"
	Square             ColorScheme = "Square"
	Tetradic           ColorScheme = "Tetradic"
	Custom             ColorScheme = "Custom"
)

func isTriadic(hues []float64, tolerance float64) bool {
	if len(hues) != 3 {
		return false
	}

	sorted := sortHues(hues)

	d1 := hueDistance(sorted[0], sorted[1])
	d2 := hueDistance(sorted[1], sorted[2])
	d3 := hueDistance(sorted[2], sorted[0])

	target := 120.0

	return math.Abs(d1-target) <= tolerance &&
		math.Abs(d2-target) <= tolerance &&
		math.Abs(d3-target) <= tolerance
}

func isSquare(hues []float64, tolerance float64) bool {
	if len(hues) != 4 {
		return false
	}

	sorted := sortHues(hues)

	target := 90.0

	for i := 0; i < 4; i++ {
		next := (i + 1) % 4
		dist := hueDistance(sorted[i], sorted[next])
		if math.Abs(dist-target) > tolerance {
			return false
		}
	}
	return true
}

func isTetradic(hues []float64, tolerance float64) bool {
	if len(hues) != 4 {
		return false
	}

	for i := 0; i < len(hues); i++ {
		for j := i + 1; j < len(hues); j++ {
			if math.Abs(hueDistance(hues[i], hues[j])-180) <= tolerance {
				var others []float64
				for k := 0; k < len(hues); k++ {
					if k != i && k != j {
						others = append(others, hues[k])
					}
				}
				if len(others) == 2 &&
					math.Abs(hueDistance(others[0], others[1])-180) <= tolerance {
					return true
				}
			}
		}
	}
	return false
}

func isSplitComplementary(hues []float64, tolerance float64) bool {
	if len(hues) != 3 {
		return false
	}

	for i := 0; i < 3; i++ {
		base := hues[i]
		var others []float64
		for j := 0; j < 3; j++ {
			if j != i {
				others = append(others, hues[j])
			}
		}

		d1 := hueDistance(base, others[0])
		d2 := hueDistance(base, others[1])

		if (math.Abs(d1-150) <= tolerance && math.Abs(d2-210) <= tolerance) ||
			(math.Abs(d1-210) <= tolerance && math.Abs(d2-150) <= tolerance) {
			return true
		}
	}

	return false
}

func sortHues(hues []float64) []float64 {
	if len(hues) == 0 {
		return nil
	}

	sorted := make([]float64, len(hues))
	copy(sorted, hues)
	sort.Float64s(sorted)
	return sorted
}
func sortHuesCircular(hues []float64) []float64 {
	if len(hues) <= 1 {
		return sortHues(hues)
	}

	sorted := sortHues(hues)

	maxGap := 0.0
	maxGapIndex := 0

	for i := 0; i < len(sorted); i++ {
		next := (i + 1) % len(sorted)
		gap := sorted[next] - sorted[i]
		if next == 0 {
			gap = (360 + sorted[next]) - sorted[i]
		}

		if gap > maxGap {
			maxGap = gap
			maxGapIndex = next
		}
	}

	if maxGapIndex != 0 {
		result := make([]float64, len(sorted))
		copy(result, sorted[maxGapIndex:])
		copy(result[len(sorted)-maxGapIndex:], sorted[:maxGapIndex])
		return result
	}

	return sorted
}
