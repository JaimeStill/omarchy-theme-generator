package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/loader"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/processor"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

// TestEnd2EndImageProcessing tests the complete pipeline from image loading to analysis
func TestEnd2EndImageProcessing(t *testing.T) {
	// Initialize all components
	s := settings.DefaultSettings()
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	ctx := context.Background()

	// Test with multiple different image types
	testImages := []struct {
		name        string
		filename    string
		expectDark  bool
		expectColor bool
	}{
		{
			name:        "Grayscale image",
			filename:    "grayscale.jpeg",
			expectDark:  true,
			expectColor: false,
		},
		{
			name:        "Colorful abstract",
			filename:    "abstract.jpeg",
			expectDark:  true,
			expectColor: true,
		},
		{
			name:        "Natural landscape",
			filename:    "mountains.jpeg",
			expectDark:  false,
			expectColor: true,
		},
		{
			name:        "Urban night scene",
			filename:    "night-city.jpeg",
			expectDark:  true,
			expectColor: true,
		},
	}

	for _, tc := range testImages {
		t.Run(tc.name, func(t *testing.T) {
			imagePath := filepath.Join("..", "images", tc.filename)
			
			// Verify file exists
			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				t.Skipf("Test image %s not found, skipping", tc.filename)
				return
			}

			// Measure processing time
			start := time.Now()

			// Load image
			img, err := l.LoadImage(ctx, imagePath)
			if err != nil {
				t.Fatalf("Failed to load image: %v", err)
			}

			// Get image info
			info, err := l.GetImageInfo(ctx, imagePath)
			if err != nil {
				t.Fatalf("Failed to get image info: %v", err)
			}

			// Process image through complete pipeline
			profile, err := p.ProcessImage(img)
			if err != nil {
				t.Fatalf("Failed to process image: %v", err)
			}

			processingTime := time.Since(start)

			// Comprehensive diagnostic logging
			t.Logf("Image: %s (%dx%d, %.1f MP)",
				tc.filename, info.Width, info.Height, float64(info.PixelCount())/1000000)
			t.Logf("Processing time: %v", processingTime)
			t.Logf("Profile characteristics:")
			t.Logf("  Mode: %s", profile.Mode)
			t.Logf("  Dominant Hue: %.1f°", profile.DominantHue)
			t.Logf("  Hue Variance: %.1f°", profile.HueVariance)
			t.Logf("  Average Luminance: %.3f", profile.AvgLuminance)
			t.Logf("  Average Saturation: %.3f", profile.AvgSaturation)
			t.Logf("  Is Grayscale: %t", profile.IsGrayscale)
			t.Logf("  Is Monochromatic: %t", profile.IsMonochromatic)
			t.Logf("  Colors extracted: %d unique from %d pixels",
				profile.Pool.UniqueColors, profile.Pool.TotalPixels)
			t.Logf("  Dominant colors: %d",
				len(profile.Pool.DominantColors))

			// Validate performance targets
			maxProcessingTime := 2 * time.Second
			if processingTime > maxProcessingTime {
				t.Errorf("Processing time %v exceeds target %v", processingTime, maxProcessingTime)
			}

			// Validate profile completeness
			if profile.Pool.UniqueColors == 0 {
				t.Error("No unique colors found")
			}

			if profile.Pool.TotalPixels == 0 {
				t.Error("Total pixels not set")
			}

			if len(profile.Pool.AllColors) == 0 {
				t.Error("No colors extracted from image")
			}

			// Validate mode detection
			isDarkMode := profile.Mode == "Dark"
			if isDarkMode != tc.expectDark {
				t.Logf("Expected dark mode: %t, got: %t (avg luminance: %.3f)",
					tc.expectDark, isDarkMode, profile.AvgLuminance)
			}

			// Validate color detection
			hasColor := !profile.IsGrayscale
			if hasColor != tc.expectColor {
				t.Logf("Expected color: %t, got: %t (is grayscale: %t)",
					tc.expectColor, hasColor, profile.IsGrayscale)
			}

			// Validate characteristic grouping
			totalGrouped := len(profile.Pool.ByLightness.Dark) + len(profile.Pool.ByLightness.Mid) + len(profile.Pool.ByLightness.Light)
			if totalGrouped == 0 {
				t.Error("No colors grouped by lightness")
			}

			// Validate statistics calculation
			if profile.Pool.Statistics.ChromaticDiversity < 0 || profile.Pool.Statistics.ChromaticDiversity > 1 {
				t.Errorf("Invalid chromatic diversity: %.3f (should be 0-1)", profile.Pool.Statistics.ChromaticDiversity)
			}

			// Log characteristic grouping results
			t.Logf("Color grouping results:")
			t.Logf("  Dark: %d, Mid: %d, Light: %d",
				len(profile.Pool.ByLightness.Dark), len(profile.Pool.ByLightness.Mid), len(profile.Pool.ByLightness.Light))
			t.Logf("  Gray: %d, Muted: %d, Normal: %d, Vibrant: %d",
				len(profile.Pool.BySaturation.Gray), len(profile.Pool.BySaturation.Muted),
				len(profile.Pool.BySaturation.Normal), len(profile.Pool.BySaturation.Vibrant))
			t.Logf("  Hue families: %d", len(profile.Pool.ByHue))

			t.Logf("Integration test passed: %d colors extracted, grouped into characteristics",
				profile.Pool.UniqueColors)
		})
	}
}

// TestSettingsIntegration validates settings loading and application
func TestSettingsIntegration(t *testing.T) {
	// Test default settings
	defaultSettings := settings.DefaultSettings()
	
	t.Logf("Testing default settings integration")
	t.Logf("Grayscale threshold: %.3f", defaultSettings.GrayscaleThreshold)
	t.Logf("Monochromatic tolerance: %.1f°", defaultSettings.MonochromaticTolerance)
	t.Logf("Theme mode threshold: %.3f", defaultSettings.ThemeModeThreshold)
	
	// Verify processor can be created with settings
	p := processor.New(defaultSettings)
	if p == nil {
		t.Fatal("Failed to create processor with default settings")
	}

	// Test settings affect processing
	ctx := context.Background()
	l := loader.NewFileLoader(defaultSettings)
	
	imagePath := filepath.Join("..", "images", "grayscale.jpeg")
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		t.Skip("Grayscale test image not found")
		return
	}

	img, err := l.LoadImage(ctx, imagePath)
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("Failed to process image: %v", err)
	}

	// Verify settings are applied correctly
	t.Logf("Profile detected grayscale: %t (threshold: %.3f, avg sat: %.3f)",
		profile.IsGrayscale, defaultSettings.GrayscaleThreshold, profile.AvgSaturation)
	
	if profile.AvgSaturation < defaultSettings.GrayscaleThreshold && !profile.IsGrayscale {
		t.Errorf("Expected grayscale detection when saturation %.3f < threshold %.3f",
			profile.AvgSaturation, defaultSettings.GrayscaleThreshold)
	}
}

// TestColorTheoryIntegration validates color theory algorithms work in integration
func TestColorTheoryIntegration(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	ctx := context.Background()

	// Test with a known colorful image
	imagePath := filepath.Join("..", "images", "abstract.jpeg")
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		t.Skip("Abstract test image not found")
		return
	}

	img, err := l.LoadImage(ctx, imagePath)
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}

	profile, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("Failed to process image: %v", err)
	}

	// Test color theory integration
	t.Logf("Testing color theory integration:")
	t.Logf("  Dominant hue: %.1f°", profile.DominantHue)
	t.Logf("  Hue variance: %.1f°", profile.HueVariance)

	// Test color theory integration with extracted colors
	if len(profile.Pool.AllColors) > 1 {
		t.Logf("  Color theory validation: %d colors extracted", len(profile.Pool.AllColors))
		
		// Color scheme identification is handled by pkg/palette, not processor
	}

	// Test accessibility calculations with dominant colors
	if len(profile.Pool.DominantColors) >= 2 {
		bg := profile.Pool.DominantColors[0].RGBA
		fg := profile.Pool.DominantColors[1].RGBA
		contrast := chromatic.ContrastRatio(bg, fg)
		accessible := chromatic.IsAccessible(bg, fg, chromatic.AA)
		
		t.Logf("  Top two colors contrast: %.2f:1", contrast)
		t.Logf("  WCAG AA compliant: %t", accessible)
		
		if contrast < 1.0 || contrast > 21.0 {
			t.Errorf("Invalid contrast ratio: %.2f", contrast)
		}
	}
}

// TestPerformanceTargets validates all performance requirements in integration
func TestPerformanceTargets(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	ctx := context.Background()

	// Test performance with different image sizes
	performanceTests := []struct {
		name           string
		filename       string
		maxTime        time.Duration
		sizeCategoryMP float64
	}{
		{
			name:           "Small image performance",
			filename:       "grayscale.jpeg",
			maxTime:        500 * time.Millisecond,
			sizeCategoryMP: 2.1,
		},
		{
			name:           "Large image performance",
			filename:       "simple.png",
			maxTime:        2 * time.Second,
			sizeCategoryMP: 14.7,
		},
	}

	for _, tc := range performanceTests {
		t.Run(tc.name, func(t *testing.T) {
			imagePath := filepath.Join("..", "images", tc.filename)
			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				t.Skipf("Test image %s not found", tc.filename)
				return
			}

			// Measure complete processing time
			start := time.Now()
			
			img, err := l.LoadImage(ctx, imagePath)
			if err != nil {
				t.Fatalf("Failed to load image: %v", err)
			}

			profile, err := p.ProcessImage(img)
			if err != nil {
				t.Fatalf("Failed to process image: %v", err)
			}

			processingTime := time.Since(start)

			t.Logf("Performance test: %s (%.1f MP)", tc.filename, tc.sizeCategoryMP)
			t.Logf("  Processing time: %v (target: <%v)", processingTime, tc.maxTime)
			t.Logf("  Colors extracted: %d", profile.Pool.UniqueColors)
			t.Logf("  Dominant colors: %d", len(profile.Pool.DominantColors))

			// Validate performance target
			if processingTime > tc.maxTime {
				t.Errorf("Processing time %v exceeds target %v for %s",
					processingTime, tc.maxTime, tc.name)
			}

			// Validate extraction completeness
			if len(profile.Pool.AllColors) == 0 {
				t.Error("No colors extracted despite successful processing")
			}

			// Overall performance should always be under 2s (project requirement)
			maxOverall := 2 * time.Second
			if processingTime > maxOverall {
				t.Errorf("Processing time %v exceeds project requirement %v",
					processingTime, maxOverall)
			}
		})
	}
}