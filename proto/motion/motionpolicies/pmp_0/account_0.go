package pmp_0

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/gov"
)

var (
	PMPAccountID = account.AccountIDFromLine(
		account.Term("pmp"),
	)
	BurnPoolAccountID = account.AccountIDFromLine(
		account.Cat(
			account.Term("pmp"),
			account.Term("burn"),
		),
	)
	TaxPoolAccountID = account.AccountIDFromLine(
		account.Cat(
			account.Term("pmp"),
			account.Term("tax"),
		),
	)
	MatchingPoolAccountID = account.AccountIDFromLine(
		account.Cat(
			account.Term("pmp"),
			account.Term("matching"),
		),
	)
)

func Boot_StageOnly(ctx context.Context, cloned gov.Cloned) {

	// create burn pool account
	account.Create_StageOnly(
		ctx,
		cloned,
		BurnPoolAccountID,
		PMPAccountID,
		fmt.Sprintf("burn account for PMP"),
	)

	// create tax pool account
	account.Create_StageOnly(
		ctx,
		cloned,
		TaxPoolAccountID,
		PMPAccountID,
		fmt.Sprintf("tax account for PMP"),
	)

	// create matching pool account
	account.Create_StageOnly(
		ctx,
		cloned,
		MatchingPoolAccountID,
		PMPAccountID,
		fmt.Sprintf("matching pool account for PMP"),
	)

}

func GetMatchFundBalance_Local(ctx context.Context, cloned gov.Cloned) float64 {
	a := account.Get_Local(ctx, cloned, MatchingPoolAccountID)
	return a.Balance(account.PluralAsset).Quantity
}
