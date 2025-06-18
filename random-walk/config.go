// Package main implements Random Walk visualization with various patterns and configurations.
package main

import (
	"fmt"
	"strings"
	"time"
)

// WalkMode represents different walking modes
type WalkMode int

// WalkMode constants
const (
	ModeSingleWalker     WalkMode = iota // Single walker
	ModeMultiWalker                      // Multiple walkers
	ModeTrailMode                        // Show trails
	ModeBrownianMotion                   // Brownian motion simulation
	ModeSelfAvoidingWalk                 // Self-avoiding walk
	ModeLevyFlight                       // Lévy flight pattern
)

// ToString returns the string representation of walk mode
func (wm WalkMode) ToString(language Language) string {
	switch wm {
	case ModeSingleWalker:
		if language == Chinese {
			return "单粒子"
		}
		return "Single Walker"
	case ModeMultiWalker:
		if language == Chinese {
			return "多粒子"
		}
		return "Multi Walker"
	case ModeTrailMode:
		if language == Chinese {
			return "轨迹模式"
		}
		return "Trail Mode"
	case ModeBrownianMotion:
		if language == Chinese {
			return "布朗运动"
		}
		return "Brownian Motion"
	case ModeSelfAvoidingWalk:
		if language == Chinese {
			return "自避行走"
		}
		return "Self-Avoiding"
	case ModeLevyFlight:
		if language == Chinese {
			return "莱维飞行"
		}
		return "Lévy Flight"
	default:
		if language == Chinese {
			return "单粒子"
		}
		return "Single Walker"
	}
}

// Direction represents movement direction
type Direction int

// Direction constants
const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
	DirectionUpLeft
	DirectionUpRight
	DirectionDownLeft
	DirectionDownRight
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
	DefaultWalkMode    = ModeSingleWalker      // Default walk mode
	DefaultWalkerCount = 3                     // Default number of walkers for multi-walker mode
	MaxWalkerCount     = 10                    // Maximum number of walkers
	DefaultTrailLength = 100                   // Default trail length
	MaxTrailLength     = 500                   // Maximum trail length

	// Colors
	DefaultWalkerColor = "#FF00FF" // Default walker color (magenta)
	DefaultTrailColor  = "#0088FF" // Default trail color (blue)
	DefaultEmptyColor  = "#000000" // Default empty cell color (black)

	// Characters
	DefaultWalkerChar = "●" // Default walker character
	DefaultTrailChar  = "·" // Default trail character
	DefaultEmptyChar  = " " // Default empty cell character

	// Walker colors for multi-walker mode
	Walker1Color  = "#FF0000" // Red
	Walker2Color  = "#00FF00" // Green
	Walker3Color  = "#0000FF" // Blue
	Walker4Color  = "#FFFF00" // Yellow
	Walker5Color  = "#FF00FF" // Magenta
	Walker6Color  = "#00FFFF" // Cyan
	Walker7Color  = "#FFA500" // Orange
	Walker8Color  = "#800080" // Purple
	Walker9Color  = "#FFC0CB" // Pink
	Walker10Color = "#A52A2A" // Brown

	// Default values
	DefaultLogFile         = "debug.log"     // Default log file path
	DefaultProfileInterval = 5 * time.Second // Default profile information output interval
	DefaultProfilePort     = 6060            // Default profile server port
)

// GetWalkerColor returns the color for a walker based on its index
func GetWalkerColor(index int) string {
	colors := []string{
		Walker1Color, Walker2Color, Walker3Color, Walker4Color, Walker5Color,
		Walker6Color, Walker7Color, Walker8Color, Walker9Color, Walker10Color,
	}
	if index >= 0 && index < len(colors) {
		return colors[index]
	}
	return DefaultWalkerColor
}

// DefaultConfig is the default configuration
var DefaultConfig = Config{
	WalkerColor: DefaultWalkerColor,
	TrailColor:  DefaultTrailColor,
	EmptyColor:  DefaultEmptyColor,
	WalkerChar:  DefaultWalkerChar,
	TrailChar:   DefaultTrailChar,
	EmptyChar:   DefaultEmptyChar,
	Language:    DefaultLanguage,
}

// Config holds all application configuration
type Config struct {
	WalkerColor string
	TrailColor  string
	EmptyColor  string
	WalkerChar  string
	TrailChar   string
	EmptyChar   string
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
	if !isValidHexColor(c.WalkerColor) {
		fmt.Printf("invalid walker color format: %s, using default\n", c.WalkerColor)
		c.WalkerColor = DefaultWalkerColor
	}
	if !isValidHexColor(c.TrailColor) {
		fmt.Printf("invalid trail color format: %s, using default\n", c.TrailColor)
		c.TrailColor = DefaultTrailColor
	}
	if !isValidHexColor(c.EmptyColor) {
		fmt.Printf("invalid empty color format: %s, using default\n", c.EmptyColor)
		c.EmptyColor = DefaultEmptyColor
	}
	if len([]rune(c.WalkerChar)) != 1 {
		fmt.Printf("invalid walker character format: %s, using default\n", c.WalkerChar)
		c.WalkerChar = DefaultWalkerChar
	}
	if len([]rune(c.TrailChar)) != 1 {
		fmt.Printf("invalid trail character format: %s, using default\n", c.TrailChar)
		c.TrailChar = DefaultTrailChar
	}
	if len([]rune(c.EmptyChar)) != 1 {
		fmt.Printf("invalid empty character format: %s, using default\n", c.EmptyChar)
		c.EmptyChar = DefaultEmptyChar
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
