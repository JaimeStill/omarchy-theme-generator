package formats_test

import (
	"testing"

	"github.com/JaimeStill/omarchy-theme-generator/pkg/formats"
)

// TestNewXYZ tests the NewXYZ constructor function
func TestNewXYZ(t *testing.T) {
	testCases := []struct {
		name string
		x, y, z float64
		expected formats.XYZ
		description string
	}{
		{
			name: "Standard values",
			x: 95.047, y: 100.000, z: 108.883,
			expected: formats.XYZ{X: 95.047, Y: 100.000, Z: 108.883},
			description: "Constructor should create XYZ with specified values",
		},
		{
			name: "Zero values",
			x: 0.0, y: 0.0, z: 0.0,
			expected: formats.XYZ{X: 0.0, Y: 0.0, Z: 0.0},
			description: "Constructor should handle zero values",
		},
		{
			name: "Negative values",
			x: -10.5, y: -20.3, z: -15.7,
			expected: formats.XYZ{X: -10.5, Y: -20.3, Z: -15.7},
			description: "Constructor should handle negative values",
		},
		{
			name: "Large values",
			x: 1000.0, y: 2000.0, z: 3000.0,
			expected: formats.XYZ{X: 1000.0, Y: 2000.0, Z: 3000.0},
			description: "Constructor should handle large values",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			xyz := formats.NewXYZ(tc.x, tc.y, tc.z)

			// Log diagnostic information
			t.Logf("Test: %s", tc.name)
			t.Logf("Input values: X=%.3f, Y=%.3f, Z=%.3f", tc.x, tc.y, tc.z)
			t.Logf("Created XYZ: %+v", xyz)
			t.Logf("Expected XYZ: %+v", tc.expected)
			t.Logf("Description: %s", tc.description)

			// Verify the created XYZ matches expected values
			if xyz.X != tc.expected.X {
				t.Errorf("X component mismatch: expected %.3f, got %.3f", tc.expected.X, xyz.X)
			}
			if xyz.Y != tc.expected.Y {
				t.Errorf("Y component mismatch: expected %.3f, got %.3f", tc.expected.Y, xyz.Y)
			}
			if xyz.Z != tc.expected.Z {
				t.Errorf("Z component mismatch: expected %.3f, got %.3f", tc.expected.Z, xyz.Z)
			}

			t.Logf("✅ XYZ constructor working correctly")
		})
	}
}

// TestD65Illuminant tests the D65 illuminant constant
func TestD65Illuminant(t *testing.T) {
	t.Logf("Testing D65Illuminant constant")
	t.Logf("D65Illuminant: X=%.3f, Y=%.3f, Z=%.3f", 
		formats.D65Illuminant.X, formats.D65Illuminant.Y, formats.D65Illuminant.Z)

	// Verify expected D65 values (CIE standard)
	expectedX := 95.047
	expectedY := 100.000
	expectedZ := 108.883

	if formats.D65Illuminant.X != expectedX {
		t.Errorf("D65Illuminant.X mismatch: expected %.3f, got %.3f", expectedX, formats.D65Illuminant.X)
	}
	if formats.D65Illuminant.Y != expectedY {
		t.Errorf("D65Illuminant.Y mismatch: expected %.3f, got %.3f", expectedY, formats.D65Illuminant.Y)
	}
	if formats.D65Illuminant.Z != expectedZ {
		t.Errorf("D65Illuminant.Z mismatch: expected %.3f, got %.3f", expectedZ, formats.D65Illuminant.Z)
	}

	t.Logf("✅ D65Illuminant constant has correct CIE standard values")
}

// TestD50Illuminant tests the D50 illuminant constant
func TestD50Illuminant(t *testing.T) {
	t.Logf("Testing D50Illuminant constant")
	t.Logf("D50Illuminant: X=%.3f, Y=%.3f, Z=%.3f", 
		formats.D50Illuminant.X, formats.D50Illuminant.Y, formats.D50Illuminant.Z)

	// Verify expected D50 values (CIE standard)
	expectedX := 96.422
	expectedY := 100.00
	expectedZ := 82.521

	if formats.D50Illuminant.X != expectedX {
		t.Errorf("D50Illuminant.X mismatch: expected %.3f, got %.3f", expectedX, formats.D50Illuminant.X)
	}
	if formats.D50Illuminant.Y != expectedY {
		t.Errorf("D50Illuminant.Y mismatch: expected %.3f, got %.3f", expectedY, formats.D50Illuminant.Y)
	}
	if formats.D50Illuminant.Z != expectedZ {
		t.Errorf("D50Illuminant.Z mismatch: expected %.3f, got %.3f", expectedZ, formats.D50Illuminant.Z)
	}

	t.Logf("✅ D50Illuminant constant has correct CIE standard values")
}

// TestXYZStruct tests that the XYZ struct can be created and accessed directly
func TestXYZStruct(t *testing.T) {
	t.Logf("Testing direct XYZ struct creation and field access")
	
	xyz := formats.XYZ{X: 50.0, Y: 75.5, Z: 25.3}
	
	t.Logf("Created XYZ struct: X=%.1f, Y=%.1f, Z=%.1f", xyz.X, xyz.Y, xyz.Z)
	
	// Verify field access
	if xyz.X != 50.0 {
		t.Errorf("X field access failed: expected 50.0, got %.1f", xyz.X)
	}
	if xyz.Y != 75.5 {
		t.Errorf("Y field access failed: expected 75.5, got %.1f", xyz.Y)
	}
	if xyz.Z != 25.3 {
		t.Errorf("Z field access failed: expected 25.3, got %.1f", xyz.Z)
	}

	t.Logf("✅ XYZ struct fields accessible and working correctly")
}