package main

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the application state
type Model struct {
	game          *GameOfLife
	rows          int
	cols          int
	language      Language
	pattern       Pattern
	boundary      BoundaryType
	paused        bool // Pause state for infinite mode
	currentStep   int
	quitting      bool // Add flag to track if we're quitting
	renderOptions RenderOptions
	// Optimized string builders with pre-allocation
	gridBuilder   strings.Builder
	statusBuilder strings.Builder
	refreshRate   time.Duration
}

// NewModel creates a new model with the given configuration
func NewModel(cfg Config) Model {
	cfg.Check()

	model := Model{
		game:          NewGameOfLife(cfg.Rows, cfg.Cols, DefaultBoundary, DefaultPattern),
		rows:          cfg.Rows,
		cols:          cfg.Cols,
		language:      cfg.Language,
		pattern:       DefaultPattern,
		boundary:      DefaultBoundary,
		paused:        false,
		currentStep:   0,
		quitting:      false,
		renderOptions: NewRenderOptions(cfg.AliveColor, cfg.DeadColor, cfg.AliveChar, cfg.DeadChar),
		refreshRate:   DefaultRefreshRate,
	}

	// Pre-allocate string builders with estimated capacity
	estimatedGridSize := cfg.Rows * (cfg.Cols + 10) // +10 for formatting
	estimatedStatusSize := 500                      // Estimated status line size

	model.gridBuilder.Grow(estimatedGridSize)
	model.statusBuilder.Grow(estimatedStatusSize)

	return model
}

// tickMsg is sent every tick for infinite mode
type tickMsg time.Time

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Always start the timer when initializing unless already quitting
	return tea.Tick(m.refreshRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Early return if quitting to prevent processing any messages
	if m.quitting {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tickMsg:
		return m.handleTick()
	}
	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.quitting = true // Set quitting flag
		return m, tea.Quit

	case " ", "enter": // Space or Enter key for pause/resume
		m.paused = !m.paused

	case "l": // Language toggle key
		if m.language == English {
			m.language = Chinese
		} else {
			m.language = English
		}

	case "+", "=": // Increase refresh rate (make it faster)
		m.refreshRate = max(m.refreshRate/2, MinRefreshRate)

	case "-", "_": // Decrease refresh rate (make it slower)
		m.refreshRate = m.refreshRate * 2

	case "p": // Cycle through patterns
		currentPattern := int(m.pattern)
		nextPattern := (currentPattern + 1) % 6 // We have 6 patterns
		m.pattern = Pattern(nextPattern)
		m.game.SetPattern(m.pattern)
		m.game.Reset()

	case "b": // Toggle boundary type
		if m.boundary == BoundaryPeriodic {
			m.boundary = BoundaryFixed
		} else {
			m.boundary = BoundaryPeriodic
		}
		// Update game boundary and reset
		m.game.boundary = m.boundary
		m.game.Reset()

	case "r": // Reset simulation
		m.currentStep = 0
		m.game.Reset()
	}

	// Continue with normal tick command unless quitting
	if !m.quitting {
		return m, tea.Tick(m.refreshRate, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}

	return m, nil
}

// handleTick processes timer ticks
func (m Model) handleTick() (tea.Model, tea.Cmd) {
	// Stop ticking if we're quitting
	if m.quitting {
		return m, nil
	}

	// Check if we should continue running (only update when not paused)
	if !m.paused {
		if m.game.Step() {
			m.currentStep = m.game.GetGeneration()
		} else {
			// Game finished naturally, set quitting flag and exit
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Continue ticking only if not quitting
	return m, tea.Tick(m.refreshRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// View renders the current state
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return m.RenderMode()
}

// RenderMode renders the complete UI mode view with enhanced layout
func (m Model) RenderMode() string {
	m.statusBuilder.Reset()

	// Build complete UI with enhanced styling
	m.statusBuilder.WriteString(GetHeaderLine(m.language))
	m.statusBuilder.WriteString("\n")
	m.statusBuilder.WriteString(GetStatusLine(
		m.language,
		m.pattern,
		m.currentStep,
		m.refreshRate,
		m.rows,
		m.cols,
		m.boundary,
		m.paused))
	m.statusBuilder.WriteString("\n")
	m.statusBuilder.WriteString(GetControlLine(m.language))
	m.statusBuilder.WriteString("\n\n")
	m.statusBuilder.WriteString(m.RenderGrid())

	return m.statusBuilder.String()
}

// RenderGrid renders the 2D grid using optimized rendering
func (m *Model) RenderGrid() string {
	m.gridBuilder.Reset()
	grid := m.game.GetCurrentGrid()
	if len(grid) == 0 {
		return ""
	}

	// Pre-calculate styled strings to avoid repeated lookups
	aliveStr := m.renderOptions.aliveStyled
	deadStr := m.renderOptions.deadStyled

	// Render all rows efficiently with minimal allocations
	lastRowIndex := len(grid) - 1
	for i, row := range grid {
		if row == nil {
			continue // Skip nil rows
		}

		m.gridBuilder.WriteString(" ")

		// Render cells in the row with optimized string operations
		for _, cell := range row {
			if cell {
				m.gridBuilder.WriteString(aliveStr)
			} else {
				m.gridBuilder.WriteString(deadStr)
			}
		}

		// Add newline except for the last row
		if i < lastRowIndex {
			m.gridBuilder.WriteByte('\n')
		}
	}

	return m.gridBuilder.String()
}
