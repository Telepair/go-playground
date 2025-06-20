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
		fmt.Fprintf(os.Stderr, "Digital Rain - A Terminal User Interface implementation of the Matrix digital rain effect\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                              # Run with default settings\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -drop-color '#FFFFFF'        # White rain drops\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -charset '01'                # Binary rain\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -lang cn                     # Run in Chinese\n", os.Args[0])
	}

	// Parse command line flags
	var dropColor = flag.String("drop-color", DefaultDropColor, "Drop color (hex)")
	var trailColor = flag.String("trail-color", DefaultTrailColor, "Trail color (hex)")
	var bgColor = flag.String("bg-color", DefaultBackgroundColor, "Background color (hex)")
	var charset = flag.String("charset", DefaultCharSet, "Character set to use")
	var minSpeed = flag.Int("min-speed", DefaultMinSpeed, "Minimum drop speed")
	var maxSpeed = flag.Int("max-speed", DefaultMaxSpeed, "Maximum drop speed")
	var dropLength = flag.Int("drop-length", DefaultDropLength, "Drop length")
	var lang = flag.String("lang", DefaultLanguage.ToString(), "Language (en/cn)")
	var enableProfiling = flag.Bool("profile", false, "Enable profiling and monitoring")
	var profilePort = flag.Int("profile-port", DefaultProfilePort, "Profiling server port")
	var profileInterval = flag.Duration("profile-interval", DefaultProfileInterval, "Profile information output interval")
	var logFile = flag.String("log-file", DefaultLogFile, "Log file path")

	flag.Parse()

	if *logFile != "" {
		_ = pkg.InitLog("debug", "text", *logFile)
	}
	slog.Debug("Digital Rain starting")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize monitoring if enabled
	if *enableProfiling {
		go pkg.StartProfile(ctx, *profilePort)
		go pkg.StartWatchdog(ctx, *profileInterval)
	}

	// Create and configure application
	config := Config{
		DropColor:       *dropColor,
		TrailColor:      *trailColor,
		BackgroundColor: *bgColor,
		CharSet:         *charset,
		MinSpeed:        *minSpeed,
		MaxSpeed:        *maxSpeed,
		DropLength:      *dropLength,
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

	slog.Debug("Digital Rain finished")
}
