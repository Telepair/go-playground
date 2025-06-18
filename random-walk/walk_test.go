package main

import (
	"testing"
)

func TestNewRandomWalk(t *testing.T) {
	rows, cols := 20, 30
	walkerCount := 3
	trailLength := 50

	tests := []struct {
		name string
		mode WalkMode
	}{
		{"Single Walker", ModeSingleWalker},
		{"Multi Walker", ModeMultiWalker},
		{"Trail Mode", ModeTrailMode},
		{"Brownian Motion", ModeBrownianMotion},
		{"Self-Avoiding Walk", ModeSelfAvoidingWalk},
		{"Lévy Flight", ModeLevyFlight},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := NewRandomWalk(rows, cols, tt.mode, walkerCount, trailLength)

			if rw == nil {
				t.Fatal("NewRandomWalk returned nil")
			}

			if rw.rows != rows {
				t.Errorf("Expected rows=%d, got %d", rows, rw.rows)
			}

			if rw.cols != cols {
				t.Errorf("Expected cols=%d, got %d", cols, rw.cols)
			}

			if rw.mode != tt.mode {
				t.Errorf("Expected mode=%v, got %v", tt.mode, rw.mode)
			}

			// Check walker count based on mode
			expectedWalkers := 1
			if tt.mode == ModeMultiWalker || tt.mode == ModeBrownianMotion {
				expectedWalkers = walkerCount
			}

			if len(rw.walkers) != expectedWalkers {
				t.Errorf("Expected %d walkers, got %d", expectedWalkers, len(rw.walkers))
			}
		})
	}
}

func TestRandomWalkStep(t *testing.T) {
	rows, cols := 10, 10
	rw := NewRandomWalk(rows, cols, ModeSingleWalker, 1, 50)

	// Get initial position
	initialWalker := rw.walkers[0]
	initialPos := initialWalker.Position

	// Step multiple times
	moved := false
	stepCount := 0
	for i := 0; i < 10; i++ {
		rw.Step()
		stepCount++
		currentPos := rw.walkers[0].Position
		if currentPos != initialPos {
			moved = true
			break
		}
	}

	if !moved {
		t.Error("Walker did not move after 10 steps")
	}

	if rw.GetSteps() != stepCount {
		t.Errorf("Expected %d steps, got %d", stepCount, rw.GetSteps())
	}
}

func TestSelfAvoidingWalk(t *testing.T) {
	rows, cols := 5, 5
	rw := NewRandomWalk(rows, cols, ModeSelfAvoidingWalk, 1, 50)

	// Run many steps
	for i := 0; i < 100; i++ {
		rw.Step()
	}

	// Check that walker visited positions are tracked
	walker := rw.walkers[0]
	if len(walker.Visited) == 0 {
		t.Error("Self-avoiding walk should track visited positions")
	}

	// Verify current position is marked as visited
	if !walker.Visited[walker.Position] {
		t.Error("Current position should be marked as visited")
	}
}

func TestMultiWalker(t *testing.T) {
	rows, cols := 20, 20
	walkerCount := 5
	rw := NewRandomWalk(rows, cols, ModeMultiWalker, walkerCount, 50)

	if len(rw.walkers) != walkerCount {
		t.Errorf("Expected %d walkers, got %d", walkerCount, len(rw.walkers))
	}

	// Check that each walker has a unique ID and color
	seenIDs := make(map[int]bool)
	seenColors := make(map[string]bool)

	for _, walker := range rw.walkers {
		if seenIDs[walker.ID] {
			t.Errorf("Duplicate walker ID: %d", walker.ID)
		}
		seenIDs[walker.ID] = true

		if seenColors[walker.Color] {
			t.Errorf("Duplicate walker color: %s", walker.Color)
		}
		seenColors[walker.Color] = true
	}
}

func TestTrailMode(t *testing.T) {
	rows, cols := 10, 10
	trailLength := 5
	rw := NewRandomWalk(rows, cols, ModeTrailMode, 1, trailLength)

	// Step multiple times
	for i := 0; i < trailLength+2; i++ {
		rw.Step()
	}

	walker := rw.walkers[0]
	if len(walker.Trail) == 0 {
		t.Error("Trail mode should maintain walker trail")
	}

	if len(walker.Trail) > trailLength {
		t.Errorf("Trail length should not exceed %d, got %d", trailLength, len(walker.Trail))
	}
}

func TestReset(t *testing.T) {
	rows, cols := 10, 10
	rw := NewRandomWalk(rows, cols, ModeSingleWalker, 1, 50)

	// Step a few times
	for i := 0; i < 5; i++ {
		rw.Step()
	}

	if rw.GetSteps() != 5 {
		t.Errorf("Expected 5 steps before reset, got %d", rw.GetSteps())
	}

	// Reset with different parameters
	newRows, newCols := 15, 15
	rw.Reset(newRows, newCols, ModeMultiWalker, 3, 100)

	if rw.GetSteps() != 0 {
		t.Errorf("Steps should be 0 after reset, got %d", rw.GetSteps())
	}

	if rw.rows != newRows || rw.cols != newCols {
		t.Errorf("Expected dimensions %dx%d, got %dx%d", newRows, newCols, rw.rows, rw.cols)
	}

	if len(rw.walkers) != 3 {
		t.Errorf("Expected 3 walkers after reset, got %d", len(rw.walkers))
	}
}

func BenchmarkRandomWalkStep(b *testing.B) {
	rows, cols := 100, 100
	rw := NewRandomWalk(rows, cols, ModeSingleWalker, 1, 50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rw.Step()
	}
}

func BenchmarkMultiWalkerStep(b *testing.B) {
	rows, cols := 100, 100
	rw := NewRandomWalk(rows, cols, ModeMultiWalker, 10, 50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rw.Step()
	}
}
