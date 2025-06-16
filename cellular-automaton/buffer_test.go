package main

import (
	"testing"
)

// TestNewGridRingBuffer tests the creation of a new GridRingBuffer
func TestNewGridRingBuffer(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		cols     int
		wantCap  int
		wantCols int
	}{
		{"normal case", 5, 10, 5, 10},
		{"zero capacity", 0, 10, 1, 10},
		{"negative capacity", -5, 10, 1, 10},
		{"zero cols", 5, 0, 5, 1},
		{"negative cols", 5, -10, 5, 1},
		{"all zero", 0, 0, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grb := NewGridRingBuffer(tt.capacity, tt.cols)
			if grb == nil {
				t.Fatal("NewGridRingBuffer returned nil")
			}
			if grb.capacity != tt.wantCap {
				t.Errorf("capacity = %d, want %d", grb.capacity, tt.wantCap)
			}
			if grb.cols != tt.wantCols {
				t.Errorf("cols = %d, want %d", grb.cols, tt.wantCols)
			}
			if grb.size != 0 {
				t.Errorf("size = %d, want 0", grb.size)
			}
			if grb.writeIndex != 0 {
				t.Errorf("writeIndex = %d, want 0", grb.writeIndex)
			}
			if grb.startIndex != 0 {
				t.Errorf("startIndex = %d, want 0", grb.startIndex)
			}
		})
	}
}

// TestAddRow tests adding rows to the ring buffer
func TestAddRow(t *testing.T) {
	grb := NewGridRingBuffer(3, 5)

	// Test adding normal rows
	row1 := []bool{true, false, true, false, true}
	grb.AddRow(row1)
	if grb.size != 1 {
		t.Errorf("size after first add = %d, want 1", grb.size)
	}

	row2 := []bool{false, true, false, true, false}
	grb.AddRow(row2)
	if grb.size != 2 {
		t.Errorf("size after second add = %d, want 2", grb.size)
	}

	// Test adding empty row
	grb.AddRow([]bool{})
	if grb.size != 2 {
		t.Errorf("size after empty row add = %d, want 2", grb.size)
	}

	// Test adding nil to nil buffer
	var nilGrb *GridRingBuffer
	nilGrb.AddRow(row1) // Should not panic

	// Test adding row with different length
	shortRow := []bool{true, false}
	grb.AddRow(shortRow)
	if grb.size != 3 {
		t.Errorf("size after short row add = %d, want 3", grb.size)
	}

	longRow := []bool{true, false, true, false, true, true, false}
	grb.AddRow(longRow)
	if grb.size != 3 { // Should still be 3 (capacity reached)
		t.Errorf("size after long row add = %d, want 3", grb.size)
	}
}

// TestAddRowCapacityOverflow tests ring buffer behavior when capacity is exceeded
func TestAddRowCapacityOverflow(t *testing.T) {
	grb := NewGridRingBuffer(2, 3)

	// Add first row
	row1 := []bool{true, false, true}
	grb.AddRow(row1)

	// Add second row
	row2 := []bool{false, true, false}
	grb.AddRow(row2)

	// Add third row (should overwrite first)
	row3 := []bool{true, true, false}
	grb.AddRow(row3)

	if grb.size != 2 {
		t.Errorf("size = %d, want 2", grb.size)
	}

	rows := grb.GetRows()
	if len(rows) != 2 {
		t.Errorf("GetRows returned %d rows, want 2", len(rows))
	}

	// Check that first row was overwritten and we have row2 and row3
	expected := [][]bool{
		{false, true, false}, // row2
		{true, true, false},  // row3
	}

	for i, row := range rows {
		for j, cell := range row {
			if cell != expected[i][j] {
				t.Errorf("row[%d][%d] = %v, want %v", i, j, cell, expected[i][j])
			}
		}
	}
}

// TestGetRows tests retrieving rows from the ring buffer
func TestGetRows(t *testing.T) {
	grb := NewGridRingBuffer(3, 4)

	// Test empty buffer
	rows := grb.GetRows()
	if rows != nil {
		t.Errorf("GetRows on empty buffer = %v, want nil", rows)
	}

	// Test nil buffer
	var nilGrb *GridRingBuffer
	rows = nilGrb.GetRows()
	if rows != nil {
		t.Errorf("GetRows on nil buffer = %v, want nil", rows)
	}

	// Add some rows
	testRows := [][]bool{
		{true, false, true, false},
		{false, true, false, true},
		{true, true, false, false},
	}

	for _, row := range testRows {
		grb.AddRow(row)
	}

	rows = grb.GetRows()
	if len(rows) != 3 {
		t.Errorf("GetRows returned %d rows, want 3", len(rows))
	}

	// Verify order and content
	for i, row := range rows {
		for j, cell := range row {
			if cell != testRows[i][j] {
				t.Errorf("row[%d][%d] = %v, want %v", i, j, cell, testRows[i][j])
			}
		}
	}

	// Test that returned rows are copies (defensive copying)
	rows[0][0] = !rows[0][0]
	newRows := grb.GetRows()
	if newRows[0][0] == rows[0][0] {
		t.Error("GetRows did not return defensive copies")
	}
}

// TestClear tests clearing the ring buffer
func TestClear(t *testing.T) {
	grb := NewGridRingBuffer(3, 4)

	// Add some rows
	grb.AddRow([]bool{true, false, true, false})
	grb.AddRow([]bool{false, true, false, true})

	// Clear buffer
	grb.Clear()

	if grb.size != 0 {
		t.Errorf("size after clear = %d, want 0", grb.size)
	}
	if grb.writeIndex != 0 {
		t.Errorf("writeIndex after clear = %d, want 0", grb.writeIndex)
	}
	if grb.startIndex != 0 {
		t.Errorf("startIndex after clear = %d, want 0", grb.startIndex)
	}

	rows := grb.GetRows()
	if rows != nil {
		t.Errorf("GetRows after clear = %v, want nil", rows)
	}

	// Test clearing nil buffer
	var nilGrb *GridRingBuffer
	nilGrb.Clear() // Should not panic
}

// TestGridRingBufferEdgeCases tests various edge cases
func TestGridRingBufferEdgeCases(t *testing.T) {
	// Test with capacity 1
	grb := NewGridRingBuffer(1, 3)
	grb.AddRow([]bool{true, false, true})
	grb.AddRow([]bool{false, true, false})

	rows := grb.GetRows()
	if len(rows) != 1 {
		t.Errorf("capacity 1 buffer has %d rows, want 1", len(rows))
	}

	expected := []bool{false, true, false}
	for i, cell := range rows[0] {
		if cell != expected[i] {
			t.Errorf("row[0][%d] = %v, want %v", i, cell, expected[i])
		}
	}

	// Test with cols 1
	grb2 := NewGridRingBuffer(3, 1)
	grb2.AddRow([]bool{true})
	grb2.AddRow([]bool{false})

	rows2 := grb2.GetRows()
	if len(rows2) != 2 {
		t.Errorf("cols 1 buffer has %d rows, want 2", len(rows2))
	}
	if len(rows2[0]) != 1 || len(rows2[1]) != 1 {
		t.Error("cols 1 buffer rows have wrong length")
	}
}

// BenchmarkNewGridRingBuffer benchmarks creating new ring buffers
func BenchmarkNewGridRingBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewGridRingBuffer(100, 80)
	}
}

// BenchmarkAddRow benchmarks adding rows to the ring buffer
func BenchmarkAddRow(b *testing.B) {
	grb := NewGridRingBuffer(1000, 80)
	row := make([]bool, 80)
	for i := 0; i < 80; i++ {
		row[i] = i%2 == 0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grb.AddRow(row)
	}
}

// BenchmarkAddRowLarge benchmarks adding rows to a large ring buffer
func BenchmarkAddRowLarge(b *testing.B) {
	grb := NewGridRingBuffer(10000, 800)
	row := make([]bool, 800)
	for i := 0; i < 800; i++ {
		row[i] = i%3 == 0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grb.AddRow(row)
	}
}

// BenchmarkGetRows benchmarks retrieving all rows from the buffer
func BenchmarkGetRows(b *testing.B) {
	grb := NewGridRingBuffer(100, 80)
	row := make([]bool, 80)
	for i := 0; i < 80; i++ {
		row[i] = i%2 == 0
	}

	// Fill the buffer
	for i := 0; i < 100; i++ {
		grb.AddRow(row)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = grb.GetRows()
	}
}

// BenchmarkGetRowsLarge benchmarks retrieving all rows from a large buffer
func BenchmarkGetRowsLarge(b *testing.B) {
	grb := NewGridRingBuffer(1000, 800)
	row := make([]bool, 800)
	for i := 0; i < 800; i++ {
		row[i] = i%3 == 0
	}

	// Fill the buffer
	for i := 0; i < 1000; i++ {
		grb.AddRow(row)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = grb.GetRows()
	}
}

// BenchmarkClear benchmarks clearing the ring buffer
func BenchmarkClear(b *testing.B) {
	grb := NewGridRingBuffer(100, 80)
	row := make([]bool, 80)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Fill some data
		for j := 0; j < 10; j++ {
			grb.AddRow(row)
		}
		grb.Clear()
	}
}

// BenchmarkRingBufferCycling benchmarks the performance when the buffer cycles
func BenchmarkRingBufferCycling(b *testing.B) {
	grb := NewGridRingBuffer(50, 80) // Small capacity to force cycling
	row := make([]bool, 80)
	for i := 0; i < 80; i++ {
		row[i] = i%2 == 0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grb.AddRow(row)
		if i%100 == 0 {
			_ = grb.GetRows() // Occasional read to simulate real usage
		}
	}
}
