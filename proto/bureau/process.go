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
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Process(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	group member.Group,
) git.ChangeNoResult {

	govOwner := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := ProcessStageOnly(ctx, govAddr, govOwner, group)
	proto.Commit(ctx, govOwner.Public.Tree(), chg)
	govOwner.Public.Push(ctx)
	return chg
}

func ProcessStageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govOwner id.OwnerCloned,
	group member.Group,
) git.ChangeNoResult {

	communityTree := govOwner.Public.Tree()

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
			fetchUserRequests(ctx, govAddr, govOwner, users[i], account).Result...)
	}

	// process requests
	for _, fetched := range fetchedReqs {
		processRequestStageOnly(ctx, govAddr, govOwner, fetched)
	}

	return git.NewChangeNoResult(
		fmt.Sprintf("Process bureau requests of users in group %v", group),
		"bureau_process",
	)
}

func processRequestStageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govOwner id.OwnerCloned,
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
				govOwner.Public.Tree(),
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
	govOwner id.OwnerCloned,
	user member.User,
	account member.Account,
) git.Change[form.Map, FetchedRequests] {

	fetched := FetchedRequests{}
	respond := func(ctx context.Context, _ mail.SeqNo, req Request, _ id.SignedPlaintext) (resp Request, err error) {
		fetched = append(fetched,
			FetchedRequest{
				User:     user,
				Address:  account.PublicAddress,
				Requests: Requests{req},
			})
		return req, nil
	}

	userPublic := git.CloneOne(ctx, git.Address(account.PublicAddress))
	recvOnly := mail.ReceiveSigned_StageOnly(
		ctx,
		govOwner,
		account.PublicAddress,
		userPublic.Tree(),
		BureauTopic,
		respond,
	)

	return git.NewChange(
		fmt.Sprintf("Fetched requests from user %v", user),
		"bureau_fetch_user_requests",
		form.Map{"user": user, "account": account},
		fetched,
		form.Forms{recvOnly},
	)
}
