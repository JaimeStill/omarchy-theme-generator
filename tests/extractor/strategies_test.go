package extractor_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/extractor"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func TestStrategySelection(t *testing.T) {
	testCases := []struct {
		image            string
		expectedStrategy string
		description      string
	}{
		{"nebula.jpeg", "saliency", "Complex space image should use saliency"},
		{"night-city.jpeg", "saliency", "High-detail urban scene should use saliency"},
		{"grayscale.jpeg", "frequency", "Grayscale image should use frequency"},
		{"mountains.jpeg", "", "Natural landscape strategy varies by characteristics"},
		{"abstract.jpeg", "", "Abstract art strategy varies by characteristics"},
	}

	options := extractor.DefaultOptions()
	options.TopColorCount = 5

	for _, tc := range testCases {
		t.Run(tc.image, func(t *testing.T) {
			imagePath := filepath.Join("..", "images", tc.image)

			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				t.Skipf("Image %s not found, skipping test", tc.image)
				return
			}

			result, err := extractor.ExtractColors(imagePath, options)
			if err != nil {
				t.Fatalf("Failed to extract colors from %s: %v", tc.image, err)
			}

			if result.SelectedStrategy == "" {
				t.Errorf("No strategy was selected for %s", tc.image)
			}

			if tc.expectedStrategy != "" && result.SelectedStrategy != tc.expectedStrategy {
				t.Errorf("Expected strategy %s for %s, got %s",
					tc.expectedStrategy, tc.image, result.SelectedStrategy)
			}

			if len(result.TopColors) == 0 {
				t.Errorf("No colors extracted from %s", tc.image)
			}

			t.Logf("%s: strategy=%s, colors=%d, dominant=%s (%.1f%%)",
				tc.image,
				result.SelectedStrategy,
				result.UniqueColors,
				formats.ToHex(result.DominantColor),
				result.TopColors[0].Percentage)
		})
	}
}

func TestThemeGenerationAnalysis(t *testing.T) {
	testImages := []string{"nebula.jpeg", "grayscale.jpeg", "mountains.jpeg"}
	options := extractor.DefaultOptions()

	for _, image := range testImages {
		t.Run(image, func(t *testing.T) {
			imagePath := filepath.Join("..", "images", image)

			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				t.Skipf("Image %s not found, skipping test", image)
				return
			}

			result, err := extractor.ExtractColors(imagePath, options)
			if err != nil {
				t.Fatalf("Failed to extract colors: %v", err)
			}

			analysis := result.AnalyzeForThemeGeneration()
			if analysis == nil {
				t.Fatal("Theme analysis returned nil")
			}

			if analysis.SuggestedStrategy == "" {
				t.Error("No suggested strategy in analysis")
			}

			// Validate analysis consistency
			if analysis.IsGrayscale && analysis.AverageSaturation > 0.1 {
				t.Errorf("Image marked as grayscale but has saturation %.3f",
					analysis.AverageSaturation)
			}

			t.Logf("%s: strategy=%s, grayscale=%v, monochromatic=%v, sat=%.3f",
				image,
				analysis.SuggestedStrategy,
				analysis.IsGrayscale,
				analysis.IsMonochromatic,
				analysis.AverageSaturation)
		})
	}
}

func TestSaliencyVsFrequency(t *testing.T) {
	// Direct comparison test for nebula problem
	nebulaPath := filepath.Join("..", "images", "nebula.jpeg")

	if _, err := os.Stat(nebulaPath); os.IsNotExist(err) {
		t.Skip("nebula.jpeg not found for comparison test")
	}

	options := extractor.DefaultOptions()
	options.TopColorCount = 5

	result, err := extractor.ExtractColors(nebulaPath, options)
	if err != nil {
		t.Fatalf("Failed to extract colors: %v", err)
	}

	// Nebula should use saliency strategy
	if result.SelectedStrategy != "saliency" {
		t.Errorf("Nebula image should use saliency strategy, got %s", result.SelectedStrategy)
	}

	// Check that we're getting saturated colors (not pure black/gray)
	hasColorful := false
	for _, cf := range result.TopColors {
		hsla := formats.RGBAToHSLA(cf.Color)
		s, l := hsla.S, hsla.L
		// Look for saturated colors with some lightness (nebula has dark but colorful regions)
		if s > 0.3 && l > 0.05 {
			hasColorful = true
			break
		}
	}

	if !hasColorful {
		t.Error("Saliency strategy should extract colorful regions from nebula, not just pure black/gray")
	}
}

func BenchmarkSaliencyStrategy(b *testing.B) {
	imagePath := filepath.Join("..", "images", "nebula.jpeg")
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		b.Skip("nebula.jpeg not found for benchmark")
	}

	options := extractor.DefaultOptions()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := extractor.ExtractColors(imagePath, options)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkFrequencyStrategy(b *testing.B) {
	imagePath := filepath.Join("..", "images", "grayscale.jpeg")
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		b.Skip("grayscale.jpeg not found for benchmark")
	}

	options := extractor.DefaultOptions()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := extractor.ExtractColors(imagePath, options)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}
