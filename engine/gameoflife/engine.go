// Package gameoflife provides the core engine implementations for Conway's Game of Life.
package gameoflife

import (
	"log/slog"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/telepair/go-playground/pkg/ui"
)

var _ ui.StepEngine = (*ConwayGameOfLife)(nil)

var (
	// HeaderCN is the Chinese header text for Conway's Game of Life
	HeaderCN = "🎮 康威生命游戏 🎮"
	// HeaderEN is the English header text for Conway's Game of Life
	HeaderEN = "🎮 Conway's Game of Life 🎮"

	// DefaultAliveColor is the default alive cell color
	DefaultAliveColor = lipgloss.Color("#00FF00")
	// DefaultDeadColor is the default dead cell color
	DefaultDeadColor = lipgloss.Color("#000000")
	// DefaultAliveChar is the default alive cell character
	DefaultAliveChar = '█'
	// DefaultDeadChar is the default dead cell character
	DefaultDeadChar = ' '
)

// Language represents the supported languages
type Language int

// Language constants
const (
	English Language = iota
	Chinese
)

// BoundaryType represents the boundary type of the Game of Life
type BoundaryType int

// BoundaryType constants
const (
	BoundaryPeriodic BoundaryType = iota // Periodic boundary (wrapping, default)
	BoundaryFixed                        // Fixed boundary (dead cells outside)
)

// ToString returns the string representation of boundary type
func (bt BoundaryType) ToString(language Language) string {
	switch bt {
	case BoundaryPeriodic:
		if language == Chinese {
			return "周期"
		}
		return "Periodic"
	case BoundaryFixed:
		if language == Chinese {
			return "固定"
		}
		return "Fixed"
	}
	if language == Chinese {
		return "周期"
	}
	return "Periodic"
}

// Pattern represents different starting patterns for Conway's Game of Life
type Pattern int

// Pattern constants
const (
	PatternRandom Pattern = iota
	PatternGlider
	PatternGliderGun
	PatternOscillator
	PatternPulsar
	PatternPentomino
)

// ToString returns the string representation of pattern type
func (p Pattern) ToString(language Language) string {
	switch p {
	case PatternRandom:
		if language == Chinese {
			return "随机"
		}
		return "random"
	case PatternGlider:
		if language == Chinese {
			return "滑翔机"
		}
		return "glider"
	case PatternGliderGun:
		if language == Chinese {
			return "滑翔机枪"
		}
		return "glider-gun"
	case PatternOscillator:
		if language == Chinese {
			return "振荡器"
		}
		return "oscillator"
	case PatternPulsar:
		if language == Chinese {
			return "脉冲星"
		}
		return "pulsar"
	case PatternPentomino:
		if language == Chinese {
			return "五格骨牌"
		}
		return "pentomino"
	default:
		if language == Chinese {
			return "随机"
		}
		return "random"
	}
}

// Config holds configuration for Conway's Game of Life
type Config struct {
	AliveColor string
	DeadColor  string
	AliveChar  string
	DeadChar   string
}

// Default values
const (
	DefaultBoundary = BoundaryPeriodic
	DefaultPattern  = PatternRandom
)

// ConwayGameOfLife represents Conway's Game of Life engine
type ConwayGameOfLife struct {
	currentGrid [][]bool
	nextGrid    [][]bool
	rows        int
	cols        int
	generation  int
	boundary    BoundaryType
	pattern     Pattern
	screen      *ui.Screen
	buf         []rune
	paused      bool
	config      Config
}

// New creates a new Conway's Game of Life instance
func New(rows, cols int, config Config) *ConwayGameOfLife {
	slog.Debug("NewConwayGameOfLife", "rows", rows, "cols", cols, "config", config)

	game := &ConwayGameOfLife{
		rows:       rows,
		cols:       cols,
		boundary:   DefaultBoundary,
		pattern:    DefaultPattern,
		generation: 0,
		paused:     false,
		config:     config,
	}
	game.initial()
	return game
}

// View returns the view of the game
func (g *ConwayGameOfLife) View() string {
	return g.screen.View()
}

// Step advances the game by one generation
func (g *ConwayGameOfLife) Step() (int, bool) {
	if g.paused {
		return g.generation, true
	}

	// Apply Conway's Game of Life rules
	for i := range g.rows {
		for j := range g.cols {
			neighbors := g.countNeighbors(i, j)
			currentCell := g.currentGrid[i][j]

			// Conway's rules
			if currentCell {
				g.nextGrid[i][j] = (neighbors == 2 || neighbors == 3)
			} else {
				g.nextGrid[i][j] = (neighbors == 3)
			}
		}
	}

	// Swap grids
	g.currentGrid, g.nextGrid = g.nextGrid, g.currentGrid
	g.generation++

	g.render()
	return g.generation, true
}

// Header returns the header text for the UI
func (g *ConwayGameOfLife) Header(lang ui.Language) string {
	if lang == ui.Chinese {
		return HeaderCN
	}
	return HeaderEN
}

// Status returns the status text for the UI
func (g *ConwayGameOfLife) Status(lang ui.Language) []ui.Status {
	// Convert ui.Language to our Language type
	var language Language
	if lang == ui.Chinese {
		language = Chinese
	} else {
		language = English
	}

	pausedStr := "▶️ Running"
	if g.paused {
		pausedStr = "⏸️ Paused"
	}
	if language == Chinese {
		if g.paused {
			pausedStr = "⏸️ 已暂停"
		} else {
			pausedStr = "▶️ 运行中"
		}
		return []ui.Status{
			{Label: "代数", Value: strconv.Itoa(g.generation)},
			{Label: "模式", Value: g.pattern.ToString(language)},
			{Label: "边界", Value: g.boundary.ToString(language)},
			{Label: "状态", Value: pausedStr},
		}
	}
	return []ui.Status{
		{Label: "Generation", Value: strconv.Itoa(g.generation)},
		{Label: "Pattern", Value: g.pattern.ToString(language)},
		{Label: "Boundary", Value: g.boundary.ToString(language)},
		{Label: "Status", Value: pausedStr},
	}
}

// HandleKeys returns the available keyboard controls
func (g *ConwayGameOfLife) HandleKeys(lang ui.Language) []ui.Control {
	if lang == ui.Chinese {
		return []ui.Control{
			{Keys: []string{"Space"}, Label: "暂停/继续"},
			{Keys: []string{"P"}, Label: "切换模式"},
			{Keys: []string{"B"}, Label: "切换边界"},
			{Keys: []string{"R"}, Label: "重置"},
		}
	}
	return []ui.Control{
		{Keys: []string{"Space"}, Label: "Pause/Resume"},
		{Keys: []string{"P"}, Label: "Pattern"},
		{Keys: []string{"B"}, Label: "Boundary"},
		{Keys: []string{"R"}, Label: "Reset"},
	}
}

// Handle handles the key press
func (g *ConwayGameOfLife) Handle(key string) (bool, error) {
	slog.Debug("ConwayGameOfLife Handle", "key", key)
	key = strings.ToLower(key)

	switch key {
	case " ", "space":
		g.paused = !g.paused
		slog.Debug("ConwayGameOfLife Handle", "key", key, "paused", g.paused)
		return true, nil

	case "p":
		// Cycle through patterns
		patterns := []Pattern{
			PatternRandom,
			PatternGlider,
			PatternGliderGun,
			PatternOscillator,
			PatternPulsar,
			PatternPentomino,
		}

		currentIdx := 0
		for i, p := range patterns {
			if p == g.pattern {
				currentIdx = i
				break
			}
		}

		nextIdx := (currentIdx + 1) % len(patterns)
		g.pattern = patterns[nextIdx]
		slog.Debug("ConwayGameOfLife Handle", "key", key, "pattern", g.pattern)
		g.initial()
		return true, nil

	case "b":
		// Toggle boundary
		if g.boundary == BoundaryPeriodic {
			g.boundary = BoundaryFixed
		} else {
			g.boundary = BoundaryPeriodic
		}
		slog.Debug("ConwayGameOfLife Handle", "key", key, "boundary", g.boundary)
		return true, nil

	case "r":
		// Reset
		slog.Debug("ConwayGameOfLife Handle", "key", key, "action", "reset")
		g.initial()
		return true, nil
	}

	slog.Debug("ConwayGameOfLife Handle", "key", key, "warning", "key not handled")
	return false, nil
}

// Reset resets the game to its initial state
func (g *ConwayGameOfLife) Reset(rows, cols int) error {
	slog.Debug("ConwayGameOfLife Reset", "rows", rows, "cols", cols)
	g.rows = rows
	g.cols = cols
	g.initial()
	return nil
}

// IsFinished returns whether the game has finished
func (g *ConwayGameOfLife) IsFinished() bool {
	return false // Conway's Game of Life runs indefinitely
}

// Stop stops the game
func (g *ConwayGameOfLife) Stop() {
	// Nothing to do for Conway's Game of Life
}

// countNeighbors counts the number of living neighbors for a cell
func (g *ConwayGameOfLife) countNeighbors(row, col int) int {
	if row < 0 || row >= g.rows || col < 0 || col >= g.cols {
		return 0
	}

	count := 0

	// Check all 8 neighbors
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue // Skip self
			}

			r, c := row+dr, col+dc

			// Handle boundary conditions
			if g.boundary == BoundaryPeriodic {
				// Wrap around
				if r < 0 {
					r = g.rows - 1
				} else if r >= g.rows {
					r = 0
				}
				if c < 0 {
					c = g.cols - 1
				} else if c >= g.cols {
					c = 0
				}
			} else {
				// Fixed boundary - treat out of bounds as dead
				if r < 0 || r >= g.rows || c < 0 || c >= g.cols {
					continue
				}
			}

			if g.currentGrid[r][c] {
				count++
			}
		}
	}

	return count
}

// initial initializes the game
func (g *ConwayGameOfLife) initial() {
	// Initialize screen
	if g.screen == nil {
		g.screen = ui.NewScreen(g.rows, g.cols)
	} else {
		g.screen.SetSize(g.cols, g.rows)
		g.screen.Reset()
	}

	// Get character runes
	aliveRune := []rune(g.config.AliveChar)[0]
	deadRune := []rune(g.config.DeadChar)[0]

	// Set colors
	g.screen.SetCharColor(aliveRune, lipgloss.Color(g.config.AliveColor))
	g.screen.SetCharColor(deadRune, lipgloss.Color(g.config.DeadColor))

	// Initialize grids
	g.currentGrid = make([][]bool, g.rows)
	g.nextGrid = make([][]bool, g.rows)
	for i := range g.rows {
		g.currentGrid[i] = make([]bool, g.cols)
		g.nextGrid[i] = make([]bool, g.cols)
	}

	g.buf = make([]rune, g.cols)
	g.generation = 0

	// Set initial pattern
	g.setInitialPattern()
	g.render()
}

// setInitialPattern sets the initial pattern
func (g *ConwayGameOfLife) setInitialPattern() {
	switch g.pattern {
	case PatternRandom:
		g.setRandomPattern()
	case PatternGlider:
		g.setGliderPattern()
	case PatternGliderGun:
		g.setGliderGunPattern()
	case PatternOscillator:
		g.setOscillatorPattern()
	case PatternPulsar:
		g.setPulsarPattern()
	case PatternPentomino:
		g.setPentominoPattern()
	default:
		g.setRandomPattern()
	}
}

// setRandomPattern creates a random initial pattern
func (g *ConwayGameOfLife) setRandomPattern() {
	seed := uint64(time.Now().UnixNano())    //nolint:gosec
	rng := rand.New(rand.NewPCG(seed, seed)) //nolint:gosec

	for i := range g.rows {
		for j := range g.cols {
			g.currentGrid[i][j] = rng.Uint32()%10 < 3 // 30% probability
		}
	}
}

// setGliderPattern creates a glider pattern
func (g *ConwayGameOfLife) setGliderPattern() {
	g.clearGrid()

	if g.rows >= 5 && g.cols >= 5 {
		startRow := 2
		startCol := 2
		pattern := [][]bool{
			{false, true, false},
			{false, false, true},
			{true, true, true},
		}
		g.placePattern(startRow, startCol, pattern)
	}
}

// setGliderGunPattern creates a Gosper glider gun pattern
func (g *ConwayGameOfLife) setGliderGunPattern() {
	g.clearGrid()

	if g.rows >= 15 && g.cols >= 40 {
		startRow := 2
		startCol := 2
		pattern := [][]bool{
			{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, true, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, false, false, true, true, false, false, false, false, false, false, true, true, false, false, false, false, false, false, false, false, false, false, false, false, true, true},
			{false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, true, false, false, false, false, true, true, false, false, false, false, false, false, false, false, false, false, false, false, true, true},
			{true, true, false, false, false, false, false, false, false, false, true, false, false, false, false, false, true, false, false, false, true, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
			{true, true, false, false, false, false, false, false, false, false, true, false, false, false, true, false, true, true, false, false, false, false, true, false, true, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, true, false, false, false, false, false, true, false, false, false, false, false, false, false, true, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, false, false, true, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		}
		g.placePattern(startRow, startCol, pattern)
	}
}

// setOscillatorPattern creates a blinker oscillator pattern
func (g *ConwayGameOfLife) setOscillatorPattern() {
	g.clearGrid()

	if g.rows >= 5 && g.cols >= 5 {
		centerRow := g.rows / 2
		centerCol := g.cols / 2
		g.currentGrid[centerRow][centerCol-1] = true
		g.currentGrid[centerRow][centerCol] = true
		g.currentGrid[centerRow][centerCol+1] = true

		if g.cols >= 10 {
			offsetCol := centerCol + 5
			g.currentGrid[centerRow-1][offsetCol] = true
			g.currentGrid[centerRow][offsetCol] = true
			g.currentGrid[centerRow+1][offsetCol] = true
		}
	}
}

// setPulsarPattern creates a pulsar oscillator pattern
func (g *ConwayGameOfLife) setPulsarPattern() {
	g.clearGrid()

	if g.rows >= 17 && g.cols >= 17 {
		centerRow := g.rows / 2
		centerCol := g.cols / 2

		offsets := [][]int{
			{-6, -4}, {-6, -3}, {-6, -2}, {-6, 2}, {-6, 3}, {-6, 4},
			{-4, -6}, {-4, -1}, {-4, 1}, {-4, 6},
			{-3, -6}, {-3, -1}, {-3, 1}, {-3, 6},
			{-2, -6}, {-2, -1}, {-2, 1}, {-2, 6},
			{-1, -4}, {-1, -3}, {-1, -2}, {-1, 2}, {-1, 3}, {-1, 4},
			{1, -4}, {1, -3}, {1, -2}, {1, 2}, {1, 3}, {1, 4},
			{2, -6}, {2, -1}, {2, 1}, {2, 6},
			{3, -6}, {3, -1}, {3, 1}, {3, 6},
			{4, -6}, {4, -1}, {4, 1}, {4, 6},
			{6, -4}, {6, -3}, {6, -2}, {6, 2}, {6, 3}, {6, 4},
		}

		for _, offset := range offsets {
			row := centerRow + offset[0]
			col := centerCol + offset[1]
			if row >= 0 && row < g.rows && col >= 0 && col < g.cols {
				g.currentGrid[row][col] = true
			}
		}
	}
}

// setPentominoPattern creates an R-pentomino pattern
func (g *ConwayGameOfLife) setPentominoPattern() {
	g.clearGrid()

	if g.rows >= 5 && g.cols >= 5 {
		centerRow := g.rows / 2
		centerCol := g.cols / 2
		pattern := [][]bool{
			{false, true, true},
			{true, true, false},
			{false, true, false},
		}
		g.placePattern(centerRow-1, centerCol-1, pattern)
	}
}

// placePattern places a pattern at the specified position
func (g *ConwayGameOfLife) placePattern(startRow, startCol int, pattern [][]bool) {
	for i, row := range pattern {
		for j, cell := range row {
			newRow := startRow + i
			newCol := startCol + j
			if newRow >= 0 && newRow < g.rows && newCol >= 0 && newCol < g.cols {
				g.currentGrid[newRow][newCol] = cell
			}
		}
	}
}

// clearGrid clears all cells
func (g *ConwayGameOfLife) clearGrid() {
	for i := range g.rows {
		for j := range g.cols {
			g.currentGrid[i][j] = false
		}
	}
}

// render renders the current grid to the screen
func (g *ConwayGameOfLife) render() {
	// Get character runes
	aliveRune := []rune(g.config.AliveChar)[0]
	deadRune := []rune(g.config.DeadChar)[0]

	// Reset screen before rendering
	g.screen.Reset()

	// Render each row
	for i := range g.rows {
		for j := range g.cols {
			if g.currentGrid[i][j] {
				g.buf[j] = aliveRune
			} else {
				g.buf[j] = deadRune
			}
		}
		g.screen.Append(g.buf)
	}
}
