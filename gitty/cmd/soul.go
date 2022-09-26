package cmd

import "github.com/spf13/cobra"

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the public and private repositories of your user",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Send a file or directory to another user",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	receiveCmd = &cobra.Command{
		Use:   "receive",
		Short: "Receive files or directories from another user",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)
