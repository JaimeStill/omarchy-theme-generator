# Testing Strategy

## Philosophy

**"If you didn't run it, it doesn't work."**

Formal testing is deferred until project completion. During development, we use execution tests that validate specific technical concepts immediately with empirical results.

## Execution Test Pattern

### Structure

```go
// tests/test_[concept].go
package main

import (
    "fmt"
    "github.com/JaimeStill/omarchy-theme-generator/pkg/[package]"
)

func main() {
    // Minimal setup - no frameworks
    // Direct function call - immediate feedback
    // Clear output - visual validation
    fmt.Printf("Result: %v\n", result)
}
```

### Characteristics
- **Minimal**: No test framework overhead
- **Focused**: One concept per test
- **Direct**: Run with `go run`
- **Visual**: Results printed to stdout
- **Fast**: Immediate feedback loop

### Transparent Test Execution

All tests must provide complete visibility into their execution:

1. **Initial State**: Display starting values of all test variables
2. **Transformations**: Show operations being performed with parameters
3. **Expected vs Actual**: Display both values with exact measurements
4. **Rationale**: Explain WHY the test passes or fails

Example of poor test output:
```
AA compliance testing: ✗
```

Example of transparent test output:
```
AA compliance testing:
  Testing: RGB(119,119,119) on RGB(255,255,255) background
  Calculated contrast: 4.48:1
  Required for AA: 4.5:1
  Result: FAIL ✗ (4.48 < 4.5, difference: 0.02)
```

This principle ensures:
- Tests are self-documenting and educational
- Failures can be diagnosed without additional debugging
- The test suite serves as living documentation of the system's behavior
- Anyone running tests understands exactly what is being validated

## Test Categories

### 1. Algorithm Validation
Test specific algorithms with known inputs and expected outputs.

```go
// tests/test_octree_quantization.go
func main() {
    colors := generateKnownPalette()
    tree := buildOctree(colors, maxDepth=8)
    reduced := tree.GetPalette(16)
    
    fmt.Printf("Input: %d colors\n", len(colors))
    fmt.Printf("Output: %d colors\n", len(reduced))
    fmt.Printf("Valid distribution? %v\n", validateDistribution(reduced))
}
```

### 2. Performance Verification
Measure execution time against targets.

```go
// tests/test_performance.go
func main() {
    img := load4KImage()
    start := time.Now()
    
    palette := extractPalette(img)
    
    elapsed := time.Since(start)
    fmt.Printf("4K processing: %v (target: <2s)\n", elapsed)
    fmt.Printf("Pass? %v\n", elapsed < 2*time.Second)
}
```

### 3. Color Theory Validation
Verify palette generation strategies.

```go
// tests/test_triadic.go
func main() {
    base := color.NewRGB(255, 128, 0) // Orange
    palette := TriadicStrategy{}.Generate(base, 6)
    
    for i, c := range palette {
        h, s, l := c.HSL()
        fmt.Printf("%d: H=%.0f° S=%.0f%% L=%.0f%%\n", 
                   i, h*360, s*100, l*100)
    }
    
    // Verify 120° separation
    h0, _, _ := palette[0].HSL()
    h1, _, _ := palette[2].HSL()
    fmt.Printf("Angle separation: %.0f°\n", (h1-h0)*360)
}
```

### 4. Integration Tests
Test component interactions.

```go
// tests/test_theme_generation.go
func main() {
    img := loadTestImage()
    config := ThemeConfig{
        SourceImage: img,
        Mode: ModeAuto,
    }
    
    theme := GenerateTheme(config)
    err := theme.Export("./test-output")
    
    fmt.Printf("Files generated: %v\n", listFiles("./test-output"))
    fmt.Printf("Valid structure? %v\n", validateStructure("./test-output"))
}
```

## Test Organization

```
tests/
├── README.md                    # Test documentation
├── helpers.go                   # Shared test utilities
├── test-color/
│   ├── main.go                  # Color type operations
│   └── README.md                # Test documentation with latest output
├── test-conversions/
│   ├── main.go                  # Color space conversions  
│   └── README.md                # Test documentation with latest output
├── test-extract/                # Basic extraction (future)
├── test-strategies/             # Palette strategies (future)
├── test-octree/                 # Octree algorithm (future)
├── test-dominant/               # Dominant color detection (future)
├── test-concurrent/             # Parallel processing (future)
├── test-contrast/               # WCAG validation (future)
├── test-generate-configs/       # Config generation (future)
├── test-theme-overrides/        # User overrides (future)
└── test-full-theme/             # Complete generation (future)
```

## Running Tests

### Individual Test
```bash
go run tests/test-color/main.go
```

### Test with Arguments
```bash
go run tests/test-extract/main.go image.jpg
```

### Test with Flags
```bash
go run tests/test-theme-overrides/main.go \
    -image=sunset.jpg \
    -mode=dark \
    -primary=#ff6b35
```

### Validation Only
```bash
go vet ./...
```

## Expected Outputs

### Color Operations
```
RGB: #ff8040
HSL: 20°, 100%, 63%
RGB (converted back): #ff8040
✓ Roundtrip successful
```

### Extraction Performance
```
Loaded: 3840x2160 image
Colors extracted: 1,847,293
Reduced to: 16
Time: 1.23s
✓ Meets <2s requirement
```

### Palette Generation
```
Monochromatic Palette:
  0: #ff8040 (base)
  1: #ffb380 (lighter)
  2: #cc5020 (darker)
  ...
✓ Single hue maintained
```

## Success Criteria

Each test should verify:
1. **Correctness**: Output matches expected values
2. **Performance**: Execution time within targets
3. **Stability**: No crashes or panics
4. **Determinism**: Consistent results across runs

## Transition to Formal Testing

After development completion:
1. Convert execution tests to proper unit tests
2. Add test coverage metrics
3. Implement integration test suite
4. Add benchmark tests
5. Set up CI/CD pipeline

## References

- Development approach: `docs/development-methodology.md`
- Technical requirements: `docs/technical-specification.md`
- Progress tracking: `PROJECT.md`
