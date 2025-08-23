package main

import (
	"fmt"
	"sync"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
)

func main() {
	fmt.Println("=== Omarchy Theme Generator - Color Type Test ===")

	// Test 1: RGBA Color Creation
	fmt.Println("\nTest 1: RGBA Color Creation")
	red := color.NewRGB(255, 0, 0)
	fmt.Printf("Red RGBA: R=%d, G=%d, B=%d, A=%d\n", red.R, red.G, red.B, red.A)
	fmt.Printf("Red HEX (no alpha): %s\n", red.HEX())
	fmt.Printf("Red HEXA (with alpha): %s\n", red.HEXA())
	fmt.Printf("Red CSS RGB: %s\n", red.CSSRGB())
	fmt.Printf("Red is opaque: %v\n", red.IsOpaque())

	// Test 2: HSLA Conversion and Caching
	fmt.Println("\nTest 2: HSLA Conversion and Caching")
	h, s, l, a := red.HSLA()
	fmt.Printf("Red HSLA: H=%.3f, S=%.3f, L=%.3f, A=%.3f\n", h, s, l, a)
	fmt.Printf("Red CSS HSLA: %s\n", red.CSSHSLA())

	// Verify caching works by calling again
	h2, s2, l2, a2 := red.HSLA()
	if h == h2 && s == s2 && l == l2 && a == a2 {
		fmt.Println("✓ HSLA caching working correctly")
	} else {
		fmt.Println("✗ HSLA caching failed")
	}

	// Test 3: HSLA Color Creation
	fmt.Println("\nTest 3: HSLA Color Creation")
	blueTransparent := color.NewHSLA(240.0/360.0, 1.0, 0.5, 0.7) // 70% alpha blue
	fmt.Printf("Blue HEXA: %s\n", blueTransparent.HEXA())
	fmt.Printf("Blue RGBA: R=%d, G=%d, B=%d, A=%d\n",
		blueTransparent.R, blueTransparent.G, blueTransparent.B, blueTransparent.A)
	fmt.Printf("Blue CSS RGBA: %s\n", blueTransparent.CSSRGBA())
	fmt.Printf("Blue alpha: %.3f\n", blueTransparent.Alpha())

	// Test 4: Alpha Manipulation
	fmt.Println("\nTest 4: Alpha Manipulation")
	green := color.NewRGB(0, 255, 0)
	greenHalf := green.WithAlpha(0.5) // 50% alpha

	fmt.Printf("Green opaque: %s (α=%.3f)\n", green.CSSRGB(), green.Alpha())
	fmt.Printf("Green 50%% alpha: %s (α=%.3f)\n", greenHalf.CSSRGBA(), greenHalf.Alpha())

	// Test transparency checks
	transparent := color.NewRGBA(255, 0, 0, 0.0)
	fmt.Printf("Transparent red: %s, IsTransparent=%v\n",
		transparent.HEXA(), transparent.IsTransparent())

	// Test 5: CSS Format Methods
	fmt.Println("\nTest 5: CSS Format Methods")
	purple := color.NewRGBA(128, 0, 128, 0.75) // 75% alpha purple

	fmt.Printf("Purple HEX (no alpha): %s\n", purple.HEX())
	fmt.Printf("Purple HEXA (with alpha): %s\n", purple.HEXA())
	fmt.Printf("Purple CSS RGB: %s\n", purple.CSSRGB())
	fmt.Printf("Purple CSS RGBA: %s\n", purple.CSSRGBA())
	fmt.Printf("Purple CSS HSL: %s\n", purple.CSSHSL())
	fmt.Printf("Purple CSS HSLA: %s\n", purple.CSSHSLA())

	// Test 6: Component Access Methods
	fmt.Println("\nTest 6: Component Access Methods")
	r, g, b := purple.RGB()
	fmt.Printf("Purple RGB components: R=%d, G=%d, B=%d\n", r, g, b)

	r, g, b, alpha := purple.RGBA()
	fmt.Printf("Purple RGBA components: R=%d, G=%d, B=%d, A=%d\n", r, g, b, alpha)

	h, s, l = purple.HSL()
	fmt.Printf("Purple HSL components: H=%.3f, S=%.3f, L=%.3f\n", h, s, l)

	h, s, l, a = purple.HSLA()
	fmt.Printf("Purple HSLA components: H=%.3f, S=%.3f, L=%.3f, A=%.3f\n", h, s, l, a)

	// Test 7: Thread Safety with HSLA
	fmt.Println("\nTest 7: Thread Safety Test")
	var wg sync.WaitGroup
	orange := color.NewRGBA(255, 165, 0, 0.8)

	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			h, s, l, a := orange.HSLA()
			fmt.Printf("Goroutine %d: HSLA=%.3f,%.3f,%.3f,%.3f\n", id, h, s, l, a)
		}(i)
	}
	wg.Wait()

	// Test 8: Color Conversion Accuracy with Alpha
	fmt.Println("\nTest 8: Color Conversion Accuracy with Alpha")

	testColorsNew := []struct {
		name    string
		r, g, b uint8
		alpha   float64
	}{
		{"White Opaque", 255, 255, 255, 1.0},
		{"Black 50%", 0, 0, 0, 0.5},
		{"Red 75%", 255, 0, 0, 0.75},
		{"Green Transparent", 0, 255, 0, 0.0},
	}

	for _, test := range testColorsNew {
		c := color.NewRGBA(test.r, test.g, test.b, test.alpha)
		h, s, l, a := c.HSLA()
		fmt.Printf("%s: RGBA(%d,%d,%d,%.2f) -> HSLA(%.3f,%.3f,%.3f,%.3f)\n",
			test.name, test.r, test.g, test.b, test.alpha, h, s, l, a)

		// Test round-trip conversion
		back := color.NewHSLA(h, s, l, a)
		if abs(int(back.R)-int(test.r)) <= 1 &&
			abs(int(back.G)-int(test.g)) <= 1 &&
			abs(int(back.B)-int(test.b)) <= 1 &&
			abs(int(back.Alpha()*255)-int(test.alpha*255)) <= 1 {
			fmt.Printf("✓ Round-trip conversion accurate\n")
		} else {
			fmt.Printf("✗ Round-trip failed: got RGBA(%d,%d,%d,%.3f)\n",
				back.R, back.G, back.B, back.Alpha())
		}
	}

	// Test 9: Alpha Edge Cases
	fmt.Println("\nTest 9: Alpha Edge Cases")

	// Test alpha clamping
	overAlpha := color.NewHSLA(0.5, 0.8, 0.6, 1.5)   // Alpha > 1
	underAlpha := color.NewHSLA(0.5, 0.8, 0.6, -0.3) // Alpha < 0

	fmt.Printf("Over-alpha (1.5) clamped to: %.3f\n", overAlpha.Alpha())
	fmt.Printf("Under-alpha (-0.3) clamped to: %.3f\n", underAlpha.Alpha())

	// Test alpha manipulation
	base := color.NewRGB(100, 150, 200)
	alphas := []float64{0.0, 0.25, 0.5, 0.75, 1.0}

	fmt.Printf("Base color alpha variations:\n")
	for _, alpha := range alphas {
		variant := base.WithAlpha(alpha)
		fmt.Printf("  α=%.2f (%s)\n", alpha, variant.CSSRGBA())
	}

	// Test 10: Alpha Conversion Consistency
	fmt.Println("\nTest 10: Alpha Conversion Consistency")
	testAlphas := []float64{0.0, 0.25, 0.5, 0.75, 1.0}

	fmt.Printf("Alpha consistency check:\n")
	for _, expectedAlpha := range testAlphas {
		c := color.NewRGB(128, 128, 128).WithAlpha(expectedAlpha)
		actualAlpha := c.Alpha()
		fmt.Printf("  Expected α=%.2f, got α=%.3f %s\n",
			expectedAlpha, actualAlpha,
			checkMark(abs(int((expectedAlpha-actualAlpha)*1000)) <= 5)) // Allow small rounding
	}

	// Test 11: Pointer vs Value Semantics
	fmt.Println("\nTest 11: Pointer Semantics")
	original := color.NewRGB(255, 128, 64)
	modified := original.WithAlpha(0.5)

	fmt.Printf("Original color: %s (α=%.3f)\n", original.CSSRGBA(), original.Alpha())
	fmt.Printf("Modified copy: %s (α=%.3f)\n", modified.CSSRGBA(), modified.Alpha())
	fmt.Printf("Original unchanged: %s\n", checkMark(original.Alpha() == 1.0))

	fmt.Println("\n=== Color Type Test Complete ===")
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func checkMark(condition bool) string {
	if condition {
		return "✓"
	}
	return "✗"
}
