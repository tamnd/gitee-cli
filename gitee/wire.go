package gitee

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// ─── wire types ──────────────────────────────────────────────────────────────

type wireUser struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	Bio         string `json:"bio"`
	Blog        string `json:"blog"`
	Weibo       string `json:"weibo"`
	Company     string `json:"company"`
	Profession  string `json:"profession"`
	Email       string `json:"email"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	Stared      int    `json:"stared"` // API typo — not "starred"
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (w wireUser) toPublic() User {
	return User{
		ID:          w.ID,
		Login:       w.Login,
		Name:        w.Name,
		AvatarURL:   w.AvatarURL,
		URL:         w.URL,
		HTMLURL:     w.HTMLURL,
		Bio:         w.Bio,
		Blog:        w.Blog,
		Weibo:       w.Weibo,
		Company:     w.Company,
		Profession:  w.Profession,
		Email:       w.Email,
		PublicRepos: w.PublicRepos,
		Followers:   w.Followers,
		Following:   w.Following,
		Stared:      w.Stared,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}
}

type wireNamespace struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	HTMLURL string `json:"html_url"`
}

// wireLicense can appear as either an object {"spdx_id":"MIT"} or a bare
// string "MIT" depending on which Gitee endpoint returns the repo.
type wireLicense struct {
	Key    string
	Name   string
	SPDXID string
}

func (l *wireLicense) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	if b[0] == '"' {
		return json.Unmarshal(b, &l.SPDXID)
	}
	var obj struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		SPDXID string `json:"spdx_id"`
	}
	if err := json.Unmarshal(b, &obj); err != nil {
		return err
	}
	l.Key = obj.Key
	l.Name = obj.Name
	l.SPDXID = obj.SPDXID
	return nil
}

type wireRepo struct {
	ID              int            `json:"id"`
	FullName        string         `json:"full_name"`
	Name            string         `json:"name"`
	Owner           *wireUser      `json:"owner"`
	Namespace       *wireNamespace `json:"namespace"`
	Description     string         `json:"description"`
	Private         bool           `json:"private"`
	Fork            bool           `json:"fork"`
	HTMLURL         string         `json:"html_url"`
	SSHURL          string         `json:"ssh_url"`
	Recommend       bool           `json:"recommend"`
	GVP             bool           `json:"gvp"`
	Language        string         `json:"language"`
	ForksCount      int            `json:"forks_count"`
	StargazersCount int            `json:"stargazers_count"`
	DefaultBranch   string         `json:"default_branch"`
	OpenIssuesCount int            `json:"open_issues_count"`
	License         *wireLicense   `json:"license"`
	PushedAt        string         `json:"pushed_at"`
	CreatedAt       string         `json:"created_at"`
	UpdatedAt       string         `json:"updated_at"`
	Parent          *wireRepo      `json:"parent"`
}

func (w wireRepo) toPublic() Repo {
	webURL := strings.TrimSuffix(w.HTMLURL, ".git")
	if webURL == "" && w.FullName != "" {
		webURL = "https://gitee.com/" + w.FullName
	}

	var owner *User
	if w.Owner != nil {
		u := w.Owner.toPublic()
		owner = &u
	}

	var ns *Namespace
	if w.Namespace != nil {
		n := Namespace{
			ID:      w.Namespace.ID,
			Type:    w.Namespace.Type,
			Name:    w.Namespace.Name,
			Path:    w.Namespace.Path,
			HTMLURL: w.Namespace.HTMLURL,
		}
		ns = &n
	}

	var lic *License
	if w.License != nil {
		l := License{
			Key:    w.License.Key,
			Name:   w.License.Name,
			SPDXID: w.License.SPDXID,
		}
		lic = &l
	}

	var parent *Repo
	if w.Parent != nil {
		p := w.Parent.toPublic()
		parent = &p
	}

	return Repo{
		ID:              w.ID,
		FullName:        w.FullName,
		Name:            w.Name,
		Owner:           owner,
		Namespace:       ns,
		Description:     w.Description,
		Private:         w.Private,
		Fork:            w.Fork,
		URL:             webURL,
		SSHURL:          w.SSHURL,
		Recommend:       w.Recommend,
		GVP:             w.GVP,
		Language:        w.Language,
		ForksCount:      w.ForksCount,
		StargazersCount: w.StargazersCount,
		DefaultBranch:   w.DefaultBranch,
		OpenIssuesCount: w.OpenIssuesCount,
		License:         lic,
		PushedAt:        w.PushedAt,
		CreatedAt:       w.CreatedAt,
		UpdatedAt:       w.UpdatedAt,
		Parent:          parent,
	}
}

type wireCommitInner struct {
	Author  CommitAuthor `json:"author"`
	Message string       `json:"message"`
}

type wireCommit struct {
	SHA    string          `json:"sha"`
	Commit wireCommitInner `json:"commit"`
	URL    string          `json:"html_url"`
}

func (w wireCommit) toPublic() Commit {
	return Commit{
		SHA:     w.SHA,
		Author:  w.Commit.Author,
		Message: w.Commit.Message,
		URL:     w.URL,
	}
}

type wireBranchCommit struct {
	SHA string `json:"sha"`
}

type wireBranch struct {
	Name      string           `json:"name"`
	Commit    wireBranchCommit `json:"commit"`
	Protected bool             `json:"protected"`
}

func (w wireBranch) toPublic() Branch {
	return Branch{
		Name:      w.Name,
		SHA:       w.Commit.SHA,
		Protected: w.Protected,
	}
}

type wireTagCommit struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

type wireTag struct {
	Name    string        `json:"name"`
	Message string        `json:"message"`
	Commit  wireTagCommit `json:"commit"`
}

func (w wireTag) toPublic() Tag {
	return Tag{
		Name:    w.Name,
		SHA:     w.Commit.SHA,
		Message: w.Message,
		URL:     w.Commit.URL,
	}
}

type wireReleaseAsset struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Size               int    `json:"size"`
	DownloadCount      int    `json:"download_count"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type wireRelease struct {
	ID              int                `json:"id"`
	TagName         string             `json:"tag_name"`
	Name            string             `json:"name"`
	Body            string             `json:"body"`
	Prerelease      bool               `json:"prerelease"`
	Assets          []wireReleaseAsset `json:"assets"`
	AssetsCount     int                `json:"assets_count"`
	Author          *wireUser          `json:"author"`
	TargetCommitish string             `json:"target_commitish"`
	CreatedAt       string             `json:"created_at"`
}

func (w wireRelease) toPublic() Release {
	var author *User
	if w.Author != nil {
		u := w.Author.toPublic()
		author = &u
	}
	assets := make([]ReleaseAsset, len(w.Assets))
	for i, a := range w.Assets {
		assets[i] = ReleaseAsset{
			ID:                 a.ID,
			Name:               a.Name,
			Size:               a.Size,
			DownloadCount:      a.DownloadCount,
			BrowserDownloadURL: a.BrowserDownloadURL,
		}
	}
	return Release{
		ID:              w.ID,
		TagName:         w.TagName,
		Name:            w.Name,
		Body:            w.Body,
		Prerelease:      w.Prerelease,
		Assets:          assets,
		AssetsCount:     w.AssetsCount,
		Author:          author,
		TargetCommitish: w.TargetCommitish,
		CreatedAt:       w.CreatedAt,
	}
}

type wireLabel struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type wireIssue struct {
	Number    string      `json:"number"`
	Title     string      `json:"title"`
	State     string      `json:"state"`
	Body      string      `json:"body"`
	User      *wireUser   `json:"user"`
	Assignee  *wireUser   `json:"assignee"`
	Labels    []wireLabel `json:"labels"`
	Comments  int         `json:"comments"`
	HTMLURL   string      `json:"html_url"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

func (w wireIssue) toPublic() Issue {
	var user *User
	if w.User != nil {
		u := w.User.toPublic()
		user = &u
	}
	var assignee *User
	if w.Assignee != nil {
		a := w.Assignee.toPublic()
		assignee = &a
	}
	labels := make([]Label, len(w.Labels))
	for i, l := range w.Labels {
		labels[i] = Label{ID: l.ID, Name: l.Name, Color: l.Color}
	}
	return Issue{
		Number:    w.Number,
		Title:     w.Title,
		State:     w.State,
		Body:      w.Body,
		User:      user,
		Assignee:  assignee,
		Labels:    labels,
		Comments:  w.Comments,
		URL:       w.HTMLURL,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}

type wirePRRef struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	SHA   string `json:"sha"`
}

type wirePullRequest struct {
	Number    int         `json:"number"`
	Title     string      `json:"title"`
	State     string      `json:"state"`
	Body      string      `json:"body"`
	Head      wirePRRef   `json:"head"`
	Base      wirePRRef   `json:"base"`
	User      *wireUser   `json:"user"`
	Labels    []wireLabel `json:"labels"`
	Merged    bool        `json:"merged"`
	Comments  int         `json:"comments"`
	HTMLURL   string      `json:"html_url"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

func (w wirePullRequest) toPublic() PullRequest {
	var user *User
	if w.User != nil {
		u := w.User.toPublic()
		user = &u
	}
	labels := make([]Label, len(w.Labels))
	for i, l := range w.Labels {
		labels[i] = Label{ID: l.ID, Name: l.Name, Color: l.Color}
	}
	return PullRequest{
		Number:    w.Number,
		Title:     w.Title,
		State:     w.State,
		Body:      w.Body,
		Head:      PRRef{Label: w.Head.Label, Ref: w.Head.Ref, SHA: w.Head.SHA},
		Base:      PRRef{Label: w.Base.Label, Ref: w.Base.Ref, SHA: w.Base.SHA},
		User:      user,
		Labels:    labels,
		Merged:    w.Merged,
		Comments:  w.Comments,
		URL:       w.HTMLURL,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}

type wireContributor struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Contributions int    `json:"contributions"`
}

func (w wireContributor) toPublic() Contributor {
	return Contributor{
		Name:          w.Name,
		Email:         w.Email,
		Contributions: w.Contributions,
	}
}

type wireTreeEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	SHA  string `json:"sha"`
	Size int    `json:"size"`
}

type wireTree struct {
	SHA       string          `json:"sha"`
	URL       string          `json:"url"`
	Tree      []wireTreeEntry `json:"tree"`
	Truncated bool            `json:"truncated"`
}

func (w wireTree) toPublic() []TreeEntry {
	out := make([]TreeEntry, len(w.Tree))
	for i, e := range w.Tree {
		out[i] = TreeEntry{
			Path: e.Path,
			Mode: e.Mode,
			Type: e.Type,
			SHA:  e.SHA,
			Size: e.Size,
		}
	}
	return out
}

type wireReadme struct {
	Name        string `json:"name"`
	SHA         string `json:"sha"`
	Size        int    `json:"size"`
	Encoding    string `json:"encoding"`
	Content     string `json:"content"`
	HTMLURL     string `json:"html_url"`
	DownloadURL string `json:"download_url"`
}

func (w wireReadme) toPublic() ReadmeFile {
	decoded := ""
	if w.Content != "" {
		clean := strings.ReplaceAll(w.Content, "\n", "")
		if b, err := base64.StdEncoding.DecodeString(clean); err == nil {
			decoded = string(b)
		}
	}
	return ReadmeFile{
		Name:           w.Name,
		SHA:            w.SHA,
		Size:           w.Size,
		DecodedContent: decoded,
		HTMLURL:        w.HTMLURL,
		DownloadURL:    w.DownloadURL,
	}
}

type wireOrg struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	Name        string `json:"name"`
	HTMLURL     string `json:"html_url"`
	Description string `json:"description"`
	Blog        string `json:"blog"`
	Email       string `json:"email"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
	CreatedAt   string `json:"created_at"`
}

func (w wireOrg) toPublic() OrgProfile {
	return OrgProfile{
		ID:          w.ID,
		Login:       w.Login,
		Name:        w.Name,
		HTMLURL:     w.HTMLURL,
		Description: w.Description,
		Blog:        w.Blog,
		Email:       w.Email,
		PublicRepos: w.PublicRepos,
		Followers:   w.Followers,
		CreatedAt:   w.CreatedAt,
	}
}
