package analysis

import (
	"image/color"
	"math"
	"sort"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

type ColorCluster struct {
	Colors []color.RGBA
	Center color.RGBA
}

func ClusterColors(colors []color.RGBA, threshold float64) []ColorCluster {
	if len(colors) == 0 {
		return nil
	}

	var clusters []ColorCluster
	processed := make([]bool, len(colors))

	for i, c := range colors {
		if processed[i] {
			continue
		}

		cluster := ColorCluster{
			Colors: []color.RGBA{c},
		}
		processed[i] = true

		for j := i + 1; j < len(colors); j++ {
			if processed[j] {
				continue
			}

			if chromatic.DistanceLAB(c, colors[j]) <= threshold {
				cluster.Colors = append(cluster.Colors, colors[j])
				processed[j] = true
			}
		}

		cluster.Center = FindRepresentativeColor(cluster.Colors)
		clusters = append(clusters, cluster)
	}

	sort.Slice(clusters, func(i, j int) bool {
		return len(clusters[i].Colors) > len(clusters[j].Colors)
	})

	return clusters
}

func FindRepresentativeColor(colors []color.RGBA) color.RGBA {
	if len(colors) == 0 {
		return color.RGBA{}
	}

	if len(colors) == 1 {
		return colors[0]
	}

	totalL, totalA, totalB := 0.0, 0.0, 0.0
	for _, c := range colors {
		lab := formats.RGBAToLAB(c)
		totalL += lab.L
		totalA += lab.A
		totalB += lab.B
	}

	avgL := totalL / float64(len(colors))
	avgA := totalA / float64(len(colors))
	avgB := totalB / float64(len(colors))
	avgLAB := formats.NewLAB(avgL, avgA, avgB)

	minDist := math.MaxFloat64
	var representative color.RGBA

	for _, c := range colors {
		lab := formats.RGBAToLAB(c)
		dist := math.Sqrt(
			math.Pow(lab.L-avgLAB.L, 2) +
				math.Pow(lab.A-avgLAB.A, 2) +
				math.Pow(lab.B-avgLAB.B, 2),
		)

		if dist < minDist {
			minDist = dist
			representative = c
		}
	}

	return representative
}
