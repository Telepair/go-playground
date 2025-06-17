package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
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
	var lang = flag.String("lang", DefaultLanguage.ToString(DefaultLanguage), "Language (en/cn)")
	var enableProfiling = flag.Bool("profile", false, "Enable profiling and monitoring")
	var profilePort = flag.String("profile-port", ":6060", "Profiling server port")

	flag.Parse()

	// Initialize monitoring if enabled
	if *enableProfiling {
		go func() {
			log.Printf("Starting pprof server on http://localhost%s/debug/pprof/", *profilePort)
			log.Println(http.ListenAndServe(*profilePort, nil)) //nolint:gosec
		}()
	}

	// Create and configure application
	config := Config{
		Rows:       *rows,
		Cols:       *cols,
		AliveColor: *aliveColor,
		DeadColor:  *deadColor,
		AliveChar:  *aliveChar,
		DeadChar:   *deadChar,
	}
	config.SetLanguage(*lang)
	config.Check()

	// Create initial model
	initialModel := NewModel(config)

	// Run the application
	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
