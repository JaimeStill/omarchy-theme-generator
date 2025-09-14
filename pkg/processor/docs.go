// Package processor provides unified image processing and color analysis
// for theme generation. It implements a characteristic-based clustering pipeline
// that extracts, analyzes, and organizes colors by visual properties and frequency.
//
// The processor combines color extraction, frequency analysis, and characteristic
// classification into a cohesive pipeline optimized for performance and
// downstream theme generation flexibility.
//
// Key Features:
//   - Characteristic-based ColorCluster system with lightness, saturation, and hue grouping
//   - Frequency-weighted color extraction sorted by visual importance
//   - Theme mode detection based on weighted luminance analysis
//   - Configurable clustering thresholds and UI color limits
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
//	// Access color clusters sorted by weight (highest first)
//	colors := profile.Colors  // []ColorCluster
//	dominant := colors[0]     // Most prominent color
//
//	// Access cluster characteristics
//	if dominant.IsNeutral && dominant.IsDark {
//	    // Handle dark neutral color
//	}
//
//	// Access analysis metadata
//	mode := profile.Mode       // Light or Dark
//	hasColor := profile.HasColor // False for grayscale images
//
// The processor enforces the settings-as-methods architectural pattern,
// requiring all operations to be performed through configured processor
// instances rather than package-level functions.
package processor