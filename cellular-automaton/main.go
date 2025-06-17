package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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
		fmt.Fprintf(os.Stderr, "  %s -rule 30                         # Run Rule 30 (Random) with default settings\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -rule 90                         # Run Rule 90 (Sierpinski Triangle) in infinite mode\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -rule 110                        # Run Rule 110 (Turing Machine) with default settings\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -rule 184 -alive-char '🚗'        # Run Rule 184 (Traffic Simulation) with custom alive character\n", os.Args[0])
	}

	// Parse command line flags
	var rule = flag.Int("rule", DefaultRule, "Cellular automaton rule number (0-255)")
	var rows = flag.Int("rows", DefaultWindowRows, "Number of rows")
	var cols = flag.Int("cols", DefaultWindowCols, "Number of columns")
	var aliveColor = flag.String("alive-color", DefaultAliveColor, "Alive cell color (hex)")
	var deadColor = flag.String("dead-color", DefaultDeadColor, "Dead cell color (hex)")
	var aliveChar = flag.String("alive-char", DefaultAliveChar, "Alive cell character")
	var deadChar = flag.String("dead-char", DefaultDeadChar, "Dead cell character")
	var lang = flag.String("lang", DefaultLanguage.ToString(DefaultLanguage), "Language (en/cn)")
	var enableProfiling = flag.Bool("profile", false, "Enable profiling and monitoring")
	var profilePort = flag.String("profile-port", ":6060", "Profiling server port")

	flag.Parse()

	if *enableProfiling {
		go func() {
			log.Printf("Starting pprof server on http://localhost%s/debug/pprof/", *profilePort)
			log.Println(http.ListenAndServe(*profilePort, nil)) //nolint:gosec
		}()
	}

	// Create and configure application
	config := Config{
		Rule:       *rule,
		Rows:       *rows,
		Cols:       *cols,
		AliveColor: *aliveColor,
		DeadColor:  *deadColor,
		AliveChar:  *aliveChar,
		DeadChar:   *deadChar,
	}
	config.SetLang(*lang)
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
