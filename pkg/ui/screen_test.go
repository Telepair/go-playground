package ui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

// TestNewScreen tests the creation of a new screen
func TestNewScreen(t *testing.T) {
	tests := []struct {
		name     string
		rows     int
		cols     int
		expected struct {
			rows int
			cols int
		}
	}{
		{
			name: "Normal dimensions",
			rows: 10,
			cols: 20,
			expected: struct {
				rows int
				cols int
			}{rows: 10, cols: 20},
		},
		{
			name: "Small dimensions",
			rows: 1,
			cols: 1,
			expected: struct {
				rows int
				cols int
			}{rows: 1, cols: 1},
		},
		{
			name: "Large dimensions",
			rows: 100,
			cols: 200,
			expected: struct {
				rows int
				cols int
			}{rows: 100, cols: 200},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := NewScreen(tt.rows, tt.cols)
			assert.NotNil(t, screen)
			assert.Equal(t, tt.expected.rows, screen.rows)
			assert.Equal(t, tt.expected.cols, screen.cols)
			assert.Equal(t, defaultZeroValue, screen.zeroValue)
			assert.NotNil(t, screen.data)
			assert.Len(t, screen.data, tt.expected.rows)
			if tt.expected.rows > 0 {
				assert.Len(t, screen.data[0], tt.expected.cols)
			}
		})
	}
}

// TestSetHeight tests the SetHeight method
func TestSetHeight(t *testing.T) {
	tests := []struct {
		name        string
		initialRows int
		initialCols int
		newHeight   int
		checkData   bool
	}{
		{
			name:        "Increase height",
			initialRows: 5,
			initialCols: 10,
			newHeight:   10,
			checkData:   true,
		},
		{
			name:        "Decrease height",
			initialRows: 10,
			initialCols: 10,
			newHeight:   5,
			checkData:   true,
		},
		{
			name:        "Same height",
			initialRows: 10,
			initialCols: 10,
			newHeight:   10,
			checkData:   false,
		},
		{
			name:        "Zero height",
			initialRows: 10,
			initialCols: 10,
			newHeight:   0,
			checkData:   false,
		},
		{
			name:        "Negative height",
			initialRows: 10,
			initialCols: 10,
			newHeight:   -5,
			checkData:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := NewScreen(tt.initialRows, tt.initialCols)

			// Fill with test data
			if tt.checkData {
				for i := 0; i < tt.initialRows && i < len(screen.data); i++ {
					for j := 0; j < tt.initialCols && j < len(screen.data[i]); j++ {
						screen.data[i][j] = rune('A' + i)
					}
				}
			}

			screen.SetHeight(tt.newHeight)

			expectedHeight := tt.newHeight
			if tt.newHeight < 0 {
				expectedHeight = 0
			}
			assert.Equal(t, expectedHeight, screen.rows)
			assert.Len(t, screen.data, expectedHeight)

			// Check data preservation
			if tt.checkData && expectedHeight > 0 {
				preservedRows := min(tt.initialRows, expectedHeight)
				for i := 0; i < preservedRows; i++ {
					assert.Equal(t, rune('A'+i), screen.data[i][0])
				}
			}
		})
	}
}

// TestSetWidth tests the SetWidth method
func TestSetWidth(t *testing.T) {
	tests := []struct {
		name        string
		initialRows int
		initialCols int
		newWidth    int
		checkData   bool
	}{
		{
			name:        "Increase width",
			initialRows: 5,
			initialCols: 10,
			newWidth:    20,
			checkData:   true,
		},
		{
			name:        "Decrease width",
			initialRows: 5,
			initialCols: 20,
			newWidth:    10,
			checkData:   true,
		},
		{
			name:        "Same width",
			initialRows: 5,
			initialCols: 10,
			newWidth:    10,
			checkData:   false,
		},
		{
			name:        "Zero width",
			initialRows: 5,
			initialCols: 10,
			newWidth:    0,
			checkData:   false,
		},
		{
			name:        "Negative width",
			initialRows: 5,
			initialCols: 10,
			newWidth:    -5,
			checkData:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := NewScreen(tt.initialRows, tt.initialCols)

			// Fill with test data
			if tt.checkData {
				for i := 0; i < tt.initialRows && i < len(screen.data); i++ {
					for j := 0; j < tt.initialCols && j < len(screen.data[i]); j++ {
						screen.data[i][j] = rune('0' + j)
					}
				}
			}

			screen.SetWidth(tt.newWidth)

			expectedWidth := tt.newWidth
			if tt.newWidth < 0 {
				expectedWidth = 0
			}
			assert.Equal(t, expectedWidth, screen.cols)

			// Check data preservation
			if tt.checkData && tt.initialRows > 0 {
				preservedCols := min(tt.initialCols, expectedWidth)
				for i := 0; i < tt.initialRows; i++ {
					assert.Len(t, screen.data[i], expectedWidth)
					for j := 0; j < preservedCols; j++ {
						assert.Equal(t, rune('0'+j), screen.data[i][j])
					}
				}
			}
		})
	}
}

// TestSetSize tests the SetSize method
func TestSetSize(t *testing.T) {
	tests := []struct {
		name      string
		initRows  int
		initCols  int
		newWidth  int
		newHeight int
	}{
		{
			name:      "Increase both dimensions",
			initRows:  5,
			initCols:  10,
			newWidth:  20,
			newHeight: 10,
		},
		{
			name:      "Decrease both dimensions",
			initRows:  10,
			initCols:  20,
			newWidth:  5,
			newHeight: 5,
		},
		{
			name:      "Mixed changes",
			initRows:  10,
			initCols:  10,
			newWidth:  5,
			newHeight: 20,
		},
		{
			name:      "No change",
			initRows:  10,
			initCols:  10,
			newWidth:  10,
			newHeight: 10,
		},
		{
			name:      "Negative values",
			initRows:  10,
			initCols:  10,
			newWidth:  -5,
			newHeight: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := NewScreen(tt.initRows, tt.initCols)
			screen.SetSize(tt.newWidth, tt.newHeight)

			expectedWidth := max(0, tt.newWidth)
			expectedHeight := max(0, tt.newHeight)

			assert.Equal(t, expectedWidth, screen.cols)
			assert.Equal(t, expectedHeight, screen.rows)
			assert.Len(t, screen.data, expectedHeight)
			if expectedHeight > 0 {
				assert.Len(t, screen.data[0], expectedWidth)
			}
		})
	}
}

// TestSetZeroValue tests the SetZeroValue method
func TestSetZeroValue(t *testing.T) {
	screen := NewScreen(5, 5)

	// Change zero value
	newZeroValue := '*'
	screen.SetZeroValue(newZeroValue)
	assert.Equal(t, newZeroValue, screen.zeroValue)

	// Reset should use new zero value
	screen.Reset()
	for i := 0; i < screen.rows; i++ {
		for j := 0; j < screen.cols; j++ {
			assert.Equal(t, newZeroValue, screen.data[i][j])
		}
	}
}

// TestReset tests the Reset method
func TestReset(t *testing.T) {
	screen := NewScreen(5, 5)

	// Fill with data
	screen.Fill('X')

	// Verify filled
	for i := 0; i < screen.rows; i++ {
		for j := 0; j < screen.cols; j++ {
			assert.Equal(t, 'X', screen.data[i][j])
		}
	}

	// Reset
	screen.Reset()

	// Verify reset to zero value
	for i := 0; i < screen.rows; i++ {
		for j := 0; j < screen.cols; j++ {
			assert.Equal(t, screen.zeroValue, screen.data[i][j])
		}
	}
}

// TestFill tests the Fill method
func TestFill(t *testing.T) {
	tests := []struct {
		name     string
		rows     int
		cols     int
		fillChar rune
	}{
		{
			name:     "Fill with letter",
			rows:     5,
			cols:     10,
			fillChar: 'A',
		},
		{
			name:     "Fill with number",
			rows:     3,
			cols:     3,
			fillChar: '9',
		},
		{
			name:     "Fill with symbol",
			rows:     4,
			cols:     8,
			fillChar: '@',
		},
		{
			name:     "Fill with unicode",
			rows:     2,
			cols:     2,
			fillChar: '★',
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := NewScreen(tt.rows, tt.cols)
			screen.Fill(tt.fillChar)

			for i := 0; i < tt.rows; i++ {
				for j := 0; j < tt.cols; j++ {
					assert.Equal(t, tt.fillChar, screen.data[i][j])
				}
			}
		})
	}
}

// TestSetColors tests color setting methods
func TestSetColors(t *testing.T) {
	screen := NewScreen(5, 5)

	// Test SetBackground
	bgColor := lipgloss.Color("#FF0000")
	screen.SetBackground(bgColor)
	assert.Equal(t, bgColor, screen.backgroundColor)

	// Test SetForeground
	fgColor := lipgloss.Color("#00FF00")
	screen.SetForeground(fgColor)
	assert.Equal(t, fgColor, screen.foregroundColor)

	// Test SetCharColor
	charColor := lipgloss.Color("#0000FF")
	screen.SetCharColor('A', charColor)
	assert.Equal(t, charColor, screen.charColors['A'])
	assert.NotNil(t, screen.charStyles['A'])

	// Test empty color (should not set)
	screen.SetBackground("")
	assert.Equal(t, bgColor, screen.backgroundColor) // Should remain unchanged

	screen.SetCharColor(0, charColor) // Zero char should not be set
	_, exists := screen.charColors[0]
	assert.False(t, exists)
}

// TestSetData tests the SetData method
func TestSetData(t *testing.T) {
	tests := []struct {
		name       string
		screenRows int
		screenCols int
		data       [][]rune
	}{
		{
			name:       "Exact size data",
			screenRows: 3,
			screenCols: 3,
			data: [][]rune{
				{'A', 'B', 'C'},
				{'D', 'E', 'F'},
				{'G', 'H', 'I'},
			},
		},
		{
			name:       "Smaller data",
			screenRows: 5,
			screenCols: 5,
			data: [][]rune{
				{'A', 'B'},
				{'C', 'D'},
			},
		},
		{
			name:       "Larger data",
			screenRows: 2,
			screenCols: 2,
			data: [][]rune{
				{'A', 'B', 'C', 'D'},
				{'E', 'F', 'G', 'H'},
				{'I', 'J', 'K', 'L'},
			},
		},
		{
			name:       "Empty data",
			screenRows: 3,
			screenCols: 3,
			data:       [][]rune{},
		},
		{
			name:       "Irregular data",
			screenRows: 3,
			screenCols: 3,
			data: [][]rune{
				{'A'},
				{'B', 'C', 'D'},
				{'E', 'F'},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := NewScreen(tt.screenRows, tt.screenCols)
			screen.SetData(tt.data)

			// Verify dimensions unchanged
			assert.Equal(t, tt.screenRows, screen.rows)
			assert.Equal(t, tt.screenCols, screen.cols)

			// Verify data copied correctly
			for i := 0; i < tt.screenRows; i++ {
				for j := 0; j < tt.screenCols; j++ {
					if i < len(tt.data) && j < len(tt.data[i]) {
						assert.Equal(t, tt.data[i][j], screen.data[i][j])
					} else {
						assert.Equal(t, screen.zeroValue, screen.data[i][j])
					}
				}
			}

			// Verify write and view lines
			assert.Equal(t, tt.screenRows-1, screen.writeLine)
			assert.Equal(t, 0, screen.viewLine)
		})
	}
}

// TestAppend tests the Append method
func TestAppend(t *testing.T) {
	tests := []struct {
		name        string
		rows        int
		cols        int
		appendRows  [][]rune
		checkScroll bool
	}{
		{
			name: "Simple append",
			rows: 5,
			cols: 5,
			appendRows: [][]rune{
				{'A', 'B', 'C', 'D', 'E'},
				{'F', 'G', 'H', 'I', 'J'},
			},
			checkScroll: false,
		},
		{
			name: "Append with wrap",
			rows: 3,
			cols: 5,
			appendRows: [][]rune{
				{'1', '2', '3', '4', '5'},
				{'6', '7', '8', '9', '0'},
				{'A', 'B', 'C', 'D', 'E'},
				{'F', 'G', 'H', 'I', 'J'}, // Should wrap
			},
			checkScroll: true,
		},
		{
			name: "Append shorter rows",
			rows: 3,
			cols: 5,
			appendRows: [][]rune{
				{'A', 'B'},
				{'C'},
				{'D', 'E', 'F'},
			},
			checkScroll: false,
		},
		{
			name: "Append longer rows",
			rows: 3,
			cols: 3,
			appendRows: [][]rune{
				{'A', 'B', 'C', 'D', 'E'},
				{'F', 'G', 'H', 'I', 'J'},
			},
			checkScroll: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := NewScreen(tt.rows, tt.cols)

			for _, row := range tt.appendRows {
				screen.Append(row)
			}

			// Verify write line position
			// writeLine increments after each append and doesn't wrap until next append
			expectedWriteLine := len(tt.appendRows)
			// Only wrap if we've exceeded rows count
			if expectedWriteLine > tt.rows {
				expectedWriteLine = expectedWriteLine % tt.rows
			}
			assert.Equal(t, expectedWriteLine, screen.writeLine)

			// Verify view line equals writeLine
			assert.Equal(t, screen.writeLine, screen.viewLine)
		})
	}
}

// TestView tests the View rendering method
func TestView(t *testing.T) {
	tests := []struct {
		name          string
		rows          int
		cols          int
		setupFunc     func(*Screen)
		expectedLines int
		contains      []string
	}{
		{
			name: "Empty screen",
			rows: 3,
			cols: 3,
			setupFunc: func(s *Screen) {
				// Empty screen
			},
			expectedLines: 3,
			contains:      []string{},
		},
		{
			name: "Filled screen",
			rows: 3,
			cols: 3,
			setupFunc: func(s *Screen) {
				s.Fill('X')
			},
			expectedLines: 3,
			contains:      []string{"XXX"},
		},
		{
			name: "Screen with colors",
			rows: 2,
			cols: 5,
			setupFunc: func(s *Screen) {
				s.SetCharColor('A', lipgloss.Color("#FF0000"))
				s.SetCharColor('B', lipgloss.Color("#00FF00"))
				s.Append([]rune{'A', 'B', 'A', 'B', 'A'})
				s.Append([]rune{'B', 'A', 'B', 'A', 'B'})
			},
			expectedLines: 2,
			contains:      []string{"A", "B"},
		},
		{
			name: "Screen with custom zero value",
			rows: 3,
			cols: 3,
			setupFunc: func(s *Screen) {
				s.SetZeroValue('.')
				s.Reset()
			},
			expectedLines: 3,
			contains:      []string{"..."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := NewScreen(tt.rows, tt.cols)
			if tt.setupFunc != nil {
				tt.setupFunc(screen)
			}

			view := screen.View()
			lines := strings.Split(strings.TrimRight(view, "\n"), "\n")

			// Check line count
			assert.GreaterOrEqual(t, len(lines), 1) // At least one line

			// Check contains
			for _, expected := range tt.contains {
				assert.Contains(t, view, expected)
			}
		})
	}
}

// BenchmarkNewScreen benchmarks screen creation
func BenchmarkNewScreen(b *testing.B) {
	sizes := []struct {
		name string
		rows int
		cols int
	}{
		{"Small-10x10", 10, 10},
		{"Medium-50x50", 50, 50},
		{"Large-100x100", 100, 100},
		{"XLarge-500x500", 500, 500},
		{"Terminal-24x80", 24, 80},
		{"HD-1080x1920", 1080, 1920},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = NewScreen(size.rows, size.cols)
			}
		})
	}
}

// BenchmarkSetSize benchmarks resizing operations
func BenchmarkSetSize(b *testing.B) {
	screen := NewScreen(100, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newSize := 50 + (i % 100)
		screen.SetSize(newSize, newSize)
	}
}

// BenchmarkAppend benchmarks append operations
func BenchmarkAppend(b *testing.B) {
	sizes := []struct {
		name string
		rows int
		cols int
	}{
		{"Small-10x50", 10, 50},
		{"Medium-50x100", 50, 100},
		{"Large-100x200", 100, 200},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			screen := NewScreen(size.rows, size.cols)
			row := make([]rune, size.cols)
			for i := range row {
				row[i] = 'A'
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				screen.Append(row)
			}
		})
	}
}

// BenchmarkView benchmarks view rendering
func BenchmarkView(b *testing.B) {
	configurations := []struct {
		name      string
		rows      int
		cols      int
		hasColors bool
	}{
		{"Small-NoColor", 10, 40, false},
		{"Small-WithColor", 10, 40, true},
		{"Medium-NoColor", 30, 80, false},
		{"Medium-WithColor", 30, 80, true},
		{"Large-NoColor", 100, 200, false},
		{"Large-WithColor", 100, 200, true},
	}

	for _, config := range configurations {
		b.Run(config.name, func(b *testing.B) {
			screen := NewScreen(config.rows, config.cols)

			// Fill with data
			for i := 0; i < config.rows; i++ {
				row := make([]rune, config.cols)
				for j := range row {
					row[j] = rune('A' + (i+j)%26)
				}
				screen.Append(row)
			}

			// Add colors if requested
			if config.hasColors {
				for i := 0; i < 26; i++ {
					color := lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", i*10, i*10, i*10))
					screen.SetCharColor(rune('A'+i), color)
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = screen.View()
			}
		})
	}
}

// BenchmarkFill benchmarks fill operations
func BenchmarkFill(b *testing.B) {
	sizes := []struct {
		name string
		rows int
		cols int
	}{
		{"Small-10x10", 10, 10},
		{"Medium-50x50", 50, 50},
		{"Large-100x100", 100, 100},
		{"XLarge-500x500", 500, 500},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			screen := NewScreen(size.rows, size.cols)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				screen.Fill(rune('A' + i%26))
			}
		})
	}
}

// BenchmarkReset benchmarks reset operations
func BenchmarkReset(b *testing.B) {
	sizes := []struct {
		name string
		rows int
		cols int
	}{
		{"Small-10x10", 10, 10},
		{"Medium-50x50", 50, 50},
		{"Large-100x100", 100, 100},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			screen := NewScreen(size.rows, size.cols)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				screen.Reset()
			}
		})
	}
}

// BenchmarkSetData benchmarks SetData operations
func BenchmarkSetData(b *testing.B) {
	sizes := []struct {
		name string
		rows int
		cols int
	}{
		{"Small-10x10", 10, 10},
		{"Medium-30x30", 30, 30},
		{"Large-50x50", 50, 50},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			screen := NewScreen(size.rows, size.cols)

			// Prepare data
			data := make([][]rune, size.rows)
			for i := range data {
				data[i] = make([]rune, size.cols)
				for j := range data[i] {
					data[i][j] = 'X'
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				screen.SetData(data)
			}
		})
	}
}
