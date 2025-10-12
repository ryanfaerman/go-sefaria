package types

import (
	"fmt"
	"net/url"
)

// BoolInt represents a boolean that encodes as "1" or "0" for URL parameters
// and other contexts where boolean values need to be represented as integers.
//
// This type is particularly useful for APIs that expect boolean values
// to be encoded as integers (1 for true, 0 for false) rather than
// the standard "true"/"false" strings.
//
// BoolInt implements several interfaces:
//   - url.Encoder for URL parameter encoding
//   - flag.Value for command-line flag parsing
//   - encoding.TextMarshaler/TextUnmarshaler for text encoding
//
// Example usage:
//
//	var b BoolInt = true
//	// When used in URL parameters, encodes as "1"
//	// When used in flags, accepts "1", "0", "true", "false"
type BoolInt bool

// EncodeValues implements url.Encoder for BoolInt.
// It encodes true as "1" and false as "0" for URL parameters.
func (b BoolInt) EncodeValues(key string, v *url.Values) error {
	if b {
		v.Add(key, "1")
	} else {
		v.Add(key, "0")
	}
	return nil
}

// String returns the string representation of BoolInt.
// Returns "1" for true, "0" for false.
func (b *BoolInt) String() string {
	if *b {
		return "1"
	}
	return "0"
}

// Set implements flag.Value for BoolInt.
// Accepts "1", "0", "true", or "false" as valid values.
func (b *BoolInt) Set(s string) error {
	switch s {
	case "1", "true":
		*b = true
	case "0", "false":
		*b = false
	default:
		return fmt.Errorf("invalid bool: %q", s)
	}
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler for BoolInt.
// Accepts "1", "0", "true", or "false" as valid values.
func (b *BoolInt) UnmarshalText(text []byte) error {
	s := string(text)
	switch s {
	case "1", "true":
		*b = true
	case "0", "false":
		*b = false
	default:
		return fmt.Errorf("invalid BoolInt: %q", s)
	}
	return nil
}

// MarshalText implements encoding.TextMarshaler for BoolInt.
// Returns "1" for true, "0" for false.
func (b BoolInt) MarshalText() ([]byte, error) {
	if b {
		return []byte("1"), nil
	}
	return []byte("0"), nil
}
