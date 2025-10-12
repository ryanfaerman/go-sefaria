package normalizer

import (
	"reflect"
	"testing"
)

func TestApply(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		normalizers []Normalizer
		expected    any
		description string
	}{
		{
			name:        "nil input",
			input:       nil,
			normalizers: []Normalizer{HTMLUnescape},
			expected:    nil,
			description: "Should handle nil input gracefully",
		},
		{
			name:        "simple string",
			input:       stringPtr("Hello &amp; World"),
			normalizers: []Normalizer{HTMLUnescape},
			expected:    stringPtr("Hello & World"),
			description: "Should normalize simple string",
		},
		{
			name:        "string pointer",
			input:       stringPtr("Test &lt; 5"),
			normalizers: []Normalizer{HTMLUnescape},
			expected:    stringPtr("Test < 5"),
			description: "Should normalize string pointer",
		},
		{
			name:        "empty string",
			input:       stringPtr(""),
			normalizers: []Normalizer{HTMLUnescape, Punctuation},
			expected:    stringPtr(""),
			description: "Should handle empty string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Apply(tt.input, tt.normalizers...)

			if tt.input == nil && tt.expected == nil {
				return
			}

			if !reflect.DeepEqual(tt.input, tt.expected) {
				t.Errorf("Apply() = %v, want %v", tt.input, tt.expected)
			}
		})
	}
}

func TestApply_Struct(t *testing.T) {
	type Person struct {
		Name    string
		Bio     string
		Age     int
		Tags    []string
		Details map[string]string
	}

	tests := []struct {
		name        string
		input       *Person
		normalizers []Normalizer
		expected    *Person
		description string
	}{
		{
			name: "struct with HTML entities and punctuation",
			input: &Person{
				Name: "John &amp; Jane",
				Bio:  "Hello " + string(rune(0x201C)) + "world" + string(rune(0x201D)) + " — nice to meet you" + string(rune(0x2026)),
				Age:  30,
				Tags: []string{"tag1", "tag2"},
				Details: map[string]string{
					"city":  "New York &amp; Co.",
					"quote": string(rune(0x201C)) + "Smart" + string(rune(0x201D)) + " quotes here",
				},
			},
			normalizers: []Normalizer{HTMLUnescape, Punctuation},
			expected: &Person{
				Name: "John & Jane",
				Bio:  "Hello \"world\" - nice to meet you...",
				Age:  30,
				Tags: []string{"tag1", "tag2"},
				Details: map[string]string{
					"city":  "New York & Co.",
					"quote": "\"Smart\" quotes here",
				},
			},
			description: "Should normalize all string fields in struct",
		},
		{
			name: "struct with unicode normalization",
			input: &Person{
				Name: "José",
				Bio:  "cafe" + string(rune(0x0301)), // decomposed form
			},
			normalizers: []Normalizer{UnicodeNFC},
			expected: &Person{
				Name: "José",
				Bio:  "café", // composed form
			},
			description: "Should normalize unicode in struct fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Apply(tt.input, tt.normalizers...)

			if !reflect.DeepEqual(tt.input, tt.expected) {
				t.Errorf("Apply() = %+v, want %+v", tt.input, tt.expected)
			}
		})
	}
}

func TestApply_Slice(t *testing.T) {
	tests := []struct {
		name        string
		input       []string
		normalizers []Normalizer
		expected    []string
		description string
	}{
		{
			name:        "string slice",
			input:       []string{"Hello &amp;", "World &lt; 5", ""},
			normalizers: []Normalizer{HTMLUnescape},
			expected:    []string{"Hello &", "World < 5", ""},
			description: "Should normalize all strings in slice",
		},
		{
			name:        "empty slice",
			input:       []string{},
			normalizers: []Normalizer{HTMLUnescape},
			expected:    []string{},
			description: "Should handle empty slice",
		},
		{
			name:        "nil slice",
			input:       []string(nil),
			normalizers: []Normalizer{HTMLUnescape},
			expected:    []string(nil),
			description: "Should handle nil slice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Apply(&tt.input, tt.normalizers...)

			if !reflect.DeepEqual(tt.input, tt.expected) {
				t.Errorf("Apply() = %v, want %v", tt.input, tt.expected)
			}
		})
	}
}

func TestApply_Map(t *testing.T) {
	tests := []struct {
		name        string
		input       map[string]string
		normalizers []Normalizer
		expected    map[string]string
		description string
	}{
		{
			name: "string map",
			input: map[string]string{
				"title": "Article &amp; News",
				"desc":  "This is " + string(rune(0x201C)) + "smart" + string(rune(0x201D)) + " text" + string(rune(0x2026)),
				"empty": "",
			},
			normalizers: []Normalizer{HTMLUnescape, Punctuation},
			expected: map[string]string{
				"title": "Article & News",
				"desc":  "This is \"smart\" text...",
				"empty": "",
			},
			description: "Should normalize all string values in map",
		},
		{
			name:        "empty map",
			input:       map[string]string{},
			normalizers: []Normalizer{HTMLUnescape},
			expected:    map[string]string{},
			description: "Should handle empty map",
		},
		{
			name:        "nil map",
			input:       map[string]string(nil),
			normalizers: []Normalizer{HTMLUnescape},
			expected:    map[string]string(nil),
			description: "Should handle nil map",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Apply(&tt.input, tt.normalizers...)

			if !reflect.DeepEqual(tt.input, tt.expected) {
				t.Errorf("Apply() = %v, want %v", tt.input, tt.expected)
			}
		})
	}
}

func TestApply_Interface(t *testing.T) {
	var str interface{} = "Hello &amp; World"
	Apply(&str, HTMLUnescape)

	expected := "Hello & World"
	if str != expected {
		t.Errorf("Apply() = %v, want %v", str, expected)
	}
}

func TestApply_Array(t *testing.T) {
	arr := [3]string{"Hello &amp;", "World &lt; 5", ""}
	Apply(&arr, HTMLUnescape)

	expected := [3]string{"Hello &", "World < 5", ""}
	if arr != expected {
		t.Errorf("Apply() = %v, want %v", arr, expected)
	}
}

func TestApply_MultipleNormalizers(t *testing.T) {
	input := stringPtr("Hello &amp; " + string(rune(0x201C)) + "world" + string(rune(0x201D)) + " — test" + string(rune(0x2026)))
	Apply(input, HTMLUnescape, Punctuation, UnicodeNFC)

	expected := "Hello & \"world\" - test..."
	if *input != expected {
		t.Errorf("Apply() = %v, want %v", *input, expected)
	}
}

func TestApply_NoNormalizers(t *testing.T) {
	input := stringPtr("Hello &amp; World")
	original := *input

	Apply(input)

	if *input != original {
		t.Errorf("Apply() with no normalizers should not change input, got %v", *input)
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
