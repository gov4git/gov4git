package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/v2"
	"github.com/gov4git/lib4git/form"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version and build information",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(os.Stdout, form.SprintJSON(gov4git.GetVersionInfo()))
		},
	}
)
