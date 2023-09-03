package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/vendor4git/github"
	"github.com/spf13/cobra"
)

var (
	vendorCmd = &cobra.Command{
		Use:   "vendor",
		Short: "Hosted git repo provisioning (creation and removal)",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	vendorCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a hosted git repo",
		Long: `Call the GitHub API to create a new repo. Example usage:

	gov4git vendor create --token=GITHUB_ACCESS_TOKEN --owner=GITHUB_USER_OR_ORG --repo=REPO_NAME 

This creates a public repo. Adding the flag --private will result in creating a private repo.
		`,
		Run: func(cmd *cobra.Command, args []string) {
			vendor := github.NewGitHubVendor(ctx, vendorGithubToken)
			repo, err := vendor.CreateRepo(ctx, vendorRepoName, vendorGithubOwner, vendorRepoPrivate)
			must.NoError(ctx, err)
			fmt.Fprint(os.Stdout, form.SprintJSON(repo))
		},
	}

	vendorRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove a hosted git repo",
		Long: `Call the GitHub API to remove a repo. Example uage:

	gov4git vendor remove --owner=GITHUB_USER_OR_ORG --repo=REPO_NAME
`,
		Run: func(cmd *cobra.Command, args []string) {
			vendor := github.NewGitHubVendor(ctx, vendorGithubToken)
			err := vendor.RemoveRepo(ctx, vendorRepoName, vendorGithubOwner)
			must.NoError(ctx, err)
		},
	}
)

var (
	vendorGithubToken string
	vendorGithubOwner string
	vendorRepoName    string
	vendorRepoPrivate bool
)

func init() {
	vendorCmd.AddCommand(vendorCreateCmd)
	vendorCreateCmd.Flags().StringVar(&vendorGithubToken, "token", "", "GitHub access token")
	vendorCreateCmd.Flags().StringVar(&vendorGithubOwner, "owner", "", "GitHub owner")
	vendorCreateCmd.Flags().StringVar(&vendorRepoName, "repo", "", "Repo name")
	vendorCreateCmd.Flags().BoolVar(&vendorRepoPrivate, "private", false, "Make private repo")
	vendorCreateCmd.MarkFlagRequired("token")
	vendorCreateCmd.MarkFlagRequired("owner")
	vendorCreateCmd.MarkFlagRequired("repo")

	vendorCmd.AddCommand(vendorRemoveCmd)
	vendorRemoveCmd.Flags().StringVar(&vendorGithubOwner, "owner", "", "GitHub owner")
	vendorRemoveCmd.Flags().StringVar(&vendorRepoName, "repo", "", "Repo name")
	vendorRemoveCmd.MarkFlagRequired("owner")
	vendorRemoveCmd.MarkFlagRequired("repo")
}
