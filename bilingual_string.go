package sefaria

import "github.com/ryanfaerman/go-sefaria/bidi"

type BilingualString struct {
	English string      `json:"en,omitempty"`
	Hebrew  bidi.String `json:"he,omitempty"`
}

func (bs BilingualString) String() string {
	return bs.English + " / " + bs.Hebrew.String()
}
