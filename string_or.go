package sefaria

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type StringOr[T any] struct {
	Value T
}

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

func (s StringOr[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Value)
}
