package bureau

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/mail"
	"github.com/gov4git/gov4git/proto/member"
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
) git.Change[mail.SeqNo] {

	govRepo := git.CloneRepo(ctx, git.Address(govAddr))
	userRepo, userTree := id.CloneOwner(ctx, userAddr)
	chg := TransferStageOnly(ctx, userAddr, govAddr, userTree, govRepo, fromUserOpt, fromBalance, toUser, toBalance, amount)
	proto.Commit(ctx, userTree.Public, chg.Msg)
	git.Push(ctx, userRepo.Public)
	return chg
}

func TransferStageOnly(
	ctx context.Context,
	userAddr id.OwnerAddress,
	govAddr gov.GovAddress,
	userTree id.OwnerTree,
	govRepo *git.Repository,
	fromUserOpt member.User,
	fromBalance balance.Balance,
	toUser member.User,
	toBalance balance.Balance,
	amount float64,
) git.Change[mail.SeqNo] {

	userCred := id.GetPublicCredentials(ctx, userTree.Public)

	// find the user name of userAddr in the community repo
	if fromUserOpt == "" {
		us := member.LookupUserByIDLocal(ctx, git.Worktree(ctx, govRepo), userCred.ID)
		switch len(us) {
		case 0:
			must.Errorf(ctx, "%s not found in community %v", userAddr.Public, govAddr)
		case 1:
			fromUserOpt = us[0]
		default:
			must.Errorf(ctx, "community %v has more than one user at address %v", govAddr, userAddr.Public)
		}
	}

	govTree := git.Worktree(ctx, govRepo)
	request := Request{
		Transfer: &TransferRequest{
			FromUser:    fromUserOpt,
			FromBalance: fromBalance,
			ToUser:      toUser,
			ToBalance:   toBalance,
			Amount:      amount,
		},
	}

	return mail.SendSignedStageOnly(ctx, userTree, govTree, BureauTopic, request)
}
