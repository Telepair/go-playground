package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"runtime"
	"time"
)

// StartProfile starts a pprof server and handles graceful shutdown
func StartProfile(ctx context.Context, port int) {
	server := &http.Server{ //nolint:gosec
		Addr: fmt.Sprintf(":%d", port),
	}

	go func() {
		slog.Info("Starting pprof server on http://localhost%s/debug/pprof/", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start pprof server", "error", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	slog.Info("Stopping pprof server")

	// Create a timeout context for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Failed to gracefully shutdown pprof server", "error", err)
	} else {
		slog.Info("Pprof server stopped gracefully")
	}
}

// StartWatchdog periodically prints runtime profile information
func StartWatchdog(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	slog.Info("Starting watchdog with interval", "interval", interval)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping watchdog")
			return
		case <-ticker.C:
			printRuntimeStats()
		}
	}
}

// printRuntimeStats prints current runtime statistics
func printRuntimeStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	logger := slog.With("module", "watchdog")
	logger.Info("Runtime Stats",
		"goroutines", runtime.NumGoroutine(),
		"alloc_mb", bToMb(m.Alloc),
		"total_alloc_mb", bToMb(m.TotalAlloc),
		"sys_mb", bToMb(m.Sys),
		"num_gc", m.NumGC,
		"next_gc_mb", bToMb(m.NextGC),
	)
}

// bToMb converts bytes to megabytes
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
