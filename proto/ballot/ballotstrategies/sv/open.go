package sv

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/member"
)

func (qv SV) Open(
	ctx context.Context,
	owner gov.OwnerCloned,
	ad *ballotproto.Advertisement,

) *ballotproto.Tally {

	return &ballotproto.Tally{
		Ad:            *ad,
		Scores:        map[string]float64{},
		ScoresByUser:  map[member.User]map[string]ballotproto.StrengthAndScore{},
		AcceptedVotes: map[member.User]ballotproto.AcceptedElections{},
		RejectedVotes: map[member.User]ballotproto.RejectedElections{},
		Charges:       map[member.User]float64{},
	}
}
