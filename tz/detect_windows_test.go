//go:build windows
// +build windows

package tz

import (
	"testing"
)

func TestSystemZoneWindows(t *testing.T) {
	// Test the Windows systemZone function
	// Note: This test will only work on Windows systems
	result := systemZone()

	if result != "" {
		t.Logf("Windows detected timezone: %s", result)
		// Common Windows timezone names
		commonWindowsTZs := []string{
			"Eastern Standard Time",
			"Central Standard Time",
			"Mountain Standard Time",
			"Pacific Standard Time",
			"UTC",
			"GMT Standard Time",
		}

		found := false
		for _, tz := range commonWindowsTZs {
			if result == tz {
				found = true
				break
			}
		}

		if !found {
			t.Logf("Warning: detected timezone '%s' is not a common Windows timezone name", result)
		}
	} else {
		t.Log("No timezone detected from Windows registry")
	}
}

func TestSystemZoneWindowsRegistryAccess(t *testing.T) {
	// Test that the function handles registry access gracefully
	// This is more of a smoke test to ensure it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("systemZone() panicked: %v", r)
		}
	}()

	result := systemZone()
	// We don't assert specific values since registry contents vary by system
	// We just ensure it returns a string (empty or valid timezone)
	_ = result
}

func TestSystemZoneWindowsConsistency(t *testing.T) {
	// Test that multiple calls return consistent results
	result1 := systemZone()
	result2 := systemZone()

	if result1 != result2 {
		t.Errorf("systemZone() inconsistent: first call = %v, second call = %v", result1, result2)
	}
}

// Benchmark test for Windows systemZone function
func BenchmarkSystemZoneWindows(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = systemZone()
	}
}
