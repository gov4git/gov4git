package bureau

import (
	"context"

	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/gov4git/mod/balance"
	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/gov4git/mod/id"
	"github.com/gov4git/gov4git/mod/mail"
	"github.com/gov4git/gov4git/mod/member"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Transfer(
	ctx context.Context,
	userAddr id.OwnerAddress,
	govAddr gov.CommunityAddress,
	fromUserOpt member.User, // optional, if empty string, a lookup forthe user is performed
	fromBalance balance.Balance,
	toUser member.User,
	toBalance balance.Balance,
	amount float64,
) git.Change[mail.SeqNo] {

	govRepo := git.CloneRepo(ctx, git.Address(govAddr))
	userRepo, userTree := id.CloneOwner(ctx, userAddr)
	chg := TransferStageOnly(ctx, userAddr, govAddr, userTree, govRepo, fromUserOpt, fromBalance, toUser, toBalance, amount)
	mod.Commit(ctx, userTree.Public, chg.Msg)
	git.Push(ctx, userRepo.Public)
	return chg
}

func TransferStageOnly(
	ctx context.Context,
	userAddr id.OwnerAddress,
	govAddr gov.CommunityAddress,
	userTree id.OwnerTree,
	govRepo *git.Repository,
	fromUserOpt member.User,
	fromBalance balance.Balance,
	toUser member.User,
	toBalance balance.Balance,
	amount float64,
) git.Change[mail.SeqNo] {

	// find the user name of userAddr in the community repo
	if fromUserOpt == "" {
		us := member.LookupUserByAddressLocal(ctx, git.Worktree(ctx, govRepo), userAddr.Public)
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
