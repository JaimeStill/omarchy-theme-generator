// Package processor provides unified image processing and color analysis
// for theme generation. It implements a single-pass pipeline that extracts,
// analyzes, and organizes colors by their intended role in theme generation.
//
// The processor combines what were previously separate analysis, extraction,
// and strategy packages into a cohesive processing pipeline optimized for
// performance and maintainability.
//
// Key Features:
//   - Single-pass processing pipeline
//   - ColorProfile composition with embedded ImageColors
//   - Role-based color assignment (background/foreground/primary/secondary/accent)
//   - Integrated analysis (grayscale, monochromatic, color scheme detection)
//   - Theme mode detection based on luminance analysis
//   - WCAG compliance validation with automatic fallbacks
//   - Performance optimization: <2s for 4K images, <100MB memory
//
// Usage:
//
//	settings := settings.DefaultSettings()
//	processor := processor.New(settings)
//	
//	profile, err := processor.ProcessImage(img)
//	if err != nil {
//	    return err
//	}
//	
//	// Access extracted colors by role
//	bg := profile.Colors.Background
//	fg := profile.Colors.Foreground
//	primary := profile.Colors.Primary
//	
//	// Access analysis metadata
//	mode := profile.Mode // Light or Dark
//	isGrayscale := profile.IsGrayscale
//
// The processor enforces the settings-as-methods architectural pattern,
// requiring all operations to be performed through configured processor
// instances rather than package-level functions.
package processor