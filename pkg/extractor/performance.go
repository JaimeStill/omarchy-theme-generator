package extractor

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// BenchmarkResult contains detailed performance metrics for extraction operations.
type BenchmarkResult struct {
	Duration       time.Duration // Total extraction time
	MemoryUsed     uint64        // Peak memory usage in bytes
	MemoryAllocs   uint64        // Total memory allocations
	PixelsPerSec   uint64        // Processing rate (pixels/second)
	ColorsPerSec   uint64        // Color extraction rate (unique colors/second)
	ImageDimensions string       // Width x Height
	UniqueColors   int           // Total unique colors found
	TotalPixels    uint32        // Total pixels processed
}

// String provides a human-readable summary of benchmark results.
func (br *BenchmarkResult) String() string {
	memoryMB := float64(br.MemoryUsed) / (1024 * 1024)
	return fmt.Sprintf(
		"Performance: %v duration, %.2f MB peak memory, %d pixels/sec, %d unique colors from %s image",
		br.Duration, memoryMB, br.PixelsPerSec, br.UniqueColors, br.ImageDimensions,
	)
}

// MeetsTarget checks if benchmark results meet performance targets.
func (br *BenchmarkResult) MeetsTarget(maxDuration time.Duration, maxMemoryMB float64) (bool, string) {
	memoryMB := float64(br.MemoryUsed) / (1024 * 1024)
	
	if br.Duration > maxDuration {
		return false, fmt.Sprintf("Duration %v exceeds target %v", br.Duration, maxDuration)
	}
	
	if memoryMB > maxMemoryMB {
		return false, fmt.Sprintf("Memory %.2f MB exceeds target %.2f MB", memoryMB, maxMemoryMB)
	}
	
	return true, "All performance targets met"
}

// BenchmarkExtraction measures extraction performance with detailed metrics.
// BenchmarkExtraction measures color extraction performance with detailed metrics.
// It captures memory usage, processing time, and throughput rates for performance analysis.
// Returns benchmark metrics, extraction results, and any error encountered.
func BenchmarkExtraction(img image.Image, options *ExtractionOptions) (*BenchmarkResult, *ExtractionResult, error) {
	if options == nil {
		options = DefaultOptions()
	}

	// Capture initial memory state
	var memBefore, memAfter runtime.MemStats
	runtime.GC() // Force garbage collection for accurate baseline
	runtime.ReadMemStats(&memBefore)

	// Perform timed extraction
	startTime := time.Now()
	result, err := ExtractFromLoadedImage(img, options)
	duration := time.Since(startTime)

	if err != nil {
		return nil, nil, fmt.Errorf("benchmark failed during extraction: %w", err)
	}

	// Capture final memory state
	runtime.ReadMemStats(&memAfter)

	// Calculate metrics
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	totalPixels := uint64(width * height)

	benchmark := &BenchmarkResult{
		Duration:        duration,
		MemoryUsed:      memAfter.HeapInuse - memBefore.HeapInuse,
		MemoryAllocs:    memAfter.TotalAlloc - memBefore.TotalAlloc,
		ImageDimensions: fmt.Sprintf("%dx%d", width, height),
		UniqueColors:    result.UniqueColors,
		TotalPixels:     result.TotalPixels,
	}

	// Calculate processing rates (avoid division by zero)
	if duration.Nanoseconds() > 0 {
		durationSec := float64(duration.Nanoseconds()) / 1e9
		benchmark.PixelsPerSec = uint64(float64(totalPixels) / durationSec)
		benchmark.ColorsPerSec = uint64(float64(result.UniqueColors) / durationSec)
	}

	return benchmark, result, nil
}

// Generate4KTestImage creates a synthetic 4K image with varied colors for benchmarking.
// This ensures consistent benchmarking across different environments.
// Generate4KTestImage creates a synthetic 4K image (3840x2160) with varied colors for benchmarking.
// The image contains gradient patterns with noise to simulate realistic color distribution.
// This ensures consistent benchmarking across different environments and validates performance targets.
func Generate4KTestImage() image.Image {
	width, height := 3840, 2160 // 4K resolution
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Create a gradient pattern with varying colors for realistic extraction testing
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Generate varied colors based on position
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(((x + y) * 255) / (width + height))
			
			// Add some noise for more realistic color distribution
			if (x+y)%7 == 0 {
				r = uint8((int(r) + 50) % 256)
			}
			if (x*y)%11 == 0 {
				g = uint8((int(g) + 30) % 256)
			}
			if (x-y)%13 == 0 {
				b = uint8((int(b) + 70) % 256)
			}

			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}

// GenerateMonochromeTestImage creates a grayscale image for testing synthesis edge cases.
func GenerateMonochromeTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Create grayscale gradient with subtle variations
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Base grayscale value
			gray := uint8((x + y) * 255 / (width + height))
			
			// Add subtle noise to create some unique colors but keep it monochrome
			if (x*y)%17 == 0 {
				gray = uint8((int(gray) + 10) % 256)
			}

			img.Set(x, y, color.RGBA{gray, gray, gray, 255})
		}
	}

	return img
}

// GenerateHighContrastTestImage creates an image with few but very distinct colors.
func GenerateHighContrastTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Define high contrast colors
	colors := []color.RGBA{
		{0, 0, 0, 255},       // Black
		{255, 255, 255, 255}, // White
		{255, 0, 0, 255},     // Red
		{0, 255, 0, 255},     // Green
		{0, 0, 255, 255},     // Blue
	}

	// Create blocks of solid colors
	blockWidth := width / len(colors)
	
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			colorIndex := x / blockWidth
			if colorIndex >= len(colors) {
				colorIndex = len(colors) - 1
			}
			
			// Add some vertical stripes for variety
			if y%50 < 10 {
				colorIndex = (colorIndex + 1) % len(colors)
			}
			
			img.Set(x, y, colors[colorIndex])
		}
	}

	return img
}

// RunPerformanceTest runs a comprehensive performance test with different image types.
func RunPerformanceTest() error {
	fmt.Println("=== Omarchy Theme Generator Performance Test ===")
	fmt.Println()

	// Test 1: 4K synthetic image
	fmt.Println("Test 1: 4K Synthetic Image (Target: <2 seconds)")
	img4K := Generate4KTestImage()
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
	imgMono := GenerateMonochromeTestImage(1920, 1080)
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
	imgContrast := GenerateHighContrastTestImage(1920, 1080)
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

// SaveImage saves an image to the specified path as PNG.
func SaveImage(img image.Image, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Create the file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer file.Close()

	// Encode as PNG
	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}

// GenerateAndSaveTestImages creates and saves all test images for documentation.
// GenerateAndSaveTestImages creates and saves all standard test images for documentation.
// It generates 4K synthetic, monochrome grayscale, and high contrast images as PNG files.
// The images are saved to the specified directory for use in test documentation.
func GenerateAndSaveTestImages(outputDir string) error {
	fmt.Printf("Generating test images in %s...\n", outputDir)

	// Generate 4K synthetic image (scaled down for README)
	fmt.Println("  Creating 4k-synthetic.png...")
	img4K := Generate4KTestImage()
	// Create a smaller version for README (960x540 = 1/4 scale)
	smallImg4K := resampleImage(img4K, 960, 540)
	if err := SaveImage(smallImg4K, filepath.Join(outputDir, "4k-synthetic.png")); err != nil {
		return fmt.Errorf("failed to save 4K synthetic image: %w", err)
	}

	// Generate monochrome grayscale image
	fmt.Println("  Creating monochrome-grayscale.png...")
	imgMono := GenerateMonochromeTestImage(400, 300)
	if err := SaveImage(imgMono, filepath.Join(outputDir, "monochrome-grayscale.png")); err != nil {
		return fmt.Errorf("failed to save monochrome image: %w", err)
	}

	// Generate high contrast image
	fmt.Println("  Creating high-contrast.png...")
	imgContrast := GenerateHighContrastTestImage(400, 300)
	if err := SaveImage(imgContrast, filepath.Join(outputDir, "high-contrast.png")); err != nil {
		return fmt.Errorf("failed to save high contrast image: %w", err)
	}

	fmt.Println("  ✅ Test images generated successfully!")
	return nil
}

// resampleImage creates a simple resampled version of an image (nearest neighbor).
func resampleImage(src image.Image, newWidth, newHeight int) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Map destination coordinates to source coordinates
			srcX := (x * srcWidth) / newWidth
			srcY := (y * srcHeight) / newHeight

			// Get source pixel and set destination pixel
			srcColor := src.At(srcX+srcBounds.Min.X, srcY+srcBounds.Min.Y)
			dst.Set(x, y, srcColor)
		}
	}

	return dst
}