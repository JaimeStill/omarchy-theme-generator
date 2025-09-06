package processor_test

import (
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/processor"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/settings"
)

func TestGetAllCategories(t *testing.T) {
	categories := processor.GetAllCategories()
	
	// Log diagnostic information
	t.Logf("Total categories found: %d", len(categories))
	for i, cat := range categories {
		t.Logf("  %d: %s", i+1, string(cat))
	}
	
	// Verify we have the expected 27 categories
	if len(categories) != 27 {
		t.Errorf("Expected 27 categories, got %d", len(categories))
	}
	
	// Verify specific core categories are present
	expectedCategories := []processor.ColorCategory{
		// Core UI Elements
		processor.CategoryBackground,
		processor.CategoryForeground,
		processor.CategoryDimForeground,
		processor.CategoryCursor,
		
		// Terminal Normal Colors (ANSI 0-7)
		processor.CategoryNormalBlack,
		processor.CategoryNormalRed,
		processor.CategoryNormalGreen,
		processor.CategoryNormalYellow,
		processor.CategoryNormalBlue,
		processor.CategoryNormalMagenta,
		processor.CategoryNormalCyan,
		processor.CategoryNormalWhite,
		
		// Terminal Bright Colors (ANSI 8-15)
		processor.CategoryBrightBlack,
		processor.CategoryBrightRed,
		processor.CategoryBrightGreen,
		processor.CategoryBrightYellow,
		processor.CategoryBrightBlue,
		processor.CategoryBrightMagenta,
		processor.CategoryBrightCyan,
		processor.CategoryBrightWhite,
		
		// Accent Colors
		processor.CategoryAccentPrimary,
		processor.CategoryAccentSecondary,
		processor.CategoryAccentTertiary,
		
		// Semantic Colors
		processor.CategoryError,
		processor.CategoryWarning,
		processor.CategorySuccess,
		processor.CategoryInfo,
	}
	
	// Create a map for efficient lookup
	categoryMap := make(map[processor.ColorCategory]bool)
	for _, cat := range categories {
		categoryMap[cat] = true
	}
	
	// Verify all expected categories are present
	missing := []processor.ColorCategory{}
	for _, expected := range expectedCategories {
		if !categoryMap[expected] {
			missing = append(missing, expected)
		}
	}
	
	if len(missing) > 0 {
		t.Errorf("Missing expected categories:")
		for _, cat := range missing {
			t.Errorf("  - %s", string(cat))
		}
	}
	
	// Verify no duplicates
	if len(categoryMap) != len(categories) {
		t.Errorf("Duplicate categories found: %d unique vs %d total", len(categoryMap), len(categories))
	}
}

func TestGetCategoryPriorityOrder(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	
	// Create a minimal color profile for testing
	profile := &processor.ColorProfile{
		Mode:        processor.Light,
		DominantHue: 120.0, // Green
		Colors: processor.ImageColors{
			TotalPixels: 1000,
		},
	}
	
	priorityOrder := p.GetCategoryPriorityOrder(profile)
	
	// Log diagnostic information
	t.Logf("Priority order length: %d", len(priorityOrder))
	for i, cat := range priorityOrder {
		t.Logf("  Priority %d: %s", i+1, string(cat))
	}
	
	// Verify priority order contains all categories
	allCategories := processor.GetAllCategories()
	if len(priorityOrder) != len(allCategories) {
		t.Errorf("Priority order length %d doesn't match total categories %d", 
			len(priorityOrder), len(allCategories))
	}
	
	// Verify background has highest priority (should be first)
	if len(priorityOrder) > 0 && priorityOrder[0] != processor.CategoryBackground {
		t.Errorf("Expected background to have highest priority, got %s", string(priorityOrder[0]))
	}
	
	// Verify no duplicates in priority order
	seen := make(map[processor.ColorCategory]bool)
	for _, cat := range priorityOrder {
		if seen[cat] {
			t.Errorf("Duplicate category in priority order: %s", string(cat))
		}
		seen[cat] = true
	}
}

func TestGetCategoryCharacteristics(t *testing.T) {
	s := settings.DefaultSettings()
	p := processor.New(s)
	
	// Create profiles for different modes
	lightProfile := &processor.ColorProfile{Mode: processor.Light}
	darkProfile := &processor.ColorProfile{Mode: processor.Dark}
	
	testCases := []struct {
		name     string
		category processor.ColorCategory
		profile  *processor.ColorProfile
	}{
		{
			name:     "Background characteristics - light mode",
			category: processor.CategoryBackground,
			profile:  lightProfile,
		},
		{
			name:     "Background characteristics - dark mode",
			category: processor.CategoryBackground,
			profile:  darkProfile,
		},
		{
			name:     "Foreground characteristics - light mode",
			category: processor.CategoryForeground,
			profile:  lightProfile,
		},
		{
			name:     "Foreground characteristics - dark mode",
			category: processor.CategoryForeground,
			profile:  darkProfile,
		},
		{
			name:     "Accent primary characteristics - light mode",
			category: processor.CategoryAccentPrimary,
			profile:  lightProfile,
		},
		{
			name:     "Error semantic characteristics - light mode",
			category: processor.CategoryError,
			profile:  lightProfile,
		},
		{
			name:     "Normal red ANSI - light mode",
			category: processor.CategoryNormalRed,
			profile:  lightProfile,
		},
		{
			name:     "Bright blue ANSI - dark mode",
			category: processor.CategoryBrightBlue,
			profile:  darkProfile,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			chars := p.GetCategoryCharacteristics(tc.category, tc.profile)
			
			// Log diagnostic information
			t.Logf("Category: %s, Mode: %s", string(tc.category), tc.profile.Mode)
			t.Logf("  MinLightness: %.3f", chars.MinLightness)
			t.Logf("  MaxLightness: %.3f", chars.MaxLightness)
			t.Logf("  MinSaturation: %.3f", chars.MinSaturation)
			t.Logf("  MaxSaturation: %.3f", chars.MaxSaturation)
			t.Logf("  MinContrast: %.1f", chars.MinContrast)
			
			if chars.HueCenter != nil {
				t.Logf("  HueCenter: %.1f°", *chars.HueCenter)
			} else {
				t.Logf("  HueCenter: nil (no constraint)")
			}
			
			if chars.HueTolerance != nil {
				t.Logf("  HueTolerance: %.1f°", *chars.HueTolerance)
			} else {
				t.Logf("  HueTolerance: nil (no constraint)")
			}
			
			// Basic validation
			if chars.MinLightness < 0 || chars.MinLightness > 1 {
				t.Errorf("MinLightness %.3f out of range [0,1]", chars.MinLightness)
			}
			
			if chars.MaxLightness < 0 || chars.MaxLightness > 1 {
				t.Errorf("MaxLightness %.3f out of range [0,1]", chars.MaxLightness)
			}
			
			if chars.MinLightness > chars.MaxLightness {
				t.Errorf("MinLightness %.3f > MaxLightness %.3f", chars.MinLightness, chars.MaxLightness)
			}
			
			if chars.MinSaturation < 0 || chars.MinSaturation > 1 {
				t.Errorf("MinSaturation %.3f out of range [0,1]", chars.MinSaturation)
			}
			
			if chars.MaxSaturation < 0 || chars.MaxSaturation > 1 {
				t.Errorf("MaxSaturation %.3f out of range [0,1]", chars.MaxSaturation)
			}
			
			if chars.MinSaturation > chars.MaxSaturation {
				t.Errorf("MinSaturation %.3f > MaxSaturation %.3f", chars.MinSaturation, chars.MaxSaturation)
			}
			
			if chars.MinContrast < 0 {
				t.Errorf("MinContrast %.1f should be non-negative", chars.MinContrast)
			}
			
			// Hue validation
			if chars.HueCenter != nil {
				hue := *chars.HueCenter
				if hue < 0 || hue >= 360 {
					t.Errorf("HueCenter %.1f° out of range [0,360)", hue)
				}
			}
			
			if chars.HueTolerance != nil {
				tolerance := *chars.HueTolerance
				if tolerance < 0 || tolerance > 180 {
					t.Errorf("HueTolerance %.1f° out of range [0,180]", tolerance)
				}
			}
			
			// Context-specific validation
			switch tc.category {
			case processor.CategoryBackground:
				t.Logf("Background category should have reasonable lightness ranges")
				
			case processor.CategoryForeground:
				if chars.MinContrast < 3.0 {
					t.Logf("Note: Foreground contrast %.1f may be below WCAG AA (4.5)", chars.MinContrast)
				}
				
			case processor.CategoryError:
				t.Logf("Error category characteristics configured for visibility")
				
			default:
				t.Logf("General category characteristics validated")
			}
		})
	}
}