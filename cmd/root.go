/*
Copyright © 2025 Liys <liys87x@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// Package cmd contains the command line interface for the go-playground application.
package cmd

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/telepair/go-playground/pkg"
	"github.com/telepair/go-playground/pkg/ui"
)

const (
	// DefaultProfilePort is the default profiling server port
	DefaultProfilePort = 6060
	// DefaultProfileInterval is the default interval for profile information output
	DefaultProfileInterval = 5 * time.Second
	// DefaultLogLevel is the default logging level
	DefaultLogLevel = "debug"
	// DefaultLogFormat is the default logging format
	DefaultLogFormat = "text"
	// DefaultLogFile is the default log file path (empty means stdout)
	DefaultLogFile = ""
)

var (
	lang            string
	refreshInterval time.Duration
	profile         bool
	profilePort     int
	profileInterval time.Duration
	logFile         string
	logLevel        string
	logFormat       string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-playground",
	Short: "A collection of terminal-based visual algorithms and simulations",
	Long: `Go Playground is a collection of terminal-based visual algorithms and simulations 
written in Go. It includes various interesting visual demonstrations such as:

- Cellular Automaton: Explore various cellular automaton rules
- Conway's Game of Life: The famous cellular automaton simulation
- Mandelbrot Set: Fractal visualization in your terminal
- Random Walk: Various random walk algorithms visualization
- Digital Rain: Matrix-style falling characters effect

Each visualization supports customization options including colors, characters,
refresh rates, and language settings (English/Chinese).`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&lang, "lang", ui.DefaultLang, "Language (en/cn)")
	rootCmd.PersistentFlags().DurationVar(&refreshInterval, "refresh-interval", ui.DefaultRefreshInterval, "Refresh interval")
	rootCmd.PersistentFlags().BoolVar(&profile, "profile", false, "Enable profiling and monitoring")
	rootCmd.PersistentFlags().IntVar(&profilePort, "profile-port", DefaultProfilePort, "Profiling server port")
	rootCmd.PersistentFlags().DurationVar(&profileInterval, "profile-interval", DefaultProfileInterval, "Profile information output interval")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", DefaultLogFile, "Log file path")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", DefaultLogLevel, "Log level (debug/info/warn/error)")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", DefaultLogFormat, "Log format (text/json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// InitLog initializes the logging
func InitLog() {
	if logFile != "" {
		if err := pkg.InitLog(logLevel, logFormat, logFile); err != nil {
			slog.Error("Failed to initialize logging", "error", err)
		}
	}
}

// InitProfile starts the profiling and watchdog
func InitProfile(ctx context.Context) {
	if profile {
		go pkg.StartProfile(ctx, profilePort)
		go pkg.StartWatchdog(ctx, profileInterval)
	}
}
