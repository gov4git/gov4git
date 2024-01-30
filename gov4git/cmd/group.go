package cmd

import (
	"github.com/gov4git/gov4git/v2/gov4git/api"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/spf13/cobra"
)

var (
	groupCmd = &cobra.Command{
		Use:   "group",
		Short: "Manage groups",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	groupAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add group to the community",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke(
				func() {
					LoadConfig()
					member.AddGroup(
						ctx,
						setup.Gov,
						member.Group(groupName),
					)
				},
			)
		},
	}

	groupRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove group from the community",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke(
				func() {
					LoadConfig()
					member.RemoveGroup(
						ctx,
						setup.Gov,
						member.Group(groupName),
					)
				},
			)
		},
	}

	groupListCmd = &cobra.Command{
		Use:   "list", // deprecated in favor of `users`
		Short: "List users in group",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					return member.ListGroupUsers(
						ctx,
						setup.Gov,
						member.Group(groupName),
					)
				},
			)
		},
	}
	groupUsersCmd = &cobra.Command{
		Use:   "users",
		Short: "List users in group",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					return member.ListGroupUsers(
						ctx,
						setup.Gov,
						member.Group(groupName),
					)
				},
			)
		},
	}
)

var (
	groupName string
)

func init() {
	groupCmd.AddCommand(groupAddCmd)
	groupAddCmd.Flags().StringVar(&groupName, "name", "", "group alias within the community")
	groupAddCmd.MarkFlagRequired("name")

	groupCmd.AddCommand(groupRemoveCmd)
	groupRemoveCmd.Flags().StringVar(&groupName, "name", "", "group alias within the community")
	groupRemoveCmd.MarkFlagRequired("name")

	groupCmd.AddCommand(groupListCmd)
	groupListCmd.Flags().StringVar(&groupName, "name", "", "group alias within the community")
	groupListCmd.MarkFlagRequired("name")
}
