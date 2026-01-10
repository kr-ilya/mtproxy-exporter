package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// MaxResponseSize is the maximum size of the response body (10MB)
	MaxResponseSize = 10 * 1024 * 1024
)

// Client represents an HTTP client for fetching MTProxy stats
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// New creates a new MTProxy client with validation
func New(baseURL string, timeout time.Duration) (*Client, error) {
	// Validate URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("invalid URL scheme: %s (expected http or https)", u.Scheme)
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:          10,
				MaxIdleConnsPerHost:   5,
				MaxConnsPerHost:       10,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				DisableKeepAlives:     false,
			},
		},
	}, nil
}

// GetStats fetches the stats from MTProxy
func (c *Client) GetStats(ctx context.Context) (string, error) {
	url := c.baseURL + "/stats"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch stats: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Limit the response size to prevent DoS attacks or bugs in MTProxy
	limitedReader := io.LimitReader(resp.Body, MaxResponseSize)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if we hit the limit
	if len(body) == MaxResponseSize {
		return "", fmt.Errorf("response size exceeded limit of %d bytes", MaxResponseSize)
	}

	return string(body), nil
}
