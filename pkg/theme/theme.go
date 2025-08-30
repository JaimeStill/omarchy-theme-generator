// Package theme provides theme generation orchestration and integration 
// with the extraction → hybrid → synthesis pipeline.
package theme

import (
	"fmt"
	"image"
	"time"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/palette"
)

// Theme represents a complete color theme generated from an image.
// It contains all necessary color information for generating config files.
type Theme struct {
	Name         string        // Theme name (derived from source image filename)
	SourcePath   string        // Path to source image file
	IsLight      bool          // True for light themes, false for dark themes
	Primary      *color.Color  // Primary accent color
	Background   *color.Color  // Background color
	Foreground   *color.Color  // Foreground/text color
	Palette      []*color.Color // Full color palette from synthesis pipeline
	Metadata     ThemeMetadata // Generation metadata and performance info
}

// ThemeMetadata contains information about how the theme was generated
// and performance metrics for validation and debugging.
type ThemeMetadata struct {
	GenerationMode   string                    // "extract", "hybrid", "synthesize"
	Strategy         string                    // Color theory strategy used
	BaseColor        *color.Color              // Base color used for synthesis
	ExtractedColors  int                       // Number of colors from extraction
	SynthesizedColors int                      // Number of colors from synthesis
	Performance      PerformanceMetrics        // Generation timing and memory usage
	Validation       *palette.ValidationResult // WCAG compliance results
	Generated        time.Time                 // When theme was generated
}

// PerformanceMetrics tracks theme generation performance for validation
// against the <2s target for 4K images.
type PerformanceMetrics struct {
	ExtractionTime time.Duration // Time spent on image extraction
	SynthesisTime  time.Duration // Time spent on color synthesis  
	ValidationTime time.Duration // Time spent on WCAG validation
	TotalTime      time.Duration // Total generation time
	MemoryUsage    int64         // Peak memory usage in bytes
	ImageSize      ImageSize     // Source image dimensions
}

// ImageSize represents the dimensions of the source image.
type ImageSize struct {
	Width  int
	Height int
}

// PixelCount returns the total number of pixels in the image.
func (is ImageSize) PixelCount() int {
	return is.Width * is.Height
}

// IsMegapixel returns true if the image is at least 1 megapixel.
func (is ImageSize) IsMegapixel() bool {
	return is.PixelCount() >= 1000000
}

// Is4K returns true if the image is approximately 4K resolution.
func (is ImageSize) Is4K() bool {
	return is.PixelCount() >= 3840*2160
}

// ThemeMode represents the light/dark mode preference.
type ThemeMode int

const (
	// ModeAuto automatically detects light/dark based on image characteristics
	ModeAuto ThemeMode = iota
	// ModeLight forces light theme generation
	ModeLight
	// ModeDark forces dark theme generation  
	ModeDark
)

// String returns the string representation of the theme mode.
func (tm ThemeMode) String() string {
	switch tm {
	case ModeAuto:
		return "auto"
	case ModeLight:
		return "light"
	case ModeDark:
		return "dark"
	default:
		return "unknown"
	}
}

// ThemeConfig contains the input configuration for theme generation.
// It integrates user preferences with image analysis.
type ThemeConfig struct {
	SourceImage   image.Image                // Source image for color extraction
	Mode          ThemeMode                  // Light/dark/auto mode preference
	Overrides     ColorOverrides            // User color overrides
	SynthesisOpts *palette.SynthesisOptions // Synthesis pipeline options
	Name          string                    // Optional theme name override
}

// ColorOverrides allows users to override specific theme colors.
// All overrides are validated for WCAG compliance.
type ColorOverrides struct {
	Primary    *color.Color // Override primary accent color
	Background *color.Color // Override background color
	Foreground *color.Color // Override foreground/text color
}

// HasOverrides returns true if any color overrides are specified.
func (co ColorOverrides) HasOverrides() bool {
	return co.Primary != nil || co.Background != nil || co.Foreground != nil
}

// String provides a human-readable summary of the theme.
func (t *Theme) String() string {
	modeStr := "dark"
	if t.IsLight {
		modeStr = "light"
	}
	
	return fmt.Sprintf(
		"Theme{name=%s, mode=%s, colors=%d, strategy=%s, performance=%v}",
		t.Name,
		modeStr,
		len(t.Palette),
		t.Metadata.Strategy,
		t.Metadata.Performance.TotalTime,
	)
}

// GetColorByRole returns a color for a specific role in the theme.
// This provides a type-safe way to access semantic colors.
func (t *Theme) GetColorByRole(role ColorRole) *color.Color {
	switch role {
	case RolePrimary:
		return t.Primary
	case RoleBackground:
		return t.Background
	case RoleForeground:
		return t.Foreground
	default:
		// For palette colors, return by index if within bounds
		if int(role) >= 0 && int(role) < len(t.Palette) {
			return t.Palette[int(role)]
		}
		return nil
	}
}

// ColorRole represents semantic color roles in a theme.
// This provides type-safe color access and prevents string-based errors.
type ColorRole int

const (
	// Core semantic roles
	RolePrimary ColorRole = iota
	RoleBackground
	RoleForeground
	
	// Extended palette roles (map to palette indices)
	RoleAccent1
	RoleAccent2
	RoleAccent3
	RoleWarning
	RoleError
	RoleSuccess
	RoleInfo
)

// String returns the string representation of the color role.
func (cr ColorRole) String() string {
	switch cr {
	case RolePrimary:
		return "primary"
	case RoleBackground:
		return "background"
	case RoleForeground:
		return "foreground"
	case RoleAccent1:
		return "accent1"
	case RoleAccent2:
		return "accent2"
	case RoleAccent3:
		return "accent3"
	case RoleWarning:
		return "warning"
	case RoleError:
		return "error"
	case RoleSuccess:
		return "success"
	case RoleInfo:
		return "info"
	default:
		return "unknown"
	}
}