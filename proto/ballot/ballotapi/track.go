package ballotapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/regime"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Track(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	addr gov.Address,
	ballotID ballotproto.BallotID,

) ballotproto.VoterStatus {

	cloned := gov.Clone(ctx, addr)
	voterOwner := id.CloneOwner(ctx, voterAddr)
	return Track_StageOnly(ctx, voterAddr, voterOwner, cloned, ballotID)
}

func FindVoterUser_Local(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	voterOwner id.OwnerCloned,
	cloned gov.Cloned,

) member.User {

	voterCred := id.GetPublicCredentials(ctx, voterOwner.Public.Tree())
	users := member.LookupUserByID_Local(ctx, cloned, voterCred.ID)
	must.Assertf(ctx, len(users) > 0, "user not found in community")
	return users[0]
}

func Track_StageOnly(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	voterOwner id.OwnerCloned,
	cloned gov.Cloned,
	ballotID ballotproto.BallotID,

) ballotproto.VoterStatus {

	ctx = regime.Dry(ctx)

	// determine the voter's username in the community
	user := member.FindClonedUser_Local(ctx, cloned, voterOwner)

	// read the ballot tally
	tally := loadTally_Local(ctx, cloned.Tree(), ballotID)

	// read the voter's log
	govCred := id.GetPublicCredentials(ctx, cloned.Tree())
	voteLogNS := ballotproto.VoteLogPath(govCred.ID, ballotID)
	voteLog, err := git.TryFromFile[ballotproto.VoteLog](ctx, voterOwner.Public.Tree(), voteLogNS)
	if git.IsNotExist(err) {
		return ballotproto.VoterStatus{
			GovID:         govCred.ID,
			GovAddress:    cloned.Address(),
			BallotID:      ballotID,
			AcceptedVotes: nil,
			RejectedVotes: nil,
			PendingVotes:  nil,
		}
	}
	must.NoError(ctx, err)

	// calculate pending votes
	pendingVotes := map[id.ID]bool{}
	for _, env := range voteLog.VoteEnvelopes {
		for _, el := range env.Elections {
			pendingVotes[el.VoteID] = true
		}
	}
	for _, acc := range tally.AcceptedVotes[user] {
		delete(pendingVotes, acc.Vote.VoteID)
	}
	for _, rej := range tally.RejectedVotes[user] {
		delete(pendingVotes, rej.Vote.VoteID)
	}

	// collect votes in order of execution
	pending := ballotproto.Elections{}
	for _, env := range voteLog.VoteEnvelopes {
		for _, el := range env.Elections {
			if pendingVotes[el.VoteID] {
				pending = append(pending, el)
			}
		}
	}

	return ballotproto.VoterStatus{
		GovID:         voteLog.GovID,
		GovAddress:    voteLog.GovAddress,
		BallotID:      ballotID,
		AcceptedVotes: tally.AcceptedVotes[user],
		RejectedVotes: tally.RejectedVotes[user],
		PendingVotes:  pending,
	}
}
