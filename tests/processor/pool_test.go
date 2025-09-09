package processor_test

import (
	"image/color"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/processor"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

func TestBuildColorPool(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Create test colors with varying weights
	colors := []color.RGBA{
		{255, 0, 0, 255},     // Red - high weight
		{0, 255, 0, 255},     // Green - medium weight
		{0, 0, 255, 255},     // Blue - low weight
		{255, 255, 0, 255},   // Yellow - medium weight
		{128, 128, 128, 255}, // Gray - low weight
	}

	// Create test image with different frequencies
	img := createTestImage(10, 1, []color.RGBA{
		colors[0], colors[0], colors[0], colors[0], // Red: 4 pixels (40%)
		colors[1], colors[1], colors[1],            // Green: 3 pixels (30%)
		colors[2], colors[2],                       // Blue: 2 pixels (20%)
		colors[3],                                  // Yellow: 1 pixel (10%)
		// Gray: 0 pixels (excluded by frequency)
	})

	t.Logf("Testing ColorPool construction with %d unique colors", len(colors))
	t.Logf("Settings: DominantColorCount=%d", s.Extraction.DominantColorCount)

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	pool := profile.Pool

	t.Logf("ColorPool results:")
	t.Logf("  TotalPixels: %d", pool.TotalPixels)
	t.Logf("  UniqueColors: %d", pool.UniqueColors)
	t.Logf("  AllColors count: %d", len(pool.AllColors))
	t.Logf("  DominantColors count: %d", len(pool.DominantColors))

	// Basic structure validation
	if pool.TotalPixels == 0 {
		t.Error("TotalPixels should be non-zero")
	}

	if pool.UniqueColors == 0 {
		t.Error("UniqueColors should be non-zero")
	}

	if len(pool.AllColors) == 0 {
		t.Error("AllColors should be non-empty")
	}

	if len(pool.DominantColors) == 0 {
		t.Error("DominantColors should be non-empty")
	}

	// Verify dominant colors are subset of all colors
	if len(pool.DominantColors) > len(pool.AllColors) {
		t.Error("DominantColors should not exceed AllColors count")
	}

	// Verify dominant colors are highest weighted
	if len(pool.DominantColors) > 1 {
		for i := 0; i < len(pool.DominantColors)-1; i++ {
			if pool.DominantColors[i].Weight < pool.DominantColors[i+1].Weight {
				t.Error("DominantColors not sorted by weight (descending)")
				break
			}
		}
	}

	// Log dominant color weights for diagnostics
	dominantWeights := make([]float64, len(pool.DominantColors))
	for i, c := range pool.DominantColors {
		dominantWeights[i] = c.Weight
	}
	t.Logf("  Dominant color weights: %v", dominantWeights)

	// Verify groupings exist
	lightnessTotal := len(pool.ByLightness.Dark) + len(pool.ByLightness.Mid) + len(pool.ByLightness.Light)
	if lightnessTotal == 0 {
		t.Error("No lightness groupings created")
	}

	saturationTotal := len(pool.BySaturation.Gray) + len(pool.BySaturation.Muted) + 
					   len(pool.BySaturation.Normal) + len(pool.BySaturation.Vibrant)
	if saturationTotal == 0 {
		t.Error("No saturation groupings created")
	}

	// Hue groupings may be empty for grayscale images
	t.Logf("  Hue families: %d", len(pool.ByHue))

	// Verify statistics are calculated
	if pool.Statistics.HueVariance < 0 {
		t.Error("Statistics should be initialized (HueVariance >= 0)")
	}
}

func TestSelectDominant(t *testing.T) {
	s := settings.DefaultSettings()
	s.Extraction.DominantColorCount = 3 // Limit to 3 dominant colors
	p := processor.New(s)

	// Create test colors with distinct weights
	colors := []color.RGBA{
		{255, 0, 0, 255},   // Will be highest weight
		{0, 255, 0, 255},   // Will be second highest 
		{0, 0, 255, 255},   // Will be third highest
		{255, 255, 0, 255}, // Will be fourth
		{128, 0, 128, 255}, // Will be lowest
	}

	// Create image with different frequencies (8 pixels total)
	img := createTestImage(8, 1, []color.RGBA{
		colors[0], colors[0], colors[0], colors[0], // Red: 4 pixels (50%)
		colors[1], colors[1], colors[1],            // Green: 3 pixels (37.5%) 
		colors[2],                                  // Blue: 1 pixel (12.5%)
		// Yellow and Purple excluded by frequency threshold
	})

	t.Logf("Testing dominant color selection with limit=%d", s.Extraction.DominantColorCount)

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	dominant := profile.Pool.DominantColors

	t.Logf("Dominant selection results:")
	t.Logf("  All colors: %d", len(profile.Pool.AllColors))
	t.Logf("  Dominant colors: %d", len(dominant))

	// Should respect the limit
	expectedLimit := min(s.Extraction.DominantColorCount, len(profile.Pool.AllColors))
	if len(dominant) > expectedLimit {
		t.Errorf("Too many dominant colors: got %d, expected max %d", len(dominant), expectedLimit)
	}

	// Should be sorted by weight (descending)
	for i := 0; i < len(dominant)-1; i++ {
		if dominant[i].Weight < dominant[i+1].Weight {
			t.Errorf("Dominant colors not sorted: position %d (%.3f) < position %d (%.3f)",
				i, dominant[i].Weight, i+1, dominant[i+1].Weight)
		}
	}

	// Log weights for diagnostics
	weights := make([]float64, len(dominant))
	for i, c := range dominant {
		weights[i] = c.Weight
	}
	t.Logf("  Dominant weights: %v", weights)

	// First dominant color should have highest weight
	if len(dominant) > 0 && len(profile.Pool.AllColors) > 0 {
		maxWeight := float64(0)
		for _, c := range profile.Pool.AllColors {
			if c.Weight > maxWeight {
				maxWeight = c.Weight
			}
		}
		if dominant[0].Weight != maxWeight {
			t.Errorf("First dominant color weight %.3f != max weight %.3f", dominant[0].Weight, maxWeight)
		}
	}
}

func TestSelectDominantEdgeCases(t *testing.T) {
	s := settings.DefaultSettings()

	tests := []struct {
		name       string
		count      int
		numColors  int
		expectAll  bool
	}{
		{"Zero count", 0, 5, true},         // Should return all
		{"Negative count", -1, 5, true},    // Should return all  
		{"Count equals colors", 5, 5, true}, // Should return all
		{"Count exceeds colors", 10, 5, true}, // Should return all
		{"Normal selection", 3, 5, false},   // Should return subset
	}

	baseColors := []color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255}, 
		{0, 0, 255, 255},
		{255, 255, 0, 255},
		{128, 128, 128, 255},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.Extraction.DominantColorCount = tt.count
			testP := processor.New(s)

			colors := baseColors[:tt.numColors]
			img := createTestImage(tt.numColors, 1, colors)

			t.Logf("Testing %s: count=%d, numColors=%d", tt.name, tt.count, tt.numColors)

			profile, err := testP.ProcessImage(img)
			if err != nil {
				t.Fatalf("ProcessImage failed: %v", err)
			}

			dominant := profile.Pool.DominantColors
			allColors := profile.Pool.AllColors

			t.Logf("  Result: %d dominant / %d total", len(dominant), len(allColors))

			if tt.expectAll {
				if len(dominant) != len(allColors) {
					t.Errorf("Expected all colors (%d), got %d dominant", len(allColors), len(dominant))
				}
			} else {
				expected := min(tt.count, len(allColors))
				if len(dominant) != expected {
					t.Errorf("Expected %d dominant colors, got %d", expected, len(dominant))
				}
			}
		})
	}
}

func TestColorPoolIntegration(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Create a complex test image
	colors := []color.RGBA{
		{255, 0, 0, 255},     // Vibrant red
		{200, 100, 100, 255}, // Muted red  
		{100, 100, 100, 255}, // Gray
		{240, 240, 240, 255}, // Light gray
		{20, 20, 20, 255},    // Dark gray
		{0, 200, 0, 255},     // Green
	}

	img := createTestImage(12, 1, []color.RGBA{
		colors[0], colors[0], colors[0], colors[0], // Red: 4 pixels
		colors[1], colors[1], colors[1],            // Muted red: 3 pixels
		colors[2], colors[2],                       // Gray: 2 pixels
		colors[3], colors[4], colors[5],            // Others: 1 each
	})

	t.Logf("Testing integrated ColorPool with complex image")

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	pool := profile.Pool

	// Comprehensive validation
	t.Logf("Integration test results:")
	t.Logf("  Total pixels: %d", pool.TotalPixels)
	t.Logf("  Unique colors: %d", pool.UniqueColors)
	t.Logf("  Lightness - Dark: %d, Mid: %d, Light: %d",
		len(pool.ByLightness.Dark), len(pool.ByLightness.Mid), len(pool.ByLightness.Light))
	t.Logf("  Saturation - Gray: %d, Muted: %d, Normal: %d, Vibrant: %d",
		len(pool.BySaturation.Gray), len(pool.BySaturation.Muted),
		len(pool.BySaturation.Normal), len(pool.BySaturation.Vibrant))
	t.Logf("  Hue families: %d", len(pool.ByHue))

	// Verify total pixels matches image
	if pool.TotalPixels != 12 {
		t.Errorf("Expected 12 total pixels, got %d", pool.TotalPixels)
	}

	// Verify all colors are accounted for in groupings
	groupedByLightness := len(pool.ByLightness.Dark) + len(pool.ByLightness.Mid) + len(pool.ByLightness.Light)
	if groupedByLightness != len(pool.AllColors) {
		t.Errorf("Lightness grouping mismatch: %d grouped != %d total", groupedByLightness, len(pool.AllColors))
	}

	groupedBySaturation := len(pool.BySaturation.Gray) + len(pool.BySaturation.Muted) +
						   len(pool.BySaturation.Normal) + len(pool.BySaturation.Vibrant)
	if groupedBySaturation != len(pool.AllColors) {
		t.Errorf("Saturation grouping mismatch: %d grouped != %d total", groupedBySaturation, len(pool.AllColors))
	}

	// Verify statistics are calculated
	if pool.Statistics.ChromaticDiversity == 0 {
		t.Error("Statistics should be calculated (ChromaticDiversity > 0 expected)")
	}
}

// Helper function to find minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}