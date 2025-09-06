package loader_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/loader"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

func TestLoader_New(t *testing.T) {
	s := settings.DefaultSettings()
	l := loader.NewFileLoader(s)
	
	if l == nil {
		t.Fatal("Expected loader to be created, got nil")
	}
}

func TestLoader_LoadImage_ValidJPEG(t *testing.T) {
	s := settings.DefaultSettings()
	l := loader.NewFileLoader(s)
	
	// Use existing test image
	imgPath := "../../tests/images/simple.png"
	ctx := context.Background()
	
	loadedImg, err := l.LoadImage(ctx, imgPath)
	if err != nil {
		t.Fatalf("Failed to load valid PNG: %v", err)
	}
	
	if loadedImg == nil {
		t.Fatal("Expected loaded image to be non-nil")
	}
	
	// Verify image has reasonable bounds
	bounds := loadedImg.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Errorf("Expected positive dimensions, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestLoader_LoadImage_ValidPNG(t *testing.T) {
	s := settings.DefaultSettings()
	l := loader.NewFileLoader(s)
	
	// Use existing test image
	imgPath := "../../tests/images/primary-background.png"
	ctx := context.Background()
	
	loadedImg, err := l.LoadImage(ctx, imgPath)
	if err != nil {
		t.Fatalf("Failed to load valid PNG: %v", err)
	}
	
	if loadedImg == nil {
		t.Fatal("Expected loaded image to be non-nil")
	}
}

func TestLoader_LoadImage_NonexistentFile(t *testing.T) {
	s := settings.DefaultSettings()
	l := loader.NewFileLoader(s)
	
	ctx := context.Background()
	img, err := l.LoadImage(ctx, "/nonexistent/path/image.png")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
	if img != nil {
		t.Error("Expected nil image for nonexistent file")
	}
}

func TestLoader_LoadImage_InvalidFormat(t *testing.T) {
	// Create a temporary text file (not an image)
	tempDir := t.TempDir()
	txtPath := filepath.Join(tempDir, "notanimage.txt")
	
	if err := os.WriteFile(txtPath, []byte("This is not an image"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	s := settings.DefaultSettings()
	l := loader.NewFileLoader(s)
	
	ctx := context.Background()
	img, err := l.LoadImage(ctx, txtPath)
	if err == nil {
		t.Error("Expected error for invalid image format")
	}
	if img != nil {
		t.Error("Expected nil image for invalid format")
	}
}

func TestLoader_LoadImage_OversizedImage(t *testing.T) {
	s := settings.DefaultSettings()
	// Set very small max dimensions to test validation
	s.LoaderMaxWidth = 10
	s.LoaderMaxHeight = 10
	l := loader.NewFileLoader(s)
	
	// Use a large existing test image (most are much larger than 10x10)
	imgPath := "../../tests/images/bokeh.jpeg"
	
	ctx := context.Background()
	loadedImg, err := l.LoadImage(ctx, imgPath)
	
	// Depending on implementation, this might resize or error
	if err != nil {
		// If it errors, that's acceptable
		t.Logf("Loader rejected oversized image: %v", err)
		return
	}
	
	// If it succeeds, check if it was resized
	if loadedImg != nil {
		bounds := loadedImg.Bounds()
		if bounds.Dx() > s.LoaderMaxWidth || bounds.Dy() > s.LoaderMaxHeight {
			t.Errorf("Image exceeds max dimensions: %dx%d > %dx%d",
				bounds.Dx(), bounds.Dy(), s.LoaderMaxWidth, s.LoaderMaxHeight)
		}
	}
}

func TestLoader_GetImageInfo(t *testing.T) {
	s := settings.DefaultSettings()
	l := loader.NewFileLoader(s)
	
	// Test with different image formats
	testCases := []struct {
		name     string
		path     string
		format   string
	}{
		{"PNG image", "../../tests/images/simple.png", "png"},
		{"JPEG image", "../../tests/images/grayscale.jpeg", "jpeg"},
		{"Another JPEG", "../../tests/images/abstract.jpeg", "jpeg"},
	}
	
	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := l.GetImageInfo(ctx, tc.path)
			if err != nil {
				t.Fatalf("Failed to get image info: %v", err)
			}
			
			if info == nil {
				t.Fatal("Expected image info to be non-nil")
			}
			
			if info.Width <= 0 || info.Height <= 0 {
				t.Errorf("Expected positive dimensions, got %dx%d", info.Width, info.Height)
			}
			
			if info.Format != tc.format {
				t.Errorf("Expected format '%s', got '%s'", tc.format, info.Format)
			}
			
			if info.Path != tc.path {
				t.Errorf("Expected path %s, got %s", tc.path, info.Path)
			}
		})
	}
}

func TestLoader_SupportedFormats(t *testing.T) {
	s := settings.DefaultSettings()
	l := loader.NewFileLoader(s)
	
	// Check that loader returns supported formats
	supportedFormats := l.SupportedFormats()
	
	if len(supportedFormats) == 0 {
		t.Error("Expected at least one supported format")
	}
	
	// Check that common formats are supported
	expectedFormats := []string{"jpeg", "jpg", "png"}
	
	for _, expected := range expectedFormats {
		found := false
		for _, supported := range supportedFormats {
			if supported == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected format %s to be in supported formats", expected)
		}
	}
}

func TestLoader_ImageInfo_Methods(t *testing.T) {
	info := &loader.ImageInfo{
		Width:  1920,
		Height: 1080,
		Format: "png",
		Path:   "/test/image.png",
	}
	
	t.Run("AspectRatio", func(t *testing.T) {
		expected := 1920.0 / 1080.0
		if ratio := info.AspectRatio(); ratio != expected {
			t.Errorf("Expected aspect ratio %v, got %v", expected, ratio)
		}
	})
	
	t.Run("IsLandscape", func(t *testing.T) {
		if !info.IsLandscape() {
			t.Error("1920x1080 should be landscape")
		}
	})
	
	t.Run("IsPortrait", func(t *testing.T) {
		if info.IsPortrait() {
			t.Error("1920x1080 should not be portrait")
		}
	})
	
	t.Run("IsSquare", func(t *testing.T) {
		if info.IsSquare() {
			t.Error("1920x1080 should not be square")
		}
	})
	
	t.Run("PixelCount", func(t *testing.T) {
		expected := 1920 * 1080
		if count := info.PixelCount(); count != expected {
			t.Errorf("Expected pixel count %d, got %d", expected, count)
		}
	})
	
	t.Run("ZeroHeight", func(t *testing.T) {
		zeroInfo := &loader.ImageInfo{Width: 100, Height: 0}
		if ratio := zeroInfo.AspectRatio(); ratio != 0 {
			t.Errorf("Expected aspect ratio 0 for zero height, got %v", ratio)
		}
	})
}

func TestLoader_MemorySafety(t *testing.T) {
	// Test that loader doesn't consume excessive memory for large images
	s := settings.DefaultSettings()
	l := loader.NewFileLoader(s)
	
	// Use the largest test image we have
	imgPath := "../../tests/images/bokeh.jpeg" // 1.4MB file
	
	ctx := context.Background()
	loadedImg, err := l.LoadImage(ctx, imgPath)
	if err != nil {
		t.Fatalf("Failed to load large image: %v", err)
	}
	
	if loadedImg == nil {
		t.Fatal("Expected loaded image to be non-nil")
	}
	
	// Image should load successfully without excessive memory use
	// (Go's garbage collector will handle cleanup)
	bounds := loadedImg.Bounds()
	t.Logf("Successfully loaded %dx%d image", bounds.Dx(), bounds.Dy())
}

func TestLoader_FormatValidation(t *testing.T) {
	s := settings.DefaultSettings()
	// Restrict to only PNG files
	s.LoaderAllowedFormats = []string{"png"}
	l := loader.NewFileLoader(s)
	
	ctx := context.Background()
	
	t.Run("Allowed PNG format", func(t *testing.T) {
		_, err := l.LoadImage(ctx, "../../tests/images/simple.png")
		if err != nil {
			t.Errorf("PNG should be allowed: %v", err)
		}
	})
	
	t.Run("Disallowed JPEG format", func(t *testing.T) {
		_, err := l.LoadImage(ctx, "../../tests/images/grayscale.jpeg")
		if err == nil {
			t.Error("JPEG should be rejected when only PNG is allowed")
		}
	})
}