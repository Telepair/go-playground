package main

import (
	"math/cmplx"
)

// MandelbrotSet represents the Mandelbrot/Julia set calculator
type MandelbrotSet struct {
	width       int         // Grid width (columns)
	height      int         // Grid height (rows)
	maxIter     int         // Maximum iterations
	zoom        float64     // Zoom level
	centerX     float64     // Center X coordinate
	centerY     float64     // Center Y coordinate
	julia       bool        // Julia set mode
	juliaC      complex128  // Julia set parameter
	grid        [][]int     // Iteration count grid
	colorScheme ColorScheme // Color scheme for rendering
}

// NewMandelbrotSet creates a new Mandelbrot set instance
func NewMandelbrotSet(config Config) *MandelbrotSet {
	juliaC, err := ParseComplexNumber(config.JuliaC)
	if err != nil {
		// Use default Julia parameter if parsing fails
		juliaC = complex(-0.7, 0.27015)
	}

	m := &MandelbrotSet{
		width:       DefaultCols,
		height:      DefaultRows,
		maxIter:     config.MaxIter,
		zoom:        config.Zoom,
		centerX:     config.CenterX,
		centerY:     config.CenterY,
		julia:       config.Julia,
		juliaC:      juliaC,
		colorScheme: config.ColorScheme,
	}

	// Initialize grid
	m.grid = make([][]int, m.height)
	for i := range m.grid {
		m.grid[i] = make([]int, m.width)
	}

	// Calculate initial set
	m.Calculate()

	return m
}

// Calculate computes the Mandelbrot or Julia set
func (m *MandelbrotSet) Calculate() {
	// Calculate the viewing window based on zoom and center
	viewWidth := 4.0 / m.zoom
	viewHeight := (4.0 * float64(m.height) / float64(m.width)) / m.zoom

	minReal := m.centerX - viewWidth/2
	maxReal := m.centerX + viewWidth/2
	minImag := m.centerY - viewHeight/2
	maxImag := m.centerY + viewHeight/2

	// Calculate step size
	stepReal := (maxReal - minReal) / float64(m.width)
	stepImag := (maxImag - minImag) / float64(m.height)

	// Compute each point in the grid
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			// Convert screen coordinates to complex plane
			realPart := minReal + float64(x)*stepReal
			imagPart := minImag + float64(y)*stepImag
			c := complex(realPart, imagPart)

			// Calculate iterations for this point
			if m.julia {
				m.grid[y][x] = m.juliaIterations(c)
			} else {
				m.grid[y][x] = m.mandelbrotIterations(c)
			}
		}
	}
}

// mandelbrotIterations calculates the number of iterations for a point in the Mandelbrot set
func (m *MandelbrotSet) mandelbrotIterations(c complex128) int {
	z := complex(0, 0)

	for i := 0; i < m.maxIter; i++ {
		if cmplx.Abs(z) > 2.0 {
			return i
		}
		z = z*z + c
	}

	return m.maxIter
}

// juliaIterations calculates the number of iterations for a point in the Julia set
func (m *MandelbrotSet) juliaIterations(z complex128) int {
	for i := 0; i < m.maxIter; i++ {
		if cmplx.Abs(z) > 2.0 {
			return i
		}
		z = z*z + m.juliaC
	}

	return m.maxIter
}

// GetGrid returns the current iteration grid
func (m *MandelbrotSet) GetGrid() [][]int {
	return m.grid
}

// SetZoom sets the zoom level and recalculates
func (m *MandelbrotSet) SetZoom(zoom float64) {
	if zoom > 0 {
		m.zoom = zoom
		m.Calculate()
	}
}

// SetCenter sets the center coordinates and recalculates
func (m *MandelbrotSet) SetCenter(x, y float64) {
	m.centerX = x
	m.centerY = y
	m.Calculate()
}

// SetMaxIterations sets the maximum iterations and recalculates
func (m *MandelbrotSet) SetMaxIterations(maxIter int) {
	if maxIter > 0 {
		m.maxIter = maxIter
		m.Calculate()
	}
}

// SetColorScheme sets the color scheme
func (m *MandelbrotSet) SetColorScheme(scheme ColorScheme) {
	m.colorScheme = scheme
}

// ToggleMode toggles between Mandelbrot and Julia set modes
func (m *MandelbrotSet) ToggleMode() {
	m.julia = !m.julia
	m.Calculate()
}

// SetJuliaParameter sets the Julia set parameter and recalculates if in Julia mode
func (m *MandelbrotSet) SetJuliaParameter(c complex128) {
	m.juliaC = c
	if m.julia {
		m.Calculate()
	}
}

// ZoomIn zooms in by a factor at the current center
func (m *MandelbrotSet) ZoomIn(factor float64) {
	m.SetZoom(m.zoom * factor)
}

// ZoomOut zooms out by a factor at the current center
func (m *MandelbrotSet) ZoomOut(factor float64) {
	m.SetZoom(m.zoom / factor)
}

// Pan moves the view by the specified offsets (in screen coordinates)
func (m *MandelbrotSet) Pan(deltaX, deltaY int) {
	// Calculate the current view dimensions
	viewWidth := 4.0 / m.zoom
	viewHeight := (4.0 * float64(m.height) / float64(m.width)) / m.zoom

	// Convert screen offset to complex plane offset
	stepReal := viewWidth / float64(m.width)
	stepImag := viewHeight / float64(m.height)

	newCenterX := m.centerX + float64(deltaX)*stepReal
	newCenterY := m.centerY + float64(deltaY)*stepImag

	m.SetCenter(newCenterX, newCenterY)
}

// Reset resets to default parameters
func (m *MandelbrotSet) Reset(height, width int) {
	m.height = height
	m.width = width
	m.zoom = DefaultZoom
	m.centerX = DefaultCenterX
	m.centerY = DefaultCenterY
	m.maxIter = DefaultMaxIterations
	m.grid = make([][]int, m.height)
	for i := range m.grid {
		m.grid[i] = make([]int, m.width)
	}
	m.julia = false
	juliaC, _ := ParseComplexNumber(DefaultJuliaC)
	m.juliaC = juliaC
	m.Calculate()
}

// GetCurrentMode returns the current mode (Mandelbrot or Julia)
func (m *MandelbrotSet) GetCurrentMode() bool {
	return m.julia
}

// GetZoom returns the current zoom level
func (m *MandelbrotSet) GetZoom() float64 {
	return m.zoom
}

// GetCenter returns the current center coordinates
func (m *MandelbrotSet) GetCenter() (float64, float64) {
	return m.centerX, m.centerY
}

// GetMaxIterations returns the current maximum iterations
func (m *MandelbrotSet) GetMaxIterations() int {
	return m.maxIter
}

// GetJuliaParameter returns the current Julia set parameter
func (m *MandelbrotSet) GetJuliaParameter() complex128 {
	return m.juliaC
}

// GetColorScheme returns the current color scheme
func (m *MandelbrotSet) GetColorScheme() ColorScheme {
	return m.colorScheme
}

// GetInterestingPoints returns a list of interesting coordinates to explore
func (m *MandelbrotSet) GetInterestingPoints() []struct {
	Name string
	X, Y float64
	Zoom float64
} {
	return []struct {
		Name string
		X, Y float64
		Zoom float64
	}{
		{"Classic View", -0.5, 0.0, 1.0},
		{"Seahorse Valley", -0.75, 0.1, 50.0},
		{"Lightning", -1.775, 0.0, 100.0},
		{"Elephant Valley", 0.25, 0.0, 10.0},
		{"Spiral", -0.1592, -1.0317, 100.0},
		{"Mini Mandelbrot", -1.25066, 0.02012, 2000.0},
		{"Feather", -0.7463, 0.1102, 200.0},
		{"Dragon", -0.7269, 0.1889, 300.0},
	}
}
