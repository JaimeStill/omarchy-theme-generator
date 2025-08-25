// Package errors provides structured error types and sentinel errors for the Omarchy Theme Generator.
//
// This package centralizes error handling across all components to prevent circular dependencies
// and provide consistent error classification for user interfaces. Errors are organized by
// domain (extractor, palette, theme, etc.) with each domain in its own file.
//
// # Error Categories
//
// The package defines both sentinel errors for programmatic checking and structured error
// types that provide detailed context:
//
//   - Sentinel errors: ErrUnsupportedFormat, ErrImageTooLarge, ErrNoColors, etc.
//   - Structured errors: ImageFormatError, ImageDimensionError, ExtractionError, etc.
//
// # Usage Patterns
//
// Check for specific error conditions using errors.Is:
//
//	if errors.Is(err, errors.ErrUnsupportedFormat) {
//	    // Handle unsupported format
//	}
//
// Extract structured error details using errors.As:
//
//	var fmtErr *errors.ImageFormatError
//	if errors.As(err, &fmtErr) {
//	    fmt.Printf("Supported formats: %v", fmtErr.Supported)
//	}
//
// # Error Chain Compatibility
//
// All error types implement the Unwrap() method to maintain compatibility with
// Go's error wrapping and unwrapping functionality, enabling both errors.Is
// and errors.As to work correctly.
//
// # Domain Organization
//
// Errors are organized by functional domain:
//   - extractor.go: Image loading, color extraction, and analysis errors
//   - palette.go: Color palette generation and synthesis errors (Session 4+)
//   - theme.go: Theme generation and template errors (Session 5+)
//   - template.go: Configuration file generation errors (Session 6+)
package errors