package main

import (
	"reflect"
	"testing"
)

// Test NewGameOfLife constructor
func TestNewGameOfLife(t *testing.T) {
	tests := []struct {
		name     string
		rows     int
		cols     int
		boundary BoundaryType
		pattern  Pattern
		wantRows int
		wantCols int
	}{
		{
			name:     "Valid dimensions",
			rows:     10,
			cols:     20,
			boundary: BoundaryPeriodic,
			pattern:  PatternRandom,
			wantRows: 10,
			wantCols: 20,
		},
		{
			name:     "Zero rows uses default",
			rows:     0,
			cols:     20,
			boundary: BoundaryPeriodic,
			pattern:  PatternRandom,
			wantRows: DefaultWindowRows,
			wantCols: 20,
		},
		{
			name:     "Zero cols uses default",
			rows:     10,
			cols:     0,
			boundary: BoundaryPeriodic,
			pattern:  PatternRandom,
			wantRows: 10,
			wantCols: DefaultWindowCols,
		},
		{
			name:     "Negative dimensions use defaults",
			rows:     -5,
			cols:     -10,
			boundary: BoundaryFixed,
			pattern:  PatternGlider,
			wantRows: DefaultWindowRows,
			wantCols: DefaultWindowCols,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGameOfLife(tt.rows, tt.cols, tt.boundary, tt.pattern)

			if game == nil {
				t.Fatal("NewGameOfLife returned nil")
			}

			if game.rows != tt.wantRows {
				t.Errorf("rows = %d, want %d", game.rows, tt.wantRows)
			}

			if game.cols != tt.wantCols {
				t.Errorf("cols = %d, want %d", game.cols, tt.wantCols)
			}

			if game.boundary != tt.boundary {
				t.Errorf("boundary = %v, want %v", game.boundary, tt.boundary)
			}

			if game.pattern != tt.pattern {
				t.Errorf("pattern = %v, want %v", game.pattern, tt.pattern)
			}

			if game.generation != 0 {
				t.Errorf("generation = %d, want 0", game.generation)
			}

			// Check that grids are initialized
			if len(game.currentGrid) != tt.wantRows {
				t.Errorf("currentGrid rows = %d, want %d", len(game.currentGrid), tt.wantRows)
			}

			if len(game.nextGrid) != tt.wantRows {
				t.Errorf("nextGrid rows = %d, want %d", len(game.nextGrid), tt.wantRows)
			}

			for i := range game.currentGrid {
				if len(game.currentGrid[i]) != tt.wantCols {
					t.Errorf("currentGrid[%d] cols = %d, want %d", i, len(game.currentGrid[i]), tt.wantCols)
				}
				if len(game.nextGrid[i]) != tt.wantCols {
					t.Errorf("nextGrid[%d] cols = %d, want %d", i, len(game.nextGrid[i]), tt.wantCols)
				}
			}
		})
	}
}

// Test countNeighbors function with different boundary types
func TestCountNeighbors(t *testing.T) {
	tests := []struct {
		name     string
		rows     int
		cols     int
		boundary BoundaryType
		grid     [][]bool
		row      int
		col      int
		want     int
	}{
		{
			name:     "Center cell with all neighbors alive",
			rows:     3,
			cols:     3,
			boundary: BoundaryFixed,
			grid: [][]bool{
				{true, true, true},
				{true, false, true},
				{true, true, true},
			},
			row:  1,
			col:  1,
			want: 8,
		},
		{
			name:     "Center cell with no neighbors alive",
			rows:     3,
			cols:     3,
			boundary: BoundaryFixed,
			grid: [][]bool{
				{false, false, false},
				{false, false, false},
				{false, false, false},
			},
			row:  1,
			col:  1,
			want: 0,
		},
		{
			name:     "Corner cell with fixed boundary",
			rows:     3,
			cols:     3,
			boundary: BoundaryFixed,
			grid: [][]bool{
				{false, true, false},
				{true, true, false},
				{false, false, false},
			},
			row:  0,
			col:  0,
			want: 3,
		},
		{
			name:     "Corner cell with periodic boundary",
			rows:     3,
			cols:     3,
			boundary: BoundaryPeriodic,
			grid: [][]bool{
				{false, false, true},
				{false, false, false},
				{true, false, true},
			},
			row:  0,
			col:  0,
			want: 3, // neighbors at (2,2), (2,0), (0,2)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &GameOfLife{
				currentGrid: tt.grid,
				rows:        tt.rows,
				cols:        tt.cols,
				boundary:    tt.boundary,
			}

			got := game.countNeighbors(tt.row, tt.col)
			if got != tt.want {
				t.Errorf("countNeighbors(%d, %d) = %d, want %d", tt.row, tt.col, got, tt.want)
			}
		})
	}
}

// Test Step function with known patterns
func TestStep(t *testing.T) {
	tests := []struct {
		name         string
		initialGrid  [][]bool
		expectedGrid [][]bool
		boundary     BoundaryType
	}{
		{
			name: "Blinker oscillator horizontal to vertical",
			initialGrid: [][]bool{
				{false, false, false, false, false},
				{false, false, false, false, false},
				{false, true, true, true, false},
				{false, false, false, false, false},
				{false, false, false, false, false},
			},
			expectedGrid: [][]bool{
				{false, false, false, false, false},
				{false, false, true, false, false},
				{false, false, true, false, false},
				{false, false, true, false, false},
				{false, false, false, false, false},
			},
			boundary: BoundaryFixed,
		},
		{
			name: "Block still life",
			initialGrid: [][]bool{
				{false, false, false, false},
				{false, true, true, false},
				{false, true, true, false},
				{false, false, false, false},
			},
			expectedGrid: [][]bool{
				{false, false, false, false},
				{false, true, true, false},
				{false, true, true, false},
				{false, false, false, false},
			},
			boundary: BoundaryFixed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows := len(tt.initialGrid)
			cols := len(tt.initialGrid[0])

			game := &GameOfLife{
				rows:       rows,
				cols:       cols,
				generation: 0,
				boundary:   tt.boundary,
			}

			// Initialize grids
			game.currentGrid = make([][]bool, rows)
			game.nextGrid = make([][]bool, rows)
			for i := range rows {
				game.currentGrid[i] = make([]bool, cols)
				game.nextGrid[i] = make([]bool, cols)
				copy(game.currentGrid[i], tt.initialGrid[i])
			}

			// Execute one step
			game.Step()

			// Check the result
			if !reflect.DeepEqual(game.currentGrid, tt.expectedGrid) {
				t.Errorf("Step() result mismatch")
				t.Errorf("Expected:")
				for i, row := range tt.expectedGrid {
					t.Errorf("  [%d] %v", i, row)
				}
				t.Errorf("Got:")
				for i, row := range game.currentGrid {
					t.Errorf("  [%d] %v", i, row)
				}
			}

			// Check generation increment
			if game.generation != 1 {
				t.Errorf("generation = %d, want 1", game.generation)
			}
		})
	}
}

// Test Reset function
func TestReset(t *testing.T) {
	game := NewGameOfLife(5, 5, BoundaryPeriodic, PatternGlider)

	// Advance a few generations
	for i := 0; i < 3; i++ {
		game.Step()
	}

	if game.generation != 3 {
		t.Errorf("generation before reset = %d, want 3", game.generation)
	}

	// Reset the game
	game.Reset()

	if game.generation != 0 {
		t.Errorf("generation after reset = %d, want 0", game.generation)
	}
}

// Test SetPattern function
func TestSetPattern(t *testing.T) {
	game := NewGameOfLife(10, 10, BoundaryPeriodic, PatternRandom)

	// Set a new pattern
	game.SetPattern(PatternGlider)

	if game.pattern != PatternGlider {
		t.Errorf("pattern = %v, want %v", game.pattern, PatternGlider)
	}

	if game.generation != 0 {
		t.Errorf("generation after SetPattern = %d, want 0", game.generation)
	}
}

// Test GetCurrentGrid function
func TestGetCurrentGrid(t *testing.T) {
	game := NewGameOfLife(3, 3, BoundaryPeriodic, PatternRandom)

	grid := game.GetCurrentGrid()

	if len(grid) != 3 {
		t.Errorf("grid rows = %d, want 3", len(grid))
	}

	for i, row := range grid {
		if len(row) != 3 {
			t.Errorf("grid[%d] cols = %d, want 3", i, len(row))
		}
	}
}

// Test GetGeneration function
func TestGetGeneration(t *testing.T) {
	game := NewGameOfLife(5, 5, BoundaryPeriodic, PatternRandom)

	if game.GetGeneration() != 0 {
		t.Errorf("initial generation = %d, want 0", game.GetGeneration())
	}

	game.Step()

	if game.GetGeneration() != 1 {
		t.Errorf("generation after step = %d, want 1", game.GetGeneration())
	}
}

// Test pattern initialization
func TestPatternInitialization(t *testing.T) {
	patterns := []Pattern{
		PatternRandom,
		PatternGlider,
		PatternGliderGun,
		PatternOscillator,
		PatternPulsar,
		PatternPentomino,
	}

	for _, pattern := range patterns {
		t.Run(pattern.ToString(DefaultLanguage), func(t *testing.T) {
			game := NewGameOfLife(20, 40, BoundaryPeriodic, pattern)

			// Check that some cells are alive (except for potentially empty patterns)
			hasAliveCell := false
			for i := range game.rows {
				for j := range game.cols {
					if game.currentGrid[i][j] {
						hasAliveCell = true
						break
					}
				}
				if hasAliveCell {
					break
				}
			}

			// Most patterns should have at least one alive cell
			// (Random pattern might rarely be all dead, but that's extremely unlikely)
			if !hasAliveCell && pattern != PatternRandom {
				t.Errorf("Pattern %s has no alive cells", pattern.ToString(DefaultLanguage))
			}
		})
	}
}

// Benchmark Step function
func BenchmarkStep(b *testing.B) {
	sizes := []struct {
		name string
		rows int
		cols int
	}{
		{"Small_10x10", 10, 10},
		{"Medium_50x50", 50, 50},
		{"Large_100x100", 100, 100},
		{"XLarge_200x200", 200, 200},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			game := NewGameOfLife(size.rows, size.cols, BoundaryPeriodic, PatternRandom)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				game.Step()
			}
		})
	}
}

// Benchmark countNeighbors function
func BenchmarkCountNeighbors(b *testing.B) {
	game := NewGameOfLife(100, 100, BoundaryPeriodic, PatternRandom)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Count neighbors for center cell
		game.countNeighbors(50, 50)
	}
}

// Benchmark different boundary types
func BenchmarkBoundaryTypes(b *testing.B) {
	boundaries := []struct {
		name     string
		boundary BoundaryType
	}{
		{"Periodic", BoundaryPeriodic},
		{"Fixed", BoundaryFixed},
	}

	for _, boundary := range boundaries {
		b.Run(boundary.name, func(b *testing.B) {
			game := NewGameOfLife(50, 50, boundary.boundary, PatternRandom)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				game.Step()
			}
		})
	}
}

// Benchmark pattern initialization
func BenchmarkPatternInitialization(b *testing.B) {
	patterns := []Pattern{
		PatternRandom,
		PatternGlider,
		PatternGliderGun,
		PatternOscillator,
		PatternPulsar,
		PatternPentomino,
	}

	for _, pattern := range patterns {
		b.Run(pattern.ToString(DefaultLanguage), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				game := NewGameOfLife(30, 60, BoundaryPeriodic, pattern)
				_ = game // Prevent compiler optimization
			}
		})
	}
}

// Benchmark Reset function
func BenchmarkReset(b *testing.B) {
	game := NewGameOfLife(50, 50, BoundaryPeriodic, PatternRandom)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		game.Reset()
	}
}

// Test edge cases for small grids
func TestSmallGrids(t *testing.T) {
	tests := []struct {
		name string
		rows int
		cols int
	}{
		{"1x1", 1, 1},
		{"1x2", 1, 2},
		{"2x1", 2, 1},
		{"2x2", 2, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGameOfLife(tt.rows, tt.cols, BoundaryPeriodic, PatternRandom)

			// Should not panic
			game.Step()
			game.Reset()

			// Should return correct dimensions
			grid := game.GetCurrentGrid()
			if len(grid) != tt.rows {
				t.Errorf("grid rows = %d, want %d", len(grid), tt.rows)
			}
			if len(grid[0]) != tt.cols {
				t.Errorf("grid cols = %d, want %d", len(grid[0]), tt.cols)
			}
		})
	}
}

// Test clearGrid function
func TestClearGrid(t *testing.T) {
	game := NewGameOfLife(5, 5, BoundaryPeriodic, PatternGlider)

	// Ensure some cells are alive initially
	hasAliveCell := false
	for i := range game.rows {
		for j := range game.cols {
			if game.currentGrid[i][j] {
				hasAliveCell = true
				break
			}
		}
		if hasAliveCell {
			break
		}
	}

	if !hasAliveCell {
		// Manually set some cells alive for testing
		game.currentGrid[2][2] = true
		game.currentGrid[2][3] = true
	}

	// Clear the grid
	game.clearGrid()

	// Check that all cells are dead
	for i := range game.rows {
		for j := range game.cols {
			if game.currentGrid[i][j] {
				t.Errorf("cell (%d, %d) is alive after clearGrid", i, j)
			}
		}
	}
}

// Test placePattern function
func TestPlacePattern(t *testing.T) {
	game := NewGameOfLife(10, 10, BoundaryPeriodic, PatternRandom)
	game.clearGrid()

	pattern := [][]bool{
		{true, false, true},
		{false, true, false},
		{true, false, true},
	}

	// Place pattern at (2, 2)
	game.placePattern(2, 2, pattern)

	// Check that pattern is placed correctly
	expectedCells := []struct{ row, col int }{
		{2, 2}, {2, 4}, {3, 3}, {4, 2}, {4, 4},
	}

	aliveCells := 0
	for i := range game.rows {
		for j := range game.cols {
			if game.currentGrid[i][j] {
				aliveCells++
			}
		}
	}

	if aliveCells != len(expectedCells) {
		t.Errorf("alive cells = %d, want %d", aliveCells, len(expectedCells))
	}

	// Check specific positions
	for _, cell := range expectedCells {
		if !game.currentGrid[cell.row][cell.col] {
			t.Errorf("cell (%d, %d) should be alive", cell.row, cell.col)
		}
	}
}
