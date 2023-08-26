package bureau

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/mail"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Transfer(
	ctx context.Context,
	userAddr id.OwnerAddress,
	govAddr gov.GovAddress,
	fromUserOpt member.User, // optional, if empty string, a lookup forthe user is performed
	fromBalance balance.Balance,
	toUser member.User,
	toBalance balance.Balance,
	amount float64,
) git.Change[form.Map, mail.RequestEnvelope[Request]] {

	govCloned := git.CloneOne(ctx, git.Address(govAddr))
	// userRepo, userTree := id.CloneOwner(ctx, userAddr)
	userOwner := id.CloneOwner(ctx, userAddr)
	chg := TransferStageOnly(ctx, userAddr, govAddr, userOwner, govCloned, fromUserOpt, fromBalance, toUser, toBalance, amount)
	proto.Commit(ctx, userOwner.Public.Tree(), chg)
	userOwner.Public.Push(ctx)
	return chg
}

func TransferStageOnly(
	ctx context.Context,
	userAddr id.OwnerAddress,
	govAddr gov.GovAddress,
	userOwner id.OwnerCloned,
	govCloned git.Cloned,
	fromUserOpt member.User,
	fromBalance balance.Balance,
	toUser member.User,
	toBalance balance.Balance,
	amount float64,
) git.Change[form.Map, mail.RequestEnvelope[Request]] {

	userCred := id.GetPublicCredentials(ctx, userOwner.Public.Tree())

	// find the user name of userAddr in the community repo
	if fromUserOpt == "" {
		us := member.LookupUserByIDLocal(ctx, govCloned.Tree(), userCred.ID)
		switch len(us) {
		case 0:
			must.Errorf(ctx, "%s not found in community %v", userAddr.Public, govAddr)
		case 1:
			fromUserOpt = us[0]
		default:
			must.Errorf(ctx, "community %v has more than one user at address %v", govAddr, userAddr.Public)
		}
	}

	request := Request{
		Transfer: &TransferRequest{
			FromUser:    fromUserOpt,
			FromBalance: fromBalance,
			ToUser:      toUser,
			ToBalance:   toBalance,
			Amount:      amount,
		},
	}

	sendOnly := mail.Request_StageOnly(ctx, userOwner, govCloned.Tree(), BureauTopic, request)
	return git.NewChange(
		"Transfer account tokens.",
		"bureau_transfer",
		form.Map{
			"from_user":    fromUserOpt,
			"from_balance": fromBalance,
			"to_user":      toUser,
			"to_balance":   toBalance,
			"amount":       amount,
		},
		sendOnly.Result,
		form.Forms{sendOnly},
	)
}
