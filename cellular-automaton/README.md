# Cellular Automaton

_[Chinese Version / 中文版本](README_CN.md)_

[Wikipedia - Cellular Automaton](https://en.wikipedia.org/wiki/Cellular_automaton)

A Terminal User Interface (TUI) implementation of 1D cellular automata with highly customizable rendering options and **automatic terminal size detection**.

[![asciicast](https://asciinema.org/a/723614.svg)](https://asciinema.org/a/723614)

## Features

- 🧬 **Multiple cellular automaton rules** (0-255) with pre-configured interesting rules
- 📐 **Auto-detect terminal size** - automatically adapts to your terminal dimensions
- 🎨 **Customizable appearance** - colors, characters, and visual styles
- 🔄 **Real-time simulation** with adjustable speed control
- 🌐 **Bilingual support** - English and Chinese interface
- 🔒 **Multiple boundary conditions** - periodic, fixed, and reflective
- ⚡ **High performance** with optimized rendering and ring buffer management

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
# Auto-detect terminal size (recommended)
./cellular-automaton
./cellular-automaton -rule 30
./cellular-automaton -rule 90 -alive-color "#00FF00" -dead-color "#FF0000"
./cellular-automaton -rule 110 -alive-char "●" -dead-char "○"
./cellular-automaton -rule 150
./cellular-automaton -rule 184 -alive-char '🚗'

# Override with specific dimensions
./cellular-automaton -rule 30 -rows 25 -cols 100
./cellular-automaton -rule 90 -rows 40 -cols 120
```

### Command Line Options

- `-rule <number>`: Cellular automaton rule number (0-255, default: 30)
- `-rows <number>`: Number of rows (auto-detect if not specified)
- `-cols <number>`: Number of columns (auto-detect if not specified)
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

## User Interface Layout

The interface features an optimized layout with status information at the top and controls at the bottom:

```
┌─────────────────────────────────────────┐
│          🧬 Cellular Automaton 🧬        │  ← Header
├─────────────────────────────────────────┤
│ 🧬 Rule: 30    ⚡ Gen: 42    🔄 Speed... │  ← Status (Top)
│ 🔒 Boundary    📐 Size       ▶️ Status   │
├─────────────────────────────────────────┤
│                                         │
│         █ █  ███ █  █ ███               │  ← Grid Content
│        ██████ ███████████               │
│       ███  █████     ████               │
│                                         │
├─────────────────────────────────────────┤
│ T Select Rule  B Boundary   +/- Speed   │  ← Controls (Bottom)
│ R Reset        L Language   Space/Q     │
└─────────────────────────────────────────┘
```

## Interesting Rules

- **Rule 30**: Chaos, pseudo-random patterns
- **Rule 90**: Sierpinski triangle pattern
- **Rule 110**: Turing complete, complex emergent behavior
- **Rule 150**: XOR pattern, create fractal structures
- **Rule 184**: Traffic simulation

## Technical Details

### Auto-Size Detection

The application automatically detects your terminal size and adapts the grid dimensions:

- **Rows**: Terminal height minus UI elements (header, status, controls)
- **Columns**: Terminal width minus padding
- **Dynamic resizing**: Automatically adjusts when terminal is resized
- **Manual override**: Use `-rows` and `-cols` to specify exact dimensions

### Boundary Types

- **Periodic**: The leftmost cell's left neighbor is the rightmost cell, and the rightmost cell's right neighbor is the leftmost cell (looping behavior)
- **Fixed**: The leftmost cell's left neighbor is always 0 (dead), and the rightmost cell's right neighbor is always 0 (dead)
- **Reflect**: The leftmost cell's left neighbor is itself, and the rightmost cell's right neighbor is itself (mirror behavior)
