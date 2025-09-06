package formats_test

import (
	"image/color"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

// TestToHex tests the ToHex function that converts RGBA to #RRGGBB format
func TestToHex(t *testing.T) {
	testCases := []struct {
		name     string
		rgba     color.RGBA
		expected string
		description string
	}{
		{
			name:     "Pure white",
			rgba:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expected: "#FFFFFF",
			description: "White should convert to #FFFFFF",
		},
		{
			name:     "Pure black",
			rgba:     color.RGBA{R: 0, G: 0, B: 0, A: 255},
			expected: "#000000",
			description: "Black should convert to #000000",
		},
		{
			name:     "Pure red",
			rgba:     color.RGBA{R: 255, G: 0, B: 0, A: 255},
			expected: "#FF0000",
			description: "Pure red should convert to #FF0000",
		},
		{
			name:     "Pure green",
			rgba:     color.RGBA{R: 0, G: 255, B: 0, A: 255},
			expected: "#00FF00",
			description: "Pure green should convert to #00FF00",
		},
		{
			name:     "Pure blue",
			rgba:     color.RGBA{R: 0, G: 0, B: 255, A: 255},
			expected: "#0000FF",
			description: "Pure blue should convert to #0000FF",
		},
		{
			name:     "Mixed color",
			rgba:     color.RGBA{R: 128, G: 64, B: 192, A: 255},
			expected: "#8040C0",
			description: "Mixed RGB values should convert correctly",
		},
		{
			name:     "Alpha ignored",
			rgba:     color.RGBA{R: 128, G: 64, B: 192, A: 0},
			expected: "#8040C0",
			description: "Alpha channel should be ignored in ToHex",
		},
		{
			name:     "Low values",
			rgba:     color.RGBA{R: 1, G: 2, B: 3, A: 255},
			expected: "#010203",
			description: "Low values should format with leading zeros",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.ToHex(tc.rgba)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("RGBA input: R=%d, G=%d, B=%d, A=%d", tc.rgba.R, tc.rgba.G, tc.rgba.B, tc.rgba.A)
			t.Logf("Hex result: %s", result)
			t.Logf("Expected result: %s", tc.expected)
			t.Logf("Description: %s", tc.description)

			// Verify the hex conversion
			if result != tc.expected {
				t.Errorf("ToHex mismatch:\nExpected: %s\nGot: %s", tc.expected, result)
			}

			t.Logf("✅ ToHex working correctly")
		})
	}
}

// TestToHexA tests the ToHexA function that converts RGBA to #RRGGBBAA format
func TestToHexA(t *testing.T) {
	testCases := []struct {
		name     string
		rgba     color.RGBA
		expected string
		description string
	}{
		{
			name:     "White opaque",
			rgba:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expected: "#FFFFFFFF",
			description: "Opaque white should include FF alpha",
		},
		{
			name:     "Black transparent",
			rgba:     color.RGBA{R: 0, G: 0, B: 0, A: 0},
			expected: "#00000000",
			description: "Transparent black should include 00 alpha",
		},
		{
			name:     "Red semi-transparent",
			rgba:     color.RGBA{R: 255, G: 0, B: 0, A: 128},
			expected: "#FF000080",
			description: "Semi-transparent red should include alpha value",
		},
		{
			name:     "Mixed with alpha",
			rgba:     color.RGBA{R: 128, G: 64, B: 192, A: 200},
			expected: "#8040C0C8",
			description: "Mixed color with alpha should format correctly",
		},
		{
			name:     "Low alpha value",
			rgba:     color.RGBA{R: 255, G: 128, B: 64, A: 1},
			expected: "#FF804001",
			description: "Low alpha values should format with leading zero",
		},
		{
			name:     "High alpha value",
			rgba:     color.RGBA{R: 64, G: 128, B: 255, A: 254},
			expected: "#4080FFFE",
			description: "High alpha values should format correctly",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formats.ToHexA(tc.rgba)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("RGBA input: R=%d, G=%d, B=%d, A=%d", tc.rgba.R, tc.rgba.G, tc.rgba.B, tc.rgba.A)
			t.Logf("HexA result: %s", result)
			t.Logf("Expected result: %s", tc.expected)
			t.Logf("Description: %s", tc.description)

			// Verify the hex conversion with alpha
			if result != tc.expected {
				t.Errorf("ToHexA mismatch:\nExpected: %s\nGot: %s", tc.expected, result)
			}

			t.Logf("✅ ToHexA working correctly")
		})
	}
}

// TestParseHex tests the ParseHex function that parses hex strings to RGBA
func TestParseHex(t *testing.T) {
	testCases := []struct {
		name     string
		hex      string
		expected color.RGBA
		hasError bool
		description string
	}{
		{
			name:     "3-digit hex with #",
			hex:      "#FFF",
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			hasError: false,
			description: "#FFF should expand to white",
		},
		{
			name:     "3-digit hex without #",
			hex:      "000",
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 255},
			hasError: false,
			description: "000 should expand to black",
		},
		{
			name:     "4-digit hex with alpha",
			hex:      "#F0A8",
			expected: color.RGBA{R: 255, G: 0, B: 170, A: 136},
			hasError: false,
			description: "#F0A8 should expand with alpha",
		},
		{
			name:     "6-digit hex",
			hex:      "#FF8040",
			expected: color.RGBA{R: 255, G: 128, B: 64, A: 255},
			hasError: false,
			description: "6-digit hex should parse with default alpha",
		},
		{
			name:     "8-digit hex with alpha",
			hex:      "#FF804080",
			expected: color.RGBA{R: 255, G: 128, B: 64, A: 128},
			hasError: false,
			description: "8-digit hex should parse with explicit alpha",
		},
		{
			name:     "Lowercase hex",
			hex:      "#ff8040",
			expected: color.RGBA{R: 255, G: 128, B: 64, A: 255},
			hasError: false,
			description: "Lowercase hex should parse correctly",
		},
		{
			name:     "Mixed case hex",
			hex:      "#Ff80aB",
			expected: color.RGBA{R: 255, G: 128, B: 171, A: 255},
			hasError: false,
			description: "Mixed case hex should parse correctly",
		},
		{
			name:     "Invalid character",
			hex:      "#GGFFFF",
			expected: color.RGBA{},
			hasError: true,
			description: "Invalid hex character should return error",
		},
		{
			name:     "Invalid length",
			hex:      "#FFFFF",
			expected: color.RGBA{},
			hasError: true,
			description: "5-character hex should return error",
		},
		{
			name:     "Empty string",
			hex:      "",
			expected: color.RGBA{},
			hasError: true,
			description: "Empty hex string should return error",
		},
		{
			name:     "Just hash",
			hex:      "#",
			expected: color.RGBA{},
			hasError: true,
			description: "Just hash symbol should return error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := formats.ParseHex(tc.hex)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Hex input: '%s'", tc.hex)
			t.Logf("RGBA result: R=%d, G=%d, B=%d, A=%d", result.R, result.G, result.B, result.A)
			t.Logf("Expected RGBA: R=%d, G=%d, B=%d, A=%d", tc.expected.R, tc.expected.G, tc.expected.B, tc.expected.A)
			t.Logf("Error occurred: %t", err != nil)
			t.Logf("Expected error: %t", tc.hasError)
			t.Logf("Description: %s", tc.description)

			// Check error expectation
			if tc.hasError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// If no error expected, check the result
			if !tc.hasError && err == nil {
				if result != tc.expected {
					t.Errorf("ParseHex result mismatch:\nExpected: %+v\nGot: %+v", tc.expected, result)
				}
			}

			if err != nil {
				t.Logf("Error details: %v", err)
			}

			t.Logf("✅ ParseHex working correctly")
		})
	}
}

// TestParseHexToHSLA tests the ParseHexToHSLA function that parses hex strings to HSLA
func TestParseHexToHSLA(t *testing.T) {
	testCases := []struct {
		name     string
		hex      string
		hasError bool
		description string
	}{
		{
			name:     "White hex to HSLA",
			hex:      "#FFFFFF",
			hasError: false,
			description: "White should convert to HSLA correctly",
		},
		{
			name:     "Black hex to HSLA",
			hex:      "#000000",
			hasError: false,
			description: "Black should convert to HSLA correctly",
		},
		{
			name:     "Red hex to HSLA",
			hex:      "#FF0000",
			hasError: false,
			description: "Pure red should convert to HSLA correctly",
		},
		{
			name:     "Mixed color hex to HSLA",
			hex:      "#8040C0",
			hasError: false,
			description: "Mixed color should convert to HSLA correctly",
		},
		{
			name:     "Short hex to HSLA",
			hex:      "#F0A",
			hasError: false,
			description: "3-digit hex should expand and convert to HSLA",
		},
		{
			name:     "Invalid hex to HSLA",
			hex:      "#GGGGGG",
			hasError: true,
			description: "Invalid hex should return error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := formats.ParseHexToHSLA(tc.hex)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Hex input: '%s'", tc.hex)
			t.Logf("HSLA result: H=%.1f°, S=%.3f, L=%.3f, A=%.3f", result.H, result.S, result.L, result.A)
			t.Logf("Error occurred: %t", err != nil)
			t.Logf("Expected error: %t", tc.hasError)
			t.Logf("Description: %s", tc.description)

			// Check error expectation
			if tc.hasError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// If no error expected, verify HSLA values are in valid ranges
			if !tc.hasError && err == nil {
				if result.H < 0 || result.H >= 360 {
					t.Errorf("H component out of range: %.1f (expected 0-359.9)", result.H)
				}
				if result.S < 0 || result.S > 1 {
					t.Errorf("S component out of range: %.3f (expected 0-1)", result.S)
				}
				if result.L < 0 || result.L > 1 {
					t.Errorf("L component out of range: %.3f (expected 0-1)", result.L)
				}
				if result.A < 0 || result.A > 1 {
					t.Errorf("A component out of range: %.3f (expected 0-1)", result.A)
				}
			}

			if err != nil {
				t.Logf("Error details: %v", err)
			}

			t.Logf("✅ ParseHexToHSLA working correctly")
		})
	}
}