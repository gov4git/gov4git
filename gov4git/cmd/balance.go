package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/form"
	man "github.com/gov4git/gov4git/man/user"
	"github.com/gov4git/gov4git/proto/cmdproto"
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/services/gov/user"
	"github.com/spf13/cobra"
)

var (
	balanceCmd = &cobra.Command{
		Use:   "balance",
		Short: "Manage user balances",
		Long:  man.GovUserBalance,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	balanceAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add (or subtract) from a user balance",
		Long:  man.GovUserBalanceAdd,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := user.GovUserService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-user-balance-add")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.BalanceAdd(ctx, &user.BalanceAddIn{
				User:    balanceUser,
				Balance: balanceBalance,
				Branch:  balanceBranch,
				Value:   balanceValue,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	balanceMulCmd = &cobra.Command{
		Use:   "mul",
		Short: "Multiply (or divide) a user balance",
		Long:  man.GovUserBalanceMul,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := user.GovUserService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-user-balance-mul")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.BalanceMul(ctx, &user.BalanceMulIn{
				User:    balanceUser,
				Balance: balanceBalance,
				Branch:  balanceBranch,
				Value:   balanceValue,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	balanceSetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set a user balance",
		Long:  man.GovUserBalanceSet,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := user.GovUserService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-user-balance-set")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.BalanceSet(ctx, &user.BalanceSetIn{
				User:    balanceUser,
				Balance: balanceBalance,
				Branch:  balanceBranch,
				Value:   balanceValue,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	balanceGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a user balance",
		Long:  man.GovUserBalanceSet,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := user.GovUserService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-user-balance-get")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.BalanceGet(ctx, &user.BalanceGetIn{
				User:    balanceUser,
				Balance: balanceBalance,
				Branch:  balanceBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}
)

var (
	balanceUser    string
	balanceBalance string
	balanceBranch  string
	balanceValue   float64
)

func init() {
	balanceCmd.AddCommand(balanceAddCmd)
	balanceCmd.AddCommand(balanceMulCmd)
	balanceCmd.AddCommand(balanceSetCmd)
	balanceCmd.AddCommand(balanceGetCmd)

	balanceAddCmd.Flags().StringVar(&balanceUser, "user", "", "name of user")
	balanceAddCmd.Flags().StringVar(&balanceBalance, "balance", "", "name of balance")
	balanceAddCmd.Flags().StringVar(&balanceBranch, "branch", "", "repo branch")
	balanceAddCmd.Flags().Float64Var(&balanceValue, "value", 0.0, "amount to add")

	balanceMulCmd.Flags().StringVar(&balanceUser, "user", "", "name of user")
	balanceMulCmd.Flags().StringVar(&balanceBalance, "balance", "", "name of balance")
	balanceMulCmd.Flags().StringVar(&balanceBranch, "branch", "", "repo branch")
	balanceMulCmd.Flags().Float64Var(&balanceValue, "value", 0.0, "amount to multiply by")

	balanceSetCmd.Flags().StringVar(&balanceUser, "user", "", "name of user")
	balanceSetCmd.Flags().StringVar(&balanceBalance, "balance", "", "name of balance")
	balanceSetCmd.Flags().StringVar(&balanceBranch, "branch", "", "repo branch")
	balanceSetCmd.Flags().Float64Var(&balanceValue, "value", 0.0, "amount to set")

	balanceGetCmd.Flags().StringVar(&balanceUser, "user", "", "name of user")
	balanceGetCmd.Flags().StringVar(&balanceBalance, "balance", "", "name of balance")
	balanceGetCmd.Flags().StringVar(&balanceBranch, "branch", "", "repo branch")
}
