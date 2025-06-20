// Package main implements Conway's Game of Life simulation with various patterns and configurations.
package main

import (
	"fmt"
	"strings"
	"time"
)

// BoundaryType represents the boundary type of the Game of Life
type BoundaryType int

// BoundaryType constants
const (
	BoundaryPeriodic BoundaryType = iota // Periodic boundary (wrapping, default)
	BoundaryFixed                        // Fixed boundary (dead cells outside)
)

// ToString returns the string representation of boundary type
func (bt BoundaryType) ToString(language Language) string {
	switch bt {
	case BoundaryPeriodic:
		if language == Chinese {
			return "周期"
		}
		return "Periodic"
	case BoundaryFixed:
		if language == Chinese {
			return "固定"
		}
		return "Fixed"
	default:
		if language == Chinese {
			return "周期"
		}
		return "Periodic"
	}
}

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

// ToString returns the string representation of pattern type
func (p Pattern) ToString(language Language) string {
	switch p {
	case PatternRandom:
		if language == Chinese {
			return "随机"
		}
		return "random"
	case PatternGlider:
		if language == Chinese {
			return "滑翔机"
		}
		return "glider"
	case PatternGliderGun:
		if language == Chinese {
			return "滑翔机枪"
		}
		return "glider-gun"
	case PatternOscillator:
		if language == Chinese {
			return "振荡器"
		}
		return "oscillator"
	case PatternPulsar:
		if language == Chinese {
			return "脉冲星"
		}
		return "pulsar"
	case PatternPentomino:
		if language == Chinese {
			return "五格骨牌"
		}
		return "pentomino"
	default:
		if language == Chinese {
			return "随机"
		}
		return "random"
	}
}

// Application constants
const (
	// Grid and display constants
	DefaultRows = 30 // Default window rows
	DefaultCols = 80 // Default window columns
	MinRows     = 10 // Minimum window rows
	MinCols     = 20 // Minimum window columns

	DefaultLanguage    = English               // Default language
	DefaultRefreshRate = 50 * time.Millisecond // Default refresh rate in milliseconds
	MinRefreshRate     = 10 * time.Millisecond // Minimum refresh rate in milliseconds
	DefaultPattern     = PatternRandom         // Default pattern
	DefaultBoundary    = BoundaryPeriodic      // Default boundary type

	// Colors
	DefaultAliveColor = "#00FF00" // Default alive cell color (green)
	DefaultDeadColor  = "#000000" // Default dead cell color (black)

	// Characters
	DefaultAliveChar = "█" // Default alive cell character
	DefaultDeadChar  = " " // Default dead cell character

	// Default values
	DefaultLogFile         = "debug.log"     // Default log file path
	DefaultProfileInterval = 5 * time.Second // Default profile information output interval
	DefaultProfilePort     = 6060            // Default profile server port
)

// DefaultConfig is the default configuration
var DefaultConfig = Config{
	AliveColor: DefaultAliveColor,
	DeadColor:  DefaultDeadColor,
	AliveChar:  DefaultAliveChar,
	DeadChar:   DefaultDeadChar,
	Language:   DefaultLanguage,
}

// Config holds all application configuration
type Config struct {
	AliveColor string
	DeadColor  string
	AliveChar  string
	DeadChar   string
	Language   Language
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
	if !isValidHexColor(c.AliveColor) {
		fmt.Printf("invalid alive color format: %s, using default\n", c.AliveColor)
		c.AliveColor = DefaultAliveColor
	}
	if !isValidHexColor(c.DeadColor) {
		fmt.Printf("invalid dead color format: %s, using default\n", c.DeadColor)
		c.DeadColor = DefaultDeadColor
	}
	if len([]rune(c.AliveChar)) != 1 {
		fmt.Printf("invalid alive character format: %s, using default\n", c.AliveChar)
		c.AliveChar = DefaultAliveChar
	}
	if len([]rune(c.DeadChar)) != 1 {
		fmt.Printf("invalid dead character format: %s, using default\n", c.DeadChar)
		c.DeadChar = DefaultDeadChar
	}
	if c.Language != English && c.Language != Chinese {
		fmt.Printf("invalid language %s, must be en or cn, using default language %s\n", c.Language.ToString(c.Language), DefaultLanguage.ToString(c.Language))
		c.Language = DefaultLanguage
	}
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
