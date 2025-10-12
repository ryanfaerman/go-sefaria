package normalizer

import "html"

// HTMLUnescape converts HTML entities to their corresponding characters.
// This function uses the standard library's html.UnescapeString to decode
// HTML entities like &amp;, &lt;, &gt;, &quot;, etc.
//
// Example:
//
//	HTMLUnescape("Hello &amp; Welcome") // Returns "Hello & Welcome"
//	HTMLUnescape("Price &lt; $100")   // Returns "Price < $100"
func HTMLUnescape(s string) string {
	return html.UnescapeString(s)
}
