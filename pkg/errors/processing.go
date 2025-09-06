package errors

import (
	"errors"
	"fmt"
)

// Extractor sentinel errors for common image processing and color extraction failures.
// These errors can be checked using errors.Is() for programmatic handling.
var (
	// ErrUnsupportedFormat indicates the image format is not supported by the extractor.
	ErrUnsupportedFormat = errors.New("unsupported image format")

	// ErrNoExtension indicates the file has no extension to determine format.
	ErrNoExtension = errors.New("no file extension found")

	// ErrImageTooLarge indicates the image exceeds configured size limits.
	ErrImageTooLarge = errors.New("image dimensions exceeded maximum")

	// ErrEmptyImage indicates the image has no pixels or invalid dimensions.
	ErrEmptyImage = errors.New("image is empty")

	// ErrNoColors indicates no colors could be extracted from the image.
	ErrNoColors = errors.New("no colors found in image")

	// ErrInsufficientColors indicates fewer colors than requested were found.
	ErrInsufficientColors = errors.New("insufficient unique colors in image")
)

// ImageFormatError provides detailed information about unsupported or missing image formats.
// It includes the specific file path, extension found, and list of supported formats.
type ImageFormatError struct {
	Path      string   // File path that caused the error
	Extension string   // File extension found (may be empty)
	Supported []string // List of supported file extensions
}

// Error returns a human-readable description of the format error.
func (e *ImageFormatError) Error() string {
	if e.Extension == "" {
		return fmt.Sprintf("no file extension found for %s: supported formats are %v", e.Path, e.Supported)
	}
	return fmt.Sprintf("unsupported format %s for file %s: supported formats are %v", e.Extension, e.Path, e.Supported)
}

// Unwrap returns the underlying sentinel error for use with errors.Is().
func (e *ImageFormatError) Unwrap() error {
	if e.Extension == "" {
		return ErrNoExtension
	}
	return ErrUnsupportedFormat
}

// ImageDimensionError provides details about image size constraint violations.
// It includes both the actual image dimensions and the configured limits.
type ImageDimensionError struct {
	Width     int // Actual image width
	Height    int // Actual image height
	MaxWidth  int // Configured maximum width
	MaxHeight int // Configured maximum height
}

// Error returns a human-readable description of the dimension constraint violation.
func (e *ImageDimensionError) Error() string {
	return fmt.Sprintf("image dimensions %dx%d exceed maximum %dx%d", e.Width, e.Height, e.MaxWidth, e.MaxHeight)
}

// Unwrap returns the underlying sentinel error for use with errors.Is().
func (e *ImageDimensionError) Unwrap() error {
	return ErrImageTooLarge
}

// ImageLoadError wraps file loading and image decoding failures with contextual information.
// It preserves the original error while adding operation context and file path details.
type ImageLoadError struct {
	Path      string // File path that failed to load
	Operation string // Operation that failed (e.g., "open", "decode")
	Err       error  // Underlying error that caused the failure
}

// Error returns a human-readable description of the load failure.
func (e *ImageLoadError) Error() string {
	return fmt.Sprintf("failed to %s image %s: %v", e.Operation, e.Path, e.Err)
}

// Unwrap returns the underlying error for use with errors.Is() and errors.As().
func (e *ImageLoadError) Unwrap() error {
	return e.Err
}

// ExtractionError represents failures during the color extraction pipeline.
// It provides stage context and detailed information about extraction failures.
type ExtractionError struct {
	Stage   string // Pipeline stage where failure occurred (e.g., "frequency", "dominant")
	Details string // Human-readable details about the failure
	Err     error  // Underlying error that caused the failure
}

// Error returns a human-readable description of the extraction failure.
func (e *ExtractionError) Error() string {
	return fmt.Sprintf("extraction failed during %s: %s: %v", e.Stage, e.Details, e.Err)
}

// Unwrap returns the underlying error for use with errors.Is() and errors.As().
func (e *ExtractionError) Unwrap() error {
	return e.Err
}

// ColorCountError indicates a mismatch between requested and available unique colors.
// This typically occurs when requesting a palette larger than the image's color diversity.
type ColorCountError struct {
	Requested int // Number of colors requested
	Available int // Number of unique colors actually found
}

// Error returns a human-readable description of the color count mismatch.
func (e *ColorCountError) Error() string {
	return fmt.Sprintf("requested %d colors but only %d unique colors available", e.Requested, e.Available)
}

// Unwrap returns the underlying sentinel error for use with errors.Is().
func (e *ColorCountError) Unwrap() error {
	return ErrInsufficientColors
}
