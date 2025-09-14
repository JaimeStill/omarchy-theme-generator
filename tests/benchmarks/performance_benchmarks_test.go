package benchmarks_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/loader"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/processor"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

// BenchmarkProcessImage_Small benchmarks processing of small images (<2MP)
func BenchmarkProcessImage_Small(b *testing.B) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	ctx := context.Background()

	// Use a small test image
	imagePath := filepath.Join("..", "images", "grayscale.jpeg")

	// Load image once for all iterations
	img, err := l.LoadImage(ctx, imagePath)
	if err != nil {
		b.Fatalf("Failed to load test image: %v", err)
	}

	// Get image info for reporting
	info, err := l.GetImageInfo(ctx, imagePath)
	if err != nil {
		b.Fatalf("Failed to get image info: %v", err)
	}

	b.Logf("Benchmarking small image: %dx%d (%.1f MP)",
		info.Width, info.Height, float64(info.PixelCount())/1000000)

	// Warm up run to ensure caches are populated
	profile, err := p.ProcessImage(img)
	if err != nil {
		b.Fatalf("Warmup processing failed: %v", err)
	}
	b.Logf("Benchmark result: %d colors, Mode=%s, HasColor=%t",
		profile.ColorCount, profile.Mode, profile.HasColor)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := p.ProcessImage(img)
		if err != nil {
			b.Fatalf("Processing failed: %v", err)
		}
	}
}

// BenchmarkProcessImage_Large benchmarks processing of large images (>8MP)
func BenchmarkProcessImage_Large(b *testing.B) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	ctx := context.Background()
	
	// Use a large test image
	imagePath := filepath.Join("..", "images", "simple.png")
	
	// Load image once for all iterations
	img, err := l.LoadImage(ctx, imagePath)
	if err != nil {
		b.Fatalf("Failed to load test image: %v", err)
	}
	
	// Get image info for reporting
	info, err := l.GetImageInfo(ctx, imagePath)
	if err != nil {
		b.Fatalf("Failed to get image info: %v", err)
	}
	
	b.Logf("Benchmarking large image: %dx%d (%.1f MP)",
		info.Width, info.Height, float64(info.PixelCount())/1000000)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := p.ProcessImage(img)
		if err != nil {
			b.Fatalf("Processing failed: %v", err)
		}
	}
}

// BenchmarkMemoryEfficiency runs a memory-focused benchmark
func BenchmarkMemoryEfficiency(b *testing.B) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	ctx := context.Background()

	// Use the largest test image available
	imagePath := filepath.Join("..", "images", "simple.png")

	img, err := l.LoadImage(ctx, imagePath)
	if err != nil {
		b.Fatalf("Failed to load test image: %v", err)
	}

	info, err := l.GetImageInfo(ctx, imagePath)
	if err != nil {
		b.Fatalf("Failed to get image info: %v", err)
	}

	b.Logf("Memory benchmark: %dx%d (%.1f MP) image processing",
		info.Width, info.Height, float64(info.PixelCount())/1000000)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		profile, err := p.ProcessImage(img)
		if err != nil {
			b.Fatalf("Processing failed: %v", err)
		}
		// Validate memory usage efficiency with simplified structure
		if profile.ColorCount == 0 {
			b.Fatal("No colors extracted - memory test invalid")
		}
	}
}