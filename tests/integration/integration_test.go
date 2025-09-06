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
			t.Logf("  Color Scheme: %s", profile.ColorScheme)
			t.Logf("  Dominant Hue: %.1f°", profile.DominantHue)
			t.Logf("  Hue Variance: %.1f°", profile.HueVariance)
			t.Logf("  Average Luminance: %.3f", profile.AvgLuminance)
			t.Logf("  Average Saturation: %.3f", profile.AvgSaturation)
			t.Logf("  Is Grayscale: %t", profile.IsGrayscale)
			t.Logf("  Is Monochromatic: %t", profile.IsMonochromatic)
			t.Logf("  Categories found: %d/27 (%.1f%%)",
				len(profile.Colors.Categories), 
				float64(len(profile.Colors.Categories))/27*100)

			// Validate performance targets
			maxProcessingTime := 2 * time.Second
			if processingTime > maxProcessingTime {
				t.Errorf("Processing time %v exceeds target %v", processingTime, maxProcessingTime)
			}

			// Validate profile completeness
			if profile.Colors.UniqueColors == 0 {
				t.Error("No unique colors found")
			}

			if profile.Colors.TotalPixels == 0 {
				t.Error("Total pixels not set")
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

			// Validate category extraction
			if len(profile.Colors.Categories) == 0 {
				t.Error("No categories extracted from image")
			}

			// Ensure core categories are present when applicable
			if _, hasBackground := profile.Colors.Categories[processor.CategoryBackground]; !hasBackground {
				t.Error("Background category not found")
			}

			// Validate category candidates
			totalCandidates := 0
			for category, candidates := range profile.Colors.CategoryCandidates {
				totalCandidates += len(candidates)
				t.Logf("Category %s has %d candidates", category, len(candidates))
			}

			if totalCandidates == 0 {
				t.Logf("Note: No category candidates found - image may have very limited color diversity")
			}

			t.Logf("Integration test passed: %d categories, %d total candidates",
				len(profile.Colors.Categories), totalCandidates)
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

	// Test color scheme identification
	t.Logf("Testing color theory integration:")
	t.Logf("  Color scheme: %s", profile.ColorScheme)
	t.Logf("  Dominant hue: %.1f°", profile.DominantHue)
	t.Logf("  Hue variance: %.1f°", profile.HueVariance)

	// Test color theory integration with extracted colors
	if len(profile.Colors.Categories) > 1 {
		t.Logf("  Color theory validation: %d categories extracted", len(profile.Colors.Categories))
		
		// Verify scheme consistency
		if profile.ColorScheme == "" {
			t.Logf("Note: No color scheme identified from extracted colors")
		}
	}

	// Test accessibility calculations
	if bg, hasBg := profile.Colors.Categories[processor.CategoryBackground]; hasBg {
		if fg, hasFg := profile.Colors.Categories[processor.CategoryForeground]; hasFg {
			contrast := chromatic.ContrastRatio(bg, fg)
			accessible := chromatic.IsAccessible(bg, fg, chromatic.AA)
			
			t.Logf("  Background-Foreground contrast: %.2f:1", contrast)
			t.Logf("  WCAG AA compliant: %t", accessible)
			
			if contrast < 1.0 || contrast > 21.0 {
				t.Errorf("Invalid contrast ratio: %.2f", contrast)
			}
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
			t.Logf("  Categories extracted: %d", len(profile.Colors.Categories))
			t.Logf("  Unique colors: %d", profile.Colors.UniqueColors)

			// Validate performance target
			if processingTime > tc.maxTime {
				t.Errorf("Processing time %v exceeds target %v for %s",
					processingTime, tc.maxTime, tc.name)
			}

			// Validate extraction completeness
			if len(profile.Colors.Categories) == 0 {
				t.Error("No categories extracted despite successful processing")
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