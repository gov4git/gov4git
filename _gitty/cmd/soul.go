package cmd

import "github.com/spf13/cobra"

var (
	soulCmd = &cobra.Command{
		Use:   "soul",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	soulInitCmd = &cobra.Command{
		Use:   "init",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	soulCmd.AddCommand(soulInitCmd)
}
