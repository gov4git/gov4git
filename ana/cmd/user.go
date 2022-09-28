package cmd

import (
	"fmt"
	"os"

	"github.com/petar/gitty/lib/base"
	"github.com/petar/gitty/lib/files"
	"github.com/petar/gitty/man"
	"github.com/petar/gitty/proto"
	"github.com/petar/gitty/services"
	"github.com/spf13/cobra"
)

var (
	// user management
	userCmd = &cobra.Command{
		Use:   "user",
		Short: "User management",
		Long:  man.GovUser,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	userAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add user to the community",
		Long:  man.GovUserAdd,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := services.GovService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-user-add")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.UserAdd(ctx, &services.GovUserAddIn{
				Name:            userName,
				URL:             userURL,
				CommunityBranch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, r.Human())
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	userRmCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove user from community",
		Long:  man.GovUserRm,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := services.GovService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-user-rm")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.UserRemove(ctx, &services.GovUserRemoveIn{
				Name:            userName,
				CommunityBranch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, r.Human())
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
			s := services.GovService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-user-set")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.UserSet(ctx, &services.GovUserSetIn{
				Name:            userName,
				Key:             userKey,
				Value:           userValue,
				CommunityBranch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, r.Human())
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
			s := services.GovService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-user-get")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.UserGet(ctx, &services.GovUserGetIn{
				Name:            userName,
				Key:             userKey,
				CommunityBranch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, r.Human())
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
	userCmd.AddCommand(userRmCmd)
	userCmd.AddCommand(userSetCmd)
	userCmd.AddCommand(userGetCmd)

	userAddCmd.Flags().StringVar(&userName, "name", "", "name of user, unique for the community")
	userAddCmd.Flags().StringVar(&userName, "url", "", "URL of user")

	userRmCmd.Flags().StringVar(&userName, "name", "", "name of user, unique for the community")

	userSetCmd.Flags().StringVar(&userKey, "key", "", "user property key")
	userSetCmd.Flags().StringVar(&userValue, "value", "", "user property value")

	userGetCmd.Flags().StringVar(&userKey, "key", "", "user property key")
}
