# Glossary

## Color Theory and Image Processing

### Color Profiles
- **Grayscale**: Images with no hue information (saturation ≈ 0). Pure black, white, and gray tones only.
- **Monotone**: Images with single hue tinting throughout (like sepia, cyanotype). No pure grayscale values present.
- **Monochromatic**: Images with single dominant hue but including pure grayscale elements (noir-style with selective color).
- **Duotone**: Images containing only 2 distinct colors.
- **Tritone**: Images containing only 3-4 distinct colors.
- **Full Color**: Standard multi-color images with diverse hue, saturation, and lightness values.

### Color Organization
- **Frequency-Based**: Traditional approach organizing colors by occurrence count in image.
- **Purpose-Driven**: Refactored approach organizing colors by their intended role in the theme.
- **Role-Based**: Colors categorized by function (background, foreground, accent) rather than frequency.

### Color Roles
- **Background**: Colors suitable for window/terminal backgrounds based on mode and lightness.
- **Foreground**: Colors suitable for text with proper contrast ratios.
- **Primary**: Main accent color used as basis for color scheme operations.
- **Secondary**: Supporting accent color for variety.
- **Terminal Colors**: ANSI color palette mapping for terminal applications.

### Color Spaces and Formats
- **RGB**: Red, Green, Blue color space - standard for digital displays.
- **HSL**: Hue, Saturation, Lightness - intuitive for color manipulation and analysis.
- **LAB**: Perceptually uniform color space for accurate color distance calculations.
- **HEXA**: Hexadecimal color format with alpha channel (#RRGGBBAA).

## Architecture and Development

### Package Organization
- **Foundation Layer**: Core utilities with no dependencies (formats, settings, config).
- **Analysis Layer**: Image understanding and profile detection.
- **Processing Layer**: Color extraction algorithms and orchestration.
- **Generation Layer**: Theme creation and output formatting.
- **Application Layer**: CLI interface and user interaction.

### Package Names
- **pkg/formats**: Handles color conversions and formatting (refactored from pkg/color).
- **pkg/schemes**: Color theory scheme generation (not pkg/palette for consistency).
- **pkg/analysis**: Image and color analysis (extracted from extractor).
- **pkg/strategies**: Pluggable extraction algorithms (extracted from extractor).

### Configuration Types
- **Settings**: System configuration controlling HOW the tool operates. Multi-layer composition from defaults → system → user → workspace → env.
- **Config**: User preferences controlling WHAT the user wants. Theme-specific overrides and customizations.
- **Preferences**: User-specific theme modifications stored with generated themes.

### Extraction Strategies
- **Frequency Strategy**: Color extraction based on occurrence count. Suitable for simple images.
- **Saliency Strategy**: Color extraction based on visual importance and contrast. Suitable for complex images.
- **Strategy Selection**: Automatic algorithm choice based on image characteristics.

## Development Methodology

### Testing Organization
- **Unit Tests**: Package-level tests in tests/ directory, not embedded with source code.
- **Integration Tests**: End-to-end workflow tests in tests/integration/.
- **Benchmark Tests**: Performance validation in tests/benchmarks/.
- **Test Images**: Real-world wallpapers in tests/images/ for realistic validation.

### Development Principles
- **Standard Library First**: Prefer Go standard types (color.RGBA) over custom implementations.
- **Settings-Driven**: All thresholds configurable, no hardcoded values in business logic.
- **Dependency Direction**: Higher layers depend on lower layers only, no circular dependencies.
- **User-Driven**: All code modifications require explicit user direction.

### Color Science Terminology
- **Contrast Ratio**: WCAG-defined measurement for text accessibility (target: 4.5:1 for AA compliance).
- **Delta-E**: Perceptual color difference measurement in LAB color space.
- **Luminance**: Perceived brightness of a color, used for mode detection.
- **Saturation Threshold**: Boundary for grayscale detection (typically < 0.05).
- **Hue Tolerance**: Angular range for monochromatic detection (±15° typical).

## Color Theory Schemes

### Standard Schemes
- **Monochromatic**: Single hue with saturation and lightness variations.
- **Analogous**: Adjacent hues (±30° on color wheel).
- **Complementary**: Opposite hues (180° separation).
- **Split-Complementary**: Base hue plus two adjacent to complement.
- **Triadic**: Three hues equally spaced (120° separation).
- **Tetradic**: Rectangle pattern on color wheel.
- **Square**: Four hues equally spaced (90° separation).

### Scheme Application
- **Calculation**: Computing missing colors using color theory algorithms and existing color data.
- **WCAG Compliance**: Automatic validation and adjustment for accessibility.
- **Mode-Aware**: Scheme adaptation based on light/dark theme mode.

## Theme Generation

### Output Formats
- **theme-gen.json**: Metadata file containing extraction data and user preferences.
- **Omarchy Integration**: Direct compatibility with Omarchy desktop environment.
- **Multi-Format**: Support for terminal emulators, editors, and system components.

### Theme Components
- **Background Colors**: Primary and alternate background colors.
- **Foreground Colors**: Text colors with proper contrast validation.
- **ANSI Colors**: Standard terminal color palette (16 colors).
- **Semantic Colors**: Success, error, warning indicators.
- **UI Elements**: Border, cursor, selection colors.

## Performance and Validation

### Performance Targets
- **4K Processing**: <2 seconds for complete extraction pipeline.
- **Memory Usage**: <100MB peak during extraction.
- **WCAG Compliance**: Automatic AA-level contrast validation (4.5:1 ratio).

### Quality Metrics
- **Extraction Accuracy**: Successful color extraction from diverse image types.
- **Theme Usability**: Generated themes suitable for extended use.
- **Accessibility**: Compliance with web accessibility standards.

### Color Science
- **Delta-E (ΔE)**: Perceptual color difference measurement in CIELAB color space.
- **LAB Color Space**: Device-independent color model based on human vision.
- **D65 Illuminant**: Standard daylight illuminant used for color calculations.
- **Relative Luminance**: Brightness measurement for contrast ratio calculations.
- **Gamma Correction**: Nonlinear transformation for perceptually uniform color representation.

## Deprecated/Changed Terms

### Terminology Updates
- **Sepia** → **Monotone**: More accurate and broadly applicable term.
- **pkg/palette** → **pkg/schemes**: Consistency with "color schemes" vocabulary.
- **Top Colors** → **Role-Based Colors**: Purpose-driven organization.
- **Session-Based Planning** → **Component-Based Architecture**: Layer-focused development.

### Architectural Changes
- **Custom Color Type** → **color.RGBA**: Standard library adoption.
- **Execution Tests** → **Standard Go Tests**: Conventional test organization.
- **Embedded Tests** → **Isolated Tests**: All tests in tests/ directory.
- **TUI-First** → **CLI-First**: Command-line interface priority.

This glossary provides consistent terminology across all project documentation and development discussions.