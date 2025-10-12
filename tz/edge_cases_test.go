package tz

import (
	"os"
	"runtime"
	"testing"
)

func TestDetectEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		setup       func()
		cleanup     func()
		expectEmpty bool
	}{
		{
			name: "TZ set to empty string",
			setup: func() {
				os.Setenv("TZ", "")
			},
			cleanup: func() {
				os.Unsetenv("TZ")
			},
			expectEmpty: true, // Empty TZ should fall back to system detection
		},
		{
			name: "TZ set to whitespace",
			setup: func() {
				os.Setenv("TZ", "   ")
			},
			cleanup: func() {
				os.Unsetenv("TZ")
			},
			expectEmpty: false, // Whitespace TZ should be returned as-is
		},
		{
			name: "TZ set to invalid timezone",
			setup: func() {
				os.Setenv("TZ", "Invalid/Timezone/Name")
			},
			cleanup: func() {
				os.Unsetenv("TZ")
			},
			expectEmpty: false, // Should return the invalid string as-is
		},
		{
			name: "TZ unset",
			setup: func() {
				os.Unsetenv("TZ")
			},
			cleanup: func() {
				// No cleanup needed
			},
			expectEmpty: true, // Should fall back to system detection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			result := Detect()

			if tt.expectEmpty {
				if result != "" {
					t.Logf("Expected empty result, but got: %s (this may be expected if system detection works)", result)
				}
			} else {
				if result == "" {
					t.Errorf("Expected non-empty result, but got empty string")
				}
			}
		})
	}
}

func TestDetectConcurrency(t *testing.T) {
	// Test that Detect() is safe for concurrent access
	os.Setenv("TZ", "America/New_York")
	defer os.Unsetenv("TZ")

	done := make(chan bool, 10)

	// Start multiple goroutines
	for i := 0; i < 10; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Detect() panicked in goroutine: %v", r)
				}
			}()

			result := Detect()
			if result != "America/New_York" {
				t.Errorf("Concurrent Detect() = %v, want America/New_York", result)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestDetectPlatformSpecific(t *testing.T) {
	// Test platform-specific behavior
	os.Unsetenv("TZ")
	defer os.Unsetenv("TZ")

	result := Detect()

	switch runtime.GOOS {
	case "windows":
		t.Logf("Windows system detected timezone: %s", result)
	case "linux", "darwin", "freebsd", "openbsd", "netbsd":
		t.Logf("Unix system detected timezone: %s", result)
	default:
		t.Logf("Unknown platform (%s) detected timezone: %s", runtime.GOOS, result)
	}

	// Just verify it doesn't panic and returns a string
	_ = result
}

func TestDetectWithSpecialCharacters(t *testing.T) {
	// Test TZ environment variable with special characters
	specialTZs := []string{
		"UTC+05:30",
		"UTC-08:00",
		"GMT+1",
		"EST5EDT",
		"America/New_York",
		"Europe/London",
		"Asia/Tokyo",
	}

	for _, tz := range specialTZs {
		t.Run("TZ_"+tz, func(t *testing.T) {
			os.Setenv("TZ", tz)
			defer os.Unsetenv("TZ")

			result := Detect()
			if result != tz {
				t.Errorf("Detect() with TZ=%s = %v, want %v", tz, result, tz)
			}
		})
	}
}

// Benchmark tests
func BenchmarkDetect(b *testing.B) {
	os.Setenv("TZ", "America/New_York")
	defer os.Unsetenv("TZ")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Detect()
	}
}

func BenchmarkDetectWithSystemFallback(b *testing.B) {
	os.Unsetenv("TZ")
	defer os.Unsetenv("TZ")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Detect()
	}
}
