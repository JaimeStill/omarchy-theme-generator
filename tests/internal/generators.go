package internal

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

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

// GenerateGrayscaleTestImage creates a grayscale image for testing synthesis edge cases.
// All pixels have equal R, G, B values (no color information).
func GenerateGrayscaleTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Create grayscale gradient with subtle variations
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Base grayscale value
			gray := uint8((x + y) * 255 / (width + height))

			// Add subtle noise to create some unique colors but keep it grayscale
			if (x*y)%17 == 0 {
				gray = uint8((int(gray) + 10) % 256)
			}

			img.Set(x, y, color.RGBA{gray, gray, gray, 255})
		}
	}

	return img
}

// GenerateMonochromaticTestImage creates a monochromatic image with single hue variations.
// Uses blue hue (240°) with varying lightness and saturation for proper monochromatic testing.
func GenerateMonochromaticTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Base hue: 240° (blue), varying lightness and saturation
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Calculate base lightness (0.2 to 0.8 range)
			lightness := 0.2 + (float64(x)/float64(width))*0.6

			// Calculate saturation (0.3 to 1.0 range)
			saturation := 0.3 + (float64(y)/float64(height))*0.7

			// Add some noise for variation while keeping within hue tolerance
			if (x+y)%15 == 0 {
				lightness = lightness * 0.8 // Darken some pixels
			}
			if (x*y)%19 == 0 {
				saturation = saturation * 0.7 // Desaturate some pixels
			}

			// Convert HSL(240°, saturation, lightness) to RGB
			// Hue = 240° = 2/3 in normalized range
			h := 2.0 / 3.0 // Blue hue (240°)
			s := saturation
			l := lightness

			// HSL to RGB conversion
			var r, g, b float64
			if s == 0 {
				r, g, b = l, l, l // Grayscale
			} else {
				hue2rgb := func(p, q, t float64) float64 {
					if t < 0 {
						t += 1
					}
					if t > 1 {
						t -= 1
					}
					if t < 1.0/6.0 {
						return p + (q-p)*6*t
					}
					if t < 1.0/2.0 {
						return q
					}
					if t < 2.0/3.0 {
						return p + (q-p)*(2.0/3.0-t)*6
					}
					return p
				}

				var q float64
				if l < 0.5 {
					q = l * (1 + s)
				} else {
					q = l + s - l*s
				}
				p := 2*l - q

				r = hue2rgb(p, q, h+1.0/3.0)
				g = hue2rgb(p, q, h)
				b = hue2rgb(p, q, h-1.0/3.0)
			}

			// Convert to 0-255 range
			rVal := uint8(r * 255)
			gVal := uint8(g * 255)
			bVal := uint8(b * 255)

			img.Set(x, y, color.RGBA{rVal, gVal, bVal, 255})
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

// GetOrGenerateTestImage checks if a test image exists in tests/samples/ and generates it if not.
// This ensures consistent test images across all test runs.
func GetOrGenerateTestImage(name string, generator func() image.Image) (image.Image, error) {
	samplesDir := "tests/samples"
	imagePath := filepath.Join(samplesDir, name+".png")

	// Check if image already exists
	if _, err := os.Stat(imagePath); err == nil {
		// Load existing image
		file, err := os.Open(imagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open existing test image: %w", err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decode test image: %w", err)
		}
		return img, nil
	}

	// Generate new image
	img := generator()

	// Save for future use
	if err := SaveImage(img, imagePath); err != nil {
		return nil, fmt.Errorf("failed to save test image: %w", err)
	}

	return img, nil
}

// resampleImage creates a simple resampled version of an image (nearest neighbor).
func ResampleImage(src image.Image, newWidth, newHeight int) image.Image {
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

// GenerateAllTestSamples generates and saves all test sample images to tests/samples/.
// This creates a consistent set of test images for validation and documentation.
func GenerateAllTestSamples() error {
	fmt.Println("Generating test sample images...")

	samples := []struct {
		name      string
		generator func() image.Image
	}{
		{"4k-synthetic", Generate4KTestImage},
		{"grayscale", func() image.Image { return GenerateGrayscaleTestImage(1920, 1080) }},
		{"monochromatic", func() image.Image { return GenerateMonochromaticTestImage(1920, 1080) }},
		{"high-contrast", func() image.Image { return GenerateHighContrastTestImage(1920, 1080) }},
		{"grayscale-small", func() image.Image { return GenerateGrayscaleTestImage(400, 300) }},
		{"monochromatic-small", func() image.Image { return GenerateMonochromaticTestImage(400, 300) }},
		{"high-contrast-small", func() image.Image { return GenerateHighContrastTestImage(400, 300) }},
	}

	for _, sample := range samples {
		fmt.Printf("  Generating %s.png...\n", sample.name)
		img := sample.generator()

		path := fmt.Sprintf("tests/samples/%s.png", sample.name)
		if err := SaveImage(img, path); err != nil {
			return fmt.Errorf("failed to save %s: %w", sample.name, err)
		}
	}

	fmt.Printf("Successfully generated %d test sample images\n", len(samples))
	return nil
}
