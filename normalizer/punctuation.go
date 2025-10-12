package normalizer

import "strings"

// Punctuation replaces fancy/smart punctuation characters with standard ASCII equivalents.
// This normalizer converts various Unicode punctuation marks to their basic ASCII counterparts:
//
//   - Smart quotes ("") → straight quotes ("")
//   - Smart apostrophes (”) → straight apostrophes (”)
//   - Em/en dashes (—–) → hyphens (-)
//   - Ellipsis (…) → three dots (...)
//   - Non-breaking space (\u00A0) → regular space
//
// This is useful for standardizing text that may contain fancy punctuation from
// word processors or web content.
//
// Example:
//
//	Punctuation(""Hello" — he said…") // Returns ""Hello" - he said..."
func Punctuation(s string) string {
	replacer := strings.NewReplacer(
		"“", `"`, "”", `"`,
		"‘", `'`, "’", `'`,
		"—", "-", "–", "-",
		"…", "...",
		"\u00A0", " ",
	)
	return replacer.Replace(s)
}
