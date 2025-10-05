package sefaria

import (
	"context"
	"net/http"
)

func (s *TextService) Versions(ctx context.Context, index string) ([]Version, error) {
	u := s.client.BaseURL.JoinPath("/texts/versions/", index)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	versions := make([]Version, 0)
	_, err = s.client.Do(req, &versions)
	return versions, err
}
