# Testing Strategy

## Philosophy

**"If you didn't run it, it doesn't work."**

The testing strategy combines focused unit tests using `go test` with real-world validation using diverse wallpaper images. Tests provide immediate feedback and serve as living documentation of system behavior.

## Current Testing Approach

### Real-World Unit Tests

```go
// tests/strategies_test.go
package extractor_test

import (
    "testing"
    "github.com/JaimeStill/omarchy-theme-generator/pkg/extractor"
)

func TestStrategySelection(t *testing.T) {
    testCases := []struct {
        image    string
        expected string
    }{
        {"nebula.jpeg", "saliency"},
        {"night-city.jpeg", "saliency"}, 
        {"grayscale.jpeg", "frequency"},
    }
    
    for _, tc := range testCases {
        t.Run(tc.image, func(t *testing.T) {
            result, err := extractor.ExtractColors(tc.image, options)
            if err != nil {
                t.Fatalf("Failed: %v", err)
            }
            if result.SelectedStrategy != tc.expected {
                t.Errorf("Expected %s, got %s", tc.expected, result.SelectedStrategy)
            }
        })
    }
}
```

### Characteristics
- **Standard Go testing**: Uses built-in `testing` package
- **Real images**: 15 diverse wallpaper samples provide validation
- **Empirical validation**: Tests actual behavior with real-world data
- **Comprehensive coverage**: Strategy selection, theme analysis, benchmarks

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

```
tests/
├── README.md                    # Test documentation and results
├── *_test.go                    # Unit test files using standard Go testing
├── images/                      # Test assets and data files
│   ├── README.md                # Generated analysis documentation
│   └── *.jpeg, *.png           # Test images for validation
├── analyze-images/              # Utility tools for test data analysis
│   ├── main.go                  # Image analysis utility
│   └── README.md                # Utility documentation
└── helpers/                     # Shared test utilities (optional)
    └── common.go                # Helper functions for tests
```

## Running Tests

### Standard Go Testing
```bash
# Run all tests
go test ./tests -v

# Run specific test suites
go test ./tests -run TestSpecificFunction -v

# Run benchmarks
go test ./tests -bench=. -v

# Run with race detection
go test ./tests -race -v
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
