package main

type Window struct {
	rows uint
	cols uint
}

// CellularAutomaton represents a 1D cellular automaton
type CellularAutomaton struct {
	grid       [][]bool
	rule       uint
	currentRow uint
	generation uint // Track actual generation number for infinite mode
	maxSteps   uint // if 0, infinite mode
	window     Window
	boundary   BoundaryType // Boundary condition type
}

// NewCellularAutomaton creates a new cellular automaton instance
func NewCellularAutomaton(rule, maxSteps uint, window *Window, boundary BoundaryType) *CellularAutomaton {
	if window == nil {
		window = &Window{
			rows: DefaultWindowRows,
			cols: DefaultWindowCols,
		}
	}

	ca := &CellularAutomaton{
		rule:       rule,
		grid:       make([][]bool, window.rows),
		maxSteps:   maxSteps,
		currentRow: 0,
		generation: 0,
		window:     *window,
		boundary:   boundary,
	}

	// Initialize grid
	for i := range ca.grid {
		ca.grid[i] = make([]bool, window.cols)
	}

	// Set initial state - single cell in the middle
	ca.grid[0][window.cols/2] = true

	return ca
}

// getRuleBit returns the bit value for a given pattern according to the rule
func (ca *CellularAutomaton) getRuleBit(left, center, right bool) bool {
	// Convert boolean pattern to integer (000 to 111)
	pattern := 0
	if left {
		pattern += 4
	}
	if center {
		pattern += 2
	}
	if right {
		pattern += 1
	}

	// Check if the corresponding bit in the rule is set
	return (ca.rule & (1 << pattern)) != 0
}

// getNeighbors returns the left, center, and right neighbors for a given position
// applying the specified boundary conditions
func (ca *CellularAutomaton) getNeighbors(i uint) (bool, bool, bool) {
	center := ca.grid[ca.currentRow][i]

	var left, right bool

	switch ca.boundary {
	case BoundaryPeriodic:
		// Periodic boundary: wrap around
		left = ca.grid[ca.currentRow][(i-1+ca.window.cols)%ca.window.cols]
		right = ca.grid[ca.currentRow][(i+1)%ca.window.cols]

	case BoundaryFixed:
		// Fixed boundary: use false (0) for boundary cells
		if i == 0 {
			left = false
		} else {
			left = ca.grid[ca.currentRow][i-1]
		}

		if i == ca.window.cols-1 {
			right = false
		} else {
			right = ca.grid[ca.currentRow][i+1]
		}

	case BoundaryReflect:
		// Reflective boundary: boundary cells reflect themselves
		if i == 0 {
			left = center
		} else {
			left = ca.grid[ca.currentRow][i-1]
		}

		if i == ca.window.cols-1 {
			right = center
		} else {
			right = ca.grid[ca.currentRow][i+1]
		}
	}

	return left, center, right
}

// Step advances the cellular automaton by one generation
func (ca *CellularAutomaton) Step() bool {
	if ca.maxSteps > 0 && ca.generation >= ca.maxSteps {
		return false
	}
	ca.generation++ // Always increment the actual generation counter
	nextRow := (ca.currentRow + 1) % ca.window.rows

	// Compute next generation
	for i := uint(0); i < ca.window.cols; i++ {
		// Get neighbors using boundary conditions
		left, center, right := ca.getNeighbors(i)

		// Apply rule
		ca.grid[nextRow][i] = ca.getRuleBit(left, center, right)
	}
	ca.currentRow = nextRow
	return true // Always return true for infinite mode
}

// GetGeneration returns the current generation number
func (ca *CellularAutomaton) GetGeneration() uint {
	return ca.generation
}

// GetGrid returns the current grid state in chronological order
// Note: The returned slice shares the underlying data with the internal grid
// Callers should not modify the returned data
func (ca *CellularAutomaton) GetGrid() [][]bool {
	// Calculate the starting row for chronological order (oldest first)
	startRow := (ca.currentRow + 1) % ca.window.rows

	// Pre-allocate result slice with known capacity
	result := make([][]bool, 0, ca.window.rows)

	// Return grid in chronological order
	for i := uint(0); i < ca.window.rows; i++ {
		rowIndex := (startRow + i) % ca.window.rows
		result = append(result, ca.grid[rowIndex])
	}
	return result
}
