package cmd

import (
	"fmt"
	"os"

	gov4git_root "github.com/gov4git/gov4git"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version and build information",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(os.Stdout, gov4git_root.VersionDotTxt)
		},
	}
)
