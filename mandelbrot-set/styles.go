package main

import (
	"fmt"
	"math"
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

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	tableBuilder strings.Builder
)

// UI text constants with enhanced formatting and icons
const (
	// Header Line
	HeaderCN = "🌀 曼德博集合 🌀"
	HeaderEN = "🌀 Mandelbrot Set 🌀"

	// Status Line
	ModeLabelCN = "🎯 模式: %s"
	ModeLabelEN = "🎯 Mode: %s"

	ZoomLabelCN = "🔍 缩放: %.2f"
	ZoomLabelEN = "🔍 Zoom: %.2f"

	CenterLabelCN = "📍 中心: (%.4f, %.4f)"
	CenterLabelEN = "📍 Center: (%.4f, %.4f)"

	IterLabelCN = "🔄 迭代: %d"
	IterLabelEN = "🔄 Iter: %d"

	ColorLabelCN = "🎨 配色: %s"
	ColorLabelEN = "🎨 Color: %s"

	StatusLabelCalculatingCN = "⚡ 计算中"
	StatusLabelCalculatingEN = "⚡ Calculating"
	StatusLabelReadyCN       = "✅ 就绪"
	StatusLabelReadyEN       = "✅ Ready"

	ModeNameMandelbrotCN = "曼德博"
	ModeNameMandelbrotEN = "Mandelbrot"
	ModeNameJuliaCN      = "朱利亚"
	ModeNameJuliaEN      = "Julia"

	// Control Line
	MoveControlLabelCN = "WASD/方向键 移动"
	MoveControlLabelEN = "WASD/Arrows Move"

	ZoomControlLabelCN = "+/- 缩放"
	ZoomControlLabelEN = "+/- Zoom"

	ModeControlLabelCN = "M 切换模式"
	ModeControlLabelEN = "M Toggle Mode"

	ColorControlLabelCN = "C 切换配色"
	ColorControlLabelEN = "C Toggle Color"

	IterControlLabelCN = "I/K 迭代+/-"
	IterControlLabelEN = "I/K Iter +/-"

	PresetControlLabelCN = "P 预设位置"
	PresetControlLabelEN = "P Preset Location"

	LanguageLabelCN = "L 切换语言"
	LanguageLabelEN = "L Switch Language"

	ResetLabelCN = "R 重置"
	ResetLabelEN = "R Reset"

	QuitLabelCN = "Q 退出"
	QuitLabelEN = "Q Quit"

	// Julia parameter line
	JuliaParamLabelCN = "🔢 朱利亚参数: %v"
	JuliaParamLabelEN = "🔢 Julia Parameter: %v"
)

// RenderOptions holds rendering configuration
type RenderOptions struct {
	colorScheme ColorScheme
}

// NewRenderOptions creates new render options
func NewRenderOptions(colorScheme ColorScheme) RenderOptions {
	return RenderOptions{
		colorScheme: colorScheme,
	}
}

// GetColorForIteration returns the color for a given iteration count
func (ro RenderOptions) GetColorForIteration(iter, maxIter int) lipgloss.Color {
	if iter >= maxIter {
		return lipgloss.Color("#000000") // Black for points in the set
	}

	// Normalize iteration count to 0-1 range
	ratio := float64(iter) / float64(maxIter)

	switch ro.colorScheme {
	case ColorSchemeClassic:
		return ro.getClassicColor(ratio)
	case ColorSchemeHot:
		return ro.getHotColor(ratio)
	case ColorSchemeCool:
		return ro.getCoolColor(ratio)
	case ColorSchemeRainbow:
		return ro.getRainbowColor(ratio)
	case ColorSchemeGrayscale:
		return ro.getGrayscaleColor(ratio)
	default:
		return ro.getClassicColor(ratio)
	}
}

// getClassicColor returns classic black and white colors
func (ro RenderOptions) getClassicColor(ratio float64) lipgloss.Color {
	if ratio < 0.1 {
		return lipgloss.Color("#000000") // Black
	} else if ratio < 0.3 {
		return lipgloss.Color("#404040") // Dark gray
	} else if ratio < 0.6 {
		return lipgloss.Color("#808080") // Gray
	} else if ratio < 0.8 {
		return lipgloss.Color("#C0C0C0") // Light gray
	}
	return lipgloss.Color("#FFFFFF") // White
}

// getHotColor returns hot color palette (red, orange, yellow)
func (ro RenderOptions) getHotColor(ratio float64) lipgloss.Color {
	if ratio < 0.2 {
		return lipgloss.Color("#000000") // Black
	} else if ratio < 0.4 {
		return lipgloss.Color("#800000") // Dark red
	} else if ratio < 0.6 {
		return lipgloss.Color("#FF0000") // Red
	} else if ratio < 0.8 {
		return lipgloss.Color("#FF8000") // Orange
	}
	return lipgloss.Color("#FFFF00") // Yellow
}

// getCoolColor returns cool color palette (blue, cyan, purple)
func (ro RenderOptions) getCoolColor(ratio float64) lipgloss.Color {
	if ratio < 0.2 {
		return lipgloss.Color("#000000") // Black
	} else if ratio < 0.4 {
		return lipgloss.Color("#000080") // Dark blue
	} else if ratio < 0.6 {
		return lipgloss.Color("#0000FF") // Blue
	} else if ratio < 0.8 {
		return lipgloss.Color("#00FFFF") // Cyan
	}
	return lipgloss.Color("#8000FF") // Purple
}

// getRainbowColor returns rainbow spectrum colors
func (ro RenderOptions) getRainbowColor(ratio float64) lipgloss.Color {
	// Use HSV color space for smooth rainbow transition
	hue := ratio * 360.0 // Full spectrum
	return ro.hsvToRGB(hue, 1.0, 1.0)
}

// getGrayscaleColor returns grayscale colors
func (ro RenderOptions) getGrayscaleColor(ratio float64) lipgloss.Color {
	intensity := int(ratio * 255)
	return lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", intensity, intensity, intensity))
}

// hsvToRGB converts HSV color to RGB hex string
func (ro RenderOptions) hsvToRGB(h, s, v float64) lipgloss.Color {
	h = math.Mod(h, 360)
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	var r, g, b float64
	if h < 60 {
		r, g, b = c, x, 0
	} else if h < 120 {
		r, g, b = x, c, 0
	} else if h < 180 {
		r, g, b = 0, c, x
	} else if h < 240 {
		r, g, b = 0, x, c
	} else if h < 300 {
		r, g, b = x, 0, c
	} else {
		r, g, b = c, 0, x
	}

	r = (r + m) * 255
	g = (g + m) * 255
	b = (b + m) * 255

	return lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", int(r), int(g), int(b)))
}

// GetCharacterForIteration returns the character to display for a given iteration count
func (ro RenderOptions) GetCharacterForIteration(iter, maxIter int) string {
	if iter >= maxIter {
		return "█" // Solid block for points in the set
	}

	// Use different characters based on iteration count
	ratio := float64(iter) / float64(maxIter)

	if ratio < 0.1 {
		return " " // Space
	} else if ratio < 0.2 {
		return "░" // Light shade
	} else if ratio < 0.4 {
		return "▒" // Medium shade
	} else if ratio < 0.6 {
		return "▓" // Dark shade
	} else if ratio < 0.8 {
		return "█" // Full block
	}
	return "█" // Full block
}

// HeaderLineView returns the header display string
func (m Model) HeaderLineView() string {
	style := headerStyle.Width(m.width)
	if m.language == Chinese {
		return style.Render(HeaderCN)
	}
	return style.Render(HeaderEN)
}

// StatusLineView returns the status display string
func (m Model) StatusLineView() string {
	var status, modeLabel, zoomLabel, centerLabel, iterLabel, colorLabel string
	var modeName string

	if m.language == Chinese {
		status = StatusLabelReadyCN
		if m.calculating {
			status = StatusLabelCalculatingCN
		}
		modeLabel = ModeLabelCN
		zoomLabel = ZoomLabelCN
		centerLabel = CenterLabelCN
		iterLabel = IterLabelCN
		colorLabel = ColorLabelCN
		if m.mandelbrotSet.GetCurrentMode() {
			modeName = ModeNameJuliaCN
		} else {
			modeName = ModeNameMandelbrotCN
		}
	} else {
		status = StatusLabelReadyEN
		if m.calculating {
			status = StatusLabelCalculatingEN
		}
		modeLabel = ModeLabelEN
		zoomLabel = ZoomLabelEN
		centerLabel = CenterLabelEN
		iterLabel = IterLabelEN
		colorLabel = ColorLabelEN
		if m.mandelbrotSet.GetCurrentMode() {
			modeName = ModeNameJuliaEN
		} else {
			modeName = ModeNameMandelbrotEN
		}
	}

	centerX, centerY := m.mandelbrotSet.GetCenter()

	tableBuilder.Reset()
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(modeLabel, modeName)))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(zoomLabel, m.mandelbrotSet.GetZoom())))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(centerLabel, centerX, centerY)))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(iterLabel, m.mandelbrotSet.GetMaxIterations())))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(fmt.Sprintf(colorLabel, m.mandelbrotSet.GetColorScheme().ToString(m.language))))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(status))

	statusLine := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(tableBuilder.String())

	// Julia parameter line (if in Julia mode)
	if m.mandelbrotSet.GetCurrentMode() {
		var juliaParamLabel string
		if m.language == Chinese {
			juliaParamLabel = JuliaParamLabelCN
		} else {
			juliaParamLabel = JuliaParamLabelEN
		}
		juliaParam := m.mandelbrotSet.GetJuliaParameter()
		juliaLine := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(
			labelStyle.Render(fmt.Sprintf(juliaParamLabel, juliaParam)))
		statusLine += "\n" + juliaLine
	}

	return statusLine
}

// ControlLineView returns the control display string
func (m Model) ControlLineView() string {
	var moveControl, zoomControl, modeControl, colorControl, iterControl, presetControl, language, reset, quit string
	if m.language == Chinese {
		moveControl = MoveControlLabelCN
		zoomControl = ZoomControlLabelCN
		modeControl = ModeControlLabelCN
		colorControl = ColorControlLabelCN
		iterControl = IterControlLabelCN
		presetControl = PresetControlLabelCN
		language = LanguageLabelCN
		reset = ResetLabelCN
		quit = QuitLabelCN
	} else {
		moveControl = MoveControlLabelEN
		zoomControl = ZoomControlLabelEN
		modeControl = ModeControlLabelEN
		colorControl = ColorControlLabelEN
		iterControl = IterControlLabelEN
		presetControl = PresetControlLabelEN
		language = LanguageLabelEN
		reset = ResetLabelEN
		quit = QuitLabelEN
	}

	tableBuilder.Reset()
	tableBuilder.WriteString(labelStyle.Render(moveControl))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(zoomControl))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(modeControl))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(colorControl))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(iterControl))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(presetControl))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(language))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(reset))
	tableBuilder.WriteString(" | ")
	tableBuilder.WriteString(labelStyle.Render(quit))

	controlLine := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(tableBuilder.String())

	// Current preset info
	if preset := m.getCurrentPresetInfo(); preset != "" {
		controlLine += "\n" + preset
	}

	return controlLine
}
