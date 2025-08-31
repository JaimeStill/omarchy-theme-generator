package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/extractor"
	"github.com/JaimeStill/omarchy-theme-generator/tests/internal"
)

func main() {
	fmt.Println("=== Omarchy Theme Generator: Image Loading Test ===")
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

	// Test 2: Monochrome edge case
	testMonochromeEdgeCase()

	// Test 3: High contrast edge case  
	testHighContrastEdgeCase()

	// Test 4: Performance summary
	testPerformanceSuite()

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
}

func test4KSynthetic() {
	fmt.Println("🎨 4K Synthetic Test")
	fmt.Println("=" + repeat("=", 50))

	// Generate 4K test image
	fmt.Println("Generating 4K synthetic image (3840x2160)...")
	img4K := internal.Generate4KTestImage()

	// Extract with benchmarking
	benchmark, result, err := internal.BenchmarkExtraction(img4K, nil)
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
}

func testMonochromeEdgeCase() {
	fmt.Println("⬜ Monochrome Edge Case Test")
	fmt.Println("=" + repeat("=", 50))

	// Generate monochrome test image
	fmt.Println("Generating monochrome image (1920x1080)...")
	imgMono := internal.GenerateMonochromeTestImage(1920, 1080)

	// Extract with benchmarking
	benchmark, result, err := internal.BenchmarkExtraction(imgMono, nil)
	if err != nil {
		fmt.Printf("❌ Monochrome extraction failed: %v\n", err)
		return
	}

	displayResults("Monochrome", result, benchmark)

	// Synthesis analysis
	analysis := result.AnalyzeForThemeGeneration()
	fmt.Printf("Synthesis Analysis:\n")
	fmt.Printf("  Is Monochrome: %v\n", analysis.IsMonochrome)
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
}

func testHighContrastEdgeCase() {
	fmt.Println("🎯 High Contrast Edge Case Test")
	fmt.Println("=" + repeat("=", 50))

	// Generate high contrast test image
	fmt.Println("Generating high contrast image (1920x1080)...")
	imgContrast := internal.GenerateHighContrastTestImage(1920, 1080)

	// Extract with benchmarking
	benchmark, result, err := internal.BenchmarkExtraction(imgContrast, nil)
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
}

func testPerformanceSuite() {
	fmt.Println("⚡ Performance Suite")
	fmt.Println("=" + repeat("=", 50))

	fmt.Println("Running comprehensive performance test...")
	err := internal.RunPerformanceTest()
	if err != nil {
		fmt.Printf("❌ Performance suite failed: %v\n", err)
		return
	}
}

func loadAndAnalyzeImage(path string) (*extractor.ExtractionResult, *internal.BenchmarkResult, error) {
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

	benchmark, _, err := internal.BenchmarkExtraction(img, nil)
	if err != nil {
		return result, nil, err
	}

	return result, benchmark, nil
}

func displayResults(testName string, result *extractor.ExtractionResult, benchmark *internal.BenchmarkResult) {
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
			i+1, cf.Color.HEX(), cf.Percentage, h, s*100, l*100)
	}

	// Analysis summary
	analysis := result.AnalyzeForThemeGeneration()
	fmt.Printf("\nTheme Generation Analysis:\n")
	fmt.Printf("  Strategy: %s\n", analysis.SuggestedStrategy)
	fmt.Printf("  Can Extract: %v\n", analysis.CanExtract)
	fmt.Printf("  Needs Synthesis: %v\n", analysis.NeedsSynthesis)
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