package types

import (
	"encoding/json"
	"time"

	"github.com/araddon/dateparse"
)

// Date represents a date with custom JSON marshaling/unmarshaling behavior.
// It uses dateparse.ParseAny for flexible date parsing and formats output
// as "2006-01-02" (date only, no time component).
//
// This type is designed to handle the various date formats that may be
// returned by the Sefaria API, providing a consistent interface for
// working with dates throughout the client.
//
// Example usage:
//
//	var d Date
//	json.Unmarshal([]byte(`"2024-01-15"`), &d)
//	json.Unmarshal([]byte(`"January 15, 2024"`), &d)
//	json.Unmarshal([]byte(`"15/01/2024"`), &d)
type Date struct {
	time.Time
}

const dateLayout = "2006-01-02"

// UnmarshalJSON implements json.Unmarshaler for Date.
// It handles various date formats by using dateparse.ParseAny,
// which can parse many common date formats automatically.
// Empty strings and "null" values are treated as zero time.
func (d *Date) UnmarshalJSON(data []byte) error {
	// Strip quotes
	s := string(data)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}

	t, err := dateparse.ParseAny(s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// MarshalJSON implements json.Marshaler for Date.
// It formats the date as "2006-01-02" (date only, no time).
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(dateLayout))
}
