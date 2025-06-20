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

	"github.com/spf13/cobra"

	"github.com/telepair/go-playground/engine/gameoflife"
	"github.com/telepair/go-playground/pkg/ui"
)

// conwayGameOfLifeCmd represents the Conway's Game of Life command
var conwayGameOfLifeCmd = &cobra.Command{
	Use:   "conway-game-of-life",
	Short: "Run Conway's Game of Life simulation",
	Long: `Run Conway's Game of Life, a cellular automaton devised by mathematician John Conway.

Conway's Game of Life is a zero-player game, meaning that its evolution is 
determined by its initial state, requiring no further input. It consists of 
a grid of cells which can be either alive or dead. Each cell evolves according 
to simple rules:

Rules:
1. Any live cell with 2-3 live neighbors survives
2. Any dead cell with exactly 3 live neighbors becomes alive  
3. All other live cells die (underpopulation or overpopulation)
4. All other dead cells remain dead

Available patterns:
- Random: Randomly distributed initial cells
- Glider: Moving pattern that travels across the grid
- Glider Gun: Produces a steady stream of gliders
- Oscillator: Patterns that oscillate between states
- Pulsar: Complex 3-period oscillator
- R-Pentomino: Chaotic pattern that evolves for 1000+ generations`,
	Run: func(cmd *cobra.Command, _ []string) {
		// Initialize logging and profiling
		InitLog()

		ctx := context.Background()
		InitProfile(ctx)

		// Get flags
		aliveColor, _ := cmd.Flags().GetString("alive-color")
		deadColor, _ := cmd.Flags().GetString("dead-color")
		aliveChar, _ := cmd.Flags().GetString("alive-char")
		deadChar, _ := cmd.Flags().GetString("dead-char")

		// Create configuration
		config := gameoflife.Config{
			AliveColor: aliveColor,
			DeadColor:  deadColor,
			AliveChar:  aliveChar,
			DeadChar:   deadChar,
		}

		// Create the Conway's Game of Life engine
		game := gameoflife.New(
			max(ui.DefaultHeight, 1),
			max(ui.DefaultWidth, 1),
			config,
		)

		// Run the UI with the engine
		if err := ui.RunModel("Conway's Game of Life", game, lang, refreshInterval); err != nil {
			slog.Error("Failed to run Conway's Game of Life", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(conwayGameOfLifeCmd)

	conwayGameOfLifeCmd.Flags().String("alive-char", "█", "Alive cell character")
	conwayGameOfLifeCmd.Flags().String("dead-char", " ", "Dead cell character")
	conwayGameOfLifeCmd.Flags().String("alive-color", "#00FF00", "Alive cell color (hex)")
	conwayGameOfLifeCmd.Flags().String("dead-color", "#000000", "Dead cell color (hex)")
}
