package extractor

import (
	"image"
	"image/color"
	"sort"

	otgcolor "github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/errors"
)

// ColorFrequency represents a color and its occurrence frequency in an image.
// It includes both raw count and percentage for convenient display and analysis.
type ColorFrequency struct {
	Color      *otgcolor.Color // The color value
	Count      uint32          // Number of pixels with this color
	Percentage float64         // Percentage of total pixels (0.0-100.0)
}

// FrequencyMap efficiently stores color occurrence counts using packed RGB as map keys.
// It ignores alpha channel for frequency counting and tracks total pixels processed.
type FrequencyMap struct {
	counts map[uint32]uint32 // packed RGB -> count mapping
	total  uint32            // total pixels counted
}

// NewFrequencyMap creates a new frequency map with optional initial capacity.
// For best performance, provide the expected unique color count as capacity.
// If capacity <= 0, uses a reasonable default of 1024.
func NewFrequencyMap(capacity int) *FrequencyMap {
	if capacity <= 0 {
		capacity = 1024
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
			Color:      otgcolor.NewRGB(r, g, b),
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
		Color:      otgcolor.NewRGB(r, g, b),
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
				Color:      otgcolor.NewRGB(r, g, b),
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

	estimatedCapacity := min(totalPixels/20, 65536)

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
