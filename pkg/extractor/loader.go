package extractor

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/errors"
)

// supportedFormats defines the image file extensions that can be processed.
var supportedFormats = []string{".jpg", ".jpeg", ".png"}

// LoadImage loads an image from the specified file path.
// It supports JPEG and PNG formats through Go's standard library image decoders.
// The image is validated for format support and non-empty content before returning.
// Returns an error if the file cannot be opened, decoded, or is empty.
func LoadImage(path string) (image.Image, error) {
	if err := ValidateImageFormat(path); err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, &errors.ImageLoadError{
			Path:      path,
			Operation: "open",
			Err:       err,
		}
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, &errors.ImageLoadError{
			Path:      path,
			Operation: "decode " + format,
			Err:       err,
		}
	}

	bounds := img.Bounds()
	if bounds.Empty() {
		return nil, &errors.ImageLoadError{
			Path:      path,
			Operation: "validate",
			Err:       errors.ErrEmptyImage,
		}
	}

	return img, nil
}

// ValidateImageFormat checks if the file extension indicates a supported image format.
// It performs case-insensitive matching against supported formats (.jpg, .jpeg, .png).
// Returns an ImageFormatError if the format is not supported or extension is missing.
func ValidateImageFormat(path string) error {
	ext := strings.ToLower(filepath.Ext(path))

	if ext == "" {
		return &errors.ImageFormatError{
			Path:      path,
			Extension: "",
			Supported: supportedFormats,
		}
	}

	if slices.Contains(supportedFormats, ext) {
		return nil
	}

	return &errors.ImageFormatError{
		Path:      path,
		Extension: ext,
		Supported: supportedFormats,
	}
}

// GetImageInfo returns basic information about an image file without fully decoding it.
// This function reads only the image header for efficient size checking and format detection.
// It's particularly useful for validation before full image processing.
// Returns width, height, format string, and any error encountered.
func GetImageInfo(path string) (width, height int, format string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, "", &errors.ImageLoadError{
			Path:      path,
			Operation: "open for info",
			Err:       err,
		}
	}
	defer file.Close()

	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, "", &errors.ImageLoadError{
			Path:      path,
			Operation: "read config",
			Err:       err,
		}
	}

	return config.Width, config.Height, format, nil
}

// LoadImageWithValidation loads an image after performing dimension validation.
// It checks image size against provided limits before full loading to prevent memory issues.
// MaxWidth and MaxHeight values of 0 are treated as unlimited (no validation).
// Returns the loaded image or appropriate error if validation fails.
func LoadImageWithValidation(path string, maxWidth, maxHeight int) (image.Image, error) {
	width, height, _, err := GetImageInfo(path)
	if err != nil {
		return nil, err
	}

	if maxWidth > 0 && maxHeight > 0 {
		if width > maxWidth || height > maxHeight {
			return nil, &errors.ImageDimensionError{
				Width:     width,
				Height:    height,
				MaxWidth:  maxWidth,
				MaxHeight: maxHeight,
			}
		}
	}

	return LoadImage(path)
}

// SupportedFormats returns a copy of the supported image file extensions.
// The returned slice is a defensive copy to prevent external modification.
// Extensions are returned in lowercase with leading dots (e.g., ".jpg", ".png").
func SupportedFormats() []string {
	return append([]string{}, supportedFormats...)
}
