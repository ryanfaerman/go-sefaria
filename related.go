package sefaria

import (
	"context"
	"net/http"

	"github.com/google/go-querystring/query"
)

type RelatedService service

type RelatedContent map[string]any

func (s *RelatedService) Get(ctx context.Context, tref string) (*RelatedContent, error) {
	u := s.client.BaseURL.JoinPath("/related", tref)
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	out := new(RelatedContent)
	_, err = s.client.Do(req, out)
	return out, err
}

type RelatedLinksOptions struct {
	WithText       bool `url:"with_text,omitempty"`
	WithSheetLinks bool `url:"with_sheet_links,omitempty"`
}

func (s *RelatedService) Links(ctx context.Context, tref string, opts *RelatedLinksOptions) (*RelatedContent, error) {
	u := s.client.BaseURL.JoinPath("/links", tref)

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

	out := new(RelatedContent)
	_, err = s.client.Do(req, out)
	return out, err
}

type Links map[string]any

func (s *RelatedService) TopicLinks(ctx context.Context, tref string) (*Links, error) {
	u := s.client.BaseURL.JoinPath("/links", tref)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	out := new(Links)
	_, err = s.client.Do(req, out)
	return out, err
}
