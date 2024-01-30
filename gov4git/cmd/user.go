package cmd

import (
	"github.com/gov4git/gov4git/v2/gov4git/api"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/gov4git/v2/proto/member"
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
			api.Invoke(
				func() {
					LoadConfig()
					member.AddUserByPublicAddress(
						ctx,
						setup.Gov,
						member.User(userName),
						id.PublicAddress{Repo: git.URL(userRepo), Branch: git.Branch(userBranch)},
					)
				},
			)
		},
	}

	userRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove user from the community",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke(
				func() {
					LoadConfig()
					member.RemoveUser(
						ctx,
						setup.Gov,
						member.User(userName),
					)
				},
			)
		},
	}

	userPropGetCmd = &cobra.Command{
		Use:   "prop-get",
		Short: "Get user property",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					return member.GetUserProp[interface{}](
						ctx,
						setup.Gov,
						member.User(userName),
						userKey,
					)
				},
			)
		},
	}

	userListGroupsCmd = &cobra.Command{
		Use:   "groups",
		Short: "List user's group memberships",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					return member.ListUserGroups(
						ctx,
						setup.Gov,
						member.User(userName),
					)
				},
			)
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
	userAddCmd.Flags().StringVar(&userRepo, "repo", "", "URL of user's public repo")
	userAddCmd.MarkFlagRequired("repo")
	userAddCmd.Flags().StringVar(&userBranch, "branch", "", "branch in user's public repo")
	userAddCmd.MarkFlagRequired("branch")

	userCmd.AddCommand(userRemoveCmd)
	userRemoveCmd.Flags().StringVar(&userName, "name", "", "user alias within the community")
	userRemoveCmd.MarkFlagRequired("name")

	userCmd.AddCommand(userPropGetCmd)
	userPropGetCmd.Flags().StringVar(&userName, "name", "", "user alias within the community")
	userPropGetCmd.MarkFlagRequired("name")
	userPropGetCmd.Flags().StringVar(&userKey, "key", "", "property key")
	userPropGetCmd.MarkFlagRequired("key")

	userCmd.AddCommand(userListGroupsCmd)
	userListGroupsCmd.Flags().StringVar(&userName, "name", "", "user alias within the community")
	userListGroupsCmd.MarkFlagRequired("name")
}
