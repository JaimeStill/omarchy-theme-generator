# Extraction Architecture Decisions
*Session 4 Enhancement - Purpose-Driven Color Extraction*

## Executive Summary

Transform the color extraction system from frequency-based to **purpose-driven extraction** that categorizes colors by their intended role in the theme (backgrounds, foregrounds, accents) rather than just their frequency. This approach ensures usable, accessible themes with proper contrast and variety.

## Core Architectural Decisions

### 1. Purpose-Driven Extraction
**Decision**: Replace frequency-based "top colors" with role-based color categorization.

**Rationale**: 
- Frequency doesn't correlate with theme usability
- Top colors are often slight variations of each other
- UI themes need specific color roles (background, foreground, accents)

**Implementation**:
```go
type ColorRole string

const (
    // Core UI roles
    RoleBackground        ColorRole = "background"
    RoleBackgroundAlt     ColorRole = "background-alt"
    RoleForeground        ColorRole = "foreground"
    RoleForegroundDim     ColorRole = "foreground-dim"
    
    // Terminal palette roles
    RoleTerminalBlack     ColorRole = "terminal-black"
    RoleTerminalRed       ColorRole = "terminal-red"
    // ... other ANSI colors
    
    // UI element roles
    RoleBorder            ColorRole = "border"
    RoleCursor            ColorRole = "cursor"
    RoleSelection         ColorRole = "selection"
    
    // Semantic roles
    RoleSuccess           ColorRole = "success"
    RoleError             ColorRole = "error"
    RoleWarning           ColorRole = "warning"
    
    // Accent roles
    RolePrimary           ColorRole = "primary"
    RoleSecondary         ColorRole = "secondary"
)
```

### 2. Mode-Aware Role Assignment
**Decision**: Role assignment must adapt based on light/dark mode detection.

**Rationale**:
- Light themes need light backgrounds, dark foregrounds
- Dark themes need dark backgrounds, light foregrounds
- Modern themes often use saturated (non-gray) backgrounds

**Implementation**:
- Detect mode from average luminance
- Apply mode-specific thresholds for role assignment
- Allow saturated backgrounds within lightness constraints

### 3. Settings-Driven Architecture
**Decision**: All thresholds and parameters must be configurable, not hardcoded.

**Rationale**:
- Enables experimentation and tuning
- Settings document tool capabilities
- Follows established pattern in `pkg/extractor/settings.go`

**Implementation**:
```go
type ThemeSettings struct {
    Extraction     ExtractionSettings     // Existing
    ModeDetection  ModeDetectionSettings  // New
    RoleAssignment RoleAssignmentSettings // New
    Clustering     ClusteringSettings     // New
    Synthesis      SynthesisSettings      // New
    Accessibility  AccessibilitySettings  // New
}
```

### 4. User Preferences System
**Decision**: Separate user preferences from global settings, stored per theme.

**Rationale**:
- Users need fine-grained control over specific themes
- Preferences should persist with the theme
- Enables reproducible theme generation

**Files**:
- `config/theme-settings.json` - Global tool configuration
- `my-theme/theme-preferences.json` - User overrides for specific theme

### 5. Edge Case Handling
**Decision**: Detect and handle minimal color variation images with specialized strategies.

**Profiles**:
- **Grayscale**: Extract temperature → synthesize complementary primary
- **Duotone/Tritone**: Use existing colors as anchors → synthesize rest
- **Monochromatic**: Preserve base hue → add complementary accents
- **Sepia**: Treat as warm monochromatic

**Implementation**:
```go
type ImageColorProfile int

const (
    FullColor     ImageColorProfile = iota
    Duotone       // 2 distinct colors
    Tritone       // 3-4 distinct colors
    Monochromatic // Single hue variations
    Sepia         // Brown/orange monochromatic
    Grayscale     // No hue information
)
```

### 6. Enhanced Extraction Pipeline
**Decision**: Multi-stage pipeline with profile detection, role assignment, and synthesis.

**Stages**:
1. Raw color extraction (existing strategies)
2. Profile detection (grayscale, duotone, etc.)
3. Mode detection (light/dark/auto)
4. Role-based categorization
5. Profile-specific processing
6. Synthesis for missing roles
7. Accessibility validation

### 7. Perceptual Clustering
**Decision**: Use perceptual distance metrics to ensure "sufficiently different" colors.

**Thresholds**:
- Background variants: ΔE > 5.0 (noticeable difference)
- Accent colors: Hue > 30° OR Saturation > 0.3
- Foreground variants: Maintain hue/saturation, adjust lightness for WCAG

---

## Updated PROJECT.md Sessions

### Session 4: Enhanced Color Extraction System (Current)
- [x] Vocabulary corrections (IsGrayscale vs IsMonochromatic)
- [ ] **NEW**: Implement `ImageColorProfile` detection
- [ ] **NEW**: Create `ThemeSettings` structure with all configurable parameters
- [ ] **NEW**: Build mode-aware role assignment system
- [ ] **NEW**: Implement perceptual clustering with Delta-E
- [ ] **NEW**: Create `ThemeColorMap` structure for role-based organization
- [ ] **NEW**: Add synthesis strategies for edge cases (grayscale, duotone)
- [ ] **NEW**: Build user preferences system (`ThemePreferences`)
- [ ] **NEW**: Implement accessibility validation and adjustment
- [ ] **NEW**: Create `theme-preferences.json` persistence
- [ ] **Test**: `tests/test-purpose-extraction/main.go`

### Session 5: Color Theory Schemes with Enhanced Foundation
- [ ] Create pkg/palette/ package with `SchemeGenerator` interface
- [ ] Implement schemes using role-based colors from Session 4
- [ ] Build scheme application on `ThemeColorMap` structure
- [ ] Integrate with synthesis strategies for edge cases
- [ ] Add scheme preferences to `ThemePreferences`
- [ ] **Test**: `tests/test-color-schemes/main.go`

### Session 6: Template Generation with Role Mapping
- [ ] Create template generators using `ThemeColorMap`
- [ ] Implement role → config mapping for each format
- [ ] Generate `theme-gen.json` with full extraction metadata
- [ ] Generate `theme-preferences.json` for user overrides
- [ ] **Test**: `tests/test-template-generation/main.go`

---

## Updated Documentation Files

### docs/technical-specification.md - New Sections

#### Color Extraction Architecture

**Purpose-Driven Extraction**
Instead of extracting colors by frequency, the system categorizes colors by their intended role in the theme:
- Background candidates (based on mode and lightness)
- Foreground candidates (ensuring WCAG compliance)
- Accent colors (saturated, distinctive hues)
- Terminal colors (mapped to ANSI color positions)

**Mode-Aware Processing**
```go
// Dark mode: darker colors → backgrounds
if mode == ModeDark && l < 0.35 {
    roles[RoleBackground] = append(...)
}

// Light mode: lighter colors → backgrounds  
if mode == ModeLight && l > 0.75 {
    roles[RoleBackground] = append(...)
}
```

**Settings Architecture**
All extraction parameters are configurable through `ThemeSettings`:
- No hardcoded thresholds in business logic
- Settings document tool capabilities
- Easy experimentation and tuning
- Stored in `config/theme-settings.json`

### docs/palette-generation.md - Enhanced Sections

#### Role-Based Color Extraction

**Color Roles**
The extraction system assigns colors to specific roles based on their characteristics:
- `background` / `background-alt` - Main and alternate backgrounds
- `foreground` / `foreground-dim` - Text colors
- `terminal-*` - ANSI terminal colors
- `border`, `cursor`, `selection` - UI elements
- `success`, `error`, `warning` - Semantic colors
- `primary`, `secondary` - Accent colors

**Profile Detection**
Images are classified into profiles for specialized handling:
- **FullColor** - Normal extraction with role assignment
- **Grayscale** - Temperature analysis → complementary synthesis
- **Monochromatic** - Preserve hue → add contrasting accents
- **Duotone/Tritone** - Use as anchors → synthesize remaining

#### User Preferences System

**Theme-Specific Overrides**
Each theme can have a `theme-preferences.json` file:
```json
{
  "mode": "dark",
  "color_overrides": {
    "background": "#1a1a2e",
    "primary": "#e94560"
  },
  "extraction_hints": {
    "preferred_background": "#1a1a2e",
    "excluded_colors": ["#000000"]
  },
  "scheme_preferences": {
    "scheme": "complementary",
    "scheme_base": "primary"
  }
}
```

### README.md - New Features Section

## Advanced Features

### Purpose-Driven Color Extraction
The generator doesn't just find frequent colors—it intelligently categorizes them by their role in the theme:
- **Mode-aware** background and foreground selection
- **Perceptual clustering** ensures color variety
- **Accessibility-first** with automatic WCAG compliance
- **Edge case handling** for grayscale, duotone, and monochromatic images

### Customization Options

#### Global Settings
Configure extraction behavior in `config/theme-settings.json`:
```json
{
  "mode_detection": {
    "dark_mode_max_luminance": 0.4
  },
  "role_assignment": {
    "accent_min_saturation": 0.5
  },
  "clustering": {
    "background_cluster_threshold": 5.0
  }
}
```

#### Theme Preferences  
Override specific colors in `theme-name/theme-preferences.json`:
```bash
# Set preferred background
omarchy-theme-gen set-preference my-theme --preferred-background "#1a1a2e"

# Override specific role
omarchy-theme-gen override my-theme --role primary --color "#e94560"
```

---

## Implementation Priority

### Phase 1: Core Infrastructure (Session 4)
1. `ImageColorProfile` detection
2. `ThemeSettings` structure
3. Mode detection algorithm
4. Role assignment logic

### Phase 2: Processing Pipeline (Session 4)
1. Perceptual clustering
2. Edge case synthesis
3. Accessibility validation
4. Preferences system

### Phase 3: Integration (Session 5-6)
1. Color schemes using role-based colors
2. Template generation with role mapping
3. CLI commands for preferences

---

## File Structure Updates

```
omarchy-theme-generator/
├── config/
│   └── theme-settings.json         # Global settings (NEW)
├── pkg/
│   ├── extractor/
│   │   ├── profiles.go            # Profile detection (NEW)
│   │   ├── roles.go               # Role assignment (NEW)
│   │   ├── clustering.go          # Perceptual clustering (NEW)
│   │   ├── synthesis.go           # Edge case synthesis (NEW)
│   │   └── theme_extractor.go     # Enhanced pipeline (NEW)
│   └── preferences/
│       └── preferences.go         # User preferences (NEW)
└── tests/
    └── test-purpose-extraction/    # New test suite (NEW)
```

---

## Testing Strategy

### New Test Requirements
1. **Profile Detection**: Validate grayscale, duotone, monochromatic detection
2. **Role Assignment**: Verify mode-aware categorization
3. **Clustering**: Ensure perceptual difference thresholds
4. **Synthesis**: Test edge case handling
5. **Preferences**: Validate override system

### Test Data Additions
Add to `tests/images/`:
- Pure grayscale image
- Duotone graphic
- Monochromatic photo
- High-saturation background example

---

## Migration Notes

### Breaking Changes
- `ExtractionResult.TopColors` → `ThemeColorMap.Primary[role]`
- Frequency-based selection → Role-based selection
- Single extraction → Multi-stage pipeline

### Backwards Compatibility
- Keep `FrequencyMap` for raw data
- Provide migration utility for existing themes
- Support legacy CLI commands with deprecation warnings
