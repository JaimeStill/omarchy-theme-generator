// Package generative provides computational image generation for testing and aesthetic purposes.
// This package implements various algorithmic image generation techniques including
// synthetic test images, aesthetic generators, and mathematical pattern creation.
package generative

import (
	"image"
	stdcolor "image/color"
	"math"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
)

// Generate4KTestImage creates a synthetic 4K image (3840x2160) with varied colors for benchmarking.
// The image contains gradient patterns with noise to simulate realistic color distribution.
// This ensures consistent benchmarking across different environments and validates performance targets.
func Generate4KTestImage() image.Image {
	width, height := 3840, 2160 // 4K resolution
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Create a gradient pattern with varying colors for realistic extraction testing
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Generate varied colors based on position
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(((x + y) * 255) / (width + height))
			
			// Add some noise for more realistic color distribution
			if (x+y)%7 == 0 {
				r = uint8((int(r) + 50) % 256)
			}
			if (x*y)%11 == 0 {
				g = uint8((int(g) + 30) % 256)
			}
			if (x-y)%13 == 0 {
				b = uint8((int(b) + 70) % 256)
			}

			img.Set(x, y, stdcolor.RGBA{r, g, b, 255})
		}
	}

	return img
}

// GenerateMonochromeTestImage creates a grayscale image for testing synthesis edge cases.
func GenerateMonochromeTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Create grayscale gradient with subtle variations
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Base grayscale value
			gray := uint8((x + y) * 255 / (width + height))
			
			// Add subtle noise to create some unique colors but keep it monochrome
			if (x*y)%17 == 0 {
				gray = uint8((int(gray) + 10) % 256)
			}

			img.Set(x, y, stdcolor.RGBA{gray, gray, gray, 255})
		}
	}

	return img
}

// GenerateHighContrastTestImage creates an image with few but very distinct colors.
func GenerateHighContrastTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Define high contrast colors
	colors := []stdcolor.RGBA{
		{0, 0, 0, 255},       // Black
		{255, 255, 255, 255}, // White
		{255, 0, 0, 255},     // Red
		{0, 255, 0, 255},     // Green
		{0, 0, 255, 255},     // Blue
	}

	// Create blocks of each color (20% each)
	colorCount := len(colors)
	pixelsPerColor := (width * height) / colorCount
	currentPixel := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			colorIndex := currentPixel / pixelsPerColor
			if colorIndex >= colorCount {
				colorIndex = colorCount - 1 // Safety check
			}
			
			img.Set(x, y, colors[colorIndex])
			currentPixel++
		}
	}

	return img
}

// Generate80sVectorImage creates an 80's synthwave-style image with neon wireframes
func Generate80sVectorImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Dark background
	darkBg := color.NewHSL(0.55, 0.8, 0.05) // Deep blue #000814
	bgR, bgG, bgB, bgA := darkBg.RGBA()
	bgColor := stdcolor.RGBA{R: uint8(bgR), G: uint8(bgG), B: uint8(bgB), A: uint8(bgA)}
	
	// Fill background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}
	
	// Synthwave palette
	neonPurple := color.NewHSL(0.83, 1.0, 0.5)  // Electric purple
	neonCyan := color.NewHSL(0.5, 1.0, 0.5)    // Cyan
	neonPink := color.NewHSL(0.92, 1.0, 0.6)   // Hot pink
	
	// Golden ratio horizon line
	horizonY := int(float64(height) * 0.618)
	
	// Draw perspective grid lines
	vanishingX := width / 2
	vanishingY := horizonY
	
	// Vertical grid lines receding to vanishing point
	for i := 0; i < 20; i++ {
		x := (i * width) / 20
		drawNeonLine(img, x, height, vanishingX, vanishingY, neonCyan)
	}
	
	// Horizontal grid lines
	for i := 0; i < 10; i++ {
		y := horizonY + ((i * (height - horizonY)) / 10)
		drawNeonLine(img, 0, y, width, y, neonPurple)
	}
	
	// Simple wireframe mountains using sine waves
	for x := 0; x < width; x++ {
		// Generate mountain silhouette using mathematical curves
		mountainHeight := int(float64(height) * 0.3 * (0.5 + 0.3*math.Sin(float64(x)*0.01) + 0.2*math.Sin(float64(x)*0.03)))
		y := horizonY - mountainHeight
		
		if y >= 0 && y < height {
			drawNeonPixel(img, x, y, neonPink)
			// Draw some vertical lines for wireframe effect
			if x%20 == 0 {
				for wy := y; wy < horizonY; wy += 10 {
					if wy < height {
						drawNeonPixel(img, x, wy, neonPink)
					}
				}
			}
		}
	}
	
	return img
}

// CassetteFuturismPalette contains the temperature-matched industrial color palette
type CassetteFuturismPalette struct {
	BackgroundDark    *color.Color
	PanelDark         *color.Color
	PanelLight        *color.Color
	ControlSurface    *color.Color
	InterfaceGray     *color.Color
	LightGray         *color.Color
	DisplayGray       *color.Color
	HighlightGray     *color.Color
	MetallicHighlight *color.Color
	AccentOrange      *color.Color
	AccentTeal        *color.Color
	AccentPrimary     *color.Color
	AccentSecondary   *color.Color
	AccentHighlight   *color.Color
	PhosphorGreen     *color.Color
	PhosphorAmber     *color.Color
	TextLight         *color.Color
	ScreenReflection  *color.Color
	CRTRedFringe      *color.Color
	CRTBlueFringe     *color.Color
}

// InterfaceZone defines a region of the cassette futurism interface
type InterfaceZone struct {
	X, Y, Width, Height int
	ZoneType           string // "main_display", "control_panel", "status_array", "terminal_screen"
}

// GenerateCassetteFuturismImage creates an enhanced cassette futurism aesthetic with authentic industrial interface elements
func GenerateCassetteFuturismImage(width, height int, accentHue float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Generate sophisticated temperature-matched palette
	palette := generateCassetteFuturismPalette(accentHue)
	
	// Fill background with dark industrial color
	fillBackground(img, 0, 0, width, height, palette.BackgroundDark)
	
	// Define interface layout zones using golden ratio proportions
	zones := generateInterfaceZones(width, height)
	
	// Render each interface zone with appropriate styling
	for _, zone := range zones {
		switch zone.ZoneType {
		case "main_display":
			renderMainDisplay(img, zone, palette)
		case "control_panel":
			renderControlPanel(img, zone, palette)
		case "status_array":
			renderStatusArray(img, zone, palette)
		case "terminal_screen":
			renderTerminalScreen(img, zone, palette)
		}
	}
	
	// Apply global effects
	applyCRTScanlines(img, 0, 0, width, height, palette)
	
	return img
}

// generateCassetteFuturismPalette creates temperature-matched industrial colors
func generateCassetteFuturismPalette(accentHue float64) *CassetteFuturismPalette {
	// Determine temperature-matched gray base
	var baseGrayHue float64
	if accentHue < 60.0/360.0 || accentHue > 300.0/360.0 {
		baseGrayHue = 30.0 / 360.0 // Warm grays for warm accents
	} else {
		baseGrayHue = 210.0 / 360.0 // Cool grays for cool accents
	}
	
	return &CassetteFuturismPalette{
		// Industrial grayscale progression with texture support
		BackgroundDark:    color.NewHSL(baseGrayHue, 0.02, 0.08),
		PanelDark:         color.NewHSL(baseGrayHue, 0.02, 0.18),
		PanelLight:        color.NewHSL(baseGrayHue, 0.02, 0.45),
		ControlSurface:    color.NewHSL(baseGrayHue, 0.02, 0.35),
		InterfaceGray:     color.NewHSL(baseGrayHue, 0.02, 0.50),
		LightGray:         color.NewHSL(baseGrayHue, 0.02, 0.65),
		DisplayGray:       color.NewHSL(baseGrayHue, 0.02, 0.80),
		HighlightGray:     color.NewHSL(baseGrayHue, 0.02, 0.92),
		MetallicHighlight: color.NewHSL(baseGrayHue, 0.05, 0.95),
		
		// Classic cassette futurism accent colors
		AccentOrange:      color.NewHSL(25.0/360.0, 0.90, 0.60), // #ff6b35 orange
		AccentTeal:        color.NewHSL(180.0/360.0, 0.80, 0.45), // #008080 teal
		
		// Strategic accent colors with industrial restraint
		AccentPrimary:     color.NewHSL(accentHue, 0.75, 0.55),
		AccentSecondary:   color.NewHSL(accentHue, 0.60, 0.40),
		AccentHighlight:   color.NewHSL(accentHue, 0.85, 0.70),
		
		// Authentic CRT phosphor colors
		PhosphorGreen:     color.NewHSL(120.0/360.0, 1.0, 0.60),
		PhosphorAmber:     color.NewHSL(45.0/360.0, 1.0, 0.65),
		
		// Interface elements
		TextLight:         color.NewHSL(baseGrayHue, 0.02, 0.90),
		ScreenReflection:  color.NewHSL(baseGrayHue, 0.10, 0.95),
		CRTRedFringe:      color.NewHSL(0.0, 0.60, 0.30),
		CRTBlueFringe:     color.NewHSL(240.0/360.0, 0.60, 0.30),
	}
}

// generateInterfaceZones creates realistic interface layout using golden ratio
func generateInterfaceZones(width, height int) []InterfaceZone {
	zones := []InterfaceZone{}
	
	// Main display (upper 38.2% using golden ratio)
	mainDisplayHeight := int(float64(height) * 0.382)
	zones = append(zones, InterfaceZone{
		X: width / 10, Y: height / 20,
		Width: width * 4 / 5, Height: mainDisplayHeight,
		ZoneType: "main_display",
	})
	
	// Control panel (left 61.8% of remaining space)
	controlY := mainDisplayHeight + height/10
	controlHeight := height - controlY - height/20
	controlWidth := int(float64(width) * 0.618)
	zones = append(zones, InterfaceZone{
		X: width / 20, Y: controlY,
		Width: controlWidth, Height: controlHeight,
		ZoneType: "control_panel",
	})
	
	// Status array (right side)
	statusX := controlWidth + width/10
	statusWidth := width - statusX - width/20
	zones = append(zones, InterfaceZone{
		X: statusX, Y: controlY,
		Width: statusWidth, Height: controlHeight / 2,
		ZoneType: "status_array",
	})
	
	// Terminal screen (lower right)
	terminalY := controlY + controlHeight/2 + height/40
	terminalHeight := height - terminalY - height/20
	zones = append(zones, InterfaceZone{
		X: statusX, Y: terminalY,
		Width: statusWidth, Height: terminalHeight,
		ZoneType: "terminal_screen",
	})
	
	return zones
}

// renderMainDisplay creates the primary interface display with data visualization
func renderMainDisplay(img *image.RGBA, zone InterfaceZone, palette *CassetteFuturismPalette) {
	// Display bezel with brushed metal texture
	drawBrushedMetalPanel(img, zone.X-2, zone.Y-2, zone.Width+4, zone.Height+4, palette.ControlSurface)
	
	// CRT screen with slight curvature simulation
	screenMargin := 8
	screenZone := InterfaceZone{
		X: zone.X + screenMargin, Y: zone.Y + screenMargin,
		Width: zone.Width - screenMargin*2, Height: zone.Height - screenMargin*2,
		ZoneType: "display",
	}
	
	// Fill screen with dark background
	fillRect(img, screenZone.X, screenZone.Y, screenZone.Width, screenZone.Height, palette.BackgroundDark)
	
	// Render oscilloscope-style waveform
	renderOscilloscopeWaveform(img, screenZone.X+10, screenZone.Y+10, screenZone.Width-20, screenZone.Height-20, palette)
	
	// Add data readouts
	renderDataReadouts(img, screenZone.X+10, screenZone.Y+screenZone.Height-40, screenZone.Width-20, 25, palette)
	
	// Screen reflection highlights
	addScreenReflection(img, screenZone.X, screenZone.Y, screenZone.Width, screenZone.Height, palette)
}

// renderControlPanel creates authentic control surface with buttons, sliders, and switches
func renderControlPanel(img *image.RGBA, zone InterfaceZone, palette *CassetteFuturismPalette) {
	// Panel background with subtle texture
	drawTexturedPanel(img, zone.X, zone.Y, zone.Width, zone.Height, palette.PanelDark)
	
	// Button matrix (4x6 grid)
	buttonSpacing := zone.Width / 7
	buttonSize := buttonSpacing / 3
	startX := zone.X + buttonSpacing/2
	startY := zone.Y + zone.Height/6
	
	for row := 0; row < 4; row++ {
		for col := 0; col < 6; col++ {
			x := startX + col*buttonSpacing
			y := startY + row*(zone.Height/6)
			
			// Alternate button states for realism
			pressed := (row+col)%3 == 0
			drawIndustrialButton(img, x, y, buttonSize, pressed, palette.ControlSurface)
			
			// LED indicators above buttons
			if row == 0 {
				active := (col+row)%4 == 0
				drawLEDIndicator(img, x+buttonSize/2, y-buttonSize, 2, active, palette.AccentOrange)
			}
		}
	}
	
	// Slider bank (right side)
	sliderX := zone.X + zone.Width*3/4
	renderSliderBank(img, sliderX, zone.Y+zone.Height/4, zone.Width/4-20, zone.Height/2, palette)
}

// renderStatusArray creates LED status indicators and seven-segment displays
func renderStatusArray(img *image.RGBA, zone InterfaceZone, palette *CassetteFuturismPalette) {
	// Status grid background
	drawTexturedPanel(img, zone.X, zone.Y, zone.Width, zone.Height, palette.PanelDark)
	
	// 8x4 LED status grid
	ledSpacing := zone.Width / 9
	ledSize := 3
	
	for row := 0; row < 4; row++ {
		for col := 0; col < 8; col++ {
			x := zone.X + (col+1)*ledSpacing
			y := zone.Y + zone.Height/6 + row*(zone.Height/6)
			
			// Simulate different system states
			active := (col*row + col + row) % 5 < 3
			ledColor := palette.AccentTeal
			if !active {
				ledColor = palette.PanelDark.AdjustLightness(-0.3)
			}
			drawLEDIndicator(img, x, y, ledSize, active, ledColor)
		}
	}
	
	// Seven-segment display
	displayY := zone.Y + zone.Height*2/3
	renderSevenSegmentDisplay(img, zone.X+10, displayY, zone.Width-20, zone.Height/4, 8, palette)
	renderSevenSegmentDisplay(img, zone.X+zone.Width/2-10, displayY, zone.Width-20, zone.Height/4, 4, palette)
}

// renderTerminalScreen creates authentic terminal interface with dot matrix text
func renderTerminalScreen(img *image.RGBA, zone InterfaceZone, palette *CassetteFuturismPalette) {
	// Terminal bezel
	drawTexturedPanel(img, zone.X-2, zone.Y-2, zone.Width+4, zone.Height+4, palette.ControlSurface)
	
	// Screen background
	fillRect(img, zone.X, zone.Y, zone.Width, zone.Height, palette.BackgroundDark)
	
	// Terminal text lines
	lines := []string{
		"SYSTEM STATUS: OPERATIONAL",
		"CPU LOAD: 23%  MEM: 45%",
		"TAPE DECK: READY",
		"LEVEL: -12dB  +/-  OK",
		"PROCESSES: 42  UP: 72:15",
		"> ANALYZING SIGNAL...",
		"FREQ RESPONSE: NOMINAL",
		"CALIBRATION: COMPLETE",
	}
	
	charHeight := 8
	charWidth := 6
	lineSpacing := charHeight + 2
	textX := zone.X + 8
	textY := zone.Y + 12
	
	for i, line := range lines {
		y := textY + i*lineSpacing
		if y+charHeight <= zone.Y+zone.Height-8 {
			renderDotMatrixText(img, textX, y, zone.Width-16, charHeight, line, palette)
		}
	}
	
	// Cursor blink simulation
	cursorX := textX + len("> ANALYZING SIGNAL...")*charWidth
	cursorY := textY + 5*lineSpacing
	drawLEDIndicator(img, cursorX, cursorY, 2, true, palette.PhosphorGreen)
}

// GenerateComplexGradientImage creates multi-stop gradients for extraction testing
func GenerateComplexGradientImage(width, height int, gradientType string) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	switch gradientType {
	case "linear-smooth":
		// Smooth HSL transition from blue to purple to pink
		for y := 0; y < height; y++ {
			progress := float64(y) / float64(height)
			
			// Smooth HSL interpolation
			hue := 0.67 - (progress * 0.25) // 240° to 180° (blue to purple to pink)
			sat := 0.8 + (progress * 0.2)   // 80% to 100%
			light := 0.3 + (progress * 0.4) // 30% to 70%
			
			c := color.NewHSL(hue, sat, light)
			r, g, b, a := c.RGBA()
			rgba := stdcolor.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
			
			for x := 0; x < width; x++ {
				img.Set(x, y, rgba)
			}
		}
		
	case "radial-complex":
		// Radial gradient with multiple color stops
		centerX, centerY := width/2, height/2
		maxRadius := math.Sqrt(float64(centerX*centerX + centerY*centerY))
		
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				dx := float64(x - centerX)
				dy := float64(y - centerY)
				distance := math.Sqrt(dx*dx + dy*dy)
				progress := distance / maxRadius
				
				// Multi-stop gradient: red -> orange -> yellow -> green -> blue
				var c *color.Color
				if progress < 0.2 {
					// Red to orange
					t := progress / 0.2
					hue := 0.0 + (t * 0.08) // 0° to 30°
					c = color.NewHSL(hue, 0.9, 0.5)
				} else if progress < 0.4 {
					// Orange to yellow  
					t := (progress - 0.2) / 0.2
					hue := 0.08 + (t * 0.08) // 30° to 60°
					c = color.NewHSL(hue, 0.9, 0.6)
				} else if progress < 0.7 {
					// Yellow to green
					t := (progress - 0.4) / 0.3
					hue := 0.17 + (t * 0.17) // 60° to 120°
					c = color.NewHSL(hue, 0.8, 0.5)
				} else {
					// Green to blue
					t := (progress - 0.7) / 0.3
					hue := 0.33 + (t * 0.34) // 120° to 240°
					c = color.NewHSL(hue, 0.7, 0.4)
				}
				
				r, g, b, a := c.RGBA()
				rgba := stdcolor.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
				img.Set(x, y, rgba)
			}
		}
		
	default: // "stepped-harsh"
		// Harsh stepped gradient with RGB transitions
		steps := 7
		stepHeight := height / steps
		
		colors := []*color.Color{
			color.NewRGB(255, 0, 0),   // Red
			color.NewRGB(255, 127, 0), // Orange
			color.NewRGB(255, 255, 0), // Yellow
			color.NewRGB(0, 255, 0),   // Green
			color.NewRGB(0, 255, 255), // Cyan
			color.NewRGB(0, 0, 255),   // Blue
			color.NewRGB(127, 0, 255), // Purple
		}
		
		for y := 0; y < height; y++ {
			stepIndex := y / stepHeight
			if stepIndex >= len(colors) {
				stepIndex = len(colors) - 1
			}
			
			c := colors[stepIndex]
			r, g, b, a := c.RGBA()
			rgba := stdcolor.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
			
			for x := 0; x < width; x++ {
				img.Set(x, y, rgba)
			}
		}
	}
	
	return img
}

// Helper functions for 80's vector graphics
func drawNeonLine(img *image.RGBA, x1, y1, x2, y2 int, neonColor *color.Color) {
	// Simple Bresenham line algorithm with neon effect
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	
	sx := 1
	if x1 > x2 {
		sx = -1
	}
	sy := 1
	if y1 > y2 {
		sy = -1
	}
	
	err := dx - dy
	x, y := x1, y1
	
	for {
		drawNeonPixel(img, x, y, neonColor)
		
		if x == x2 && y == y2 {
			break
		}
		
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

func drawNeonPixel(img *image.RGBA, x, y int, neonColor *color.Color) {
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return
	}
	
	// Create glow effect by setting multiple pixels with decreasing intensity
	r, g, b, a := neonColor.RGBA()
	coreColor := stdcolor.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	
	img.Set(x, y, coreColor)
	
	// Add subtle glow around the core pixel
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			
			gx, gy := x+dx, y+dy
			if gx >= bounds.Min.X && gx < bounds.Max.X && gy >= bounds.Min.Y && gy < bounds.Max.Y {
				// Lighter version for glow
				h, s, l := neonColor.HSL()
				glowColor := color.NewHSL(h, s*0.6, l*0.7)
				gr, gg, gb, ga := glowColor.RGBA()
				glow := stdcolor.RGBA{R: uint8(gr), G: uint8(gg), B: uint8(gb), A: uint8(ga)}
				
				// Blend with existing pixel
				existing := img.RGBAAt(gx, gy)
				if existing.R == 0 && existing.G == 0 && existing.B == 0 { // Only set on background
					img.Set(gx, gy, glow)
				}
			}
		}
	}
}

// fillBackground fills a rectangular region with a solid color
func fillBackground(img *image.RGBA, x, y, w, h int, c *color.Color) {
	rgba := toStdColor(c)
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			if x+dx < img.Bounds().Dx() && y+dy < img.Bounds().Dy() {
				img.Set(x+dx, y+dy, rgba)
			}
		}
	}
}

// fillRect fills a rectangle with enhanced industrial shading
func fillRect(img *image.RGBA, x, y, w, h int, baseColor *color.Color) {
	rgba := toStdColor(baseColor)
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			if x+dx < img.Bounds().Dx() && y+dy < img.Bounds().Dy() {
				img.Set(x+dx, y+dy, rgba)
			}
		}
	}
}

// drawBrushedMetalPanel creates a realistic brushed metal texture effect
func drawBrushedMetalPanel(img *image.RGBA, x, y, w, h int, baseColor *color.Color) {
	// Base metallic surface
	fillRect(img, x, y, w, h, baseColor)
	
	// Add brushed texture with horizontal lines
	for dy := 0; dy < h; dy += 2 {
		for dx := 0; dx < w; dx++ {
			if x+dx < img.Bounds().Dx() && y+dy < img.Bounds().Dy() {
				// Alternate between slightly darker and lighter tones
				var c *color.Color
				if (dx%4) < 2 {
					c = baseColor.AdjustLightness(-0.05)
				} else {
					c = baseColor.AdjustLightness(0.03)
				}
				img.Set(x+dx, y+dy, toStdColor(c))
			}
		}
	}
	
	// Add metallic highlights
	for dy := 0; dy < h; dy += 8 {
		for dx := 0; dx < w; dx += 3 {
			if x+dx < img.Bounds().Dx() && y+dy < img.Bounds().Dy() {
				c := baseColor.AdjustLightness(0.3)
				img.Set(x+dx, y+dy, toStdColor(c))
			}
		}
	}
}

// drawTexturedPanel creates an industrial plastic panel texture
func drawTexturedPanel(img *image.RGBA, x, y, w, h int, baseColor *color.Color) {
	// Base panel color
	fillRect(img, x, y, w, h, baseColor)
	
	// Add subtle texture pattern
	for dy := 0; dy < h; dy += 3 {
		for dx := 0; dx < w; dx += 3 {
			if x+dx < img.Bounds().Dx() && y+dy < img.Bounds().Dy() {
				// Create subtle texture variation
				lightness := -0.02 + (float64((dx+dy)%6)/6.0)*0.04
				c := baseColor.AdjustLightness(lightness)
				img.Set(x+dx, y+dy, toStdColor(c))
			}
		}
	}
}

// addScreenReflection simulates CRT screen reflection effect
func addScreenReflection(img *image.RGBA, x, y, w, h int, palette *CassetteFuturismPalette) {
	// Add curved reflection highlight in top-left
	for dy := 0; dy < h/3; dy++ {
		for dx := 0; dx < w/3; dx++ {
			if x+dx < img.Bounds().Dx() && y+dy < img.Bounds().Dy() {
				// Create curved reflection effect
				distance := float64(dx*dx + dy*dy)
				maxDist := float64((w/3)*(w/3) + (h/3)*(h/3))
				opacity := 0.3 * (1.0 - distance/maxDist)
				if opacity > 0 {
					c := palette.ScreenReflection.WithAlpha(opacity)
					// Blend with existing pixel
					existing := img.RGBAAt(x+dx, y+dy)
					blended := blendColors(colorFromRGBA(existing), &c)
					img.Set(x+dx, y+dy, toStdColor(blended))
				}
			}
		}
	}
}

// renderOscilloscopeWaveform creates a realistic oscilloscope display
func renderOscilloscopeWaveform(img *image.RGBA, x, y, w, h int, palette *CassetteFuturismPalette) {
	// Draw grid background
	gridSpacing := 20
	for gx := 0; gx <= w; gx += gridSpacing {
		for gy := 0; gy < h; gy++ {
			if x+gx < img.Bounds().Dx() && y+gy < img.Bounds().Dy() {
				c := palette.AccentOrange.WithAlpha(0.2)
				img.Set(x+gx, y+gy, toStdColor(&c))
			}
		}
	}
	for gy := 0; gy <= h; gy += gridSpacing {
		for gx := 0; gx < w; gx++ {
			if x+gx < img.Bounds().Dx() && y+gy < img.Bounds().Dy() {
				c := palette.AccentOrange.WithAlpha(0.2)
				img.Set(x+gx, y+gy, toStdColor(&c))
			}
		}
	}
	
	// Draw waveform
	centerY := y + h/2
	frequency := 0.1
	amplitude := float64(h) * 0.3
	
	for dx := 0; dx < w-1; dx++ {
		// Generate sine wave with some harmonics
		wave1 := math.Sin(float64(dx) * frequency)
		wave2 := math.Sin(float64(dx) * frequency * 3) * 0.3
		combined := wave1 + wave2
		
		y1 := centerY + int(combined*amplitude)
		y2 := centerY + int((math.Sin(float64(dx+1)*frequency) + math.Sin(float64(dx+1)*frequency*3)*0.3)*amplitude)
		
		// Draw line from current point to next
		drawBresenhamLine(img, x+dx, y1, x+dx+1, y2, palette.AccentOrange)
	}
	
	// Add glow effect around waveform
	for dx := 0; dx < w; dx++ {
		wave := math.Sin(float64(dx) * frequency) + math.Sin(float64(dx) * frequency * 3) * 0.3
		waveY := centerY + int(wave*amplitude)
		
		// Add glow pixels above and below the main line
		for glowOffset := -2; glowOffset <= 2; glowOffset++ {
			glowY := waveY + glowOffset
			if x+dx < img.Bounds().Dx() && glowY >= y && glowY < y+h {
				opacity := 0.3 * (1.0 - float64(absInt(glowOffset))/3.0)
				c := palette.AccentOrange.WithAlpha(opacity)
				existing := img.RGBAAt(x+dx, glowY)
				blended := blendColors(colorFromRGBA(existing), &c)
				img.Set(x+dx, glowY, toStdColor(blended))
			}
		}
	}
}

// renderDataReadouts creates realistic data display elements
func renderDataReadouts(img *image.RGBA, x, y, w, h int, palette *CassetteFuturismPalette) {
	// Background
	fillBackground(img, x, y, w, h, palette.PanelDark)
	
	// Add decorative elements
	for dx := 0; dx < w; dx += 40 {
		for dy := 0; dy < 3; dy++ {
			if x+dx < img.Bounds().Dx() && y+h-5+dy < img.Bounds().Dy() {
				img.Set(x+dx, y+h-5+dy, toStdColor(palette.AccentTeal))
			}
		}
	}
}

// drawIndustrialButton creates a realistic industrial button
func drawIndustrialButton(img *image.RGBA, x, y, size int, pressed bool, buttonColor *color.Color) {
	// Button base
	radius := size / 2
	
	for dy := 0; dy < size; dy++ {
		for dx := 0; dx < size; dx++ {
			distance := math.Sqrt(float64((dx-radius)*(dx-radius) + (dy-radius)*(dy-radius)))
			if distance <= float64(radius) {
				if x+dx < img.Bounds().Dx() && y+dy < img.Bounds().Dy() {
					c := buttonColor
					
					// Add depth shading
					if !pressed {
						// Raised button - lighter on top
						if dy < size/2 {
							c = buttonColor.AdjustLightness(0.2)
						} else {
							c = buttonColor.AdjustLightness(-0.1)
						}
					} else {
						// Pressed button - darker overall
						c = buttonColor.AdjustLightness(-0.3)
					}
					
					img.Set(x+dx, y+dy, toStdColor(c))
				}
			}
		}
	}
	
	// Add highlight ring
	for angle := 0.0; angle < 2*math.Pi; angle += 0.1 {
		ringX := x + radius + int(float64(radius-2)*math.Cos(angle))
		ringY := y + radius + int(float64(radius-2)*math.Sin(angle))
		if ringX >= 0 && ringX < img.Bounds().Dx() && ringY >= 0 && ringY < img.Bounds().Dy() {
			c := buttonColor.AdjustLightness(0.4)
			img.Set(ringX, ringY, toStdColor(c))
		}
	}
}

// drawLEDIndicator creates a realistic LED indicator
func drawLEDIndicator(img *image.RGBA, x, y, radius int, active bool, ledColor *color.Color) {
	for dy := 0; dy < radius*2; dy++ {
		for dx := 0; dx < radius*2; dx++ {
			distance := math.Sqrt(float64((dx-radius)*(dx-radius) + (dy-radius)*(dy-radius)))
			if distance <= float64(radius) {
				if x+dx < img.Bounds().Dx() && y+dy < img.Bounds().Dy() {
					var c *color.Color
					if active {
						// Active LED - bright with glow
						intensity := 1.0 - (distance / float64(radius))
						c = ledColor.AdjustLightness(intensity * 0.5)
					} else {
						// Inactive LED - dim
						c = ledColor.AdjustLightness(-0.7)
					}
					img.Set(x+dx, y+dy, toStdColor(c))
				}
			}
		}
	}
}

// renderSliderBank creates a bank of industrial sliders
func renderSliderBank(img *image.RGBA, x, y, w, h int, palette *CassetteFuturismPalette) {
	sliderCount := 8
	sliderWidth := (w - (sliderCount+1)*5) / sliderCount
	sliderHeight := h - 20
	
	for slider := 0; slider < sliderCount; slider++ {
		sliderX := x + 5 + slider*(sliderWidth+5)
		sliderY := y + 10
		
		// Slider track
		fillRect(img, sliderX, sliderY, sliderWidth, sliderHeight, palette.PanelDark)
		
		// Slider handle position (varies per slider)
		handlePos := float64(slider) / float64(sliderCount-1)
		handleY := sliderY + int(handlePos*float64(sliderHeight-20))
		
		// Draw slider handle
		fillRect(img, sliderX-2, handleY, sliderWidth+4, 15, palette.AccentTeal)
		
		// Add scale markings
		for mark := 0; mark < 5; mark++ {
			markY := sliderY + (mark * sliderHeight / 4)
			for px := 0; px < 3; px++ {
				if sliderX+sliderWidth+2+px < img.Bounds().Dx() && markY < img.Bounds().Dy() {
					img.Set(sliderX+sliderWidth+2+px, markY, toStdColor(palette.TextLight))
				}
			}
		}
	}
}

// renderSevenSegmentDisplay creates a 7-segment display digit
func renderSevenSegmentDisplay(img *image.RGBA, x, y, w, h, digit int, palette *CassetteFuturismPalette) {
	// Define 7-segment patterns for digits 0-9
	segments := [][]bool{
		{true, true, true, false, true, true, true},     // 0
		{false, false, true, false, false, true, false}, // 1
		{true, false, true, true, true, false, true},    // 2
		{true, false, true, true, false, true, true},    // 3
		{false, true, true, true, false, true, false},   // 4
		{true, true, false, true, false, true, true},    // 5
		{true, true, false, true, true, true, true},     // 6
		{true, false, true, false, false, true, false},  // 7
		{true, true, true, true, true, true, true},      // 8
		{true, true, true, true, false, true, true},     // 9
	}
	
	if digit < 0 || digit > 9 {
		return
	}
	
	pattern := segments[digit]
	segmentColor := palette.AccentOrange
	inactiveColor := palette.AccentOrange.AdjustLightness(-0.8)
	
	// Draw segments (simplified rectangles)
	segmentThickness := w / 8
	
	// Top horizontal (segment 0)
	c := inactiveColor
	if pattern[0] {
		c = segmentColor
	}
	fillRect(img, x+segmentThickness, y, w-2*segmentThickness, segmentThickness, c)
	
	// Top-right vertical (segment 1)
	c = inactiveColor
	if pattern[1] {
		c = segmentColor
	}
	fillRect(img, x+w-segmentThickness, y+segmentThickness, segmentThickness, h/2-segmentThickness, c)
	
	// Bottom-right vertical (segment 2)
	c = inactiveColor
	if pattern[2] {
		c = segmentColor
	}
	fillRect(img, x+w-segmentThickness, y+h/2, segmentThickness, h/2-segmentThickness, c)
	
	// Bottom horizontal (segment 3)
	c = inactiveColor
	if pattern[3] {
		c = segmentColor
	}
	fillRect(img, x+segmentThickness, y+h-segmentThickness, w-2*segmentThickness, segmentThickness, c)
	
	// Bottom-left vertical (segment 4)
	c = inactiveColor
	if pattern[4] {
		c = segmentColor
	}
	fillRect(img, x, y+h/2, segmentThickness, h/2-segmentThickness, c)
	
	// Top-left vertical (segment 5)
	c = inactiveColor
	if pattern[5] {
		c = segmentColor
	}
	fillRect(img, x, y+segmentThickness, segmentThickness, h/2-segmentThickness, c)
	
	// Middle horizontal (segment 6)
	c = inactiveColor
	if pattern[6] {
		c = segmentColor
	}
	fillRect(img, x+segmentThickness, y+h/2-segmentThickness/2, w-2*segmentThickness, segmentThickness, c)
}

// renderDotMatrixText creates dot matrix style text display
func renderDotMatrixText(img *image.RGBA, x, y, w, h int, text string, palette *CassetteFuturismPalette) {
	// Simple dot matrix rendering for terminal text
	charWidth := 6
	
	for i, char := range text {
		charX := x + i*charWidth
		if charX+charWidth > x+w {
			break
		}
		
		// Simple character rendering (just draw dots for readability)
		if char != ' ' {
			// Draw a simple dot pattern for non-space characters
			for row := 0; row < 5; row++ {
				for col := 0; col < 3; col++ {
					// Simple pattern for visibility
					if (row+col+int(char))%3 == 0 {
						dotX := charX + col
						dotY := y + row
						
						if dotX < img.Bounds().Dx() && dotY < img.Bounds().Dy() {
							img.Set(dotX, dotY, toStdColor(palette.PhosphorGreen))
						}
					}
				}
			}
		}
	}
}

// applyCRTScanlines adds authentic CRT display scanline effect
func applyCRTScanlines(img *image.RGBA, x, y, w, h int, palette *CassetteFuturismPalette) {
	// Add horizontal scanlines every 2 pixels
	for dy := 0; dy < h; dy += 2 {
		for dx := 0; dx < w; dx++ {
			if x+dx < img.Bounds().Dx() && y+dy < img.Bounds().Dy() {
				// Darken scanline
				existing := img.RGBAAt(x+dx, y+dy)
				existingColor := colorFromRGBA(existing)
				darkened := existingColor.AdjustLightness(-0.2)
				img.Set(x+dx, y+dy, toStdColor(darkened))
			}
		}
	}
	
	// Add subtle RGB color fringing effect
	for dy := 0; dy < h; dy += 4 {
		for dx := 1; dx < w-1; dx += 3 {
			if x+dx < img.Bounds().Dx() && y+dy < img.Bounds().Dy() {
				// Slight color separation
				img.Set(x+dx-1, y+dy, toStdColor(palette.CRTRedFringe))
				img.Set(x+dx+1, y+dy, toStdColor(palette.CRTBlueFringe))
			}
		}
	}
}

// Helper functions

// toStdColor converts our color.Color to standard library image/color.RGBA
func toStdColor(c *color.Color) stdcolor.RGBA {
	r, g, b, a := c.RGBA()
	return stdcolor.RGBA{R: r, G: g, B: b, A: a}
}

// drawBresenhamLine draws a line using Bresenham's algorithm
func drawBresenhamLine(img *image.RGBA, x0, y0, x1, y1 int, c *color.Color) {
	dx := absInt(x1 - x0)
	dy := absInt(y1 - y0)
	sx, sy := 1, 1
	
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	
	err := dx - dy
	x, y := x0, y0
	
	for {
		if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
			img.Set(x, y, toStdColor(c))
		}
		
		if x == x1 && y == y1 {
			break
		}
		
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// absInt returns the absolute value of an integer
func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// blendColors blends two colors with alpha transparency
func blendColors(base, overlay *color.Color) *color.Color {
	// Simple alpha blending
	baseH, baseS, baseL := base.HSL()
	overlayH, overlayS, overlayL := overlay.HSL()
	
	// Get alpha values from RGBA
	_, _, _, baseA := base.RGBA()
	_, _, _, overlayA := overlay.RGBA()
	
	alpha := float64(overlayA) / 255.0
	invAlpha := 1.0 - alpha
	
	// Blend in HSL space
	blendedH := baseH*invAlpha + overlayH*alpha
	blendedS := baseS*invAlpha + overlayS*alpha
	blendedL := baseL*invAlpha + overlayL*alpha
	blendedA := float64(baseA) / 255.0 // Keep original alpha
	
	return color.NewHSLA(blendedH, blendedS, blendedL, blendedA)
}

// colorFromRGBA converts image/color.RGBA to our color.Color
func colorFromRGBA(rgba stdcolor.RGBA) *color.Color {
	a := float64(rgba.A) / 255.0
	return color.NewRGBA(rgba.R, rgba.G, rgba.B, a)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}