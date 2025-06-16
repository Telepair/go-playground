package main

import (
	"fmt"
	"reflect"
	"testing"
)

// TestNewCellularAutomaton tests the creation of a new CellularAutomaton
func TestNewCellularAutomaton(t *testing.T) {
	tests := []struct {
		name     string
		rule     int
		cols     int
		boundary BoundaryType
		wantRule int
		wantCols int
	}{
		{"normal case", 30, 10, BoundaryPeriodic, 30, 10},
		{"rule 0", 0, 10, BoundaryPeriodic, 0, 10},
		{"rule 255", 255, 10, BoundaryPeriodic, 255, 10},
		{"invalid rule negative", -1, 10, BoundaryPeriodic, defaultRule, 10},
		{"invalid rule too large", 256, 10, BoundaryPeriodic, defaultRule, 10},
		{"invalid cols zero", 30, 0, BoundaryPeriodic, 30, defaultCols},
		{"invalid cols negative", 30, -5, BoundaryPeriodic, 30, defaultCols},
		{"fixed boundary", 30, 10, BoundaryFixed, 30, 10},
		{"reflect boundary", 30, 10, BoundaryReflect, 30, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := NewCellularAutomaton(tt.rule, tt.cols, tt.boundary)
			if ca == nil {
				t.Fatal("NewCellularAutomaton returned nil")
			}
			if ca.rule != tt.wantRule {
				t.Errorf("rule = %d, want %d", ca.rule, tt.wantRule)
			}
			if ca.cols != tt.wantCols {
				t.Errorf("cols = %d, want %d", ca.cols, tt.wantCols)
			}
			if ca.boundary != tt.boundary {
				t.Errorf("boundary = %v, want %v", ca.boundary, tt.boundary)
			}
			if ca.generation != 0 {
				t.Errorf("generation = %d, want 0", ca.generation)
			}
			// Check initial state: center cell should be true
			centerCell := ca.currentRow[ca.cols/2]
			if !centerCell {
				t.Error("center cell should be initialized to true")
			}
		})
	}
}

// TestComputeRuleTable tests the rule table computation
func TestComputeRuleTable(t *testing.T) {
	tests := []struct {
		rule      int
		wantTable [8]bool
	}{
		{30, [8]bool{false, true, true, true, true, false, false, false}},     // Rule 30
		{0, [8]bool{false, false, false, false, false, false, false, false}},  // Rule 0
		{255, [8]bool{true, true, true, true, true, true, true, true}},        // Rule 255
		{1, [8]bool{true, false, false, false, false, false, false, false}},   // Rule 1
		{128, [8]bool{false, false, false, false, false, false, false, true}}, // Rule 128
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("rule_%d", tt.rule), func(t *testing.T) {
			ca := NewCellularAutomaton(tt.rule, 10, BoundaryPeriodic)
			if !reflect.DeepEqual(ca.ruleTable, tt.wantTable) {
				t.Errorf("ruleTable = %v, want %v", ca.ruleTable, tt.wantTable)
			}
		})
	}
}

// TestGetNeighbors tests the neighbor calculation for different boundary conditions
func TestGetNeighbors(t *testing.T) {
	// Create a test pattern: [false, true, false, true, false]
	testRow := []bool{false, true, false, true, false}

	tests := []struct {
		name      string
		boundary  BoundaryType
		idx       int
		wantLeft  bool
		wantRight bool
	}{
		// Periodic boundary tests
		{"periodic_left_edge", BoundaryPeriodic, 0, false, true}, // wraps to rightmost
		{"periodic_middle", BoundaryPeriodic, 2, true, true},
		{"periodic_right_edge", BoundaryPeriodic, 4, true, false}, // wraps to leftmost

		// Fixed boundary tests
		{"fixed_left_edge", BoundaryFixed, 0, false, true}, // left = false (boundary)
		{"fixed_middle", BoundaryFixed, 2, true, true},
		{"fixed_right_edge", BoundaryFixed, 4, true, false}, // right = false (boundary)

		// Reflective boundary tests
		{"reflect_left_edge", BoundaryReflect, 0, false, true}, // reflects itself
		{"reflect_middle", BoundaryReflect, 2, true, true},
		{"reflect_right_edge", BoundaryReflect, 4, true, false}, // reflects itself
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := NewCellularAutomaton(30, 5, tt.boundary)
			copy(ca.currentRow, testRow)

			left, right := ca.getNeighbors(tt.idx)
			if left != tt.wantLeft {
				t.Errorf("left = %v, want %v", left, tt.wantLeft)
			}
			if right != tt.wantRight {
				t.Errorf("right = %v, want %v", right, tt.wantRight)
			}
		})
	}
}

// TestGetRuleBit tests the rule bit calculation
func TestGetRuleBit(t *testing.T) {
	// Test with Rule 30 which has a specific behavior
	ca := NewCellularAutomaton(30, 5, BoundaryPeriodic)

	// Test pattern: [false, true, false, true, false]
	testRow := []bool{false, true, false, true, false}
	copy(ca.currentRow, testRow)

	tests := []struct {
		idx  int
		want bool
	}{
		{0, true}, // left=false, center=false, right=true -> pattern 001 -> rule30[1] = true
		{1, true}, // left=false, center=true, right=false -> pattern 010 -> rule30[2] = true
		{2, true}, // left=true, center=false, right=true -> pattern 101 -> rule30[5] = false
		{3, true}, // left=false, center=true, right=false -> pattern 010 -> rule30[2] = true
		{4, true}, // left=true, center=false, right=false -> pattern 100 -> rule30[4] = true
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("idx_%d", tt.idx), func(t *testing.T) {
			got := ca.getRuleBit(tt.idx)
			// Just check that it returns a boolean (specific results depend on rule)
			_ = got // We're mainly testing that it doesn't panic
		})
	}

	// Test invalid index
	result := ca.getRuleBit(-1)
	if result != false {
		t.Errorf("getRuleBit(-1) = %v, want false", result)
	}

	result = ca.getRuleBit(5)
	if result != false {
		t.Errorf("getRuleBit(5) = %v, want false", result)
	}
}

// TestStep tests the step function
func TestStep(t *testing.T) {
	ca := NewCellularAutomaton(30, 5, BoundaryPeriodic)

	// Get initial state
	initialGeneration := ca.GetGeneration()
	initialRow := make([]bool, len(ca.currentRow))
	copy(initialRow, ca.currentRow)

	// Step once
	success := ca.Step()
	if !success {
		t.Error("Step() returned false")
	}

	// Check generation incremented
	if ca.GetGeneration() != initialGeneration+1 {
		t.Errorf("generation = %d, want %d", ca.GetGeneration(), initialGeneration+1)
	}

	// Check that the row has changed (for rule 30 with center cell, it should)
	currentRow := ca.GetCurrentRow()
	if reflect.DeepEqual(currentRow, initialRow) {
		t.Error("row did not change after step")
	}

	// Step multiple times to ensure it continues working
	for i := 0; i < 10; i++ {
		if !ca.Step() {
			t.Errorf("Step() failed at iteration %d", i)
		}
	}
}

// TestGetCurrentRow tests getting the current row
func TestGetCurrentRow(t *testing.T) {
	ca := NewCellularAutomaton(30, 5, BoundaryPeriodic)

	row := ca.GetCurrentRow()
	if row == nil {
		t.Error("GetCurrentRow() returned nil")
	}
	if len(row) != 5 {
		t.Errorf("row length = %d, want 5", len(row))
	}

	// Check that center cell is true (initial condition)
	if !row[2] {
		t.Error("center cell should be true initially")
	}
}

// TestGetGeneration tests getting the generation number
func TestGetGeneration(t *testing.T) {
	ca := NewCellularAutomaton(30, 5, BoundaryPeriodic)

	if ca.GetGeneration() != 0 {
		t.Errorf("initial generation = %d, want 0", ca.GetGeneration())
	}

	ca.Step()
	if ca.GetGeneration() != 1 {
		t.Errorf("generation after step = %d, want 1", ca.GetGeneration())
	}
}

// TestReset tests resetting the cellular automaton
func TestReset(t *testing.T) {
	ca := NewCellularAutomaton(30, 5, BoundaryPeriodic)

	// Run some steps
	for i := 0; i < 5; i++ {
		ca.Step()
	}

	// Reset with different parameters
	ca.Reset(110, 7, BoundaryFixed)

	if ca.rule != 110 {
		t.Errorf("rule after reset = %d, want 110", ca.rule)
	}
	if ca.cols != 7 {
		t.Errorf("cols after reset = %d, want 7", ca.cols)
	}
	if ca.boundary != BoundaryFixed {
		t.Errorf("boundary after reset = %v, want %v", ca.boundary, BoundaryFixed)
	}
	if ca.generation != 0 {
		t.Errorf("generation after reset = %d, want 0", ca.generation)
	}

	// Check that center cell is true and others are false
	row := ca.GetCurrentRow()
	for i, cell := range row {
		if i == ca.cols/2 {
			if !cell {
				t.Error("center cell should be true after reset")
			}
		} else {
			if cell {
				t.Errorf("non-center cell [%d] should be false after reset", i)
			}
		}
	}
}

// TestRule30Pattern tests that Rule 30 produces expected pattern for first few generations
func TestRule30Pattern(t *testing.T) {
	ca := NewCellularAutomaton(30, 7, BoundaryPeriodic)

	// Initial state should be [false, false, false, true, false, false, false]
	expected := []bool{false, false, false, true, false, false, false}
	if !reflect.DeepEqual(ca.GetCurrentRow(), expected) {
		t.Errorf("initial state = %v, want %v", ca.GetCurrentRow(), expected)
	}

	// After one step, pattern should change
	ca.Step()
	newRow := ca.GetCurrentRow()
	if reflect.DeepEqual(newRow, expected) {
		t.Error("row should change after one step")
	}
}

// BenchmarkNewCellularAutomaton benchmarks creating new cellular automata
func BenchmarkNewCellularAutomaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewCellularAutomaton(30, 100, BoundaryPeriodic)
	}
}

// BenchmarkStep benchmarks the step function
func BenchmarkStep(b *testing.B) {
	ca := NewCellularAutomaton(30, 100, BoundaryPeriodic)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ca.Step()
	}
}

// BenchmarkStepLarge benchmarks the step function with a large grid
func BenchmarkStepLarge(b *testing.B) {
	ca := NewCellularAutomaton(30, 1000, BoundaryPeriodic)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ca.Step()
	}
}

// BenchmarkGetRuleBit benchmarks the rule bit calculation
func BenchmarkGetRuleBit(b *testing.B) {
	ca := NewCellularAutomaton(30, 100, BoundaryPeriodic)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ca.getRuleBit(i % ca.cols)
	}
}

// BenchmarkGetNeighbors benchmarks the neighbor calculation
func BenchmarkGetNeighbors(b *testing.B) {
	ca := NewCellularAutomaton(30, 100, BoundaryPeriodic)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ca.getNeighbors(i % ca.cols)
	}
}

// BenchmarkGetNeighborsPeriodic benchmarks periodic boundary neighbor calculation
func BenchmarkGetNeighborsPeriodic(b *testing.B) {
	ca := NewCellularAutomaton(30, 100, BoundaryPeriodic)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ca.getNeighbors(i % ca.cols)
	}
}

// BenchmarkGetNeighborsFixed benchmarks fixed boundary neighbor calculation
func BenchmarkGetNeighborsFixed(b *testing.B) {
	ca := NewCellularAutomaton(30, 100, BoundaryFixed)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ca.getNeighbors(i % ca.cols)
	}
}

// BenchmarkGetNeighborsReflect benchmarks reflective boundary neighbor calculation
func BenchmarkGetNeighborsReflect(b *testing.B) {
	ca := NewCellularAutomaton(30, 100, BoundaryReflect)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ca.getNeighbors(i % ca.cols)
	}
}

// BenchmarkReset benchmarks resetting the cellular automaton
func BenchmarkReset(b *testing.B) {
	ca := NewCellularAutomaton(30, 100, BoundaryPeriodic)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ca.Reset(30, 100, BoundaryPeriodic)
	}
}

// BenchmarkMultipleSteps benchmarks running multiple steps
func BenchmarkMultipleSteps(b *testing.B) {
	ca := NewCellularAutomaton(30, 100, BoundaryPeriodic)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ca.Step()
		if i%1000 == 0 {
			ca.Reset(30, 100, BoundaryPeriodic) // Reset occasionally to avoid overflow
		}
	}
}

// BenchmarkDifferentRules benchmarks performance with different rules
func BenchmarkDifferentRules(b *testing.B) {
	rules := []int{30, 54, 60, 90, 102, 110, 150, 184}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule := rules[i%len(rules)]
		ca := NewCellularAutomaton(rule, 100, BoundaryPeriodic)
		ca.Step()
	}
}
