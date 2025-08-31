package extractor

import (
	"fmt"
	"image"
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/errors"
)

type SaliencyStrategy struct{
	Settings *Settings
}

type SaliencyMap struct {
	scores [][]float64
	width  int
	height int
}

func (s *SaliencyStrategy) CanHandle(characteristics *ImageCharacteristics) bool {
	// Primary condition: High detail images always use saliency
	if characteristics.Type == HighDetail {
		return true
	}
	
	// Edge-based detection
	if characteristics.EdgeDensity > s.Settings.Strategy.SaliencyEdgeThreshold {
		return true
	}
	
	// Color complexity detection: Visually rich images benefit from saliency
	if characteristics.ColorComplexity > s.Settings.Strategy.SaliencyColorComplexity && 
	   characteristics.AverageSaturation > s.Settings.Strategy.SaliencySaturationThreshold {
		return true
	}
	
	return false
}

func (s *SaliencyStrategy) Priority(characteristics *ImageCharacteristics) int {
	switch characteristics.Type {
	case HighDetail:
		return 100
	case Complex:
		return 80
	default:
		return 10
	}
}

func (s *SaliencyStrategy) Name() string {
	return "saliency"
}

func (s *SaliencyStrategy) Extract(img image.Image, options *ExtractionOptions) (*ExtractionResult, error) {
	fm, err := ExtractFromImage(img)
	if err != nil {
		return nil, fmt.Errorf("saliency extraction failed: %w", err)
	}

	topColors := fm.GetTopColors(options.TopColorCount * 2)

	if len(topColors) == 0 {
		return nil, &errors.ExtractionError{
			Stage:   "extraction",
			Details: "no colors extracted from image",
			Err:     errors.ErrNoColors,
		}
	}

	saliencyMap := s.generateSaliencyMap(img)

	saliencyWeightedColors := s.weightColorsBySaliency(img, topColors, saliencyMap)

	finalCount := min(len(saliencyWeightedColors), options.TopColorCount)
	finalColors := saliencyWeightedColors[:finalCount]

	if len(finalColors) == 0 {
		return nil, &errors.ExtractionError{
			Stage:   "saliency",
			Details: "no salient colors found",
			Err:     errors.ErrNoColors,
		}
	}

	result := &ExtractionResult{
		Image:            img,
		FrequencyMap:     fm,
		DominantColor:    finalColors[0].Color,
		TopColors:        finalColors,
		UniqueColors:     fm.Size(),
		TotalPixels:      fm.Total(),
		SelectedStrategy: s.Name(),
	}

	return result, nil
}

func (sm *SaliencyMap) At(x, y int) float64 {
	if x < 0 || y < 0 || x >= sm.width || y >= sm.height {
		return 0.0
	}

	return sm.scores[y][x]
}

func (s *SaliencyStrategy) generateSaliencyMap(img image.Image) *SaliencyMap {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	saliencyMap := &SaliencyMap{
		scores: make([][]float64, height),
		width:  width,
		height: height,
	}

	for y := range saliencyMap.scores {
		saliencyMap.scores[y] = make([]float64, width)
	}

	sampleRate := s.Settings.Saliency.SaliencyMapSampleRate
	for y := bounds.Min.Y + 2; y < bounds.Max.Y-2; y += sampleRate {
		for x := bounds.Min.X + 2; x < bounds.Max.X-2; x += sampleRate {
			saliency := s.calculateLocalSaliency(img, x, y)

			mapY := y - bounds.Min.Y
			mapX := x - bounds.Min.X

			if mapY >= 0 && mapY < height && mapX >= 0 && mapX < width {
				saliencyMap.scores[mapY][mapX] = saliency

				radius := s.Settings.Saliency.SaliencyMapSpreadRadius
				for dy := -radius; dy <= radius; dy++ {
					for dx := -radius; dx <= radius; dx++ {
						ny, nx := mapY+dy, mapX+dx
						if ny >= 0 && ny < height && nx >= 0 && nx < width {
							distance := math.Sqrt(float64(dx*dx + dy*dy))
							weight := math.Max(0, 1.0-distance/3.0)
							saliencyMap.scores[ny][nx] = math.Max(saliencyMap.scores[ny][nx], saliency*weight)
						}
					}
				}
			}
		}
	}

	return saliencyMap
}

func (s *SaliencyStrategy) calculateLocalSaliency(img image.Image, x, y int) float64 {
	localContrast := s.calculateLocalContrast(img, x, y)
	edgeStrength := s.calculateEdgeStrength(img, x, y)
	colorUniqueness := s.calculateColorUniqueness(img, x, y)

	saliency := s.Settings.Saliency.LocalContrastWeight*localContrast + 
				s.Settings.Saliency.EdgeStrengthWeight*edgeStrength + 
				s.Settings.Saliency.ColorUniquenessWeight*colorUniqueness

	return saliency
}

func (s *SaliencyStrategy) calculateLocalContrast(img image.Image, centerX, centerY int) float64 {
	centerGray := getGrayscaleValue(img.At(centerX, centerY))

	totalDiff := 0.0
	sampleCount := 0

	radius := s.Settings.Saliency.ContrastWindowRadius
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}

			neighborGray := getGrayscaleValue(img.At(centerX+dx, centerY+dy))
			diff := math.Abs(float64(centerGray - neighborGray))
			totalDiff += diff
			sampleCount++
		}
	}

	if sampleCount == 0 {
		return 0.0
	}

	avgContrast := totalDiff / float64(sampleCount)
	return avgContrast / 255.0
}

func (s *SaliencyStrategy) calculateEdgeStrength(img image.Image, x, y int) float64 {
	gx := getGrayscaleValue(img.At(x+1, y)) - getGrayscaleValue(img.At(x-1, y))
	gy := getGrayscaleValue(img.At(x, y+1)) - getGrayscaleValue(img.At(x, y-1))

	edgeStrength := math.Sqrt(float64(gx*gx + gy*gy))
	return math.Min(edgeStrength/255.0, 1.0)
}

func (s *SaliencyStrategy) calculateColorUniqueness(img image.Image, centerX, centerY int) float64 {
	centerColor := img.At(centerX, centerY)
	cr, cg, cb, _ := centerColor.RGBA()
	cr, cg, cb = cr>>8, cg>>8, cb>>8

	similarCount := 8
	totalSamples := 0

	radius := s.Settings.Saliency.UniquenessWindowRadius
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}

			neightborColor := img.At(centerX+dx, centerY+dy)
			nr, ng, nb, _ := neightborColor.RGBA()
			nr, ng, nb = nr>>8, ng>>8, nb>>8

			colorDist := math.Sqrt(float64((cr-nr)*(cr-nr) + (cg-ng)*(cg-ng) + (cb-nb)*(cb-nb)))

			if colorDist < s.Settings.Saliency.ColorSimilarityThreshold {
				similarCount++
			}

			totalSamples++
		}
	}

	if totalSamples == 0 {
		return 0.5
	}

	uniqueness := 1.0 - (float64(similarCount) / float64(totalSamples))
	return uniqueness
}

func (s *SaliencyStrategy) weightColorsBySaliency(img image.Image, topColors []*ColorFrequency, saliencyMap *SaliencyMap) []*ColorFrequency {
	bounds := img.Bounds()

	colorSaliencies := make(map[uint32]float64)
	colorPixelCounts := make(map[uint32]int)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 4 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 4 {
			rgba := img.At(x, y)
			r, g, b, _ := rgba.RGBA()
			r, g, b = r>>8, g>>8, b>>8

			packed := uint32(r)<<16 | uint32(g)<<8 | uint32(b)

			mapX := x - bounds.Min.X
			mapY := y - bounds.Min.Y
			saliency := saliencyMap.At(mapX, mapY)

			colorSaliencies[packed] += saliency
			colorPixelCounts[packed]++
		}
	}

	for packed := range colorSaliencies {
		if count := colorPixelCounts[packed]; count > 0 {
			colorSaliencies[packed] /= float64(count)
		}
	}

	weightedColors := make([]*ColorFrequency, 0, len(topColors))

	for _, colorFreq := range topColors {
		r, g, b := colorFreq.Color.RGB()
		packed := uint32(r)<<16 | uint32(g)<<8 | uint32(b)

		avgSaliency := colorSaliencies[packed]

		combinedScore := s.Settings.Saliency.FrequencyWeight*colorFreq.Percentage + 
						 s.Settings.Saliency.SaliencyWeight*avgSaliency*100.0

		weightedColors = append(weightedColors, &ColorFrequency{
			Color:      colorFreq.Color,
			Count:      colorFreq.Count,
			Percentage: combinedScore,
		})
	}

	for i := 0; i < len(weightedColors); i++ {
		for j := i + 1; j < len(weightedColors); j++ {
			if weightedColors[j].Percentage > weightedColors[i].Percentage {
				weightedColors[i], weightedColors[j] = weightedColors[j], weightedColors[i]
			}
		}
	}

	return weightedColors
}
