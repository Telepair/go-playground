package main

import (
	"math/rand/v2"
	"time"
)

var (
	defaultRows = 30
	defaultCols = 60
)

// GameOfLife represents Conway's Game of Life
type GameOfLife struct {
	currentGrid [][]bool
	nextGrid    [][]bool
	rows        int
	cols        int
	generation  int
	boundary    BoundaryType
	pattern     Pattern
}

// NewGameOfLife creates a new Game of Life instance
func NewGameOfLife(rows, cols int, boundary BoundaryType, pattern Pattern) *GameOfLife {
	if rows <= 0 {
		rows = defaultRows
	}
	if cols <= 0 {
		cols = defaultCols
	}

	game := &GameOfLife{
		rows:       rows,
		cols:       cols,
		generation: 0,
		boundary:   boundary,
		pattern:    pattern,
	}

	// Initialize grids
	game.currentGrid = make([][]bool, rows)
	game.nextGrid = make([][]bool, rows)
	for i := range rows {
		game.currentGrid[i] = make([]bool, cols)
		game.nextGrid[i] = make([]bool, cols)
	}

	// Set initial pattern
	game.setInitialPattern()

	return game
}

// setInitialPattern sets the initial pattern based on the selected pattern type
func (g *GameOfLife) setInitialPattern() {
	switch g.pattern {
	case PatternRandom:
		g.setRandomPattern()
	case PatternGlider:
		g.setGliderPattern()
	case PatternGliderGun:
		g.setGliderGunPattern()
	case PatternOscillator:
		g.setOscillatorPattern()
	case PatternPulsar:
		g.setPulsarPattern()
	case PatternPentomino:
		g.setPentominoPattern()
	default:
		g.setRandomPattern()
	}
}

// setRandomPattern creates a random initial pattern
func (g *GameOfLife) setRandomPattern() {
	// Use time-based seeding for game randomization (not cryptographic)
	// #nosec G115 - Conversion is safe for our use case
	seed := uint64(time.Now().UnixNano())

	// Create a new local random generator with the time-based seed
	// #nosec G404 - Using math/rand for game simulation, not cryptography
	rng := rand.New(rand.NewPCG(seed, seed))

	for i := range g.rows {
		for j := range g.cols {
			g.currentGrid[i][j] = rng.Float32() < 0.3 // 30% probability of being alive
		}
	}
}

// setGliderPattern creates a glider pattern
func (g *GameOfLife) setGliderPattern() {
	// Clear the grid first
	g.clearGrid()

	// Place glider in upper left area
	if g.rows >= 5 && g.cols >= 5 {
		startRow := 2
		startCol := 2
		// Glider pattern:
		//  X
		//   X
		// XXX
		pattern := [][]bool{
			{false, true, false},
			{false, false, true},
			{true, true, true},
		}
		g.placePattern(startRow, startCol, pattern)
	}
}

// setGliderGunPattern creates a Gosper glider gun pattern
func (g *GameOfLife) setGliderGunPattern() {
	// Clear the grid first
	g.clearGrid()

	// Only place if we have enough space (gun is 38 wide x 11 tall)
	if g.rows >= 15 && g.cols >= 40 {
		startRow := 2
		startCol := 2
		// Simplified glider gun pattern (part of the full gun)
		pattern := [][]bool{
			{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, true, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, false, false, true, true, false, false, false, false, false, false, true, true, false, false, false, false, false, false, false, false, false, false, false, false, true, true},
			{false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, true, false, false, false, false, true, true, false, false, false, false, false, false, false, false, false, false, false, false, true, true},
			{true, true, false, false, false, false, false, false, false, false, true, false, false, false, false, false, true, false, false, false, true, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
			{true, true, false, false, false, false, false, false, false, false, true, false, false, false, true, false, true, true, false, false, false, false, true, false, true, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, true, false, false, false, false, false, true, false, false, false, false, false, false, false, true, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, false, false, true, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		}
		g.placePattern(startRow, startCol, pattern)
	}
}

// setOscillatorPattern creates a blinker oscillator pattern
func (g *GameOfLife) setOscillatorPattern() {
	// Clear the grid first
	g.clearGrid()

	// Place multiple blinkers
	if g.rows >= 5 && g.cols >= 5 {
		// Horizontal blinker
		centerRow := g.rows / 2
		centerCol := g.cols / 2
		g.currentGrid[centerRow][centerCol-1] = true
		g.currentGrid[centerRow][centerCol] = true
		g.currentGrid[centerRow][centerCol+1] = true

		// Vertical blinker (offset)
		if g.cols >= 10 {
			offsetCol := centerCol + 5
			g.currentGrid[centerRow-1][offsetCol] = true
			g.currentGrid[centerRow][offsetCol] = true
			g.currentGrid[centerRow+1][offsetCol] = true
		}
	}
}

// setPulsarPattern creates a pulsar oscillator pattern
func (g *GameOfLife) setPulsarPattern() {
	// Clear the grid first
	g.clearGrid()

	// Pulsar is 15x15, only place if we have enough space
	if g.rows >= 17 && g.cols >= 17 {
		centerRow := g.rows / 2
		centerCol := g.cols / 2

		// Pulsar pattern (simplified version)
		offsets := [][]int{
			{-6, -4}, {-6, -3}, {-6, -2}, {-6, 2}, {-6, 3}, {-6, 4},
			{-4, -6}, {-4, -1}, {-4, 1}, {-4, 6},
			{-3, -6}, {-3, -1}, {-3, 1}, {-3, 6},
			{-2, -6}, {-2, -1}, {-2, 1}, {-2, 6},
			{-1, -4}, {-1, -3}, {-1, -2}, {-1, 2}, {-1, 3}, {-1, 4},
			{1, -4}, {1, -3}, {1, -2}, {1, 2}, {1, 3}, {1, 4},
			{2, -6}, {2, -1}, {2, 1}, {2, 6},
			{3, -6}, {3, -1}, {3, 1}, {3, 6},
			{4, -6}, {4, -1}, {4, 1}, {4, 6},
			{6, -4}, {6, -3}, {6, -2}, {6, 2}, {6, 3}, {6, 4},
		}

		for _, offset := range offsets {
			row := centerRow + offset[0]
			col := centerCol + offset[1]
			if row >= 0 && row < g.rows && col >= 0 && col < g.cols {
				g.currentGrid[row][col] = true
			}
		}
	}
}

// setPentominoPattern creates an R-pentomino pattern
func (g *GameOfLife) setPentominoPattern() {
	// Clear the grid first
	g.clearGrid()

	if g.rows >= 5 && g.cols >= 5 {
		centerRow := g.rows / 2
		centerCol := g.cols / 2
		// R-pentomino pattern:
		//  XX
		// XX
		//  X
		pattern := [][]bool{
			{false, true, true},
			{true, true, false},
			{false, true, false},
		}
		g.placePattern(centerRow-1, centerCol-1, pattern)
	}
}

// placePattern places a pattern at the specified position
func (g *GameOfLife) placePattern(startRow, startCol int, pattern [][]bool) {
	for i, row := range pattern {
		for j, cell := range row {
			newRow := startRow + i
			newCol := startCol + j
			if newRow >= 0 && newRow < g.rows && newCol >= 0 && newCol < g.cols {
				g.currentGrid[newRow][newCol] = cell
			}
		}
	}
}

// clearGrid clears all cells in the grid
func (g *GameOfLife) clearGrid() {
	for i := range g.rows {
		for j := range g.cols {
			g.currentGrid[i][j] = false
		}
	}
}

// countNeighbors counts the number of living neighbors for a cell
func (g *GameOfLife) countNeighbors(row, col int) int {
	count := 0

	// Check all 8 neighbors
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue // Skip the cell itself
			}

			neighborRow := row + dr
			neighborCol := col + dc

			// Handle boundary conditions
			if g.boundary == BoundaryPeriodic {
				// Wrap around for periodic boundary
				neighborRow = (neighborRow + g.rows) % g.rows
				neighborCol = (neighborCol + g.cols) % g.cols
			} else if g.boundary == BoundaryFixed {
				// Out of bounds cells are considered dead for fixed boundary
				if neighborRow < 0 || neighborRow >= g.rows || neighborCol < 0 || neighborCol >= g.cols {
					continue
				}
			}

			// Count living neighbors
			if neighborRow >= 0 && neighborRow < g.rows && neighborCol >= 0 && neighborCol < g.cols {
				if g.currentGrid[neighborRow][neighborCol] {
					count++
				}
			}
		}
	}

	return count
}

// Step advances the Game of Life by one generation
func (g *GameOfLife) Step() bool {
	// Apply Conway's Game of Life rules
	for i := range g.rows {
		for j := range g.cols {
			neighbors := g.countNeighbors(i, j)
			currentCell := g.currentGrid[i][j]

			// Conway's Game of Life rules:
			// 1. Any live cell with 2 or 3 live neighbors survives
			// 2. Any dead cell with exactly 3 live neighbors becomes a live cell
			// 3. All other live cells die, and all other dead cells stay dead

			if currentCell {
				// Cell is currently alive
				g.nextGrid[i][j] = (neighbors == 2 || neighbors == 3)
			} else {
				// Cell is currently dead
				g.nextGrid[i][j] = (neighbors == 3)
			}
		}
	}

	// Swap current and next grids
	g.currentGrid, g.nextGrid = g.nextGrid, g.currentGrid

	g.generation++
	return true
}

// GetCurrentGrid returns the current grid state
func (g *GameOfLife) GetCurrentGrid() [][]bool {
	return g.currentGrid
}

// GetGeneration returns the current generation number
func (g *GameOfLife) GetGeneration() int {
	return g.generation
}

// Reset resets the game to its initial state
func (g *GameOfLife) Reset() {
	// Clear both grids
	for i := range g.rows {
		for j := range g.cols {
			g.currentGrid[i][j] = false
			g.nextGrid[i][j] = false
		}
	}

	g.generation = 0
	g.setInitialPattern()
}

// SetPattern updates the pattern and resets the game
func (g *GameOfLife) SetPattern(pattern Pattern) {
	g.pattern = pattern
	g.Reset()
}
