package template

import (
	"fmt"
	"strings"
	"sync"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/color"
	"github.com/JaimeStill/omarchy-theme-generator/pkg/theme"
)

// ColorFormat represents different output formats for colors in configuration files.
type ColorFormat string

const (
	FormatHex     ColorFormat = "hex"      // #RRGGBB or #RRGGBBAA
	FormatCSS     ColorFormat = "css"      // rgb(r,g,b) or rgba(r,g,b,a)
	FormatHSL     ColorFormat = "hsl"      // hsl(h,s,l) or hsla(h,s,l,a)
	FormatQuoted  ColorFormat = "quoted"   // "value" (wraps any format in quotes)
	FormatUnixHex ColorFormat = "unix_hex" // 0xRRGGBB format for some configs
)

// String returns the string representation of the color format.
func (cf ColorFormat) String() string {
	return string(cf)
}

// TemplateData contains all color information pre-formatted for template execution.
// This eliminates the need for complex formatting logic within templates themselves.
type TemplateData struct {
	// Theme information
	Theme    *theme.Theme
	ThemeName string
	
	// Primary semantic colors (pre-formatted)
	Primary    FormattedColor
	Background FormattedColor
	Foreground FormattedColor
	
	// Terminal color mapping
	Terminal TerminalColorData
	
	// Metadata for template logic
	IsLight    bool
	IsDark     bool
	Strategy   string
	Generated  string // Human-readable generation timestamp
}

// FormattedColor contains a color in multiple output formats for template flexibility.
type FormattedColor struct {
	Hex     string // #RRGGBB
	HexAlpha string // #RRGGBBAA (if alpha != 255)
	CSS     string // rgb(r,g,b) or rgba(r,g,b,a)
	HSL     string // hsl(h,s,l) or hsla(h,s,l,a)
	Quoted  string // "hex_value"
	UnixHex string // 0xRRGGBB
	
	// Raw values for mathematical operations in templates
	R, G, B, A uint8
	H, S, L    float64
}

// TerminalColorData contains the complete 16-color terminal palette in various formats.
type TerminalColorData struct {
	Normal [8]FormattedColor // ANSI colors 0-7
	Bright [8]FormattedColor // ANSI colors 8-15
	
	// Named access for common terminal colors
	Black   FormattedColor
	Red     FormattedColor
	Green   FormattedColor
	Yellow  FormattedColor
	Blue    FormattedColor
	Magenta FormattedColor
	Cyan    FormattedColor
	White   FormattedColor
}

// ColorFormatter handles the conversion of colors into various string formats
// optimized for template execution performance.
type ColorFormatter struct {
	// DefaultFormat specifies the default format when no format is specified
	DefaultFormat ColorFormat
	
	// Pool for reusing string builders to reduce allocations
	builderPool sync.Pool
}

// NewColorFormatter creates a formatter with the specified default format.
func NewColorFormatter(defaultFormat ColorFormat) *ColorFormatter {
	cf := &ColorFormatter{
		DefaultFormat: defaultFormat,
	}
	
	// Initialize string builder pool
	cf.builderPool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
	
	return cf
}

// FormatColor converts a color to the specified format.
func (cf *ColorFormatter) FormatColor(c *color.Color, format ColorFormat) string {
	if c == nil {
		return ""
	}
	
	switch format {
	case FormatHex:
		return c.HEX()
	case FormatCSS:
		return c.CSSRGBA()
	case FormatHSL:
		return c.CSSHSLA()
	case FormatQuoted:
		hex := c.HEX()
		return fmt.Sprintf(`"%s"`, hex)
	case FormatUnixHex:
		return cf.formatUnixHex(c)
	default:
		return cf.FormatColor(c, cf.DefaultFormat)
	}
}

// formatUnixHex formats a color as 0xRRGGBB for Unix-style configurations.
func (cf *ColorFormatter) formatUnixHex(c *color.Color) string {
	// Use direct RGB fields which are already uint8
	return fmt.Sprintf("0x%02x%02x%02x", c.R, c.G, c.B)
}

// CreateFormattedColor converts a color into a FormattedColor with all formats pre-computed.
func (cf *ColorFormatter) CreateFormattedColor(c *color.Color) FormattedColor {
	if c == nil {
		return FormattedColor{}
	}
	
	h, s, l := c.HSL()
	
	// Use direct RGB fields which are already uint8
	r8 := c.R
	g8 := c.G
	b8 := c.B
	a8 := c.A
	
	hex := c.HEX()
	
	return FormattedColor{
		Hex:     hex,
		HexAlpha: cf.getHexWithAlpha(c),
		CSS:     c.CSSRGBA(),
		HSL:     c.CSSHSLA(),
		Quoted:  fmt.Sprintf(`"%s"`, hex),
		UnixHex: cf.formatUnixHex(c),
		R:       r8,
		G:       g8,
		B:       b8,
		A:       a8,
		H:       h,
		S:       s,
		L:       l,
	}
}

// getHexWithAlpha returns hex format including alpha channel if not fully opaque.
func (cf *ColorFormatter) getHexWithAlpha(c *color.Color) string {
	if c.A == 255 {
		return c.HEX() // No alpha needed
	}
	
	// Include alpha in hex format using direct fields
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}

// PrepareTemplateData creates a complete TemplateData structure from a theme,
// with all colors pre-formatted for efficient template execution.
func (cf *ColorFormatter) PrepareTemplateData(t *theme.Theme) (*TemplateData, error) {
	if t == nil {
		return nil, fmt.Errorf("theme is nil")
	}
	
	// Create terminal color mapping
	mapper := NewTerminalMapper()
	terminalMapping, err := mapper.MapToTerminal(t)
	if err != nil {
		return nil, fmt.Errorf("terminal color mapping failed: %w", err)
	}
	
	// Prepare template data
	data := &TemplateData{
		Theme:      t,
		ThemeName:  t.Name,
		Primary:    cf.CreateFormattedColor(t.Primary),
		Background: cf.CreateFormattedColor(t.Background),
		Foreground: cf.CreateFormattedColor(t.Foreground),
		Terminal:   cf.createTerminalColorData(terminalMapping),
		IsLight:    t.IsLight,
		IsDark:     !t.IsLight,
		Strategy:   t.Metadata.Strategy,
		Generated:  t.Metadata.Generated.Format("2006-01-02 15:04:05"),
	}
	
	return data, nil
}

// createTerminalColorData converts the terminal color mapping into formatted color data.
func (cf *ColorFormatter) createTerminalColorData(mapping TerminalColorMap) TerminalColorData {
	data := TerminalColorData{}
	
	// Format normal colors (0-7)
	normalColors := mapping.GetNormal()
	for i, c := range normalColors {
		data.Normal[i] = cf.CreateFormattedColor(c)
	}
	
	// Format bright colors (8-15)
	brightColors := mapping.GetBright()
	for i, c := range brightColors {
		data.Bright[i] = cf.CreateFormattedColor(c)
	}
	
	// Set named colors for easy template access
	data.Black = data.Normal[0]
	data.Red = data.Normal[1]
	data.Green = data.Normal[2]
	data.Yellow = data.Normal[3]
	data.Blue = data.Normal[4]
	data.Magenta = data.Normal[5]
	data.Cyan = data.Normal[6]
	data.White = data.Normal[7]
	
	return data
}

// TemplateDataPool provides object pooling for TemplateData to reduce allocations
// during frequent template generation operations.
type TemplateDataPool struct {
	pool sync.Pool
}

// NewTemplateDataPool creates a new template data pool.
func NewTemplateDataPool() *TemplateDataPool {
	return &TemplateDataPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &TemplateData{}
			},
		},
	}
}

// Get retrieves a TemplateData instance from the pool.
func (tdp *TemplateDataPool) Get() *TemplateData {
	return tdp.pool.Get().(*TemplateData)
}

// Put returns a TemplateData instance to the pool for reuse.
func (tdp *TemplateDataPool) Put(data *TemplateData) {
	// Reset the data to avoid memory leaks
	*data = TemplateData{}
	tdp.pool.Put(data)
}

// Global pool instance for template data
var globalTemplateDataPool = NewTemplateDataPool()

// GetPooledTemplateData retrieves a TemplateData instance from the global pool
// and populates it with the given theme data.
func GetPooledTemplateData(t *theme.Theme) (*TemplateData, error) {
	formatter := NewColorFormatter(FormatHex)
	data := globalTemplateDataPool.Get()
	
	// Populate with theme data
	populatedData, err := formatter.PrepareTemplateData(t)
	if err != nil {
		globalTemplateDataPool.Put(data)
		return nil, err
	}
	
	*data = *populatedData
	return data, nil
}

// ReturnPooledTemplateData returns a TemplateData instance to the global pool.
func ReturnPooledTemplateData(data *TemplateData) {
	if data != nil {
		globalTemplateDataPool.Put(data)
	}
}