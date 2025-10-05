package sefaria

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type TopicService service

type Topic map[string]any

func (s *TopicService) All(ctx context.Context, limit int) ([]Topic, error) {
	u := s.client.BaseURL.JoinPath("/topics")

	q := u.Query()
	q.Set("limit", strconv.Itoa(limit))
	u.RawQuery = q.Encode()

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	out := make([]Topic, 0)
	_, err = s.client.Do(req, &out)
	return out, err
}

func (s *TopicService) Get(ctx context.Context, topic string) (*Topic, error) {
	u := s.client.BaseURL.JoinPath("/v2/topics/", topic)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	out := new(Topic)
	_, err = s.client.Do(req, &out)
	return out, err
}

func (s *TopicService) Graph(ctx context.Context, topic string, linkType string) ([]Topic, error) {
	u := s.client.BaseURL.JoinPath("/topics-graph/", topic)

	q := u.Query()
	if linkType == "" {
		linkType = "is-a"
	}
	q.Set("link_type", linkType)
	u.RawQuery = q.Encode()

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	out := make([]Topic, 0)
	_, err = s.client.Do(req, &out)
	return out, err
}

func (s *TopicService) Recommended(ctx context.Context, refs ...string) ([]Topic, error) {
	if len(refs) == 0 {
		return nil, fmt.Errorf("at least one ref is required")
	}
	refList := strings.Join(refs, "+")
	u := s.client.BaseURL.JoinPath("/recommend/topics", refList)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	out := make([]Topic, 0)
	_, err = s.client.Do(req, &out)
	return out, err
}

func (s *TopicService) Random(ctx context.Context) ([]Topic, error) {
	u := s.client.BaseURL.JoinPath("/texts/random-by-topic")

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	out := make([]Topic, 0)
	_, err = s.client.Do(req, &out)
	return out, err
}
