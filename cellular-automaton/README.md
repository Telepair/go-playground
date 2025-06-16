# Cellular Automaton

_[Chinese Version / 中文版本](README_CN.md)_

[Wikipedia - Cellular Automaton](https://en.wikipedia.org/wiki/Cellular_automaton)

A Terminal User Interface (TUI) implementation of 1D cellular automata with highly customizable rendering options.

[![asciicast](https://asciinema.org/a/723358.svg)](https://asciinema.org/a/723358)

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
./cellular-automaton
./cellular-automaton -rule 30
./cellular-automaton -rule 90 -alive-color "#00FF00" -dead-color "#FF0000"
./cellular-automaton -rule 110 -alive-char "●" -dead-char "○"
./cellular-automaton -rule 150
./cellular-automaton -rule 184 -alive-char '🚗' -rows 30 -cols 80
```

### Command Line Options

- `-rule <number>`: Cellular automaton rule number (0-255, default: 30)
- `-rows <number>`: Number of rows (default: 30)
- `-cols <number>`: Number of columns (default: 80)
- `-alive-color <color>`: Alive cell color (hex)
- `-dead-color <color>`: Dead cell color (hex)
- `-alive-char <char>`: Alive cell character
- `-dead-char <char>`: Dead cell character
- `-lang <en/cn>`: Language (default: en)

## Control Keys

- `t`: Toggle rule selection modal (T for "Type" rule)
- `b`: Toggle boundary selection modal (B for "Boundary" selection)
- `r`: Reset simulation to initial state
- `l`: Toggle language (English/Chinese)
- `+` or `=`: Increase refresh rate (speed up simulation)
- `-` or `_`: Decrease refresh rate (slow down simulation)
- `space` or `enter`: Pause/resume simulation
- `q` or `Ctrl+C`: Quit application

## Interesting Rules

- **Rule 30**: Chaos, pseudo-random patterns
- **Rule 90**: Sierpinski triangle pattern
- **Rule 110**: Turing complete, complex emergent behavior
- **Rule 150**: XOR pattern, create fractal structures
- **Rule 184**: Traffic simulation

## Technical Details

- **Boundary Types**:
  - **Periodic**: The leftmost cell's left neighbor is the rightmost cell, and the rightmost cell's right neighbor is the leftmost cell (looping behavior)
  - **Fixed**: The leftmost cell's left neighbor is always 0 (dead), and the rightmost cell's right neighbor is always 0 (dead)
  - **Reflect**: The leftmost cell's left neighbor is itself, and the rightmost cell's right neighbor is itself (mirror behavior)
