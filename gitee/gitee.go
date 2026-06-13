// Package gitee is the library behind the gitee command: the HTTP client,
// request shaping, and the typed data models for Gitee.
//
// One API: the official Gitee REST v5 endpoint at https://gitee.com/api/v5.
// The endpoints used here are open and require no authentication token.
package gitee

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DefaultUserAgent identifies the client to Gitee.
const DefaultUserAgent = "gitee/dev (+https://github.com/tamnd/gitee-cli)"

// ErrNotFound is returned when the API returns null for an object.
var ErrNotFound = errors.New("not found")

// Config holds constructor parameters.
type Config struct {
	BaseURL   string
	UserAgent string
	Rate      time.Duration
	Retries   int
	Timeout   time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://gitee.com/api/v5",
		UserAgent: DefaultUserAgent,
		Rate:      200 * time.Millisecond,
		Retries:   3,
		Timeout:   30 * time.Second,
	}
}

// Client talks to the Gitee API.
type Client struct {
	cfg  Config
	http *http.Client
	mu   sync.Mutex
	last time.Time
}

// NewClient returns a Client with the given config.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: cfg.Timeout},
	}
}

// get fetches a URL with pacing and retries.
func (c *Client) get(ctx context.Context, rawURL string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}
		body, retry, err := c.do(ctx, rawURL)
		if err == nil {
			return body, nil
		}
		lastErr = err
		if !retry {
			return nil, err
		}
	}
	return nil, fmt.Errorf("get %s: %w", rawURL, lastErr)
}

func (c *Client) do(ctx context.Context, rawURL string) ([]byte, bool, error) {
	c.pace()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, true, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return nil, true, fmt.Errorf("http %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("http %d", resp.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, true, err
	}
	return b, false, nil
}

func (c *Client) pace() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cfg.Rate <= 0 {
		return
	}
	if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}

func backoff(attempt int) time.Duration {
	d := time.Duration(attempt) * 500 * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}

// getJSON fetches and JSON-decodes into v. Returns ErrNotFound when the body is null.
func (c *Client) getJSON(ctx context.Context, rawURL string, v any) error {
	body, err := c.get(ctx, rawURL)
	if err != nil {
		return err
	}
	trimmed := strings.TrimSpace(string(body))
	if trimmed == "null" {
		return ErrNotFound
	}
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("decode %s: %w", rawURL, err)
	}
	return nil
}

// ─── API methods ─────────────────────────────────────────────────────────────

// SearchRepos searches Gitee repositories. sort: "stars", "forks", "updated".
// The Gitee search endpoint returns a plain JSON array (not a wrapper object).
func (c *Client) SearchRepos(ctx context.Context, query, lang, sort string, limit int) ([]Repo, error) {
	if limit <= 0 {
		limit = 20
	}
	giteeSort := "stars_count"
	switch sort {
	case "forks":
		giteeSort = "forks_count"
	case "updated":
		giteeSort = "updated"
	}
	q := query
	if lang != "" {
		q += "+language:" + lang
	}

	pageSize := limit
	if pageSize > 100 {
		pageSize = 100
	}

	var out []Repo
	page := 1
	for {
		params := url.Values{}
		params.Set("q", q)
		params.Set("sort", giteeSort)
		params.Set("order", "desc")
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))

		rawURL := c.cfg.BaseURL + "/search/repositories?" + params.Encode()
		var repos []wireRepo
		if err := c.getJSON(ctx, rawURL, &repos); err != nil {
			return out, err
		}
		for _, wr := range repos {
			out = append(out, wireRepoToRepo(wr, len(out)+1))
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(repos) == 0 {
			break
		}
		page++
	}
	return out, nil
}

// GetRepo fetches a single repository.
func (c *Client) GetRepo(ctx context.Context, owner, repo string) (Repo, error) {
	rawURL := c.cfg.BaseURL + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo)
	var wr wireRepo
	if err := c.getJSON(ctx, rawURL, &wr); err != nil {
		return Repo{}, err
	}
	return wireRepoToRepo(wr, 0), nil
}

// TrendingRepos fetches trending repos from the explore endpoint.
// sort: "stars" (default), "newest", "updated".
func (c *Client) TrendingRepos(ctx context.Context, lang, sort string, limit int) ([]Repo, error) {
	if limit <= 0 {
		limit = 20
	}
	giteeSort := "most_stars"
	switch sort {
	case "newest":
		giteeSort = "newest"
	case "updated":
		giteeSort = "recently_updated"
	}

	pageSize := limit
	if pageSize > 100 {
		pageSize = 100
	}

	var out []Repo
	page := 1
	for {
		params := url.Values{}
		params.Set("sort", giteeSort)
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))
		if lang != "" {
			params.Set("language", lang)
		}

		rawURL := c.cfg.BaseURL + "/repos/explore?" + params.Encode()
		var repos []wireRepo
		if err := c.getJSON(ctx, rawURL, &repos); err != nil {
			return out, err
		}
		for _, wr := range repos {
			out = append(out, wireRepoToRepo(wr, len(out)+1))
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(repos) == 0 {
			break
		}
		page++
	}
	return out, nil
}

// GetUser fetches a user profile.
func (c *Client) GetUser(ctx context.Context, username string) (User, error) {
	rawURL := c.cfg.BaseURL + "/users/" + url.PathEscape(username)
	var wu wireUser
	if err := c.getJSON(ctx, rawURL, &wu); err != nil {
		return User{}, err
	}
	return wireUserToUser(wu), nil
}

// UserRepos fetches a user's public repositories.
func (c *Client) UserRepos(ctx context.Context, username string, limit int) ([]Repo, error) {
	if limit <= 0 {
		limit = 20
	}
	pageSize := limit
	if pageSize > 100 {
		pageSize = 100
	}

	var out []Repo
	page := 1
	for {
		params := url.Values{}
		params.Set("sort", "updated")
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))

		rawURL := c.cfg.BaseURL + "/users/" + url.PathEscape(username) + "/repos?" + params.Encode()
		var repos []wireRepo
		if err := c.getJSON(ctx, rawURL, &repos); err != nil {
			return out, err
		}
		for _, wr := range repos {
			out = append(out, wireRepoToRepo(wr, len(out)+1))
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(repos) == 0 {
			break
		}
		page++
	}
	return out, nil
}

// ListReleases fetches releases for a repository.
func (c *Client) ListReleases(ctx context.Context, owner, repo string, limit int) ([]Release, error) {
	if limit <= 0 {
		limit = 10
	}
	pageSize := limit
	if pageSize > 100 {
		pageSize = 100
	}

	var out []Release
	page := 1
	for {
		params := url.Values{}
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))

		rawURL := c.cfg.BaseURL + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/releases?" + params.Encode()
		var releases []wireRelease
		if err := c.getJSON(ctx, rawURL, &releases); err != nil {
			return out, err
		}
		for _, wr := range releases {
			out = append(out, wireReleaseToRelease(wr, owner, repo, len(out)+1))
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(releases) == 0 {
			break
		}
		page++
	}
	return out, nil
}
