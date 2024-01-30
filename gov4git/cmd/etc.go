package cmd

import (
	"encoding/json"
	"io"
	"os"

	"github.com/gov4git/gov4git/v2/gov4git/api"
	"github.com/gov4git/gov4git/v2/proto/etc"
	"github.com/gov4git/lib4git/must"
	"github.com/spf13/cobra"
)

var (
	etcCmd = &cobra.Command{
		Use:   "etc",
		Short: "Manage system settings",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	etcGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get system settings",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					settings := etc.GetSettings(ctx, setup.Gov)
					return settings
				},
			)
		},
	}

	etcSetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set system settings",
		Long:  `System settings must be given as JSON on the standard input.`,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke(
				func() {
					LoadConfig()

					jsonData, err := io.ReadAll(os.Stdin)
					must.NoError(ctx, err)

					var settings etc.Settings
					err = json.Unmarshal(jsonData, &settings)
					must.NoError(ctx, err)

					etc.SetSettings(ctx, setup.Gov, settings)
				},
			)
		},
	}
)

func init() {
	etcCmd.AddCommand(etcGetCmd)
	etcCmd.AddCommand(etcSetCmd)
}
