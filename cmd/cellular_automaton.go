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

	cellularautomaton "github.com/telepair/go-playground/engine/cellularautomaton"
	"github.com/telepair/go-playground/pkg/ui"
)

// cellularAutomatonCmd represents the cellular automaton command
var cellularAutomatonCmd = &cobra.Command{
	Use:   "cellular-automaton",
	Short: "Run a 1D cellular automaton simulation",
	Long: `Run a 1D cellular automaton simulation with various rules and boundary conditions.

The cellular automaton is a mathematical model of a grid of cells, each of which 
can be in one of a finite number of states. The state of each cell evolves over 
time according to a set of rules based on the states of neighboring cells.

Example rules:
- Rule 30: Chaotic pattern generator
- Rule 90: Sierpinski triangle pattern
- Rule 110: Complex patterns (proven to be Turing complete)
- Rule 184: Traffic flow simulation`,
	Run: func(cmd *cobra.Command, _ []string) {
		// Initialize logging and profiling
		InitLog()

		ctx := context.Background()
		InitProfile(ctx)

		// Get flags
		rule, _ := cmd.Flags().GetInt("rule")
		boundary, _ := cmd.Flags().GetInt("boundary")

		// Create the cellular automaton engine
		ca := cellularautomaton.New(
			rule,
			max(ui.DefaultHeight, 1),
			max(ui.DefaultWidth, 1),
			cellularautomaton.BoundaryType(boundary),
		)

		// Run the UI with the engine
		if err := ui.RunModel("Cellular Automaton", ca, lang, refreshInterval); err != nil {
			slog.Error("Failed to run cellular automaton", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cellularAutomatonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cellularAutomatonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cellularAutomatonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Local flags specific to cellular automaton
	cellularAutomatonCmd.Flags().Int("rule", 30, "Cellular automaton rule number (0-255)")
	cellularAutomatonCmd.Flags().Int("boundary", 0, "Boundary condition type (0=Periodic, 1=Fixed, 2=Reflect)")
}
