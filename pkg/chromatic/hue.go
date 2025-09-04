package chromatic

import (
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func FindDominantHue(hslas []formats.HSLA) float64 {
	if len(hslas) == 0 {
		return math.NaN()
	}

	sinSum := 0.0
	cosSum := 0.0

	for _, hsla := range hslas {
		rad := hsla.H * math.Pi / 180.0
		sinSum += math.Sin(rad)
		cosSum += math.Cos(rad)
	}

	avgRad := math.Atan2(sinSum, cosSum)
	avgDeg := avgRad * 180.0 / math.Pi

	if avgDeg < 0 {
		avgDeg += 360
	}

	return avgDeg
}

func CalculateHueVariance(hslas []formats.HSLA) float64 {
	if len(hslas) <= 1 {
		return 0.0
	}

	dominantHue := FindDominantHue(hslas)
	sumSquaredDiff := 0.0

	for _, hsla := range hslas {
		diff := hueDistance(hsla.H, dominantHue)
		sumSquaredDiff += diff * diff
	}

	return math.Sqrt(sumSquaredDiff / float64(len(hslas)))
}
