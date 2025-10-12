// Package normalizer provides text normalization utilities for cleaning and standardizing text data.
//
// The normalizer package is designed to work with structured data, automatically applying
// normalization functions to string fields within complex data structures like structs,
// slices, arrays, and maps. This is particularly useful when processing data from external
// sources that may contain inconsistent formatting, encoding, or special characters.
//
// # Core Concepts
//
// A Normalizer is a function that takes a string and returns a normalized string.
// The package provides several built-in normalizers and a mechanism to apply them
// recursively to complex data structures.
//
// # Built-in Normalizers
//
// The package includes several common text normalization functions:
//
//   - HTMLUnescape: Converts HTML entities to their corresponding characters
//   - UnicodeNFC: Normalizes Unicode text to NFC (Canonical Decomposed, then Canonical Composed) form
//   - Punctuation: Replaces fancy/smart punctuation with standard ASCII equivalents
//
// # Usage Example
//
// The most common usage pattern is to define a set of normalizers and apply them
// to structured data:
//
//	normalizers := []normalizer.Normalizer{
//		normalizer.HTMLUnescape,
//		normalizer.UnicodeNFC,
//		normalizer.Punctuation,
//	}
//
//	type Document struct {
//		Title   string
//		Content string
//		Tags    []string
//	}
//
//	doc := &Document{
//		Title:   "Hello &amp; Welcome",
//		Content: "This is "smart" text with fancy punctuation…",
//		Tags:    []string{"tag1", "tag2"},
//	}
//
//	normalizer.Apply(doc, normalizers...)
//	// doc.Title is now "Hello & Welcome"
//	// doc.Content is now "This is "smart" text with fancy punctuation..."
//	// doc.Tags remain unchanged
//
// # Custom Normalizers
//
// You can create custom normalizers by implementing the Normalizer function type:
//
//	func TrimWhitespace(s string) string {
//		return strings.TrimSpace(s)
//	}
//
//	normalizers := []normalizer.Normalizer{
//		TrimWhitespace,
//		normalizer.HTMLUnescape,
//	}
//
// # Supported Data Types
//
// The Apply function works with the following data types:
//
//   - Strings: Direct normalization
//   - Structs: Recursively applies to all string fields
//   - Pointers: Dereferences and applies to the underlying value
//   - Interfaces: Applies to the underlying concrete value
//   - Slices/Arrays: Applies to each element
//   - Maps: Applies to string values and recursively to other values
//
// # Performance Considerations
//
// The Apply function uses reflection to traverse data structures, which has some
// performance overhead. For high-performance scenarios, consider applying
// normalizers directly to known string fields rather than using the generic Apply function.
//
// # Thread Safety
//
// All normalizer functions are stateless and thread-safe. The Apply function
// modifies data in-place and should not be called concurrently on the same
// data structure.
package normalizer
