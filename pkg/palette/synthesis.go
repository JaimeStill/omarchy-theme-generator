package palette

import (
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
)

// SynthesisStrategy defines the interface for all color synthesis strategies.
// Each strategy generates a palette based on color theory principles.
type SynthesisStrategy interface {
	// Generate creates a palette of the specified size using the strategy's color theory rules.
	// The baseColor serves as the starting point for generation.
	Generate(baseColor *color.Color, size int) []*color.Color
	
	// Name returns the human-readable name of the strategy.
	Name() string
	
	// Description provides a brief explanation of the color theory behind the strategy.
	Description() string
}

// SynthesisMode determines how synthesis integrates with extraction.
type SynthesisMode int

const (
	// ModeExtract uses only extracted colors from the image.
	ModeExtract SynthesisMode = iota
	
	// ModeHybrid combines extracted colors with synthesized ones.
	ModeHybrid
	
	// ModeSynthesize generates colors purely from color theory.
	ModeSynthesize
)

// PaletteGenerator orchestrates the extraction-synthesis pipeline.
// It automatically selects the appropriate mode based on image analysis.
type PaletteGenerator struct {
	strategies map[string]SynthesisStrategy
	options    *SynthesisOptions
}

// SynthesisOptions configures the palette generation behavior.
type SynthesisOptions struct {
	// PreferredStrategy specifies which synthesis strategy to use.
	PreferredStrategy string
	
	// FallbackStrategy used when preferred strategy cannot generate enough colors.
	FallbackStrategy string
	
	// BaseHue for synthesis when image is grayscale (0-360 degrees).
	BaseHue float64
	
	// BaseSaturation for synthesis (0.0-1.0).
	BaseSaturation float64
	
	// BaseLightness for synthesis (0.0-1.0).
	BaseLightness float64
	
	// IncludeTemperatureMatchedGrays adds grays that match the color temperature.
	IncludeTemperatureMatchedGrays bool
	
	// MinContrast ensures all colors meet this WCAG contrast ratio.
	MinContrast float64
	
	// PaletteSize is the target number of colors to generate.
	PaletteSize int
}

// DefaultSynthesisOptions provides sensible defaults for palette generation.
func DefaultSynthesisOptions() *SynthesisOptions {
	return &SynthesisOptions{
		PreferredStrategy:              "complementary",
		FallbackStrategy:               "analogous",
		BaseHue:                        220.0, // A pleasant blue
		BaseSaturation:                 0.7,
		BaseLightness:                  0.5,
		IncludeTemperatureMatchedGrays: true,
		MinContrast:                    4.5, // WCAG AA
		PaletteSize:                    16,
	}
}

// NewPaletteGenerator creates a generator with all available strategies registered.
func NewPaletteGenerator(options *SynthesisOptions) *PaletteGenerator {
	if options == nil {
		options = DefaultSynthesisOptions()
	}
	
	gen := &PaletteGenerator{
		strategies: make(map[string]SynthesisStrategy),
		options:    options,
	}
	
	// Register all available strategies
	gen.RegisterStrategy(&MonochromaticStrategy{})
	gen.RegisterStrategy(&AnalogousStrategy{})
	gen.RegisterStrategy(&ComplementaryStrategy{})
	gen.RegisterStrategy(&TriadicStrategy{})
	gen.RegisterStrategy(&TetradicStrategy{})
	gen.RegisterStrategy(&SplitComplementaryStrategy{})
	
	return gen
}

// RegisterStrategy adds a synthesis strategy to the generator.
func (pg *PaletteGenerator) RegisterStrategy(strategy SynthesisStrategy) {
	pg.strategies[strategy.Name()] = strategy
}

// GenerateFromBase creates a palette using the specified base color and strategy.
func (pg *PaletteGenerator) GenerateFromBase(baseColor *color.Color, strategyName string) ([]*color.Color, error) {
	strategy, exists := pg.strategies[strategyName]
	if !exists {
		// Fall back to default strategy
		strategy = pg.strategies[pg.options.FallbackStrategy]
		if strategy == nil {
			strategy = &ComplementaryStrategy{}
		}
	}
	
	palette := strategy.Generate(baseColor, pg.options.PaletteSize)
	
	// Add temperature-matched grays if requested
	if pg.options.IncludeTemperatureMatchedGrays {
		grays := GenerateTemperatureMatchedGrays(baseColor, 4)
		palette = append(palette, grays...)
	}
	
	return palette, nil
}

// GenerateTemperatureMatchedGrays creates grayscale colors that harmonize with the given hue.
// Warm hues get warm grays, cool hues get cool grays.
func GenerateTemperatureMatchedGrays(baseColor *color.Color, count int) []*color.Color {
	h, _, _ := baseColor.HSL()
	grays := make([]*color.Color, count)
	
	// Determine if hue is warm (0-60, 300-360) or cool (60-300)
	hDegrees := h * 360
	isWarm := hDegrees < 60 || hDegrees > 300
	
	// Create grays with slight tint matching the temperature
	for i := 0; i < count; i++ {
		// Lightness from dark to light
		lightness := float64(i+1) / float64(count+1)
		
		// Very low saturation for gray appearance
		saturation := 0.02 // 2% saturation for subtle tint
		
		// Use the base hue for temperature matching
		if isWarm {
			// Warm grays - slight red/yellow tint
			if hDegrees > 300 {
				h = 350.0 / 360.0 // Reddish
			} else {
				h = 40.0 / 360.0 // Yellowish
			}
		} else {
			// Cool grays - slight blue tint
			h = 210.0 / 360.0 // Bluish
		}
		
		grays[i] = color.NewHSL(h, saturation, lightness)
	}
	
	return grays
}

// NormalizeHue ensures a hue value is within 0-1 range.
func NormalizeHue(h float64) float64 {
	h = math.Mod(h, 1.0)
	if h < 0 {
		h += 1.0
	}
	return h
}

// DegreesToHue converts degrees (0-360) to hue value (0-1).
func DegreesToHue(degrees float64) float64 {
	return NormalizeHue(degrees / 360.0)
}

// HueToDegrees converts hue value (0-1) to degrees (0-360).
func HueToDegrees(hue float64) float64 {
	return NormalizeHue(hue) * 360.0
}