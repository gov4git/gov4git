package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/mail"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

type FetchedVote struct {
	Voter     member.User      `json:"voter_user"`
	Address   id.PublicAddress `json:"voter_address"`
	Elections common.Elections `json:"voter_elections"`
}

type FetchedVotes []FetchedVote

func (x FetchedVotes) Len() int           { return len(x) }
func (x FetchedVotes) Less(i, j int) bool { return x[i].Voter < x[j].Voter }
func (x FetchedVotes) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func fetchedVotesToElections(fv FetchedVotes) map[member.User]common.Elections {
	ue := map[member.User]common.Elections{}
	for _, fv := range fv {
		ue[fv.Voter] = append(ue[fv.Voter], fv.Elections...)
	}
	return ue
}

func fetchVotes(
	ctx context.Context,
	govAddr gov.GovOwnerAddress,
	govOwner gov.GovOwnerCloned,
	ballotName common.BallotName,
	user member.User,
	account member.Account,
) git.Change[form.Map, FetchedVotes] {
	userCloned := git.CloneOne(ctx, git.Address(account.PublicAddress))
	return fetchVotesCloned(ctx, govAddr, govOwner, ballotName, user, account, userCloned)
}

func fetchVotesCloned(
	ctx context.Context,
	govAddr gov.GovOwnerAddress,
	govOwner gov.GovOwnerCloned,
	ballotName common.BallotName,
	user member.User,
	account member.Account,
	userCloned git.Cloned,
) git.Change[form.Map, FetchedVotes] {

	fetched := FetchedVotes{}
	var respond mail.Responder[common.VoteEnvelope, common.VoteEnvelope] = func(
		ctx context.Context,
		_ mail.SeqNo,
		req common.VoteEnvelope,
	) (resp common.VoteEnvelope, err error) {

		if !req.VerifyConsistency() {
			return common.VoteEnvelope{}, fmt.Errorf("vote envelope is not valid")
		}
		fetched = append(fetched,
			FetchedVote{
				Voter:     user,
				Address:   account.PublicAddress,
				Elections: req.Elections,
			})
		return req, nil
	}

	voterPublicTree := userCloned.Tree()
	mail.Respond_StageOnly[common.VoteEnvelope, common.VoteEnvelope](
		ctx,
		govOwner.IDOwnerCloned(),
		account.PublicAddress,
		voterPublicTree,
		common.BallotTopic(ballotName),
		respond,
	)

	return git.NewChange(
		fmt.Sprintf("Fetched votes from user %v on ballot %v", user, ballotName),
		"ballot_fetch_votes",
		form.Map{"ballot_name": ballotName, "user": user, "account": account},
		fetched,
		nil,
	)
}
