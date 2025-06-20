package main

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	keepWidth  = 4
	keepHeight = 8
)

// Model represents the application state
type Model struct {
	mandelbrotSet *MandelbrotSet

	language Language

	width         int
	gridHeight    int
	gridWidth     int
	refreshRate   time.Duration
	calculating   bool
	currentPreset int
	// String builders for performance
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
		mandelbrotSet: NewMandelbrotSet(cfg),
		width:         DefaultCols,
		gridHeight:    gridHeight,
		gridWidth:     gridWidth,
		language:      cfg.Language,
		renderOptions: NewRenderOptions(cfg.ColorScheme),
		refreshRate:   DefaultRefreshRate,
		calculating:   false,
		currentPreset: 0,
		logger:        slog.With("module", "ui"),
	}

	return model
}

// calculationMsg is sent when calculation is complete
type calculationMsg time.Time

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.logger.Debug("Key pressed", "key", msg.String())
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.logger.Debug("Window size", "width", msg.Width, "height", msg.Height)
		return m.handleWindowResize(msg)
	case calculationMsg:
		m.logger.Debug("Calculation complete", "time", msg)
		m.calculating = false
		return m, nil
	}
	return m, nil
}

// View renders the current state
func (m Model) View() string {
	m.logger.Debug("Model View",
		"width", m.width,
		"gridWidth", m.gridWidth,
		"gridHeight", m.gridHeight,
		"language", m.language,
		"calculating", m.calculating,
		"currentPreset", m.currentPreset,
		"refreshRate", m.refreshRate)
	return m.RenderMode()
}

// handleWindowResize processes terminal window size changes
func (m Model) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.gridWidth = msg.Width - keepWidth
	m.gridHeight = msg.Height - keepHeight
	m.mandelbrotSet.Reset(m.gridHeight, m.gridWidth)
	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.calculating {
		// Ignore input while calculating except quit
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
		return m, nil
	}

	switch msg.String() {
	case "ctrl+c", "q", "esc":
		return m, tea.Quit

	// Pan controls
	case "up", "w":
		m.mandelbrotSet.Pan(0, -5)
		return m.recalculate()
	case "down", "s":
		m.mandelbrotSet.Pan(0, 5)
		return m.recalculate()
	case "left", "a":
		m.mandelbrotSet.Pan(-5, 0)
		return m.recalculate()
	case "right", "d":
		m.mandelbrotSet.Pan(5, 0)
		return m.recalculate()

	// Zoom controls
	case "+", "=":
		m.mandelbrotSet.ZoomIn(2.0)
		return m.recalculate()
	case "-", "_":
		m.mandelbrotSet.ZoomOut(2.0)
		return m.recalculate()

	// Mode controls
	case "m", "M":
		m.mandelbrotSet.ToggleMode()
		return m.recalculate()

	// Color scheme controls
	case "c", "C":
		currentScheme := m.mandelbrotSet.GetColorScheme()
		nextScheme := (currentScheme + 1) % 5
		m.mandelbrotSet.SetColorScheme(nextScheme)
		m.renderOptions.colorScheme = nextScheme

	// Iteration controls
	case "i", "I":
		currentIter := m.mandelbrotSet.GetMaxIterations()
		newIter := min(currentIter+10, MaxMaxIterations)
		m.mandelbrotSet.SetMaxIterations(newIter)
		return m.recalculate()
	case "k", "K":
		currentIter := m.mandelbrotSet.GetMaxIterations()
		newIter := max(currentIter-10, MinMaxIterations)
		m.mandelbrotSet.SetMaxIterations(newIter)
		return m.recalculate()

	// Preset locations
	case "p", "P":
		return m.goToNextPreset()

	// Language toggle
	case "l", "L":
		if m.language == English {
			m.language = Chinese
		} else {
			m.language = English
		}

	// Reset
	case "r", "R":
		m.mandelbrotSet.Reset(m.gridHeight, m.gridWidth)
		m.currentPreset = 0
		return m.recalculate()

	// Fine pan controls
	case "shift+up":
		m.mandelbrotSet.Pan(0, -1)
		return m.recalculate()
	case "shift+down":
		m.mandelbrotSet.Pan(0, 1)
		return m.recalculate()
	case "shift+left":
		m.mandelbrotSet.Pan(-1, 0)
		return m.recalculate()
	case "shift+right":
		m.mandelbrotSet.Pan(1, 0)
		return m.recalculate()
	}

	return m, nil
}

// recalculate starts a new calculation
func (m Model) recalculate() (tea.Model, tea.Cmd) {
	m.calculating = true
	return m, tea.Tick(m.refreshRate, func(t time.Time) tea.Msg {
		return calculationMsg(t)
	})
}

// goToNextPreset goes to the next interesting preset location
func (m Model) goToNextPreset() (tea.Model, tea.Cmd) {
	presets := m.mandelbrotSet.GetInterestingPoints()
	if len(presets) == 0 {
		return m, nil
	}

	m.currentPreset = (m.currentPreset + 1) % len(presets)
	preset := presets[m.currentPreset]

	m.mandelbrotSet.SetCenter(preset.X, preset.Y)
	m.mandelbrotSet.SetZoom(preset.Zoom)
	return m.recalculate()
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

// RenderGrid renders the fractal grid
func (m Model) RenderGrid() string {
	// Main fractal grid
	if m.calculating {
		return m.renderCalculatingMessage()
	}

	m.gridBuffer.Reset()
	grid := m.mandelbrotSet.GetGrid()
	maxIter := m.mandelbrotSet.GetMaxIterations()

	for y := 0; y < m.gridHeight; y++ {
		for x := 0; x < m.gridWidth; x++ {
			if y < len(grid) && x < len(grid[y]) {
				iter := grid[y][x]
				char := m.renderOptions.GetCharacterForIteration(iter, maxIter)
				color := m.renderOptions.GetColorForIteration(iter, maxIter)

				// Create styled character
				style := lipgloss.NewStyle().Foreground(color)
				m.gridBuffer.WriteString(style.Render(char))
			} else {
				m.gridBuffer.WriteString(" ")
			}
		}
		if y < m.gridHeight-1 {
			m.gridBuffer.WriteString("\n")
		}
	}

	return m.gridBuffer.String()
}

// renderCalculatingMessage renders the calculating message
func (m Model) renderCalculatingMessage() string {
	var msg string
	if m.language == Chinese {
		msg = "正在计算分形图案..."
	} else {
		msg = "Calculating fractal pattern..."
	}

	// Center the message in the grid area
	lines := make([]string, m.gridHeight)
	midRow := m.gridHeight / 2

	for i := range lines {
		if i == midRow {
			padding := (m.width - len(msg)) / 2
			if padding > 0 {
				lines[i] = strings.Repeat(" ", padding) + msg
			} else {
				lines[i] = msg
			}
		} else {
			lines[i] = strings.Repeat(" ", m.width)
		}
	}

	return strings.Join(lines, "\n")
}

// getCurrentPresetInfo returns information about the current preset
func (m Model) getCurrentPresetInfo() string {
	presets := m.mandelbrotSet.GetInterestingPoints()
	if len(presets) == 0 || m.currentPreset >= len(presets) {
		return ""
	}

	preset := presets[m.currentPreset]
	if m.language == Chinese {
		return helpStyle.Render(fmt.Sprintf("当前预设: %s (%d/%d)", preset.Name, m.currentPreset+1, len(presets)))
	}
	return helpStyle.Render(fmt.Sprintf("Current Preset: %s (%d/%d)", preset.Name, m.currentPreset+1, len(presets)))
}
