package main

import (
	"fmt"
	"math"
	"sync"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/tests"
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
	fmt.Printf("Input: Red RGB(255,0,0)\n")
	fmt.Printf("Converted HSLA: H=%.3f, S=%.3f, L=%.3f, A=%.3f\n", h, s, l, a)
	fmt.Printf("CSS Format: %s\n", red.CSSHSLA())

	// Verify caching works by calling again
	h2, s2, l2, a2 := red.HSLA()
	cachingOK := h == h2 && s == s2 && l == l2 && a == a2
	fmt.Printf("Caching test: First call=(%.3f,%.3f,%.3f,%.3f)\n", h, s, l, a)
	fmt.Printf("              Second call=(%.3f,%.3f,%.3f,%.3f)\n", h2, s2, l2, a2)
	fmt.Printf("Result: %s - Values %s\n",
		tests.CheckMark(cachingOK),
		map[bool]string{true: "identical (cached)", false: "different (not cached)"}[cachingOK])

	// Test 3: HSLA Color Creation
	fmt.Println("\nTest 3: HSLA Color Creation")
	fmt.Printf("Input: HSLA(240°, 100%%, 50%%, 70%%)\n")
	fmt.Printf("Creating: NewHSLA(%.3f, %.3f, %.3f, %.3f)\n", 240.0/360.0, 1.0, 0.5, 0.7)

	blueTransparent := color.NewHSLA(240.0/360.0, 1.0, 0.5, 0.7) // 70% alpha blue

	fmt.Printf("Result RGB: R=%d, G=%d, B=%d, A=%d\n",
		blueTransparent.R, blueTransparent.G, blueTransparent.B, blueTransparent.A)
	fmt.Printf("Hex format: %s\n", blueTransparent.HEXA())
	fmt.Printf("CSS format: %s\n", blueTransparent.CSSRGBA())
	fmt.Printf("Alpha value: %.3f (expected 0.7±0.01)\n", blueTransparent.Alpha())

	alphaOK := math.Abs(blueTransparent.Alpha()-0.7) < 0.01
	fmt.Printf("Conversion: %s - Alpha preserved correctly\n", tests.CheckMark(alphaOK))

	// Test 4: Alpha Manipulation
	fmt.Println("\nTest 4: Alpha Manipulation")
	green := color.NewRGB(0, 255, 0)
	fmt.Printf("Original: RGB(0,255,0), Alpha=%.3f\n", green.Alpha())

	greenHalf := green.WithAlpha(0.5) // 50% alpha
	fmt.Printf("Operation: WithAlpha(0.5)\n")
	fmt.Printf("Result: %s, Alpha=%.3f\n", greenHalf.CSSRGBA(), greenHalf.Alpha())

	// Verify original unchanged (immutability)
	immutableOK := green.Alpha() == 1.0
	fmt.Printf("Original unchanged: Alpha=%.3f %s\n", green.Alpha(), tests.CheckMark(immutableOK))

	// Test transparency checks
	transparent := color.NewRGBA(255, 0, 0, 0.0)
	opaque := color.NewRGBA(255, 0, 0, 1.0)
	fmt.Printf("\nTransparency checks:\n")
	fmt.Printf("  RGBA(255,0,0,0.0): IsTransparent=%v, IsOpaque=%v\n",
		transparent.IsTransparent(), transparent.IsOpaque())
	fmt.Printf("  RGBA(255,0,0,1.0): IsTransparent=%v, IsOpaque=%v\n",
		opaque.IsTransparent(), opaque.IsOpaque())

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
		fmt.Printf("\n%s test:\n", test.name)
		fmt.Printf("  Input: RGBA(%d,%d,%d,%.2f)\n", test.r, test.g, test.b, test.alpha)

		c := color.NewRGBA(test.r, test.g, test.b, test.alpha)
		h, s, l, a := c.HSLA()
		fmt.Printf("  → HSLA: (%.3f,%.3f,%.3f,%.3f)\n", h, s, l, a)

		// Test round-trip conversion
		back := color.NewHSLA(h, s, l, a)
		fmt.Printf("  → Back to RGB: (%d,%d,%d,%.3f)\n", back.R, back.G, back.B, back.Alpha())

		rDiff := int(math.Abs(float64(back.R) - float64(test.r)))
		gDiff := int(math.Abs(float64(back.G) - float64(test.g)))
		bDiff := int(math.Abs(float64(back.B) - float64(test.b)))
		aDiff := int(math.Abs(float64(back.Alpha()*255) - float64(test.alpha*255)))

		roundTripOK := rDiff <= 1 && gDiff <= 1 && bDiff <= 1 && aDiff <= 1
		fmt.Printf("  Differences: R=%d, G=%d, B=%d, A=%d (tolerance ≤1)\n", rDiff, gDiff, bDiff, aDiff)
		fmt.Printf("  Result: %s - Round-trip %s\n",
			tests.CheckMark(roundTripOK),
			map[bool]string{true: "accurate", false: "failed"}[roundTripOK])
	}

	// Test 9: Alpha Edge Cases
	fmt.Println("\nTest 9: Alpha Edge Cases")

	// Test alpha clamping
	fmt.Printf("Alpha clamping test:\n")
	overAlpha := color.NewHSLA(0.5, 0.8, 0.6, 1.5)   // Alpha > 1
	underAlpha := color.NewHSLA(0.5, 0.8, 0.6, -0.3) // Alpha < 0

	fmt.Printf("  Input alpha=1.5 → Result: %.3f (expected 1.0) %s\n",
		overAlpha.Alpha(), tests.CheckMark(overAlpha.Alpha() == 1.0))
	fmt.Printf("  Input alpha=-0.3 → Result: %.3f (expected 0.0) %s\n",
		underAlpha.Alpha(), tests.CheckMark(underAlpha.Alpha() == 0.0))

	// Test alpha manipulation
	fmt.Printf("\nAlpha variation test:\n")
	base := color.NewRGB(100, 150, 200)
	alphas := []float64{0.0, 0.25, 0.5, 0.75, 1.0}

	fmt.Printf("  Base color: RGB(100,150,200)\n")
	for _, alpha := range alphas {
		variant := base.WithAlpha(alpha)
		actualAlpha := variant.Alpha()
		fmt.Printf("  WithAlpha(%.2f): %s, actual=%.3f %s\n",
			alpha, variant.CSSRGBA(), actualAlpha,
			tests.CheckMark(math.Abs(actualAlpha-alpha) < 0.01))
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
			tests.CheckMark(int(math.Abs((expectedAlpha-actualAlpha)*1000)) <= 5)) // Allow small rounding
	}

	// Test 11: Pointer vs Value Semantics
	fmt.Println("\nTest 11: Immutability Test")
	original := color.NewRGB(255, 128, 64)
	fmt.Printf("Original color: %s (α=%.3f)\n", original.CSSRGBA(), original.Alpha())

	modified := original.WithAlpha(0.5)
	fmt.Printf("After WithAlpha(0.5):\n")
	fmt.Printf("  Original: %s (α=%.3f)\n", original.CSSRGBA(), original.Alpha())
	fmt.Printf("  Modified: %s (α=%.3f)\n", modified.CSSRGBA(), modified.Alpha())

	immutabilityOK := original.Alpha() == 1.0
	fmt.Printf("Result: %s - Original %s\n",
		tests.CheckMark(immutabilityOK),
		map[bool]string{true: "unchanged (immutable)", false: "modified (mutable)"}[immutabilityOK])

	// Test 12: HEXA Parsing
	fmt.Println("\nTest 12: HEXA Parsing")

	hexaTests := []struct {
		name        string
		input       string
		expectedR   uint8
		expectedG   uint8
		expectedB   uint8
		expectedA   uint8
		shouldError bool
	}{
		{"Standard HEXA with #", "#ff8000c0", 255, 128, 0, 192, false},
		{"Standard HEXA without #", "ff8000c0", 255, 128, 0, 192, false},
		{"Full opacity", "#00ff00ff", 0, 255, 0, 255, false},
		{"Full transparency", "#0000ff00", 0, 0, 255, 0, false},
		{"White semi-transparent", "#ffffff80", 255, 255, 255, 128, false},
		{"Black opaque", "#000000ff", 0, 0, 0, 255, false},
		{"Too short", "#ff800", 0, 0, 0, 0, true},
		{"Too long", "#ff8000c0ff", 0, 0, 0, 0, true},
		{"Invalid characters", "#gghhiijj", 0, 0, 0, 0, true},
	}

	for _, test := range hexaTests {
		fmt.Printf("  %s: %s\n", test.name, test.input)

		parsed, err := color.ParseHEXA(test.input)
		if test.shouldError {
			hasError := err != nil
			fmt.Printf("    Expected error: %s %s\n",
				tests.CheckMark(hasError),
				map[bool]string{true: "✓ Got expected error", false: "✗ Should have errored"}[hasError])
			if err != nil {
				fmt.Printf("    Error: %v\n", err)
			}
		} else {
			if err != nil {
				fmt.Printf("    ✗ Unexpected error: %v\n", err)
				continue
			}

			r, g, b, a := parsed.RGBA()
			rgbaMatch := r == test.expectedR && g == test.expectedG && b == test.expectedB && a == test.expectedA
			fmt.Printf("    Expected: R=%d, G=%d, B=%d, A=%d\n", test.expectedR, test.expectedG, test.expectedB, test.expectedA)
			fmt.Printf("    Got:      R=%d, G=%d, B=%d, A=%d\n", r, g, b, a)
			fmt.Printf("    Result:   %s RGBA values %s\n",
				tests.CheckMark(rgbaMatch),
				map[bool]string{true: "match", false: "don't match"}[rgbaMatch])

			// Test round-trip conversion
			roundTrip := parsed.HEXA()
			expected := test.input
			if len(expected) > 0 && expected[0] != '#' {
				expected = "#" + expected
			}
			roundTripMatch := roundTrip == expected
			fmt.Printf("    Round-trip: %s → %s %s\n",
				expected, roundTrip,
				tests.CheckMark(roundTripMatch))
		}
		fmt.Println()
	}

	// Test 13: HEX Parsing (6-digit)
	fmt.Println("Test 13: HEX Parsing (6-digit)")

	hexTests := []struct {
		name        string
		input       string
		expectedR   uint8
		expectedG   uint8
		expectedB   uint8
		shouldError bool
	}{
		{"Red", "#ff0000", 255, 0, 0, false},
		{"Green without #", "00ff00", 0, 255, 0, false},
		{"Blue", "#0000ff", 0, 0, 255, false},
		{"White", "#ffffff", 255, 255, 255, false},
		{"Black", "#000000", 0, 0, 0, false},
		{"Gray", "#808080", 128, 128, 128, false},
		{"Too short", "#ff00", 0, 0, 0, true},
		{"Too long", "#ff0000ff", 0, 0, 0, true},
		{"Invalid", "#gghhii", 0, 0, 0, true},
	}

	for _, test := range hexTests {
		fmt.Printf("  %s: %s\n", test.name, test.input)

		parsed, err := color.ParseHEX(test.input)
		if test.shouldError {
			hasError := err != nil
			fmt.Printf("    Expected error: %s %s\n",
				tests.CheckMark(hasError),
				map[bool]string{true: "✓ Got expected error", false: "✗ Should have errored"}[hasError])
		} else {
			if err != nil {
				fmt.Printf("    ✗ Unexpected error: %v\n", err)
				continue
			}

			r, g, b, a := parsed.RGBA()
			rgbaMatch := r == test.expectedR && g == test.expectedG && b == test.expectedB && a == 255
			fmt.Printf("    Expected: R=%d, G=%d, B=%d, A=255\n", test.expectedR, test.expectedG, test.expectedB)
			fmt.Printf("    Got:      R=%d, G=%d, B=%d, A=%d\n", r, g, b, a)
			fmt.Printf("    Result:   %s RGBA values %s\n",
				tests.CheckMark(rgbaMatch),
				map[bool]string{true: "match", false: "don't match"}[rgbaMatch])
		}
		fmt.Println()
	}

	fmt.Println("=== Color Type Test Complete ===")
}
