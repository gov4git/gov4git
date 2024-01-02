package ballotproto

import (
	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/kv"
)

var (
	BallotNS = proto.RootNS.Append("ballot")
	BallotKV = kv.KV[BallotID, struct{}]{}
)

var (
	AdFilebase      = "ballot_ad.json"
	TallyFilebase   = "ballot_tally.json"
	OutcomeFilebase = "ballot_outcome.json"
	PolicyFilebase  = "ballot_policy.json" // policy instance state
)

var (
	VoteLogNS = proto.RootNS.Append("votes") // namespace in voter's repo for recording votes
)
