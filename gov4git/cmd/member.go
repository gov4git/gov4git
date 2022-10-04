package cmd

import (
	"fmt"
	"os"

	"github.com/petar/gov4git/lib/base"
	"github.com/petar/gov4git/lib/files"
	man "github.com/petar/gov4git/man/member"
	"github.com/petar/gov4git/proto"
	"github.com/petar/gov4git/services/gov/member"
	"github.com/spf13/cobra"
)

var (
	// member management
	memberCmd = &cobra.Command{
		Use:   "member",
		Short: "Manage memberships",
		Long:  man.GovMember,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	memberAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Make a user a member of a group",
		Long:  man.GovMemberAdd,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := member.GovMemberService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-member-add")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.MemberAdd(ctx, &member.GovMemberAddIn{
				Group:           memberGroup,
				User:            memberUser,
				CommunityBranch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, r.Human(cmd.Context()))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	memberRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove a member of a group",
		Long:  man.GovMemberRemove,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := member.GovMemberService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-member-rm")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.MemberRemove(ctx, &member.GovMemberRemoveIn{
				Group:           memberGroup,
				User:            memberUser,
				CommunityBranch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, r.Human(cmd.Context()))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	memberListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all members of a group, or all groups that a user is a member of",
		Long:  man.GovMemberList,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := member.GovMemberService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-member-list-members")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.MemberList(ctx, &member.GovMemberListIn{
				User:            memberUser,
				Group:           memberGroup,
				CommunityBranch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, r.Human(cmd.Context()))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}
)

var (
	memberGroup string
	memberUser  string
)

func init() {
	memberCmd.AddCommand(memberAddCmd)
	memberCmd.AddCommand(memberRemoveCmd)
	memberCmd.AddCommand(memberListCmd)

	memberAddCmd.Flags().StringVar(&memberGroup, "group", "", "group name")
	memberAddCmd.Flags().StringVar(&memberUser, "user", "", "user name")

	memberRemoveCmd.Flags().StringVar(&memberGroup, "group", "", "group name")
	memberRemoveCmd.Flags().StringVar(&memberUser, "user", "", "user name")

	memberListCmd.Flags().StringVar(&memberGroup, "group", "", "group name")
	memberListCmd.Flags().StringVar(&memberUser, "user", "", "user name")
}
