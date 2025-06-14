package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Parse command line flags
	var rule = flag.Uint("rule", DefaultRule, "Cellular automaton rule number (0-255)")
	var steps = flag.Uint("steps", DefaultSteps, "Number of steps (0 or negative for infinite mode)")
	var sizeStr = flag.String("size", "60x30", "Window size (format: WIDTHxHEIGHT, e.g.: 60x30)")
	var cellSize = flag.Uint("cellsize", DefaultCellSize, "Cell size (1-3)")
	var aliveColor = flag.String("alive-color", DefaultAliveColor, "Alive cell color (hex)")
	var deadColor = flag.String("dead-color", DefaultDeadColor, "Dead cell color (hex)")
	var aliveChar = flag.String("alive-char", DefaultAliveChar, "Alive cell character")
	var deadChar = flag.String("dead-char", DefaultDeadChar, "Dead cell character")
	var refreshRate = flag.Float64("refresh", DefaultRefreshRate, "Refresh rate (seconds, minimum 0.001)")
	var lang = flag.String("lang", DefaultLanguage, "Language (en/cn)")
	var boundary = flag.String("boundary", DefaultBoundary.String(), "Boundary type (periodic/fixed/reflect)")

	flag.Parse()

	// Create and configure application
	config := NewConfig()
	config.SetRule(*rule)
	config.SetSteps(*steps)
	config.SetCellSize(*cellSize)
	config.SetRefreshRate(*refreshRate)
	config.SetLanguage(*lang)
	config.SetBoundary(*boundary)
	config.SetWindowSize(*sizeStr)
	config.AliveColor = *aliveColor
	config.DeadColor = *deadColor
	config.AliveChar = *aliveChar
	config.DeadChar = *deadChar

	// Create initial model
	initialModel := NewModel(config)

	// Run the application
	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
