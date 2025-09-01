# Architecture Documentation

## System Architecture

The omarchy-theme-generator uses a layered architecture with clear dependencies and separation of concerns. Each layer depends only on layers below it, preventing circular dependencies and ensuring maintainable code.

```
┌─────────────────────────────────────────────────────────┐
│                Application Layer (CLI)                  │
├─────────────────────────────────────────────────────────┤
│                  Generation Layer                       │
│         ┌──────────────┐      ┌──────────────┐          │
│         │    schemes   │      │    theme     │          │
│         └──────────────┘      └──────────────┘          │
├─────────────────────────────────────────────────────────┤
│                  Processing Layer                       │
│         ┌──────────────┐      ┌──────────────┐          │
│         │   extractor  │      │  strategies  │          │
│         └──────────────┘      └──────────────┘          │
├─────────────────────────────────────────────────────────┤
│                   Analysis Layer                        │
│                ┌──────────────────┐                     │
│                │     analysis     │                     │
│                └──────────────────┘                     │
├─────────────────────────────────────────────────────────┤
│                 Foundation Layer                        │
│      ┌──────────┐   ┌──────────┐   ┌──────────┐         │
│      │ formats  │   │ settings │   │  config  │         │
│      └──────────┘   └──────────┘   └──────────┘         │
└─────────────────────────────────────────────────────────┘
```

## Package Responsibilities

### Foundation Layer

**pkg/formats**
- Color space conversions (RGB↔HSL)
- Format utilities (ToHex, ParseHex)
- Type definitions (ColorRole, ThemeMode)
- No dependencies - pure functions only

**pkg/settings**
- System configuration and tool behavior
- Multi-layer composition (defaults → system → user → workspace → env)
- Empirical thresholds and operational parameters
- Global tool configuration

**pkg/config**
- User preferences and theme-specific overrides
- Theme-gen.json integration for metadata
- Per-theme color overrides and extraction hints
- User customization layer

### Analysis Layer

**pkg/analysis**
- Image characteristic detection (edge density, complexity)
- Profile classification (Grayscale, Monotone, Monochromatic, Duotone/Tritone)
- Mode detection (light/dark based on luminance)
- Role assignment logic for purpose-driven extraction
- Perceptual clustering and color grouping

### Processing Layer

**pkg/extractor**
- Orchestrates extraction pipeline
- Coordinates strategies and analysis
- Aggregates results and validates output
- Simplified to orchestration after refactoring

**pkg/strategies**
- Pluggable extraction algorithms
- Strategy interface for extensibility
- Frequency strategy for simple images
- Saliency strategy for complex images
- Strategy selection based on image characteristics

### Generation Layer

**pkg/schemes**
- Color theory scheme generation
- Edge case synthesis for minimal-color images
- Color harmony validation and WCAG compliance
- Role-based scheme application

**pkg/theme**
- Template processing and theme file generation
- Configuration generation for supported formats
- Format-specific color conversion
- Metadata creation and management

### Application Layer

**cmd/omarchy-theme-gen**
- CLI interface and command handling
- User interaction and workflow management
- Integration of all lower layers
- Command implementation (generate, set-scheme, etc.)

## Data Flow

1. **Input** → Image file provided by user
2. **Analysis** → Image characteristics and profile detection
3. **Strategy Selection** → Choose optimal extraction algorithm
4. **Extraction** → Raw color data from image
5. **Role Assignment** → Categorize colors by purpose (background, foreground, accents)
6. **Profile Processing** → Apply profile-specific handling
7. **Calculation** → Calculate missing colors using color theory algorithms
8. **Scheme Application** → Apply color theory if requested
9. **Validation** → Ensure WCAG compliance and accessibility
10. **Generation** → Create theme configuration files
11. **Output** → Complete theme package with metadata

## Design Principles

### 1. Separation of Concerns
Each package has a single, well-defined responsibility:
- **Formats**: Data transformation only
- **Analysis**: Image understanding only  
- **Extraction**: Color gathering only
- **Generation**: Output creation only

### 2. Dependency Direction
Dependencies flow downward only:
- Higher layers depend on lower layers
- Lower layers never depend on higher layers
- Same-layer dependencies are minimized
- No circular dependencies allowed

### 3. Settings-Driven Configuration
- No hardcoded values in business logic
- All thresholds and parameters configurable
- Multi-layer settings composition
- Clear separation between system settings and user preferences

### 4. Standard Library First
- Use Go standard library types where possible
- Prefer `color.RGBA` over custom color types
- Minimize external dependencies
- Leverage proven, well-tested implementations

### 5. Purpose-Driven Organization
- Colors organized by role, not frequency
- Background/foreground/accent categorization
- Mode-aware role assignment (light/dark themes)
- Edge case handling through profiles

## Configuration Architecture

### Settings vs Config Distinction

**Settings** (`pkg/settings`) - **HOW** the tool operates:
- Extraction thresholds and parameters
- Algorithm behavior configuration  
- Performance and accuracy trade-offs
- System-wide operational values
- Multi-layer composition from various sources

**Config** (`pkg/config`) - **WHAT** the user wants:
- Theme-specific color overrides
- User preferences and customizations
- Per-theme extraction hints
- Output format preferences
- Stored with generated themes

### Settings Composition Order
```
defaults → system → user → workspace → env
```
Later values override earlier ones, allowing flexible configuration at multiple levels.

## Profile Detection System

### Image Profiles
- **Grayscale**: No hue information (s ≈ 0) → requires color synthesis
- **Monotone**: Single hue tinting throughout → extract tint, enhance variation  
- **Monochromatic**: Single dominant hue with pure grayscale elements → extract accent
- **Duotone/Tritone**: 2-3 distinct colors only → use as anchors, synthesize rest
- **Full Color**: Standard multi-color image → normal extraction pipeline

### Profile-Specific Processing
Each profile can specify its own processing pipeline while reusing common components. Profiles are designed to be extensible for future edge cases.

## Testing Strategy

### Package-Level Testing
- Standard Go test files (`*_test.go`)
- Unit tests for each package in isolation
- Mock dependencies for layer testing
- Comprehensive coverage of public APIs

### Integration Testing
- End-to-end pipeline validation
- Real image processing tests
- Performance benchmarking
- Cross-layer interaction validation

### Test Organization
- `tests/internal/` - Shared test utilities
- `tests/samples/` - Reusable test images
- Package-specific tests co-located with source
- Benchmark tests for performance validation

## Extension Points

### Adding New Profiles
1. Define profile detection logic in `pkg/analysis`
2. Implement profile-specific processing
3. Register with profile detector
4. Add corresponding tests

### Adding New Strategies  
1. Implement `Strategy` interface in `pkg/strategies`
2. Add strategy-specific configuration to settings
3. Update strategy selector logic
4. Validate with diverse test images

### Adding New Output Formats
1. Create format-specific template in `pkg/theme`
2. Implement color mapping for the format
3. Add format detection and selection
4. Test with generated themes

## Performance Considerations

### Current Targets
- 4K image processing: <2 seconds
- Memory usage: <100MB peak
- Extraction strategies: Multiple algorithms available
- WCAG compliance: Automatic validation

### Optimization Strategies
- Efficient color space operations
- Minimal memory allocations
- Concurrent processing where beneficial
- Early termination for edge cases

This architecture provides a solid foundation for the purpose-driven theme generation system while maintaining flexibility for future enhancements.
