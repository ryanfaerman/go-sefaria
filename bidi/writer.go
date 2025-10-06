package bidi

import (
	"bytes"
	"io"
	"unicode"
)

// Writer is an io.Writer that processes bidirectional text by reversing
// RTL sequences marked with Unicode directional markers. It can also
// optionally reverse all RTL text when the flipRTL flag is enabled.
type Writer struct {
	w       io.Writer
	flipRTL bool // optional flag to flip text even without markers
}

// NewWriter creates a new Writer that wraps the provided io.Writer.
// If flipRTL is true, all RTL text will be reversed even without directional markers.
// If false, only text between RLM/LRM markers will be reversed.
func NewWriter(w io.Writer, flipRTL bool) *Writer {
	return &Writer{w: w, flipRTL: flipRTL}
}

// Write implements io.Writer, processing bidirectional text by reversing
// RTL sequences. Text between RLM (U+200F) and LRM (U+200E) markers
// is reversed, and if flipRTL is enabled, all RTL text is reversed.
func (bw *Writer) Write(p []byte) (int, error) {
	// Convert to runes
	runes := []rune(string(p))
	var out bytes.Buffer

	i := 0
	for i < len(runes) {
		switch runes[i] {
		case '\u200F': // RLM start
			start := i + 1
			end := start
			for end < len(runes) && runes[end] != '\u200E' { // until LRM
				end++
			}
			// flip the text between RLM and LRM
			for j := end - 1; j >= start; j-- {
				out.WriteRune(runes[j])
			}
			i = end + 1 // skip past LRM
		default:
			// optionally flip all RTL text if flag is set
			if bw.flipRTL && isRTL(runes[i]) {
				start := i
				for i < len(runes) && isRTL(runes[i]) {
					i++
				}
				for j := i - 1; j >= start; j-- {
					out.WriteRune(runes[j])
				}
			} else {
				out.WriteRune(runes[i])
				i++
			}
		}
	}

	n, err := bw.w.Write(out.Bytes())
	if n != len(out.Bytes()) && err == nil {
		err = io.ErrShortWrite
	}
	return len(p), err
}

// isRTL detects if a rune belongs to a right-to-left script.
// Currently supports Hebrew and Arabic scripts.
func isRTL(r rune) bool {
	return unicode.In(r, unicode.Hebrew, unicode.Arabic)
}
