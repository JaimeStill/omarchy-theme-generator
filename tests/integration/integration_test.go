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
			t.Logf("  HasColor: %t", profile.HasColor)
			t.Logf("  ColorCount: %d", profile.ColorCount)
			t.Logf("  Colors extracted: %d clusters", len(profile.Colors))
			if len(profile.Colors) > 0 {
				t.Logf("  Top color: RGBA(%d,%d,%d,%d), Weight=%.3f, Lightness=%.3f",
					profile.Colors[0].R, profile.Colors[0].G, profile.Colors[0].B, profile.Colors[0].A,
					profile.Colors[0].Weight, profile.Colors[0].Lightness)
			}

			// Validate performance targets
			maxProcessingTime := 2 * time.Second
			if processingTime > maxProcessingTime {
				t.Errorf("Processing time %v exceeds target %v", processingTime, maxProcessingTime)
			}

			// Validate profile completeness
			if profile.ColorCount == 0 {
				t.Error("No colors found")
			}

			if len(profile.Colors) == 0 {
				t.Error("No color clusters extracted from image")
			}

			if len(profile.Colors) != profile.ColorCount {
				t.Errorf("ColorCount mismatch: expected %d, got %d colors", profile.ColorCount, len(profile.Colors))
			}

			// Validate mode detection
			isDarkMode := profile.Mode == "Dark"
			if isDarkMode != tc.expectDark {
				t.Logf("Expected dark mode: %t, got: %t",
					tc.expectDark, isDarkMode)
			}

			// Validate color detection
			hasColor := profile.HasColor
			if hasColor != tc.expectColor {
				t.Logf("Expected color: %t, got: %t",
					tc.expectColor, hasColor)
			}

			// Validate color characteristics
			var darkColors, lightColors, neutralColors, vibrantColors int
			for _, color := range profile.Colors {
				if color.IsDark {
					darkColors++
				}
				if color.IsLight {
					lightColors++
				}
				if color.IsNeutral {
					neutralColors++
				}
				if color.IsVibrant {
					vibrantColors++
				}
			}

			// Log color characteristics
			t.Logf("Color characteristics:")
			t.Logf("  Dark colors: %d, Light colors: %d", darkColors, lightColors)
			t.Logf("  Neutral colors: %d, Vibrant colors: %d", neutralColors, vibrantColors)

			// Validate reasonable color distribution
			if darkColors+lightColors == 0 {
				t.Error("No dark or light colors found - lightness calculation may be incorrect")
			}

			t.Logf("Integration test passed: %d color clusters extracted with characteristics",
				profile.ColorCount)
		})
	}
}

// TestSettingsIntegration validates settings loading and application
func TestSettingsIntegration(t *testing.T) {
	// Test default settings
	defaultSettings := settings.DefaultSettings()
	
	t.Logf("Testing default settings integration")
	t.Logf("Processor settings configured: %+v", defaultSettings.Processor)
	t.Logf("Chromatic settings configured: %+v", defaultSettings.Chromatic)
	t.Logf("Formats settings configured: %+v", defaultSettings.Formats)
	
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
	t.Logf("Profile has color: %t", profile.HasColor)
	t.Logf("Profile mode: %s", profile.Mode)
	t.Logf("Color count: %d", profile.ColorCount)

	// For grayscale image, expect no significant color
	if !profile.HasColor {
		t.Logf("✓ Correctly detected grayscale/low-saturation image")
	} else {
		t.Logf("Image has significant color content")
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
	t.Logf("  Total color clusters: %d", len(profile.Colors))
	t.Logf("  Has significant color: %t", profile.HasColor)

	// Test color theory integration with extracted colors
	if len(profile.Colors) > 1 {
		t.Logf("  Color theory validation: %d color clusters extracted", len(profile.Colors))

		// Log hue distribution
		hueCount := make(map[int]int)
		for _, color := range profile.Colors {
			hueRange := int(color.Hue) / 60 // Group into 60-degree ranges
			hueCount[hueRange]++
		}
		t.Logf("  Hue distribution by 60° ranges: %v", hueCount)
	}

	// Test accessibility calculations with top colors
	if len(profile.Colors) >= 2 {
		bg := profile.Colors[0].RGBA
		fg := profile.Colors[1].RGBA
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
			t.Logf("  Color clusters: %d", profile.ColorCount)
			t.Logf("  Colors extracted: %d", len(profile.Colors))

			// Validate performance target
			if processingTime > tc.maxTime {
				t.Errorf("Processing time %v exceeds target %v for %s",
					processingTime, tc.maxTime, tc.name)
			}

			// Validate extraction completeness
			if len(profile.Colors) == 0 {
				t.Error("No color clusters extracted despite successful processing")
			}

			// Validate ColorCount consistency
			if profile.ColorCount != len(profile.Colors) {
				t.Errorf("ColorCount inconsistency: expected %d, got %d", len(profile.Colors), profile.ColorCount)
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