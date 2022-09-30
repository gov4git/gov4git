package cmd

import (
	"fmt"
	"os"

	"github.com/petar/gitty/lib/base"
	"github.com/petar/gitty/lib/files"
	man "github.com/petar/gitty/man/user"
	"github.com/petar/gitty/proto"
	"github.com/petar/gitty/services/user"
	"github.com/spf13/cobra"
)

var (
	// user management
	userCmd = &cobra.Command{
		Use:   "user",
		Short: "Manage users",
		Long:  man.GovUser,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	userAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add user to the community",
		Long:  man.GovUserAdd,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := user.GovUserService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-user-add")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.UserAdd(ctx, &user.GovUserAddIn{
				Name:            userName,
				URL:             userURL,
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

	userRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove user from community",
		Long:  man.GovUserRemove,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := user.GovUserService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-user-rm")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.UserRemove(ctx, &user.GovUserRemoveIn{
				Name:            userName,
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

	userSetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set user property",
		Long:  man.GovUserSet,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := user.GovUserService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-user-set")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.UserSet(ctx, &user.GovUserSetIn{
				Name:            userName,
				Key:             userKey,
				Value:           userValue,
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

	userGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get user property",
		Long:  man.GovUserGet,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := user.GovUserService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-user-get")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.UserGet(ctx, &user.GovUserGetIn{
				Name:            userName,
				Key:             userKey,
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

	userListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all community users",
		Long:  man.GovUserList,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := user.GovUserService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-user-list")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.UserList(ctx, &user.GovUserListIn{
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
	userName  string
	userURL   string
	userKey   string
	userValue string
)

func init() {
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userRemoveCmd)
	userCmd.AddCommand(userSetCmd)
	userCmd.AddCommand(userGetCmd)
	userCmd.AddCommand(userListCmd)

	userAddCmd.Flags().StringVar(&userName, "name", "", "name of user, unique for the community")
	userAddCmd.Flags().StringVar(&userURL, "url", "", "URL of user")

	userRemoveCmd.Flags().StringVar(&userName, "name", "", "name of user, unique for the community")

	userSetCmd.Flags().StringVar(&userName, "name", "", "name of user")
	userSetCmd.Flags().StringVar(&userKey, "key", "", "user property key")
	userSetCmd.Flags().StringVar(&userValue, "value", "", "user property value")

	userGetCmd.Flags().StringVar(&userName, "name", "", "name of user")
	userGetCmd.Flags().StringVar(&userKey, "key", "", "user property key")
}
