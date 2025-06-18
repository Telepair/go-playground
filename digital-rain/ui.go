package main

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	keepWidth  = 2
	keepHeight = 7
)

// Model represents the application state
type Model struct {
	rain          *DigitalRain
	language      Language
	paused        bool
	refreshRate   time.Duration
	width         int
	height        int
	gridWidth     int
	gridHeight    int
	buffer        strings.Builder
	renderOptions RenderOptions
	config        Config
	logger        *slog.Logger
}

// NewModel creates a new model with the given configuration
func NewModel(cfg Config) Model {
	cfg.Check()

	gridHeight := DefaultRows - keepHeight
	gridWidth := DefaultCols - keepWidth

	model := Model{
		rain:          NewDigitalRain(gridWidth, gridHeight, cfg.CharSet, cfg.MinSpeed, cfg.MaxSpeed, cfg.DropLength),
		language:      cfg.Language,
		paused:        false,
		refreshRate:   DefaultRefreshRate,
		width:         DefaultCols,
		height:        DefaultRows,
		gridWidth:     gridWidth,
		gridHeight:    gridHeight,
		renderOptions: NewRenderOptions(cfg.DropColor, cfg.TrailColor, cfg.BackgroundColor),
		config:        cfg,
		logger:        slog.With("module", "ui"),
	}

	return model
}

// tickMsg is sent every tick
type tickMsg time.Time

// Init initializes the model
func (m Model) Init() tea.Cmd {
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
		"language", m.language,
		"paused", m.paused,
		"refreshRate", m.refreshRate)
	return m.renderUI()
}

// handleWindowResize processes terminal window size changes
func (m Model) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	m.gridWidth = msg.Width - keepWidth
	m.gridHeight = msg.Height - keepHeight
	m.rain.Reset(m.gridWidth, m.gridHeight)
	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		return m, tea.Quit

	case " ", "enter": // Pause/resume
		m.paused = !m.paused

	case "l": // Language toggle
		if m.language == English {
			m.language = Chinese
		} else {
			m.language = English
		}

	case "+", "=", "up": // Increase speed
		m.refreshRate = max(m.refreshRate/2, MinRefreshRate)

	case "-", "_", "down": // Decrease speed
		m.refreshRate = m.refreshRate * 2

	case "r": // Reset
		m.rain.Reset(m.gridWidth, m.gridHeight)

	case "d": // Increase drop length
		if m.config.DropLength < 20 {
			m.config.DropLength++
			m.rain = NewDigitalRain(m.gridWidth, m.gridHeight, m.config.CharSet,
				m.config.MinSpeed, m.config.MaxSpeed, m.config.DropLength)
		}

	case "D": // Decrease drop length
		if m.config.DropLength > 3 {
			m.config.DropLength--
			m.rain = NewDigitalRain(m.gridWidth, m.gridHeight, m.config.CharSet,
				m.config.MinSpeed, m.config.MaxSpeed, m.config.DropLength)
		}

	case "s": // Increase max speed
		if m.config.MaxSpeed < 10 {
			m.config.MaxSpeed++
			m.rain = NewDigitalRain(m.gridWidth, m.gridHeight, m.config.CharSet,
				m.config.MinSpeed, m.config.MaxSpeed, m.config.DropLength)
		}

	case "S": // Decrease max speed
		if m.config.MaxSpeed > m.config.MinSpeed {
			m.config.MaxSpeed--
			m.rain = NewDigitalRain(m.gridWidth, m.gridHeight, m.config.CharSet,
				m.config.MinSpeed, m.config.MaxSpeed, m.config.DropLength)
		}
	}

	return m, nil
}

// handleTick processes timer ticks
func (m Model) handleTick() (tea.Model, tea.Cmd) {
	if !m.paused {
		m.rain.Step()
	}

	return m, tea.Tick(m.refreshRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// renderUI renders the complete UI
func (m Model) renderUI() string {
	m.buffer.Reset()

	// Header
	m.buffer.WriteString(m.renderHeader())
	m.buffer.WriteString("\n")
	m.buffer.WriteString(m.renderStatus())
	m.buffer.WriteString("\n\n")

	// Main grid
	m.buffer.WriteString(m.renderGrid())
	m.buffer.WriteString("\n\n")

	// Controls
	m.buffer.WriteString(m.renderControls())

	return m.buffer.String()
}

// renderHeader renders the header
func (m Model) renderHeader() string {
	title := "Digital Rain"
	if m.language == Chinese {
		title = "数字雨"
	}
	return headerStyle().Render(title)
}

// renderStatus renders the status line
func (m Model) renderStatus() string {
	status := fmt.Sprintf("Speed: %v | Drop Length: %d | Max Speed: %d",
		m.refreshRate, m.config.DropLength, m.config.MaxSpeed)

	if m.language == Chinese {
		status = fmt.Sprintf("速度: %v | 雨滴长度: %d | 最大速度: %d",
			m.refreshRate, m.config.DropLength, m.config.MaxSpeed)
	}

	if m.paused {
		if m.language == Chinese {
			status += " | [暂停]"
		} else {
			status += " | [PAUSED]"
		}
	}

	return statusStyle().Render(status)
}

// renderControls renders the control instructions
func (m Model) renderControls() string {
	var controls string
	if m.language == Chinese {
		controls = "空格: 暂停/继续 | +/-: 调整速度 | d/D: 雨滴长度 | s/S: 最大速度 | r: 重置 | l: 语言 | q: 退出"
	} else {
		controls = "Space: Pause | +/-: Speed | d/D: Drop Length | s/S: Max Speed | r: Reset | l: Language | q: Quit"
	}
	return helpStyle().Render(controls)
}

// renderGrid renders the rain grid
func (m Model) renderGrid() string {
	var sb strings.Builder
	grid := m.rain.GetGrid()
	trail := m.rain.GetTrail()

	if len(grid) == 0 {
		return ""
	}

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] != 0 {
				// Character with appropriate style based on trail intensity
				char := string(grid[i][j])
				intensity := trail[i][j]
				if intensity > 200 {
					sb.WriteString(m.renderOptions.dropStyle.Render(char))
				} else {
					style := m.renderOptions.GetTrailStyle(intensity)
					sb.WriteString(style.Render(char))
				}
			} else {
				sb.WriteString(" ")
			}
		}
		if i < len(grid)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
