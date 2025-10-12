package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// StringOr is a generic type that can unmarshal JSON as either a string
// representation of type T or as the actual type T directly.
//
// This type is useful for handling APIs that inconsistently return values
// as either strings or their native types. For example, an API might return
// a boolean as "true" in one response and as true in another.
//
// Supported types for T:
//   - int: parses string numbers to integers
//   - bool: parses "true"/"false" strings to booleans
//   - float64: parses string numbers to floats
//   - string: passes through string values
//
// Example usage:
//
//	var priority StringOr[float32]
//	json.Unmarshal([]byte(`"3.14"`), &priority)  // Works
//	json.Unmarshal([]byte(`3.14`), &priority)    // Also works
//
//	var enabled StringOr[bool]
//	json.Unmarshal([]byte(`"true"`), &enabled)   // Works
//	json.Unmarshal([]byte(`true`), &enabled)     // Also works
type StringOr[T any] struct {
	Value T
}

// UnmarshalJSON implements json.Unmarshaler for StringOr.
// It first attempts to unmarshal the JSON as a string, then converts
// that string to type T. If that fails, it attempts to unmarshal
// directly as type T.
//
// Empty strings and "null" values are treated as the zero value of T.
func (s *StringOr[T]) UnmarshalJSON(data []byte) error {
	// Try to decode as a raw string
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str == "" || str == "null" {
			var zero T
			s.Value = zero
			return nil
		}
		var v any
		switch any(*new(T)).(type) {
		case int:
			i, err := strconv.Atoi(str)
			if err != nil {
				return fmt.Errorf("invalid int string %q: %w", str, err)
			}
			v = i
		case bool:
			b, err := strconv.ParseBool(str)
			if err != nil {
				return fmt.Errorf("invalid bool string %q: %w", str, err)
			}
			v = b
		case float64:
			f, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return fmt.Errorf("invalid float string %q: %w", str, err)
			}
			v = f
		case string:
			v = str
		default:
			return fmt.Errorf("unsupported type for StringOr: %T", *new(T))
		}
		s.Value = v.(T)
		return nil
	}

	// Try direct decode as T
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("StringOr: cannot decode %s: %w", string(data), err)
	}
	s.Value = v
	return nil
}

// MarshalJSON implements json.Marshaler for StringOr.
// It marshals the underlying Value directly.
func (s StringOr[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Value)
}
