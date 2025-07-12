package httpadapter

import (
	"RSSHub/internal/domain/models"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"
)

var (
	blackList = []string{
		"https://news.ycombinator.com/rss",
	}
)

// HTTPClient defines the interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Adapter wraps an HTTP client to make it easier to test/make requests.
type Adapter struct {
	client HTTPClient
}

// NewClient creates a new HTTP adapter with optional custom timeout.
func NewClient(timeout time.Duration) *Adapter {
	return &Adapter{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (a *Adapter) FetchRSSFeed(ctx context.Context, url string) (*models.RSSFeed, error) {
	if slices.Contains(blackList, url) {
		return nil, errors.New("BLACK LIST NIGGA")
	}

	body, err := a.fetch(ctx, url)
	if err != nil {
		return nil, err
	}

	feed, err := Parse(body)
	if err != nil {
		return nil, err
	}

	feed.CreatedAt = time.Now()

	if feed.Channel.Link == "" {
		feed.Channel.Link = url
	}

	return feed, nil
}

// Fetch makes a GET request to the specified URL and returns the response body.
func (a *Adapter) fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("httpadapter: failed to create request: %w", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpadapter: request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("httpadapter: bad status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func Parse(r io.Reader) (*models.RSSFeed, error) {
	feed := new(models.RSSFeed)
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(feed); err != nil {
		return nil, fmt.Errorf("rssparser: failed to decode RSS XML: %w", err)
	}
	return feed, nil
}
