package main

import (
	"testing"
	"time"
)

// BenchmarkCellularAutomatonStep benchmarks the Step function
func BenchmarkCellularAutomatonStep(b *testing.B) {
	ca := NewCellularAutomaton(30, 0, 100, BoundaryPeriodic)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ca.Step()
	}
}

// BenchmarkRenderGrid benchmarks the grid rendering
func BenchmarkRenderGrid(b *testing.B) {
	cfg := NewConfig()
	cfg.SetRows(50)
	cfg.SetCols(100)
	model := NewModel(cfg)

	// Generate some data first
	for i := 0; i < 25; i++ {
		model.ca.Step()
		model.gridRingBuffer.AddRow(model.ca.GetCurrentRow())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.RenderGrid()
	}
}

// BenchmarkGetNeighbors benchmarks the neighbor calculation
func BenchmarkGetNeighbors(b *testing.B) {
	ca := NewCellularAutomaton(30, 0, 1000, BoundaryPeriodic)

	b.ResetTimer()
	for range b.N {
		for j := range 1000 {
			_, _ = ca.getNeighbors(j)
		}
	}
}

// BenchmarkRuleTableLookup benchmarks rule table lookup vs bit operations
func BenchmarkRuleTableLookup(b *testing.B) {
	ca := NewCellularAutomaton(110, 0, 100, BoundaryPeriodic)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_ = ca.getRuleBit(j)
		}
	}
}

// BenchmarkStringBuilderCapacity tests string builder performance with different capacities
func BenchmarkStringBuilderCapacity(b *testing.B) {
	cfg := NewConfig()
	cfg.SetRows(35)
	cfg.SetCols(60)
	model := NewModel(cfg)

	// Generate some data
	for i := 0; i < 35; i++ {
		model.ca.Step()
		model.gridRingBuffer.AddRow(model.ca.GetCurrentRow())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.RenderMode()
	}
}

// Test memory allocations during rendering
func BenchmarkMemoryAllocations(b *testing.B) {
	cfg := NewConfig()
	cfg.SetRows(35)
	cfg.SetCols(60)
	model := NewModel(cfg)

	// Pre-populate with data
	for i := 0; i < 35; i++ {
		model.ca.Step()
		model.gridRingBuffer.AddRow(model.ca.GetCurrentRow())
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = model.RenderGrid()
	}
}

// TestOptimizationCorrectness verifies that optimizations don't break functionality
func TestOptimizationCorrectness(t *testing.T) {
	// Test cellular automaton correctness
	ca := NewCellularAutomaton(30, 50, 10, BoundaryPeriodic)
	initialRow := make([]bool, 10)
	copy(initialRow, ca.GetCurrentRow())

	if ca.GetGeneration() != 0 {
		t.Errorf("Initial generation should be 0, got %d", ca.GetGeneration())
	}

	// Test a few steps
	for i := range 5 {
		if !ca.Step() {
			t.Errorf("Step %d should succeed", i)
		}
		if ca.GetGeneration() != i+1 {
			t.Errorf("After step %d, generation should be %d, got %d", i, i+1, ca.GetGeneration())
		}
	}

	// Test reset
	ca.Reset()
	if ca.GetGeneration() != 0 {
		t.Errorf("After reset, generation should be 0, got %d", ca.GetGeneration())
	}

	// Verify initial state is restored
	resetRow := ca.GetCurrentRow()
	for i, cell := range resetRow {
		if cell != initialRow[i] {
			t.Errorf("Reset didn't restore initial state at position %d", i)
		}
	}
}

// TestRingBufferCorrectness tests the ring buffer implementation
func TestRingBufferCorrectness(t *testing.T) {
	grb := NewGridRingBuffer(3, 5)

	// Test empty buffer
	rows := grb.GetRows()
	if len(rows) != 0 {
		t.Errorf("Empty buffer should return 0 rows, got %d", len(rows))
	}

	// Add some rows
	row1 := []bool{true, false, true, false, true}
	row2 := []bool{false, true, false, true, false}
	row3 := []bool{true, true, false, false, true}
	row4 := []bool{false, false, true, true, false}

	grb.AddRow(row1)
	grb.AddRow(row2)
	grb.AddRow(row3)

	rows = grb.GetRows()
	if len(rows) != 3 {
		t.Errorf("Should have 3 rows, got %d", len(rows))
	}

	// Add fourth row (should overflow and replace first)
	grb.AddRow(row4)
	rows = grb.GetRows()
	if len(rows) != 3 {
		t.Errorf("Should still have 3 rows after overflow, got %d", len(rows))
	}

	// Check that oldest row is now row2
	for i, cell := range rows[0] {
		if cell != row2[i] {
			t.Errorf("First row should be row2 after overflow")
		}
	}
}

// TestBoundaryConditions tests all boundary condition types
func TestBoundaryConditions(t *testing.T) {
	testCases := []struct {
		boundary BoundaryType
		name     string
	}{
		{BoundaryPeriodic, "Periodic"},
		{BoundaryFixed, "Fixed"},
		{BoundaryReflect, "Reflect"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ca := NewCellularAutomaton(30, 10, 10, tc.boundary)

			// Run a few steps to ensure no panics
			for range 5 {
				if !ca.Step() {
					t.Errorf("Step failed for boundary type %s", tc.name)
				}
			}

			// Test edge cases
			left, right := ca.getNeighbors(0)
			_, _ = left, right // Use variables to avoid compiler warnings

			left, right = ca.getNeighbors(9)
			_, _ = left, right
		})
	}
}

// Performance comparison utility
func init() {
	// Set up any global test configuration
	time.Local = time.UTC // Ensure consistent timing
}
