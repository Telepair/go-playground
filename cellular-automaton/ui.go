package main

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the application state
type Model struct {
	ca            *CellularAutomaton
	rule          int
	rows          int
	cols          int
	language      Language
	paused        bool // Pause state for infinite mode
	currentStep   int
	refreshRate   time.Duration
	boundary      BoundaryType
	renderOptions RenderOptions
	// Optimized string builders with pre-allocation
	gridBuilder    strings.Builder
	statusBuilder  strings.Builder
	gridRingBuffer *GridRingBuffer
	quitting       bool // Flag to track if we're quitting
}

// NewModel creates a new model with the given configuration
func NewModel(cfg Config) Model {
	cfg.Check()

	model := Model{
		ca:             NewCellularAutomaton(cfg.Rule, cfg.Cols, DefaultBoundary),
		rule:           cfg.Rule,
		rows:           cfg.Rows,
		cols:           cfg.Cols,
		language:       cfg.Language,
		refreshRate:    DefaultRefreshRate,
		boundary:       DefaultBoundary,
		renderOptions:  NewRenderOptions(cfg.AliveColor, cfg.DeadColor, cfg.AliveChar, cfg.DeadChar),
		gridRingBuffer: NewGridRingBuffer(cfg.Rows, cfg.Cols),
		quitting:       false,
	}

	// Pre-allocate string builders with estimated capacity
	estimatedGridSize := cfg.Rows * (cfg.Cols + 10) // +10 for formatting
	estimatedStatusSize := 500                      // Estimated status line size

	model.gridBuilder.Grow(estimatedGridSize)
	model.statusBuilder.Grow(estimatedStatusSize)

	// Initialize the ring buffer with the initial state - add safety check
	if currentRow := model.ca.GetCurrentRow(); currentRow != nil {
		model.gridRingBuffer.AddRow(currentRow)
	}

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
	// Handle normal application keys when no modal is active
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.quitting = true
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
		m.ca.Reset(m.rule, m.cols, m.boundary)
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
		m.ca.Reset(m.rule, m.cols, m.boundary)
		m.gridRingBuffer.Clear()
		m.gridRingBuffer.AddRow(m.ca.GetCurrentRow())
	case "l": // Language toggle key
		if m.language == English {
			m.language = Chinese
		} else {
			m.language = English
		}

	case "r": // Reset simulation
		m.ca.Reset(m.rule, m.cols, m.boundary)
		m.gridRingBuffer.Clear()
		m.gridRingBuffer.AddRow(m.ca.GetCurrentRow())

	case "+", "=": // Increase refresh rate (make it faster)
		m.refreshRate = max(m.refreshRate/2, MinRefreshRate)

	case "-", "_": // Decrease refresh rate (make it slower)
		m.refreshRate = m.refreshRate * 2
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
	// Check if we're quitting
	if m.quitting {
		return m, nil
	}

	// Check if we should continue running
	if !m.paused {
		if m.ca.Step() {
			m.currentStep = m.ca.GetGeneration()
			m.gridRingBuffer.AddRow(m.ca.GetCurrentRow())
		} else {
			// Cellular automaton has finished (reached max steps)
			// In finite mode, quit the program
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
		m.rule,
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

// RenderGrid renders the grid using the optimized ring buffer
func (m Model) RenderGrid() string {
	m.gridBuilder.Reset()
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

		m.gridBuilder.WriteString("  ")

		// Render cells in the row
		for _, cell := range row {
			if cell {
				m.gridBuilder.WriteString(aliveStr)
			} else {
				m.gridBuilder.WriteString(deadStr)
			}
		}

		// Add newline except for the last row
		if i < len(rows)-1 {
			m.gridBuilder.WriteByte('\n')
		}
	}

	return m.gridBuilder.String()
}
