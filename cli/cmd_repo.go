package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func (a *App) repoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "repo <owner/repo>",
		Short: "Show single repo details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			owner, repo, err := parseOwnerRepo(args[0])
			if err != nil {
				return codeError(exitUsage, err)
			}
			a.progressf("fetching %s/%s...", owner, repo)
			r, err := a.client.GetRepo(cmd.Context(), owner, repo)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.render(r)
		},
	}
}

// parseOwnerRepo accepts "owner/repo" or a full gitee.com URL.
func parseOwnerRepo(s string) (owner, repo string, err error) {
	s = strings.TrimPrefix(s, "https://gitee.com/")
	s = strings.TrimPrefix(s, "http://gitee.com/")
	s = strings.TrimSuffix(s, "/")
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repo %q: want owner/repo", s)
	}
	return parts[0], parts[1], nil
}
