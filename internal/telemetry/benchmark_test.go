package telemetry

import (
	"fmt"
	"testing"
	"time"
)

// BenchmarkClientTrack measures the performance of event tracking
func BenchmarkClientTrack(b *testing.B) {
	client := NewClient("")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create new properties map for each iteration to avoid race conditions
		properties := map[string]interface{}{
			"command":     "benchmark_test",
			"duration_ms": 100,
			"success":     true,
			"timestamp":   time.Now().Unix(),
		}
		client.Track("benchmark_event", properties)
	}
}

// BenchmarkClientTrackCommand measures command tracking performance
func BenchmarkClientTrackCommand(b *testing.B) {
	client := NewClient("")
	duration := 50 * time.Millisecond

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.TrackCommand("benchmark_command", duration, true)
	}
}

// BenchmarkClientTrackError measures error tracking performance
func BenchmarkClientTrackError(b *testing.B) {
	client := NewClient("")
	err := fmt.Errorf("benchmark test error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create new properties map for each iteration to avoid race conditions
		properties := map[string]interface{}{
			"context": "benchmark",
		}
		client.TrackError(err, "benchmark_command", properties)
	}
}

// BenchmarkMiddlewareTrack measures global middleware performance
func BenchmarkMiddlewareTrack(b *testing.B) {
	Initialize("")
	defer Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create new properties map for each iteration to avoid race conditions
		properties := map[string]interface{}{
			"source": "middleware_benchmark",
		}
		Track("middleware_benchmark_event", properties)
	}
}

// BenchmarkAnonymousIDGeneration measures ID generation performance
func BenchmarkAnonymousIDGeneration(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generateAnonymousID()
	}
}

// BenchmarkSessionIDGeneration measures session ID generation performance
func BenchmarkSessionIDGeneration(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generateSessionID()
	}
}

// BenchmarkClientCreation measures client initialization performance
func BenchmarkClientCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := NewClient("https://test@sentry.io/123")
		_ = client
	}
}

// BenchmarkDisabledClient measures performance when telemetry is disabled
func BenchmarkDisabledClient(b *testing.B) {
	// Create a disabled client
	client := &Client{
		enabled:   false,
		userID:    "test_user",
		sessionID: "test_session",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create new properties map for each iteration to avoid race conditions
		properties := map[string]interface{}{
			"test": "disabled_benchmark",
		}
		client.Track("disabled_event", properties)
	}
}

// BenchmarkCommandStartEnd measures start/end command tracking
func BenchmarkCommandStartEnd(b *testing.B) {
	client := NewClient("") // No external services for benchmark
	duration := 25 * time.Millisecond

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.StartCommand("benchmark_command")
		client.EndCommand("benchmark_command", true, duration)
	}
}
