// Package tz provides timezone detection functionality for cross-platform systems.
//
// The package automatically detects the system's timezone using a hierarchical
// approach that works across Unix-like systems (Linux, macOS, BSD) and Windows.
//
// Detection Order:
//  1. TZ environment variable (if set)
//  2. Go's time.Local timezone
//  3. System-specific detection:
//     - Unix: /etc/localtime symlink analysis
//     - Windows: Registry-based detection
//
// The detection is cached in-memory and is DST-aware, automatically invalidating
// the cache when UTC offset changes occur.
//
// Example usage:
//
//	import "github.com/yourorg/go-sefaria/tz"
//
//	timezone := tz.Detect()
//	if timezone != "" {
//		fmt.Printf("Detected timezone: %s\n", timezone)
//	} else {
//		fmt.Println("Could not detect timezone")
//	}
//
// Supported timezone formats:
// - IANA timezone names (e.g., "America/New_York", "Europe/London")
// - Legacy timezone names (Windows fallback, e.g., "Eastern Standard Time")
// - Any timezone string set via the TZ environment variable
package tz
