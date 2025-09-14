// Package chromatic provides color theory algorithms and perceptual analysis
// for the theme generation pipeline. It implements color similarity detection,
// contrast ratio calculations, and accessibility compliance validation.
//
// The package serves as the foundation layer for color processing, providing
// accurate perceptual distance calculations using multiple color spaces and
// specialized handling for neutral colors.
//
// Key Features:
//   - Perceptual color similarity using LAB color space
//   - Specialized neutral color clustering with lightness thresholds
//   - WCAG 2.1 accessibility compliance (AA/AAA levels)
//   - Multiple distance metrics (RGB, HSL, LAB)
//   - Hue analysis and variance calculations
//
// Color Similarity:
//
// The ColorsSimilar method implements a two-tier approach:
//   - Neutral colors (low saturation) use lightness difference thresholds
//   - Saturated colors use LAB color space distance for perceptual accuracy
//
// Usage:
//
//	settings := settings.DefaultSettings()
//	chroma := chromatic.NewChroma(settings)
//
//	// Test color similarity for clustering
//	if chroma.ColorsSimilar(color1, color2) {
//	    // Colors should be clustered together
//	}
//
//	// Check accessibility compliance
//	if chromatic.IsAccessible(fg, bg, chromatic.AA) {
//	    // Meets WCAG 2.1 AA requirements
//	}
//
//	// Calculate perceptual distance
//	distance := chromatic.DistanceLAB(color1, color2)
//
// The package follows the settings-as-methods pattern, requiring configuration
// through the Chroma struct for operations requiring thresholds.
package chromatic
