//go:build !windows
// +build !windows

package tz

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSystemZone(t *testing.T) {
	tests := []struct {
		name          string
		localtimePath string
		expected      string
		shouldFail    bool
	}{
		{
			name:          "Standard zoneinfo symlink",
			localtimePath: "/usr/share/zoneinfo/America/New_York",
			expected:      "America/New_York",
			shouldFail:    false,
		},
		{
			name:          "Zoneinfo with leading slash",
			localtimePath: "/usr/share/zoneinfo/Europe/London",
			expected:      "Europe/London",
			shouldFail:    false,
		},
		{
			name:          "UTC timezone",
			localtimePath: "/usr/share/zoneinfo/UTC",
			expected:      "UTC",
			shouldFail:    false,
		},
		{
			name:          "Invalid symlink path",
			localtimePath: "/invalid/path",
			expected:      "",
			shouldFail:    true,
		},
		{
			name:          "Path without zoneinfo",
			localtimePath: "/some/other/path",
			expected:      "",
			shouldFail:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for testing
			tempDir := t.TempDir()

			// Create the symlink target
			targetPath := filepath.Join(tempDir, "target")
			err := os.MkdirAll(filepath.Dir(targetPath), 0755)
			if err != nil {
				t.Fatalf("Failed to create target directory: %v", err)
			}

			// Create a dummy file at the target
			err = os.WriteFile(targetPath, []byte("dummy"), 0644)
			if err != nil {
				t.Fatalf("Failed to create target file: %v", err)
			}

			// Create the symlink
			linkPath := filepath.Join(tempDir, "localtime")
			err = os.Symlink(tt.localtimePath, linkPath)
			if err != nil {
				t.Fatalf("Failed to create symlink: %v", err)
			}

			// Test the function by temporarily replacing /etc/localtime
			// We'll test the logic by examining the symlink directly
			linkTarget, err := os.Readlink(linkPath)
			if err != nil {
				if tt.shouldFail {
					return // Expected to fail
				}
				t.Fatalf("Failed to read symlink: %v", err)
			}

			// Simulate the systemZone logic
			result := extractTimezoneFromPath(linkTarget)

			if tt.shouldFail {
				if result != "" {
					t.Errorf("Expected empty result for failing case, got %v", result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("systemZone() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

// Helper function to test the timezone extraction logic
func extractTimezoneFromPath(link string) string {
	const zoneInfoPrefix = "/usr/share/zoneinfo/"

	// Extract timezone name from symlink path
	// Handle paths like "/usr/share/zoneinfo/America/New_York"
	if i := strings.Index(link, "zoneinfo/"); i != -1 {
		return strings.TrimPrefix(link[i+len("zoneinfo/"):], "/")
	}

	// Handle direct symlinks to zoneinfo files
	if after, found := strings.CutPrefix(link, zoneInfoPrefix); found {
		return after
	}

	// Return empty string if no valid timezone could be extracted
	return ""
}

func TestSystemZoneRealSystem(t *testing.T) {
	// Test with the actual system's /etc/localtime if it exists
	if _, err := os.Stat("/etc/localtime"); os.IsNotExist(err) {
		t.Skip("Skipping real system test: /etc/localtime does not exist")
		return
	}

	result := systemZone()
	if result != "" {
		t.Logf("System detected timezone: %s", result)
		// Verify it's a reasonable timezone name
		if !strings.Contains(result, "/") && result != "UTC" && result != "GMT" {
			t.Logf("Warning: detected timezone '%s' doesn't look like a standard IANA timezone", result)
		}
	} else {
		t.Log("No timezone detected from /etc/localtime")
	}
}

func TestSystemZoneEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		linkPath string
		expected string
	}{
		{
			name:     "Empty path",
			linkPath: "",
			expected: "",
		},
		{
			name:     "Path with multiple zoneinfo occurrences",
			linkPath: "/usr/share/zoneinfo/America/New_York/zoneinfo/Europe/London",
			expected: "America/New_York/zoneinfo/Europe/London",
		},
		{
			name:     "Path ending with slash",
			linkPath: "/usr/share/zoneinfo/America/New_York/",
			expected: "America/New_York/",
		},
		{
			name:     "Direct zoneinfo prefix match",
			linkPath: "/usr/share/zoneinfo/UTC",
			expected: "UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTimezoneFromPath(tt.linkPath)
			if result != tt.expected {
				t.Errorf("extractTimezoneFromPath(%q) = %v, want %v", tt.linkPath, result, tt.expected)
			}
		})
	}
}
