# Conway's Game of Life

_[Chinese Version / 中文版本](README_CN.md)_

[Wikipedia - Conway's Game of Life](https://en.wikipedia.org/wiki/Conway's_Game_of_Life)

A Terminal User Interface (TUI) implementation of Conway's Game of Life with multiple predefined patterns and highly customizable rendering options.

[![asciicast](https://asciinema.org/a/723612.svg)](https://asciinema.org/a/723612)

## Features

- **Classic Game of Life Rules**: Faithful implementation of Conway's original cellular automaton
- **Multiple Starting Patterns**:
  - Random: Randomly distributed initial cells
  - Glider: The famous glider pattern that moves across the grid
  - Glider Gun: Gosper's glider gun that continuously produces gliders
  - Oscillator: Blinker patterns that oscillate between states
  - Pulsar: Period-3 oscillator with complex behavior
  - R-Pentomino: Chaotic pattern that evolves for over 1000 generations
- **Dual Boundary Conditions**:
  - Periodic: Wrapping edges (torus topology)
  - Fixed: Dead cells beyond boundaries
- **Enhanced User Interface**:
  - 🎮 Modern header with game branding
  - ⚡ Real-time status display with generation count and speed
  - 🎨 Interactive pattern switching
  - 🔄 Pause/resume functionality
  - 📐 Customizable cell rendering and colors
- **Real-time Controls**: Change patterns, boundary conditions, and speed without restart
- **Bilingual Support**: English and Chinese interface
- **Performance Optimized**: Efficient 2D grid computation and rendering

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd conway-game-of-life

# Build the application
go build -o conway-game-of-life
```

## Usage

### Basic Commands

```bash
# Run with default settings (random pattern)
./conway-game-of-life

# Start with a glider pattern
./conway-game-of-life -pattern glider

# Custom colors and characters
./conway-game-of-life -alive-char "🟢" -dead-char "⚫"
```

### Command Line Options

- `-rows <number>`: Number of rows in the grid (default: 30)
- `-cols <number>`: Number of columns in the grid (default: 60)
- `-alive-color <color>`: Alive cell color in hex format (default: #00FF00)
- `-dead-color <color>`: Dead cell color in hex format (default: #000000)
- `-alive-char <char>`: Character for alive cells (default: █)
- `-dead-char <char>`: Character for dead cells (default: space)
- `-lang <en/cn>`: Interface language (default: en)

### Example Commands

```bash
# On a large grid
./conway-game-of-life -rows 50 -cols 100

# Custom colors
./conway-game-of-life -alive-color "#FF0000" -dead-color "#000033"

# Emoji-based visualization
./conway-game-of-life -alive-char "🔴" -dead-char "⚪"

# Chinese interface
./conway-game-of-life -lang cn
```

## Controls

### Universal Controls

- **Space** or **Enter**: Pause/Resume the simulation
- **q** or **Ctrl+C**: Quit the application
- **l**: Toggle language (English/Chinese)

### Interactive Controls

- **p**: Cycle through different patterns (random → glider → glider-gun → oscillator → pulsar → pentomino)
- **b**: Toggle boundary conditions (periodic ↔ fixed)
- **+** or **=**: Increase speed (decrease refresh rate)
- **-** or **\_**: Decrease speed (increase refresh rate)

## Patterns

### Glider

A 5-cell pattern that moves diagonally across the grid every 4 generations.

```
  █
   █
███
```

### Glider Gun

Produces a steady stream of gliders. Demonstrates that Game of Life can have patterns with unlimited growth.

### Oscillators

- **Blinker**: Simple 2-period oscillator
- **Pulsar**: Complex 3-period oscillator with period 3

### R-Pentomino

A methuselah pattern that evolves chaotically for 1103 generations before stabilizing.

## Game Rules

Conway's Game of Life follows these simple rules:

1. **Survival**: Any live cell with 2 or 3 live neighbors survives to the next generation
2. **Birth**: Any dead cell with exactly 3 live neighbors becomes a live cell
3. **Death**: All other live cells die (underpopulation or overpopulation)
4. **Stasis**: All other dead cells remain dead

## Technical Details

### Boundary Conditions

- **Periodic**: The grid wraps around like a torus - cells at the edges interact with cells on the opposite side
- **Fixed**: Cells outside the grid boundaries are considered permanently dead

### Performance

- **Efficient Computation**: Optimized neighbor counting with boundary condition handling
- **Memory Management**: Pre-allocated grids with grid swapping to minimize allocations
- **Rendering Optimization**: Cached styled strings and efficient string building
- **Real-time Updates**: Dynamic parameter adjustment without restart

### Pattern Complexity Classes

- **Still Lifes**: Patterns that don't change (achieved after evolution)
- **Oscillators**: Patterns that return to their initial state after a fixed number of generations
- **Spaceships**: Patterns that translate themselves across the grid
- **Methuselahs**: Patterns that take a long time to stabilize
- **Infinite Growth**: Patterns that grow without bound

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source. Please check the license file for details.

## Interesting Facts

- Conway's Game of Life is Turing complete - it can simulate any computable function
- The Game of Life has been used to build computers, calculators, and even to simulate itself
- Some patterns take thousands of generations to stabilize or enter repeating cycles
- The glider gun was the first discovered pattern with infinite growth
