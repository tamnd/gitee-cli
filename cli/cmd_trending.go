package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) trendingCmd() *cobra.Command {
	var (
		lang string
		sort string
	)
	cmd := &cobra.Command{
		Use:   "trending",
		Short: "Explore trending Gitee repositories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			n := a.effectiveLimit(20)
			a.progressf("fetching trending repos (sort=%s)...", sort)
			repos, err := a.client.TrendingRepos(cmd.Context(), lang, sort, n)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(repos, len(repos))
		},
	}
	cmd.Flags().StringVar(&lang, "lang", "", "filter by programming language")
	cmd.Flags().StringVar(&sort, "sort", "stars", "sort order: stars|newest|updated")
	return cmd
}
