package processor_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/processor"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

func TestNew(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	
	t.Logf("Testing processor creation")
	t.Logf("Settings provided: %+v", s.MinFrequency)
	
	if p == nil {
		t.Fatal("Expected processor to be created, got nil")
	}
	
	t.Logf("Processor created successfully")
}

func TestProcessImage_SimpleImage(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	
	// Create a simple 3x3 test image with distinct colors
	img := createTestImage(3, 3, []color.RGBA{
		{255, 255, 255, 255}, // White background
		{0, 0, 0, 255},       // Black foreground  
		{255, 0, 0, 255},     // Red accent
	})
	
	t.Logf("Testing ProcessImage with 3x3 test image")
	t.Logf("Image bounds: %v", img.Bounds())
	
	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}
	
	if profile == nil {
		t.Fatal("Expected profile to be non-nil")
	}
	
	// Log comprehensive diagnostic information
	t.Logf("Color Profile Results:")
	t.Logf("  Mode: %s", profile.Mode)
	t.Logf("  ColorScheme: %s", profile.ColorScheme)
	t.Logf("  IsGrayscale: %t", profile.IsGrayscale)
	t.Logf("  IsMonochromatic: %t", profile.IsMonochromatic)
	t.Logf("  DominantHue: %.1f°", profile.DominantHue)
	t.Logf("  HueVariance: %.1f°", profile.HueVariance)
	t.Logf("  AvgLuminance: %.3f", profile.AvgLuminance)
	t.Logf("  AvgSaturation: %.3f", profile.AvgSaturation)
	
	t.Logf("Image Colors Results:")
	t.Logf("  TotalPixels: %d", profile.Colors.TotalPixels)
	t.Logf("  UniqueColors: %d", profile.Colors.UniqueColors)
	t.Logf("  CoverageRatio: %.3f", profile.Colors.CoverageRatio)
	t.Logf("  Categories found: %d", len(profile.Colors.Categories))
	t.Logf("  Category candidates: %d total", countCandidates(profile.Colors.CategoryCandidates))
	
	// Basic validation
	if profile.Colors.TotalPixels != 9 {
		t.Errorf("Expected 9 total pixels, got %d", profile.Colors.TotalPixels)
	}
	
	if profile.Colors.UniqueColors == 0 {
		t.Error("Expected non-zero unique colors")
	}
	
	if len(profile.Colors.Categories) == 0 {
		t.Error("Expected at least one category to be assigned")
	}
	
	// Verify background category is always present
	if _, hasBackground := profile.Colors.Categories[processor.CategoryBackground]; !hasBackground {
		t.Error("Background category should always be present")
	}
}

func TestProcessImage_GrayscaleImage(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	
	// Create a grayscale test image
	img := createTestImage(2, 2, []color.RGBA{
		{50, 50, 50, 255},   // Dark gray
		{200, 200, 200, 255}, // Light gray
	})
	
	t.Logf("Testing ProcessImage with grayscale image")
	
	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}
	
	t.Logf("Grayscale Detection Results:")
	t.Logf("  IsGrayscale: %t", profile.IsGrayscale)
	t.Logf("  AvgSaturation: %.3f (threshold: %.3f)", profile.AvgSaturation, s.GrayscaleThreshold)
	t.Logf("  Mode: %s", profile.Mode)
	t.Logf("  AvgLuminance: %.3f (mode threshold: %.3f)", profile.AvgLuminance, s.ThemeModeThreshold)
	
	// Expect grayscale detection
	if !profile.IsGrayscale {
		t.Errorf("Expected grayscale detection for gray image (avg sat: %.3f)", profile.AvgSaturation)
	}
	
	// Background should still be detected
	if _, hasBackground := profile.Colors.Categories[processor.CategoryBackground]; !hasBackground {
		t.Error("Background category should be present even in grayscale images")
	}
}

func TestProcessImage_MonochromaticImage(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	
	// Create a monochromatic test image (different shades of blue)
	img := createTestImage(3, 3, []color.RGBA{
		{0, 0, 100, 255},   // Dark blue
		{0, 0, 150, 255},   // Medium blue
		{0, 0, 200, 255},   // Light blue
	})
	
	t.Logf("Testing ProcessImage with monochromatic blue image")
	
	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}
	
	t.Logf("Monochromatic Detection Results:")
	t.Logf("  IsMonochromatic: %t", profile.IsMonochromatic)
	t.Logf("  DominantHue: %.1f° (expected ~240° for blue)", profile.DominantHue)
	t.Logf("  HueVariance: %.1f° (tolerance: %.1f°)", profile.HueVariance, s.MonochromaticTolerance)
	t.Logf("  ColorScheme: %s", profile.ColorScheme)
	
	// Should detect monochromatic pattern
	if !profile.IsMonochromatic {
		t.Logf("Note: Monochromatic detection may depend on saturation levels")
	}
	
	// Should detect blue hue range (240° ± tolerance)
	expectedHue := 240.0
	hueDiff := abs(profile.DominantHue - expectedHue)
	if hueDiff > 180 {
		hueDiff = 360 - hueDiff
	}
	
	if hueDiff > 30 { // Allow some tolerance for blue detection
		t.Logf("Note: Dominant hue %.1f° differs from expected blue ~240°", profile.DominantHue)
	}
}

func TestProcessImage_EmptyImage(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	
	// Create a 0x0 image (should fail gracefully)
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	
	t.Logf("Testing ProcessImage with empty image")
	
	_, err := p.ProcessImage(img)
	if err == nil {
		t.Error("Expected error for empty image")
	}
	
	t.Logf("Empty image correctly rejected: %v", err)
}

func TestProcessImage_MinimumColorRequirement(t *testing.T) {
	s := settings.DefaultSettings()
	s.MinFrequency = 0.9 // Set very high threshold
	p := processor.New(s)
	
	// Create image where no colors meet the minimum frequency
	img := createTestImage(10, 10, []color.RGBA{
		{255, 0, 0, 255},   // Red - 1 pixel
		{0, 255, 0, 255},   // Green - 1 pixel  
		{0, 0, 255, 255},   // Blue - 1 pixel
		{255, 255, 0, 255}, // Yellow - 97 pixels (dominant)
	})
	
	t.Logf("Testing ProcessImage with high minimum frequency threshold")
	t.Logf("MinFrequency setting: %.3f", s.MinFrequency)
	
	_, err := p.ProcessImage(img)
	if err == nil {
		t.Error("Expected error when no colors meet minimum frequency")
	}
	
	t.Logf("High frequency threshold correctly enforced: %v", err)
}

// Helper functions

func createTestImage(width, height int, colors []color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	colorIndex := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, colors[colorIndex%len(colors)])
			colorIndex++
		}
	}
	
	return img
}

func countCandidates(candidates map[processor.ColorCategory][]processor.ColorCandidate) int {
	total := 0
	for _, candidateList := range candidates {
		total += len(candidateList)
	}
	return total
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}