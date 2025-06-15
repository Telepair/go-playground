# Cellular Automaton

_[Chinese Version / 中文版本](README_CN.md)_

[Wikipedia - Cellular Automaton](https://en.wikipedia.org/wiki/Cellular_automaton)

A Terminal User Interface (TUI) implementation of 1D cellular automata with highly customizable rendering options.

[![asciicast](https://asciinema.org/a/723316.svg)](https://asciinema.org/a/723316)

More demos: [Demos](../demos/cellular-automaton/README.md)

## Features

- **Rule-based Generation**: Support for all 256 elementary cellular automaton rules (0-255)
- **Dual Operating Modes**:
  - Finite mode: Run for a specific number of steps
  - Infinite mode: Continuous generation with real-time visualization
- **Enhanced User Interface**:
  - 🧬 Modern header with icons and styled branding
  - ⚡ Status bar with organized information groups and visual indicators
  - 🎮 Categorized control panel with grouped commands
  - 📐 Card-like layout with rounded borders and enhanced styling
  - 🔄 Real-time status indicators (Running/Paused with visual feedback)
- **Auto Window Detection**: Automatically detects terminal size or allows manual specification
- **Highly Customizable Rendering**:
  - Configurable cell size (1-3x)
  - Custom colors (hex format)
  - Custom characters for rendering
- **Dynamic Refresh**: Configurable refresh rate (default 1s, minimum 1ms)
- **Flexible Window Size**: Auto-detection or manual specification using WIDTHxHEIGHT format (e.g., 100x80)
- **Bilingual Support**: English and Chinese interface with improved formatting

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd cellular-automaton

# Build the application
go build -o cellular-automaton
```

## Usage

### Basic Commands

```bash
# Run with default settings (Rule 30, auto window size)
./cellular-automaton

# Specific rule with custom parameters
./cellular-automaton -rule 110 -steps 100 -size 120x60

# Infinite mode with auto window detection
./cellular-automaton -rule 30 -steps 0 -size auto

# Custom rendering style
./cellular-automaton -rule 90 -cellsize 3 -alive-char "●" -dead-char "○"
```

### Command Line Options

- `-rule <number>`: Cellular automaton rule (0-255, default: 30)
- `-steps <number>`: Number of steps (0 or negative for infinite mode, default: 50)
- `-size <size>`: Window size (format: WIDTHxHEIGHT, e.g.: 100x80, or 'auto' for auto-detection, default: auto)
- `-cellsize <size>`: Cell rendering size (1-3, default: 2)
- `-alive-color <color>`: Alive cell color in hex format (default: #FFFFFF)
- `-dead-color <color>`: Dead cell color in hex format (default: #000000)
- `-alive-char <char>`: Character for alive cells (default: █)
- `-dead-char <char>`: Character for dead cells (default: space)
- `-refresh <seconds>`: Refresh rate in seconds, minimum 0.001 (default: 0.1)
- `-boundary <type>`: Boundary condition type (periodic/fixed/reflect, default: periodic)
- `-lang <en/cn>`: Interface language (default: en)

### Makefile

```bash
# Show help
make help

# Build the cellular automaton
make build

# Run demos
make cellular-automaton-basic
make cellular-automaton-sierpinski
make cellular-automaton-turing
make cellular-automaton-traffic
make cellular-automaton-infinite
make cellular-automaton-colorful
make cellular-automaton-fixed
make cellular-automaton-periodic
make cellular-automaton-reflect
```

### Example Commands

```bash
# Rule 30 with auto-detected window size
./cellular-automaton -rule 30 -size auto

# Infinite mode with fast refresh
./cellular-automaton -rule 30 -steps 0 -refresh 0.1

# Custom characters for ASCII art style
./cellular-automaton -rule 184 -alive-char "■" -dead-char "□" -cellsize 1

# Large window size for detailed patterns
./cellular-automaton -rule 110 -size 200x100 -steps 150

# Fixed boundary conditions (no wrapping)
./cellular-automaton -rule 30 -boundary fixed

# Reflective boundary conditions
./cellular-automaton -rule 110 -boundary reflect

# Chinese interface
./cellular-automaton -rule 30 -lang cn
```

## Controls

### All Modes

- **q** or **Ctrl+C**: Quit the application
- **l**: Toggle language (English/Chinese)

### Infinite Mode Only

- **Space** or **Enter**: Pause/Resume the simulation

### Advanced Controls (Infinite Mode)

- **r**: Reset simulation to initial state
- **+** or **=**: Increase refresh rate (make simulation faster)
- **-** or **\_**: Decrease refresh rate (make simulation slower)
- **1**, **2**, **3**: Change cell rendering size (1x, 2x, 3x)

## Interesting Rules to Try

- **Rule 30**: Chaotic, pseudo-random patterns
- **Rule 90**: Sierpinski triangle pattern
- **Rule 110**: Turing complete, complex emergent behavior
- **Rule 184**: Traffic flow simulation
- **Rule 150**: XOR pattern, creates fractal structures

## Technical Details

### Boundary Conditions

The cellular automaton supports three boundary condition types:

- **Periodic (default)**: The left neighbor of the leftmost cell is the rightmost cell, and the right neighbor of the rightmost cell is the leftmost cell (wrapping behavior)
- **Fixed**: The left neighbor of the leftmost cell is always 0 (dead), and the right neighbor of the rightmost cell is always 0 (dead)
- **Reflect**: The left neighbor of the leftmost cell is the cell itself, and the right neighbor of the rightmost cell is the cell itself (mirror behavior)

### Window Size Detection

The application can automatically detect your terminal window size:

- Use `-size auto` (default) for automatic detection
- The application reserves space for UI elements (header, controls)
- Minimum height is enforced to ensure proper display
- Manual size specification overrides auto-detection

### Performance

- **Refresh Rate**: Supports high-frequency updates (minimum 1ms)
- **Memory Efficient**: Only stores necessary grid data, with zero-allocation cell computation
- **Optimized Rendering**: Uses pre-computed styled strings and efficient string building
- **Buffer Reuse**: Minimizes garbage collection through buffer pooling
- **Real-time Controls**: Dynamic adjustment of refresh rate and cell size without restart

### Pattern Analysis

Different rules produce distinct pattern types:

- **Class 1**: Evolves to homogeneous state (e.g., Rule 0, 32)
- **Class 2**: Evolves to simple periodic structures (e.g., Rule 4, 108)
- **Class 3**: Chaotic, aperiodic behavior (e.g., Rule 30, 45)
- **Class 4**: Complex, localized structures (e.g., Rule 110, 124)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source. Please check the license file for details.
