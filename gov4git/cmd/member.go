package cmd

import (
	"github.com/gov4git/gov4git/mod/member"
	"github.com/spf13/cobra"
)

var (
	memberCmd = &cobra.Command{
		Use:   "member",
		Short: "Manage members",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	memberAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add member",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			member.AddMember(
				ctx,
				setup.Community,
				member.User(memberUser),
				member.Group(memberGroup),
			)
		},
	}

	memberRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove member from the community",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			member.RemoveMember(
				ctx,
				setup.Community,
				member.User(memberUser),
				member.Group(memberGroup),
			)
		},
	}
)

var (
	memberUser  string
	memberGroup string
)

func init() {
	memberCmd.AddCommand(memberAddCmd)
	memberAddCmd.Flags().StringVar(&memberUser, "user", "", "user alias within the community")
	memberAddCmd.Flags().StringVar(&memberGroup, "group", "", "group alias within the community")

	memberCmd.AddCommand(memberRemoveCmd)
	memberRemoveCmd.Flags().StringVar(&memberUser, "user", "", "user alias within the community")
	memberRemoveCmd.Flags().StringVar(&memberGroup, "group", "", "group alias within the community")
}
