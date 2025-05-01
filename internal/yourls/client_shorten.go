package yourls

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

// Shorten takes in a url and returns a shortened url with the title set internally on YOURLS.
//
// Since YOURLS does not support expiry by default, the shortened URL
// should be deleted manually / by this program after the MinIO presigned url expires.
func (c *Client) Shorten(ctx context.Context, input string, title string) (string, error) {
	u, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("invalid input provided: %w", err)
	}

	uid, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("could not generate a uuid for keyword: %w", err)
	}

	v := make(map[string]string)
	v["signature"] = c.signature
	v["action"] = "shorturl"
	v["format"] = "json"
	v["url"] = u.String()
	v["title"] = title
	v["keyword"] = uid.String()

	resp, err := c.doAPIRequest(ctx, v)
	if err != nil {
		return "", fmt.Errorf("failed to do api request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"unexpected status code: %d (status: %s)",
			resp.StatusCode,
			resp.Status,
		)
	}

	res := &shortenURLResponse{}
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return "", fmt.Errorf("could not decode json response: %w", err)
	}

	return res.Shorturl, nil
}

type shortenURLResponse struct {
	URL struct {
		Keyword string `json:"keyword"`
		URL     string `json:"url"`
		Title   string `json:"title"`
		Date    string `json:"date"`
		IP      string `json:"ip"`
	} `json:"url"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Title      string `json:"title"`
	Shorturl   string `json:"shorturl"`
	StatusCode string `json:"statusCode"`
}
