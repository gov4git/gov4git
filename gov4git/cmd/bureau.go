package cmd

import (
	"github.com/gov4git/gov4git/mod/balance"
	"github.com/gov4git/gov4git/mod/bureau"
	"github.com/gov4git/gov4git/mod/member"
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
			bureau.Transfer(
				ctx,
				setup.Member,
				setup.Community,
				member.User(bureauFromUser),
				balance.Balance(bureauFromBalance),
				member.User(bureauToUser),
				balance.Balance(bureauToBalance),
				bureauAmount,
			)
		},
	}
)

var (
	bureauGroup       string
	bureauFromUser    string
	bureauFromBalance string
	bureauToUser      string
	bureauToBalance   string
	bureauAmount      float64
)

func init() {
	bureauCmd.AddCommand(bureauProcessCmd)
	bureauProcessCmd.Flags().StringVar(&bureauGroup, "group", "", "group of users to process requests from")
	bureauProcessCmd.MarkFlagRequired("group")

	bureauCmd.AddCommand(bureauTransferCmd)
	bureauTransferCmd.Flags().StringVar(&bureauFromUser, "from-user", "", "transfer from user")
	bureauTransferCmd.Flags().StringVar(&bureauFromBalance, "from-balance", "", "transfer from balance")
	bureauTransferCmd.MarkFlagRequired("from-balance")
	bureauTransferCmd.Flags().StringVar(&bureauToUser, "to-user", "", "transfer to user")
	bureauTransferCmd.Flags().StringVar(&bureauToBalance, "to-balance", "", "transfer to balance")
	bureauTransferCmd.MarkFlagRequired("to-balance")
	bureauTransferCmd.Flags().Float64Var(&bureauAmount, "amount", 0, "transfer amount")
	bureauTransferCmd.MarkFlagRequired("amount")
}
