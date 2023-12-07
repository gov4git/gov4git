package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/lib4git/form"
	"github.com/spf13/cobra"
)

var (
	accountCmd = &cobra.Command{
		Use:   "account",
		Short: "Manage accounts",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	accountDepositCmd = &cobra.Command{
		Use:   "deposit",
		Short: "Make a deposit",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			account.Deposit(
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
	}

	accountWithdrawCmd = &cobra.Command{
		Use:   "withdraw",
		Short: "Make a withdrawal",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			account.Withdraw(
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
	}

	accountTransferCmd = &cobra.Command{
		Use:   "transfer",
		Short: "Transfer from one account to another",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
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
	}

	accountListCmd = &cobra.Command{
		Use:   "list",
		Short: "List accounts",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			v := account.List(
				ctx,
				setup.Gov,
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(v))
		},
	}

	accountShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show account",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			v := account.Get(
				ctx,
				setup.Gov,
				account.AccountID(accountID),
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(v))
		},
	}

	accountBalanceCmd = &cobra.Command{
		Use:   "balance",
		Short: "Show account balance for a given asset",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			v := account.Get(
				ctx,
				setup.Gov,
				account.AccountID(accountID),
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(v.Balance(account.Asset(accountAsset)).Quantity))
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
	// deposit
	accountCmd.AddCommand(accountDepositCmd)
	accountDepositCmd.Flags().StringVar(&accountToID, "to", "", "to account id")
	accountDepositCmd.MarkFlagRequired("to")
	accountDepositCmd.Flags().StringVarP(&accountAsset, "asset", "a", "", "asset")
	accountDepositCmd.MarkFlagRequired("asset")
	accountDepositCmd.Flags().Float64VarP(&accountQuantity, "quantity", "q", 0.0, "quantity")
	accountDepositCmd.MarkFlagRequired("quantity")
	accountDepositCmd.Flags().StringVarP(&accountNote, "Note", "n", "manual", "note")
	// withdraw
	accountCmd.AddCommand(accountWithdrawCmd)
	accountWithdrawCmd.Flags().StringVar(&accountFromID, "from", "", "from account id")
	accountWithdrawCmd.MarkFlagRequired("from")
	accountWithdrawCmd.Flags().StringVarP(&accountAsset, "asset", "a", "", "asset")
	accountWithdrawCmd.MarkFlagRequired("asset")
	accountWithdrawCmd.Flags().Float64VarP(&accountQuantity, "quantity", "q", 0.0, "quantity")
	accountWithdrawCmd.MarkFlagRequired("quantity")
	accountWithdrawCmd.Flags().StringVarP(&accountNote, "Note", "n", "manual", "note")
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
