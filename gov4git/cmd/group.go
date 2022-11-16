package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/mod/member"
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
				setup.Community,
				member.Group(groupName),
			)
			// fmt.Fprint(os.Stdout, form.Pretty(chg.Result))
		},
	}

	groupRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove group from the community",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			member.RemoveGroup(
				ctx,
				setup.Community,
				member.Group(groupName),
			)
			// fmt.Fprint(os.Stdout, form.Pretty(chg.Result))
		},
	}

	groupListCmd = &cobra.Command{
		Use:   "list",
		Short: "List users in group",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			l := member.ListGroupUsers(
				ctx,
				setup.Community,
				member.Group(groupName),
			)
			fmt.Fprint(os.Stdout, form.Pretty(l))
		},
	}
)

var (
	groupName string
)

func init() {
	groupCmd.AddCommand(groupAddCmd)
	groupAddCmd.Flags().StringVar(&groupName, "name", "", "group alias within the community")

	groupCmd.AddCommand(groupRemoveCmd)
	groupRemoveCmd.Flags().StringVar(&groupName, "name", "", "group alias within the community")

	groupCmd.AddCommand(groupListCmd)
	groupListCmd.Flags().StringVar(&groupName, "name", "", "group alias within the community")
}
