package processor_test

import (
	"context"
	"image/color"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/loader"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/processor"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
)

func TestProcessor_New(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	
	if p == nil {
		t.Fatal("Expected processor to be created, got nil")
	}
}

func TestProcessor_ProcessImage_SimpleImage(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	
	// Use simple test image
	ctx := context.Background()
	img, err := l.LoadImage(ctx, "../../tests/images/simple.png")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	
	result, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("Expected no error for simple image, got: %v", err)
	}
	
	if result == nil {
		t.Fatal("Expected result to be non-nil")
	}
	
	// Verify all color roles are assigned
	if result.Colors.Background == (color.RGBA{}) {
		t.Error("Expected background color to be assigned")
	}
	if result.Colors.Foreground == (color.RGBA{}) {
		t.Error("Expected foreground color to be assigned")
	}
	if result.Colors.Primary == (color.RGBA{}) {
		t.Error("Expected primary color to be assigned")
	}
}

func TestProcessor_ProcessImage_GrayscaleImage(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	
	// Use grayscale test image
	ctx := context.Background()
	img, err := l.LoadImage(ctx, "../../tests/images/grayscale.jpeg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	
	result, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("Expected no error for grayscale image, got: %v", err)
	}
	
	// Verify contrast between background and foreground
	contrast := chromatic.ContrastRatio(result.Colors.Background, result.Colors.Foreground)
	if contrast < s.MinContrastRatio {
		t.Errorf("Expected contrast ratio >= %v, got %v", s.MinContrastRatio, contrast)
	}
}

func TestProcessor_ProcessImage_MonochromeImage(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	
	// Use monochrome test image
	ctx := context.Background()
	img, err := l.LoadImage(ctx, "../../tests/images/monochrome.jpeg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	
	result, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("Expected no error for monochrome image, got: %v", err)
	}
	
	// Verify that colors were extracted successfully
	if result.Colors.Primary == (color.RGBA{}) {
		t.Error("Expected primary color to be assigned for monochrome image")
	}
}

func TestProcessor_ProcessImage_LowFrequencyFiltering(t *testing.T) {
	s := settings.DefaultSettings()
	s.MinFrequency = 0.9 // 90% minimum frequency - very strict
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	
	// Use complex multi-color image that won't meet strict frequency requirements
	ctx := context.Background()
	img, err := l.LoadImage(ctx, "../../tests/images/abstract.jpeg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	
	result, err := p.ProcessImage(img)
	if err == nil {
		t.Error("Expected error when no colors meet frequency threshold")
	}
	if result != nil {
		t.Error("Expected nil result when processing fails")
	}
}

// Test theme mode detection with known images
func TestProcessor_ThemeMode_Detection(t *testing.T) {
	testCases := []struct {
		name         string
		imagePath    string
		expectedMode processor.ThemeMode
		description  string
	}{
		{
			name:         "Dark night city suggests Dark theme",
			imagePath:    "../../tests/images/night-city.jpeg",
			expectedMode: processor.Dark,
			description:  "Night scene should pair with dark theme",
		},
		{
			name:         "Bright coastal scene suggests Light theme", 
			imagePath:    "../../tests/images/coast.jpeg",
			expectedMode: processor.Light,
			description:  "Bright coastal scene should pair with light theme",
		},
		{
			name:         "Warm image suggests appropriate theme",
			imagePath:    "../../tests/images/warm.jpeg",
			expectedMode: processor.Dark, // Based on actual analysis: bg luminance 0.036
			description:  "Warm image is actually quite dark overall",
		},
	}
	
	s := settings.DefaultSettings()
	l := loader.NewFileLoader(s)
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			img, err := l.LoadImage(ctx, tc.imagePath)
			if err != nil {
				t.Fatalf("Failed to load test image: %v", err)
			}
			
			p := processor.New(s)
			result, err := p.ProcessImage(img)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			
			// Calculate actual luminance metrics for diagnosis
			bgLuminance := chromatic.Luminance(result.Colors.Background)
			fgLuminance := chromatic.Luminance(result.Colors.Foreground)
			primaryLuminance := chromatic.Luminance(result.Colors.Primary)
			
			// Log diagnostic information
			t.Logf("Image: %s", tc.imagePath)
			t.Logf("Background luminance: %v", bgLuminance)
			t.Logf("Foreground luminance: %v", fgLuminance)
			t.Logf("Primary luminance: %v", primaryLuminance)
			t.Logf("Theme mode threshold: %v", s.ThemeModeThreshold)
			
			// Determine actual theme mode based on background selection
			actualMode := processor.Dark
			if bgLuminance >= 0.5 {
				actualMode = processor.Light
			}
			
			t.Logf("Expected mode: %v, Actual mode: %v", tc.expectedMode, actualMode)
			
			if actualMode != tc.expectedMode {
				t.Errorf("Theme mode mismatch for %s: expected %v, got %v (bg luminance: %v, threshold: %v)",
					tc.imagePath, tc.expectedMode, actualMode, bgLuminance, s.ThemeModeThreshold)
			}
		})
	}
}

// Test fallback color parsing
func TestProcessor_FallbackColors(t *testing.T) {
	s := settings.DefaultSettings()
	
	// Test valid hex colors
	s.LightBackgroundFallback = "#ffffff"
	s.DarkBackgroundFallback = "#000000"
	s.PrimaryFallback = "#ff5733"
	
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	
	// Use a mid-tone image that might trigger fallback paths
	ctx := context.Background()
	img, err := l.LoadImage(ctx, "../../tests/images/sepia.jpeg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	
	result, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Verify colors are assigned (fallbacks may or may not be used)
	if result.Colors.Background == (color.RGBA{}) {
		t.Error("Expected background color to be assigned")
	}
	if result.Colors.Primary == (color.RGBA{}) {
		t.Error("Expected primary color to be assigned")
	}
}

func TestProcessor_InvalidFallbackColors(t *testing.T) {
	s := settings.DefaultSettings()
	
	// Set invalid hex colors
	s.LightBackgroundFallback = "invalid-hex"
	s.DarkBackgroundFallback = "not-a-color"
	s.PrimaryFallback = "#gggggg" // Invalid hex digits
	
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	
	// Use any test image
	ctx := context.Background()
	img, err := l.LoadImage(ctx, "../../tests/images/simple.png")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	
	result, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("Processor should handle invalid hex gracefully, got: %v", err)
	}
	
	// Should still get valid colors from hardcoded fallbacks
	if result.Colors.Background == (color.RGBA{}) {
		t.Error("Expected hardcoded fallback background color")
	}
	if result.Colors.Primary == (color.RGBA{}) {
		t.Error("Expected hardcoded fallback primary color")
	}
}

func TestProcessor_ColorExtractionQuality(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	
	testCases := []struct {
		name      string
		imagePath string
		checkFunc func(*testing.T, *processor.ColorProfile)
	}{
		{
			name:      "Portal image extracts vibrant colors",
			imagePath: "../../tests/images/portal.jpeg",
			checkFunc: func(t *testing.T, result *processor.ColorProfile) {
				// Portal image should have distinctive colors
				if result.Colors.Primary == result.Colors.Background {
					t.Error("Primary color should be different from background")
				}
			},
		},
		{
			name:      "Nebula image handles complex colors",
			imagePath: "../../tests/images/nebula.jpeg",
			checkFunc: func(t *testing.T, result *processor.ColorProfile) {
				// Nebula should extract meaningful colors
				contrast := chromatic.ContrastRatio(result.Colors.Background, result.Colors.Foreground)
				if contrast < s.MinContrastRatio {
					t.Errorf("Nebula should maintain contrast ratio >= %v, got %v", s.MinContrastRatio, contrast)
				}
			},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			img, err := l.LoadImage(ctx, tc.imagePath)
			if err != nil {
				t.Fatalf("Failed to load test image: %v", err)
			}
			
			result, err := p.ProcessImage(img)
			if err != nil {
				t.Fatalf("Unexpected error processing %s: %v", tc.imagePath, err)
			}
			
			tc.checkFunc(t, result)
		})
	}
}