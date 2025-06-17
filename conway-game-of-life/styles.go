package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
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
const (
	// Header Line
	HeaderEN = "🎮 Conway's Game of Life 🎮"
	HeaderCN = "🎮 康威生命游戏 🎮"

	// Status Line
	GenerationIcon    = "⚡"
	GenerationLabelCN = "代数"
	GenerationLabelEN = "Gen"

	SpeedIcon    = "🔄"
	SpeedLabelCN = "刷新"
	SpeedLabelEN = "Speed"

	SizeIcon    = "📐"
	SizeLabelCN = "尺寸"
	SizeLabelEN = "Size"

	BoundaryIcon    = "🔒"
	BoundaryLabelCN = "边界"
	BoundaryLabelEN = "Boundary"

	PatternIcon    = "🎨"
	PatternLabelCN = "模式"
	PatternLabelEN = "Pattern"

	StatusLabelCN = "状态"
	StatusLabelEN = "Status"
	PlayingIcon   = "▶️"
	PlayingCN     = "运行中"
	PlayingEN     = "Running"
	PausedIcon    = "⏸️"
	PausedEN      = "Paused"
	PausedCN      = "已暂停"

	ResetKey     = "R"
	ResetLabelCN = "重置"
	ResetLabelEN = "Reset"

	SpeedControlKey     = "+/-"
	SpeedControlLabelCN = "加速/减速"
	SpeedControlLabelEN = "Speed Up/Down"

	QuitKey     = "Space/Q"
	QuitLabelCN = "暂停/退出"
	QuitLabelEN = "Pause/Quit"

	LanguageKey     = "L"
	LanguageLabelCN = "切换语言"
	LanguageLabelEN = "Switch Language"

	SelectPatternKey     = "P"
	SelectPatternLabelCN = "选择模式"
	SelectPatternLabelEN = "Select Pattern"

	SelectBoundaryKey     = "B"
	SelectBoundaryLabelCN = "选择边界"
	SelectBoundaryLabelEN = "Select Boundary"
)

// RenderOptions contains rendering configuration with cached styles
type RenderOptions struct {
	aliveStyled string // Cached styled alive cell
	deadStyled  string // Cached styled dead cell
}

// NewRenderOptions creates optimized render options with pre-computed styles
func NewRenderOptions(aliveColor, deadColor, aliveChar, deadChar string) RenderOptions {
	return RenderOptions{
		aliveStyled: lipgloss.NewStyle().Foreground(lipgloss.Color(aliveColor)).Render(aliveChar),
		deadStyled:  lipgloss.NewStyle().Foreground(lipgloss.Color(deadColor)).Render(deadChar),
	}
}

// formatTableCell formats a table cell with icon, label, and value
func formatStatus(icon, label, value string) string {
	return tableCellStyle.Render(fmt.Sprintf("%s %s: %s", icon, tableLabelStyle.Render(label), tableValueStyle.Render(value)))
}

func formatControl(key, label string) string {
	return tableCellStyle.Render(fmt.Sprintf("%s %s", controlKeyStyle.Render(key), tableLabelStyle.Render(label)))
}

// GetHeaderLine returns the header display string
func GetHeaderLine(language Language) string {
	style := headerStyle.Inherit(tableStyle)
	if language == Chinese {
		return style.Render(HeaderCN)
	}
	return style.Render(HeaderEN)
}

// GetStatusLine returns the status display string for the first row
func GetStatusLine(language Language, pattern Pattern, generation int, speed time.Duration, rows, cols int, boundary BoundaryType, paused bool) string {
	style := statusTableStyle.Inherit(tableStyle)
	if language == Chinese {
		status := PlayingCN
		statusIcon := PlayingIcon
		if paused {
			status = PausedCN
			statusIcon = PausedIcon
		}
		tableContent := lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top,
				formatStatus(GenerationIcon, GenerationLabelCN, fmt.Sprintf("%d", generation)),
				formatStatus(SpeedIcon, SpeedLabelCN, speed.String()),
				formatStatus(SizeIcon, SizeLabelCN, fmt.Sprintf("%d×%d", rows, cols)),
			),
			lipgloss.JoinHorizontal(lipgloss.Top,
				formatStatus(BoundaryIcon, BoundaryLabelCN, boundary.ToString(language)),
				formatStatus(PatternIcon, PatternLabelCN, pattern.ToString(language)),
				formatStatus(statusIcon, StatusLabelCN, status),
			),
		)
		return style.Render(tableContent)
	}

	status := PlayingEN
	statusIcon := PlayingIcon
	if paused {
		status = PausedEN
		statusIcon = PausedIcon
	}

	tableContent := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top,
			formatStatus(GenerationIcon, GenerationLabelEN, fmt.Sprintf("%d", generation)),
			formatStatus(SpeedIcon, SpeedLabelEN, speed.String()),
			formatStatus(SizeIcon, SizeLabelEN, fmt.Sprintf("%d×%d", rows, cols)),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			formatStatus(BoundaryIcon, BoundaryLabelEN, boundary.ToString(language)),
			formatStatus(PatternIcon, PatternLabelEN, pattern.ToString(language)),
			formatStatus(statusIcon, StatusLabelEN, status),
		),
	)
	return style.Render(tableContent)
}

// GetControlLine returns the control display string: T,B,R + Space, L, Q
func GetControlLine(language Language) string {
	style := controlTableStyle.Inherit(tableStyle)
	if language == Chinese {
		tableContent := lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top,
				formatControl(ResetKey, ResetLabelCN),
				formatControl(SpeedControlKey, SpeedControlLabelCN),
				formatControl(QuitKey, QuitLabelCN),
			),
			lipgloss.JoinHorizontal(lipgloss.Top,
				formatControl(LanguageKey, LanguageLabelCN),
				formatControl(SelectPatternKey, SelectPatternLabelCN),
				formatControl(SelectBoundaryKey, SelectBoundaryLabelCN),
			),
		)
		return style.Render(tableContent)
	}

	tableContent := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top,
			formatControl(ResetKey, ResetLabelEN),
			formatControl(SpeedControlKey, SpeedControlLabelEN),
			formatControl(QuitKey, QuitLabelEN),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			formatControl(LanguageKey, LanguageLabelEN),
			formatControl(SelectPatternKey, SelectPatternLabelEN),
			formatControl(SelectBoundaryKey, SelectBoundaryLabelEN),
		),
	)
	return style.Render(tableContent)
}
