package main

import (
	"fmt"
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/palette"
	"github.com/JaimeStill/omarchy-theme-generator/tests"
)

func main() {
	fmt.Println("=== Omarchy Theme Generator - Color Synthesis Validation ===")

	// Test 1: Color Theory Algorithm Accuracy
	fmt.Println("\nTest 1: Color Theory Algorithm Accuracy")
	
	// Base color: Red (0° hue)
	redBase := color.NewRGB(255, 0, 0) // H=0°, S=100%, L=50%
	h, s, l := redBase.HSL()
	fmt.Printf("Base color: Red RGB(255,0,0) → HSL(%.0f°, %.0f%%, %.0f%%)\n", 
		h*360, s*100, l*100)

	// Test each synthesis strategy
	strategies := []struct {
		name     string
		strategy palette.SynthesisStrategy
		expected []float64 // Expected hue values in degrees
	}{
		{
			name:     "Monochromatic",
			strategy: &palette.MonochromaticStrategy{},
			expected: []float64{0, 0, 0, 0}, // Same hue variations
		},
		{
			name:     "Complementary", 
			strategy: &palette.ComplementaryStrategy{},
			expected: []float64{0, 180}, // Base + complement
		},
		{
			name:     "Triadic",
			strategy: &palette.TriadicStrategy{},
			expected: []float64{0, 120, 240}, // 120° apart
		},
		{
			name:     "Analogous",
			strategy: &palette.AnalogousStrategy{},
			expected: []float64{0}, // ±30° range around base
		},
		{
			name:     "Split-Complementary",
			strategy: &palette.SplitComplementaryStrategy{},
			expected: []float64{0, 150, 210}, // Base + ±30° from complement
		},
	}

	for _, test := range strategies {
		fmt.Printf("\n%s Strategy:\n", test.name)
		
		generated := test.strategy.Generate(redBase, 6)
		
		// Extract hues from generated palette
		var hues []float64
		hueMap := make(map[int]bool) // Track unique hues (rounded to nearest 10°)
		
		for i, c := range generated {
			h, s, l := c.HSL()
			hue := h * 360
			hues = append(hues, hue)
			hueKey := int(math.Round(hue/10) * 10) // Round to nearest 10°
			hueMap[hueKey] = true
			
			fmt.Printf("  Color %d: H=%.0f°, S=%.0f%%, L=%.0f%%\n", 
				i+1, hue, s*100, l*100)
		}
		
		// Validate color theory principles
		switch test.name {
		case "Monochromatic":
			// All colors should have same hue (±5° tolerance)
			baseHue := hues[0]
			monochromatic := true
			for _, hue := range hues[1:] {
				if math.Abs(hue-baseHue) > 5 {
					monochromatic = false
					break
				}
			}
			fmt.Printf("  Validation: %s - All hues within ±5° of base\n", 
				tests.CheckMark(monochromatic))
			
		case "Complementary":
			// Should find both base (0°) and complement (180°)
			hasBase := false
			hasComplement := false
			for hue := range hueMap {
				if math.Abs(float64(hue)) < 10 { // Near 0°
					hasBase = true
				}
				if math.Abs(float64(hue)-180) < 10 { // Near 180°
					hasComplement = true
				}
			}
			valid := hasBase && hasComplement
			fmt.Printf("  Validation: %s - Contains base (0°) and complement (180°)\n",
				tests.CheckMark(valid))
			
		case "Triadic":
			// Should find hues near 0°, 120°, 240°
			expectedHues := []float64{0, 120, 240}
			found := 0
			for _, expected := range expectedHues {
				for hue := range hueMap {
					if math.Abs(float64(hue)-expected) < 15 { // 15° tolerance
						found++
						break
					}
				}
			}
			valid := found >= 3
			fmt.Printf("  Validation: %s - Contains triadic hues (0°, 120°, 240°) ± 15°\n",
				tests.CheckMark(valid))
			
		case "Analogous":
			// All hues should be within ±30° of base
			baseHue := 0.0
			analogous := true
			for _, hue := range hues {
				distance := math.Min(math.Abs(hue-baseHue), 360-math.Abs(hue-baseHue))
				if distance > 35 { // 35° tolerance for analogous
					analogous = false
					break
				}
			}
			fmt.Printf("  Validation: %s - All hues within ±35° of base\n",
				tests.CheckMark(analogous))
			
		case "Split-Complementary":
			// Should find base (0°) and split complements (150°, 210°)
			hasBase := false
			hasSplit1 := false
			hasSplit2 := false
			for hue := range hueMap {
				if math.Abs(float64(hue)) < 10 {
					hasBase = true
				}
				if math.Abs(float64(hue)-150) < 15 {
					hasSplit1 = true
				}
				if math.Abs(float64(hue)-210) < 15 {
					hasSplit2 = true
				}
			}
			valid := hasBase && (hasSplit1 || hasSplit2)
			fmt.Printf("  Validation: %s - Contains base and split complement hues\n",
				tests.CheckMark(valid))
		}
		
		fmt.Printf("  Description: %s\n", test.strategy.Description())
	}

	// Test 2: Temperature-Matched Grays
	fmt.Println("\nTest 2: Temperature-Matched Grays")
	
	warmBase := color.NewHSLA(30.0/360.0, 0.8, 0.6, 1.0) // Orange (warm)
	coolBase := color.NewHSLA(210.0/360.0, 0.8, 0.6, 1.0) // Blue (cool)
	
	warmGrays := palette.GenerateTemperatureMatchedGrays(warmBase, 4)
	coolGrays := palette.GenerateTemperatureMatchedGrays(coolBase, 4)
	
	fmt.Printf("Warm base (30°): %s\n", warmBase.CSSHSL())
	fmt.Printf("Generated warm grays:\n")
	for i, gray := range warmGrays {
		h, s, l := gray.HSL()
		fmt.Printf("  Gray %d: H=%.0f°, S=%.1f%%, L=%.0f%%\n", 
			i+1, h*360, s*100, l*100)
	}
	
	fmt.Printf("\nCool base (210°): %s\n", coolBase.CSSHSL())
	fmt.Printf("Generated cool grays:\n")
	for i, gray := range coolGrays {
		h, s, l := gray.HSL()
		fmt.Printf("  Gray %d: H=%.0f°, S=%.1f%%, L=%.0f%%\n", 
			i+1, h*360, s*100, l*100)
	}
	
	// Validate temperature matching
	warmH, _, _ := warmGrays[0].HSL()
	coolH, _, _ := coolGrays[0].HSL()
	
	warmTemp := (warmH*360 < 60) || (warmH*360 > 300) // Red/yellow range
	coolTemp := (coolH*360 >= 180) && (coolH*360 <= 240) // Blue range
	
	fmt.Printf("\nTemperature validation:\n")
	fmt.Printf("  Warm grays use warm hue: %s (H=%.0f°)\n", 
		tests.CheckMark(warmTemp), warmH*360)
	fmt.Printf("  Cool grays use cool hue: %s (H=%.0f°)\n", 
		tests.CheckMark(coolTemp), coolH*360)

	// Test 3: WCAG Compliance Integration
	fmt.Println("\nTest 3: WCAG Compliance Integration")
	
	generator := palette.NewPaletteGenerator(palette.DefaultSynthesisOptions())
	
	// Generate palette and test contrast
	testBase := color.NewHSLA(220.0/360.0, 0.7, 0.5, 1.0) // Blue base
	generatedPalette, err := generator.GenerateFromBase(testBase, "complementary")
	
	if err != nil {
		fmt.Printf("Error generating palette: %v\n", err)
		return
	}
	
	fmt.Printf("Generated palette (%d colors):\n", len(generatedPalette))
	
	// Test against white background
	white := color.NewRGB(255, 255, 255)
	aaCompliant := 0
	
	for i, c := range generatedPalette[:6] { // Test first 6 colors
		contrast := c.ContrastRatio(white)
		compliant := contrast >= 4.5 // WCAG AA
		if compliant {
			aaCompliant++
		}
		
		fmt.Printf("  Color %d: %s, contrast=%.2f:1 %s\n", 
			i+1, c.CSSHSL(), contrast, 
			tests.CheckMark(compliant))
	}
	
	complianceRate := float64(aaCompliant) / 6.0
	fmt.Printf("WCAG AA compliance: %.0f%% (%d/6 colors) %s\n", 
		complianceRate*100, aaCompliant,
		tests.CheckMark(complianceRate >= 0.5)) // At least 50% should be compliant

	// Test 4: Hue Distance and Normalization
	fmt.Println("\nTest 4: Hue Distance and Normalization")
	
	// Test hue normalization edge cases
	testCases := []struct {
		input    float64
		expected float64
	}{
		{1.5, 0.5},   // > 1.0
		{-0.3, 0.7},  // < 0.0  
		{0.5, 0.5},   // normal
		{2.7, 0.7},   // > 2.0
	}
	
	fmt.Printf("Hue normalization tests:\n")
	for _, test := range testCases {
		result := palette.NormalizeHue(test.input)
		correct := math.Abs(result-test.expected) < 0.01
		fmt.Printf("  Input: %.1f → Output: %.3f (expected %.1f) %s\n", 
			test.input, result, test.expected, tests.CheckMark(correct))
	}
	
	// Test degree conversions
	degrees := 270.0
	hue := palette.DegreesToHue(degrees)
	backToDegrees := palette.HueToDegrees(hue)
	
	fmt.Printf("\nDegree conversion test:\n")
	fmt.Printf("  270° → %.3f → %.0f° %s\n", 
		hue, backToDegrees, 
		tests.CheckMark(math.Abs(backToDegrees-degrees) < 0.1))

	// Test 5: Strategy Registration and Fallback
	fmt.Println("\nTest 5: Strategy Registration and Fallback")
	
	// Test with invalid strategy name
	fallbackPalette, err := generator.GenerateFromBase(testBase, "invalid-strategy")
	if err != nil {
		fmt.Printf("Error with fallback: %v\n", err)
	} else {
		fmt.Printf("Fallback test: %s - Generated %d colors with invalid strategy name\n",
			tests.CheckMark(len(fallbackPalette) > 0), len(fallbackPalette))
	}
	
	// List all registered strategies
	fmt.Printf("\nRegistered strategies validation:\n")
	expectedStrategies := []string{
		"monochromatic", "analogous", "complementary", 
		"triadic", "tetradic", "split-complementary",
	}
	
	for _, strategyName := range expectedStrategies {
		palette, err := generator.GenerateFromBase(testBase, strategyName)
		success := err == nil && len(palette) > 0
		fmt.Printf("  %s: %s (%d colors)\n", 
			strategyName, tests.CheckMark(success), len(palette))
	}

	fmt.Println("\n=== Color Synthesis Validation Complete ===")
}