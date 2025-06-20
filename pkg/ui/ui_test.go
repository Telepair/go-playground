package ui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStepEngine is a mock implementation of StepEngine for testing
type MockStepEngine struct {
	mock.Mock
}

func (m *MockStepEngine) Step() (int, bool) {
	args := m.Called()
	return args.Int(0), args.Bool(1)
}

func (m *MockStepEngine) Header(lang Language) string {
	args := m.Called(lang)
	return args.String(0)
}

func (m *MockStepEngine) Status(lang Language) []Status {
	args := m.Called(lang)
	return args.Get(0).([]Status)
}

func (m *MockStepEngine) HandleKeys(lang Language) []Control {
	args := m.Called(lang)
	return args.Get(0).([]Control)
}

func (m *MockStepEngine) Handle(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockStepEngine) Reset(rows, cols int) error {
	args := m.Called(rows, cols)
	return args.Error(0)
}

func (m *MockStepEngine) IsFinished() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockStepEngine) Stop() {
	m.Called()
}

func (m *MockStepEngine) View() string {
	args := m.Called()
	return args.String(0)
}

// TestLanguageConversion tests the language conversion functions
func TestLanguageConversion(t *testing.T) {
	tests := []struct {
		input    string
		expected Language
	}{
		{"en", English},
		{"EN", English},
		{"english", English},
		{"English", English},
		{"zh", Chinese},
		{"ZH", Chinese},
		{"chinese", Chinese},
		{"Chinese", Chinese},
		{"", English},        // Default
		{"unknown", English}, // Unknown defaults to English
		{"fr", English},      // Unsupported language
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToLanguage(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestModelInitialization tests the initialization of the Model
func TestModelInitialization(t *testing.T) {
	mockEngine := new(MockStepEngine)
	mockEngine.On("View").Return("test view")
	mockEngine.On("Header", English).Return("Test Header")
	mockEngine.On("Status", English).Return([]Status{{Label: "Test", Value: "Value"}})
	mockEngine.On("HandleKeys", English).Return([]Control{{Keys: []string{"T"}, Label: "Test"}})
	mockEngine.On("IsFinished").Return(false)

	model := &Model{
		engine:      mockEngine,
		language:    English,
		refreshRate: DefaultRefreshInterval,
		width:       DefaultWidth,
		height:      DefaultHeight,
		buffer:      strings.Builder{},
		controlKeys: make(map[string]struct{}),
		logger:      slog.Default(),
	}

	// Test Init
	cmd := model.Init()
	assert.NotNil(t, cmd, "Init should return a tick command")

	// Test View
	view := model.View()
	assert.NotEmpty(t, view, "View should not be empty")
	assert.Contains(t, view, "Test Header", "View should contain the header")
}

// TestHandleWindowResize tests window resize handling
func TestHandleWindowResize(t *testing.T) {
	mockEngine := new(MockStepEngine)
	mockEngine.On("Reset", mock.Anything, mock.Anything).Return(nil)
	mockEngine.On("View").Return("test view")
	mockEngine.On("Header", English).Return("Test Header")
	mockEngine.On("Status", English).Return([]Status{})
	mockEngine.On("HandleKeys", English).Return([]Control{})
	mockEngine.On("IsFinished").Return(false)

	model := &Model{
		engine:      mockEngine,
		language:    English,
		refreshRate: DefaultRefreshInterval,
		width:       DefaultWidth,
		height:      DefaultHeight,
		buffer:      strings.Builder{},
		controlKeys: make(map[string]struct{}),
		logger:      slog.Default(),
	}

	// Test window resize
	resizeMsg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updatedModel, cmd := model.handleWindowResize(resizeMsg)

	assert.Equal(t, 100, updatedModel.(*Model).width)
	assert.Equal(t, 30, updatedModel.(*Model).height)
	assert.Nil(t, cmd)

	// Verify Reset was called with correct dimensions
	expectedHeight := 30 - keepHeight
	expectedWidth := 100
	mockEngine.AssertCalled(t, "Reset", expectedHeight, expectedWidth)
}

// TestHandleKeyPress tests keyboard input handling
func TestHandleKeyPress(t *testing.T) {
	mockEngine := new(MockStepEngine)
	mockEngine.On("Stop").Return()
	mockEngine.On("Handle", mock.Anything).Return(false, nil)
	mockEngine.On("Reset", mock.Anything, mock.Anything).Return(nil)
	mockEngine.On("View").Return("test view")
	mockEngine.On("Header", mock.Anything).Return("Test Header")
	mockEngine.On("Status", mock.Anything).Return([]Status{})
	mockEngine.On("HandleKeys", mock.Anything).Return([]Control{})
	mockEngine.On("IsFinished").Return(false)

	tests := []struct {
		name           string
		key            string
		initialPaused  bool
		initialLang    Language
		initialRate    time.Duration
		expectedPaused bool
		expectedLang   Language
		expectedRate   time.Duration
		shouldQuit     bool
	}{
		{
			name:           "Pause with space",
			key:            " ",
			initialPaused:  false,
			initialLang:    English,
			initialRate:    DefaultRefreshInterval,
			expectedPaused: true,
			expectedLang:   English,
			expectedRate:   DefaultRefreshInterval,
			shouldQuit:     false,
		},
		{
			name:           "Resume with enter",
			key:            "enter",
			initialPaused:  true,
			initialLang:    English,
			initialRate:    DefaultRefreshInterval,
			expectedPaused: false,
			expectedLang:   English,
			expectedRate:   DefaultRefreshInterval,
			shouldQuit:     false,
		},
		{
			name:           "Toggle language",
			key:            "l",
			initialPaused:  false,
			initialLang:    English,
			initialRate:    DefaultRefreshInterval,
			expectedPaused: false,
			expectedLang:   Chinese,
			expectedRate:   DefaultRefreshInterval,
			shouldQuit:     false,
		},
		{
			name:           "Increase speed",
			key:            "+",
			initialPaused:  false,
			initialLang:    English,
			initialRate:    DefaultRefreshInterval,
			expectedPaused: false,
			expectedLang:   English,
			expectedRate:   DefaultRefreshInterval / 2,
			shouldQuit:     false,
		},
		{
			name:           "Decrease speed",
			key:            "-",
			initialPaused:  false,
			initialLang:    English,
			initialRate:    DefaultRefreshInterval,
			expectedPaused: false,
			expectedLang:   English,
			expectedRate:   DefaultRefreshInterval * 2,
			shouldQuit:     false,
		},
		{
			name:           "Quit with q",
			key:            "q",
			initialPaused:  false,
			initialLang:    English,
			initialRate:    DefaultRefreshInterval,
			expectedPaused: false,
			expectedLang:   English,
			expectedRate:   DefaultRefreshInterval,
			shouldQuit:     true,
		},
		{
			name:           "Reset with r",
			key:            "r",
			initialPaused:  false,
			initialLang:    English,
			initialRate:    DefaultRefreshInterval,
			expectedPaused: false,
			expectedLang:   English,
			expectedRate:   DefaultRefreshInterval,
			shouldQuit:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &Model{
				engine:      mockEngine,
				language:    tt.initialLang,
				refreshRate: tt.initialRate,
				paused:      tt.initialPaused,
				width:       DefaultWidth,
				height:      DefaultHeight,
				buffer:      strings.Builder{},
				controlKeys: make(map[string]struct{}),
				logger:      slog.Default(),
			}

			_, cmd := model.handleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)})

			if tt.shouldQuit {
				assert.NotNil(t, cmd, "Should return quit command")
			} else {
				assert.Equal(t, tt.expectedPaused, model.paused)
				assert.Equal(t, tt.expectedLang, model.language)
				assert.Equal(t, tt.expectedRate, model.refreshRate)
			}
		})
	}
}

// TestHandleTick tests the tick handling
func TestHandleTick(t *testing.T) {
	tests := []struct {
		name              string
		initialPaused     bool
		stepReturns       []interface{}
		isFinishedReturns bool
		expectedPaused    bool
		expectedStep      int
	}{
		{
			name:              "Normal step when not paused",
			initialPaused:     false,
			stepReturns:       []interface{}{5, true},
			isFinishedReturns: false,
			expectedPaused:    false,
			expectedStep:      5,
		},
		{
			name:              "Pause when step returns false",
			initialPaused:     false,
			stepReturns:       []interface{}{10, false},
			isFinishedReturns: false,
			expectedPaused:    true,
			expectedStep:      10,
		},
		{
			name:              "Pause when finished",
			initialPaused:     false,
			stepReturns:       []interface{}{15, true},
			isFinishedReturns: true,
			expectedPaused:    true,
			expectedStep:      15,
		},
		{
			name:              "No step when paused",
			initialPaused:     true,
			stepReturns:       nil, // Step shouldn't be called
			isFinishedReturns: false,
			expectedPaused:    true,
			expectedStep:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEngine := new(MockStepEngine)
			if tt.stepReturns != nil {
				mockEngine.On("Step").Return(tt.stepReturns...).Once()
			}
			mockEngine.On("IsFinished").Return(tt.isFinishedReturns)

			model := &Model{
				engine:      mockEngine,
				paused:      tt.initialPaused,
				currentStep: 0,
				logger:      slog.Default(),
			}

			_, cmd := model.handleTick()

			assert.NotNil(t, cmd, "Should return a tick command")
			assert.Equal(t, tt.expectedPaused, model.paused)
			assert.Equal(t, tt.expectedStep, model.currentStep)

			if tt.stepReturns != nil {
				mockEngine.AssertCalled(t, "Step")
			} else {
				mockEngine.AssertNotCalled(t, "Step")
			}
		})
	}
}

// TestRenderMethods tests the various render methods
func TestRenderMethods(t *testing.T) {
	mockEngine := new(MockStepEngine)
	mockEngine.On("View").Return("Engine View Content")
	mockEngine.On("Header", English).Return("Test Application")
	mockEngine.On("Header", Chinese).Return("测试应用")
	mockEngine.On("Status", English).Return([]Status{
		{Label: "Rule", Value: "30"},
		{Label: "Cells", Value: "100"},
	})
	mockEngine.On("Status", Chinese).Return([]Status{
		{Label: "规则", Value: "30"},
		{Label: "细胞", Value: "100"},
	})
	mockEngine.On("HandleKeys", English).Return([]Control{
		{Keys: []string{"T"}, Label: "Toggle"},
	})
	mockEngine.On("HandleKeys", Chinese).Return([]Control{
		{Keys: []string{"T"}, Label: "切换"},
	})
	mockEngine.On("IsFinished").Return(false)

	t.Run("English rendering", func(t *testing.T) {
		model := &Model{
			engine:        mockEngine,
			language:      English,
			refreshRate:   100 * time.Millisecond,
			currentStep:   5,
			paused:        false,
			width:         80,
			height:        24,
			buffer:        strings.Builder{},
			statusBuffer:  strings.Builder{},
			controlBuffer: strings.Builder{},
			controlKeys:   make(map[string]struct{}),
			logger:        slog.Default(),
		}

		// Test RenderHeader
		header := model.RenderHeader()
		assert.Contains(t, header, "Test Application")

		// Test RenderStatus
		status := model.RenderStatus()
		assert.Contains(t, status, "Rule")
		assert.Contains(t, status, "30")
		assert.Contains(t, status, "Gen")
		assert.Contains(t, status, "5")
		assert.Contains(t, status, "Running")

		// Test RenderControlLine
		controls := model.RenderControlLine()
		assert.Contains(t, controls, "T")
		assert.Contains(t, controls, "Toggle")
		assert.Contains(t, controls, "Language")
		assert.Contains(t, controls, "Quit")

		// Verify control keys were registered
		assert.Contains(t, model.controlKeys, "t")
	})

	t.Run("Chinese rendering", func(t *testing.T) {
		model := &Model{
			engine:        mockEngine,
			language:      Chinese,
			refreshRate:   200 * time.Millisecond,
			currentStep:   10,
			paused:        true,
			width:         100,
			height:        30,
			buffer:        strings.Builder{},
			statusBuffer:  strings.Builder{},
			controlBuffer: strings.Builder{},
			controlKeys:   make(map[string]struct{}),
			logger:        slog.Default(),
		}

		// Test RenderHeader
		header := model.RenderHeader()
		assert.Contains(t, header, "测试应用")

		// Test RenderStatus
		status := model.RenderStatus()
		assert.Contains(t, status, "规则")
		assert.Contains(t, status, "30")
		assert.Contains(t, status, "代数")
		assert.Contains(t, status, "10")
		assert.Contains(t, status, "已暂停")

		// Test RenderControlLine
		controls := model.RenderControlLine()
		assert.Contains(t, controls, "T")
		assert.Contains(t, controls, "切换")
		assert.Contains(t, controls, "语言")
		assert.Contains(t, controls, "退出")
	})
}

// TestStatusAndControls tests the Status and Controls methods
func TestStatusAndControls(t *testing.T) {
	model := &Model{
		refreshRate: 100 * time.Millisecond,
		currentStep: 42,
		paused:      false,
		width:       120,
		height:      40,
		logger:      slog.Default(),
	}

	t.Run("English Status", func(t *testing.T) {
		status := model.Status(English)
		assert.Len(t, status, 4)
		assert.Equal(t, "Gen", status[0].Label)
		assert.Equal(t, "42", status[0].Value)
		assert.Equal(t, "Speed", status[1].Label)
		assert.Equal(t, "100ms", status[1].Value)
		assert.Equal(t, "Size", status[2].Label)
		assert.Contains(t, status[2].Value, "120×40")
		assert.Equal(t, "Status", status[3].Label)
		assert.Equal(t, "Running", status[3].Value)
	})

	t.Run("Chinese Status", func(t *testing.T) {
		model.paused = true
		status := model.Status(Chinese)
		assert.Len(t, status, 4)
		assert.Equal(t, "代数", status[0].Label)
		assert.Equal(t, "42", status[0].Value)
		assert.Equal(t, "刷新", status[1].Label)
		assert.Equal(t, "尺寸", status[2].Label)
		assert.Equal(t, "状态", status[3].Label)
		assert.Equal(t, "已暂停", status[3].Value)
	})

	t.Run("English Controls", func(t *testing.T) {
		controls := model.Controls(English)
		assert.Len(t, controls, 6)

		// Check for specific controls
		var hasLanguage, hasSpeed, hasPause bool
		for _, control := range controls {
			switch control.Label {
			case "Language":
				hasLanguage = true
				assert.Contains(t, control.Keys, "L")
			case "Speed +":
				hasSpeed = true
			case "Pause/Continue":
				hasPause = true
				assert.Contains(t, control.Keys, "Space")
			}
		}
		assert.True(t, hasLanguage)
		assert.True(t, hasSpeed)
		assert.True(t, hasPause)
	})

	t.Run("Chinese Controls", func(t *testing.T) {
		controls := model.Controls(Chinese)
		assert.Len(t, controls, 6)

		// Check for specific controls
		var hasLanguage, hasSpeed, hasPause bool
		for _, control := range controls {
			switch control.Label {
			case "语言":
				hasLanguage = true
			case "加速":
				hasSpeed = true
			case "暂停/继续":
				hasPause = true
			}
		}
		assert.True(t, hasLanguage)
		assert.True(t, hasSpeed)
		assert.True(t, hasPause)
	})
}

// TestRunModelValidation tests RunModel parameter validation
func TestRunModelValidation(t *testing.T) {
	mockEngine := new(MockStepEngine)

	tests := []struct {
		name          string
		appName       string
		engine        StepEngine
		lang          string
		refreshRate   time.Duration
		expectedError string
	}{
		{
			name:          "Empty app name",
			appName:       "",
			engine:        mockEngine,
			lang:          "en",
			refreshRate:   DefaultRefreshInterval,
			expectedError: "appName cannot be empty",
		},
		{
			name:          "Nil engine",
			appName:       "TestApp",
			engine:        nil,
			lang:          "en",
			refreshRate:   DefaultRefreshInterval,
			expectedError: "engine cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Only test parameter validation, skip actual TUI execution
			if tt.expectedError != "" {
				err := RunModel(tt.appName, tt.engine, tt.lang, tt.refreshRate)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

// TestEngineHandleIntegration tests the integration with engine's Handle method
func TestEngineHandleIntegration(t *testing.T) {
	mockEngine := new(MockStepEngine)
	mockEngine.On("HandleKeys", English).Return([]Control{
		{Keys: []string{"A", "a"}, Label: "Action A"},
		{Keys: []string{"B"}, Label: "Action B"},
		{Keys: []string{"D"}, Label: "Action D"},
	})
	mockEngine.On("Handle", "a").Return(true, nil).Once()
	mockEngine.On("Handle", "b").Return(true, nil).Once()
	mockEngine.On("Handle", "c").Return(false, nil).Once()
	mockEngine.On("Handle", "d").Return(false, fmt.Errorf("error handling key")).Once()
	mockEngine.On("Stop").Return()
	mockEngine.On("View").Return("test view")
	mockEngine.On("Header", English).Return("Test Header")
	mockEngine.On("Status", English).Return([]Status{})
	mockEngine.On("IsFinished").Return(false)

	model := &Model{
		engine:        mockEngine,
		language:      English,
		width:         80,
		height:        24,
		buffer:        strings.Builder{},
		statusBuffer:  strings.Builder{},
		controlBuffer: strings.Builder{},
		controlKeys:   make(map[string]struct{}),
		logger:        slog.Default(),
	}

	// Initialize control keys
	model.RenderControlLine()

	tests := []struct {
		name       string
		key        string
		shouldQuit bool
	}{
		{"Handle registered key a", "a", false},
		{"Handle registered key b", "b", false},
		{"Handle unregistered key c", "c", false},
		{"Handle error key d", "d", true}, // d is registered and returns error, so quit
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, cmd := model.handleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)})
			if tt.shouldQuit {
				assert.NotNil(t, cmd)
			}
		})
	}
}

// Benchmark tests

// BenchmarkRenderMode benchmarks the complete rendering process
func BenchmarkRenderMode(b *testing.B) {
	mockEngine := new(MockStepEngine)
	mockEngine.On("View").Return(strings.Repeat("X", 1000))
	mockEngine.On("Header", mock.Anything).Return("Benchmark Application Header")
	mockEngine.On("Status", mock.Anything).Return([]Status{
		{Label: "Status1", Value: "Value1"},
		{Label: "Status2", Value: "Value2"},
		{Label: "Status3", Value: "Value3"},
	})
	mockEngine.On("HandleKeys", mock.Anything).Return([]Control{
		{Keys: []string{"A", "a"}, Label: "Action A"},
		{Keys: []string{"B", "b"}, Label: "Action B"},
		{Keys: []string{"C", "c"}, Label: "Action C"},
	})
	mockEngine.On("IsFinished").Return(false)

	model := &Model{
		engine:        mockEngine,
		language:      English,
		refreshRate:   100 * time.Millisecond,
		currentStep:   1000,
		paused:        false,
		width:         120,
		height:        40,
		buffer:        strings.Builder{},
		statusBuffer:  strings.Builder{},
		controlBuffer: strings.Builder{},
		controlKeys:   make(map[string]struct{}),
		logger:        slog.Default(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.RenderMode()
	}
}

// BenchmarkHandleKeyPress benchmarks key press handling
func BenchmarkHandleKeyPress(b *testing.B) {
	mockEngine := new(MockStepEngine)
	mockEngine.On("Handle", mock.Anything).Return(false, nil)
	mockEngine.On("Stop").Return()
	mockEngine.On("Reset", mock.Anything, mock.Anything).Return(nil)

	model := &Model{
		engine:      mockEngine,
		language:    English,
		refreshRate: DefaultRefreshInterval,
		paused:      false,
		width:       80,
		height:      24,
		controlKeys: make(map[string]struct{}),
		logger:      slog.Default(),
	}

	keys := []string{" ", "l", "+", "-", "r", "x", "t", "q"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			model.handleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
		}
	}
}

// BenchmarkRenderStatus benchmarks status line rendering
func BenchmarkRenderStatus(b *testing.B) {
	mockEngine := new(MockStepEngine)
	mockEngine.On("Status", mock.Anything).Return([]Status{
		{Label: "Status1", Value: "Value1"},
		{Label: "Status2", Value: "Value2"},
		{Label: "Status3", Value: "Value3"},
		{Label: "Status4", Value: "Value4"},
		{Label: "Status5", Value: "Value5"},
	})

	model := &Model{
		engine:       mockEngine,
		language:     English,
		refreshRate:  100 * time.Millisecond,
		currentStep:  999,
		paused:       false,
		width:        150,
		height:       50,
		statusBuffer: strings.Builder{},
		logger:       slog.Default(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.RenderStatus()
	}
}

// BenchmarkRenderControlLine benchmarks control line rendering
func BenchmarkRenderControlLine(b *testing.B) {
	mockEngine := new(MockStepEngine)
	mockEngine.On("HandleKeys", mock.Anything).Return([]Control{
		{Keys: []string{"A", "a"}, Label: "Action A"},
		{Keys: []string{"B", "b"}, Label: "Action B"},
		{Keys: []string{"C", "c"}, Label: "Action C"},
		{Keys: []string{"D", "d"}, Label: "Action D"},
		{Keys: []string{"E", "e"}, Label: "Action E"},
	})

	model := &Model{
		engine:        mockEngine,
		language:      English,
		width:         150,
		controlBuffer: strings.Builder{},
		controlKeys:   make(map[string]struct{}),
		logger:        slog.Default(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Clear control keys for each iteration
		model.controlKeys = make(map[string]struct{})
		_ = model.RenderControlLine()
	}
}

// BenchmarkHandleTick benchmarks tick handling
func BenchmarkHandleTick(b *testing.B) {
	mockEngine := new(MockStepEngine)
	mockEngine.On("Step").Return(1, true)
	mockEngine.On("IsFinished").Return(false)

	model := &Model{
		engine:      mockEngine,
		paused:      false,
		currentStep: 0,
		logger:      slog.Default(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.handleTick()
	}
}

// BenchmarkLanguageToggle benchmarks language switching
func BenchmarkLanguageToggle(b *testing.B) {
	model := &Model{
		language: English,
		logger:   slog.Default(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if model.language == English {
			model.language = Chinese
		} else {
			model.language = English
		}
	}
}

// BenchmarkRefreshRateAdjustment benchmarks refresh rate calculations
func BenchmarkRefreshRateAdjustment(b *testing.B) {
	model := &Model{
		refreshRate: DefaultRefreshInterval,
		logger:      slog.Default(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Increase rate
		model.refreshRate = max(model.refreshRate/2, MinRefreshInterval)
		// Decrease rate
		model.refreshRate = model.refreshRate * 2
	}
}
