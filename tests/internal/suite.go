package internal

import (
	"fmt"
	"image"
	"time"
)

// RunPerformanceTest runs a comprehensive performance test with different image types.
func RunPerformanceTest() error {
	fmt.Println("=== Omarchy Theme Generator Performance Test ===")
	fmt.Println()

	// Test 1: 4K synthetic image
	fmt.Println("Test 1: 4K Synthetic Image (Target: <2 seconds)")
	img4K, err := GetOrGenerateTestImage("4k-synthetic", Generate4KTestImage)
	if err != nil {
		return fmt.Errorf("failed to get 4K test image: %w", err)
	}

	bench4K, result4K, err := BenchmarkExtraction(img4K, nil)
	if err != nil {
		return fmt.Errorf("4K test failed: %w", err)
	}

	fmt.Printf("  %s\n", bench4K.String())
	meets4K, msg4K := bench4K.MeetsTarget(2*time.Second, 100.0)
	fmt.Printf("  Target Status: %s\n", msg4K)
	if meets4K {
		fmt.Println("  ✅ 4K performance target met!")
	} else {
		fmt.Println("  ❌ 4K performance target missed")
	}
	fmt.Printf("  Analysis: %s\n", result4K.AnalyzeForThemeGeneration().SuggestedStrategy)
	fmt.Println()

	// Test 2: Monochrome image (synthesis test case)
	fmt.Println("Test 2: Monochrome Image (1920x1080)")
	imgMono, err := GetOrGenerateTestImage("monochrome-grayscale", func() image.Image {
		return GenerateMonochromeTestImage(1920, 1080)
	})
	if err != nil {
		return fmt.Errorf("failed to get monochrome test image: %w", err)
	}

	benchMono, resultMono, err := BenchmarkExtraction(imgMono, nil)
	if err != nil {
		return fmt.Errorf("monochrome test failed: %w", err)
	}

	fmt.Printf("  %s\n", benchMono.String())
	analysis := resultMono.AnalyzeForThemeGeneration()
	fmt.Printf("  Analysis: %s (monochrome: %v, avg saturation: %.3f)\n",
		analysis.SuggestedStrategy, analysis.IsMonochrome, analysis.AverageSaturation)
	fmt.Println()

	// Test 3: High contrast image
	fmt.Println("Test 3: High Contrast Image (1920x1080)")
	imgContrast, err := GetOrGenerateTestImage("high-contrast", func() image.Image {
		return GenerateHighContrastTestImage(1920, 1080)
	})
	if err != nil {
		return fmt.Errorf("failed to get high contrast test image: %w", err)
	}

	benchContrast, resultContrast, err := BenchmarkExtraction(imgContrast, nil)
	if err != nil {
		return fmt.Errorf("high contrast test failed: %w", err)
	}

	fmt.Printf("  %s\n", benchContrast.String())
	analysisContrast := resultContrast.AnalyzeForThemeGeneration()
	fmt.Printf("  Analysis: %s (dominant coverage: %.1f%%)\n",
		analysisContrast.SuggestedStrategy, analysisContrast.DominantCoverage)
	fmt.Println()

	fmt.Println("=== Performance Test Complete ===")
	return nil
}

// GenerateTestSamples generates and saves all test sample images to tests/samples/.
// This creates a consistent set of test images for validation and documentation.
func GenerateTestSamples() error {
	fmt.Println("Generating test sample images...")

	samples := []struct {
		name      string
		generator func() image.Image
	}{
		{"4k-synthetic", Generate4KTestImage},
		{"monochrome-grayscale", func() image.Image { return GenerateMonochromeTestImage(1920, 1080) }},
		{"high-contrast", func() image.Image { return GenerateHighContrastTestImage(1920, 1080) }},
		{"small-monochrome", func() image.Image { return GenerateMonochromeTestImage(400, 300) }},
		{"small-high-contrast", func() image.Image { return GenerateHighContrastTestImage(400, 300) }},
	}

	for _, sample := range samples {
		fmt.Printf("  Generating %s.png...\n", sample.name)
		img := sample.generator()
		path := fmt.Sprintf("tests/samples/%s.png", sample.name)
		if err := SaveImage(img, path); err != nil {
			return fmt.Errorf("failed to save %s: %w", sample.name, err)
		}
	}

	// Create a smaller version of 4K for documentation
	fmt.Println("  Generating 4k-synthetic-small.png for documentation...")
	img4K := Generate4KTestImage()
	smallImg := ResampleImage(img4K, 960, 540) // 1/4 scale
	if err := SaveImage(smallImg, "tests/samples/4k-synthetic-small.png"); err != nil {
		return fmt.Errorf("failed to save small 4K image: %w", err)
	}

	fmt.Println("✅ Test samples generated successfully!")
	return nil
}