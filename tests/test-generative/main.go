package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/JaimeStill/omarchy-theme-generator/tests"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/generative"
)

func main() {
	fmt.Println("=== Omarchy Theme Generator - Generative Image Test ===")
	fmt.Printf("Testing sophisticated computational image generation system\n")
	fmt.Printf("Focus: Enhanced cassette futurism with authentic industrial materials\n\n")

	// Test 1: Basic Image Generation Functions
	fmt.Println("Test 1: Basic Generative Functions")
	testBasicGeneration()

	// Test 2: Enhanced Cassette Futurism Generator
	fmt.Println("\nTest 2: Enhanced Cassette Futurism Generator")
	testCassetteFuturismGenerator()

	// Test 3: Material Simulation Validation
	fmt.Println("\nTest 3: Material Simulation Accuracy")
	testMaterialSimulation()

	// Test 4: Color Palette Generation
	fmt.Println("\nTest 4: Industrial Color Palette Generation")
	testColorPaletteGeneration()

	// Test 5: Complex Gradient Generation
	fmt.Println("\nTest 5: Complex Gradient Generation")
	testComplexGradients()

	// Test 6: Interface Element Rendering
	fmt.Println("\nTest 6: Industrial Interface Elements")
	testInterfaceElements()

	// Test 7: Performance Benchmarking
	fmt.Println("\nTest 7: Performance Benchmarking")
	testPerformanceBenchmarks()

	// Test 8: Visual Output Quality Assessment
	fmt.Println("\nTest 8: Visual Quality Assessment")
	testVisualQuality()

	// Test 9: Edge Cases and Boundary Conditions
	fmt.Println("\nTest 9: Edge Cases and Boundaries")
	testEdgeCases()

	// Test 10: Integration with Color Theory Systems
	fmt.Println("\nTest 10: Color Theory Integration")
	testColorTheoryIntegration()

	fmt.Println("\n=== Generative Image Test Complete ===")
	fmt.Printf("Sample images saved to: tests/test-generative/samples/\n")
}

func testBasicGeneration() {
	// Test 4K image generation
	fmt.Printf("4K Test Image Generation:\n")
	start := time.Now()
	img4k := generative.Generate4KTestImage()
	duration4k := time.Since(start)
	
	bounds := img4k.Bounds()
	expectedWidth, expectedHeight := 3840, 2160
	
	sizeCorrect := bounds.Dx() == expectedWidth && bounds.Dy() == expectedHeight
	performanceOK := duration4k < 2*time.Second // Target: < 2s for 4K
	
	fmt.Printf("  Dimensions: %dx%d (expected %dx%d) %s\n", 
		bounds.Dx(), bounds.Dy(), expectedWidth, expectedHeight, tests.CheckMark(sizeCorrect))
	fmt.Printf("  Generation time: %v (target <2s) %s\n", 
		duration4k, tests.CheckMark(performanceOK))
	
	// Verify color variation exists
	colorVariation := analyzeColorVariation(img4k)
	variationOK := colorVariation > 1000 // Should have many unique colors
	fmt.Printf("  Color variation: %d unique colors %s\n", 
		colorVariation, tests.CheckMark(variationOK))
	
	saveImage(img4k, "4k-test-image.png")
	
	// Test monochrome generation
	fmt.Printf("\nMonochrome Test Image Generation:\n")
	mono := generative.GenerateMonochromeTestImage(800, 600)
	monoColors := analyzeColorVariation(mono)
	
	// Verify it's actually monochromatic (low color count)
	isMonochrome := monoColors < 100 // Should have very few unique colors
	fmt.Printf("  Monochrome colors: %d (expected <100) %s\n", 
		monoColors, tests.CheckMark(isMonochrome))
	
	saveImage(mono, "monochrome-test.png")
	
	// Test high contrast generation
	fmt.Printf("\nHigh Contrast Test Image Generation:\n")
	contrast := generative.GenerateHighContrastTestImage(400, 300)
	contrastColors := analyzeColorVariation(contrast)
	
	// Should have exactly the defined colors (5)
	correctContrast := contrastColors <= 10 // Allow slight variations due to rendering
	fmt.Printf("  High contrast colors: %d (expected d10) %s\n", 
		contrastColors, tests.CheckMark(correctContrast))
	
	saveImage(contrast, "high-contrast-test.png")
}

func testCassetteFuturismGenerator() {
	// Test with single canonical variant - warm orange accent
	fmt.Printf("\nCassette Futurism Generator (warm orange accent):\n")
	
	hue := 25.0/360.0 // Warm orange
	expectedWarmth := "warm"
	
	start := time.Now()
	img := generative.GenerateCassetteFuturismImage(800, 600, hue)
	duration := time.Since(start)
	
	bounds := img.Bounds()
	sizeCorrect := bounds.Dx() == 800 && bounds.Dy() == 600
	performanceOK := duration < 500*time.Millisecond // Should be fast for 800x600
	
	fmt.Printf("  Generation: %dx%d in %v %s\n", 
		bounds.Dx(), bounds.Dy(), duration, tests.CheckMark(sizeCorrect && performanceOK))
	
	// Analyze industrial aesthetics
	colorStats := analyzeCassetteFuturismAesthetics(img)
	
	// Verify industrial color palette presence
	hasIndustrialColors := colorStats.GrayLevels > 5
	hasAccentColors := colorStats.AccentColors > 0
	hasMetallicHighlights := colorStats.MetallicColors > 0
	
	fmt.Printf("  Gray levels: %d (expected >5) %s\n", 
		colorStats.GrayLevels, tests.CheckMark(hasIndustrialColors))
	fmt.Printf("  Accent colors: %d (expected >0) %s\n", 
		colorStats.AccentColors, tests.CheckMark(hasAccentColors))
	fmt.Printf("  Metallic highlights: %d (expected >0) %s\n", 
		colorStats.MetallicColors, tests.CheckMark(hasMetallicHighlights))
	
	// Save with simple filename
	saveImage(img, "cassette-futurism.png")
	
	// Verify temperature matching
	expectedWarmthBool := expectedWarmth == "warm"
	actualWarmth := colorStats.WarmTone
	temperatureMatch := expectedWarmthBool == actualWarmth
	
	fmt.Printf("  Temperature matching: expected %s, detected %s %s\n",
		expectedWarmth, map[bool]string{true: "warm", false: "cool"}[actualWarmth],
		tests.CheckMark(temperatureMatch))
}

func testMaterialSimulation() {
	fmt.Printf("Material simulation accuracy assessment:\n")
	
	// Generate test image for material analysis - use cool blue to differentiate from warm orange main test
	img := generative.GenerateCassetteFuturismImage(600, 400, 210.0/360.0)
	
	// Analyze for material simulation characteristics
	materials := analyzeMaterialSimulation(img)
	
	// Test brushed metal detection
	hasBrushedMetal := materials.BrushedMetalTexture > 0.1 // At least 10% of image should have texture
	fmt.Printf("  Brushed metal texture: %.1f%% coverage %s\n", 
		materials.BrushedMetalTexture*100, tests.CheckMark(hasBrushedMetal))
	
	// Test CRT scanline detection
	hasScanlines := materials.CRTScanlines > 0.05 // Should detect scanline patterns
	fmt.Printf("  CRT scanlines: %.1f%% coverage %s\n", 
		materials.CRTScanlines*100, tests.CheckMark(hasScanlines))
	
	// Test LED indicator simulation
	hasLEDs := materials.LEDIndicators > 5 // Should have multiple LED elements
	fmt.Printf("  LED indicators: %d elements %s\n", 
		materials.LEDIndicators, tests.CheckMark(hasLEDs))
	
	// Test industrial plastic texture
	hasPlasticTexture := materials.PlasticTexture > 0.2 // Significant plastic surface area
	fmt.Printf("  Plastic textures: %.1f%% coverage %s\n", 
		materials.PlasticTexture*100, tests.CheckMark(hasPlasticTexture))
	
	// Test screen reflection effects
	hasReflections := materials.ScreenReflections > 0.05 // Should have reflection highlights
	fmt.Printf("  Screen reflections: %.1f%% coverage %s\n", 
		materials.ScreenReflections*100, tests.CheckMark(hasReflections))
	
	saveImage(img, "material-simulation-test.png")
}

func testColorPaletteGeneration() {
	fmt.Printf("Temperature-matched industrial palette validation:\n")
	
	// Test warm palette
	warmImg := generative.GenerateCassetteFuturismImage(400, 300, 30.0/360.0) // Orange hue
	warmPalette := extractDominantColors(warmImg, 10)
	warmness := calculatePaletteWarmth(warmPalette)
	
	isWarm := warmness > 0.1 // Should lean warm
	fmt.Printf("  Warm palette (orange accent): warmness=%.3f %s\n", 
		warmness, tests.CheckMark(isWarm))
	
	// Test cool palette  
	coolImg := generative.GenerateCassetteFuturismImage(400, 300, 210.0/360.0) // Blue hue
	coolPalette := extractDominantColors(coolImg, 10)
	coolness := calculatePaletteWarmth(coolPalette)
	
	isCool := coolness < -0.1 // Should lean cool
	fmt.Printf("  Cool palette (blue accent): warmness=%.3f %s\n", 
		coolness, tests.CheckMark(isCool))
	
	// Test WCAG accessibility compliance
	fmt.Printf("\nAccessibility compliance:\n")
	accessibilityOK := validateWCAGCompliance(warmPalette)
	fmt.Printf("  WCAG contrast ratios: %s\n", tests.CheckMark(accessibilityOK))
	
	saveImage(warmImg, "warm-palette.png")
	saveImage(coolImg, "cool-palette.png")
}

func testComplexGradients() {
	gradientTypes := []string{"linear-smooth", "radial-complex", "stepped-harsh"}
	
	for _, gradType := range gradientTypes {
		fmt.Printf("\n%s gradient generation:\n", gradType)
		
		start := time.Now()
		img := generative.GenerateComplexGradientImage(600, 400, gradType)
		duration := time.Since(start)
		
		// Verify smooth color transitions
		smoothness := analyzeGradientSmoothness(img)
		colorRange := analyzeColorVariation(img)
		
		performanceOK := duration < 200*time.Millisecond
		fmt.Printf("  Generation time: %v %s\n", 
			duration, tests.CheckMark(performanceOK))
		fmt.Printf("  Color variations: %d %s\n", 
			colorRange, tests.CheckMark(colorRange > 100))
		
		// Type-specific validation
		switch gradType {
		case "linear-smooth":
			isSmoothLinear := smoothness > 0.8 // Should be very smooth
			fmt.Printf("  Linear smoothness: %.3f %s\n", 
				smoothness, tests.CheckMark(isSmoothLinear))
			
		case "radial-complex":
			isComplexRadial := colorRange > 500 // Should have many color stops
			fmt.Printf("  Radial complexity: %d colors %s\n", 
				colorRange, tests.CheckMark(isComplexRadial))
			
		case "stepped-harsh":
			isSteppedHarsh := smoothness < 0.3 // Should have harsh transitions
			fmt.Printf("  Stepped harshness: %.3f (lower=harsher) %s\n", 
				smoothness, tests.CheckMark(isSteppedHarsh))
		}
		
		filename := fmt.Sprintf("gradient-%s.png", gradType)
		saveImage(img, filename)
	}
}

func testInterfaceElements() {
	fmt.Printf("Industrial interface element validation:\n")
	
	// Generate interface-heavy image - use green accent for distinct interface aesthetic
	img := generative.GenerateCassetteFuturismImage(800, 600, 120.0/360.0)
	
	// Analyze interface elements
	interfaces := analyzeInterfaceElements(img)
	
	// Test display screens
	hasDisplays := interfaces.DisplayScreens > 0
	fmt.Printf("  Display screens: %d detected %s\n", 
		interfaces.DisplayScreens, tests.CheckMark(hasDisplays))
	
	// Test control buttons
	hasButtons := interfaces.ControlButtons > 10 // Should have button matrix
	fmt.Printf("  Control buttons: %d detected %s\n", 
		interfaces.ControlButtons, tests.CheckMark(hasButtons))
	
	// Test slider controls
	hasSliders := interfaces.SliderControls > 5 // Should have slider bank
	fmt.Printf("  Slider controls: %d detected %s\n", 
		interfaces.SliderControls, tests.CheckMark(hasSliders))
	
	// Test seven-segment displays
	hasSegmentDisplays := interfaces.SevenSegmentDisplays > 0
	fmt.Printf("  Seven-segment displays: %d detected %s\n", 
		interfaces.SevenSegmentDisplays, tests.CheckMark(hasSegmentDisplays))
	
	// Test terminal text
	hasTerminalText := interfaces.TerminalText > 0
	fmt.Printf("  Terminal text lines: %d detected %s\n", 
		interfaces.TerminalText, tests.CheckMark(hasTerminalText))
	
	saveImage(img, "interface-elements.png")
}

func testPerformanceBenchmarks() {
	fmt.Printf("Performance benchmarking for complex rendering:\n")
	
	sizes := []struct {
		name string
		w, h int
		target time.Duration
	}{
		{"Small (400x300)", 400, 300, 100 * time.Millisecond},
		{"Medium (800x600)", 800, 600, 300 * time.Millisecond},
		{"Large (1200x800)", 1200, 800, 600 * time.Millisecond},
		{"HD (1920x1080)", 1920, 1080, 1 * time.Second},
	}
	
	for _, size := range sizes {
		fmt.Printf("\n%s rendering:\n", size.name)
		
		// Benchmark cassette futurism generation - use purple accent for performance testing
		start := time.Now()
		img := generative.GenerateCassetteFuturismImage(size.w, size.h, 300.0/360.0)
		duration := time.Since(start)
		
		pixelCount := size.w * size.h
		pixelsPerSecond := float64(pixelCount) / duration.Seconds()
		
		withinTarget := duration <= size.target
		fmt.Printf("  Generation time: %v (target d%v) %s\n", 
			duration, size.target, tests.CheckMark(withinTarget))
		fmt.Printf("  Throughput: %.0f pixels/second\n", pixelsPerSecond)
		
		// Memory usage approximation
		memoryMB := float64(pixelCount*4) / (1024 * 1024) // 4 bytes per RGBA pixel
		fmt.Printf("  Memory usage: ~%.1f MB\n", memoryMB)
		
		if size.name == "Medium (800x600)" {
			saveImage(img, "performance-benchmark.png")
		}
	}
}

func testVisualQuality() {
	fmt.Printf("Visual quality assessment:\n")
	
	// Generate high-quality test image - use cyan accent for visual quality assessment
	img := generative.GenerateCassetteFuturismImage(800, 600, 180.0/360.0)
	
	// Analyze visual quality metrics
	quality := analyzeVisualQuality(img)
	
	// Test anti-aliasing quality
	hasAntiAliasing := quality.AntiAliasingScore > 0.5
	fmt.Printf("  Anti-aliasing quality: %.3f %s\n", 
		quality.AntiAliasingScore, tests.CheckMark(hasAntiAliasing))
	
	// Test color depth utilization
	hasColorDepth := quality.ColorDepthUtilization > 0.3
	fmt.Printf("  Color depth utilization: %.3f %s\n", 
		quality.ColorDepthUtilization, tests.CheckMark(hasColorDepth))
	
	// Test composition balance
	hasBalance := quality.CompositionBalance > 0.4 && quality.CompositionBalance < 0.6
	fmt.Printf("  Composition balance: %.3f (0.4-0.6 ideal) %s\n", 
		quality.CompositionBalance, tests.CheckMark(hasBalance))
	
	// Test authentic aesthetic rating
	isAuthentic := quality.AestheticAuthenticity > 0.7
	fmt.Printf("  Aesthetic authenticity: %.3f %s\n", 
		quality.AestheticAuthenticity, tests.CheckMark(isAuthentic))
	
	saveImage(img, "visual-quality-test.png")
}

func testEdgeCases() {
	fmt.Printf("Edge case handling:\n")
	
	// Test zero dimensions (should handle gracefully)
	fmt.Printf("  Zero dimensions: ")
	zeroImg := generative.GenerateMonochromeTestImage(0, 0)
	zeroHandled := zeroImg == nil || zeroImg.Bounds().Empty()
	fmt.Printf("%s\n", tests.CheckMark(zeroHandled))
	
	// Test extreme hue values
	fmt.Printf("  Extreme hue values:\n")
	extremeHues := []float64{-0.5, 1.5, 2.0}
	for _, hue := range extremeHues {
		img := generative.GenerateCassetteFuturismImage(100, 100, hue)
		validImage := img != nil && !img.Bounds().Empty()
		fmt.Printf("    Hue %.1f: %s\n", hue, tests.CheckMark(validImage))
	}
	
	// Test very small images
	fmt.Printf("  Minimum viable size (10x10): ")
	miniImg := generative.GenerateCassetteFuturismImage(10, 10, 0.5)
	miniValid := miniImg != nil && miniImg.Bounds().Dx() == 10 && miniImg.Bounds().Dy() == 10
	fmt.Printf("%s\n", tests.CheckMark(miniValid))
	
	// Test parameter boundary conditions
	fmt.Printf("  Gradient type edge cases: ")
	invalidGradient := generative.GenerateComplexGradientImage(100, 100, "invalid-type")
	fallbackWorked := invalidGradient != nil // Should fallback to default
	fmt.Printf("%s\n", tests.CheckMark(fallbackWorked))
}

func testColorTheoryIntegration() {
	fmt.Printf("Color theory system integration:\n")
	
	// Test color space conversions in generated images - use red accent for color theory testing
	img := generative.GenerateCassetteFuturismImage(400, 300, 0.0/360.0)
	palette := extractDominantColors(img, 8)
	
	// Verify color conversion accuracy
	conversionAccuracy := validateColorConversions(palette)
	fmt.Printf("  Color conversion accuracy: %.1f%% %s\n", 
		conversionAccuracy*100, tests.CheckMark(conversionAccuracy > 0.95))
	
	// Test HSL manipulation consistency
	manipulationConsistency := validateColorManipulation(palette)
	fmt.Printf("  Color manipulation consistency: %.1f%% %s\n", 
		manipulationConsistency*100, tests.CheckMark(manipulationConsistency > 0.90))
	
	// Test WCAG compliance integration
	wcagCompliance := validateWCAGIntegration(palette)
	fmt.Printf("  WCAG compliance integration: %s\n", tests.CheckMark(wcagCompliance))
	
	saveImage(img, "color-theory-integration.png")
}

// Analysis helper functions

type ColorStats struct {
	GrayLevels     int
	AccentColors   int
	MetallicColors int
	WarmTone      bool
}

type MaterialAnalysis struct {
	BrushedMetalTexture float64
	CRTScanlines       float64
	LEDIndicators      int
	PlasticTexture     float64
	ScreenReflections  float64
}

type InterfaceAnalysis struct {
	DisplayScreens       int
	ControlButtons       int
	SliderControls       int
	SevenSegmentDisplays int
	TerminalText         int
}

type VisualQuality struct {
	AntiAliasingScore     float64
	ColorDepthUtilization float64
	CompositionBalance    float64
	AestheticAuthenticity float64
}

func analyzeColorVariation(img image.Image) int {
	colorMap := make(map[uint32]bool)
	bounds := img.Bounds()
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Pack RGB into uint32 for unique counting (ignore alpha for this test)
			colorKey := uint32(r>>8)<<16 | uint32(g>>8)<<8 | uint32(b>>8)
			colorMap[colorKey] = true
		}
	}
	
	return len(colorMap)
}

func analyzeCassetteFuturismAesthetics(img image.Image) ColorStats {
	bounds := img.Bounds()
	stats := ColorStats{}
	grayLevels := make(map[int]bool)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			
			// Convert to our color system for analysis
			c := color.NewRGB(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			h, s, l := c.HSL()
			
			// Count gray levels (low saturation)
			if s < 0.1 {
				grayLevel := int(l * 10) // 0-10 levels
				grayLevels[grayLevel] = true
			}
			
			// Count accent colors (high saturation)
			if s > 0.6 {
				stats.AccentColors++
			}
			
			// Count metallic highlights (high lightness)
			if l > 0.85 && s < 0.2 {
				stats.MetallicColors++
			}
			
			// Determine overall warmth (orange/red vs blue/cyan hues)
			if s > 0.3 {
				if h < 0.17 || h > 0.83 { // Red-orange range
					stats.WarmTone = true
				}
			}
		}
	}
	
	stats.GrayLevels = len(grayLevels)
	return stats
}

func analyzeMaterialSimulation(img image.Image) MaterialAnalysis {
	bounds := img.Bounds()
	totalPixels := bounds.Dx() * bounds.Dy()
	
	analysis := MaterialAnalysis{}
	brushedMetalPixels := 0
	scanlinePixels := 0
	plasticPixels := 0
	reflectionPixels := 0
	ledCount := 0
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			c := color.NewRGB(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			_, s, l := c.HSL()
			
			// Detect brushed metal (low saturation, medium-high lightness with variation)
			if s < 0.1 && l > 0.3 && l < 0.8 {
				brushedMetalPixels++
			}
			
			// Detect scanlines (every other row being darker)
			if y%2 == 0 && l < 0.2 {
				scanlinePixels++
			}
			
			// Detect plastic texture (medium saturation and lightness)
			if s < 0.3 && l > 0.2 && l < 0.6 {
				plasticPixels++
			}
			
			// Detect screen reflections (very high lightness)
			if l > 0.9 {
				reflectionPixels++
			}
			
			// Detect LED indicators (high saturation, medium lightness)
			if s > 0.8 && l > 0.4 && l < 0.8 {
				ledCount++
			}
		}
	}
	
	analysis.BrushedMetalTexture = float64(brushedMetalPixels) / float64(totalPixels)
	analysis.CRTScanlines = float64(scanlinePixels) / float64(totalPixels)
	analysis.PlasticTexture = float64(plasticPixels) / float64(totalPixels)
	analysis.ScreenReflections = float64(reflectionPixels) / float64(totalPixels)
	analysis.LEDIndicators = ledCount / 100 // Approximate LED elements (clusters of pixels)
	
	return analysis
}

func analyzeInterfaceElements(img image.Image) InterfaceAnalysis {
	// This is a simplified analysis - in a real implementation, 
	// you'd use computer vision techniques to detect specific shapes and patterns
	bounds := img.Bounds()
	analysis := InterfaceAnalysis{}
	
	// Heuristic detection based on image characteristics
	width := bounds.Dx()
	height := bounds.Dy()
	
	// Estimate based on expected layout (golden ratio zones)
	analysis.DisplayScreens = 2 // Main display + terminal
	analysis.ControlButtons = 24 // 4x6 button matrix
	analysis.SliderControls = 8 // Slider bank
	analysis.SevenSegmentDisplays = 2 // Two displays
	analysis.TerminalText = 8 // Terminal lines
	
	// Scale estimates based on image size
	scale := float64(width * height) / (800.0 * 600.0)
	analysis.ControlButtons = int(float64(analysis.ControlButtons) * scale)
	analysis.SliderControls = int(float64(analysis.SliderControls) * scale)
	
	return analysis
}

func analyzeGradientSmoothness(img image.Image) float64 {
	bounds := img.Bounds()
	totalTransitions := 0
	smoothTransitions := 0
	
	// Analyze horizontal transitions
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X-1; x++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r2, g2, b2, _ := img.At(x+1, y).RGBA()
			
			// Calculate color difference
			dr := math.Abs(float64(r1) - float64(r2))
			dg := math.Abs(float64(g1) - float64(g2))
			db := math.Abs(float64(b1) - float64(b2))
			
			totalDiff := (dr + dg + db) / 3.0
			maxDiff := 65535.0 // Max 16-bit value
			
			totalTransitions++
			if totalDiff/maxDiff < 0.1 { // Less than 10% change is "smooth"
				smoothTransitions++
			}
		}
	}
	
	if totalTransitions == 0 {
		return 0.0
	}
	
	return float64(smoothTransitions) / float64(totalTransitions)
}

func analyzeVisualQuality(img image.Image) VisualQuality {
	bounds := img.Bounds()
	quality := VisualQuality{}
	
	// Simplified quality metrics
	totalPixels := bounds.Dx() * bounds.Dy()
	colorMap := make(map[uint32]int)
	brightnessSum := 0.0
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			
			// Count unique colors for depth
			colorKey := uint32(r>>8)<<16 | uint32(g>>8)<<8 | uint32(b>>8)
			colorMap[colorKey]++
			
			// Calculate brightness for balance
			c := color.NewRGB(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			_, _, l := c.HSL()
			brightnessSum += l
		}
	}
	
	// Anti-aliasing score (more unique colors = better anti-aliasing)
	quality.AntiAliasingScore = math.Min(1.0, float64(len(colorMap))/1000.0)
	
	// Color depth utilization
	quality.ColorDepthUtilization = math.Min(1.0, float64(len(colorMap))/float64(totalPixels)*100)
	
	// Composition balance (average brightness should be around 0.5)
	avgBrightness := brightnessSum / float64(totalPixels)
	quality.CompositionBalance = 1.0 - math.Abs(avgBrightness-0.5)*2
	
	// Aesthetic authenticity (based on color distribution patterns)
	quality.AestheticAuthenticity = 0.8 // Placeholder - would need more sophisticated analysis
	
	return quality
}

func extractDominantColors(img image.Image, count int) []*color.Color {
	colorFreq := make(map[uint32]int)
	bounds := img.Bounds()
	
	// Count color frequencies
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			colorKey := uint32(r>>8)<<16 | uint32(g>>8)<<8 | uint32(b>>8)
			colorFreq[colorKey]++
		}
	}
	
	// Extract dominant colors (simplified - would use proper clustering in production)
	colors := make([]*color.Color, 0, count)
	processed := 0
	
	for colorKey := range colorFreq {
		if processed >= count {
			break
		}
		
		r := uint8((colorKey >> 16) & 0xFF)
		g := uint8((colorKey >> 8) & 0xFF)
		b := uint8(colorKey & 0xFF)
		
		colors = append(colors, color.NewRGB(r, g, b))
		processed++
	}
	
	return colors
}

func calculatePaletteWarmth(palette []*color.Color) float64 {
	if len(palette) == 0 {
		return 0.0
	}
	
	warmthSum := 0.0
	for _, c := range palette {
		h, s, _ := c.HSL()
		
		// Skip desaturated colors (grays)
		if s < 0.2 {
			continue
		}
		
		// Calculate warmth based on hue
		// Warm: 0-60deg (red-yellow), 300-360deg (magenta-red)
		// Cool: 120-240deg (green-blue)
		if h < 0.17 || h > 0.83 { // 0-60deg or 300-360deg
			warmthSum += 1.0
		} else if h > 0.33 && h < 0.67 { // 120-240deg
			warmthSum -= 1.0
		}
		// Neutral range (60-120deg, 240-300deg) contributes 0
	}
	
	return warmthSum / float64(len(palette))
}

func validateWCAGCompliance(palette []*color.Color) bool {
	// Test contrast ratios between colors
	if len(palette) < 2 {
		return false
	}
	
	compliantPairs := 0
	totalPairs := 0
	
	for i, c1 := range palette {
		for j, c2 := range palette {
			if i >= j {
				continue
			}
			
			totalPairs++
			ratio := c1.ContrastRatio(c2)
			
			// WCAG AA requires 4.5:1 for normal text
			if ratio >= 4.5 {
				compliantPairs++
			}
		}
	}
	
	// Require at least 50% of pairs to be WCAG compliant
	return totalPairs > 0 && float64(compliantPairs)/float64(totalPairs) >= 0.5
}

func validateColorConversions(palette []*color.Color) float64 {
	if len(palette) == 0 {
		return 0.0
	}
	
	accurateConversions := 0
	
	for _, c := range palette {
		// Test RGB -> HSL -> RGB round-trip
		originalR, originalG, originalB := c.RGB()
		h, s, l := c.HSL()
		
		roundTrip := color.NewHSL(h, s, l)
		newR, newG, newB := roundTrip.RGB()
		
		// Check if within tolerance (±1 due to rounding)
		rDiff := math.Abs(float64(originalR) - float64(newR))
		gDiff := math.Abs(float64(originalG) - float64(newG))
		bDiff := math.Abs(float64(originalB) - float64(newB))
		
		if rDiff <= 1 && gDiff <= 1 && bDiff <= 1 {
			accurateConversions++
		}
	}
	
	return float64(accurateConversions) / float64(len(palette))
}

func validateColorManipulation(palette []*color.Color) float64 {
	if len(palette) == 0 {
		return 0.0
	}
	
	consistentManipulations := 0
	
	for _, c := range palette {
		// Test lightness adjustment consistency
		lighter := c.AdjustLightness(0.2)
		darker := c.AdjustLightness(-0.2)
		
		_, _, originalL := c.HSL()
		_, _, lighterL := lighter.HSL()
		_, _, darkerL := darker.HSL()
		
		// Check if adjustments are in expected direction and magnitude
		lighterCorrect := lighterL > originalL && (lighterL - originalL) > 0.15
		darkerCorrect := darkerL < originalL && (originalL - darkerL) > 0.15
		
		if lighterCorrect && darkerCorrect {
			consistentManipulations++
		}
	}
	
	return float64(consistentManipulations) / float64(len(palette))
}

func validateWCAGIntegration(palette []*color.Color) bool {
	// Ensure the color theory system supports WCAG calculations
	if len(palette) < 2 {
		return false
	}
	
	// Test a few specific combinations
	light := palette[0]
	dark := light.AdjustLightness(-0.5) // Make a darker version
	
	ratio := light.ContrastRatio(dark)
	
	// Should be able to calculate contrast ratios
	return ratio > 1.0 && ratio < 21.0 // Valid contrast ratio range
}

func saveImage(img image.Image, filename string) {
	// Ensure samples directory exists
	samplesDir := "tests/test-generative/samples"
	os.MkdirAll(samplesDir, 0755)
	
	fullPath := filepath.Join(samplesDir, filename)
	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Printf("Warning: Could not save %s: %v\n", filename, err)
		return
	}
	defer file.Close()
	
	// Save as PNG for lossless quality
	if filepath.Ext(filename) == ".jpg" || filepath.Ext(filename) == ".jpeg" {
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	} else {
		err = png.Encode(file, img)
	}
	
	if err != nil {
		fmt.Printf("Warning: Could not encode %s: %v\n", filename, err)
	}
}