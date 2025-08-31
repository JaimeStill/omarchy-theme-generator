# Color Type Test

## Overview
Comprehensive test of the core `Color` type functionality, validating RGBA storage, HSLA conversions, alpha handling, thread safety, and format output methods.

## Purpose
This test validates the fundamental color type implementation that serves as the foundation for all color operations in the theme generator. It ensures:

- **Color Creation**: RGBA and HSLA constructors work correctly
- **Format Output**: Hex, CSS RGB/RGBA, and CSS HSL/HSLA formats
- **Color Parsing**: HEXA (8-digit) and HEX (6-digit) string parsing with validation
- **Thread Safety**: Cached HSLA conversion works safely across goroutines
- **Alpha Handling**: Proper alpha storage, manipulation, and immutability
- **Conversion Accuracy**: Round-trip RGB↔HSLA conversions maintain precision

## Usage
```bash
go run tests/test-color/main.go
```

## Expected Results
- ✅ All 13 test sections should pass with checkmarks
- ✅ Round-trip conversions should show ≤1 unit differences (within tolerance)
- ✅ Thread safety test should show identical HSLA values across all goroutines
- ✅ Immutability tests should confirm original colors remain unchanged
- ✅ HEXA parsing should correctly handle 8-digit hex strings with alpha
- ✅ HEX parsing should correctly handle 6-digit hex strings (full opacity)
- ✅ Color parsing error handling should catch invalid formats

## Latest Test Output
```
=== Omarchy Theme Generator - Color Type Test ===

Test 1: RGBA Color Creation
Red RGBA: R=255, G=0, B=0, A=255
Red HEX (no alpha): #ff0000
Red HEXA (with alpha): #ff0000ff
Red CSS RGB: rgb(255, 0, 0)
Red is opaque: true

Test 2: HSLA Conversion and Caching
Input: Red RGB(255,0,0)
Converted HSLA: H=0.000, S=1.000, L=0.500, A=1.000
CSS Format: hsla(0.0, 100.0, 50.0, 1.000)
Caching test: First call=(0.000,1.000,0.500,1.000)
              Second call=(0.000,1.000,0.500,1.000)
Result: ✓ - Values identical (cached)

Test 3: HSLA Color Creation
Input: HSLA(240°, 100%, 50%, 70%)
Creating: NewHSLA(0.667, 1.000, 0.500, 0.700)
Result RGB: R=0, G=0, B=255, A=178
Hex format: #0000ffb2
CSS format: rgba(0, 0, 255, 0.698)
Alpha value: 0.698 (expected 0.7±0.01)
Conversion: ✓ - Alpha preserved correctly

Test 4: Alpha Manipulation
Original: RGB(0,255,0), Alpha=1.000
Operation: WithAlpha(0.5)
Result: rgba(0, 255, 0, 0.498), Alpha=0.498
Original unchanged: Alpha=1.000 ✓

Transparency checks:
  RGBA(255,0,0,0.0): IsTransparent=true, IsOpaque=false
  RGBA(255,0,0,1.0): IsTransparent=false, IsOpaque=true

Test 5: CSS Format Methods
Purple HEX (no alpha): #800080
Purple HEXA (with alpha): #800080bf
Purple CSS RGB: rgb(128, 0, 128)
Purple CSS RGBA: rgba(128, 0, 128, 0.749)
Purple CSS HSL: hsl(300.0, 100.0, 25.1)
Purple CSS HSLA: hsla(300.0, 100.0, 25.1, 0.749)

Test 6: Component Access Methods
Purple RGB components: R=128, G=0, B=128
Purple RGBA components: R=128, G=0, B=128, A=191
Purple HSL components: H=0.833, S=1.000, L=0.251
Purple HSLA components: H=0.833, S=1.000, L=0.251, A=0.749

Test 7: Thread Safety Test
Goroutine 9: HSLA=0.108,1.000,0.500,0.800
Goroutine 0: HSLA=0.108,1.000,0.500,0.800
Goroutine 1: HSLA=0.108,1.000,0.500,0.800
Goroutine 2: HSLA=0.108,1.000,0.500,0.800
Goroutine 3: HSLA=0.108,1.000,0.500,0.800
Goroutine 4: HSLA=0.108,1.000,0.500,0.800
Goroutine 5: HSLA=0.108,1.000,0.500,0.800
Goroutine 6: HSLA=0.108,1.000,0.500,0.800
Goroutine 7: HSLA=0.108,1.000,0.500,0.800
Goroutine 8: HSLA=0.108,1.000,0.500,0.800

Test 8: Color Conversion Accuracy with Alpha

White Opaque test:
  Input: RGBA(255,255,255,1.00)
  → HSLA: (0.000,0.000,1.000,1.000)
  → Back to RGB: (255,255,255,1.000)
  Differences: R=0, G=0, B=0, A=0 (tolerance ≤1)
  Result: ✓ - Round-trip accurate

Black 50% test:
  Input: RGBA(0,0,0,0.50)
  → HSLA: (0.000,0.000,0.000,0.498)
  → Back to RGB: (0,0,0,0.494)
  Differences: R=0, G=0, B=0, A=1 (tolerance ≤1)
  Result: ✓ - Round-trip accurate

Red 75% test:
  Input: RGBA(255,0,0,0.75)
  → HSLA: (0.000,1.000,0.500,0.749)
  → Back to RGB: (255,0,0,0.745)
  Differences: R=0, G=0, B=0, A=1 (tolerance ≤1)
  Result: ✓ - Round-trip accurate

Green Transparent test:
  Input: RGBA(0,255,0,0.00)
  → HSLA: (0.333,1.000,0.500,0.000)
  → Back to RGB: (0,255,0,0.000)
  Differences: R=0, G=0, B=0, A=0 (tolerance ≤1)
  Result: ✓ - Round-trip accurate

Test 9: Alpha Edge Cases
Alpha clamping test:
  Input alpha=1.5 → Result: 1.000 (expected 1.0) ✓
  Input alpha=-0.3 → Result: 0.000 (expected 0.0) ✓

Alpha variation test:
  Base color: RGB(100,150,200)
  WithAlpha(0.00): rgba(100, 150, 200, 0.000), actual=0.000 ✓
  WithAlpha(0.25): rgba(100, 150, 200, 0.247), actual=0.247 ✓
  WithAlpha(0.50): rgba(100, 150, 200, 0.498), actual=0.498 ✓
  WithAlpha(0.75): rgba(100, 150, 200, 0.749), actual=0.749 ✓
  WithAlpha(1.00): rgba(100, 150, 200, 1.000), actual=1.000 ✓

Test 10: Alpha Conversion Consistency
Alpha consistency check:
  Expected α=0.00, got α=0.000 ✓
  Expected α=0.25, got α=0.247 ✓
  Expected α=0.50, got α=0.498 ✓
  Expected α=0.75, got α=0.749 ✓
  Expected α=1.00, got α=1.000 ✓

Test 11: Immutability Test
Original color: rgba(255, 128, 64, 1.000) (α=1.000)
After WithAlpha(0.5):
  Original: rgba(255, 128, 64, 1.000) (α=1.000)
  Modified: rgba(255, 128, 64, 0.498) (α=0.498)
Result: ✓ - Original unchanged (immutable)

Test 12: HEXA Parsing
  Standard HEXA with #: #ff8000c0
    Expected: R=255, G=128, B=0, A=192
    Got:      R=255, G=128, B=0, A=192
    Result:   ✓ RGBA values match
    Round-trip: #ff8000c0 → #ff8000c0 ✓

  Standard HEXA without #: ff8000c0
    Expected: R=255, G=128, B=0, A=192
    Got:      R=255, G=128, B=0, A=192
    Result:   ✓ RGBA values match
    Round-trip: #ff8000c0 → #ff8000c0 ✓

  Full opacity: #00ff00ff
    Expected: R=0, G=255, B=0, A=255
    Got:      R=0, G=255, B=0, A=255
    Result:   ✓ RGBA values match
    Round-trip: #00ff00ff → #00ff00ff ✓

  Full transparency: #0000ff00
    Expected: R=0, G=0, B=255, A=0
    Got:      R=0, G=0, B=255, A=0
    Result:   ✓ RGBA values match
    Round-trip: #0000ff00 → #0000ff00 ✓

  Too short: #ff800
    Expected error: ✓ ✓ Got expected error
    Error: invalid HEXA format: expected 8 characters, got 5

  Invalid characters: #gghhiijj
    Expected error: ✓ ✓ Got expected error
    Error: invalid red component in HEXA: strconv.ParseUint: parsing "gg": invalid syntax

Test 13: HEX Parsing (6-digit)
  Red: #ff0000
    Expected: R=255, G=0, B=0, A=255
    Got:      R=255, G=0, B=0, A=255
    Result:   ✓ RGBA values match

  Green without #: 00ff00
    Expected: R=0, G=255, B=0, A=255
    Got:      R=0, G=255, B=0, A=255
    Result:   ✓ RGBA values match

  Too short: #ff00
    Expected error: ✓ ✓ Got expected error

  Too long: #ff0000ff
    Expected error: ✓ ✓ Got expected error

  Invalid: #gghhii
    Expected error: ✓ ✓ Got expected error

=== Color Type Test Complete ===
```

## Validation Criteria
✅ **All tests passing**: Every test shows ✓ checkmark indicating success  
✅ **Round-trip accuracy**: RGB↔HSLA conversions maintain ≤1 unit precision  
✅ **Thread safety**: Identical HSLA values across all concurrent goroutines  
✅ **Immutability**: Original colors unchanged after operations  
✅ **Format consistency**: All CSS and hex output formats properly formatted  
✅ **HEXA parsing accuracy**: 8-digit hex strings correctly parsed with alpha channel  
✅ **HEX parsing accuracy**: 6-digit hex strings correctly parsed (full opacity)  
✅ **Error handling**: Invalid hex formats properly rejected with clear error messages  
✅ **Round-trip hex conversion**: Parsed colors regenerate identical hex strings