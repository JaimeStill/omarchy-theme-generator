# Implementation Roadmap - Architecture Refactoring

## Overview
This roadmap captures all architectural adjustments from the planning session, organized into discrete implementation phases that can be executed over multiple development sessions.

---

## Phase 1: Foundation Refactoring

### 1.1 pkg/color → pkg/formats
**Goal**: Simplify to only what's actually used, leverage standard library

**Actions**:
```go
// DELETE these unused features:
- LAB color space conversions
- Distance calculations (except what's needed)
- Color manipulation methods (lighten, darken, etc.)
- Cached HSLA (premature optimization)
- Custom Color type

// KEEP only these functions:
- RGBToHSL(c color.Color) (h, s, l float64)
- ContrastRatio(c1, c2 color.Color) float64
- ToHex(c color.Color) string
- ToHexA(c color.Color) string  
- ParseHex(hex string) (color.RGBA, error)
```

**Migration**:
```go
// Update all references:
tcolor.Color → color.RGBA
c.HEX() → formats.ToHex(c)
c.HSL() → formats.RGBToHSL(c)
```

### 1.2 Create pkg/settings and pkg/config
**Goal**: Separate system settings from user configuration

**pkg/settings** - System operation values:
```go
// settings/settings.go
type Settings struct {
    Extraction  ExtractionSettings
    Analysis    AnalysisSettings  
    Clustering  ClusteringSettings
    Synthesis   SynthesisSettings
    Accessibility AccessibilitySettings
}

// Load from: config/settings.json (global)
// Purpose: HOW the tool operates
```

**pkg/config** - User theme configuration:
```go
// config/preferences.go
type ThemePreferences struct {
    Mode           ThemeMode
    ColorOverrides map[ColorRole]string
    ExtractionHints ExtractionHints
    SchemePreferences SchemePreferences
}

// Load from: theme-name/theme-preferences.json
// Purpose: WHAT the user wants
```

---

## Phase 2: Extractor Decomposition

### 2.1 Extract pkg/analysis from extractor
**Goal**: Separate image analysis from color extraction

**Move to pkg/analysis**:
```go
// From extractor/image_analysis.go → analysis/image.go
- ImageCharacteristics struct
- AnalyzeImageCharacteristics()
- Edge detection functions
- Color complexity analysis

// From extractor/extractor.go → analysis/profile.go
- ImageColorProfile enum
- DetectColorProfile()
- IsGrayscale(), IsMonochromatic()
- Theme mode detection
```

### 2.2 Extract pkg/strategies from extractor
**Goal**: Make strategies pluggable and testable

**Move to pkg/strategies**:
```go
// strategies/interface.go
type Strategy interface {
    Extract(img image.Image, opts *Options) (*Result, error)
    CanHandle(characteristics *analysis.ImageCharacteristics) bool
    Priority(characteristics *analysis.ImageCharacteristics) int
    Name() string
}

// strategies/frequency.go (from extractor/strategy_frequency.go)
// strategies/saliency.go (from extractor/strategy_saliency.go)
// strategies/selector.go (from extractor/strategies.go)
```

### 2.3 Simplify pkg/extractor
**Goal**: Extractor becomes orchestrator only

**Keep in pkg/extractor**:
```go
// extractor/extractor.go
type Extractor struct {
    settings  *settings.Settings
    analyzer  *analysis.Analyzer
    selector  *strategies.Selector
}

func (e *Extractor) Extract(img image.Image) (*ThemeColorMap, error) {
    // 1. Analyze image
    // 2. Select strategy
    // 3. Extract colors
    // 4. Assign roles
    // 5. Return organized colors
}
```

---

## Phase 3: Purpose-Driven Extraction

### 3.1 Implement ThemeColorMap
**Goal**: Replace frequency-based with role-based organization

```go
// formats/theme.go
type ThemeColorMap struct {
    // Role-based organization
    Primary    map[ColorRole]color.RGBA
    Candidates map[ColorRole][]ScoredColor
    
    // Metadata
    Profile    ImageProfile
    Mode       ThemeMode
    Luminance  float64
    
    // Synthesis requirements
    NeedsSynthesis map[ColorRole]bool
    SynthesisBase  color.RGBA
}

type ScoredColor struct {
    Color       color.RGBA
    Frequency   float64
    Suitability float64
    Distance    float64
}
```

### 3.2 Mode-Aware Role Assignment
**Goal**: Assign colors to roles based on light/dark mode

```go
// analysis/roles.go
func AssignColorRoles(colors []color.RGBA, mode ThemeMode, settings *settings.Settings) map[ColorRole][]ScoredColor {
    // Dark mode: dark colors → backgrounds
    // Light mode: light colors → backgrounds
    // Saturated colors → accents
    // Apply perceptual clustering
}
```

### 3.3 Edge Case Synthesis
**Goal**: Handle minimal color images gracefully

```go
// palette/synthesis.go
func SynthesizeFromProfile(profile ImageProfile, baseColor color.RGBA) *ThemeColorMap {
    switch profile {
    case ProfileGrayscale:
        // Extract temperature → complementary primary
    case ProfileDuotone:
        // Use as anchors → synthesize rest
    case ProfileMonochromatic:
        // Preserve hue → add complements
    }
}
```

---

## Phase 4: Documentation Updates

### 4.1 Update CLAUDE.md
```markdown
## Current Implementation Status
- ✅ Multi-strategy extraction (frequency/saliency)
- ✅ Settings-driven configuration
- 🔄 Purpose-driven extraction (in progress)
- ⏳ Color scheme generation (pending)

## Key Technical Decisions
- Standard library color.RGBA instead of custom type
- Role-based color organization instead of frequency
- Separate settings (system) from config (user)
- Layered architecture with clear dependencies
```

### 4.2 Update PROJECT.md
Use new structure from artifact:
- Current Implementation
- Current Work
- Components & Features (by layer)
- Remove session-based planning

### 4.3 Update README.md
```markdown
## Architecture
- **pkg/formats** - Color conversion and formatting
- **pkg/analysis** - Image and color analysis
- **pkg/extractor** - Color extraction orchestration
- **pkg/strategies** - Extraction strategies
- **pkg/palette** - Color scheme generation
- **pkg/theme** - Theme file generation
- **pkg/settings** - System configuration
- **pkg/config** - User preferences
```

### 4.4 Update docs/
- `technical-specification.md` - Add purpose-driven extraction
- `palette-generation.md` - Update with role-based approach
- `architecture.md` - NEW: Document layered architecture

---

## Phase 5: Testing Updates

### 5.1 Refactor Existing Tests
```go
// Update imports
tcolor.Color → color.RGBA
extractor.AnalyzeImageCharacteristics → analysis.AnalyzeImage

// Update test structure
tests/
├── formats_test.go     # Test color conversions
├── analysis_test.go    # Test image analysis
├── strategies_test.go  # Test extraction strategies
├── roles_test.go       # Test role assignment (NEW)
└── synthesis_test.go   # Test edge cases (NEW)
```

### 5.2 Add New Tests
- Role assignment validation
- Mode detection accuracy
- Synthesis for edge cases
- Settings/config loading

---

## Implementation Order

### Session 1: Foundation
1. Create pkg/formats from pkg/color
2. Update all color references
3. Create pkg/settings and pkg/config
4. Run tests to ensure nothing breaks

### Session 2: Decomposition
1. Extract pkg/analysis from extractor
2. Extract pkg/strategies from extractor
3. Simplify pkg/extractor to orchestrator
4. Update imports and tests

### Session 3: Purpose-Driven
1. Implement ThemeColorMap structure
2. Add role assignment logic
3. Implement mode detection
4. Add perceptual clustering

### Session 4: Edge Cases
1. Implement profile detection
2. Add synthesis strategies
3. Test with minimal color images
4. Validate accessibility

### Session 5: Integration
1. Update all documentation
2. Create end-to-end tests
3. Build CLI commands
4. Final validation

---

## File Structure (Final State)

```
omarchy-theme-generator/
├── config/
│   └── settings.json           # Global system settings
├── pkg/
│   ├── formats/               # Data formatting and types
│   │   ├── color.go          # Color conversions
│   │   ├── theme.go          # ThemeColorMap structure
│   │   └── types.go          # ColorRole, ThemeMode, etc.
│   ├── analysis/              # Image and color analysis
│   │   ├── image.go          # Image characteristics
│   │   ├── profile.go        # Color profile detection
│   │   └── roles.go          # Role assignment
│   ├── strategies/            # Extraction strategies
│   │   ├── interface.go      # Strategy interface
│   │   ├── frequency.go      # Frequency strategy
│   │   ├── saliency.go       # Saliency strategy
│   │   └── selector.go       # Strategy selection
│   ├── extractor/             # Extraction orchestration
│   │   └── extractor.go      # Main extraction pipeline
│   ├── palette/               # Color scheme generation
│   │   ├── schemes.go        # Color theory schemes
│   │   └── synthesis.go      # Edge case synthesis
│   ├── theme/                 # Theme file generation
│   │   ├── templates/        # Config file templates
│   │   └── generator.go      # Template processing
│   ├── settings/              # System configuration
│   │   └── settings.go       # Settings structures
│   └── config/                # User preferences
│       └── preferences.go     # User config structures
├── cmd/
│   └── omarchy-theme-gen/
│       └── main.go
└── tests/
    ├── formats_test.go
    ├── analysis_test.go
    ├── strategies_test.go
    └── integration_test.go
```

---

## Success Criteria

### Code Quality
- [ ] No unused code in pkg/formats
- [ ] Clear dependency layers (no circular deps)
- [ ] All thresholds in settings, not hardcoded
- [ ] Standard library types where possible

### Functionality
- [ ] Handles all image types (grayscale to full-color)
- [ ] Mode-aware color assignment
- [ ] WCAG compliance validation
- [ ] User preferences override system

### Performance
- [ ] <2s for 4K images
- [ ] <100MB memory usage
- [ ] Efficient color clustering

### Documentation
- [ ] Clear component purposes
- [ ] Updated examples
- [ ] Migration guide from old structure

---

## Migration Checklist

### Breaking Changes
- `pkg/color` → `pkg/formats`
- Custom Color type → `color.RGBA`
- Frequency-based → Role-based extraction
- Settings location change

### Compatibility
- [ ] Provide migration tool for existing themes
- [ ] Document all breaking changes
- [ ] Update CLI to handle old/new formats

---

## Notes

### Why These Changes?
1. **Simplification**: Remove 90% of unused color code
2. **Standard Library**: Better interoperability
3. **Purpose-Driven**: Colors organized by role, not frequency
4. **Separation**: System settings vs user preferences
5. **Modularity**: Clear dependency layers

### Risk Mitigation
- Test each phase independently
- Keep old code until new code works
- Document all changes thoroughly
- Provide migration tools
