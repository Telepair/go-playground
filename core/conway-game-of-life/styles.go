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
	HeaderCN = "🎮 康威生命游戏 🎮"
	HeaderEN = "🎮 Conway's Game of Life 🎮"

	// Status Line
	GenerationLabelCN = "⚡ 代数: %d"
	GenerationLabelEN = "⚡ Gen: %d"

	SpeedLabelCN = "🔄 刷新: %s"
	SpeedLabelEN = "🔄 Speed: %s"

	SizeLabelCN = "📐 尺寸: %d×%d"
	SizeLabelEN = "📐 Size: %d×%d"

	BoundaryLabelCN = "🔒 边界: %s"
	BoundaryLabelEN = "🔒 Boundary: %s"

	PatternLabelCN = "🎨 模式: %s"
	PatternLabelEN = "🎨 Pattern: %s"

	StatusLabelPlayingCN = "▶️ 运行中"
	StatusLabelPlayingEN = "▶️ Running"
	StatusLabelPausedCN  = "⏸️ 已暂停"
	StatusLabelPausedEN  = "⏸️ Paused"

	// Control Line
	SelectPatternLabelCN = "P 选择模式"
	SelectPatternLabelEN = "P Select Pattern"

	SelectBoundaryLabelCN = "B 选择边界"
	SelectBoundaryLabelEN = "B Select Boundary"

	LanguageLabelCN = "L 切换语言"
	LanguageLabelEN = "L Switch Language"

	SpeedControlLabelCN = "+/- 加速/减速"
	SpeedControlLabelEN = "+/- Speed Up/Down"

	SpaceControlLabelCN = "Space 暂停"
	SpaceControlLabelEN = "Space Pause"

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
	style := headerStyle.Width(m.width)
	if m.language == Chinese {
		return style.Render(HeaderCN)
	}
	return style.Render(HeaderEN)
}

// StatusLineView returns the status display string for the first row
func (m Model) StatusLineView() string {
	var status, generationLabel, speedLabel, boundaryLabel, sizeLabel, patternLabel string

	if m.language == Chinese {
		status = StatusLabelPlayingCN
		if m.paused {
			status = StatusLabelPausedCN
		}
		generationLabel = GenerationLabelCN
		speedLabel = SpeedLabelCN
		sizeLabel = SizeLabelCN
		boundaryLabel = BoundaryLabelCN
		patternLabel = PatternLabelCN
	} else {
		status = StatusLabelPlayingEN
		if m.paused {
			status = StatusLabelPausedEN
		}
		generationLabel = GenerationLabelEN
		speedLabel = SpeedLabelEN
		sizeLabel = SizeLabelEN
		boundaryLabel = BoundaryLabelEN
		patternLabel = PatternLabelEN
	}

	tableBuilder.Reset()
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(generationLabel, m.currentStep)))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(speedLabel, m.refreshRate.String())))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(sizeLabel, m.gridHeight, m.gridWidth)))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(boundaryLabel, m.boundary.ToString(m.language))))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(patternLabel, m.pattern.ToString(m.language))))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(status))

	return lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(tableBuilder.String())
}

// ControlLineView returns the control display string: T,B,R + Space, L, Q
func (m Model) ControlLineView() string {
	var selectPattern, selectBoundary, speedControl, language, space, reset, quit string
	if m.language == Chinese {
		selectPattern = SelectPatternLabelCN
		selectBoundary = SelectBoundaryLabelCN
		language = LanguageLabelCN
		speedControl = SpeedControlLabelCN
		space = SpaceControlLabelCN
		reset = ResetLabelCN
		quit = QuitLabelCN
	} else {
		selectPattern = SelectPatternLabelEN
		selectBoundary = SelectBoundaryLabelEN
		language = LanguageLabelEN
		speedControl = SpeedControlLabelEN
		space = SpaceControlLabelEN
		reset = ResetLabelEN
		quit = QuitLabelEN
	}
	tableBuilder.Reset()
	tableBuilder.WriteString(labelStyle.Render(selectPattern))
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
