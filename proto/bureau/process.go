package bureau

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/mail"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Process(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	group member.Group,
) git.ChangeNoResult {

	govRepo, govTree := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := ProcessStageOnly(ctx, govAddr, govRepo, govTree, group)
	proto.Commit(ctx, git.Worktree(ctx, govRepo.Home), chg.Msg)
	git.Push(ctx, govRepo.Home)
	return chg
}

func ProcessStageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	group member.Group,
) git.ChangeNoResult {

	communityTree := govTree.Home

	// list participating users
	users := member.ListGroupUsersLocal(ctx, communityTree, group)

	// get user accounts
	accounts := make([]member.Account, len(users))
	for i, user := range users {
		accounts[i] = member.GetUserLocal(ctx, communityTree, user)
	}

	// fetch user requests
	var fetchedReqs FetchedRequests
	for i, account := range accounts {
		fetchedReqs = append(fetchedReqs,
			fetchUserRequests(ctx, govAddr, govRepo, govTree, users[i], account).Result...)
	}

	// process requests
	for _, fetched := range fetchedReqs {
		processRequestStageOnly(ctx, govAddr, govRepo, govTree, fetched)
	}

	return git.ChangeNoResult{
		Msg: fmt.Sprintf("Process bureau requests of users in group %v", group),
	}
}

func processRequestStageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	fetched FetchedRequest,
) {
	for _, req := range fetched.Requests {
		if req.Transfer == nil {
			continue
		}
		if req.Transfer.FromUser != fetched.User {
			base.Infof("bureau: invalid transfer request from user %v; origin of transfer is not the requesting user", fetched.User)
			continue
		}
		err := must.Try(func() {
			balance.TransferStageOnly(
				ctx,
				govTree.Home,
				req.Transfer.FromUser, req.Transfer.FromBalance,
				req.Transfer.ToUser, req.Transfer.ToBalance,
				req.Transfer.Amount,
			)
		})
		if err != nil {
			base.Infof("bureau: transfer error (%v)", err)
			continue
		}
		base.Infof("bureau: transferred %v from %v:%v to %v:%v",
			req.Transfer.Amount,
			req.Transfer.FromUser, req.Transfer.FromBalance,
			req.Transfer.ToUser, req.Transfer.ToBalance,
		)
	}
}

func fetchUserRequests(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	user member.User,
	account member.Account,
) git.Change[FetchedRequests] {

	fetched := FetchedRequests{}
	respond := func(ctx context.Context, req Request, _ id.SignedPlaintext) (resp Request, err error) {
		fetched = append(fetched,
			FetchedRequest{
				User:     user,
				Address:  account.Home,
				Requests: Requests{req},
			})
		return req, nil
	}

	_, userPublicTree := git.Clone(ctx, git.Address(account.Home))
	mail.ReceiveSignedStageOnly(
		ctx,
		govTree,
		account.Home,
		userPublicTree,
		BureauTopic,
		respond,
	)

	return git.Change[FetchedRequests]{
		Result: fetched,
		Msg:    fmt.Sprintf("Fetched requests from user %v", user),
	}
}
