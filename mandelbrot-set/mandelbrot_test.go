package main

import (
	"testing"
)

func TestNewMandelbrotSet(t *testing.T) {
	config := DefaultConfig
	mandelbrot := NewMandelbrotSet(config)

	if mandelbrot == nil {
		t.Fatal("NewMandelbrotSet returned nil")
	}

	if mandelbrot.GetZoom() != DefaultZoom {
		t.Errorf("Expected zoom %f, got %f", DefaultZoom, mandelbrot.GetZoom())
	}

	centerX, centerY := mandelbrot.GetCenter()
	if centerX != DefaultCenterX || centerY != DefaultCenterY {
		t.Errorf("Expected center (%f, %f), got (%f, %f)", DefaultCenterX, DefaultCenterY, centerX, centerY)
	}
}

func TestMandelbrotIterations(t *testing.T) {
	config := DefaultConfig
	mandelbrot := NewMandelbrotSet(config)

	// Test a point known to be in the set (should return max iterations)
	result := mandelbrot.mandelbrotIterations(complex(0, 0))
	if result != DefaultMaxIterations {
		t.Errorf("Point (0,0) should be in the set, expected %d iterations, got %d", DefaultMaxIterations, result)
	}

	// Test a point known to diverge quickly
	result = mandelbrot.mandelbrotIterations(complex(2, 2))
	if result >= DefaultMaxIterations {
		t.Errorf("Point (2,2) should diverge quickly, got %d iterations", result)
	}
}

func TestComplexNumberParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected complex128
		hasError bool
	}{
		{"0.5+0.5i", complex(0.5, 0.5), false},
		{"-0.7+0.27015i", complex(-0.7, 0.27015), false},
		{"0.285+0.01i", complex(0.285, 0.01), false},
		{"-0.5-0.5i", complex(-0.5, -0.5), false},
		// New test cases for edge cases
		{"", complex(0, 0), true},          // Empty string
		{"1.5", complex(1.5, 0), false},    // Pure real number
		{"-2.3", complex(-2.3, 0), false},  // Negative real number
		{"0+0i", complex(0, 0), false},     // Zero complex number
		{"+1.0+1.0i", complex(0, 0), true}, // Invalid leading plus
		{"1.0++1.0i", complex(0, 0), true}, // Double plus
		{"1.0+-1.0i", complex(0, 0), true}, // Mixed operators
		{"abc+def", complex(0, 0), true},   // Invalid characters
		{"1.0+i", complex(0, 0), true},     // Missing imaginary value
		{"invalid", 0, true},
	}

	for _, test := range tests {
		result, err := ParseComplexNumber(test.input)
		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for input '%s', but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("For input '%s', expected %v, got %v", test.input, test.expected, result)
			}
		}
	}
}

func TestColorSchemes(t *testing.T) {
	renderOptions := NewRenderOptions(ColorSchemeClassic)

	// Test color scheme functions don't panic
	for scheme := ColorSchemeClassic; scheme <= ColorSchemeGrayscale; scheme++ {
		renderOptions.colorScheme = scheme
		color := renderOptions.GetColorForIteration(10, 50)
		if color == "" {
			t.Errorf("Color scheme %d returned empty color", scheme)
		}
	}
}

func BenchmarkMandelbrotCalculation(b *testing.B) {
	config := Config{
		MaxIter:     50,
		Zoom:        1.0,
		CenterX:     -0.5,
		CenterY:     0.0,
		ColorScheme: ColorSchemeClassic,
		Julia:       false,
		JuliaC:      DefaultJuliaC,
		Language:    DefaultLanguage,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mandelbrot := NewMandelbrotSet(config)
		_ = mandelbrot.GetGrid()
	}
}

// Test numerical stability with extreme values
func TestMandelbrotNumericalStability(t *testing.T) {
	config := Config{
		MaxIter:     100,
		Zoom:        1.0,
		CenterX:     0.0,
		CenterY:     0.0,
		ColorScheme: ColorSchemeClassic,
		Julia:       false,
		JuliaC:      DefaultJuliaC,
		Language:    DefaultLanguage,
	}
	mandelbrot := NewMandelbrotSet(config)

	// Test with extremely large values that could cause overflow
	extremeValues := []complex128{
		complex(1e10, 1e10),   // Very large numbers
		complex(-1e10, -1e10), // Very large negative numbers
		complex(1e-10, 1e-10), // Very small numbers
	}

	for _, c := range extremeValues {
		iterations := mandelbrot.mandelbrotIterations(c)
		if iterations < 0 || iterations > config.MaxIter {
			t.Errorf("Invalid iteration count %d for complex number %v", iterations, c)
		}

		// Julia set test
		juliaIterations := mandelbrot.juliaIterations(c)
		if juliaIterations < 0 || juliaIterations > config.MaxIter {
			t.Errorf("Invalid Julia iteration count %d for complex number %v", juliaIterations, c)
		}
	}
}
