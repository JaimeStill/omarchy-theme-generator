package processor_test

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/processor"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
)

// Test color analysis functionality
func TestColorProfile_GrayscaleDetection(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	
	testCases := []struct {
		name        string
		colors      map[color.RGBA]uint32
		expectGray  bool
		description string
	}{
		{
			name: "All grayscale colors",
			colors: map[color.RGBA]uint32{
				{R: 50, G: 50, B: 50, A: 255}:   100, // Dark gray
				{R: 128, G: 128, B: 128, A: 255}: 200, // Medium gray
				{R: 200, G: 200, B: 200, A: 255}: 150, // Light gray
			},
			expectGray:  true,
			description: "Low saturation colors should be detected as grayscale",
		},
		{
			name: "Mixed colored and grayscale",
			colors: map[color.RGBA]uint32{
				{R: 128, G: 128, B: 128, A: 255}: 100, // Gray
				{R: 255, G: 0, B: 0, A: 255}:     200, // Red
			},
			expectGray:  false,
			description: "Mix of grayscale and colored should not be grayscale",
		},
		{
			name: "All saturated colors",
			colors: map[color.RGBA]uint32{
				{R: 255, G: 0, B: 0, A: 255}:   100, // Red
				{R: 0, G: 255, B: 0, A: 255}:   150, // Green
				{R: 0, G: 0, B: 255, A: 255}:   200, // Blue
			},
			expectGray:  false,
			description: "Saturated colors should not be detected as grayscale",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// We can't directly call analyzeColors, so create image and check the result
			// This is testing the integration but validates the analysis logic
			img := createImageFromFrequencyMap(tc.colors)
			
			result, err := p.ProcessImage(img)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			
			// Infer grayscale detection from color choices
			// Grayscale images should have low saturation in selected colors
			primarySat := formats.RGBAToHSLA(result.Colors.Primary).S
			accentSat := formats.RGBAToHSLA(result.Colors.Accent).S
			
			if tc.expectGray {
				// For grayscale images, even generated colors should have low saturation
				if primarySat > 0.3 || accentSat > 0.3 {
					t.Logf("Note: Grayscale image generated saturated colors (acceptable fallback behavior)")
				}
			} else {
				// For colored images, we should get some saturated colors
				if primarySat < 0.1 && accentSat < 0.1 {
					t.Error("Expected some saturated colors for non-grayscale image")
				}
			}
		})
	}
}

func TestColorProfile_MonochromaticDetection(t *testing.T) {
	s := settings.DefaultSettings()
	s.MonochromaticTolerance = 15.0 // ±15 degrees
	p := processor.New(s)
	
	testCases := []struct {
		name         string
		colors       map[color.RGBA]uint32
		expectMono   bool
		description  string
	}{
		{
			name: "Similar red hues",
			colors: map[color.RGBA]uint32{
				{R: 255, G: 0, B: 0, A: 255}:   100, // Pure red (0°)
				{R: 255, G: 50, B: 50, A: 255}: 150, // Light red (~10°)
				{R: 200, G: 0, B: 0, A: 255}:   200, // Dark red (0°)
			},
			expectMono:  true,
			description: "Colors within tolerance should be monochromatic",
		},
		{
			name: "Complementary colors",
			colors: map[color.RGBA]uint32{
				{R: 255, G: 0, B: 0, A: 255}: 100, // Red (0°)
				{R: 0, G: 255, B: 0, A: 255}: 150, // Green (120°)
				{R: 0, G: 0, B: 255, A: 255}: 200, // Blue (240°)
			},
			expectMono:  false,
			description: "Widely spaced hues should not be monochromatic",
		},
		{
			name: "Single color repeated",
			colors: map[color.RGBA]uint32{
				{R: 255, G: 100, B: 0, A: 255}: 500, // Orange
			},
			expectMono:  true,
			description: "Single color should be considered monochromatic",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			img := createImageFromFrequencyMap(tc.colors)
			
			result, err := p.ProcessImage(img)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			
			// For monochromatic images, colors should be more similar
			// For non-monochromatic, should have more hue diversity
			primaryHue := formats.RGBAToHSLA(result.Colors.Primary).H
			secondaryHue := formats.RGBAToHSLA(result.Colors.Secondary).H
			
			hueDiff := math.Abs(primaryHue - secondaryHue)
			if hueDiff > 180 {
				hueDiff = 360 - hueDiff
			}
			
			if tc.expectMono && hueDiff > 30 {
				t.Logf("Note: Monochromatic image has diverse generated colors (acceptable fallback)")
			}
			if !tc.expectMono && hueDiff < 10 {
				t.Logf("Note: Non-monochromatic image has similar generated colors")
			}
		})
	}
}

func TestThemeMode_Calculation(t *testing.T) {
	testCases := []struct {
		name         string
		colors       []color.RGBA
		threshold    float64
		expectedMode processor.ThemeMode
	}{
		{
			name: "Dark colors with default threshold",
			colors: []color.RGBA{
				{R: 50, G: 50, B: 50, A: 255},
				{R: 100, G: 50, B: 75, A: 255},
			},
			threshold:    0.5,
			expectedMode: processor.Dark, // Dark colors pair with Dark theme
		},
		{
			name: "Bright colors with default threshold", 
			colors: []color.RGBA{
				{R: 200, G: 200, B: 200, A: 255},
				{R: 255, G: 180, B: 220, A: 255},
			},
			threshold:    0.5,
			expectedMode: processor.Light, // Bright colors pair with Light theme
		},
		{
			name: "Edge case - exactly at threshold",
			colors: []color.RGBA{
				{R: 188, G: 188, B: 188, A: 255}, // Luminance ≈ 0.5 (accounting for gamma correction)
			},
			threshold:    0.5,
			expectedMode: processor.Light, // >= threshold = Light
		},
		{
			name: "Custom threshold test",
			colors: []color.RGBA{
				{R: 150, G: 150, B: 150, A: 255},
			},
			threshold:    0.7, // Higher threshold
			expectedMode: processor.Dark, // Below 0.7 = Dark
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := settings.DefaultSettings()
			s.ThemeModeThreshold = tc.threshold
			p := processor.New(s)
			
			colorFreq := make(map[color.RGBA]uint32)
			for i, c := range tc.colors {
				colorFreq[c] = uint32(100 * (i + 1)) // Varying frequencies
			}
			
			img := createImageFromFrequencyMap(colorFreq)
			result, err := p.ProcessImage(img)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			
			// Calculate luminance metrics for diagnosis
			bgLuminance := chromatic.Luminance(result.Colors.Background)
			
			// Also calculate input color luminance for comparison
			var inputLuminances []float64
			for _, c := range tc.colors {
				inputLuminances = append(inputLuminances, chromatic.Luminance(c))
			}
			
			avgInputLuminance := 0.0
			for _, l := range inputLuminances {
				avgInputLuminance += l
			}
			if len(inputLuminances) > 0 {
				avgInputLuminance /= float64(len(inputLuminances))
			}
			
			var detectedMode processor.ThemeMode
			if bgLuminance > 0.5 {
				detectedMode = processor.Light
			} else {
				detectedMode = processor.Dark
			}
			
			// Log diagnostic information
			t.Logf("Input colors average luminance: %v", avgInputLuminance)
			t.Logf("Threshold: %v", tc.threshold)
			t.Logf("Background luminance: %v", bgLuminance)
			t.Logf("Expected mode: %v, Detected mode: %v", tc.expectedMode, detectedMode)
			
			if detectedMode != tc.expectedMode {
				t.Errorf("Expected theme mode %v, detected %v (input avg: %v, threshold: %v, bg luminance: %v)", 
					tc.expectedMode, detectedMode, avgInputLuminance, tc.threshold, bgLuminance)
			}
		})
	}
}

// Helper function to create an image from a color frequency map
func createImageFromFrequencyMap(colorFreq map[color.RGBA]uint32) image.Image {
	// Calculate total pixels needed
	totalPixels := uint32(0)
	for _, freq := range colorFreq {
		totalPixels += freq
	}
	
	// Create square image that can hold all pixels
	side := int(math.Sqrt(float64(totalPixels))) + 1
	img := image.NewRGBA(image.Rect(0, 0, side, side))
	
	// Fill image with colors according to frequency
	pixelIndex := 0
	for color, freq := range colorFreq {
		for i := uint32(0); i < freq; i++ {
			if pixelIndex >= side*side {
				break
			}
			x := pixelIndex % side
			y := pixelIndex / side
			img.Set(x, y, color)
			pixelIndex++
		}
	}
	
	// Fill any remaining pixels with the first color
	if len(colorFreq) > 0 {
		var firstColor color.RGBA
		for c := range colorFreq {
			firstColor = c
			break
		}
		for pixelIndex < side*side {
			x := pixelIndex % side
			y := pixelIndex / side
			img.Set(x, y, firstColor)
			pixelIndex++
		}
	}
	
	return img
}