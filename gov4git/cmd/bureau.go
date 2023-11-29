package cmd

import (
	"github.com/gov4git/gov4git/proto/bureau"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/spf13/cobra"
)

var (
	bureauCmd = &cobra.Command{
		Use:   "bureau",
		Short: "Handle requests to the community initiated by individual members",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	bureauProcessCmd = &cobra.Command{
		Use:   "process",
		Short: "Fetch and process requests from community members",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			bureau.Process(
				ctx,
				setup.Organizer,
				member.Group(bureauGroup),
			)
		},
	}

	bureauTransferCmd = &cobra.Command{
		Use:   "transfer",
		Short: "Make a transfer request to the community governance",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			bureau.Transfer(
				ctx,
				setup.Member,
				setup.Gov,
				member.User(bureauFromUser),
				member.User(bureauToUser),
				bureauAmount,
			)
		},
	}
)

var (
	bureauGroup    string
	bureauFromUser string
	bureauToUser   string
	bureauAmount   float64
)

func init() {
	bureauCmd.AddCommand(bureauProcessCmd)
	bureauProcessCmd.Flags().StringVar(&bureauGroup, "group", "", "group of users to process requests from")
	bureauProcessCmd.MarkFlagRequired("group")

	bureauCmd.AddCommand(bureauTransferCmd)
	bureauTransferCmd.Flags().StringVar(&bureauFromUser, "from", "", "transfer from user")
	bureauTransferCmd.Flags().StringVar(&bureauToUser, "to", "", "transfer to user")
	bureauTransferCmd.Flags().Float64Var(&bureauAmount, "amount", 0, "transfer amount")
	bureauTransferCmd.MarkFlagRequired("amount")
}
