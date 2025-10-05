package sefaria

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

type RandomTextOptions struct {
	Titles     []string
	Categories []string
}

func (s *TextService) Random(ctx context.Context, opts *RandomTextOptions) (*Text, error) {
	u := s.client.BaseURL.JoinPath("/texts/random")
	if opts != nil {
		// Convert struct to query string
		v, err := query.Values(opts)
		if err != nil {
			return nil, err
		}

		// Sefaria expects "|" separator for slices
		// Override default repeated key style
		q := u.Query()
		for key, vals := range v {
			q.Del(key) // remove any existing
			q.Set(key, strings.Join(vals, "|"))
		}

		u.RawQuery = q.Encode()
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	text := new(Text)
	_, err = s.client.Do(req, text)
	return text, err
}
