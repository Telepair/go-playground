package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
)

// Boundary type constants
type BoundaryType uint

const (
	BoundaryPeriodic BoundaryType = iota // Periodic boundary (default)
	BoundaryFixed                        // Fixed boundary (0 values)
	BoundaryReflect                      // Reflective boundary
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
	DefaultWindowRows = 30 // Default window rows
	DefaultWindowCols = 60 // Default window columns
	MinWindowRows     = 10 // Minimum window rows
	MinWindowCols     = 20 // Minimum window columns
	UISpaceReserved   = 8  // Space reserved for UI elements

	// Rule validation
	MinRule = 0   // Minimum rule number
	MaxRule = 255 // Maximum rule number

	// Cell size validation
	MinCellSize = 1 // Minimum cell size
	MaxCellSize = 3 // Maximum cell size

	// Timing constants
	MinRefreshRate     = 1  // Minimum refresh rate in milliseconds
	FiniteModeTickRate = 50 // Tick rate for finite mode in milliseconds

	// Default values
	DefaultRule        = 30               // Default cellular automaton rule
	DefaultSteps       = 30               // Default number of steps
	DefaultCellSize    = 2                // Default cell size
	DefaultRefreshRate = 0.1              // Default refresh rate in seconds
	DefaultLanguage    = "en"             // Default language
	DefaultBoundary    = BoundaryPeriodic // Default boundary type

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
	Window       Window
	Rule         uint
	Steps        uint
	CellSize     uint
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
		Window: Window{
			rows: DefaultWindowRows,
			cols: DefaultWindowCols,
		},
		Rule:        DefaultRule,
		Steps:       DefaultSteps,
		CellSize:    DefaultCellSize,
		AliveColor:  DefaultAliveColor,
		DeadColor:   DefaultDeadColor,
		AliveChar:   DefaultAliveChar,
		DeadChar:    DefaultDeadChar,
		RefreshRate: time.Duration(DefaultRefreshRate * float64(time.Second)),
		Language:    English,
		Boundary:    DefaultBoundary,
	}
}

// SetRule sets the cellular automaton rule
func (c *Config) SetRule(rule uint) {
	c.Rule = rule
	if c.Rule < MinRule || c.Rule > MaxRule {
		fmt.Printf("Invalid rule: %d, using default rule: %d\n", c.Rule, DefaultRule)
		c.Rule = DefaultRule
	}
}

// SetSteps sets the number of steps and sets the infinite mode based on the number of steps
func (c *Config) SetSteps(steps uint) {
	c.Steps = steps
	c.InfiniteMode = c.Steps == 0
}

// ValidateCellSize validates the cell size
func (c *Config) SetCellSize(cellSize uint) {
	c.CellSize = cellSize
	if c.CellSize < MinCellSize || c.CellSize > MaxCellSize {
		fmt.Printf("Invalid cell size: %d, using default cell size: %d\n", c.CellSize, DefaultCellSize)
		c.CellSize = DefaultCellSize
	}
}

// SetRefreshRate sets the refresh rate
func (c *Config) SetRefreshRate(seconds float64) {
	duration := time.Duration(seconds * float64(time.Second))
	if duration < time.Millisecond*MinRefreshRate {
		duration = time.Millisecond * MinRefreshRate
	}
	c.RefreshRate = duration
}

// SetLanguage sets the language
func (c *Config) SetLanguage(lang string) {
	if lang == "cn" || lang == "zh" {
		c.Language = Chinese
	} else {
		c.Language = English
	}
}

// SetBoundary sets the boundary type
func (c *Config) SetBoundary(boundary string) {
	c.Boundary = ParseBoundaryType(boundary)
}

// SetSize sets the width and height of the grid
func (c *Config) SetWindowSize(sizeStr string) {
	sizeStr = strings.ToLower(strings.TrimSpace(sizeStr))
	// Parse manual size specification
	parts := strings.Split(sizeStr, "x")
	if len(parts) != 2 {
		fmt.Printf("Invalid size format: %s, using default size: %dx%d\n", sizeStr, DefaultWindowCols, DefaultWindowRows)
		c.Window.cols = DefaultWindowCols
		c.Window.rows = DefaultWindowRows
		return
	}

	c.Window.cols = cast.ToUint(parts[0])
	c.Window.rows = cast.ToUint(parts[1])

	// Validate minimum sizes
	if c.Window.cols < MinWindowCols {
		c.Window.cols = MinWindowCols
	}
	if c.Window.rows < MinWindowRows {
		c.Window.rows = MinWindowRows
	}
}
