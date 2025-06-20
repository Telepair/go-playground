package pkg

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInitLog tests the InitLog function with various configurations
func TestInitLog(t *testing.T) {
	// Create a temporary directory for test log files
	tempDir, err := os.MkdirTemp("", "log_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name      string
		level     string
		format    string
		file      string
		wantLevel slog.Level
		wantErr   bool
	}{
		{
			name:      "Default settings with stdout",
			level:     "",
			format:    "",
			file:      "",
			wantLevel: slog.LevelInfo,
			wantErr:   false,
		},
		{
			name:      "Debug level with JSON format",
			level:     "debug",
			format:    "json",
			file:      "",
			wantLevel: slog.LevelDebug,
			wantErr:   false,
		},
		{
			name:      "Info level with text format",
			level:     "info",
			format:    "text",
			file:      "",
			wantLevel: slog.LevelInfo,
			wantErr:   false,
		},
		{
			name:      "Warn level with file output",
			level:     "warn",
			format:    "json",
			file:      filepath.Join(tempDir, "test.log"),
			wantLevel: slog.LevelWarn,
			wantErr:   false,
		},
		{
			name:      "Error level with uppercase",
			level:     "ERROR",
			format:    "JSON",
			file:      "",
			wantLevel: slog.LevelError,
			wantErr:   false,
		},
		{
			name:      "Invalid level defaults to info",
			level:     "invalid",
			format:    "text",
			file:      "",
			wantLevel: slog.LevelInfo,
			wantErr:   false,
		},
		{
			name:      "Invalid format defaults to text",
			level:     "info",
			format:    "invalid",
			file:      "",
			wantLevel: slog.LevelInfo,
			wantErr:   false,
		},
		{
			name:      "Trimmed level with spaces",
			level:     "  debug  ",
			format:    "text",
			file:      "",
			wantLevel: slog.LevelDebug,
			wantErr:   false,
		},
		{
			name:      "Mixed case format",
			level:     "info",
			format:    "JsOn",
			file:      "",
			wantLevel: slog.LevelInfo,
			wantErr:   false,
		},
		{
			name:      "Invalid file path",
			level:     "info",
			format:    "text",
			file:      "/invalid/path/that/does/not/exist/test.log",
			wantLevel: slog.LevelInfo,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InitLog(tt.level, tt.format, tt.file)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify the log file was created if specified
				if tt.file != "" {
					_, err := os.Stat(tt.file)
					assert.NoError(t, err, "Log file should be created")
				}
			}
		})
	}
}

// TestLogLevelBehavior tests that the log level filtering works correctly
func TestLogLevelBehavior(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "log_level_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name          string
		level         string
		testLogLevel  slog.Level
		shouldContain bool
	}{
		{
			name:          "Debug level allows debug messages",
			level:         "debug",
			testLogLevel:  slog.LevelDebug,
			shouldContain: true,
		},
		{
			name:          "Info level filters out debug messages",
			level:         "info",
			testLogLevel:  slog.LevelDebug,
			shouldContain: false,
		},
		{
			name:          "Warn level filters out info messages",
			level:         "warn",
			testLogLevel:  slog.LevelInfo,
			shouldContain: false,
		},
		{
			name:          "Error level allows error messages",
			level:         "error",
			testLogLevel:  slog.LevelError,
			shouldContain: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logFile := filepath.Join(tempDir, tt.name+".log")
			err := InitLog(tt.level, "text", logFile)
			require.NoError(t, err)

			// Write a test log message
			testMessage := "test log message"
			switch tt.testLogLevel {
			case slog.LevelDebug:
				slog.Debug(testMessage)
			case slog.LevelInfo:
				slog.Info(testMessage)
			case slog.LevelWarn:
				slog.Warn(testMessage)
			case slog.LevelError:
				slog.Error(testMessage)
			}

			// Read the log file and check if the message was written
			content, err := os.ReadFile(logFile)
			require.NoError(t, err)

			if tt.shouldContain {
				assert.Contains(t, string(content), testMessage)
			} else {
				assert.NotContains(t, string(content), testMessage)
			}
		})
	}
}

// TestLogFormatOutput tests that different log formats produce expected output
func TestLogFormatOutput(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "log_format_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name           string
		format         string
		expectedFormat func(string) bool
	}{
		{
			name:   "JSON format",
			format: "json",
			expectedFormat: func(content string) bool {
				return strings.Contains(content, `"msg":"test message"`) &&
					strings.Contains(content, `"level":"INFO"`)
			},
		},
		{
			name:   "Text format",
			format: "text",
			expectedFormat: func(content string) bool {
				return strings.Contains(content, "level=INFO") &&
					strings.Contains(content, `msg="test message"`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logFile := filepath.Join(tempDir, tt.name+".log")
			err := InitLog("info", tt.format, logFile)
			require.NoError(t, err)

			// Write a test log message
			slog.Info("test message", "key", "value")

			// Read the log file
			content, err := os.ReadFile(logFile)
			require.NoError(t, err)

			// Check format
			assert.True(t, tt.expectedFormat(string(content)),
				"Log format should match expected pattern")
		})
	}
}

// BenchmarkInitLog benchmarks the InitLog function
func BenchmarkInitLog(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "log_bench_*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	benchmarks := []struct {
		name   string
		level  string
		format string
		file   string
	}{
		{"Stdout-Text", "info", "text", ""},
		{"Stdout-JSON", "info", "json", ""},
		{"File-Text", "info", "text", filepath.Join(tempDir, "bench_text.log")},
		{"File-JSON", "info", "json", filepath.Join(tempDir, "bench_json.log")},
		{"Debug-Text", "debug", "text", ""},
		{"Error-JSON", "error", "json", ""},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err := InitLog(bm.level, bm.format, bm.file)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkLogging benchmarks actual logging performance with different configurations
func BenchmarkLogging(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "log_perf_bench_*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	configurations := []struct {
		name   string
		level  string
		format string
		file   string
	}{
		{"TextStdout", "info", "text", ""},
		{"JSONStdout", "info", "json", ""},
		{"TextFile", "info", "text", filepath.Join(tempDir, "perf_text.log")},
		{"JSONFile", "info", "json", filepath.Join(tempDir, "perf_json.log")},
	}

	for _, config := range configurations {
		b.Run(config.name, func(b *testing.B) {
			err := InitLog(config.level, config.format, config.file)
			require.NoError(b, err)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				slog.Info("benchmark log message",
					"iteration", i,
					"config", config.name,
					"extra_field", "some value")
			}
		})
	}
}

// BenchmarkLogLevelFiltering benchmarks the performance impact of log level filtering
func BenchmarkLogLevelFiltering(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "log_filter_bench_*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	levels := []string{"debug", "info", "warn", "error"}

	for _, level := range levels {
		b.Run("Level-"+level, func(b *testing.B) {
			logFile := filepath.Join(tempDir, "filter_"+level+".log")
			err := InitLog(level, "text", logFile)
			require.NoError(b, err)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Log at all levels - some will be filtered
				slog.Debug("debug message", "i", i)
				slog.Info("info message", "i", i)
				slog.Warn("warn message", "i", i)
				slog.Error("error message", "i", i)
			}
		})
	}
}
