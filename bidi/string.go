package bidi

import (
	"bytes"
	"encoding/json"
	"unicode"
)

type String string

// String implements fmt.Stringer, wrapping RTL text automatically.
func (s String) String() string {
	return wrapRTL(string(s))
}

// MarshalJSON implements json.Marshaler, applying RTL wrapping when marshaling.
func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON detects Hebrew (or RTL) content and stores it as-is.
func (s *String) UnmarshalJSON(data []byte) error {
	// Trim quotes and unmarshal
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*s = String(raw)
	return nil
}

// wrapRTL wraps any contiguous RTL text with RLM/LRM marks
func wrapRTL(str string) string {
	var buf bytes.Buffer
	runes := []rune(str)
	i := 0
	for i < len(runes) {
		if isRTL(runes[i]) {
			buf.WriteRune('\u200F') // RLM
			// Find the end of the RTL sequence, including spaces between RTL words
			start := i
			for i < len(runes) && (isRTL(runes[i]) || unicode.IsSpace(runes[i])) {
				// If we hit a space, check if it's followed by more RTL text
				if unicode.IsSpace(runes[i]) {
					// Look ahead to see if there's more RTL text after spaces
					j := i + 1
					for j < len(runes) && unicode.IsSpace(runes[j]) {
						j++
					}
					if j < len(runes) && isRTL(runes[j]) {
						i = j // include the spaces and continue with RTL text
					} else {
						break // stop if spaces are followed by LTR text
					}
				} else {
					i++
				}
			}
			// Write the RTL text as-is (the Writer will handle reversal)
			for j := start; j < i; j++ {
				buf.WriteRune(runes[j])
			}
			buf.WriteRune('\u200E') // LRM
		} else {
			buf.WriteRune(runes[i])
			i++
		}
	}
	return buf.String()
}
