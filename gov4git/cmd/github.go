package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/github"
	"github.com/gov4git/lib4git/form"
	"github.com/spf13/cobra"
)

var (
	githubCmd = &cobra.Command{
		Use:   "github",
		Short: "Import and export GitHub issues and pull requests",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	githubImportCmd = &cobra.Command{
		Use:   "import",
		Short: "Import GitHub issues and pull requests",
		Long: `Import GitHub issues and pull requests. Example usage:

	gov4git github import --token=GITHUB_ACCESS_TOKEN --owner=GITHUB_USER_OR_ORG --repo=GITHUB_REPO
`,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			repo := github.GithubRepo{Owner: githubOwner, Name: githubRepo}
			github.SetTokenSource(ctx, repo, github.MakeStaticTokenSource(ctx, githubToken))
			importedIssues := github.ImportIssuesForPrioritization(ctx, repo, setup.Gov)
			fmt.Fprint(os.Stdout, form.SprintJSON(importedIssues))
		},
	}
)

var (
	githubToken string
	githubOwner string
	githubRepo  string
)

func init() {
	githubCmd.AddCommand(githubImportCmd)
	githubImportCmd.Flags().StringVar(&githubToken, "token", "", "GitHub access token")
	githubImportCmd.Flags().StringVar(&githubOwner, "owner", "", "GitHub owner")
	githubImportCmd.Flags().StringVar(&githubRepo, "repo", "", "GitHub repo")
	githubImportCmd.MarkFlagRequired("token")
	githubImportCmd.MarkFlagRequired("owner")
	githubImportCmd.MarkFlagRequired("repo")
}
