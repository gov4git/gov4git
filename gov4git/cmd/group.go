package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
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
			member.AddGroup(
				ctx,
				setup.Gov,
				member.Group(groupName),
			)
			// fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	groupRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove group from the community",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			member.RemoveGroup(
				ctx,
				setup.Gov,
				member.Group(groupName),
			)
			// fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	groupListCmd = &cobra.Command{
		Use:   "list", // deprecated in favor of `users`
		Short: "List users in group",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			l := member.ListGroupUsers(
				ctx,
				setup.Gov,
				member.Group(groupName),
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(l))
		},
	}
	groupUsersCmd = &cobra.Command{
		Use:   "users",
		Short: "List users in group",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			l := member.ListGroupUsers(
				ctx,
				setup.Gov,
				member.Group(groupName),
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(l))
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
