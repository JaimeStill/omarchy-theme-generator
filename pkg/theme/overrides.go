package theme

import (
	"fmt"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/palette"
)

// OverrideProcessor handles safe application of user color overrides
// while maintaining WCAG compliance and providing clear feedback.
type OverrideProcessor struct {
	validator *palette.PaletteValidator
	minContrast float64 // Minimum acceptable contrast ratio
}

// NewOverrideProcessor creates a processor that validates overrides
// against the specified background color and contrast requirements.
func NewOverrideProcessor(backgroundColor *color.Color) *OverrideProcessor {
	if backgroundColor == nil {
		backgroundColor = color.NewRGB(255, 255, 255) // Default white
	}
	
	return &OverrideProcessor{
		validator:   palette.NewPaletteValidator(backgroundColor),
		minContrast: 4.5, // WCAG AA standard
	}
}

// OverrideResult contains the results of applying user color overrides,
// including any adjustments made for accessibility compliance.
type OverrideResult struct {
	Primary      *color.Color  // Final primary color (may be adjusted)
	Background   *color.Color  // Final background color
	Foreground   *color.Color  // Final foreground color (may be adjusted)
	Adjustments  []Adjustment  // List of adjustments made for compliance
	Issues       []string      // Human-readable warnings about override conflicts
}

// Adjustment describes a color modification made for WCAG compliance.
type Adjustment struct {
	Role        ColorRole     // Which color role was adjusted
	Original    *color.Color  // Original user-provided color
	Adjusted    *color.Color  // Color after WCAG adjustment
	Reason      string        // Why the adjustment was necessary
	ContrastRatio float64     // Achieved contrast ratio
}

// ApplyOverrides safely applies user color overrides with automatic
// WCAG compliance adjustment and detailed feedback.
func (op *OverrideProcessor) ApplyOverrides(
	overrides ColorOverrides,
	defaultPrimary, defaultBackground, defaultForeground *color.Color,
) (*OverrideResult, error) {
	
	result := &OverrideResult{
		Primary:    defaultPrimary,
		Background: defaultBackground,
		Foreground: defaultForeground,
		Adjustments: make([]Adjustment, 0),
		Issues:      make([]string, 0),
	}
	
	// Apply overrides in order, validating each step
	if overrides.Background != nil {
		result.Background = overrides.Background
		// Update validator background for subsequent validations
		op.validator = palette.NewPaletteValidator(overrides.Background)
	}
	
	if overrides.Primary != nil {
		adjusted, needsAdjustment := op.validateAndAdjust(
			overrides.Primary, 
			result.Background,
			RolePrimary,
		)
		
		result.Primary = adjusted
		
		if needsAdjustment {
			result.Adjustments = append(result.Adjustments, Adjustment{
				Role:     RolePrimary,
				Original: overrides.Primary,
				Adjusted: adjusted,
				Reason:   fmt.Sprintf("Insufficient contrast ratio against background (required: %.1f:1)", op.minContrast),
				ContrastRatio: adjusted.ContrastRatio(result.Background),
			})
		}
	}
	
	if overrides.Foreground != nil {
		adjusted, needsAdjustment := op.validateAndAdjust(
			overrides.Foreground,
			result.Background, 
			RoleForeground,
		)
		
		result.Foreground = adjusted
		
		if needsAdjustment {
			result.Adjustments = append(result.Adjustments, Adjustment{
				Role:     RoleForeground,
				Original: overrides.Foreground,
				Adjusted: adjusted,
				Reason:   fmt.Sprintf("Insufficient contrast ratio against background (required: %.1f:1)", op.minContrast),
				ContrastRatio: adjusted.ContrastRatio(result.Background),
			})
		}
	}
	
	// Final validation to ensure all critical pairs meet standards
	err := op.validateCriticalPairs(result)
	if err != nil {
		return nil, fmt.Errorf("critical validation failed: %w", err)
	}
	
	return result, nil
}

// validateAndAdjust checks if a color meets contrast requirements and adjusts if necessary.
func (op *OverrideProcessor) validateAndAdjust(
	userColor *color.Color, 
	background *color.Color, 
	role ColorRole,
) (*color.Color, bool) {
	
	// Check if user color already meets requirements
	if userColor.ContrastRatio(background) >= op.minContrast {
		return userColor, false // No adjustment needed
	}
	
	// Attempt to adjust the color while preserving hue and saturation
	adjusted := op.adjustForContrast(userColor, background)
	return adjusted, true
}

// adjustForContrast adjusts a color's lightness to meet WCAG contrast requirements
// while preserving hue and saturation as much as possible.
func (op *OverrideProcessor) adjustForContrast(userColor *color.Color, background *color.Color) *color.Color {
	h, s, l := userColor.HSL()
	bgLuminance := background.RelativeLuminance()
	
	// Determine if we need to make the color lighter or darker
	needsLighter := bgLuminance < 0.5
	
	// Binary search for optimal lightness
	var minL, maxL float64
	if needsLighter {
		minL = l
		maxL = 1.0
	} else {
		minL = 0.0
		maxL = l
	}
	
	// Binary search with up to 20 iterations for precision
	for i := 0; i < 20; i++ {
		testL := (minL + maxL) / 2
		testColor := color.NewHSL(h, s, testL)
		
		if testColor.ContrastRatio(background) >= op.minContrast {
			if needsLighter {
				maxL = testL
			} else {
				minL = testL
			}
		} else {
			if needsLighter {
				minL = testL
			} else {
				maxL = testL
			}
		}
	}
	
	// Use the final result
	finalL := (minL + maxL) / 2
	return color.NewHSL(h, s, finalL)
}

// validateCriticalPairs ensures all important color combinations meet WCAG standards.
func (op *OverrideProcessor) validateCriticalPairs(result *OverrideResult) error {
	// Check foreground-background contrast (most critical)
	fgBgRatio := result.Foreground.ContrastRatio(result.Background)
	if fgBgRatio < op.minContrast {
		return fmt.Errorf("foreground-background contrast ratio %.2f:1 is below minimum %.1f:1", 
			fgBgRatio, op.minContrast)
	}
	
	// Check primary-background contrast (important for UI elements)
	primaryBgRatio := result.Primary.ContrastRatio(result.Background)
	if primaryBgRatio < op.minContrast {
		result.Issues = append(result.Issues, 
			fmt.Sprintf("primary-background contrast ratio %.2f:1 is below recommended %.1f:1",
				primaryBgRatio, op.minContrast))
	}
	
	return nil
}

// ValidateTheme performs a comprehensive accessibility audit of a theme.
// This is useful for validating themes before config generation.
func (op *OverrideProcessor) ValidateTheme(theme *Theme) (*ThemeValidation, error) {
	validation := &ThemeValidation{
		Theme:      theme,
		Compliance: make(map[string]ContrastResult),
		Issues:     make([]string, 0),
		Passed:     true,
	}
	
	// Test critical color pairs
	pairs := []struct {
		name string
		fg   *color.Color
		bg   *color.Color
	}{
		{"foreground-background", theme.Foreground, theme.Background},
		{"primary-background", theme.Primary, theme.Background},
	}
	
	for _, pair := range pairs {
		ratio := pair.fg.ContrastRatio(pair.bg)
		result := ContrastResult{
			Ratio:  ratio,
			Passes: ratio >= op.minContrast,
			Level:  op.getComplianceLevel(ratio),
		}
		
		validation.Compliance[pair.name] = result
		
		if !result.Passes {
			validation.Passed = false
			validation.Issues = append(validation.Issues,
				fmt.Sprintf("%s contrast %.2f:1 fails WCAG AA (required: %.1f:1)",
					pair.name, ratio, op.minContrast))
		}
	}
	
	return validation, nil
}

// ThemeValidation contains the results of accessibility validation.
type ThemeValidation struct {
	Theme      *Theme                     // The theme that was validated
	Compliance map[string]ContrastResult  // Contrast test results
	Issues     []string                   // List of accessibility issues found
	Passed     bool                       // True if all tests passed
}

// ContrastResult represents the result of a contrast ratio test.
type ContrastResult struct {
	Ratio  float64 // Actual contrast ratio
	Passes bool    // Whether it meets WCAG AA standard
	Level  string  // Compliance level: "AAA", "AA", "Fail"
}

// getComplianceLevel determines WCAG compliance level from contrast ratio.
func (op *OverrideProcessor) getComplianceLevel(ratio float64) string {
	if ratio >= 7.0 {
		return "AAA"
	} else if ratio >= 4.5 {
		return "AA"
	} else {
		return "Fail"
	}
}

// String provides a summary of the override result.
func (or *OverrideResult) String() string {
	if len(or.Adjustments) == 0 {
		return "All user overrides applied successfully without adjustments"
	}
	
	return fmt.Sprintf("Applied %d user overrides with %d adjustments for WCAG compliance",
		or.countAppliedOverrides(), len(or.Adjustments))
}

// countAppliedOverrides returns the number of overrides that were actually applied.
func (or *OverrideResult) countAppliedOverrides() int {
	count := 0
	if or.Primary != nil {
		count++
	}
	if or.Background != nil {
		count++
	}
	if or.Foreground != nil {
		count++
	}
	return count
}