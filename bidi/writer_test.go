package bidi

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWriter(t *testing.T) {
	tests := []struct {
		name     string
		w        io.Writer
		flipRTL  bool
		expected *Writer
	}{
		{
			name:    "with flipRTL enabled",
			w:       &bytes.Buffer{},
			flipRTL: true,
			expected: &Writer{
				w:       &bytes.Buffer{},
				flipRTL: true,
			},
		},
		{
			name:    "with flipRTL disabled",
			w:       &bytes.Buffer{},
			flipRTL: false,
			expected: &Writer{
				w:       &bytes.Buffer{},
				flipRTL: false,
			},
		},
		{
			name:    "with nil writer",
			w:       nil,
			flipRTL: true,
			expected: &Writer{
				w:       nil,
				flipRTL: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewWriter(tt.w, tt.flipRTL)
			assert.Equal(t, tt.expected.w, result.w)
			assert.Equal(t, tt.expected.flipRTL, result.flipRTL)
		})
	}
}

func TestWriter_Write_WithMarkers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		flipRTL  bool
	}{
		{
			name:     "empty input",
			input:    "",
			expected: "",
			flipRTL:  false,
		},
		{
			name:     "LTR only text",
			input:    "Hello World",
			expected: "Hello World",
			flipRTL:  false,
		},
		{
			name:     "Hebrew text with markers",
			input:    "\u200Fשלום\u200E",
			expected: "םולש",
			flipRTL:  false,
		},
		{
			name:     "Arabic text with markers",
			input:    "\u200Fمرحبا\u200E",
			expected: "ابحرم",
			flipRTL:  false,
		},
		{
			name:     "mixed LTR and Hebrew with markers",
			input:    "Hello \u200Fשלום\u200E World",
			expected: "Hello םולש World",
			flipRTL:  false,
		},
		{
			name:     "multiple RTL sequences with markers",
			input:    "Hello \u200Fשלום\u200E World \u200Fمرحبا\u200E Test",
			expected: "Hello םולש World ابحرم Test",
			flipRTL:  false,
		},
		{
			name:     "RTL text with spaces in markers",
			input:    "\u200Fשלום עולם\u200E",
			expected: "םלוע םולש",
			flipRTL:  false,
		},
		{
			name:     "RTL text with tabs in markers",
			input:    "\u200Fשלום\tעולם\u200E",
			expected: "םלוע\tםולש",
			flipRTL:  false,
		},
		{
			name:     "RTL text with newlines in markers",
			input:    "\u200Fשלום\nעולם\u200E",
			expected: "םלוע\nםולש",
			flipRTL:  false,
		},
		{
			name:     "malformed marker sequence - missing LRM",
			input:    "\u200Fשלום",
			expected: "םולש",
			flipRTL:  false,
		},
		{
			name:     "malformed marker sequence - missing RLM",
			input:    "שלום\u200E",
			expected: "שלום\u200E",
			flipRTL:  false,
		},
		{
			name:     "empty marker sequence",
			input:    "\u200F\u200E",
			expected: "",
			flipRTL:  false,
		},
		{
			name:     "marker sequence with only spaces",
			input:    "\u200F   \u200E",
			expected: "   ",
			flipRTL:  false,
		},
		{
			name:     "complex mixed text with markers",
			input:    "The word \u200Fשלום\u200E means 'hello' in Hebrew, and \u200Fمرحبا\u200E means 'hello' in Arabic.",
			expected: "The word םולש means 'hello' in Hebrew, and ابحرم means 'hello' in Arabic.",
			flipRTL:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewWriter(&buf, tt.flipRTL)

			n, err := writer.Write([]byte(tt.input))
			require.NoError(t, err)
			assert.Equal(t, len(tt.input), n)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestWriter_Write_WithFlipRTL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		flipRTL  bool
	}{
		{
			name:     "Hebrew text without markers - flipRTL enabled",
			input:    "שלום",
			expected: "םולש",
			flipRTL:  true,
		},
		{
			name:     "Arabic text without markers - flipRTL enabled",
			input:    "مرحبا",
			expected: "ابحرم",
			flipRTL:  true,
		},
		{
			name:     "mixed LTR and Hebrew - flipRTL enabled",
			input:    "Hello שלום World",
			expected: "Hello םולש World",
			flipRTL:  true,
		},
		{
			name:     "multiple RTL sequences - flipRTL enabled",
			input:    "Hello שלום World مرحبا Test",
			expected: "Hello םולש World ابحرم Test",
			flipRTL:  true,
		},
		{
			name:     "RTL text with spaces - flipRTL enabled",
			input:    "שלום עולם",
			expected: "םולש םלוע",
			flipRTL:  true,
		},
		{
			name:     "RTL text with tabs - flipRTL enabled",
			input:    "שלום\tעולם",
			expected: "םולש\tםלוע",
			flipRTL:  true,
		},
		{
			name:     "RTL text with newlines - flipRTL enabled",
			input:    "שלום\nעולם",
			expected: "םולש\nםלוע",
			flipRTL:  true,
		},
		{
			name:     "Hebrew text without markers - flipRTL disabled",
			input:    "שלום",
			expected: "שלום",
			flipRTL:  false,
		},
		{
			name:     "Arabic text without markers - flipRTL disabled",
			input:    "مرحبا",
			expected: "مرحبا",
			flipRTL:  false,
		},
		{
			name:     "mixed LTR and Hebrew - flipRTL disabled",
			input:    "Hello שלום World",
			expected: "Hello שלום World",
			flipRTL:  false,
		},
		{
			name:     "LTR only text - flipRTL enabled",
			input:    "Hello World",
			expected: "Hello World",
			flipRTL:  true,
		},
		{
			name:     "numbers and punctuation - flipRTL enabled",
			input:    "Price: 100₪",
			expected: "Price: 100₪",
			flipRTL:  true,
		},
		{
			name:     "complex mixed text - flipRTL enabled",
			input:    "The word שלום means 'hello' in Hebrew, and مرحبا means 'hello' in Arabic.",
			expected: "The word םולש means 'hello' in Hebrew, and ابحرم means 'hello' in Arabic.",
			flipRTL:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewWriter(&buf, tt.flipRTL)

			n, err := writer.Write([]byte(tt.input))
			require.NoError(t, err)
			assert.Equal(t, len(tt.input), n)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestWriter_Write_MultipleWrites(t *testing.T) {
	tests := []struct {
		name     string
		inputs   []string
		expected string
		flipRTL  bool
	}{
		{
			name:     "multiple writes with markers",
			inputs:   []string{"Hello ", "\u200Fשלום\u200E", " World"},
			expected: "Hello םולש World",
			flipRTL:  false,
		},
		{
			name:     "multiple writes with flipRTL",
			inputs:   []string{"Hello ", "שלום", " World"},
			expected: "Hello םולש World",
			flipRTL:  true,
		},
		{
			name:     "marker split across writes",
			inputs:   []string{"\u200Fשל", "ום\u200E"},
			expected: "לשום\u200E",
			flipRTL:  false,
		},
		{
			name:     "RTL text split across writes with flipRTL",
			inputs:   []string{"של", "ום"},
			expected: "לשםו",
			flipRTL:  true,
		},
		{
			name:     "empty writes",
			inputs:   []string{"", "Hello", "", " World", ""},
			expected: "Hello World",
			flipRTL:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewWriter(&buf, tt.flipRTL)

			totalWritten := 0
			for _, input := range tt.inputs {
				n, err := writer.Write([]byte(input))
				require.NoError(t, err)
				totalWritten += n
			}

			expectedTotal := 0
			for _, input := range tt.inputs {
				expectedTotal += len(input)
			}

			assert.Equal(t, expectedTotal, totalWritten)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestWriter_Write_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		flipRTL  bool
	}{
		{
			name:     "only RLM marker",
			input:    "\u200F",
			expected: "",
			flipRTL:  false,
		},
		{
			name:     "only LRM marker",
			input:    "\u200E",
			expected: "\u200E",
			flipRTL:  false,
		},
		{
			name:     "consecutive RLM markers",
			input:    "\u200F\u200Fשלום\u200E",
			expected: "םולש\u200F",
			flipRTL:  false,
		},
		{
			name:     "consecutive LRM markers",
			input:    "\u200Fשלום\u200E\u200E",
			expected: "םולש\u200E",
			flipRTL:  false,
		},
		{
			name:     "nested markers",
			input:    "\u200Fשל\u200Fום\u200E\u200E",
			expected: "םו\u200Fלש\u200E",
			flipRTL:  false,
		},
		{
			name:     "markers with LTR text inside",
			input:    "\u200Fשלום Hello עולם\u200E",
			expected: "םלוע olleH םולש",
			flipRTL:  false,
		},
		{
			name:     "very long RTL text with markers",
			input:    "\u200F" + strings.Repeat("שלום", 100) + "\u200E",
			expected: strings.Repeat("םולש", 100),
			flipRTL:  false,
		},
		{
			name:     "very long RTL text without markers - flipRTL",
			input:    strings.Repeat("שלום", 100),
			expected: strings.Repeat("םולש", 100),
			flipRTL:  true,
		},
		{
			name:     "Unicode emoji with RTL text",
			input:    "Hello \u200Fשלום\u200E 😀 World",
			expected: "Hello םולש 😀 World",
			flipRTL:  false,
		},
		{
			name:     "numbers and currency with RTL",
			input:    "Price: \u200F100₪\u200E",
			expected: "Price: ₪001",
			flipRTL:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewWriter(&buf, tt.flipRTL)

			n, err := writer.Write([]byte(tt.input))
			require.NoError(t, err)
			assert.Equal(t, len(tt.input), n)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestWriter_Write_ErrorHandling(t *testing.T) {
	t.Run("short write error", func(t *testing.T) {
		// Create a writer that only writes partial data
		shortWriter := &shortWriter{maxBytes: 5}
		writer := NewWriter(shortWriter, false)

		input := "Hello World"
		n, err := writer.Write([]byte(input))

		assert.Error(t, err)
		assert.Equal(t, io.ErrShortWrite, err)
		assert.Equal(t, len(input), n) // Should return the input length
	})

	t.Run("nil writer", func(t *testing.T) {
		writer := NewWriter(nil, false)

		input := "Hello World"
		n, err := writer.Write([]byte(input))

		assert.Error(t, err)
		assert.Equal(t, io.ErrClosedPipe, err)
		assert.Equal(t, len(input), n) // Should return the input length
	})

	t.Run("nil writer with RTL text", func(t *testing.T) {
		writer := NewWriter(nil, true)

		input := "שלום עולם"
		n, err := writer.Write([]byte(input))

		assert.Error(t, err)
		assert.Equal(t, io.ErrClosedPipe, err)
		assert.Equal(t, len(input), n) // Should return the input length
	})

	t.Run("nil writer with markers", func(t *testing.T) {
		writer := NewWriter(nil, false)

		input := "\u200Fשלום\u200E"
		n, err := writer.Write([]byte(input))

		assert.Error(t, err)
		assert.Equal(t, io.ErrClosedPipe, err)
		assert.Equal(t, len(input), n) // Should return the input length
	})
}

// shortWriter is a test helper that only writes a limited number of bytes
type shortWriter struct {
	maxBytes int
	written  int
}

func (sw *shortWriter) Write(p []byte) (int, error) {
	remaining := sw.maxBytes - sw.written
	if remaining <= 0 {
		return 0, io.ErrShortWrite
	}

	toWrite := len(p)
	if toWrite > remaining {
		toWrite = remaining
	}

	sw.written += toWrite
	return toWrite, nil
}

func TestWriter_IntegrationWithString(t *testing.T) {
	tests := []struct {
		name     string
		input    String
		expected string
		flipRTL  bool
	}{
		{
			name:     "String with Hebrew text",
			input:    String("שלום"),
			expected: "םולש",
			flipRTL:  true,
		},
		{
			name:     "String with mixed text",
			input:    String("Hello שלום World"),
			expected: "Hello םולש World",
			flipRTL:  true,
		},
		{
			name:     "String with markers (should not double-reverse)",
			input:    String("\u200Fשלום\u200E"),
			expected: "םולש",
			flipRTL:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewWriter(&buf, tt.flipRTL)

			// Write the String's output
			n, err := writer.Write([]byte(tt.input.String()))
			require.NoError(t, err)
			assert.Equal(t, len(tt.input.String()), n)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

// Benchmark tests for performance
func BenchmarkWriter_Write_WithMarkers(b *testing.B) {
	testCases := []struct {
		name  string
		input string
	}{
		{"LTR only", "Hello World"},
		{"Hebrew with markers", "\u200Fשלום עולם\u200E"},
		{"Arabic with markers", "\u200Fمرحبا بالعالم\u200E"},
		{"Mixed with markers", "Hello \u200Fשלום\u200E World \u200Fمرحبا\u200E Test"},
		{"Long mixed text", "This is a very long text with Hebrew \u200Fשלום עולם\u200E and Arabic \u200Fمرحبا بالعالم\u200E mixed together for performance testing"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			var buf bytes.Buffer
			writer := NewWriter(&buf, false)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				buf.Reset()
				writer.Write([]byte(tc.input))
			}
		})
	}
}

func BenchmarkWriter_Write_WithFlipRTL(b *testing.B) {
	testCases := []struct {
		name  string
		input string
	}{
		{"LTR only", "Hello World"},
		{"Hebrew only", "שלום עולם"},
		{"Arabic only", "مرحبا بالعالم"},
		{"Mixed text", "Hello שלום World مرحبا Test"},
		{"Long mixed text", "This is a very long text with Hebrew שלום עולם and Arabic مرحبا بالعالم mixed together for performance testing"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			var buf bytes.Buffer
			writer := NewWriter(&buf, true)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				buf.Reset()
				writer.Write([]byte(tc.input))
			}
		})
	}
}

func BenchmarkWriter_Write_MultipleWrites(b *testing.B) {
	inputs := []string{"Hello ", "\u200Fשלום\u200E", " World"}

	b.Run("with markers", func(b *testing.B) {
		var buf bytes.Buffer
		writer := NewWriter(&buf, false)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			for _, input := range inputs {
				writer.Write([]byte(input))
			}
		}
	})

	b.Run("with flipRTL", func(b *testing.B) {
		inputs := []string{"Hello ", "שלום", " World"}
		var buf bytes.Buffer
		writer := NewWriter(&buf, true)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			for _, input := range inputs {
				writer.Write([]byte(input))
			}
		}
	})
}
