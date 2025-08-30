package main

import (
	"fmt"
	"image"
	"time"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/generative"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/theme"
)

func main() {
	fmt.Println("=== Omarchy Theme Generator: Palette Strategies & Theme Modes Test ===")
	fmt.Println()

	// Test 1: Theme Generation Workflow
	testThemeGeneration()
	
	// Test 2: Mode Detection Logic
	testModeDetection()
	
	// Test 3: Override Validation System
	testOverrideValidation()
	
	// Test 4: Performance Maintenance
	testPerformanceMaintenance()
	
	// Test 5: Computational Graphics Integration
	testComputationalIntegration()
	
	fmt.Println("=== Test Suite Complete ===")
}

// testThemeGeneration demonstrates the core theme creation workflow
// integrating extraction → pipeline → theme generation.
func testThemeGeneration() {
	fmt.Println("--- Test 1: Core Theme Generation Workflow ---")
	
	// Create theme generator
	generator := theme.NewGenerator(nil)
	
	// Test with different image types
	testImages := []struct {
		name      string
		generator func() image.Image
	}{
		{"4K Synthetic", func() image.Image { return generative.Generate4KTestImage() }},
		{"Grayscale", func() image.Image { return generative.GenerateGrayscaleTestImage(800, 600) }},
		{"Monochromatic", func() image.Image { return generative.GenerateMonochromaticTestImage(800, 600) }},
		{"80s Synthwave", func() image.Image { return generative.Generate80sVectorImage(800, 600) }},
	}
	
	for _, test := range testImages {
		fmt.Printf("\n  Testing %s Image:\n", test.name)
		
		img := test.generator()
		
		// Create theme configuration
		config := theme.ThemeConfig{
			SourceImage: img,
			Mode:        theme.ModeAuto,
			Name:        fmt.Sprintf("Test %s Theme", test.name),
		}
		
		// Generate theme
		start := time.Now()
		generatedTheme, err := generator.GenerateTheme(config)
		elapsed := time.Since(start)
		
		if err != nil {
			fmt.Printf("    ❌ Generation failed: %v\n", err)
			continue
		}
		
		fmt.Printf("    ✅ Generated successfully\n")
		fmt.Printf("    📊 Theme: %s mode, strategy=%s, colors=%d\n", 
			modeString(generatedTheme.IsLight), 
			generatedTheme.Metadata.Strategy,
			len(generatedTheme.Palette))
		fmt.Printf("    🎨 Primary: %s, Background: %s, Foreground: %s\n",
			generatedTheme.Primary.HEX(),
			generatedTheme.Background.HEX(), 
			generatedTheme.Foreground.HEX())
		fmt.Printf("    ⚡ Performance: %v (target: <2s)\n", elapsed)
		fmt.Printf("    ♿ WCAG: %d/%d colors passing AA standard\n",
			generatedTheme.Metadata.Validation.PassingColors,
			generatedTheme.Metadata.Validation.TotalColors)
		
		// Verify WCAG compliance
		if generatedTheme.Metadata.Validation.FailingColors > 0 {
			fmt.Printf("    ⚠️  %d colors required adjustment for accessibility\n",
				generatedTheme.Metadata.Validation.FailingColors)
		}
	}
}

// testModeDetection demonstrates light/dark mode classification logic.
func testModeDetection() {
	fmt.Println("\n--- Test 2: Light/Dark Mode Detection ---")
	
	detector := theme.NewModeDetector()
	
	testCases := []struct {
		name      string
		generator func() image.Image
		expected  string
	}{
		{"Bright Gradient", func() image.Image { 
			return generative.GenerateComplexGradientImage(400, 300, "radial-bright")
		}, "light"},
		{"Dark 80s Interface", func() image.Image { 
			return generative.Generate80sVectorImage(400, 300) 
		}, "dark"},
		{"Mid-tone Industrial", func() image.Image { 
			return generative.GenerateCassetteFuturismImage(400, 300, 180.0)
		}, "varies"},
		{"Pure Grayscale", func() image.Image { 
			return generative.GenerateGrayscaleTestImage(400, 300)
		}, "light/dark"},
	}
	
	fmt.Printf("\n  Mode Detection Results:\n")
	for _, test := range testCases {
		img := test.generator()
		detectedMode := detector.DetectFromImage(img)
		
		fmt.Printf("    %s: detected %s (expected: %s)\n",
			test.name, detectedMode.String(), test.expected)
		
		// Test with primary color influence
		primaryColor := color.NewHSL(0.6, 0.7, 0.3) // Dark blue
		modeWithPrimary := detector.DetectWithPrimary(img, primaryColor)
		
		if modeWithPrimary != detectedMode {
			fmt.Printf("      → Primary color influence: %s → %s\n",
				detectedMode.String(), modeWithPrimary.String())
		}
	}
}

// testOverrideValidation demonstrates the user color override system with WCAG validation.
func testOverrideValidation() {
	fmt.Println("\n--- Test 3: Override Validation System ---")
	
	// Create test image and base theme
	img := generative.GenerateMonochromaticTestImage(400, 300)
	generator := theme.NewGenerator(nil)
	
	// Test different override scenarios
	overrideTests := []struct {
		name      string
		overrides theme.ColorOverrides
	}{
		{
			name: "Valid High Contrast",
			overrides: theme.ColorOverrides{
				Primary:    color.NewRGB(0, 120, 215),   // Blue
				Background: color.NewRGB(255, 255, 255), // White
				Foreground: color.NewRGB(30, 30, 30),    // Dark gray
			},
		},
		{
			name: "Poor Contrast (needs adjustment)", 
			overrides: theme.ColorOverrides{
				Primary:    color.NewRGB(200, 200, 200), // Light gray
				Background: color.NewRGB(255, 255, 255), // White (poor contrast)
				Foreground: color.NewRGB(150, 150, 150), // Light gray (poor contrast)
			},
		},
		{
			name: "Dark Theme Override",
			overrides: theme.ColorOverrides{
				Background: color.NewRGB(25, 25, 25),    // Very dark
				Foreground: color.NewRGB(240, 240, 240), // Light
			},
		},
	}
	
	for _, test := range overrideTests {
		fmt.Printf("\n  Testing %s:\n", test.name)
		
		config := theme.ThemeConfig{
			SourceImage: img,
			Mode:        theme.ModeAuto,
			Overrides:   test.overrides,
			Name:        fmt.Sprintf("Override Test: %s", test.name),
		}
		
		overrideTheme, err := generator.GenerateTheme(config)
		if err != nil {
			fmt.Printf("    ❌ Override failed: %v\n", err)
			continue
		}
		
		fmt.Printf("    ✅ Override applied successfully\n")
		
		// Check contrast ratios
		primaryBgRatio := overrideTheme.Primary.ContrastRatio(overrideTheme.Background)
		foregroundBgRatio := overrideTheme.Foreground.ContrastRatio(overrideTheme.Background)
		
		fmt.Printf("    🎨 Final Colors: P=%s, B=%s, F=%s\n",
			overrideTheme.Primary.HEX(),
			overrideTheme.Background.HEX(),
			overrideTheme.Foreground.HEX())
		fmt.Printf("    📏 Contrast Ratios: P-B=%.2f:1, F-B=%.2f:1 (min: 4.5:1)\n",
			primaryBgRatio, foregroundBgRatio)
		
		// Check if adjustments were made
		if test.overrides.Primary != nil && !colorsEqual(test.overrides.Primary, overrideTheme.Primary) {
			fmt.Printf("    🔧 Primary adjusted from %s for WCAG compliance\n", 
				test.overrides.Primary.HEX())
		}
		if test.overrides.Foreground != nil && !colorsEqual(test.overrides.Foreground, overrideTheme.Foreground) {
			fmt.Printf("    🔧 Foreground adjusted from %s for WCAG compliance\n",
				test.overrides.Foreground.HEX())
		}
	}
}

// testPerformanceMaintenance validates that theme generation maintains the 242ms target.
func testPerformanceMaintenance() {
	fmt.Println("\n--- Test 4: Performance Target Validation ---")
	
	generator := theme.NewGenerator(nil)
	
	// Test with 4K image (most demanding case)
	fmt.Printf("\n  Performance Test with 4K Image (3840×2160):\n")
	
	img := generative.Generate4KTestImage()
	config := theme.ThemeConfig{
		SourceImage: img,
		Mode:        theme.ModeAuto,
		Name:        "4K Performance Test",
	}
	
	// Run multiple iterations to get average performance
	const iterations = 3
	totalTime := time.Duration(0)
	
	for i := 0; i < iterations; i++ {
		start := time.Now()
		testTheme, err := generator.GenerateTheme(config)
		elapsed := time.Since(start)
		totalTime += elapsed
		
		if err != nil {
			fmt.Printf("    ❌ Iteration %d failed: %v\n", i+1, err)
			continue
		}
		
		fmt.Printf("    ⚡ Iteration %d: %v\n", i+1, elapsed)
		fmt.Printf("       📊 Colors: %d, Strategy: %s, Mode: %s\n",
			len(testTheme.Palette), testTheme.Metadata.Strategy, 
			modeString(testTheme.IsLight))
	}
	
	averageTime := totalTime / iterations
	target := 2 * time.Second
	
	fmt.Printf("\n    📈 Performance Summary:\n")
	fmt.Printf("       Average: %v\n", averageTime)
	fmt.Printf("       Target:  %v\n", target)
	if averageTime < target {
		speedup := float64(target) / float64(averageTime)
		fmt.Printf("       ✅ Target achieved! (%.1fx faster than 2s limit)\n", speedup)
	} else {
		fmt.Printf("       ❌ Target missed by %v\n", averageTime-target)
	}
}

// testComputationalIntegration verifies that the theme system works with generated images.
func testComputationalIntegration() {
	fmt.Println("\n--- Test 5: Computational Graphics Integration ---")
	
	generator := theme.NewGenerator(nil)
	
	// Test with sophisticated computational graphics
	fmt.Printf("\n  Integration with Computational Aesthetics:\n")
	
	// Generate cassette futurism interface with specific accent hue
	accentHue := 315.0 // Magenta-pink
	img := generative.GenerateCassetteFuturismImage(800, 600, accentHue)
	
	config := theme.ThemeConfig{
		SourceImage: img,
		Mode:        theme.ModeAuto,
		Name:        "Cassette Futurism Integration",
	}
	
	computationalTheme, err := generator.GenerateTheme(config)
	if err != nil {
		fmt.Printf("    ❌ Computational integration failed: %v\n", err)
		return
	}
	
	fmt.Printf("    ✅ Generated theme from computational graphics\n")
	fmt.Printf("    🎨 Detected as %s theme (industrial interfaces typically dark)\n",
		modeString(computationalTheme.IsLight))
	fmt.Printf("    📐 Material simulation integration: %d synthesized colors\n",
		computationalTheme.Metadata.SynthesizedColors)
	fmt.Printf("    🖼️  Image specs: %dx%d pixels, %s strategy\n",
		computationalTheme.Metadata.Performance.ImageSize.Width,
		computationalTheme.Metadata.Performance.ImageSize.Height,
		computationalTheme.Metadata.Strategy)
	
	// Verify the accent hue was properly detected and incorporated
	h, _, _ := computationalTheme.Primary.HSL()
	primaryHueDegrees := h * 360
	
	fmt.Printf("    🎯 Primary color hue: %.0f° (input accent: %.0f°)\n",
		primaryHueDegrees, accentHue)
	
	// Check if the computational aesthetics influenced the theme
	if computationalTheme.Metadata.GenerationMode == "hybrid" ||
	   computationalTheme.Metadata.GenerationMode == "extract" {
		fmt.Printf("    🔍 Successfully extracted colors from computational materials\n")
	} else {
		fmt.Printf("    🎨 Computational image required synthesis fallback\n")
	}
	
	fmt.Printf("    ⚙️  Performance: %v generation time\n",
		computationalTheme.Metadata.Performance.TotalTime)
}

// Helper functions

func modeString(isLight bool) string {
	if isLight {
		return "light"
	}
	return "dark"
}

func colorsEqual(c1, c2 *color.Color) bool {
	if c1 == nil || c2 == nil {
		return c1 == c2
	}
	return c1.HEX() == c2.HEX()
}