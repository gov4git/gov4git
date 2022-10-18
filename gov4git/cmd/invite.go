package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/form"
	man "github.com/gov4git/gov4git/man/gov"
	"github.com/gov4git/gov4git/proto/cmdproto"
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/services/gov"
	"github.com/spf13/cobra"
)

var (
	inviteCmd = &cobra.Command{
		Use:   "invite",
		Short: "Invite someone to the community",
		Long:  man.Invite,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := gov.GovService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-invite")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			if config.SMTPPlainAuth == nil {
				err := fmt.Errorf("smtp not configured")
				fmt.Fprint(os.Stderr, err.Error())
				return err
			}
			if config.InviteEmailFrom == nil {
				err := fmt.Errorf("invite from email not set")
				fmt.Fprint(os.Stderr, err.Error())
				return err
			}
			r, err := s.Invite(ctx, &gov.InviteIn{
				SMTP:            *config.SMTPPlainAuth,
				From:            *config.InviteEmailFrom,
				To:              cmdproto.EmailAddress{Name: inviteeName, Address: inviteeAddress},
				CommunityURL:    communityURL,
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
	inviteeName    string
	inviteeAddress string
)

func init() {
	inviteCmd.Flags().StringVar(&inviteeName, "name", "", "greeting name of invitee")
	inviteCmd.Flags().StringVar(&inviteeAddress, "address", "", "email address of invitee")
}
