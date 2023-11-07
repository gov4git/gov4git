package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/proto/docket/docket"
	"github.com/gov4git/lib4git/form"
	"github.com/spf13/cobra"
)

var (
	motionCmd = &cobra.Command{
		Use:   "concern",
		Short: "Manage concerns",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	motionOpenCmd = &cobra.Command{
		Use:   "open",
		Short: "Open a new concern",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			docket.OpenMotion(
				ctx,
				setup.Gov,
				docket.MotionID(motionName),
				motionTitle,
				motionDesc,
				docket.ParseMotionType(ctx, motionType),
				motionTrackerURL,
				nil,
			)
			// fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	motionCloseCmd = &cobra.Command{
		Use:   "close",
		Short: "Close a concern",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			docket.CloseMotion(
				ctx,
				setup.Gov,
				docket.MotionID(motionName),
			)
			// fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	motionListCmd = &cobra.Command{
		Use:   "list",
		Short: "List concerns",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			l := docket.ListMotions(ctx, setup.Gov)
			fmt.Fprint(os.Stdout, form.SprintJSON(l))
		},
	}
)

var (
	motionName       string
	motionTitle      string
	motionDesc       string
	motionType       string
	motionTrackerURL string
)

func init() {
	motionCmd.AddCommand(motionOpenCmd)
	motionOpenCmd.Flags().StringVar(&motionName, "name", "", "unique name for motion")
	motionOpenCmd.MarkFlagRequired("name")
	motionOpenCmd.Flags().StringVar(&motionTitle, "title", "", "title for motion")
	motionOpenCmd.Flags().StringVar(&motionDesc, "desc", "", "description for motion")
	motionOpenCmd.Flags().StringVar(&motionType, "type", "concern", "type of motion (concern, proposal)")
	motionOpenCmd.Flags().StringVar(&motionTrackerURL, "tracking", "", "tracking URL for motion")

	motionCmd.AddCommand(motionCloseCmd)
	motionCloseCmd.Flags().StringVar(&motionName, "name", "", "name of motion")
	motionCloseCmd.MarkFlagRequired("name")

	motionCmd.AddCommand(motionListCmd)
	motionListCmd.Flags().StringVar(&motionName, "name", "", "name of motion")
	motionListCmd.MarkFlagRequired("name")
}
