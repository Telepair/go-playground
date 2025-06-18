package main

import (
	"log/slog"
	"math"
	"math/rand/v2"
	"time"
)

// Position represents a 2D position
type Position struct {
	X, Y int
}

// Walker represents a single walker
type Walker struct {
	ID       int
	Position Position
	Trail    []Position
	Color    string
	Visited  map[Position]bool // For self-avoiding walk
}

// RandomWalk represents the random walk simulation
type RandomWalk struct {
	grid        [][]int // Grid to store walker IDs (0 = empty, >0 = walker ID)
	trails      [][]int // Grid to store trail intensities
	walkers     []*Walker
	rows        int
	cols        int
	steps       int
	mode        WalkMode
	trailLength int
	rng         *rand.Rand
}

// NewRandomWalk creates a new random walk instance
func NewRandomWalk(rows, cols int, mode WalkMode, walkerCount int, trailLength int) *RandomWalk {
	slog.Debug("NewRandomWalk", "rows", rows, "cols", cols, "mode", mode, "walkerCount", walkerCount, "trailLength", trailLength)

	// Use time-based seeding for randomization
	// #nosec G115 - Conversion is safe for our use case
	seed := uint64(time.Now().UnixNano())
	// #nosec G404 - Using math/rand for simulation, not cryptography
	rng := rand.New(rand.NewPCG(seed, seed))

	rw := &RandomWalk{
		rows:        rows,
		cols:        cols,
		mode:        mode,
		trailLength: trailLength,
		steps:       0,
		rng:         rng,
	}
	rw.Init(walkerCount)
	return rw
}

// Init initializes the random walk
func (rw *RandomWalk) Init(walkerCount int) {
	slog.Debug("RandomWalk Init", "rows", rw.rows, "cols", rw.cols, "mode", rw.mode, "walkerCount", walkerCount)

	// Initialize grids
	rw.grid = make([][]int, rw.rows)
	rw.trails = make([][]int, rw.rows)
	for i := range rw.rows {
		rw.grid[i] = make([]int, rw.cols)
		rw.trails[i] = make([]int, rw.cols)
	}

	// Initialize walkers based on mode
	rw.walkers = make([]*Walker, 0)

	switch rw.mode {
	case ModeSingleWalker, ModeTrailMode, ModeSelfAvoidingWalk, ModeLevyFlight:
		// Single walker starting at center
		walker := &Walker{
			ID:       1,
			Position: Position{X: rw.cols / 2, Y: rw.rows / 2},
			Trail:    make([]Position, 0, rw.trailLength),
			Color:    DefaultWalkerColor,
			Visited:  make(map[Position]bool),
		}
		walker.Visited[walker.Position] = true
		rw.walkers = append(rw.walkers, walker)
		rw.grid[walker.Position.Y][walker.Position.X] = walker.ID

	case ModeMultiWalker, ModeBrownianMotion:
		// Multiple walkers starting at random positions
		if walkerCount > MaxWalkerCount {
			walkerCount = MaxWalkerCount
		}
		if walkerCount < 1 {
			walkerCount = DefaultWalkerCount
		}

		for i := 0; i < walkerCount; i++ {
			// Random starting position
			x := rw.rng.IntN(rw.cols)
			y := rw.rng.IntN(rw.rows)

			walker := &Walker{
				ID:       i + 1,
				Position: Position{X: x, Y: y},
				Trail:    make([]Position, 0, rw.trailLength),
				Color:    GetWalkerColor(i),
				Visited:  make(map[Position]bool),
			}
			walker.Visited[walker.Position] = true
			rw.walkers = append(rw.walkers, walker)
			rw.grid[walker.Position.Y][walker.Position.X] = walker.ID
		}
	}
}

// Step advances the random walk by one step
func (rw *RandomWalk) Step() bool {
	for _, walker := range rw.walkers {
		rw.moveWalker(walker)
	}

	rw.steps++
	rw.updateTrails()

	return true
}

// moveWalker moves a single walker according to the walk mode
func (rw *RandomWalk) moveWalker(walker *Walker) {
	// Clear current position
	if rw.grid[walker.Position.Y][walker.Position.X] == walker.ID {
		rw.grid[walker.Position.Y][walker.Position.X] = 0
	}

	// Add current position to trail
	if rw.mode == ModeTrailMode || rw.mode == ModeBrownianMotion {
		walker.Trail = append(walker.Trail, walker.Position)
		if len(walker.Trail) > rw.trailLength {
			walker.Trail = walker.Trail[1:]
		}
	}

	// Calculate next position based on mode
	var newPos Position

	switch rw.mode {
	case ModeSelfAvoidingWalk:
		newPos = rw.getSelfAvoidingNextPosition(walker)
		if newPos == walker.Position {
			// No valid moves, walker is stuck
			return
		}

	case ModeLevyFlight:
		newPos = rw.getLevyFlightNextPosition(walker)

	case ModeBrownianMotion:
		// Brownian motion with smaller steps
		angle := rw.rng.Float64() * 2 * math.Pi
		distance := rw.rng.Float64() * 2
		dx := int(math.Round(distance * math.Cos(angle)))
		dy := int(math.Round(distance * math.Sin(angle)))
		newPos = Position{
			X: walker.Position.X + dx,
			Y: walker.Position.Y + dy,
		}

	default:
		// Regular random walk (4 or 8 directions)
		directions := rw.getDirections()
		dir := directions[rw.rng.IntN(len(directions))]
		newPos = rw.applyDirection(walker.Position, dir)
	}

	// Wrap around boundaries
	newPos.X = (newPos.X + rw.cols) % rw.cols
	newPos.Y = (newPos.Y + rw.rows) % rw.rows

	// Update walker position
	walker.Position = newPos
	walker.Visited[newPos] = true

	// Update grid
	rw.grid[walker.Position.Y][walker.Position.X] = walker.ID
}

// getSelfAvoidingNextPosition returns the next position for self-avoiding walk
func (rw *RandomWalk) getSelfAvoidingNextPosition(walker *Walker) Position {
	directions := rw.getDirections()
	validMoves := make([]Direction, 0)

	// Find all valid moves (not visited positions)
	for _, dir := range directions {
		newPos := rw.applyDirection(walker.Position, dir)
		// Wrap around
		newPos.X = (newPos.X + rw.cols) % rw.cols
		newPos.Y = (newPos.Y + rw.rows) % rw.rows

		if !walker.Visited[newPos] {
			validMoves = append(validMoves, dir)
		}
	}

	// If no valid moves, walker is stuck
	if len(validMoves) == 0 {
		return walker.Position
	}

	// Choose random valid move
	dir := validMoves[rw.rng.IntN(len(validMoves))]
	newPos := rw.applyDirection(walker.Position, dir)
	newPos.X = (newPos.X + rw.cols) % rw.cols
	newPos.Y = (newPos.Y + rw.rows) % rw.rows

	return newPos
}

// getLevyFlightNextPosition returns the next position for Lévy flight
func (rw *RandomWalk) getLevyFlightNextPosition(walker *Walker) Position {
	// Lévy flight: occasional long jumps
	angle := rw.rng.Float64() * 2 * math.Pi

	// Lévy distribution approximation
	u := rw.rng.Float64()
	if u > 0.9 { // 10% chance of long jump
		distance := rw.rng.Float64() * float64(min(rw.rows, rw.cols)) / 4
		dx := int(math.Round(distance * math.Cos(angle)))
		dy := int(math.Round(distance * math.Sin(angle)))
		return Position{
			X: walker.Position.X + dx,
			Y: walker.Position.Y + dy,
		}
	}

	// Regular short step
	directions := rw.getDirections()
	dir := directions[rw.rng.IntN(len(directions))]
	return rw.applyDirection(walker.Position, dir)
}

// getDirections returns available directions based on walk mode
func (rw *RandomWalk) getDirections() []Direction {
	// For most modes, use 8 directions
	return []Direction{
		DirectionUp, DirectionDown, DirectionLeft, DirectionRight,
		DirectionUpLeft, DirectionUpRight, DirectionDownLeft, DirectionDownRight,
	}
}

// applyDirection applies a direction to a position
func (rw *RandomWalk) applyDirection(pos Position, dir Direction) Position {
	switch dir {
	case DirectionUp:
		return Position{X: pos.X, Y: pos.Y - 1}
	case DirectionDown:
		return Position{X: pos.X, Y: pos.Y + 1}
	case DirectionLeft:
		return Position{X: pos.X - 1, Y: pos.Y}
	case DirectionRight:
		return Position{X: pos.X + 1, Y: pos.Y}
	case DirectionUpLeft:
		return Position{X: pos.X - 1, Y: pos.Y - 1}
	case DirectionUpRight:
		return Position{X: pos.X + 1, Y: pos.Y - 1}
	case DirectionDownLeft:
		return Position{X: pos.X - 1, Y: pos.Y + 1}
	case DirectionDownRight:
		return Position{X: pos.X + 1, Y: pos.Y + 1}
	default:
		return pos
	}
}

// updateTrails updates the trail intensity grid
func (rw *RandomWalk) updateTrails() {
	// Decay existing trails
	for i := range rw.rows {
		for j := range rw.cols {
			if rw.trails[i][j] > 0 {
				rw.trails[i][j]--
			}
		}
	}

	// Add current walker trails
	for _, walker := range rw.walkers {
		for i, pos := range walker.Trail {
			intensity := (i + 1) * 255 / len(walker.Trail) // Gradient intensity
			if intensity > rw.trails[pos.Y][pos.X] {
				rw.trails[pos.Y][pos.X] = intensity
			}
		}
	}
}

// GetGrid returns the current grid state
func (rw *RandomWalk) GetGrid() [][]int {
	return rw.grid
}

// GetTrails returns the trail intensity grid
func (rw *RandomWalk) GetTrails() [][]int {
	return rw.trails
}

// GetWalkers returns all walkers
func (rw *RandomWalk) GetWalkers() []*Walker {
	return rw.walkers
}

// GetSteps returns the number of steps taken
func (rw *RandomWalk) GetSteps() int {
	return rw.steps
}

// Reset resets the random walk
func (rw *RandomWalk) Reset(rows, cols int, mode WalkMode, walkerCount int, trailLength int) {
	slog.Debug("RandomWalk Reset", "rows", rows, "cols", cols, "mode", mode, "walkerCount", walkerCount, "trailLength", trailLength)
	rw.rows = rows
	rw.cols = cols
	rw.mode = mode
	rw.trailLength = trailLength
	rw.steps = 0
	rw.Init(walkerCount)
}
