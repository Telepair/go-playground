package main

import (
	"fmt"
	"time"
)

// BoundaryType represents the boundary type of the cellular automaton
type BoundaryType int

// BoundaryType constants
const (
	BoundaryPeriodic BoundaryType = iota // Periodic boundary (default)
	BoundaryFixed                        // Fixed boundary (0 values)
	BoundaryReflect                      // Reflective boundary (mirror)
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
	case BoundaryReflect:
		if language == Chinese {
			return "反射"
		}
		return "Reflect"
	}
	if language == Chinese {
		return "周期"
	}
	return "Periodic"
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

// Application constants
const (
	DefaultRows = 30 // Default window rows
	DefaultCols = 80 // Default window columns
	MinRows     = 10 // Minimum window rows
	MinCols     = 20 // Minimum window columns

	// Rule validation
	DefaultRule = 30  // Default cellular automaton rule
	MinRule     = 0   // Minimum rule number
	MaxRule     = 255 // Maximum rule number

	// Timing constants
	DefaultRefreshRate = 200 * time.Millisecond // Default refresh rate in milliseconds
	MinRefreshRate     = 10 * time.Millisecond  // Minimum refresh rate in milliseconds

	// Default values
	DefaultLanguage = English          // Default language
	DefaultBoundary = BoundaryPeriodic // Default boundary type

	// Colors
	DefaultAliveColor = "#FFFFFF" // Default alive cell color
	DefaultDeadColor  = "#000000" // Default dead cell color

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
	Rule:       DefaultRule,
	AliveColor: DefaultAliveColor,
	DeadColor:  DefaultDeadColor,
	AliveChar:  DefaultAliveChar,
	DeadChar:   DefaultDeadChar,
	Language:   DefaultLanguage,
}

// Config holds all application configuration
type Config struct {
	Rule       int
	AliveColor string
	DeadColor  string
	AliveChar  string
	DeadChar   string
	Language   Language
}

// SetLang sets the language
func (c *Config) SetLang(lang string) {
	if lang == "cn" || lang == "zh" {
		c.Language = Chinese
	} else {
		c.Language = English
	}
}

// Check validates the configuration
func (c *Config) Check() {
	if c.Rule < MinRule || c.Rule > MaxRule {
		fmt.Printf("invalid rule %d, must be between %d and %d, using default rule %d\n", c.Rule, MinRule, MaxRule, DefaultRule)
		c.Rule = DefaultRule
	}

	if c.Language != English && c.Language != Chinese {
		fmt.Printf("invalid language %s, must be en or cn, using default language %s\n", c.Language.ToString(c.Language), DefaultLanguage.ToString(c.Language))
		c.Language = DefaultLanguage
	}
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
