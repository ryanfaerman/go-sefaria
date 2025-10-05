package sefaria

import (
	"context"
	"net/http"
)

type Translation struct {
	Category     string `json:"category"`
	Name         string `json:"name"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	VersionTitle string `json:"versionTitle"`
	RTLLanguage  string `json:"rtlLanguage"`
}

func (s *TextService) Translations(ctx context.Context, lang string) ([]Translation, error) {
	u := s.client.BaseURL.JoinPath("/texts/translations", lang)
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	results := make(map[string]map[string][]Translation)

	_, err = s.client.Do(req, &results)
	if err != nil {
		return nil, err
	}

	output := make([]Translation, 0)

	for categoryName, categoryContents := range results {
		for name, translations := range categoryContents {
			for _, t := range translations {
				t.Category = categoryName
				t.Name = name
				output = append(output, t)
			}
		}
	}

	return output, nil
}
