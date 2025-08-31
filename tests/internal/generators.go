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