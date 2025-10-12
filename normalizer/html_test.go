package normalizer

import (
	"fmt"
	"testing"
)

func TestHTMLUnescape(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		description string
	}{
		{
			name:        "basic HTML entities",
			input:       "Hello &amp; Welcome",
			expected:    "Hello & Welcome",
			description: "Should decode &amp; to &",
		},
		{
			name:        "less than and greater than",
			input:       "Price &lt; $100 &gt; $50",
			expected:    "Price < $100 > $50",
			description: "Should decode &lt; and &gt;",
		},
		{
			name:        "quotes",
			input:       "He said &quot;Hello&quot;",
			expected:    "He said \"Hello\"",
			description: "Should decode &quot; to double quotes",
		},
		{
			name:        "apostrophe",
			input:       "Don&apos;t do that",
			expected:    "Don't do that",
			description: "Should decode &apos; to apostrophe",
		},
		{
			name:        "multiple entities",
			input:       "&amp;&lt;&gt;&quot;&apos;",
			expected:    "&<>\"'",
			description: "Should decode multiple HTML entities",
		},
		{
			name:        "numeric entities",
			input:       "Copyright &#169; 2023",
			expected:    "Copyright © 2023",
			description: "Should decode numeric HTML entities",
		},
		{
			name:        "hex entities",
			input:       "Euro symbol &#x20AC;",
			expected:    "Euro symbol €",
			description: "Should decode hexadecimal HTML entities",
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			description: "Should handle empty string",
		},
		{
			name:        "no entities",
			input:       "Plain text without entities",
			expected:    "Plain text without entities",
			description: "Should return unchanged text when no entities present",
		},
		{
			name:        "mixed content",
			input:       "Text with &amp; entities &lt; and &gt; symbols",
			expected:    "Text with & entities < and > symbols",
			description: "Should decode entities in mixed content",
		},
		{
			name:        "invalid entity",
			input:       "Invalid &invalid; entity",
			expected:    "Invalid &invalid; entity",
			description: "Should leave invalid entities unchanged",
		},
		{
			name:        "incomplete entity",
			input:       "Incomplete &amp",
			expected:    "Incomplete &",
			description: "Should leave incomplete entities unchanged",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HTMLUnescape(tt.input)
			if result != tt.expected {
				t.Errorf("HTMLUnescape(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHTMLUnescape_EdgeCases(t *testing.T) {
	t.Run("very long string", func(t *testing.T) {
		// Create a long string with many entities
		longInput := ""
		expected := ""
		for i := 0; i < 1000; i++ {
			longInput += "&amp;"
			expected += "&"
		}

		result := HTMLUnescape(longInput)
		if result != expected {
			t.Errorf("HTMLUnescape() failed on long string")
		}
	})

	t.Run("unicode characters", func(t *testing.T) {
		input := "Hello &amp; 世界"
		expected := "Hello & 世界"

		result := HTMLUnescape(input)
		if result != expected {
			t.Errorf("HTMLUnescape(%q) = %q, want %q", input, result, expected)
		}
	})

	t.Run("newlines and special characters", func(t *testing.T) {
		input := "Line 1 &amp; Line 2\nTab\t&amp; Space"
		expected := "Line 1 & Line 2\nTab\t& Space"

		result := HTMLUnescape(input)
		if result != expected {
			t.Errorf("HTMLUnescape(%q) = %q, want %q", input, result, expected)
		}
	})
}

// ExampleHTMLUnescape demonstrates basic usage of HTMLUnescape
func ExampleHTMLUnescape() {
	text := "Hello &amp; Welcome to our &lt;website&gt;"
	normalized := HTMLUnescape(text)
	fmt.Println(normalized)
	// Output: Hello & Welcome to our <website>
}

// ExampleHTMLUnescape_advanced demonstrates handling multiple entity types
func ExampleHTMLUnescape_advanced() {
	text := "He said &quot;Hello &amp; goodbye&quot; &lt; 5 minutes"
	normalized := HTMLUnescape(text)
	fmt.Println(normalized)
	// Output: He said "Hello & goodbye" < 5 minutes
}
