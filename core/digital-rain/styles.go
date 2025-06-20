package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderOptions holds rendering options for the digital rain
type RenderOptions struct {
	dropStyle lipgloss.Style
	bgStyle   lipgloss.Style
	// Pre-computed trail styles for different intensities
	trailStyles map[int]lipgloss.Style
}

// NewRenderOptions creates new render options with the given colors
func NewRenderOptions(dropColor, trailColor, backgroundColor string) RenderOptions {
	opts := RenderOptions{
		dropStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color(dropColor)),
		bgStyle:     lipgloss.NewStyle().Background(lipgloss.Color(backgroundColor)),
		trailStyles: make(map[int]lipgloss.Style),
	}

	// Pre-compute trail styles for different intensities
	dropRGB := hexToRGB(dropColor)
	trailRGB := hexToRGB(trailColor)

	for i := 0; i <= 10; i++ {
		intensity := float64(i) / 10.0
		r := int(float64(trailRGB.r) + (float64(dropRGB.r-trailRGB.r) * intensity))
		g := int(float64(trailRGB.g) + (float64(dropRGB.g-trailRGB.g) * intensity))
		b := int(float64(trailRGB.b) + (float64(dropRGB.b-trailRGB.b) * intensity))
		color := fmt.Sprintf("#%02x%02x%02x", r, g, b)
		opts.trailStyles[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	}

	return opts
}

// RGB represents an RGB color
type RGB struct {
	r, g, b int
}

// hexToRGB converts hex color to RGB
// Returns black (0,0,0) for invalid input
func hexToRGB(hex string) RGB {
	// Remove # prefix if present
	hex = strings.TrimPrefix(hex, "#")

	// Validate hex string length
	if len(hex) != 6 {
		// Return black for invalid length
		return RGB{0, 0, 0}
	}

	// Parse RGB components with error handling
	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return RGB{0, 0, 0}
	}

	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return RGB{0, 0, 0}
	}

	b, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		return RGB{0, 0, 0}
	}

	return RGB{int(r), int(g), int(b)}
}

// GetTrailStyle returns the appropriate style for a given trail intensity
func (ro *RenderOptions) GetTrailStyle(intensity int) lipgloss.Style {
	// Map intensity (0-255) to style index (0-10)
	idx := intensity * 10 / 255
	if idx > 10 {
		idx = 10
	}
	if idx < 0 {
		idx = 0
	}
	return ro.trailStyles[idx]
}

// headerStyle returns the style for headers
func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF00")).
		MarginBottom(1)
}

// statusStyle returns the style for status text
func statusStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00AA00"))
}

// helpStyle returns the style for help text
func helpStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))
}
