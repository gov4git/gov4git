package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
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
			ops.OpenMotion(
				ctx,
				setup.Gov,
				schema.MotionID(motionName),
				schema.PolicyName(motionPolicy),
				motionTitle,
				motionDesc,
				schema.ParseMotionType(ctx, motionType),
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
			ops.CloseMotion(
				ctx,
				setup.Gov,
				schema.MotionID(motionName),
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
			l := ops.ListMotions(ctx, setup.Gov)
			fmt.Fprint(os.Stdout, form.SprintJSON(l))
		},
	}
)

var (
	motionName       string
	motionPolicy     string
	motionTitle      string
	motionDesc       string
	motionType       string
	motionTrackerURL string
)

func init() {
	motionCmd.AddCommand(motionOpenCmd)
	motionOpenCmd.Flags().StringVar(&motionName, "name", "", "unique name for motion")
	motionOpenCmd.MarkFlagRequired("name")
	motionOpenCmd.Flags().StringVar(&motionPolicy, "policy", "", "policy ("+strings.Join(policy.Policies(), ", ")+")")
	motionOpenCmd.MarkFlagRequired("policy")
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
