package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gitty",
		Short: "gitty is a command-line tool for XXX",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	govCmd = &cobra.Command{
		Use:   "gov",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(receiveCmd)
	rootCmd.AddCommand(govCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
