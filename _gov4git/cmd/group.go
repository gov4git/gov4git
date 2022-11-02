package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/form"
	man "github.com/gov4git/gov4git/man/group"
	"github.com/gov4git/gov4git/proto/cmdproto"
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/services/gov/group"
	"github.com/spf13/cobra"
)

var (
	// group management
	groupCmd = &cobra.Command{
		Use:   "group",
		Short: "Manage groups",
		Long:  man.GovGroup,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	groupAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add group to the community",
		Long:  man.GovGroupAdd,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := group.GovGroupService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-group-add")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.Add(ctx, &group.AddIn{
				Name:            groupName,
				CommunityBranch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	groupRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove group from community",
		Long:  man.GovGroupRemove,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := group.GovGroupService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-group-rm")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.Remove(ctx, &group.RemoveIn{
				Name:            groupName,
				CommunityBranch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	groupSetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set group property",
		Long:  man.GovGroupSet,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := group.GovGroupService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-group-set")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.Set(ctx, &group.SetIn{
				Name:            groupName,
				Key:             groupKey,
				Value:           groupValue,
				CommunityBranch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	groupGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get group property",
		Long:  man.GovGroupGet,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := group.GovGroupService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-group-get")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.Get(ctx, &group.GetIn{
				Name:            groupName,
				Key:             groupKey,
				CommunityBranch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	groupListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all community groups",
		Long:  man.GovGroupList,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := group.GovGroupService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-group-list")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.List(ctx, &group.ListIn{
				CommunityBranch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}
)

var (
	groupName  string
	groupKey   string
	groupValue string
)

func init() {
	groupCmd.AddCommand(groupAddCmd)
	groupCmd.AddCommand(groupRemoveCmd)
	groupCmd.AddCommand(groupSetCmd)
	groupCmd.AddCommand(groupGetCmd)
	groupCmd.AddCommand(groupListCmd)

	groupAddCmd.Flags().StringVar(&groupName, "name", "", "name of group, unique for the community")

	groupRemoveCmd.Flags().StringVar(&groupName, "name", "", "name of group, unique for the community")

	groupSetCmd.Flags().StringVar(&groupName, "name", "", "name of group")
	groupSetCmd.Flags().StringVar(&groupKey, "key", "", "group property key")
	groupSetCmd.Flags().StringVar(&groupValue, "value", "", "group property value")

	groupGetCmd.Flags().StringVar(&groupName, "name", "", "name of group")
	groupGetCmd.Flags().StringVar(&groupKey, "key", "", "group property key")
}
