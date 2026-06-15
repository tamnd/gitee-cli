package gitee

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ─── Exported record types ────────────────────────────────────────────────────

// Repo is the canonical output record for any repository surface.
type Repo struct {
	Rank        int    `json:"rank"                  table:"rank"`
	FullName    string `json:"full_name"  kit:"id"   table:"full_name"`
	Description string `json:"description"           table:"description"`
	Language    string `json:"language"              table:"language"`
	Stars       int    `json:"stars"                 table:"stars"`
	Forks       int    `json:"forks"                 table:"forks"`
	UpdatedAt   string `json:"updated_at"            table:"updated_at"`
	URL         string `json:"url"                   table:"url,url"`
}

// User is the output record for a Gitee user profile.
type User struct {
	Login     string `json:"login"    kit:"id"  table:"login"`
	Name      string `json:"name"               table:"name"`
	Followers int    `json:"followers"          table:"followers"`
	Following int    `json:"following"          table:"following"`
	Repos     int    `json:"repos"              table:"repos"`
	Blog      string `json:"blog"               table:"blog"`
	URL       string `json:"url"                table:"url,url"`
}

// Release is the output record for a repository release.
type Release struct {
	Rank       int    `json:"rank"                 table:"rank"`
	TagName    string `json:"tag_name"  kit:"id"   table:"tag_name"`
	Name       string `json:"name"                 table:"name"`
	Prerelease bool   `json:"prerelease"           table:"prerelease"`
	CreatedAt  string `json:"created_at"           table:"created_at"`
	URL        string `json:"url"                  table:"url,url"`
}

// ─── Wire types (internal) ────────────────────────────────────────────────────

type wireRepo struct {
	ID              int          `json:"id"`
	FullName        string       `json:"full_name"`
	HumanName       string       `json:"human_name"`
	Name            string       `json:"name"`
	Path            string       `json:"path"`
	Description     string       `json:"description"`
	Private         bool         `json:"private"`
	Fork            bool         `json:"fork"`
	HTMLURL         string       `json:"html_url"`
	ForksCount      int          `json:"forks_count"`
	StargazersCount int          `json:"stargazers_count"`
	WatchersCount   int          `json:"watchers_count"`
	OpenIssuesCount int          `json:"open_issues_count"`
	DefaultBranch   string       `json:"default_branch"`
	Language        string       `json:"language"`
	License         *wireLicense `json:"license"`
	PushedAt        string       `json:"pushed_at"`
	CreatedAt       string       `json:"created_at"`
	UpdatedAt       string       `json:"updated_at"`
	Owner           wireOwner    `json:"owner"`
}

// wireLicense can appear as either an object {"spdx_id":"MIT"} or a bare
// string "MIT" depending on which Gitee endpoint returns the repo.
type wireLicense struct {
	SPDXID string
}

func (l *wireLicense) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	// bare string form
	if b[0] == '"' {
		return json.Unmarshal(b, &l.SPDXID)
	}
	// object form
	var obj struct {
		SPDXID string `json:"spdx_id"`
	}
	if err := json.Unmarshal(b, &obj); err != nil {
		return err
	}
	l.SPDXID = obj.SPDXID
	return nil
}

type wireOwner struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

type wireUser struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
	Blog        string `json:"blog"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	PublicRepos int    `json:"public_repos"`
	HTMLURL     string `json:"html_url"`
	CreatedAt   string `json:"created_at"`
}

type wireRelease struct {
	ID         int       `json:"id"`
	TagName    string    `json:"tag_name"`
	Name       string    `json:"name"`
	Prerelease bool      `json:"prerelease"`
	Draft      bool      `json:"draft"`
	Body       string    `json:"body"`
	CreatedAt  string    `json:"created_at"`
	Author     wireOwner `json:"author"`
}

// ─── Conversion helpers ───────────────────────────────────────────────────────

func wireRepoToRepo(wr wireRepo, rank int) Repo {
	// The API html_url includes a trailing .git; strip it for the web URL.
	webURL := strings.TrimSuffix(wr.HTMLURL, ".git")
	if webURL == "" && wr.FullName != "" {
		webURL = "https://gitee.com/" + wr.FullName
	}
	return Repo{
		Rank:        rank,
		FullName:    wr.FullName,
		Description: wr.Description,
		Language:    wr.Language,
		Stars:       wr.StargazersCount,
		Forks:       wr.ForksCount,
		UpdatedAt:   wr.UpdatedAt,
		URL:         webURL,
	}
}

func wireUserToUser(wu wireUser) User {
	return User{
		Login:     wu.Login,
		Name:      wu.Name,
		Followers: wu.Followers,
		Following: wu.Following,
		Repos:     wu.PublicRepos,
		Blog:      wu.Blog,
		URL:       fmt.Sprintf("https://gitee.com/%s", wu.Login),
	}
}

func wireReleaseToRelease(wr wireRelease, owner, repo string, rank int) Release {
	return Release{
		Rank:       rank,
		TagName:    wr.TagName,
		Name:       wr.Name,
		Prerelease: wr.Prerelease,
		CreatedAt:  wr.CreatedAt,
		URL:        fmt.Sprintf("https://gitee.com/%s/%s/releases/tag/%s", owner, repo, wr.TagName),
	}
}
