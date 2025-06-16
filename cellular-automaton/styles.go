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
	HeaderCN = "🧬 元胞自动机 🧬"
	HeaderEN = "🧬 Cellular Automaton 🧬"

	// Status Line
	RuleIcon    = "🧬"
	RuleLabelCN = "规则"
	RuleLabelEN = "Rule"

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

	PausedLabelCN = "状态"
	PausedLabelEN = "Status"
	PlayingIcon   = "▶️"
	PlayingCN     = "运行"
	PlayingEN     = "Running"
	PausedIcon    = "⏸️"
	PausedEN      = "Paused"
	PausedCN      = "已暂停"

	// Control Line
	SelectRuleKey     = "T"
	SelectRuleLabelCN = "选择规则"
	SelectRuleLabelEN = "Select Rule"

	SelectBoundaryKey     = "B"
	SelectBoundaryLabelCN = "选择边界"
	SelectBoundaryLabelEN = "Select Boundary"

	SpeedControlKey     = "+/-"
	SpeedControlLabelCN = "加速/减速"
	SpeedControlLabelEN = "Speed Up/Down"

	ResetKey     = "R"
	ResetLabelCN = "重置"
	ResetLabelEN = "Reset"

	LanguageKey     = "L"
	LanguageLabelCN = "切换语言"
	LanguageLabelEN = "Switch Language"

	QuitKey     = "Space/Q"
	QuitLabelCN = "暂停/退出"
	QuitLabelEN = "Pause/Quit"
)

// RenderOptions contains rendering configuration with cached styles
type RenderOptions struct {
	aliveStyled string         // Cached styled alive cell
	deadStyled  string         // Cached styled dead cell
	aliveStyle  lipgloss.Style // Cached alive style
	deadStyle   lipgloss.Style // Cached dead style
}

// NewRenderOptions creates optimized render options with pre-computed styles
func NewRenderOptions(aliveColor, deadColor, aliveChar, deadChar string) RenderOptions {
	return RenderOptions{
		aliveStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color(aliveColor)),
		deadStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color(deadColor)),
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
func GetStatusLine(language Language, rule int, generation int, speed time.Duration, rows, cols int, boundary BoundaryType, paused bool) string {
	style := statusTableStyle.Inherit(tableStyle)
	if language == Chinese {
		pausedStatus := PlayingCN
		if paused {
			pausedStatus = PausedCN
		}
		tableContent := lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top,
				formatStatus(RuleIcon, RuleLabelCN, fmt.Sprintf("%d", rule)),
				formatStatus(GenerationIcon, GenerationLabelCN, fmt.Sprintf("%d", generation)),
				formatStatus(SpeedIcon, SpeedLabelCN, speed.String()),
			),
			lipgloss.JoinHorizontal(lipgloss.Top,
				formatStatus(BoundaryIcon, BoundaryLabelCN, boundary.ToString(language)),
				formatStatus(SizeIcon, SizeLabelCN, fmt.Sprintf("%d×%d", rows, cols)),
				formatStatus(PausedIcon, PausedLabelCN, pausedStatus),
			),
		)
		return style.Render(tableContent)
	}

	pausedStatus := PlayingEN
	if paused {
		pausedStatus = PausedEN
	}

	tableContent := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top,
			formatStatus(RuleIcon, RuleLabelEN, fmt.Sprintf("%d", rule)),
			formatStatus(GenerationIcon, GenerationLabelEN, fmt.Sprintf("%d", generation)),
			formatStatus(SpeedIcon, SpeedLabelEN, speed.String()),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			formatStatus(BoundaryIcon, BoundaryLabelEN, boundary.ToString(language)),
			formatStatus(SizeIcon, SizeLabelEN, fmt.Sprintf("%d×%d", rows, cols)),
			formatStatus(PausedIcon, PausedLabelEN, pausedStatus),
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
				formatControl(SelectRuleKey, SelectRuleLabelCN),
				formatControl(SelectBoundaryKey, SelectBoundaryLabelCN),
				formatControl(SpeedControlKey, SpeedControlLabelCN),
			),
			lipgloss.JoinHorizontal(lipgloss.Top,
				formatControl(ResetKey, ResetLabelCN),
				formatControl(LanguageKey, LanguageLabelCN),
				formatControl(QuitKey, QuitLabelCN),
			),
		)
		return style.Render(tableContent)
	}

	tableContent := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top,
			formatControl(SelectRuleKey, SelectRuleLabelEN),
			formatControl(SelectBoundaryKey, SelectBoundaryLabelEN),
			formatControl(SpeedControlKey, SpeedControlLabelEN),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			formatControl(ResetKey, ResetLabelEN),
			formatControl(LanguageKey, LanguageLabelEN),
			formatControl(QuitKey, QuitLabelEN),
		),
	)
	return style.Render(tableContent)
}
