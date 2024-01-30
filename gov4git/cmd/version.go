package cmd

import (
	"github.com/gov4git/gov4git/v2"
	"github.com/gov4git/gov4git/v2/gov4git/api"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version and build information",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					return gov4git.GetVersionInfo()
				},
			)
		},
	}
)
