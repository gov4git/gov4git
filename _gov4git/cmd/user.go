package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/form"
	man "github.com/gov4git/gov4git/man/user"
	"github.com/gov4git/gov4git/proto/cmdproto"
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/services/gov/user"
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
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-user-add")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.Add(ctx, &user.AddIn{
				Name:            userName,
				URL:             userURL,
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

	userRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove user from community",
		Long:  man.GovUserRemove,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := user.GovUserService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-user-rm")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.Remove(ctx, &user.RemoveIn{
				Name:            userName,
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

	userSetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set user property",
		Long:  man.GovUserSet,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := user.GovUserService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-user-set")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.Set(ctx, &user.SetIn{
				Name:   userName,
				Key:    userKey,
				Value:  userValue,
				Branch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
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
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-user-get")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.Get(ctx, &user.GetIn{
				Name:   userName,
				Key:    userKey,
				Branch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
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
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-user-list")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.List(ctx, &user.ListIn{
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
	userCmd.AddCommand(balanceCmd)

	userAddCmd.Flags().StringVar(&userName, "name", "", "name of user, unique for the community")
	userAddCmd.Flags().StringVar(&userURL, "url", "", "URL of user")

	userRemoveCmd.Flags().StringVar(&userName, "name", "", "name of user, unique for the community")

	userSetCmd.Flags().StringVar(&userName, "name", "", "name of user")
	userSetCmd.Flags().StringVar(&userKey, "key", "", "user property key")
	userSetCmd.Flags().StringVar(&userValue, "value", "", "user property value")

	userGetCmd.Flags().StringVar(&userName, "name", "", "name of user")
	userGetCmd.Flags().StringVar(&userKey, "key", "", "user property key")
}
