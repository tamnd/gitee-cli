package gitee

import (
	"context"
	"net/url"
	"strconv"
)

const maxPerPage = 100

// GetUser fetches a user profile.
func (c *Client) GetUser(ctx context.Context, username string) (User, error) {
	rawURL := c.cfg.BaseURL + "/users/" + url.PathEscape(username)
	var w wireUser
	if err := c.getJSON(ctx, rawURL, &w); err != nil {
		return User{}, err
	}
	return w.toPublic(), nil
}

// UserRepos fetches a user's public repositories.
func (c *Client) UserRepos(ctx context.Context, username, sort, direction, typ string, limit int) ([]Repo, error) {
	if limit <= 0 {
		limit = 20
	}
	if sort == "" {
		sort = "updated"
	}
	if direction == "" {
		direction = "desc"
	}
	pageSize := limit
	if pageSize > maxPerPage {
		pageSize = maxPerPage
	}
	var out []Repo
	for page := 1; ; page++ {
		params := url.Values{}
		params.Set("sort", sort)
		params.Set("direction", direction)
		if typ != "" {
			params.Set("type", typ)
		}
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))
		rawURL := c.cfg.BaseURL + "/users/" + url.PathEscape(username) + "/repos?" + params.Encode()
		var repos []wireRepo
		if err := c.getJSON(ctx, rawURL, &repos); err != nil {
			return out, err
		}
		for _, wr := range repos {
			out = append(out, wr.toPublic())
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(repos) < pageSize {
			break
		}
	}
	return out, nil
}

// Followers fetches a user's followers.
func (c *Client) Followers(ctx context.Context, username string, limit int) ([]User, error) {
	return c.fetchUsers(ctx, "/users/"+url.PathEscape(username)+"/followers", limit)
}

// Following fetches users that a user follows.
func (c *Client) Following(ctx context.Context, username string, limit int) ([]User, error) {
	return c.fetchUsers(ctx, "/users/"+url.PathEscape(username)+"/following", limit)
}

func (c *Client) fetchUsers(ctx context.Context, path string, limit int) ([]User, error) {
	if limit <= 0 {
		limit = 20
	}
	pageSize := limit
	if pageSize > maxPerPage {
		pageSize = maxPerPage
	}
	var out []User
	for page := 1; ; page++ {
		params := url.Values{}
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))
		rawURL := c.cfg.BaseURL + path + "?" + params.Encode()
		var users []wireUser
		if err := c.getJSON(ctx, rawURL, &users); err != nil {
			return out, err
		}
		for _, w := range users {
			out = append(out, w.toPublic())
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(users) < pageSize {
			break
		}
	}
	return out, nil
}

// GetRepo fetches a single repository.
func (c *Client) GetRepo(ctx context.Context, owner, repo string) (Repo, error) {
	rawURL := c.cfg.BaseURL + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo)
	var w wireRepo
	if err := c.getJSON(ctx, rawURL, &w); err != nil {
		return Repo{}, err
	}
	return w.toPublic(), nil
}

// Commits fetches commits for a repository.
func (c *Client) Commits(ctx context.Context, owner, repo, sha, path, author, since, until string, limit int) ([]Commit, error) {
	if limit <= 0 {
		limit = 20
	}
	pageSize := limit
	if pageSize > maxPerPage {
		pageSize = maxPerPage
	}
	var out []Commit
	for page := 1; ; page++ {
		params := url.Values{}
		if sha != "" {
			params.Set("sha", sha)
		}
		if path != "" {
			params.Set("path", path)
		}
		if author != "" {
			params.Set("author", author)
		}
		if since != "" {
			params.Set("since", since)
		}
		if until != "" {
			params.Set("until", until)
		}
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))
		rawURL := c.cfg.BaseURL + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/commits?" + params.Encode()
		var commits []wireCommit
		if err := c.getJSON(ctx, rawURL, &commits); err != nil {
			return out, err
		}
		for _, w := range commits {
			out = append(out, w.toPublic())
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(commits) < pageSize {
			break
		}
	}
	return out, nil
}

// Branches fetches branches for a repository.
func (c *Client) Branches(ctx context.Context, owner, repo string, limit int) ([]Branch, error) {
	if limit <= 0 {
		limit = 20
	}
	pageSize := limit
	if pageSize > maxPerPage {
		pageSize = maxPerPage
	}
	var out []Branch
	for page := 1; ; page++ {
		params := url.Values{}
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))
		rawURL := c.cfg.BaseURL + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/branches?" + params.Encode()
		var branches []wireBranch
		if err := c.getJSON(ctx, rawURL, &branches); err != nil {
			return out, err
		}
		for _, w := range branches {
			out = append(out, w.toPublic())
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(branches) < pageSize {
			break
		}
	}
	return out, nil
}

// Tags fetches tags for a repository.
func (c *Client) Tags(ctx context.Context, owner, repo, sort, direction string, limit int) ([]Tag, error) {
	if limit <= 0 {
		limit = 20
	}
	pageSize := limit
	if pageSize > maxPerPage {
		pageSize = maxPerPage
	}
	var out []Tag
	for page := 1; ; page++ {
		params := url.Values{}
		if sort != "" {
			params.Set("sort_by", sort)
		}
		if direction != "" {
			params.Set("direction", direction)
		}
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))
		rawURL := c.cfg.BaseURL + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/tags?" + params.Encode()
		var tags []wireTag
		if err := c.getJSON(ctx, rawURL, &tags); err != nil {
			return out, err
		}
		for _, w := range tags {
			out = append(out, w.toPublic())
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(tags) < pageSize {
			break
		}
	}
	return out, nil
}

// Releases fetches releases for a repository.
func (c *Client) Releases(ctx context.Context, owner, repo string, limit int) ([]Release, error) {
	if limit <= 0 {
		limit = 10
	}
	pageSize := limit
	if pageSize > maxPerPage {
		pageSize = maxPerPage
	}
	var out []Release
	for page := 1; ; page++ {
		params := url.Values{}
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))
		rawURL := c.cfg.BaseURL + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/releases?" + params.Encode()
		var releases []wireRelease
		if err := c.getJSON(ctx, rawURL, &releases); err != nil {
			return out, err
		}
		for _, w := range releases {
			out = append(out, w.toPublic())
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(releases) < pageSize {
			break
		}
	}
	return out, nil
}

// Issues fetches issues for a repository.
func (c *Client) Issues(ctx context.Context, owner, repo, state, sort, direction string, limit int) ([]Issue, error) {
	if limit <= 0 {
		limit = 20
	}
	if state == "" {
		state = "open"
	}
	if sort == "" {
		sort = "created"
	}
	if direction == "" {
		direction = "desc"
	}
	pageSize := limit
	if pageSize > maxPerPage {
		pageSize = maxPerPage
	}
	var out []Issue
	for page := 1; ; page++ {
		params := url.Values{}
		params.Set("state", state)
		params.Set("sort", sort)
		params.Set("direction", direction)
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))
		rawURL := c.cfg.BaseURL + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/issues?" + params.Encode()
		var issues []wireIssue
		if err := c.getJSON(ctx, rawURL, &issues); err != nil {
			return out, err
		}
		for _, w := range issues {
			out = append(out, w.toPublic())
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(issues) < pageSize {
			break
		}
	}
	return out, nil
}

// Pulls fetches pull requests for a repository.
func (c *Client) Pulls(ctx context.Context, owner, repo, state, sort, direction string, limit int) ([]PullRequest, error) {
	if limit <= 0 {
		limit = 20
	}
	if state == "" {
		state = "open"
	}
	if sort == "" {
		sort = "created"
	}
	if direction == "" {
		direction = "desc"
	}
	pageSize := limit
	if pageSize > maxPerPage {
		pageSize = maxPerPage
	}
	var out []PullRequest
	for page := 1; ; page++ {
		params := url.Values{}
		params.Set("state", state)
		params.Set("sort", sort)
		params.Set("direction", direction)
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))
		rawURL := c.cfg.BaseURL + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/pulls?" + params.Encode()
		var pulls []wirePullRequest
		if err := c.getJSON(ctx, rawURL, &pulls); err != nil {
			return out, err
		}
		for _, w := range pulls {
			out = append(out, w.toPublic())
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(pulls) < pageSize {
			break
		}
	}
	return out, nil
}

// Readme fetches the README for a repository.
func (c *Client) Readme(ctx context.Context, owner, repo, ref string) (ReadmeFile, error) {
	params := url.Values{}
	if ref != "" {
		params.Set("ref", ref)
	}
	rawURL := c.cfg.BaseURL + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/readme"
	if len(params) > 0 {
		rawURL += "?" + params.Encode()
	}
	var w wireReadme
	if err := c.getJSON(ctx, rawURL, &w); err != nil {
		return ReadmeFile{}, err
	}
	return w.toPublic(), nil
}

// Tree fetches the git tree for a repository.
func (c *Client) Tree(ctx context.Context, owner, repo, ref string, recursive bool) ([]TreeEntry, error) {
	if ref == "" {
		ref = "HEAD"
	}
	params := url.Values{}
	if recursive {
		params.Set("recursive", "1")
	}
	rawURL := c.cfg.BaseURL + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/git/trees/" + url.PathEscape(ref)
	if len(params) > 0 {
		rawURL += "?" + params.Encode()
	}
	var w wireTree
	if err := c.getJSON(ctx, rawURL, &w); err != nil {
		return nil, err
	}
	return w.toPublic(), nil
}

// Stargazers fetches users who starred a repository.
func (c *Client) Stargazers(ctx context.Context, owner, repo string, limit int) ([]User, error) {
	return c.fetchUsers(ctx, "/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(repo)+"/stargazers", limit)
}

// Forks fetches forks of a repository.
func (c *Client) Forks(ctx context.Context, owner, repo, sort string, limit int) ([]Repo, error) {
	if limit <= 0 {
		limit = 20
	}
	if sort == "" {
		sort = "newest"
	}
	pageSize := limit
	if pageSize > maxPerPage {
		pageSize = maxPerPage
	}
	var out []Repo
	for page := 1; ; page++ {
		params := url.Values{}
		params.Set("sort", sort)
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))
		rawURL := c.cfg.BaseURL + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/forks?" + params.Encode()
		var repos []wireRepo
		if err := c.getJSON(ctx, rawURL, &repos); err != nil {
			return out, err
		}
		for _, w := range repos {
			out = append(out, w.toPublic())
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(repos) < pageSize {
			break
		}
	}
	return out, nil
}

// Contributors fetches contributors to a repository.
func (c *Client) Contributors(ctx context.Context, owner, repo, typ string, limit int) ([]Contributor, error) {
	if limit <= 0 {
		limit = 20
	}
	pageSize := limit
	if pageSize > maxPerPage {
		pageSize = maxPerPage
	}
	var out []Contributor
	for page := 1; ; page++ {
		params := url.Values{}
		if typ != "" {
			params.Set("type", typ)
		}
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))
		rawURL := c.cfg.BaseURL + "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/contributors?" + params.Encode()
		var contributors []wireContributor
		if err := c.getJSON(ctx, rawURL, &contributors); err != nil {
			return out, err
		}
		for _, w := range contributors {
			out = append(out, w.toPublic())
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(contributors) < pageSize {
			break
		}
	}
	return out, nil
}

// SearchRepos searches Gitee repositories using Elasticsearch-style from/size.
func (c *Client) SearchRepos(ctx context.Context, query, sort string, limit int) ([]Repo, error) {
	if limit <= 0 {
		limit = 20
	}
	if sort == "" {
		sort = "stars_count"
	} else {
		switch sort {
		case "forks":
			sort = "forks_count"
		case "updated":
			sort = "updated"
		case "stars":
			sort = "stars_count"
		}
	}
	pageSize := limit
	if pageSize > maxPerPage {
		pageSize = maxPerPage
	}
	var out []Repo
	for from := 0; ; from += pageSize {
		params := url.Values{}
		params.Set("q", query)
		params.Set("sort", sort)
		params.Set("order", "desc")
		params.Set("from", strconv.Itoa(from))
		params.Set("size", strconv.Itoa(pageSize))
		rawURL := c.cfg.BaseURL + "/search/repositories?" + params.Encode()
		var repos []wireRepo
		if err := c.getJSON(ctx, rawURL, &repos); err != nil {
			return out, err
		}
		for _, w := range repos {
			out = append(out, w.toPublic())
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(repos) < pageSize {
			break
		}
	}
	return out, nil
}

// SearchUsers searches Gitee users using Elasticsearch-style from/size.
func (c *Client) SearchUsers(ctx context.Context, query string, limit int) ([]User, error) {
	if limit <= 0 {
		limit = 20
	}
	pageSize := limit
	if pageSize > maxPerPage {
		pageSize = maxPerPage
	}
	var out []User
	for from := 0; ; from += pageSize {
		params := url.Values{}
		params.Set("q", query)
		params.Set("from", strconv.Itoa(from))
		params.Set("size", strconv.Itoa(pageSize))
		rawURL := c.cfg.BaseURL + "/search/users?" + params.Encode()
		var users []wireUser
		if err := c.getJSON(ctx, rawURL, &users); err != nil {
			return out, err
		}
		for _, w := range users {
			out = append(out, w.toPublic())
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(users) < pageSize {
			break
		}
	}
	return out, nil
}

// GetOrg fetches an organization profile.
func (c *Client) GetOrg(ctx context.Context, name string) (OrgProfile, error) {
	rawURL := c.cfg.BaseURL + "/orgs/" + url.PathEscape(name)
	var w wireOrg
	if err := c.getJSON(ctx, rawURL, &w); err != nil {
		return OrgProfile{}, err
	}
	return w.toPublic(), nil
}

// OrgRepos fetches repositories for an organization.
func (c *Client) OrgRepos(ctx context.Context, name, typ, sort, direction string, limit int) ([]Repo, error) {
	if limit <= 0 {
		limit = 20
	}
	if sort == "" {
		sort = "updated"
	}
	if direction == "" {
		direction = "desc"
	}
	pageSize := limit
	if pageSize > maxPerPage {
		pageSize = maxPerPage
	}
	var out []Repo
	for page := 1; ; page++ {
		params := url.Values{}
		if typ != "" {
			params.Set("type", typ)
		}
		params.Set("sort", sort)
		params.Set("direction", direction)
		params.Set("page", strconv.Itoa(page))
		params.Set("per_page", strconv.Itoa(pageSize))
		rawURL := c.cfg.BaseURL + "/orgs/" + url.PathEscape(name) + "/repos?" + params.Encode()
		var repos []wireRepo
		if err := c.getJSON(ctx, rawURL, &repos); err != nil {
			return out, err
		}
		for _, w := range repos {
			out = append(out, w.toPublic())
			if len(out) >= limit {
				return out, nil
			}
		}
		if len(repos) < pageSize {
			break
		}
	}
	return out, nil
}
