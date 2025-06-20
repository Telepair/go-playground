package gameoflife

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telepair/go-playground/pkg/ui"
)

// TestNewConwayGameOfLife tests the creation of a new game
func TestNewConwayGameOfLife(t *testing.T) {
	config := Config{
		AliveColor: "#00FF00",
		DeadColor:  "#000000",
		AliveChar:  "█",
		DeadChar:   " ",
	}
	game := New(30, 80, config)

	assert.NotNil(t, game)
	assert.Equal(t, 30, game.rows)
	assert.Equal(t, 80, game.cols)
	assert.Equal(t, 0, game.generation)
	assert.False(t, game.paused)
	assert.NotNil(t, game.screen)
	assert.NotNil(t, game.currentGrid)
	assert.NotNil(t, game.nextGrid)
}

// TestStep tests the step function
func TestStep(t *testing.T) {
	config := Config{
		AliveColor: "#00FF00",
		DeadColor:  "#000000",
		AliveChar:  "█",
		DeadChar:   " ",
	}
	game := New(10, 10, config)

	// Test multiple steps
	for i := 1; i <= 5; i++ {
		gen, ok := game.Step()
		assert.True(t, ok)
		assert.Equal(t, i, gen)
	}

	// Test paused state
	game.paused = true
	gen, ok := game.Step()
	assert.True(t, ok)
	assert.Equal(t, 5, gen) // Should not advance when paused
}

// TestHeader tests the header generation
func TestHeader(t *testing.T) {
	config := Config{
		AliveColor: "#00FF00",
		DeadColor:  "#000000",
		AliveChar:  "█",
		DeadChar:   " ",
	}
	game := New(10, 10, config)

	// Test English header
	assert.Equal(t, HeaderEN, game.Header(ui.English))

	// Test Chinese header
	assert.Equal(t, HeaderCN, game.Header(ui.Chinese))
}

// TestStatus tests the status generation
func TestStatus(t *testing.T) {
	config := Config{
		AliveColor: "#00FF00",
		DeadColor:  "#000000",
		AliveChar:  "█",
		DeadChar:   " ",
	}
	game := New(10, 10, config)

	// Test English status
	status := game.Status(ui.English)
	assert.Len(t, status, 4)
	assert.Equal(t, "Generation", status[0].Label)
	assert.Equal(t, "0", status[0].Value)
	assert.Equal(t, "Pattern", status[1].Label)
	assert.Equal(t, "Boundary", status[2].Label)
	assert.Equal(t, "Status", status[3].Label)

	// Test Chinese status
	status = game.Status(ui.Chinese)
	assert.Len(t, status, 4)
	assert.Equal(t, "代数", status[0].Label)
	assert.Equal(t, "0", status[0].Value)
	assert.Equal(t, "模式", status[1].Label)
	assert.Equal(t, "边界", status[2].Label)
	assert.Equal(t, "状态", status[3].Label)
}

// TestHandleKeys tests the keyboard controls
func TestHandleKeys(t *testing.T) {
	config := Config{
		AliveColor: "#00FF00",
		DeadColor:  "#000000",
		AliveChar:  "█",
		DeadChar:   " ",
	}
	game := New(10, 10, config)

	// Test English controls
	controls := game.HandleKeys(ui.English)
	assert.Len(t, controls, 4)

	// Test Chinese controls
	controls = game.HandleKeys(ui.Chinese)
	assert.Len(t, controls, 4)
}

// TestHandle tests key handling
func TestHandle(t *testing.T) {
	config := Config{
		AliveColor: "#00FF00",
		DeadColor:  "#000000",
		AliveChar:  "█",
		DeadChar:   " ",
	}
	game := New(10, 10, config)

	// Test pause/resume
	handled, err := game.Handle(" ")
	assert.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, game.paused)

	handled, err = game.Handle("space")
	assert.NoError(t, err)
	assert.True(t, handled)
	assert.False(t, game.paused)

	// Test pattern cycling
	originalPattern := game.pattern
	handled, err = game.Handle("p")
	assert.NoError(t, err)
	assert.True(t, handled)
	assert.NotEqual(t, originalPattern, game.pattern)

	// Test boundary toggle
	originalBoundary := game.boundary
	handled, err = game.Handle("b")
	assert.NoError(t, err)
	assert.True(t, handled)
	assert.NotEqual(t, originalBoundary, game.boundary)

	// Test reset
	game.generation = 10
	handled, err = game.Handle("r")
	assert.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, 0, game.generation)

	// Test unhandled key
	handled, err = game.Handle("x")
	assert.NoError(t, err)
	assert.False(t, handled)
}

// TestReset tests the reset functionality
func TestReset(t *testing.T) {
	config := Config{
		AliveColor: "#00FF00",
		DeadColor:  "#000000",
		AliveChar:  "█",
		DeadChar:   " ",
	}
	game := New(10, 10, config)

	// Advance a few generations
	for i := 0; i < 5; i++ {
		game.Step()
	}

	// Reset with new dimensions
	err := game.Reset(20, 30)
	assert.NoError(t, err)
	assert.Equal(t, 20, game.rows)
	assert.Equal(t, 30, game.cols)
	assert.Equal(t, 0, game.generation)
}

// TestIsFinished tests the IsFinished method
func TestIsFinished(t *testing.T) {
	config := Config{
		AliveColor: "#00FF00",
		DeadColor:  "#000000",
		AliveChar:  "█",
		DeadChar:   " ",
	}
	game := New(10, 10, config)

	// Conway's Game of Life never finishes
	assert.False(t, game.IsFinished())

	// Even after many steps
	for i := 0; i < 100; i++ {
		game.Step()
	}
	assert.False(t, game.IsFinished())
}

// TestStop tests the Stop method
func TestStop(t *testing.T) {
	config := Config{
		AliveColor: "#00FF00",
		DeadColor:  "#000000",
		AliveChar:  "█",
		DeadChar:   " ",
	}
	game := New(10, 10, config)

	// Should not panic
	assert.NotPanics(t, func() {
		game.Stop()
	})
}

// TestCountNeighbors tests neighbor counting
func TestCountNeighbors(t *testing.T) {
	config := Config{
		AliveColor: "#00FF00",
		DeadColor:  "#000000",
		AliveChar:  "█",
		DeadChar:   " ",
	}
	game := New(5, 5, config)

	// Create a known pattern
	game.currentGrid = [][]bool{
		{true, true, false, false, false},
		{true, true, false, false, false},
		{false, false, false, false, false},
		{false, false, false, true, true},
		{false, false, false, true, true},
	}

	// Test fixed boundary first (simpler case)
	game.boundary = BoundaryFixed
	assert.Equal(t, 3, game.countNeighbors(0, 0))
	assert.Equal(t, 3, game.countNeighbors(0, 1))
	assert.Equal(t, 3, game.countNeighbors(1, 0))
	assert.Equal(t, 3, game.countNeighbors(1, 1))
	assert.Equal(t, 3, game.countNeighbors(3, 3))
	assert.Equal(t, 3, game.countNeighbors(4, 4))

	// Test periodic boundary
	game.boundary = BoundaryPeriodic
	// For interior cells, periodic shouldn't change the count
	assert.Equal(t, 3, game.countNeighbors(1, 1))

	// For edge cells, periodic boundary wraps around
	// Cell (0,0) with periodic boundary sees (4,4) which is true
	assert.Equal(t, 4, game.countNeighbors(0, 0))

	// Test empty grid
	game.currentGrid = [][]bool{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
	}
	assert.Equal(t, 0, game.countNeighbors(2, 2))
}

// Benchmark tests
func BenchmarkStep(b *testing.B) {
	config := Config{
		AliveColor: "#00FF00",
		DeadColor:  "#000000",
		AliveChar:  "█",
		DeadChar:   " ",
	}
	game := New(100, 100, config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		game.Step()
	}
}

func BenchmarkCountNeighbors(b *testing.B) {
	config := Config{
		AliveColor: "#00FF00",
		DeadColor:  "#000000",
		AliveChar:  "█",
		DeadChar:   " ",
	}
	game := New(100, 100, config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		game.countNeighbors(50, 50)
	}
}

func BenchmarkRender(b *testing.B) {
	config := Config{
		AliveColor: "#00FF00",
		DeadColor:  "#000000",
		AliveChar:  "█",
		DeadChar:   " ",
	}
	game := New(100, 100, config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		game.render()
	}
}
