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
	ca            *CellularAutomaton
	cfg           *Config
	paused        bool // Pause state for infinite mode
	currentStep   int
	renderOptions *RenderOptions
	// Optimized string builders with pre-allocation
	gridBuilder   strings.Builder
	statusBuilder strings.Builder
	// Grid rendering optimization
	gridRingBuffer *GridRingBuffer
}

// TableData represents a 2x3 table structure
type TableData struct {
	Rows [][]string // 2 rows, 3 columns each
}

// GridRingBuffer and related functions are now in buffer.go

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

	ca := NewCellularAutomaton(cfg.Rule, cfg.Steps, cfg.Cols, cfg.Boundary)
	if ca == nil {
		// Handle case where cellular automaton creation fails
		configLogger.Printf("Failed to create cellular automaton, using default settings")
		ca = NewCellularAutomaton(DefaultRule, DefaultSteps, DefaultWindowCols, BoundaryPeriodic)
	}

	model := Model{
		ca:             ca,
		cfg:            cfg,
		renderOptions:  NewRenderOptions(cfg),
		gridRingBuffer: NewGridRingBuffer(cfg.Rows, cfg.Cols),
	}

	// Pre-allocate string builders with estimated capacity
	estimatedGridSize := cfg.Rows * (cfg.Cols*cfg.CellSize + 10) // +10 for formatting
	estimatedStatusSize := 500                                   // Estimated status line size

	model.gridBuilder.Grow(estimatedGridSize)
	model.statusBuilder.Grow(estimatedStatusSize)

	// Initialize the ring buffer with the initial state - add safety check
	if currentRow := ca.GetCurrentRow(); currentRow != nil {
		model.gridRingBuffer.AddRow(currentRow)
	}

	return model
}

// infiniteMsg is sent every tick for infinite mode
type tickMsg time.Time

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Only start ticking if we're in infinite mode or have steps to run
	if m.cfg.InfiniteMode || m.cfg.Steps > 0 {
		return tea.Tick(m.cfg.RefreshRate, func(t time.Time) tea.Msg {
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
		return m, tea.Quit

	case " ", "enter": // Space or Enter key for pause/resume in infinite mode
		m.paused = !m.paused

	case "l": // Language toggle key
		if m.cfg.Language == English {
			m.cfg.Language = Chinese
		} else {
			m.cfg.Language = English
		}

	case "r": // Reset simulation
		if m.ca != nil {
			m.ca.Reset()
			m.currentStep = 0
			// Reset ring buffer safely
			if m.gridRingBuffer != nil {
				m.gridRingBuffer.Clear()
			} else {
				m.gridRingBuffer = NewGridRingBuffer(m.cfg.Rows, m.cfg.Cols)
			}
			// Add initial row with safety check
			if currentRow := m.ca.GetCurrentRow(); currentRow != nil {
				m.gridRingBuffer.AddRow(currentRow)
			}
		}

	case "+", "=": // Increase refresh rate (make it faster)
		newRate := m.cfg.RefreshRate / 2
		if newRate >= time.Millisecond {
			m.cfg.RefreshRate = newRate
		}

	case "-", "_": // Decrease refresh rate (make it slower)
		m.cfg.RefreshRate = m.cfg.RefreshRate * 2

	case "1", "2", "3": // Change cell size
		if cellSize := int(msg.String()[0] - '0'); cellSize >= 1 && cellSize <= 3 {
			_ = m.cfg.SetCellSize(cellSize) // Error already logged
			m.renderOptions = NewRenderOptions(m.cfg)
			// Re-estimate and grow string builder capacity
			estimatedGridSize := m.cfg.Rows * (m.cfg.Cols*m.cfg.CellSize + 10)
			m.gridBuilder.Reset()
			m.gridBuilder.Grow(estimatedGridSize)
		}
	}
	return m, nil
}

// handleTick processes timer ticks
func (m Model) handleTick() (tea.Model, tea.Cmd) {
	// Check if we should continue running
	if !m.paused {
		if m.ca != nil && m.ca.Step() {
			m.currentStep = m.ca.GetGeneration()
			// Add new row to ring buffer with safety check
			if currentRow := m.ca.GetCurrentRow(); currentRow != nil && m.gridRingBuffer != nil {
				m.gridRingBuffer.AddRow(currentRow)
			}
		} else {
			// Cellular automaton has finished (reached max steps)
			// In finite mode, quit the program
			if !m.cfg.InfiniteMode {
				return m, tea.Quit
			}
		}
	}

	// Continue ticking only if we're not finished or in infinite mode
	if m.cfg.InfiniteMode || (m.ca != nil && !m.ca.IsFinished()) {
		return m, tea.Tick(m.cfg.RefreshRate, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}

	// In finite mode and finished, stop ticking but don't quit immediately
	// Let the user see the final result
	return m, nil
}

// View renders the current state
func (m Model) View() string {
	return m.RenderMode()
}

// RenderOptions, Language, and styling constants are now in styles.go

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

	return TableData{
		Rows: [][]string{
			{
				formatTableCell(generationIcon, labels.generation, fmt.Sprintf("%d", m.currentStep)),
				formatTableCell(speedIcon, labels.refresh, m.cfg.RefreshRate.String()),
				formatTableCell("", labels.status, statusDisplay),
			},
			{
				formatTableCell(sizeIcon, labels.size, fmt.Sprintf("%d×%d", m.cfg.Rows, m.cfg.Cols)),
				formatTableCell(cellIcon, labels.cellSize, fmt.Sprintf("%d", m.cfg.CellSize)),
				formatTableCell(boundaryIcon, labels.boundary, m.cfg.Boundary.String()),
			},
		},
	}
}

// StatusLabels holds localized status labels
type StatusLabels struct {
	generation, refresh, size, cellSize, boundary, status string
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
			status:     statusPausedLabelCN,
		}
	}
	return StatusLabels{
		generation: statusGenerationLabelEN,
		refresh:    statusRefreshLabelEN,
		size:       statusSizeLabelEN,
		cellSize:   statusCellSizeLabelEN,
		boundary:   statusBoundaryLabelEN,
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

	if !m.cfg.InfiniteMode && m.ca != nil && m.ca.IsFinished() {
		if isChinese {
			return fmt.Sprintf("%s %s", finishedIcon, statusFinishedCN)
		}
		return fmt.Sprintf("%s %s", finishedIcon, statusFinishedEN)
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
					fmt.Sprintf("%s %s", controlKeyStyle.Render("L"), tableLabelStyle.Render("切换语言")),
					fmt.Sprintf("%s %s", controlKeyStyle.Render("Q"), tableLabelStyle.Render("退出")),
				},
				{
					fmt.Sprintf("%s %s", controlKeyStyle.Render("R"), tableLabelStyle.Render("重置")),
					fmt.Sprintf("%s %s", controlKeyStyle.Render("+/-"), tableLabelStyle.Render("调节速度")),
					fmt.Sprintf("%s %s", controlKeyStyle.Render("1-3"), tableLabelStyle.Render("元胞大小")),
				},
			},
		}
	}
	return TableData{
		Rows: [][]string{
			{
				fmt.Sprintf("%s %s", controlKeyStyle.Render("Space"), tableLabelStyle.Render("Pause/Resume")),
				fmt.Sprintf("%s %s", controlKeyStyle.Render("L"), tableLabelStyle.Render("Language")),
				fmt.Sprintf("%s %s", controlKeyStyle.Render("Q"), tableLabelStyle.Render("Quit")),
			},
			{
				fmt.Sprintf("%s %s", controlKeyStyle.Render("R"), tableLabelStyle.Render("Reset")),
				fmt.Sprintf("%s %s", controlKeyStyle.Render("+/-"), tableLabelStyle.Render("Speed")),
				fmt.Sprintf("%s %s", controlKeyStyle.Render("1-3"), tableLabelStyle.Render("Cell Size")),
			},
		},
	}
}

// RenderGrid renders the grid using the optimized ring buffer with performance optimizations
func (m *Model) RenderGrid() string {
	m.gridBuilder.Reset()

	// Safety checks
	if m.gridRingBuffer == nil {
		return ""
	}

	rows := m.gridRingBuffer.GetRows()
	if len(rows) == 0 {
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
	lastRowIndex := len(rows) - 1
	for i, row := range rows {
		if row == nil {
			continue // Skip nil rows
		}

		m.gridBuilder.WriteString("\t")

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
		title = fmt.Sprintf(headerTitleFmtCN, m.cfg.Rule)
	} else {
		title = fmt.Sprintf(headerTitleFmtEN, m.cfg.Rule)
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
	m.statusBuilder.WriteString("\n")
	m.statusBuilder.WriteString(m.RenderGrid())

	return m.statusBuilder.String()
}
