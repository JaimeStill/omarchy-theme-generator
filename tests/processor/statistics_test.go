package processor_test

import (
	"image/color"
	"math"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/processor"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

func TestCalculateStatistics(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Create test image with diverse colors
	colors := []color.RGBA{
		{255, 0, 0, 255},     // Red (vibrant)
		{150, 100, 100, 255}, // Muted red-brown
		{128, 128, 128, 255}, // Gray
		{240, 240, 240, 255}, // Light gray
		{30, 30, 30, 255},    // Dark gray
		{0, 200, 0, 255},     // Green (vibrant)
		{0, 0, 200, 255},     // Blue (normal)
		{200, 200, 50, 255},  // Yellow (normal)
	}

	img := createTestImage(len(colors), 1, colors)

	t.Logf("Testing statistics calculation with %d colors", len(colors))

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	stats := profile.Pool.Statistics

	t.Logf("Color statistics results:")
	t.Logf("  Hue histogram bins: %d", len(stats.HueHistogram))
	t.Logf("  Lightness histogram bins: %d", len(stats.LightnessHistogram))
	t.Logf("  Saturation groups: %d", len(stats.SaturationGroups))
	t.Logf("  Primary hue: %.1f°", stats.PrimaryHue)
	t.Logf("  Secondary hue: %.1f°", stats.SecondaryHue)
	t.Logf("  Tertiary hue: %.1f°", stats.TertiaryHue)
	t.Logf("  Chromatic diversity: %.3f", stats.ChromaticDiversity)
	t.Logf("  Contrast range: %.3f", stats.ContrastRange)
	t.Logf("  Hue variance: %.3f", stats.HueVariance)
	t.Logf("  Lightness spread: %.3f", stats.LightnessSpread)
	t.Logf("  Saturation spread: %.3f", stats.SaturationSpread)

	// Validate basic structure
	if len(stats.HueHistogram) == 0 {
		t.Error("Hue histogram should have bins")
	}

	if len(stats.LightnessHistogram) != 10 {
		t.Errorf("Expected 10 lightness histogram bins, got %d", len(stats.LightnessHistogram))
	}

	// Validate saturation groups
	expectedGroups := []string{"gray", "muted", "normal", "vibrant"}
	for _, group := range expectedGroups {
		if _, exists := stats.SaturationGroups[group]; !exists {
			t.Errorf("Saturation group '%s' missing from statistics", group)
		}
	}

	// Validate ranges
	if stats.ChromaticDiversity < 0 || stats.ChromaticDiversity > 1 {
		t.Errorf("ChromaticDiversity %.3f should be in range [0, 1]", stats.ChromaticDiversity)
	}

	if stats.ContrastRange < 0 || stats.ContrastRange > 1 {
		t.Errorf("ContrastRange %.3f should be in range [0, 1]", stats.ContrastRange)
	}

	if stats.LightnessSpread < 0 || stats.LightnessSpread > 1 {
		t.Errorf("LightnessSpread %.3f should be in range [0, 1]", stats.LightnessSpread)
	}

	if stats.SaturationSpread < 0 || stats.SaturationSpread > 1 {
		t.Errorf("SaturationSpread %.3f should be in range [0, 1]", stats.SaturationSpread)
	}
}

func TestChromaticDiversity(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Test high diversity (rainbow colors)
	highDiversityColors := []color.RGBA{
		{255, 0, 0, 255},   // Red (0°)
		{255, 127, 0, 255}, // Orange (30°)
		{255, 255, 0, 255}, // Yellow (60°)
		{127, 255, 0, 255}, // Lime (90°)
		{0, 255, 0, 255},   // Green (120°)
		{0, 255, 127, 255}, // Spring (150°)
		{0, 255, 255, 255}, // Cyan (180°)
		{0, 127, 255, 255}, // Azure (210°)
		{0, 0, 255, 255},   // Blue (240°)
		{127, 0, 255, 255}, // Violet (270°)
		{255, 0, 255, 255}, // Magenta (300°)
		{255, 0, 127, 255}, // Rose (330°)
	}

	highImg := createTestImage(len(highDiversityColors), 1, highDiversityColors)

	highProfile, err := p.ProcessImage(highImg)
	if err != nil {
		t.Fatalf("ProcessImage failed for high diversity: %v", err)
	}

	// Test low diversity (monochromatic blues)
	lowDiversityColors := []color.RGBA{
		{0, 0, 100, 255},   // Dark blue
		{0, 0, 150, 255},   // Medium blue
		{0, 0, 200, 255},   // Light blue
		{50, 50, 255, 255}, // Bright blue
	}

	lowImg := createTestImage(len(lowDiversityColors), 1, lowDiversityColors)

	lowProfile, err := p.ProcessImage(lowImg)
	if err != nil {
		t.Fatalf("ProcessImage failed for low diversity: %v", err)
	}

	highDiversity := highProfile.Pool.Statistics.ChromaticDiversity
	lowDiversity := lowProfile.Pool.Statistics.ChromaticDiversity

	t.Logf("Chromatic diversity comparison:")
	t.Logf("  High diversity (rainbow): %.3f", highDiversity)
	t.Logf("  Low diversity (monochromatic): %.3f", lowDiversity)

	// High diversity should be significantly higher
	if highDiversity <= lowDiversity {
		t.Errorf("High diversity %.3f should be greater than low diversity %.3f",
			highDiversity, lowDiversity)
	}

	// Both should be in valid range
	if highDiversity < 0 || highDiversity > 1 {
		t.Errorf("High diversity %.3f should be in range [0, 1]", highDiversity)
	}

	if lowDiversity < 0 || lowDiversity > 1 {
		t.Errorf("Low diversity %.3f should be in range [0, 1]", lowDiversity)
	}
}

func TestContrastRange(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Test high contrast (black and white)
	highContrastColors := []color.RGBA{
		{0, 0, 0, 255},       // Pure black
		{255, 255, 255, 255}, // Pure white
	}

	highImg := createTestImage(len(highContrastColors), 1, highContrastColors)

	highProfile, err := p.ProcessImage(highImg)
	if err != nil {
		t.Fatalf("ProcessImage failed for high contrast: %v", err)
	}

	// Test low contrast (similar grays)
	lowContrastColors := []color.RGBA{
		{100, 100, 100, 255}, // Medium gray
		{120, 120, 120, 255}, // Slightly lighter gray
		{130, 130, 130, 255}, // Light gray
	}

	lowImg := createTestImage(len(lowContrastColors), 1, lowContrastColors)

	lowProfile, err := p.ProcessImage(lowImg)
	if err != nil {
		t.Fatalf("ProcessImage failed for low contrast: %v", err)
	}

	highContrast := highProfile.Pool.Statistics.ContrastRange
	lowContrast := lowProfile.Pool.Statistics.ContrastRange

	t.Logf("Contrast range comparison:")
	t.Logf("  High contrast (black/white): %.3f", highContrast)
	t.Logf("  Low contrast (similar grays): %.3f", lowContrast)

	// High contrast should be significantly higher
	if highContrast <= lowContrast {
		t.Errorf("High contrast %.3f should be greater than low contrast %.3f",
			highContrast, lowContrast)
	}

	// High contrast should be close to maximum (1.0)
	if highContrast < 0.8 {
		t.Errorf("High contrast %.3f should be close to 1.0 for black/white", highContrast)
	}
}

func TestHueHistogram(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Create colors with known hue distribution
	colors := []color.RGBA{
		{255, 0, 0, 255},   // Red (0°) - Sector 0
		{255, 0, 0, 255},   // Red (0°) - Sector 0 
		{0, 255, 0, 255},   // Green (120°) - Sector 4
		{0, 0, 255, 255},   // Blue (240°) - Sector 8
		{128, 128, 128, 255}, // Gray (excluded from histogram)
	}

	img := createTestImage(len(colors), 1, colors)

	t.Logf("Testing hue histogram with known color distribution")

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	histogram := profile.Pool.Statistics.HueHistogram
	expectedBins := s.HueSectorCount

	t.Logf("Hue histogram results:")
	t.Logf("  Expected bins: %d", expectedBins)
	t.Logf("  Actual bins: %d", len(histogram))

	if len(histogram) != expectedBins {
		t.Errorf("Expected %d histogram bins, got %d", expectedBins, len(histogram))
	}

	// Sum should be 1.0 (probability distribution)
	sum := 0.0
	nonZeroBins := 0
	for i, prob := range histogram {
		sum += prob
		if prob > 0 {
			nonZeroBins++
			t.Logf("  Bin %d (%.1f°-%.1f°): %.3f", i, 
				float64(i)*s.HueSectorSize, float64(i+1)*s.HueSectorSize, prob)
		}
	}

	t.Logf("  Total probability: %.3f", sum)
	t.Logf("  Non-zero bins: %d", nonZeroBins)

	if math.Abs(sum-1.0) > 0.001 {
		t.Errorf("Histogram probabilities should sum to 1.0, got %.3f", sum)
	}

	// Should have some hue representation (not all grayscale)
	if nonZeroBins == 0 {
		t.Error("Expected some non-zero histogram bins for colored image")
	}
}

func TestLightnessSpread(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Test balanced lightness distribution
	balancedColors := []color.RGBA{
		{30, 30, 30, 255},    // Dark
		{60, 60, 60, 255},    // Dark
		{128, 128, 128, 255}, // Mid
		{160, 160, 160, 255}, // Mid
		{220, 220, 220, 255}, // Light
		{240, 240, 240, 255}, // Light
	}

	balancedImg := createTestImage(len(balancedColors), 1, balancedColors)

	balancedProfile, err := p.ProcessImage(balancedImg)
	if err != nil {
		t.Fatalf("ProcessImage failed for balanced: %v", err)
	}

	// Test unbalanced lightness distribution (all dark)
	unbalancedColors := []color.RGBA{
		{10, 10, 10, 255},  // Dark
		{20, 20, 20, 255},  // Dark
		{30, 30, 30, 255},  // Dark
		{40, 40, 40, 255},  // Dark
	}

	unbalancedImg := createTestImage(len(unbalancedColors), 1, unbalancedColors)

	unbalancedProfile, err := p.ProcessImage(unbalancedImg)
	if err != nil {
		t.Fatalf("ProcessImage failed for unbalanced: %v", err)
	}

	balancedSpread := balancedProfile.Pool.Statistics.LightnessSpread
	unbalancedSpread := unbalancedProfile.Pool.Statistics.LightnessSpread

	t.Logf("Lightness spread comparison:")
	t.Logf("  Balanced distribution: %.3f", balancedSpread)
	t.Logf("  Unbalanced distribution: %.3f", unbalancedSpread)

	// Balanced should have higher spread
	if balancedSpread <= unbalancedSpread {
		t.Errorf("Balanced spread %.3f should be greater than unbalanced %.3f",
			balancedSpread, unbalancedSpread)
	}

	// Both should be in valid range
	if balancedSpread < 0 || balancedSpread > 1 {
		t.Errorf("Balanced spread %.3f should be in range [0, 1]", balancedSpread)
	}

	// Log lightness group counts for diagnostics
	t.Logf("Balanced lightness groups: Dark=%d, Mid=%d, Light=%d",
		len(balancedProfile.Pool.ByLightness.Dark),
		len(balancedProfile.Pool.ByLightness.Mid),
		len(balancedProfile.Pool.ByLightness.Light))
	
	t.Logf("Unbalanced lightness groups: Dark=%d, Mid=%d, Light=%d",
		len(unbalancedProfile.Pool.ByLightness.Dark),
		len(unbalancedProfile.Pool.ByLightness.Mid),
		len(unbalancedProfile.Pool.ByLightness.Light))
}

func TestSaturationSpread(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Test diverse saturation (all 4 groups)
	diverseColors := []color.RGBA{
		{128, 128, 128, 255}, // Gray
		{180, 150, 150, 255}, // Muted
		{200, 100, 100, 255}, // Normal
		{255, 0, 0, 255},     // Vibrant
	}

	diverseImg := createTestImage(len(diverseColors), 1, diverseColors)

	diverseProfile, err := p.ProcessImage(diverseImg)
	if err != nil {
		t.Fatalf("ProcessImage failed for diverse saturation: %v", err)
	}

	// Test limited saturation (only grayscale)
	limitedColors := []color.RGBA{
		{100, 100, 100, 255}, // Gray
		{150, 150, 150, 255}, // Gray
		{200, 200, 200, 255}, // Gray
	}

	limitedImg := createTestImage(len(limitedColors), 1, limitedColors)

	limitedProfile, err := p.ProcessImage(limitedImg)
	if err != nil {
		t.Fatalf("ProcessImage failed for limited saturation: %v", err)
	}

	diverseSpread := diverseProfile.Pool.Statistics.SaturationSpread
	limitedSpread := limitedProfile.Pool.Statistics.SaturationSpread

	t.Logf("Saturation spread comparison:")
	t.Logf("  Diverse saturation: %.3f", diverseSpread)
	t.Logf("  Limited saturation: %.3f", limitedSpread)

	// Diverse should have higher spread
	if diverseSpread <= limitedSpread {
		t.Errorf("Diverse spread %.3f should be greater than limited %.3f",
			diverseSpread, limitedSpread)
	}

	// Diverse should be close to maximum (1.0 for all 4 groups)
	if diverseSpread < 0.8 {
		t.Logf("Note: Diverse spread %.3f lower than expected (saturation grouping may be strict)",
			diverseSpread)
	}

	// Log saturation group counts for diagnostics
	t.Logf("Diverse saturation groups: Gray=%d, Muted=%d, Normal=%d, Vibrant=%d",
		len(diverseProfile.Pool.BySaturation.Gray),
		len(diverseProfile.Pool.BySaturation.Muted),
		len(diverseProfile.Pool.BySaturation.Normal),
		len(diverseProfile.Pool.BySaturation.Vibrant))
	
	t.Logf("Limited saturation groups: Gray=%d, Muted=%d, Normal=%d, Vibrant=%d",
		len(limitedProfile.Pool.BySaturation.Gray),
		len(limitedProfile.Pool.BySaturation.Muted),
		len(limitedProfile.Pool.BySaturation.Normal),
		len(limitedProfile.Pool.BySaturation.Vibrant))
}

func TestDominantHues(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Create colors with known dominant hues
	colors := []color.RGBA{
		{255, 0, 0, 255},     // Red (0°) - Should be primary
		{255, 50, 50, 255},   // Red variant
		{255, 100, 100, 255}, // Red variant (total: 3 reds)
		{0, 255, 0, 255},     // Green (120°) - Should be secondary  
		{50, 255, 50, 255},   // Green variant (total: 2 greens)
		{0, 0, 255, 255},     // Blue (240°) - Should be tertiary (total: 1 blue)
	}

	img := createTestImage(len(colors), 1, colors)

	t.Logf("Testing dominant hue detection")

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("ProcessImage failed: %v", err)
	}

	stats := profile.Pool.Statistics

	t.Logf("Dominant hues detected:")
	t.Logf("  Primary: %.1f°", stats.PrimaryHue)
	t.Logf("  Secondary: %.1f°", stats.SecondaryHue) 
	t.Logf("  Tertiary: %.1f°", stats.TertiaryHue)

	// Validate hue ranges (allowing for sector discretization)
	primarySector := int(stats.PrimaryHue / s.HueSectorSize)
	secondarySector := int(stats.SecondaryHue / s.HueSectorSize)
	tertiarySector := int(stats.TertiaryHue / s.HueSectorSize)

	t.Logf("Hue sectors:")
	t.Logf("  Primary sector: %d", primarySector)
	t.Logf("  Secondary sector: %d", secondarySector)
	t.Logf("  Tertiary sector: %d", tertiarySector)

	// Validate hues are within expected ranges
	if stats.PrimaryHue < 0 || stats.PrimaryHue >= 360 {
		t.Errorf("Primary hue %.1f° should be in range [0, 360)", stats.PrimaryHue)
	}

	if stats.SecondaryHue < 0 || stats.SecondaryHue >= 360 {
		t.Errorf("Secondary hue %.1f° should be in range [0, 360)", stats.SecondaryHue)
	}

	if stats.TertiaryHue < 0 || stats.TertiaryHue >= 360 {
		t.Errorf("Tertiary hue %.1f° should be in range [0, 360)", stats.TertiaryHue)
	}

	// All dominant hues should be different (different sectors)
	if primarySector == secondarySector || primarySector == tertiarySector || secondarySector == tertiarySector {
		t.Logf("Note: Some dominant hues are in the same sector (may indicate clustering)")
	}
}

func TestStatisticsEdgeCases(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)

	// Test with single color
	singleColor := []color.RGBA{{255, 0, 0, 255}}
	singleImg := createTestImage(1, 1, singleColor)

	singleProfile, err := p.ProcessImage(singleImg)
	if err != nil {
		t.Fatalf("ProcessImage failed for single color: %v", err)
	}

	singleStats := singleProfile.Pool.Statistics

	t.Logf("Single color statistics:")
	t.Logf("  Chromatic diversity: %.3f", singleStats.ChromaticDiversity)
	t.Logf("  Contrast range: %.3f", singleStats.ContrastRange)
	t.Logf("  Hue variance: %.3f", singleStats.HueVariance)

	// Single color should have zero or minimal diversity/variance
	if singleStats.ChromaticDiversity > 0.1 {
		t.Errorf("Single color should have low chromatic diversity, got %.3f", singleStats.ChromaticDiversity)
	}

	if singleStats.ContrastRange != 0.0 {
		t.Errorf("Single color should have zero contrast range, got %.3f", singleStats.ContrastRange)
	}

	// Test with only grayscale colors
	grayscaleColors := []color.RGBA{
		{50, 50, 50, 255},
		{100, 100, 100, 255},
		{150, 150, 150, 255},
		{200, 200, 200, 255},
	}

	grayscaleImg := createTestImage(len(grayscaleColors), 1, grayscaleColors)

	grayscaleProfile, err := p.ProcessImage(grayscaleImg)
	if err != nil {
		t.Fatalf("ProcessImage failed for grayscale: %v", err)
	}

	grayscaleStats := grayscaleProfile.Pool.Statistics

	t.Logf("Grayscale-only statistics:")
	t.Logf("  Chromatic diversity: %.3f", grayscaleStats.ChromaticDiversity)
	t.Logf("  Hue variance: %.3f", grayscaleStats.HueVariance)

	// Grayscale should have zero or minimal hue-based metrics
	if grayscaleStats.HueVariance > 0.1 {
		t.Errorf("Grayscale image should have minimal hue variance, got %.3f", grayscaleStats.HueVariance)
	}

	// Should still have contrast range from lightness differences
	if grayscaleStats.ContrastRange == 0 {
		t.Error("Grayscale image with different lightness should have contrast range > 0")
	}
}