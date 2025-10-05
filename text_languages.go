package sefaria

import (
	"context"
	"net/http"
)

func (s *TextService) Languages(ctx context.Context) ([]string, error) {
	u := s.client.BaseURL.JoinPath("/texts/translations")

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	langs := make([]string, 0)
	_, err = s.client.Do(req, &langs)
	return langs, err
}
