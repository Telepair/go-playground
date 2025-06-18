package main

import (
	"log/slog"
	"math/rand/v2"
	"time"
)

// Use constants from config.go instead of separate variables

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
	slog.Debug("NewGameOfLife", "rows", rows, "cols", cols, "boundary", boundary, "pattern", pattern)
	game := &GameOfLife{
		rows:       rows,
		cols:       cols,
		boundary:   boundary,
		pattern:    pattern,
		generation: 0,
	}
	game.Init()
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

	// Optimized random generation - use Uint32 for better performance
	for i := range g.rows {
		for j := range g.cols {
			// Use bit manipulation for 30% probability (faster than float comparison)
			g.currentGrid[i][j] = rng.Uint32()%10 < 3 // 30% probability of being alive
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
// Optimized version with direct neighbor checking
func (g *GameOfLife) countNeighbors(row, col int) int {
	// Input validation to prevent index out of bounds
	if row < 0 || row >= g.rows || col < 0 || col >= g.cols || g.currentGrid == nil {
		return 0
	}

	count := 0

	// Define neighbor offsets for direct access
	neighbors := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, offset := range neighbors {
		neighborRow := row + offset[0]
		neighborCol := col + offset[1]

		// Handle boundary conditions
		if g.boundary == BoundaryPeriodic {
			// Wrap around for periodic boundary
			neighborRow = (neighborRow + g.rows) % g.rows
			neighborCol = (neighborCol + g.cols) % g.cols

			// Additional bounds check for safety after modulo operation
			// Note: modulo operation can still result in negative values in some Go versions
			if neighborRow >= 0 && neighborRow < g.rows &&
				neighborCol >= 0 && neighborCol < g.cols &&
				g.currentGrid != nil && len(g.currentGrid) > neighborRow &&
				len(g.currentGrid[neighborRow]) > neighborCol {
				if g.currentGrid[neighborRow][neighborCol] {
					count++
				}
			}
		} else {
			// Fixed boundary: skip out-of-bounds cells
			if neighborRow >= 0 && neighborRow < g.rows &&
				neighborCol >= 0 && neighborCol < g.cols {
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

// Init initializes the game of life
func (g *GameOfLife) Init() {
	slog.Debug("GameOfLife Init", "rows", g.rows, "cols", g.cols, "boundary", g.boundary, "pattern", g.pattern)
	if g.rows <= MinRows {
		slog.Warn("GameOfLife rows is less than MinRows, using default rows", "rows", g.rows, "minRows", MinRows, "defaultRows", DefaultRows)
		g.rows = DefaultRows
	}
	if g.cols <= MinCols {
		slog.Warn("GameOfLife cols is less than MinCols, using default cols", "cols", g.cols, "minCols", MinCols, "defaultCols", DefaultCols)
		g.cols = DefaultCols
	}
	g.currentGrid = make([][]bool, g.rows)
	g.nextGrid = make([][]bool, g.rows)
	for i := range g.rows {
		g.currentGrid[i] = make([]bool, g.cols)
		g.nextGrid[i] = make([]bool, g.cols)
	}
	g.setInitialPattern()
}

// Reset resets the game to its initial state
func (g *GameOfLife) Reset(rows, cols int, boundary BoundaryType, pattern Pattern) {
	slog.Debug("GameOfLife Reset", "rows", rows, "cols", cols, "boundary", boundary, "pattern", pattern)
	g.rows = rows
	g.cols = cols
	g.boundary = boundary
	g.pattern = pattern
	g.generation = 0
	g.Init()
}
