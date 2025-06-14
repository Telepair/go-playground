package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Cellular Automaton - A Terminal User Interface implementation of 1D cellular automata\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -rule 30                         # Run Rule 30 with default settings\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -rule 110 -steps 100             # Run Rule 110 for 100 steps\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -rule 90 -steps 0                # Run Rule 90 in infinite mode\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -rule 184 -alive-char '🚗'       # Traffic simulation\n", os.Args[0])
	}

	// Parse command line flags
	var rule = flag.Int("rule", DefaultRule, "Cellular automaton rule number (0-255)")
	var steps = flag.Int("steps", DefaultSteps, "Number of steps (0 or negative for infinite mode)")
	var rows = flag.Int("rows", DefaultWindowRows, "Number of rows")
	var cols = flag.Int("cols", DefaultWindowCols, "Number of columns")
	var cellSize = flag.Int("cellsize", DefaultCellSize, "Cell size (1-3)")
	var aliveColor = flag.String("alive-color", DefaultAliveColor, "Alive cell color (hex)")
	var deadColor = flag.String("dead-color", DefaultDeadColor, "Dead cell color (hex)")
	var aliveChar = flag.String("alive-char", DefaultAliveChar, "Alive cell character")
	var deadChar = flag.String("dead-char", DefaultDeadChar, "Dead cell character")
	var refreshRate = flag.Duration("refresh", DefaultRefreshRate, "Refresh rate (minimum 1ms)")
	var lang = flag.String("lang", DefaultLanguage, "Language (en/cn)")
	var boundary = flag.String("boundary", DefaultBoundary.String(), "Boundary type (periodic/fixed/reflect)")

	flag.Parse()

	// Create and configure application
	config := NewConfig()
	var errors []error

	// Set configuration values with error collection
	if err := config.SetRule(*rule); err != nil {
		errors = append(errors, err)
	}

	config.SetSteps(*steps)

	if err := config.SetCellSize(*cellSize); err != nil {
		errors = append(errors, err)
	}

	if err := config.SetRefreshRate(*refreshRate); err != nil {
		errors = append(errors, err)
	}

	config.SetLanguage(*lang)
	config.SetBoundary(*boundary)

	if err := config.SetRows(*rows); err != nil {
		errors = append(errors, err)
	}

	if err := config.SetCols(*cols); err != nil {
		errors = append(errors, err)
	}

	config.AliveColor = *aliveColor
	config.DeadColor = *deadColor
	config.AliveChar = *aliveChar
	config.DeadChar = *deadChar

	// Validate colors
	if err := config.ValidateColors(); err != nil {
		errors = append(errors, err)
	}

	// Log configuration errors if any (they've been logged already but we count them)
	if len(errors) > 0 {
		configLogger.Printf("Configuration completed with %d warnings", len(errors))
	}

	// Create initial model
	initialModel := NewModel(config)

	// Run the application
	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
