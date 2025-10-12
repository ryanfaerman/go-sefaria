package tz

import (
	"os"
	"testing"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name     string
		tzEnv    string
		expected string
	}{
		{
			name:     "TZ environment variable set",
			tzEnv:    "America/New_York",
			expected: "America/New_York",
		},
		{
			name:     "TZ environment variable with UTC offset",
			tzEnv:    "UTC+5",
			expected: "UTC+5",
		},
		{
			name:     "TZ environment variable with legacy timezone",
			tzEnv:    "EST",
			expected: "EST",
		},
		{
			name:     "Empty TZ environment variable",
			tzEnv:    "",
			expected: "", // Will fall back to system detection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.tzEnv != "" {
				os.Setenv("TZ", tt.tzEnv)
			} else {
				os.Unsetenv("TZ")
			}

			// Test detection
			result := Detect()

			// For TZ environment variable tests, we expect exact match
			if tt.tzEnv != "" {
				if result != tt.expected {
					t.Errorf("Detect() = %v, want %v", result, tt.expected)
				}
			} else {
				// For empty TZ, we just verify it doesn't panic and returns a string
				// (could be empty or a detected timezone)
				if result == "" {
					t.Logf("No timezone detected (this may be expected on some systems)")
				} else {
					t.Logf("Detected timezone: %s", result)
				}
			}

			// Clean up
			os.Unsetenv("TZ")
		})
	}
}

func TestDetectWithTZEnvPriority(t *testing.T) {
	// Test that TZ environment variable takes priority over system detection
	os.Setenv("TZ", "America/Los_Angeles")
	defer os.Unsetenv("TZ")

	result := Detect()
	if result != "America/Los_Angeles" {
		t.Errorf("Detect() with TZ env var = %v, want America/Los_Angeles", result)
	}
}

func TestDetectConsistency(t *testing.T) {
	// Test that Detect() returns consistent results
	os.Setenv("TZ", "Europe/London")
	defer os.Unsetenv("TZ")

	result1 := Detect()
	result2 := Detect()

	if result1 != result2 {
		t.Errorf("Detect() inconsistent: first call = %v, second call = %v", result1, result2)
	}
}

func TestDetectEmptyResult(t *testing.T) {
	// Test behavior when no timezone can be detected
	os.Unsetenv("TZ")

	result := Detect()
	// Result could be empty or a detected timezone depending on system
	// We just verify it doesn't panic
	if result == "" {
		t.Log("No timezone detected (expected on some systems)")
	} else {
		t.Logf("System detected timezone: %s", result)
	}
}
