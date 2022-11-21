package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
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
					Home: id.HomeAddress{Repo: git.URL(userRepo), Branch: git.Branch(userBranch)},
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

	userPropGetCmd = &cobra.Command{
		Use:   "prop-get",
		Short: "Get user property",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			v := member.GetUserProp[interface{}](
				ctx,
				setup.Community,
				member.User(userName),
				userKey,
			)
			fmt.Fprint(os.Stdout, v)
		},
	}
)

var (
	userName   string
	userRepo   string
	userBranch string
	userKey    string
	userValue  string
)

func init() {
	userCmd.AddCommand(userAddCmd)
	userAddCmd.Flags().StringVar(&userName, "name", "", "user alias within the community")
	userAddCmd.MarkFlagRequired("name")
	userAddCmd.Flags().StringVar(&userRepo, "repo", "", "repo URL of user public identity")
	userAddCmd.MarkFlagRequired("repo")
	userAddCmd.Flags().StringVar(&userBranch, "branch", "", "branch of user public identity")
	userAddCmd.MarkFlagRequired("branch")

	userCmd.AddCommand(userRemoveCmd)
	userRemoveCmd.Flags().StringVar(&userName, "name", "", "user alias within the community")
	userRemoveCmd.MarkFlagRequired("name")

	userCmd.AddCommand(userPropGetCmd)
	userPropGetCmd.Flags().StringVar(&userName, "name", "", "user alias within the community")
	userPropGetCmd.MarkFlagRequired("name")
	userPropGetCmd.Flags().StringVar(&userKey, "key", "", "property key")
	userPropGetCmd.MarkFlagRequired("key")
}
