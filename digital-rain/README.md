# Digital Rain

[Chinese Version / 中文版本](README_CN.md)

[Wikipedia - Matrix Digital Rain](https://en.wikipedia.org/wiki/Matrix_digital_rain)

A Terminal User Interface implementation of the famous Matrix digital rain effect.

## Features

- **Matrix-style Animation**: Characters fall vertically with trailing effects
- **Customizable Colors**: Configure drop, trail, and background colors
- **Variable Speed**: Adjust animation speed and drop characteristics
- **Character Sets**: Default base64 charset or customize your own
- **Bilingual Support**: Interface available in English and Chinese
- **Interactive Controls**: Pause, adjust parameters, and reset in real-time

## Installation

```bash
cd digital-rain
go build
```

## Usage

Run with default settings:

```bash
./digital-rain
```

### Command Line Options

- `-drop-color`: Drop color in hex format (default: "#00FF00")
- `-trail-color`: Trail color in hex format (default: "#008800")
- `-bg-color`: Background color in hex format (default: "#000000")
- `-charset`: Character set to use (default: base64 characters)
- `-min-speed`: Minimum drop speed (default: 1)
- `-max-speed`: Maximum drop speed (default: 5)
- `-drop-length`: Drop length (default: 10)
- `-lang`: Language (en/cn) (default: "en")
- `-profile`: Enable profiling and monitoring
- `-log-file`: Log file path for debugging

### Examples

Binary rain effect:

```bash
./digital-rain -charset '01'
```

White rain on dark background:

```bash
./digital-rain -drop-color '#FFFFFF' -trail-color '#888888'
```

Fast rain with long drops:

```bash
./digital-rain -max-speed 10 -drop-length 20
```

## Controls

- **Space/Enter**: Pause/Resume animation
- **+/-** or **↑/↓**: Increase/Decrease animation speed
- **d/D**: Increase/Decrease drop length
- **s/S**: Increase/Decrease maximum speed
- **r**: Reset animation
- **l**: Toggle language (English/Chinese)
- **q/Esc/Ctrl+C**: Quit

## Implementation Details

The digital rain effect consists of multiple "drops" falling down the screen:

1. **Drop System**: Each column has an independent drop with its own speed and length
2. **Trail Effect**: Characters fade gradually using color intensity mapping
3. **Character Variation**: Characters randomly change as they fall
4. **Spawn Control**: New drops appear at random intervals after previous ones disappear

## Performance

The application is optimized for smooth animation:

- Efficient grid rendering with string builders
- Pre-computed color styles for trail effects
- Concurrent-safe operations with mutex protection
- Configurable refresh rate (minimum 10ms)

## Testing

Run the test suite:

```bash
go test
```

Run benchmarks:

```bash
go test -bench=.
```
