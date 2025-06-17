// Package main implements a Mandelbrot and Julia set fractal visualization
// using a terminal user interface with interactive navigation capabilities.
package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Language represents the supported languages
type Language int

// Language constants
const (
	English Language = iota
	Chinese
)

// ToString returns the string representation of language
func (l Language) ToString(language Language) string {
	switch l {
	case English:
		if language == Chinese {
			return "英文"
		}
		return "en"
	case Chinese:
		if language == Chinese {
			return "中文"
		}
		return "cn"
	}
	if language == Chinese {
		return "英文"
	}
	return "en"
}

// ColorScheme represents different color schemes for rendering
type ColorScheme int

// ColorScheme constants
const (
	ColorSchemeClassic   ColorScheme = iota // Black and white classic
	ColorSchemeHot                          // Hot colors (red, orange, yellow)
	ColorSchemeCool                         // Cool colors (blue, cyan, purple)
	ColorSchemeRainbow                      // Rainbow spectrum
	ColorSchemeGrayscale                    // Grayscale gradient
)

// ToString returns the string representation of color scheme
func (cs ColorScheme) ToString(language Language) string {
	switch cs {
	case ColorSchemeClassic:
		if language == Chinese {
			return "经典"
		}
		return "Classic"
	case ColorSchemeHot:
		if language == Chinese {
			return "热色"
		}
		return "Hot"
	case ColorSchemeCool:
		if language == Chinese {
			return "冷色"
		}
		return "Cool"
	case ColorSchemeRainbow:
		if language == Chinese {
			return "彩虹"
		}
		return "Rainbow"
	case ColorSchemeGrayscale:
		if language == Chinese {
			return "灰度"
		}
		return "Grayscale"
	default:
		if language == Chinese {
			return "经典"
		}
		return "Classic"
	}
}

// Application constants
const (
	// Grid and display constants
	DefaultWindowRows = 30 // Default window rows
	DefaultWindowCols = 80 // Default window columns
	MinWindowRows     = 10 // Minimum window rows
	MinWindowCols     = 20 // Minimum window columns

	// Mandelbrot set constants
	DefaultMaxIterations = 50              // Default maximum iterations
	MinMaxIterations     = 10              // Minimum iterations
	MaxMaxIterations     = 1000            // Maximum iterations
	DefaultZoom          = 1.0             // Default zoom level
	DefaultCenterX       = -0.5            // Default center X coordinate
	DefaultCenterY       = 0.0             // Default center Y coordinate
	DefaultJuliaC        = "-0.7+0.27015i" // Default Julia set parameter

	// Timing constants
	DefaultRefreshRate = 100 * time.Millisecond // Default refresh rate
	MinRefreshRate     = 10 * time.Millisecond  // Minimum refresh rate

	// Default values
	DefaultLanguage    = English            // Default language
	DefaultColorScheme = ColorSchemeClassic // Default color scheme
)

// DefaultConfig is the default configuration
var DefaultConfig = Config{
	Rows:        DefaultWindowRows,
	Cols:        DefaultWindowCols,
	MaxIter:     DefaultMaxIterations,
	Zoom:        DefaultZoom,
	CenterX:     DefaultCenterX,
	CenterY:     DefaultCenterY,
	ColorScheme: DefaultColorScheme,
	Julia:       false,
	JuliaC:      DefaultJuliaC,
	Language:    DefaultLanguage,
}

// Config holds all application configuration
type Config struct {
	Rows        int
	Cols        int
	MaxIter     int
	Zoom        float64
	CenterX     float64
	CenterY     float64
	ColorScheme ColorScheme
	Julia       bool
	JuliaC      string
	Language    Language
}

// SetLanguage sets the language
func (c *Config) SetLanguage(lang string) {
	langLower := strings.ToLower(lang)
	if langLower == "cn" || langLower == "zh" {
		c.Language = Chinese
	} else {
		c.Language = English
	}
}

// Check validates the configuration
func (c *Config) Check() {
	if c.Rows < MinWindowRows {
		fmt.Printf("invalid number of rows %d, must be at least %d, using default number of rows %d\n", c.Rows, MinWindowRows, DefaultWindowRows)
		c.Rows = DefaultWindowRows
	}
	if c.Cols < MinWindowCols {
		fmt.Printf("invalid number of columns %d, must be at least %d, using default number of columns %d\n", c.Cols, MinWindowCols, DefaultWindowCols)
		c.Cols = DefaultWindowCols
	}
	if c.MaxIter < MinMaxIterations || c.MaxIter > MaxMaxIterations {
		fmt.Printf("invalid max iterations %d, must be between %d and %d, using default %d\n", c.MaxIter, MinMaxIterations, MaxMaxIterations, DefaultMaxIterations)
		c.MaxIter = DefaultMaxIterations
	}
	if c.Zoom <= 0 {
		fmt.Printf("invalid zoom level %f, must be positive, using default %f\n", c.Zoom, DefaultZoom)
		c.Zoom = DefaultZoom
	}
	if c.ColorScheme < ColorSchemeClassic || c.ColorScheme > ColorSchemeGrayscale {
		fmt.Printf("invalid color scheme %d, must be between 0 and 4, using default %d\n", c.ColorScheme, DefaultColorScheme)
		c.ColorScheme = DefaultColorScheme
	}
}

// ParseComplexNumber parses a complex number string in the format "a+bi" or "a-bi"
func ParseComplexNumber(s string) (complex128, error) {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "i", "")

	// Handle different formats
	if strings.Contains(s, "+") {
		parts := strings.Split(s, "+")
		if len(parts) == 2 {
			realPart, err1 := strconv.ParseFloat(parts[0], 64)
			imagPart, err2 := strconv.ParseFloat(parts[1], 64)
			if err1 == nil && err2 == nil {
				return complex(realPart, imagPart), nil
			}
		}
	} else if strings.Count(s, "-") == 2 && strings.HasPrefix(s, "-") {
		// Handle negative real part like "-0.7-0.27015"
		s = s[1:] // Remove first minus
		parts := strings.Split(s, "-")
		if len(parts) == 2 {
			realPart, err1 := strconv.ParseFloat("-"+parts[0], 64)
			imagPart, err2 := strconv.ParseFloat("-"+parts[1], 64)
			if err1 == nil && err2 == nil {
				return complex(realPart, imagPart), nil
			}
		}
	} else if strings.Contains(s, "-") && !strings.HasPrefix(s, "-") {
		parts := strings.Split(s, "-")
		if len(parts) == 2 {
			realPart, err1 := strconv.ParseFloat(parts[0], 64)
			imagPart, err2 := strconv.ParseFloat("-"+parts[1], 64)
			if err1 == nil && err2 == nil {
				return complex(realPart, imagPart), nil
			}
		}
	}

	return 0, fmt.Errorf("invalid complex number format: %s", s)
}
