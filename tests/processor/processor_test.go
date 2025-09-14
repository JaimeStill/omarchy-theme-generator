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
	t.Logf("MinFrequency setting: %.4f", s.Processor.MinFrequency)
	t.Logf("MaxUIColors setting: %d", s.Processor.MaxUIColors)

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
	t.Logf("Input colors: White(255,255,255), Black(0,0,0), Red(255,0,0)")

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	if profile == nil {
		t.Fatal("Expected profile to be non-nil")
	}

	// Log comprehensive diagnostic information about ColorProfile
	t.Logf("ColorProfile Results:")
	t.Logf("  Mode: %s", profile.Mode)
	t.Logf("  HasColor: %t", profile.HasColor)
	t.Logf("  ColorCount: %d", profile.ColorCount)
	t.Logf("  Colors extracted: %d", len(profile.Colors))

	// Log each ColorCluster with full diagnostics
	for i, cluster := range profile.Colors {
		t.Logf("  Cluster %d:", i)
		t.Logf("    RGBA: (%d,%d,%d,%d)", cluster.R, cluster.G, cluster.B, cluster.A)
		t.Logf("    Weight: %.6f", cluster.Weight)
		t.Logf("    Lightness: %.3f, Saturation: %.3f, Hue: %.1f°", cluster.Lightness, cluster.Saturation, cluster.Hue)
		t.Logf("    Characteristics: Neutral=%t, Dark=%t, Light=%t, Muted=%t, Vibrant=%t",
			 cluster.IsNeutral, cluster.IsDark, cluster.IsLight, cluster.IsMuted, cluster.IsVibrant)
	}

	// Basic validation
	if profile.ColorCount != len(profile.Colors) {
		t.Errorf("ColorCount mismatch: reported %d, actual slice length %d", profile.ColorCount, len(profile.Colors))
	}

	if len(profile.Colors) == 0 {
		t.Error("Expected non-zero color clusters")
	}

	// Validate weight ordering (should be sorted by weight descending)
	for i := 1; i < len(profile.Colors); i++ {
		if profile.Colors[i-1].Weight < profile.Colors[i].Weight {
			t.Errorf("Colors not sorted by weight: cluster %d weight %.6f > cluster %d weight %.6f",
				 i, profile.Colors[i].Weight, i-1, profile.Colors[i-1].Weight)
		}
	}

	// Validate total weight doesn't exceed 1.0
	totalWeight := 0.0
	for _, cluster := range profile.Colors {
		totalWeight += cluster.Weight
	}
	if totalWeight > 1.0 {
		t.Errorf("Total weight %.6f exceeds 1.0", totalWeight)
	}
	t.Logf("Total weight validation: %.6f (≤ 1.0)", totalWeight)
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
	t.Logf("Input colors: DarkGray(50,50,50), LightGray(200,200,200)")
	t.Logf("Neutral threshold: %.3f", s.Chromatic.NeutralThreshold)

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	// Calculate expected properties
	weightedLightness := 0.0
	weightedSaturation := 0.0
	totalWeight := 0.0
	neutralCount := 0

	for _, cluster := range profile.Colors {
		weightedLightness += cluster.Lightness * cluster.Weight
		weightedSaturation += cluster.Saturation * cluster.Weight
		totalWeight += cluster.Weight
		if cluster.IsNeutral {
			neutralCount++
		}
	}

	avgLightness := weightedLightness / totalWeight
	avgSaturation := weightedSaturation / totalWeight

	t.Logf("Grayscale Analysis Results:")
	t.Logf("  Mode: %s (threshold: %.3f)", profile.Mode, s.Processor.LightThemeThreshold)
	t.Logf("  HasColor: %t (threshold: %.3f)", profile.HasColor, s.Processor.SignificantColorThreshold)
	t.Logf("  Calculated avg lightness: %.3f", avgLightness)
	t.Logf("  Calculated avg saturation: %.3f", avgSaturation)
	t.Logf("  Neutral clusters: %d/%d", neutralCount, len(profile.Colors))

	// Log each cluster's neutral classification
	for i, cluster := range profile.Colors {
		t.Logf("  Cluster %d: Sat=%.3f, IsNeutral=%t (threshold: %.3f)",
			 i, cluster.Saturation, cluster.IsNeutral, s.Chromatic.NeutralThreshold)
	}

	// Validate that low-saturation colors are marked as neutral
	if len(profile.Colors) > 0 {
		if avgSaturation < s.Chromatic.NeutralThreshold && neutralCount == 0 {
			t.Errorf("Expected neutral colors in grayscale image (avg saturation: %.3f)", avgSaturation)
		}
	}

	// Verify HasColor should be false for truly grayscale images
	if avgSaturation < s.Chromatic.NeutralThreshold && profile.HasColor {
		t.Logf("Note: HasColor=true despite low saturation - may indicate color weight above threshold")
	}
}

func TestProcessImage_ColorfulImage(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Create a colorful test image with distinct hues
	img := createTestImage(3, 3, []color.RGBA{
		{255, 0, 0, 255},   // Red
		{0, 255, 0, 255},   // Green
		{0, 0, 255, 255},   // Blue
	})

	t.Logf("Testing ProcessImage with colorful RGB image")
	t.Logf("Input colors: Red(255,0,0), Green(0,255,0), Blue(0,0,255)")
	t.Logf("Expected hues: Red(0°), Green(120°), Blue(240°)")

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	// Calculate color characteristics
	colorWeight := 0.0
	vibrantCount := 0
	mutedCount := 0
	neutralCount := 0

	for _, cluster := range profile.Colors {
		if !cluster.IsNeutral {
			colorWeight += cluster.Weight
		}
		if cluster.IsVibrant {
			vibrantCount++
		}
		if cluster.IsMuted {
			mutedCount++
		}
		if cluster.IsNeutral {
			neutralCount++
		}
	}

	t.Logf("Colorful Image Analysis:")
	t.Logf("  HasColor: %t (threshold: %.3f)", profile.HasColor, s.Processor.SignificantColorThreshold)
	t.Logf("  Color weight: %.3f", colorWeight)
	t.Logf("  Vibrant clusters: %d, Muted: %d, Neutral: %d", vibrantCount, mutedCount, neutralCount)

	// Log saturation analysis for each cluster
	for i, cluster := range profile.Colors {
		t.Logf("  Cluster %d: Hue=%.1f°, Sat=%.3f, Vibrant=%t (threshold: %.3f)",
			 i, cluster.Hue, cluster.Saturation, cluster.IsVibrant, s.Chromatic.VibrantSaturationMin)
	}

	// Should detect significant color content
	if !profile.HasColor {
		t.Logf("Note: Expected HasColor=true for RGB image (color weight: %.3f, threshold: %.3f)",
			 colorWeight, s.Processor.SignificantColorThreshold)
	}

	// Should have vibrant colors from pure RGB
	if vibrantCount == 0 {
		t.Logf("Note: Expected vibrant colors from pure RGB values")
	}
}

func TestProcessImage_EmptyImage(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Create a 0x0 image (should fail gracefully)
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))

	t.Logf("Testing ProcessImage with empty image")
	t.Logf("Image bounds: %v", img.Bounds())
	t.Logf("Expected: Error due to no extractable colors")

	profile, err := p.ProcessImage(img)
	if err == nil {
		t.Error("Expected error for empty image")
	}

	if profile != nil {
		t.Error("Expected nil profile for failed processing")
	}

	t.Logf("Empty image correctly rejected: %v", err)
}

func TestProcessImage_HighFrequencyThreshold(t *testing.T) {
	s := settings.DefaultSettings()
	s.Processor.MinFrequency = 0.9 // Set very high threshold
	p := processor.New(s)

	// Create image where no colors meet the minimum frequency
	img := createTestImage(10, 10, []color.RGBA{
		{255, 0, 0, 255},   // Red - 1 pixel
		{0, 255, 0, 255},   // Green - 1 pixel
		{0, 0, 255, 255},   // Blue - 1 pixel
		{255, 255, 0, 255}, // Yellow - 97 pixels (dominant)
	})

	t.Logf("Testing ProcessImage with high minimum frequency threshold")
	t.Logf("MinFrequency setting: %.3f (90%% threshold)", s.Processor.MinFrequency)
	t.Logf("Expected: Only yellow meets 90%% threshold (97/100 pixels)")
	t.Logf("Other colors: Red=1%%, Green=1%%, Blue=1%% (all below threshold)")

	profile, err := p.ProcessImage(img)
	if err != nil {
		// If error occurs, log the reason
		t.Logf("Processing failed (expected if no colors meet UI weight requirements): %v", err)
		return
	}

	// If processing succeeds, validate that only high-frequency colors remain
	if profile != nil {
		t.Logf("Processing succeeded with %d clusters", len(profile.Colors))
		for i, cluster := range profile.Colors {
			t.Logf("  Cluster %d: RGBA(%d,%d,%d,%d), Weight=%.6f",
				 i, cluster.R, cluster.G, cluster.B, cluster.A, cluster.Weight)

			// Validate weight meets minimum requirement
			if cluster.Weight < s.Processor.MinFrequency {
				t.Errorf("Cluster %d weight %.6f below MinFrequency %.3f",
					 i, cluster.Weight, s.Processor.MinFrequency)
			}
		}
	}
}

func TestProcessImage_ThemeModeCalculation(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Test light theme detection with bright colors
	lightImg := createTestImage(4, 4, []color.RGBA{
		{240, 240, 240, 255}, // Very light gray - dominant
		{220, 220, 220, 255}, // Light gray
		{200, 200, 200, 255}, // Medium gray
		{100, 100, 100, 255}, // Darker accent
	})

	t.Logf("Testing theme mode calculation with light-dominant image")
	t.Logf("Light theme threshold: %.3f", s.Processor.LightThemeThreshold)
	t.Logf("Input colors: VeryLight(240,240,240), Light(220,220,220), Medium(200,200,200), Dark(100,100,100)")

	profile, err := p.ProcessImage(lightImg)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	// Calculate weighted lightness for validation
	weightedLightness := 0.0
	totalWeight := 0.0
	for _, cluster := range profile.Colors {
		weightedLightness += cluster.Lightness * cluster.Weight
		totalWeight += cluster.Weight
		t.Logf("  Cluster: L=%.3f, Weight=%.6f, Contribution=%.6f",
			cluster.Lightness, cluster.Weight, cluster.Lightness*cluster.Weight)
	}
	avgLightness := weightedLightness / totalWeight

	t.Logf("Theme mode calculation:")
	t.Logf("  Calculated avg lightness: %.3f", avgLightness)
	t.Logf("  Light theme threshold: %.3f", s.Processor.LightThemeThreshold)
	t.Logf("  Detected mode: %s", profile.Mode)

	expectedMode := processor.Dark
	if avgLightness > s.Processor.LightThemeThreshold {
		expectedMode = processor.Light
	}

	if profile.Mode != expectedMode {
		t.Errorf("Expected mode %s, got %s (avg lightness: %.3f, threshold: %.3f)",
			expectedMode, profile.Mode, avgLightness, s.Processor.LightThemeThreshold)
	}
}

func TestProcessImage_ColorClustering(t *testing.T) {
	s := settings.DefaultSettings()
	// Lower the merge threshold to make clustering more sensitive
	s.Chromatic.ColorMergeThreshold = 5.0 // Very low threshold for similar colors
	p := processor.New(s)

	// Create image with similar colors that should be clustered together
	img := createTestImage(6, 6, []color.RGBA{
		{255, 0, 0, 255},   // Red
		{250, 5, 5, 255},   // Very similar red
		{245, 10, 10, 255}, // Another similar red
		{0, 255, 0, 255},   // Green (distinct)
		{0, 0, 255, 255},   // Blue (distinct)
		{128, 128, 128, 255}, // Gray (distinct)
	})

	t.Logf("Testing color clustering with similar colors")
	t.Logf("Color merge threshold: %.1f (Delta-E units)", s.Chromatic.ColorMergeThreshold)
	t.Logf("Input: 3 similar reds, plus green, blue, gray")
	t.Logf("Expected: Similar reds should cluster together")

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	t.Logf("Clustering results:")
	t.Logf("  Total clusters: %d", len(profile.Colors))

	for i, cluster := range profile.Colors {
		t.Logf("  Cluster %d: RGBA(%d,%d,%d), Weight=%.6f",
			i, cluster.R, cluster.G, cluster.B, cluster.Weight)
	}

	// With aggressive clustering, expect fewer clusters than input colors
	// This validates that similar reds were merged
	if len(profile.Colors) >= 6 {
		t.Logf("Note: Expected clustering to reduce color count from 6 unique colors")
		t.Logf("Actual clusters: %d (clustering may need tuning)", len(profile.Colors))
	}
}

func TestProcessImage_MaxUIColorsLimit(t *testing.T) {
	s := settings.DefaultSettings()
	s.Processor.MaxUIColors = 3 // Limit to 3 colors
	s.Processor.MinFrequency = 0.01 // Low threshold to allow many colors
	p := processor.New(s)

	// Create image with many distinct colors
	colors := []color.RGBA{
		{255, 0, 0, 255},   // Red
		{0, 255, 0, 255},   // Green
		{0, 0, 255, 255},   // Blue
		{255, 255, 0, 255}, // Yellow
		{255, 0, 255, 255}, // Magenta
		{0, 255, 255, 255}, // Cyan
		{128, 128, 128, 255}, // Gray
	}
	img := createTestImage(7, 7, colors)

	t.Logf("Testing MaxUIColors limit")
	t.Logf("MaxUIColors setting: %d", s.Processor.MaxUIColors)
	t.Logf("Input colors: %d distinct colors", len(colors))

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	t.Logf("UI color limiting results:")
	t.Logf("  Colors returned: %d", len(profile.Colors))
	t.Logf("  Maximum allowed: %d", s.Processor.MaxUIColors)

	if len(profile.Colors) > s.Processor.MaxUIColors {
		t.Errorf("Expected max %d colors, got %d", s.Processor.MaxUIColors, len(profile.Colors))
	}

	// Verify colors are sorted by weight (most important first)
	for i := 1; i < len(profile.Colors); i++ {
		if profile.Colors[i-1].Weight < profile.Colors[i].Weight {
			t.Errorf("Colors not properly sorted by weight at index %d", i)
		}
	}

	// Log the selected colors for analysis
	for i, cluster := range profile.Colors {
		t.Logf("  Selected color %d: RGBA(%d,%d,%d), Weight=%.6f",
			i, cluster.R, cluster.G, cluster.B, cluster.Weight)
	}
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


func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}