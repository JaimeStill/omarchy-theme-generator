package extractor

import (
	"fmt"
	"image"
	"runtime"
	"time"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/generative"
)

// BenchmarkResult contains detailed performance metrics for extraction operations.
type BenchmarkResult struct {
	Duration       time.Duration // Total extraction time
	PeakMemoryMB   float64       // Peak memory usage in MB
	PixelsPerSec   uint64        // Pixels processed per second
	ColorsPerSec   uint64        // Unique colors identified per second
	TotalPixels    uint64        // Total pixels in the image
	UniqueColors   int           // Number of unique colors found
	ImageDimensions string       // Width x Height for context
}

// String returns a formatted string representation of the benchmark results.
func (br *BenchmarkResult) String() string {
	return fmt.Sprintf("Performance: %v duration, %.2f MB peak memory, %d pixels/sec, %d unique colors from %s image",
		br.Duration, br.PeakMemoryMB, br.PixelsPerSec, br.UniqueColors, br.ImageDimensions)
}

// MeetsTarget checks if the benchmark meets the specified performance targets.
func (br *BenchmarkResult) MeetsTarget(maxDuration time.Duration, maxMemoryMB float64) (bool, string) {
	durationMet := br.Duration <= maxDuration
	memoryMet := br.PeakMemoryMB <= maxMemoryMB

	var msg string
	if durationMet && memoryMet {
		msg = fmt.Sprintf("All performance targets met (%.2fs, %.2fMB)", 
			br.Duration.Seconds(), br.PeakMemoryMB)
	} else if !durationMet && !memoryMet {
		msg = fmt.Sprintf("Both targets exceeded: %.2fs > %.2fs, %.2fMB > %.2fMB", 
			br.Duration.Seconds(), maxDuration.Seconds(), br.PeakMemoryMB, maxMemoryMB)
	} else if !durationMet {
		msg = fmt.Sprintf("Duration target exceeded: %.2fs > %.2fs", 
			br.Duration.Seconds(), maxDuration.Seconds())
	} else {
		msg = fmt.Sprintf("Memory target exceeded: %.2fMB > %.2fMB", 
			br.PeakMemoryMB, maxMemoryMB)
	}

	return durationMet && memoryMet, msg
}

// BenchmarkExtraction performs a benchmarked color extraction with detailed performance metrics.
func BenchmarkExtraction(img image.Image, options *ExtractionOptions) (*BenchmarkResult, *ExtractionResult, error) {
	// Memory tracking
	var memBefore, memAfter runtime.MemStats
	runtime.GC() // Force garbage collection for accurate measurement
	runtime.ReadMemStats(&memBefore)
	
	// Time the extraction
	startTime := time.Now()
	result, err := ExtractColorsFromImage(img, options)
	duration := time.Since(startTime)
	
	// Capture post-extraction memory
	runtime.ReadMemStats(&memAfter)
	
	if err != nil {
		return nil, nil, err
	}
	
	// Calculate metrics
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	totalPixels := uint64(width * height)
	
	// Peak memory calculation (in MB)
	peakMemoryBytes := memAfter.Alloc - memBefore.Alloc
	if memAfter.Sys > memBefore.Sys {
		// If system memory increased, use that as a better indicator
		peakMemoryBytes = memAfter.Sys - memBefore.Sys
	}
	peakMemoryMB := float64(peakMemoryBytes) / (1024 * 1024)
	
	benchmark := &BenchmarkResult{
		Duration:       duration,
		PeakMemoryMB:   peakMemoryMB,
		TotalPixels:    totalPixels,
		UniqueColors:   result.UniqueColors,
		ImageDimensions: fmt.Sprintf("%dx%d", width, height),
	}
	
	// Calculate per-second metrics (avoid division by zero)
	if duration > 0 {
		durationSec := duration.Seconds()
		benchmark.PixelsPerSec = uint64(float64(totalPixels) / durationSec)
		benchmark.ColorsPerSec = uint64(float64(result.UniqueColors) / durationSec)
	}

	return benchmark, result, nil
}

// RunPerformanceTest runs a comprehensive performance test with different image types.
func RunPerformanceTest() error {
	fmt.Println("=== Omarchy Theme Generator Performance Test ===")
	fmt.Println()

	// Test 1: 4K synthetic image
	fmt.Println("Test 1: 4K Synthetic Image (Target: <2 seconds)")
	img4K := generative.Generate4KTestImage()
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
	imgMono := generative.GenerateMonochromeTestImage(1920, 1080)
	benchMono, resultMono, err := BenchmarkExtraction(imgMono, nil)
	if err != nil {
		return fmt.Errorf("monochrome test failed: %w", err)
	}
	
	fmt.Printf("  %s\n", benchMono.String())
	analysis := resultMono.AnalyzeForThemeGeneration()
	fmt.Printf("  Analysis: %s (grayscale: %v, avg saturation: %.3f)\n", 
		analysis.SuggestedStrategy, analysis.IsGrayscale, analysis.AverageSaturation)
	fmt.Println()

	// Test 3: High contrast image
	fmt.Println("Test 3: High Contrast Image (1920x1080)")
	imgContrast := generative.GenerateHighContrastTestImage(1920, 1080)
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