package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the application state
type Model struct {
	mandelbrotSet *MandelbrotSet
	rows          int
	cols          int
	language      Language
	renderOptions RenderOptions
	refreshRate   time.Duration
	calculating   bool
	currentPreset int
	quitting      bool
	// String builders for performance
	gridBuilder   strings.Builder
	statusBuilder strings.Builder
}

// NewModel creates a new model with the given configuration
func NewModel(cfg Config) Model {
	cfg.Check()

	model := Model{
		mandelbrotSet: NewMandelbrotSet(cfg),
		rows:          cfg.Rows,
		cols:          cfg.Cols,
		language:      cfg.Language,
		renderOptions: NewRenderOptions(cfg.ColorScheme),
		refreshRate:   DefaultRefreshRate,
		calculating:   false,
		currentPreset: 0,
		quitting:      false,
	}

	// Pre-allocate string builders
	estimatedGridSize := cfg.Rows * (cfg.Cols + 10)
	estimatedStatusSize := 500

	model.gridBuilder.Grow(estimatedGridSize)
	model.statusBuilder.Grow(estimatedStatusSize)

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
	if m.quitting {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case calculationMsg:
		m.calculating = false
		return m, nil
	}
	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.calculating {
		// Ignore input while calculating except quit
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil
	}

	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.quitting = true
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
		newIter := minInt(currentIter+10, MaxMaxIterations)
		m.mandelbrotSet.SetMaxIterations(newIter)
		return m.recalculate()
	case "k", "K":
		currentIter := m.mandelbrotSet.GetMaxIterations()
		newIter := maxInt(currentIter-10, MinMaxIterations)
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
		m.mandelbrotSet.Reset()
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

// View renders the current state
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return m.RenderFullUI()
}

// RenderFullUI renders the complete user interface
func (m Model) RenderFullUI() string {
	m.statusBuilder.Reset()

	// Header
	m.statusBuilder.WriteString(GetHeaderLine(m.language))
	m.statusBuilder.WriteString("\n")

	// Status line
	m.statusBuilder.WriteString(GetStatusLine(m.mandelbrotSet, m.language))
	m.statusBuilder.WriteString("\n")

	// Julia parameter line (if in Julia mode)
	if juliaLine := GetJuliaParameterLine(m.mandelbrotSet, m.language); juliaLine != "" {
		m.statusBuilder.WriteString(juliaLine)
		m.statusBuilder.WriteString("\n")
	}

	// Main fractal grid
	if m.calculating {
		m.statusBuilder.WriteString(m.renderCalculatingMessage())
	} else {
		m.statusBuilder.WriteString(m.RenderGrid())
	}

	// Help line
	m.statusBuilder.WriteString("\n")
	m.statusBuilder.WriteString(GetHelpLine(m.language))

	// Current preset info
	if preset := m.getCurrentPresetInfo(); preset != "" {
		m.statusBuilder.WriteString("\n")
		m.statusBuilder.WriteString(preset)
	}

	return m.statusBuilder.String()
}

// RenderGrid renders the fractal grid
func (m Model) RenderGrid() string {
	m.gridBuilder.Reset()
	grid := m.mandelbrotSet.GetGrid()
	maxIter := m.mandelbrotSet.GetMaxIterations()

	for y := 0; y < m.rows; y++ {
		for x := 0; x < m.cols; x++ {
			if y < len(grid) && x < len(grid[y]) {
				iter := grid[y][x]
				char := m.renderOptions.GetCharacterForIteration(iter, maxIter)
				color := m.renderOptions.GetColorForIteration(iter, maxIter)

				// Create styled character
				style := lipgloss.NewStyle().Foreground(color)
				m.gridBuilder.WriteString(style.Render(char))
			} else {
				m.gridBuilder.WriteString(" ")
			}
		}
		if y < m.rows-1 {
			m.gridBuilder.WriteString("\n")
		}
	}

	return m.gridBuilder.String()
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
	lines := make([]string, m.rows)
	midRow := m.rows / 2

	for i := range lines {
		if i == midRow {
			padding := (m.cols - len(msg)) / 2
			if padding > 0 {
				lines[i] = strings.Repeat(" ", padding) + msg
			} else {
				lines[i] = msg
			}
		} else {
			lines[i] = strings.Repeat(" ", m.cols)
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

// Helper functions for min/max since they might not be available in older Go versions
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
