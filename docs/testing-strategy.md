# Testing Strategy

## Philosophy

**"If you didn't run it, it doesn't work."**

The testing strategy combines focused unit tests using `go test` with real-world validation using diverse wallpaper images. Tests provide immediate feedback and serve as living documentation of system behavior.

## Current Testing Approach

### Package-Level Unit Tests

```go
// tests/formats_test.go
package formats_test

import (
    "image/color"
    "testing"
    "github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

func TestColorConversion(t *testing.T) {
    testCases := []struct {
        name     string
        input    color.RGBA
        expected struct{ h, s, l float64 }
    }{
        {"red", color.RGBA{255, 0, 0, 255}, struct{ h, s, l float64 }{0.0, 1.0, 0.5}},
        {"gray", color.RGBA{128, 128, 128, 255}, struct{ h, s, l float64 }{0.0, 0.0, 0.5}},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            h, s, l := formats.RGBToHSL(tc.input)
            if h != tc.expected.h || s != tc.expected.s || l != tc.expected.l {
                t.Errorf("Expected HSL(%.1f, %.1f, %.1f), got HSL(%.1f, %.1f, %.1f)", 
                         tc.expected.h, tc.expected.s, tc.expected.l, h, s, l)
            }
        })
    }
}
```

### Characteristics
- **Standard Go testing**: Uses built-in `testing` package with `*_test.go` files
- **Layered testing**: Each package tested in isolation with clear dependencies
- **Real world validation**: Integration tests with actual image samples
- **Comprehensive coverage**: Unit tests for all public APIs and critical functions

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
Test core algorithms with known inputs and expected outputs.

```go
func TestAlgorithmBehavior(t *testing.T) {
    testCases := []struct {
        name     string
        input    InputType
        expected OutputType
    }{
        // Test cases with predictable results
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := Algorithm(tc.input)
            if !reflect.DeepEqual(result, tc.expected) {
                t.Errorf("Expected %v, got %v", tc.expected, result)
            }
        })
    }
}
```

### 2. Performance Verification
Measure execution time against requirements.

```go
func BenchmarkSystemPerformance(b *testing.B) {
    input := prepareTestInput()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        result := ProcessInput(input)
        if result == nil {
            b.Fatal("Unexpected nil result")
        }
    }
}
```

### 3. Behavioral Validation
Verify system behavior with real-world scenarios.

```go
func TestSystemBehavior(t *testing.T) {
    // Test with diverse inputs that represent actual usage
    testInputs := prepareRealWorldInputs()
    
    for _, input := range testInputs {
        t.Run(input.Name, func(t *testing.T) {
            result, err := ProcessInput(input.Data)
            if err != nil {
                t.Fatalf("Processing failed: %v", err)
            }
            
            if !validateOutput(result) {
                t.Error("Output validation failed")
            }
        })
    }
}
```

### 4. Integration Tests
Test interactions between major components.

```go
func TestComponentIntegration(t *testing.T) {
    // Setup test environment
    system := SetupTestSystem()
    defer system.Cleanup()
    
    // Test full workflow
    input := prepareTestInput()
    result, err := system.ProcessComplete(input)
    
    if err != nil {
        t.Fatalf("Integration failed: %v", err)
    }
    
    if !validateIntegrationResult(result) {
        t.Error("Integration result validation failed")
    }
}
```

## Test Organization

### Package-Specific Test Structure
Tests are organized by package in the tests/ directory:

```
tests/
├── formats/                     # Unit tests for pkg/formats
│   ├── formats_hsla_test.go     # HSLA color space conversions
│   ├── formats_hex_test.go      # Hex color parsing and formatting
│   ├── formats_contrast_test.go # WCAG accessibility calculations
│   └── formats_analysis_test.go # Color analysis utilities
├── extractor/                   # Unit tests for pkg/extractor
│   └── strategies_test.go       # Strategy selection and analysis
├── images/                      # Real-world wallpaper test images
│   ├── README.md                # Image analysis documentation
│   ├── grayscale.jpeg           # Pure grayscale test image
│   ├── nebula.jpeg              # Complex space image
│   ├── night-city.jpeg          # High-detail urban scene
│   ├── mountains.jpeg           # Natural landscape
│   ├── abstract.jpeg            # Abstract art
│   └── *.jpg, *.png             # Additional test wallpapers
└── analyze-images/              # Test analysis utility
    ├── main.go                  # Image characteristic analysis tool
    └── README.md                # Utility documentation
```

## Running Tests

### Standard Go Testing
```bash
# Run all unit tests
go test ./tests/formats ./tests/extractor -v

# Run specific package tests
go test ./tests/formats -run TestParseHex -v
go test ./tests/extractor -run TestStrategySelection -v

# Run with race detection
go test ./tests/formats ./tests/extractor -race -v
```

### Code Validation
```bash
# Check for compilation errors and vet issues
go vet ./...

# Format code consistently
go fmt ./...
```

### Utility Tools
```bash
# Generate test documentation
go run tests/analyze-images/main.go
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

## Test Coverage Goals

### Unit Test Coverage
- **pkg/formats**: 100% coverage of public functions
- **pkg/analysis**: Profile detection and mode classification
- **pkg/strategies**: Strategy selection and extraction algorithms  
- **pkg/extractor**: Pipeline orchestration
- **pkg/schemes**: Color theory scheme generation
- **pkg/theme**: Template processing and output generation

### Integration Test Coverage
- End-to-end extraction pipeline with real images
- Complete theme generation workflow
- Settings and configuration integration
- Profile detection with diverse image types

### Benchmark Coverage
- Color extraction performance with 4K images
- Color space conversion efficiency
- Profile detection speed
- Memory usage optimization

## References

- Development approach: `docs/development-methodology.md`
- Architecture details: `docs/architecture.md`
- Progress tracking: `PROJECT.md`
