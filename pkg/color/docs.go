// Package color provides high-performance color representation and manipulation
// for theme generation. Colors use RGBA storage with lazy-cached HSLA conversion
// for optimal image processing performance while supporting color theory operations.
//
// The Color type is thread-safe and supports both programmatic alpha values (0.0-1.0)
// and user-friendly opacity percentages (0-100%). All CSS color formats are supported
// including hex, rgb(), rgba(), hsl(), and hsla().
//
// Key features:
//   - Zero-allocation RGBA storage for image extraction
//   - Thread-safe HSLA caching using sync.Once
//   - Bidirectional RGB ↔ HSL conversion following CSS Color Module Level 3
//   - Comprehensive output formatting for web and terminal applications
package color
