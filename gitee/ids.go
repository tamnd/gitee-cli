package gitee

import (
	"fmt"
	"strings"

	"github.com/tamnd/any-cli/kit/errs"
)

type RepoSlug struct {
	Owner string
	Repo  string
}

func ParseRepoSlug(s string) (RepoSlug, error) {
	// Strip common URL prefixes
	s = strings.TrimPrefix(s, "https://gitee.com/")
	s = strings.TrimPrefix(s, "http://gitee.com/")
	s = strings.TrimSuffix(s, "/")
	s = strings.TrimSpace(s)

	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return RepoSlug{}, errs.Usage("invalid repo %q: expected owner/repo", s)
	}
	return RepoSlug{Owner: parts[0], Repo: parts[1]}, nil
}

func (r RepoSlug) String() string { return fmt.Sprintf("%s/%s", r.Owner, r.Repo) }
