package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/telepair/go-playground/pkg"
)

func main() {
	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Cellular Automaton - A Terminal User Interface implementation of 1D cellular automata\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -rule 30                         # Run Rule 30 (Random) with auto-detected size\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -rule 90                         # Run Rule 90 (Sierpinski Triangle) with auto-detected size\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -rule 110                        # Run Rule 110 (Turing Machine) with auto-detected size\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -rule 184 -alive-char '🚗'        # Run Rule 184 (Traffic Simulation) with custom alive character\n", os.Args[0])
	}

	// Parse command line flags
	var rule = flag.Int("rule", DefaultRule, "Cellular automaton rule number (0-255)")
	var aliveColor = flag.String("alive-color", DefaultAliveColor, "Alive cell color (hex)")
	var deadColor = flag.String("dead-color", DefaultDeadColor, "Dead cell color (hex)")
	var aliveChar = flag.String("alive-char", DefaultAliveChar, "Alive cell character")
	var deadChar = flag.String("dead-char", DefaultDeadChar, "Dead cell character")
	var lang = flag.String("lang", DefaultLanguage.ToString(DefaultLanguage), "Language (en/cn)")
	var enableProfiling = flag.Bool("profile", false, "Enable profiling and monitoring")
	var profilePort = flag.Int("profile-port", DefaultProfilePort, "Profiling server port")
	var profileInterval = flag.Duration("profile-interval", DefaultProfileInterval, "Profile information output interval")
	var logFile = flag.String("log-file", DefaultLogFile, "Log file path")

	flag.Parse()

	if *logFile != "" {
		_ = pkg.InitLog("debug", "text", *logFile)
	}
	slog.Debug("Cellular Automaton starting")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if *enableProfiling {
		go pkg.StartProfile(ctx, *profilePort)
		go pkg.StartWatchdog(ctx, *profileInterval)
	}

	// Create and configure application
	config := Config{
		Rule:       *rule,
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
		slog.Error("Error running program", "error", err)
		os.Exit(1)
	}

	slog.Debug("Cellular Automaton finished")
}
