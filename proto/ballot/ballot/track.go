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
	return TrackStageOnly(ctx, voterAddr, govAddr, voterOwner, govCloned, ballotName)
}

func TrackStageOnly(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	govAddr gov.GovAddress,
	voterOwner id.OwnerCloned,
	govCloned git.Cloned,
	ballotName ns.NS,
) common.VoterStatus {

	// determine the voter's username in the community
	voterCred := id.GetPublicCredentials(ctx, voterOwner.Public.Tree())
	user := member.LookupUserByIDLocal(ctx, govCloned.Tree(), voterCred.ID)
	must.Assertf(ctx, len(user) > 0, "user not found in community")

	// read the ballot tally
	tally := LoadTally(ctx, govCloned.Tree(), ballotName)

	// read the voter's log
	govCred := id.GetPublicCredentials(ctx, govCloned.Tree())
	voteLogNS := common.VoteLogPath(govCred.ID, ballotName)
	voteLog := git.FromFile[common.VoteLog](ctx, voterOwner.Public.Tree(), voteLogNS.Path())

	XXX

	return common.VoterStatus{XXX}
}
