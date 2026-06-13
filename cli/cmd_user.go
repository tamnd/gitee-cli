package cli

import (
	"github.com/spf13/cobra"
	"github.com/tamnd/gitee-cli/gitee"
)

func (a *App) userCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "user <username>",
		Short: "Show a Gitee user profile and their repos",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]
			a.progressf("fetching user %q...", username)
			user, err := a.client.GetUser(cmd.Context(), username)
			if err != nil {
				return mapFetchErr(err)
			}
			if err := a.render([]gitee.User{user}); err != nil {
				return err
			}
			n := a.effectiveLimit(20)
			a.progressf("fetching repos for %q (limit %d)...", username, n)
			repos, err := a.client.UserRepos(cmd.Context(), username, n)
			if err != nil {
				return mapFetchErr(err)
			}
			if len(repos) > 0 {
				return a.render(repos)
			}
			return nil
		},
	}
}
