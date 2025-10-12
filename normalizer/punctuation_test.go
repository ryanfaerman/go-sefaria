package normalizer

import (
	"fmt"
	"testing"
)

func TestPunctuation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		description string
	}{
		{
			name:        "smart quotes",
			input:       string(rune(0x201C)) + "Hello" + string(rune(0x201D)) + " and " + string(rune(0x2018)) + "world" + string(rune(0x2019)),
			expected:    "\"Hello\" and 'world'",
			description: "Should convert smart quotes to straight quotes",
		},
		{
			name:        "smart apostrophes",
			input:       "It" + string(rune(0x2019)) + "s a " + string(rune(0x2018)) + "smart" + string(rune(0x2019)) + " apostrophe",
			expected:    "It's a 'smart' apostrophe",
			description: "Should convert smart apostrophes to straight apostrophes",
		},
		{
			name:        "em and en dashes",
			input:       "Em dash " + string(rune(0x2014)) + " and en dash " + string(rune(0x2013)),
			expected:    "Em dash - and en dash -",
			description: "Should convert em and en dashes to hyphens",
		},
		{
			name:        "ellipsis",
			input:       "Wait for it" + string(rune(0x2026)),
			expected:    "Wait for it...",
			description: "Should convert ellipsis to three dots",
		},
		{
			name:        "non-breaking space",
			input:       "Non-breaking" + string(rune(0x00A0)) + "space",
			expected:    "Non-breaking space",
			description: "Should convert non-breaking space to regular space",
		},
		{
			name:        "mixed punctuation",
			input:       string(rune(0x201C)) + "Hello" + string(rune(0x201D)) + " " + string(rune(0x2014)) + " he said" + string(rune(0x2026)) + " " + string(rune(0x2018)) + "Really?" + string(rune(0x2019)),
			expected:    "\"Hello\" - he said... 'Really?'",
			description: "Should normalize all types of punctuation",
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			description: "Should handle empty string",
		},
		{
			name:        "no fancy punctuation",
			input:       "Plain text without fancy punctuation",
			expected:    "Plain text without fancy punctuation",
			description: "Should return unchanged text when no fancy punctuation present",
		},
		{
			name:        "multiple occurrences",
			input:       string(rune(0x201C)) + "Hello" + string(rune(0x201D)) + " " + string(rune(0x201C)) + "world" + string(rune(0x201D)) + " " + string(rune(0x2014)) + " test" + string(rune(0x2026)) + " and more" + string(rune(0x2026)),
			expected:    "\"Hello\" \"world\" - test... and more...",
			description: "Should handle multiple occurrences of each punctuation type",
		},
		{
			name:        "unicode text with punctuation",
			input:       "Hola " + string(rune(0x201C)) + "mundo" + string(rune(0x201D)) + " " + string(rune(0x2014)) + " ¡Hola!",
			expected:    "Hola \"mundo\" - ¡Hola!",
			description: "Should normalize punctuation in unicode text",
		},
		{
			name:        "quotes within quotes",
			input:       "He said " + string(rune(0x201C)) + "She said " + string(rune(0x2018)) + "Hello" + string(rune(0x2019)) + string(rune(0x201D)),
			expected:    "He said \"She said 'Hello'\"",
			description: "Should handle nested quotes correctly",
		},
		{
			name:        "dashes in different contexts",
			input:       "Range: 1" + string(rune(0x2014)) + "10, Date: 2020" + string(rune(0x2013)) + "2023",
			expected:    "Range: 1-10, Date: 2020-2023",
			description: "Should convert dashes in various contexts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Punctuation(tt.input)
			if result != tt.expected {
				t.Errorf("Punctuation(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPunctuation_EdgeCases(t *testing.T) {
	t.Run("very long string", func(t *testing.T) {
		// Create a long string with many punctuation marks
		longInput := ""
		expected := ""
		for i := 0; i < 1000; i++ {
			longInput += string(rune(0x2026))
			expected += "..."
		}

		result := Punctuation(longInput)
		if result != expected {
			t.Errorf("Punctuation() failed on long string")
		}
	})

	t.Run("mixed unicode and punctuation", func(t *testing.T) {
		input := "中文" + string(rune(0x201C)) + "标点" + string(rune(0x201D)) + "符号" + string(rune(0x2026))
		expected := "中文\"标点\"符号..."

		result := Punctuation(input)
		if result != expected {
			t.Errorf("Punctuation(%q) = %q, want %q", input, result, expected)
		}
	})

	t.Run("newlines and tabs with punctuation", func(t *testing.T) {
		input := "Line 1 " + string(rune(0x201C)) + "quotes" + string(rune(0x201D)) + "\nLine 2 " + string(rune(0x2014)) + " dash\tLine 3" + string(rune(0x2026))
		expected := "Line 1 \"quotes\"\nLine 2 - dash\tLine 3..."

		result := Punctuation(input)
		if result != expected {
			t.Errorf("Punctuation(%q) = %q, want %q", input, result, expected)
		}
	})

	t.Run("only punctuation characters", func(t *testing.T) {
		input := string(rune(0x201C)) + string(rune(0x201D)) + string(rune(0x2018)) + string(rune(0x2019)) + string(rune(0x2014)) + string(rune(0x2013)) + string(rune(0x2026)) + string(rune(0x00A0))
		expected := "\"\"''--... "

		result := Punctuation(input)
		if result != expected {
			t.Errorf("Punctuation(%q) = %q, want %q", input, result, expected)
		}
	})
}

func TestPunctuation_Performance(t *testing.T) {
	// Test that punctuation normalization is reasonably fast
	input := "This is a " + string(rune(0x201C)) + "test" + string(rune(0x201D)) + " string with " + string(rune(0x2014)) + " various" + string(rune(0x2026)) + " punctuation marks " + string(rune(0x2018)) + "and" + string(rune(0x2019)) + " quotes"

	// Run multiple times to ensure consistent performance
	for i := 0; i < 1000; i++ {
		result := Punctuation(input)
		if len(result) == 0 {
			t.Error("Punctuation returned empty result")
		}
	}
}

// ExamplePunctuation demonstrates basic usage of Punctuation
func ExamplePunctuation() {
	text := string(rune(0x201C)) + "Hello" + string(rune(0x201D)) + " " + string(rune(0x2014)) + " he said" + string(rune(0x2026)) + " " + string(rune(0x2018)) + "Really?" + string(rune(0x2019))
	normalized := Punctuation(text)
	fmt.Println(normalized)
	// Output: "Hello" - he said... 'Really?'
}

// ExamplePunctuation_mixed demonstrates handling various punctuation types
func ExamplePunctuation_mixed() {
	text := "Smart " + string(rune(0x201C)) + "quotes" + string(rune(0x201D)) + " and " + string(rune(0x2018)) + "apostrophes" + string(rune(0x2019)) + " " + string(rune(0x2014)) + " with dashes" + string(rune(0x2026)) + " and spaces" + string(rune(0x00A0)) + "here"
	normalized := Punctuation(text)
	fmt.Println(normalized)
	// Output: Smart "quotes" and 'apostrophes' - with dashes... and spaces here
}
