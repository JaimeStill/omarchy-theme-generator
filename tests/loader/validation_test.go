package loader_test

import (
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/loader"
)

func TestValidateImageFormat(t *testing.T) {
	supportedFormats := []string{"jpeg", "jpg", "png", "webp"}

	testCases := []struct {
		name           string
		path           string
		expectedError  bool
		errorContains  string
	}{
		{
			name:          "Valid JPEG extension",
			path:          "/path/to/image.jpg",
			expectedError: false,
		},
		{
			name:          "Valid JPEG extension (alternative)",
			path:          "/path/to/image.jpeg",
			expectedError: false,
		},
		{
			name:          "Valid PNG extension",
			path:          "/path/to/image.png",
			expectedError: false,
		},
		{
			name:          "Valid WebP extension",
			path:          "/path/to/image.webp",
			expectedError: false,
		},
		{
			name:          "Invalid extension",
			path:          "/path/to/image.gif",
			expectedError: true,
			errorContains: "unsupported format",
		},
		{
			name:          "No extension",
			path:          "/path/to/image",
			expectedError: true,
			errorContains: "no file extension",
		},
		{
			name:          "Empty path",
			path:          "",
			expectedError: true,
			errorContains: "no file extension",
		},
		{
			name:          "Case sensitivity test",
			path:          "/path/to/image.JPG",
			expectedError: false,
		},
		{
			name:          "Multiple dots in filename",
			path:          "/path/to/my.image.file.png",
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := loader.ValidateImageFormat(tc.path, supportedFormats)

			// Comprehensive diagnostic logging
			t.Logf("Path: %s", tc.path)
			t.Logf("Supported formats: %v", supportedFormats)
			t.Logf("Expected error: %t", tc.expectedError)
			if err != nil {
				t.Logf("Actual error: %s", err.Error())
			} else {
				t.Logf("Actual error: nil")
			}

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error for path '%s', got nil", tc.path)
				} else if tc.errorContains != "" && !containsIgnoreCase(err.Error(), tc.errorContains) {
					t.Errorf("Expected error containing '%s', got '%s'", tc.errorContains, err.Error())
				} else {
					t.Logf("✓ Correctly returned expected error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for path '%s', got: %v", tc.path, err)
				} else {
					t.Logf("✓ Correctly accepted valid format")
				}
			}
		})
	}
}

func TestValidateImageFormatName(t *testing.T) {
	supportedFormats := []string{"jpeg", "jpg", "png", "webp"}

	testCases := []struct {
		name           string
		formatName     string
		expectedError  bool
		errorContains  string
	}{
		{
			name:          "Valid format - jpg",
			formatName:    "jpg",
			expectedError: false,
		},
		{
			name:          "Valid format - jpeg",
			formatName:    "jpeg",
			expectedError: false,
		},
		{
			name:          "Valid format - png",
			formatName:    "png",
			expectedError: false,
		},
		{
			name:          "Valid format - webp",
			formatName:    "webp",
			expectedError: false,
		},
		{
			name:          "Invalid format",
			formatName:    "gif",
			expectedError: true,
			errorContains: "unsupported format",
		},
		{
			name:          "Empty format",
			formatName:    "",
			expectedError: true,
			errorContains: "no file extension",
		},
		{
			name:          "Case sensitivity - uppercase",
			formatName:    "PNG",
			expectedError: true, // Implementation is case sensitive
			errorContains: "unsupported format",
		},
		{
			name:          "Case sensitivity - mixed",
			formatName:    "JpEg",
			expectedError: true, // Implementation is case sensitive
			errorContains: "unsupported format",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := loader.ValidateImageFormatName(tc.formatName, supportedFormats)

			// Comprehensive diagnostic logging
			t.Logf("Format name: '%s'", tc.formatName)
			t.Logf("Supported formats: %v", supportedFormats)
			t.Logf("Expected error: %t", tc.expectedError)
			if err != nil {
				t.Logf("Actual error: %s", err.Error())
			} else {
				t.Logf("Actual error: nil")
			}

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error for format '%s', got nil", tc.formatName)
				} else if tc.errorContains != "" && !containsIgnoreCase(err.Error(), tc.errorContains) {
					t.Errorf("Expected error containing '%s', got '%s'", tc.errorContains, err.Error())
				} else {
					t.Logf("✓ Correctly returned expected error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for format '%s', got: %v", tc.formatName, err)
				} else {
					t.Logf("✓ Correctly accepted valid format name")
				}
			}
		})
	}
}

func TestValidateImageDimensions(t *testing.T) {
	testCases := []struct {
		name          string
		width         int
		height        int
		maxWidth      int
		maxHeight     int
		expectedError bool
		errorContains string
	}{
		{
			name:          "Valid dimensions within limits",
			width:         1920,
			height:        1080,
			maxWidth:      2048,
			maxHeight:     2048,
			expectedError: false,
		},
		{
			name:          "Exact maximum dimensions",
			width:         2048,
			height:        2048,
			maxWidth:      2048,
			maxHeight:     2048,
			expectedError: false,
		},
		{
			name:          "Width exceeds maximum",
			width:         4096,
			height:        1080,
			maxWidth:      2048,
			maxHeight:     2048,
			expectedError: true,
			errorContains: "width 4096 exceeds maximum",
		},
		{
			name:          "Height exceeds maximum",
			width:         1920,
			height:        4096,
			maxWidth:      2048,
			maxHeight:     2048,
			expectedError: true,
			errorContains: "height 4096 exceeds maximum",
		},
		{
			name:          "Both dimensions exceed maximum",
			width:         4096,
			height:        4096,
			maxWidth:      2048,
			maxHeight:     2048,
			expectedError: true,
			errorContains: "width 4096 exceeds maximum", // Should catch width first
		},
		{
			name:          "Zero width",
			width:         0,
			height:        1080,
			maxWidth:      2048,
			maxHeight:     2048,
			expectedError: true,
			errorContains: "invalid dimensions",
		},
		{
			name:          "Zero height",
			width:         1920,
			height:        0,
			maxWidth:      2048,
			maxHeight:     2048,
			expectedError: true,
			errorContains: "invalid dimensions",
		},
		{
			name:          "Negative width",
			width:         -1920,
			height:        1080,
			maxWidth:      2048,
			maxHeight:     2048,
			expectedError: true,
			errorContains: "invalid dimensions",
		},
		{
			name:          "Negative height",
			width:         1920,
			height:        -1080,
			maxWidth:      2048,
			maxHeight:     2048,
			expectedError: true,
			errorContains: "invalid dimensions",
		},
		{
			name:          "Very small valid image",
			width:         1,
			height:        1,
			maxWidth:      2048,
			maxHeight:     2048,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := loader.ValidateImageDimensions(tc.width, tc.height, tc.maxWidth, tc.maxHeight)

			// Comprehensive diagnostic logging
			t.Logf("Dimensions: %dx%d", tc.width, tc.height)
			t.Logf("Maximum allowed: %dx%d", tc.maxWidth, tc.maxHeight)
			t.Logf("Expected error: %t", tc.expectedError)
			if err != nil {
				t.Logf("Actual error: %s", err.Error())
			} else {
				t.Logf("Actual error: nil")
			}

			// Calculate aspect ratio for additional context
			if tc.width > 0 && tc.height > 0 {
				aspectRatio := float64(tc.width) / float64(tc.height)
				t.Logf("Aspect ratio: %.3f", aspectRatio)
			}

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error for dimensions %dx%d (max: %dx%d), got nil",
						tc.width, tc.height, tc.maxWidth, tc.maxHeight)
				} else if tc.errorContains != "" && !containsIgnoreCase(err.Error(), tc.errorContains) {
					t.Errorf("Expected error containing '%s', got '%s'", tc.errorContains, err.Error())
				} else {
					t.Logf("✓ Correctly returned expected error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for dimensions %dx%d (max: %dx%d), got: %v",
						tc.width, tc.height, tc.maxWidth, tc.maxHeight, err)
				} else {
					t.Logf("✓ Correctly accepted valid dimensions")
				}
			}
		})
	}
}

func TestValidateImageInfo(t *testing.T) {
	supportedFormats := []string{"jpeg", "jpg", "png", "webp"}
	maxWidth := 2048
	maxHeight := 2048

	testCases := []struct {
		name          string
		info          *loader.ImageInfo
		expectedError bool
		errorContains string
	}{
		{
			name: "Valid image info",
			info: &loader.ImageInfo{
				Format: "jpeg",
				Width:  1920,
				Height: 1080,
				Path:   "/path/to/image.jpg",
			},
			expectedError: false,
		},
		{
			name:          "Nil image info",
			info:          nil,
			expectedError: true,
			errorContains: "image info cannot be nil",
		},
		{
			name: "Invalid format",
			info: &loader.ImageInfo{
				Format: "gif",
				Width:  1920,
				Height: 1080,
				Path:   "/path/to/image.gif",
			},
			expectedError: true,
			errorContains: "unsupported format",
		},
		{
			name: "Oversized width",
			info: &loader.ImageInfo{
				Format: "png",
				Width:  4096,
				Height: 1080,
				Path:   "/path/to/large.png",
			},
			expectedError: true,
			errorContains: "width 4096 exceeds maximum",
		},
		{
			name: "Oversized height",
			info: &loader.ImageInfo{
				Format: "png",
				Width:  1920,
				Height: 4096,
				Path:   "/path/to/tall.png",
			},
			expectedError: true,
			errorContains: "height 4096 exceeds maximum",
		},
		{
			name: "Invalid dimensions - zero width",
			info: &loader.ImageInfo{
				Format: "png",
				Width:  0,
				Height: 1080,
				Path:   "/path/to/invalid.png",
			},
			expectedError: true,
			errorContains: "invalid dimensions",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := loader.ValidateImageInfo(tc.info, maxWidth, maxHeight, supportedFormats)

			// Comprehensive diagnostic logging
			if tc.info != nil {
				t.Logf("Image info:")
				t.Logf("  Format: %s", tc.info.Format)
				t.Logf("  Dimensions: %dx%d", tc.info.Width, tc.info.Height)
				t.Logf("  Path: %s", tc.info.Path)
				t.Logf("  Aspect ratio: %.3f", tc.info.AspectRatio())
				t.Logf("  Pixel count: %d", tc.info.PixelCount())
				t.Logf("  Is portrait: %t", tc.info.IsPortrait())
				t.Logf("  Is landscape: %t", tc.info.IsLandscape())
				t.Logf("  Is square: %t", tc.info.IsSquare())
			} else {
				t.Logf("Image info: nil")
			}
			t.Logf("Maximum allowed: %dx%d", maxWidth, maxHeight)
			t.Logf("Supported formats: %v", supportedFormats)
			t.Logf("Expected error: %t", tc.expectedError)
			if err != nil {
				t.Logf("Actual error: %s", err.Error())
			} else {
				t.Logf("Actual error: nil")
			}

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error for image info, got nil")
				} else if tc.errorContains != "" && !containsIgnoreCase(err.Error(), tc.errorContains) {
					t.Errorf("Expected error containing '%s', got '%s'", tc.errorContains, err.Error())
				} else {
					t.Logf("✓ Correctly returned expected error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for valid image info, got: %v", err)
				} else {
					t.Logf("✓ Correctly accepted valid image info")
				}
			}
		})
	}
}

func TestValidation_EdgeCases(t *testing.T) {
	t.Run("Empty supported formats list", func(t *testing.T) {
		emptyFormats := []string{}

		// Any format should be invalid with empty supported formats
		err := loader.ValidateImageFormatName("jpg", emptyFormats)

		t.Logf("Testing with empty supported formats list")
		t.Logf("Error: %v", err)

		if err == nil {
			t.Error("Expected error with empty supported formats, got nil")
		} else {
			t.Logf("✓ Correctly rejected format when no formats are supported")
		}
	})

	t.Run("Nil supported formats list", func(t *testing.T) {
		// Test with nil formats list
		err := loader.ValidateImageFormatName("jpg", nil)

		t.Logf("Testing with nil supported formats list")
		t.Logf("Error: %v", err)

		if err == nil {
			t.Error("Expected error with nil supported formats, got nil")
		} else {
			t.Logf("✓ Correctly handled nil supported formats")
		}
	})

	t.Run("Very large valid dimensions", func(t *testing.T) {
		// Test dimensions just under the limit
		maxDim := 8192
		err := loader.ValidateImageDimensions(maxDim-1, maxDim-1, maxDim, maxDim)

		t.Logf("Testing dimensions %dx%d (max: %dx%d)", maxDim-1, maxDim-1, maxDim, maxDim)
		t.Logf("Error: %v", err)

		if err != nil {
			t.Errorf("Expected no error for valid large dimensions, got: %v", err)
		} else {
			t.Logf("✓ Correctly accepted large but valid dimensions")
		}
	})
}

// Helper function to check if a string contains another string (case insensitive)
func containsIgnoreCase(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}

	// Simple contains check - convert both to lowercase and check
	lowerS := stringToLower(s)
	lowerSubstr := stringToLower(substr)

	for i := 0; i <= len(lowerS)-len(lowerSubstr); i++ {
		if lowerS[i:i+len(lowerSubstr)] == lowerSubstr {
			return true
		}
	}
	return false
}

// Simple case conversion for testing
func stringToLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			result[i] = s[i] + 32
		} else {
			result[i] = s[i]
		}
	}
	return string(result)
}