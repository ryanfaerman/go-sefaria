package tz

import (
	"os"
	"time"
)

// Detect returns the system timezone name using a hierarchical detection approach.
//
// The detection follows this order of precedence:
//  1. TZ environment variable (if set)
//  2. Go's time.Local timezone (if not "Local")
//  3. System-specific detection:
//     - Unix: /etc/localtime symlink analysis
//     - Windows: Registry-based detection
//
// The detection is cached in-memory and is DST-aware, automatically invalidating
// the cache when UTC offset changes occur.
//
// Returns:
//   - The detected timezone name (e.g., "America/New_York", "Europe/London")
//   - Empty string if no timezone could be detected
//
// Example:
//
//	timezone := tz.Detect()
//	if timezone != "" {
//	    fmt.Printf("Detected timezone: %s\n", timezone)
//	}
func Detect() string {
	// 1. Check TZ environment variable first (highest priority)
	if tz := os.Getenv("TZ"); tz != "" {
		return tz
	}

	// 2. Try Go's time.Local timezone (if it's not the generic "Local")
	if loc := time.Now().Location(); loc != nil && loc.String() != "Local" {
		return loc.String()
	}

	// 3. Fall back to system-specific detection
	if tz := systemZone(); tz != "" {
		return tz
	}

	// 4. Return empty string if no timezone could be detected
	return ""
}
