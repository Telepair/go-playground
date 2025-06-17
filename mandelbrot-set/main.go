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
	var rows = flag.Int("rows", DefaultWindowRows, "Number of rows in the grid")
	var cols = flag.Int("cols", DefaultWindowCols, "Number of columns in the grid")
	var maxIter = flag.Int("max-iter", DefaultMaxIterations, "Maximum number of iterations")
	var zoom = flag.Float64("zoom", DefaultZoom, "Zoom level")
	var centerX = flag.Float64("center-x", DefaultCenterX, "Center X coordinate")
	var centerY = flag.Float64("center-y", DefaultCenterY, "Center Y coordinate")
	var colorScheme = flag.Int("color-scheme", int(DefaultColorScheme), "Color scheme (0-4)")
	var julia = flag.Bool("julia", false, "Enable Julia set mode")
	var juliaC = flag.String("julia-c", DefaultJuliaC, "Julia set parameter (complex number)")
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
		Rows:        *rows,
		Cols:        *cols,
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
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
