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
)

var (
	communityURL    string
	communityBranch string
)

func init() {
	govCmd.PersistentFlags().StringVar(&communityURL, "community", "", "community repo URL")
	govCmd.PersistentFlags().StringVar(&communityBranch, "branch", "", "branch in community repo to work with")

	govCmd.AddCommand(userCmd)
	govCmd.AddCommand(groupCmd)
	govCmd.AddCommand(memberCmd)
	govCmd.AddCommand(policyCmd)
	govCmd.AddCommand(ballotCmd)
	govCmd.AddCommand(voteCmd)
	govCmd.AddCommand(tallyCmd)
	govCmd.AddCommand(sealCmd)
	govCmd.AddCommand(listCmd)
	govCmd.AddCommand(inviteCmd)
}
