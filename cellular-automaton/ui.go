package main

import (
	"log/slog"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	keepWidth  = 4
	keepHeight = 6
)

// Model represents the application state
type Model struct {
	ca *CellularAutomaton

	rule     int
	language Language

	paused         bool // Pause state for infinite mode
	currentStep    int
	refreshRate    time.Duration
	boundary       BoundaryType
	width          int
	gridHeight     int
	gridWidth      int
	buffer         strings.Builder
	gridBuffer     strings.Builder
	gridRingBuffer *GridRingBuffer
	renderOptions  RenderOptions
	logger         *slog.Logger
}

// NewModel creates a new model with the given configuration
func NewModel(cfg Config) Model {
	cfg.Check()

	gridHeight := DefaultRows - keepHeight
	gridWidth := DefaultCols - keepWidth
	model := Model{
		ca:             NewCellularAutomaton(cfg.Rule, DefaultCols, DefaultBoundary),
		rule:           cfg.Rule,
		language:       cfg.Language,
		refreshRate:    DefaultRefreshRate,
		boundary:       DefaultBoundary,
		width:          DefaultCols,
		gridHeight:     gridHeight,
		gridWidth:      gridWidth,
		gridRingBuffer: NewGridRingBuffer(gridHeight, gridWidth),
		renderOptions:  NewRenderOptions(cfg.AliveColor, cfg.DeadColor, cfg.AliveChar, cfg.DeadChar),
		logger:         slog.With("module", "ui"),
	}

	// Initialize the ring buffer with the initial state - add safety check
	model.gridRingBuffer.AddRow(model.ca.GetCurrentRow())

	return model
}

// infiniteMsg is sent every tick for infinite mode
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
		"width", m.width,
		"gridWidth", m.gridWidth,
		"gridHeight", m.gridHeight,
		"rule", m.rule,
		"boundary", m.boundary,
		"language", m.language,
		"paused", m.paused,
		"currentStep", m.currentStep,
		"refreshRate", m.refreshRate)
	return m.RenderMode()
}

// handleWindowResize processes terminal window size changes
func (m Model) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.gridWidth = msg.Width - keepWidth
	m.gridHeight = msg.Height - keepHeight
	m.logger.Debug("Window size changed", "width", m.width, "gridWidth", m.gridWidth, "gridHeight", m.gridHeight)
	m.ca.Reset(m.rule, m.gridWidth, m.boundary)
	m.gridRingBuffer = NewGridRingBuffer(m.gridHeight, m.gridWidth)
	m.gridRingBuffer.AddRow(m.ca.GetCurrentRow())
	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keyStr := msg.String()

	// Handle normal application keys when no modal is active
	switch keyStr {
	case "ctrl+c", "q", "esc":
		return m, tea.Quit

	case " ", "enter": // Space or Enter key for pause/resume in infinite mode
		m.paused = !m.paused

	case "t": // Toggle rule selection modal (T for "Type" rule)
		switch m.rule {
		case 30:
			m.rule = 90
		case 90:
			m.rule = 110
		case 110:
			m.rule = 154
		case 154:
			m.rule = 184
		case 184:
			m.rule = 30
		default:
			m.rule = 30
		}
		m.ca.Reset(m.rule, m.width, m.boundary)
		m.gridRingBuffer.Clear()
		m.gridRingBuffer.AddRow(m.ca.GetCurrentRow())

	case "b": // Show boundary selection modal
		switch m.boundary {
		case BoundaryPeriodic:
			m.boundary = BoundaryFixed
		case BoundaryFixed:
			m.boundary = BoundaryReflect
		case BoundaryReflect:
			m.boundary = BoundaryPeriodic
		}
		m.ca.Reset(m.rule, m.width, m.boundary)
		m.gridRingBuffer.Clear()
		m.gridRingBuffer.AddRow(m.ca.GetCurrentRow())
	case "l": // Language toggle key
		if m.language == English {
			m.language = Chinese
		} else {
			m.language = English
		}

	case "r": // Reset simulation
		m.ca.Reset(m.rule, m.width, m.boundary)
		m.gridRingBuffer.Clear()
		m.gridRingBuffer.AddRow(m.ca.GetCurrentRow())

	case "+", "=", "up": // Increase refresh rate (make it faster)
		m.refreshRate = max(m.refreshRate/2, MinRefreshRate)

	case "-", "_", "down": // Decrease refresh rate (make it slower)
		m.refreshRate = m.refreshRate * 2
	}
	return m, nil
}

// handleTick processes timer ticks
func (m Model) handleTick() (tea.Model, tea.Cmd) {
	if !m.paused && m.ca.Step() {
		m.currentStep = m.ca.GetGeneration()
		m.gridRingBuffer.AddRow(m.ca.GetCurrentRow())
	}

	// Continue ticking only if not quitting
	return m, tea.Tick(m.refreshRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// RenderMode renders the complete UI mode view with enhanced layout
// Layout: Header -> Status -> Grid -> Control (bottom)
func (m Model) RenderMode() string {
	m.buffer.Reset()

	// Build complete UI with enhanced styling - new layout order
	m.buffer.WriteString(m.HeaderLineView())
	m.buffer.WriteString("\n")
	m.buffer.WriteString(m.StatusLineView())
	m.buffer.WriteString("\n\n")
	m.buffer.WriteString(m.RenderGrid())
	m.buffer.WriteString("\n\n")
	m.buffer.WriteString(m.ControlLineView())

	return m.buffer.String()
}

// RenderGrid renders the grid using the optimized ring buffer
func (m Model) RenderGrid() string {
	m.gridBuffer.Reset()
	rows := m.gridRingBuffer.GetRows()
	if len(rows) == 0 {
		return ""
	}

	// Pre-calculate styled strings to avoid repeated lookups
	aliveStr := m.renderOptions.aliveStyled
	deadStr := m.renderOptions.deadStyled

	// Render all rows efficiently
	for i, row := range rows {
		if row == nil {
			continue // Skip nil rows
		}

		m.gridBuffer.WriteString("  ")

		// Render cells in the row
		for _, cell := range row {
			if cell {
				m.gridBuffer.WriteString(aliveStr)
			} else {
				m.gridBuffer.WriteString(deadStr)
			}
		}

		// Add newline except for the last row
		if i < len(rows)-1 {
			m.gridBuffer.WriteByte('\n')
		}
	}

	for i := 0; i < m.gridHeight-len(rows); i++ {
		m.gridBuffer.WriteString("\n")
	}

	return m.gridBuffer.String()
}
