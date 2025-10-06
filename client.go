package sefaria

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-retryablehttp"
)

type Client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *retryablehttp.Client
	clientMu   sync.Mutex

	validate *validator.Validate

	logger *slog.Logger

	common service

	Text     *TextService
	Index    *IndexService
	Related  *RelatedService
	Calendar *CalendarService
	Lexicon  *LexiconService
	Topics   *TopicService
	Terms    *TermService
}

type ClientOption func(*Client)

func NewClient(opts ...ClientOption) *Client {
	u, err := url.Parse("https://www.sefaria.org/api")
	if err != nil {
		// This should _never_ happen.
		panic(err)
	}
	c := &Client{
		BaseURL:    u,
		httpClient: retryablehttp.NewClient(),
		UserAgent:  "go-sefaria/v1",
		validate:   validator.New(validator.WithRequiredStructEnabled()),
	}

	c.httpClient.RetryMax = 3
	c.httpClient.RetryWaitMin = 150 * time.Millisecond
	c.httpClient.RetryWaitMax = 1 * time.Second
	c.httpClient.Logger = nil // disable default logging

	c.common.client = c
	c.Text = (*TextService)(&c.common)
	c.Index = (*IndexService)(&c.common)
	c.Related = (*RelatedService)(&c.common)
	c.Calendar = (*CalendarService)(&c.common)
	c.Lexicon = (*LexiconService)(&c.common)
	c.Topics = (*TopicService)(&c.common)
	c.Terms = (*TermService)(&c.common)

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithAPIEndpoint(endpoint string) ClientOption {
	u, err := url.Parse(endpoint)
	return func(c *Client) {
		if err != nil {
			c.log(slog.LevelError, "invalid API endpoint URL", slog.String("url", endpoint))
			return
		}
		c.BaseURL = u
	}
}

func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
		c.httpClient.Logger = logger
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient.HTTPClient = httpClient
	}
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

func (c *Client) validateStruct(s any) error {
	if err := c.validate.Struct(s); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			msgs := make([]string, 0, len(ve))
			for _, fe := range ve {
				// Simple message: "Field <Field> failed <Tag> validation"
				msg := fmt.Sprintf("Field '%s' failed validation: %s", fe.Field(), fe.Tag())
				if fe.Param() != "" {
					msg += fmt.Sprintf(" (expected %s)", fe.Param())
				}
				msgs = append(msgs, msg)
			}
			return errors.New(strings.Join(msgs, "; "))
		}
		return err
	}
	return nil
}

func (c *Client) log(level slog.Level, msg string, attrs ...slog.Attr) {
	if c.logger != nil {
		c.logger.LogAttrs(context.Background(), level, msg, attrs...)
	}
}

type service struct {
	client *Client
}
