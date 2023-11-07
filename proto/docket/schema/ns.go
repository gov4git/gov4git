package schema

import (
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/kv"
)

var (
	DocketNS = proto.RootNS.Append("docket")
	MotionNS = DocketNS.Append("motion")
	MotionKV = kv.KV[MotionID, Motion]{}
)

func MotionPollBallotName(id MotionID) common.BallotName {
	return common.BallotName{"docket", "motion", "poll", id.String()}
}
