package template

import (
	"fmt"
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/theme"
)

// TerminalColorIndex represents the standard 16-color terminal palette positions.
// This provides type safety and prevents errors from magic number indices.
type TerminalColorIndex int

const (
	// Normal colors (ANSI 0-7)
	ColorBlack TerminalColorIndex = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	
	// Bright colors (ANSI 8-15)
	ColorBrightBlack
	ColorBrightRed
	ColorBrightGreen
	ColorBrightYellow
	ColorBrightBlue
	ColorBrightMagenta
	ColorBrightCyan
	ColorBrightWhite
)

// String returns the ANSI color name for this index.
func (tci TerminalColorIndex) String() string {
	names := []string{
		"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
		"bright_black", "bright_red", "bright_green", "bright_yellow", 
		"bright_blue", "bright_magenta", "bright_cyan", "bright_white",
	}
	if int(tci) >= 0 && int(tci) < len(names) {
		return names[tci]
	}
	return "unknown"
}

// IsNormal returns true if this is a normal (non-bright) color index.
func (tci TerminalColorIndex) IsNormal() bool {
	return tci >= ColorBlack && tci <= ColorWhite
}

// IsBright returns true if this is a bright color index.
func (tci TerminalColorIndex) IsBright() bool {
	return tci >= ColorBrightBlack && tci <= ColorBrightWhite
}

// NormalEquivalent returns the normal color equivalent for bright colors.
// For normal colors, returns itself.
func (tci TerminalColorIndex) NormalEquivalent() TerminalColorIndex {
	if tci.IsBright() {
		return tci - 8
	}
	return tci
}

// BrightEquivalent returns the bright color equivalent for normal colors.
// For bright colors, returns itself.
func (tci TerminalColorIndex) BrightEquivalent() TerminalColorIndex {
	if tci.IsNormal() {
		return tci + 8
	}
	return tci
}

// TerminalColorMap maps terminal color indices to actual colors.
type TerminalColorMap map[TerminalColorIndex]*color.Color

// Get returns the color for the given index, or nil if not found.
func (tcm TerminalColorMap) Get(index TerminalColorIndex) *color.Color {
	return tcm[index]
}

// Set assigns a color to the given terminal index.
func (tcm TerminalColorMap) Set(index TerminalColorIndex, c *color.Color) {
	tcm[index] = c
}

// GetNormal returns all 8 normal colors (ANSI 0-7) in order.
func (tcm TerminalColorMap) GetNormal() [8]*color.Color {
	var colors [8]*color.Color
	for i := 0; i < 8; i++ {
		colors[i] = tcm[TerminalColorIndex(i)]
	}
	return colors
}

// GetBright returns all 8 bright colors (ANSI 8-15) in order.
func (tcm TerminalColorMap) GetBright() [8]*color.Color {
	var colors [8]*color.Color
	for i := 0; i < 8; i++ {
		colors[i] = tcm[TerminalColorIndex(i+8)]
	}
	return colors
}

// TerminalMapper handles the intelligent mapping of theme colors to 
// terminal color slots based on color theory and synthesis strategies.
type TerminalMapper struct {
	// PreserveExtractionOrder determines whether to maintain the original
	// color order from image extraction when possible
	PreserveExtractionOrder bool
	
	// ColorTemperatureTolerance specifies the maximum allowed color temperature
	// variation within the terminal palette (in Kelvin)
	ColorTemperatureTolerance float64
	
	// MinContrastRatio specifies the minimum contrast ratio for terminal colors
	// against the background (WCAG compliance)
	MinContrastRatio float64
}

// NewTerminalMapper creates a terminal mapper with default settings optimized
// for accessibility and color harmony.
func NewTerminalMapper() *TerminalMapper {
	return &TerminalMapper{
		PreserveExtractionOrder:   false, // Prioritize color harmony over extraction order
		ColorTemperatureTolerance: 500.0, // ±500K temperature variation allowed
		MinContrastRatio:          4.5,   // WCAG AA standard
	}
}

// MapToTerminal intelligently maps a theme's color palette to the 16 standard
// terminal color positions, taking into account the synthesis strategy used
// to generate the palette.
func (tm *TerminalMapper) MapToTerminal(t *theme.Theme) (TerminalColorMap, error) {
	if t == nil {
		return nil, fmt.Errorf("theme is nil")
	}
	
	if len(t.Palette) < 8 {
		return nil, fmt.Errorf("theme palette has insufficient colors (%d, need at least 8)", len(t.Palette))
	}
	
	mapping := make(TerminalColorMap, 16)
	
	// Apply strategy-specific mapping logic
	switch t.Metadata.Strategy {
	case "monochromatic":
		return tm.mapMonochromatic(t, mapping)
	case "analogous":
		return tm.mapAnalogous(t, mapping)
	case "complementary":
		return tm.mapComplementary(t, mapping)
	case "triadic":
		return tm.mapTriadic(t, mapping)
	case "tetradic":
		return tm.mapTetradic(t, mapping)
	case "split-complementary":
		return tm.mapSplitComplementary(t, mapping)
	default:
		return tm.mapDefault(t, mapping)
	}
}

// mapMonochromatic maps colors for monochromatic themes by organizing them
// along the lightness dimension while maintaining hue consistency.
func (tm *TerminalMapper) mapMonochromatic(t *theme.Theme, mapping TerminalColorMap) (TerminalColorMap, error) {
	palette := t.Palette
	
	// Sort colors by lightness for monochromatic progression
	sortedColors := tm.sortByLightness(palette)
	
	// Map darker colors to normal positions (0-7)
	normalColors := tm.selectColorsForNormalRange(sortedColors)
	for i, c := range normalColors {
		mapping[TerminalColorIndex(i)] = c
	}
	
	// Generate bright variants using perceptual lightening
	for i := 0; i < 8; i++ {
		baseColor := mapping[TerminalColorIndex(i)]
		if baseColor != nil {
			brightColor := tm.createBrightVariant(baseColor)
			mapping[TerminalColorIndex(i+8)] = brightColor
		}
	}
	
	return mapping, nil
}

// mapComplementary maps colors for complementary themes by distributing
// primary and complementary color families across normal/bright groups.
func (tm *TerminalMapper) mapComplementary(t *theme.Theme, mapping TerminalColorMap) (TerminalColorMap, error) {
	palette := t.Palette
	
	// Split palette into primary and complementary hue groups
	primaryGroup, complementaryGroup := tm.splitByHueDistance(palette, t.Primary)
	
	// Map primary group to positions 0,1,4,5 (traditional red/blue positions)
	primaryPositions := []TerminalColorIndex{ColorBlack, ColorRed, ColorBlue, ColorMagenta}
	for i, pos := range primaryPositions {
		if i < len(primaryGroup) {
			mapping[pos] = primaryGroup[i]
		}
	}
	
	// Map complementary group to positions 2,3,6,7 (traditional green/yellow positions)
	complementaryPositions := []TerminalColorIndex{ColorGreen, ColorYellow, ColorCyan, ColorWhite}
	for i, pos := range complementaryPositions {
		if i < len(complementaryGroup) {
			mapping[pos] = complementaryGroup[i]
		}
	}
	
	// Generate bright variants
	for i := 0; i < 8; i++ {
		if mapping[TerminalColorIndex(i)] != nil {
			brightColor := tm.createBrightVariant(mapping[TerminalColorIndex(i)])
			mapping[TerminalColorIndex(i+8)] = brightColor
		}
	}
	
	return mapping, nil
}

// mapTriadic maps colors for triadic themes by distributing three hue families
// evenly across the terminal color space.
func (tm *TerminalMapper) mapTriadic(t *theme.Theme, mapping TerminalColorMap) (TerminalColorMap, error) {
	palette := t.Palette
	
	// Group colors into three hue families (120° apart)
	hueGroups := tm.groupByTriadicHues(palette, t.Primary)
	
	// Distribute groups across terminal positions
	// Group 1: positions 0,1,2
	// Group 2: positions 3,4,5  
	// Group 3: positions 6,7 + neutral position
	
	positionGroups := [][]TerminalColorIndex{
		{ColorBlack, ColorRed, ColorGreen},
		{ColorYellow, ColorBlue, ColorMagenta},
		{ColorCyan, ColorWhite},
	}
	
	for groupIdx, group := range hueGroups {
		if groupIdx < len(positionGroups) {
			positions := positionGroups[groupIdx]
			for i, pos := range positions {
				if i < len(group) {
					mapping[pos] = group[i]
				}
			}
		}
	}
	
	// Generate bright variants
	for i := 0; i < 8; i++ {
		if mapping[TerminalColorIndex(i)] != nil {
			brightColor := tm.createBrightVariant(mapping[TerminalColorIndex(i)])
			mapping[TerminalColorIndex(i+8)] = brightColor
		}
	}
	
	return mapping, nil
}

// mapTetradic maps colors for tetradic (square) themes using four hue families.
func (tm *TerminalMapper) mapTetradic(t *theme.Theme, mapping TerminalColorMap) (TerminalColorMap, error) {
	palette := t.Palette
	
	// Group into four hue families (90° apart)
	hueGroups := tm.groupByTetradicHues(palette, t.Primary)
	
	// Map each group to 2 positions
	positionGroups := [][]TerminalColorIndex{
		{ColorBlack, ColorRed},
		{ColorGreen, ColorYellow},
		{ColorBlue, ColorMagenta},
		{ColorCyan, ColorWhite},
	}
	
	// First pass: assign colors from groups
	for groupIdx, group := range hueGroups {
		if groupIdx < len(positionGroups) {
			positions := positionGroups[groupIdx]
			for i, pos := range positions {
				if i < len(group) {
					mapping[pos] = group[i]
				}
			}
		}
	}
	
	// Second pass: fill any missing positions with fallback colors
	allColors := make([]*color.Color, 0, len(palette))
	for _, group := range hueGroups {
		allColors = append(allColors, group...)
	}
	
	// Ensure we have enough colors by falling back to palette order if needed
	if len(allColors) < len(palette) {
		allColors = palette
	}
	
	colorIdx := 0
	for i := 0; i < 8; i++ {
		pos := TerminalColorIndex(i)
		if mapping[pos] == nil && colorIdx < len(allColors) {
			mapping[pos] = allColors[colorIdx]
			colorIdx++
		}
	}
	
	// Generate bright variants
	for i := 0; i < 8; i++ {
		if mapping[TerminalColorIndex(i)] != nil {
			brightColor := tm.createBrightVariant(mapping[TerminalColorIndex(i)])
			mapping[TerminalColorIndex(i+8)] = brightColor
		} else {
			// Fallback: create a brighter version of the primary color
			brightColor := tm.createBrightVariant(t.Primary)
			mapping[TerminalColorIndex(i+8)] = brightColor
		}
	}
	
	return mapping, nil
}

// mapSplitComplementary maps colors for split-complementary themes.
func (tm *TerminalMapper) mapSplitComplementary(t *theme.Theme, mapping TerminalColorMap) (TerminalColorMap, error) {
	// Similar to complementary but with two colors opposite the primary
	return tm.mapComplementary(t, mapping) // Simplified implementation for now
}

// mapAnalogous maps colors for analogous themes (adjacent hues).
func (tm *TerminalMapper) mapAnalogous(t *theme.Theme, mapping TerminalColorMap) (TerminalColorMap, error) {
	// Similar to monochromatic but with slight hue variations
	return tm.mapMonochromatic(t, mapping) // Simplified implementation for now
}

// mapDefault provides fallback mapping when strategy is unknown.
func (tm *TerminalMapper) mapDefault(t *theme.Theme, mapping TerminalColorMap) (TerminalColorMap, error) {
	palette := t.Palette
	
	// Simple sequential mapping for first 8 colors
	for i := 0; i < 8 && i < len(palette); i++ {
		mapping[TerminalColorIndex(i)] = palette[i]
	}
	
	// Use remaining colors for bright positions, or generate bright variants
	for i := 8; i < 16; i++ {
		if i < len(palette) {
			mapping[TerminalColorIndex(i)] = palette[i]
		} else {
			// Generate bright variant of corresponding normal color
			normalIdx := i - 8
			if mapping[TerminalColorIndex(normalIdx)] != nil {
				brightColor := tm.createBrightVariant(mapping[TerminalColorIndex(normalIdx)])
				mapping[TerminalColorIndex(i)] = brightColor
			}
		}
	}
	
	return mapping, nil
}

// Helper methods for color mapping logic

// sortByLightness sorts colors from darkest to lightest using HSL lightness.
func (tm *TerminalMapper) sortByLightness(colors []*color.Color) []*color.Color {
	sorted := make([]*color.Color, len(colors))
	copy(sorted, colors)
	
	// Simple bubble sort by lightness (efficient for small arrays)
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			_, _, l1 := sorted[j].HSL()
			_, _, l2 := sorted[j+1].HSL()
			if l1 > l2 {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}
	
	return sorted
}

// selectColorsForNormalRange selects 8 colors from the palette for normal positions,
// prioritizing darker colors that will contrast well with bright variants.
func (tm *TerminalMapper) selectColorsForNormalRange(colors []*color.Color) [8]*color.Color {
	var selected [8]*color.Color
	
	// Take first 8 colors, padding with nil if insufficient
	for i := 0; i < 8; i++ {
		if i < len(colors) {
			selected[i] = colors[i]
		}
	}
	
	return selected
}

// createBrightVariant creates a perceptually brighter version of the given color
// using LAB color space for uniform lightness adjustment.
func (tm *TerminalMapper) createBrightVariant(c *color.Color) *color.Color {
	if c == nil {
		return nil
	}
	
	// For now, use HSL lightening as we haven't implemented LAB-based lightening yet
	// This will be enhanced when we add the LAB perceptual methods
	h, s, l := c.HSL()
	
	// Increase lightness by 20%, ensuring we don't exceed 95%
	newL := math.Min(0.95, l + 0.20)
	
	return color.NewHSL(h, s, newL)
}

// splitByHueDistance splits colors into two groups based on hue distance from primary.
func (tm *TerminalMapper) splitByHueDistance(colors []*color.Color, primary *color.Color) ([]*color.Color, []*color.Color) {
	if primary == nil || len(colors) == 0 {
		return colors, nil
	}
	
	primaryH, _, _ := primary.HSL()
	var primaryGroup, complementaryGroup []*color.Color
	
	for _, c := range colors {
		h, _, _ := c.HSL()
		
		// Calculate hue distance (accounting for circular nature of hue)
		hueDist := math.Min(math.Abs(h-primaryH), 1.0-math.Abs(h-primaryH))
		
		// Colors within 90° (0.25 in normalized hue) belong to primary group
		if hueDist <= 0.25 {
			primaryGroup = append(primaryGroup, c)
		} else {
			complementaryGroup = append(complementaryGroup, c)
		}
	}
	
	return primaryGroup, complementaryGroup
}

// groupByTriadicHues groups colors into three hue families for triadic harmony.
func (tm *TerminalMapper) groupByTriadicHues(colors []*color.Color, primary *color.Color) [][]*color.Color {
	if primary == nil || len(colors) == 0 {
		return [][]*color.Color{colors}
	}
	
	primaryH, _, _ := primary.HSL()
	
	// Calculate triadic hue positions (120° apart)
	hue1 := primaryH
	hue2 := math.Mod(primaryH + 0.333, 1.0) // +120°
	hue3 := math.Mod(primaryH + 0.667, 1.0) // +240°
	
	var group1, group2, group3 []*color.Color
	
	for _, c := range colors {
		h, _, _ := c.HSL()
		
		// Find closest triadic hue
		dist1 := math.Min(math.Abs(h-hue1), 1.0-math.Abs(h-hue1))
		dist2 := math.Min(math.Abs(h-hue2), 1.0-math.Abs(h-hue2))
		dist3 := math.Min(math.Abs(h-hue3), 1.0-math.Abs(h-hue3))
		
		if dist1 <= dist2 && dist1 <= dist3 {
			group1 = append(group1, c)
		} else if dist2 <= dist3 {
			group2 = append(group2, c)
		} else {
			group3 = append(group3, c)
		}
	}
	
	return [][]*color.Color{group1, group2, group3}
}

// groupByTetradicHues groups colors into four hue families for tetradic harmony.
func (tm *TerminalMapper) groupByTetradicHues(colors []*color.Color, primary *color.Color) [][]*color.Color {
	if primary == nil || len(colors) == 0 {
		return [][]*color.Color{colors}
	}
	
	primaryH, _, _ := primary.HSL()
	
	// Calculate tetradic hue positions (90° apart)
	hue1 := primaryH
	hue2 := math.Mod(primaryH + 0.25, 1.0)  // +90°
	hue3 := math.Mod(primaryH + 0.50, 1.0)  // +180°
	hue4 := math.Mod(primaryH + 0.75, 1.0)  // +270°
	
	var group1, group2, group3, group4 []*color.Color
	
	for _, c := range colors {
		h, _, _ := c.HSL()
		
		// Find closest tetradic hue
		distances := []float64{
			math.Min(math.Abs(h-hue1), 1.0-math.Abs(h-hue1)),
			math.Min(math.Abs(h-hue2), 1.0-math.Abs(h-hue2)),
			math.Min(math.Abs(h-hue3), 1.0-math.Abs(h-hue3)),
			math.Min(math.Abs(h-hue4), 1.0-math.Abs(h-hue4)),
		}
		
		minIdx := 0
		minDist := distances[0]
		for i, dist := range distances[1:] {
			if dist < minDist {
				minDist = dist
				minIdx = i + 1
			}
		}
		
		switch minIdx {
		case 0:
			group1 = append(group1, c)
		case 1:
			group2 = append(group2, c)
		case 2:
			group3 = append(group3, c)
		case 3:
			group4 = append(group4, c)
		}
	}
	
	return [][]*color.Color{group1, group2, group3, group4}
}