package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/loader"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/processor"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

type ImageResult struct {
	Name         string
	Width        int
	Height       int
	PixelCount   int
	LoadTime     time.Duration
	ProcessTime  time.Duration
	TotalTime    time.Duration
	MemoryUsed   float64
	Success      bool
	Error        error
}

func main() {
	// Parse command line flags
	imagesDir := flag.String("images", "tests/images", "Directory containing test images")
	flag.Parse()

	fmt.Printf("Comprehensive Performance Test\n")
	fmt.Printf("Target: < 2 seconds for 4K images (4096x2160 = 8.8MP)\n")
	fmt.Printf("Target: < 100MB peak memory usage\n\n")

	// Find all image files
	entries, err := os.ReadDir(*imagesDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	var images []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			images = append(images, name)
		}
	}

	sort.Strings(images)

	// Initialize processor
	s := settings.DefaultSettings()
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	ctx := context.Background()

	var results []ImageResult
	var totalProcessingTime time.Duration

	fmt.Printf("Processing %d images...\n\n", len(images))

	for _, imageName := range images {
		imagePath := filepath.Join(*imagesDir, imageName)
		result := ImageResult{Name: imageName}

		// Get memory stats before
		var m1 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		start := time.Now()

		// Load image
		img, err := l.LoadImage(ctx, imagePath)
		if err != nil {
			result.Error = err
			result.Success = false
			results = append(results, result)
			continue
		}

		result.LoadTime = time.Since(start)

		// Get image info
		info, err := l.GetImageInfo(ctx, imagePath)
		if err != nil {
			result.Error = err
			result.Success = false
			results = append(results, result)
			continue
		}

		result.Width = info.Width
		result.Height = info.Height
		result.PixelCount = info.PixelCount()

		// Process image
		processStart := time.Now()
		_, err = p.ProcessImage(img)
		if err != nil {
			result.Error = err
			result.Success = false
			results = append(results, result)
			continue
		}
		result.ProcessTime = time.Since(processStart)
		result.TotalTime = time.Since(start)

		// Get memory stats after
		var m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m2)

		result.MemoryUsed = float64(m2.Sys-m1.Sys) / 1024 / 1024
		result.Success = true

		totalProcessingTime += result.TotalTime

		// Print individual result
		fmt.Printf("%-20s %dx%-4d %6.1fMP %6.0fms %6.1fMB %s\n",
			result.Name,
			result.Width, result.Height,
			float64(result.PixelCount)/1000000,
			float64(result.TotalTime.Nanoseconds())/1000000,
			result.MemoryUsed,
			getPerformanceStatus(result))

		results = append(results, result)
	}

	// Statistical Analysis
	fmt.Printf("\n%s\n", strings.Repeat("=", 70))
	fmt.Printf("PERFORMANCE ANALYSIS SUMMARY\n")
	fmt.Printf("%s\n\n", strings.Repeat("=", 70))

	successful := filterSuccessful(results)
	if len(successful) == 0 {
		fmt.Printf("No images processed successfully!\n")
		return
	}

	// Time statistics
	times := extractTimes(successful)
	sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
	
	avgTime := totalProcessingTime / time.Duration(len(successful))
	minTime := times[0]
	maxTime := times[len(times)-1]
	medianTime := times[len(times)/2]

	fmt.Printf("Processing Time Statistics:\n")
	fmt.Printf("  Average: %6.0fms\n", float64(avgTime.Nanoseconds())/1000000)
	fmt.Printf("  Median:  %6.0fms\n", float64(medianTime.Nanoseconds())/1000000)
	fmt.Printf("  Min:     %6.0fms\n", float64(minTime.Nanoseconds())/1000000)
	fmt.Printf("  Max:     %6.0fms\n", float64(maxTime.Nanoseconds())/1000000)
	fmt.Printf("\n")

	// Memory statistics
	memories := extractMemories(successful)
	sort.Float64s(memories)
	
	avgMemory := average(memories)
	minMemory := memories[0]
	maxMemory := memories[len(memories)-1]
	medianMemory := memories[len(memories)/2]

	fmt.Printf("Memory Usage Statistics:\n")
	fmt.Printf("  Average: %6.1fMB\n", avgMemory)
	fmt.Printf("  Median:  %6.1fMB\n", medianMemory)
	fmt.Printf("  Min:     %6.1fMB\n", minMemory)
	fmt.Printf("  Max:     %6.1fMB\n", maxMemory)
	fmt.Printf("\n")

	// Performance by image size
	fmt.Printf("Performance by Image Size:\n")
	small, medium, large := categorizeBySize(successful)
	
	fmt.Printf("  Small (<2MP):   %d images, avg %4.0fms\n", 
		len(small), float64(averageTime(small).Nanoseconds())/1000000)
	fmt.Printf("  Medium (2-8MP): %d images, avg %4.0fms\n", 
		len(medium), float64(averageTime(medium).Nanoseconds())/1000000)
	fmt.Printf("  Large (>8MP):   %d images, avg %4.0fms\n", 
		len(large), float64(averageTime(large).Nanoseconds())/1000000)
	fmt.Printf("\n")

	// Target compliance
	fmt.Printf("Performance Target Analysis:\n")
	timeCompliant := 0
	memoryCompliant := 0
	
	for _, r := range successful {
		if r.TotalTime < 2*time.Second {
			timeCompliant++
		}
		if r.MemoryUsed < 100 {
			memoryCompliant++
		}
	}

	fmt.Printf("  Time Target (< 2s):     %d/%d images (%.1f%%)\n", 
		timeCompliant, len(successful), 100.0*float64(timeCompliant)/float64(len(successful)))
	fmt.Printf("  Memory Target (< 100MB): %d/%d images (%.1f%%)\n", 
		memoryCompliant, len(successful), 100.0*float64(memoryCompliant)/float64(len(successful)))

	if timeCompliant == len(successful) && memoryCompliant == len(successful) {
		fmt.Printf("  🎉 ALL PERFORMANCE TARGETS MET!\n")
	} else {
		fmt.Printf("  ⚠️  Some targets not met - see individual results above\n")
	}

	// Errors summary
	errors := filterErrors(results)
	if len(errors) > 0 {
		fmt.Printf("\nErrors (%d):\n", len(errors))
		for _, r := range errors {
			fmt.Printf("  %s: %v\n", r.Name, r.Error)
		}
	}
}

func getPerformanceStatus(r ImageResult) string {
	if !r.Success {
		return "❌ ERROR"
	}
	
	timeOk := r.TotalTime < 2*time.Second
	memOk := r.MemoryUsed < 100
	
	if timeOk && memOk {
		return "✅ PASS"
	} else if !timeOk && !memOk {
		return "❌ FAIL (T+M)"
	} else if !timeOk {
		return "❌ FAIL (T)"
	} else {
		return "❌ FAIL (M)"
	}
}

func filterSuccessful(results []ImageResult) []ImageResult {
	var successful []ImageResult
	for _, r := range results {
		if r.Success {
			successful = append(successful, r)
		}
	}
	return successful
}

func filterErrors(results []ImageResult) []ImageResult {
	var errors []ImageResult
	for _, r := range results {
		if !r.Success {
			errors = append(errors, r)
		}
	}
	return errors
}

func extractTimes(results []ImageResult) []time.Duration {
	times := make([]time.Duration, len(results))
	for i, r := range results {
		times[i] = r.TotalTime
	}
	return times
}

func extractMemories(results []ImageResult) []float64 {
	memories := make([]float64, len(results))
	for i, r := range results {
		memories[i] = r.MemoryUsed
	}
	return memories
}

func average(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func averageTime(results []ImageResult) time.Duration {
	if len(results) == 0 {
		return time.Duration(0)
	}
	total := time.Duration(0)
	for _, r := range results {
		total += r.TotalTime
	}
	return total / time.Duration(len(results))
}

func categorizeBySize(results []ImageResult) (small, medium, large []ImageResult) {
	for _, r := range results {
		mp := float64(r.PixelCount) / 1000000
		if mp < 2 {
			small = append(small, r)
		} else if mp < 8 {
			medium = append(medium, r)
		} else {
			large = append(large, r)
		}
	}
	return
}