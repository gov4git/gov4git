package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/proto/collab"
	"github.com/gov4git/lib4git/form"
	"github.com/spf13/cobra"
)

var (
	concernCmd = &cobra.Command{
		Use:   "concern",
		Short: "Manage concerns",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	concernOpenCmd = &cobra.Command{
		Use:   "open",
		Short: "Open a new concern",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			collab.OpenConcern(
				ctx,
				setup.Gov,
				collab.Name(concernName),
				concernTitle,
				concernDesc,
				concernTrackerURL,
			)
			// fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	concernCloseCmd = &cobra.Command{
		Use:   "close",
		Short: "Close a concern",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			collab.CloseConcern(
				ctx,
				setup.Gov,
				collab.Name(concernName),
			)
			// fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	concernListCmd = &cobra.Command{
		Use:   "list",
		Short: "List concerns",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			l := collab.ListConcerns(ctx, setup.Gov)
			fmt.Fprint(os.Stdout, form.SprintJSON(l))
		},
	}
)

var (
	concernName       string
	concernTitle      string
	concernDesc       string
	concernTrackerURL string
)

func init() {
	concernCmd.AddCommand(concernOpenCmd)
	concernOpenCmd.Flags().StringVar(&concernName, "name", "", "unique name for concern")
	concernOpenCmd.MarkFlagRequired("name")
	concernOpenCmd.Flags().StringVar(&concernTitle, "title", "", "title for concern")
	concernOpenCmd.Flags().StringVar(&concernDesc, "desc", "", "description for concern")
	concernOpenCmd.Flags().StringVar(&concernTrackerURL, "tracking", "", "tracking URL for concern")

	concernCmd.AddCommand(concernCloseCmd)
	concernCloseCmd.Flags().StringVar(&concernName, "name", "", "name of concern")
	concernCloseCmd.MarkFlagRequired("name")

	concernCmd.AddCommand(concernListCmd)
	concernListCmd.Flags().StringVar(&concernName, "name", "", "name of concern")
	concernListCmd.MarkFlagRequired("name")
}
