package gitee

import "fmt"

// ─── Exported record types ────────────────────────────────────────────────────

// Repo is the canonical output record for any repository surface.
type Repo struct {
	Rank        int    `json:"rank"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	UpdatedAt   string `json:"updated_at"`
	URL         string `json:"url"`
}

// User is the output record for a Gitee user profile.
type User struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
	Repos     int    `json:"repos"`
	Blog      string `json:"blog"`
	URL       string `json:"url"`
}

// Release is the output record for a repository release.
type Release struct {
	Rank       int    `json:"rank"`
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Prerelease bool   `json:"prerelease"`
	CreatedAt  string `json:"created_at"`
	URL        string `json:"url"`
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

type wireLicense struct {
	SPDXID string `json:"spdx_id"`
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

type searchResp struct {
	TotalCount int        `json:"total_count"`
	Items      []wireRepo `json:"items"`
}

// ─── Conversion helpers ───────────────────────────────────────────────────────

func wireRepoToRepo(wr wireRepo, rank int) Repo {
	return Repo{
		Rank:        rank,
		FullName:    wr.FullName,
		Description: wr.Description,
		Language:    wr.Language,
		Stars:       wr.StargazersCount,
		Forks:       wr.ForksCount,
		UpdatedAt:   wr.UpdatedAt,
		URL:         wr.HTMLURL,
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
