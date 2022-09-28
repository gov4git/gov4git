package cmd

import (
	"fmt"
	"os"

	"github.com/petar/gitty/lib/base"
	"github.com/petar/gitty/lib/files"
	"github.com/petar/gitty/proto"
	"github.com/petar/gitty/services"
	"github.com/spf13/cobra"
)

var (
	// user management
	userCmd = &cobra.Command{
		Use:   "user",
		Short: "User management",
		Long:  `Add and remove users. Set and get user properties.`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	userAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add user to the community",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := services.GovService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-user-add")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.AddUser(ctx, &services.GovAddUserIn{
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
		Short: "Remove user",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	userSetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set user property",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	userGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get user property",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

var (
	userName string
	userURL  string
)

func init() {
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userRmCmd)
	userCmd.AddCommand(userSetCmd)
	userCmd.AddCommand(userGetCmd)

	userAddCmd.Flags().StringVar(&userName, "name", "", "name of user, unique for the community")
	userAddCmd.Flags().StringVar(&userName, "url", "", "URL of user")
}
