package ui

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	// DefaultLang is the default language setting
	DefaultLang = "en"
	// DefaultRefreshInterval is the default refresh interval for the UI
	DefaultRefreshInterval = 200 * time.Millisecond
	// MinRefreshInterval is the minimum allowed refresh interval
	MinRefreshInterval = 10 * time.Millisecond
	// DefaultAliveColor is the default color for alive/active cells
	DefaultAliveColor = "#00FF00"
	// DefaultDeadColor is the default color for dead/inactive cells
	DefaultDeadColor = "#000000"
	// DefaultAliveChar is the default character for alive/active cells
	DefaultAliveChar rune = '█'
	// DefaultDeadChar is the default character for dead/inactive cells
	DefaultDeadChar rune = ' '

	// DefaultWidth is the default terminal width
	DefaultWidth = 80
	// DefaultHeight is the default terminal height
	DefaultHeight = 24
)

var (
	keepHeight      = 6 // header (2), status, control line, screen (2)
	headerLineStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#16213E")).
			MarginBottom(1).
			Align(lipgloss.Center)
	statusLineStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(lipgloss.Color("#94A3B8")).
			Background(lipgloss.Color("#0F3460")).
			Bold(true)
	controlLineStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Foreground(lipgloss.Color("#94A3B8")).
				Background(lipgloss.Color("#0F3460")).
				Bold(true)
	statusKVSplit    = ": "
	statusItemSplit  = " | "
	controlKVSplit   = ": "
	controlItemSplit = " | "
)

// Model represents the application state
type Model struct {
	engine StepEngine

	language    Language
	refreshRate time.Duration

	currentStep int
	paused      bool
	height      int
	width       int

	buffer        strings.Builder
	statusBuffer  strings.Builder
	controlBuffer strings.Builder
	controlKeys   map[string]struct{}
	logger        *slog.Logger
}

// RunModel creates a new model with the given configuration
func RunModel(appName string, engine StepEngine, defaultLang string, defaultRefreshInterval time.Duration) error {
	if appName == "" {
		return fmt.Errorf("appName cannot be empty")
	}
	if engine == nil {
		return fmt.Errorf("engine cannot be nil")
	}

	logger := slog.With("app", appName)

	model := &Model{
		engine:      engine,
		language:    ToLanguage(defaultLang),
		refreshRate: defaultRefreshInterval,
		currentStep: 0,
		paused:      false,
		width:       DefaultWidth,
		height:      DefaultHeight,
		buffer:      strings.Builder{},
		controlKeys: make(map[string]struct{}),
		logger:      logger,
	}

	// Run the TUI application
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Error("Error running program", "error", err)
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	logger.Debug("Finished")
	return nil
}

// tickMsg is sent every tick for infinite mode
type tickMsg time.Time

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	// Always start the timer when initializing unless already quitting
	return tea.Tick(m.refreshRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.logger.Debug("Window size changed", "width", msg.Width, "height", msg.Height)
		return m.handleWindowResize(msg)
	case tea.KeyMsg:
		m.logger.Debug("Key pressed", "key", msg.String())
		return m.handleKeyPress(msg)
	case tickMsg:
		m.logger.Debug("Tick", "time", msg)
		return m.handleTick()
	}
	return m, nil
}

// View renders the current state
func (m *Model) View() string {
	m.logger.Debug("Model View",
		"width", m.width,
		"height", m.height,
		"language", m.language,
		"paused", m.paused,
		"currentStep", m.currentStep,
		"refreshRate", m.refreshRate,
		"status", m.engine.Status(English),
		"isFinished", m.engine.IsFinished())
	return m.RenderMode()
}

// handleWindowResize processes terminal window size changes
func (m *Model) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.height = msg.Height
	m.width = msg.Width

	// Ensure minimum dimensions
	engineHeight := max(m.height-keepHeight, 1)
	engineWidth := max(m.width, 1)

	if err := m.engine.Reset(engineHeight, engineWidth); err != nil {
		m.logger.Error("Failed to reset engine", "height", engineHeight, "width", engineWidth, "error", err)
	} else {
		m.logger.Debug("Reset engine", "height", engineHeight, "width", engineWidth)
	}
	return m, nil
}

// handleKeyPress processes keyboard input
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := strings.ToLower(msg.String())
	m.logger.Debug("Key pressed", "key", key)

	if _, ok := m.controlKeys[key]; ok {
		if handled, err := m.engine.Handle(key); err != nil {
			m.logger.Error("Key not handled", "key", key, "error", err)
			return m, tea.Quit
		} else if handled {
			m.logger.Debug("Key handled", "key", key)
			return m, nil
		}
	}

	switch key {
	case "ctrl+c", "q", "esc":
		m.engine.Stop()
		return m, tea.Quit

	case " ", "enter": // Space or Enter key for pause/resume
		m.paused = !m.paused

	case "l": // Language toggle key
		if m.language == English {
			m.language = Chinese
		} else {
			m.language = English
		}

	case "+", "=", "up": // Increase refresh rate (make it faster)
		m.refreshRate = max(m.refreshRate/2, MinRefreshInterval)

	case "-", "_", "down": // Decrease refresh rate (make it slower)
		m.refreshRate = m.refreshRate * 2

	case "r": // Reset simulation
		m.currentStep = 0
		engineHeight := max(m.height-keepHeight, 1)
		engineWidth := max(m.width, 1)
		if err := m.engine.Reset(engineHeight, engineWidth); err != nil {
			m.logger.Error("Failed to reset engine", "height", engineHeight, "width", engineWidth, "error", err)
		} else {
			m.logger.Debug("Reset engine", "height", engineHeight, "width", engineWidth)
		}
	}

	return m, nil
}

// handleTick processes timer ticks
func (m *Model) handleTick() (tea.Model, tea.Cmd) {
	// Check if we should continue running (only update when not paused)
	if !m.paused {
		currentStep, ok := m.engine.Step()
		if !ok {
			m.paused = true
		}
		m.currentStep = currentStep
		// Only pause if the simulation is finished
		if m.engine.IsFinished() {
			m.paused = true
		}
	}

	// Continue ticking only if not quitting
	return m, tea.Tick(m.refreshRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// RenderMode renders the complete UI mode view with enhanced layout
func (m *Model) RenderMode() string {
	m.buffer.Reset()
	m.buffer.WriteString(m.RenderHeader())
	m.buffer.WriteString("\n")
	m.buffer.WriteString(m.RenderStatus())
	m.buffer.WriteString("\n")
	m.buffer.WriteString(m.engine.View())
	m.buffer.WriteString("\n")
	m.buffer.WriteString(m.RenderControlLine())
	return m.buffer.String()
}

// RenderHeader renders the header
func (m *Model) RenderHeader() string {
	return headerLineStyle.Width(m.width).Render(m.engine.Header(m.language))
}

// RenderStatus renders the status line
func (m *Model) RenderStatus() string {
	m.statusBuffer.Reset()

	// Add engine status items
	engineStatus := m.engine.Status(m.language)
	modelStatus := m.Status(m.language)
	allStatus := append(engineStatus, modelStatus...)
	if len(allStatus) == 0 {
		return ""
	}

	m.statusBuffer.WriteString(allStatus[0].Label)
	m.statusBuffer.WriteString(statusKVSplit)
	m.statusBuffer.WriteString(allStatus[0].Value)
	for _, item := range allStatus[1:] {
		m.statusBuffer.WriteString(statusItemSplit)
		m.statusBuffer.WriteString(item.Label)
		m.statusBuffer.WriteString(statusKVSplit)
		m.statusBuffer.WriteString(item.Value)
	}

	return statusLineStyle.Width(m.width).Render(m.statusBuffer.String())
}

// Status returns the status line items
func (m *Model) Status(lang Language) []Status {
	var statusText string
	if m.paused {
		if lang == Chinese {
			statusText = "已暂停"
		} else {
			statusText = "Paused"
		}
	} else {
		if lang == Chinese {
			statusText = "运行中"
		} else {
			statusText = "Running"
		}
	}

	engineHeight := max(m.height, 1)
	engineWidth := max(m.width, 1)

	if lang == Chinese {
		return []Status{
			{Label: "代数", Value: strconv.Itoa(m.currentStep)},
			{Label: "刷新", Value: m.refreshRate.String()},
			{Label: "尺寸", Value: fmt.Sprintf("%d×%d", engineWidth, engineHeight)},
			{Label: "状态", Value: statusText},
		}
	}
	return []Status{
		{Label: "Gen", Value: strconv.Itoa(m.currentStep)},
		{Label: "Speed", Value: m.refreshRate.String()},
		{Label: "Size", Value: fmt.Sprintf("%d×%d", engineWidth, engineHeight)},
		{Label: "Status", Value: statusText},
	}
}

// RenderControlLine renders the control line
func (m *Model) RenderControlLine() string {
	m.controlBuffer.Reset()

	first := true

	// Get engine-specific controls
	for _, item := range m.engine.HandleKeys(m.language) {
		for _, key := range item.Keys {
			m.controlKeys[strings.ToLower(key)] = struct{}{}
		}
		if !first {
			m.controlBuffer.WriteString(controlItemSplit)
		}
		m.controlBuffer.WriteString(strings.Join(item.Keys, "/"))
		m.controlBuffer.WriteString(controlKVSplit)
		m.controlBuffer.WriteString(item.Label)
		first = false
	}
	// Get common controls
	commonControls := m.Controls(m.language)
	for _, item := range commonControls {
		if !first {
			m.controlBuffer.WriteString(controlItemSplit)
		}
		m.controlBuffer.WriteString(strings.Join(item.Keys, "/"))
		m.controlBuffer.WriteString(controlKVSplit)
		m.controlBuffer.WriteString(item.Label)
		first = false
	}

	return controlLineStyle.Width(m.width).Render(m.controlBuffer.String())
}

// Controls returns the control line items
func (m *Model) Controls(lang Language) []Control {
	if lang == Chinese {
		return []Control{
			{Keys: []string{"L"}, Label: "语言"},
			{Keys: []string{"+", "up"}, Label: "加速"},
			{Keys: []string{"-", "down"}, Label: "减速"},
			{Keys: []string{"Space"}, Label: "暂停/继续"},
			{Keys: []string{"R"}, Label: "重置"},
			{Keys: []string{"Q"}, Label: "退出"},
		}
	}
	return []Control{
		{Keys: []string{"L"}, Label: "Language"},
		{Keys: []string{"+", "Up"}, Label: "Speed +"},
		{Keys: []string{"-", "Down"}, Label: "Speed -"},
		{Keys: []string{"Space"}, Label: "Pause/Continue"},
		{Keys: []string{"R"}, Label: "Reset"},
		{Keys: []string{"Q"}, Label: "Quit"},
	}
}
