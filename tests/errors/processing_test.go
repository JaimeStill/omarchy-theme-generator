package errors_test

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/errors"
)

func TestImageFormatError(t *testing.T) {
	testCases := []struct {
		name        string
		format      string
		supported   []string
		expectedMsg string
		description string
	}{
		{
			name:        "BMP format not supported",
			format:      "bmp",
			supported:   []string{"jpeg", "png", "webp"},
			expectedMsg: "unsupported format bmp for file test.file: supported formats are [jpeg png webp]",
			description: "Error should list supported formats when format is unsupported",
		},
		{
			name:        "TIFF format not supported",
			format:      "tiff",
			supported:   []string{"jpeg", "png"},
			expectedMsg: "unsupported format tiff for file test.file: supported formats are [jpeg png]",
			description: "Error should handle different supported format lists",
		},
		{
			name:        "Empty format",
			format:      "",
			supported:   []string{"jpeg", "png"},
			expectedMsg: "no file extension found for test.file: supported formats are [jpeg png]",
			description: "Error should handle empty format strings",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := &errors.ImageFormatError{
				Path:      "test.file",
				Extension: tc.format,
				Supported: tc.supported,
			}

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Format: '%s'", tc.format)
			t.Logf("Supported: %v", tc.supported)
			t.Logf("Expected message: %s", tc.expectedMsg)
			t.Logf("Actual message: %s", err.Error())
			t.Logf("Description: %s", tc.description)

			// Verify error message
			if err.Error() != tc.expectedMsg {
				t.Errorf("Error message mismatch:\nExpected: %s\nActual: %s", tc.expectedMsg, err.Error())
			}

			// Verify it's the correct error type
			var formatErr *errors.ImageFormatError
			if !stderrors.As(err, &formatErr) {
				t.Error("Error should be of type *ImageFormatError")
			} else {
				t.Logf("✅ Correct error type: *ImageFormatError")
			}

			// Test error unwrapping - check for the appropriate sentinel
			if tc.format == "" {
				if !stderrors.Is(err, errors.ErrNoExtension) {
					t.Error("Error should unwrap to ErrNoExtension sentinel for empty format")
				} else {
					t.Logf("✅ Correctly unwraps to ErrNoExtension sentinel")
				}
			} else {
				if !stderrors.Is(err, errors.ErrUnsupportedFormat) {
					t.Error("Error should unwrap to ErrUnsupportedFormat sentinel")
				} else {
					t.Logf("✅ Correctly unwraps to ErrUnsupportedFormat sentinel")
				}
			}
		})
	}
}

func TestImageDimensionError(t *testing.T) {
	testCases := []struct {
		name        string
		width       int
		height      int
		maxWidth    int
		maxHeight   int
		expectedMsg string
		description string
	}{
		{
			name:        "Width exceeds limit",
			width:       5000,
			height:      3000,
			maxWidth:    4000,
			maxHeight:   4000,
			expectedMsg: "image dimensions 5000x3000 exceed maximum 4000x4000",
			description: "Error when width exceeds the maximum allowed",
		},
		{
			name:        "Height exceeds limit",
			width:       3000,
			height:      5000,
			maxWidth:    4000,
			maxHeight:   4000,
			expectedMsg: "image dimensions 3000x5000 exceed maximum 4000x4000",
			description: "Error when height exceeds the maximum allowed",
		},
		{
			name:        "Both dimensions exceed limit",
			width:       6000,
			height:      5000,
			maxWidth:    4000,
			maxHeight:   4000,
			expectedMsg: "image dimensions 6000x5000 exceed maximum 4000x4000",
			description: "Error when both dimensions exceed limits",
		},
		{
			name:        "Zero dimensions",
			width:       0,
			height:      100,
			maxWidth:    1000,
			maxHeight:   1000,
			expectedMsg: "image dimensions 0x100 exceed maximum 1000x1000",
			description: "Error should handle zero dimensions",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := &errors.ImageDimensionError{
				Width:     tc.width,
				Height:    tc.height,
				MaxWidth:  tc.maxWidth,
				MaxHeight: tc.maxHeight,
			}

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Image dimensions: %dx%d", tc.width, tc.height)
			t.Logf("Maximum allowed: %dx%d", tc.maxWidth, tc.maxHeight)
			t.Logf("Expected message: %s", tc.expectedMsg)
			t.Logf("Actual message: %s", err.Error())
			t.Logf("Description: %s", tc.description)

			// Verify error message
			if err.Error() != tc.expectedMsg {
				t.Errorf("Error message mismatch:\nExpected: %s\nActual: %s", tc.expectedMsg, err.Error())
			}

			// Verify it's the correct error type
			var dimErr *errors.ImageDimensionError
			if !stderrors.As(err, &dimErr) {
				t.Error("Error should be of type *ImageDimensionError")
			} else {
				t.Logf("✅ Correct error type: *ImageDimensionError")
			}

			// Test error unwrapping
			if !stderrors.Is(err, errors.ErrImageTooLarge) {
				t.Error("Error should unwrap to ErrImageTooLarge sentinel")
			} else {
				t.Logf("✅ Correctly unwraps to ErrImageTooLarge sentinel")
			}
		})
	}
}

func TestImageLoadError(t *testing.T) {
	testCases := []struct {
		name        string
		path        string
		operation   string
		cause       error
		expectedMsg string
		description string
	}{
		{
			name:        "File not found",
			path:        "/path/to/nonexistent.jpg",
			operation:   "open",
			cause:       fmt.Errorf("no such file or directory"),
			expectedMsg: "failed to open image /path/to/nonexistent.jpg: no such file or directory",
			description: "Error when image file doesn't exist",
		},
		{
			name:        "Permission denied",
			path:        "/restricted/image.png",
			operation:   "decode",
			cause:       fmt.Errorf("permission denied"),
			expectedMsg: "failed to decode image /restricted/image.png: permission denied",
			description: "Error when access to image file is denied",
		},
		{
			name:        "Corrupt file",
			path:        "/path/to/corrupt.jpg",
			operation:   "decode",
			cause:       fmt.Errorf("invalid JPEG format"),
			expectedMsg: "failed to decode image /path/to/corrupt.jpg: invalid JPEG format",
			description: "Error when image file is corrupted",
		},
		{
			name:        "Empty path",
			path:        "",
			operation:   "open",
			cause:       fmt.Errorf("empty path"),
			expectedMsg: "failed to open image : empty path",
			description: "Error should handle empty file paths",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := &errors.ImageLoadError{
				Path:      tc.path,
				Operation: tc.operation,
				Err:       tc.cause,
			}

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Path: '%s'", tc.path)
			t.Logf("Operation: '%s'", tc.operation)
			t.Logf("Cause: %v", tc.cause)
			t.Logf("Expected message: %s", tc.expectedMsg)
			t.Logf("Actual message: %s", err.Error())
			t.Logf("Description: %s", tc.description)

			// Verify error message
			if err.Error() != tc.expectedMsg {
				t.Errorf("Error message mismatch:\nExpected: %s\nActual: %s", tc.expectedMsg, err.Error())
			}

			// Verify it's the correct error type
			var loadErr *errors.ImageLoadError
			if !stderrors.As(err, &loadErr) {
				t.Error("Error should be of type *ImageLoadError")
			} else {
				t.Logf("✅ Correct error type: *ImageLoadError")
			}

			// Test error unwrapping to cause
			unwrapped := stderrors.Unwrap(err)
			if unwrapped == nil {
				t.Error("Error should unwrap to original cause")
			} else if unwrapped.Error() != tc.cause.Error() {
				t.Errorf("Expected unwrapped error %q, got %q", tc.cause.Error(), unwrapped.Error())
			} else {
				t.Logf("✅ Correctly unwraps to original cause")
			}
		})
	}
}

func TestExtractionError(t *testing.T) {
	testCases := []struct {
		name        string
		stage       string
		details     string
		cause       error
		expectedMsg string
		description string
	}{
		{
			name:        "Color frequency extraction failed",
			stage:       "frequency",
			details:     "insufficient colors found",
			cause:       fmt.Errorf("insufficient memory"),
			expectedMsg: "extraction failed during frequency: insufficient colors found: insufficient memory",
			description: "Error during color frequency analysis",
		},
		{
			name:        "Dominant color analysis failed",
			stage:       "dominant",
			details:     "color clustering failed",
			cause:       fmt.Errorf("no suitable colors found"),
			expectedMsg: "extraction failed during dominant: color clustering failed: no suitable colors found",
			description: "Error during dominant color extraction",
		},
		{
			name:        "Palette analysis failed",
			stage:       "palette",
			details:     "color space conversion failed",
			cause:       fmt.Errorf("invalid color space"),
			expectedMsg: "extraction failed during palette: color space conversion failed: invalid color space",
			description: "Error during color palette analysis",
		},
		{
			name:        "Empty stage",
			stage:       "",
			details:     "unknown failure",
			cause:       fmt.Errorf("unknown error"),
			expectedMsg: "extraction failed during : unknown failure: unknown error",
			description: "Error should handle empty stage strings",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := &errors.ExtractionError{
				Stage:   tc.stage,
				Details: tc.details,
				Err:     tc.cause,
			}

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Stage: '%s'", tc.stage)
			t.Logf("Details: '%s'", tc.details)
			t.Logf("Cause: %v", tc.cause)
			t.Logf("Expected message: %s", tc.expectedMsg)
			t.Logf("Actual message: %s", err.Error())
			t.Logf("Description: %s", tc.description)

			// Verify error message
			if err.Error() != tc.expectedMsg {
				t.Errorf("Error message mismatch:\nExpected: %s\nActual: %s", tc.expectedMsg, err.Error())
			}

			// Verify it's the correct error type
			var extErr *errors.ExtractionError
			if !stderrors.As(err, &extErr) {
				t.Error("Error should be of type *ExtractionError")
			} else {
				t.Logf("✅ Correct error type: *ExtractionError")
			}

			// Test error unwrapping to cause
			unwrapped := stderrors.Unwrap(err)
			if unwrapped == nil {
				t.Error("Error should unwrap to original cause")
			} else if unwrapped.Error() != tc.cause.Error() {
				t.Errorf("Expected unwrapped error %q, got %q", tc.cause.Error(), unwrapped.Error())
			} else {
				t.Logf("✅ Correctly unwraps to original cause")
			}
		})
	}
}

func TestColorCountError(t *testing.T) {
	testCases := []struct {
		name        string
		requested   int
		available   int
		expectedMsg string
		description string
	}{
		{
			name:        "No colors found",
			requested:   10,
			available:   0,
			expectedMsg: "requested 10 colors but only 0 unique colors available",
			description: "Error when no colors are extracted from image",
		},
		{
			name:        "Too few colors",
			requested:   10,
			available:   3,
			expectedMsg: "requested 10 colors but only 3 unique colors available",
			description: "Error when color count is below minimum threshold",
		},
		{
			name:        "Edge case - need 1, found 0",
			requested:   1,
			available:   0,
			expectedMsg: "requested 1 colors but only 0 unique colors available",
			description: "Edge case with minimal requirements",
		},
		{
			name:        "Large numbers",
			requested:   10000,
			available:   5000,
			expectedMsg: "requested 10000 colors but only 5000 unique colors available",
			description: "Error should handle large color counts",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := &errors.ColorCountError{
				Requested: tc.requested,
				Available: tc.available,
			}

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Colors requested: %d", tc.requested)
			t.Logf("Colors available: %d", tc.available)
			t.Logf("Expected message: %s", tc.expectedMsg)
			t.Logf("Actual message: %s", err.Error())
			t.Logf("Description: %s", tc.description)

			// Verify error message
			if err.Error() != tc.expectedMsg {
				t.Errorf("Error message mismatch:\nExpected: %s\nActual: %s", tc.expectedMsg, err.Error())
			}

			// Verify it's the correct error type
			var countErr *errors.ColorCountError
			if !stderrors.As(err, &countErr) {
				t.Error("Error should be of type *ColorCountError")
			} else {
				t.Logf("✅ Correct error type: *ColorCountError")
			}

			// Test error unwrapping
			if !stderrors.Is(err, errors.ErrInsufficientColors) {
				t.Error("Error should unwrap to ErrInsufficientColors sentinel")
			} else {
				t.Logf("✅ Correctly unwraps to ErrInsufficientColors sentinel")
			}
		})
	}
}

func TestErrorUnwrapping(t *testing.T) {
	// Test complex error unwrapping scenarios
	t.Run("Nested error unwrapping", func(t *testing.T) {
		originalErr := fmt.Errorf("disk full")
		loadErr := &errors.ImageLoadError{
			Path:      "/path/image.jpg",
			Operation: "open",
			Err:       originalErr,
		}
		wrappedErr := fmt.Errorf("processing failed: %w", loadErr)

		// Log diagnostic information
		t.Logf("Original error: %v", originalErr)
		t.Logf("Image load error: %v", loadErr)
		t.Logf("Wrapped error: %v", wrappedErr)

		// Test unwrapping through multiple layers
		if !stderrors.Is(wrappedErr, originalErr) {
			t.Error("Should be able to unwrap to original error through multiple layers")
		} else {
			t.Logf("✅ Successfully unwrapped through multiple layers to original error")
		}

		// Test As() with nested errors
		var imgLoadErr *errors.ImageLoadError
		if !stderrors.As(wrappedErr, &imgLoadErr) {
			t.Error("Should be able to extract ImageLoadError from wrapped error")
		} else {
			t.Logf("✅ Successfully extracted ImageLoadError from wrapped error")
		}
	})
}

func TestErrorChaining(t *testing.T) {
	// Test that errors can be properly chained for context
	t.Run("Error chaining scenario", func(t *testing.T) {
		// Simulate a complex error scenario
		step1Err := fmt.Errorf("network timeout")
		step2Err := &errors.ImageLoadError{
			Path:      "/remote/image.jpg",
			Operation: "download",
			Err:       step1Err,
		}
		step3Err := &errors.ExtractionError{
			Stage:   "remote processing",
			Details: "network image failed",
			Err:     step2Err,
		}
		finalErr := fmt.Errorf("theme generation failed: %w", step3Err)

		// Log the complete error chain
		t.Logf("Error chain:")
		t.Logf("  1. Network: %v", step1Err)
		t.Logf("  2. Load: %v", step2Err)
		t.Logf("  3. Extract: %v", step3Err)
		t.Logf("  4. Final: %v", finalErr)

		// Verify we can unwrap to any level
		if !stderrors.Is(finalErr, step1Err) {
			t.Error("Should unwrap to original network error")
		} else {
			t.Logf("✅ Unwraps to original network error")
		}

		// Extract specific error types from the chain
		var loadErr *errors.ImageLoadError
		var extractErr *errors.ExtractionError

		if !stderrors.As(finalErr, &loadErr) {
			t.Error("Should extract ImageLoadError from chain")
		} else {
			t.Logf("✅ Extracted ImageLoadError from chain")
		}

		if !stderrors.As(finalErr, &extractErr) {
			t.Error("Should extract ExtractionError from chain")
		} else {
			t.Logf("✅ Extracted ExtractionError from chain")
		}
	})
}

func TestSentinelErrors(t *testing.T) {
	// Test that sentinel errors are properly defined and unique
	sentinels := []error{
		errors.ErrUnsupportedFormat,
		errors.ErrNoExtension,
		errors.ErrImageTooLarge,
		errors.ErrEmptyImage,
		errors.ErrNoColors,
		errors.ErrInsufficientColors,
	}

	sentinelNames := []string{
		"ErrUnsupportedFormat",
		"ErrNoExtension", 
		"ErrImageTooLarge",
		"ErrEmptyImage",
		"ErrNoColors",
		"ErrInsufficientColors",
	}

	// Log all sentinel errors
	t.Logf("Testing sentinel errors:")
	for i, sentinel := range sentinels {
		t.Logf("  %s: %v", sentinelNames[i], sentinel)

		// Verify sentinel is not nil
		if sentinel == nil {
			t.Errorf("Sentinel %s should not be nil", sentinelNames[i])
		}

		// Verify sentinel has meaningful message
		if sentinel.Error() == "" {
			t.Errorf("Sentinel %s should have non-empty error message", sentinelNames[i])
		}
	}

	// Verify sentinels are unique (don't equal each other)
	for i, sentinel1 := range sentinels {
		for j, sentinel2 := range sentinels {
			if i != j && stderrors.Is(sentinel1, sentinel2) {
				t.Errorf("Sentinels should be unique: %s equals %s", 
					sentinelNames[i], sentinelNames[j])
			}
		}
	}

	t.Logf("✅ All sentinel errors are properly defined and unique")
}