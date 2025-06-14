package main

var (
	defaultRule     = 30
	defaultMaxSteps = 0
	defaultCols     = 1
)

// CellularAutomaton represents a 1D cellular automaton
type CellularAutomaton struct {
	currentRow []bool
	nextRow    []bool
	rule       int
	generation int // Track actual generation number for infinite mode
	maxSteps   int // if 0, infinite mode
	cols       int
	boundary   BoundaryType // Boundary condition type
	ruleTable  [8]bool      // Pre-computed rule table for better performance
}

// NewCellularAutomaton creates a new cellular automaton instance
func NewCellularAutomaton(rule, maxSteps, cols int, boundary BoundaryType) *CellularAutomaton {
	if cols <= 0 {
		cols = defaultCols
	}
	if rule < 0 || rule > 255 {
		rule = defaultRule
	}
	if maxSteps < 0 {
		maxSteps = defaultMaxSteps
	}
	ca := &CellularAutomaton{
		rule:       rule,
		generation: 0,
		maxSteps:   maxSteps,
		boundary:   boundary,
		cols:       cols,
		currentRow: make([]bool, cols),
		nextRow:    make([]bool, cols),
	}

	// Pre-compute rule table for better performance
	ca.computeRuleTable()

	// Set initial seed in the current row so it displays immediately
	ca.currentRow[cols/2] = true
	return ca
}

// computeRuleTable pre-computes the rule lookup table for better performance
func (ca *CellularAutomaton) computeRuleTable() {
	for i := range 8 {
		ca.ruleTable[i] = (ca.rule & (1 << i)) != 0
	}
}

// getNeighbors returns the left and right neighbors for a given cell index
// This function handles all boundary conditions in one place for clarity
func (ca *CellularAutomaton) getNeighbors(idx int) (left, right bool) {
	switch ca.boundary {
	case BoundaryPeriodic:
		// Periodic boundary: wrap around
		leftIdx := (idx - 1 + ca.cols) % ca.cols
		rightIdx := (idx + 1) % ca.cols
		left = ca.currentRow[leftIdx]
		right = ca.currentRow[rightIdx]

	case BoundaryFixed:
		// Fixed boundary: use false (0) for boundary cells
		if idx > 0 {
			left = ca.currentRow[idx-1]
		}
		if idx < ca.cols-1 {
			right = ca.currentRow[idx+1]
		}

	case BoundaryReflect:
		// Reflective boundary: boundary cells reflect themselves
		if idx == 0 {
			left = ca.currentRow[idx] // Reflect itself
		} else {
			left = ca.currentRow[idx-1]
		}

		if idx == ca.cols-1 {
			right = ca.currentRow[idx] // Reflect itself
		} else {
			right = ca.currentRow[idx+1]
		}
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
	if ca.maxSteps > 0 && ca.generation >= ca.maxSteps {
		return false
	}

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
// This is more efficient than creating new slices
func (ca *CellularAutomaton) Reset() {
	// Clear both rows efficiently
	for i := range ca.currentRow {
		ca.currentRow[i] = false
		ca.nextRow[i] = false
	}

	// Set initial seed in current row
	ca.currentRow[ca.cols/2] = true
	ca.generation = 0
}

// SetRule updates the rule and recomputes the rule table
func (ca *CellularAutomaton) SetRule(rule int) {
	ca.rule = rule
	ca.computeRuleTable()
}

// IsFinished returns true if the cellular automaton has reached its maximum steps
func (ca *CellularAutomaton) IsFinished() bool {
	return ca.maxSteps > 0 && ca.generation >= ca.maxSteps
}
