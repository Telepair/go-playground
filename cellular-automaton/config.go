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

// BoundaryType represents the boundary type of the cellular automaton
type BoundaryType int

// BoundaryType constants
const (
	BoundaryPeriodic BoundaryType = iota // Periodic boundary (default)
	BoundaryFixed                        // Fixed boundary (0 values)
	BoundaryReflect                      // Reflective boundary (mirror)
)

// String returns the string representation of boundary type
func (bt BoundaryType) String() string {
	switch bt {
	case BoundaryPeriodic:
		return "periodic"
	case BoundaryFixed:
		return "fixed"
	case BoundaryReflect:
		return "reflect"
	default:
		return "periodic"
	}
}

// ParseBoundaryType parses string to BoundaryType
func ParseBoundaryType(s string) BoundaryType {
	switch strings.ToLower(s) {
	case "periodic":
		return BoundaryPeriodic
	case "fixed":
		return BoundaryFixed
	case "reflect", "reflective":
		return BoundaryReflect
	default:
		return BoundaryPeriodic
	}
}

// Application constants
const (
	// Grid and display constants
	DefaultWindowRows = 35 // Default window rows
	DefaultWindowCols = 60 // Default window columns
	MinWindowRows     = 10 // Minimum window rows
	MinWindowCols     = 20 // Minimum window columns

	// Rule validation
	DefaultRule = 30  // Default cellular automaton rule
	MinRule     = 0   // Minimum rule number
	MaxRule     = 255 // Maximum rule number

	// Cell size validation
	DefaultCellSize = 1 // Default cell size
	MinCellSize     = 1 // Minimum cell size
	MaxCellSize     = 3 // Maximum cell size

	// Timing constants
	DefaultRefreshRate = 200 * time.Millisecond // Default refresh rate in milliseconds
	MinRefreshRate     = 1 * time.Millisecond   // Minimum refresh rate in milliseconds

	// Default values
	DefaultSteps    = 30               // Default number of steps
	DefaultLanguage = "en"             // Default language
	DefaultBoundary = BoundaryPeriodic // Default boundary type

	// Colors
	DefaultAliveColor = "#FFFFFF" // Default alive cell color
	DefaultDeadColor  = "#000000" // Default dead cell color
	HeaderBgColor     = "#874BFD" // Header background color
	HeaderFgColor     = "#FFFFFF" // Header foreground color

	// Characters
	DefaultAliveChar   = "█" // Default alive cell character
	DefaultDeadChar    = " " // Default dead cell character
	ProgressFilledChar = "█" // Progress bar filled character
	ProgressEmptyChar  = "░" // Progress bar empty character
)

// Config holds all application configuration
type Config struct {
	Rule         int
	Rows         int
	Cols         int
	Steps        int
	CellSize     int
	AliveColor   string
	DeadColor    string
	AliveChar    string
	DeadChar     string
	RefreshRate  time.Duration
	Language     Language
	Boundary     BoundaryType
	InfiniteMode bool
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		Rule:        DefaultRule,
		Rows:        DefaultWindowRows,
		Cols:        DefaultWindowCols,
		Steps:       DefaultSteps,
		CellSize:    DefaultCellSize,
		AliveColor:  DefaultAliveColor,
		DeadColor:   DefaultDeadColor,
		AliveChar:   DefaultAliveChar,
		DeadChar:    DefaultDeadChar,
		RefreshRate: DefaultRefreshRate,
		Language:    English,
		Boundary:    DefaultBoundary,
	}
}

// SetRule sets the cellular automaton rule with validation
func (c *Config) SetRule(rule int) error {
	if rule < MinRule || rule > MaxRule {
		err := fmt.Errorf("invalid rule %d, must be between %d and %d", rule, MinRule, MaxRule)
		configLogger.Printf("Warning: %v, using default rule %d", err, DefaultRule)
		c.Rule = DefaultRule
		return err
	}
	c.Rule = rule
	return nil
}

// SetSteps sets the number of steps and sets the infinite mode based on the number of steps
func (c *Config) SetSteps(steps int) {
	c.Steps = steps
	c.InfiniteMode = c.Steps == 0
}

// SetCellSize sets the cell size with validation
func (c *Config) SetCellSize(cellSize int) error {
	if cellSize < MinCellSize || cellSize > MaxCellSize {
		err := fmt.Errorf("invalid cell size %d, must be between %d and %d", cellSize, MinCellSize, MaxCellSize)
		configLogger.Printf("Warning: %v, using default cell size %d", err, DefaultCellSize)
		c.CellSize = DefaultCellSize
		return err
	}
	c.CellSize = cellSize
	return nil
}

// SetRefreshRate sets the refresh rate with validation
func (c *Config) SetRefreshRate(duration time.Duration) error {
	if duration < MinRefreshRate {
		err := fmt.Errorf("invalid refresh rate %s, must be at least %s", duration, MinRefreshRate)
		configLogger.Printf("Warning: %v, using default refresh rate %s", err, DefaultRefreshRate)
		c.RefreshRate = DefaultRefreshRate
		return err
	}
	c.RefreshRate = duration
	return nil
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

// SetBoundary sets the boundary type
func (c *Config) SetBoundary(boundary string) {
	c.Boundary = ParseBoundaryType(boundary)
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
