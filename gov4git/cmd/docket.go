package cmd

import (
	"github.com/spf13/cobra"
)

var (
	docketCmd = &cobra.Command{
		Use:   "docket",
		Short: "Collaboration tools",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}
)

func init() {
	docketCmd.AddCommand(motionCmd)
}
