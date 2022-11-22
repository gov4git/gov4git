package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/spf13/cobra"
)

var (
	initIDCmd = &cobra.Command{
		Use:   "init-id",
		Short: "Initialize public and private repositories of your identity",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			chg := id.Init(ctx, setup.Member)
			fmt.Fprint(os.Stdout, form.Pretty(chg.Result))
		},
	}

	initGovCmd = &cobra.Command{
		Use:   "init-gov",
		Short: "Initialize public and private repositories of your governance",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			chg := id.Init(ctx, id.OwnerAddress(setup.Organizer))
			fmt.Fprint(os.Stdout, form.Pretty(chg.Result))
		},
	}
)
