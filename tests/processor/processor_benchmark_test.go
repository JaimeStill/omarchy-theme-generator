package processor_test

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
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := p.ProcessImage(img)
		if err != nil {
			b.Fatalf("Processing failed: %v", err)
		}
	}
}

// BenchmarkProcessImage_Medium benchmarks processing of medium images (2-8MP)
func BenchmarkProcessImage_Medium(b *testing.B) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	ctx := context.Background()
	
	// Use a medium test image
	imagePath := filepath.Join("..", "images", "abstract.jpeg")
	
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
	
	b.Logf("Benchmarking medium image: %dx%d (%.1f MP)",
		info.Width, info.Height, float64(info.PixelCount())/1000000)
	
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

// BenchmarkCategoryExtraction benchmarks the category extraction specifically
func BenchmarkCategoryExtraction(b *testing.B) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	ctx := context.Background()
	
	// Use a diverse test image
	imagePath := filepath.Join("..", "images", "nebula.jpeg")
	
	img, err := l.LoadImage(ctx, imagePath)
	if err != nil {
		b.Fatalf("Failed to load test image: %v", err)
	}
	
	b.Logf("Benchmarking category extraction with nebula.jpeg")
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		profile, err := p.ProcessImage(img)
		if err != nil {
			b.Fatalf("Processing failed: %v", err)
		}
		
		// Ensure categories were extracted
		if len(profile.Colors.Categories) == 0 {
			b.Fatalf("No categories extracted")
		}
	}
}

// BenchmarkColorSpaceConversions benchmarks color space conversion performance
func BenchmarkColorSpaceConversions(b *testing.B) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	l := loader.NewFileLoader(s)
	ctx := context.Background()
	
	// Use a colorful test image
	imagePath := filepath.Join("..", "images", "warm.jpeg")
	
	img, err := l.LoadImage(ctx, imagePath)
	if err != nil {
		b.Fatalf("Failed to load test image: %v", err)
	}
	
	b.Logf("Benchmarking color space conversions with warm.jpeg")
	
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
	
	b.Logf("Benchmarking memory efficiency with large image")
	
	b.ReportAllocs()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := p.ProcessImage(img)
		if err != nil {
			b.Fatalf("Processing failed: %v", err)
		}
	}
}