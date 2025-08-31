package extractor

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sort"

	tcolor "github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/errors"
)

type FrequencyStrategy struct{
	Settings *Settings
}

func (f *FrequencyStrategy) Extract(img image.Image, options *ExtractionOptions) (*ExtractionResult, error) {
	fm, err := ExtractFromImage(img)
	if err != nil {
		return nil, fmt.Errorf("frequency extraction failed: %w", err)
	}

	topColors := fm.GetTopColors(options.TopColorCount)
	if len(topColors) == 0 {
		return nil, &errors.ExtractionError{
			Stage:   "extraction",
			Details: "no colors extracted from image",
			Err:     errors.ErrNoColors,
		}
	}

	characteristics := AnalyzeImageCharacteristicsWithSettings(img, f.Settings)

	primaryColor := f.selectPrimaryColorByImportance(topColors, characteristics)
	if primaryColor == nil {
		return nil, &errors.ExtractionError{
			Stage:   "primary",
			Details: "no primary color found",
			Err:     errors.ErrNoColors,
		}
	}

	result := &ExtractionResult{
		Image:            img,
		FrequencyMap:     fm,
		DominantColor:    primaryColor,
		TopColors:        topColors,
		UniqueColors:     fm.Size(),
		TotalPixels:      fm.Total(),
		SelectedStrategy: f.Name(),
	}

	return result, nil
}

func (f *FrequencyStrategy) CanHandle(characteristics *ImageCharacteristics) bool {
	return true
}

func (f *FrequencyStrategy) Priority(characteristics *ImageCharacteristics) int {
	switch characteristics.Type {
	case LowDetail:
		return 100
	case Smooth:
		if characteristics.ColorComplexity < 50 {
			return 80
		}
		return 30
	case HighDetail:
		return 20
	case Complex:
		return 25
	default:
		return 50
	}
}

func (f *FrequencyStrategy) Name() string {
	return "frequency"
}

// FrequencyMap stores color occurrence counts with optimized memory layout.
//
// Thread Safety: NOT safe for concurrent use. Create separate instances
// for each goroutine or use external synchronization.
//
// Memory Optimization:
//   - Uses packed uint32 RGB keys (8 bytes per entry vs 20+ for Color keys)
//   - Ignores alpha channel to maximize frequency consolidation
//   - Pre-allocates capacity based on image size estimates
//   - Typical memory: 4KB for simple images, <1MB for complex 4K images
//
// Performance Characteristics:
//   - Insertion: O(1) average, amortized across map growth
//   - Lookup: O(1) average for frequency queries
//   - Sorting (GetTopColors): O(n log n) where n = unique colors
//   - Memory overhead: ~32 bytes per unique color
type FrequencyMap struct {
	counts map[uint32]uint32 // packed RGB -> count mapping
	total  uint32            // total pixels counted
}

// NewFrequencyMap creates a new frequency map with optional initial capacity.
// For best performance, provide the expected unique color count as capacity.
// If capacity <= 0, uses a reasonable default of 1024.
func NewFrequencyMap(capacity int) *FrequencyMap {
	if capacity <= 0 {
		capacity = CurrentSettings.Extraction.InitialMapCapacity
	}
	return &FrequencyMap{
		counts: make(map[uint32]uint32, capacity),
		total:  0,
	}
}

func (fm *FrequencyMap) Add(c color.Color) {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	packed := packRGB(rgba.R, rgba.G, rgba.B)
	fm.counts[packed]++
	fm.total++
}

func (fm *FrequencyMap) AddRGBA(r, g, b, a uint8) {
	packed := packRGB(r, g, b)
	fm.counts[packed]++
	fm.total++
}

func (fm *FrequencyMap) Size() int {
	return len(fm.counts)
}

func (fm *FrequencyMap) Total() uint32 {
	return fm.total
}

func (fm *FrequencyMap) GetTopColors(n int) []*ColorFrequency {
	if n <= 0 || n > len(fm.counts) {
		n = len(fm.counts)
	}

	colors := make([]*ColorFrequency, 0, len(fm.counts))
	for packed, count := range fm.counts {
		r, g, b := unpackRGB(packed)
		colors = append(colors, &ColorFrequency{
			Color:      tcolor.NewRGB(r, g, b),
			Count:      count,
			Percentage: float64(count) / float64(fm.total) * 100.0,
		})
	}

	sort.Slice(colors, func(i, j int) bool {
		return colors[i].Count > colors[j].Count
	})

	if n < len(colors) {
		return colors[:n]
	}
	return colors
}

func (fm *FrequencyMap) GetDominantColor() *ColorFrequency {
	if len(fm.counts) == 0 {
		return nil
	}

	var maxCount uint32
	var maxPacked uint32

	for packed, count := range fm.counts {
		if count > maxCount {
			maxCount = count
			maxPacked = packed
		}
	}

	r, g, b := unpackRGB(maxPacked)
	return &ColorFrequency{
		Color:      tcolor.NewRGB(r, g, b),
		Count:      maxCount,
		Percentage: float64(maxCount) / float64(fm.total) * 100.0,
	}
}

func (fm *FrequencyMap) FilterByThreshold(thresholdPercent float64) []*ColorFrequency {
	if thresholdPercent <= 0 {
		return fm.GetTopColors(0)
	}

	minCount := uint32(float64(fm.total) * thresholdPercent / 100.0)
	results := make([]*ColorFrequency, 0)

	for packed, count := range fm.counts {
		if count >= minCount {
			r, g, b := unpackRGB(packed)
			results = append(results, &ColorFrequency{
				Color:      tcolor.NewRGB(r, g, b),
				Count:      count,
				Percentage: float64(count) / float64(fm.total) * 100.0,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Count > results[j].Count
	})

	return results
}

// ExtractFromImage builds a frequency map from all pixels in an image.
// It uses optimized extraction paths for different image types (RGBA, NRGBA, generic)
// and pre-allocates capacity based on image characteristics for best performance.
// This is the main entry point for color frequency analysis.
func ExtractFromImage(img image.Image) (*FrequencyMap, error) {
	bounds := img.Bounds()
	if bounds.Empty() {
		return nil, &errors.ImageLoadError{
			Path:      "memory",
			Operation: "extract colors",
			Err:       errors.ErrEmptyImage,
		}
	}

	width := bounds.Dx()
	height := bounds.Dy()
	totalPixels := width * height

	estimatedCapacity := min(totalPixels/20, CurrentSettings.Extraction.MaxMapCapacity)

	fm := NewFrequencyMap(estimatedCapacity)

	switch img := img.(type) {
	case *image.RGBA:
		extractRGBA(fm, img)
	case *image.NRGBA:
		extractNRGBA(fm, img)
	default:
		extractGeneric(fm, img)
	}

	if fm.Size() == 0 {
		return nil, &errors.ExtractionError{
			Stage:   "frequency",
			Details: "no colors extracted from image",
			Err:     errors.ErrNoColors,
		}
	}

	return fm, nil
}

func packRGB(r, g, b uint8) uint32 {
	return uint32(r)<<16 | uint32(g)<<8 | uint32(b)
}

func unpackRGB(packed uint32) (r, g, b uint8) {
	r = uint8((packed >> 16) & 0xFF)
	g = uint8((packed >> 8) & 0xFF)
	b = uint8(packed & 0xFF)
	return
}

func extractRGBA(fm *FrequencyMap, img *image.RGBA) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			fm.AddRGBA(
				img.Pix[i],
				img.Pix[i+1],
				img.Pix[i+2],
				img.Pix[i+3],
			)
		}
	}
}

func extractNRGBA(fm *FrequencyMap, img *image.NRGBA) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			fm.AddRGBA(
				img.Pix[i],
				img.Pix[i+1],
				img.Pix[i+2],
				img.Pix[i+3],
			)
		}
	}
}

func extractGeneric(fm *FrequencyMap, img image.Image) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			fm.Add(img.At(x, y))
		}
	}
}

func (f *FrequencyStrategy) selectPrimaryColorByImportance(topColors []*ColorFrequency, characteristics *ImageCharacteristics) *tcolor.Color {
	if len(topColors) == 0 {
		return nil
	}

	bestColor := topColors[0].Color
	bestScore := f.calculateVisualImportance(topColors[0], characteristics, topColors)

	maxCandidates := min(len(topColors), 20)

	for i := 1; i < maxCandidates; i++ {
		score := f.calculateVisualImportance(topColors[i], characteristics, topColors)

		if score > bestScore {
			bestScore = score
			bestColor = topColors[i].Color
		}
	}

	return bestColor
}

type valueWeights struct {
	frequency  float64
	saturation float64
	lightness  float64
	contrast   float64
}

func (f *FrequencyStrategy) calculateVisualImportance(colorFreq *ColorFrequency, characteristics *ImageCharacteristics, allColors []*ColorFrequency) float64 {
	_, s, l := colorFreq.Color.HSL()
	freq := f.Settings.Frequency

	frequencyScore := colorFreq.Percentage / 100.0
	saturationScore := s

	lightnessScore := 1.0
	if (l < freq.DarkLightThreshold || l > freq.BrightLightThreshold) && colorFreq.Percentage < 60.0 {
		lightnessScore = freq.ExtremeLightnessPenalty
	} else if l > 0.2 && l < 0.8 {
		lightnessScore = freq.OptimalLightnessBonus
	}

	contrastScore := f.calculateContrastImportance(colorFreq, allColors)

	var weights [4]float64
	switch characteristics.Type {
	case HighDetail:
		weights = freq.HighDetailWeights
	case LowDetail:
		weights = freq.LowDetailWeights
	case Smooth:
		weights = freq.SmoothWeights
	case Complex:
		weights = freq.ComplexWeights
	default:
		weights = freq.DefaultWeights
	}

	totalScore := weights[0]*frequencyScore +
		weights[1]*saturationScore +
		weights[2]*lightnessScore +
		weights[3]*contrastScore

	return totalScore
}

func (f *FrequencyStrategy) calculateContrastImportance(colorFreq *ColorFrequency, allColors []*ColorFrequency) float64 {
	if len(allColors) < 2 {
		return 0.5
	}

	maxContrast := 0.0

	background := allColors[0]
	if background != colorFreq {
		contrast := colorFreq.Color.ContrastRatio(background.Color)
		maxContrast = math.Max(maxContrast, contrast)
	}

	checkCount := min(len(allColors), f.Settings.Frequency.MaxContrastSamples)
	for i := 1; i < checkCount; i++ {
		if allColors[i] != colorFreq {
			contrast := colorFreq.Color.ContrastRatio(allColors[i].Color)
			maxContrast = math.Max(maxContrast, contrast)
		}
	}

	return math.Min(maxContrast/21.0, 1.0)
}
