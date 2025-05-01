package yourls

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// TODO: make http client configurable
// via options on NewClient
type Client struct {
	endpoint  string
	signature string
	client    *http.Client
}

func NewClient(endpoint string, signature string) *Client {
	return &Client{endpoint, signature, http.DefaultClient}
}

// since all api calls are post we can simply build
// a request using a context and url.Values
func (c *Client) doAPIRequest(
	ctx context.Context,
	values map[string]string,
) (*http.Response, error) {
	v, err := convertValues(values)
	if err != nil {
		return nil, fmt.Errorf("could not convert values: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, valuesToReader(v))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not get response: %w", err)
	}

	return resp, nil
}

func convertValues(values map[string]string) (url.Values, error) {
	v := url.Values{}
	for key, value := range values {
		v.Set(key, value)
	}

	return v, nil
}

func valuesToReader(values url.Values) io.Reader {
	return strings.NewReader(values.Encode())
}
