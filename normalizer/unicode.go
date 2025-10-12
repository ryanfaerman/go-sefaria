package normalizer

import "golang.org/x/text/unicode/norm"

// UnicodeNFC normalizes Unicode text to NFC (Canonical Decomposed, then Canonical Composed) form.
// This function uses the golang.org/x/text/unicode/norm package to ensure consistent
// Unicode representation by decomposing characters and then recomposing them in
// canonical order.
//
// NFC normalization is important for:
//   - Consistent string comparison and searching
//   - Preventing duplicate entries that differ only in Unicode representation
//   - Ensuring compatibility across different systems and platforms
//
// Example:
//
//	UnicodeNFC("café") // Returns properly normalized "café"
//	UnicodeNFC("cafe\u0301") // Also returns "café" (normalized form)
func UnicodeNFC(s string) string {
	return norm.NFC.String(s)
}
