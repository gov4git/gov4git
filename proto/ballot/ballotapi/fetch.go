package ballotapi

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/gov4git/v2/proto/mail"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

type FetchedVote struct {
	Voter     member.User           `json:"voter_user"`
	Address   id.PublicAddress      `json:"voter_address"`
	Elections ballotproto.Elections `json:"voter_elections"`
}

type FetchedVotes []FetchedVote

func (x FetchedVotes) Len() int           { return len(x) }
func (x FetchedVotes) Less(i, j int) bool { return x[i].Voter < x[j].Voter }
func (x FetchedVotes) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func fetchedVotesToElections(fv FetchedVotes) map[member.User]ballotproto.Elections {
	ue := map[member.User]ballotproto.Elections{}
	for _, fv := range fv {
		ue[fv.Voter] = append(ue[fv.Voter], fv.Elections...)
	}
	return ue
}

func fetchVotes(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id ballotproto.BallotID,
	user member.User,
	account member.UserProfile,
) git.Change[form.Map, FetchedVotes] {
	userCloned := git.CloneOne(ctx, git.Address(account.PublicAddress))
	return fetchVotesCloned(ctx, cloned, id, user, account, userCloned)
}

func fetchVotesCloned(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id ballotproto.BallotID,
	user member.User,
	account member.UserProfile,
	userCloned git.Cloned,
) git.Change[form.Map, FetchedVotes] {

	fetched := FetchedVotes{}
	var respond mail.Responder[ballotproto.VoteEnvelope, ballotproto.VoteEnvelope] = func(
		ctx context.Context,
		_ mail.SeqNo,
		req ballotproto.VoteEnvelope,
	) (resp ballotproto.VoteEnvelope, err error) {

		if !req.VerifyConsistency() {
			return ballotproto.VoteEnvelope{}, fmt.Errorf("vote envelope is not valid")
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
	mail.Respond_StageOnly[ballotproto.VoteEnvelope, ballotproto.VoteEnvelope](
		ctx,
		cloned.IDOwnerCloned(),
		account.PublicAddress,
		voterPublicTree,
		ballotproto.BallotTopic(id),
		respond,
	)

	return git.NewChange(
		fmt.Sprintf("Fetched votes from user %v on ballot %v", user, id),
		"ballot_fetch_votes",
		form.Map{"id": id, "user": user, "account": account},
		fetched,
		nil,
	)
}
