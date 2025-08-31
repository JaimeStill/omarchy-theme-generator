package main

import (
	"fmt"
	"github.com/JaimeStill/omarchy-theme-generator/tests/internal"
)

func main() {
	fmt.Println("Regenerating test samples with correct vocabulary...")

	// Generate individual samples
	samples := map[string]func() error{
		"grayscale": func() error {
			img := internal.GenerateGrayscaleTestImage(1920, 1080)
			return internal.SaveImage(img, "samples/grayscale.png")
		},
		"monochromatic": func() error {
			img := internal.GenerateMonochromaticTestImage(1920, 1080)
			return internal.SaveImage(img, "samples/monochromatic.png")
		},
		"grayscale-small": func() error {
			img := internal.GenerateGrayscaleTestImage(400, 300)
			return internal.SaveImage(img, "samples/grayscale-small.png")
		},
		"monochromatic-small": func() error {
			img := internal.GenerateMonochromaticTestImage(400, 300)
			return internal.SaveImage(img, "samples/monochromatic-small.png")
		},
	}

	for name, generator := range samples {
		fmt.Printf("  Generating %s.png...\n", name)
		if err := generator(); err != nil {
			fmt.Printf("Error generating %s: %v\n", name, err)
			return
		}
	}

	fmt.Println("Test samples regenerated successfully!")
}
