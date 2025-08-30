package generative

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

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

// GenerateAndSaveComputationalImages creates and saves computationally generated aesthetic images.
// These images serve as proof of concept for the broader aesthetic generation system.
func GenerateAndSaveComputationalImages(outputDir string) error {
	fmt.Printf("Generating computationally generated images in %s...\n", outputDir)

	// Generate 80's Vector Graphics aesthetic
	fmt.Println("  Creating 80s-vector-graphics.png...")
	img80s := Generate80sVectorImage(800, 450) // 16:9 ratio for documentation
	if err := SaveImage(img80s, filepath.Join(outputDir, "80s-vector-graphics.png")); err != nil {
		return fmt.Errorf("failed to save 80's vector graphics image: %w", err)
	}

	// Generate Cassette Futurism aesthetic with orange accent
	fmt.Println("  Creating cassette-futurism.png...")
	imgCassette := GenerateCassetteFuturismImage(800, 450, 30.0/360.0) // Orange accent
	if err := SaveImage(imgCassette, filepath.Join(outputDir, "cassette-futurism.png")); err != nil {
		return fmt.Errorf("failed to save cassette futurism image: %w", err)
	}

	// Generate Complex Gradients (3 variations)
	gradientTypes := []string{"linear-smooth", "radial-complex", "stepped-harsh"}
	
	for _, gradientType := range gradientTypes {
		fmt.Printf("  Creating gradient-%s.png...\n", gradientType)
		imgGradient := GenerateComplexGradientImage(600, 400, gradientType)
		filename := fmt.Sprintf("gradient-%s.png", gradientType)
		if err := SaveImage(imgGradient, filepath.Join(outputDir, filename)); err != nil {
			return fmt.Errorf("failed to save %s gradient image: %w", gradientType, err)
		}
	}

	fmt.Println("  ✅ Computationally generated images created successfully!")
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
			
			srcColor := src.At(srcX+srcBounds.Min.X, srcY+srcBounds.Min.Y)
			dst.Set(x, y, srcColor)
		}
	}

	return dst
}