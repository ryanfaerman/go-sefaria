package sefaria

import (
	"context"
	"net/http"

	"github.com/google/go-querystring/query"
)

type CalendarService service

type CalendarGetOptions struct {
	Diaspora bool   `url:"diaspora,omitempty"`
	Custom   string `url:"custom,omitempty"`
	Year     int    `url:"year,omitempty"`
	Month    int    `url:"month,omitempty"`
	Day      int    `url:"day,omitempty"`
	TimeZone string `url:"timezone,omitempty"`
}

type Calendar map[string]any

func (s *CalendarService) Get(ctx context.Context, opts *CalendarGetOptions) (*Calendar, error) {
	u := s.client.BaseURL.JoinPath("/calendars")
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

	calendar := new(Calendar)
	_, err = s.client.Do(req, calendar)
	return calendar, err
}

type Parsha map[string]any

func (s *CalendarService) NextRead(ctx context.Context, parsha string) (*Parsha, error) {
	u := s.client.BaseURL.JoinPath("/calendars/next-read", parsha)
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	out := new(Parsha)
	_, err = s.client.Do(req, out)
	return out, err
}
