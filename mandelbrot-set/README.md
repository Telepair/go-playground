# Mandelbrot Set - Fractal Visualization

A beautiful Terminal User Interface (TUI) implementation of the Mandelbrot and Julia sets, showcasing the fascinating world of fractals through interactive visualization.

## Features

- **Mandelbrot Set**: Explore the classic fractal set with infinite complexity
- **Julia Set**: Switch to Julia set mode with customizable parameters
- **Interactive Navigation**: Pan, zoom, and explore the fractal landscape
- **Multiple Color Schemes**: 5 different color palettes for stunning visuals
- **Preset Locations**: Quick access to interesting fractal features
- **Bilingual Support**: English and Chinese interface
- **Real-time Calculation**: Dynamic fractal generation as you explore
- **Keyboard Controls**: Intuitive navigation without mouse dependency

## Mathematical Background

The Mandelbrot set is defined as the set of complex numbers `c` for which the sequence:

```
z₀ = 0
zₙ₊₁ = zₙ² + c
```

remains bounded (i.e., |z| ≤ 2) as n approaches infinity.

The Julia set, for a given complex parameter `c`, is defined similarly but starts with `z₀ = z` (the point being tested):

```
zₙ₊₁ = zₙ² + c
```

## Installation

```bash
# Navigate to the mandelbrot-set directory
cd mandelbrot-set

# Run the program
go run .
```

## Usage

### Command Line Options

```bash
# Basic usage
go run .

# Custom settings
go run . -zoom 2.0 -center-x -0.5 -center-y 0.0

# High iteration count with custom colors
go run . -max-iter 100 -color-scheme 2

# Julia set mode
go run . -julia -julia-c "0.285+0.01i"

# Chinese interface
go run . -lang cn
```

### Keyboard Controls

| Key                    | Action                                   |
| ---------------------- | ---------------------------------------- |
| `Arrow Keys` / `WASD`  | Pan around the fractal                   |
| `Shift + Arrow Keys`   | Fine panning                             |
| `+` / `=`              | Zoom in                                  |
| `-` / `_`              | Zoom out                                 |
| `M`                    | Toggle between Mandelbrot and Julia sets |
| `C`                    | Cycle through color schemes              |
| `I`                    | Increase maximum iterations              |
| `K`                    | Decrease maximum iterations              |
| `P`                    | Go to next preset location               |
| `L`                    | Toggle language (English/Chinese)        |
| `R`                    | Reset to default view                    |
| `Q` / `Ctrl+C` / `Esc` | Quit                                     |

### Color Schemes

1. **Classic**: Traditional black and white
2. **Hot**: Warm colors (red, orange, yellow)
3. **Cool**: Cool colors (blue, cyan, purple)
4. **Rainbow**: Full spectrum colors
5. **Grayscale**: Smooth grayscale gradient

### Preset Locations

The program includes several interesting preset locations:

- **Classic View**: The standard Mandelbrot set view
- **Seahorse Valley**: Beautiful seahorse-like structures
- **Lightning**: Electric-like branching patterns
- **Elephant Valley**: Elephant-trunk-like formations
- **Spiral**: Spiral arms and structures
- **Mini Mandelbrot**: Self-similar smaller copies
- **Feather**: Delicate feather-like patterns
- **Dragon**: Dragon-curve-like structures

## Technical Details

- **Language**: Go
- **UI Framework**: Bubble Tea (TUI)
- **Styling**: Lip Gloss
- **Complex Math**: Native Go complex128 type
- **Performance**: Optimized with string builders and efficient rendering

## Configuration

### Command Line Parameters

| Parameter       | Default         | Description                           |
| --------------- | --------------- | ------------------------------------- |
| `-rows`         | 30              | Number of rows in the display grid    |
| `-cols`         | 80              | Number of columns in the display grid |
| `-max-iter`     | 50              | Maximum number of iterations          |
| `-zoom`         | 1.0             | Initial zoom level                    |
| `-center-x`     | -0.5            | Initial center X coordinate           |
| `-center-y`     | 0.0             | Initial center Y coordinate           |
| `-color-scheme` | 0               | Color scheme (0-4)                    |
| `-julia`        | false           | Start in Julia set mode               |
| `-julia-c`      | "-0.7+0.27015i" | Julia set parameter                   |
| `-lang`         | "en"            | Language (en/cn)                      |

## Examples

### Exploring the Mandelbrot Set

1. Start the program: `go run .`
2. Use arrow keys to pan around
3. Use `+` and `-` to zoom in and out
4. Press `P` to jump to interesting locations
5. Press `C` to cycle through color schemes

### Julia Set Exploration

1. Start with Julia set: `go run . -julia`
2. Or toggle mode with `M` key
3. The Julia set uses a fixed parameter `c`
4. Different `c` values create different Julia sets

### High-Detail Rendering

For detailed exploration:

```bash
go run . -max-iter 200 -zoom 100 -center-x -0.7463 -center-y 0.1102
```

## Mathematical Interest

The Mandelbrot set exhibits several fascinating properties:

- **Self-similarity**: Zooming in reveals similar structures at different scales
- **Infinite complexity**: The boundary has infinite detail
- **Connectedness**: The entire set is connected
- **Julia sets**: Each point in the complex plane has an associated Julia set

## Performance Notes

- Calculation time increases with iteration count and zoom level
- Higher zoom levels may require more iterations for detail
- The program uses efficient algorithms but very high zoom levels will be slower
- Modern multi-core systems will benefit from parallel computation

## Contributing

Feel free to contribute improvements, additional color schemes, or new preset locations!

## License

This project is part of the go-playground repository and follows the same license terms.
