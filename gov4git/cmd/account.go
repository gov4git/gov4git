package cmd

import (
	"github.com/gov4git/gov4git/v2/gov4git/api"
	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/spf13/cobra"
)

var (
	accountCmd = &cobra.Command{
		Use:   "account",
		Short: "Manage accounts",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	accountIssueCmd = &cobra.Command{
		Use:   "issue",
		Short: "Issue to account",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke(
				func() {
					LoadConfig()
					account.Issue(
						ctx,
						setup.Gov,
						account.AccountID(accountToID),
						account.H(
							account.Asset(accountAsset),
							accountQuantity,
						),
						accountNote,
					)
				},
			)
		},
	}

	accountBurnCmd = &cobra.Command{
		Use:   "burn",
		Short: "Burn from account",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke(
				func() {
					LoadConfig()
					account.Burn(
						ctx,
						setup.Gov,
						account.AccountID(accountFromID),
						account.H(
							account.Asset(accountAsset),
							accountQuantity,
						),
						accountNote,
					)
				},
			)
		},
	}

	accountTransferCmd = &cobra.Command{
		Use:   "transfer",
		Short: "Transfer from one account to another",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke(
				func() {
					LoadConfig()
					account.Transfer(
						ctx,
						setup.Gov,
						account.AccountID(accountFromID),
						account.AccountID(accountToID),
						account.H(
							account.Asset(accountAsset),
							accountQuantity,
						),
						accountNote,
					)
				},
			)
		},
	}

	accountListCmd = &cobra.Command{
		Use:   "list",
		Short: "List accounts",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					return account.List(
						ctx,
						setup.Gov,
					)
				},
			)
		},
	}

	accountShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show account",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					return account.Get(
						ctx,
						setup.Gov,
						account.AccountID(accountID),
					)
				},
			)
		},
	}

	accountBalanceCmd = &cobra.Command{
		Use:   "balance",
		Short: "Show account balance for a given asset",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					v := account.Get(
						ctx,
						setup.Gov,
						account.AccountID(accountID),
					)
					return v.Balance(account.Asset(accountAsset)).Quantity
				},
			)
		},
	}
)

var (
	accountID       string
	accountFromID   string
	accountToID     string
	accountAsset    string
	accountQuantity float64
	accountNote     string
)

func init() {
	// issue
	accountCmd.AddCommand(accountIssueCmd)
	accountIssueCmd.Flags().StringVar(&accountToID, "to", "", "to account id")
	accountIssueCmd.MarkFlagRequired("to")
	accountIssueCmd.Flags().StringVarP(&accountAsset, "asset", "a", "", "asset")
	accountIssueCmd.MarkFlagRequired("asset")
	accountIssueCmd.Flags().Float64VarP(&accountQuantity, "quantity", "q", 0.0, "quantity")
	accountIssueCmd.MarkFlagRequired("quantity")
	accountIssueCmd.Flags().StringVarP(&accountNote, "Note", "n", "manual", "note")
	// burn
	accountCmd.AddCommand(accountBurnCmd)
	accountBurnCmd.Flags().StringVar(&accountFromID, "from", "", "from account id")
	accountBurnCmd.MarkFlagRequired("from")
	accountBurnCmd.Flags().StringVarP(&accountAsset, "asset", "a", "", "asset")
	accountBurnCmd.MarkFlagRequired("asset")
	accountBurnCmd.Flags().Float64VarP(&accountQuantity, "quantity", "q", 0.0, "quantity")
	accountBurnCmd.MarkFlagRequired("quantity")
	accountBurnCmd.Flags().StringVarP(&accountNote, "Note", "n", "manual", "note")
	// transfer
	accountCmd.AddCommand(accountTransferCmd)
	accountTransferCmd.Flags().StringVar(&accountFromID, "from", "", "from account id")
	accountTransferCmd.MarkFlagRequired("from")
	accountTransferCmd.Flags().StringVar(&accountToID, "to", "", "to account id")
	accountTransferCmd.MarkFlagRequired("to")
	accountTransferCmd.Flags().StringVarP(&accountAsset, "asset", "a", "", "asset")
	accountTransferCmd.MarkFlagRequired("asset")
	accountTransferCmd.Flags().Float64VarP(&accountQuantity, "quantity", "q", 0.0, "quantity")
	accountTransferCmd.MarkFlagRequired("quantity")
	accountTransferCmd.Flags().StringVarP(&accountNote, "Note", "n", "manual", "note")
	// list
	accountCmd.AddCommand(accountListCmd)
	// show
	accountCmd.AddCommand(accountShowCmd)
	accountShowCmd.Flags().StringVar(&accountID, "id", "", "account id")
	accountShowCmd.MarkFlagRequired("id")
	// balance
	accountCmd.AddCommand(accountBalanceCmd)
	accountBalanceCmd.Flags().StringVar(&accountID, "id", "", "account id")
	accountBalanceCmd.MarkFlagRequired("id")
	accountBalanceCmd.Flags().StringVarP(&accountAsset, "asset", "a", "", "asset")
	accountBalanceCmd.MarkFlagRequired("asset")
}
