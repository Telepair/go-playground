package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
			Align(lipgloss.Center)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#4A5568")).
			Padding(0, 1).
			Bold(true)

	tableBuilder strings.Builder
)

// UI text constants with enhanced formatting and icons
const (
	// Header Line
	HeaderCN = "🧬 元胞自动机 🧬"
	HeaderEN = "🧬 Cellular Automaton 🧬"

	// Status Line
	RuleLabelCN = "🧬 规则: %d"
	RuleLabelEN = "🧬 Rule: %d"

	GenerationLabelCN = "⚡ 代数: %d"
	GenerationLabelEN = "⚡ Gen: %d"

	SpeedLabelCN = "🔄 刷新: %s"
	SpeedLabelEN = "🔄 Speed: %s"

	SizeLabelCN = "📐 尺寸: %d×%d"
	SizeLabelEN = "📐 Size: %d×%d"

	BoundaryLabelCN = "🔒 边界: %s"
	BoundaryLabelEN = "🔒 Boundary: %s"

	StatusLabelPlayingCN = "▶️ 运行中"
	StatusLabelPlayingEN = "▶️ Running"
	StatusLabelPausedCN  = "⏸️ 已暂停"
	StatusLabelPausedEN  = "⏸️ Paused"

	// Control Line
	SelectRuleLabelCN = "T 选择规则"
	SelectRuleLabelEN = "T Select Rule"

	SelectBoundaryLabelCN = "B 选择边界"
	SelectBoundaryLabelEN = "B Select Boundary"

	SpeedControlLabelCN = "+/- 加速/减速"
	SpeedControlLabelEN = "+/- Speed Up/Down"

	LanguageLabelCN = "L 切换语言"
	LanguageLabelEN = "L Switch Language"

	SpaceLabelCN = "Space 暂停"
	SpaceLabelEN = "Space Pause"

	ResetLabelCN = "R 重置"
	ResetLabelEN = "R Reset"

	QuitLabelCN = "Q 退出"
	QuitLabelEN = "Q Quit"
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

// HeaderLineView returns the header display string
func (m Model) HeaderLineView() string {
	// Dynamically set the width of the header to the screen width for centering.
	style := headerStyle.Width(m.width)
	if m.language == Chinese {
		return style.Render(HeaderCN)
	}
	return style.Render(HeaderEN)
}

// StatusLineView returns the status display string for the first row
func (m Model) StatusLineView() string {
	var status, ruleLabel, generationLabel, speedLabel, boundaryLabel, sizeLabel string

	if m.language == Chinese {
		status = StatusLabelPlayingCN
		if m.paused {
			status = StatusLabelPausedCN
		}
		ruleLabel = RuleLabelCN
		generationLabel = GenerationLabelCN
		speedLabel = SpeedLabelCN
		boundaryLabel = BoundaryLabelCN
		sizeLabel = SizeLabelCN
	} else {
		status = StatusLabelPlayingEN
		if m.paused {
			status = StatusLabelPausedEN
		}
		ruleLabel = RuleLabelEN
		generationLabel = GenerationLabelEN
		speedLabel = SpeedLabelEN
		boundaryLabel = BoundaryLabelEN
		sizeLabel = SizeLabelEN
	}

	tableBuilder.Reset()
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(ruleLabel, m.rule)))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(generationLabel, m.currentStep)))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(speedLabel, m.refreshRate.String())))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(boundaryLabel, m.boundary.ToString(m.language))))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(sizeLabel, m.height, m.width)))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(status))

	return lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(tableBuilder.String())
}

// ControlLineView returns the control display string: T,B,R + Space, L, Q
func (m Model) ControlLineView() string {
	var selectRule, selectBoundary, speedControl, language, space, reset, quit string
	if m.language == Chinese {
		selectRule = SelectRuleLabelCN
		selectBoundary = SelectBoundaryLabelCN
		speedControl = SpeedControlLabelCN
		language = LanguageLabelCN
		space = SpaceLabelCN
		reset = ResetLabelCN
		quit = QuitLabelCN
	} else {
		selectRule = SelectRuleLabelEN
		selectBoundary = SelectBoundaryLabelEN
		speedControl = SpeedControlLabelEN
		language = LanguageLabelEN
		space = SpaceLabelEN
		reset = ResetLabelEN
		quit = QuitLabelEN
	}

	tableBuilder.Reset()
	tableBuilder.WriteString(labelStyle.Render(selectRule))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(selectBoundary))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(speedControl))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(language))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(space))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(reset))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(quit))

	return lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(tableBuilder.String())
}
