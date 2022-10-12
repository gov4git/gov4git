package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/proto/identityproto"
	"github.com/gov4git/gov4git/services/identity"
	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the public and private repositories of your soul",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := identity.IdentityService{
				IdentityConfig: identityproto.IdentityConfig{
					PublicURL:  publicURL,
					PrivateURL: privateURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(govproto.LocalAgentTempPath, "init")
			base.AssertNoErr(err)
			r, err := s.Init(files.WithWorkDir(cmd.Context(), workDir), &identity.InitIn{})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}
)
