//go:build !windows
// +build !windows

package tz

import (
	"os"
	"strings"
)

// systemZone detects the system timezone on Unix-like systems by analyzing
// the /etc/localtime symlink.
//
// This function reads the symlink target of /etc/localtime and extracts the
// timezone name from the path. It handles both direct symlinks to zoneinfo
// files and symlinks within the /usr/share/zoneinfo/ directory structure.
//
// Returns:
//   - The timezone name extracted from the symlink (e.g., "America/New_York")
//   - Empty string if the symlink cannot be read or doesn't contain a valid timezone
//
// This is a Unix-specific implementation that works on Linux, macOS, and BSD systems.
func systemZone() string {
	const zoneInfoPrefix = "/usr/share/zoneinfo/"

	// Read the symlink target of /etc/localtime
	link, err := os.Readlink("/etc/localtime")
	if err != nil {
		return ""
	}

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
