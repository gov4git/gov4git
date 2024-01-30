package cmd

import (
	"github.com/gov4git/gov4git/v2/gov4git/api"
	"github.com/gov4git/gov4git/v2/proto/boot"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/spf13/cobra"
)

var (
	initIDCmd = &cobra.Command{
		Use:   "init-id",
		Short: "Initialize public and private repositories of your identity",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					chg := id.Init(ctx, setup.Member)
					return chg.Result
				},
			)
		},
	}

	initGovCmd = &cobra.Command{
		Use:   "init-gov",
		Short: "Initialize public and private repositories of your governance",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					chg := boot.Boot(ctx, setup.Organizer)
					return chg.Result
				},
			)
		},
	}
)
