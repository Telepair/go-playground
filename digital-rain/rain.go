//nolint:gosec // This file uses math/rand for visual effects, not security
package main

import (
	"math/rand"
	"sync"
)

// Drop represents a single falling character column
type Drop struct {
	X        int    // Column position
	Y        int    // Current head position
	Speed    int    // Fall speed (cells per update)
	Length   int    // Length of the trail
	Chars    []rune // Characters in this drop
	NextMove int    // Counter for next move
	Active   bool   // Whether this drop is active
}

// DigitalRain manages the digital rain effect
type DigitalRain struct {
	mu       sync.RWMutex
	width    int
	height   int
	drops    []*Drop
	grid     [][]rune
	trail    [][]int // Trail intensity (0-255)
	charSet  []rune
	minSpeed int
	maxSpeed int
	dropLen  int
}

// NewDigitalRain creates a new digital rain instance
func NewDigitalRain(width, height int, charSet string, minSpeed, maxSpeed, dropLen int) *DigitalRain {
	dr := &DigitalRain{
		width:    width,
		height:   height,
		charSet:  []rune(charSet),
		minSpeed: minSpeed,
		maxSpeed: maxSpeed,
		dropLen:  dropLen,
	}
	dr.Reset(width, height)
	return dr
}

// Reset reinitializes the digital rain with new dimensions
func (dr *DigitalRain) Reset(width, height int) {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	dr.width = width
	dr.height = height

	// Initialize grid
	dr.grid = make([][]rune, height)
	dr.trail = make([][]int, height)
	for i := 0; i < height; i++ {
		dr.grid[i] = make([]rune, width)
		dr.trail[i] = make([]int, width)
	}

	// Initialize drops (one per column)
	dr.drops = make([]*Drop, width)
	for i := 0; i < width; i++ {
		dr.drops[i] = dr.createNewDrop(i)
	}
}

// createNewDrop creates a new drop at the given column
func (dr *DigitalRain) createNewDrop(x int) *Drop {
	length := dr.dropLen + rand.Intn(5) - 2 // Add some variation
	if length < 3 {
		length = 3
	}

	drop := &Drop{
		X:        x,
		Y:        -length, // Start above the screen
		Speed:    dr.minSpeed + rand.Intn(dr.maxSpeed-dr.minSpeed+1),
		Length:   length,
		Chars:    make([]rune, length),
		NextMove: 0,
		Active:   true,
	}

	// Fill with random characters
	for i := 0; i < length; i++ {
		drop.Chars[i] = dr.charSet[rand.Intn(len(dr.charSet))]
	}

	return drop
}

// Step advances the animation by one frame
func (dr *DigitalRain) Step() {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	// Clear grid
	for i := 0; i < dr.height; i++ {
		for j := 0; j < dr.width; j++ {
			dr.grid[i][j] = 0
			// Fade trail
			if dr.trail[i][j] > 0 {
				dr.trail[i][j] -= 20
				if dr.trail[i][j] < 0 {
					dr.trail[i][j] = 0
				}
			}
		}
	}

	// Update drops
	for i, drop := range dr.drops {
		if !drop.Active {
			// Randomly restart inactive drops
			if rand.Float32() < 0.01 {
				dr.drops[i] = dr.createNewDrop(i)
			}
			continue
		}

		// Move drop
		drop.NextMove++
		if drop.NextMove >= drop.Speed {
			drop.NextMove = 0
			drop.Y++
		}

		// Randomly change characters
		if rand.Float32() < 0.1 {
			idx := rand.Intn(drop.Length)
			drop.Chars[idx] = dr.charSet[rand.Intn(len(dr.charSet))]
		}

		// Draw drop
		for j := 0; j < drop.Length; j++ {
			y := drop.Y - j
			if y >= 0 && y < dr.height {
				dr.grid[y][drop.X] = drop.Chars[j]
				// Set trail intensity (brighter at head)
				intensity := 255 - (j * 255 / drop.Length)
				if intensity > dr.trail[y][drop.X] {
					dr.trail[y][drop.X] = intensity
				}
			}
		}

		// Check if drop is off screen
		if drop.Y-drop.Length >= dr.height {
			drop.Active = false
		}
	}
}

// GetGrid returns the current character grid
func (dr *DigitalRain) GetGrid() [][]rune {
	dr.mu.RLock()
	defer dr.mu.RUnlock()
	return dr.grid
}

// GetTrail returns the current trail intensity grid
func (dr *DigitalRain) GetTrail() [][]int {
	dr.mu.RLock()
	defer dr.mu.RUnlock()
	return dr.trail
}

// GetDimensions returns the current width and height
func (dr *DigitalRain) GetDimensions() (int, int) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()
	return dr.width, dr.height
}
