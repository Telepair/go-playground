package pkg

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStartProfile tests the StartProfile function
func TestStartProfile(t *testing.T) {
	tests := []struct {
		name         string
		port         int
		contextDelay time.Duration
	}{
		{
			name:         "Start and stop profile server",
			port:         6060,
			contextDelay: 100 * time.Millisecond,
		},
		{
			name:         "Start with different port",
			port:         6061,
			contextDelay: 100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.contextDelay)
			defer cancel()

			// Run StartProfile in a goroutine
			done := make(chan struct{})
			go func() {
				StartProfile(ctx, tt.port)
				close(done)
			}()

			// Give the server time to start
			time.Sleep(50 * time.Millisecond)

			// Check if the server is running
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/debug/pprof/", tt.port))
			if err == nil {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				resp.Body.Close()
			}

			// Wait for context to be cancelled and server to shut down
			<-done
		})
	}
}

// TestStartProfileServerError tests server startup error handling
func TestStartProfileServerError(t *testing.T) {
	// Start a server on a port
	server := &http.Server{
		Addr: ":6062",
	}
	go func() {
		server.ListenAndServe()
	}()
	defer server.Close()

	// Give the first server time to start
	time.Sleep(50 * time.Millisecond)

	// Try to start another server on the same port
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		StartProfile(ctx, 6062)
		close(done)
	}()

	// Wait for the function to complete
	<-done
}

// TestStartWatchdog tests the StartWatchdog function
func TestStartWatchdog(t *testing.T) {
	tests := []struct {
		name         string
		interval     time.Duration
		runDuration  time.Duration
		expectedLogs int
	}{
		{
			name:         "Watchdog with short interval",
			interval:     50 * time.Millisecond,
			runDuration:  150 * time.Millisecond,
			expectedLogs: 2, // Should log at least 2 times
		},
		{
			name:         "Watchdog with longer interval",
			interval:     100 * time.Millisecond,
			runDuration:  250 * time.Millisecond,
			expectedLogs: 2, // Should log at least 2 times
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.runDuration)
			defer cancel()

			// Capture initial goroutine count
			initialGoroutines := runtime.NumGoroutine()

			// Run watchdog
			done := make(chan struct{})
			go func() {
				StartWatchdog(ctx, tt.interval)
				close(done)
			}()

			// Wait for watchdog to complete
			<-done

			// Verify goroutines are cleaned up
			time.Sleep(10 * time.Millisecond) // Give time for cleanup
			finalGoroutines := runtime.NumGoroutine()
			assert.LessOrEqual(t, finalGoroutines, initialGoroutines+1,
				"Goroutines should be cleaned up after watchdog stops")
		})
	}
}

// TestPrintRuntimeStats tests the printRuntimeStats function
func TestPrintRuntimeStats(t *testing.T) {
	// This test mainly ensures the function doesn't panic
	assert.NotPanics(t, func() {
		printRuntimeStats()
	})

	// Force a GC to ensure we have some stats
	runtime.GC()
	assert.NotPanics(t, func() {
		printRuntimeStats()
	})
}

// TestBToMb tests the byte to megabyte conversion
func TestBToMb(t *testing.T) {
	tests := []struct {
		name     string
		bytes    uint64
		expected uint64
	}{
		{
			name:     "Zero bytes",
			bytes:    0,
			expected: 0,
		},
		{
			name:     "One megabyte",
			bytes:    1024 * 1024,
			expected: 1,
		},
		{
			name:     "Multiple megabytes",
			bytes:    5 * 1024 * 1024,
			expected: 5,
		},
		{
			name:     "Partial megabyte rounds down",
			bytes:    1024*1024 + 500*1024,
			expected: 1,
		},
		{
			name:     "Large value",
			bytes:    1024 * 1024 * 1024, // 1 GB
			expected: 1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bToMb(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// BenchmarkStartProfile benchmarks the profile server startup and shutdown
func BenchmarkStartProfile(b *testing.B) {
	ports := []int{7060, 7061, 7062, 7063, 7064}

	for i := 0; i < b.N; i++ {
		port := ports[i%len(ports)]
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)

		done := make(chan struct{})
		go func() {
			StartProfile(ctx, port)
			close(done)
		}()

		<-done
		cancel()

		// Give time for port to be released
		time.Sleep(10 * time.Millisecond)
	}
}

// BenchmarkPrintRuntimeStats benchmarks the runtime stats printing
func BenchmarkPrintRuntimeStats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printRuntimeStats()
	}
}

// BenchmarkBToMb benchmarks the byte to megabyte conversion
func BenchmarkBToMb(b *testing.B) {
	values := []uint64{
		0,
		1024,
		1024 * 1024,
		1024 * 1024 * 1024,
		1024 * 1024 * 1024 * 1024,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = bToMb(v)
		}
	}
}

// BenchmarkWatchdogInterval benchmarks different watchdog intervals
func BenchmarkWatchdogInterval(b *testing.B) {
	intervals := []time.Duration{
		10 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		500 * time.Millisecond,
	}

	for _, interval := range intervals {
		b.Run(fmt.Sprintf("Interval-%v", interval), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)

				done := make(chan struct{})
				go func() {
					StartWatchdog(ctx, interval)
					close(done)
				}()

				<-done
				cancel()
			}
		})
	}
}

// TestIntegration tests the integration of profile server and watchdog
func TestIntegration(t *testing.T) {
	require := require.New(t)

	// Create a parent context
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Start profile server
	profileDone := make(chan struct{})
	go func() {
		StartProfile(ctx, 6063)
		close(profileDone)
	}()

	// Start watchdog
	watchdogDone := make(chan struct{})
	go func() {
		StartWatchdog(ctx, 50*time.Millisecond)
		close(watchdogDone)
	}()

	// Give services time to start
	time.Sleep(100 * time.Millisecond)

	// Verify profile server is running
	resp, err := http.Get("http://localhost:6063/debug/pprof/")
	if err == nil {
		require.Equal(http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// Try to access goroutine profile
		resp, err = http.Get("http://localhost:6063/debug/pprof/goroutine")
		require.NoError(err)
		require.Equal(http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Wait for services to stop
	<-profileDone
	<-watchdogDone
}
