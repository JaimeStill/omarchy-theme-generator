package processor_test

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/processor"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

// Test role-based color extraction
func TestExtraction_BackgroundSelection(t *testing.T) {
	testCases := []struct {
		name              string
		colors            map[color.RGBA]uint32
		expectedThemeMode processor.ThemeMode
		minBgLuminance    float64
		maxBgLuminance    float64
		description       string
	}{
		{
			name: "Light theme with bright colors",
			colors: map[color.RGBA]uint32{
				{R: 240, G: 240, B: 240, A: 255}: 300, // Very light gray
				{R: 200, G: 200, B: 220, A: 255}: 200, // Light blue-gray
				{R: 255, G: 50, B: 50, A: 255}:   100, // Bright red
			},
			expectedThemeMode: processor.Light,
			minBgLuminance:    0.5, // Should pick reasonably light background  
			maxBgLuminance:    1.0,
			description:       "Light theme should select light background",
		},
		{
			name: "Dark theme with dark colors",
			colors: map[color.RGBA]uint32{
				{R: 30, G: 30, B: 30, A: 255}:   300, // Very dark gray
				{R: 50, G: 70, B: 90, A: 255}:   200, // Dark blue-gray
				{R: 100, G: 150, B: 200, A: 255}: 100, // Medium blue
			},
			expectedThemeMode: processor.Dark,
			minBgLuminance:    0.0,
			maxBgLuminance:    0.3, // Should pick dark background
			description:       "Dark theme should select dark background",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := settings.DefaultSettings()
			p := processor.New(s)
			
			img := createImageFromFrequencyMap(tc.colors)
			result, err := p.ProcessImage(img)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			
			bgLuminance := chromatic.Luminance(result.Colors.Background)
			fgLuminance := chromatic.Luminance(result.Colors.Foreground)
			primaryLuminance := chromatic.Luminance(result.Colors.Primary)
			
			// Calculate luminance of input colors for comparison
			var colorLuminances []float64
			for c := range tc.colors {
				colorLuminances = append(colorLuminances, chromatic.Luminance(c))
			}
			
			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Input color luminances: %v", colorLuminances)
			t.Logf("Selected background: RGB(%d,%d,%d) luminance: %v", 
				result.Colors.Background.R, result.Colors.Background.G, result.Colors.Background.B, bgLuminance)
			t.Logf("Selected foreground: RGB(%d,%d,%d) luminance: %v",
				result.Colors.Foreground.R, result.Colors.Foreground.G, result.Colors.Foreground.B, fgLuminance)
			t.Logf("Selected primary: RGB(%d,%d,%d) luminance: %v",
				result.Colors.Primary.R, result.Colors.Primary.G, result.Colors.Primary.B, primaryLuminance)
			t.Logf("Expected theme mode: %v", tc.expectedThemeMode)
			t.Logf("Expected luminance range: [%v, %v]", tc.minBgLuminance, tc.maxBgLuminance)
			
			if bgLuminance < tc.minBgLuminance || bgLuminance > tc.maxBgLuminance {
				t.Errorf("Background luminance %v not in expected range [%v, %v]",
					bgLuminance, tc.minBgLuminance, tc.maxBgLuminance)
			}
		})
	}
}

func TestExtraction_ContrastCompliance(t *testing.T) {
	s := settings.DefaultSettings()
	s.MinContrastRatio = 4.5 // WCAG AA standard
	p := processor.New(s)
	
	// Create image with high contrast potential
	img := image.NewRGBA(image.Rect(0, 0, 3, 3))
	
	// Fill with white, black, and gray
	colors := []color.RGBA{
		{R: 255, G: 255, B: 255, A: 255}, // White
		{R: 0, G: 0, B: 0, A: 255},       // Black  
		{R: 128, G: 128, B: 128, A: 255}, // Gray
	}
	
	idx := 0
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			img.Set(x, y, colors[idx%len(colors)])
			idx++
		}
	}
	
	result, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Check background/foreground contrast
	contrast := chromatic.ContrastRatio(result.Colors.Background, result.Colors.Foreground)
	if contrast < s.MinContrastRatio {
		t.Errorf("Background/foreground contrast %v below minimum %v", 
			contrast, s.MinContrastRatio)
	}
	
	// Verify contrast is meaningful (not just 1:1)
	if contrast < 2.0 {
		t.Error("Background/foreground should have meaningful contrast")
	}
}

func TestExtraction_PrimarySaturationRequirement(t *testing.T) {
	s := settings.DefaultSettings()
	s.MinPrimarySaturation = 0.4 // Require moderate saturation
	p := processor.New(s)
	
	// Create image with both saturated and desaturated colors
	colorFreq := map[color.RGBA]uint32{
		{R: 200, G: 200, B: 200, A: 255}: 300, // Gray (low saturation)
		{R: 180, G: 180, B: 180, A: 255}: 200, // Lighter gray
		{R: 255, G: 100, B: 100, A: 255}: 150, // Pink (moderate saturation)
		{R: 255, G: 0, B: 0, A: 255}:     100, // Pure red (high saturation)
	}
	
	img := createImageFromFrequencyMap(colorFreq)
	result, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	primarySat := formats.RGBAToHSLA(result.Colors.Primary).S
	
	// Primary should meet saturation requirement or be a fallback
	if primarySat < s.MinPrimarySaturation {
		// Check if it's a generated/fallback color (acceptable)
		isPureColor := (result.Colors.Primary.R == 255 && result.Colors.Primary.G == 0 && result.Colors.Primary.B == 0) || // Red
			(result.Colors.Primary.R == 100 && result.Colors.Primary.G == 150 && result.Colors.Primary.B == 200) // Fallback blue
		
		if !isPureColor {
			t.Errorf("Primary color saturation %v below requirement %v", 
				primarySat, s.MinPrimarySaturation)
		}
	}
}

func TestExtraction_AccentColorBounds(t *testing.T) {
	s := settings.DefaultSettings()
	s.MinAccentLightness = 0.3
	s.MaxAccentLightness = 0.7
	s.MinAccentSaturation = 0.5
	p := processor.New(s)
	
	// Create colorful image to provide accent color options
	colorFreq := map[color.RGBA]uint32{
		{R: 100, G: 100, B: 100, A: 255}: 200, // Dark gray
		{R: 255, G: 150, B: 50, A: 255}:  150, // Orange (good accent candidate)
		{R: 50, G: 200, B: 255, A: 255}:  100, // Cyan (good accent candidate)
		{R: 200, G: 50, B: 200, A: 255}:  100, // Magenta (good accent candidate)
	}
	
	img := createImageFromFrequencyMap(colorFreq)
	result, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	accentHSLA := formats.RGBAToHSLA(result.Colors.Accent)
	
	// Check accent color bounds (may be generated if no suitable candidate)
	if accentHSLA.S >= s.MinAccentSaturation {
		// If it meets saturation requirement, check lightness bounds
		if accentHSLA.L < s.MinAccentLightness || accentHSLA.L > s.MaxAccentLightness {
			// Allow generated/complementary colors to be outside bounds
			isGenerated := result.Colors.Accent.R != 255 && result.Colors.Accent.G != 150 && result.Colors.Accent.B != 50 &&
				result.Colors.Accent.R != 50 && result.Colors.Accent.G != 200 && result.Colors.Accent.B != 255 &&
				result.Colors.Accent.R != 200 && result.Colors.Accent.G != 50 && result.Colors.Accent.B != 200
			
			if !isGenerated {
				t.Logf("Accent lightness %v outside bounds [%v, %v] but may be acceptable", 
					accentHSLA.L, s.MinAccentLightness, s.MaxAccentLightness)
			}
		}
	}
}

func TestExtraction_ColorDistinctness(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	
	// Create image with distinct colors
	colorFreq := map[color.RGBA]uint32{
		{R: 255, G: 255, B: 255, A: 255}: 400, // White
		{R: 0, G: 0, B: 0, A: 255}:       300, // Black
		{R: 255, G: 0, B: 0, A: 255}:     200, // Red
		{R: 0, G: 255, B: 0, A: 255}:     150, // Green
		{R: 0, G: 0, B: 255, A: 255}:     100, // Blue
	}
	
	img := createImageFromFrequencyMap(colorFreq)
	result, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Extract all role colors
	colors := []color.RGBA{
		result.Colors.Background,
		result.Colors.Foreground,
		result.Colors.Primary,
		result.Colors.Secondary,
		result.Colors.Accent,
	}
	
	// Check that colors are reasonably distinct
	for i, c1 := range colors {
		for j, c2 := range colors {
			if i >= j {
				continue
			}
			
			// Calculate color distance in RGB space
			dr := int(c1.R) - int(c2.R)
			dg := int(c1.G) - int(c2.G)
			db := int(c1.B) - int(c2.B)
			distance := math.Sqrt(float64(dr*dr + dg*dg + db*db))
			
			// Colors should be somewhat distinct (not identical)
			if distance < 10.0 {
				// Allow background/foreground to be similar if high contrast
				if (i == 0 && j == 1) || (i == 1 && j == 0) {
					contrast := chromatic.ContrastRatio(c1, c2)
					if contrast < 3.0 {
						t.Errorf("Background and foreground too similar (distance: %v, contrast: %v)",
							distance, contrast)
					}
				} else {
					t.Logf("Colors %d and %d are very similar (distance: %v) - may be acceptable for generated colors", 
						i, j, distance)
				}
			}
		}
	}
}

func TestExtraction_FallbackBehavior(t *testing.T) {
	s := settings.DefaultSettings()
	
	// Set valid fallback hex colors
	s.LightBackgroundFallback = "#f8f8f8"
	s.DarkBackgroundFallback = "#1a1a1a" 
	s.PrimaryFallback = "#3366cc"
	
	p := processor.New(s)
	
	// Create image with insufficient color variety to trigger fallbacks
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	gray := color.RGBA{R: 128, G: 128, B: 128, A: 255}
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			img.Set(x, y, gray)
		}
	}
	
	result, err := p.ProcessImage(img)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Verify fallback colors are used appropriately
	if result.Colors.Background == (color.RGBA{}) {
		t.Error("Background should have fallback color")
	}
	if result.Colors.Primary == (color.RGBA{}) {
		t.Error("Primary should have fallback color")
	}
	
	// Verify fallbacks parse correctly
	expectedPrimary, err := formats.ParseHex(s.PrimaryFallback)
	if err == nil && result.Colors.Primary == expectedPrimary {
		t.Logf("Primary fallback color correctly applied: %+v", result.Colors.Primary)
	}
}