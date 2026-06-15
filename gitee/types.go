package gitee

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrRateLimit    = errors.New("rate limit exceeded")
)

type User struct {
	ID          int    `json:"id"            table:"ID"`
	Login       string `json:"login"         table:"Login"`
	Name        string `json:"name"          table:"Name"`
	AvatarURL   string `json:"avatar_url"    table:"-"`
	URL         string `json:"url"           table:"-"`
	HTMLURL     string `json:"html_url"      table:"URL"      kit:"url"`
	Bio         string `json:"bio"           table:"-"`
	Blog        string `json:"blog"          table:"Blog"`
	Weibo       string `json:"weibo"         table:"-"`
	Company     string `json:"company"       table:"Company"`
	Profession  string `json:"profession"    table:"-"`
	Email       string `json:"email"         table:"Email"`
	PublicRepos int    `json:"public_repos"  table:"Repos"`
	Followers   int    `json:"followers"     table:"Followers"`
	Following   int    `json:"following"     table:"Following"`
	Stared      int    `json:"stared"        table:"Starred"` // API typo — frozen
	CreatedAt   string `json:"created_at"    table:"Created"`
	UpdatedAt   string `json:"updated_at"    table:"-"`
}

type Repo struct {
	ID              int        `json:"id"                    table:"ID"`
	FullName        string     `json:"full_name"             table:"Full Name"`
	Name            string     `json:"name"                  table:"Name"`
	Owner           *User      `json:"owner"                 table:"-"`
	Namespace       *Namespace `json:"namespace"             table:"-"`
	Description     string     `json:"description"           table:"Description"`
	Private         bool       `json:"private"               table:"-"`
	Fork            bool       `json:"fork"                  table:"Fork"`
	URL             string     `json:"url"                   table:"URL"           kit:"url"` // html_url .git stripped
	SSHURL          string     `json:"ssh_url"               table:"-"`
	Recommend       bool       `json:"recommend"             table:"-"`
	GVP             bool       `json:"gvp"                   table:"GVP"`
	Language        string     `json:"language"              table:"Lang"`
	ForksCount      int        `json:"forks_count"           table:"Forks"`
	StargazersCount int        `json:"stargazers_count"      table:"Stars"`
	DefaultBranch   string     `json:"default_branch"        table:"Branch"`
	OpenIssuesCount int        `json:"open_issues_count"     table:"Issues"`
	License         *License   `json:"license"               table:"-"`
	PushedAt        string     `json:"pushed_at"             table:"Pushed"`
	CreatedAt       string     `json:"created_at"            table:"Created"`
	UpdatedAt       string     `json:"updated_at"            table:"-"`
	Parent          *Repo      `json:"parent"                table:"-"`
}

type Namespace struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	HTMLURL string `json:"html_url"`
}

type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SPDXID string `json:"spdx_id"`
}

type Commit struct {
	SHA     string       `json:"sha"     table:"SHA"`
	Author  CommitAuthor `json:"author"  table:"-"`
	Message string       `json:"message" table:"Message"`
	URL     string       `json:"url"     table:"URL"    kit:"url"`
}

type CommitAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

type Branch struct {
	Name      string `json:"name"      table:"Name"`
	SHA       string `json:"sha"       table:"SHA"`
	Protected bool   `json:"protected" table:"Protected"`
}

type Tag struct {
	Name    string `json:"name"    table:"Name"`
	SHA     string `json:"sha"     table:"SHA"`
	Message string `json:"message" table:"Message"`
	URL     string `json:"url"     table:"URL"    kit:"url"`
}

type Release struct {
	ID              int            `json:"id"               table:"ID"`
	TagName         string         `json:"tag_name"         table:"Tag"`
	Name            string         `json:"name"             table:"Name"`
	Body            string         `json:"body"             table:"-"`
	Prerelease      bool           `json:"prerelease"       table:"Pre"`
	Assets          []ReleaseAsset `json:"assets"           table:"-"`
	AssetsCount     int            `json:"assets_count"     table:"Assets"`
	Author          *User          `json:"author"           table:"-"`
	TargetCommitish string         `json:"target_commitish" table:"Branch"`
	CreatedAt       string         `json:"created_at"       table:"Created"`
}

type ReleaseAsset struct {
	ID                 int    `json:"id"                   table:"ID"`
	Name               string `json:"name"                 table:"Name"`
	Size               int    `json:"size"                 table:"Size"`
	DownloadCount      int    `json:"download_count"       table:"Downloads"`
	BrowserDownloadURL string `json:"browser_download_url" table:"URL"    kit:"url"`
}

// Issue.Number is STRING like "IJUFI0" — Gitee API contract, NOT integer.
type Issue struct {
	Number    string  `json:"number"     table:"Number"`
	Title     string  `json:"title"      table:"Title"`
	State     string  `json:"state"      table:"State"`
	Body      string  `json:"body"       table:"-"`
	User      *User   `json:"user"       table:"-"`
	Assignee  *User   `json:"assignee"   table:"-"`
	Labels    []Label `json:"labels"     table:"-"`
	Comments  int     `json:"comments"   table:"Comments"`
	URL       string  `json:"html_url"   table:"URL"     kit:"url"`
	CreatedAt string  `json:"created_at" table:"Created"`
	UpdatedAt string  `json:"updated_at" table:"-"`
}

type Label struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type PullRequest struct {
	Number    int     `json:"number"     table:"Number"`
	Title     string  `json:"title"      table:"Title"`
	State     string  `json:"state"      table:"State"`
	Body      string  `json:"body"       table:"-"`
	Head      PRRef   `json:"head"       table:"-"`
	Base      PRRef   `json:"base"       table:"-"`
	User      *User   `json:"user"       table:"-"`
	Labels    []Label `json:"labels"     table:"-"`
	Merged    bool    `json:"merged"     table:"Merged"`
	Comments  int     `json:"comments"   table:"Comments"`
	URL       string  `json:"html_url"   table:"URL"     kit:"url"`
	CreatedAt string  `json:"created_at" table:"Created"`
	UpdatedAt string  `json:"updated_at" table:"-"`
}

type PRRef struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	SHA   string `json:"sha"`
}

type Contributor struct {
	Name          string `json:"name"          table:"Name"`
	Email         string `json:"email"         table:"Email"`
	Contributions int    `json:"contributions" table:"Commits"`
}

type TreeEntry struct {
	Path string `json:"path" table:"Path"`
	Mode string `json:"mode" table:"Mode"`
	Type string `json:"type" table:"Type"`
	SHA  string `json:"sha"  table:"SHA"`
	Size int    `json:"size" table:"Size"`
}

type ReadmeFile struct {
	Name           string `json:"name"            table:"Name"`
	SHA            string `json:"sha"             table:"SHA"`
	Size           int    `json:"size"            table:"Size"`
	DecodedContent string `json:"decoded_content" table:"-"`
	HTMLURL        string `json:"html_url"        table:"URL"    kit:"url"`
	DownloadURL    string `json:"download_url"    table:"-"`
}

type OrgProfile struct {
	ID          int    `json:"id"           table:"ID"`
	Login       string `json:"login"        table:"Login"`
	Name        string `json:"name"         table:"Name"`
	HTMLURL     string `json:"html_url"     table:"URL"     kit:"url"`
	Description string `json:"description"  table:"Description"`
	Blog        string `json:"blog"         table:"Blog"`
	Email       string `json:"email"        table:"Email"`
	PublicRepos int    `json:"public_repos" table:"Repos"`
	Followers   int    `json:"followers"    table:"Followers"`
	CreatedAt   string `json:"created_at"   table:"Created"`
}
