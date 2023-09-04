package cmd

import (
	"fmt"
	"os"

	govgh "github.com/gov4git/gov4git/github"
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

	gov4git github import --token=GITHUB_ACCESS_TOKEN --project=PROJECT_OWNER/PROJECT_REPO

You must be the organizer of the community to run this command. In particular, both public and private repos of
the community must be present in your local config file, as well as their respective access tokens.
`,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			repo := govgh.ParseGithubRepo(ctx, githubProject)
			govgh.SetTokenSource(ctx, repo, govgh.MakeStaticTokenSource(ctx, githubToken))
			importedIssues := govgh.ImportIssuesForPrioritization(ctx, repo, nil, setup.Organizer)
			fmt.Fprint(os.Stdout, form.SprintJSON(importedIssues))
		},
	}

	githubDeployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy governance for a GitHub project repo",
		Long: `Example usage:

	gov4git github deploy \
		--token=GITHUB_ACCESS_TOKEN \
		--project=PROJECT_OWNER/PROJECT_REPO \
		--gov=GOV_OWNER/GOV_REPO_PREFIX

--token is a GitHub access token which has read access to the project repo's issues and pull requests; and
create and write access to the governance repos.

--project is the GitHub project_owner/project_repo of the project repository to be governed.

--gov is the GitHub owner and repo name prefix (in the form owner/repo_prefix) of the public and private
governance repositories to be created. Their names will be repo_prefix:gov.public and repo_prefix:gov.private,
respectively. If --gov is not specified, their names will default to project_repo:gov.public and
project_repo:gov.private, respectively.

Therefore, aside for debugging purposes, users should deploy with:

	gov4git github deploy --token=GITHUB_ACCESS_TOKEN --project=PROJECT_OWNER/PROJECT_REPO

`,
		Run: func(cmd *cobra.Command, args []string) {
			project := govgh.ParseGithubRepo(ctx, githubProject)
			var govPrefix govgh.GithubRepo
			if githubGov == "" {
				govPrefix = project
			} else {
				govPrefix = govgh.ParseGithubRepo(ctx, githubGov)
			}

			// deploy governance on GitHub (by way of placing GitHub actions in the public governance repo)
			config := govgh.Deploy(ctx, githubToken, project, govPrefix)
			fmt.Fprint(os.Stdout, form.SprintJSON(config))
		},
	}
)

var (
	githubToken   string
	githubProject string
	githubGov     string
)

func init() {
	githubCmd.AddCommand(githubImportCmd)
	githubImportCmd.Flags().StringVar(&githubToken, "token", "", "GitHub access token")
	githubImportCmd.Flags().StringVar(&githubProject, "project", "", "GitHub project owner/repo")
	githubImportCmd.MarkFlagRequired("token")
	githubImportCmd.MarkFlagRequired("project")

	githubCmd.AddCommand(githubDeployCmd)
	githubDeployCmd.Flags().StringVar(&githubToken, "token", "", "GitHub access token")
	githubDeployCmd.Flags().StringVar(&githubProject, "project", "", "GitHub project owner/repo")
	githubDeployCmd.Flags().StringVar(&githubGov, "gov", "", "governance Github owner/repo_prefix")
	githubDeployCmd.MarkFlagRequired("token")
	githubDeployCmd.MarkFlagRequired("project")
	githubDeployCmd.MarkFlagRequired("gov")
}
