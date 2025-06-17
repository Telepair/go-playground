package main

import (
	"log/slog"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the application state
type Model struct {
	game *GameOfLife

	language Language
	pattern  Pattern

	paused        bool // Pause state for infinite mode
	currentStep   int
	refreshRate   time.Duration
	boundary      BoundaryType
	height        int
	width         int
	buffer        strings.Builder
	gridBuffer    strings.Builder
	renderOptions RenderOptions
	logger        *slog.Logger
}

// NewModel creates a new model with the given configuration
func NewModel(cfg Config) Model {
	cfg.Check()

	model := Model{
		game:          NewGameOfLife(DefaultRows, DefaultCols, DefaultBoundary, DefaultPattern),
		language:      cfg.Language,
		pattern:       DefaultPattern,
		boundary:      DefaultBoundary,
		paused:        false,
		currentStep:   0,
		renderOptions: NewRenderOptions(cfg.AliveColor, cfg.DeadColor, cfg.AliveChar, cfg.DeadChar),
		refreshRate:   DefaultRefreshRate,
		logger:        slog.With("module", "ui"),
	}

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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.logger.Debug("Window size changed", "width", msg.Width, "height", msg.Height)
		return m.handleWindowResize(msg)
	case tea.KeyMsg:
		m.logger.Debug("Key pressed", "key", msg.String())
		return m.handleKeyPress(msg)
	case tickMsg:
		m.logger.Debug("Tick", "time", msg)
		return m.handleTick()
	}
	return m, nil
}

// View renders the current state
func (m Model) View() string {
	m.logger.Debug("Model View",
		"height", m.height,
		"width", m.width,
		"pattern", m.pattern,
		"boundary", m.boundary,
		"language", m.language,
		"paused", m.paused,
		"currentStep", m.currentStep,
		"refreshRate", m.refreshRate)
	return m.RenderMode()
}

// handleWindowResize processes terminal window size changes
func (m Model) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width - 4
	m.height = msg.Height - 6
	m.game.Reset(m.height, m.width, m.boundary, m.pattern)
	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		return m, tea.Quit

	case " ", "enter": // Space or Enter key for pause/resume
		m.paused = !m.paused

	case "l": // Language toggle key
		if m.language == English {
			m.language = Chinese
		} else {
			m.language = English
		}

	case "+", "=", "up": // Increase refresh rate (make it faster)
		m.refreshRate = max(m.refreshRate/2, MinRefreshRate)

	case "-", "_", "down": // Decrease refresh rate (make it slower)
		m.refreshRate = m.refreshRate * 2

	case "p": // Cycle through patterns
		m.pattern = Pattern((int(m.pattern) + 1) % 6) // We have 6 patterns
		m.game.Reset(m.height, m.width, m.boundary, m.pattern)

	case "b": // Toggle boundary type
		if m.boundary == BoundaryPeriodic {
			m.boundary = BoundaryFixed
		} else {
			m.boundary = BoundaryPeriodic
		}
		m.game.Reset(m.height, m.width, m.boundary, m.pattern)

	case "r": // Reset simulation
		m.currentStep = 0
		m.game.Reset(m.height, m.width, m.boundary, m.pattern)
	}

	return m, nil
}

// handleTick processes timer ticks
func (m Model) handleTick() (tea.Model, tea.Cmd) {
	// Check if we should continue running (only update when not paused)
	if !m.paused && m.game.Step() {
		m.currentStep = m.game.GetGeneration()
	}

	// Continue ticking only if not quitting
	return m, tea.Tick(m.refreshRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// RenderMode renders the complete UI mode view with enhanced layout
func (m Model) RenderMode() string {
	m.buffer.Reset()

	// Build complete UI with enhanced styling
	m.buffer.WriteString(m.HeaderLineView())
	m.buffer.WriteString("\n")
	m.buffer.WriteString(m.StatusLineView())
	m.buffer.WriteString("\n\n")
	m.buffer.WriteString(m.RenderGrid())
	m.buffer.WriteString("\n\n")
	m.buffer.WriteString(m.ControlLineView())

	return m.buffer.String()
}

// RenderGrid renders the 2D grid using optimized rendering
func (m *Model) RenderGrid() string {
	m.gridBuffer.Reset()
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

		m.gridBuffer.WriteString(" ")

		// Render cells in the row with optimized string operations
		for _, cell := range row {
			if cell {
				m.gridBuffer.WriteString(aliveStr)
			} else {
				m.gridBuffer.WriteString(deadStr)
			}
		}

		// Add newline except for the last row
		if i < lastRowIndex {
			m.gridBuffer.WriteByte('\n')
		}
	}

	return m.gridBuffer.String()
}
