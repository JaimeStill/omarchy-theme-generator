package internal

import (
	"fmt"
	"image"
	"runtime"
	"time"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/extractor"
)

// BenchmarkResult contains detailed performance metrics for extraction operations.
type BenchmarkResult struct {
	Duration        time.Duration // Total extraction time
	MemoryUsed      uint64        // Peak memory usage in bytes
	MemoryAllocs    uint64        // Total memory allocations
	PixelsPerSec    uint64        // Processing rate (pixels/second)
	ColorsPerSec    uint64        // Color extraction rate (unique colors/second)
	ImageDimensions string        // Width x Height
	UniqueColors    int           // Total unique colors found
	TotalPixels     uint32        // Total pixels processed
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

// BenchmarkExtraction measures color extraction performance with detailed metrics.
// It captures memory usage, processing time, and throughput rates for performance analysis.
// Returns benchmark metrics, extraction results, and any error encountered.
func BenchmarkExtraction(img image.Image, options *extractor.ExtractionOptions) (*BenchmarkResult, *extractor.ExtractionResult, error) {
	if options == nil {
		options = extractor.DefaultOptions()
	}

	// Capture initial memory state
	var memBefore, memAfter runtime.MemStats
	runtime.GC() // Force garbage collection for accurate baseline
	runtime.ReadMemStats(&memBefore)

	// Perform timed extraction
	startTime := time.Now()
	result, err := extractor.ExtractFromLoadedImage(img, options)
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
