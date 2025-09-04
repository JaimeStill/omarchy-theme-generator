package loader

import (
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

	if slices.Contains(supportedFormats, ext) {
		return nil
	}

	return &errors.ImageFormatError{
		Path:      path,
		Extension: ext,
		Supported: supportedFormats,
	}
}

func ValidateImageDimensions(width, height, maxWidth, maxHeight int) error {
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
	if err := ValidateImageFormat(info.Format, supportedFormats); err != nil {
		return &errors.ImageFormatError{
			Path:      info.Path,
			Extension: info.Format,
			Supported: supportedFormats,
		}
	}

	return ValidateImageDimensions(info.Width, info.Height, maxWidth, maxHeight)
}
