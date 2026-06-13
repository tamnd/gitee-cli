package gitee

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestClient(baseURL string) *Client {
	cfg := DefaultConfig()
	cfg.BaseURL = baseURL
	cfg.Rate = 0
	return NewClient(cfg)
}

func TestGetSendsUserAgent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") == "" {
			t.Error("request carried no User-Agent")
		}
		_, _ = w.Write([]byte(`"hello"`))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	body, err := c.get(context.Background(), srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != `"hello"` {
		t.Errorf("body = %q", body)
	}
}

func TestGetRetriesOn503(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if hits < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte(`"recovered"`))
	}))
	defer srv.Close()

	cfg := DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0
	cfg.Retries = 5
	c := NewClient(cfg)

	start := time.Now()
	body, err := c.get(context.Background(), srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != `"recovered"` {
		t.Errorf("body = %q after retries", body)
	}
	if hits != 3 {
		t.Errorf("server saw %d hits, want 3", hits)
	}
	if time.Since(start) < 500*time.Millisecond {
		t.Error("retries did not back off")
	}
}

func TestGetNullReturnsNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("null"))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	var v any
	err := c.getJSON(context.Background(), srv.URL, &v)
	if err != ErrNotFound {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestSearchRepos(t *testing.T) {
	resp := searchResp{
		TotalCount: 2,
		Items: []wireRepo{
			{FullName: "foo/bar", StargazersCount: 100, HTMLURL: "https://gitee.com/foo/bar"},
			{FullName: "baz/qux", StargazersCount: 50, HTMLURL: "https://gitee.com/baz/qux"},
		},
	}
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if calls > 1 {
			_ = json.NewEncoder(w).Encode(searchResp{})
			return
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	repos, err := c.SearchRepos(context.Background(), "test", "", "stars", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != 2 {
		t.Fatalf("got %d repos, want 2", len(repos))
	}
	if repos[0].FullName != "foo/bar" {
		t.Errorf("first repo = %q, want %q", repos[0].FullName, "foo/bar")
	}
	if repos[0].Stars != 100 {
		t.Errorf("stars = %d, want 100", repos[0].Stars)
	}
	if repos[0].Rank != 1 {
		t.Errorf("rank = %d, want 1", repos[0].Rank)
	}
}

func TestGetRepo(t *testing.T) {
	wr := wireRepo{
		FullName:        "gitee/gitee",
		StargazersCount: 999,
		HTMLURL:         "https://gitee.com/gitee/gitee",
		Language:        "Ruby",
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(wr)
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	repo, err := c.GetRepo(context.Background(), "gitee", "gitee")
	if err != nil {
		t.Fatal(err)
	}
	if repo.FullName != "gitee/gitee" {
		t.Errorf("full_name = %q", repo.FullName)
	}
	if repo.URL != "https://gitee.com/gitee/gitee" {
		t.Errorf("url = %q", repo.URL)
	}
	if repo.Stars != 999 {
		t.Errorf("stars = %d", repo.Stars)
	}
}

func TestTrendingRepos(t *testing.T) {
	repos := []wireRepo{
		{FullName: "a/b", StargazersCount: 200, HTMLURL: "https://gitee.com/a/b"},
		{FullName: "c/d", StargazersCount: 100, HTMLURL: "https://gitee.com/c/d"},
	}
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if calls > 1 {
			_ = json.NewEncoder(w).Encode([]wireRepo{})
			return
		}
		_ = json.NewEncoder(w).Encode(repos)
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	out, err := c.TrendingRepos(context.Background(), "", "stars", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("got %d repos, want 2", len(out))
	}
	if out[0].Rank != 1 {
		t.Errorf("rank = %d, want 1", out[0].Rank)
	}
	if out[1].Rank != 2 {
		t.Errorf("rank = %d, want 2", out[1].Rank)
	}
}

func TestGetUser(t *testing.T) {
	wu := wireUser{
		Login:       "testuser",
		Name:        "Test User",
		Followers:   42,
		Following:   10,
		PublicRepos: 7,
		Blog:        "https://example.com",
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(wu)
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	user, err := c.GetUser(context.Background(), "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if user.Login != "testuser" {
		t.Errorf("login = %q", user.Login)
	}
	if user.Followers != 42 {
		t.Errorf("followers = %d", user.Followers)
	}
	if user.URL != "https://gitee.com/testuser" {
		t.Errorf("url = %q", user.URL)
	}
}

func TestUserRepos(t *testing.T) {
	repos := []wireRepo{
		{FullName: "testuser/alpha", StargazersCount: 10, HTMLURL: "https://gitee.com/testuser/alpha"},
		{FullName: "testuser/beta", StargazersCount: 5, HTMLURL: "https://gitee.com/testuser/beta"},
	}
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if calls > 1 {
			_ = json.NewEncoder(w).Encode([]wireRepo{})
			return
		}
		_ = json.NewEncoder(w).Encode(repos)
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	out, err := c.UserRepos(context.Background(), "testuser", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("got %d repos, want 2", len(out))
	}
	if out[0].FullName != "testuser/alpha" {
		t.Errorf("first repo = %q", out[0].FullName)
	}
}

func TestListReleases(t *testing.T) {
	releases := []wireRelease{
		{TagName: "v1.0.0", Name: "First Release", Prerelease: false, CreatedAt: "2024-01-01T00:00:00+08:00"},
		{TagName: "v0.9.0", Name: "Beta", Prerelease: true, CreatedAt: "2023-12-01T00:00:00+08:00"},
	}
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if calls > 1 {
			_ = json.NewEncoder(w).Encode([]wireRelease{})
			return
		}
		_ = json.NewEncoder(w).Encode(releases)
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	out, err := c.ListReleases(context.Background(), "owner", "repo", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("got %d releases, want 2", len(out))
	}
	if out[0].TagName != "v1.0.0" {
		t.Errorf("tag = %q", out[0].TagName)
	}
	wantURL := "https://gitee.com/owner/repo/releases/tag/v1.0.0"
	if out[0].URL != wantURL {
		t.Errorf("url = %q, want %q", out[0].URL, wantURL)
	}
	if out[1].Prerelease != true {
		t.Errorf("prerelease = %v, want true", out[1].Prerelease)
	}
}
