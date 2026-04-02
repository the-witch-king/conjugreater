package wanikani

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	DefaultBaseURL = "https://api.wanikani.com/v2"
	Revision       = "20170710"
)

// Client is an HTTP client for the WaniKani API.
type Client struct {
	Token      string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new WaniKani API client.
func NewClient(token string) *Client {
	return &Client{
		Token:      token,
		BaseURL:    DefaultBaseURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// FetchSubjects fetches all vocabulary subjects of the given types.
// subjectTypes should be e.g. []string{"vocabulary", "kana_vocabulary"}.
// If updatedAfter is non-empty, only subjects updated after that timestamp are returned.
func (c *Client) FetchSubjects(ctx context.Context, subjectTypes []string, updatedAfter string) ([]Resource[SubjectData], string, error) {
	url := c.BaseURL + "/subjects?types=" + joinTypes(subjectTypes)
	if updatedAfter != "" {
		url += "&updated_after=" + updatedAfter
	}
	return fetchAll[Resource[SubjectData]](c, ctx, url)
}

// FetchAssignments fetches all assignments for the given subject types.
// If updatedAfter is non-empty, only assignments updated after that timestamp are returned.
func (c *Client) FetchAssignments(ctx context.Context, subjectTypes []string, updatedAfter string) ([]Resource[AssignmentData], string, error) {
	url := c.BaseURL + "/assignments?subject_types=" + joinTypes(subjectTypes)
	if updatedAfter != "" {
		url += "&updated_after=" + updatedAfter
	}
	return fetchAll[Resource[AssignmentData]](c, ctx, url)
}

// fetchAll paginates through all pages and returns the collected data plus the
// data_updated_at timestamp from the first page response.
func fetchAll[T any](c *Client, ctx context.Context, initialURL string) ([]T, string, error) {
	var all []T
	var dataUpdatedAt string
	url := initialURL
	page := 1

	for url != "" {
		var coll Collection[T]
		if err := c.doRequest(ctx, url, &coll); err != nil {
			return nil, "", fmt.Errorf("page %d: %w", page, err)
		}

		if page == 1 {
			dataUpdatedAt = coll.DataUpdatedAt
			log.Printf("Fetching %d items...", coll.TotalCount)
		}

		all = append(all, coll.Data...)
		log.Printf("  page %d: got %d items (%d/%d)", page, len(coll.Data), len(all), coll.TotalCount)

		if coll.Pages.NextURL != nil {
			url = *coll.Pages.NextURL
		} else {
			url = ""
		}
		page++
	}

	return all, dataUpdatedAt, nil
}

func (c *Client) doRequest(ctx context.Context, url string, target any) error {
	const maxRetries = 3

	for attempt := range maxRetries {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("creating request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+c.Token)
		req.Header.Set("Wanikani-Revision", Revision)

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return fmt.Errorf("executing request: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("reading response: %w", err)
		}

		switch resp.StatusCode {
		case http.StatusOK:
			if err := json.Unmarshal(body, target); err != nil {
				return fmt.Errorf("decoding response: %w", err)
			}
			return nil

		case http.StatusUnauthorized:
			return fmt.Errorf("invalid API token (401). Get yours at https://www.wanikani.com/settings/personal_access_tokens")

		case http.StatusTooManyRequests:
			retryAfter := 5
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if parsed, err := strconv.Atoi(ra); err == nil {
					retryAfter = parsed
				}
			}
			if attempt < maxRetries-1 {
				log.Printf("Rate limited, waiting %ds (attempt %d/%d)...", retryAfter, attempt+1, maxRetries)
				select {
				case <-time.After(time.Duration(retryAfter) * time.Second):
				case <-ctx.Done():
					return ctx.Err()
				}
				continue
			}
			return fmt.Errorf("rate limited after %d retries", maxRetries)

		default:
			return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		}
	}

	return fmt.Errorf("exhausted retries")
}

func joinTypes(types []string) string {
	result := ""
	for i, t := range types {
		if i > 0 {
			result += ","
		}
		result += t
	}
	return result
}
