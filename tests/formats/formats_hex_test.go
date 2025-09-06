package formats_test

import (
	"image/color"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func TestToHex(t *testing.T) {
	testCases := []struct {
		name     string
		input    color.RGBA
		expected string
	}{
		{
			name:     "Black",
			input:    color.RGBA{R: 0, G: 0, B: 0, A: 255},
			expected: "#000000",
		},
		{
			name:     "White",
			input:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expected: "#FFFFFF",
		},
		{
			name:     "Red",
			input:    color.RGBA{R: 255, G: 0, B: 0, A: 255},
			expected: "#FF0000",
		},
		{
			name:     "Green",
			input:    color.RGBA{R: 0, G: 255, B: 0, A: 255},
			expected: "#00FF00",
		},
		{
			name:     "Blue",
			input:    color.RGBA{R: 0, G: 0, B: 255, A: 255},
			expected: "#0000FF",
		},
		{
			name:     "Orange",
			input:    color.RGBA{R: 255, G: 165, B: 0, A: 255},
			expected: "#FFA500",
		},
		{
			name:     "With Alpha (ignored)",
			input:    color.RGBA{R: 255, G: 128, B: 0, A: 128},
			expected: "#FF8000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.ToHex(tc.input)
			
			// Diagnostic logging
			t.Logf("Input RGBA: R=%d, G=%d, B=%d, A=%d", tc.input.R, tc.input.G, tc.input.B, tc.input.A)
			t.Logf("Expected hex: %s", tc.expected)
			t.Logf("Actual hex: %s", result)
			
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestToHexA(t *testing.T) {
	testCases := []struct {
		name     string
		input    color.RGBA
		expected string
	}{
		{
			name:     "Opaque Black",
			input:    color.RGBA{R: 0, G: 0, B: 0, A: 255},
			expected: "#000000FF",
		},
		{
			name:     "Opaque White",
			input:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expected: "#FFFFFFFF",
		},
		{
			name:     "Semi-transparent Red",
			input:    color.RGBA{R: 255, G: 0, B: 0, A: 128},
			expected: "#FF000080",
		},
		{
			name:     "Fully Transparent",
			input:    color.RGBA{R: 255, G: 255, B: 255, A: 0},
			expected: "#FFFFFF00",
		},
		{
			name:     "Quarter Opacity Blue",
			input:    color.RGBA{R: 0, G: 0, B: 255, A: 64},
			expected: "#0000FF40",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.ToHexA(tc.input)
			
			// Diagnostic logging
			t.Logf("Input RGBA: R=%d, G=%d, B=%d, A=%d", tc.input.R, tc.input.G, tc.input.B, tc.input.A)
			t.Logf("Expected HEXA: %s", tc.expected)
			t.Logf("Actual HEXA: %s", result)
			
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestParseHex(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expected  color.RGBA
		shouldErr bool
	}{
		// Valid 6-character hex
		{
			name:     "Black #000000",
			input:    "#000000",
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		},
		{
			name:     "White FFFFFF",
			input:    "FFFFFF",
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:     "Red #FF0000",
			input:    "#FF0000",
			expected: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		// Valid 3-character hex
		{
			name:     "Short Black #000",
			input:    "#000",
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		},
		{
			name:     "Short White FFF",
			input:    "FFF",
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:     "Short Red #F00",
			input:    "#F00",
			expected: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		// Valid 8-character hex with alpha
		{
			name:     "Semi-transparent Red #FF000080",
			input:    "#FF000080",
			expected: color.RGBA{R: 255, G: 0, B: 0, A: 128},
		},
		{
			name:     "Fully transparent FFFFFF00",
			input:    "FFFFFF00",
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 0},
		},
		// Valid 4-character hex with alpha
		{
			name:     "Short semi-transparent #F008",
			input:    "#F008",
			expected: color.RGBA{R: 255, G: 0, B: 0, A: 136},
		},
		// Lowercase
		{
			name:     "Lowercase hex #ff8800",
			input:    "#ff8800",
			expected: color.RGBA{R: 255, G: 136, B: 0, A: 255},
		},
		// Invalid cases
		{
			name:      "Invalid characters",
			input:     "#GGGGGG",
			shouldErr: true,
		},
		{
			name:      "Empty string",
			input:     "",
			shouldErr: true,
		},
		{
			name:      "Too long",
			input:     "#FF0000FF00",
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := formats.ParseHex(tc.input)

			// Diagnostic logging
			t.Logf("Input hex string: %s", tc.input)
			t.Logf("Should error: %t", tc.shouldErr)
			if err != nil {
				t.Logf("Parse error: %v", err)
			} else {
				t.Logf("Parsed RGBA: R=%d, G=%d, B=%d, A=%d", result.R, result.G, result.B, result.A)
				t.Logf("Expected RGBA: R=%d, G=%d, B=%d, A=%d", tc.expected.R, tc.expected.G, tc.expected.B, tc.expected.A)
			}

			if tc.shouldErr {
				if err == nil {
					t.Errorf("Expected error for input %s, but got none", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tc.expected {
					t.Errorf("Expected RGBA(%d, %d, %d, %d), got RGBA(%d, %d, %d, %d)",
						tc.expected.R, tc.expected.G, tc.expected.B, tc.expected.A,
						result.R, result.G, result.B, result.A)
				}
			}
		})
	}
}

func TestHSLAToHex(t *testing.T) {
	testCases := []struct {
		name     string
		input    formats.HSLA
		expected string
	}{
		{
			name:     "Red",
			input:    formats.NewHSLA(0, 1.0, 0.5, 1.0),
			expected: "#FF0000",
		},
		{
			name:     "Green",
			input:    formats.NewHSLA(120, 1.0, 0.5, 1.0),
			expected: "#00FF00",
		},
		{
			name:     "Blue",
			input:    formats.NewHSLA(240, 1.0, 0.5, 1.0),
			expected: "#0000FF",
		},
		{
			name:     "Gray",
			input:    formats.NewHSLA(0, 0.0, 0.5, 1.0),
			expected: "#808080",
		},
		{
			name:     "Semi-transparent (alpha ignored)",
			input:    formats.NewHSLA(0, 1.0, 0.5, 0.5),
			expected: "#FF0000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.HSLAToHex(tc.input)
			
			// Diagnostic logging
			t.Logf("Input HSLA: H=%.1f, S=%.3f, L=%.3f, A=%.3f", tc.input.H, tc.input.S, tc.input.L, tc.input.A)
			t.Logf("Expected hex: %s", tc.expected)
			t.Logf("Actual hex: %s", result)
			
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestHSLAToHexA(t *testing.T) {
	testCases := []struct {
		name     string
		input    formats.HSLA
		expected string
	}{
		{
			name:     "Opaque Red",
			input:    formats.NewHSLA(0, 1.0, 0.5, 1.0),
			expected: "#FF0000FF",
		},
		{
			name:     "Semi-transparent Blue",
			input:    formats.NewHSLA(240, 1.0, 0.5, 0.5),
			expected: "#0000FF80",
		},
		{
			name:     "Transparent Gray",
			input:    formats.NewHSLA(0, 0.0, 0.5, 0.0),
			expected: "#80808000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.HSLAToHexA(tc.input)
			
			// Diagnostic logging
			t.Logf("Input HSLA: H=%.1f, S=%.3f, L=%.3f, A=%.3f", tc.input.H, tc.input.S, tc.input.L, tc.input.A)
			t.Logf("Expected HEXA: %s", tc.expected)
			t.Logf("Actual HEXA: %s", result)
			
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestParseHexToHSLA(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:  "Valid hex #FF0000",
			input: "#FF0000",
		},
		{
			name:  "Valid short hex #F00",
			input: "#F00",
		},
		{
			name:      "Invalid hex",
			input:     "#GGGGGG",
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := formats.ParseHexToHSLA(tc.input)

			if tc.shouldErr {
				if err == nil {
					t.Errorf("Expected error, got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// Verify it's a valid HSLA
				if result.H < 0 || result.H >= 360 {
					t.Errorf("Invalid hue: %f", result.H)
				}
				if result.S < 0 || result.S > 1 {
					t.Errorf("Invalid saturation: %f", result.S)
				}
				if result.L < 0 || result.L > 1 {
					t.Errorf("Invalid lightness: %f", result.L)
				}
				if result.A < 0 || result.A > 1 {
					t.Errorf("Invalid alpha: %f", result.A)
				}
			}
		})
	}
}
