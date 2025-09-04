package loader

import (
	"context"
	"image"
	"os"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/errors"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

type ImageInfo struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Format string `json:"format"`
	Path   string `json:"path"`
}

func (info *ImageInfo) AspectRatio() float64 {
	if info.Height == 0 {
		return 0
	}

	return float64(info.Width) / float64(info.Height)
}

func (info *ImageInfo) IsPortrait() bool {
	return info.Height > info.Width
}

func (info *ImageInfo) IsLandscape() bool {
	return info.Width > info.Height
}

func (info *ImageInfo) IsSquare() bool {
	return info.Width == info.Height
}

func (info *ImageInfo) PixelCount() int {
	return info.Width * info.Height
}

type ImageLoader interface {
	LoadImage(ctx context.Context, path string) (image.Image, error)
	GetImageInfo(ctx context.Context, path string) (*ImageInfo, error)
	SupportedFormats() []string
}

type FileLoader struct {
	supportedFormats []string
	maxWidth         int
	maxHeight        int
}

func NewFileLoader(s *settings.Settings) *FileLoader {
	if s == nil {
		s = settings.DefaultSettings()
	}

	return &FileLoader{
		supportedFormats: s.LoaderAllowedFormats,
		maxHeight:        s.LoaderMaxHeight,
		maxWidth:         s.LoaderMaxWidth,
	}
}

func (fl *FileLoader) LoadImage(ctx context.Context, path string) (image.Image, error) {
	info, err := fl.GetImageInfo(ctx, path)
	if err != nil {
		return nil, err
	}

	if err := ValidateImageInfo(info, fl.maxWidth, fl.maxHeight, fl.supportedFormats); err != nil {
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

func (fl *FileLoader) GetImageInfo(ctx context.Context, path string) (*ImageInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, &errors.ImageLoadError{
			Path:      path,
			Operation: "open for info",
			Err:       err,
		}
	}
	defer file.Close()

	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, &errors.ImageLoadError{
			Path:      path,
			Operation: "read config",
			Err:       err,
		}
	}

	return &ImageInfo{
		Width:  config.Width,
		Height: config.Height,
		Format: format,
		Path:   path,
	}, nil
}

func (fl *FileLoader) SupportedFormats() []string {
	return append([]string{}, fl.supportedFormats...)
}
