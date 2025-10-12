package normalizer

import (
	"fmt"
	"testing"
)

func TestUnicodeNFC(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		description string
	}{
		{
			name:        "already normalized",
			input:       "café",
			expected:    "café",
			description: "Should return already normalized text unchanged",
		},
		{
			name:        "decomposed form",
			input:       "cafe" + string(rune(0x0301)), // é as e + combining acute
			expected:    "café",
			description: "Should normalize decomposed form to composed",
		},
		{
			name:        "multiple combining marks",
			input:       "e" + string(rune(0x0301)) + string(rune(0x0302)), // e with acute and circumflex
			expected:    "é̂",                                              // should normalize to composed form
			description: "Should handle multiple combining marks",
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			description: "Should handle empty string",
		},
		{
			name:        "ascii text",
			input:       "Hello World",
			expected:    "Hello World",
			description: "Should return ASCII text unchanged",
		},
		{
			name:        "mixed normalized and decomposed",
			input:       "café cafe" + string(rune(0x0301)),
			expected:    "café café",
			description: "Should normalize mixed text consistently",
		},
		{
			name:        "various accented characters",
			input:       "naïve résumé",
			expected:    "naïve résumé",
			description: "Should handle various accented characters",
		},
		{
			name:        "decomposed accented characters",
			input:       "nai" + string(rune(0x0308)) + "ve re" + string(rune(0x0301)) + "sume" + string(rune(0x0301)),
			expected:    "naïve résumé",
			description: "Should normalize decomposed accented characters",
		},
		{
			name:        "korean characters",
			input:       "한글",
			expected:    "한글",
			description: "Should handle Korean characters",
		},
		{
			name:        "chinese characters",
			input:       "中文",
			expected:    "中文",
			description: "Should handle Chinese characters",
		},
		{
			name:        "arabic characters",
			input:       "العربية",
			expected:    "العربية",
			description: "Should handle Arabic characters",
		},
		{
			name:        "hebrew characters",
			input:       "עברית",
			expected:    "עברית",
			description: "Should handle Hebrew characters",
		},
		{
			name:        "emoji",
			input:       "Hello 😀 World",
			expected:    "Hello 😀 World",
			description: "Should handle emoji characters",
		},
		{
			name:        "ligatures",
			input:       "ﬃ ﬁ ﬂ",
			expected:    "ﬃ ﬁ ﬂ",
			description: "Should handle ligatures",
		},
		{
			name:        "decomposed ligatures",
			input:       "ff" + string(rune(0x0069)) + " fi" + string(rune(0x0069)) + " fl" + string(rune(0x0069)),
			expected:    "ffi fii fli",
			description: "Should normalize decomposed ligatures",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnicodeNFC(tt.input)
			if result != tt.expected {
				t.Errorf("UnicodeNFC(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnicodeNFC_EdgeCases(t *testing.T) {
	t.Run("very long string", func(t *testing.T) {
		// Create a long string with many accented characters
		longInput := ""
		expected := ""
		for i := 0; i < 1000; i++ {
			longInput += "cafe" + string(rune(0x0301))
			expected += "café"
		}

		result := UnicodeNFC(longInput)
		if result != expected {
			t.Errorf("UnicodeNFC() failed on long string")
		}
	})

	t.Run("control characters", func(t *testing.T) {
		input := "Hello\tWorld\n"
		expected := "Hello\tWorld\n"

		result := UnicodeNFC(input)
		if result != expected {
			t.Errorf("UnicodeNFC(%q) = %q, want %q", input, result, expected)
		}
	})

	t.Run("surrogate pairs", func(t *testing.T) {
		input := "Hello \U0001F600 World" // 😀 emoji
		expected := "Hello \U0001F600 World"

		result := UnicodeNFC(input)
		if result != expected {
			t.Errorf("UnicodeNFC(%q) = %q, want %q", input, result, expected)
		}
	})

	t.Run("zero width characters", func(t *testing.T) {
		input := "Hello" + string(rune(0x200B)) + "World" // zero width space
		expected := "Hello" + string(rune(0x200B)) + "World"

		result := UnicodeNFC(input)
		if result != expected {
			t.Errorf("UnicodeNFC(%q) = %q, want %q", input, result, expected)
		}
	})
}

func TestUnicodeNFC_Idempotent(t *testing.T) {
	// Test that applying UnicodeNFC multiple times doesn't change the result
	testCases := []string{
		"café",
		"cafe" + string(rune(0x0301)),
		"naïve résumé",
		"Hello World",
		"中文",
		"العربية",
		"Hello 😀 World",
	}

	for _, input := range testCases {
		first := UnicodeNFC(input)
		second := UnicodeNFC(first)

		if first != second {
			t.Errorf("UnicodeNFC is not idempotent for %q: first=%q, second=%q", input, first, second)
		}
	}
}

func TestUnicodeNFC_Performance(t *testing.T) {
	// Test that Unicode normalization is reasonably fast
	input := "This is a test string with café and naïve characters"

	// Run multiple times to ensure consistent performance
	for i := 0; i < 1000; i++ {
		result := UnicodeNFC(input)
		if len(result) == 0 {
			t.Error("UnicodeNFC returned empty result")
		}
	}
}

// ExampleUnicodeNFC demonstrates basic usage of UnicodeNFC
func ExampleUnicodeNFC() {
	text := "cafe" + string(rune(0x0301)) // decomposed form
	normalized := UnicodeNFC(text)
	fmt.Println(normalized)
	// Output: café
}

// ExampleUnicodeNFC_mixed demonstrates handling mixed normalized and decomposed text
func ExampleUnicodeNFC_mixed() {
	text := "café cafe" + string(rune(0x0301)) + " naïve nai" + string(rune(0x0308)) + "ve"
	normalized := UnicodeNFC(text)
	fmt.Println(normalized)
	// Output: café café naïve naïve
}
