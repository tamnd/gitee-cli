package gitee

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

// domain.go exposes gitee as a kit Domain: a driver that a multi-domain
// host enables with a single blank import,
//
//	import _ "github.com/tamnd/gitee-cli/gitee"
//
// The init below registers it; the host then dereferences gitee:// URIs by
// routing to the operations Register installs.
func init() { kit.Register(Domain{}) }

// Domain is the gitee kit driver. It carries no state.
type Domain struct{}

// Info describes the scheme, the hostnames a pasted link is matched against,
// and the identity reused for the binary's help and version.
func (Domain) Info() kit.DomainInfo {
	return kit.DomainInfo{
		Scheme: "gitee",
		Hosts:  []string{"gitee.com"},
		Identity: kit.Identity{
			Binary: "gitee",
			Short:  "A command line for Gitee.",
			Long: `A command line for Gitee.

gitee reads public Gitee data over plain HTTPS, shapes it into
clean records, and prints output that pipes into the rest of your tools.
No API key, nothing to run alongside it.`,
			Site: "gitee.com",
			Repo: "https://github.com/tamnd/gitee-cli",
		},
	}
}

// Register installs the client factory and every operation onto app.
func (Domain) Register(app *kit.App) {
	app.SetClient(newKitClient)

	// trending: list trending repos from the Gitee explore endpoint.
	kit.Handle(app, kit.OpMeta{
		Name:    "trending",
		Group:   "read",
		List:    true,
		Summary: "List trending Gitee repositories",
		URIType: "repo",
	}, listTrending)

	// search: search repos via the Gitee search API.
	kit.Handle(app, kit.OpMeta{
		Name:    "search",
		Group:   "read",
		List:    true,
		Summary: "Search Gitee repositories",
		URIType: "repo",
		Args:    []kit.Arg{{Name: "query", Help: "search query"}},
	}, listSearch)

	// user: fetch a single user profile.
	kit.Handle(app, kit.OpMeta{
		Name:     "user",
		Group:    "read",
		Single:   true,
		Summary:  "Get a Gitee user profile",
		URIType:  "user",
		Resolver: true,
		Args:     []kit.Arg{{Name: "username", Help: "Gitee username"}},
	}, getUser)

	// repo: fetch a single repository.
	kit.Handle(app, kit.OpMeta{
		Name:     "repo",
		Group:    "read",
		Single:   true,
		Summary:  "Get a Gitee repository",
		URIType:  "repo",
		Resolver: true,
		Args:     []kit.Arg{{Name: "ref", Help: "owner/repo or Gitee URL"}},
	}, getRepo)
}

// newKitClient builds the client from the kit-resolved config.
func newKitClient(_ context.Context, cfg kit.Config) (any, error) {
	c := DefaultConfig()
	if cfg.UserAgent != "" {
		c.UserAgent = cfg.UserAgent
	}
	if cfg.Rate > 0 {
		c.Rate = cfg.Rate
	}
	if cfg.Retries > 0 {
		c.Retries = cfg.Retries
	}
	if cfg.Timeout > 0 {
		c.Timeout = cfg.Timeout
	}
	return NewClient(c), nil
}

// --- inputs ---

type trendingInput struct {
	Lang   string  `kit:"flag" help:"filter by language"`
	Sort   string  `kit:"flag" help:"sort: stars|newest|updated"`
	Limit  int     `kit:"flag,inherit" help:"max results"`
	Client *Client `kit:"inject"`
}

type searchInput struct {
	Query  string  `kit:"arg" help:"search query"`
	Lang   string  `kit:"flag" help:"filter by language"`
	Sort   string  `kit:"flag" help:"sort: stars|forks|updated"`
	Limit  int     `kit:"flag,inherit" help:"max results"`
	Client *Client `kit:"inject"`
}

type userInput struct {
	Username string  `kit:"arg" help:"Gitee username"`
	Client   *Client `kit:"inject"`
}

type repoInput struct {
	Ref    string  `kit:"arg" help:"owner/repo or Gitee URL"`
	Client *Client `kit:"inject"`
}

// --- handlers ---

func listTrending(ctx context.Context, in trendingInput, emit func(*Repo) error) error {
	repos, err := in.Client.TrendingRepos(ctx, in.Lang, in.Sort, in.Limit)
	if err != nil {
		return err
	}
	for i := range repos {
		if err := emit(&repos[i]); err != nil {
			return err
		}
	}
	return nil
}

func listSearch(ctx context.Context, in searchInput, emit func(*Repo) error) error {
	if strings.TrimSpace(in.Query) == "" {
		return errs.Usage("query is required")
	}
	repos, err := in.Client.SearchRepos(ctx, in.Query, in.Lang, in.Sort, in.Limit)
	if err != nil {
		return err
	}
	for i := range repos {
		if err := emit(&repos[i]); err != nil {
			return err
		}
	}
	return nil
}

func getUser(ctx context.Context, in userInput, emit func(*User) error) error {
	if strings.TrimSpace(in.Username) == "" {
		return errs.Usage("username is required")
	}
	u, err := in.Client.GetUser(ctx, in.Username)
	if err != nil {
		return err
	}
	return emit(&u)
}

func getRepo(ctx context.Context, in repoInput, emit func(*Repo) error) error {
	owner, name, err := parseRepoRef(in.Ref)
	if err != nil {
		return errs.Usage("%s", err.Error())
	}
	r, err := in.Client.GetRepo(ctx, owner, name)
	if err != nil {
		return err
	}
	return emit(&r)
}

// --- Resolver: pure string functions, no network ---

// Classify turns a Gitee URL or owner/repo path into the canonical (type, id).
func (Domain) Classify(input string) (uriType, id string, err error) {
	input = strings.TrimSpace(input)
	if u, parseErr := url.Parse(input); parseErr == nil && (u.Scheme == "http" || u.Scheme == "https") {
		// Could be a user or repo URL.
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		switch len(parts) {
		case 1:
			if parts[0] != "" {
				return "user", parts[0], nil
			}
		case 2:
			if parts[0] != "" && parts[1] != "" {
				return "repo", parts[0] + "/" + parts[1], nil
			}
		}
		return "", "", errs.Usage("unrecognized Gitee URL: %q", input)
	}
	// bare owner/repo or username
	parts := strings.Split(strings.Trim(input, "/"), "/")
	switch len(parts) {
	case 1:
		if parts[0] != "" {
			return "user", parts[0], nil
		}
	case 2:
		if parts[0] != "" && parts[1] != "" {
			return "repo", parts[0] + "/" + parts[1], nil
		}
	}
	return "", "", errs.Usage("unrecognized Gitee reference: %q", input)
}

// Locate is the inverse: the live https URL for a (type, id).
func (Domain) Locate(uriType, id string) (string, error) {
	switch uriType {
	case "user":
		return "https://gitee.com/" + id, nil
	case "repo":
		return "https://gitee.com/" + id, nil
	default:
		return "", errs.Usage("gitee has no resource type %q", uriType)
	}
}

// --- helpers ---

// parseRepoRef splits a ref like "owner/repo" or a full URL into owner and repo.
func parseRepoRef(ref string) (owner, name string, err error) {
	ref = strings.TrimSpace(ref)
	if u, parseErr := url.Parse(ref); parseErr == nil && (u.Scheme == "http" || u.Scheme == "https") {
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) >= 2 && parts[0] != "" && parts[1] != "" {
			return parts[0], parts[1], nil
		}
		return "", "", fmt.Errorf("cannot extract owner/repo from URL: %q", ref)
	}
	parts := strings.SplitN(strings.Trim(ref, "/"), "/", 2)
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}
	return "", "", fmt.Errorf("expected owner/repo, got: %q", ref)
}
