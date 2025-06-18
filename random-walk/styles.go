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
			Background(lipgloss.Color("#6B46C1")).
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
	HeaderCN = "🚶 随机游走可视化 🚶"
	HeaderEN = "🚶 Random Walk Visualization 🚶"

	// Status Line
	StepsLabelCN = "📍 步数: %d"
	StepsLabelEN = "📍 Steps: %d"

	SpeedLabelCN = "🔄 刷新: %s"
	SpeedLabelEN = "🔄 Speed: %s"

	SizeLabelCN = "📐 尺寸: %d×%d"
	SizeLabelEN = "📐 Size: %d×%d"

	ModeLabelCN = "🎨 模式: %s"
	ModeLabelEN = "🎨 Mode: %s"

	WalkersLabelCN = "👥 粒子数: %d"
	WalkersLabelEN = "👥 Walkers: %d"

	TrailLabelCN = "🌟 轨迹长度: %d"
	TrailLabelEN = "🌟 Trail: %d"

	StatusLabelPlayingCN = "▶️ 运行中"
	StatusLabelPlayingEN = "▶️ Running"
	StatusLabelPausedCN  = "⏸️ 已暂停"
	StatusLabelPausedEN  = "⏸️ Paused"

	// Control Line
	SelectModeLabelCN = "M 切换模式"
	SelectModeLabelEN = "M Change Mode"

	WalkerControlLabelCN = "W/w 粒子数 +/-"
	WalkerControlLabelEN = "W/w Walkers +/-"

	TrailControlLabelCN = "T/t 轨迹长度 +/-"
	TrailControlLabelEN = "T/t Trail +/-"

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
	walkerStyled string // Cached styled walker
	trailStyled  string // Cached styled trail
	emptyStyled  string // Cached styled empty cell
	walkerChar   string
	trailChar    string
	emptyChar    string
}

// NewRenderOptions creates optimized render options with pre-computed styles
func NewRenderOptions(walkerColor, trailColor, emptyColor, walkerChar, trailChar, emptyChar string) RenderOptions {
	return RenderOptions{
		walkerStyled: lipgloss.NewStyle().Foreground(lipgloss.Color(walkerColor)).Render(walkerChar),
		trailStyled:  lipgloss.NewStyle().Foreground(lipgloss.Color(trailColor)).Render(trailChar),
		emptyStyled:  lipgloss.NewStyle().Foreground(lipgloss.Color(emptyColor)).Render(emptyChar),
		walkerChar:   walkerChar,
		trailChar:    trailChar,
		emptyChar:    emptyChar,
	}
}

// getWalkerStyled returns a styled walker with custom color
func (ro RenderOptions) getWalkerStyled(color, char string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(char)
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
	var status, stepsLabel, speedLabel, sizeLabel, modeLabel, walkersLabel, trailLabel string

	if m.language == Chinese {
		status = StatusLabelPlayingCN
		if m.paused {
			status = StatusLabelPausedCN
		}
		stepsLabel = StepsLabelCN
		speedLabel = SpeedLabelCN
		sizeLabel = SizeLabelCN
		modeLabel = ModeLabelCN
		walkersLabel = WalkersLabelCN
		trailLabel = TrailLabelCN
	} else {
		status = StatusLabelPlayingEN
		if m.paused {
			status = StatusLabelPausedEN
		}
		stepsLabel = StepsLabelEN
		speedLabel = SpeedLabelEN
		sizeLabel = SizeLabelEN
		modeLabel = ModeLabelEN
		walkersLabel = WalkersLabelEN
		trailLabel = TrailLabelEN
	}

	tableBuilder.Reset()
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(stepsLabel, m.currentStep)))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(speedLabel, m.refreshRate.String())))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(sizeLabel, m.gridHeight, m.gridWidth)))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(modeLabel, m.mode.ToString(m.language))))

	// Show walker count for multi-walker modes
	if m.mode == ModeMultiWalker || m.mode == ModeBrownianMotion {
		tableBuilder.WriteString(" | ")
		tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(walkersLabel, m.walkerCount)))
	}

	// Show trail length for trail modes
	if m.mode == ModeTrailMode || m.mode == ModeBrownianMotion {
		tableBuilder.WriteString(" | ")
		tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(trailLabel, m.trailLength)))
	}

	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(status))

	return lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(tableBuilder.String())
}

// ControlLineView returns the control display string
func (m Model) ControlLineView() string {
	var selectMode, walkerControl, trailControl, speedControl, language, space, reset, quit string
	if m.language == Chinese {
		selectMode = SelectModeLabelCN
		walkerControl = WalkerControlLabelCN
		trailControl = TrailControlLabelCN
		language = LanguageLabelCN
		speedControl = SpeedControlLabelCN
		space = SpaceControlLabelCN
		reset = ResetLabelCN
		quit = QuitLabelCN
	} else {
		selectMode = SelectModeLabelEN
		walkerControl = WalkerControlLabelEN
		trailControl = TrailControlLabelEN
		language = LanguageLabelEN
		speedControl = SpeedControlLabelEN
		space = SpaceControlLabelEN
		reset = ResetLabelEN
		quit = QuitLabelEN
	}

	tableBuilder.Reset()
	tableBuilder.WriteString(labelStyle.Render(selectMode))

	// Show walker control for multi-walker modes
	if m.mode == ModeMultiWalker || m.mode == ModeBrownianMotion {
		tableBuilder.WriteString(" | ")
		tableBuilder.WriteString(labelStyle.Render(walkerControl))
	}

	// Show trail control for trail modes
	if m.mode == ModeTrailMode || m.mode == ModeBrownianMotion {
		tableBuilder.WriteString(" | ")
		tableBuilder.WriteString(labelStyle.Render(trailControl))
	}

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
