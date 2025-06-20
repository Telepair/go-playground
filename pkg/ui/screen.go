package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	defaultZeroValue = ' '
	// defaultBackground = lipgloss.Color("#1A1A2E")
	// defaultForeground = lipgloss.Color("#E0E0E0")

	screenStyle = lipgloss.NewStyle().Padding(1, 1, 1, 1)
)

// Screen represents a terminal screen buffer with styling capabilities
type Screen struct {
	rows int
	cols int // Store column count for consistency checks

	zeroValue       rune
	backgroundColor lipgloss.Color
	foregroundColor lipgloss.Color
	screenStyle     lipgloss.Style
	charColors      map[rune]lipgloss.Color
	charStyles      map[rune]lipgloss.Style

	data      [][]rune
	writeLine int
	viewLine  int
	buf       strings.Builder
	lineBuf   strings.Builder
}

// NewScreen creates a new screen with the specified dimensions
func NewScreen(rows, cols int) *Screen {
	gs := &Screen{
		rows:            rows,
		cols:            cols,
		zeroValue:       defaultZeroValue,
		backgroundColor: "",
		foregroundColor: "",
		screenStyle:     screenStyle,
		charColors:      make(map[rune]lipgloss.Color),
		charStyles:      make(map[rune]lipgloss.Style),
		writeLine:       0,
		viewLine:        0,
		buf:             strings.Builder{},
		lineBuf:         strings.Builder{},
	}
	gs.Reset()

	return gs
}

// SetHeight sets the height of the screen, preserving existing data where possible
func (gs *Screen) SetHeight(height int) {
	if height == gs.rows {
		return
	}
	if height < 0 {
		height = 0
	}
	data := make([][]rune, height)
	rows := min(height, gs.rows)
	for i := range rows {
		data[i] = make([]rune, gs.cols)
		copy(data[i], gs.data[i][:gs.cols])
	}
	for i := rows; i < height; i++ {
		data[i] = make([]rune, gs.cols)
		for j := range gs.cols {
			data[i][j] = gs.zeroValue
		}
	}
	gs.data = data
	gs.rows = height
}

// SetWidth sets the width of the screen, preserving existing data where possible
func (gs *Screen) SetWidth(width int) {
	if width == gs.cols {
		return
	}
	if width < 0 {
		width = 0
	}
	data := make([][]rune, gs.rows)
	for i := range gs.rows {
		data[i] = make([]rune, width)
		cols := min(width, gs.cols)
		for j := range cols {
			data[i][j] = gs.data[i][j]
		}
		for j := cols; j < width; j++ {
			data[i][j] = gs.zeroValue
		}
	}
	gs.data = data
	gs.cols = width
}

// SetSize sets both width and height of the screen, preserving existing data where possible
func (gs *Screen) SetSize(width, height int) {
	if width == gs.cols && height == gs.rows {
		return
	}
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	data := make([][]rune, height)
	rows := min(height, gs.rows)
	for i := range rows {
		data[i] = make([]rune, width)
		cols := min(width, gs.cols)
		for j := range cols {
			data[i][j] = gs.data[i][j]
		}
		for j := cols; j < width; j++ {
			data[i][j] = gs.zeroValue
		}
	}
	for i := rows; i < height; i++ {
		data[i] = make([]rune, width)
		for j := range width {
			data[i][j] = gs.zeroValue
		}
	}
	gs.data = data
	gs.rows = height
	gs.cols = width
}

// SetZeroValue sets the default character used for empty cells
func (gs *Screen) SetZeroValue(zeroValue rune) {
	gs.zeroValue = zeroValue
}

// Reset resets the entire screen
func (gs *Screen) Reset() {
	if gs.data == nil {
		gs.data = make([][]rune, gs.rows)
	}
	for i := range gs.rows {
		if gs.data[i] == nil {
			gs.data[i] = make([]rune, gs.cols)
		}
		for j := range gs.cols {
			gs.data[i][j] = gs.zeroValue
		}
	}
}

// Fill fills the entire screen with the specified character
func (gs *Screen) Fill(char rune) {
	if gs.data == nil {
		gs.data = make([][]rune, gs.rows)
	}
	for i := range gs.rows {
		if gs.data[i] == nil {
			gs.data[i] = make([]rune, gs.cols)
		}
		for j := range gs.cols {
			gs.data[i][j] = char
		}
	}
}

// SetBackground sets the background color for the screen
func (gs *Screen) SetBackground(background lipgloss.Color) {
	if background == "" {
		return
	}
	gs.backgroundColor = background
	gs.screenStyle = gs.screenStyle.Background(background)
}

// SetForeground sets the foreground color for the screen
func (gs *Screen) SetForeground(foreground lipgloss.Color) {
	if foreground == "" {
		return
	}
	gs.foregroundColor = foreground
	gs.screenStyle = gs.screenStyle.Foreground(foreground)
}

// SetCharColor sets a specific color for a character when rendered
func (gs *Screen) SetCharColor(char rune, color lipgloss.Color) {
	if color == "" {
		return
	}
	if char == 0 {
		return
	}
	gs.charColors[char] = color
	gs.charStyles[char] = lipgloss.NewStyle().Foreground(color)
}

// SetData sets the screen data from a 2D rune array
func (gs *Screen) SetData(data [][]rune) {
	if gs.data == nil {
		gs.data = make([][]rune, gs.rows)
	}
	rows := min(len(data), gs.rows)
	for i := range rows {
		if gs.data[i] == nil {
			gs.data[i] = make([]rune, gs.cols)
		}
		cols := min(len(data[i]), gs.cols)
		for j := range cols {
			gs.data[i][j] = data[i][j]
		}
		for j := cols; j < gs.cols; j++ {
			gs.data[i][j] = gs.zeroValue
		}
	}
	for i := rows; i < gs.rows; i++ {
		if gs.data[i] == nil {
			gs.data[i] = make([]rune, gs.cols)
		}
		for j := range gs.cols {
			gs.data[i][j] = gs.zeroValue
		}
	}
	gs.writeLine = gs.rows - 1
	gs.viewLine = 0
}

// Append adds a new row to the screen, scrolling if necessary
func (gs *Screen) Append(row []rune) {
	// Ensure data is initialized
	if gs.data == nil {
		gs.Reset()
	}

	if gs.writeLine >= gs.rows {
		gs.writeLine = 0
	}

	cols := min(len(row), gs.cols)
	if gs.data[gs.writeLine] == nil {
		gs.data[gs.writeLine] = make([]rune, gs.cols)
	}
	copy(gs.data[gs.writeLine][:cols], row[:cols])
	for j := cols; j < gs.cols; j++ {
		gs.data[gs.writeLine][j] = gs.zeroValue
	}
	gs.writeLine++
	gs.viewLine = gs.writeLine
}

// View renders the screen content as a styled string
func (gs *Screen) View() string {
	gs.buf.Reset()
	var line int
	for i := range gs.rows {
		line = (gs.viewLine + i) % gs.rows
		gs.lineBuf.Reset()
		for j := range gs.cols {
			style, ok := gs.charStyles[gs.data[line][j]]
			if ok {
				gs.lineBuf.WriteString(style.Render(string(gs.data[line][j])))
			} else {
				gs.lineBuf.WriteRune(gs.data[line][j])
			}
		}
		gs.buf.WriteString(gs.lineBuf.String())
		if i < gs.rows-1 {
			gs.buf.WriteRune('\n')
		}
	}
	return gs.screenStyle.Render(gs.buf.String())
}
