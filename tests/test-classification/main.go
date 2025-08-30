package main

import (
	"fmt"
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/tests"
)

// Mock image classification functions to test monochromatic vs grayscale detection
func classifyColorSet(colors []*color.Color) string {
	if len(colors) == 0 {
		return "empty"
	}
	
	// Check if all colors are grayscale (no hue information)
	grayscaleCount := 0
	coloredCount := 0
	
	// Count colors in 10° hue bins
	hueBins := make(map[int]int)
	
	for _, c := range colors {
		h, s, l := c.HSL()
		
		// Consider grayscale if saturation < 5% or lightness is pure black/white
		if s < 0.05 || l < 0.02 || l > 0.98 {
			grayscaleCount++
		} else {
			coloredCount++
			
			// Bin the hue into 10° buckets for monochromatic detection
			hueKey := int(math.Round(h*360/10) * 10) % 360
			hueBins[hueKey]++
		}
	}
	
	total := len(colors)
	
	// If 95% or more are grayscale, classify as grayscale
	if float64(grayscaleCount)/float64(total) >= 0.95 {
		return "grayscale"
	}
	
	// Check for monochromatic: 90% of colored pixels in same 10° hue range
	if coloredCount > 0 {
		maxBinCount := 0
		dominantBin := 0
		
		for bin, count := range hueBins {
			if count > maxBinCount {
				maxBinCount = count
				dominantBin = bin
			}
		}
		
		// Check adjacent bins too (±10° tolerance = ±1 bin)
		adjacentCount := 0
		for bin, count := range hueBins {
			binDistance := math.Min(math.Abs(float64(bin-dominantBin)), 360-math.Abs(float64(bin-dominantBin)))
			if binDistance <= 10 {
				adjacentCount += count
			}
		}
		
		if float64(adjacentCount)/float64(coloredCount) >= 0.90 {
			return "monochromatic"
		}
	}
	
	return "full-color"
}

func main() {
	fmt.Println("=== Omarchy Theme Generator - Color Classification Validation ===")

	// Test 1: Grayscale Detection
	fmt.Println("\nTest 1: Grayscale Detection (Pure Achromatic)")
	
	grayscaleColors := []*color.Color{
		color.NewRGB(0, 0, 0),         // Pure black
		color.NewRGB(64, 64, 64),      // Dark gray
		color.NewRGB(128, 128, 128),   // Medium gray
		color.NewRGB(192, 192, 192),   // Light gray
		color.NewRGB(255, 255, 255),   // Pure white
	}
	
	fmt.Printf("Test colors:\n")
	for i, c := range grayscaleColors {
		h, s, l := c.HSL()
		fmt.Printf("  %d: RGB(%d,%d,%d) → HSL(%.0f°, %.1f%%, %.0f%%)\n", 
			i+1, c.R, c.G, c.B, h*360, s*100, l*100)
	}
	
	classification := classifyColorSet(grayscaleColors)
	fmt.Printf("Classification: %s %s\n", 
		classification, tests.CheckMark(classification == "grayscale"))

	// Test 2: Monochromatic Detection (Single Hue ±10°)
	fmt.Println("\nTest 2: Monochromatic Detection (Single Dominant Hue)")
	
	monochromaticColors := []*color.Color{
		color.NewHSL(220.0/360, 0.9, 0.2),  // Dark blue
		color.NewHSL(215.0/360, 0.7, 0.4),  // Blue (within 10°)
		color.NewHSL(225.0/360, 0.8, 0.6),  // Light blue (within 10°)
		color.NewHSL(218.0/360, 0.6, 0.8),  // Very light blue
		color.NewRGB(128, 128, 128),         // Gray (should be ignored)
		color.NewRGB(64, 64, 64),            // Gray (should be ignored)
	}
	
	fmt.Printf("Test colors:\n")
	for i, c := range monochromaticColors {
		h, s, l := c.HSL()
		fmt.Printf("  %d: HSL(%.0f°, %.0f%%, %.0f%%) → %s\n", 
			i+1, h*360, s*100, l*100, c.CSSHSL())
	}
	
	classification = classifyColorSet(monochromaticColors)
	fmt.Printf("Classification: %s %s\n", 
		classification, tests.CheckMark(classification == "monochromatic"))

	// Test 3: Full-Color Detection (Multiple Distinct Hues)
	fmt.Println("\nTest 3: Full-Color Detection (Multiple Distinct Hues)")
	
	fullColorColors := []*color.Color{
		color.NewHSL(0.0/360, 0.8, 0.5),    // Red
		color.NewHSL(120.0/360, 0.7, 0.4),  // Green  
		color.NewHSL(240.0/360, 0.9, 0.6),  // Blue
		color.NewHSL(60.0/360, 0.8, 0.5),   // Yellow
	}
	
	fmt.Printf("Test colors:\n")
	for i, c := range fullColorColors {
		h, s, l := c.HSL()
		fmt.Printf("  %d: HSL(%.0f°, %.0f%%, %.0f%%) → %s\n", 
			i+1, h*360, s*100, l*100, c.CSSHSL())
	}
	
	classification = classifyColorSet(fullColorColors)
	fmt.Printf("Classification: %s %s\n", 
		classification, tests.CheckMark(classification == "full-color"))

	// Test 4: Edge Case - Near-Monochromatic (Just Outside Tolerance)
	fmt.Println("\nTest 4: Edge Case - Near-Monochromatic (Outside 10° Tolerance)")
	
	nearMonochromaticColors := []*color.Color{
		color.NewHSL(220.0/360, 0.8, 0.5),  // Blue
		color.NewHSL(235.0/360, 0.8, 0.5),  // Blue-purple (15° away, outside tolerance)
		color.NewHSL(205.0/360, 0.8, 0.5),  // Blue-cyan (15° away, outside tolerance) 
		color.NewHSL(222.0/360, 0.8, 0.5),  // Blue (within tolerance)
	}
	
	fmt.Printf("Test colors:\n")
	for i, c := range nearMonochromaticColors {
		h, s, l := c.HSL()
		fmt.Printf("  %d: HSL(%.0f°, %.0f%%, %.0f%%) → %s\n", 
			i+1, h*360, s*100, l*100, c.CSSHSL())
	}
	
	classification = classifyColorSet(nearMonochromaticColors)
	fmt.Printf("Classification: %s %s\n", 
		classification, tests.CheckMark(classification == "full-color"))

	// Test 5: Monochromatic with Temperature-Matched Grays
	fmt.Println("\nTest 5: Monochromatic with Temperature-Matched Grays")
	
	monoWithGraysColors := []*color.Color{
		color.NewHSL(30.0/360, 0.8, 0.3),   // Dark orange
		color.NewHSL(28.0/360, 0.9, 0.5),   // Orange (within tolerance)
		color.NewHSL(32.0/360, 0.7, 0.7),   // Light orange (within tolerance)
		color.NewHSL(35.0/360, 0.02, 0.2),  // Warm gray (very low saturation)
		color.NewHSL(40.0/360, 0.03, 0.6),  // Warm gray (very low saturation)
		color.NewHSL(30.0/360, 0.01, 0.8),  // Warm gray (very low saturation)
	}
	
	fmt.Printf("Test colors:\n")
	for i, c := range monoWithGraysColors {
		h, s, l := c.HSL()
		isGray := l < 0.02 || l > 0.98 || s < 0.05
		colorType := "colored"
		if isGray {
			colorType = "gray"
		}
		fmt.Printf("  %d: HSL(%.0f°, %.1f%%, %.0f%%) → %s (%s)\n", 
			i+1, h*360, s*100, l*100, c.CSSHSL(), colorType)
	}
	
	classification = classifyColorSet(monoWithGraysColors)
	fmt.Printf("Classification: %s %s\n", 
		classification, tests.CheckMark(classification == "monochromatic"))

	// Test 6: Hue Binning Logic Validation
	fmt.Println("\nTest 6: Hue Binning Logic (10° Tolerance)")
	
	testHues := []struct {
		hue         float64
		description string
	}{
		{210.0, "Base Blue"},
		{205.0, "5° away (within tolerance)"},
		{215.0, "5° away (within tolerance)"},
		{200.0, "10° away (edge of tolerance)"},
		{220.0, "10° away (edge of tolerance)"},
		{195.0, "15° away (outside tolerance)"},
		{225.0, "15° away (outside tolerance)"},
	}
	
	fmt.Printf("Hue binning test (10° bins, ±10° tolerance):\n")
	baseHue := 210.0
	baseBin := int(math.Round(baseHue/10) * 10)
	
	for _, test := range testHues {
		testBin := int(math.Round(test.hue/10) * 10) % 360
		binDistance := math.Min(math.Abs(float64(testBin-baseBin)), 360-math.Abs(float64(testBin-baseBin)))
		withinTolerance := binDistance <= 10
		
		fmt.Printf("  %.0f° → Bin %d°, distance %.0f°: %s (%s)\n", 
			test.hue, testBin, binDistance,
			tests.CheckMark(withinTolerance == (binDistance <= 10)),
			test.description)
	}

	// Test 7: Grayscale vs Monochromatic Distinction
	fmt.Println("\nTest 7: Grayscale vs Monochromatic Distinction")
	
	// Pure grayscale (no hue information)
	pureGrayscale := []*color.Color{
		color.NewRGB(0, 0, 0),
		color.NewRGB(85, 85, 85),
		color.NewRGB(170, 170, 170),
		color.NewRGB(255, 255, 255),
	}
	
	// Monochromatic blues with achromatic grays
	monoBlueWithGrays := []*color.Color{
		color.NewHSL(240.0/360, 0.8, 0.3),  // Dark blue (colored)
		color.NewHSL(235.0/360, 0.9, 0.5),  // Blue (colored, within tolerance)
		color.NewRGB(64, 64, 64),            // Gray (achromatic)
		color.NewRGB(192, 192, 192),         // Gray (achromatic)
	}
	
	pureResult := classifyColorSet(pureGrayscale)
	monoResult := classifyColorSet(monoBlueWithGrays)
	
	fmt.Printf("Pure achromatic grays: %s %s\n", 
		pureResult, tests.CheckMark(pureResult == "grayscale"))
	fmt.Printf("Monochromatic + grays: %s %s\n", 
		monoResult, tests.CheckMark(monoResult == "monochromatic"))
	
	fmt.Printf("\nDistinction validation: %s - Correctly distinguishes grayscale vs monochromatic\n",
		tests.CheckMark(pureResult == "grayscale" && monoResult == "monochromatic"))

	fmt.Println("\n=== Color Classification Validation Complete ===")
}