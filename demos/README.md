# Demos

[中文版本 / Chinese Version](README_CN.md)

This directory contains demonstration recordings and GIF files for the Go Playground projects.

## Recording Demonstrations

### Using Asciinema

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
   asciinema rec ./demos/cellular-automaton.cast --title "Cellular Automaton" --command "./bin/cellular-automaton"
   ```

3. Play the demo:

   ```bash
   asciinema play ./demos/cellular-automaton.cast
   ```

4. Upload to asciinema.org (optional):

   ```bash
   asciinema upload ./demos/cellular-automaton.cast
   ```
