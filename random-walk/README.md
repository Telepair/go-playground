# Random Walk Visualization

[Chinese Version / 中文版本](README_CN.md)

[Wikipedia - Random Walk](https://en.wikipedia.org/wiki/Random_walk)

A terminal-based visualization of various random walk algorithms, implemented in Go using the Bubble Tea framework.

[![asciicast](https://asciinema.org/a/723662.svg)](https://asciinema.org/a/723662)

## Features

- **Multiple Walk Modes**:

  - **Single Walker**: Classic random walk with one particle
  - **Multi Walker**: Multiple particles walking simultaneously
  - **Trail Mode**: Single walker with visible trail
  - **Brownian Motion**: Simulates Brownian motion with continuous movement
  - **Self-Avoiding Walk**: Walker cannot revisit previously visited positions
  - **Lévy Flight**: Random walk with occasional long jumps

- **Interactive Controls**:
  - Real-time visualization with adjustable speed
  - Pause/resume functionality
  - Dynamic walker count adjustment (for multi-walker modes)
  - Configurable trail length
  - Bilingual support (English/Chinese)

## Installation

### Prerequisites

- Go 1.22 or higher
- Terminal with Unicode support

### Building from Source

```bash
# Clone the repository
git clone https://github.com/telepair/go-playground.git
cd go-playground

# Build the random walk visualization
make build-random-walk

# Or build directly
go build -o ./bin/random-walk ./random-walk
```

## Usage

### Basic Usage

```bash
# Run with default settings
./bin/random-walk

# Or use make
make random-walk
```

### Command Line Options

```bash
./bin/random-walk [options]

Options:
  -walker-color string    Walker color in hex format (default "#FF00FF")
  -trail-color string     Trail color in hex format (default "#0088FF")
  -empty-color string     Empty cell color in hex format (default "#000000")
  -walker-char string     Character for walker (default "●")
  -trail-char string      Character for trail (default "·")
  -empty-char string      Character for empty cells (default " ")
  -lang string           Language: en or cn (default "en")
  -profile               Enable profiling and monitoring
  -profile-port int      Profiling server port (default 6060)
  -log-file string       Log file path for debugging
```

### Examples

```bash
# Use custom walker and trail characters
./bin/random-walk -walker-char '🐾' -trail-char '·'

# Custom colors
./bin/random-walk -walker-color '#FF00FF' -trail-color '#00FFFF'

# Run in Chinese
./bin/random-walk -lang cn

# Enable profiling
./bin/random-walk -profile -log-file debug.log
```

## Controls

| Key                | Action                                              |
| ------------------ | --------------------------------------------------- |
| `M`                | Cycle through walk modes                            |
| `W/w`              | Increase/decrease walker count (multi-walker modes) |
| `T/t`              | Increase/decrease trail length (trail modes)        |
| `+/-` or `↑/↓`     | Speed up/slow down                                  |
| `Space` or `Enter` | Pause/resume                                        |
| `L`                | Switch language (English/Chinese)                   |
| `R`                | Reset simulation                                    |
| `Q` or `Esc`       | Quit                                                |

## Walk Modes Explained

### Single Walker

A classic random walk where a single particle moves randomly in 8 directions (including diagonals).

### Multi Walker

Multiple particles walking simultaneously, each with a unique color. Useful for studying collision and coverage patterns.

### Trail Mode

Shows the path taken by a single walker, with the trail gradually fading over time.

### Brownian Motion

Simulates Brownian motion with continuous movement in random directions and distances.

### Self-Avoiding Walk

The walker cannot revisit positions it has already visited. The walk may become stuck if no valid moves are available.

### Lévy Flight

A random walk where the walker occasionally makes long jumps, simulating Lévy flight patterns found in nature.

## Technical Details

### Implementation

- Written in Go using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework
- Uses [Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling
- Implements efficient grid rendering with minimal allocations
- Supports terminal resizing and maintains aspect ratio

### Performance

- Optimized for smooth animation at 50ms refresh rate
- Efficient trail rendering using intensity decay
- Minimal memory allocations during rendering

## Development

### Running Tests

```bash
# Run all tests
make test

# Run benchmarks
make bench
```

### Project Structure

```
random-walk/
├── main.go          # Entry point
├── config.go        # Configuration and constants
├── walk.go          # Core random walk logic
├── ui.go            # UI and interaction logic
├── styles.go        # Visual styles and rendering
├── walk_test.go     # Unit tests
├── README.md        # English documentation
└── README_CN.md     # Chinese documentation
```

## License

This project is part of the go-playground collection and follows the same license terms.
