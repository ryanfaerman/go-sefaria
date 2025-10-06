package bidi

import (
	"bytes"
	"encoding/json"
	"unicode"
)

// String is a custom string type that provides automatic bidirectional text
// support. It wraps RTL (right-to-left) text sequences with Unicode directional
// markers to ensure proper display in mixed-language contexts.
type String string

// String implements fmt.Stringer, automatically wrapping RTL text with
// Unicode directional markers (RLM/LRM) for proper bidirectional display.
func (s String) String() string {
	return wrapRTL(string(s))
}

// MarshalJSON implements json.Marshaler, applying RTL wrapping when marshaling
// to JSON. This ensures that RTL text is properly marked for bidirectional
// display when serialized.
func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON implements json.Unmarshaler, storing RTL content as-is without
// applying directional markers. The markers will be added automatically when
// the string is displayed via String() or MarshalJSON().
func (s *String) UnmarshalJSON(data []byte) error {
	// Trim quotes and unmarshal
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*s = String(raw)
	return nil
}

// wrapRTL wraps contiguous RTL text sequences with Unicode directional markers.
// It adds RLM (Right-to-Left Mark, U+200F) at the start and LRM (Left-to-Right
// Mark, U+200E) at the end of each RTL sequence, including spaces between
// RTL words. This ensures proper bidirectional text rendering.
// It also prevents double-wrapping by detecting existing markers.
func wrapRTL(str string) string {
	var buf bytes.Buffer
	runes := []rune(str)
	i := 0
	for i < len(runes) {
		// Check if we're already at a marker sequence
		if runes[i] == '\u200F' { // RLM
			// Skip the entire marked sequence
			start := i
			i++ // skip RLM
			// Find the matching LRM
			foundLRM := false
			for i < len(runes) && runes[i] != '\u200E' {
				i++
			}
			if i < len(runes) && runes[i] == '\u200E' {
				i++ // skip LRM
				foundLRM = true
			}
			// Copy the entire marked sequence as-is
			for j := start; j < i; j++ {
				buf.WriteRune(runes[j])
			}
			// If we didn't find a matching LRM, add one
			if !foundLRM {
				buf.WriteRune('\u200E')
			}
		} else if runes[i] == '\u200E' { // LRM without matching RLM
			// This is a malformed sequence, treat the LRM as regular text
			buf.WriteRune(runes[i])
			i++
		} else if isRTL(runes[i]) {
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
