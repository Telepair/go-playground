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
	game          *GameOfLife
	cfg           *Config
	paused        bool // Pause state for infinite mode
	currentStep   int
	renderOptions *RenderOptions
	// Optimized string builders with pre-allocation
	gridBuilder   strings.Builder
	statusBuilder strings.Builder
	refreshRate   time.Duration
	pattern       Pattern
	boundary      BoundaryType
	quitting      bool // Add flag to track if we're quitting
}

// TableData represents a 2x3 table structure
type TableData struct {
	Rows [][]string // 2 rows, 3 columns each
}

// NewModel creates a new model with the given configuration
func NewModel(cfg *Config) Model {
	// Add safety checks
	if cfg == nil {
		cfg = NewConfig() // Use default config if nil
	}

	// Ensure minimum dimensions
	if cfg.Rows < 1 {
		cfg.Rows = 1
	}
	if cfg.Cols < 1 {
		cfg.Cols = 1
	}

	game := NewGameOfLife(cfg.Rows, cfg.Cols, DefaultBoundary, DefaultPattern)
	if game == nil {
		// Handle case where game creation fails
		configLogger.Printf("Failed to create game, using default settings")
		game = NewGameOfLife(DefaultWindowRows, DefaultWindowCols, DefaultBoundary, DefaultPattern)
	}

	model := Model{
		game:          game,
		cfg:           cfg,
		renderOptions: NewRenderOptions(cfg),
		refreshRate:   DefaultRefreshRate,
		pattern:       DefaultPattern,
		boundary:      DefaultBoundary,
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
	// Start ticking for Game of Life only if not quitting
	if !m.quitting {
		return tea.Tick(m.refreshRate, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		if m.cfg.Language == English {
			m.cfg.Language = Chinese
		} else {
			m.cfg.Language = English
		}

	case "+", "=": // Increase refresh rate (make it faster)
		newRate := m.refreshRate / 2
		if newRate >= MinRefreshRate {
			m.refreshRate = newRate
		} else {
			m.refreshRate = MinRefreshRate
		}

	case "-", "_": // Decrease refresh rate (make it slower)
		m.refreshRate = m.refreshRate * 2

	case "p": // Cycle through patterns
		currentPattern := int(m.pattern)
		nextPattern := (currentPattern + 1) % 6 // We have 6 patterns
		m.pattern = Pattern(nextPattern)
		if m.game != nil {
			m.game.SetPattern(m.pattern)
		}

	case "b": // Toggle boundary type
		if m.boundary == BoundaryPeriodic {
			m.boundary = BoundaryFixed
		} else {
			m.boundary = BoundaryPeriodic
		}
		// Reset game with new boundary
		if m.game != nil {
			m.game = NewGameOfLife(m.cfg.Rows, m.cfg.Cols, m.boundary, m.pattern)
		}
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
	if !m.paused && m.game != nil {
		m.game.Step()
		m.currentStep = m.game.GetGeneration()
	}

	// Continue ticking regardless of pause state (for UI responsiveness)
	// but only if not quitting
	if !m.quitting {
		return m, tea.Tick(m.refreshRate, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}

	return m, nil
}

// View renders the current state
func (m Model) View() string {
	return m.RenderMode()
}

// renderTable renders a 2x3 table with given data and style
func renderTable(data TableData, containerStyle lipgloss.Style) string {
	if len(data.Rows) != 2 || len(data.Rows[0]) != 3 || len(data.Rows[1]) != 3 {
		return containerStyle.Render("Invalid table data")
	}

	// Render first row
	firstRow := lipgloss.JoinHorizontal(lipgloss.Top,
		tableCellStyle.Render(data.Rows[0][0]),
		tableCellStyle.Render(data.Rows[0][1]),
		tableCellStyle.Render(data.Rows[0][2]),
	)

	// Render second row
	secondRow := lipgloss.JoinHorizontal(lipgloss.Top,
		tableCellStyle.Render(data.Rows[1][0]),
		tableCellStyle.Render(data.Rows[1][1]),
		tableCellStyle.Render(data.Rows[1][2]),
	)

	// Join rows vertically
	tableContent := lipgloss.JoinVertical(lipgloss.Left, firstRow, secondRow)

	return containerStyle.Render(tableContent)
}

// createStatusTableData creates table data for status display
func (m Model) createStatusTableData() TableData {
	// Determine if using Chinese language
	isChinese := m.cfg.Language == Chinese

	// Get localized labels
	labels := getStatusLabels(isChinese)

	// Get status display string
	statusDisplay := m.getStatusDisplay(isChinese)

	// Get pattern display string based on language
	var patternDisplay string
	if isChinese {
		patternDisplay = m.pattern.ChineseString()
	} else {
		patternDisplay = m.pattern.String()
	}

	// Get boundary display string based on language
	var boundaryDisplay string
	if isChinese {
		boundaryDisplay = m.boundary.ChineseString()
	} else {
		boundaryDisplay = m.boundary.String()
	}

	return TableData{
		Rows: [][]string{
			{
				formatTableCell(generationIcon, labels.generation, fmt.Sprintf("%d", m.currentStep)),
				formatTableCell(speedIcon, labels.refresh, m.refreshRate.String()),
				formatTableCell("", labels.status, statusDisplay),
			},
			{
				formatTableCell(sizeIcon, labels.size, fmt.Sprintf("%d×%d", m.cfg.Rows, m.cfg.Cols)),
				formatTableCell(patternIcon, labels.pattern, patternDisplay),
				formatTableCell(boundaryIcon, labels.boundary, boundaryDisplay),
			},
		},
	}
}

// StatusLabels holds localized status labels
type StatusLabels struct {
	generation, refresh, size, cellSize, boundary, pattern, status string
}

// getStatusLabels returns localized status labels
func getStatusLabels(isChinese bool) StatusLabels {
	if isChinese {
		return StatusLabels{
			generation: statusGenerationLabelCN,
			refresh:    statusRefreshLabelCN,
			size:       statusSizeLabelCN,
			cellSize:   statusCellSizeLabelCN,
			boundary:   statusBoundaryLabelCN,
			pattern:    statusPatternLabelCN,
			status:     statusPausedLabelCN,
		}
	}
	return StatusLabels{
		generation: statusGenerationLabelEN,
		refresh:    statusRefreshLabelEN,
		size:       statusSizeLabelEN,
		cellSize:   statusCellSizeLabelEN,
		boundary:   statusBoundaryLabelEN,
		pattern:    statusPatternLabelEN,
		status:     statusPausedLabelEN,
	}
}

// getStatusDisplay returns the formatted status display string
func (m Model) getStatusDisplay(isChinese bool) string {
	if m.paused {
		if isChinese {
			return fmt.Sprintf("%s %s", pausedIcon, statusPausedCN)
		}
		return fmt.Sprintf("%s %s", pausedIcon, statusPausedEN)
	}

	if isChinese {
		return fmt.Sprintf("%s %s", playingIcon, statusRunningCN)
	}
	return fmt.Sprintf("%s %s", playingIcon, statusRunningEN)
}

// formatTableCell formats a table cell with icon, label, and value
func formatTableCell(icon, label, value string) string {
	if icon == "" {
		return fmt.Sprintf("%s: %s", tableLabelStyle.Render(label), tableValueStyle.Render(value))
	}
	return fmt.Sprintf("%s %s: %s", icon, tableLabelStyle.Render(label), tableValueStyle.Render(value))
}

// createControlTableData creates table data for control display
func (m Model) createControlTableData() TableData {
	if m.cfg.Language == Chinese {
		return TableData{
			Rows: [][]string{
				{
					fmt.Sprintf("%s %s", controlKeyStyle.Render("空格"), tableLabelStyle.Render("暂停/继续")),
					fmt.Sprintf("%s %s", controlKeyStyle.Render("+/-"), tableLabelStyle.Render("调节速度")),
					fmt.Sprintf("%s %s", controlKeyStyle.Render("Q"), tableLabelStyle.Render("退出")),
				},
				{
					fmt.Sprintf("%s %s", controlKeyStyle.Render("L"), tableLabelStyle.Render("切换语言")),
					fmt.Sprintf("%s %s", controlKeyStyle.Render("P"), tableLabelStyle.Render("切换模式")),
					fmt.Sprintf("%s %s", controlKeyStyle.Render("B"), tableLabelStyle.Render("边界类型")),
				},
			},
		}
	}
	return TableData{
		Rows: [][]string{
			{
				fmt.Sprintf("%s %s", controlKeyStyle.Render("Space"), tableLabelStyle.Render("Pause/Resume")),
				fmt.Sprintf("%s %s", controlKeyStyle.Render("+/-"), tableLabelStyle.Render("Speed")),
				fmt.Sprintf("%s %s", controlKeyStyle.Render("Q"), tableLabelStyle.Render("Quit")),
			},
			{
				fmt.Sprintf("%s %s", controlKeyStyle.Render("L"), tableLabelStyle.Render("Language")),
				fmt.Sprintf("%s %s", controlKeyStyle.Render("P"), tableLabelStyle.Render("Pattern")),
				fmt.Sprintf("%s %s", controlKeyStyle.Render("B"), tableLabelStyle.Render("Boundary")),
			},
		},
	}
}

// RenderGrid renders the 2D grid using optimized rendering
func (m *Model) RenderGrid() string {
	m.gridBuilder.Reset()

	// Safety checks
	if m.game == nil {
		return ""
	}

	grid := m.game.GetCurrentGrid()
	if len(grid) == 0 {
		return ""
	}

	// Pre-calculate styled strings to avoid repeated lookups
	var aliveStr, deadStr string
	if m.renderOptions != nil {
		aliveStr = m.renderOptions.aliveStyled
		deadStr = m.renderOptions.deadStyled
	} else {
		aliveStr = "█"
		deadStr = " "
	}

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

// RenderHeaderTitle returns the enhanced header title
func (m Model) RenderHeaderTitle() string {
	var title string
	if m.cfg.Language == Chinese {
		title = headerTitleFmtCN
	} else {
		title = headerTitleFmtEN
	}
	return headerStyle.Render(title)
}

// RenderStatusLine returns the status table (2x3)
func (m Model) RenderStatusLine() string {
	statusData := m.createStatusTableData()
	return renderTable(statusData, statusTableStyle.Inherit(tableStyle))
}

// RenderControlPanel returns the control table (2x3)
func (m Model) RenderControlPanel() string {
	controlData := m.createControlTableData()
	return renderTable(controlData, controlTableStyle.Inherit(tableStyle))
}

// RenderMode renders the complete UI mode view with enhanced layout
func (m Model) RenderMode() string {
	m.statusBuilder.Reset()

	// Build complete UI with enhanced styling
	m.statusBuilder.WriteString(m.RenderHeaderTitle())
	m.statusBuilder.WriteString("\n")
	m.statusBuilder.WriteString(m.RenderStatusLine())
	m.statusBuilder.WriteString("\n")
	m.statusBuilder.WriteString(m.RenderControlPanel())
	m.statusBuilder.WriteString("\n\n")
	m.statusBuilder.WriteString(m.RenderGrid())

	return m.statusBuilder.String()
}
