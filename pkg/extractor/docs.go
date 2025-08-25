// Package extractor provides high-performance image color extraction and analysis functionality
// for the Omarchy Theme Generator.
//
// This package handles the complete pipeline from image loading to color frequency analysis,
// with optimized processing for different image types and comprehensive performance benchmarking.
// It replaces traditional validation-based approaches with intelligent analysis that guides
// synthesis strategies for edge cases.
//
// # Core Components
//
// The package is organized into several focused modules:
//
//   - Loading (loader.go): Image file loading with format validation and error handling
//   - Frequency Analysis (frequency.go): Efficient color counting with optimized pixel access
//   - Extraction Pipeline (extractor.go): Complete extraction workflow with analysis
//   - Performance Testing (performance.go): Benchmarking and test image generation
//
// # Usage Patterns
//
// Basic extraction from file:
//
//	result, err := extractor.ExtractColors("image.jpg", nil)
//	if err != nil {
//	    return err
//	}
//	
//	analysis := result.AnalyzeForThemeGeneration()
//	fmt.Printf("Strategy: %s", analysis.SuggestedStrategy)
//
// Custom extraction options:
//
//	options := &extractor.ExtractionOptions{
//	    TopColorCount: 20,
//	    MinThreshold: 0.5,
//	    MaxImageDimension: 4096,
//	}
//	result, err := extractor.ExtractColors("image.jpg", options)
//
// Performance benchmarking:
//
//	img := extractor.Generate4KTestImage()
//	benchmark, result, err := extractor.BenchmarkExtraction(img, nil)
//	fmt.Printf("Processing time: %v", benchmark.Duration)
//
// # Performance Characteristics
//
// The extractor is optimized for high-resolution images with the following targets:
//   - 4K images (3840x2160): < 2 seconds processing time
//   - Memory usage: < 100MB peak during extraction
//   - Processing rate: > 30M pixels/second on modern hardware
//
// Performance optimizations include:
//   - Type-specific pixel access for RGBA and NRGBA images
//   - Packed RGB keys for efficient frequency mapping
//   - Pre-allocated data structures based on image characteristics
//   - Cache-friendly sequential pixel iteration patterns
//
// # Analysis vs Validation
//
// Unlike traditional validation approaches that fail on edge cases, this package
// provides comprehensive analysis that guides synthesis strategies:
//
//   - "extract": Sufficient color diversity for direct extraction
//   - "hybrid": Some extraction possible, synthesis needed for completeness
//   - "synthesize": Minimal color information, requires color theory generation
//
// # Edge Case Handling
//
// The package gracefully handles challenging image types:
//   - Grayscale images: Detected via saturation analysis
//   - High-contrast images: Analyzed for dominance patterns
//   - Monochromatic images: Distinguished from grayscale
//   - Large images: Processed efficiently within memory constraints
//
// # Integration with Synthesis
//
// Analysis results provide essential information for Session 4's synthesis implementation:
//   - Primary non-grayscale color detection for synthesis seeds
//   - Color distribution metrics for strategy selection  
//   - Dominance analysis for hybrid approaches
//   - Performance benchmarking for synthesis target validation
package extractor