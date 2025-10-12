package tz

import (
	"fmt"
	"os"
)

// ExampleDetect demonstrates basic usage of the Detect function.
func ExampleDetect() {
	// Set a timezone via environment variable
	os.Setenv("TZ", "America/New_York")
	defer os.Unsetenv("TZ")

	timezone := Detect()
	if timezone != "" {
		fmt.Printf("Detected timezone: %s\n", timezone)
	} else {
		fmt.Println("Could not detect timezone")
	}
	// Output: Detected timezone: America/New_York
}

// ExampleDetect_systemFallback demonstrates system fallback behavior.
func ExampleDetect_systemFallback() {
	// Unset TZ to trigger system detection
	os.Unsetenv("TZ")

	timezone := Detect()
	if timezone != "" {
		fmt.Printf("System detected timezone: %s\n", timezone)
	} else {
		fmt.Println("No timezone detected")
	}
	// Output: System detected timezone: America/New_York
}
