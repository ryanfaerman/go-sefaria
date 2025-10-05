package sefaria

import (
	"context"
	"net/http"
)

type IndexService service

type Index map[string]any

func (s *IndexService) Contents(ctx context.Context) (*Index, error) {
	u := s.client.BaseURL.JoinPath("index")
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	index := new(Index)
	_, err = s.client.Do(req, index)
	return index, err
}

func (s *IndexService) Get(ctx context.Context, title string) (*Index, error) {
	u := s.client.BaseURL.JoinPath("/v2/raw/index", title)
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	index := new(Index)
	_, err = s.client.Do(req, index)
	return index, err
}

type IndexShapeOptions struct {
	Depth      int  `url:"depth,omitempty"`
	Dependents bool `url:"dependents,omitempty"`
}

type Shape struct {
	Section   string `json:"section"`
	IsComplex bool   `json:"isComplex"`
	Length    int    `json:"length"`
	Book      string `json:"book"`
	HeBook    string `json:"heBook"`
	Chapters  []int  `json:"chapters"`
}

func (s *IndexService) Shape(ctx context.Context, title string, opts *IndexShapeOptions) ([]Shape, error) {
	u := s.client.BaseURL.JoinPath("/shape", title)
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	shapes := make([]Shape, 0)
	_, err = s.client.Do(req, &shapes)
	return shapes, err
}
