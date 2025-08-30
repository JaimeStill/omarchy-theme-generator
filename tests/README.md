# Tests Directory

## Overview
This directory contains execution tests that validate specific technical concepts immediately with empirical results. Following the "If you didn't run it, it doesn't work" philosophy, these tests provide transparent output showing initial state, transformations, and detailed explanations.

## Testing Philosophy
- **Minimal Structure**: No test frameworks, direct Go execution
- **Focused Validation**: One concept per test directory
- **Immediate Feedback**: Run with simple `go run` commands
- **Educational Output**: Tests serve as living documentation
- **Transparent Execution**: Show calculations, expected vs actual, and reasoning

## Available Tests

### Core Functionality Tests
- **[test-color/](test-color/)** - Core Color type operations and conversions
- **[test-conversions/](test-conversions/)** - Advanced color manipulation and analysis

### Image Processing & Extraction Tests
- **[test-load-image/](test-load-image/)** - Image loading, extraction, and synthesis pipeline (Session 3)

### Color Synthesis & Classification Tests (Session 4)
- **[test-synthesis/](test-synthesis/)** - Color theory algorithms, WCAG compliance, and mathematical validation
- **[test-classification/](test-classification/)** - Grayscale vs monochromatic vs full-color classification accuracy
- **[test-generative/](test-generative/)** - Computational image generation, material simulation, and aesthetic validation

### Theme Orchestration Tests (Session 5)
- **[test-palette-strategies/](test-palette-strategies/)** - Complete theme generation with mode detection, overrides, and WCAG compliance

### Template Generation Tests (Session 6)
- **[test-generate-alacritty/](test-generate-alacritty/)** - Alacritty terminal template generation with strategy-aware color mapping, WCAG compliance, and performance validation
- **[test-generate-alacritty/](test-generate-alacritty/)** - Alacritty TOML generator with strategy-aware color mapping and WCAG compliance

### Planned Tests (Future Sessions)
- **test-octree/** - Octree quantization implementation (Session 7)
- **test-generate-configs/** - Multiple configuration generators (Sessions 11-15)
- **test-performance/** - Benchmarking and optimization validation (Session 28)

## Usage

### Run Individual Tests
```bash
# Core color functionality
go run tests/test-color/main.go

# Advanced color conversions and analysis
go run tests/test-conversions/main.go

# Image loading and synthesis pipeline
go run tests/test-load-image/main.go

# Color synthesis and mathematical validation
go run tests/test-synthesis/main.go

# Color classification accuracy
go run tests/test-classification/main.go

# Computational image generation and material simulation
go run tests/test-generative/main.go

# Theme generation with mode detection and overrides
go run tests/test-palette-strategies/main.go

# Template generation with Alacritty TOML output
go run tests/test-generate-alacritty/main.go
```

### Run All Current Tests
```bash
# Validate all functionality (Core + Synthesis + Classification + Generation + Theme + Template)
go run tests/test-color/main.go && \
go run tests/test-conversions/main.go && \
go run tests/test-load-image/main.go && \
go run tests/test-synthesis/main.go && \
go run tests/test-classification/main.go && \
go run tests/test-generative/main.go && \
go run tests/test-palette-strategies/main.go && \
go run tests/test-generate-alacritty/main.go
```

### Test with Custom Images
```bash
# Test extraction and synthesis with your own image
go run tests/test-load-image/main.go path/to/your/image.jpg
```

### Code Quality Validation
```bash
# Type checking and validation
go vet ./tests/...

# Code formatting
go fmt ./tests/...
```

## Test Structure

Each test directory contains:
- **main.go** - The executable test implementation
- **README.md** - Documentation with latest output and validation criteria

## Shared Utilities

The `helpers.go` file provides common test utilities:
- `CheckMark(condition bool)` - Returns ✓ or ✗ based on test condition
- Additional formatting and validation helpers as needed

## Expected Outputs

### Successful Test Characteristics
- **Clear Headers**: Each test section clearly labeled
- **Step-by-Step Output**: Show initial state, operations, and results
- **Validation Explanations**: Explain why tests pass or fail with specific values
- **Performance Metrics**: Include timing where relevant
- **Checkmark Indicators**: ✓ for passing tests, ✗ for failures with detailed reasons

### Example Test Output Format
```
=== Test Section Name ===
Initial state: [values]
Operation: [what is being tested]
Expected: [expected result]
Actual: [actual result]
Result: ✓ - [explanation of why it passes]
```

## Development Workflow

1. **Implementation**: User implements functionality based on Claude's guides
2. **Test Creation**: Claude creates or updates execution tests
3. **Validation**: User runs tests and reports results
4. **Documentation**: Test output captured in README.md files

## Session Integration

Tests are created progressively throughout development sessions:
- Each session includes relevant execution tests
- Tests validate session objectives and technical requirements
- README.md files updated with latest output after successful implementation

## Quality Standards

All tests must meet these criteria:
- **Correctness**: Output matches expected values and specifications
- **Performance**: Execution time within established targets
- **Stability**: No crashes, panics, or inconsistent behavior
- **Determinism**: Consistent results across multiple runs
- **Educational Value**: Clear explanations of what is being validated and why

This testing approach ensures immediate validation of implementations while serving as comprehensive documentation of the system's behavior and capabilities.