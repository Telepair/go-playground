package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Language represents the supported languages
type Language int

// Language constants
const (
	English Language = iota
	Chinese
)

// UI layout constants for consistent alignment
const (
	PanelWidth     = 90 // Standard width for all UI panels
	TableCellWidth = 28 // Each cell takes 1/3 of panel width (90/3 ≈ 30, minus padding)
)

// Enhanced UI styles for better visual appearance
var (
	// Header styles
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#874BFD")).
			Padding(0, 2).
			MarginBottom(1).
			Align(lipgloss.Center).
			Width(PanelWidth)

	controlKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#4A5568")).
			Padding(0, 1).
			Bold(true)

	// Table styles for metadata display
	tableStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4A5568")).
			Width(PanelWidth).
			MarginBottom(1)

	statusTableStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#2D3748")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(1, 2)

	controlTableStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#1A202C")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(1, 2)

	// Table cell styles
	tableCellStyle = lipgloss.NewStyle().
			Width(TableCellWidth).
			Align(lipgloss.Left).
			Padding(0, 1)

	tableLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E2E8F0")).
			Bold(false)

	tableValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#68D391")).
			Bold(true)
)

// UI text constants with enhanced formatting and icons
var (
	// Status icons and labels
	generationIcon = "⚡"
	speedIcon      = "🔄"
	sizeIcon       = "📐"
	cellIcon       = "🔹"
	boundaryIcon   = "🔒"
	pausedIcon     = "⏸️"
	playingIcon    = "▶️"
	finishedIcon   = "✅"

	// Chinese text templates with enhanced formatting
	headerTitleFmtCN = "🧬 元胞自动机 - 规则 %d 🧬"

	// Status line components
	statusGenerationLabelCN = "代数"
	statusRefreshLabelCN    = "刷新"
	statusSizeLabelCN       = "尺寸"
	statusCellSizeLabelCN   = "元胞"
	statusBoundaryLabelCN   = "边界"
	statusPausedLabelCN     = "状态"

	// English text templates with enhanced formatting
	headerTitleFmtEN = "🧬 Cellular Automaton - Rule %d 🧬"

	// Status line components
	statusGenerationLabelEN = "Gen"
	statusRefreshLabelEN    = "Speed"
	statusSizeLabelEN       = "Size"
	statusCellSizeLabelEN   = "Cell"
	statusBoundaryLabelEN   = "Boundary"
	statusPausedLabelEN     = "Status"

	// Status text messages
	statusPausedCN   = "已暂停"
	statusRunningCN  = "运行中"
	statusFinishedCN = "已完成"
	statusPausedEN   = "PAUSED"
	statusRunningEN  = "RUNNING"
	statusFinishedEN = "FINISHED"
)

// RenderOptions contains rendering configuration with cached styles
type RenderOptions struct {
	CellSize   int    // Size of each cell (1-3)
	AliveColor string // Color for alive cells
	DeadColor  string // Color for dead cells
	AliveChar  string // Character for alive cells
	DeadChar   string // Character for dead cells
	// Cached styled strings for better performance
	aliveStyled string         // Cached styled alive cell
	deadStyled  string         // Cached styled dead cell
	aliveStyle  lipgloss.Style // Cached alive style
	deadStyle   lipgloss.Style // Cached dead style
}

// NewRenderOptions creates optimized render options with pre-computed styles
func NewRenderOptions(cfg *Config) *RenderOptions {
	if cfg == nil {
		cfg = NewConfig() // Use default config if nil
	}

	options := &RenderOptions{
		CellSize:   cfg.CellSize,
		AliveColor: cfg.AliveColor,
		DeadColor:  cfg.DeadColor,
		AliveChar:  cfg.AliveChar,
		DeadChar:   cfg.DeadChar,
	}

	// Pre-compute styles for better performance
	options.aliveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(options.AliveColor))
	options.deadStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(options.DeadColor))

	// Pre-compute styled strings with repeated characters
	aliveCell := strings.Repeat(options.AliveChar, options.CellSize)
	deadCell := strings.Repeat(options.DeadChar, options.CellSize)

	options.aliveStyled = options.aliveStyle.Render(aliveCell)
	options.deadStyled = options.deadStyle.Render(deadCell)

	return options
}
