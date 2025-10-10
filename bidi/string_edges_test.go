package bidi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString_Edgecases(t *testing.T) {
	examples := map[string]struct {
		input    string
		expected string
	}{
		"punct": {
			input:    "בראשית ל״ה (35)",
			expected: "\u200fבראשית ל״ה\u200e (35)",
		},
		"intra-quote": {
			input:    "הרמב\"ם היומי",
			expected: "\u200fהרמב\"ם היומי\u200e",
		},
		"name with quote": {
			input:    "יחזקאל ל״ח:י״ח-ל״ט:ט״ז",
			expected: "\u200fיחזקאל ל״ח:י״ח-ל״ט:ט״ז\u200e",
		},
		"name-number": {
			input:    "בראשית ל׳ (30)",
			expected: "\u200fבראשית ל׳\u200e (30)",
		},
	}

	for name, ex := range examples {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, ex.expected, String(ex.input).String(), name)
		})
	}
}
