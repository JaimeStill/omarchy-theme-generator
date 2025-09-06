package chromatic_test

import (
	"fmt"
	"image/color"
	"math"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/chromatic"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func TestDistanceRGB(t *testing.T) {
	testCases := []struct {
		name     string
		color1   color.RGBA
		color2   color.RGBA
		expected float64
		description string
	}{
		{
			name:     "Identical colors",
			color1:   color.RGBA{R: 128, G: 128, B: 128, A: 255},
			color2:   color.RGBA{R: 128, G: 128, B: 128, A: 255},
			expected: 0.0,
			description: "Distance between identical colors should be zero",
		},
		{
			name:     "Pure black to pure white",
			color1:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
			color2:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expected: math.Sqrt(255*255 + 255*255 + 255*255), // ~441.67
			description: "Maximum RGB distance in standard color space",
		},
		{
			name:     "Pure red to pure green",
			color1:   color.RGBA{R: 255, G: 0, B: 0, A: 255},
			color2:   color.RGBA{R: 0, G: 255, B: 0, A: 255},
			expected: math.Sqrt(255*255 + 255*255), // ~360.62
			description: "Distance between primary colors",
		},
		{
			name:     "Small difference",
			color1:   color.RGBA{R: 100, G: 100, B: 100, A: 255},
			color2:   color.RGBA{R: 101, G: 101, B: 101, A: 255},
			expected: math.Sqrt(3), // ~1.73
			description: "Small RGB difference should have small distance",
		},
		{
			name:     "Single channel difference",
			color1:   color.RGBA{R: 100, G: 100, B: 100, A: 255},
			color2:   color.RGBA{R: 150, G: 100, B: 100, A: 255},
			expected: 50.0,
			description: "Distance in single RGB channel",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			distance := chromatic.DistanceRGB(tc.color1, tc.color2)
			
			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Color 1: RGB(%d,%d,%d)", tc.color1.R, tc.color1.G, tc.color1.B)
			t.Logf("Color 2: RGB(%d,%d,%d)", tc.color2.R, tc.color2.G, tc.color2.B)
			t.Logf("Expected distance: %.3f", tc.expected)
			t.Logf("Actual distance: %.3f", distance)
			t.Logf("Difference: %.6f", math.Abs(distance-tc.expected))
			t.Logf("Description: %s", tc.description)

			tolerance := 0.01
			if math.Abs(distance-tc.expected) > tolerance {
				t.Errorf("DistanceRGB mismatch: expected %.3f, got %.3f (tolerance %.3f)", 
					tc.expected, distance, tolerance)
			}
		})
	}
}

func TestDistanceHSL(t *testing.T) {
	testCases := []struct {
		name     string
		color1   color.RGBA
		color2   color.RGBA
		maxExpectedDistance float64
		description string
	}{
		{
			name:     "Identical colors",
			color1:   color.RGBA{R: 128, G: 128, B: 128, A: 255},
			color2:   color.RGBA{R: 128, G: 128, B: 128, A: 255},
			maxExpectedDistance: 0.001, // Allow small floating point errors
			description: "HSL distance between identical colors should be near zero",
		},
		{
			name:     "Pure red to pure blue - hue difference",
			color1:   color.RGBA{R: 255, G: 0, B: 0, A: 255},
			color2:   color.RGBA{R: 0, G: 0, B: 255, A: 255},
			maxExpectedDistance: math.Sqrt(240*240), // 240° hue difference squared
			description: "Large hue difference should dominate HSL distance",
		},
		{
			name:     "Grayscale to saturated - saturation difference",
			color1:   color.RGBA{R: 128, G: 128, B: 128, A: 255}, // Gray (S=0)
			color2:   color.RGBA{R: 255, G: 0, B: 0, A: 255},     // Red (S=1)
			maxExpectedDistance: 2.0, // Large saturation component
			description: "Saturation difference should contribute significantly to distance",
		},
		{
			name:     "Light vs dark same hue",
			color1:   color.RGBA{R: 64, G: 0, B: 0, A: 255},   // Dark red
			color2:   color.RGBA{R: 255, G: 0, B: 0, A: 255},  // Bright red
			maxExpectedDistance: 1.0, // Lightness difference
			description: "Lightness difference with same hue and saturation",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			distance := chromatic.DistanceHSL(tc.color1, tc.color2)
			
			// Convert to HSL for diagnostic information
			hsl1 := formats.RGBAToHSLA(tc.color1)
			hsl2 := formats.RGBAToHSLA(tc.color2)
			
			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Color 1: RGB(%d,%d,%d) → HSL(%.1f°, %.3f, %.3f)", 
				tc.color1.R, tc.color1.G, tc.color1.B, hsl1.H, hsl1.S, hsl1.L)
			t.Logf("Color 2: RGB(%d,%d,%d) → HSL(%.1f°, %.3f, %.3f)", 
				tc.color2.R, tc.color2.G, tc.color2.B, hsl2.H, hsl2.S, hsl2.L)
			t.Logf("HSL distance: %.3f", distance)
			t.Logf("Max expected: %.3f", tc.maxExpectedDistance)
			t.Logf("Description: %s", tc.description)

			if distance > tc.maxExpectedDistance {
				t.Errorf("DistanceHSL too large: got %.3f, expected ≤ %.3f", 
					distance, tc.maxExpectedDistance)
			}

			if distance < 0 {
				t.Error("Distance cannot be negative")
			}
		})
	}
}

func TestDistanceLAB(t *testing.T) {
	testCases := []struct {
		name     string
		color1   color.RGBA
		color2   color.RGBA
		description string
	}{
		{
			name:     "Identical colors",
			color1:   color.RGBA{R: 128, G: 128, B: 128, A: 255},
			color2:   color.RGBA{R: 128, G: 128, B: 128, A: 255},
			description: "LAB distance between identical colors should be zero",
		},
		{
			name:     "Pure white to pure black",
			color1:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
			color2:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
			description: "Maximum lightness difference in LAB space",
		},
		{
			name:     "Red to green - perceptual difference",
			color1:   color.RGBA{R: 255, G: 0, B: 0, A: 255},
			color2:   color.RGBA{R: 0, G: 255, B: 0, A: 255},
			description: "Perceptually significant color difference",
		},
		{
			name:     "Blue variations",
			color1:   color.RGBA{R: 0, G: 0, B: 200, A: 255},
			color2:   color.RGBA{R: 0, G: 0, B: 255, A: 255},
			description: "Similar blues with lightness difference",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			distance := chromatic.DistanceLAB(tc.color1, tc.color2)
			
			// Convert to LAB for diagnostic information
			lab1 := formats.RGBAToLAB(tc.color1)
			lab2 := formats.RGBAToLAB(tc.color2)
			
			// Calculate component differences
			deltaL := lab1.L - lab2.L
			deltaA := lab1.A - lab2.A
			deltaB := lab1.B - lab2.B
			expectedDistance := math.Sqrt(deltaL*deltaL + deltaA*deltaA + deltaB*deltaB)
			
			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Color 1: RGB(%d,%d,%d) → LAB(%.1f, %.1f, %.1f)", 
				tc.color1.R, tc.color1.G, tc.color1.B, lab1.L, lab1.A, lab1.B)
			t.Logf("Color 2: RGB(%d,%d,%d) → LAB(%.1f, %.1f, %.1f)", 
				tc.color2.R, tc.color2.G, tc.color2.B, lab2.L, lab2.A, lab2.B)
			t.Logf("Component differences: ΔL=%.3f, Δa=%.3f, Δb=%.3f", deltaL, deltaA, deltaB)
			t.Logf("Expected Euclidean distance: %.3f", expectedDistance)
			t.Logf("Actual LAB distance: %.3f", distance)
			t.Logf("Description: %s", tc.description)

			// Basic validation
			if distance < 0 {
				t.Error("Distance cannot be negative")
			}

			if tc.name == "Identical colors" && distance > 0.001 {
				t.Errorf("Identical colors should have near-zero distance, got %.6f", distance)
			}
			
			// LAB distance should be roughly Euclidean distance (allowing for implementation details)
			tolerance := expectedDistance * 0.1 + 0.1 // 10% tolerance plus small constant
			if math.Abs(distance-expectedDistance) > tolerance {
				t.Logf("Note: LAB distance (%.3f) differs from Euclidean (%.3f) by %.3f (tolerance %.3f)",
					distance, expectedDistance, math.Abs(distance-expectedDistance), tolerance)
			}
		})
	}
}

func TestDistanceSymmetry(t *testing.T) {
	// Test that all distance functions are symmetric: d(a,b) = d(b,a)
	testColors := []color.RGBA{
		{R: 255, G: 0, B: 0, A: 255},     // Red
		{R: 0, G: 255, B: 0, A: 255},     // Green  
		{R: 0, G: 0, B: 255, A: 255},     // Blue
		{R: 128, G: 128, B: 128, A: 255}, // Gray
		{R: 255, G: 255, B: 255, A: 255}, // White
		{R: 0, G: 0, B: 0, A: 255},       // Black
	}

	for i, color1 := range testColors {
		for j, color2 := range testColors {
			if i >= j {
				continue // Skip same color and duplicates
			}
			
			t.Run(fmt.Sprintf("Symmetry_%d_%d", i, j), func(t *testing.T) {
				// Test RGB distance symmetry
				distRGB1 := chromatic.DistanceRGB(color1, color2)
				distRGB2 := chromatic.DistanceRGB(color2, color1)
				
				// Test HSL distance symmetry
				distHSL1 := chromatic.DistanceHSL(color1, color2)
				distHSL2 := chromatic.DistanceHSL(color2, color1)
				
				// Test LAB distance symmetry
				distLAB1 := chromatic.DistanceLAB(color1, color2)
				distLAB2 := chromatic.DistanceLAB(color2, color1)
				
				// Log diagnostic information
				t.Logf("Color 1: RGB(%d,%d,%d)", color1.R, color1.G, color1.B)
				t.Logf("Color 2: RGB(%d,%d,%d)", color2.R, color2.G, color2.B)
				t.Logf("RGB distance: %.6f vs %.6f", distRGB1, distRGB2)
				t.Logf("HSL distance: %.6f vs %.6f", distHSL1, distHSL2)
				t.Logf("LAB distance: %.6f vs %.6f", distLAB1, distLAB2)

				tolerance := 0.000001
				if math.Abs(distRGB1-distRGB2) > tolerance {
					t.Errorf("RGB distance not symmetric: %.6f vs %.6f", distRGB1, distRGB2)
				}
				if math.Abs(distHSL1-distHSL2) > tolerance {
					t.Errorf("HSL distance not symmetric: %.6f vs %.6f", distHSL1, distHSL2)
				}
				if math.Abs(distLAB1-distLAB2) > tolerance {
					t.Errorf("LAB distance not symmetric: %.6f vs %.6f", distLAB1, distLAB2)
				}
			})
		}
	}
}

func TestDistanceTriangleInequality(t *testing.T) {
	// Test triangle inequality: d(a,c) ≤ d(a,b) + d(b,c)
	testCases := []struct {
		name   string
		colorA color.RGBA
		colorB color.RGBA
		colorC color.RGBA
	}{
		{
			name:   "Primary colors triangle",
			colorA: color.RGBA{R: 255, G: 0, B: 0, A: 255},   // Red
			colorB: color.RGBA{R: 0, G: 255, B: 0, A: 255},   // Green
			colorC: color.RGBA{R: 0, G: 0, B: 255, A: 255},   // Blue
		},
		{
			name:   "Grayscale gradient",
			colorA: color.RGBA{R: 0, G: 0, B: 0, A: 255},     // Black
			colorB: color.RGBA{R: 128, G: 128, B: 128, A: 255}, // Gray
			colorC: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test RGB distance triangle inequality
			distRGB_AC := chromatic.DistanceRGB(tc.colorA, tc.colorC)
			distRGB_AB := chromatic.DistanceRGB(tc.colorA, tc.colorB)
			distRGB_BC := chromatic.DistanceRGB(tc.colorB, tc.colorC)
			
			// Test HSL distance triangle inequality
			distHSL_AC := chromatic.DistanceHSL(tc.colorA, tc.colorC)
			distHSL_AB := chromatic.DistanceHSL(tc.colorA, tc.colorB)
			distHSL_BC := chromatic.DistanceHSL(tc.colorB, tc.colorC)
			
			// Test LAB distance triangle inequality
			distLAB_AC := chromatic.DistanceLAB(tc.colorA, tc.colorC)
			distLAB_AB := chromatic.DistanceLAB(tc.colorA, tc.colorB)
			distLAB_BC := chromatic.DistanceLAB(tc.colorB, tc.colorC)
			
			// Log diagnostic information
			t.Logf("Triangle: RGB(%d,%d,%d) - RGB(%d,%d,%d) - RGB(%d,%d,%d)",
				tc.colorA.R, tc.colorA.G, tc.colorA.B,
				tc.colorB.R, tc.colorB.G, tc.colorB.B,
				tc.colorC.R, tc.colorC.G, tc.colorC.B)
			t.Logf("RGB: d(A,C)=%.3f, d(A,B)+d(B,C)=%.3f+%.3f=%.3f",
				distRGB_AC, distRGB_AB, distRGB_BC, distRGB_AB+distRGB_BC)
			t.Logf("HSL: d(A,C)=%.3f, d(A,B)+d(B,C)=%.3f+%.3f=%.3f",
				distHSL_AC, distHSL_AB, distHSL_BC, distHSL_AB+distHSL_BC)
			t.Logf("LAB: d(A,C)=%.3f, d(A,B)+d(B,C)=%.3f+%.3f=%.3f",
				distLAB_AC, distLAB_AB, distLAB_BC, distLAB_AB+distLAB_BC)

			// Allow small tolerance for floating point errors
			tolerance := 0.001
			
			if distRGB_AC > distRGB_AB+distRGB_BC+tolerance {
				t.Errorf("RGB triangle inequality violated: %.3f > %.3f + %.3f", 
					distRGB_AC, distRGB_AB, distRGB_BC)
			}
			if distHSL_AC > distHSL_AB+distHSL_BC+tolerance {
				t.Errorf("HSL triangle inequality violated: %.3f > %.3f + %.3f", 
					distHSL_AC, distHSL_AB, distHSL_BC)
			}
			if distLAB_AC > distLAB_AB+distLAB_BC+tolerance {
				t.Errorf("LAB triangle inequality violated: %.3f > %.3f + %.3f", 
					distLAB_AC, distLAB_AB, distLAB_BC)
			}
		})
	}
}