package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/v2/proto/panorama"
	"github.com/gov4git/lib4git/form"
	"github.com/spf13/cobra"
)

var (
	panoramaCmd = &cobra.Command{
		Use:   "panorama",
		Short: "Panoramic view of user and motions, capturing pending votes",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			pano := panorama.Panorama(ctx, setup.Gov, setup.Member)
			fmt.Fprint(os.Stdout, form.SprintJSON(pano))
		},
	}
)
