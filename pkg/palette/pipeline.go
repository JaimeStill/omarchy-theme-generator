package palette

import (
	"fmt"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/extractor"
)

// PipelineResult contains the final palette and metadata about how it was generated.
type PipelineResult struct {
	Palette           []*color.Color
	Mode              SynthesisMode
	Strategy          string
	BaseColor         *color.Color
	ExtractedColors   int
	SynthesizedColors int
	ValidationResult  *ValidationResult
}

// GenerationPipeline implements the extraction → hybrid → synthesis failover logic.
type GenerationPipeline struct {
	generator *PaletteGenerator
	validator *PaletteValidator
	options   *SynthesisOptions
}

// NewGenerationPipeline creates a pipeline with the specified options.
func NewGenerationPipeline(options *SynthesisOptions, background *color.Color) *GenerationPipeline {
	if options == nil {
		options = DefaultSynthesisOptions()
	}
	
	return &GenerationPipeline{
		generator: NewPaletteGenerator(options),
		validator: NewPaletteValidator(background),
		options:   options,
	}
}

// GenerateFromExtraction creates a palette based on extraction results and analysis.
// It automatically selects the appropriate mode based on the image characteristics.
func (gp *GenerationPipeline) GenerateFromExtraction(result *extractor.ExtractionResult) (*PipelineResult, error) {
	analysis := result.AnalyzeForThemeGeneration()
	
	// Determine mode based on analysis
	var mode SynthesisMode
	switch analysis.SuggestedStrategy {
	case "extract":
		mode = ModeExtract
	case "hybrid":
		mode = ModeHybrid
	case "synthesize":
		mode = ModeSynthesize
	default:
		mode = ModeHybrid
	}
	
	// Get base color for synthesis
	var baseColor *color.Color
	if analysis.IsGrayscale {
		// Use configured base color for grayscale images
		baseColor = color.NewHSL(
			gp.options.BaseHue/360.0,
			gp.options.BaseSaturation,
			gp.options.BaseLightness,
		)
	} else if analysis.IsMonochromatic {
		// Use the dominant hue for monochromatic images
		baseColor = color.NewHSL(
			analysis.DominantHue/360.0,
			gp.options.BaseSaturation,
			gp.options.BaseLightness,
		)
	} else {
		// Use the dominant color from extraction
		baseColor = result.DominantColor
		if baseColor == nil {
			// Fallback to first top color
			if len(result.TopColors) > 0 {
				baseColor = result.TopColors[0].Color
			} else {
				// Ultimate fallback
				baseColor = color.NewHSL(
					gp.options.BaseHue/360.0,
					gp.options.BaseSaturation,
					gp.options.BaseLightness,
				)
			}
		}
	}
	
	// Generate palette based on mode
	pipelineResult := &PipelineResult{
		Mode:      mode,
		BaseColor: baseColor,
	}
	
	switch mode {
	case ModeExtract:
		pipelineResult.Palette = gp.extractPalette(result)
		pipelineResult.ExtractedColors = len(pipelineResult.Palette)
		pipelineResult.Strategy = "extraction"
		
	case ModeHybrid:
		extracted, synthesized := gp.hybridPalette(result, baseColor)
		pipelineResult.Palette = append(extracted, synthesized...)
		pipelineResult.ExtractedColors = len(extracted)
		pipelineResult.SynthesizedColors = len(synthesized)
		pipelineResult.Strategy = gp.selectStrategyForImage(analysis)
		
	case ModeSynthesize:
		pipelineResult.Palette = gp.synthesizePalette(baseColor, analysis)
		pipelineResult.SynthesizedColors = len(pipelineResult.Palette)
		pipelineResult.Strategy = gp.selectStrategyForImage(analysis)
	}
	
	// Validate the palette
	pipelineResult.ValidationResult = gp.validator.Validate(pipelineResult.Palette)
	
	// Ensure contrast if needed
	if pipelineResult.ValidationResult.FailingColors > 0 {
		pipelineResult.Palette = gp.validator.EnsureContrast(pipelineResult.Palette)
		// Re-validate after adjustment
		pipelineResult.ValidationResult = gp.validator.Validate(pipelineResult.Palette)
	}
	
	return pipelineResult, nil
}

// extractPalette creates a palette purely from extracted colors.
func (gp *GenerationPipeline) extractPalette(result *extractor.ExtractionResult) []*color.Color {
	targetSize := gp.options.PaletteSize
	
	// Get top colors up to target size
	topColors := result.TopColors
	if len(topColors) > targetSize {
		topColors = topColors[:targetSize]
	}
	
	palette := make([]*color.Color, len(topColors))
	for i, cf := range topColors {
		palette[i] = cf.Color
	}
	
	return palette
}

// hybridPalette combines extracted colors with synthesized ones.
func (gp *GenerationPipeline) hybridPalette(result *extractor.ExtractionResult, baseColor *color.Color) ([]*color.Color, []*color.Color) {
	targetSize := gp.options.PaletteSize
	
	// Use top 1/3 of target size from extraction
	extractCount := targetSize / 3
	if extractCount > len(result.TopColors) {
		extractCount = len(result.TopColors)
	}
	
	extracted := make([]*color.Color, extractCount)
	for i := 0; i < extractCount; i++ {
		extracted[i] = result.TopColors[i].Color
	}
	
	// Synthesize the remaining colors
	synthesizeCount := targetSize - extractCount
	strategy := gp.options.PreferredStrategy
	synthesized, _ := gp.generator.GenerateFromBase(baseColor, strategy)
	
	if len(synthesized) > synthesizeCount {
		synthesized = synthesized[:synthesizeCount]
	}
	
	return extracted, synthesized
}

// synthesizePalette creates a palette purely from color theory.
func (gp *GenerationPipeline) synthesizePalette(baseColor *color.Color, analysis *extractor.ThemeGenerationAnalysis) []*color.Color {
	strategy := gp.selectStrategyForImage(analysis)
	palette, _ := gp.generator.GenerateFromBase(baseColor, strategy)
	return palette
}

// selectStrategyForImage chooses the best synthesis strategy based on image analysis.
func (gp *GenerationPipeline) selectStrategyForImage(analysis *extractor.ThemeGenerationAnalysis) string {
	if analysis.IsGrayscale {
		// For grayscale, use complementary for contrast
		return "complementary"
	} else if analysis.IsMonochromatic {
		// For monochromatic, use the monochromatic strategy
		return "monochromatic"
	} else if analysis.DominantCoverage > 60 {
		// High dominance benefits from split-complementary
		return "split-complementary"
	} else if analysis.UniqueColors < 5 {
		// Low diversity benefits from triadic
		return "triadic"
	}
	
	// Default to configured preference
	return gp.options.PreferredStrategy
}

// GenerateWithOverrides allows manual specification of base color and strategy.
func (gp *GenerationPipeline) GenerateWithOverrides(
	baseColor *color.Color,
	strategy string,
	mode SynthesisMode,
) (*PipelineResult, error) {
	if baseColor == nil {
		baseColor = color.NewHSL(
			gp.options.BaseHue/360.0,
			gp.options.BaseSaturation,
			gp.options.BaseLightness,
		)
	}
	
	result := &PipelineResult{
		Mode:      mode,
		BaseColor: baseColor,
		Strategy:  strategy,
	}
	
	// Generate based on mode
	switch mode {
	case ModeSynthesize:
		palette, err := gp.generator.GenerateFromBase(baseColor, strategy)
		if err != nil {
			return nil, fmt.Errorf("synthesis failed: %w", err)
		}
		result.Palette = palette
		result.SynthesizedColors = len(palette)
		
	default:
		return nil, fmt.Errorf("override generation only supports synthesis mode")
	}
	
	// Validate
	result.ValidationResult = gp.validator.Validate(result.Palette)
	
	// Ensure contrast if needed
	if result.ValidationResult.FailingColors > 0 {
		result.Palette = gp.validator.EnsureContrast(result.Palette)
		result.ValidationResult = gp.validator.Validate(result.Palette)
	}
	
	return result, nil
}

// String provides a summary of the pipeline result.
func (pr *PipelineResult) String() string {
	modeStr := "unknown"
	switch pr.Mode {
	case ModeExtract:
		modeStr = "extraction"
	case ModeHybrid:
		modeStr = "hybrid"
	case ModeSynthesize:
		modeStr = "synthesis"
	}
	
	return fmt.Sprintf(
		"PipelineResult: mode=%s, strategy=%s, colors=%d (extracted=%d, synthesized=%d), passing_wcag=%d/%d",
		modeStr,
		pr.Strategy,
		len(pr.Palette),
		pr.ExtractedColors,
		pr.SynthesizedColors,
		pr.ValidationResult.PassingColors,
		pr.ValidationResult.TotalColors,
	)
}