package treasury

import (
	"context"

	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/gov"
)

var (
	OwnerID = account.AccountIDFromLine(
		account.Term("treasury"),
	)
	BurnAccountID = account.AccountIDFromLine(
		account.Cat(
			account.Term("treasury"),
			account.Term("burn"),
		),
	)
)

func Boot_StageOnly(ctx context.Context, cloned gov.Cloned) {

	// create burn pool account
	account.Create_StageOnly(
		ctx,
		cloned,
		BurnAccountID,
		account.OwnerID(OwnerID),
	)

}
