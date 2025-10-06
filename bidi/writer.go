package bidi

import (
	"bytes"
	"io"
	"unicode"
)

type Writer struct {
	w       io.Writer
	flipRTL bool // optional flag to flip text even without markers
}

func NewWriter(w io.Writer, flipRTL bool) *Writer {
	return &Writer{w: w, flipRTL: flipRTL}
}

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

// Detect Hebrew / Arabic
func isRTL(r rune) bool {
	return unicode.In(r, unicode.Hebrew, unicode.Arabic)
}
