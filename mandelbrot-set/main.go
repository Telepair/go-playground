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
		fmt.Fprintf(os.Stderr, "Mandelbrot Set - A Terminal User Interface implementation of the Mandelbrot fractal set\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                                  # Run with default settings\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -zoom 2.0 -center-x -0.5        # Zoom into a specific area\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -max-iter 100 -color-scheme 2   # High iteration with different colors\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -julia -julia-c '0.285+0.01i'   # Julia set mode with custom parameter\n", os.Args[0])
	}

	// Parse command line flags
	var maxIter = flag.Int("max-iter", DefaultMaxIterations, "Maximum number of iterations")
	var zoom = flag.Float64("zoom", DefaultZoom, "Zoom level")
	var centerX = flag.Float64("center-x", DefaultCenterX, "Center X coordinate")
	var centerY = flag.Float64("center-y", DefaultCenterY, "Center Y coordinate")
	var colorScheme = flag.Int("color-scheme", int(DefaultColorScheme), "Color scheme (0-4)")
	var julia = flag.Bool("julia", false, "Enable Julia set mode")
	var juliaC = flag.String("julia-c", DefaultJuliaC, "Julia set parameter (complex number)")
	var lang = flag.String("lang", DefaultLanguage.ToString(DefaultLanguage), "Language (en/cn)")
	var enableProfiling = flag.Bool("profile", false, "Enable profiling and monitoring")
	var profilePort = flag.Int("profile-port", DefaultProfilePort, "Profiling server port")
	var profileInterval = flag.Duration("profile-interval", DefaultProfileInterval, "Profile information output interval")
	var logFile = flag.String("log-file", DefaultLogFile, "Log file path")

	flag.Parse()

	if *logFile != "" {
		_ = pkg.InitLog("debug", "text", *logFile)
	}
	slog.Debug("Mandelbrot Set starting")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if *enableProfiling {
		go pkg.StartProfile(ctx, *profilePort)
		go pkg.StartWatchdog(ctx, *profileInterval)
	}

	// Create and configure application
	config := Config{
		MaxIter:     *maxIter,
		Zoom:        *zoom,
		CenterX:     *centerX,
		CenterY:     *centerY,
		ColorScheme: ColorScheme(*colorScheme),
		Julia:       *julia,
		JuliaC:      *juliaC,
	}
	config.SetLanguage(*lang)
	config.Check()

	// Create initial model
	initialModel := NewModel(config)

	// Run the application
	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		slog.Error("Error running program", "error", err)
		os.Exit(1)
	}

	slog.Debug("Mandelbrot Set finished")
}
