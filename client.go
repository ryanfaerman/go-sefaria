package sefaria

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/hashicorp/go-retryablehttp"
)

type Client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *retryablehttp.Client
	clientMu   sync.Mutex

	common service

	Text     *TextService
	Index    *IndexService
	Calendar *CalendarService
	Related  *RelatedService
	Lexicon  *LexiconService
	Topics   *TopicService
	Terms    *TermService
}

func NewClient() *Client {
	u, err := url.Parse("https://www.sefaria.org/api")
	if err != nil {
		panic(err)
	}
	c := &Client{
		BaseURL:    u,
		httpClient: retryablehttp.NewClient(),
		UserAgent:  "go-sefaria/v1",
	}

	c.common.client = c
	c.Text = (*TextService)(&c.common)
	c.Index = (*IndexService)(&c.common)
	c.Calendar = (*CalendarService)(&c.common)
	c.Related = (*RelatedService)(&c.common)
	c.Lexicon = (*LexiconService)(&c.common)
	c.Topics = (*TopicService)(&c.common)
	c.Terms = (*TermService)(&c.common)

	return c
}

type RequestOption func(req *http.Request)

func (c *Client) NewRequest(ctx context.Context, method string, u *url.URL, body any, opts ...RequestOption) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v any) (*http.Response, error) {
	rreq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(rreq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// TODO: handle http errors here

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, res.Body)
	default:
		decErr := json.NewDecoder(res.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors from empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}

	return res, err
}

type service struct {
	client *Client
}
