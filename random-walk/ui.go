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
	walk *RandomWalk

	language    Language
	mode        WalkMode
	walkerCount int
	trailLength int

	paused        bool
	currentStep   int
	refreshRate   time.Duration
	width         int
	gridHeight    int
	gridWidth     int
	buffer        strings.Builder
	gridBuffer    strings.Builder
	renderOptions RenderOptions
	logger        *slog.Logger
}

// NewModel creates a new model with the given configuration
func NewModel(cfg Config) Model {
	cfg.Check()

	gridHeight := DefaultRows - keepHeight
	gridWidth := DefaultCols - keepWidth

	model := Model{
		walk:          NewRandomWalk(gridHeight, gridWidth, DefaultWalkMode, DefaultWalkerCount, DefaultTrailLength),
		language:      cfg.Language,
		mode:          DefaultWalkMode,
		walkerCount:   DefaultWalkerCount,
		trailLength:   DefaultTrailLength,
		width:         DefaultCols,
		gridHeight:    gridHeight,
		gridWidth:     gridWidth,
		paused:        false,
		currentStep:   0,
		renderOptions: NewRenderOptions(cfg.WalkerColor, cfg.TrailColor, cfg.EmptyColor, cfg.WalkerChar, cfg.TrailChar, cfg.EmptyChar),
		refreshRate:   DefaultRefreshRate,
		logger:        slog.With("module", "ui"),
	}

	return model
}

// tickMsg is sent every tick for infinite mode
type tickMsg time.Time

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Always start the timer when initializing
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
		"mode", m.mode,
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
	m.walk.Reset(m.gridHeight, m.gridWidth, m.mode, m.walkerCount, m.trailLength)
	m.currentStep = 0
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

	case "m": // Cycle through walk modes
		m.mode = WalkMode((int(m.mode) + 1) % 6) // We have 6 modes
		m.walk.Reset(m.gridHeight, m.gridWidth, m.mode, m.walkerCount, m.trailLength)
		m.currentStep = 0

	case "w": // Increase walker count (for multi-walker modes)
		if m.walkerCount < MaxWalkerCount {
			m.walkerCount++
			if m.mode == ModeMultiWalker || m.mode == ModeBrownianMotion {
				m.walk.Reset(m.gridHeight, m.gridWidth, m.mode, m.walkerCount, m.trailLength)
				m.currentStep = 0
			}
		}

	case "W": // Decrease walker count (for multi-walker modes)
		if m.walkerCount > 1 {
			m.walkerCount--
			if m.mode == ModeMultiWalker || m.mode == ModeBrownianMotion {
				m.walk.Reset(m.gridHeight, m.gridWidth, m.mode, m.walkerCount, m.trailLength)
				m.currentStep = 0
			}
		}

	case "t": // Increase trail length
		if m.trailLength < MaxTrailLength {
			m.trailLength += 10
			m.walk.Reset(m.gridHeight, m.gridWidth, m.mode, m.walkerCount, m.trailLength)
			m.currentStep = 0
		}

	case "T": // Decrease trail length
		if m.trailLength > 10 {
			m.trailLength -= 10
			m.walk.Reset(m.gridHeight, m.gridWidth, m.mode, m.walkerCount, m.trailLength)
			m.currentStep = 0
		}

	case "r": // Reset simulation
		m.currentStep = 0
		m.walk.Reset(m.gridHeight, m.gridWidth, m.mode, m.walkerCount, m.trailLength)
	}

	return m, nil
}

// handleTick processes timer ticks
func (m Model) handleTick() (tea.Model, tea.Cmd) {
	// Check if we should continue running (only update when not paused)
	if !m.paused && m.walk.Step() {
		m.currentStep = m.walk.GetSteps()
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
	grid := m.walk.GetGrid()
	trails := m.walk.GetTrails()
	walkers := m.walk.GetWalkers()

	if len(grid) == 0 {
		return ""
	}

	// Pre-calculate styled strings to avoid repeated lookups
	emptyStr := m.renderOptions.emptyStyled
	trailStr := m.renderOptions.trailStyled

	// Create walker styled strings
	walkerStyles := make(map[int]string)
	for _, walker := range walkers {
		walkerStyles[walker.ID] = m.renderOptions.getWalkerStyled(walker.Color, m.renderOptions.walkerChar)
	}

	// Render all rows efficiently with minimal allocations
	lastRowIndex := len(grid) - 1
	for i, row := range grid {
		if row == nil {
			continue // Skip nil rows
		}

		m.gridBuffer.WriteString(" ")

		// Render cells in the row with optimized string operations
		for j, cell := range row {
			if cell > 0 {
				// Walker at this position
				m.gridBuffer.WriteString(walkerStyles[cell])
			} else if trails[i][j] > 0 {
				// Trail at this position
				m.gridBuffer.WriteString(trailStr)
			} else {
				// Empty cell
				m.gridBuffer.WriteString(emptyStr)
			}
		}

		// Add newline except for the last row
		if i < lastRowIndex {
			m.gridBuffer.WriteByte('\n')
		}
	}

	return m.gridBuffer.String()
}
