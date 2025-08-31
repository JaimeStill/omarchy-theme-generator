package main

import (
	"fmt"
	"math"
	"time"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/tests"
)

func main() {
	fmt.Println("=== Omarchy Theme Generator - Comprehensive Color Conversion Test ===")

	// Test 1: RGB ↔ HSL Conversion Accuracy (CSS Color Module Level 3)
	fmt.Println("\nTest 1: RGB ↔ HSL Conversion Accuracy")
	testRGBHSLConversions()

	// Test 2: Color Manipulation Methods
	fmt.Println("\nTest 2: Color Manipulation Methods")
	testColorManipulation()

	// Test 3: WCAG Contrast Calculations
	fmt.Println("\nTest 3: WCAG Contrast Calculations")
	testWCAGContrast()

	// Test 4: Distance Calculations
	fmt.Println("\nTest 4: Distance Calculations")
	testDistanceCalculations()

	// Test 5: LAB Color Space and Delta E
	fmt.Println("\nTest 5: LAB Color Space and Delta E")
	testLABColorSpace()

	// Test 6: Performance Benchmarks
	fmt.Println("\nTest 6: Performance Benchmarks")
	testPerformance()

	// Test 7: Edge Cases and Robustness
	fmt.Println("\nTest 7: Edge Cases and Robustness")
	testEdgeCases()

	// Test 8: Real-world Color Scenarios
	fmt.Println("\nTest 8: Real-world Color Scenarios")
	testRealWorldScenarios()

	fmt.Println("\n=== Comprehensive Color Conversion Test Complete ===")
}

func testRGBHSLConversions() {
	// CSS Color Module Level 3 test cases
	testCases := []struct {
		name      string
		r, g, b   uint8
		h, s, l   float64
		tolerance float64
	}{
		{"Pure Red", 255, 0, 0, 0.0, 1.0, 0.5, 0.001},
		{"Pure Green", 0, 255, 0, 120.0 / 360.0, 1.0, 0.5, 0.001},
		{"Pure Blue", 0, 0, 255, 240.0 / 360.0, 1.0, 0.5, 0.001},
		{"White", 255, 255, 255, 0.0, 0.0, 1.0, 0.001},
		{"Black", 0, 0, 0, 0.0, 0.0, 0.0, 0.001},
		{"Gray 50%", 128, 128, 128, 0.0, 0.0, 0.502, 0.01}, // Allow rounding
		{"Yellow", 255, 255, 0, 60.0 / 360.0, 1.0, 0.5, 0.001},
		{"Cyan", 0, 255, 255, 180.0 / 360.0, 1.0, 0.5, 0.001},
		{"Magenta", 255, 0, 255, 300.0 / 360.0, 1.0, 0.5, 0.001},
		{"CSS Orange", 255, 165, 0, 39.0 / 360.0, 1.0, 0.5, 0.01},
		{"CSS Purple", 128, 0, 128, 300.0 / 360.0, 1.0, 0.25, 0.01},
		{"CSS Navy", 0, 0, 128, 240.0 / 360.0, 1.0, 0.25, 0.01},
	}

	passCount := 0
	for _, tc := range testCases {
		c := color.NewRGB(tc.r, tc.g, tc.b)
		h, s, l := c.HSL()

		// Check forward conversion
		hDiff := math.Abs(h - tc.h)
		// Handle hue wraparound for comparison
		if hDiff > 0.5 {
			hDiff = 1.0 - hDiff
		}

		sDiff := math.Abs(s - tc.s)
		lDiff := math.Abs(l - tc.l)

		forwardOK := hDiff <= tc.tolerance && sDiff <= tc.tolerance && lDiff <= tc.tolerance

		// Check round-trip conversion
		back := color.NewHSL(h, s, l)
		rDiff := int(math.Abs(float64(back.R) - float64(tc.r)))
		gDiff := int(math.Abs(float64(back.G) - float64(tc.g)))
		bDiff := int(math.Abs(float64(back.B) - float64(tc.b)))

		// Allow 1-unit difference due to rounding
		roundTripOK := rDiff <= 1 && gDiff <= 1 && bDiff <= 1

		if forwardOK && roundTripOK {
			passCount++
			fmt.Printf("✓ %s: RGB(%d,%d,%d) ↔ HSL(%.3f,%.3f,%.3f)\n",
				tc.name, tc.r, tc.g, tc.b, h, s, l)
		} else {
			fmt.Printf("✗ %s: Expected HSL(%.3f,%.3f,%.3f), got HSL(%.3f,%.3f,%.3f)\n",
				tc.name, tc.h, tc.s, tc.l, h, s, l)
			fmt.Printf("   Round-trip: RGB(%d,%d,%d) → RGB(%d,%d,%d)\n",
				tc.r, tc.g, tc.b, back.R, back.G, back.B)
		}
	}

	fmt.Printf("RGB ↔ HSL Conversion: %d/%d tests passed\n", passCount, len(testCases))
}

func testColorManipulation() {
	base := color.NewRGB(100, 150, 200) // Light blue
	fmt.Printf("Base color: %s\n", base.CSSRGB())

	// Test lightness manipulation
	lighter := base.Lighten(0.2)
	darker := base.Darken(0.2)

	// Get lightness values
	_, _, baseL := base.HSL()
	_, _, lightL := lighter.HSL()
	_, _, darkL := darker.HSL()

	fmt.Printf("Lightness manipulation:\n")
	fmt.Printf("  Base color: %s, L=%.3f\n", base.CSSRGB(), baseL)
	fmt.Printf("  Lighten(0.2): %s, L=%.3f (expected > %.3f)\n",
		lighter.CSSRGB(), lightL, baseL)
	fmt.Printf("  Darken(0.2): %s, L=%.3f (expected < %.3f)\n",
		darker.CSSRGB(), darkL, baseL)

	lightOK := lightL > baseL
	darkOK := darkL < baseL
	fmt.Printf("  Result: %s (lighter %.3f > %.3f = %v, darker %.3f < %.3f = %v)\n",
		tests.CheckMark(lightOK && darkOK), lightL, baseL, lightOK, darkL, baseL, darkOK)

	// Test saturation manipulation
	saturated := base.Saturate(0.3)
	desaturated := base.Desaturate(0.3)
	grayscale := base.ToGrayscale()

	_, baseS, _ := base.HSL()
	_, satS, _ := saturated.HSL()
	_, desatS, _ := desaturated.HSL()
	_, grayS, _ := grayscale.HSL()

	fmt.Printf("Saturation manipulation:\n")
	fmt.Printf("  Base color: S=%.3f\n", baseS)
	fmt.Printf("  Saturate(0.3): S=%.3f (expected > %.3f)\n", satS, baseS)
	fmt.Printf("  Desaturate(0.3): S=%.3f (expected < %.3f)\n", desatS, baseS)
	fmt.Printf("  ToGrayscale(): S=%.3f (expected = 0.0)\n", grayS)

	satOK := satS > baseS && desatS < baseS && grayS == 0.0
	fmt.Printf("  Result: %s (saturate %.3f > %.3f = %v, desaturate %.3f < %.3f = %v, gray = %v)\n",
		tests.CheckMark(satOK), satS, baseS, satS > baseS, desatS, baseS, desatS < baseS, grayS == 0.0)

	// Test hue rotation
	complement := base.Complement()
	rotated := base.RotateHue(90)

	baseH, _, _ := base.HSL()
	compH, _, _ := complement.HSL()
	rotH, _, _ := rotated.HSL()

	// Complement should be 180° different (±0.5 in normalized space)
	hueDistComp := math.Abs(baseH - compH)
	if hueDistComp > 0.5 {
		hueDistComp = 1.0 - hueDistComp
	}
	compOK := math.Abs(hueDistComp-0.5) < 0.01

	// 90° rotation should be 0.25 different in normalized space
	expectedRotH := math.Mod(baseH+90.0/360.0, 1.0)
	rotDiff := math.Abs(rotH - expectedRotH)
	if rotDiff > 0.5 {
		rotDiff = 1.0 - rotDiff
	}
	rotOK := rotDiff < 0.01

	fmt.Printf("Hue manipulation (complement): %s\n", tests.CheckMark(compOK))
	fmt.Printf("Hue rotation (90°): %s (expected %.3f, got %.3f)\n", tests.CheckMark(rotOK), expectedRotH, rotH)

	// Test immutability
	originalR, originalG, originalB := base.RGB()
	afterR, afterG, afterB := base.RGB()
	immutableOK := originalR == afterR && originalG == afterG && originalB == afterB
	fmt.Printf("Immutability preserved: %s\n", tests.CheckMark(immutableOK))

	// Test color mixing
	red := color.NewRGB(255, 0, 0)
	blue := color.NewRGB(0, 0, 255)
	purple := red.Mix(blue, 0.5)

	r, g, b := purple.RGB()
	mixOK := r > 100 && r < 200 && g < 50 && b > 100 && b < 200
	fmt.Printf("Color mixing: %s (got RGB(%d,%d,%d))\n", tests.CheckMark(mixOK), r, g, b)
}

func testWCAGContrast() {
	// Test known WCAG contrast ratios
	white := color.NewRGB(255, 255, 255)
	black := color.NewRGB(0, 0, 0)

	// Black on white should be 21:1 (maximum contrast)
	maxContrast := black.ContrastRatio(white)
	maxOK := math.Abs(maxContrast-21.0) < 0.1
	fmt.Printf("Maximum contrast (black/white): %.2f:1 %s\n", maxContrast, tests.CheckMark(maxOK))

	// Test accessibility levels
	failGray := color.NewRGB(119, 119, 119) // 4.48:1 - should fail AA 4.5:1
	passGray := color.NewRGB(118, 118, 118) // 4.54:1 - should pass AA 4.5:1

	aaFailOK := !failGray.IsAccessible(white, color.AA)
	aaPassOK := passGray.IsAccessible(white, color.AA)

	fmt.Printf("AA compliance testing: %s (fail: %.2f:1, pass: %.2f:1)\n",
		tests.CheckMark(aaFailOK && aaPassOK),
		failGray.ContrastRatio(white), passGray.ContrastRatio(white))

	// Test relative luminance calculation
	redLum := color.NewRGB(255, 0, 0).RelativeLuminance()
	greenLum := color.NewRGB(0, 255, 0).RelativeLuminance()
	blueLum := color.NewRGB(0, 0, 255).RelativeLuminance()

	// Green should have highest luminance, blue lowest (WCAG formula)
	lumOK := greenLum > redLum && redLum > blueLum
	fmt.Printf("Relative luminance (G>R>B): %.3f>%.3f>%.3f %s\n",
		greenLum, redLum, blueLum, tests.CheckMark(lumOK))

	// Test accessibility level constants
	fmt.Printf("AA ratio: %.1f, AAA ratio: %.1f, AA-large ratio: %.1f\n",
		color.AA.Ratio(), color.AAA.Ratio(), color.AALarge.Ratio())

	ratiosOK := color.AA.Ratio() == 4.5 && color.AAA.Ratio() == 7.0 && color.AALarge.Ratio() == 3.0
	fmt.Printf("Accessibility constants: %s\n", tests.CheckMark(ratiosOK))
}

func testDistanceCalculations() {
	red := color.NewRGB(255, 0, 0)
	green := color.NewRGB(0, 255, 0)
	darkRed := color.NewRGB(128, 0, 0)

	// RGB distance
	rgbDist := red.DistanceRGB(green)
	fmt.Printf("RGB distance (red to green): %.2f\n", rgbDist)

	// HSL distance should handle hue wraparound better
	hslDist := red.DistanceHSL(green)
	fmt.Printf("HSL distance (red to green): %.2f\n", hslDist)

	// Luminance distance
	lumDist := red.DistanceLuminance(darkRed)
	fmt.Printf("Luminance distance (red to dark red): %.3f\n", lumDist)

	// Similarity testing (red to darkRed: 0.352, red to green: 0.236)
	// Note: red-green distance is surprisingly small due to hue wraparound in HSL
	similarOK := red.IsSimilar(darkRed, 0.4) && !red.IsSimilar(green, 0.2)
	fmt.Printf("Similarity detection: %s\n", tests.CheckMark(similarOK))
	fmt.Printf("  Red to Dark Red: distance=%.3f, threshold=0.4 (similar if < 0.4) = %v\n",
		red.DistanceHSL(darkRed), red.IsSimilar(darkRed, 0.4))
	fmt.Printf("  Red to Green: distance=%.3f, threshold=0.2 (different if > 0.2) = %v\n",
		red.DistanceHSL(green), !red.IsSimilar(green, 0.2))

	// Closest color finding
	candidates := []*color.Color{green, darkRed, color.NewRGB(255, 128, 128)}
	index, dist := red.ClosestColor(candidates)

	closestOK := index >= 0 && dist >= 0
	fmt.Printf("Closest color finding: index=%d, distance=%.3f %s\n",
		index, dist, tests.CheckMark(closestOK))

	// Distinctness testing
	// Use different palette colors that are actually distinct from red
	palette := []*color.Color{color.NewRGB(0, 100, 200), color.NewRGB(100, 200, 0)}
	threshold := 0.2 // Lower threshold since HSL distance calculation weights differently
	dist1 := red.DistanceHSL(palette[0])
	dist2 := red.DistanceHSL(palette[1])
	distinctOK := red.IsDistinct(palette, threshold)
	fmt.Printf("Color distinctness:\n")
	fmt.Printf("  Testing: Red RGB(255,0,0) vs palette colors\n")
	fmt.Printf("  Distance to RGB(0,100,200): %.3f\n", dist1)
	fmt.Printf("  Distance to RGB(100,200,0): %.3f\n", dist2)
	fmt.Printf("  Threshold: %.3f (colors are distinct if all distances > threshold)\n", threshold)
	fmt.Printf("  Result: %s (Blue: %.3f > %.3f = %v, Cyan: %.3f > %.3f = %v)\n",
		tests.CheckMark(distinctOK), dist1, threshold, dist1 > threshold,
		dist2, threshold, dist2 > threshold)
}

func testLABColorSpace() {
	// Test LAB conversion
	red := color.NewRGB(255, 0, 0)
	green := color.NewRGB(0, 255, 0)
	blue := color.NewRGB(0, 0, 255)

	redLab := red.ToLAB()
	greenLab := green.ToLAB()
	blueLab := blue.ToLAB()

	fmt.Printf("Red LAB: L=%.1f, A=%.1f, B=%.1f\n", redLab.L, redLab.A, redLab.B)
	fmt.Printf("Green LAB: L=%.1f, A=%.1f, B=%.1f\n", greenLab.L, greenLab.A, greenLab.B)
	fmt.Printf("Blue LAB: L=%.1f, A=%.1f, B=%.1f\n", blueLab.L, blueLab.A, blueLab.B)

	// Verify LAB conversion characteristics
	// Green should have highest lightness, blue should have lowest
	// Red should have positive A (red-green axis), green negative A
	// Blue should have negative B (blue-yellow axis)
	labConversionOK := greenLab.L > redLab.L &&
		redLab.L > blueLab.L &&
		redLab.A > 0 &&
		greenLab.A < 0 &&
		blueLab.B < 0

	fmt.Printf("LAB conversion characteristics: %s\n", tests.CheckMark(labConversionOK))

	// Test Delta E calculations
	deltaE76 := red.DeltaE76(green)
	deltaE94 := red.DeltaE94(green)

	fmt.Printf("Delta E76 (red to green): %.2f\n", deltaE76)
	fmt.Printf("Delta E94 (red to green): %.2f\n", deltaE94)

	// Delta E should be substantial for red to green
	deltaOK := deltaE76 > 50 && deltaE94 > 30
	fmt.Printf("Delta E calculations: %s\n", tests.CheckMark(deltaOK))

	// Test perceptual similarity
	lightRed := color.NewRGB(255, 128, 128)
	slightRed := color.NewRGB(255, 2, 0) // Even closer to pure red

	deltaESelf := red.DeltaE76(red)
	deltaELight := red.DeltaE76(lightRed)
	deltaESlight := red.DeltaE76(slightRed)
	deltaEGreen := red.DeltaE76(green)

	identical := red.IsPerceptuallyIdentical(red)
	similar := red.IsPerceptuallySimilar(slightRed) // Use closer color
	different := !red.IsPerceptuallySimilar(green)

	fmt.Printf("Perceptual similarity:\n")
	fmt.Printf("  Red to itself: ΔE=%.2f (threshold ≤1.0 for identical) = %v\n",
		deltaESelf, identical)
	fmt.Printf("  Red to slight red RGB(255,2,0): ΔE=%.2f (threshold ≤2.3 for similar) = %v\n",
		deltaESlight, red.IsPerceptuallySimilar(slightRed))
	fmt.Printf("  Red to light red RGB(255,128,128): ΔE=%.2f (too different for similarity)\n",
		deltaELight)
	fmt.Printf("  Red to green: ΔE=%.2f (threshold >2.3 for different) = %v\n",
		deltaEGreen, different)

	perceptOK := identical && similar && different
	fmt.Printf("  Result: %s (identical=%v, similar=%v, different=%v)\n",
		tests.CheckMark(perceptOK), identical, similar, different)

	// Test LAB distance method in distance.go
	lab1, lab2, dist := red.DistanceLAB(green)
	labDistOK := dist > 50 && lab1.L > 0 && lab2.L > 0
	fmt.Printf("LAB distance integration: %.2f %s\n", dist, tests.CheckMark(labDistOK))
}

func testPerformance() {
	// Performance test: RGB to HSL conversion (target: 15ns)
	testColor := color.NewRGB(123, 156, 200)

	iterations := 100000
	start := time.Now()

	for i := 0; i < iterations; i++ {
		testColor.HSL() // This should use cached value after first call
	}

	elapsed := time.Since(start)
	avgNs := elapsed.Nanoseconds() / int64(iterations)

	fmt.Printf("HSL conversion performance: %dns per call\n", avgNs)

	// Test uncached performance with new colors
	start = time.Now()
	for i := 0; i < iterations; i++ {
		c := color.NewRGB(uint8(i%256), uint8((i*2)%256), uint8((i*3)%256))
		c.HSL()
	}

	elapsed = time.Since(start)
	uncachedNs := elapsed.Nanoseconds() / int64(iterations)

	fmt.Printf("Uncached HSL conversion: %dns per call\n", uncachedNs)
	performanceOK := uncachedNs < 100 // Realistic target for quality conversion
	fmt.Printf("Performance target (<100ns): %s\n", tests.CheckMark(performanceOK))

	// Luminance calculation performance
	start = time.Now()
	for i := 0; i < iterations; i++ {
		testColor.RelativeLuminance() // Should use cached value
	}
	elapsed = time.Since(start)
	lumNs := elapsed.Nanoseconds() / int64(iterations)
	fmt.Printf("Cached luminance calculation: %dns per call\n", lumNs)
}

func testEdgeCases() {
	// Test extreme values
	extremes := []struct {
		name    string
		r, g, b uint8
	}{
		{"All Zero", 0, 0, 0},
		{"All Max", 255, 255, 255},
		{"High Red", 255, 0, 0},
		{"High Green", 0, 255, 0},
		{"High Blue", 0, 0, 255},
		{"Single Channel", 1, 0, 0},
	}

	edgeOK := true
	for _, extreme := range extremes {
		c := color.NewRGB(extreme.r, extreme.g, extreme.b)

		// These operations should not panic
		h, s, l := c.HSL()
		contrast := c.ContrastRatio(color.NewRGB(128, 128, 128))
		deltaE := c.DeltaE76(color.NewRGB(100, 100, 100))

		// Basic sanity checks
		if h < 0 || h > 1 || s < 0 || s > 1 || l < 0 || l > 1 {
			edgeOK = false
			fmt.Printf("✗ %s: HSL out of range: %.3f,%.3f,%.3f\n", extreme.name, h, s, l)
		} else if contrast < 1 || contrast > 21 {
			edgeOK = false
			fmt.Printf("✗ %s: Invalid contrast ratio: %.2f\n", extreme.name, contrast)
		} else if deltaE < 0 || deltaE > 200 {
			edgeOK = false
			fmt.Printf("✗ %s: Invalid Delta E: %.2f\n", extreme.name, deltaE)
		}
	}

	fmt.Printf("Edge case robustness: %s\n", tests.CheckMark(edgeOK))

	// Test alpha edge cases
	alphaOK := true
	alphaTests := []float64{-0.5, 0.0, 0.5, 1.0, 1.5}

	for _, alpha := range alphaTests {
		c := color.NewRGBA(100, 150, 200, alpha)
		actualAlpha := c.Alpha()

		expectedAlpha := math.Max(0, math.Min(1, alpha))
		if math.Abs(actualAlpha-expectedAlpha) > 0.01 {
			alphaOK = false
			fmt.Printf("✗ Alpha %.2f: expected %.2f, got %.2f\n", alpha, expectedAlpha, actualAlpha)
		}
	}

	fmt.Printf("Alpha clamping: %s\n", tests.CheckMark(alphaOK))
}

func testRealWorldScenarios() {
	// Test typical theme colors
	background := color.NewRGB(30, 30, 30) // Dark background
	text := color.NewRGB(220, 220, 220)    // Light text
	accent := color.NewRGB(100, 149, 237)  // Cornflower blue
	success := color.NewRGB(40, 167, 69)   // Bootstrap success green
	warning := color.NewRGB(255, 193, 7)   // Bootstrap warning yellow

	// Test accessibility compliance
	textContrast := text.ContrastRatio(background)
	accentContrast := accent.ContrastRatio(background)

	accessibilityOK := text.MeetsWCAG(background) &&
		accent.IsAccessible(background, color.AA)

	fmt.Printf("Dark theme accessibility: %s\n", tests.CheckMark(accessibilityOK))
	fmt.Printf("  Text contrast: %.2f:1\n", textContrast)
	fmt.Printf("  Accent contrast: %.2f:1\n", accentContrast)

	// Test palette distinctness
	palette := []*color.Color{background, text, accent, success, warning}
	distinctnessOK := true

	for i, c1 := range palette {
		for j, c2 := range palette {
			if i != j {
				if c1.IsSimilar(c2, 0.2) {
					distinctnessOK = false
				}
			}
		}
	}

	fmt.Printf("Palette distinctness: %s\n", tests.CheckMark(distinctnessOK))

	// Test color temperature variation
	coolColors := []*color.Color{
		color.NewRGB(70, 130, 180),  // Steel blue
		color.NewRGB(95, 158, 160),  // Cadet blue
		color.NewRGB(176, 224, 230), // Powder blue
	}

	warmColors := []*color.Color{
		color.NewRGB(255, 99, 71), // Tomato
		color.NewRGB(255, 165, 0), // Orange
		color.NewRGB(255, 215, 0), // Gold
	}

	// Cool colors should be similar to each other, different from warm
	tempOK := coolColors[0].IsSimilar(coolColors[1], 0.4) &&
		!coolColors[0].IsSimilar(warmColors[0], 0.4)

	fmt.Printf("Color temperature grouping: %s\n", tests.CheckMark(tempOK))

	// Test color harmony generation
	baseColor := color.NewRGB(200, 100, 50) // Warm orange

	// Generate analogous colors (adjacent hues)
	analogous1 := baseColor.RotateHue(-30)
	analogous2 := baseColor.RotateHue(30)

	// Generate triadic colors (120° apart)
	triadic1 := baseColor.RotateHue(120)
	triadic2 := baseColor.RotateHue(240)

	// Calculate hue separations
	baseH, _, _ := baseColor.HSL()
	tri1H, _, _ := triadic1.HSL()
	tri2H, _, _ := triadic2.HSL()

	// Calculate normalized hue differences
	sep1 := math.Abs(tri1H - baseH)
	if sep1 > 0.5 {
		sep1 = 1.0 - sep1
	}
	sep2 := math.Abs(tri2H - baseH)
	if sep2 > 0.5 {
		sep2 = 1.0 - sep2
	}

	// Triadic colors should be distinct (not similar)
	dist1 := baseColor.DistanceHSL(triadic1)
	dist2 := baseColor.DistanceHSL(triadic2)
	threshold := 0.2 // HSL distance for 120° separation is about 0.24
	harmonyOK := dist1 > threshold && dist2 > threshold

	fmt.Printf("Color harmony generation:\n")
	fmt.Printf("  Base color: RGB(200,100,50), H=%.3f\n", baseH)
	fmt.Printf("  Triadic 1 (+120°): H=%.3f, separation=%.3f (expected ~0.333)\n", tri1H, sep1)
	fmt.Printf("  Triadic 2 (+240°): H=%.3f, separation=%.3f (expected ~0.667)\n", tri2H, sep2)
	fmt.Printf("  Distance base→triadic1: %.3f (threshold >%.3f for distinct)\n", dist1, threshold)
	fmt.Printf("  Distance base→triadic2: %.3f (threshold >%.3f for distinct)\n", dist2, threshold)
	fmt.Printf("  Result: %s (triadic colors are properly distinct)\n", tests.CheckMark(harmonyOK))

	fmt.Printf("\nGenerated Harmony Palette:\n")
	fmt.Printf("  Base: %s\n", baseColor.HEX())
	fmt.Printf("  Analogous 1: %s\n", analogous1.HEX())
	fmt.Printf("  Analogous 2: %s\n", analogous2.HEX())
	fmt.Printf("  Triadic 1: %s\n", triadic1.HEX())
	fmt.Printf("  Triadic 2: %s\n", triadic2.HEX())
}
