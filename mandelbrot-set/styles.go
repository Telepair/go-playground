package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
func (r RenderOptions) GetColorForIteration(iter, maxIter int) lipgloss.Color {
	if iter >= maxIter {
		return lipgloss.Color("#000000") // Black for points in the set
	}

	// Normalize iteration count to 0-1 range
	ratio := float64(iter) / float64(maxIter)

	switch r.colorScheme {
	case ColorSchemeClassic:
		return r.getClassicColor(ratio)
	case ColorSchemeHot:
		return r.getHotColor(ratio)
	case ColorSchemeCool:
		return r.getCoolColor(ratio)
	case ColorSchemeRainbow:
		return r.getRainbowColor(ratio)
	case ColorSchemeGrayscale:
		return r.getGrayscaleColor(ratio)
	default:
		return r.getClassicColor(ratio)
	}
}

// getClassicColor returns classic black and white colors
func (r RenderOptions) getClassicColor(ratio float64) lipgloss.Color {
	if ratio < 0.1 {
		return lipgloss.Color("#000000") // Black
	} else if ratio < 0.3 {
		return lipgloss.Color("#404040") // Dark gray
	} else if ratio < 0.6 {
		return lipgloss.Color("#808080") // Gray
	} else if ratio < 0.8 {
		return lipgloss.Color("#C0C0C0") // Light gray
	} else {
		return lipgloss.Color("#FFFFFF") // White
	}
}

// getHotColor returns hot color palette (red, orange, yellow)
func (r RenderOptions) getHotColor(ratio float64) lipgloss.Color {
	if ratio < 0.2 {
		return lipgloss.Color("#000000") // Black
	} else if ratio < 0.4 {
		return lipgloss.Color("#800000") // Dark red
	} else if ratio < 0.6 {
		return lipgloss.Color("#FF0000") // Red
	} else if ratio < 0.8 {
		return lipgloss.Color("#FF8000") // Orange
	} else {
		return lipgloss.Color("#FFFF00") // Yellow
	}
}

// getCoolColor returns cool color palette (blue, cyan, purple)
func (r RenderOptions) getCoolColor(ratio float64) lipgloss.Color {
	if ratio < 0.2 {
		return lipgloss.Color("#000000") // Black
	} else if ratio < 0.4 {
		return lipgloss.Color("#000080") // Dark blue
	} else if ratio < 0.6 {
		return lipgloss.Color("#0000FF") // Blue
	} else if ratio < 0.8 {
		return lipgloss.Color("#00FFFF") // Cyan
	} else {
		return lipgloss.Color("#8000FF") // Purple
	}
}

// getRainbowColor returns rainbow spectrum colors
func (r RenderOptions) getRainbowColor(ratio float64) lipgloss.Color {
	// Use HSV color space for smooth rainbow transition
	hue := ratio * 360.0 // Full spectrum
	return r.hsvToRGB(hue, 1.0, 1.0)
}

// getGrayscaleColor returns grayscale colors
func (r RenderOptions) getGrayscaleColor(ratio float64) lipgloss.Color {
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
func (r RenderOptions) GetCharacterForIteration(iter, maxIter int) string {
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
	} else {
		return "█" // Full block
	}
}

// Style definitions for different UI elements
var (
	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5A5A5A")).
			Padding(0, 1)

	// Status style
	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#3A3A3A")).
			Padding(0, 1)

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)

// GetHeaderLine returns the formatted header line
func GetHeaderLine(language Language) string {
	if language == Chinese {
		return titleStyle.Render("曼德博集合 - 分形可视化")
	}
	return titleStyle.Render("Mandelbrot Set - Fractal Visualization")
}

// GetStatusLine returns the formatted status line
func GetStatusLine(m *MandelbrotSet, language Language) string {
	var status strings.Builder

	if language == Chinese {
		if m.GetCurrentMode() {
			status.WriteString("模式: 朱利亚集合")
		} else {
			status.WriteString("模式: 曼德博集合")
		}
		status.WriteString(fmt.Sprintf(" | 缩放: %.2f", m.GetZoom()))
		centerX, centerY := m.GetCenter()
		status.WriteString(fmt.Sprintf(" | 中心: (%.4f, %.4f)", centerX, centerY))
		status.WriteString(fmt.Sprintf(" | 迭代: %d", m.GetMaxIterations()))
		status.WriteString(fmt.Sprintf(" | 配色: %s", m.GetColorScheme().ToString(language)))
	} else {
		if m.GetCurrentMode() {
			status.WriteString("Mode: Julia Set")
		} else {
			status.WriteString("Mode: Mandelbrot Set")
		}
		status.WriteString(fmt.Sprintf(" | Zoom: %.2f", m.GetZoom()))
		centerX, centerY := m.GetCenter()
		status.WriteString(fmt.Sprintf(" | Center: (%.4f, %.4f)", centerX, centerY))
		status.WriteString(fmt.Sprintf(" | Iterations: %d", m.GetMaxIterations()))
		status.WriteString(fmt.Sprintf(" | Colors: %s", m.GetColorScheme().ToString(language)))
	}

	return statusStyle.Render(status.String())
}

// GetHelpLine returns the formatted help line
func GetHelpLine(language Language) string {
	if language == Chinese {
		return helpStyle.Render("控制: 方向键(移动) | +/-(缩放) | M(模式) | C(配色) | I/K(迭代) | P(预设) | L(语言) | R(重置) | Q(退出)")
	}
	return helpStyle.Render("Controls: Arrow Keys(Pan) | +/-(Zoom) | M(Mode) | C(Colors) | I/K(Iterations) | P(Presets) | L(Language) | R(Reset) | Q(Quit)")
}

// GetJuliaParameterLine returns the Julia parameter information
func GetJuliaParameterLine(m *MandelbrotSet, language Language) string {
	if !m.GetCurrentMode() {
		return ""
	}

	juliaC := m.GetJuliaParameter()
	if language == Chinese {
		return helpStyle.Render(fmt.Sprintf("朱利亚参数: %.4f + %.4fi", real(juliaC), imag(juliaC)))
	}
	return helpStyle.Render(fmt.Sprintf("Julia Parameter: %.4f + %.4fi", real(juliaC), imag(juliaC)))
}
