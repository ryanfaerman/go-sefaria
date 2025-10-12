// Package types provides custom types with specialized marshaling and unmarshaling behavior
// for the Sefaria API client. These types handle various encoding concerns including:
//
//   - Date formatting and parsing with flexible input handling
//   - Boolean-to-integer conversion for URL parameters
//   - Generic string-or-type parsing for API responses
//
// The types in this package are designed to handle the quirks and inconsistencies
// found in the Sefaria API responses and parameter encoding requirements.
package types
