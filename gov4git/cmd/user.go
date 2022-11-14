package cmd

import (
	"github.com/gov4git/gov4git/mod/id"
	"github.com/gov4git/gov4git/mod/member"
	"github.com/gov4git/lib4git/git"
	"github.com/spf13/cobra"
)

var (
	userCmd = &cobra.Command{
		Use:   "user",
		Short: "Manage users",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	userAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add user to the community",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			member.AddUser(
				ctx,
				setup.Community,
				member.User(userName),
				member.Account{
					Home: id.PublicAddress{Repo: git.URL(userRepo), Branch: git.Branch(userBranch)},
				},
			)
			// fmt.Fprint(os.Stdout, form.Pretty(chg.Result))
		},
	}

	userRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove user from the community",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			member.RemoveUser(
				ctx,
				setup.Community,
				member.User(userName),
			)
			// fmt.Fprint(os.Stdout, form.Pretty(chg.Result))
		},
	}
)

var (
	userName   string
	userRepo   string
	userBranch string
	// userKey    string
	// userValue  string
)

func init() {
	userCmd.AddCommand(userAddCmd)
	userAddCmd.Flags().StringVar(&userName, "name", "", "user alias within the community")
	userAddCmd.Flags().StringVar(&userRepo, "repo", "", "repo URL of user public identity")
	userAddCmd.Flags().StringVar(&userBranch, "branch", "", "branch of user public identity")

	userCmd.AddCommand(userRemoveCmd)
	userRemoveCmd.Flags().StringVar(&userName, "name", "", "user alias within the community")
}
