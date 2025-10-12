//go:build windows
// +build windows

package tz

import (
	"golang.org/x/sys/windows/registry"
)

// systemZone detects the system timezone on Windows by reading from the
// Windows Registry.
//
// This function queries the Windows Registry at the location:
// HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\TimeZoneInformation
//
// It first tries to read the modern "TimeZoneKeyName" value (Windows 10+
// with IANA timezone names), then falls back to the legacy "StandardName"
// value for older Windows versions.
//
// Returns:
//   - The timezone name from the registry (e.g., "Eastern Standard Time")
//   - Empty string if the registry cannot be accessed or no timezone is found
//
// This is a Windows-specific implementation that requires the
// golang.org/x/sys/windows/registry package.
func systemZone() string {
	// Open the Windows Registry key for timezone information
	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\TimeZoneInformation`,
		registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer key.Close()

	// Try modern Windows timezone detection (Windows 10+ with IANA names)
	tz, _, err := key.GetStringValue("TimeZoneKeyName")
	if err == nil && tz != "" {
		return tz
	}

	// Fall back to legacy Windows timezone detection
	tz, _, err = key.GetStringValue("StandardName")
	if err == nil && tz != "" {
		return tz
	}

	// Return empty string if no timezone could be detected
	return ""
}
