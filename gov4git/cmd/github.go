package cmd

import (
	"fmt"
	"os"

	govgh "github.com/gov4git/gov4git/github"
	"github.com/gov4git/gov4git/github/deploy/tools"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/vendor4git/github"
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
		Short: "Import GitHub artifacts for governance",
		Long: `Import GitHub issues and pull requests. Example usage:

	gov4git github import --token=GITHUB_ACCESS_TOKEN --project=PROJECT_OWNER/PROJECT_REPO

You must be the organizer of the community to run this command. In particular, both public and private repos of
the community must be present in your local config file, as well as their respective access tokens.
`,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			repo := govgh.ParseRepo(ctx, githubProject)
			govgh.SetTokenSource(ctx, repo, govgh.MakeStaticTokenSource(ctx, githubToken))
			ghc := govgh.GetGithubClient(ctx, repo)
			importedIssues := govgh.Import(ctx, repo, ghc, setup.Organizer)
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
		--release=GOV4GIT_RELEASE \
		--gov=GOV_OWNER/GOV_REPO_PREFIX

--token is a GitHub access token which has read access to the project repo's issues and pull requests; and
create and write access to the governance repos.

--project is the GitHub project_owner/project_repo of the project repository to be governed.

--release specifies the GitHub release of gov4git to use for automation.

--gov is the GitHub owner and repo name prefix (in the form owner/repo_prefix) of the public and private
governance repositories to be created. Their names will be repo_prefix:gov.public and repo_prefix:gov.private,
respectively. If --gov is not specified, their names will default to project_repo:gov.public and
project_repo:gov.private, respectively.

Therefore, aside for debugging purposes, users should deploy with:

	gov4git github deploy \
		--token=GITHUB_ACCESS_TOKEN \
		--project=PROJECT_OWNER/PROJECT_REPO \
		--release=GOV4GIT_RELEASE

`,
		Run: func(cmd *cobra.Command, args []string) {
			must.Assertf(ctx, githubRelease != "", "github release must be specified")

			project := govgh.ParseRepo(ctx, githubProject)

			var govPrefix govgh.Repo
			if githubGov == "" {
				govPrefix = project
			} else {
				govPrefix = govgh.ParseRepo(ctx, githubGov)
			}

			// deploy governance on GitHub (by way of placing GitHub actions in the public governance repo)
			config := govgh.Deploy(ctx, githubToken, project, govPrefix, githubRelease)
			fmt.Fprint(os.Stdout, form.SprintJSON(config))
		},
	}

	githubCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a git repo hosted on GitHub",
		Long: `Call the GitHub API to create a new repo. Example usage:

	gov4git github create --token=GITHUB_ACCESS_TOKEN --repo=GITHUB_OWNER/GITHUB_REPO

This creates a public repo. Adding the flag --private will result in creating a private repo.
		`,
		Run: func(cmd *cobra.Command, args []string) {
			ghRepo := govgh.ParseRepo(ctx, githubRepo)
			vendor := github.NewGitHubVendor(ctx, githubToken)
			repo, err := vendor.CreateRepo(ctx, ghRepo.Name, ghRepo.Owner, githubPrivate)
			must.NoError(ctx, err)
			fmt.Fprint(os.Stdout, form.SprintJSON(repo))
		},
	}

	githubRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove a git repo hosted on GitHub",
		Long: `Call the GitHub API to remove a repo. Example uage:

	gov4git github remove --token=GITHUB_ACCESS_TOKEN --repo=GITHUB_OWNER/GITHUB_REPO
`,
		Run: func(cmd *cobra.Command, args []string) {
			ghRepo := govgh.ParseRepo(ctx, githubRepo)
			vendor := github.NewGitHubVendor(ctx, githubToken)
			err := vendor.RemoveRepo(ctx, ghRepo.Name, ghRepo.Owner)
			must.NoError(ctx, err)
		},
	}

	githubClearCommentsCmd = &cobra.Command{
		Use:   "clear-comments",
		Short: "Delete all comments from an issue or PR",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ghRepo := govgh.ParseRepo(ctx, githubRepo)
			tools.ClearComments(
				ctx,
				githubToken,
				ghRepo,
				githubIssueNo,
			)
		},
	}
)

var (
	githubToken   string
	githubProject string
	githubRelease string
	githubGov     string
	githubRepo    string
	githubPrivate bool
	githubIssueNo int64
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
	githubDeployCmd.Flags().StringVar(&githubRelease, "release", "", "GitHub release of gov4git to use for automation")
	githubDeployCmd.Flags().StringVar(&githubGov, "gov", "", "governance Github owner/repo_prefix")
	githubDeployCmd.MarkFlagRequired("token")
	githubDeployCmd.MarkFlagRequired("project")
	githubDeployCmd.MarkFlagRequired("release")

	githubCmd.AddCommand(githubCreateCmd)
	githubCreateCmd.Flags().StringVar(&githubToken, "token", "", "GitHub access token")
	githubCreateCmd.Flags().StringVar(&githubRepo, "repo", "", "GitHub owner/repo")
	githubCreateCmd.Flags().BoolVar(&githubPrivate, "private", false, "Make private repo")
	githubCreateCmd.MarkFlagRequired("token")
	githubCreateCmd.MarkFlagRequired("repo")

	githubCmd.AddCommand(githubRemoveCmd)
	githubRemoveCmd.Flags().StringVar(&githubToken, "token", "", "GitHub access token")
	githubRemoveCmd.Flags().StringVar(&githubRepo, "repo", "", "GitHub owner/repo")
	githubRemoveCmd.MarkFlagRequired("token")
	githubRemoveCmd.MarkFlagRequired("repo")

	githubCmd.AddCommand(githubClearCommentsCmd)
	githubClearCommentsCmd.Flags().StringVar(&githubToken, "token", "", "GitHub access token")
	githubClearCommentsCmd.Flags().StringVar(&githubRepo, "repo", "", "GitHub owner/repo")
	githubClearCommentsCmd.Flags().Int64Var(&githubIssueNo, "issue", -1, "GitHub issue or PR number (0 means all)")
	githubClearCommentsCmd.MarkFlagRequired("token")
	githubClearCommentsCmd.MarkFlagRequired("project")
	githubClearCommentsCmd.MarkFlagRequired("issue")

}
