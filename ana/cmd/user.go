package cmd

import "github.com/spf13/cobra"

var (
	// user management
	userCmd = &cobra.Command{
		Use:   "user",
		Short: "User management",
		Long:  `Add and remove users. Set and get user properties.`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	userAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add user",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	userRmCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove user",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	userSetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set user property",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	userGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get user property",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userRmCmd)
	userCmd.AddCommand(userSetCmd)
	userCmd.AddCommand(userGetCmd)
}
