package cmd

import (
	"github.com/spf13/cobra"
)

var (
	collabCmd = &cobra.Command{
		Use:   "collab",
		Short: "Collaboration tools",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}
)

func init() {
	collabCmd.AddCommand(concernCmd)
}
