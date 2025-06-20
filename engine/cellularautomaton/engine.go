// Package cellularautomaton provides the core engine implementations for cellular automata.
package cellularautomaton

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/telepair/go-playground/pkg/ui"
)

var _ ui.StepEngine = (*CellularAutomaton)(nil)

// CellularAutomaton represents a 1D cellular automaton
type CellularAutomaton struct {
	rule        Rule
	boundary    BoundaryType // Boundary condition type
	currentRule int
	ruleList    []Rule

	rows       int
	cols       int
	generation int // Track actual generation number for infinite mode
	currentRow []bool
	nextRow    []bool
	screen     *ui.Screen
	ruleTable  [8]bool // Pre-computed rule table for better performance
	buf        []rune
}

// New creates a new cellular automaton instance
func New(cfg Config) *CellularAutomaton {
	slog.Debug("NewCellularAutomaton", "cfg", cfg)
	cfg.Init()
	ca := &CellularAutomaton{
		rule:     cfg.GetRule(),
		rows:     defaultRows,
		cols:     defaultCols,
		boundary: cfg.GetBoundary(),
		ruleList: cfg.GetRuleList(),
	}
	// Set currentRule to point to the current rule in ruleList
	for i, rule := range ca.ruleList {
		if rule.Value == ca.rule.Value {
			ca.currentRule = i
			break
		}
	}
	ca.initial()
	return ca
}

// View returns the view of the cellular automaton
func (ca *CellularAutomaton) View() string {
	return ca.screen.View()
}

// Step advances the cellular automaton by one generation
func (ca *CellularAutomaton) Step() (int, bool) {
	// Increment generation counter
	ca.generation++

	// Handle first cell
	ca.nextRow[0] = ca.getRuleBit(0)

	// Handle middle cells with direct neighbor access for better performance
	for i := 1; i < ca.cols-1; i++ {
		// For middle cells, we can directly access neighbors without boundary checks
		left := ca.currentRow[i-1]
		center := ca.currentRow[i]
		right := ca.currentRow[i+1]

		// Convert boolean triplet to pattern using bit operations
		pattern := 0
		if left {
			pattern |= 4
		}
		if center {
			pattern |= 2
		}
		if right {
			pattern |= 1
		}

		ca.nextRow[i] = ca.ruleTable[pattern]
	}

	// Handle last cell
	if ca.cols > 1 {
		ca.nextRow[ca.cols-1] = ca.getRuleBit(ca.cols - 1)
	}
	// Swap current and next rows for next iteration (more efficient than copying)
	ca.currentRow, ca.nextRow = ca.nextRow, ca.currentRow

	ca.append()

	// Return generation and true (not finished, as cellular automaton runs indefinitely)
	return ca.generation, true
}

// Header returns the header text for the UI in the specified language
func (ca *CellularAutomaton) Header(lang ui.Language) string {
	if lang == ui.Chinese {
		return HeaderCN
	}
	return HeaderEN
}

// Status returns the status text for the UI in the specified language
func (ca *CellularAutomaton) Status(lang ui.Language) []ui.Status {
	if lang == ui.Chinese {
		return []ui.Status{
			{Label: "规则", Value: strconv.Itoa(ca.rule.Value)},
			{Label: "边界", Value: ca.boundary.ToString(lang)},
		}
	}
	return []ui.Status{
		{Label: "Rule", Value: strconv.Itoa(ca.rule.Value)},
		{Label: "Boundary", Value: ca.boundary.ToString(lang)},
	}
}

// HandleKeys returns the available keyboard controls for the cellular automaton
func (ca *CellularAutomaton) HandleKeys(lang ui.Language) []ui.Control {
	if lang == ui.Chinese {
		return []ui.Control{
			{Keys: []string{"T"}, Label: "规则"},
			{Keys: []string{"B"}, Label: "边界"},
		}
	}
	return []ui.Control{
		{Keys: []string{"T"}, Label: "Rule"},
		{Keys: []string{"B"}, Label: "Boundary"},
	}
}

// Handle handles the key press
func (ca *CellularAutomaton) Handle(key string) (bool, error) {
	slog.Debug("CellularAutomaton Handle", "key", key)
	key = strings.ToLower(key)
	switch key {
	case "t":
		ca.currentRule = (ca.currentRule + 1) % len(ca.ruleList)
		ca.rule = ca.ruleList[ca.currentRule]
		slog.Debug("CellularAutomaton Rule Changed", "key", key, "rule", ca.rule)
		ca.initial()
		return true, nil
	case "b":
		ca.boundary = (ca.boundary + 1) % 3
		slog.Debug("CellularAutomaton Boundary Changed", "key", key, "boundary", ca.boundary)
		return true, nil
	}
	slog.Debug("CellularAutomaton Unhandled Key", "key", key)
	return false, nil
}

// Reset resets the cellular automaton to its initial state
func (ca *CellularAutomaton) Reset(rows, cols int) error {
	slog.Debug("CellularAutomaton Reset", "rows", rows, "cols", cols)
	ca.rows = rows
	ca.cols = cols
	ca.initial()
	return nil
}

// IsFinished returns whether the cellular automaton has finished execution
func (ca *CellularAutomaton) IsFinished() bool {
	return false
}

// Stop stops the cellular automaton execution
func (ca *CellularAutomaton) Stop() {}

// computeRuleTable pre-computes the rule lookup table for better performance
func (ca *CellularAutomaton) computeRuleTable() {
	for i := range 8 {
		ca.ruleTable[i] = (ca.rule.Value & (1 << i)) != 0
	}
}

// getNeighbors returns the left and right neighbors for a given cell index
// Optimized with pre-computed indices to reduce branching
func (ca *CellularAutomaton) getNeighbors(idx int) (left, right bool) {
	// Input validation to prevent index out of bounds
	if idx < 0 || idx >= ca.cols || ca.currentRow == nil {
		return false, false
	}

	switch ca.boundary {
	case BoundaryPeriodic:
		// Optimized periodic boundary without modulo
		leftIdx := idx - 1
		if leftIdx < 0 {
			leftIdx = ca.cols - 1
		}

		rightIdx := idx + 1
		if rightIdx >= ca.cols {
			rightIdx = 0
		}

		return ca.currentRow[leftIdx], ca.currentRow[rightIdx]

	case BoundaryReflect:
		// Reflective boundary: mirror the grid at boundaries
		leftIdx := idx - 1
		if leftIdx < 0 {
			// For left boundary, reflect to position 1 (mirror of -1 around 0)
			leftIdx = 1
			if leftIdx >= ca.cols {
				leftIdx = ca.cols - 1 // Fallback for very small grids
			}
		}

		rightIdx := idx + 1
		if rightIdx >= ca.cols {
			// For right boundary, reflect to position cols-2 (mirror of cols around cols-1)
			rightIdx = max(ca.cols-2, 0)
		}

		return ca.currentRow[leftIdx], ca.currentRow[rightIdx]

	default: // BoundaryFixed
		// Fixed boundary: return false for out-of-bounds
		left = idx > 0 && ca.currentRow[idx-1]
		right = idx < ca.cols-1 && ca.currentRow[idx+1]
		return left, right
	}
}

// getRuleBit returns the next state for a cell based on its neighborhood
func (ca *CellularAutomaton) getRuleBit(idx int) bool {
	// Input validation
	if idx < 0 || idx >= ca.cols {
		return false
	}

	// Get neighbors using the optimized neighbor function
	left, right := ca.getNeighbors(idx)
	center := ca.currentRow[idx]

	// Convert boolean triplet to integer pattern (000 to 111) using bit operations
	pattern := 0
	if left {
		pattern |= 4 // Left bit (most significant)
	}
	if center {
		pattern |= 2 // Center bit
	}
	if right {
		pattern |= 1 // Right bit (least significant)
	}

	// Use pre-computed rule table for better performance
	return ca.ruleTable[pattern]
}

func (ca *CellularAutomaton) initial() {
	if ca.screen == nil {
		ca.screen = ui.NewScreen(ca.rows, ca.cols)
	} else {
		ca.screen.SetSize(ca.cols, ca.rows)
		ca.screen.Reset()
	}
	if ca.rule.ActiveChar == 0 {
		ca.rule.ActiveChar = DefaultAliveChar
	}
	if ca.rule.DeadChar == 0 {
		ca.rule.DeadChar = DefaultDeadChar
	}
	if ca.rule.ActiveColor == "" {
		ca.rule.ActiveColor = DefaultAliveColor
	}
	if ca.rule.DeadColor == "" {
		ca.rule.DeadColor = DefaultDeadColor
	}
	ca.screen.SetCharColor(ca.rule.ActiveChar, ca.rule.ActiveColor)
	ca.screen.SetCharColor(ca.rule.DeadChar, ca.rule.DeadColor)
	ca.screen.Reset()
	ca.currentRow = make([]bool, ca.cols)
	ca.nextRow = make([]bool, ca.cols)
	ca.buf = make([]rune, ca.cols)
	ca.generation = 0
	ca.computeRuleTable()
	ca.currentRow[ca.cols/2] = true
	ca.append()
}

func (ca *CellularAutomaton) append() {
	for i := range ca.cols {
		if ca.currentRow[i] {
			ca.buf[i] = ca.rule.ActiveChar
		} else {
			ca.buf[i] = ca.rule.DeadChar
		}
	}
	ca.screen.Append(ca.buf)
}
