package cellularautomaton

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/telepair/go-playground/pkg/ui"
)

var (
	// HeaderCN is the Chinese header text for cellular automaton
	HeaderCN = "🚀 元胞自动机 🚀"
	// HeaderEN is the English header text for cellular automaton
	HeaderEN = "🚀 Cellular Automaton 🚀"

	// DefaultAliveColor is the default alive cell color
	DefaultAliveColor = lipgloss.Color("#FFFFFF")
	// DefaultDeadColor is the default dead cell color
	DefaultDeadColor = lipgloss.Color("#000000")
	// DefaultAliveChar is the default alive cell character
	DefaultAliveChar = '█'
	// DefaultDeadChar is the default dead cell character
	DefaultDeadChar = ' '

	// Rules is the default rules for cellular automaton
	Rules = map[int]Rule{
		30:  {Value: 30, ActiveChar: DefaultAliveChar, DeadChar: DefaultDeadChar, ActiveColor: DefaultAliveColor, DeadColor: DefaultDeadColor},
		90:  {Value: 90, ActiveChar: DefaultAliveChar, DeadChar: DefaultDeadChar, ActiveColor: lipgloss.Color("#00FF00"), DeadColor: lipgloss.Color("#FF0000")},
		110: {Value: 110, ActiveChar: DefaultAliveChar, DeadChar: DefaultDeadChar, ActiveColor: DefaultAliveColor, DeadColor: DefaultDeadColor},
		150: {Value: 150, ActiveChar: DefaultAliveChar, DeadChar: DefaultDeadChar, ActiveColor: DefaultAliveColor, DeadColor: DefaultDeadColor},
		184: {Value: 184, ActiveChar: '🚗', DeadChar: DefaultDeadChar, ActiveColor: DefaultAliveColor, DeadColor: DefaultDeadColor},
	}

	defaultRows = 20
	defaultCols = 40
)

// Config represents the configuration for a cellular automaton
type Config struct {
	Rule       int
	Boundary   int
	AliveChar  string
	DeadChar   string
	AliveColor string
	DeadColor  string

	rule     Rule
	boundary BoundaryType
	ruleList []Rule
}

// Init initializes the configuration with default values and validates the settings
func (c *Config) Init() {
	if rule, ok := Rules[c.Rule]; ok {
		c.rule = rule
	} else {
		// Default to rule 30 if specified rule doesn't exist
		c.rule = Rules[30]
	}
	if c.AliveColor != "" {
		c.rule.ActiveColor = lipgloss.Color(c.AliveColor)
	}
	if c.DeadColor != "" {
		c.rule.DeadColor = lipgloss.Color(c.DeadColor)
	}
	if len(c.AliveChar) > 0 {
		c.rule.ActiveChar = rune(c.AliveChar[0])
	}
	if len(c.DeadChar) > 0 {
		c.rule.DeadChar = rune(c.DeadChar[0])
	}
	c.boundary = BoundaryType(c.Boundary)
	c.ruleList = make([]Rule, 0, len(Rules))
	// Create sorted list of rule values for deterministic order
	ruleValues := make([]int, 0, len(Rules))
	for ruleValue := range Rules {
		ruleValues = append(ruleValues, ruleValue)
	}
	// Sort rule values to ensure consistent order
	for i := 0; i < len(ruleValues); i++ {
		for j := i + 1; j < len(ruleValues); j++ {
			if ruleValues[i] > ruleValues[j] {
				ruleValues[i], ruleValues[j] = ruleValues[j], ruleValues[i]
			}
		}
	}
	// Build rule list in sorted order
	for _, ruleValue := range ruleValues {
		rule := Rules[ruleValue]
		rule.ActiveColor = c.rule.ActiveColor
		rule.DeadColor = c.rule.DeadColor
		rule.ActiveChar = c.rule.ActiveChar
		rule.DeadChar = c.rule.DeadChar
		c.ruleList = append(c.ruleList, rule)
	}
}

// GetRule returns the configured cellular automaton rule
func (c *Config) GetRule() Rule {
	return c.rule
}

// GetRuleList returns the list of available cellular automaton rules
func (c *Config) GetRuleList() []Rule {
	return c.ruleList
}

// GetBoundary returns the configured boundary type for the cellular automaton
func (c *Config) GetBoundary() BoundaryType {
	return BoundaryType(c.Boundary)
}

// BoundaryType represents the boundary type of the cellular automaton
type BoundaryType int

// BoundaryType constants
const (
	BoundaryPeriodic BoundaryType = iota // Periodic boundary (default)
	BoundaryFixed                        // Fixed boundary (0 values)
	BoundaryReflect                      // Reflective boundary (mirror)
)

// ToString returns the string representation of boundary type
func (bt BoundaryType) ToString(language ui.Language) string {
	switch bt {
	case BoundaryPeriodic:
		if language == ui.Chinese {
			return "周期"
		}
		return "Periodic"
	case BoundaryFixed:
		if language == ui.Chinese {
			return "固定"
		}
		return "Fixed"
	case BoundaryReflect:
		if language == ui.Chinese {
			return "反射"
		}
		return "Reflect"
	}
	if language == ui.Chinese {
		return "周期"
	}
	return "Periodic"
}

// Rule represents a cellular automaton rule
type Rule struct {
	Value       int
	ActiveChar  rune
	DeadChar    rune
	ActiveColor lipgloss.Color
	DeadColor   lipgloss.Color
}
