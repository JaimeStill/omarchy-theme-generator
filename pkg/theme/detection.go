package theme

import (
	"image"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
)

// ModeDetector analyzes images and colors to determine optimal light/dark theme mode.
// It uses WCAG-accurate relative luminance calculations for perceptually correct results.
type ModeDetector struct {
	// LuminanceThreshold is the threshold for classifying as light vs dark theme
	// 0.5 means images with average luminance above 50% become light themes
	LuminanceThreshold float64
	
	// PrimaryWeight is how much the primary color influences mode detection
	// relative to the overall image luminance (0.0 to 1.0)
	PrimaryWeight float64
}

// NewModeDetector creates a detector with sensible defaults for theme mode detection.
func NewModeDetector() *ModeDetector {
	return &ModeDetector{
		LuminanceThreshold: 0.5,  // 50% luminance threshold
		PrimaryWeight:      0.3,  // Primary color has 30% weight in decision
	}
}

// DetectFromImage analyzes an image to determine the optimal theme mode.
// It calculates average WCAG relative luminance and considers image variance
// for more accurate light/dark classification.
func (md *ModeDetector) DetectFromImage(img image.Image) ThemeMode {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	if width == 0 || height == 0 {
		return ModeDark // Default for invalid images
	}
	
	// Calculate average luminance using WCAG-accurate relative luminance
	totalLuminance := 0.0
	pixelCount := 0
	
	// Sample the image - for performance, we can sample every Nth pixel for large images
	stepX := max(1, width/200)  // Sample at most 200 points across width
	stepY := max(1, height/200) // Sample at most 200 points across height
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y += stepY {
		for x := bounds.Min.X; x < bounds.Max.X; x += stepX {
			r, g, b, _ := img.At(x, y).RGBA()
			
			// Convert to 0-255 range
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			
			// Create color and get WCAG relative luminance
			c := color.NewRGB(r8, g8, b8)
			luminance := c.RelativeLuminance()
			
			totalLuminance += luminance
			pixelCount++
		}
	}
	
	if pixelCount == 0 {
		return ModeDark // Default for empty images
	}
	
	averageLuminance := totalLuminance / float64(pixelCount)
	
	// Calculate variance for low-contrast image detection
	totalVariance := 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y += stepY {
		for x := bounds.Min.X; x < bounds.Max.X; x += stepX {
			r, g, b, _ := img.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			
			c := color.NewRGB(r8, g8, b8)
			luminance := c.RelativeLuminance()
			diff := luminance - averageLuminance
			totalVariance += diff * diff
		}
	}
	
	variance := totalVariance / float64(pixelCount)
	
	// Adjust threshold based on image characteristics
	adjustedThreshold := md.LuminanceThreshold
	
	// Low contrast images (variance < 0.01) bias toward light mode for readability
	if variance < 0.01 {
		adjustedThreshold += 0.1
	}
	
	// High contrast images (variance > 0.1) can use standard threshold
	if variance > 0.1 {
		adjustedThreshold = md.LuminanceThreshold
	}
	
	// Make decision based on adjusted luminance
	if averageLuminance > adjustedThreshold {
		return ModeLight
	}
	
	return ModeDark
}

// DetectFromColors determines theme mode based on dominant colors extracted from an image.
// This is used when we have color extraction results but want to re-evaluate mode.
func (md *ModeDetector) DetectFromColors(colors []*color.Color) ThemeMode {
	if len(colors) == 0 {
		return ModeDark // Default for empty color list
	}
	
	// Calculate weighted average luminance of the colors
	totalLuminance := 0.0
	totalWeight := 0.0
	
	for i, c := range colors {
		if c == nil {
			continue
		}
		
		// Weight colors by their position (earlier colors are more dominant)
		weight := 1.0 / float64(i+1)
		totalLuminance += c.RelativeLuminance() * weight
		totalWeight += weight
	}
	
	if totalWeight == 0 {
		return ModeDark // Default for invalid colors
	}
	
	averageLuminance := totalLuminance / totalWeight
	
	if averageLuminance > md.LuminanceThreshold {
		return ModeLight
	}
	
	return ModeDark
}

// DetectWithPrimary determines theme mode considering both image and primary color.
// This provides the most accurate mode detection by combining image analysis
// with the specific primary color that will be used in the theme.
func (md *ModeDetector) DetectWithPrimary(img image.Image, primaryColor *color.Color) ThemeMode {
	// Get base mode from image analysis
	imageMode := md.DetectFromImage(img)
	
	if primaryColor == nil {
		return imageMode
	}
	
	// Get primary color luminance
	primaryLuminance := primaryColor.RelativeLuminance()
	
	// If image suggests light mode but primary is very dark, consider dark mode
	if imageMode == ModeLight && primaryLuminance < 0.2 {
		return ModeDark
	}
	
	// If image suggests dark mode but primary is very light, consider light mode
	if imageMode == ModeDark && primaryLuminance > 0.8 {
		return ModeLight
	}
	
	// For borderline cases, use weighted decision
	imageLuminance := md.calculateImageAverageLuminance(img)
	
	// Combine image and primary color luminance with weighting
	combinedLuminance := imageLuminance*(1-md.PrimaryWeight) + primaryLuminance*md.PrimaryWeight
	
	if combinedLuminance > md.LuminanceThreshold {
		return ModeLight
	}
	
	return ModeDark
}

// calculateImageAverageLuminance is a helper that efficiently calculates
// average luminance for an image (used by DetectWithPrimary).
func (md *ModeDetector) calculateImageAverageLuminance(img image.Image) float64 {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	if width == 0 || height == 0 {
		return 0.0
	}
	
	// Sample for performance
	stepX := max(1, width/100)
	stepY := max(1, height/100)
	
	totalLuminance := 0.0
	pixelCount := 0
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y += stepY {
		for x := bounds.Min.X; x < bounds.Max.X; x += stepX {
			r, g, b, _ := img.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			
			c := color.NewRGB(r8, g8, b8)
			totalLuminance += c.RelativeLuminance()
			pixelCount++
		}
	}
	
	if pixelCount == 0 {
		return 0.0
	}
	
	return totalLuminance / float64(pixelCount)
}

// max returns the maximum of two integers (Go 1.25 doesn't have builtin max for int)
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}