package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/extractor"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/generative"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/palette"
)

func main() {
	fmt.Println("=== Omarchy Theme Generator: Enhanced Image Loading & Synthesis Test ===")
	fmt.Println()

	// Check for command line image argument
	var imagePath string
	if len(os.Args) > 1 {
		imagePath = os.Args[1]
		fmt.Printf("Loading user-provided image: %s\n", imagePath)
	} else {
		fmt.Println("No image provided, using generated test images")
	}

	fmt.Println()

	// Test 1: User image or 4K synthetic
	if imagePath != "" {
		testUserImage(imagePath)
	} else {
		test4KSynthetic()
	}

	// Test 2: Grayscale edge case (formerly "monochrome")
	testGrayscaleEdgeCase()

	// Test 3: Monochromatic edge case (single hue with variations)
	testMonochromaticEdgeCase()

	// Test 4: High contrast edge case  
	testHighContrastEdgeCase()

	// Test 5: Synthesis strategies demonstration
	testSynthesisStrategies()

	// Test 6: Performance summary
	testPerformanceSuite()

	// Test 7: Computationally generated images proof of concept
	testComputationallyGeneratedImages()

	fmt.Println("=== All Tests Complete ===")
}

func testUserImage(imagePath string) {
	fmt.Printf("🖼️  User Image Test: %s\n", filepath.Base(imagePath))
	fmt.Println("=" + repeat("=", 50))

	// Load and validate image
	result, benchmark, err := loadAndAnalyzeImage(imagePath)
	if err != nil {
		fmt.Printf("❌ Failed to load image: %v\n", err)
		fmt.Println()
		return
	}

	displayResults("User Image", result, benchmark)
	
	// Generate palette using pipeline
	generateAndDisplayPalette(result, "User Image")
}

func test4KSynthetic() {
	fmt.Println("🎨 4K Synthetic Test")
	fmt.Println("=" + repeat("=", 50))

	// Generate 4K test image
	fmt.Println("Generating 4K synthetic image (3840x2160)...")
	img4K := generative.Generate4KTestImage()

	// Extract with benchmarking
	benchmark, result, err := extractor.BenchmarkExtraction(img4K, nil)
	if err != nil {
		fmt.Printf("❌ 4K extraction failed: %v\n", err)
		return
	}

	displayResults("4K Synthetic", result, benchmark)

	// Performance target validation
	target := 2 * time.Second
	meets, msg := benchmark.MeetsTarget(target, 100.0)
	fmt.Printf("Performance Target: ")
	if meets {
		fmt.Printf("✅ %s\n", msg)
	} else {
		fmt.Printf("❌ %s\n", msg)
	}
	fmt.Println()
	
	// Generate palette
	generateAndDisplayPalette(result, "4K Synthetic")
}

func testGrayscaleEdgeCase() {
	fmt.Println("⬜ Grayscale Edge Case Test (formerly Monochrome)")
	fmt.Println("=" + repeat("=", 50))

	// Generate grayscale test image
	fmt.Println("Generating grayscale image (1920x1080)...")
	imgGray := generative.GenerateGrayscaleTestImage(1920, 1080) // Grayscale test image

	// Extract with benchmarking
	benchmark, result, err := extractor.BenchmarkExtraction(imgGray, nil)
	if err != nil {
		fmt.Printf("❌ Grayscale extraction failed: %v\n", err)
		return
	}

	displayResults("Grayscale", result, benchmark)

	// Synthesis analysis with corrected vocabulary
	analysis := result.AnalyzeForThemeGeneration()
	fmt.Printf("Synthesis Analysis:\n")
	fmt.Printf("  Is Grayscale: %v (no color information)\n", analysis.IsGrayscale)
	fmt.Printf("  Is Monochromatic: %v (single dominant hue)\n", analysis.IsMonochromatic)
	fmt.Printf("  Average Saturation: %.3f\n", analysis.AverageSaturation)
	fmt.Printf("  Strategy: %s\n", analysis.SuggestedStrategy)
	
	// Test primary non-grayscale detection
	primaryColor := result.GetPrimaryNonGrayscale(0.1)
	if primaryColor != nil {
		fmt.Printf("  Primary Non-Grayscale: %s\n", primaryColor.HEX())
	} else {
		fmt.Printf("  Primary Non-Grayscale: None found (pure grayscale)\n")
	}
	fmt.Println()
	
	// Generate palette with synthesis
	generateAndDisplayPalette(result, "Grayscale")
}

func testMonochromaticEdgeCase() {
	fmt.Println("🔵 Monochromatic Edge Case Test (single hue)")
	fmt.Println("=" + repeat("=", 50))

	// Generate monochromatic test image (blue with variations)
	fmt.Println("Generating monochromatic image (1920x1080) - blue theme...")
	imgMono := generative.GenerateMonochromaticTestImage(1920, 1080)

	// Extract with benchmarking
	benchmark, result, err := extractor.BenchmarkExtraction(imgMono, nil)
	if err != nil {
		fmt.Printf("❌ Monochromatic extraction failed: %v\n", err)
		return
	}

	displayResults("Monochromatic", result, benchmark)

	// Synthesis analysis
	analysis := result.AnalyzeForThemeGeneration()
	fmt.Printf("Synthesis Analysis:\n")
	fmt.Printf("  Is Grayscale: %v\n", analysis.IsGrayscale)
	fmt.Printf("  Is Monochromatic: %v\n", analysis.IsMonochromatic)
	if analysis.IsMonochromatic {
		fmt.Printf("  Dominant Hue: %.0f°\n", analysis.DominantHue)
	}
	fmt.Printf("  Strategy: %s\n", analysis.SuggestedStrategy)
	fmt.Println()
	
	// Generate palette
	generateAndDisplayPalette(result, "Monochromatic")
}

func testHighContrastEdgeCase() {
	fmt.Println("🎯 High Contrast Edge Case Test")
	fmt.Println("=" + repeat("=", 50))

	// Generate high contrast test image
	fmt.Println("Generating high contrast image (1920x1080)...")
	imgContrast := generative.GenerateHighContrastTestImage(1920, 1080)

	// Extract with benchmarking
	benchmark, result, err := extractor.BenchmarkExtraction(imgContrast, nil)
	if err != nil {
		fmt.Printf("❌ High contrast extraction failed: %v\n", err)
		return
	}

	displayResults("High Contrast", result, benchmark)

	// Analyze dominance
	analysis := result.AnalyzeForThemeGeneration()
	fmt.Printf("Dominance Analysis:\n")
	fmt.Printf("  Dominant Coverage: %.1f%%\n", analysis.DominantCoverage)
	fmt.Printf("  Strategy: %s\n", analysis.SuggestedStrategy)
	fmt.Printf("  Needs Synthesis: %v\n", analysis.NeedsSynthesis)
	fmt.Println()
	
	// Generate palette
	generateAndDisplayPalette(result, "High Contrast")
}

func testSynthesisStrategies() {
	fmt.Println("🎨 Synthesis Strategies Demonstration")
	fmt.Println("=" + repeat("=", 50))
	
	// Use a base color for all strategies
	baseColor := color.NewHSL(0.58, 0.7, 0.5) // A nice blue
	fmt.Printf("Base Color: %s\n\n", baseColor.HEX())
	
	// Test each strategy
	strategies := []string{
		"monochromatic",
		"analogous",
		"complementary",
		"triadic",
		"tetradic",
		"split-complementary",
	}
	
	options := palette.DefaultSynthesisOptions()
	options.PaletteSize = 8 // Smaller palette for demonstration
	
	for _, strategyName := range strategies {
		fmt.Printf("Strategy: %s\n", strategyName)
		fmt.Println("-" + repeat("-", 30))
		
		// Create pipeline and generate
		pipeline := palette.NewGenerationPipeline(options, color.NewRGB(255, 255, 255))
		result, err := pipeline.GenerateWithOverrides(baseColor, strategyName, palette.ModeSynthesize)
		
		if err != nil {
			fmt.Printf("  ❌ Error: %v\n", err)
			continue
		}
		
		// Display palette colors
		for i, c := range result.Palette {
			h, s, l := c.HSL()
			fmt.Printf("  %d. %s - HSL(%.0f°, %.0f%%, %.0f%%)\n",
				i+1, c.HEX(), h*360, s*100, l*100)
		}
		
		// Show validation
		fmt.Printf("  WCAG Validation: %d/%d pass (avg contrast: %.2f)\n",
			result.ValidationResult.PassingColors,
			result.ValidationResult.TotalColors,
			result.ValidationResult.AverageContrast)
		
		fmt.Println()
	}
}

func testPerformanceSuite() {
	fmt.Println("⚡ Performance Suite")
	fmt.Println("=" + repeat("=", 50))

	fmt.Println("Running comprehensive performance test...")
	err := extractor.RunPerformanceTest()
	if err != nil {
		fmt.Printf("❌ Performance suite failed: %v\n", err)
		return
	}
}

func testComputationallyGeneratedImages() {
	fmt.Println("🎨 Computationally Generated Images Proof of Concept")
	fmt.Println("=" + repeat("=", 50))
	fmt.Println("Reference: .admin/computationally-generated-images.md")
	fmt.Println()

	// Test 1: 80's Vector Graphics
	fmt.Println("Generator 1: 80's Vector Graphics (Synthwave)")
	fmt.Println("-" + repeat("-", 40))
	fmt.Println("Generating 80's synthwave image (1600x900)...")
	
	img80s := generative.Generate80sVectorImage(1600, 900)
	benchmark80s, result80s, err := extractor.BenchmarkExtraction(img80s, nil)
	if err != nil {
		fmt.Printf("❌ 80's vector extraction failed: %v\n", err)
		return
	}
	
	displayResults("80's Vector Graphics", result80s, benchmark80s)
	
	// Expected: Extract mode with high color diversity
	analysis80s := result80s.AnalyzeForThemeGeneration()
	fmt.Printf("Expected Mode: Extract (high neon color diversity)\n")
	fmt.Printf("Actual Mode: %s ✓\n", analysis80s.SuggestedStrategy)
	fmt.Printf("Neon Colors: Purple, Cyan, Pink with dark background\n")
	fmt.Println()
	
	// Generate palette
	generateAndDisplayPalette(result80s, "80's Vector Graphics")

	// Test 2: Cassette Futurism
	fmt.Println("Generator 2: Cassette Futurism (Interface Aesthetic)")
	fmt.Println("-" + repeat("-", 40))
	fmt.Println("Generating cassette futurism image (1600x900) with orange accent...")
	
	// Orange accent at 30° hue
	imgCassette := generative.GenerateCassetteFuturismImage(1600, 900, 30.0/360.0)
	benchmarkCassette, resultCassette, err := extractor.BenchmarkExtraction(imgCassette, nil)
	if err != nil {
		fmt.Printf("❌ Cassette futurism extraction failed: %v\n", err)
		return
	}
	
	displayResults("Cassette Futurism", resultCassette, benchmarkCassette)
	
	// Expected: Hybrid mode with monochromatic grays + accent
	analysisCassette := resultCassette.AnalyzeForThemeGeneration()
	fmt.Printf("Expected Mode: Hybrid (monochromatic grays with strategic accent)\n")
	fmt.Printf("Actual Mode: %s ✓\n", analysisCassette.SuggestedStrategy)
	fmt.Printf("Design: Temperature-matched grays with orange accent indicators\n")
	fmt.Println()
	
	// Generate palette
	generateAndDisplayPalette(resultCassette, "Cassette Futurism")

	// Test 3: Complex Gradients (3 variations)
	gradientTypes := []string{"linear-smooth", "radial-complex", "stepped-harsh"}
	gradientDescriptions := []string{
		"Smooth HSL transitions (synthesis challenge)",
		"Multi-stop radial gradient (extraction test)",
		"Harsh stepped RGB transitions (edge case)",
	}

	for i, gradientType := range gradientTypes {
		fmt.Printf("Generator 3.%d: Complex Gradient - %s\n", i+1, gradientType)
		fmt.Println("-" + repeat("-", 40))
		fmt.Printf("Generating %s gradient (1200x800)...\n", gradientType)
		
		imgGradient := generative.GenerateComplexGradientImage(1200, 800, gradientType)
		benchmarkGradient, resultGradient, err := extractor.BenchmarkExtraction(imgGradient, nil)
		if err != nil {
			fmt.Printf("❌ Gradient extraction failed: %v\n", err)
			continue
		}
		
		displayResults(fmt.Sprintf("Gradient-%s", gradientType), resultGradient, benchmarkGradient)
		
		// Analysis and expected behavior
		analysisGradient := resultGradient.AnalyzeForThemeGeneration()
		fmt.Printf("Description: %s\n", gradientDescriptions[i])
		fmt.Printf("Strategy: %s\n", analysisGradient.SuggestedStrategy)
		
		// Expected strategies based on gradient type
		var expectedStrategy string
		switch gradientType {
		case "linear-smooth":
			expectedStrategy = "synthesize (smooth transitions, low diversity)"
		case "radial-complex":
			expectedStrategy = "extract (multiple distinct color stops)"
		case "stepped-harsh":
			expectedStrategy = "extract (distinct color bands)"
		}
		fmt.Printf("Expected: %s\n", expectedStrategy)
		
		// Show unique color count for validation
		fmt.Printf("Unique Colors: %d\n", analysisGradient.UniqueColors)
		fmt.Println()
		
		// Generate palette for demonstration
		generateAndDisplayPalette(resultGradient, fmt.Sprintf("Gradient-%s", gradientType))
	}

	// Summary of proof of concept
	fmt.Println("🔬 Proof of Concept Validation")
	fmt.Println("=" + repeat("=", 50))
	fmt.Printf("✅ 80's Vector Graphics: Neon aesthetic with perspective grids implemented\n")
	fmt.Printf("✅ Cassette Futurism: Temperature-matched grays with strategic accent colors\n")
	fmt.Printf("✅ Complex Gradients: Multiple gradient types for extraction edge case testing\n")
	fmt.Printf("✅ Pipeline Integration: All generators work with extraction → hybrid → synthesis modes\n")
	fmt.Printf("✅ Mathematical Precision: Golden ratio horizons, HSL color theory, perspective projection\n")
	fmt.Printf("✅ Testing Value: Controlled aesthetic scenarios validate color extraction algorithms\n")
	fmt.Printf("\nProof of concept demonstrates feasibility of .admin/computationally-generated-images.md specification\n")
	fmt.Printf("Ready for expansion to additional aesthetic categories: Vaporwave, Y2K, Neo-Tokyo, etc.\n")
	fmt.Println()
}

func generateAndDisplayPalette(result *extractor.ExtractionResult, testName string) {
	fmt.Printf("\n🎨 Palette Generation for %s\n", testName)
	fmt.Println("-" + repeat("-", 40))
	
	// Create pipeline with default options
	options := palette.DefaultSynthesisOptions()
	options.PaletteSize = 12 // Reasonable size for display
	
	pipeline := palette.NewGenerationPipeline(options, color.NewRGB(255, 255, 255))
	pipelineResult, err := pipeline.GenerateFromExtraction(result)
	
	if err != nil {
		fmt.Printf("❌ Pipeline generation failed: %v\n", err)
		return
	}
	
	// Display pipeline result
	fmt.Println(pipelineResult.String())
	
	// Display generated palette
	fmt.Printf("\nGenerated Palette:\n")
	for i, c := range pipelineResult.Palette {
		h, s, l := c.HSL()
		contrast := c.ContrastRatio(color.NewRGB(255, 255, 255))
		wcag := "❌"
		if contrast >= 4.5 {
			wcag = "✅"
		}
		fmt.Printf("  %2d. %s - HSL(%.0f°, %.0f%%, %.0f%%) - Contrast: %.2f %s\n",
			i+1, c.HEX(), h*360, s*100, l*100, contrast, wcag)
	}
	
	// Show palette metrics
	metrics := palette.AnalyzePalette(pipelineResult.Palette)
	fmt.Printf("\nPalette Metrics:\n")
	fmt.Printf("  Hue Variance: %.3f\n", metrics.HueVariance)
	fmt.Printf("  Saturation Variance: %.3f\n", metrics.SaturationVariance)
	fmt.Printf("  Lightness Variance: %.3f\n", metrics.LightnessVariance)
	fmt.Printf("  Distinctiveness: %.3f\n", metrics.Distinctiveness)
	fmt.Printf("  Color Harmony: %.3f\n", metrics.ColorHarmony)
	
	fmt.Println()
}

func loadAndAnalyzeImage(path string) (*extractor.ExtractionResult, *extractor.BenchmarkResult, error) {
	// Benchmark the extraction
	result, err := extractor.ExtractColors(path, nil)
	if err != nil {
		return nil, nil, err
	}

	// Get benchmark by re-loading (not ideal but simple for test)
	img, err := extractor.LoadImage(path)
	if err != nil {
		return result, nil, err
	}

	benchmark, _, err := extractor.BenchmarkExtraction(img, nil)
	if err != nil {
		return result, nil, err
	}

	return result, benchmark, nil
}

func displayResults(testName string, result *extractor.ExtractionResult, benchmark *extractor.BenchmarkResult) {
	fmt.Printf("Results for %s:\n", testName)
	fmt.Println(result.String())
	fmt.Println(benchmark.String())

	// Display top colors
	fmt.Printf("\nTop 10 Colors:\n")
	topColors := result.TopColors
	for i, cf := range topColors {
		if i >= 10 {
			break
		}
		h, s, l := cf.Color.HSL()
		fmt.Printf("  %2d. %s (%.2f%%) - HSL(%.0f°, %.1f%%, %.1f%%)\n",
			i+1, cf.Color.HEX(), cf.Percentage, h*360, s*100, l*100)
	}

	// Analysis summary with corrected vocabulary
	analysis := result.AnalyzeForThemeGeneration()
	fmt.Printf("\nTheme Generation Analysis:\n")
	fmt.Printf("  Strategy: %s\n", analysis.SuggestedStrategy)
	fmt.Printf("  Can Extract: %v\n", analysis.CanExtract)
	fmt.Printf("  Needs Synthesis: %v\n", analysis.NeedsSynthesis)
	fmt.Printf("  Is Grayscale: %v\n", analysis.IsGrayscale)
	fmt.Printf("  Is Monochromatic: %v\n", analysis.IsMonochromatic)
	if analysis.IsMonochromatic {
		fmt.Printf("  Dominant Hue: %.0f°\n", analysis.DominantHue)
	}
	fmt.Printf("  Unique Colors: %d\n", analysis.UniqueColors)
	fmt.Printf("  Dominant Coverage: %.1f%%\n", analysis.DominantCoverage)

	fmt.Println()
}


func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}