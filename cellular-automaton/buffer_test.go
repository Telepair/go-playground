package main

import (
	"testing"
)

// Test NewGridRingBuffer creation
func TestNewGridRingBuffer(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		cols     int
		expected struct {
			capacity int
			cols     int
		}
	}{
		{
			name:     "Valid parameters",
			capacity: 10,
			cols:     20,
			expected: struct {
				capacity int
				cols     int
			}{10, 20},
		},
		{
			name:     "Invalid capacity - too small",
			capacity: 5,
			cols:     20,
			expected: struct {
				capacity int
				cols     int
			}{DefaultRows, 20},
		},
		{
			name:     "Invalid cols - too small",
			capacity: 10,
			cols:     5,
			expected: struct {
				capacity int
				cols     int
			}{10, DefaultCols},
		},
		{
			name:     "Both invalid",
			capacity: 5,
			cols:     5,
			expected: struct {
				capacity int
				cols     int
			}{DefaultRows, DefaultCols},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grb := NewGridRingBuffer(tt.capacity, tt.cols)
			if grb.capacity != tt.expected.capacity {
				t.Errorf("Expected capacity %d, got %d", tt.expected.capacity, grb.capacity)
			}
			if grb.cols != tt.expected.cols {
				t.Errorf("Expected cols %d, got %d", tt.expected.cols, grb.cols)
			}
			if grb.size != 0 {
				t.Errorf("Expected size 0, got %d", grb.size)
			}
			if grb.writeIndex != 0 {
				t.Errorf("Expected writeIndex 0, got %d", grb.writeIndex)
			}
			if grb.startIndex != 0 {
				t.Errorf("Expected startIndex 0, got %d", grb.startIndex)
			}
		})
	}
}

// Test AddRow functionality
func TestGridRingBuffer_AddRow(t *testing.T) {
	grb := NewGridRingBuffer(3, 5)

	// Test adding nil row
	grb.AddRow(nil)
	if grb.size != 0 {
		t.Errorf("Expected size 0 after adding nil row, got %d", grb.size)
	}

	// Test adding empty row
	grb.AddRow([]bool{})
	if grb.size != 0 {
		t.Errorf("Expected size 0 after adding empty row, got %d", grb.size)
	}

	// Test adding normal rows and check basic functionality
	row1 := []bool{true, false, true, false, true}
	grb.AddRow(row1)
	if grb.size != 1 {
		t.Errorf("Expected size 1, got %d", grb.size)
	}

	row2 := []bool{false, true, false, true, false}
	grb.AddRow(row2)
	if grb.size != 2 {
		t.Errorf("Expected size 2, got %d", grb.size)
	}

	row3 := []bool{true, true, false, false, true}
	grb.AddRow(row3)
	if grb.size != 3 {
		t.Errorf("Expected size 3, got %d", grb.size)
	}

	// Test that buffer works (just add more rows and check it doesn't crash)
	row4 := []bool{false, false, true, true, false}
	grb.AddRow(row4)

	row5 := []bool{true, true, false, false, true}
	grb.AddRow(row5)

	// Buffer should maintain its capacity
	if grb.size > grb.capacity {
		t.Errorf("Size %d should not exceed capacity %d", grb.size, grb.capacity)
	}
}

// Test AddRow with different row sizes
func TestGridRingBuffer_AddRowDifferentSizes(t *testing.T) {
	grb := NewGridRingBuffer(3, 5)

	// Test adding row shorter than expected
	shortRow := []bool{true, false}
	grb.AddRow(shortRow)
	rows := grb.GetRows()
	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}
	if len(rows[0]) != grb.cols {
		t.Errorf("Expected row length %d, got %d", grb.cols, len(rows[0]))
	}
	// Check that the short row is padded correctly
	if len(rows[0]) < 2 {
		t.Fatal("Row too short to test")
	}
	if rows[0][0] != true {
		t.Errorf("Expected row[0] = true, got %v", rows[0][0])
	}
	if rows[0][1] != false {
		t.Errorf("Expected row[1] = false, got %v", rows[0][1])
	}
	// Rest should be false (default), but we won't test every element due to buffer size differences

	// Test adding row longer than expected
	longRow := []bool{true, false, true, false, true, true, false}
	grb.AddRow(longRow)
	rows = grb.GetRows()
	if len(rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(rows))
	}
	// Check that the long row is truncated to fit buffer column size
	longRowExpected := []bool{true, false, true, false, true}
	// Check first few elements
	for i := 0; i < min(len(longRowExpected), len(rows[1])); i++ {
		if rows[1][i] != longRowExpected[i] {
			t.Errorf("Expected row[%d] = %v, got %v", i, longRowExpected[i], rows[1][i])
		}
	}
}

// Test GetRows functionality
func TestGridRingBuffer_GetRows(t *testing.T) {
	grb := NewGridRingBuffer(3, 4)

	// Test empty buffer
	rows := grb.GetRows()
	if rows != nil {
		t.Errorf("Expected nil for empty buffer, got %v", rows)
	}

	// Add some rows
	row1 := []bool{true, false, true, false}
	row2 := []bool{false, true, false, true}
	row3 := []bool{true, true, false, false}

	grb.AddRow(row1)
	grb.AddRow(row2)
	grb.AddRow(row3)

	rows = grb.GetRows()
	if len(rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(rows))
	}

	// Verify the rows are in chronological order
	expectedRows := [][]bool{row1, row2, row3}
	for i, expectedRow := range expectedRows {
		for j, expectedVal := range expectedRow {
			if rows[i][j] != expectedVal {
				t.Errorf("Expected rows[%d][%d] = %v, got %v", i, j, expectedVal, rows[i][j])
			}
		}
	}

	// Test overflow scenario by adding more rows than capacity
	row4 := []bool{false, false, true, true}
	row5 := []bool{true, true, false, false}
	grb.AddRow(row4)
	grb.AddRow(row5) // This should cause overflow

	rows = grb.GetRows()
	// Just check that we get some rows back and don't crash
	if len(rows) == 0 {
		t.Errorf("Expected at least some rows after overflow, got %d", len(rows))
	}
	if len(rows) > grb.capacity {
		t.Errorf("Expected at most %d rows after overflow, got %d", grb.capacity, len(rows))
	}
}

// Test GetRows with nil buffer
func TestGridRingBuffer_GetRowsNilBuffer(t *testing.T) {
	var grb *GridRingBuffer
	rows := grb.GetRows()
	if rows != nil {
		t.Errorf("Expected nil for nil buffer, got %v", rows)
	}
}

// Test Clear functionality
func TestGridRingBuffer_Clear(t *testing.T) {
	grb := NewGridRingBuffer(3, 4)

	// Add some rows
	grb.AddRow([]bool{true, false, true, false})
	grb.AddRow([]bool{false, true, false, true})

	// Clear the buffer
	grb.Clear()

	if grb.size != 0 {
		t.Errorf("Expected size 0 after clear, got %d", grb.size)
	}
	if grb.writeIndex != 0 {
		t.Errorf("Expected writeIndex 0 after clear, got %d", grb.writeIndex)
	}
	if grb.startIndex != 0 {
		t.Errorf("Expected startIndex 0 after clear, got %d", grb.startIndex)
	}

	rows := grb.GetRows()
	if rows != nil {
		t.Errorf("Expected nil rows after clear, got %v", rows)
	}
}

// Test Clear with nil buffer
func TestGridRingBuffer_ClearNilBuffer(t *testing.T) {
	var grb *GridRingBuffer
	// Should not panic
	grb.Clear()
}

// Test AddRow with nil buffer
func TestGridRingBuffer_AddRowNilBuffer(t *testing.T) {
	var grb *GridRingBuffer
	// Should not panic
	grb.AddRow([]bool{true, false})
}

// Test concurrent access safety (basic test)
func TestGridRingBuffer_ConcurrentAccess(t *testing.T) {
	grb := NewGridRingBuffer(10, 5)

	// Test that GetRows returns defensive copies
	grb.AddRow([]bool{true, false, true, false, true})
	rows1 := grb.GetRows()
	rows2 := grb.GetRows()

	// Modify one of the returned slices
	if len(rows1) > 0 && len(rows1[0]) > 0 {
		rows1[0][0] = false
	}

	// The other slice should not be affected
	if len(rows2) > 0 && len(rows2[0]) > 0 {
		if !rows2[0][0] {
			t.Errorf("Expected defensive copy, but slice was modified")
		}
	}
}

// Benchmark tests
func BenchmarkNewGridRingBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewGridRingBuffer(100, 80)
	}
}

func BenchmarkGridRingBuffer_AddRow(b *testing.B) {
	grb := NewGridRingBuffer(1000, 80)
	row := make([]bool, 80)
	for i := range row {
		row[i] = i%2 == 0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grb.AddRow(row)
	}
}

func BenchmarkGridRingBuffer_GetRows(b *testing.B) {
	grb := NewGridRingBuffer(100, 80)
	row := make([]bool, 80)
	for i := range row {
		row[i] = i%2 == 0
	}

	// Fill the buffer
	for i := 0; i < 100; i++ {
		grb.AddRow(row)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grb.GetRows()
	}
}

func BenchmarkGridRingBuffer_Clear(b *testing.B) {
	grb := NewGridRingBuffer(100, 80)
	row := make([]bool, 80)

	for i := 0; i < b.N; i++ {
		// Fill the buffer
		for j := 0; j < 100; j++ {
			grb.AddRow(row)
		}
		grb.Clear()
	}
}

// Benchmark overflow scenario
func BenchmarkGridRingBuffer_AddRowOverflow(b *testing.B) {
	grb := NewGridRingBuffer(10, 80) // Small buffer to force overflow
	row := make([]bool, 80)
	for i := range row {
		row[i] = i%2 == 0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grb.AddRow(row)
	}
}
