package main

var (
	defaultRule = 30
	defaultCols = 80 // Should match DefaultWindowCols for consistency
)

// CellularAutomaton represents a 1D cellular automaton
type CellularAutomaton struct {
	currentRow []bool
	nextRow    []bool
	rule       int
	generation int // Track actual generation number for infinite mode
	cols       int
	boundary   BoundaryType // Boundary condition type
	ruleTable  [8]bool      // Pre-computed rule table for better performance
}

// NewCellularAutomaton creates a new cellular automaton instance
func NewCellularAutomaton(rule, cols int, boundary BoundaryType) *CellularAutomaton {
	ca := &CellularAutomaton{}
	ca.Reset(rule, cols, boundary)
	return ca
}

// computeRuleTable pre-computes the rule lookup table for better performance
func (ca *CellularAutomaton) computeRuleTable() {
	for i := range 8 {
		ca.ruleTable[i] = (ca.rule & (1 << i)) != 0
	}
}

// getNeighbors returns the left and right neighbors for a given cell index
// Optimized with early returns to reduce branching
func (ca *CellularAutomaton) getNeighbors(idx int) (left, right bool) {
	if ca.boundary == BoundaryPeriodic {
		// Periodic boundary: wrap around (most efficient case first)
		leftIdx := (idx - 1 + ca.cols) % ca.cols
		rightIdx := (idx + 1) % ca.cols
		return ca.currentRow[leftIdx], ca.currentRow[rightIdx]
	}

	// Calculate neighbors for fixed and reflective boundaries
	if idx > 0 {
		left = ca.currentRow[idx-1]
	} else if ca.boundary == BoundaryReflect {
		left = ca.currentRow[idx] // Reflect itself
	}

	if idx < ca.cols-1 {
		right = ca.currentRow[idx+1]
	} else if ca.boundary == BoundaryReflect {
		right = ca.currentRow[idx] // Reflect itself
	}

	return left, right
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

	// Convert boolean triplet to integer pattern (000 to 111)
	pattern := 0
	if left {
		pattern += 4 // Left bit (most significant)
	}
	if center {
		pattern += 2 // Center bit
	}
	if right {
		pattern++ // Right bit (least significant)
	}

	// Use pre-computed rule table for better performance
	return ca.ruleTable[pattern]
}

// Step advances the cellular automaton by one generation
func (ca *CellularAutomaton) Step() bool {
	// Calculate next generation based on current row
	for i := range ca.cols {
		ca.nextRow[i] = ca.getRuleBit(i)
	}

	// Swap current and next rows for next iteration (more efficient than copying)
	ca.currentRow, ca.nextRow = ca.nextRow, ca.currentRow

	ca.generation++ // Increment generation counter after computing
	return true
}

// GetCurrentRow returns the current row of the cellular automaton
func (ca *CellularAutomaton) GetCurrentRow() []bool {
	return ca.currentRow
}

// GetGeneration returns the current generation number
func (ca *CellularAutomaton) GetGeneration() int {
	return ca.generation
}

// Reset resets the cellular automaton to its initial state
func (ca *CellularAutomaton) Reset(rule, cols int, boundary BoundaryType) {
	// Input validation with defaults
	if cols <= 0 {
		cols = defaultCols
	}
	if rule < 0 || rule > 255 {
		rule = defaultRule
	}

	ca.rule = rule
	ca.cols = cols
	ca.boundary = boundary
	ca.currentRow = make([]bool, ca.cols)
	ca.nextRow = make([]bool, ca.cols)
	ca.generation = 0
	ca.computeRuleTable()

	// Initialize with center cell alive
	ca.currentRow[ca.cols/2] = true
}
