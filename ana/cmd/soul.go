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
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the public and private repositories of your soul",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := services.SoulService{
				SoulConfig: proto.SoulConfig{
					PublicURL:  publicURL,
					PrivateURL: privateURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "init")
			base.AssertNoErr(err)
			r, err := s.Init(files.WithWorkDir(cmd.Context(), workDir), &services.SoulInitIn{})
			if err == nil {
				fmt.Fprint(os.Stdout, r.Human())
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}
)
