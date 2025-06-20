// Package ui provides the user interface components and engine abstractions for terminal-based visualizations.
package ui

// StepEngine represents a step engine that can be rendered by the UI
type StepEngine interface {
	// View returns the view of the engine
	View() string
	// Step executes one iteration, returns the current iteration number and whether it was successful
	Step() (int, bool)
	// Header returns the header text for the UI in the specified language
	Header(lang Language) string
	// Status returns the internal status information in the specified language
	Status(lang Language) []Status
	// HandleKeys returns all available keyboard operations in the specified language
	HandleKeys(lang Language) []Control
	// Handle handles keyboard input operations, returns true if the key was handled
	Handle(key string) (bool, error)
	// Reset resets the model
	Reset(height, width int) error
	// IsFinished returns whether the model has finished execution
	IsFinished() bool
	// Stop stops the model
	Stop()
}

// Language represents the supported languages
type Language int

// Language constants
const (
	English Language = iota // English is the default language
	Chinese                 // Chinese is the Chinese language
)

// String returns the string representation of the language
func (l Language) String() string {
	switch l {
	case English:
		return "en"
	case Chinese:
		return "cn"
	}
	return "en"
}

// ToLanguage converts a string to a language
func ToLanguage(lang string) Language {
	switch lang {
	case "en", "EN", "english", "English":
		return English
	case "cn", "CN", "zh", "ZH", "chinese", "Chinese":
		return Chinese
	}
	return English
}

// Status represents a status item with icon, key and value
type Status struct {
	Label string // Label is the label of the status
	Value string // Value is the value of the status
}

// Control represents a control item with keys and label
type Control struct {
	Keys  []string // Key is the key of the control
	Label string   // Label is the label of the control
}
