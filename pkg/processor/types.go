package processor

import (
	"image/color"
)

type ThemeMode string

const (
	Light ThemeMode = "Light"
	Dark  ThemeMode = "Dark"
)

// ColorCluster represents a visually distinct color group with UI-relevant metadata
type ColorCluster struct {
	color.RGBA                   // The representative color
	Weight      float64          // Combined weight (0.0-1.0)
	Lightness   float64          // Pre-calculated HSL lightness for efficiency
	Saturation  float64          // Pre-calculated HSL saturation for efficiency
	Hue         float64          // Hue in degrees (0-360)
	IsNeutral   bool            // Grayscale or very low saturation
	IsDark      bool            // L < 0.3
	IsLight     bool            // L > 0.7
	IsMuted     bool            // S < 0.3
	IsVibrant   bool            // S > 0.7
}

// ColorProfile is the minimal data needed for theme generation
type ColorProfile struct {
	Mode       ThemeMode      // Light or Dark theme base
	Colors     []ColorCluster // Distinct colors, sorted by weight
	HasColor   bool          // False if image is essentially grayscale
	ColorCount int           // Number of distinct colors found
}

// WeightedColor is an internal type for processing
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