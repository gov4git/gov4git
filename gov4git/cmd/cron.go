package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/gov4git/gov4git/v2/github"
	govgh "github.com/gov4git/gov4git/v2/github"
	"github.com/gov4git/gov4git/v2/proto/cron"
	"github.com/gov4git/lib4git/form"
	"github.com/spf13/cobra"
)

var (
	cronCmd = &cobra.Command{
		Use:   "cron",
		Short: "cron performs time-dependent update operations to the governance system",
		Long: `
This command is intended as a target for a cronjob which runs every couple of minutes.
It will ensure that:
- Governance is synchronized with the issues and pull requests of a GitHub project at a configurable frequency, and
- Votes from community members are incorporated in governance ballots at a configurable frequency.
`,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			repo := govgh.ParseRepo(ctx, githubProject)
			govgh.SetTokenSource(ctx, repo, govgh.MakeStaticTokenSource(ctx, githubToken))
			ghc := govgh.GetGithubClient(ctx, repo)
			result := cron.Cron(
				ctx,
				repo,
				ghc,
				setup.Organizer,
				time.Duration(cronGithubFreqSeconds)*time.Second,
				time.Duration(cronCommunityFreqSeconds)*time.Second,
				syncFetchPar,
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(result))
		},
	}
)

var (
	cronGithubFreqSeconds    int
	cronCommunityFreqSeconds int
)

func init() {
	cronCmd.Flags().StringVar(&githubProject, "project", "", "GitHub project owner/repo")
	cronCmd.Flags().StringVar(&githubToken, "token", "", "GitHub access token")
	cronCmd.Flags().IntVar(&cronGithubFreqSeconds, "github_freq", github.DefaultGithubFreq, "frequency of GitHub import, in seconds")
	cronCmd.Flags().IntVar(&cronCommunityFreqSeconds, "community_freq", github.DefaultCommunityFreq, "frequency of community tallies, in seconds")
	cronCmd.Flags().IntVar(&syncFetchPar, "fetch_par", github.DefaultFetchParallelism, "parallelism while clonging member repos for vote collection")

	cronCmd.MarkFlagRequired("project")
	cronCmd.MarkFlagRequired("token")
	cronCmd.MarkFlagRequired("github_freq")
	cronCmd.MarkFlagRequired("community_freq")
	cronCmd.MarkFlagRequired("fetch_par")
}
