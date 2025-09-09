package processor

import (
	"image/color"
)

type ThemeMode string

const (
	Light ThemeMode = "Light"
	Dark  ThemeMode = "Dark"
)

type ColorProfile struct {
	Mode            ThemeMode
	IsGrayscale     bool
	IsMonochromatic bool
	DominantHue     float64
	HueVariance     float64
	AvgLuminance    float64
	AvgSaturation   float64
	Pool            ColorPool
}

type WeightedColor struct {
	color.RGBA
	Frequency uint32
	Weight    float64
}

func NewWeightedColor(c color.RGBA, freq, total uint32) WeightedColor {
	return WeightedColor{
		RGBA:      c,
		Frequency: freq,
		Weight:    float64(freq) / float64(total),
	}
}

type ColorPool struct {
	AllColors      []WeightedColor
	DominantColors []WeightedColor

	ByLightness  LightnessGroups
	BySaturation SaturationGroups
	ByHue        HueFamilies

	TotalPixels  uint32
	UniqueColors int

	Statistics ColorStatistics
}

type LightnessGroups struct {
	Dark  []WeightedColor
	Mid   []WeightedColor
	Light []WeightedColor
}

func NewLightnessGroups() LightnessGroups {
	return LightnessGroups{
		Dark:  make([]WeightedColor, 0),
		Mid:   make([]WeightedColor, 0),
		Light: make([]WeightedColor, 0),
	}
}

type SaturationGroups struct {
	Gray    []WeightedColor
	Muted   []WeightedColor
	Normal  []WeightedColor
	Vibrant []WeightedColor
}

func NewSaturationGroups() SaturationGroups {
	return SaturationGroups{
		Gray:    make([]WeightedColor, 0),
		Muted:   make([]WeightedColor, 0),
		Normal:  make([]WeightedColor, 0),
		Vibrant: make([]WeightedColor, 0),
	}
}

type HueFamilies map[int][]WeightedColor

type ColorStatistics struct {
	HueHistogram       []float64
	LightnessHistogram []float64
	SaturationGroups   map[string]float64

	PrimaryHue         float64
	SecondaryHue       float64
	TertiaryHue        float64
	ChromaticDiversity float64
	ContrastRange      float64

	HueVariance      float64
	LightnessSpread  float64
	SaturationSpread float64
}
