package cmd

import "github.com/spf13/cobra"

var (
	// governance operations
	govCmd = &cobra.Command{
		Use:   "gov",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// configuration of governance policy
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	configAddCmd = &cobra.Command{
		Use:   "add",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	configRmCmd = &cobra.Command{
		Use:   "rm",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	configUpCmd = &cobra.Command{
		Use:   "up",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// approval of change proceedings

	// propose a change and begin governance review proceedings (e.g. referendum)
	proposeCmd = &cobra.Command{
		Use:   "propose",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	// cast a vote on referendum
	voteCmd = &cobra.Command{
		Use:   "vote",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	// approve a proposed change
	approveCmd = &cobra.Command{
		Use:   "approve",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	// tally results and prepare proof of compliance
	tallyCmd = &cobra.Command{
		Use:   "tally",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

var (
	communityURL    string
	communityBranch string
)

func init() {
	govCmd.PersistentFlags().StringVar(&communityURL, "community", "", "community repo URL")
	govCmd.PersistentFlags().StringVar(&communityBranch, "branch", "", "branch in community repo to work with")

	govCmd.AddCommand(configCmd)
	govCmd.AddCommand(userCmd)
	govCmd.AddCommand(groupCmd)
	govCmd.AddCommand(memberCmd)
	govCmd.AddCommand(proposeCmd)
	govCmd.AddCommand(voteCmd)
	govCmd.AddCommand(approveCmd)
	govCmd.AddCommand(tallyCmd)

	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configRmCmd)
	configCmd.AddCommand(configUpCmd)
}
