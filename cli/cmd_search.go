package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) searchCmd() *cobra.Command {
	var (
		lang string
		sort string
	)
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search Gitee repositories",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n := a.effectiveLimit(20)
			a.progressf("searching for %q...", args[0])
			repos, err := a.client.SearchRepos(cmd.Context(), args[0], lang, sort, n)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(repos, len(repos))
		},
	}
	cmd.Flags().StringVar(&lang, "lang", "", "filter by programming language")
	cmd.Flags().StringVar(&sort, "sort", "stars", "sort order: stars|forks|updated")
	return cmd
}
