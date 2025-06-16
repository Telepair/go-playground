// Package main implements Conway's Game of Life simulation with various patterns and configurations.
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Logger for configuration warnings and errors
var configLogger *log.Logger

func init() {
	// Initialize logger to write to stderr to avoid UI pollution
	configLogger = log.New(os.Stderr, "[CONFIG] ", log.LstdFlags|log.Lshortfile)
}

// BoundaryType represents the boundary type of the Game of Life
type BoundaryType int

// BoundaryType constants
const (
	BoundaryPeriodic BoundaryType = iota // Periodic boundary (wrapping, default)
	BoundaryFixed                        // Fixed boundary (dead cells outside)
)

// String returns the string representation of boundary type
func (bt BoundaryType) String() string {
	switch bt {
	case BoundaryPeriodic:
		return "periodic"
	case BoundaryFixed:
		return "fixed"
	default:
		return "periodic"
	}
}

// ChineseString returns the Chinese string representation of boundary type
func (bt BoundaryType) ChineseString() string {
	switch bt {
	case BoundaryPeriodic:
		return "周期"
	case BoundaryFixed:
		return "固定"
	default:
		return "周期"
	}
}

// ParseBoundaryType parses string to BoundaryType
func ParseBoundaryType(s string) BoundaryType {
	switch strings.ToLower(s) {
	case "periodic":
		return BoundaryPeriodic
	case "fixed":
		return BoundaryFixed
	default:
		return BoundaryPeriodic
	}
}

// Application constants
const (
	// Grid and display constants
	DefaultWindowRows = 24 // Default window rows
	DefaultWindowCols = 90 // Default window columns
	MinWindowRows     = 10 // Minimum window rows
	MinWindowCols     = 20 // Minimum window columns

	DefaultLanguage    = "en"                  // Default language
	DefaultRefreshRate = 50 * time.Millisecond // Default refresh rate in milliseconds
	MinRefreshRate     = 10 * time.Millisecond // Minimum refresh rate in milliseconds
	DefaultPattern     = PatternRandom         // Default pattern
	DefaultBoundary    = BoundaryPeriodic      // Default boundary type

	// Colors
	DefaultAliveColor = "#00FF00" // Default alive cell color (green)
	DefaultDeadColor  = "#000000" // Default dead cell color (black)
	HeaderBgColor     = "#874BFD" // Header background color
	HeaderFgColor     = "#FFFFFF" // Header foreground color

	// Characters
	DefaultAliveChar   = "█" // Default alive cell character
	DefaultDeadChar    = " " // Default dead cell character
	ProgressFilledChar = "█" // Progress bar filled character
	ProgressEmptyChar  = "░" // Progress bar empty character
)

// Pattern represents different starting patterns for Conway's Game of Life
type Pattern int

// Pattern constants
const (
	PatternRandom Pattern = iota
	PatternGlider
	PatternGliderGun
	PatternOscillator
	PatternPulsar
	PatternPentomino
)

// String returns the string representation of pattern type
func (p Pattern) String() string {
	switch p {
	case PatternRandom:
		return "random"
	case PatternGlider:
		return "glider"
	case PatternGliderGun:
		return "glider-gun"
	case PatternOscillator:
		return "oscillator"
	case PatternPulsar:
		return "pulsar"
	case PatternPentomino:
		return "pentomino"
	default:
		return "random"
	}
}

// ChineseString returns the Chinese string representation of pattern type
func (p Pattern) ChineseString() string {
	switch p {
	case PatternRandom:
		return "随机"
	case PatternGlider:
		return "滑翔机"
	case PatternGliderGun:
		return "滑翔机枪"
	case PatternOscillator:
		return "振荡器"
	case PatternPulsar:
		return "脉冲星"
	case PatternPentomino:
		return "五格骨牌"
	default:
		return "随机"
	}
}

// ParsePattern parses string to Pattern
func ParsePattern(s string) Pattern {
	switch strings.ToLower(s) {
	case "random":
		return PatternRandom
	case "glider":
		return PatternGlider
	case "glider-gun":
		return PatternGliderGun
	case "oscillator":
		return PatternOscillator
	case "pulsar":
		return PatternPulsar
	case "pentomino":
		return PatternPentomino
	default:
		return PatternRandom
	}
}

// Config holds all application configuration
type Config struct {
	Rows       int
	Cols       int
	AliveColor string
	DeadColor  string
	AliveChar  string
	DeadChar   string
	Language   Language
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		Rows:       DefaultWindowRows,
		Cols:       DefaultWindowCols,
		AliveColor: DefaultAliveColor,
		DeadColor:  DefaultDeadColor,
		AliveChar:  DefaultAliveChar,
		DeadChar:   DefaultDeadChar,
		Language:   English,
	}
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

// SetRows sets the number of rows with validation
func (c *Config) SetRows(rows int) error {
	if rows < MinWindowRows {
		err := fmt.Errorf("invalid number of rows %d, must be at least %d", rows, MinWindowRows)
		configLogger.Printf("Warning: %v, using default number of rows %d", err, DefaultWindowRows)
		c.Rows = DefaultWindowRows
		return err
	}
	c.Rows = rows
	return nil
}

// SetCols sets the number of columns with validation
func (c *Config) SetCols(cols int) error {
	if cols < MinWindowCols {
		err := fmt.Errorf("invalid number of columns %d, must be at least %d", cols, MinWindowCols)
		configLogger.Printf("Warning: %v, using default number of columns %d", err, DefaultWindowCols)
		c.Cols = DefaultWindowCols
		return err
	}
	c.Cols = cols
	return nil
}

// ValidateColors validates color format
func (c *Config) ValidateColors() error {
	// Simple hex color validation
	if !isValidHexColor(c.AliveColor) {
		err := fmt.Errorf("invalid alive color format: %s", c.AliveColor)
		configLogger.Printf("Warning: %v, using default", err)
		c.AliveColor = DefaultAliveColor
		return err
	}
	if !isValidHexColor(c.DeadColor) {
		err := fmt.Errorf("invalid dead color format: %s", c.DeadColor)
		configLogger.Printf("Warning: %v, using default", err)
		c.DeadColor = DefaultDeadColor
		return err
	}
	return nil
}

// isValidHexColor checks if a string is a valid hex color
func isValidHexColor(color string) bool {
	if len(color) != 7 || color[0] != '#' {
		return false
	}
	for _, c := range color[1:] {
		if (c < '0' || c > '9') && (c < 'A' || c > 'F') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}
