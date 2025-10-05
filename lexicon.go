package sefaria

import (
	"context"
	"net/http"

	"github.com/google/go-querystring/query"
)

type LexiconService service

type LexiconGetOptions struct {
	LookupRef        string `url:"lookup_ref,omitempty"`
	NeverSplit       bool   `url:"never_split,omitempty"`
	AlwaysSplit      bool   `url:"always_split,omitempty"`
	AlwaysConsonants bool   `url:"always_consonants,omitempty"`
}

type DictionaryEntry map[string]any

func (s *LexiconService) Get(ctx context.Context, word string, opts *LexiconGetOptions) ([]DictionaryEntry, error) {
	u := s.client.BaseURL.JoinPath("/words", word)

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

	out := make([]DictionaryEntry, 0)
	_, err = s.client.Do(req, &out)
	return out, err
}

type LexiconCompletionsOptions struct {
	Limit int `url:"limit,omitempty"`
}

func (s *LexiconService) Completions(ctx context.Context, word string, lexicon string, opts *LexiconCompletionsOptions) ([][]string, error) {
	u := s.client.BaseURL.JoinPath("/words/completion", word, lexicon)

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

	out := make([][]string, 0)
	_, err = s.client.Do(req, &out)
	return out, err
}
