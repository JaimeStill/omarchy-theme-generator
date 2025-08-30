package theme

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/extractor"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/palette"
)

// Generator orchestrates the complete theme generation process.
// It integrates extraction, synthesis, mode detection, and validation
// to produce complete themes ready for config generation.
type Generator struct {
	pipeline  *palette.GenerationPipeline
	detector  *ModeDetector
	validator *palette.PaletteValidator
}

// NewGenerator creates a theme generator with the specified background color
// for WCAG validation. Use nil for default white background.
func NewGenerator(backgroundColor *color.Color) *Generator {
	if backgroundColor == nil {
		backgroundColor = color.NewRGB(255, 255, 255) // Default white background
	}
	
	return &Generator{
		detector:  NewModeDetector(),
		validator: palette.NewPaletteValidator(backgroundColor),
	}
}

// GenerateTheme creates a complete theme from the provided configuration.
// This is the main entry point for theme generation, integrating all
// existing pipeline components with mode detection and user overrides.
func (g *Generator) GenerateTheme(config ThemeConfig) (*Theme, error) {
	startTime := time.Now()
	
	// Validate input
	if config.SourceImage == nil {
		return nil, fmt.Errorf("source image is required")
	}
	
	// Extract basic image information
	bounds := config.SourceImage.Bounds()
	imageSize := ImageSize{
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
	}
	
	// Derive theme name from source if not provided
	themeName := config.Name
	if themeName == "" {
		// This will be empty for generated images, but that's okay
		themeName = "Generated Theme"
	}
	
	// Step 1: Extract colors from image
	extractionStart := time.Now()
	extractionResult, err := extractor.ExtractColorsFromImage(config.SourceImage, nil)
	if err != nil {
		return nil, fmt.Errorf("color extraction failed: %w", err)
	}
	extractionTime := time.Since(extractionStart)
	
	// Step 2: Initialize synthesis pipeline with options
	synthesisOpts := config.SynthesisOpts
	if synthesisOpts == nil {
		synthesisOpts = palette.DefaultSynthesisOptions()
	}
	
	// Create pipeline with background color for validation
	backgroundColor := g.getBackgroundColor(config.Overrides)
	g.pipeline = palette.NewGenerationPipeline(synthesisOpts, backgroundColor)
	
	// Step 3: Generate palette using existing pipeline
	synthesisStart := time.Now()
	pipelineResult, err := g.pipeline.GenerateFromExtraction(extractionResult)
	if err != nil {
		return nil, fmt.Errorf("palette generation failed: %w", err)
	}
	synthesisTime := time.Since(synthesisStart)
	
	// Step 4: Determine theme mode
	var themeMode ThemeMode
	if config.Mode == ModeAuto {
		themeMode = g.detector.DetectWithPrimary(config.SourceImage, pipelineResult.BaseColor)
	} else {
		themeMode = config.Mode
	}
	
	isLight := (themeMode == ModeLight)
	
	// Step 5: Derive semantic colors from palette
	primary, background, foreground := g.deriveSemanticColors(
		pipelineResult.Palette, 
		pipelineResult.BaseColor, 
		isLight,
	)
	
	// Step 6: Apply user overrides with validation
	validationStart := time.Now()
	if config.Overrides.HasOverrides() {
		primary, background, foreground, err = g.applyOverrides(
			config.Overrides, 
			primary, 
			background, 
			foreground,
		)
		if err != nil {
			return nil, fmt.Errorf("override validation failed: %w", err)
		}
	}
	
	// Final validation of all colors
	allColors := append([]*color.Color{primary, background, foreground}, pipelineResult.Palette...)
	validationResult := g.validator.Validate(allColors)
	validationTime := time.Since(validationStart)
	
	// Ensure any failing colors are adjusted
	if validationResult.FailingColors > 0 {
		adjustedColors := g.validator.EnsureContrast(allColors)
		primary = adjustedColors[0]
		background = adjustedColors[1]
		foreground = adjustedColors[2]
		// Re-validate after adjustment
		validationResult = g.validator.Validate(adjustedColors)
	}
	
	totalTime := time.Since(startTime)
	
	// Create complete theme
	theme := &Theme{
		Name:       themeName,
		SourcePath: "", // Will be set by caller if from file
		IsLight:    isLight,
		Primary:    primary,
		Background: background,
		Foreground: foreground,
		Palette:    pipelineResult.Palette,
		Metadata: ThemeMetadata{
			GenerationMode:    getSynthesisModeString(pipelineResult.Mode),
			Strategy:          pipelineResult.Strategy,
			BaseColor:         pipelineResult.BaseColor,
			ExtractedColors:   pipelineResult.ExtractedColors,
			SynthesizedColors: pipelineResult.SynthesizedColors,
			Performance: PerformanceMetrics{
				ExtractionTime: extractionTime,
				SynthesisTime:  synthesisTime,
				ValidationTime: validationTime,
				TotalTime:      totalTime,
				ImageSize:      imageSize,
			},
			Validation: validationResult,
			Generated:  startTime,
		},
	}
	
	return theme, nil
}

// GenerateFromFile is a convenience method for generating themes from image files.
// It handles file loading and automatically derives the theme name from the filename.
func (g *Generator) GenerateFromFile(imagePath string, mode ThemeMode, overrides ColorOverrides) (*Theme, error) {
	// Load image
	img, err := extractor.LoadImage(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load image %s: %w", imagePath, err)
	}
	
	// Extract filename without extension for theme name
	filename := filepath.Base(imagePath)
	ext := filepath.Ext(filename)
	themeName := strings.TrimSuffix(filename, ext)
	
	// Create config
	config := ThemeConfig{
		SourceImage: img,
		Mode:        mode,
		Overrides:   overrides,
		Name:        themeName,
	}
	
	// Generate theme
	theme, err := g.GenerateTheme(config)
	if err != nil {
		return nil, err
	}
	
	// Set source path
	theme.SourcePath = imagePath
	
	return theme, nil
}

// deriveSemanticColors extracts primary, background, and foreground colors
// from the generated palette based on theme mode and color theory.
func (g *Generator) deriveSemanticColors(
	palette []*color.Color, 
	baseColor *color.Color, 
	isLight bool,
) (*color.Color, *color.Color, *color.Color) {
	
	// Use base color as primary, fallback to first palette color
	primary := baseColor
	if primary == nil && len(palette) > 0 {
		primary = palette[0]
	}
	if primary == nil {
		// Ultimate fallback
		primary = color.NewRGB(100, 100, 200)
	}
	
	// Derive background and foreground based on mode
	var background, foreground *color.Color
	
	if isLight {
		// Light theme: light background, dark foreground
		h, s, _ := primary.HSL()
		
		// Light background with subtle tint of primary hue
		background = color.NewHSL(h, s*0.1, 0.97) // Very light with minimal saturation
		
		// Dark foreground that contrasts well
		foreground = color.NewHSL(h, s*0.2, 0.15) // Dark with slight primary tint
		
	} else {
		// Dark theme: dark background, light foreground  
		h, s, _ := primary.HSL()
		
		// Dark background with subtle tint of primary hue
		background = color.NewHSL(h, s*0.3, 0.08) // Very dark with some saturation
		
		// Light foreground that contrasts well
		foreground = color.NewHSL(h, s*0.1, 0.92) // Light with minimal primary tint
	}
	
	return primary, background, foreground
}

// applyOverrides applies user color overrides while maintaining WCAG compliance.
func (g *Generator) applyOverrides(
	overrides ColorOverrides,
	primary, background, foreground *color.Color,
) (*color.Color, *color.Color, *color.Color, error) {
	
	// Apply overrides
	if overrides.Primary != nil {
		primary = overrides.Primary
	}
	if overrides.Background != nil {
		background = overrides.Background
	}
	if overrides.Foreground != nil {
		foreground = overrides.Foreground
	}
	
	// Create temporary validator with the override background
	tempValidator := palette.NewPaletteValidator(background)
	
	// Check critical contrast ratios
	colors := []*color.Color{primary, foreground}
	validationResult := tempValidator.Validate(colors)
	
	// If validation fails, auto-adjust colors
	if validationResult.FailingColors > 0 {
		adjustedColors := tempValidator.EnsureContrast(colors)
		primary = adjustedColors[0]
		foreground = adjustedColors[1]
		
		// Note: We don't adjust background from overrides to respect user intent
		// The foreground is adjusted instead to maintain readability
	}
	
	return primary, background, foreground, nil
}

// getBackgroundColor determines the background color for validation,
// considering user overrides and theme mode defaults.
func (g *Generator) getBackgroundColor(overrides ColorOverrides) *color.Color {
	if overrides.Background != nil {
		return overrides.Background
	}
	
	// Default white background for validation
	return color.NewRGB(255, 255, 255)
}

// getSynthesisModeString converts SynthesisMode to string for theme metadata
func getSynthesisModeString(sm palette.SynthesisMode) string {
	switch sm {
	case palette.ModeExtract:
		return "extract"
	case palette.ModeHybrid:
		return "hybrid"
	case palette.ModeSynthesize:
		return "synthesize"
	default:
		return "unknown"
	}
}