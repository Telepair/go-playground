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
		fmt.Fprintf(os.Stderr, "Conway's Game of Life - A Terminal User Interface implementation of Conway's cellular automaton\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                                  # Run with default settings\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -pattern glider                  # Start with a glider pattern\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -pattern glider-gun -size 30x80  # Glider gun in custom size\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -alive-char '🟢' -dead-char '⚫' # Custom emoji cells\n", os.Args[0])
	}

	// Parse command line flags
	var rows = flag.Int("rows", DefaultWindowRows, "Number of rows in the grid")
	var cols = flag.Int("cols", DefaultWindowCols, "Number of columns in the grid")
	var aliveColor = flag.String("alive-color", DefaultAliveColor, "Alive cell color (hex)")
	var deadColor = flag.String("dead-color", DefaultDeadColor, "Dead cell color (hex)")
	var aliveChar = flag.String("alive-char", DefaultAliveChar, "Alive cell character")
	var deadChar = flag.String("dead-char", DefaultDeadChar, "Dead cell character")
	var lang = flag.String("lang", DefaultLanguage, "Language (en/cn)")

	flag.Parse()

	// Create and configure application
	config := NewConfig()
	var errors []error

	config.SetLanguage(*lang)

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
