// Package bidi provides bidirectional text support for Hebrew and Arabic text.
// It handles automatic text direction detection and proper rendering of mixed
// left-to-right (LTR) and right-to-left (RTL) content.
//
// The package includes:
//   - String: A custom string type that automatically wraps RTL text with
//     Unicode directional markers for proper display
//   - Writer: An io.Writer that processes bidirectional text and reverses
//     RTL sequences for correct rendering
//
// This is particularly useful for applications dealing with Hebrew or Arabic
// text that needs to be displayed correctly in mixed-language contexts.
package bidi
