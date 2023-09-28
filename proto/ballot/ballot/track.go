package ballot

import (
	"context"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func Track(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	govAddr gov.GovAddress,
	ballotName ns.NS,
) common.VoterStatus {

	govCloned := git.CloneOne(ctx, git.Address(govAddr))
	voterOwner := id.CloneOwner(ctx, voterAddr)
	return Track_StageOnly(ctx, voterAddr, govAddr, voterOwner, govCloned, ballotName)
}

func Track_StageOnly(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	govAddr gov.GovAddress,
	voterOwner id.OwnerCloned,
	govCloned git.Cloned,
	ballotName ns.NS,
) common.VoterStatus {

	// determine the voter's username in the community
	voterCred := id.GetPublicCredentials(ctx, voterOwner.Public.Tree())
	users := member.LookupUserByID_Local(ctx, govCloned.Tree(), voterCred.ID)
	must.Assertf(ctx, len(users) > 0, "user not found in community")
	user := users[0]

	// read the ballot tally
	tally := LoadTally(ctx, govCloned.Tree(), ballotName)

	// read the voter's log
	govCred := id.GetPublicCredentials(ctx, govCloned.Tree())
	voteLogNS := common.VoteLogPath(govCred.ID, ballotName)
	voteLog := git.FromFile[common.VoteLog](ctx, voterOwner.Public.Tree(), voteLogNS.Path())

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
	pending := common.Elections{}
	for _, env := range voteLog.VoteEnvelopes {
		for _, el := range env.Elections {
			if pendingVotes[el.VoteID] {
				pending = append(pending, el)
			}
		}
	}

	return common.VoterStatus{
		GovID:         voteLog.GovID,
		GovAddress:    voteLog.GovAddress,
		BallotName:    ballotName,
		AcceptedVotes: tally.AcceptedVotes[user],
		RejectedVotes: tally.RejectedVotes[user],
		PendingVotes:  pending,
	}
}
