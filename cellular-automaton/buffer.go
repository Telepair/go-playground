package main

// GridRingBuffer efficiently manages grid history with a circular buffer
type GridRingBuffer struct {
	buffer     [][]bool
	capacity   int
	size       int
	writeIndex int
	startIndex int
	cols       int // Store column count for consistency checks
}

// NewGridRingBuffer creates a new ring buffer for grid history with validation
func NewGridRingBuffer(capacity, cols int) *GridRingBuffer {
	// Ensure minimum valid values
	if capacity < 1 {
		capacity = 1
	}
	if cols < 1 {
		cols = 1
	}

	// Pre-allocate all buffer rows to avoid allocations during runtime
	buffer := make([][]bool, capacity)
	for i := range buffer {
		buffer[i] = make([]bool, cols)
	}

	return &GridRingBuffer{
		buffer:     buffer,
		capacity:   capacity,
		cols:       cols,
		size:       0,
		writeIndex: 0,
		startIndex: 0,
	}
}

// AddRow adds a new row to the ring buffer with comprehensive bounds checking
func (grb *GridRingBuffer) AddRow(row []bool) {
	if len(row) == 0 || grb == nil {
		return
	}

	// Ensure write index is within bounds
	if grb.writeIndex >= grb.capacity {
		grb.writeIndex = 0
	}

	// Safely copy row data with proper bounds checking
	copyLen := min(len(row), grb.cols, len(grb.buffer[grb.writeIndex]))
	copy(grb.buffer[grb.writeIndex][:copyLen], row[:copyLen])

	// Fill remaining cells with false if row is shorter than expected
	for i := copyLen; i < grb.cols && i < len(grb.buffer[grb.writeIndex]); i++ {
		grb.buffer[grb.writeIndex][i] = false
	}

	// Update ring buffer indices
	grb.writeIndex = (grb.writeIndex + 1) % grb.capacity
	if grb.size < grb.capacity {
		grb.size++
	} else {
		grb.startIndex = (grb.startIndex + 1) % grb.capacity
	}
}

// GetRows returns all rows in chronological order with safe copying
func (grb *GridRingBuffer) GetRows() [][]bool {
	if grb == nil || grb.size == 0 {
		return nil
	}

	result := make([][]bool, grb.size)
	for i := 0; i < grb.size; i++ {
		idx := (grb.startIndex + i) % grb.capacity
		if idx < len(grb.buffer) && grb.buffer[idx] != nil {
			// Create defensive copy to prevent data races
			result[i] = make([]bool, len(grb.buffer[idx]))
			copy(result[i], grb.buffer[idx])
		}
	}
	return result
}

// Clear resets the ring buffer to empty state
func (grb *GridRingBuffer) Clear() {
	if grb == nil {
		return
	}
	grb.size = 0
	grb.writeIndex = 0
	grb.startIndex = 0
}
