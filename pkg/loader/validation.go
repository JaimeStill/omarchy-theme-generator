package loader

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/errors"
)

func ValidateImageFormat(path string, supportedFormats []string) error {
	ext := strings.ToLower(filepath.Ext(path))

	if ext == "" {
		return &errors.ImageFormatError{
			Path:      path,
			Extension: "",
			Supported: supportedFormats,
		}
	}

	// Convert extension to format name (remove the dot)
	format := strings.TrimPrefix(ext, ".")
	
	if slices.Contains(supportedFormats, format) {
		return nil
	}

	return &errors.ImageFormatError{
		Path:      path,
		Extension: ext,
		Supported: supportedFormats,
	}
}

func ValidateImageFormatName(formatName string, supportedFormats []string) error {
	if slices.Contains(supportedFormats, formatName) {
		return nil
	}

	return &errors.ImageFormatError{
		Path:      "",
		Extension: formatName,
		Supported: supportedFormats,
	}
}

func ValidateImageDimensions(width, height, maxWidth, maxHeight int) error {
	// Check for invalid dimensions first
	if width <= 0 || height <= 0 {
		return &errors.ImageDimensionError{
			Width:     width,
			Height:    height,
			MaxWidth:  maxWidth,
			MaxHeight: maxHeight,
		}
	}

	if maxWidth > 0 && width > maxWidth {
		return &errors.ImageDimensionError{
			Width:     width,
			Height:    height,
			MaxWidth:  maxWidth,
			MaxHeight: maxHeight,
		}
	}

	if maxHeight > 0 && height > maxHeight {
		return &errors.ImageDimensionError{
			Width:     width,
			Height:    height,
			MaxWidth:  maxWidth,
			MaxHeight: maxHeight,
		}
	}

	return nil
}

func ValidateImageInfo(info *ImageInfo, maxWidth, maxHeight int, supportedFormats []string) error {
	if info == nil {
		return fmt.Errorf("image info cannot be nil")
	}

	if err := ValidateImageFormatName(info.Format, supportedFormats); err != nil {
		return &errors.ImageFormatError{
			Path:      info.Path,
			Extension: info.Format,
			Supported: supportedFormats,
		}
	}

	return ValidateImageDimensions(info.Width, info.Height, maxWidth, maxHeight)
}
