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
	err           error
	currentStep   uint
	renderOptions *RenderOptions
}

func NewModel(cfg *Config) Model {
	ca := NewCellularAutomaton(cfg.Rule, cfg.Steps, &cfg.Window, cfg.Boundary)
	return Model{
		ca:            ca,
		cfg:           cfg,
		renderOptions: NewRenderOptions(cfg),
	}
}

// infiniteMsg is sent every tick for infinite mode
type infiniteMsg time.Time

// finiteMsg is sent every tick for finite mode
type finiteMsg time.Time

// Init initializes the model
func (m Model) Init() tea.Cmd {
	if m.cfg.InfiniteMode {
		return tea.Tick(m.cfg.RefreshRate, func(t time.Time) tea.Msg {
			return infiniteMsg(t)
		})
	}
	// For finite mode, start computing steps
	return tea.Tick(time.Millisecond*FiniteModeTickRate, func(t time.Time) tea.Msg {
		return finiteMsg(t)
	})
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// Show final screen before quitting
			return m, tea.Quit
		case " ", "enter": // Space or Enter key for pause/resume in infinite mode
			if m.cfg.InfiniteMode {
				m.paused = !m.paused
			} else {
				return m, tea.Quit
			}
		case "l": // Language toggle key
			if m.cfg.Language == English {
				m.cfg.Language = Chinese
			} else {
				m.cfg.Language = English
			}
		case "r": // Reset simulation
			m.ca = NewCellularAutomaton(m.cfg.Rule, m.cfg.Steps, &m.cfg.Window, m.cfg.Boundary)
			m.currentStep = 0
		case "+", "=": // Increase refresh rate (make it faster)
			if m.cfg.InfiniteMode {
				newRate := m.cfg.RefreshRate / 2
				if newRate >= time.Millisecond {
					m.cfg.RefreshRate = newRate
				}
			}
		case "-", "_": // Decrease refresh rate (make it slower)
			if m.cfg.InfiniteMode {
				newRate := m.cfg.RefreshRate * 2
				if newRate <= time.Second*5 { // Max 5 seconds
					m.cfg.RefreshRate = newRate
				}
			}
		case "1", "2", "3": // Change cell size
			if cellSize := uint(msg.String()[0] - '0'); cellSize >= 1 && cellSize <= 3 {
				m.cfg.CellSize = cellSize
				m.renderOptions = NewRenderOptions(m.cfg)
			}
		}

	case infiniteMsg:
		if !m.paused && m.ca.Step() {
			m.currentStep = m.ca.GetGeneration()
		}
		// Still send tick even when paused to maintain the loop
		return m, tea.Tick(m.cfg.RefreshRate, func(t time.Time) tea.Msg {
			return infiniteMsg(t)
		})

	case finiteMsg:
		if m.ca.Step() {
			m.currentStep = m.ca.GetGeneration()
		}
		return m, tea.Tick(time.Millisecond*FiniteModeTickRate, func(t time.Time) tea.Msg {
			return finiteMsg(t)
		})
	}

	return m, nil
}

// View renders the current state
func (m Model) View() string {
	return RenderMode(m)
}

// RenderOptions contains rendering configuration
type RenderOptions struct {
	CellSize   uint   // Size of each cell (1-3)
	AliveColor string // Color for alive cells
	DeadColor  string // Color for dead cells
	AliveChar  string // Character for alive cells
	DeadChar   string // Character for dead cells
	// Add cached styled strings for better performance
	aliveStyled string         // Cached styled alive cell
	deadStyled  string         // Cached styled dead cell
	aliveStyle  lipgloss.Style // Cached alive style
	deadStyle   lipgloss.Style // Cached dead style
}

// NewRenderOptions creates default render options
func NewRenderOptions(cfg *Config) *RenderOptions {
	if cfg == nil {
		return &RenderOptions{
			CellSize:   DefaultCellSize,
			AliveColor: DefaultAliveColor,
			DeadColor:  DefaultDeadColor,
			AliveChar:  DefaultAliveChar,
			DeadChar:   DefaultDeadChar,
		}
	}

	options := &RenderOptions{
		CellSize:   cfg.CellSize,
		AliveColor: cfg.AliveColor,
		DeadColor:  cfg.DeadColor,
		AliveChar:  cfg.AliveChar,
		DeadChar:   cfg.DeadChar,
	}

	// Pre-compute styles and styled strings for better performance
	options.aliveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(options.AliveColor))
	options.deadStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(options.DeadColor))

	aliveCell := strings.Repeat(options.AliveChar, int(options.CellSize))
	deadCell := strings.Repeat(options.DeadChar, int(options.CellSize))

	options.aliveStyled = options.aliveStyle.Render(aliveCell)
	options.deadStyled = options.deadStyle.Render(deadCell)

	return options
}

// Language represents the supported languages
type Language int

const (
	English Language = iota
	Chinese
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(HeaderFgColor)).
			Background(lipgloss.Color(HeaderBgColor)).
			Padding(0, 1)

	headerTitleFmt         = "\t元胞自动机 - 规则 %d"
	statusLineFmt          = "\t代数: %d | 刷新: %.3fs | 窗口: %d x %d | 元胞大小: %d | 边界: %s"
	pausedStatusLineFmt    = "\t代数: %d | 刷新: %.3fs | 窗口: %d x %d | 元胞大小: %d | 边界: %s | 已暂停"
	finiteStatusLineFmt    = "\t代数: %d/%d | 刷新: %.3fs | 窗口: %d x %d | 元胞大小: %d | 边界: %s"
	controlLineFmt         = "\t按键: [空格/回车]暂停 | [l]切换语言 | [q]退出"
	finiteControlLineFmt   = "\t按键: [l]切换语言 | [空格/回车/q]退出"
	extendedControlLineFmt = "\t高级: [r]重置 | [+/-]调速度 | [1-3]元胞大小"
	errorPrefixFmt         = "错误: %v"

	headerTitleFmtEn         = "\tCellular Automaton - Rule %d"
	statusLineFmtEn          = "\tGeneration: %d | Refresh: %.3fs | Window: %d x %d | Cell Size: %d | Boundary: %s"
	pausedStatusLineFmtEn    = "\tGeneration: %d | Refresh: %.3fs | Window: %d x %d | Cell Size: %d | Boundary: %s | Paused"
	finiteStatusLineFmtEn    = "\tGeneration: %d/%d | Refresh: %.3fs | Window: %d x %d | Cell Size: %d | Boundary: %s"
	controlLineFmtEn         = "\tControls: [Space/Enter]Pause | [l]Language | [q]Quit"
	finiteControlLineFmtEn   = "\tControls: [l]Language | [Space/Enter/q]Quit"
	extendedControlLineFmtEn = "\tAdvanced: [r]Reset | [+/-]Speed | [1-3]Cell Size"
	errorPrefixFmtEn         = "Error: %v"
)

// RenderGrid renders the grid as a string using specified render options
func RenderGrid(grid [][]bool, options *RenderOptions) string {
	if len(grid) == 0 {
		return ""
	}

	if options == nil {
		options = NewRenderOptions(nil)
	}

	// Pre-allocate string builder with estimated capacity
	estimatedSize := len(grid) * (len(grid[0])*int(options.CellSize) + 1) // +1 for newline
	var result strings.Builder
	result.Grow(estimatedSize)

	for _, row := range grid {
		result.WriteString("\t")
		for _, cell := range row {
			if cell {
				result.WriteString(options.aliveStyled)
			} else {
				result.WriteString(options.deadStyled)
			}
		}
		result.WriteByte('\n')
	}
	return result.String()
}

// RenderMode renders the mode view
func RenderMode(m Model) string {
	if m.err != nil {
		return fmt.Sprintf("%s: %v\n", RenderErrorPrefix(m.cfg), m.err)
	}

	var s strings.Builder

	// Header
	s.WriteString(RenderHeaderTitle(m.cfg) + "\n\n")
	s.WriteString(RenderStatusLine(m.cfg, m.ca.GetGeneration(), m.paused) + "\n")
	s.WriteString(RenderControlLine(m.cfg) + "\n")
	s.WriteString(RenderExtendedControlLine(m.cfg) + "\n\n")

	// Grid visualization
	s.WriteString(RenderGrid(m.ca.GetGrid(), m.renderOptions))

	return s.String()
}

// RenderHeaderTitle returns the header title for the current mode
func RenderHeaderTitle(cfg *Config) string {
	if cfg.Language == Chinese {
		return headerStyle.Render(fmt.Sprintf(headerTitleFmt, cfg.Rule))
	}
	return headerStyle.Render(fmt.Sprintf(headerTitleFmtEn, cfg.Rule))
}

// RenderStatusLine returns the status line for the current mode
func RenderStatusLine(cfg *Config, currentStep uint, paused bool) string {
	if cfg.Language == Chinese {
		if cfg.InfiniteMode {
			if paused {
				return fmt.Sprintf(pausedStatusLineFmt, currentStep, cfg.RefreshRate.Seconds(), cfg.Window.cols, cfg.Window.rows, cfg.CellSize, cfg.Boundary)
			}
			return fmt.Sprintf(statusLineFmt, currentStep, cfg.RefreshRate.Seconds(), cfg.Window.cols, cfg.Window.rows, cfg.CellSize, cfg.Boundary)
		}
		return fmt.Sprintf(finiteStatusLineFmt, currentStep, cfg.Steps, cfg.RefreshRate.Seconds(), cfg.Window.cols, cfg.Window.rows, cfg.CellSize, cfg.Boundary)
	}

	if cfg.InfiniteMode {
		if paused {
			return fmt.Sprintf(pausedStatusLineFmtEn, currentStep, cfg.RefreshRate.Seconds(), cfg.Window.cols, cfg.Window.rows, cfg.CellSize, cfg.Boundary)
		}
		return fmt.Sprintf(statusLineFmtEn, currentStep, cfg.RefreshRate.Seconds(), cfg.Window.cols, cfg.Window.rows, cfg.CellSize, cfg.Boundary)
	}
	return fmt.Sprintf(finiteStatusLineFmtEn, currentStep, cfg.Steps, cfg.RefreshRate.Seconds(), cfg.Window.cols, cfg.Window.rows, cfg.CellSize, cfg.Boundary)
}

// RenderControlLine returns the control line for the current mode
func RenderControlLine(cfg *Config) string {
	if cfg.Language == Chinese {
		if cfg.InfiniteMode {
			return controlLineFmt
		}
		return finiteControlLineFmt
	}
	if cfg.InfiniteMode {
		return controlLineFmtEn
	}
	return finiteControlLineFmtEn
}

// RenderExtendedControlLine returns the extended control line for the current mode
func RenderExtendedControlLine(cfg *Config) string {
	if cfg.Language == Chinese {
		return extendedControlLineFmt
	}
	return extendedControlLineFmtEn
}

// RenderErrorPrefix returns the error prefix for the current mode
func RenderErrorPrefix(cfg *Config) string {
	if cfg.Language == Chinese {
		return errorPrefixFmt
	}
	return errorPrefixFmtEn
}
