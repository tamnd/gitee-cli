package gitee

import (
	"context"

	"github.com/tamnd/any-cli/kit"
)

func registerOps(app *kit.App) {
	registerUser(app)
	registerUserRepos(app)
	registerFollowers(app)
	registerFollowing(app)
	registerRepo(app)
	registerCommits(app)
	registerBranches(app)
	registerTags(app)
	registerReleases(app)
	registerIssues(app)
	registerPulls(app)
	registerReadme(app)
	registerTree(app)
	registerStargazers(app)
	registerForks(app)
	registerContributors(app)
	registerSearchRepos(app)
	registerSearchUsers(app)
	registerOrg(app)
	registerOrgRepos(app)
}

// user -----------------------------------------------------------------------

type userIn struct {
	Session  *Session `kit:"inject"`
	Username string   `kit:"arg"`
	Quiet    bool     `kit:"flag,inherit"`
}

func registerUser(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "user",
		Group:   "read",
		Single:  true,
		Summary: "Fetch a Gitee user profile",
		Args:    []kit.Arg{{Name: "username", Help: "Gitee login handle"}},
	}, func(ctx context.Context, in userIn, emit func(User) error) error {
		in.Session.Quiet = in.Quiet
		in.Session.Progressf("fetching user %s", in.Username)
		u, err := in.Session.Client.GetUser(ctx, in.Username)
		if err != nil {
			return MapErr(err)
		}
		return emit(u)
	})
}

// user repos -----------------------------------------------------------------

type userReposIn struct {
	Session   *Session `kit:"inject"`
	Username  string   `kit:"arg"`
	Sort      string   `kit:"flag"`
	Direction string   `kit:"flag"`
	Type      string   `kit:"flag"`
	Limit     int      `kit:"flag,inherit"`
	Quiet     bool     `kit:"flag,inherit"`
}

func registerUserRepos(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "repos",
		Parent:  "user",
		Group:   "read",
		List:    true,
		Summary: "List a user's public repositories",
		Args:    []kit.Arg{{Name: "username", Help: "Gitee login handle"}},
	}, func(ctx context.Context, in userReposIn, emit func(Repo) error) error {
		in.Session.Quiet = in.Quiet
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		in.Session.Progressf("fetching repos for %s (limit %d)", in.Username, limit)
		repos, err := in.Session.Client.UserRepos(ctx, in.Username, in.Sort, in.Direction, in.Type, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(repos, emit)
	})
}

// followers ------------------------------------------------------------------

type followersIn struct {
	Session  *Session `kit:"inject"`
	Username string   `kit:"arg"`
	Limit    int      `kit:"flag,inherit"`
	Quiet    bool     `kit:"flag,inherit"`
}

func registerFollowers(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "followers",
		Group:   "read",
		List:    true,
		Summary: "List a user's followers",
		Args:    []kit.Arg{{Name: "username", Help: "Gitee login handle"}},
	}, func(ctx context.Context, in followersIn, emit func(User) error) error {
		in.Session.Quiet = in.Quiet
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		in.Session.Progressf("fetching followers for %s", in.Username)
		users, err := in.Session.Client.Followers(ctx, in.Username, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(users, emit)
	})
}

// following ------------------------------------------------------------------

type followingIn struct {
	Session  *Session `kit:"inject"`
	Username string   `kit:"arg"`
	Limit    int      `kit:"flag,inherit"`
	Quiet    bool     `kit:"flag,inherit"`
}

func registerFollowing(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "following",
		Group:   "read",
		List:    true,
		Summary: "List users that a user follows",
		Args:    []kit.Arg{{Name: "username", Help: "Gitee login handle"}},
	}, func(ctx context.Context, in followingIn, emit func(User) error) error {
		in.Session.Quiet = in.Quiet
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		in.Session.Progressf("fetching following for %s", in.Username)
		users, err := in.Session.Client.Following(ctx, in.Username, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(users, emit)
	})
}

// repo -----------------------------------------------------------------------

type repoIn struct {
	Session *Session `kit:"inject"`
	Slug    string   `kit:"arg"`
	Quiet   bool     `kit:"flag,inherit"`
}

func registerRepo(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "repo",
		Group:   "read",
		Single:  true,
		Summary: "Fetch a repository",
		Args:    []kit.Arg{{Name: "repo", Help: "owner/repo slug or gitee.com URL"}},
	}, func(ctx context.Context, in repoIn, emit func(Repo) error) error {
		in.Session.Quiet = in.Quiet
		slug, err := ParseRepoSlug(in.Slug)
		if err != nil {
			return err
		}
		in.Session.Progressf("fetching %s", slug)
		r, err := in.Session.Client.GetRepo(ctx, slug.Owner, slug.Repo)
		if err != nil {
			return MapErr(err)
		}
		return emit(r)
	})
}

// commits --------------------------------------------------------------------

type commitsIn struct {
	Session   *Session `kit:"inject"`
	Slug      string   `kit:"arg"`
	SHA       string   `kit:"flag"`
	Path      string   `kit:"flag"`
	Author    string   `kit:"flag"`
	Since     string   `kit:"flag"`
	Until     string   `kit:"flag"`
	Limit     int      `kit:"flag,inherit"`
	Quiet     bool     `kit:"flag,inherit"`
}

func registerCommits(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "commits",
		Group:   "read",
		List:    true,
		Summary: "List commits for a repository",
		Args:    []kit.Arg{{Name: "repo", Help: "owner/repo slug"}},
	}, func(ctx context.Context, in commitsIn, emit func(Commit) error) error {
		in.Session.Quiet = in.Quiet
		slug, err := ParseRepoSlug(in.Slug)
		if err != nil {
			return err
		}
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		in.Session.Progressf("fetching commits for %s", slug)
		commits, err := in.Session.Client.Commits(ctx, slug.Owner, slug.Repo, in.SHA, in.Path, in.Author, in.Since, in.Until, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(commits, emit)
	})
}

// branches -------------------------------------------------------------------

type branchesIn struct {
	Session *Session `kit:"inject"`
	Slug    string   `kit:"arg"`
	Limit   int      `kit:"flag,inherit"`
	Quiet   bool     `kit:"flag,inherit"`
}

func registerBranches(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "branches",
		Group:   "read",
		List:    true,
		Summary: "List branches for a repository",
		Args:    []kit.Arg{{Name: "repo", Help: "owner/repo slug"}},
	}, func(ctx context.Context, in branchesIn, emit func(Branch) error) error {
		in.Session.Quiet = in.Quiet
		slug, err := ParseRepoSlug(in.Slug)
		if err != nil {
			return err
		}
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		in.Session.Progressf("fetching branches for %s", slug)
		branches, err := in.Session.Client.Branches(ctx, slug.Owner, slug.Repo, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(branches, emit)
	})
}

// tags -----------------------------------------------------------------------

type tagsIn struct {
	Session   *Session `kit:"inject"`
	Slug      string   `kit:"arg"`
	Sort      string   `kit:"flag"`
	Direction string   `kit:"flag"`
	Limit     int      `kit:"flag,inherit"`
	Quiet     bool     `kit:"flag,inherit"`
}

func registerTags(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "tags",
		Group:   "read",
		List:    true,
		Summary: "List tags for a repository",
		Args:    []kit.Arg{{Name: "repo", Help: "owner/repo slug"}},
	}, func(ctx context.Context, in tagsIn, emit func(Tag) error) error {
		in.Session.Quiet = in.Quiet
		slug, err := ParseRepoSlug(in.Slug)
		if err != nil {
			return err
		}
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		in.Session.Progressf("fetching tags for %s", slug)
		tags, err := in.Session.Client.Tags(ctx, slug.Owner, slug.Repo, in.Sort, in.Direction, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(tags, emit)
	})
}

// releases -------------------------------------------------------------------

type releasesIn struct {
	Session *Session `kit:"inject"`
	Slug    string   `kit:"arg"`
	Limit   int      `kit:"flag,inherit"`
	Quiet   bool     `kit:"flag,inherit"`
}

func registerReleases(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "releases",
		Group:   "read",
		List:    true,
		Summary: "List releases for a repository",
		Args:    []kit.Arg{{Name: "repo", Help: "owner/repo slug"}},
	}, func(ctx context.Context, in releasesIn, emit func(Release) error) error {
		in.Session.Quiet = in.Quiet
		slug, err := ParseRepoSlug(in.Slug)
		if err != nil {
			return err
		}
		limit := in.Limit
		if limit <= 0 {
			limit = 10
		}
		in.Session.Progressf("fetching releases for %s", slug)
		releases, err := in.Session.Client.Releases(ctx, slug.Owner, slug.Repo, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(releases, emit)
	})
}

// issues ---------------------------------------------------------------------

type issuesIn struct {
	Session   *Session `kit:"inject"`
	Slug      string   `kit:"arg"`
	State     string   `kit:"flag"`
	Sort      string   `kit:"flag"`
	Direction string   `kit:"flag"`
	Limit     int      `kit:"flag,inherit"`
	Quiet     bool     `kit:"flag,inherit"`
}

func registerIssues(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "issues",
		Group:   "read",
		List:    true,
		Summary: "List issues for a repository",
		Args:    []kit.Arg{{Name: "repo", Help: "owner/repo slug"}},
	}, func(ctx context.Context, in issuesIn, emit func(Issue) error) error {
		in.Session.Quiet = in.Quiet
		slug, err := ParseRepoSlug(in.Slug)
		if err != nil {
			return err
		}
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		state := in.State
		if state == "" {
			state = "open"
		}
		sort := in.Sort
		if sort == "" {
			sort = "created"
		}
		direction := in.Direction
		if direction == "" {
			direction = "desc"
		}
		in.Session.Progressf("fetching issues for %s", slug)
		issues, err := in.Session.Client.Issues(ctx, slug.Owner, slug.Repo, state, sort, direction, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(issues, emit)
	})
}

// pulls ----------------------------------------------------------------------

type pullsIn struct {
	Session   *Session `kit:"inject"`
	Slug      string   `kit:"arg"`
	State     string   `kit:"flag"`
	Sort      string   `kit:"flag"`
	Direction string   `kit:"flag"`
	Limit     int      `kit:"flag,inherit"`
	Quiet     bool     `kit:"flag,inherit"`
}

func registerPulls(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "pulls",
		Group:   "read",
		List:    true,
		Summary: "List pull requests for a repository",
		Args:    []kit.Arg{{Name: "repo", Help: "owner/repo slug"}},
	}, func(ctx context.Context, in pullsIn, emit func(PullRequest) error) error {
		in.Session.Quiet = in.Quiet
		slug, err := ParseRepoSlug(in.Slug)
		if err != nil {
			return err
		}
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		state := in.State
		if state == "" {
			state = "open"
		}
		sort := in.Sort
		if sort == "" {
			sort = "created"
		}
		direction := in.Direction
		if direction == "" {
			direction = "desc"
		}
		in.Session.Progressf("fetching pull requests for %s", slug)
		pulls, err := in.Session.Client.Pulls(ctx, slug.Owner, slug.Repo, state, sort, direction, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(pulls, emit)
	})
}

// readme ---------------------------------------------------------------------

type readmeIn struct {
	Session *Session `kit:"inject"`
	Slug    string   `kit:"arg"`
	Ref     string   `kit:"flag"`
	Quiet   bool     `kit:"flag,inherit"`
}

func registerReadme(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "readme",
		Group:   "read",
		Single:  true,
		Summary: "Fetch the README for a repository",
		Args:    []kit.Arg{{Name: "repo", Help: "owner/repo slug"}},
	}, func(ctx context.Context, in readmeIn, emit func(ReadmeFile) error) error {
		in.Session.Quiet = in.Quiet
		slug, err := ParseRepoSlug(in.Slug)
		if err != nil {
			return err
		}
		in.Session.Progressf("fetching README for %s", slug)
		readme, err := in.Session.Client.Readme(ctx, slug.Owner, slug.Repo, in.Ref)
		if err != nil {
			return MapErr(err)
		}
		return emit(readme)
	})
}

// tree -----------------------------------------------------------------------

type treeIn struct {
	Session   *Session `kit:"inject"`
	Slug      string   `kit:"arg"`
	Ref       string   `kit:"flag"`
	Recursive bool     `kit:"flag"`
	Quiet     bool     `kit:"flag,inherit"`
}

func registerTree(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "tree",
		Group:   "read",
		List:    true,
		Summary: "List the git tree for a repository",
		Args:    []kit.Arg{{Name: "repo", Help: "owner/repo slug"}},
	}, func(ctx context.Context, in treeIn, emit func(TreeEntry) error) error {
		in.Session.Quiet = in.Quiet
		slug, err := ParseRepoSlug(in.Slug)
		if err != nil {
			return err
		}
		ref := in.Ref
		if ref == "" {
			ref = "HEAD"
		}
		in.Session.Progressf("fetching tree for %s @ %s", slug, ref)
		entries, err := in.Session.Client.Tree(ctx, slug.Owner, slug.Repo, ref, in.Recursive)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(entries, emit)
	})
}

// stargazers -----------------------------------------------------------------

type stargazersIn struct {
	Session *Session `kit:"inject"`
	Slug    string   `kit:"arg"`
	Limit   int      `kit:"flag,inherit"`
	Quiet   bool     `kit:"flag,inherit"`
}

func registerStargazers(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "stargazers",
		Group:   "read",
		List:    true,
		Summary: "List users who starred a repository",
		Args:    []kit.Arg{{Name: "repo", Help: "owner/repo slug"}},
	}, func(ctx context.Context, in stargazersIn, emit func(User) error) error {
		in.Session.Quiet = in.Quiet
		slug, err := ParseRepoSlug(in.Slug)
		if err != nil {
			return err
		}
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		in.Session.Progressf("fetching stargazers for %s", slug)
		users, err := in.Session.Client.Stargazers(ctx, slug.Owner, slug.Repo, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(users, emit)
	})
}

// forks ----------------------------------------------------------------------

type forksIn struct {
	Session *Session `kit:"inject"`
	Slug    string   `kit:"arg"`
	Sort    string   `kit:"flag"`
	Limit   int      `kit:"flag,inherit"`
	Quiet   bool     `kit:"flag,inherit"`
}

func registerForks(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "forks",
		Group:   "read",
		List:    true,
		Summary: "List forks of a repository",
		Args:    []kit.Arg{{Name: "repo", Help: "owner/repo slug"}},
	}, func(ctx context.Context, in forksIn, emit func(Repo) error) error {
		in.Session.Quiet = in.Quiet
		slug, err := ParseRepoSlug(in.Slug)
		if err != nil {
			return err
		}
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		in.Session.Progressf("fetching forks for %s", slug)
		repos, err := in.Session.Client.Forks(ctx, slug.Owner, slug.Repo, in.Sort, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(repos, emit)
	})
}

// contributors ---------------------------------------------------------------

type contributorsIn struct {
	Session *Session `kit:"inject"`
	Slug    string   `kit:"arg"`
	Type    string   `kit:"flag"`
	Limit   int      `kit:"flag,inherit"`
	Quiet   bool     `kit:"flag,inherit"`
}

func registerContributors(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "contributors",
		Group:   "read",
		List:    true,
		Summary: "List contributors to a repository",
		Args:    []kit.Arg{{Name: "repo", Help: "owner/repo slug"}},
	}, func(ctx context.Context, in contributorsIn, emit func(Contributor) error) error {
		in.Session.Quiet = in.Quiet
		slug, err := ParseRepoSlug(in.Slug)
		if err != nil {
			return err
		}
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		in.Session.Progressf("fetching contributors for %s", slug)
		contributors, err := in.Session.Client.Contributors(ctx, slug.Owner, slug.Repo, in.Type, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(contributors, emit)
	})
}

// search repos ---------------------------------------------------------------

type searchReposIn struct {
	Session *Session `kit:"inject"`
	Query   string   `kit:"arg"`
	Sort    string   `kit:"flag"`
	Limit   int      `kit:"flag,inherit"`
	Quiet   bool     `kit:"flag,inherit"`
}

func registerSearchRepos(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "repos",
		Parent:  "search",
		Group:   "read",
		List:    true,
		Summary: "Search Gitee repositories",
		Args:    []kit.Arg{{Name: "query", Help: "search keywords"}},
	}, func(ctx context.Context, in searchReposIn, emit func(Repo) error) error {
		in.Session.Quiet = in.Quiet
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		in.Session.Progressf("searching repos for %q", in.Query)
		repos, err := in.Session.Client.SearchRepos(ctx, in.Query, in.Sort, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(repos, emit)
	})
}

// search users ---------------------------------------------------------------

type searchUsersIn struct {
	Session *Session `kit:"inject"`
	Query   string   `kit:"arg"`
	Limit   int      `kit:"flag,inherit"`
	Quiet   bool     `kit:"flag,inherit"`
}

func registerSearchUsers(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "users",
		Parent:  "search",
		Group:   "read",
		List:    true,
		Summary: "Search Gitee users",
		Args:    []kit.Arg{{Name: "query", Help: "search keywords"}},
	}, func(ctx context.Context, in searchUsersIn, emit func(User) error) error {
		in.Session.Quiet = in.Quiet
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		in.Session.Progressf("searching users for %q", in.Query)
		users, err := in.Session.Client.SearchUsers(ctx, in.Query, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(users, emit)
	})
}

// org ------------------------------------------------------------------------

type orgIn struct {
	Session *Session `kit:"inject"`
	Name    string   `kit:"arg"`
	Quiet   bool     `kit:"flag,inherit"`
}

func registerOrg(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "org",
		Group:   "read",
		Single:  true,
		Summary: "Fetch a Gitee organization profile",
		Args:    []kit.Arg{{Name: "name", Help: "organization login name"}},
	}, func(ctx context.Context, in orgIn, emit func(OrgProfile) error) error {
		in.Session.Quiet = in.Quiet
		in.Session.Progressf("fetching org %s", in.Name)
		org, err := in.Session.Client.GetOrg(ctx, in.Name)
		if err != nil {
			return MapErr(err)
		}
		return emit(org)
	})
}

// org repos ------------------------------------------------------------------

type orgReposIn struct {
	Session   *Session `kit:"inject"`
	Name      string   `kit:"arg"`
	Type      string   `kit:"flag"`
	Sort      string   `kit:"flag"`
	Direction string   `kit:"flag"`
	Limit     int      `kit:"flag,inherit"`
	Quiet     bool     `kit:"flag,inherit"`
}

func registerOrgRepos(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "repos",
		Parent:  "org",
		Group:   "read",
		List:    true,
		Summary: "List repositories for an organization",
		Args:    []kit.Arg{{Name: "name", Help: "organization login name"}},
	}, func(ctx context.Context, in orgReposIn, emit func(Repo) error) error {
		in.Session.Quiet = in.Quiet
		limit := in.Limit
		if limit <= 0 {
			limit = 20
		}
		in.Session.Progressf("fetching repos for org %s", in.Name)
		repos, err := in.Session.Client.OrgRepos(ctx, in.Name, in.Type, in.Sort, in.Direction, limit)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(repos, emit)
	})
}

// helpers --------------------------------------------------------------------

func emitAll[T any](items []T, emit func(T) error) error {
	for _, item := range items {
		if err := emit(item); err != nil {
			return err
		}
	}
	return nil
}
