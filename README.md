# Go Playground

_[中文版本 / Chinese Version](README_CN.md)_

A collection of interesting programs implemented in Go language. Each sub-project is an independent program that demonstrates different programming concepts, algorithms, or fascinating implementations.

## Project List

### 🧬 [Cellular Automaton](./cellular-automaton/)

An interactive one-dimensional cellular automaton program with a beautiful terminal user interface built using Bubble Tea.

[Wikipedia - Cellular Automaton](https://en.wikipedia.org/wiki/Cellular_automaton)

**Demo:**

[![asciicast](https://asciinema.org/a/723524.svg)](https://asciinema.org/a/723524)

### 🎮 [Conway's Game of Life](./conway-game-of-life/)

A terminal user interface (TUI) implementation of Conway's Game of Life with multiple predefined patterns and highly customizable rendering options.

[Wikipedia - Conway's Game of Life](https://en.wikipedia.org/wiki/Conway's_Game_of_Life)

**Demo:**

[![asciicast](https://asciinema.org/a/723376.svg)](https://asciinema.org/a/723376)

### 📊 [Mandelbrot Set](./mandelbrot-set/)

An interactive Mandelbrot set terminal user interface (TUI) implementation.

[Wikipedia - Mandelbrot Set](https://en.wikipedia.org/wiki/Mandelbrot_set)

**Demo:**

[![asciicast](https://asciinema.org/a/723615.svg)](https://asciinema.org/a/723615)

## Project Structure

```
go-playground/
├── README.md                    # Main project documentation
├── cellular-automaton/          # Cellular Automaton
├── conway-game-of-life/         # Conway Game of Life
├── LICENSE                     # Project license
└── .gitignore                 # Git ignore file
```

## Using Asciinema to record demos

1. Install asciinema:

   ```bash
   # macOS
   brew install asciinema

   # Linux
   pip install asciinema
   ```

2. Record a demo:

   ```bash
   # Start recording
   # Note: After the program finishes running, press 'Q' to quit the program and complete the recording

   # Cellular Automaton
   asciinema rec ./demos/cellular-automaton.cast --title "Cellular Automaton" --command "./bin/cellular-automaton"

   # Conway Game of Life
   asciinema rec ./demos/conway-game-of-life.cast --title "Conway Game of Life" --command "./bin/conway-game-of-life"
   ```

3. Play the demo:

   ```bash
   # Cellular Automaton
   asciinema play ./demos/cellular-automaton.cast

   # Conway Game of Life
   asciinema play ./demos/conway-game-of-life.cast
   ```

4. Upload to asciinema.org (optional):

   ```bash
   # Cellular Automaton
    asciinema upload ./demos/cellular-automaton.cast

   # Conway Game of Life
   asciinema upload ./demos/conway-game-of-life.cast
   ```

## Technical Features

- **Modern Go Development**: Uses the latest features of Go 1.24+
- **Elegant User Interface**: Beautiful terminal interfaces built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Independent Module Design**: Each sub-project has its own `go.mod` for easy management and usage
- **Clear Code Structure**: Focus on code readability and maintainability
- **Comprehensive Documentation**: Each project comes with complete usage instructions and examples

## Planned Projects

Interesting projects that may be added in the future:

- 🧮 **Mandelbrot Set** - Mandelbrot set visualization
- 🎵 **Audio Visualizer** - Audio spectrum visualization
- 🌊 **Wave Function Collapse** - Wave Function Collapse algorithm
- 🎲 **Random Walk** - Random walk visualization
- 📊 **Data Structures Visualization** - Data structure visualization
- 🔍 **Algorithm Visualization** - Sorting and searching algorithm visualization

## Contributing

Issues and Pull Requests are welcome! If you have interesting ideas or find bugs, please feel free to contact us.

### Contribution Guidelines

1. Fork this project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Requirements

- Go 1.24.4 or higher
- Unicode-enabled terminal (modern terminals like iTerm2, Windows Terminal, etc. are recommended)

## License

This project is licensed under [License Name]. See the [LICENSE](LICENSE) file for details.

---

**Enjoy exploring various interesting concepts with Go!** 🚀
