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

	// configuration of governance rules

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
	// tally results and prepare proof of compliance
	tallyCmd = &cobra.Command{
		Use:   "tally",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	govCmd.AddCommand(proposeCmd)
	govCmd.AddCommand(voteCmd)
	govCmd.AddCommand(tallyCmd)
}
