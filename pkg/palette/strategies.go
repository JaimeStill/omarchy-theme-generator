package palette

import (
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
)

// MonochromaticStrategy generates colors using a single hue with variations in saturation and lightness.
// This creates harmonious palettes perfect for minimalist designs.
type MonochromaticStrategy struct{}

func (ms *MonochromaticStrategy) Name() string {
	return "monochromatic"
}

func (ms *MonochromaticStrategy) Description() string {
	return "Single hue with saturation and lightness variations"
}

func (ms *MonochromaticStrategy) Generate(baseColor *color.Color, size int) []*color.Color {
	h, s, l := baseColor.HSL()
	palette := make([]*color.Color, size)
	
	// Always include the base color
	palette[0] = baseColor
	
	// Generate variations
	for i := 1; i < size; i++ {
		// Create variations by adjusting saturation and lightness
		factor := float64(i) / float64(size-1)
		
		// Alternate between lighter and darker, more and less saturated
		var newS, newL float64
		
		if i%2 == 0 {
			// Even indices: adjust lightness more
			newL = l + (factor * 0.4) - 0.2 // Range: -0.2 to +0.2
			newS = s * (1.0 - factor*0.3)   // Reduce saturation slightly
		} else {
			// Odd indices: adjust saturation more
			newS = s * (1.0 - factor*0.5) // Vary saturation more
			newL = l - (factor * 0.3) + 0.15 // Slightly darker
		}
		
		// Clamp values to valid range
		newS = math.Max(0, math.Min(1, newS))
		newL = math.Max(0, math.Min(1, newL))
		
		palette[i] = color.NewHSL(h, newS, newL)
	}
	
	return palette
}

// AnalogousStrategy generates colors using adjacent hues on the color wheel.
// This creates natural, pleasing palettes often found in nature.
type AnalogousStrategy struct{}

func (as *AnalogousStrategy) Name() string {
	return "analogous"
}

func (as *AnalogousStrategy) Description() string {
	return "Adjacent hues on the color wheel (±30°)"
}

func (as *AnalogousStrategy) Generate(baseColor *color.Color, size int) []*color.Color {
	h, s, l := baseColor.HSL()
	palette := make([]*color.Color, size)
	
	// Include base color at center
	palette[0] = baseColor
	
	// Spread colors across ±30 degrees
	angleStep := 60.0 / float64(size-1) // Total 60° range
	
	for i := 1; i < size; i++ {
		// Calculate hue offset from base
		offset := (float64(i) - float64(size)/2) * angleStep
		newH := NormalizeHue(h + offset/360.0)
		
		// Slight variations in saturation and lightness for interest
		satVariation := 0.1 * math.Sin(float64(i))
		lightVariation := 0.1 * math.Cos(float64(i))
		
		newS := math.Max(0, math.Min(1, s+satVariation))
		newL := math.Max(0, math.Min(1, l+lightVariation))
		
		palette[i] = color.NewHSL(newH, newS, newL)
	}
	
	return palette
}

// ComplementaryStrategy generates colors using opposite hues on the color wheel.
// This creates high-contrast, vibrant palettes with strong visual impact.
type ComplementaryStrategy struct{}

func (cs *ComplementaryStrategy) Name() string {
	return "complementary"
}

func (cs *ComplementaryStrategy) Description() string {
	return "Opposite hues on the color wheel (180° apart)"
}

func (cs *ComplementaryStrategy) Generate(baseColor *color.Color, size int) []*color.Color {
	h, s, l := baseColor.HSL()
	palette := make([]*color.Color, size)
	
	// Base color
	palette[0] = baseColor
	
	// Complementary color (opposite on wheel)
	complementH := NormalizeHue(h + 0.5)
	
	// Distribute colors between base and complement
	halfSize := size / 2
	
	// First half: variations of base color
	for i := 1; i <= halfSize && i < size; i++ {
		factor := float64(i) / float64(halfSize+1)
		newS := s * (1.0 - factor*0.3)
		newL := l + (factor*0.4 - 0.2)
		
		newS = math.Max(0, math.Min(1, newS))
		newL = math.Max(0, math.Min(1, newL))
		
		palette[i] = color.NewHSL(h, newS, newL)
	}
	
	// Second half: variations of complement
	for i := halfSize + 1; i < size; i++ {
		factor := float64(i-halfSize) / float64(size-halfSize)
		newS := s * (1.0 - factor*0.3)
		newL := l + (factor*0.4 - 0.2)
		
		newS = math.Max(0, math.Min(1, newS))
		newL = math.Max(0, math.Min(1, newL))
		
		palette[i] = color.NewHSL(complementH, newS, newL)
	}
	
	return palette
}

// TriadicStrategy generates colors using three equally spaced hues on the color wheel.
// This creates vibrant, balanced palettes while maintaining harmony.
type TriadicStrategy struct{}

func (ts *TriadicStrategy) Name() string {
	return "triadic"
}

func (ts *TriadicStrategy) Description() string {
	return "Three equally spaced hues (120° apart)"
}

func (ts *TriadicStrategy) Generate(baseColor *color.Color, size int) []*color.Color {
	h, s, l := baseColor.HSL()
	palette := make([]*color.Color, size)
	
	// Three main hues
	hue1 := h
	hue2 := NormalizeHue(h + 1.0/3.0)
	hue3 := NormalizeHue(h + 2.0/3.0)
	
	// Distribute colors across the three hues
	colorsPerHue := size / 3
	remainder := size % 3
	
	idx := 0
	
	// Generate variations for each hue
	for hueIdx, hue := range []float64{hue1, hue2, hue3} {
		count := colorsPerHue
		if hueIdx < remainder {
			count++
		}
		
		for i := 0; i < count && idx < size; i++ {
			factor := float64(i) / float64(count)
			
			// Vary saturation and lightness
			newS := s * (0.7 + factor*0.3)
			newL := l + (factor*0.4 - 0.2)
			
			newS = math.Max(0, math.Min(1, newS))
			newL = math.Max(0, math.Min(1, newL))
			
			palette[idx] = color.NewHSL(hue, newS, newL)
			idx++
		}
	}
	
	return palette
}

// TetradicStrategy generates colors using four hues arranged in a rectangle on the color wheel.
// This creates rich, diverse palettes with multiple color relationships.
type TetradicStrategy struct{}

func (ts *TetradicStrategy) Name() string {
	return "tetradic"
}

func (ts *TetradicStrategy) Description() string {
	return "Four hues in a rectangle pattern on the wheel"
}

func (ts *TetradicStrategy) Generate(baseColor *color.Color, size int) []*color.Color {
	h, s, l := baseColor.HSL()
	palette := make([]*color.Color, size)
	
	// Four hues: base, complement, and two others at 60° and 240°
	hue1 := h
	hue2 := NormalizeHue(h + 60.0/360.0)
	hue3 := NormalizeHue(h + 180.0/360.0)
	hue4 := NormalizeHue(h + 240.0/360.0)
	
	// Distribute colors across the four hues
	colorsPerHue := size / 4
	remainder := size % 4
	
	idx := 0
	
	// Generate variations for each hue
	for hueIdx, hue := range []float64{hue1, hue2, hue3, hue4} {
		count := colorsPerHue
		if hueIdx < remainder {
			count++
		}
		
		for i := 0; i < count && idx < size; i++ {
			factor := float64(i) / float64(count)
			
			// Create variations
			newS := s * (0.6 + factor*0.4)
			newL := l + (factor*0.5 - 0.25)
			
			newS = math.Max(0, math.Min(1, newS))
			newL = math.Max(0, math.Min(1, newL))
			
			palette[idx] = color.NewHSL(hue, newS, newL)
			idx++
		}
	}
	
	return palette
}

// SplitComplementaryStrategy generates colors using the base hue and two hues adjacent to its complement.
// This creates vibrant palettes with less tension than pure complementary.
type SplitComplementaryStrategy struct{}

func (scs *SplitComplementaryStrategy) Name() string {
	return "split-complementary"
}

func (scs *SplitComplementaryStrategy) Description() string {
	return "Base hue plus two hues adjacent to its complement"
}

func (scs *SplitComplementaryStrategy) Generate(baseColor *color.Color, size int) []*color.Color {
	h, s, l := baseColor.HSL()
	palette := make([]*color.Color, size)
	
	// Three main hues: base and two split complements
	hue1 := h
	hue2 := NormalizeHue(h + 150.0/360.0) // 30° before complement
	hue3 := NormalizeHue(h + 210.0/360.0) // 30° after complement
	
	// Distribute colors across the three hues
	colorsPerHue := size / 3
	remainder := size % 3
	
	idx := 0
	
	// Generate variations for each hue
	for hueIdx, hue := range []float64{hue1, hue2, hue3} {
		count := colorsPerHue
		if hueIdx < remainder {
			count++
		}
		
		for i := 0; i < count && idx < size; i++ {
			factor := float64(i) / float64(count)
			
			// Base hue gets more variations, split complements get fewer
			var newS, newL float64
			if hueIdx == 0 {
				// More variation for base hue
				newS = s * (0.5 + factor*0.5)
				newL = l + (factor*0.6 - 0.3)
			} else {
				// Less variation for split complements
				newS = s * (0.7 + factor*0.3)
				newL = l + (factor*0.4 - 0.2)
			}
			
			newS = math.Max(0, math.Min(1, newS))
			newL = math.Max(0, math.Min(1, newL))
			
			palette[idx] = color.NewHSL(hue, newS, newL)
			idx++
		}
	}
	
	return palette
}