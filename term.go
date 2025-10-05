package sefaria

import (
	"context"
	"net/http"

	"github.com/google/go-querystring/query"
)

type TermService service

type Term map[string]any

func (s *TermService) Get(ctx context.Context, term string) (*Term, error) {
	u := s.client.BaseURL.JoinPath("/terms", term)
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	out := new(Term)
	_, err = s.client.Do(req, out)
	return out, err
}

type TermNameOptions struct {
	Limit int `url:"limit,omitempty"`
	Type  int `url:"type,omitempty"`
}

func (s *TermService) Name(ctx context.Context, name string, opts *TermNameOptions) (*Term, error) {
	u := s.client.BaseURL.JoinPath("/name", name)
	if opts != nil {
		v, err := query.Values(opts)
		if err != nil {
			return nil, err
		}
		u.RawQuery = v.Encode()
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	out := new(Term)
	_, err = s.client.Do(req, out)
	return out, err
}
