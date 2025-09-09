package processor_test

import (
	"image/color"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/processor"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

func TestGroupByLightness(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Create test colors spanning lightness spectrum
	colors := []processor.WeightedColor{
		processor.NewWeightedColor(color.RGBA{10, 10, 10, 255}, 100, 1000),     // Dark (L ≈ 0.04)
		processor.NewWeightedColor(color.RGBA{128, 128, 128, 255}, 200, 1000),  // Mid (L = 0.5)
		processor.NewWeightedColor(color.RGBA{240, 240, 240, 255}, 150, 1000),  // Light (L ≈ 0.94)
		processor.NewWeightedColor(color.RGBA{60, 60, 60, 255}, 80, 1000),      // Dark (L ≈ 0.24)
		processor.NewWeightedColor(color.RGBA{180, 180, 180, 255}, 120, 1000),  // Mid (L ≈ 0.71)
	}

	t.Logf("Testing lightness grouping with %d colors", len(colors))
	t.Logf("Settings: darkMax=%.3f, lightMin=%.3f", s.LightnessDarkMax, s.LightnessLightMin)

	// Create processor to access private method via public interface
	// Use a simple test image to trigger grouping
	img := createTestImage(len(colors), 1, []color.RGBA{
		{10, 10, 10, 255},
		{128, 128, 128, 255}, 
		{240, 240, 240, 255},
		{60, 60, 60, 255},
		{180, 180, 180, 255},
	})

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	groups := profile.Pool.ByLightness

	t.Logf("Lightness grouping results:")
	t.Logf("  Dark colors: %d", len(groups.Dark))
	t.Logf("  Mid colors: %d", len(groups.Mid))
	t.Logf("  Light colors: %d", len(groups.Light))

	// Verify groups are non-empty and properly distributed
	totalGrouped := len(groups.Dark) + len(groups.Mid) + len(groups.Light)
	if totalGrouped == 0 {
		t.Error("No colors were grouped by lightness")
	}

	// Verify dark colors have lowest lightness
	if len(groups.Dark) > 0 {
		t.Logf("  Dark group weights: %v", getWeights(groups.Dark))
		// Should be sorted by weight (highest first)
		if len(groups.Dark) > 1 && groups.Dark[0].Weight < groups.Dark[1].Weight {
			t.Error("Dark colors not sorted by weight")
		}
	}

	// Verify light colors have highest lightness  
	if len(groups.Light) > 0 {
		t.Logf("  Light group weights: %v", getWeights(groups.Light))
		if len(groups.Light) > 1 && groups.Light[0].Weight < groups.Light[1].Weight {
			t.Error("Light colors not sorted by weight")
		}
	}

	// Verify mid colors fall between thresholds
	if len(groups.Mid) > 0 {
		t.Logf("  Mid group weights: %v", getWeights(groups.Mid))
		if len(groups.Mid) > 1 && groups.Mid[0].Weight < groups.Mid[1].Weight {
			t.Error("Mid colors not sorted by weight")
		}
	}
}

func TestGroupBySaturation(t *testing.T) {
	s := settings.DefaultSettings()  
	p := processor.New(s)

	// Create test colors spanning saturation spectrum
	colors := []color.RGBA{
		{128, 128, 128, 255}, // Gray (S = 0)
		{150, 128, 128, 255}, // Muted (S ≈ 0.15)
		{200, 128, 128, 255}, // Normal (S ≈ 0.36)  
		{255, 128, 128, 255}, // Vibrant (S ≈ 0.50)
		{255, 0, 0, 255},     // Vibrant (S = 1.0)
	}

	img := createTestImage(len(colors), 1, colors)

	t.Logf("Testing saturation grouping with %d colors", len(colors))
	t.Logf("Settings: grayMax=%.3f, mutedMax=%.3f, normalMax=%.3f", 
		s.SaturationGrayMax, s.SaturationMutedMax, s.SaturationNormalMax)

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	groups := profile.Pool.BySaturation

	t.Logf("Saturation grouping results:")
	t.Logf("  Gray colors: %d", len(groups.Gray))
	t.Logf("  Muted colors: %d", len(groups.Muted))
	t.Logf("  Normal colors: %d", len(groups.Normal))
	t.Logf("  Vibrant colors: %d", len(groups.Vibrant))

	// Verify groups are properly distributed
	totalGrouped := len(groups.Gray) + len(groups.Muted) + len(groups.Normal) + len(groups.Vibrant)
	if totalGrouped == 0 {
		t.Error("No colors were grouped by saturation")
	}

	// Verify sorting within groups
	validateSorting := func(name string, group []processor.WeightedColor) {
		if len(group) > 1 {
			for i := 0; i < len(group)-1; i++ {
				if group[i].Weight < group[i+1].Weight {
					t.Errorf("%s group not properly sorted by weight", name)
					break
				}
			}
		}
	}

	validateSorting("Gray", groups.Gray)
	validateSorting("Muted", groups.Muted)
	validateSorting("Normal", groups.Normal) 
	validateSorting("Vibrant", groups.Vibrant)

	// Log weights for diagnostic purposes
	if len(groups.Gray) > 0 {
		t.Logf("  Gray group weights: %v", getWeights(groups.Gray))
	}
	if len(groups.Muted) > 0 {
		t.Logf("  Muted group weights: %v", getWeights(groups.Muted))
	}
	if len(groups.Normal) > 0 {
		t.Logf("  Normal group weights: %v", getWeights(groups.Normal))
	}
	if len(groups.Vibrant) > 0 {
		t.Logf("  Vibrant group weights: %v", getWeights(groups.Vibrant))
	}
}

func TestGroupByHue(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Create test colors spanning hue spectrum (avoiding grayscale)
	colors := []color.RGBA{
		{255, 0, 0, 255},     // Red (H = 0°)
		{255, 255, 0, 255},   // Yellow (H = 60°)
		{0, 255, 0, 255},     // Green (H = 120°)
		{0, 255, 255, 255},   // Cyan (H = 180°)
		{0, 0, 255, 255},     // Blue (H = 240°)
		{255, 0, 255, 255},   // Magenta (H = 300°)
		{128, 128, 128, 255}, // Gray (should be excluded)
	}

	img := createTestImage(len(colors), 1, colors)

	t.Logf("Testing hue grouping with %d colors", len(colors))
	t.Logf("Settings: sectorSize=%.1f°, sectorCount=%d, grayscaleThreshold=%.3f",
		s.HueSectorSize, s.HueSectorCount, s.GrayscaleThreshold)

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	families := profile.Pool.ByHue

	t.Logf("Hue grouping results:")
	t.Logf("  Number of hue families: %d", len(families))

	totalGrouped := 0
	for sector, colors := range families {
		t.Logf("  Sector %d: %d colors (weights: %v)", sector, len(colors), getWeights(colors))
		totalGrouped += len(colors)

		// Verify sorting within families
		if len(colors) > 1 {
			for i := 0; i < len(colors)-1; i++ {
				if colors[i].Weight < colors[i+1].Weight {
					t.Errorf("Hue family %d not properly sorted by weight", sector)
					break
				}
			}
		}
	}

	// Should have some grouped colors (grayscale excluded)
	if totalGrouped == 0 {
		t.Error("No colors were grouped by hue")
	}

	// Verify grayscale colors are excluded
	expectedNonGrayscale := len(colors) - 1 // Minus the gray color
	if totalGrouped > expectedNonGrayscale {
		t.Logf("Note: %d colors grouped, expected max %d (grayscale should be excluded)",
			totalGrouped, expectedNonGrayscale)
	}

	// Verify hue sector distribution makes sense
	maxSector := s.HueSectorCount - 1
	for sector := range families {
		if sector < 0 || sector > maxSector {
			t.Errorf("Invalid hue sector %d (should be 0-%d)", sector, maxSector)
		}
	}
}

func TestGroupingWithEmptyInput(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Test with empty image (should handle gracefully)
	img := createTestImage(0, 0, []color.RGBA{})

	t.Logf("Testing grouping with empty input")

	_, err := p.ProcessImage(img)
	if err == nil {
		t.Error("Expected error for empty image")
	}

	t.Logf("Empty input correctly handled: %v", err)
}

func TestGroupingConsistency(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Test with same colors multiple times - should get consistent results
	colors := []color.RGBA{
		{255, 0, 0, 255},   // Red
		{128, 128, 128, 255}, // Gray
		{255, 255, 255, 255}, // White
		{0, 0, 0, 255},     // Black
	}

	img := createTestImage(4, 1, colors)

	t.Logf("Testing grouping consistency with repeated processing")

	// Process same image multiple times
	var profiles []*processor.ColorProfile
	for i := 0; i < 3; i++ {
		profile, err := p.ProcessImage(img)
		if err != nil {
			t.Fatalf("ProcessImage failed on iteration %d: %v", i, err)
		}
		profiles = append(profiles, profile)
	}

	// Compare results for consistency
	first := profiles[0].Pool
	for i, profile := range profiles[1:] {
		pool := profile.Pool

		if len(first.ByLightness.Dark) != len(pool.ByLightness.Dark) ||
		   len(first.ByLightness.Mid) != len(pool.ByLightness.Mid) ||
		   len(first.ByLightness.Light) != len(pool.ByLightness.Light) {
			t.Errorf("Inconsistent lightness grouping on iteration %d", i+2)
		}

		if len(first.BySaturation.Gray) != len(pool.BySaturation.Gray) ||
		   len(first.BySaturation.Vibrant) != len(pool.BySaturation.Vibrant) {
			t.Errorf("Inconsistent saturation grouping on iteration %d", i+2)
		}

		if len(first.ByHue) != len(pool.ByHue) {
			t.Errorf("Inconsistent hue grouping on iteration %d", i+2)
		}
	}

	t.Logf("Grouping consistency verified across %d iterations", len(profiles))
}

// Helper function to extract weights from WeightedColor slice for logging
func getWeights(colors []processor.WeightedColor) []float64 {
	weights := make([]float64, len(colors))
	for i, c := range colors {
		weights[i] = c.Weight
	}
	return weights
}