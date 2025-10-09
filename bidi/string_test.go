package bidi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString_DoubleWrappingPrevention(t *testing.T) {
	tests := []struct {
		name     string
		input    String
		expected string
	}{
		{
			name:     "already wrapped Hebrew",
			input:    String("\u200Fשלום\u200E"),
			expected: "\u200Fשלום\u200E",
		},
		{
			name:     "already wrapped Arabic",
			input:    String("\u200Fمرحبا\u200E"),
			expected: "\u200Fمرحبا\u200E",
		},
		{
			name:     "mixed wrapped and unwrapped",
			input:    String("Hello \u200Fשלום\u200E World مرحبا"),
			expected: "Hello \u200Fשלום\u200E World \u200Fمرحبا\u200E",
		},
		{
			name:     "multiple wrapped sequences",
			input:    String("\u200Fשלום\u200E Hello \u200Fمرحبا\u200E World"),
			expected: "\u200Fשלום\u200E Hello \u200Fمرحبا\u200E World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that multiple calls to String() produce the same result
			firstCall := tt.input.String()
			secondCall := tt.input.String()

			assert.Equal(t, tt.expected, firstCall)
			assert.Equal(t, firstCall, secondCall)
		})
	}
}

func TestString_String(t *testing.T) {
	tests := []struct {
		name     string
		input    String
		expected string
	}{
		{
			name:     "empty string",
			input:    String(""),
			expected: "",
		},
		{
			name:     "LTR only text",
			input:    String("Hello World"),
			expected: "Hello World",
		},
		{
			name:     "Hebrew text only",
			input:    String("שלום"),
			expected: "\u200Fשלום\u200E",
		},
		{
			name:     "Arabic text only",
			input:    String("مرحبا"),
			expected: "\u200Fمرحبا\u200E",
		},
		{
			name:     "Mixed LTR and Hebrew",
			input:    String("Hello שלום World"),
			expected: "Hello \u200Fשלום\u200E World",
		},
		{
			name:     "Mixed LTR and Arabic",
			input:    String("Hello مرحبا World"),
			expected: "Hello \u200Fمرحبا\u200E World",
		},
		{
			name:     "Multiple RTL sequences",
			input:    String("Hello שלום World مرحبا Test"),
			expected: "Hello \u200Fשלום\u200E World \u200Fمرحبا\u200E Test",
		},
		{
			name:     "RTL text with spaces",
			input:    String("שלום עולם"),
			expected: "\u200Fשלום עולם\u200E",
		},
		{
			name:     "RTL text with multiple spaces",
			input:    String("שלום   עולם"),
			expected: "\u200Fשלום   עולם\u200E",
		},
		{
			name:     "RTL text with tabs and newlines",
			input:    String("שלום\tעולם\n"),
			expected: "\u200Fשלום\tעולם\u200E\n",
		},
		{
			name:     "Mixed with RTL spaces",
			input:    String("Hello שלום עולם World"),
			expected: "Hello \u200Fשלום עולם\u200E World",
		},
		{
			name:     "RTL text followed by LTR",
			input:    String("שלום Hello"),
			expected: "\u200Fשלום\u200E Hello",
		},
		{
			name:     "LTR text followed by RTL",
			input:    String("Hello שלום"),
			expected: "Hello \u200Fשלום\u200E",
		},
		{
			name:     "Numbers and punctuation with RTL",
			input:    String("Price: 100₪"),
			expected: "Price: 100₪",
		},
		{
			name:     "Complex mixed text",
			input:    String("The word שלום means 'hello' in Hebrew, and مرحبا means 'hello' in Arabic."),
			expected: "The word \u200Fשלום\u200E means 'hello' in Hebrew, and \u200Fمرحبا\u200E means 'hello' in Arabic.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    String
		expected string
	}{
		{
			name:     "empty string",
			input:    String(""),
			expected: `""`,
		},
		{
			name:     "LTR only text",
			input:    String("Hello World"),
			expected: `"Hello World"`,
		},
		{
			name:     "Hebrew text",
			input:    String("שלום"),
			expected: `"‏שלום‎"`,
		},
		{
			name:     "Arabic text",
			input:    String("مرحبا"),
			expected: `"‏مرحبا‎"`,
		},
		{
			name:     "Mixed text",
			input:    String("Hello שלום World"),
			expected: `"Hello ‏שלום‎ World"`,
		},
		{
			name:     "Text with quotes",
			input:    String(`He said "שלום"`),
			expected: `"He said \"‏שלום‎\""`,
		},
		{
			name:     "Text with backslashes",
			input:    String(`Path: C:\שלום`),
			expected: `"Path: C:\\‏שלום‎"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    String
		expectError bool
	}{
		{
			name:        "empty string",
			input:       `""`,
			expected:    String(""),
			expectError: false,
		},
		{
			name:        "LTR only text",
			input:       `"Hello World"`,
			expected:    String("Hello World"),
			expectError: false,
		},
		{
			name:        "Hebrew text",
			input:       `"שלום"`,
			expected:    String("שלום"),
			expectError: false,
		},
		{
			name:        "Arabic text",
			input:       `"مرحبا"`,
			expected:    String("مرحبا"),
			expectError: false,
		},
		{
			name:        "Mixed text",
			input:       `"Hello שלום World"`,
			expected:    String("Hello שלום World"),
			expectError: false,
		},
		{
			name:        "Text with escaped quotes",
			input:       `"He said \"שלום\""`,
			expected:    String(`He said "שלום"`),
			expectError: false,
		},
		{
			name:        "Text with escaped backslashes",
			input:       `"Path: C:\\שלום"`,
			expected:    String(`Path: C:\שלום`),
			expectError: false,
		},
		{
			name:        "Text with Unicode escapes",
			input:       `"שלום \u05d0\u05d1\u05d2"`,
			expected:    String("שלום אבג"),
			expectError: false,
		},
		{
			name:        "Invalid JSON - missing quotes",
			input:       `Hello World`,
			expected:    String(""),
			expectError: true,
		},
		{
			name:        "Invalid JSON - unclosed string",
			input:       `"Hello World`,
			expected:    String(""),
			expectError: true,
		},
		{
			name:        "Invalid JSON - invalid escape",
			input:       `"Hello \x World"`,
			expected:    String(""),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result String
			err := result.UnmarshalJSON([]byte(tt.input))

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestString_JSONRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input String
	}{
		{
			name:  "empty string",
			input: String(""),
		},
		{
			name:  "LTR only text",
			input: String("Hello World"),
		},
		{
			name:  "Hebrew text",
			input: String("שלום"),
		},
		{
			name:  "Arabic text",
			input: String("مرحبا"),
		},
		{
			name:  "Mixed text",
			input: String("Hello שלום World"),
		},
		{
			name:  "Complex mixed text",
			input: String("The word שלום means 'hello' in Hebrew, and مرحبا means 'hello' in Arabic."),
		},
		{
			name:  "Text with special characters",
			input: String(`He said "שלום" and the path is C:\שלום`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := tt.input.MarshalJSON()
			require.NoError(t, err)

			// Unmarshal back to String
			var result String
			err = result.UnmarshalJSON(jsonData)
			require.NoError(t, err)

			// Now test that multiple calls to String() don't double-wrap
			firstCall := result.String()
			secondCall := result.String()

			// Both calls should produce the same result (no double-wrapping)
			assert.Equal(t, firstCall, secondCall)

			// The result should be the same as the original String() output
			assert.Equal(t, tt.input.String(), result.String())
		})
	}
}

func TestWrapRTL_DoubleWrappingPrevention(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "already wrapped Hebrew",
			input:    "\u200Fשלום\u200E",
			expected: "\u200Fשלום\u200E",
		},
		{
			name:     "already wrapped Arabic",
			input:    "\u200Fمرحبا\u200E",
			expected: "\u200Fمرحبا\u200E",
		},
		{
			name:     "mixed wrapped and unwrapped",
			input:    "Hello \u200Fשלום\u200E World مرحبا",
			expected: "Hello \u200Fשלום\u200E World \u200Fمرحبا\u200E",
		},
		{
			name:     "multiple wrapped sequences",
			input:    "\u200Fשלום\u200E Hello \u200Fمرحبا\u200E World",
			expected: "\u200Fשלום\u200E Hello \u200Fمرحبا\u200E World",
		},
		{
			name:     "wrapped with spaces",
			input:    "\u200Fשלום עולם\u200E",
			expected: "\u200Fשלום עולם\u200E",
		},
		{
			name:     "malformed marker sequence - missing LRM",
			input:    "\u200Fשלום",
			expected: "\u200Fשלום\u200E",
		},
		{
			name:     "malformed marker sequence - missing RLM",
			input:    "שלום\u200E",
			expected: "\u200Fשלום\u200E\u200E",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapRTL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWrapRTL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "LTR only",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "Hebrew only",
			input:    "שלום",
			expected: "\u200Fשלום\u200E",
		},
		{
			name:     "Arabic only",
			input:    "مرحبا",
			expected: "\u200Fمرحبا\u200E",
		},
		{
			name:     "Mixed LTR and Hebrew",
			input:    "Hello שלום World",
			expected: "Hello \u200Fשלום\u200E World",
		},
		{
			name:     "Multiple RTL sequences",
			input:    "Hello שלום World مرحبا Test",
			expected: "Hello \u200Fשלום\u200E World \u200Fمرحبا\u200E Test",
		},
		{
			name:     "RTL with spaces",
			input:    "שלום עולם",
			expected: "\u200Fשלום עולם\u200E",
		},
		{
			name:     "RTL with multiple spaces",
			input:    "שלום   עולם",
			expected: "\u200Fשלום   עולם\u200E",
		},
		{
			name:     "RTL with tabs",
			input:    "שלום\tעולם",
			expected: "\u200Fשלום\tעולם\u200E",
		},
		{
			name:     "RTL with newlines",
			input:    "שלום\nעולם",
			expected: "\u200Fשלום\nעולם\u200E",
		},
		{
			name:     "RTL spaces between LTR",
			input:    "Hello שלום עולם World",
			expected: "Hello \u200Fשלום עולם\u200E World",
		},
		{
			name:     "LTR spaces between RTL sequences",
			input:    "שלום Hello مرحبا",
			expected: "\u200Fשלום\u200E Hello \u200Fمرحبا\u200E",
		},
		{
			name:     "Numbers with RTL currency",
			input:    "Price: 100₪",
			expected: "Price: 100₪",
		},
		{
			name:     "Complex punctuation",
			input:    "שלום! עולם?",
			expected: "\u200Fשלום! עולם\u200E?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapRTL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsRTL(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{
			name:     "Hebrew alef",
			input:    'א',
			expected: true,
		},
		{
			name:     "Hebrew shin",
			input:    'ש',
			expected: true,
		},
		{
			name:     "Arabic alif",
			input:    'ا',
			expected: true,
		},
		{
			name:     "Arabic mim",
			input:    'م',
			expected: true,
		},
		{
			name:     "Latin A",
			input:    'A',
			expected: false,
		},
		{
			name:     "Latin z",
			input:    'z',
			expected: false,
		},
		{
			name:     "Digit",
			input:    '5',
			expected: false,
		},
		{
			name:     "Space",
			input:    ' ',
			expected: false,
		},
		{
			name:     "Punctuation",
			input:    '!',
			expected: false,
		},
		{
			name:     "Hebrew punctuation - geresh",
			input:    '׳',
			expected: true,
		},
		{
			name:     "Arabic punctuation - comma",
			input:    '،',
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRTL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Benchmark tests for performance
func BenchmarkString_String(b *testing.B) {
	testCases := []struct {
		name  string
		input String
	}{
		{"LTR only", String("Hello World")},
		{"Hebrew only", String("שלום עולם")},
		{"Arabic only", String("مرحبا بالعالم")},
		{"Mixed text", String("Hello שלום World مرحبا Test")},
		{"Long mixed text", String("This is a very long text with Hebrew שלום עולם and Arabic مرحبا بالعالم mixed together for performance testing")},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tc.input.String()
			}
		})
	}
}

func BenchmarkWrapRTL(b *testing.B) {
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
			for i := 0; i < b.N; i++ {
				_ = wrapRTL(tc.input)
			}
		})
	}
}
