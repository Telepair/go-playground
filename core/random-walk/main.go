package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	_ "net/http/pprof" //nolint:gosec
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/telepair/go-playground/pkg"
)

func main() {
	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Random Walk Visualization - A Terminal User Interface implementation of various random walk algorithms\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                                  # Run with default settings\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -walker-char '🐾' -trail-char '·' # Custom walker and trail characters\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -walker-color '#FF00FF'          # Custom walker color\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -lang cn                         # Run in Chinese\n", os.Args[0])
	}

	// Parse command line flags
	var walkerColor = flag.String("walker-color", DefaultWalkerColor, "Walker color (hex)")
	var trailColor = flag.String("trail-color", DefaultTrailColor, "Trail color (hex)")
	var emptyColor = flag.String("empty-color", DefaultEmptyColor, "Empty cell color (hex)")
	var walkerChar = flag.String("walker-char", DefaultWalkerChar, "Walker character")
	var trailChar = flag.String("trail-char", DefaultTrailChar, "Trail character")
	var emptyChar = flag.String("empty-char", DefaultEmptyChar, "Empty cell character")
	var lang = flag.String("lang", DefaultLanguage.ToString(DefaultLanguage), "Language (en/cn)")
	var enableProfiling = flag.Bool("profile", false, "Enable profiling and monitoring")
	var profilePort = flag.Int("profile-port", DefaultProfilePort, "Profiling server port")
	var profileInterval = flag.Duration("profile-interval", DefaultProfileInterval, "Profile information output interval")
	var logFile = flag.String("log-file", DefaultLogFile, "Log file path")

	flag.Parse()

	if *logFile != "" {
		_ = pkg.InitLog("debug", "text", *logFile)
	}
	slog.Debug("Random Walk Visualization starting")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize monitoring if enabled
	if *enableProfiling {
		go pkg.StartProfile(ctx, *profilePort)
		go pkg.StartWatchdog(ctx, *profileInterval)
	}

	// Create and configure application
	config := Config{
		WalkerColor: *walkerColor,
		TrailColor:  *trailColor,
		EmptyColor:  *emptyColor,
		WalkerChar:  *walkerChar,
		TrailChar:   *trailChar,
		EmptyChar:   *emptyChar,
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

	slog.Debug("Random Walk Visualization finished")
}
