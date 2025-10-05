package sefaria

import (
	"context"
	"net/http"
)

type Manuscript struct {
	ManuscriptSlug    string   `json:"manuscript_slug"`
	PageID            string   `json:"page_id"`
	ImageURL          string   `json:"image_url"`
	ThumbnailURL      string   `json:"thumbnail_url"`
	AnchorRef         string   `json:"anchorRef"`
	AnchorRefExpanded []string `json:"anchorRefExpanded"`
	Manuscript        struct {
		Slug          string `json:"slug"`
		Title         string `json:"title"`
		HeTitle       string `json:"he_title"`
		Source        string `json:"source"`
		Description   string `json:"description"`
		HeDescription string `json:"he_description"`
	} `json:"manuscript"`
}

func (s *TextService) Manuscripts(ctx context.Context, tref string) ([]Manuscript, error) {
	u := s.client.BaseURL.JoinPath("/manuscripts/", tref)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	manuscripts := make([]Manuscript, 0)
	_, err = s.client.Do(req, &manuscripts)
	return manuscripts, err
}
